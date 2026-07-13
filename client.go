package n2k

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/adapter"
	"github.com/open-ships/n2k/internal/claiming"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/internal/transport"
	"github.com/open-ships/n2k/pgn"
)

// defaultClaimTimeout is the maximum time NewClient blocks waiting for address
// claiming to complete.
const defaultClaimTimeout = 1500 * time.Millisecond

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
	// frame; for bus clients it writes to the hardware.
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

	// mu guards writtenFrames, closed state, and sourceAddr.
	mu            sync.Mutex
	writtenFrames []can.Frame // captured frames (replay/testing)
	closed        bool

	bus     Bus
	claimer *claiming.Claimer
	addrErr error // set by OnFatalError during address claiming

	// addrReady is closed once address claiming completes (or immediately for replay).
	addrReady chan struct{}

	// msgCh delivers decoded messages from the internal read loop to the read API.
	// nil for replay clients.
	msgCh chan pgn.Message

	// pipeline is the persistent read pipeline (pre-filter -> assembly ->
	// decode -> unknown-PGN policy -> post-filter -> msgCh) for the internal
	// read loop. Only used for bus clients.
	pipeline *readPipeline
}

type writeJob struct {
	msg    pgn.Message
	result *WriteResult
}

// NewClient creates a Client that can read and write NMEA 2000 messages.
// At least one source option (CAN, USB, Replay, or Bus) must be provided.
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

	// Validate the device NAME eagerly so a malformed NAME fails fast.
	if cfg.deviceName != nil {
		if err := cfg.deviceName.Validate(); err != nil {
			return nil, fmt.Errorf("n2k: invalid device name: %w", err)
		}
	}

	// Compile the CEL filter eagerly (and discard the result) so a bad
	// expression fails fast regardless of bus/replay mode. The read
	// pipeline (built in initBus for bus clients, or per-Scanner call for
	// replay clients) compiles its own copy; this preserves NewClient's
	// eager-error contract without threading the compiled filter through.
	if cfg.filterExpr != "" {
		if _, err := compileFilter(cfg.filterExpr); err != nil {
			return nil, fmt.Errorf("n2k: compiling filter: %w", err)
		}
	}

	ctx, cancel := context.WithCancel(ctx)

	// Determine device NAME.
	var deviceName uint64
	if cfg.deviceName != nil {
		deviceName = cfg.deviceName.Pack(true)
	} else {
		deviceName = DefaultDeviceName().Pack(true)
	}

	// Determine if we have a hardware bus or replay-only.
	hasBus := cfg.bus != nil
	if !hasBus {
		for _, src := range cfg.sources {
			if _, ok := src.(*replaySource); !ok {
				hasBus = true
				break
			}
		}
	}

	// Determine source address.
	var sourceAddr uint8
	if cfg.sourceAddress != nil {
		sourceAddr = *cfg.sourceAddress
	} else if hasBus {
		sourceAddr = 253
	}
	// else: replay default is 0

	c := &Client{
		cfg:        cfg,
		ctx:        ctx,
		cancel:     cancel,
		log:        cfg.logger,
		sourceAddr: sourceAddr,
		deviceName: deviceName,
		addrReady:  make(chan struct{}),
	}

	if hasBus {
		if err := c.initBus(cfg); err != nil {
			cancel()
			return nil, err
		}
	} else {
		// Replay path: capture frames in memory.
		close(c.addrReady)

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
	}

	// Start the single writer goroutine for FIFO ordering.
	c.writeCh = make(chan writeJob, 64)
	c.writeWg.Add(1)
	go c.writeLoop()

	return c, nil
}

