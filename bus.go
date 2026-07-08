package n2k

import (
	"context"

	"github.com/brutella/can"
)

// Bus is a physical or virtual CAN bus. External callers can implement Bus
// to inject fake hardware for testing, or to adapt transports this library
// does not ship. Pass an implementation via WithBus.
type Bus interface {
	// Run opens the bus and delivers every incoming frame to handler until
	// ctx is cancelled or an unrecoverable error occurs.
	Run(ctx context.Context, handler func(can.Frame)) error
	// WriteFrame sends one frame.
	WriteFrame(frame can.Frame) error
	// Close releases resources. Safe to call even if Run was never called.
	Close() error
}
