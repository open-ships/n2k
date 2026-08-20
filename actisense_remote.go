package n2k

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"sync"
	"sync/atomic"
	"time"

	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/pgn"
	"github.com/open-ships/n2k/raw"
)

const (
	actisenseRemotePGN       = uint32(126720)
	actisenseRemotePriority  = uint8(3)
	maxRemotePendingRequests = 64
)

var ErrActisenseRemoteEpochChanged = errors.New("n2k: Actisense remote request canceled because the local address or connection epoch changed")

type actisenseRemoteConfig struct {
	timeout    time.Duration
	inactivity time.Duration
}

// ActisenseRemoteOption configures a remote-device handle.
type ActisenseRemoteOption interface{ applyActisenseRemote(*actisenseRemoteConfig) }

type actisenseRemoteOptionFunc func(*actisenseRemoteConfig)

func (f actisenseRemoteOptionFunc) applyActisenseRemote(config *actisenseRemoteConfig) { f(config) }

func WithActisenseRemoteTimeout(timeout time.Duration) ActisenseRemoteOption {
	return actisenseRemoteOptionFunc(func(config *actisenseRemoteConfig) { config.timeout = timeout })
}

func WithActisenseRemoteMultiReplyInactivity(timeout time.Duration) ActisenseRemoteOption {
	return actisenseRemoteOptionFunc(func(config *actisenseRemoteConfig) { config.inactivity = timeout })
}

// ActisenseRemoteDevice issues the same typed commands as a local gateway
// through Actisense's addressed PGN-126720 envelope.
type ActisenseRemoteDevice struct {
	*ActisenseDevice
	SourceAddress uint8
	manager       *actisenseRemoteManager
}

// ActisenseRemoteDevice returns a handle for one addressed Actisense device.
// Incoming replies must target this Client's source snapshot and are canceled
// across connection or address-claim epochs.
func (c *Client) ActisenseRemoteDevice(source uint8, options ...ActisenseRemoteOption) (*ActisenseRemoteDevice, error) {
	if c == nil {
		return nil, ErrClientClosed
	}
	if source > 251 {
		return nil, fmt.Errorf("n2k: Actisense remote source address %d is outside 0-251", source)
	}
	config := actisenseRemoteConfig{timeout: defaultActisenseSessionTimeout, inactivity: defaultActisenseSessionInactivity}
	for _, option := range options {
		if option == nil {
			return nil, errors.New("n2k: nil ActisenseRemoteOption")
		}
		option.applyActisenseRemote(&config)
	}
	if config.timeout <= 0 || config.inactivity <= 0 {
		return nil, errors.New("n2k: Actisense remote timeouts must be positive")
	}
	c.mu.Lock()
	manager := c.actisenseRemote
	closed := c.closed
	terminalErr := c.terminalErr
	c.mu.Unlock()
	if terminalErr != nil {
		return nil, terminalErr
	}
	if closed || manager == nil {
		return nil, ErrClientClosed
	}
	requester := &actisenseRemoteRequester{manager: manager, source: source}
	return &ActisenseRemoteDevice{
		ActisenseDevice: &ActisenseDevice{CommandSet: actisense.NewCommandSet(requester, actisense.CommandSetConfig{
			Timeout: config.timeout, MultiInactivity: config.inactivity, Remote: true,
		})},
		SourceAddress: source,
		manager:       manager,
	}, nil
}

func (d *ActisenseRemoteDevice) Diagnostics() iter.Seq2[ActisenseDiagnostic, error] {
	return func(yield func(ActisenseDiagnostic, error) bool) {
		if d == nil || d.manager == nil {
			yield(ActisenseDiagnostic{}, ErrClientClosed)
			return
		}
		subscription := d.manager.diagnostics.subscribe()
		defer subscription.unsubscribe()
		for diagnostic := range subscription.ch {
			if diagnostic.Response.Origin.Path != ActisenseBEMRemote || diagnostic.Response.Origin.Source != d.SourceAddress {
				continue
			}
			if !yield(cloneActisenseDiagnostic(diagnostic), nil) {
				return
			}
		}
		if err := subscription.terminalError(); err != nil {
			yield(ActisenseDiagnostic{}, err)
		}
	}
}

// Metrics returns cumulative BEM correlation and latency counters shared by
// all remote-device handles on this Client.
func (d *ActisenseRemoteDevice) Metrics() ActisenseProtocolMetrics {
	if d == nil || d.manager == nil {
		return ActisenseProtocolMetrics{BSTFrames: make(map[byte]uint64)}
	}
	return d.manager.metrics.snapshot()
}

type actisenseRemoteKey struct {
	source          uint8
	destination     uint8
	connectionEpoch uint64
	claimEpoch      uint64
	bstID           byte
	bemID           byte
}

