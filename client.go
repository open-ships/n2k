package n2k

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"reflect"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/claiming"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/internal/transport"
	"github.com/open-ships/n2k/pgn"
	"github.com/open-ships/n2k/raw"
)

// defaultClaimTimeout is the maximum time NewClient blocks waiting for address
// claiming to complete.
const defaultClaimTimeout = 1500 * time.Millisecond

// defaultReadyTimeout bounds transport opening and readiness negotiation. It
// must remain longer than built-in gateway command deadlines so their typed
// errors reach NewClient before this outer guard fires.
const defaultReadyTimeout = 5 * time.Second

const defaultWriteQueue = 64

const pgnISOCommandedAddress uint32 = 65240

var (
	// ErrWriteQueueFull reports that an asynchronous write could not be admitted
	// without blocking the caller. Callers may retry or apply their own policy.
	ErrWriteQueueFull = errors.New("n2k: write queue full")
	// ErrClientClosed reports an operation attempted after Client.Close.
	ErrClientClosed = errors.New("n2k: client closed")
)

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
	// protocolTx owns priority admission for automatic protocol traffic.
	protocolTx   *protocolTransmitter
	wire         *wireTransmitter
	epochCtx     context.Context
	epochCancel  context.CancelCauseFunc
	dispatchMu   sync.Mutex
	deviceInfoCh chan deviceInfoRequest
	rejoinCh     chan rejoinRequest
	backgroundWg sync.WaitGroup

	// mu guards writtenFrames, closed state, and sourceAddr.
	mu                   sync.Mutex
	writtenFrames        []can.Frame // bounded capture ring (replay/testing)
	writtenFramesStart   int
	writtenFramesDropped uint64
	closed               bool
	terminalErr          error
	claimed              bool
	claimEpoch           uint64
	txReady              chan struct{}
	connected            bool
	connectionEpoch      uint64
	connectionChanged    chan struct{}
	rejoining            bool
	rejoinMu             sync.Mutex

	bus        Bus
	busDone    chan struct{}
	closeDone  chan struct{}
	closeErr   error
	busStarted bool
	claimer    *claiming.Claimer
	addrErr    error // set by OnFatalError during address claiming

	// addrReady is closed once address claiming completes (or immediately for replay).
	addrReady chan struct{}

	// msgHub delivers decoded live messages without allowing application
	// backpressure to stall protocol processing. It is nil for replay clients.
	msgHub *messageHub
	// observationHub independently delivers owned transport observations.
	observationHub *observationHub

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
	// actisenseRemote owns destination- and epoch-safe BEM correlation over
	// proprietary PGN 126720.
	actisenseRemote *actisenseRemoteManager

	// bMu guards broadcasters, the active periodic transmissions by PGN.
	bMu              sync.Mutex
	broadcasters     map[uint32]*broadcaster
	broadcastWorkers map[*broadcaster]struct{}
	broadcastSeq     atomic.Uint64

	// registry tracks devices observed on the bus, keyed by NAME. Only set
	// for bus clients.
	registry *registry

	applicationWritesAccepted  atomic.Uint64
	applicationWritesCompleted atomic.Uint64
	applicationWritesFailed    atomic.Uint64
	applicationWritesRejected  atomic.Uint64
	framesReceived             atomic.Uint64
	framesTransmitted          atomic.Uint64
	messagesObserved           atomic.Uint64
	decodeErrorsObserved       atomic.Uint64
	gatewayEventsObserved      atomic.Uint64
	transportErrorsObserved    atomic.Uint64
	deviceInfoRequestsDropped  atomic.Uint64
}

type writeJob struct {
	msg           pgn.Message
	result        *WriteResult
	protocol      bool
	protocolClass protocolWriteClass
	operation     string
	ctx           context.Context
	stop          func()
}

type runtimeBusError struct{ err error }

func (e *runtimeBusError) Error() string { return "n2k: bus write: " + e.err.Error() }
func (e *runtimeBusError) Unwrap() error { return e.err }

func wrapRuntimeBusError(err error) error {
	if err == nil {
		return nil
	}
	return &runtimeBusError{err: err}
}

func validateClientConfig(cfg config) error {
	if cfg.bus != nil && len(cfg.sources) > 0 {
		return errors.New("n2k: WithBus cannot be combined with source options")
	}
	busSources := 0
	for _, src := range cfg.sources {
		switch value := src.(type) {
		case *tcpSource:
			if value.format == FormatActisense || value.format == FormatActisenseN2KASCII {
				return ErrActisenseGatewaySessionRequired
			}
		case *serialSource:
			if value.format == FormatActisense || value.format == FormatActisenseN2KASCII {
				return ErrActisenseGatewaySessionRequired
			}
		}
		if _, ok := src.(busBacked); ok {
			busSources++
		}
	}
	if busSources > 1 || (busSources == 1 && len(cfg.sources) != 1) {
		return errors.New("n2k: Client requires exactly one writable bus source")
	}
	if cfg.bus == nil && busSources == 0 {
		for _, src := range cfg.sources {
			if _, ok := src.(*replaySource); !ok {
				return errors.New("n2k: File, EBL, and UDP are read-only; use Receive or NewScanner")
			}
		}
	}
	if cfg.productInfo != nil {
		fields := map[string]string{
			"model ID":         cfg.productInfo.ModelID,
			"software version": cfg.productInfo.SoftwareVersion,
			"model version":    cfg.productInfo.ModelVersion,
			"serial number":    cfg.productInfo.SerialNumber,
		}
		for name, value := range fields {
			if len([]byte(value)) > 32 {
				return fmt.Errorf("n2k: product %s exceeds 32 bytes", name)
			}
		}
	}
	if cfg.configInfo != nil {
		for name, value := range map[string]string{
			"installation description 1": cfg.configInfo.InstallationDescription1,
			"installation description 2": cfg.configInfo.InstallationDescription2,
			"manufacturer information":   cfg.configInfo.ManufacturerInformation,
		} {
			if len([]byte(value)) > 253 {
				return fmt.Errorf("n2k: configuration %s exceeds 253 bytes", name)
			}
		}
	}
	return nil
}

