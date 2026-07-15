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

// Client is a read/write NMEA 2000 bus node. It composes address claiming,
// transport protocol, encoding, and framing into a single type that can both
// receive and transmit PGN messages.
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

	// system decodes protocol PGNs (product info, group functions, request
	// responses) independently of the user filter. Only used for bus clients.
	system *systemRouter

	// heartbeat transmits PGN 126993 periodically. Only set for bus clients
	// (nil for replay clients).
	heartbeat *heartbeater

	// productInfo and configInfo identify this device to the network
	// (PGNs 126996 and 126998).
	productInfo ProductInfo
	configInfo  ConfigInfo

	// correlator matches system messages to in-flight Request calls. Only
	// set for bus clients.
	correlator *correlator

	// bMu guards broadcasters, the active periodic transmissions by PGN.
	bMu          sync.Mutex
	broadcasters map[uint32]*broadcaster

	// registry tracks devices observed on the bus, keyed by NAME. Only set
	// for bus clients.
	registry *registry
}

type writeJob struct {
	msg    pgn.Message
	result *WriteResult
}

// NewClient creates a Client that can read and write NMEA 2000 messages.
// Provide CAN, USB, TCP, Replay, or WithBus for writable clients; File and UDP
// are read-only sources for Receive/NewScanner.
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

	c.productInfo = defaultProductInfo(UnpackDeviceName(deviceName).IdentityNumber)
	if cfg.productInfo != nil {
		c.productInfo = *cfg.productInfo
	}
	c.configInfo = defaultConfigInfo()
	if cfg.configInfo != nil {
		c.configInfo = *cfg.configInfo
	}

	// Start the single writer goroutine for FIFO ordering. It must exist
	// before initBus so the client's own protocol goroutines (heartbeat,
	// info responses) can write as soon as the address claim completes.
	c.writeCh = make(chan writeJob, 64)
	c.writeWg.Add(1)
	go c.writeLoop()

	if hasBus {
		if err := c.initBus(cfg); err != nil {
			cancel()
			close(c.writeCh)
			c.writeWg.Wait()
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

	return c, nil
}

// initBus sets up the bus runtime: bus interface, address claiming, transport
// protocol, and the internal read/decode pipeline.
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
		return errors.New("n2k: no writable bus available — File and UDP sources are read-only; use Receive/NewScanner with them, or give the client a CAN, USB, or TCP source or WithBus")
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
			c.system.handleAssembled(info, data)
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

	// Set up the protocol-message router and start its dispatch loop.
	sys, err := newSystemRouter(c.ctx, cfg)
	if err != nil {
		return err
	}
	c.system = sys
	c.correlator = newCorrelator()
	c.registry = newRegistry()
	sys.addHandler(c.correlator.observe)
	sys.addHandler(c.handleGroupFunction)
	sys.addHandler(c.registry.observe)
	go sys.run()

	// Set up the heartbeat (PGN 126993). It waits for addrReady before its
	// first transmission.
	hbInterval := defaultHeartbeatInterval
	if cfg.heartbeatInterval != nil {
		hbInterval = *cfg.heartbeatInterval
	}
	c.heartbeat = newHeartbeater(hbInterval, c.Write)
	go c.heartbeat.run(c.ctx, c.addrReady)

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

	// Enumerate the bus: ask every device to (re-)announce its address claim
	// so the registry fills without waiting for spontaneous traffic.
	enumerate := uint64(framer.PGNISOAddressClaim)
	c.Write(&pgn.IsoRequest{Pgn: &enumerate})

	return nil
}

// requestDeviceInfo asks a newly seen device for its product and
// configuration info once the client itself is ready to transmit.
func (c *Client) requestDeviceInfo(addr uint8) {
	go func() {
		select {
		case <-c.addrReady:
		case <-c.ctx.Done():
			return
		}
		for _, requested := range []uint32{126996, 126998} {
			pgnNum := uint64(requested)
			c.Write(&pgn.IsoRequest{
				Info: pgn.MessageInfo{TargetId: pgn.Target(addr)},
				Pgn:  &pgnNum,
			})
		}
	}()
}

// handleBusFrame is the central frame router called for every incoming CAN frame.
func (c *Client) handleBusFrame(frame can.Frame) {
	info := adapter.NewPacketInfo(&frame)

	// Route address claim frames (PGN 60928) to the claimer and the device
	// registry.
	if info.PGN == framer.PGNISOAddressClaim && frame.Length == 8 {
		name := binary.LittleEndian.Uint64(frame.Data[:])
		c.claimer.HandleAddressClaim(info.SourceId, name)
		if c.registry.handleClaim(info.SourceId, name, info.Timestamp) {
			c.requestDeviceInfo(info.SourceId)
		}
	} else {
		c.registry.touch(info.SourceId, info.Timestamp)
	}

	// Route ISO requests (PGN 59904) to the request responder.
	if info.PGN == framer.PGNISORequest {
		c.handleISORequest(info, frame)
	}

	// Route transport protocol frames to the TP manager.
	if info.PGN == transport.PGNCM || info.PGN == transport.PGNDT {
		c.tp.HandleFrame(frame)
	}

	// Decode for the read API using the persistent pipeline (pre-filter moved
	// inside HandleFrame).
	c.pipeline.HandleFrame(frame)

	// Decode protocol PGNs for the client's own use, independent of the user
	// filter.
	c.system.handleFrame(frame, info.PGN)
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

	// Message-oriented buses (see MessageWriter) take whole payloads and do
	// their own wire framing; anything larger than one message still goes
	// frame-by-frame below via ISO-TP.
	if mw, ok := c.bus.(MessageWriter); ok && len(payload) <= 223 {
		return mw.WriteMessage(pgnNum, priority, srcAddr, destination, payload)
	}

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

	// Stop protocol writers before closing the write channel so their final
	// writes cannot race the shutdown.
	if c.heartbeat != nil {
		c.heartbeat.stop()
	}
	c.stopBroadcasters()

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
