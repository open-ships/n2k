package actisense

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
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

// Session serializes one Actisense connection epoch. It owns the sole wire
// writer and a bounded BEM correlation table with at most one request per
// response/BEM key.
type Session struct {
	write func([]byte) error

	writeMu sync.Mutex
	mu      sync.Mutex
	pending map[byte]chan responseResult
	done    chan struct{}
	doneErr error
	once    sync.Once

	onDatagram    func(Datagram)
	onDiagnostic  func(Diagnostic)
	onDecodeError func(DecodeError)
}

type SessionConfig struct {
	Write         func([]byte) error
	OnDatagram    func(Datagram)
	OnDiagnostic  func(Diagnostic)
	OnDecodeError func(DecodeError)
}

func NewSession(config SessionConfig) *Session {
	return &Session{
		write:         config.Write,
		pending:       make(map[byte]chan responseResult),
		done:          make(chan struct{}),
		onDatagram:    config.OnDatagram,
		onDiagnostic:  config.OnDiagnostic,
		onDecodeError: config.OnDecodeError,
	}
}

// Run reads one connection until EOF or failure. It must be called once.
func (s *Session) Run(reader io.Reader) error {
	parser := NewParser()
	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			parser.Feed(buf[:n], s.handleDatagram, s.handleDecodeError)
		}
		if err != nil {
			parser.End(s.handleDecodeError)
			if errors.Is(err, io.EOF) {
				err = nil
			}
			s.finish(err)
			return err
		}
	}
}

func (s *Session) handleDatagram(datagram Datagram) {
	if response, ok, err := DecodeBEMResponse(datagram); ok {
		if err != nil {
			s.handleDecodeError(DecodeError{Kind: DecodeLength, ID: datagram.ID, Err: err})
		} else if diagnostic, isDiagnostic, diagnosticErr := DecodeDiagnostic(response); isDiagnostic {
			if diagnosticErr != nil {
				s.handleDecodeError(DecodeError{Kind: DecodeLength, ID: datagram.ID, Err: diagnosticErr})
			} else {
				if diagnostic.NegativeAck != nil {
					s.failNegativeAck(*diagnostic.NegativeAck)
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

func (s *Session) failNegativeAck(nack NegativeAck) {
	rejected := byte(nack.UniqueCommandID)
	s.mu.Lock()
	command := rejected
	ch := s.pending[command]
	if ch == nil && len(s.pending) == 1 {
		for soleCommand, sole := range s.pending {
			command, ch = soleCommand, sole
		}
	}
	if ch != nil {
		delete(s.pending, command)
	}
	s.mu.Unlock()
	if ch != nil {
		ch <- responseResult{err: &NegativeAckError{
			Command: command, UniqueCommandID: nack.UniqueCommandID, DeviceCode: nack.ErrorCode,
		}}
	}
}

func (s *Session) handleDecodeError(err DecodeError) {
	if s.onDecodeError != nil {
		s.onDecodeError(err)
	}
}

func (s *Session) deliverResponse(response BEMResponse) {
	s.mu.Lock()
	ch := s.pending[response.BEMID]
	if ch != nil {
		delete(s.pending, response.BEMID)
	}
	s.mu.Unlock()
	if ch == nil {
		return
	}
	result := responseResult{response: response}
	if response.ErrorCode != 0 {
		result.err = &DeviceError{Command: response.BEMID, Code: response.ErrorCode}
	}
	ch <- result
}

func (s *Session) finish(err error) {
	s.once.Do(func() {
		s.mu.Lock()
		if err == nil {
			err = ErrSessionClosed
		}
		s.doneErr = err
		pending := s.pending
		s.pending = make(map[byte]chan responseResult)
		close(s.done)
		s.mu.Unlock()
		for _, ch := range pending {
			ch <- responseResult{err: err}
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
		return errors.New("actisense: session has no writer")
	}
	return write(buf)
}

// Request sends one local BEM command and waits for its A0 response. The
// correlation key intentionally excludes the response sequence byte: gateway
// firmware uses the response BST group, BEM verb, and origin, and local A1/A0
// sessions permit one in-flight request per verb.
func (s *Session) Request(ctx context.Context, command byte, data []byte) (BEMResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	responseCh := make(chan responseResult, 1)
	s.mu.Lock()
	select {
	case <-s.done:
		err := s.doneErr
		s.mu.Unlock()
		if err == nil {
			err = ErrSessionClosed
		}
		return BEMResponse{}, err
	default:
	}
	if _, exists := s.pending[command]; exists {
		s.mu.Unlock()
		return BEMResponse{}, fmt.Errorf("%w for BEM 0x%02X", ErrRequestInFlight, command)
	}
	if len(s.pending) >= maxPendingBEMRequests {
		s.mu.Unlock()
		return BEMResponse{}, errors.New("actisense: BEM request table is full")
	}
	s.pending[command] = responseCh
	s.mu.Unlock()

	encoded, err := EncodeBEMRequest(command, data)
	if err == nil {
		err = s.Write(encoded)
	}
	if err != nil {
		s.removePending(command, responseCh)
		return BEMResponse{}, err
	}

	select {
	case result := <-responseCh:
		return result.response, result.err
	case <-ctx.Done():
		s.removePending(command, responseCh)
		return BEMResponse{}, fmt.Errorf("actisense: BEM 0x%02X: %w", command, ctx.Err())
	case <-s.done:
		s.removePending(command, responseCh)
		s.mu.Lock()
		doneErr := s.doneErr
		s.mu.Unlock()
		if doneErr == nil {
			doneErr = ErrSessionClosed
		}
		return BEMResponse{}, doneErr
	}
}

func (s *Session) removePending(command byte, expected chan responseResult) {
	s.mu.Lock()
	if s.pending[command] == expected {
		delete(s.pending, command)
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