// NewClient creates a Client that can read and write NMEA 2000 messages.
// Provide CAN, USB, Serial, TCP, Replay, or WithBus for writable clients;
// File, EBL, and UDP are read-only sources for Receive/NewScanner.
func NewClient(ctx context.Context, opts ...Option) (*Client, error) {
	cfg := config{}
	for _, o := range opts {
		if o == nil {
			return nil, errors.New("n2k: nil Option")
		}
		o.apply(&cfg)
	}
	cfg.applyReconnect()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	if cfg.logger == nil {
		cfg.logger = slog.Default()
	}
	if err := validateClientConfig(cfg); err != nil {
		return nil, err
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
	arbitraryAddressCapable := cfg.sourceAddress == nil
	if cfg.deviceName != nil {
		deviceName = cfg.deviceName.Pack(arbitraryAddressCapable)
	} else {
		deviceName = DefaultDeviceName().Pack(arbitraryAddressCapable)
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
	} else if cfg.preferredAddress != nil {
		sourceAddr = *cfg.preferredAddress
	} else if hasBus {
		sourceAddr = 251
	}
	// else: replay default is 0

	c := &Client{
		cfg:               cfg,
		ctx:               ctx,
		cancel:            cancel,
		log:               cfg.logger,
		sourceAddr:        sourceAddr,
		deviceName:        deviceName,
		addrReady:         make(chan struct{}),
		busDone:           make(chan struct{}),
		closeDone:         make(chan struct{}),
		connectionChanged: make(chan struct{}),
	}
	c.txReady = c.addrReady
	c.epochCtx, c.epochCancel = context.WithCancelCause(ctx)

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
	writeQueue := defaultWriteQueue
	if cfg.writeQueue != nil {
		writeQueue = *cfg.writeQueue
	}
	c.writeCh = make(chan writeJob, writeQueue)
	c.protocolTx = newProtocolTransmitter(c.log)
	c.writeWg.Add(2)
	go c.writeLoop()
	go c.protocolWriteLoop()

	if hasBus {
		if err := c.initBus(cfg); err != nil {
			// cancel() unblocks the protocol goroutines so none issues a new
			// write; closing the bus releases any write still parked waiting for
			// an auto-reconnect (e.g. the initial claim against an unreachable
			// gateway) and stops the read loop, so teardown cannot deadlock.
			// Marking closed and draining in-flight senders before closing the
			// channel keeps the teardown symmetric with Close.
			cancel()
			c.mu.Lock()
			c.closed = true
			c.terminalErr = err
			c.mu.Unlock()
			if c.bus != nil {
				_ = c.bus.Close()
			}
			if c.tp != nil {
				c.tp.Close()
			}
			c.writeWg.Wait()
			if c.wire != nil {
				<-c.wire.done
			}
			if c.busStarted {
				<-c.busDone
			}
			if c.system != nil {
				<-c.system.done
			}
			c.backgroundWg.Wait()
			return nil, err
		}
	} else {
		// Replay path: capture a bounded recent frame window in memory.
		c.claimed = true
		close(c.addrReady)

		c.writeFrame = func(f can.Frame) error {
			c.mu.Lock()
			defer c.mu.Unlock()
			if c.closed {
				return errors.New("n2k: client closed")
			}
			c.captureReplayFrameLocked(f)
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
		return errors.New("n2k: no writable bus available — File, EBL, and UDP sources are read-only; use Receive/NewScanner with them, or give the client a CAN, USB, Serial, or TCP source or WithBus")
	}

	// Set writeFrame to delegate to the bus.
	c.wire = newWireTransmitter(c)
	c.writeFrame = func(f can.Frame) error { return c.writeFrameContext(c.ctx, f) }

	// Create transport manager with OnComplete that feeds decoded TP messages
	// through the read pipeline.
	c.tp = transport.NewManager(transport.ManagerConfig{
		WriteFrame:        c.writeFrame,
		WriteFrameContext: c.writeFrameContext,
		LocalAddress: func() uint8 {
			c.mu.Lock()
			defer c.mu.Unlock()
			return c.sourceAddr
		},
		OnCompleteInfo: func(info pgn.MessageInfo, data []byte) {
			if !c.currentMessageEpoch(info) {
				return
			}
			destination := uint8(framer.BroadcastAddr)
			if info.TargetId != nil {
				destination = *info.TargetId
			}
			if info.PGN == pgnISOCommandedAddress {
				c.handleCommandedAddressTransfer(destination, data)
			}
			c.system.handleAssembled(info, data)
			c.pipeline.InjectAssembled(info, data)
		},
		Logger: c.log,
	})

	// Set up the persistent read pipeline (pre-filter -> assembly -> decode ->
	// unknown-PGN policy -> post-filter -> non-blocking subscription hub).
	receiveBuffer := defaultReceiveBuffer
	if cfg.receiveBuffer != nil {
		receiveBuffer = *cfg.receiveBuffer
	}
	c.msgHub = newMessageHub(receiveBuffer)
	c.observationHub = newObservationHub(receiveBuffer)
	p, err := newReadPipeline(c.ctx, cfg, c.msgHub.publish)
	if err != nil {
		return err
	}
	c.pipeline = p
	c.pipeline.setObservationOutput(c.publishObservation)

	// Set up the protocol-message router and start its dispatch loop.
	sys, err := newSystemRouter(c.ctx, cfg)
	if err != nil {
		return err
	}
	c.system = sys
	c.correlator = newCorrelator()
	c.actisenseRemote = newActisenseRemoteManager(c)
	c.registry = newRegistry()
	sys.addHandler(c.correlator.observe)
	sys.addHandler(c.handleGroupFunction)
	sys.addHandler(c.registry.observe)
	sys.addObservationHandler(c.actisenseRemote.observe)
	sys.dispatchMu = &c.dispatchMu
	sys.current = c.currentMessageEpoch
	sys.onError = c.fail
	go sys.run()
	c.deviceInfoCh = make(chan deviceInfoRequest, 252)
	c.rejoinCh = make(chan rejoinRequest, 1)
	c.backgroundWg.Add(2)
	go c.deviceInfoLoop()
	go c.rejoinLoop()

	// Set up the heartbeat (PGN 126993). It waits for addrReady before its
	// first transmission.
	hbInterval := defaultHeartbeatInterval
	if cfg.heartbeatInterval != nil {
		hbInterval = *cfg.heartbeatInterval
	}
	c.heartbeat = newHeartbeater(hbInterval, func(msg pgn.Message) *WriteResult {
		return c.writeProtocol("heartbeat", protocolRequired, msg)
	})
	go func() {
		defer c.recoverGoroutine("heartbeat")
		c.heartbeat.run(c.ctx, c.addrReady)
	}()

	if lifecycleBus, ok := c.bus.(ConnectionLifecycleBus); ok {
		lifecycleBus.SetConnectionObserver(c.handleConnectionChange)
	} else {
		c.mu.Lock()
		c.connected = true
		c.connectionEpoch = 1
		c.mu.Unlock()
		c.resetReadEpoch()
	}

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
			c.handleAddressChange(newAddr)
		},
		OnFatalError: func(err error) {
			c.log.Error("address claiming fatal error", "error", err)
			c.mu.Lock()
			c.addrErr = err
			active := c.claimed
			c.mu.Unlock()
			if active {
				c.fail(fmt.Errorf("n2k: address claiming failed: %w", err))
			}
		},
		Logger: c.log,
	})

	// Start the bus read loop goroutine.
	c.busStarted = true
	go c.busReadLoop()

	claimTimeout := defaultClaimTimeout
	if cfg.claimTimeout != nil {
		claimTimeout = *cfg.claimTimeout
	}
	readyTimeout := defaultReadyTimeout
	if cfg.readyTimeout != nil {
		readyTimeout = *cfg.readyTimeout
	}
	readyTimer := time.NewTimer(readyTimeout)
	defer readyTimer.Stop()
	if readyBus, ok := c.bus.(ReadyBus); ok {
		select {
		case <-readyBus.Ready():
		case <-c.busDone:
			return c.currentTerminalError("n2k: bus stopped before becoming ready")
		case <-readyTimer.C:
			return errors.New("n2k: bus did not become ready within ready timeout")
		case <-c.ctx.Done():
			return c.currentTerminalError(c.ctx.Err().Error())
		}
	}
	// Lifecycle transports must identify the actual connection before any
	// claim is admitted. Readiness and connection opening share one deadline.
	var initialConnectionEpoch uint64
	for {
		c.mu.Lock()
		if c.connected {
			initialConnectionEpoch = c.connectionEpoch
			c.mu.Unlock()
			break
		}
		changed := c.connectionChanged
		c.mu.Unlock()
		select {
		case <-changed:
		case <-c.busDone:
			return c.currentTerminalError("n2k: bus stopped before connecting")
		case <-readyTimer.C:
			return errors.New("n2k: bus did not connect within ready timeout")
		case <-c.ctx.Done():
			return c.currentTerminalError(c.ctx.Err().Error())
		}
	}
	readyTimer.Stop()

	// Send the initial address claim, bounded by the claim deadline. With an
	// auto-reconnect policy the first write blocks until the gateway connection
	// is established, so a gateway that is unreachable at startup must fail
	// NewClient rather than hang it forever. Start the deadline only after the
	// claim goroutine is scheduled: scheduler delay is not gateway delay and can
	// otherwise consume a short claim window on a loaded machine.
	startErr := make(chan error, 1)
	claimStarted := make(chan struct{})
	c.backgroundWg.Add(1)
	go func() {
		defer c.backgroundWg.Done()
		// This goroutine is owned by n2k, so a panic from a misbehaving Bus
		// (e.g. WriteFrame panicking instead of returning an error) would
		// otherwise be unrecoverable by the caller and crash the whole
		// process. Convert it into the same error return the caller already
		// handles for a claim that fails to start.
		defer func() {
			if r := recover(); r != nil {
				c.log.Error("panic during initial address claim",
					"panic", r, "stack", string(debug.Stack()))
				startErr <- fmt.Errorf("n2k: panic during address claim: %v", r)
			}
		}()
		close(claimStarted)
		startErr <- c.claimer.Start()
	}()
	select {
	case <-claimStarted:
	case <-c.ctx.Done():
		return c.currentTerminalError(c.ctx.Err().Error())
	}

	timer := time.NewTimer(claimTimeout)
	defer timer.Stop()
	select {
	case err := <-startErr:
		if err != nil {
			return fmt.Errorf("n2k: starting address claim: %w", err)
		}
	case <-timer.C:
		return errors.New("n2k: gateway unreachable — initial connection not established within claim timeout")
	case <-c.ctx.Done():
		return c.currentTerminalError(c.ctx.Err().Error())
	}

	// Wait for the remainder of the claim window to allow the network to respond.
	select {
	case <-timer.C:
	case <-c.ctx.Done():
		return c.currentTerminalError(c.ctx.Err().Error())
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
	if !c.connected || c.connectionEpoch != initialConnectionEpoch {
		c.mu.Unlock()
		return fmt.Errorf("n2k: connection changed during initial address claim: %w", ErrEpochChanged)
	}
	c.sourceAddr = addr
	c.claimed = true
	c.mu.Unlock()

	close(c.addrReady)

	// Enumerate the bus: ask every device to (re-)announce its address claim
	// so the registry fills without waiting for spontaneous traffic.
	enumerate := uint64(framer.PGNISOAddressClaim)
	c.writeProtocol("bus enumeration", protocolAdvisory, &pgn.IsoRequest{Pgn: &enumerate})

	return nil
}

