package n2k

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/gateway"
	"github.com/open-ships/n2k/raw"
	"go.bug.st/serial/enumerator"
)

type serialSource struct {
	port   string
	format StreamFormat
}

func (s *serialSource) run(ctx context.Context, log *slog.Logger, handler func(raw.Observation)) error {
	bus := s.newBus(log)
	if bus == nil {
		return fmt.Errorf("n2k: serial format %d is not supported", s.format)
	}
	defer func() { _ = bus.Close() }()
	if observationBus, ok := bus.(ObservationBus); ok {
		return observationBus.RunObservations(ctx, handler)
	}
	return bus.Run(ctx, func(frame can.Frame) {
		handler(frameObservation(frame, "serial:"+s.port, s.port, raw.DirectionReceived))
	})
}

func (s *serialSource) newBus(log *slog.Logger) Bus {
	switch s.format {
	case FormatActisense:
		return gateway.NewActisenseSerialBus(log, s.port)
	case FormatActisenseRaw:
		return gateway.NewActisenseRawSerialBus(log, s.port)
	default:
		return nil
	}
}

// SerialDevice describes one OS serial port. USB identity fields are populated
// when the platform exposes them.
type SerialDevice struct {
	Name         string
	IsUSB        bool
	VID          string
	PID          string
	SerialNumber string
	Product      string
}

// SerialDevices enumerates serial ports using the host operating system.
func SerialDevices() ([]SerialDevice, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return nil, fmt.Errorf("n2k: enumerating serial devices: %w", err)
	}
	devices := make([]SerialDevice, 0, len(ports))
	for _, port := range ports {
		if port == nil {
			continue
		}
		devices = append(devices, SerialDevice{
			Name:         port.Name,
			IsUSB:        port.IsUSB,
			VID:          port.VID,
			PID:          port.PID,
			SerialNumber: port.SerialNumber,
			Product:      port.Product,
		})
	}
	return devices, nil
}
