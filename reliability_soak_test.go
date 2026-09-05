package n2k

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
)

type soakSample struct {
	ElapsedSeconds float64 `json:"elapsedSeconds"`
	Cycles         uint64  `json:"cycles"`
	HeapBytes      uint64  `json:"heapBytes"`
	Goroutines     int     `json:"goroutines"`
}

type soakStatistics struct {
	Kind                string       `json:"kind"`
	Status              string       `json:"status"`
	GoVersion           string       `json:"goVersion"`
	GOOS                string       `json:"goos"`
	GOARCH              string       `json:"goarch"`
	TestBinarySHA256    string       `json:"testBinarySHA256"`
	RequestedDuration   string       `json:"requestedDuration"`
	StartedAt           time.Time    `json:"startedAt"`
	ElapsedSeconds      float64      `json:"elapsedSeconds"`
	Cycles              uint64       `json:"cycles"`
	FramesInjected      uint64       `json:"framesInjected"`
	CorruptFrames       uint64       `json:"corruptFrames"`
	ReplayMessages      uint64       `json:"replayMessages"`
	ApplicationWrites   uint64       `json:"applicationWrites"`
	QueueOverflows      uint64       `json:"queueOverflows"`
	SubscriberOverflows uint64       `json:"subscriberOverflows"`
	Reconnects          uint64       `json:"reconnects"`
	ClientStarts        uint64       `json:"clientStarts"`
	InjectedFailures    uint64       `json:"injectedFailures"`
	Baseline            soakSample   `json:"baseline"`
	PeakHeapBytes       uint64       `json:"peakHeapBytes"`
	PeakGoroutines      int          `json:"peakGoroutines"`
	HeapAllowanceBytes  uint64       `json:"heapAllowanceBytes"`
	GoroutineAllowance  int          `json:"goroutineAllowance"`
	Samples             []soakSample `json:"recentSamples"`
}