// requestDeviceInfo asks a newly seen device for its product and
// configuration info once the client itself is ready to transmit.
func (c *Client) requestDeviceInfo(addr uint8) {
	device, exists := c.DeviceAt(addr)
	if !exists {
		return
	}
	c.mu.Lock()
	request := deviceInfoRequest{address: addr, name: device.RawName, connection: c.connectionEpoch, claim: c.claimEpoch}
	c.mu.Unlock()
	select {
	case c.deviceInfoCh <- request:
	default:
		c.deviceInfoRequestsDropped.Add(1)
		c.log.Warn("device information queue full", "address", addr)
	}
}

// handleBusFrame is the central frame router called for every incoming CAN frame.
func (c *Client) handleBusFrame(frame can.Frame) {
	c.handleBusObservation(frameObservation(frame, "bus", "bus", raw.DirectionReceived))
}

// handleBusObservation is the central transport router. The frame observation
// is published before synchronous protocol handling.
func (c *Client) handleBusObservation(observation raw.Observation) {
	observation = normalizeObservation(observation)
	c.mu.Lock()
	observation.ConnectionEpoch, observation.ClaimEpoch = c.connectionEpoch, c.claimEpoch
	c.mu.Unlock()
	c.publishObservation(observation)
	if observation.Frame == nil {
		if observation.Kind == raw.KindMessage {
			c.system.handleObservation(observation, observation.PGN)
			c.pipeline.HandleObservation(observation)
		}
		return
	}
	frame := *observation.Frame
	info := messageInfoForObservation(observation)

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
		c.tp.HandleFrameWithInfo(frame, info)
	}

	// Decode protocol PGNs for the client's own use, independent of the user
	// filter.
	c.system.handleObservation(observation, info.PGN)

	// Decode for the read API after all synchronous protocol routing. User
	// delivery itself is non-blocking through msgHub.
	c.pipeline.HandleObservation(observation)
}