// initBus sets up the CAN bus integration: bus interface, address claiming,
// transport protocol, and the internal read/decode pipeline.
func (c *Client) initBus(cfg config) error {
	// Get or construct the bus.
	if cfg.bus != nil {
		c.bus = cfg.bus
	} else {
		// Construct from first source backed by real hardware.
		for _, src := range cfg.sources {
			if bb, ok := src.(busBacked); ok {
				c.bus = bb.newBus(c.log)
				break
			}
		}
	}
	if c.bus == nil {
		return errors.New("n2k: no writable bus available — File, TCP, and UDP sources are read-only; use Receive/NewScanner with them, or give the client a CAN/USB source or WithBus")
	}

	// Set writeFrame to delegate to the bus.
	c.writeFrame = func(f can.Frame) error {
		c.mu.Lock()
		closed := c.closed
		c.mu.Unlock()
		if closed {
			return errors.New("n2k: client closed")
		}
		return c.bus.WriteFrame(f)
	}

	// Create transport manager with OnComplete that feeds decoded TP messages
	// through the read pipeline.
	c.tp = transport.NewManager(transport.ManagerConfig{
		WriteFrame: c.writeFrame,
		OnComplete: func(tpPGN uint32, source uint8, destination uint8, data []byte) {
			priority := uint8(6)
			info := pgn.MessageInfo{Timestamp: time.Now(), PGN: tpPGN, SourceId: source, Priority: &priority}
			if destination != framer.BroadcastAddr {
				info.TargetId = &destination
			}
			c.pipeline.InjectAssembled(info, data)
		},
		Logger: c.log,
	})

	// Set up the persistent read pipeline (pre-filter -> assembly -> decode
	// -> unknown-PGN policy -> post-filter -> msgCh).
	c.msgCh = make(chan pgn.Message, 64)
	p, err := newReadPipeline(c.ctx, cfg, c.msgCh)
	if err != nil {
		return err
	}
	c.pipeline = p

	// Determine claiming mode.
	mode := claiming.ModeAuto
	if cfg.sourceAddress != nil {
		mode = claiming.ModeExplicit
	}

	// Create the claimer.
	c.claimer = claiming.New(claiming.Config{
		Mode:       mode,
		Address:    c.sourceAddr,
		Name:       c.deviceName,
		WriteFrame: c.writeFrame,
		OnAddressChange: func(newAddr uint8) {
			c.mu.Lock()
			c.sourceAddr = newAddr
			c.mu.Unlock()
		},
		OnFatalError: func(err error) {
			c.log.Error("address claiming fatal error", "error", err)
			c.mu.Lock()
			c.addrErr = err
			c.mu.Unlock()
		},
		Logger: c.log,
	})

	// Start the bus read loop goroutine.
	go c.busReadLoop()

	// Send the initial address claim.
	if err := c.claimer.Start(); err != nil {
		return fmt.Errorf("n2k: starting address claim: %w", err)
	}

	// Wait for the claim timeout to allow the network to respond.
	claimTimeout := defaultClaimTimeout
	if cfg.claimTimeout != nil {
		claimTimeout = *cfg.claimTimeout
	}

	timer := time.NewTimer(claimTimeout)
	defer timer.Stop()

	select {
	case <-timer.C:
	case <-c.ctx.Done():
		return c.ctx.Err()
	}

	// Check if a fatal error occurred during the claim window.
	c.mu.Lock()
	if c.addrErr != nil {
		err := c.addrErr
		c.mu.Unlock()
		return fmt.Errorf("n2k: address claiming failed: %w", err)
	}
	c.mu.Unlock()

	// Check if we got a valid address.
	addr := c.claimer.Address()
	if addr == 254 {
		return errors.New("n2k: address claiming failed — all addresses exhausted or contention in explicit mode")
	}

	// Update sourceAddr with the final claimed address.
	c.mu.Lock()
	c.sourceAddr = addr
	c.mu.Unlock()

	close(c.addrReady)
	return nil
}

// handleBusFrame is the central frame router called for every incoming CAN frame.
func (c *Client) handleBusFrame(frame can.Frame) {
	info := adapter.NewPacketInfo(&frame)

	// Route address claim frames (PGN 60928) to the claimer.
	if info.PGN == framer.PGNISOAddressClaim && frame.Length == 8 {
		name := binary.LittleEndian.Uint64(frame.Data[:])
		c.claimer.HandleAddressClaim(info.SourceId, name)
	}

	// Route ISO requests (PGN 59904) for address claim to the claimer.
	if info.PGN == framer.PGNISORequest && frame.Length >= 3 {
		requestedPGN := uint32(frame.Data[0]) | uint32(frame.Data[1])<<8 | uint32(frame.Data[2])<<16
		if requestedPGN == framer.PGNISOAddressClaim {
			c.claimer.HandleISORequest()
		}
	}

	// Route transport protocol frames to the TP manager.
	if info.PGN == transport.PGNCM || info.PGN == transport.PGNDT {
		c.tp.HandleFrame(frame)
	}

	// Decode for the read API using the persistent pipeline (pre-filter moved
	// inside HandleFrame).
	c.pipeline.HandleFrame(frame)
}

