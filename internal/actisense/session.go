package actisense

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

var (
	ErrSessionClosed   = errors.New("actisense: session closed")
	ErrRequestInFlight = errors.New("actisense: request already in flight")
)

// NegativeAckError reports an asynchronous F4 rejection correlated to an
// in-flight command.
type NegativeAckError struct {
	Command         byte
	UniqueCommandID uint32
	DeviceCode      int32
}

func (e *NegativeAckError) Error() string {
	return fmt.Sprintf("actisense: device rejected BEM command 0x%02X (id 0x%08X, error %d)", e.Command, e.UniqueCommandID, e.DeviceCode)
}

type responseResult struct {
	response BEMResponse
	err      error
}

type responseKey struct {
	bstID  byte
	bemID  byte
	origin BEMOrigin
}

type pendingRequest struct {
	results   chan responseResult
	multi     bool
	delivered int
}

// Session serializes one Actisense connection epoch. It owns the sole wire
// writer and a bounded BEM correlation table with at most one request per
// response/BEM key.
type Session struct {
	write func([]byte) error

	writeMu sync.Mutex
	mu      sync.Mutex
	pending map[responseKey]*pendingRequest
	done    chan struct{}
	doneErr error
	once    sync.Once

	onDatagram    func(Datagram)
	onDiagnostic  func(Diagnostic)
	onDecodeError func(DecodeError)
	onWireBytes   func(WireDirection, time.Time, []byte)
	onUnframed    func([]byte)
	metrics       sessionMetrics
}

type WireDirection uint8

const (
	WireReceived WireDirection = iota
	WireTransmitted
)

type SessionConfig struct {
	Write         func([]byte) error
	OnDatagram    func(Datagram)
	OnDiagnostic  func(Diagnostic)
	OnDecodeError func(DecodeError)
	OnWireBytes   func(WireDirection, time.Time, []byte)
	OnUnframed    func([]byte)
}

func NewSession(config SessionConfig) *Session {
	return &Session{
		write:         config.Write,
		pending:       make(map[responseKey]*pendingRequest),
		done:          make(chan struct{}),
		onDatagram:    config.OnDatagram,
		onDiagnostic:  config.OnDiagnostic,
		onDecodeError: config.OnDecodeError,
		onWireBytes:   config.OnWireBytes,
		onUnframed:    config.OnUnframed,
	}
}

// Metrics returns a concurrency-safe point-in-time snapshot.
func (s *Session) Metrics() SessionMetrics {
	if s == nil {
		return SessionMetrics{BSTFrames: make(map[byte]uint64)}
	}
	return s.metrics.snapshot()
}

// Run reads one connection until EOF or failure. It must be called once.
func (s *Session) Run(reader io.Reader) error {
	parser := NewParser()
	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		s.metrics.transportReadCalls.Add(1)
		if n > 0 {
			s.metrics.transportReadBytes.Add(uint64(n))
			if s.onWireBytes != nil {
				s.onWireBytes(WireReceived, time.Now(), append([]byte(nil), buf[:n]...))
			}
			parser.FeedWithUnframed(buf[:n], s.handleDatagram, s.handleDecodeError, s.handleUnframed)
		}
		if err != nil {
			parser.EndWithUnframed(s.handleDecodeError, s.handleUnframed)
			if errors.Is(err, io.EOF) {
				err = nil
			} else {
				s.metrics.transportReadErrors.Add(1)
			}
			s.finish(err)
			return err
		}
	}
}

func (s *Session) handleDatagram(datagram Datagram) {
	s.metrics.datagrams.Add(1)
	s.metrics.bstFrames[datagram.ID].Add(1)
	if response, ok, err := DecodeBEMResponse(datagram); ok {
		s.metrics.bemResponses.Add(1)
		if err != nil {
			s.handleDecodeError(DecodeError{Kind: DecodeLength, ID: datagram.ID, Err: err})
		} else if diagnostic, isDiagnostic, diagnosticErr := DecodeDiagnostic(response); isDiagnostic {
			if diagnosticErr != nil {
				s.handleDecodeError(DecodeError{Kind: DecodeLength, ID: datagram.ID, Err: diagnosticErr})
			} else {
				if diagnostic.NegativeAck != nil {
					s.failNegativeAck(response, *diagnostic.NegativeAck)
				}
				if s.onDiagnostic != nil {
					s.onDiagnostic(diagnostic)
				}
			}
		} else {
			s.deliverResponse(response)
		}
	}
	if s.onDatagram != nil {
		s.onDatagram(datagram.Clone())
	}
}