// busReadLoop runs the bus and closes the message channel when done.
// recoverGoroutine converts a panic on an n2k-owned goroutine into a logged
// error with a stack trace, so a fault in user-supplied code (a Bus method or
// message handler) degrades the client instead of crashing the host process.
// Use as `defer c.recoverGoroutine(name)` at the top of a goroutine, or around
// a single unit of work that should not abort the surrounding loop.
func (c *Client) recoverGoroutine(name string) {
	if r := recover(); r != nil {
		err := fmt.Errorf("n2k: panic in %s: %v", name, r)
		c.log.Error("recovered panic in "+name,
			"panic", r, "stack", string(debug.Stack()))
		if c.ctx.Err() == nil {
			c.fail(err)
		}
	}
}

func (c *Client) busReadLoop() {
	defer close(c.busDone)
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("n2k: panic in bus read loop: %v", r)
			c.log.Error("recovered panic in bus read loop", "panic", r, "stack", string(debug.Stack()))
			c.fail(err)
		}
	}()

	// Handling an inbound frame can write synchronously on this goroutine — the
	// claimer defends its address by re-broadcasting a claim in response to a
	// contender or ISO request — so a panicking Bus surfaces here, not only on
	// the write loop. Recover per frame so one bad frame is logged and skipped
	// rather than tearing down the read loop (and with it the whole process).
	var handleMu sync.Mutex
	handle := func(observation raw.Observation) {
		handleMu.Lock()
		defer handleMu.Unlock()
		defer c.recoverGoroutine("bus frame handler")
		c.handleBusObservation(observation)
	}
	var err error
	if observationBus, ok := c.bus.(ObservationBus); ok {
		err = observationBus.RunObservations(c.ctx, handle)
	} else {
		err = c.bus.Run(c.ctx, func(frame can.Frame) {
			handle(frameObservation(frame, "bus", "bus", raw.DirectionReceived))
		})
	}
	err = wrapActisenseModeError(err)
	if c.ctx.Err() == nil {
		if err == nil {
			err = errors.New("bus stopped unexpectedly")
		}
		wrapped := fmt.Errorf("n2k: bus read loop: %w", err)
		c.log.Error("bus read loop error", "error", err)
		c.fail(wrapped)
		return
	}
	if c.msgHub != nil {
		c.msgHub.close(nil)
	}
	if c.observationHub != nil {
		c.observationHub.close(nil)
	}
}