type actisenseRemoteResult struct {
	response actisense.BEMResponse
	err      error
}

type actisenseRemotePending struct {
	results   chan actisenseRemoteResult
	multi     bool
	delivered int
}

type actisenseRemoteManager struct {
	client *Client

	mu          sync.Mutex
	pending     map[actisenseRemoteKey]*actisenseRemotePending
	closed      bool
	diagnostics *actisenseDiagnosticHub
	metrics     actisenseRemoteMetrics
}

type actisenseRemoteMetrics struct {
	requests          atomic.Uint64
	responses         atomic.Uint64
	completed         atomic.Uint64
	correlationMisses atomic.Uint64
	duplicates        atomic.Uint64
	timeouts          atomic.Uint64
	deviceErrors      atomic.Uint64
	negativeAcks      atomic.Uint64
	overflows         atomic.Uint64
	inFlight          atomic.Uint64
	maxInFlight       atomic.Uint64
	latencyTotal      atomic.Uint64
	latencyMinimum    atomic.Uint64
	latencyMaximum    atomic.Uint64
}

func (m *actisenseRemoteMetrics) incrementInFlight() {
	current := m.inFlight.Add(1)
	for maximum := m.maxInFlight.Load(); current > maximum; maximum = m.maxInFlight.Load() {
		if m.maxInFlight.CompareAndSwap(maximum, current) {
			break
		}
	}
}

func (m *actisenseRemoteMetrics) observeLatency(elapsed time.Duration) {
	nanos := uint64(max(elapsed.Nanoseconds(), 0))
	m.latencyTotal.Add(nanos)
	for current := m.latencyMinimum.Load(); current == 0 || nanos < current; current = m.latencyMinimum.Load() {
		if m.latencyMinimum.CompareAndSwap(current, nanos) {
			break
		}
	}
	for current := m.latencyMaximum.Load(); nanos > current; current = m.latencyMaximum.Load() {
		if m.latencyMaximum.CompareAndSwap(current, nanos) {
			break
		}
	}
}

func (m *actisenseRemoteMetrics) snapshot() ActisenseProtocolMetrics {
	completed := m.completed.Load()
	result := ActisenseProtocolMetrics{
		BSTFrames:   make(map[byte]uint64),
		BEMRequests: m.requests.Load(), BEMResponses: m.responses.Load(), BEMCompleted: completed,
		BEMCorrelationMisses: m.correlationMisses.Load(), BEMDuplicateRequests: m.duplicates.Load(),
		BEMTimeouts: m.timeouts.Load(), BEMDeviceErrors: m.deviceErrors.Load(), BEMNegativeAcks: m.negativeAcks.Load(),
		BEMResponseTrainOverflows: m.overflows.Load(), BEMInFlight: m.inFlight.Load(), BEMMaxInFlight: m.maxInFlight.Load(),
		BEMLatencyMinimum: time.Duration(m.latencyMinimum.Load()), BEMLatencyMaximum: time.Duration(m.latencyMaximum.Load()),
	}
	if completed != 0 {
		result.BEMLatencyAverage = time.Duration(m.latencyTotal.Load() / completed)
	}
	return result
}

func newActisenseRemoteManager(client *Client) *actisenseRemoteManager {
	return &actisenseRemoteManager{
		client: client, pending: make(map[actisenseRemoteKey]*actisenseRemotePending),
		diagnostics: newActisenseDiagnosticHub(defaultReceiveBuffer),
	}
}

type actisenseRemoteRequester struct {
	manager *actisenseRemoteManager
	source  uint8
}

// actisenseRemoteMessage avoids the generated PGN-126720 variant's fixed
// 223-byte Data field. The manufacturer envelope is variable length on the
// wire; padding a short BEM command changes the observable payload.
type actisenseRemoteMessage struct {
	info    pgn.MessageInfo
	payload []byte
}

func (m *actisenseRemoteMessage) PGNNumber() uint32                   { return actisenseRemotePGN }
func (m *actisenseRemoteMessage) MessageInfo() pgn.MessageInfo        { return m.info }
func (m *actisenseRemoteMessage) SetMessageInfo(info pgn.MessageInfo) { m.info = info }
func (m *actisenseRemoteMessage) DecodePayload(payload []uint8) error {
	m.payload = append(m.payload[:0], payload...)
	return nil
}
func (m *actisenseRemoteMessage) EncodePayload() ([]uint8, error) {
	return append([]byte(nil), m.payload...), nil
}

func (r *actisenseRemoteRequester) Request(ctx context.Context, command byte, data []byte) (actisense.BEMResponse, error) {
	responses, err := r.manager.request(ctx, r.source, command, data, 0, nil)
	if err != nil {
		if len(responses) != 0 {
			return responses[0], err
		}
		return actisense.BEMResponse{}, err
	}
	return responses[0], nil
}

