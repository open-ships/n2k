package n2k

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/internal/gateway"
	"github.com/open-ships/n2k/pgn"
	"github.com/open-ships/n2k/raw"
)

const (
	defaultActisenseSessionTimeout    = 5 * time.Second
	defaultActisenseSessionInactivity = 500 * time.Millisecond
)

// ActisenseByteStream is a full-duplex byte transport for one Actisense
// connection epoch. Close must unblock Read and Write.
type ActisenseByteStream interface {
	io.Reader
	io.Writer
	io.Closer
}

// ActisenseOpenFunc opens one custom Actisense byte-stream connection epoch.
// It must return promptly when its context ends. It may be called again when
// reconnection is configured.
type ActisenseOpenFunc func(context.Context) (ActisenseByteStream, error)

type actisenseSessionConfig struct {
	logger         *slog.Logger
	reconnect      *ReconnectPolicy
	mode           ActisenseOperatingMode
	preserveMode   bool
	commandTimeout time.Duration
	readyTimeout   time.Duration
	inactivity     time.Duration
	buffer         int
	wireTrace      *ActisenseEBLTrace
}

func defaultActisenseSessionConfig() actisenseSessionConfig {
	return actisenseSessionConfig{
		logger: slog.Default(), mode: ActisenseModeTransferReceiveAll,
		commandTimeout: defaultActisenseSessionTimeout,
		readyTimeout:   defaultActisenseSessionTimeout,
		inactivity:     defaultActisenseSessionInactivity,
		buffer:         defaultReceiveBuffer,
	}
}

// ActisenseSessionOption configures a public gateway-owned session.
type ActisenseSessionOption interface{ applyActisenseSession(*actisenseSessionConfig) }

type actisenseSessionOptionFunc func(*actisenseSessionConfig)

func (f actisenseSessionOptionFunc) applyActisenseSession(config *actisenseSessionConfig) { f(config) }

// WithActisenseSessionLogger selects the session logger; nil uses slog.Default.
func WithActisenseSessionLogger(logger *slog.Logger) ActisenseSessionOption {
	return actisenseSessionOptionFunc(func(config *actisenseSessionConfig) { config.logger = logger })
}

// WithActisenseSessionReconnect enables reconnects after a dropped session.
// Each acknowledged connection starts a new epoch and cancels old requests.
// Zero backoff fields use the ReconnectPolicy defaults.
func WithActisenseSessionReconnect(policy ReconnectPolicy) ActisenseSessionOption {
	return actisenseSessionOptionFunc(func(config *actisenseSessionConfig) { config.reconnect = &policy })
}

// WithActisenseSessionMode selects gateway-owned BST message mode 1 or 2.
// Source-authoritative mode 5 belongs to NewClient via ActisenseTCP or
// ActisenseSerial and is intentionally rejected here.
func WithActisenseSessionMode(mode ActisenseOperatingMode) ActisenseSessionOption {
	return actisenseSessionOptionFunc(func(config *actisenseSessionConfig) {
		config.mode, config.preserveMode = mode, false
	})
}

// WithActisensePreserveOperatingMode opens a control session that reads and
// preserves the device's mode on every connection. It sends no mode setter
// during startup or Close. The last mode/preserve option wins. PGN sends still
// require mode 1 or 2; raw BST and BEM commands remain available in other modes.
func WithActisensePreserveOperatingMode() ActisenseSessionOption {
	return actisenseSessionOptionFunc(func(config *actisenseSessionConfig) { config.preserveMode = true })
}

// WithActisenseCommandTimeout bounds each local command, physical send, and
// configuration restoration attempt. It must be positive; the default is five
// seconds. An earlier caller deadline also applies to the original operation.
func WithActisenseCommandTimeout(timeout time.Duration) ActisenseSessionOption {
	return actisenseSessionOptionFunc(func(config *actisenseSessionConfig) { config.commandTimeout = timeout })
}