func (c *Client) currentTerminalError(fallback string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.terminalErr != nil {
		return c.terminalErr
	}
	return errors.New(fallback)
}

func (c *Client) fail(err error) {
	if err == nil {
		return
	}
	c.mu.Lock()
	if c.terminalErr == nil {
		c.terminalErr = err
	}
	c.connected = false
	c.rejoining = false
	if c.epochCancel != nil {
		c.epochCancel(err)
	}
	if c.correlator != nil {
		c.correlator.invalidate(err)
	}
	c.mu.Unlock()
	if c.actisenseRemote != nil {
		c.actisenseRemote.close(err)
	}
	if c.msgHub != nil {
		c.msgHub.close(err)
	}
	if c.observationHub != nil {
		c.observationHub.close(err)
	}
	c.cancel()
}

// handleConnectionChange gates application and protocol writes across TCP
// connection epochs. The gateway calls the connected notification before it
// publishes the new connection to blocked writers, closing the reconnect race.
func (c *Client) handleConnectionChange(connected bool, epoch uint64) {
	c.dispatchMu.Lock()
	defer c.dispatchMu.Unlock()
	c.mu.Lock()
	if c.closed || c.terminalErr != nil || epoch < c.connectionEpoch {
		c.mu.Unlock()
		return
	}
	previousConnectionEpoch := c.connectionEpoch
	changed := c.connected != connected || c.connectionEpoch != epoch
	c.connected = connected
	c.connectionEpoch = epoch
	if changed && c.connectionChanged != nil {
		close(c.connectionChanged)
		c.connectionChanged = make(chan struct{})
	}
	if c.actisenseRemote != nil {
		// Hold the Client epoch lock through invalidation. Remote request
		// registration takes the same Client -> manager lock order, so a
		// request is atomically assigned to either the old or new epoch.
		c.actisenseRemote.invalidate(ErrActisenseRemoteEpochChanged)
	}
	active := c.claimed
	if !active {
		if changed && previousConnectionEpoch != 0 {
			c.invalidateEpochLocked(ErrEpochChanged)
			if c.registry != nil {
				c.registry.reset()
			}
		}
		c.mu.Unlock()
		if changed {
			c.resetReadEpoch()
		}
		return
	}
	if !connected {
		c.invalidateEpochLocked(ErrEpochChanged)
		c.txReady = make(chan struct{})
		c.rejoining = true
		if c.registry != nil {
			c.registry.reset()
		}
		c.mu.Unlock()
		c.resetReadEpoch()
		return
	}
	if !c.rejoining {
		c.invalidateEpochLocked(ErrEpochChanged)
		c.txReady = make(chan struct{})
	}
	ready := c.txReady
	c.rejoining = true
	c.mu.Unlock()

	c.resetReadEpoch()
	c.queueRejoin(rejoinRequest{epoch: epoch, ready: ready})
}

func (c *Client) rejoinNetwork(epoch uint64, reconnectReady chan struct{}) {
	c.rejoinMu.Lock()
	defer c.rejoinMu.Unlock()

	c.mu.Lock()
	if c.closed || c.terminalErr != nil || !c.connected || c.connectionEpoch != epoch {
		c.mu.Unlock()
		return
	}
	epochCtx := c.epochCtx
	c.mu.Unlock()

	if err := c.claimer.Reclaim(); err != nil {
		if epochCtx.Err() == nil {
			c.fail(fmt.Errorf("n2k: reclaiming address after reconnect: %w", err))
		}
		return
	}

	timer := time.NewTimer(250 * time.Millisecond)
	select {
	case <-timer.C:
	case <-c.ctx.Done():
		timer.Stop()
		return
	}

	// Contention can replace a readiness gate within this connection. Never
	// park on an obsolete gate: its epoch cancellation wakes us to recheck.
	for {
		c.mu.Lock()
		if c.closed || c.terminalErr != nil || !c.connected || c.connectionEpoch != epoch {
			c.mu.Unlock()
			return
		}
		ready := c.txReady
		epochCtx = c.epochCtx
		if ready == reconnectReady {
			select {
			case <-ready:
			default:
				close(ready)
			}
		}
		c.mu.Unlock()
		select {
		case <-ready:
		case <-epochCtx.Done():
			continue
		}
		c.mu.Lock()
		current := c.connected && c.connectionEpoch == epoch && c.txReady == ready
		if current {
			c.rejoining = false
		}
		c.mu.Unlock()
		if current {
			break
		}
	}

	enumerate := uint64(framer.PGNISOAddressClaim)
	c.retryAdvisoryProtocolContext(epochCtx, "post-reconnect bus enumeration", &pgn.IsoRequest{Pgn: &enumerate})
	if c.heartbeat != nil && c.heartbeat.currentInterval() > 0 {
		c.heartbeat.sendNow()
	}
}