func (s *Session) failNegativeAck(response BEMResponse, nack NegativeAck) {
	s.metrics.bemNegativeAcks.Add(1)
	rejected := byte(nack.UniqueCommandID)
	s.mu.Lock()
	key := responseKey{bstID: response.BSTID, bemID: rejected, origin: response.Origin}
	pending := s.pending[key]
	if pending == nil {
		var candidateCount int
		for candidateKey, candidate := range s.pending {
			if candidateKey.bstID == response.BSTID && candidateKey.origin == response.Origin {
				candidateCount++
				key, pending = candidateKey, candidate
			}
		}
		if candidateCount != 1 {
			pending = nil
		}
	}
	if pending != nil {
		delete(s.pending, key)
	}
	s.mu.Unlock()
	if pending != nil {
		pending.results <- responseResult{err: &NegativeAckError{
			Command: key.bemID, UniqueCommandID: nack.UniqueCommandID, DeviceCode: nack.ErrorCode,
		}}
	}
}

func (s *Session) handleDecodeError(err DecodeError) {
	switch err.Kind {
	case DecodeFraming:
		s.metrics.framingErrors.Add(1)
	case DecodeChecksum:
		s.metrics.checksumErrors.Add(1)
	case DecodeLength:
		s.metrics.lengthErrors.Add(1)
	case DecodeOversize:
		s.metrics.oversizeErrors.Add(1)
	}
	if s.onDecodeError != nil {
		s.onDecodeError(err)
	}
}

func (s *Session) handleUnframed(data []byte) {
	s.metrics.unframedBytes.Add(uint64(len(data)))
	if s.onUnframed != nil {
		s.onUnframed(data)
	}
}

func (s *Session) deliverResponse(response BEMResponse) {
	key := responseKey{bstID: response.BSTID, bemID: response.BEMID, origin: response.Origin}
	s.mu.Lock()
	pending := s.pending[key]
	var overflow bool
	if pending != nil {
		pending.delivered++
		overflow = pending.delivered > maxBEMResponseTrain
		if !pending.multi || overflow {
			delete(s.pending, key)
		}
	}
	s.mu.Unlock()
	if pending == nil {
		s.metrics.bemCorrelationMisses.Add(1)
		return
	}
	result := responseResult{response: response}
	if overflow {
		s.metrics.bemResponseTrainOverflows.Add(1)
		result = responseResult{err: fmt.Errorf("actisense: BEM 0x%02X response train exceeds %d records", response.BEMID, maxBEMResponseTrain)}
	}
	if response.ErrorCode != 0 {
		s.metrics.bemDeviceErrors.Add(1)
		result.err = &DeviceError{Command: response.BEMID, Code: response.ErrorCode}
	}
	pending.results <- result
}

func (s *Session) finish(err error) {
	s.once.Do(func() {
		s.mu.Lock()
		if err == nil {
			err = ErrSessionClosed
		}
		s.doneErr = err
		pending := s.pending
		s.pending = make(map[responseKey]*pendingRequest)
		close(s.done)
		s.mu.Unlock()
		for _, request := range pending {
			request.results <- responseResult{err: err}
		}
	})
}

// Close marks the session terminal and releases pending requests. The owner of
// the underlying transport remains responsible for closing it to unblock Run.
func (s *Session) Close(err error) { s.finish(err) }

// Write serializes one already-framed BDTP unit with all BEM and bus traffic.
func (s *Session) Write(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	s.mu.Lock()
	select {
	case <-s.done:
		err := s.doneErr
		s.mu.Unlock()
		if err == nil {
			return ErrSessionClosed
		}
		return err
	default:
	}
	write := s.write
	s.mu.Unlock()
	if write == nil {
		s.metrics.transportWriteErrors.Add(1)
		return errors.New("actisense: session has no writer")
	}
	s.metrics.transportWriteCalls.Add(1)
	if s.onWireBytes != nil {
		s.onWireBytes(WireTransmitted, time.Now(), append([]byte(nil), buf...))
	}
	err := write(buf)
	if err != nil {
		s.metrics.transportWriteErrors.Add(1)
	} else {
		s.metrics.transportWriteBytes.Add(uint64(len(buf)))
	}
	return err
}

// Request sends one local BEM command and waits for its A0 response. The
// correlation key intentionally excludes the response sequence byte: gateway
// firmware uses the response BST group, BEM verb, and origin, and local A1/A0
// sessions permit one in-flight request per verb.
func (s *Session) Request(ctx context.Context, command byte, data []byte) (BEMResponse, error) {
	responses, err := s.request(ctx, command, data, 0, nil)
	if err != nil {
		if len(responses) != 0 {
			return responses[0], err
		}
		return BEMResponse{}, err
	}
	return responses[0], nil
}

// RequestMulti sends one local BEM command and collects a bounded response
// train. complete is evaluated in wire order after every response and must
// return true for the final response. inactivity bounds the gap between
// responses; the caller context remains the overall deadline.
func (s *Session) RequestMulti(ctx context.Context, command byte, data []byte, inactivity time.Duration, complete func([]BEMResponse) (bool, error)) ([]BEMResponse, error) {
	if complete == nil {
		return nil, errors.New("actisense: multi-response request requires a completion function")
	}
	if inactivity <= 0 {
		return nil, errors.New("actisense: multi-response inactivity timeout must be positive")
	}
	return s.request(ctx, command, data, inactivity, complete)
}