// WithActisenseMultiReplyInactivity bounds the gap between local BEM replies.
// It must be positive; the default is 500 milliseconds. The command timeout
// separately bounds the whole response train.
func WithActisenseMultiReplyInactivity(timeout time.Duration) ActisenseSessionOption {
	return actisenseSessionOptionFunc(func(config *actisenseSessionConfig) { config.inactivity = timeout })
}

// WithActisenseSessionReadyTimeout bounds opening and mode negotiation.
// It must be positive; the default is five seconds.
func WithActisenseSessionReadyTimeout(timeout time.Duration) ActisenseSessionOption {
	return actisenseSessionOptionFunc(func(config *actisenseSessionConfig) { config.readyTimeout = timeout })
}

// WithActisenseSessionBuffer sets the positive per-subscriber observation and
// diagnostic capacity; the default is 64. Overflow terminates the affected
// subscription with ErrObservationOverflow without blocking the session reader.
func WithActisenseSessionBuffer(size int) ActisenseSessionOption {
	return actisenseSessionOptionFunc(func(config *actisenseSessionConfig) { config.buffer = size })
}

// WithActisenseWireTrace records exact transmitted bytes and received wire
// evidence to an EBL trace. Valid received BDTP frames are stored in the
// checksum-stripped BSTRaw representation used by Actisense tooling.
func WithActisenseWireTrace(trace *ActisenseEBLTrace) ActisenseSessionOption {
	return actisenseSessionOptionFunc(func(config *actisenseSessionConfig) { config.wireTrace = trace })
}

func applyActisenseSessionOptions(options []ActisenseSessionOption) (actisenseSessionConfig, error) {
	config := defaultActisenseSessionConfig()
	for _, option := range options {
		if option == nil {
			return config, errors.New("n2k: nil ActisenseSessionOption")
		}
		option.applyActisenseSession(&config)
	}
	if config.logger == nil {
		config.logger = slog.Default()
	}
	if !config.preserveMode && config.mode != ActisenseModeTransferNormal && config.mode != ActisenseModeTransferReceiveAll {
		return config, fmt.Errorf("n2k: gateway-owned Actisense session mode must be 1 or 2; got %d", config.mode)
	}
	if config.commandTimeout <= 0 || config.readyTimeout <= 0 || config.inactivity <= 0 {
		return config, errors.New("n2k: Actisense session timeouts must be positive")
	}
	if config.buffer <= 0 {
		return config, errors.New("n2k: Actisense session buffer must be positive")
	}
	return config, nil
}

// ActisenseSessionStatus honestly describes the gateway-owned connection. A
// session never claims an independent NMEA 2000 source identity.
type ActisenseSessionStatus struct {
	Connected              bool
	ConnectionEpoch        uint64
	Closed                 bool
	TerminalError          error
	OperatingMode          ActisenseOperatingMode
	SourceAuthoritative    bool
	ReceiveAll             bool
	ISOControlPGNsVisible  bool
	ObservationSubscribers int
	DiagnosticSubscribers  int
	DeviceCapabilities     ActisenseDeviceCapabilities
	WireTraceError         error
	Metrics                ActisenseSessionMetrics
	// GatewaySourceAddress is learned from a nonce-verified remote reply,
	// never from the stored CAN configuration. It is nil until verified.
	GatewaySourceAddress *uint8
	IdentityEpoch        uint64
	RemoteMetrics        ActisenseProtocolMetrics
}

// ActisenseTxPGNConfiguration is one staged transmit-list change. PGN is the
// message number and Flag selects enabled, disabled, or respond-only behavior.
// Rate is milliseconds: zero means event-driven; nil or a value >= 65535 leaves
// the current rate unchanged. The pointed-to rate must remain stable during
// ConfigureTransmitPGNs.
type ActisenseTxPGNConfiguration struct {
	PGN  uint32
	Flag ActisensePGNEnableFlag
	Rate *uint32
}

