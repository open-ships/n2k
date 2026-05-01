package n2k

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/internal/transport"
	"github.com/open-ships/n2k/pgn"
)

// Client is the central integration point for NMEA 2000 communication. It
// composes address claiming, transport protocol, encoding, and framing into a
// single type that can both read and write PGN messages.
type Client struct {
	cfg    config
	ctx    context.Context
	cancel context.CancelFunc
	log    *slog.Logger

	// sourceAddr is the address used as the CAN source for outgoing frames.
	sourceAddr uint8

	// deviceName is the packed 64-bit ISO 11783 NAME.
	deviceName uint64

	// writeFrame sends a single CAN frame. For replay sources this records the
	// frame; for real buses it writes to the hardware.
	writeFrame func(can.Frame) error

	// tp is the transport protocol manager used for multi-frame messages that
	// exceed fast-packet capacity (> 223 bytes).
	tp *transport.Manager

	// fastSeq is an atomic counter for fast-packet sequence IDs (mod 8).
	fastSeq atomic.Uint32

	// writeCh serializes all write operations through a single goroutine to
	// guarantee FIFO ordering per bus for multi-frame messages.
	writeCh chan writeJob
	writeWg sync.WaitGroup

	// opts preserves the original Option slice so Receive and Scanner can
	// reconstruct the same configuration.
	opts []Option

	// mu guards writtenFrames and closed state.
	mu            sync.Mutex
	writtenFrames []can.Frame // captured frames (replay/testing)
	closed        bool
}

type writeJob struct {
	msg    any
	result *WriteResult
}

// NewClient creates a Client that can read and write NMEA 2000 messages.
// At least one source option (CAN, USB, or Replay) must be provided.
func NewClient(ctx context.Context, opts ...Option) (*Client, error) {
	cfg := config{}
	for _, o := range opts {
		o.apply(&cfg)
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	if cfg.logger == nil {
		cfg.logger = slog.Default()
	}

	ctx, cancel := context.WithCancel(ctx)

	// Determine device NAME.
	var deviceName uint64
	if cfg.deviceName != nil {
		deviceName = cfg.deviceName.Pack(true)
	} else {
		deviceName = DefaultDeviceName().Pack(true)
	}

	// Determine source address.
	var sourceAddr uint8
	if cfg.sourceAddress != nil {
		sourceAddr = *cfg.sourceAddress
	} else {
		// Default to 0 for replay, 253 for real buses.
		sourceAddr = 0
		for _, src := range cfg.sources {
			if _, ok := src.(*replaySource); !ok {
				sourceAddr = 253
				break
			}
		}
	}

	c := &Client{
		cfg:        cfg,
		ctx:        ctx,
		cancel:     cancel,
		log:        cfg.logger,
		sourceAddr: sourceAddr,
		deviceName: deviceName,
		opts:       opts,
	}

	// Determine if we have a real bus or replay-only.
	hasRealBus := false
	for _, src := range cfg.sources {
		if _, ok := src.(*replaySource); !ok {
			hasRealBus = true
			break
		}
	}

	if hasRealBus {
		// TODO: Real bus integration — create canbus.Interface, run it, set
		// up address claiming. For now, return an error indicating this path
		// is not yet implemented.
		cancel()
		return nil, errors.New("n2k: real bus (CAN/USB) client not yet implemented; use Replay for testing")
	}

	// Replay path: no real bus, capture frames in memory.
	c.writeFrame = func(f can.Frame) error {
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.closed {
			return errors.New("n2k: client closed")
		}
		c.writtenFrames = append(c.writtenFrames, f)
		return nil
	}

	// Create transport manager for BAM/RTS-CTS sends.
	c.tp = transport.NewManager(transport.ManagerConfig{
		WriteFrame: c.writeFrame,
		Logger:     c.log,
	})

	// Start the single writer goroutine for FIFO ordering.
	c.writeCh = make(chan writeJob, 64)
	c.writeWg.Add(1)
	go c.writeLoop()

	return c, nil
}

// Write asynchronously encodes and transmits a PGN message. The message must
// be a pointer to a PGN struct (e.g. *pgn.VesselHeading). The returned
// WriteResult can be used to wait for completion and check for errors.
// Writes are serialized through a single goroutine to guarantee FIFO ordering.
func (c *Client) Write(msg any) *WriteResult {
	wr := newWriteResult()

	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		wr.complete(errors.New("n2k: client closed"))
		return wr
	}
	c.mu.Unlock()

	c.writeCh <- writeJob{msg: msg, result: wr}
	return wr
}