// handleAddressChange moves application writes behind a fresh contention
// window whenever an active auto-claiming client changes address.
func (c *Client) handleAddressChange(newAddr uint8) {
	c.dispatchMu.Lock()
	defer c.dispatchMu.Unlock()
	c.mu.Lock()
	c.sourceAddr = newAddr
	if c.actisenseRemote != nil {
		c.actisenseRemote.invalidate(ErrActisenseRemoteEpochChanged)
	}
	if !c.claimed {
		c.mu.Unlock()
		return
	}
	if newAddr == 254 {
		c.mu.Unlock()
		c.fail(errors.New("n2k: address claiming failed: no address available"))
		return
	}
	c.invalidateEpochLocked(ErrEpochChanged)
	epoch := c.claimEpoch
	ready := make(chan struct{})
	c.txReady = ready
	c.mu.Unlock()
	c.resetReadEpoch()

	time.AfterFunc(250*time.Millisecond, func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		if !c.closed && c.claimEpoch == epoch && c.txReady == ready {
			close(ready)
		}
	})
}

// handleCommandedAddressTransfer validates the transfer-level requirements
// that are lost after generic PGN decoding, then delegates the address state
// transition to the claiming module. PGN 65240 is acted on only when it arrives
// as a broadcast ISO transport transfer with its exact nine-byte payload.
func (c *Client) handleCommandedAddressTransfer(destination uint8, data []byte) {
	if destination != framer.BroadcastAddr || len(data) != 9 {
		return
	}

	commandedName := binary.LittleEndian.Uint64(data[:8])
	changed, err := c.claimer.HandleCommandedAddress(commandedName, data[8])
	if err != nil {
		c.fail(fmt.Errorf("n2k: commanded address failed: %w", err))
		return
	}
	if changed {
		c.log.Info("address changed by PGN 65240", "address", data[8])
	}
}

// Write snapshots a PGN message and asynchronously transmits it. After Write
// returns the caller may reuse the message. Application jobs retain FIFO
// order; protocol traffic can run between their wire frames.
func (c *Client) Write(msg pgn.Message) *WriteResult {
	return c.WriteContext(context.Background(), msg)
}

// WriteContext cancels queued work and in-progress transmission when ctx ends.
// Cancellation after a physical write attempt may leave a partial transfer;
// WriteError reports that uncertainty. Writes are not replayed on reconnect.
func (c *Client) WriteContext(ctx context.Context, msg pgn.Message) *WriteResult {
	wr := newWriteResult()
	if c == nil {
		wr.complete(ErrClientClosed)
		return wr
	}
	if ctx == nil {
		ctx = context.Background()
	}
	c.mu.Lock()
	if err := c.readyErrorLocked(); err != nil {
		c.mu.Unlock()
		c.applicationWritesRejected.Add(1)
		wr.complete(err)
		return wr
	}
	jobCtx, stop := c.writeContextLocked(ctx, false)
	c.mu.Unlock()
	snapshot, err := snapshotMessage(msg)
	if err == nil {
		err = c.checkWriteContext(jobCtx)
	}
	if err != nil {
		stop()
		c.applicationWritesRejected.Add(1)
		wr.complete(err)
		return wr
	}
	c.mu.Lock()
	if err := c.readyErrorLocked(); err != nil {
		c.mu.Unlock()
		stop()
		c.applicationWritesRejected.Add(1)
		wr.complete(err)
		return wr
	}
	select {
	case c.writeCh <- writeJob{msg: snapshot, result: wr, ctx: jobCtx, stop: stop}:
		c.applicationWritesAccepted.Add(1)
	default:
		stop()
		c.applicationWritesRejected.Add(1)
		wr.complete(ErrWriteQueueFull)
	}
	c.mu.Unlock()
	return wr
}

func (c *Client) writeLoop() {
	defer c.writeWg.Done()
	for {
		select {
		case <-c.ctx.Done():
			c.drainWriteQueues(c.operationError())
			return
		case job := <-c.writeCh:
			c.runWriteJob(job)
		}
	}
}

func (c *Client) protocolWriteLoop() {
	defer c.writeWg.Done()
	for {
		if c.ctx.Err() != nil {
			c.drainWriteQueues(c.operationError())
			return
		}
		if job, ok := c.protocolTx.takeReady(); ok {
			c.runWriteJob(job)
			continue
		}
		select {
		case <-c.ctx.Done():
			c.drainWriteQueues(c.operationError())
			return
		case job := <-c.protocolTx.required:
			c.runWriteJob(job)
		case job := <-c.protocolTx.advisory:
			c.runWriteJob(job)
		}
	}
}