// ActisenseGatewaySession owns a gateway identity, a sole-reader BEM session,
// observations, typed diagnostics, and explicit PGN-list transactions. It is
// deliberately not a Bus and does not run address claiming.
type ActisenseGatewaySession struct {
	*ActisenseDevice
	transport *gateway.ActisenseGatewaySession
	mode      ActisenseOperatingMode

	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}

	observations   *observationHub
	diagnostics    *actisenseDiagnosticHub
	wireTrace      *ActisenseEBLTrace
	commandTimeout time.Duration
	remote         *actisenseRemoteManager
	probeGate      chan struct{}
	txGate         chan struct{}

	mu                 sync.Mutex
	connected          bool
	epoch              uint64
	closed             bool
	terminalErr        error
	closeErr           error
	txOriginal         map[uint32]ActisenseTxPGNState
	txEpoch            uint64
	txCancel           context.CancelFunc
	closeOnce          sync.Once
	remoteAddress      uint8
	remoteAddressKnown bool
	identityEpoch      uint64
	remoteProbe        *actisenseGatewayProbe
}

func gatewayReconnectPolicy(policy *ReconnectPolicy) *gateway.ReconnectPolicy {
	if policy == nil {
		return nil
	}
	return &gateway.ReconnectPolicy{InitialBackoff: policy.InitialBackoff, MaxBackoff: policy.MaxBackoff}
}

// NewActisenseTCPSession opens an honest gateway-owned BST-93/94 and BEM
// session over TCP.
func NewActisenseTCPSession(ctx context.Context, address string, options ...ActisenseSessionOption) (*ActisenseGatewaySession, error) {
	if address == "" {
		return nil, errors.New("n2k: Actisense TCP session address cannot be empty")
	}
	config, err := applyActisenseSessionOptions(options)
	if err != nil {
		return nil, err
	}
	transport := gateway.NewActisenseTCPGatewaySession(config.logger, address, gatewayReconnectPolicy(config.reconnect), config.mode)
	return startActisenseGatewaySession(ctx, transport, config)
}

// NewActisenseSerialSession opens a gateway-owned session with configurable
// serial settings. The zero serial config defaults to 115200 8N1.
func NewActisenseSerialSession(ctx context.Context, port string, serialConfig ActisenseSerialConfig, options ...ActisenseSessionOption) (*ActisenseGatewaySession, error) {
	if port == "" {
		return nil, errors.New("n2k: Actisense serial session port cannot be empty")
	}
	serialConfig = normalizeActisenseSerialConfig(serialConfig)
	if err := serialConfig.validate(); err != nil {
		return nil, err
	}
	config, err := applyActisenseSessionOptions(options)
	if err != nil {
		return nil, err
	}
	transport := gateway.NewActisenseSerialGatewaySession(config.logger, port, gatewayActisenseSerialSettings(serialConfig), config.mode)
	return startActisenseGatewaySession(ctx, transport, config)
}

// NewActisenseGatewaySession opens a gateway-owned session over a custom
// full-duplex byte-stream Adapter.
func NewActisenseGatewaySession(ctx context.Context, endpoint string, open ActisenseOpenFunc, options ...ActisenseSessionOption) (*ActisenseGatewaySession, error) {
	if open == nil {
		return nil, errors.New("n2k: Actisense custom session requires an open function")
	}
	if endpoint == "" {
		endpoint = "custom"
	}
	config, err := applyActisenseSessionOptions(options)
	if err != nil {
		return nil, err
	}
	transport := gateway.NewActisenseCustomGatewaySession(config.logger, endpoint, "actisense:"+endpoint, func(openCtx context.Context) (gateway.ActisenseConnection, error) {
		return open(openCtx)
	}, gatewayReconnectPolicy(config.reconnect), config.mode)
	return startActisenseGatewaySession(ctx, transport, config)
}

