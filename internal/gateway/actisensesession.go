package gateway

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/raw"
)

// ActisenseGatewaySession exposes the gateway-owned BST-93/94 and BEM
// session without claiming that it is a source-authoritative CAN Bus.
type ActisenseGatewaySession struct {
	stream *actisenseStreamBus
}

var ErrActisenseNotReady = errors.New("actisense: gateway session is not ready")

// These hooks are installed before Run. They execute on the sole reader.
func (s *ActisenseGatewaySession) PreserveOperatingMode() {
	s.stream.preserveMode = true
}

func (s *ActisenseGatewaySession) SetReconnectPolicy(policy *ReconnectPolicy) {
	s.stream.reconnect = policy
}

func (s *ActisenseGatewaySession) SetMessageObserver(observer func(actisense.Message)) {
	s.stream.mu.Lock()
	s.stream.messageObserver = observer
	s.stream.mu.Unlock()
}

func (s *ActisenseGatewaySession) SetModeObserver(observer func(actisense.OperatingMode)) {
	s.stream.mu.Lock()
	s.stream.modeObserver = observer
	s.stream.mu.Unlock()
}

func NewActisenseTCPGatewaySession(log *slog.Logger, address string, reconnect *ReconnectPolicy, mode actisense.OperatingMode) *ActisenseGatewaySession {
	return &ActisenseGatewaySession{stream: newActisenseTCPStreamBus(log, address, reconnect, mode, false)}
}

func NewActisenseSerialGatewaySession(log *slog.Logger, port string, settings ActisenseSerialSettings, mode actisense.OperatingMode) *ActisenseGatewaySession {
	return &ActisenseGatewaySession{stream: newActisenseSerialStreamBus(log, port, settings, mode, false)}
}

func NewActisenseCustomGatewaySession(log *slog.Logger, endpoint, adapterID string, open ActisenseOpen, reconnect *ReconnectPolicy, mode actisense.OperatingMode) *ActisenseGatewaySession {
	internalOpen := func(ctx context.Context) (actisenseConnection, error) {
		return open(ctx)
	}
	return &ActisenseGatewaySession{stream: newActisenseStreamBus(log, endpoint, adapterID, internalOpen, reconnect, mode, false)}
}

func (s *ActisenseGatewaySession) Run(ctx context.Context, handler func(can.Frame)) error {
	return s.stream.Run(ctx, handler)
}

func (s *ActisenseGatewaySession) RunObservations(ctx context.Context, handler func(raw.Observation)) error {
	return s.stream.RunObservations(ctx, handler)
}

func (s *ActisenseGatewaySession) Ready() <-chan struct{} { return s.stream.Ready() }

func (s *ActisenseGatewaySession) SetConnectionObserver(observer func(bool, uint64)) {
	s.stream.SetConnectionObserver(observer)
}

func (s *ActisenseGatewaySession) SetDiagnosticObserver(observer func(actisense.Diagnostic)) {
	s.stream.SetDiagnosticObserver(observer)
}

func (s *ActisenseGatewaySession) SetWireObserver(observer func(actisense.WireDirection, time.Time, []byte)) {
	s.stream.SetWireObserver(observer)
}

func (s *ActisenseGatewaySession) SetCommandTimeout(timeout time.Duration) error {
	return s.stream.SetCommandTimeout(timeout)
}

func (s *ActisenseGatewaySession) Metrics() ActisenseGatewayMetrics {
	return s.stream.Metrics()
}

func (s *ActisenseGatewaySession) Request(ctx context.Context, command byte, data []byte) (actisense.BEMResponse, error) {
	epoch, err := s.currentEpoch(0)
	if err != nil {
		return actisense.BEMResponse{}, err
	}
	return epoch.session.Request(ctx, command, data)
}

func (s *ActisenseGatewaySession) currentEpoch(number uint64) (*actisenseEpoch, error) {
	s.stream.mu.Lock()
	defer s.stream.mu.Unlock()
	if s.stream.closed {
		return nil, actisense.ErrSessionClosed
	}
	if s.stream.epoch == nil || (number != 0 && s.stream.epochNum != number) {
		return nil, ErrActisenseNotReady
	}
	return s.stream.epoch, nil
}

// EpochRequester binds a transaction to one acknowledged connection. Its
// requests fail with that session rather than moving onto a reconnect.
func (s *ActisenseGatewaySession) EpochRequester(epoch uint64) (actisense.Requester, error) {
	s.stream.mu.Lock()
	defer s.stream.mu.Unlock()
	if s.stream.closed || s.stream.epoch == nil || s.stream.epochNum != epoch {
		return nil, fmt.Errorf("actisense: connection epoch %d is no longer ready", epoch)
	}
	return s.stream.epoch.session, nil
}

func (s *ActisenseGatewaySession) RequestMulti(ctx context.Context, command byte, data []byte, inactivity time.Duration, complete func([]actisense.BEMResponse) (bool, error)) ([]actisense.BEMResponse, error) {
	epoch, err := s.currentEpoch(0)
	if err != nil {
		return nil, err
	}
	return epoch.session.RequestMulti(ctx, command, data, inactivity, complete)
}

func (s *ActisenseGatewaySession) WriteMessage(pgn uint32, priority, destination uint8, payload []byte) error {
	return s.WriteMessageContext(context.Background(), pgn, priority, destination, payload)
}

func (s *ActisenseGatewaySession) WriteMessageContext(ctx context.Context, pgn uint32, priority, destination uint8, payload []byte) error {
	return s.WriteMessageEpoch(ctx, 0, pgn, priority, destination, payload)
}

func (s *ActisenseGatewaySession) WriteMessageEpoch(ctx context.Context, number uint64, pgn uint32, priority, destination uint8, payload []byte) error {
	epoch, err := s.currentEpoch(number)
	if err != nil {
		return err
	}
	mode := actisense.OperatingMode(epoch.mode.Load())
	if mode != actisense.ModeTransferNormal && mode != actisense.ModeTransferReceiveAll {
		return fmt.Errorf("actisense: gateway PGN sends require operating mode 1 or 2; current mode is %d", mode)
	}
	wire, err := actisense.EncodeMessage94(actisense.Message{PGN: pgn, Priority: priority, Destination: destination, Data: payload})
	if err != nil {
		return err
	}
	return epoch.session.WriteContext(ctx, wire)
}

// WriteContext captures the current connection once and never waits for or
// retries on a reconnect. The public caller owns and bounds the wire snapshot.
func (s *ActisenseGatewaySession) WriteContext(ctx context.Context, wire []byte) error {
	epoch, err := s.currentEpoch(0)
	if err != nil {
		return err
	}
	return epoch.session.WriteContext(ctx, wire)
}

func (s *ActisenseGatewaySession) Close() error { return s.stream.Close() }
