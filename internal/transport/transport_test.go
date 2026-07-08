package transport

import (
	"sync"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/stretchr/testify/assert"
)

// fakeTimer is a no-op Timer used by fakeClock.
type fakeTimer struct{ stopped bool }

func (f *fakeTimer) Stop() bool { f.stopped = true; return true }

// fakeClock captures scheduled callbacks so tests fire them synchronously
// instead of sleeping real wall-clock time.
type fakeClock struct {
	mu        sync.Mutex
	callbacks []func()
}

func (c *fakeClock) AfterFunc(d time.Duration, f func()) Timer {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.callbacks = append(c.callbacks, f)
	return &fakeTimer{}
}

func (c *fakeClock) fireAll() {
	c.mu.Lock()
	cbs := c.callbacks
	c.callbacks = nil
	c.mu.Unlock()
	for _, f := range cbs {
		f()
	}
}

func TestBAMReceiveTimeoutDeterministic(t *testing.T) {
	clock := &fakeClock{}
	var written []can.Frame
	mgr := NewManager(ManagerConfig{
		WriteFrame: func(f can.Frame) error { written = append(written, f); return nil },
		AfterFunc:  clock.AfterFunc,
		Sleep:      func(time.Duration) {},
	})
	defer mgr.Close()

	// Start a BAM receive session (CM frame), then fire the DT timeout without sleeping.
	cm := buildCMBAMFrame(14, 2, 130816, 10)
	mgr.HandleFrame(cm)
	if got := activeSessionCount(mgr); got != 1 {
		t.Fatalf("sessions = %d, want 1", got)
	}
	clock.fireAll()
	if got := activeSessionCount(mgr); got != 0 {
		t.Fatalf("sessions after timeout = %d, want 0", got)
	}
}

func activeSessionCount(m *Manager) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.sessions)
}

func TestOnCompleteRunsOutsideLock(t *testing.T) {
	release := make(chan struct{})
	entered := make(chan struct{})
	mgr := NewManager(ManagerConfig{
		WriteFrame: func(can.Frame) error { return nil },
		OnComplete: func(uint32, uint8, uint8, []byte) {
			close(entered)
			<-release // simulate a slow consumer
		},
	})
	defer mgr.Close()

	// Complete a 1-frame BAM: CM announcing 7 bytes / 1 frame, then DT 1.
	mgr.HandleFrame(buildCMBAMFrame(7, 1, 130816, 10))
	done := make(chan struct{})
	go func() {
		var dt [7]byte
		mgr.HandleFrame(makeDTFrame(1, dt, 10, BroadcastAddr))
		close(done)
	}()
	<-entered
	// While OnComplete is blocked, another goroutine must still be able to
	// take the manager lock (e.g. a new CM frame).
	locked := make(chan struct{})
	go func() {
		mgr.HandleFrame(buildCMBAMFrame(7, 1, 130817, 11))
		close(locked)
	}()
	select {
	case <-locked:
	case <-time.After(time.Second):
		t.Fatal("Manager lock held across OnComplete")
	}
	close(release)
	<-done
}

func TestManager_RoutesCMToCorrectHandler(t *testing.T) {
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{
		WriteFrame: h.writeFrame,
		OnComplete: h.onComplete,
	})
	defer mgr.Close()

	source := uint8(42)
	transportedPGN := uint32(126996)

	// Send a CM_BAM frame. This should create a BAM receive session.
	bamFrame := buildTestCMFrame(ControlBAM, 14, 2, transportedPGN, source, BroadcastAddr)
	mgr.HandleFrame(bamFrame)

	mgr.mu.Lock()
	bamKey := sessionKey{source: source, destination: BroadcastAddr, pgn: transportedPGN}
	_, bamExists := mgr.sessions[bamKey]
	mgr.mu.Unlock()
	assert.True(t, bamExists, "BAM session should have been created")

	// Send an RTS frame. This should create an RTS/CTS receive session.
	destination := uint8(10)
	rtsPGN := uint32(126998)
	rtsFrame := buildTestCMFrame(ControlRTS, 14, 2, rtsPGN, source, destination)
	rtsFrame.Data[4] = 0xFF
	mgr.HandleFrame(rtsFrame)

	mgr.mu.Lock()
	rtsKey := sessionKey{source: source, destination: destination, pgn: rtsPGN}
	_, rtsExists := mgr.sessions[rtsKey]
	mgr.mu.Unlock()
	assert.True(t, rtsExists, "RTS/CTS session should have been created")

	// Verify a CTS was sent in response to the RTS (not for the BAM).
	sentFrames := h.getSentFrames()
	assert.Len(t, sentFrames, 1, "should have sent exactly one CTS for the RTS")
	assert.Equal(t, ControlCTS, sentFrames[0].Data[0], "response should be a CTS")
}

func TestManager_IgnoresNonTPFrames(t *testing.T) {
	h := newTestHelper()
	mgr := NewManager(ManagerConfig{
		WriteFrame: h.writeFrame,
		OnComplete: h.onComplete,
	})
	defer mgr.Close()

	// Send a frame with a non-TP PGN (e.g., VesselHeading 127250).
	canID := framer.BuildCANID(127250, 2, 42, 255)
	nonTPFrame := can.Frame{
		ID:     canID,
		Length: 8,
		Data:   [8]uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}
	mgr.HandleFrame(nonTPFrame)

	// No sessions should have been created.
	mgr.mu.Lock()
	assert.Empty(t, mgr.sessions, "non-TP frame should not create any sessions")
	mgr.mu.Unlock()

	// No frames should have been sent.
	assert.Empty(t, h.getSentFrames(), "non-TP frame should not trigger any response")
}

func TestParseCANID_RoundTrip(t *testing.T) {
	tests := []struct {
		name        string
		pgn         uint32
		source      uint8
		destination uint8
	}{
		{"TP.CM broadcast", PGNCM, 42, BroadcastAddr},
		{"TP.CM addressed", PGNCM, 42, 10},
		{"TP.DT broadcast", PGNDT, 42, BroadcastAddr},
		{"TP.DT addressed", PGNDT, 42, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			canID := framer.BuildCANID(tt.pgn, TPPriority, tt.source, tt.destination)
			c := framer.ParseCANID(canID)

			assert.Equal(t, tt.pgn, c.PGN, "PGN should round-trip")
			assert.Equal(t, tt.source, c.Source, "source should round-trip")
			assert.Equal(t, tt.destination, c.Destination, "destination should round-trip")
		})
	}
}

func TestEncodePGN_RoundTrip(t *testing.T) {
	tests := []uint32{0, 60160, 60416, 126996, 130306, 0x1FFFF}
	for _, pgn := range tests {
		buf := make([]byte, 3)
		encodePGN(buf, pgn)
		var data [8]uint8
		copy(data[5:8], buf)
		got := extractPGN(data)
		assert.Equal(t, pgn, got, "PGN 0x%X should round-trip", pgn)
	}
}