func startActisenseGatewaySession(parent context.Context, transport *gateway.ActisenseGatewaySession, config actisenseSessionConfig) (*ActisenseGatewaySession, error) {
	if parent == nil {
		parent = context.Background()
	}
	if err := transport.SetCommandTimeout(config.commandTimeout); err != nil {
		return nil, err
	}
	if config.preserveMode {
		transport.PreserveOperatingMode()
	}
	transport.SetReconnectPolicy(gatewayReconnectPolicy(config.reconnect))
	runCtx, cancel := context.WithCancel(parent)
	session := &ActisenseGatewaySession{
		transport: transport, mode: config.mode, ctx: runCtx, cancel: cancel, done: make(chan struct{}),
		observations: newObservationHub(config.buffer), diagnostics: newActisenseDiagnosticHub(config.buffer),
		txOriginal: make(map[uint32]ActisenseTxPGNState), wireTrace: config.wireTrace,
		commandTimeout: config.commandTimeout,
		probeGate:      make(chan struct{}, 1),
		txGate:         make(chan struct{}, 1),
	}
	session.remote = newActisenseRemoteManager(nil)
	session.remote.gateway = session
	session.ActisenseDevice = &ActisenseDevice{CommandSet: actisense.NewCommandSet(transport, actisense.CommandSetConfig{
		Timeout: config.commandTimeout, MultiInactivity: config.inactivity,
	})}
	transport.SetConnectionObserver(session.handleConnection)
	transport.SetModeObserver(session.handleMode)
	transport.SetMessageObserver(session.handleGatewayMessage)
	transport.SetDiagnosticObserver(session.diagnostics.publish)
	if config.wireTrace != nil {
		transport.SetWireObserver(config.wireTrace.trace)
	}
	go session.run()

	timer := time.NewTimer(config.readyTimeout)
	defer timer.Stop()
	select {
	case <-transport.Ready():
		return session, nil
	case <-session.done:
		err := session.Err()
		if err == nil {
			err = errors.New("n2k: Actisense session stopped before becoming ready")
		}
		return nil, wrapActisenseModeError(err)
	case <-timer.C:
		_ = session.Close()
		return nil, errors.New("n2k: Actisense session did not become ready within ready timeout")
	case <-parent.Done():
		_ = session.Close()
		return nil, parent.Err()
	}
}

func (s *ActisenseGatewaySession) run() {
	err := s.transport.RunObservations(s.ctx, func(observation raw.Observation) {
		s.observations.publish(normalizeObservation(observation))
	})
	if errors.Is(err, context.Canceled) || s.ctx.Err() != nil {
		err = nil
	}
	s.mu.Lock()
	if err != nil && s.terminalErr == nil {
		s.terminalErr = wrapActisenseModeError(err)
	}
	terminalErr := s.terminalErr
	s.connected = false
	s.invalidateRemoteLocked(terminalErr)
	s.mu.Unlock()
	s.remote.close(terminalErr)
	s.observations.close(terminalErr)
	s.diagnostics.close(terminalErr)
	close(s.done)
}

func (s *ActisenseGatewaySession) handleConnection(connected bool, epoch uint64) {
	s.mu.Lock()
	if !connected && s.txEpoch == epoch {
		// Volatile list changes belonged to the lost epoch and cannot be
		// restored through a later connection.
		clear(s.txOriginal)
		s.txEpoch = 0
	}
	s.connected, s.epoch = connected, epoch
	s.invalidateRemoteLocked(ErrActisenseRemoteEpochChanged)
	s.mu.Unlock()
}

func (s *ActisenseGatewaySession) handleMode(mode ActisenseOperatingMode) {
	s.mu.Lock()
	if s.mode != mode {
		s.invalidateRemoteLocked(ErrActisenseRemoteEpochChanged)
	}
	s.mode = mode
	s.mu.Unlock()
}