// runWriteJob executes one queued write, recovering a panicking Bus.WriteFrame
// so the write completes with an error and the loop keeps serving. Letting the
// panic escape would crash the process and, short of that, wedge every future
// Write blocked on the send to writeCh.
func (c *Client) runWriteJob(job writeJob) {
	if job.stop != nil {
		defer job.stop()
	}
	if job.ctx == nil {
		job.ctx = c.ctx
	}
	defer func() {
		if r := recover(); r != nil {
			c.log.Error("recovered panic in write loop",
				"panic", r, "stack", string(debug.Stack()))
			err := fmt.Errorf("n2k: panic writing message: %v", r)
			c.finishWriteJob(job, err)
			if c.ctx.Err() == nil {
				c.fail(err)
			}
		}
	}()
	c.finishWriteJob(job, c.doWriteContext(job.ctx, job.msg))
}

func (c *Client) finishWriteJob(job writeJob, err error) {
	if err != nil && job.ctx != nil {
		if stamp, ok := job.ctx.Value(writeStampKey{}).(writeStamp); ok && stamp.progress != nil {
			err = &WriteError{Err: err, CompletedRecords: stamp.progress.completed.Load(), TransmissionUncertain: stamp.progress.attempted.Load() > 0}
		}
	}
	job.result.complete(err)
	if job.protocol && c.protocolTx != nil {
		c.protocolTx.finish(err)
	} else if err != nil {
		c.applicationWritesFailed.Add(1)
	} else {
		c.applicationWritesCompleted.Add(1)
	}
	if err == nil || c.ctx.Err() != nil || errors.Is(err, ErrEpochChanged) || errors.Is(err, context.Canceled) {
		return
	}
	var busErr *runtimeBusError
	if errors.As(err, &busErr) {
		c.fail(busErr)
		return
	}
	if job.protocol {
		c.fail(fmt.Errorf("n2k: protocol transmission %s failed: %w", job.operation, err))
	}
}

// doWrite performs the synchronous work of encoding and framing a PGN message.
func (c *Client) doWriteContext(ctx context.Context, msg pgn.Message) error {
	if err := c.checkWriteContext(ctx); err != nil {
		return err
	}
	if msg == nil {
		return errors.New("n2k: cannot write a nil message")
	}
	v := reflect.ValueOf(msg)
	if v.Kind() == reflect.Pointer && v.IsNil() {
		return fmt.Errorf("n2k: cannot write a nil %T", msg)
	}

	// Wait for address readiness before sending.
	c.mu.Lock()
	ready := c.txReady
	terminalErr := c.terminalErr
	c.mu.Unlock()
	if terminalErr != nil {
		return terminalErr
	}
	select {
	case <-ready:
	case <-ctx.Done():
		return context.Cause(ctx)
	case <-c.ctx.Done():
		return c.currentTerminalError(c.ctx.Err().Error())
	}
	c.mu.Lock()
	closed := c.closed
	terminalErr = c.terminalErr
	c.mu.Unlock()
	if terminalErr != nil {
		return terminalErr
	}
	if closed {
		return ErrClientClosed
	}
	if err := c.checkWriteContext(ctx); err != nil {
		return err
	}

	pgnNum := msg.PGNNumber()

	// Check that msg implements pgn.PGN and extract its MessageInfo.
	pgnMsg, ok := msg.(pgn.PGN)
	if !ok {
		return fmt.Errorf("n2k: %T does not implement pgn.PGN", msg)
	}
	if pgnNum > 0x3FFFF {
		return fmt.Errorf("n2k: PGN %d exceeds the 18-bit range", pgnNum)
	}
	if pgnNum&0x20000 != 0 {
		return fmt.Errorf("n2k: PGN %d sets the reserved CAN-ID bit", pgnNum)
	}
	pduFormat := uint8((pgnNum >> 8) & 0xFF)
	if pduFormat < 240 && pgnNum&0xFF != 0 {
		return fmt.Errorf("n2k: PDU1 PGN %d must have a zero group-extension byte", pgnNum)
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
	if priority > 7 {
		return fmt.Errorf("n2k: priority %d is outside 0-7", priority)
	}

	var destination uint8 = 255
	if info.TargetId != nil {
		destination = *info.TargetId
	}
	if info.TargetId != nil && destination != framer.BroadcastAddr && pduFormat >= 240 {
		return fmt.Errorf("n2k: PDU2 PGN %d cannot target address %d", pgnNum, destination)
	}

	// Read sourceAddr under lock.
	c.mu.Lock()
	srcAddr := c.sourceAddr
	c.mu.Unlock()

	// Message-oriented buses (see MessageWriter) take whole payloads and do
	// their own wire framing; anything larger than one message still goes
	// frame-by-frame below via ISO-TP.
	if mw, ok := c.bus.(MessageWriter); ok && len(payload) <= 223 {
		if err := c.wire.send(ctx, func(ctx context.Context) error {
			if writer, ok := c.bus.(ContextMessageWriter); ok {
				return writer.WriteMessageContext(ctx, pgnNum, priority, srcAddr, destination, payload)
			}
			done := make(chan struct{})
			stop := context.AfterFunc(ctx, func() { _ = c.bus.Close(); close(done) })
			err := mw.WriteMessage(pgnNum, priority, srcAddr, destination, payload)
			if !stop() {
				<-done
			}
			return err
		}); err != nil {
			return err
		}
		destinationCopy := destination
		now := time.Now()
		stamp, _ := ctx.Value(writeStampKey{}).(writeStamp)
		c.publishObservation(raw.Observation{
			ConnectionEpoch: stamp.connection,
			ClaimEpoch:      stamp.claim,
			Kind:            raw.KindMessage,
			Timestamp:       now,
			ReceivedAt:      now,
			AdapterID:       "client",
			NetworkID:       "bus",
			Direction:       raw.DirectionTransmitted,
			PGN:             pgnNum,
			Priority:        priority,
			Source:          srcAddr,
			Destination:     &destinationCopy,
			Payload:         payload,
		})
		return nil
	}

	canID := framer.BuildCANID(pgnNum, priority, srcAddr, destination)

	isFast := false
	if infos, ok := pgn.PgnInfoLookup[pgnNum]; ok && len(infos) > 0 {
		isFast = infos[0].Fast
	}

	if !isFast && len(payload) <= 8 {
		frame := framer.FrameSingle(canID, payload)
		return c.writeFrameContext(ctx, frame)
	}

	if isFast && len(payload) <= 223 {
		seqID := uint8(c.fastSeq.Add(1) % 8)
		frames := framer.FrameFastPacket(canID, payload, seqID)
		for _, f := range frames {
			if err := c.writeFrameContext(ctx, f); err != nil {
				return err
			}
		}
		return nil
	}

	if destination != framer.BroadcastAddr {
		return c.tp.SendRTSCTSContext(ctx, pgnNum, srcAddr, destination, payload)
	}
	return c.tp.SendBAMContext(ctx, pgnNum, srcAddr, payload)
}

// Receive returns an iterator of decoded NMEA 2000 messages. For bus clients
// it reads from the internal message channel; for replay clients it builds a
// fresh Scanner over the client's config so each call gets a full replay.
func (c *Client) Receive() iter.Seq2[pgn.Message, error] {
	if c.msgHub != nil {
		return func(yield func(pgn.Message, error) bool) {
			sub := c.msgHub.subscribe()
			defer sub.unsubscribe()
			for msg := range sub.ch {
				if !yield(msg, nil) {
					return
				}
			}
			if err := sub.terminalError(); err != nil {
				yield(nil, err)
			}
		}
	}
	return func(yield func(pgn.Message, error) bool) {
		s := c.newReplayScanner()
		defer func() { _ = s.Close() }()
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
	if c.msgHub != nil {
		return &Scanner{ctx: c.ctx, cfg: c.cfg, sub: c.msgHub.subscribe()}
	}
	return c.newReplayScanner()
}

// newReplayScanner builds a Scanner over the client's already-parsed config,
// so replay clients share the exact pipeline code without re-parsing options.
func (c *Client) newReplayScanner() *Scanner {
	ctx, cancel := context.WithCancel(c.ctx)
	receiveBuffer := defaultReceiveBuffer
	if c.cfg.receiveBuffer != nil {
		receiveBuffer = *c.cfg.receiveBuffer
	}
	return &Scanner{ctx: ctx, cancel: cancel, cfg: c.cfg, ch: make(chan pgn.Message, receiveBuffer)}
}

// WrittenFrames returns an owned snapshot of the newest ReplayFrameCapacity
// frames written in replay mode, ordered oldest to newest. Once the capture
// fills, new writes evict its oldest frame and still succeed. Status exposes
// ReplayFramesDropped so capture loss is observable. Live clients return no
// captured frames; use Observe for live traffic.
func (c *Client) WrittenFrames() []can.Frame {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]can.Frame, len(c.writtenFrames))
	for i := range out {
		out[i] = c.writtenFrames[(c.writtenFramesStart+i)%len(c.writtenFrames)]
	}
	return out
}