// busReadLoop runs the bus and closes the message channel when done.
func (c *Client) busReadLoop() {
	defer close(c.msgCh)
	if err := c.bus.Run(c.ctx, c.handleBusFrame); err != nil && c.ctx.Err() == nil {
		c.log.Error("bus read loop error", "error", err)
	}
}

// Write asynchronously encodes and transmits a PGN message. The message must
// be a pointer to a PGN struct (e.g. *pgn.VesselHeading). The returned
// WriteResult can be used to wait for completion and check for errors.
// Writes are serialized through a single goroutine to guarantee FIFO ordering.
func (c *Client) Write(msg pgn.Message) *WriteResult {
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
func (c *Client) doWrite(msg pgn.Message) error {
	// Wait for address readiness before sending.
	select {
	case <-c.addrReady:
	case <-c.ctx.Done():
		return c.ctx.Err()
	}

	pgnNum := msg.PGNNumber()

	// Check that msg implements pgn.PGN and extract its MessageInfo.
	pgnMsg, ok := msg.(pgn.PGN)
	if !ok {
		return fmt.Errorf("n2k: %T does not implement pgn.PGN", msg)
	}

	payload, err := pgn.EncodeMessage(msg)
	if err != nil {
		return fmt.Errorf("n2k: encode PGN %d: %w", pgnNum, err)
	}

	info := pgnMsg.MessageInfo()

	var priority uint8 = 6
	if info.Priority != nil {
		priority = *info.Priority
	}

	var destination uint8 = 255
	if info.TargetId != nil {
		destination = *info.TargetId
	}

	// Read sourceAddr under lock.
	c.mu.Lock()
	srcAddr := c.sourceAddr
	c.mu.Unlock()

	canID := framer.BuildCANID(pgnNum, priority, srcAddr, destination)

	isFast := false
	if infos, ok := pgn.PgnInfoLookup[pgnNum]; ok && len(infos) > 0 {
		isFast = infos[0].Fast
	}

	if !isFast && len(payload) <= 8 {
		frame := framer.FrameSingle(canID, payload)
		return c.writeFrame(frame)
	}

	if isFast && len(payload) <= 223 {
		seqID := uint8(c.fastSeq.Add(1) % 8)
		frames := framer.FrameFastPacket(canID, payload, seqID)
		for _, f := range frames {
			if err := c.writeFrame(f); err != nil {
				return err
			}
		}
		return nil
	}

	return c.tp.SendBAM(pgnNum, srcAddr, payload)
}

// Receive returns an iterator of decoded NMEA 2000 messages. For bus clients
// it reads from the internal message channel; for replay clients it builds a
// fresh Scanner over the client's config so each call gets a full replay.
func (c *Client) Receive() iter.Seq2[pgn.Message, error] {
	if c.msgCh != nil {
		return func(yield func(pgn.Message, error) bool) {
			for msg := range c.msgCh {
				if !yield(msg, nil) {
					return
				}
			}
		}
	}
	return func(yield func(pgn.Message, error) bool) {
		s := c.newReplayScanner()
		for s.Next() {
			if !yield(s.Message(), nil) {
				return
			}
		}
		if s.Err() != nil {
			yield(nil, s.Err())
		}
	}
}

// Scanner creates a new Scanner that reads from this client. For bus clients
// it reads from the internal message channel; for replay clients it builds a
// fresh Scanner over the client's config.
func (c *Client) Scanner() *Scanner {
	if c.msgCh != nil {
		s := &Scanner{
			ctx: c.ctx,
			cfg: c.cfg,
			ch:  c.msgCh,
		}
		// ch is already live (fed by the client's own read pipeline), so
		// suppress Next()'s lazy-start goroutine: firing once.Do here means
		// the real closure passed to Next()'s once.Do never runs.
		s.once.Do(func() {})
		return s
	}
	return c.newReplayScanner()
}

// newReplayScanner builds a Scanner over the client's already-parsed config,
// so replay clients share the exact pipeline code without re-parsing options.
func (c *Client) newReplayScanner() *Scanner {
	return &Scanner{ctx: c.ctx, cfg: c.cfg, ch: make(chan pgn.Message, 64)}
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
	if c.bus != nil {
		return c.bus.Close()
	}
	return nil
}