// Status returns a concurrency-safe, owned snapshot of connection readiness,
// identity, subscribers, and cumulative metrics. A nil session is closed.
func (s *ActisenseGatewaySession) Status() ActisenseSessionStatus {
	if s == nil {
		return ActisenseSessionStatus{Closed: true, TerminalError: errors.New("n2k: nil Actisense session")}
	}
	s.mu.Lock()
	status := ActisenseSessionStatus{
		Connected: s.connected, ConnectionEpoch: s.epoch, Closed: s.closed, TerminalError: s.terminalErr,
		OperatingMode: s.mode, ReceiveAll: s.mode == ActisenseModeTransferReceiveAll,
		DeviceCapabilities: s.DeviceCapabilities(),
		IdentityEpoch:      s.identityEpoch,
	}
	if s.remoteAddressKnown {
		address := s.remoteAddress
		status.GatewaySourceAddress = &address
	}
	s.mu.Unlock()
	status.RemoteMetrics = s.remote.metrics.snapshot()
	status.ISOControlPGNsVisible = !status.ReceiveAll || !status.DeviceCapabilities.ReceiveAllOmitsISOControlPGNs
	status.ObservationSubscribers = s.observations.subscriberCount()
	status.DiagnosticSubscribers = s.diagnostics.subscriberCount()
	metrics := s.transport.Metrics()
	status.Metrics = ActisenseSessionMetrics{
		ConnectionEpochs: metrics.ConnectionEpochs, Reconnects: metrics.Reconnects,
		GatewayResets: metrics.GatewayResets, Protocol: metrics.Protocol,
	}
	if s.wireTrace != nil {
		status.WireTraceError = s.wireTrace.Err()
	}
	return status
}

// Err returns the terminal session error, or nil while healthy or after normal
// shutdown. It is safe concurrently with other session methods.
func (s *ActisenseGatewaySession) Err() error {
	if s == nil {
		return errors.New("n2k: nil Actisense session")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.terminalErr
}

// Observations subscribes when iteration starts and yields owned transport
// records until shutdown, error, or early loop exit. Independent bounded
// subscriptions fail with ErrObservationOverflow if a consumer falls behind.
func (s *ActisenseGatewaySession) Observations() iter.Seq2[Observation, error] {
	return func(yield func(Observation, error) bool) {
		if s == nil || s.observations == nil {
			yield(Observation{}, errors.New("n2k: nil Actisense session"))
			return
		}
		subscription := s.observations.subscribe()
		defer subscription.unsubscribe()
		for observation := range subscription.ch {
			if !yield(observation.Clone(), nil) {
				return
			}
		}
		if err := subscription.terminalError(); err != nil {
			yield(Observation{}, err)
		}
	}
}

// Diagnostics yields owned local BEM diagnostic events, including retained
// startup diagnostics. Subscriptions are independent and bounded by the session
// buffer; slow consumers terminate with ErrObservationOverflow.
func (s *ActisenseGatewaySession) Diagnostics() iter.Seq2[ActisenseDiagnostic, error] {
	return func(yield func(ActisenseDiagnostic, error) bool) {
		if s == nil || s.diagnostics == nil {
			yield(ActisenseDiagnostic{}, errors.New("n2k: nil Actisense session"))
			return
		}
		subscription := s.diagnostics.subscribe()
		defer subscription.unsubscribe()
		for diagnostic := range subscription.ch {
			if !yield(cloneActisenseDiagnostic(diagnostic), nil) {
				return
			}
		}
		if err := subscription.terminalError(); err != nil {
			yield(ActisenseDiagnostic{}, err)
		}
	}
}

// SendPGN transmits one assembled PGN under the gateway's own claimed source
// address. The message must implement pgn.PGN; nil or read-only messages return
// an error. It never enables or activates a Tx list implicitly.
func (s *ActisenseGatewaySession) SendPGN(ctx context.Context, message pgn.Message) error {
	if message == nil {
		return errors.New("n2k: cannot send a nil PGN through an Actisense session")
	}
	typed, ok := message.(pgn.PGN)
	if !ok {
		return fmt.Errorf("n2k: %T does not implement pgn.PGN", message)
	}
	payload, err := pgn.EncodeMessage(message)
	if err != nil {
		return err
	}
	info := typed.MessageInfo()
	priority := uint8(6)
	if info.Priority != nil {
		priority = *info.Priority
	}
	destination := uint8(255)
	if info.TargetId != nil {
		destination = *info.TargetId
	}
	return s.SendRawPGN(ctx, message.PGNNumber(), priority, destination, payload)
}

// SendRawPGN synchronously sends up to 223 payload bytes using the gateway
// identity. Priority must be 0-7; broadcast PGNs require destination 255. It
// copies the payload, honors caller and command deadlines, and rejects an
// unready session. Success confirms the transport write, not remote acceptance.
// PGN lists are never changed implicitly.
func (s *ActisenseGatewaySession) SendRawPGN(ctx context.Context, pgnNumber uint32, priority, destination uint8, payload []byte) error {
	if s == nil {
		return actisense.ErrSessionClosed
	}
	s.mu.Lock()
	closed := s.closed
	s.mu.Unlock()
	if closed {
		return actisense.ErrSessionClosed
	}
	if err := actisenseValidateMessage(pgnNumber, priority, destination, payload); err != nil {
		return err
	}
	writeCtx, cancel := s.writeContext(ctx)
	defer cancel()
	return s.transport.WriteMessageContext(writeCtx, pgnNumber, priority, destination, append([]byte(nil), payload...))
}

func (s *ActisenseGatewaySession) writeContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, s.commandTimeout)
}

