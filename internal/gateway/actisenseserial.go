package gateway

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/raw"
	"go.bug.st/serial"
)

func newActisenseSerialStreamBus(log *slog.Logger, port string, mode actisense.OperatingMode, rawCAN bool) *actisenseStreamBus {
	open := func(context.Context) (actisenseConnection, error) {
		connection, err := serial.Open(port, &serial.Mode{
			BaudRate: 115200,
			DataBits: 8,
			Parity:   serial.NoParity,
			StopBits: serial.OneStopBit,
		})
		if err != nil {
			return nil, fmt.Errorf("actisense: opening serial port %s: %w", port, err)
		}
		return connection, nil
	}
	return newActisenseStreamBus(log, port, "serial:"+port, open, nil, mode, rawCAN)
}

// ActisenseSerialBus is the gateway-owned BST-93/94 message session over a
// direct 115200 8N1 serial connection.
type ActisenseSerialBus struct{ stream *actisenseStreamBus }

func NewActisenseSerialBus(log *slog.Logger, port string) *ActisenseSerialBus {
	return &ActisenseSerialBus{stream: newActisenseSerialStreamBus(log, port, actisense.ModeTransferReceiveAll, false)}
}

func (b *ActisenseSerialBus) Run(ctx context.Context, handler func(can.Frame)) error {
	return b.stream.Run(ctx, handler)
}
func (b *ActisenseSerialBus) RunObservations(ctx context.Context, handler func(raw.Observation)) error {
	return b.stream.RunObservations(ctx, handler)
}
func (b *ActisenseSerialBus) WriteFrame(frame can.Frame) error { return b.stream.writeFrame(frame) }
func (b *ActisenseSerialBus) WriteMessage(pgnNumber uint32, priority, _ uint8, destination uint8, payload []byte) error {
	return b.stream.writeMessage(pgnNumber, priority, destination, payload)
}
func (b *ActisenseSerialBus) Ready() <-chan struct{} { return b.stream.Ready() }
func (b *ActisenseSerialBus) Close() error           { return b.stream.Close() }

// ActisenseRawSerialBus is a source-authoritative BST-95 raw CAN Bus over a
// direct 115200 8N1 serial connection.
type ActisenseRawSerialBus struct{ stream *actisenseStreamBus }

func NewActisenseRawSerialBus(log *slog.Logger, port string) *ActisenseRawSerialBus {
	return &ActisenseRawSerialBus{stream: newActisenseSerialStreamBus(log, port, actisense.ModeCANPacket, true)}
}

func (b *ActisenseRawSerialBus) Run(ctx context.Context, handler func(can.Frame)) error {
	return b.stream.Run(ctx, handler)
}
func (b *ActisenseRawSerialBus) RunObservations(ctx context.Context, handler func(raw.Observation)) error {
	return b.stream.RunObservations(ctx, handler)
}
func (b *ActisenseRawSerialBus) WriteFrame(frame can.Frame) error { return b.stream.writeFrame(frame) }
func (b *ActisenseRawSerialBus) Ready() <-chan struct{}           { return b.stream.Ready() }
func (b *ActisenseRawSerialBus) Close() error                     { return b.stream.Close() }
