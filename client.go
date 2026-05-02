package n2k

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/adapter"
	"github.com/open-ships/n2k/internal/canbus"
	"github.com/open-ships/n2k/internal/claiming"
	"github.com/open-ships/n2k/internal/decoder"
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

	// opts preserves the original Option slice so Receive and Scanner can
	// reconstruct the same configuration.
	opts []Option

	// mu guards writtenFrames, closed state, and sourceAddr.
	mu            sync.Mutex
	writtenFrames []can.Frame // captured frames (replay/testing)
	closed        bool

	bus     canbus.Interface
	claimer *claiming.Claimer
	addrErr error // set by OnFatalError during address claiming

	// addrReady is closed once address claiming completes (or immediately for replay).
	addrReady chan struct{}

	// msgCh delivers decoded messages from the internal read loop to the read API.
	// nil for replay clients.
	msgCh chan pgn.Message

	// readAdapter and readDecoder are the persistent decode pipeline for the
	// internal read loop. Only used for bus clients.
	readAdapter *adapter.CANAdapter
	readDecoder *decoder.Decoder

	// readFilter holds the compiled CEL filter for the read API (nil if no filter).
	readFilter *filter
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

	// Compile CEL filter early if configured.
	var readFilter *filter
	if cfg.filterExpr != "" {
		var err error
		readFilter, err = compileFilter(cfg.filterExpr)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("n2k: compiling filter: %w", err)
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
		addrReady:  make(chan struct{}),
		readFilter: readFilter,
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
	// Get or construct the bus interface.
	if cfg.bus != nil {
		c.bus = cfg.bus
	} else {
		// Construct from first non-replay source.
		for _, src := range cfg.sources {
			switch s := src.(type) {
			case *socketCANSource:
				c.bus = canbus.NewSocketCAN(c.log, s.iface, c.handleBusFrame)
			case *usbCANSource:
				c.bus = canbus.NewUSBCAN(c.log, s.port, c.handleBusFrame)
			}
			if c.bus != nil {
				break
			}
		}
	}
	if c.bus == nil {
		return errors.New("n2k: could not construct bus from sources")
	}

	// If the bus supports setting the handler after construction, set it now.
	if hs, ok := c.bus.(canbus.HandlerSettable); ok {
		hs.SetHandler(c.handleBusFrame)
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
			// Build a synthetic CAN frame for the reassembled TP message
			// and feed it through the decode pipeline.
			canID := framer.BuildCANID(tpPGN, 6, source, destination)
			f := can.Frame{
				ID:     canID,
				Length: uint8(min(len(data), 8)),
			}
			copy(f.Data[:], data)
			// The adapter + decoder pipeline handles the rest.
			c.readAdapter.HandleMessage(&f)
		},
		Logger: c.log,
	})

	// Set up the persistent decode pipeline.
	c.readAdapter = adapter.NewCANAdapter()
	c.readDecoder = decoder.New()
	c.readDecoder.SetOutput(&clientDecoderHandler{client: c})
	c.readAdapter.SetOutput(c.readDecoder)

	// Create the message channel for the read API.
	c.msgCh = make(chan pgn.Message, 64)

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
	if info.PGN == 60928 && frame.Length == 8 {
		name := binary.LittleEndian.Uint64(frame.Data[:])
		c.claimer.HandleAddressClaim(info.SourceId, name)
	}

	// Route ISO requests (PGN 59904) for address claim to the claimer.
	if info.PGN == 59904 && frame.Length >= 3 {
		requestedPGN := uint32(frame.Data[0]) | uint32(frame.Data[1])<<8 | uint32(frame.Data[2])<<16
		if requestedPGN == 60928 {
			c.claimer.HandleISORequest()
		}
	}

	// Route transport protocol frames to the TP manager.
	if info.PGN == 60416 || info.PGN == 60160 {
		c.tp.HandleFrame(frame)
	}

	// Pre-filter: skip decoding if metadata doesn't match.
	if c.readFilter != nil && !c.readFilter.evalPre(info) {
		return
	}

	// Decode for the read API using the persistent pipeline.
	c.readAdapter.HandleMessage(&frame)
}

// busReadLoop runs the bus and closes the message channel when done.
func (c *Client) busReadLoop() {
	defer close(c.msgCh)
	if err := c.bus.Run(c.ctx); err != nil && c.ctx.Err() == nil {
		c.log.Error("bus read loop error", "error", err)
	}
}

// clientDecoderHandler receives decoded messages from the decoder pipeline and
// delivers them to the client's msgCh.
type clientDecoderHandler struct {
	client *Client
}

func (h *clientDecoderHandler) HandleStruct(msg pgn.Message) {
	if msg == nil {
		return
	}

	if u, ok := msg.(*pgn.UnknownPGN); ok {
		if !h.client.cfg.includeUnknown {
			h.client.cfg.logger.Debug("dropping unknown PGN",
				"pgn", u.Info.PGN, "reason", u.Reason)
			return
		}
	}

	if h.client.readFilter != nil && h.client.readFilter.hasPost {
		fields := structToFilterMap(msg)
		rv := reflect.ValueOf(msg)
		if rv.Kind() == reflect.Pointer {
			rv = rv.Elem()
		}
		info := rv.FieldByName("Info").Interface().(pgn.MessageInfo)
		if !h.client.readFilter.evalPostWithInfo(info, fields) {
			return
		}
	}

	select {
	case h.client.msgCh <- msg:
	case <-h.client.ctx.Done():
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

	encoder, ok := pgn.EncoderLookup[pgnNum]
	if !ok {
		return fmt.Errorf("n2k: no encoder registered for PGN %d", pgnNum)
	}

	payload, err := encoder(msg)
	if err != nil {
		return fmt.Errorf("n2k: encode PGN %d: %w", pgnNum, err)
	}

	// Extract MessageInfo from the struct's exported Info field.
	info := reflect.ValueOf(msg).Elem().FieldByName("Info").Interface().(pgn.MessageInfo)

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

	if len(payload) <= 8 {
		frame := framer.FrameSingle(canID, payload)
		return c.writeFrame(frame)
	}

	isFast := false
	if infos, ok := pgn.PgnInfoLookup[pgnNum]; ok && len(infos) > 0 {
		isFast = infos[0].Fast
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
// it reads from the internal message channel; for replay clients it delegates
// to the top-level Receive function.
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
	return Receive(c.ctx, c.opts...)
}

// Scanner creates a new Scanner that reads from this client. For bus clients
// it reads from the internal message channel; for replay clients it delegates
// to NewScanner.
func (c *Client) Scanner() *Scanner {
	if c.msgCh != nil {
		s := &Scanner{
			ctx: c.ctx,
			cfg: c.cfg,
			ch:  c.msgCh,
		}
		s.once.Do(func() {}) // prevent lazy start since ch is already live
		return s
	}
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
	if c.bus != nil {
		return c.bus.Close()
	}
	return nil
}
