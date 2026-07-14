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

// MessageWriter is optionally implemented by Bus implementations that
// transmit whole assembled PGN messages rather than raw CAN frames —
// message-oriented gateways such as Actisense-format streams, where the
// gateway performs fast-packet fragmentation itself. When a client's bus
// implements MessageWriter, writes that fit in one message (payloads up to
// 223 bytes) bypass CAN framing and use WriteMessage; larger ISO-TP
// transfers and protocol frames (address claims, ISO requests) still go
// frame-by-frame through WriteFrame.
//
// Implementations may not be able to honor the source address (an Actisense
// gateway stamps its own claimed address); it is provided so implementations
// that can control it (custom transports) have it.
type MessageWriter interface {
	WriteMessage(pgnNum uint32, priority, source, destination uint8, payload []byte) error
}