func (r *actisenseRemoteRequester) RequestMulti(ctx context.Context, command byte, data []byte, inactivity time.Duration, complete func([]actisense.BEMResponse) (bool, error)) ([]actisense.BEMResponse, error) {
	if complete == nil || inactivity <= 0 {
		return nil, errors.New("n2k: Actisense remote multi-response request requires a completion function and positive inactivity timeout")
	}
	return r.manager.request(ctx, r.source, command, data, inactivity, complete)
}

func encodeActisenseRemoteCommand(command byte, data []byte) ([]byte, error) {
	// The two proprietary bytes leave 221 bytes for inner BST in a 223-byte
	// fast packet. Inner Type-1 length is one BEM verb plus its data.
	if len(data) > 218 {
		return nil, fmt.Errorf("n2k: Actisense remote BEM data is %d bytes; maximum is 218", len(data))
	}
	inner := make([]byte, 3, 3+len(data))
	inner[0] = actisense.BSTBEMCommand
	inner[1] = byte(1 + len(data))
	inner[2] = command
	return append(inner, data...), nil
}

func (m *actisenseRemoteManager) request(ctx context.Context, source, command byte, data []byte, inactivity time.Duration, complete func([]actisense.BEMResponse) (bool, error)) (responses []actisense.BEMResponse, err error) {
	started := time.Now()
	m.metrics.requests.Add(1)
	defer func() {
		m.metrics.completed.Add(1)
		m.metrics.observeLatency(time.Since(started))
		if errors.Is(err, context.DeadlineExceeded) {
			m.metrics.timeouts.Add(1)
		}
	}()
	if ctx == nil {
		ctx = context.Background()
	}
	inner, err := encodeActisenseRemoteCommand(command, data)
	if err != nil {
		return nil, err
	}
	m.client.mu.Lock()
	if m.client.closed || m.client.terminalErr != nil {
		err := m.client.terminalErr
		m.client.mu.Unlock()
		if err == nil {
			err = ErrClientClosed
		}
		return nil, err
	}
	local := m.client.sourceAddr
	connectionEpoch := m.client.connectionEpoch
	claimEpoch := m.client.claimEpoch
	key := actisenseRemoteKey{
		source: source, destination: local, connectionEpoch: connectionEpoch, claimEpoch: claimEpoch,
		bstID: actisense.BSTBEMResponse, bemID: command,
	}
	pending := &actisenseRemotePending{results: make(chan actisenseRemoteResult, 257), multi: complete != nil}
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		m.client.mu.Unlock()
		return nil, ErrClientClosed
	}
	if _, exists := m.pending[key]; exists {
		m.mu.Unlock()
		m.client.mu.Unlock()
		m.metrics.duplicates.Add(1)
		return nil, fmt.Errorf("%w for remote source %d BEM 0x%02X", actisense.ErrRequestInFlight, source, command)
	}
	if len(m.pending) >= maxRemotePendingRequests {
		m.mu.Unlock()
		m.client.mu.Unlock()
		return nil, errors.New("n2k: Actisense remote request table is full")
	}
	m.pending[key] = pending
	m.mu.Unlock()
	m.client.mu.Unlock()
	m.metrics.incrementInFlight()
	defer m.metrics.inFlight.Add(^uint64(0))
	remove := func() {
		m.mu.Lock()
		if m.pending[key] == pending {
			delete(m.pending, key)
		}
		m.mu.Unlock()
	}
	defer remove()

	priority := actisenseRemotePriority
	destination := source
	payload := make([]byte, 2, 2+len(inner))
	payload[0], payload[1] = 0x11, 0x99
	payload = append(payload, inner...)
	message := &actisenseRemoteMessage{
		info:    pgn.MessageInfo{PGN: actisenseRemotePGN, Priority: &priority, TargetId: &destination},
		payload: payload,
	}
	if err := m.client.Write(message).WaitContext(ctx); err != nil {
		return nil, err
	}

	var timer *time.Timer
	var inactivityC <-chan time.Time
	if complete != nil {
		timer = time.NewTimer(inactivity)
		defer timer.Stop()
		inactivityC = timer.C
	}
	responses = make([]actisense.BEMResponse, 0, 4)
	for {
		select {
		case result := <-pending.results:
			if result.err != nil {
				if result.response.BEMID != 0 {
					responses = append(responses, result.response)
				}
				return responses, result.err
			}
			responses = append(responses, result.response)
			if complete == nil {
				return responses, nil
			}
			done, completionErr := complete(responses)
			if completionErr != nil || done {
				return responses, completionErr
			}
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(inactivity)
		case <-inactivityC:
			return responses, fmt.Errorf("n2k: Actisense remote BEM 0x%02X response train: %w", command, context.DeadlineExceeded)
		case <-ctx.Done():
			return responses, fmt.Errorf("n2k: Actisense remote BEM 0x%02X: %w", command, ctx.Err())
		}
	}
}