// Close shuts down the client, cancels the context, and releases resources.
// It is safe to call Close multiple times.
func (c *Client) Close() error {
	c.mu.Lock()
	if c.closed {
		done := c.closeDone
		c.mu.Unlock()
		if done != nil {
			<-done
		}
		c.mu.Lock()
		err := c.closeErr
		c.mu.Unlock()
		return err
	}
	c.closed = true
	c.connected = false
	c.rejoining = false
	if c.epochCancel != nil {
		c.epochCancel(ErrClientClosed)
	}
	if c.correlator != nil {
		c.correlator.invalidate(ErrClientClosed)
	}
	c.mu.Unlock()
	if c.actisenseRemote != nil {
		c.actisenseRemote.close(ErrClientClosed)
	}

	// Cancel first and close the bus so an active hardware write is released.
	// Write never blocks on queue admission and the write channel is not closed,
	// so concurrent callers cannot race a send-on-closed-channel panic.
	c.cancel()
	var busErr error
	if c.bus != nil {
		busErr = c.bus.Close()
	}

	if c.heartbeat != nil {
		c.heartbeat.stop()
	}
	c.stopBroadcasters()
	if c.tp != nil {
		c.tp.Close()
	}
	c.writeWg.Wait()
	if c.wire != nil {
		<-c.wire.done
	}
	if c.busStarted {
		<-c.busDone
	}
	if c.system != nil {
		<-c.system.done
	}
	c.backgroundWg.Wait()
	if c.msgHub != nil {
		c.msgHub.close(nil)
	}
	if c.observationHub != nil {
		c.observationHub.close(nil)
	}
	c.mu.Lock()
	c.closeErr = busErr
	if c.closeDone != nil {
		close(c.closeDone)
	}
	c.mu.Unlock()
	return busErr
}