// TestReliabilitySoak is software evidence only. It deliberately skips in the
// ordinary suite; just reliability-soak supplies an explicit duration.
func TestReliabilitySoak(t *testing.T) {
	requested := os.Getenv("N2K_SOAK_DURATION")
	if requested == "" {
		t.Skip("set N2K_SOAK_DURATION or run just reliability-soak; software soak not run")
	}
	duration, err := time.ParseDuration(requested)
	if err != nil || duration <= 0 || duration > 24*time.Hour {
		t.Fatalf("N2K_SOAK_DURATION must be positive and at most 24h: %q", requested)
	}
	stats := soakStatistics{
		Kind: "software-simulation", Status: "running", GoVersion: runtime.Version(),
		GOOS: runtime.GOOS, GOARCH: runtime.GOARCH,
		RequestedDuration: requested, StartedAt: time.Now().UTC(),
		HeapAllowanceBytes: 64 << 20, GoroutineAllowance: 64,
	}
	artifactDirectory := os.Getenv("N2K_SOAK_ARTIFACT_DIR")
	if artifactDirectory == "" {
		artifactDirectory = "conformance-artifacts"
	}
	if err := os.MkdirAll(artifactDirectory, 0o750); err != nil {
		t.Fatal(err)
	}
	defer func() {
		stats.ElapsedSeconds = time.Since(stats.StartedAt).Seconds()
		stats.Status = "pass"
		if t.Failed() {
			stats.Status = "fail"
		}
		encoded, err := json.MarshalIndent(stats, "", "  ")
		if err == nil {
			err = os.WriteFile(filepath.Join(artifactDirectory, "soak-stats.json"), append(encoded, '\n'), 0o600)
		}
		if err != nil {
			t.Errorf("write software soak evidence: %v", err)
		}
	}()
	stats.TestBinarySHA256, err = soakExecutableSHA256()
	if err != nil {
		t.Fatalf("identify executed software soak test binary: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration+30*time.Second)
	defer cancel()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	var client *Client
	var bus *soakBus
	defer func() {
		if client != nil {
			if err := client.Close(); err != nil {
				t.Errorf("close soak client: %v", err)
			}
		}
	}()
	heading := &pgn.VesselHeading{}
	heading.SetHeadingValue(1.25)
	payload, err := heading.EncodePayload()
	if err != nil {
		t.Fatal(err)
	}
	valid := framer.FrameSingle(framer.BuildCANID(127250, 3, 42, 255), payload)
	frames := make([]can.Frame, 0, 68)
	for i := 0; i < 64; i++ {
		frames = append(frames, valid)
	}
	frames = append(frames,
		can.Frame{ID: valid.ID, Length: 1, Data: [8]byte{0xff}},
		can.Frame{ID: framer.BuildCANID(60416, 6, 43, 255), Length: 8, Data: [8]byte{32, 16, 0, 2, 255, 0x16, 0xf0, 1}},
		can.Frame{ID: framer.BuildCANID(126996, 6, 44, 255), Length: 8, Data: [8]byte{0, 223, 1, 2, 3, 4, 5, 6}},
		can.Frame{ID: framer.BuildCANID(126996, 6, 44, 255), Length: 8, Data: [8]byte{31, 7, 8, 9, 10, 11, 12, 13}},
	)
	sampleInterval := max(time.Second, min(duration/60, 30*time.Second))
	nextSample := time.Now()
	violations := 0
	for time.Since(stats.StartedAt) < duration {
		if client == nil {
			bus = newSoakBus()
			client, err = NewClient(ctx, WithBus(bus), WithClaimTimeout(50*time.Millisecond), WithHeartbeatInterval(0), WithReceiveBuffer(8), WithWriteQueue(4), WithLogger(logger))
			if err != nil {
				t.Fatal(err)
			}
			stats.ClientStarts++
		}
		if err := soakCycle(ctx, client, bus, heading, frames, &stats); err != nil {
			t.Fatalf("software soak cycle %d: %v", stats.Cycles, err)
		}
		stats.Cycles++
		if stats.Cycles%4 == 0 {
			if err := bus.send(ctx, soakEvent{operation: "disconnect"}); err != nil {
				t.Fatal(err)
			}
			if err := bus.send(ctx, soakEvent{operation: "connect"}); err != nil {
				t.Fatal(err)
			}
			if err := soakWait(ctx, func() bool { status := client.Status(); return status.Connected && !status.Rejoining }); err != nil {
				t.Fatalf("reconnect did not restore readiness: %v, status=%+v", err, client.Status())
			}
			stats.Reconnects++
		}
		if stats.Cycles%16 == 0 {
			if err := bus.send(ctx, soakEvent{operation: "fail"}); err != nil {
				t.Fatal(err)
			}
			if err := soakWait(ctx, func() bool { return client.Err() != nil }); err != nil {
				t.Fatalf("injected bus failure was not observable: %v", err)
			}
			if err := client.Close(); err != nil {
				t.Fatal(err)
			}
			client = nil
			stats.InjectedFailures++
		}
		if !time.Now().Before(nextSample) {
			runtime.GC()
			var memory runtime.MemStats
			runtime.ReadMemStats(&memory)
			sample := soakSample{ElapsedSeconds: time.Since(stats.StartedAt).Seconds(), Cycles: stats.Cycles, HeapBytes: memory.HeapAlloc, Goroutines: runtime.NumGoroutine()}
			if stats.Baseline.Cycles == 0 || (stats.Baseline.Cycles < 8 && stats.Cycles >= 8) {
				// Take a fresh baseline once repeated initialization, disconnect,
				// and replay paths have warmed their bounded caches.
				stats.Baseline = sample
			}
			stats.PeakHeapBytes = max(stats.PeakHeapBytes, sample.HeapBytes)
			stats.PeakGoroutines = max(stats.PeakGoroutines, sample.Goroutines)
			if len(stats.Samples) == 120 {
				copy(stats.Samples, stats.Samples[1:])
				stats.Samples = stats.Samples[:119]
			}
			stats.Samples = append(stats.Samples, sample)
			if sample.HeapBytes > stats.Baseline.HeapBytes+stats.HeapAllowanceBytes || sample.Goroutines > stats.Baseline.Goroutines+stats.GoroutineAllowance {
				violations++
			} else {
				violations = 0
			}
			t.Logf("software soak elapsed=%.1fs cycles=%d heap=%d goroutines=%d reconnects=%d failures=%d", sample.ElapsedSeconds, stats.Cycles, sample.HeapBytes, sample.Goroutines, stats.Reconnects, stats.InjectedFailures)
			if violations >= 3 {
				t.Fatalf("resource growth exceeded the warm baseline allowances for three consecutive post-GC samples: baseline=%+v latest=%+v", stats.Baseline, sample)
			}
			nextSample = time.Now().Add(sampleInterval)
		}
	}
	if stats.Cycles == 0 {
		t.Fatal("software soak executed no cycles")
	}
}

func soakExecutableSHA256() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	executable, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = executable.Close() }()
	digest := sha256.New()
	if _, err := io.Copy(digest, executable); err != nil {
		return "", err
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}

func soakCycle(ctx context.Context, client *Client, bus *soakBus, heading *pgn.VesselHeading, frames []can.Frame, stats *soakStatistics) error {
	cycleCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	slow := client.Scanner()
	defer func() { _ = slow.Close() }()
	if err := bus.send(cycleCtx, soakEvent{frames: frames}); err != nil {
		return err
	}
	for slow.Next() {
	}
	if !errors.Is(slow.Err(), ErrReceiveOverflow) {
		return fmt.Errorf("slow subscriber must fail explicitly: %v", slow.Err())
	}
	stats.SubscriberOverflows++
	stats.FramesInjected += uint64(len(frames))
	stats.CorruptFrames += 4
	// Replay uses the public source/scanner API and allocates only one finite
	// capture per cycle. The mock never retains transmitted frame histories.
	replay := NewScanner(cycleCtx, Replay(frames[:2]), WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	for replay.Next() {
		stats.ReplayMessages++
	}
	if err := replay.Err(); err != nil {
		_ = replay.Close()
		return err
	}
	_ = replay.Close()
	bus.blockApplications()
	defer bus.releaseApplications()
	results := make([]*WriteResult, 0, 33)
	results = append(results, client.Write(heading))
	select {
	case <-bus.writeEntered:
	case <-cycleCtx.Done():
		return fmt.Errorf("application write did not reach the bus: %w", cycleCtx.Err())
	}
	for i := 0; i < 32; i++ {
		results = append(results, client.Write(heading))
	}
	// Release the single physical write before asking for a claim response;
	// a blocked physical frame cannot be preempted by any protocol scheduler.
	// This verifies recovery after saturation without an impossible bus model.
	bus.releaseApplications()
	request := framer.FrameSingle(framer.BuildCANID(59904, 6, 42, 255), []byte{0, 0xee, 0})
	if err := bus.send(cycleCtx, soakEvent{frames: []can.Frame{request}}); err != nil {
		return err
	}
	overflows := uint64(0)
	for _, result := range results {
		err := result.WaitContext(cycleCtx)
		if errors.Is(err, ErrWriteQueueFull) {
			overflows++
		} else if err != nil {
			return err
		} else {
			stats.ApplicationWrites++
		}
	}
	if overflows == 0 {
		return errors.New("saturated application queue reported no admission failures")
	}
	stats.QueueOverflows += overflows
	return client.Err()
}

func soakWait(ctx context.Context, ready func() bool) error {
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()
	for !ready() {
		select {
		case <-ticker.C:
		case <-timer.C:
			return errors.New("software soak lifecycle transition timed out")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

type soakEvent struct {
	operation string
	frames    []can.Frame
	done      chan struct{}
}

type soakBus struct {
	events       chan soakEvent
	done         chan struct{}
	ready        chan struct{}
	writeEntered chan struct{}
	once         sync.Once
	writes       atomic.Uint64
	mu           sync.Mutex
	observer     func(bool, uint64)
	connected    bool
	epoch        uint64
	changed      chan struct{}
	writeBlock   chan struct{}
}

func newSoakBus() *soakBus {
	return &soakBus{events: make(chan soakEvent, 8), done: make(chan struct{}), ready: make(chan struct{}), writeEntered: make(chan struct{}, 1), changed: make(chan struct{})}
}

func (b *soakBus) Ready() <-chan struct{} { return b.ready }

func (b *soakBus) SetConnectionObserver(observer func(bool, uint64)) {
	b.mu.Lock()
	b.observer = observer
	b.mu.Unlock()
}

func (b *soakBus) transition(connected bool) {
	b.mu.Lock()
	if connected {
		b.epoch++
	} else {
		b.connected = false
	}
	observer, epoch := b.observer, b.epoch
	b.mu.Unlock()
	if observer != nil {
		observer(connected, epoch)
	}
	b.mu.Lock()
	b.connected = connected
	close(b.changed)
	b.changed = make(chan struct{})
	b.mu.Unlock()
}

func (b *soakBus) Run(ctx context.Context, handler func(can.Frame)) error {
	b.transition(true)
	close(b.ready)
	for {
		select {
		case event := <-b.events:
			switch event.operation {
			case "connect":
				b.transition(true)
			case "disconnect":
				b.transition(false)
			case "fail":
				close(event.done)
				return errors.New("software soak injected bus disconnect")
			default:
				for _, frame := range event.frames {
					handler(frame)
				}
			}
			close(event.done)
		case <-ctx.Done():
			return ctx.Err()
		case <-b.done:
			return nil
		}
	}
}

func (b *soakBus) send(ctx context.Context, event soakEvent) error {
	event.done = make(chan struct{})
	select {
	case b.events <- event:
	case <-ctx.Done():
		return ctx.Err()
	case <-b.done:
		return errors.New("software soak bus closed")
	}
	select {
	case <-event.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-b.done:
		return errors.New("software soak bus closed")
	}
}

func (b *soakBus) WriteFrame(frame can.Frame) error {
	return b.WriteFrameContext(context.Background(), frame)
}

func (b *soakBus) WriteFrameContext(ctx context.Context, frame can.Frame) error {
	for {
		b.mu.Lock()
		connected, changed, block := b.connected, b.changed, b.writeBlock
		b.mu.Unlock()
		if connected {
			if block != nil && framer.ParseCANID(frame.ID).PGN == 127250 {
				select {
				case b.writeEntered <- struct{}{}:
				default:
				}
				select {
				case <-block:
				case <-ctx.Done():
					return context.Cause(ctx)
				case <-b.done:
					return errors.New("software soak bus closed")
				}
			}
			b.writes.Add(1)
			return nil
		}
		select {
		case <-changed:
		case <-ctx.Done():
			return context.Cause(ctx)
		case <-b.done:
			return errors.New("software soak bus closed")
		}
	}
}

func (b *soakBus) blockApplications() {
	b.mu.Lock()
	b.writeBlock = make(chan struct{})
	b.mu.Unlock()
}

func (b *soakBus) releaseApplications() {
	b.mu.Lock()
	if b.writeBlock != nil {
		close(b.writeBlock)
		b.writeBlock = nil
	}
	b.mu.Unlock()
}

func (b *soakBus) Close() error {
	b.once.Do(func() { close(b.done) })
	return nil
}