func (m *actisenseRemoteManager) observe(observation raw.Observation) {
	if observation.Kind != raw.KindMessage || observation.PGN != actisenseRemotePGN || len(observation.Payload) < 5 {
		return
	}
	if observation.Payload[0] != 0x11 || observation.Payload[1] != 0x99 || observation.Destination == nil {
		return
	}
	m.client.mu.Lock()
	local := m.client.sourceAddr
	connectionEpoch := m.client.connectionEpoch
	claimEpoch := m.client.claimEpoch
	m.client.mu.Unlock()
	if *observation.Destination != local {
		return
	}
	datagram, err := actisense.DecodeRaw(observation.Payload[2:])
	if err != nil {
		return
	}
	response, ok, err := actisense.DecodeBEMResponse(datagram)
	if !ok || err != nil {
		return
	}
	response.Origin = actisense.BEMOrigin{Path: actisense.BEMPathRemote, Source: observation.Source}
	m.metrics.responses.Add(1)
	if diagnostic, diagnosticOK, diagnosticErr := actisense.DecodeDiagnostic(response); diagnosticOK {
		if diagnosticErr == nil {
			if diagnostic.NegativeAck != nil {
				m.failNegativeAck(response, *diagnostic.NegativeAck, local, connectionEpoch, claimEpoch)
			}
			m.diagnostics.publish(diagnostic)
		}
		return
	}
	key := actisenseRemoteKey{
		source: observation.Source, destination: local, connectionEpoch: connectionEpoch, claimEpoch: claimEpoch,
		bstID: response.BSTID, bemID: response.BEMID,
	}
	m.deliver(key, response)
}

func (m *actisenseRemoteManager) deliver(key actisenseRemoteKey, response actisense.BEMResponse) {
	m.mu.Lock()
	pending := m.pending[key]
	overflow := false
	if pending != nil {
		pending.delivered++
		overflow = pending.delivered > 256
		if !pending.multi || overflow {
			delete(m.pending, key)
		}
	}
	m.mu.Unlock()
	if pending == nil {
		m.metrics.correlationMisses.Add(1)
		return
	}
	result := actisenseRemoteResult{response: response}
	if overflow {
		m.metrics.overflows.Add(1)
		result = actisenseRemoteResult{err: errors.New("n2k: Actisense remote response train exceeds 256 records")}
	} else if response.ErrorCode != 0 {
		m.metrics.deviceErrors.Add(1)
		result.err = &actisense.DeviceError{Command: response.BEMID, Code: response.ErrorCode}
	}
	pending.results <- result
}

func (m *actisenseRemoteManager) failNegativeAck(response actisense.BEMResponse, nack actisense.NegativeAck, destination uint8, connectionEpoch, claimEpoch uint64) {
	m.metrics.negativeAcks.Add(1)
	key := actisenseRemoteKey{
		source: response.Origin.Source, destination: destination, connectionEpoch: connectionEpoch, claimEpoch: claimEpoch,
		bstID: response.BSTID, bemID: byte(nack.UniqueCommandID),
	}
	m.mu.Lock()
	pending := m.pending[key]
	if pending == nil {
		var candidates int
		for candidateKey, candidate := range m.pending {
			if candidateKey.source == key.source && candidateKey.destination == key.destination &&
				candidateKey.connectionEpoch == key.connectionEpoch && candidateKey.claimEpoch == key.claimEpoch && candidateKey.bstID == key.bstID {
				candidates++
				key, pending = candidateKey, candidate
			}
		}
		if candidates != 1 {
			pending = nil
		}
	}
	if pending != nil {
		delete(m.pending, key)
	}
	m.mu.Unlock()
	if pending != nil {
		pending.results <- actisenseRemoteResult{err: &actisense.NegativeAckError{
			Command: key.bemID, UniqueCommandID: nack.UniqueCommandID, DeviceCode: nack.ErrorCode,
		}}
	}
}

func (m *actisenseRemoteManager) invalidate(err error) {
	if m == nil {
		return
	}
	if err == nil {
		err = ErrActisenseRemoteEpochChanged
	}
	m.mu.Lock()
	pending := m.pending
	m.pending = make(map[actisenseRemoteKey]*actisenseRemotePending)
	m.mu.Unlock()
	for _, request := range pending {
		request.results <- actisenseRemoteResult{err: err}
	}
}

func (m *actisenseRemoteManager) close(err error) {
	if m == nil {
		return
	}
	m.mu.Lock()
	m.closed = true
	m.mu.Unlock()
	m.invalidate(err)
	m.diagnostics.close(err)
}