func actisenseValidateMessage(pgnNumber uint32, priority, destination uint8, payload []byte) error {
	if pgnNumber > 0x3FFFF || pgnNumber&0x20000 != 0 {
		return fmt.Errorf("n2k: PGN %d is outside the NMEA 2000 range", pgnNumber)
	}
	if priority > 7 {
		return fmt.Errorf("n2k: priority %d is outside 0-7", priority)
	}
	if pgnNumber>>8&0xFF < 240 && pgnNumber&0xFF != 0 {
		return fmt.Errorf("n2k: PDU1 PGN %d must have a zero group-extension byte", pgnNumber)
	}
	if pgnNumber>>8&0xFF >= 240 && destination != 255 {
		return fmt.Errorf("n2k: PDU2 PGN %d cannot target address %d", pgnNumber, destination)
	}
	if len(payload) > 223 {
		return fmt.Errorf("n2k: Actisense gateway PGN payload is %d bytes; maximum is 223", len(payload))
	}
	return nil
}

// ConfigureTransmitPGNs snapshots every affected session entry, stages all
// changes, then activates once. Transactions are serialized. On failure, a
// separate cleanup deadline attempts restoration on the same connection epoch;
// the returned error includes any restoration failure. The earliest original
// entries are retained for another restoration attempt by Close. Cancellation
// may therefore take up to the command timeout to finish cleanup.
// Direct Tx-list commands must not run concurrently with this transaction.
func (s *ActisenseGatewaySession) ConfigureTransmitPGNs(ctx context.Context, configurations []ActisenseTxPGNConfiguration) error {
	if s == nil {
		return actisense.ErrSessionClosed
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if len(configurations) == 0 {
		return nil
	}
	select {
	case s.txGate <- struct{}{}:
		defer func() { <-s.txGate }()
	case <-ctx.Done():
		return ctx.Err()
	case <-s.ctx.Done():
		return actisense.ErrSessionClosed
	}
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return actisense.ErrSessionClosed
	}
	epoch := s.epoch
	ctx, cancel := context.WithCancel(ctx)
	s.txCancel = cancel
	s.mu.Unlock()
	stop := context.AfterFunc(s.ctx, cancel)
	defer func() {
		stop()
		cancel()
		s.mu.Lock()
		s.txCancel = nil
		s.mu.Unlock()
	}()
	requester, err := s.transport.EpochRequester(epoch)
	if err != nil {
		return err
	}
	commands := actisense.NewCommandSet(requester, actisense.CommandSetConfig{Timeout: s.commandTimeout})

	originals := make(map[uint32]ActisenseTxPGNState, len(configurations))
	for _, configuration := range configurations {
		if _, duplicate := originals[configuration.PGN]; duplicate {
			return fmt.Errorf("n2k: duplicate Tx PGN %d in one Actisense transaction", configuration.PGN)
		}
		state, err := commands.GetTxPGN(ctx, configuration.PGN)
		if err != nil {
			return fmt.Errorf("n2k: snapshot Actisense Tx PGN %d: %w", configuration.PGN, err)
		}
		originals[configuration.PGN] = state
	}
	// Retain the original settings before any physical mutation, including
	// failed transactions whose first rollback might itself fail.
	s.mu.Lock()
	if s.closed || !s.connected || s.epoch != epoch {
		s.mu.Unlock()
		return actisense.ErrSessionClosed
	}
	if s.txEpoch != epoch {
		clear(s.txOriginal)
		s.txEpoch = epoch
	}
	for number, state := range originals {
		if _, exists := s.txOriginal[number]; !exists {
			s.txOriginal[number] = state
		}
	}
	s.mu.Unlock()
	for _, configuration := range configurations {
		if _, err := commands.SetTxPGN(ctx, configuration.PGN, configuration.Flag, configuration.Rate); err != nil {
			return errors.Join(fmt.Errorf("n2k: stage Actisense Tx PGN %d: %w", configuration.PGN, err), s.rollbackTransmitPGNs(commands, originals))
		}
	}
	if err := commands.ActivatePGNLists(ctx); err != nil {
		return errors.Join(fmt.Errorf("n2k: activate Actisense Tx PGN transaction: %w", err), s.rollbackTransmitPGNs(commands, originals))
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.connected || s.epoch != epoch {
		return errors.New("n2k: Actisense connection epoch changed during Tx PGN transaction")
	}
	return nil
}

func (s *ActisenseGatewaySession) rollbackTransmitPGNs(commands *actisense.CommandSet, originals map[uint32]ActisenseTxPGNState) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.commandTimeout)
	defer cancel()
	numbers := make([]uint32, 0, len(originals))
	for number := range originals {
		numbers = append(numbers, number)
	}
	slices.Sort(numbers)
	var failures []error
	for _, number := range numbers {
		state := originals[number]
		rate := state.Rate
		if _, err := commands.SetTxPGN(ctx, number, ActisensePGNEnableFlag(state.Enabled), &rate); err != nil {
			failures = append(failures, fmt.Errorf("restore Tx PGN %d: %w", number, err))
		}
		if ctx.Err() != nil {
			break
		}
	}
	if err := commands.ActivatePGNLists(ctx); err != nil {
		failures = append(failures, fmt.Errorf("activate restored Tx PGNs: %w", err))
	}
	if err := errors.Join(failures...); err != nil {
		return fmt.Errorf("n2k: Actisense Tx PGN restoration is uncertain: %w", err)
	}
	return nil
}