func (c *Client) writeLoop() {
	defer c.writeWg.Done()
	for job := range c.writeCh {
		job.result.complete(c.doWrite(job.msg))
	}
}

// doWrite performs the synchronous work of encoding and framing a PGN message.
func (c *Client) doWrite(msg any) error {
	// Extract MessageInfo from the struct's Info field using reflect.
	rv := reflect.ValueOf(msg)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("n2k: expected pointer to PGN struct, got %T", msg)
	}

	infoField := rv.FieldByName("Info")
	if !infoField.IsValid() {
		return fmt.Errorf("n2k: type %T has no Info field; not a PGN struct", msg)
	}
	info, ok := infoField.Interface().(pgn.MessageInfo)
	if !ok {
		return fmt.Errorf("n2k: Info field of %T is not pgn.MessageInfo", msg)
	}

	pgnNum := info.PGN
	if pgnNum == 0 {
		return fmt.Errorf("n2k: PGN number is 0 in %T; set Info.PGN before writing", msg)
	}

	// Look up the encoder.
	encoder, ok := pgn.EncoderLookup[pgnNum]
	if !ok {
		return fmt.Errorf("n2k: no encoder registered for PGN %d", pgnNum)
	}

	// Encode the struct to bytes.
	payload, err := encoder(msg)
	if err != nil {
		return fmt.Errorf("n2k: encode PGN %d: %w", pgnNum, err)
	}

	// Determine priority: use Info.Priority if non-zero, else default 6.
	priority := info.Priority
	if priority == 0 {
		priority = 6
	}

	// Determine destination: use Info.TargetId if non-zero, else 255 (broadcast).
	destination := info.TargetId
	if destination == 0 {
		destination = 255
	}

	// Build the CAN ID.
	canID := framer.BuildCANID(pgnNum, priority, c.sourceAddr, destination)

	// Determine framing strategy.
	if len(payload) <= 8 {
		// Single frame.
		frame := framer.FrameSingle(canID, payload)
		return c.writeFrame(frame)
	}

	// Check if this is a fast-packet PGN.
	isFast := false
	if infos, ok := pgn.PgnInfoLookup[pgnNum]; ok && len(infos) > 0 {
		isFast = infos[0].Fast
	}

	if isFast && len(payload) <= 223 {
		// Fast-packet.
		seqID := uint8(c.fastSeq.Add(1) % 8)
		frames := framer.FrameFastPacket(canID, payload, seqID)
		for _, f := range frames {
			if err := c.writeFrame(f); err != nil {
				return err
			}
		}
		return nil
	}

	// Large payload: use transport protocol BAM (broadcast).
	return c.tp.SendBAM(pgnNum, c.sourceAddr, payload)
}

// Receive returns an iterator of decoded NMEA 2000 messages from the
// configured sources. It delegates to the top-level n2k.Receive function
// with the same options used to create this Client.
func (c *Client) Receive() iter.Seq2[any, error] {
	return Receive(c.ctx, c.opts...)
}

// Scanner creates a new Scanner that reads from the configured sources. It
// delegates to n2k.NewScanner with the same options used to create this Client.
func (c *Client) Scanner() *Scanner {
	return NewScanner(c.ctx, c.opts...)
}

// WrittenFrames returns a copy of all CAN frames written through this client.
// This is primarily useful in replay/testing mode to inspect what was sent.
func (c *Client) WrittenFrames() []can.Frame {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]can.Frame, len(c.writtenFrames))
	copy(out, c.writtenFrames)
	return out
}

// Close shuts down the client, cancels the context, and releases resources.
// It is safe to call Close multiple times.
func (c *Client) Close() error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil
	}
	c.closed = true
	c.mu.Unlock()

	close(c.writeCh)
	c.writeWg.Wait()
	c.cancel()
	if c.tp != nil {
		c.tp.Close()
	}
	return nil
}
