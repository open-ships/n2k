package gateway

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/internal/serialio"
	"github.com/open-ships/n2k/raw"
	"go.bug.st/serial"
)

// ActisenseSerialSettings is the transport-neutral serial configuration passed
// by the public Adapter.
type ActisenseSerialSettings struct {
	BaudRate int
	DataBits int
	Parity   uint8
	StopBits uint8
}

func openActisenseSerialConnection(path string, mode *serial.Mode) (actisenseConnection, error) {
	return serialio.Open(path, mode)
}

func normalizeActisenseSerialSettings(settings []ActisenseSerialSettings) ActisenseSerialSettings {
	if len(settings) == 0 {
		return ActisenseSerialSettings{BaudRate: 115200, DataBits: 8}
	}
	return settings[0]
}

func newActisenseSerialStreamBus(log *slog.Logger, port string, settings ActisenseSerialSettings, mode actisense.OperatingMode, rawCAN bool) *actisenseStreamBus {
	open := func(context.Context) (actisenseConnection, error) {
		connection, err := openActisenseSerialConnection(port, &serial.Mode{
			BaudRate: settings.BaudRate,
			DataBits: settings.DataBits,
			Parity:   serial.Parity(settings.Parity),
			StopBits: serial.StopBits(settings.StopBits),
		})
		if err != nil {
			return nil, fmt.Errorf("actisense: opening serial port %s: %w", port, err)
		}
		return connection, nil
	}
	return newActisenseStreamBus(log, port, "serial:"+port, open, nil, mode, rawCAN)
}

func runPassiveActisenseSerial(ctx context.Context, port string, settings ActisenseSerialSettings, read func(io.Reader, func(raw.Observation)) error, handler func(raw.Observation)) error {
	connection, err := serialio.Open(port, &serial.Mode{
		BaudRate: settings.BaudRate,
		DataBits: settings.DataBits,
		Parity:   serial.Parity(settings.Parity),
		StopBits: serial.StopBits(settings.StopBits),
	})
	if err != nil {
		return fmt.Errorf("actisense: opening serial port %s: %w", port, err)
	}
	defer func() { _ = connection.Close() }()
	watchDone := make(chan struct{})
	defer close(watchDone)
	go func() {
		select {
		case <-ctx.Done():
			_ = connection.Close()
		case <-watchDone:
		}
	}()
	err = read(connection, func(observation raw.Observation) {
		observation.AdapterID = "serial:" + port
		observation.NetworkID = port
		if handler != nil {
			handler(observation)
		}
	})
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}

// RunPassiveActisenseSerialObservations reads binary BST without changing the
// device's operating mode.
func RunPassiveActisenseSerialObservations(ctx context.Context, port string, settings ActisenseSerialSettings, handler func(raw.Observation)) error {
	return runPassiveActisenseSerial(ctx, port, settings, ReadActisenseObservations, handler)
}

// RunPassiveActisenseN2KASCIISerialObservations reads an already-configured
// assembled N2K ASCII output without a BEM mode mutation.
func RunPassiveActisenseN2KASCIISerialObservations(ctx context.Context, port string, settings ActisenseSerialSettings, handler func(raw.Observation)) error {
	return runPassiveActisenseSerial(ctx, port, settings, ReadActisenseN2KASCIIObservations, handler)
}

// ActisenseRawSerialBus is a source-authoritative BST-95 raw CAN Bus over a
// direct 115200 8N1 serial connection.
type ActisenseRawSerialBus struct{ stream *actisenseStreamBus }

func NewActisenseRawSerialBus(log *slog.Logger, port string, settings ...ActisenseSerialSettings) *ActisenseRawSerialBus {
	return &ActisenseRawSerialBus{stream: newActisenseSerialStreamBus(log, port, normalizeActisenseSerialSettings(settings), actisense.ModeCANPacket, true)}
}

func (b *ActisenseRawSerialBus) Run(ctx context.Context, handler func(can.Frame)) error {
	return b.stream.Run(ctx, handler)
}
func (b *ActisenseRawSerialBus) RunObservations(ctx context.Context, handler func(raw.Observation)) error {
	return b.stream.RunObservations(ctx, handler)
}
func (b *ActisenseRawSerialBus) WriteFrame(frame can.Frame) error { return b.stream.writeFrame(frame) }
func (b *ActisenseRawSerialBus) WriteFrameContext(ctx context.Context, frame can.Frame) error {
	return b.stream.writeFrameContext(ctx, frame)
}
func (b *ActisenseRawSerialBus) Ready() <-chan struct{} { return b.stream.Ready() }
func (b *ActisenseRawSerialBus) Close() error           { return b.stream.Close() }
func (b *ActisenseRawSerialBus) Metrics() ActisenseGatewayMetrics {
	return b.stream.Metrics()
}

// ActisenseCANASCIISerialBus is a source-authoritative mode-6 CAN ASCII Bus.
type ActisenseCANASCIISerialBus struct{ stream *actisenseStreamBus }

func NewActisenseCANASCIISerialBus(log *slog.Logger, port string, settings ...ActisenseSerialSettings) *ActisenseCANASCIISerialBus {
	return &ActisenseCANASCIISerialBus{stream: newActisenseSerialStreamBus(log, port, normalizeActisenseSerialSettings(settings), actisense.ModeCANPacketASCII, true)}
}

func (b *ActisenseCANASCIISerialBus) Run(ctx context.Context, handler func(can.Frame)) error {
	return b.stream.Run(ctx, handler)
}
func (b *ActisenseCANASCIISerialBus) RunObservations(ctx context.Context, handler func(raw.Observation)) error {
	return b.stream.RunObservations(ctx, handler)
}
func (b *ActisenseCANASCIISerialBus) WriteFrame(frame can.Frame) error {
	return b.stream.writeFrame(frame)
}
func (b *ActisenseCANASCIISerialBus) WriteFrameContext(ctx context.Context, frame can.Frame) error {
	return b.stream.writeFrameContext(ctx, frame)
}
func (b *ActisenseCANASCIISerialBus) Ready() <-chan struct{} { return b.stream.Ready() }
func (b *ActisenseCANASCIISerialBus) Close() error           { return b.stream.Close() }
func (b *ActisenseCANASCIISerialBus) Metrics() ActisenseGatewayMetrics {
	return b.stream.Metrics()
}