func (s *ActisenseGatewaySession) restoreTransmitPGNs() error {
	s.mu.Lock()
	if !s.connected || len(s.txOriginal) == 0 || s.txEpoch != s.epoch {
		s.mu.Unlock()
		return nil
	}
	originals := make(map[uint32]ActisenseTxPGNState, len(s.txOriginal))
	for number, state := range s.txOriginal {
		originals[number] = state
	}
	epoch := s.txEpoch
	s.mu.Unlock()
	requester, err := s.transport.EpochRequester(epoch)
	if err != nil {
		return err
	}
	commands := actisense.NewCommandSet(requester, actisense.CommandSetConfig{Timeout: s.commandTimeout})
	return s.rollbackTransmitPGNs(commands, originals)
}

// Close cancels an active configuration transaction, waits for its cleanup,
// attempts same-epoch restoration, closes the transport, and joins the reader.
// Restoration has its own command timeout; failed restoration is included in
// the returned error. Close is concurrent-safe, idempotent, and nil-safe.
func (s *ActisenseGatewaySession) Close() error {
	if s == nil {
		return nil
	}
	s.closeOnce.Do(func() {
		s.mu.Lock()
		s.closed = true
		s.invalidateRemoteLocked(actisense.ErrSessionClosed)
		if s.txCancel != nil {
			s.txCancel()
		}
		s.mu.Unlock()
		s.txGate <- struct{}{}
		restoreErr := s.restoreTransmitPGNs()
		closeErr := errors.Join(restoreErr, s.transport.Close())
		<-s.txGate
		s.cancel()
		<-s.done
		if s.wireTrace != nil {
			if traceErr := s.wireTrace.Flush(); closeErr == nil {
				closeErr = traceErr
			}
		}
		s.mu.Lock()
		s.closeErr = closeErr
		s.mu.Unlock()
	})
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closeErr
}