func (s *Session) request(ctx context.Context, command byte, data []byte, inactivity time.Duration, complete func([]BEMResponse) (bool, error)) (responses []BEMResponse, err error) {
	started := time.Now()
	s.metrics.bemRequests.Add(1)
	defer func() {
		s.metrics.bemCompleted.Add(1)
		s.metrics.observeLatency(time.Since(started))
		if errors.Is(err, context.DeadlineExceeded) {
			s.metrics.bemTimeouts.Add(1)
		}
	}()
	if ctx == nil {
		ctx = context.Background()
	}
	key := responseKey{bstID: BSTBEMResponse, bemID: command, origin: LocalBEMOrigin}
	request := &pendingRequest{
		results: make(chan responseResult, maxBEMResponseTrain+1),
		multi:   complete != nil,
	}
	s.mu.Lock()
	select {
	case <-s.done:
		err := s.doneErr
		s.mu.Unlock()
		if err == nil {
			err = ErrSessionClosed
		}
		return nil, err
	default:
	}
	if _, exists := s.pending[key]; exists {
		s.mu.Unlock()
		s.metrics.bemDuplicateRequests.Add(1)
		return nil, fmt.Errorf("%w for BST 0x%02X BEM 0x%02X origin %+v", ErrRequestInFlight, key.bstID, command, key.origin)
	}
	if len(s.pending) >= maxPendingBEMRequests {
		s.mu.Unlock()
		return nil, errors.New("actisense: BEM request table is full")
	}
	s.pending[key] = request
	s.mu.Unlock()
	s.metrics.incrementInFlight()
	defer s.metrics.bemInFlight.Add(^uint64(0))

	encoded, err := EncodeBEMRequest(command, data)
	if err == nil {
		err = s.Write(encoded)
	}
	if err != nil {
		s.removePending(key, request)
		return nil, err
	}

	var timer *time.Timer
	var inactivityC <-chan time.Time
	if complete != nil {
		timer = time.NewTimer(inactivity)
		defer timer.Stop()
		inactivityC = timer.C
	}
	responses = make([]BEMResponse, 0, 4)
	defer s.removePending(key, request)
	for {
		select {
		case result := <-request.results:
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
			if completionErr != nil {
				return nil, completionErr
			}
			if done {
				return responses, nil
			}
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(inactivity)
		case <-inactivityC:
			return responses, fmt.Errorf("actisense: BEM 0x%02X response train: %w", command, context.DeadlineExceeded)
		case <-ctx.Done():
			return responses, fmt.Errorf("actisense: BEM 0x%02X: %w", command, ctx.Err())
		case <-s.done:
			s.mu.Lock()
			doneErr := s.doneErr
			s.mu.Unlock()
			if doneErr == nil {
				doneErr = ErrSessionClosed
			}
			return responses, doneErr
		}
	}
}

func (s *Session) removePending(key responseKey, expected *pendingRequest) {
	s.mu.Lock()
	if s.pending[key] == expected {
		delete(s.pending, key)
	}
	s.mu.Unlock()
}

func (s *Session) GetOperatingMode(ctx context.Context) (OperatingMode, error) {
	response, err := s.Request(ctx, BEMOperatingMode, OperatingModeRequest())
	if err != nil {
		return 0, err
	}
	return DecodeOperatingMode(response)
}

func (s *Session) SetOperatingMode(ctx context.Context, mode OperatingMode) error {
	response, err := s.Request(ctx, BEMOperatingMode, OperatingModeSet(mode))
	if err != nil {
		return err
	}
	accepted, err := DecodeOperatingMode(response)
	if err != nil {
		return err
	}
	if accepted != mode {
		return fmt.Errorf("actisense: gateway acknowledged operating mode %d, requested %d", accepted, mode)
	}
	return nil
}

func (s *Session) GetTxPGN(ctx context.Context, pgn uint32) (TxPGNState, error) {
	response, err := s.Request(ctx, BEMTxPGNEnable, TxPGNEnableGet(pgn))
	if err != nil {
		return TxPGNState{}, err
	}
	return DecodeTxPGNState(response)
}

func (s *Session) SetTxPGN(ctx context.Context, pgn uint32, enabled bool) error {
	response, err := s.Request(ctx, BEMTxPGNEnable, TxPGNEnableSet(pgn, enabled))
	if err != nil {
		return err
	}
	state, err := DecodeTxPGNState(response)
	if err != nil {
		return err
	}
	want := uint8(0)
	if enabled {
		want = 1
	}
	if state.PGN != pgn || state.Enabled != want {
		return fmt.Errorf("actisense: gateway acknowledged Tx PGN %d state %d; requested PGN %d state %d", state.PGN, state.Enabled, pgn, want)
	}
	return nil
}

func (s *Session) ActivatePGNLists(ctx context.Context) error {
	_, err := s.Request(ctx, BEMActivatePGNLists, nil)
	return err
}
