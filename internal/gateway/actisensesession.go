package gateway

import (
	"context"
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
	return s.stream.Request(ctx, command, data)
}

func (s *ActisenseGatewaySession) RequestMulti(ctx context.Context, command byte, data []byte, inactivity time.Duration, complete func([]actisense.BEMResponse) (bool, error)) ([]actisense.BEMResponse, error) {
	return s.stream.RequestMulti(ctx, command, data, inactivity, complete)
}

func (s *ActisenseGatewaySession) WriteMessage(pgn uint32, priority, destination uint8, payload []byte) error {
	return s.stream.writeMessage(pgn, priority, destination, payload)
}

func (s *ActisenseGatewaySession) WriteMessageContext(ctx context.Context, pgn uint32, priority, destination uint8, payload []byte) error {
	return s.stream.writeMessageContext(ctx, pgn, priority, destination, payload)
}

func (s *ActisenseGatewaySession) Close() error { return s.stream.Close() }