type actisenseDiagnosticHub struct {
	mu          sync.Mutex
	buffer      int
	subscribers map[*actisenseDiagnosticSubscription]struct{}
	backlog     []ActisenseDiagnostic
	closed      bool
	err         error
}

type actisenseDiagnosticSubscription struct {
	hub  *actisenseDiagnosticHub
	ch   chan ActisenseDiagnostic
	once sync.Once
	mu   sync.Mutex
	err  error
}

func newActisenseDiagnosticHub(buffer int) *actisenseDiagnosticHub {
	return &actisenseDiagnosticHub{buffer: buffer, subscribers: make(map[*actisenseDiagnosticSubscription]struct{})}
}

func cloneActisenseDiagnostic(value ActisenseDiagnostic) ActisenseDiagnostic {
	value.Response.Data = append([]byte(nil), value.Response.Data...)
	if value.ErrorReport != nil {
		copyReport := *value.ErrorReport
		copyReport.ContextData = append([]byte(nil), value.ErrorReport.ContextData...)
		value.ErrorReport = &copyReport
	}
	if value.System != nil {
		copyStatus := *value.System
		copyStatus.Individual = append([]actisense.IndividualBufferStatus(nil), value.System.Individual...)
		copyStatus.Unified = append([]actisense.UnifiedBufferStatus(nil), value.System.Unified...)
		value.System = &copyStatus
	}
	return value
}

func (h *actisenseDiagnosticHub) publish(value actisense.Diagnostic) {
	value = cloneActisenseDiagnostic(value)
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed {
		return
	}
	if len(h.subscribers) == 0 {
		if len(h.backlog) == h.buffer {
			copy(h.backlog, h.backlog[1:])
			h.backlog[len(h.backlog)-1] = value
		} else {
			h.backlog = append(h.backlog, value)
		}
		return
	}
	for subscription := range h.subscribers {
		select {
		case subscription.ch <- cloneActisenseDiagnostic(value):
		default:
			subscription.failLocked(ErrObservationOverflow)
			delete(h.subscribers, subscription)
		}
	}
}

func (h *actisenseDiagnosticHub) subscribe() *actisenseDiagnosticSubscription {
	h.mu.Lock()
	defer h.mu.Unlock()
	subscription := &actisenseDiagnosticSubscription{hub: h, ch: make(chan ActisenseDiagnostic, h.buffer)}
	for _, value := range h.backlog {
		subscription.ch <- cloneActisenseDiagnostic(value)
	}
	h.backlog = nil
	if h.closed {
		subscription.failLocked(h.err)
		return subscription
	}
	h.subscribers[subscription] = struct{}{}
	return subscription
}

func (h *actisenseDiagnosticHub) close(err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed {
		return
	}
	h.closed, h.err = true, err
	for subscription := range h.subscribers {
		subscription.failLocked(err)
		delete(h.subscribers, subscription)
	}
}

func (h *actisenseDiagnosticHub) subscriberCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.subscribers)
}

func (s *actisenseDiagnosticSubscription) failLocked(err error) {
	s.mu.Lock()
	s.err = err
	s.mu.Unlock()
	s.once.Do(func() { close(s.ch) })
}

func (s *actisenseDiagnosticSubscription) unsubscribe() {
	s.hub.mu.Lock()
	delete(s.hub.subscribers, s)
	s.once.Do(func() { close(s.ch) })
	s.hub.mu.Unlock()
}

func (s *actisenseDiagnosticSubscription) terminalError() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.err
}
