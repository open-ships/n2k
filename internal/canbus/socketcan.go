package canbus

import (
	"context"
	"log/slog"

	"github.com/brutella/can"
	"github.com/pkg/errors"
)

// socketCANChannelOptions contains the configuration required to open and operate a SocketCAN channel.
//
// SocketCAN is the native Linux kernel CAN bus interface. It represents CAN controllers as
// regular network interfaces (e.g., "can0", "can1") that can be configured using standard
// networking tools (ip link) and accessed through socket APIs.
type socketCANChannelOptions struct {
	// InterfaceName is the Linux network interface name for the CAN controller, e.g. "can0".
	// This corresponds to the interface shown by `ip link show` and is typically assigned
	// by the kernel when a CAN controller driver (such as MCP2515 over SPI) is loaded.
	InterfaceName string `json:"interfaceName"`

	// MessageHandler is the callback function invoked for each CAN frame received on the bus.
	// The handler receives a can.Frame containing the CAN ID and up to 8 bytes of payload.
	MessageHandler can.HandlerFunc `json:"messageHandler"`
}

// socketCANChannel represents a single SocketCAN-based canbus channel for sending/receiving CAN frames.
// It uses the brutella/can library for CAN socket I/O on an already-configured Linux CAN interface.
type socketCANChannel struct {
	// options holds the user-provided configuration (interface name, handler, etc.)
	options socketCANChannelOptions

	// bus is the underlying brutella/can bus object that manages the CAN socket connection.
	// It handles subscribing to incoming frames and publishing outgoing frames.
	bus *can.Bus

	// busHandler wraps the user's MessageHandler callback into a can.Handler interface
	// so it can be registered with the brutella/can bus subscription system.
	busHandler can.Handler

	// log is the structured logger for diagnostic output about interface state changes and errors.
	log *slog.Logger
}

// newSocketCANChannel creates and returns a new socketCANChannel configured with the given options.
// The channel is not opened until Run() is called.
//
// The CAN interface must already be configured and up before calling Run().
//
// Parameters:
//   - log: structured logger for diagnostic output
//   - options: required configuration including interface name and message handler
//
// Returns a *socketCANChannel (concrete type, not Interface) because SocketCAN-specific
// callers may need access to SocketCAN-specific functionality.
func newSocketCANChannel(log *slog.Logger, options socketCANChannelOptions) *socketCANChannel {
	c := socketCANChannel{
		options: options,
		log:     log,
	}

	return &c
}

// Run opens the SocketCAN interface and starts listening for CAN frames.
// The CAN interface must already be configured and up (e.g., via `ip link set can0 up type can bitrate 250000`).
//
// This method blocks until an error occurs or the connection is closed.
func (c *socketCANChannel) Run(ctx context.Context) error {
	// Open a CAN socket using the brutella/can library. This creates a raw CAN socket
	// bound to the specified network interface and provides a higher-level pub/sub API.
	bus, err := can.NewBusForInterfaceWithName(c.options.InterfaceName)
	if err != nil {
		return err
	}

	c.bus = bus

	// Wrap the user's handler callback and subscribe it to receive all incoming CAN frames.
	c.busHandler = can.NewHandler(c.options.MessageHandler)
	c.bus.Subscribe(c.busHandler)

	c.log.Info("Opened SocketCAN and listening", "interfaceName", c.options.InterfaceName)

	// ConnectAndPublish opens the socket and enters a blocking read loop.
	// Each received CAN frame is published to all subscribed handlers.
	// This call blocks until the socket is closed or an error occurs.
	return bus.ConnectAndPublish()
}

// Close shuts down the SocketCAN channel by unsubscribing the frame handler and disconnecting
// the underlying CAN socket. It is safe to call Close() if the bus was never opened (bus is nil).
func (c *socketCANChannel) Close() error {
	if c.bus == nil {
		return nil
	}

	// Unsubscribe first to stop receiving frames, then disconnect the socket.
	c.bus.Unsubscribe(c.busHandler)
	if err := c.bus.Disconnect(); err != nil {
		return errors.Wrap(err, "close underlying bus connection")
	}

	return nil
}

// WriteFrame sends a single CAN frame out on the SocketCAN bus.
// The brutella/can library handles encoding the frame into the Linux SocketCAN wire format
// and writing it to the raw CAN socket. Returns an error if the bus is not yet open.
func (c *socketCANChannel) WriteFrame(frame can.Frame) error {
	if c.bus == nil {
		return errors.New("socketCAN: bus not open (interface not available or Run not called)")
	}
	return c.bus.Publish(frame)
}

// NewSocketCAN creates a SocketCAN Interface for the given Linux CAN interface name.
// The handler callback receives each incoming CAN frame. The interface is not opened
// until Run() is called.
func NewSocketCAN(log *slog.Logger, iface string, handler func(can.Frame)) Interface {
	var frameHandler can.HandlerFunc = func(frame can.Frame) {
		if handler != nil {
			handler(frame)
		}
	}
	return newSocketCANChannel(log, socketCANChannelOptions{
		InterfaceName:  iface,
		MessageHandler: frameHandler,
	})
}

// RunSocketCAN creates a SocketCAN channel for the given interface and runs it,
// calling handler for each received CAN frame. The interface must already be configured and up.
// Blocks until error or context done.
func RunSocketCAN(ctx context.Context, log *slog.Logger, iface string, handler func(can.Frame)) error {
	var frameHandler can.HandlerFunc = func(frame can.Frame) {
		handler(frame)
	}

	ch := newSocketCANChannel(log, socketCANChannelOptions{
		InterfaceName:  iface,
		MessageHandler: frameHandler,
	})

	return ch.Run(ctx)
}
