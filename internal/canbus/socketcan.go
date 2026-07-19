package canbus

import (
	"context"
	"log/slog"
	"sync"

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
	newBus        func(string) (*can.Bus, error)
}

// socketCANChannel represents a single SocketCAN-based canbus channel for sending/receiving CAN frames.
// It uses the brutella/can library for CAN socket I/O on an already-configured Linux CAN interface.
type socketCANChannel struct {
	// options holds the user-provided configuration (interface name).
	options socketCANChannelOptions

	// bus is the underlying brutella/can bus object that manages the CAN socket connection.
	// It handles subscribing to incoming frames and publishing outgoing frames.
	bus       *can.Bus
	mu        sync.Mutex
	ready     chan struct{}
	readyOnce sync.Once

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
//   - options: required configuration including the interface name (the message
//     handler is passed separately to Run)
//
// Returns a *socketCANChannel (catalog type, not Interface) because SocketCAN-specific
// callers may need access to SocketCAN-specific functionality.
func newSocketCANChannel(log *slog.Logger, options socketCANChannelOptions) *socketCANChannel {
	c := socketCANChannel{
		options: options,
		log:     log,
		ready:   make(chan struct{}),
	}

	return &c
}

// Run opens the SocketCAN interface and starts listening for CAN frames,
// delivering each one to handler. The CAN interface must already be
// configured and up (e.g., via `ip link set can0 up type can bitrate 250000`).
//
// This method blocks until an error occurs or the connection is closed.
func (c *socketCANChannel) Run(ctx context.Context, handler func(can.Frame)) error {
	// Open a CAN socket using the brutella/can library. This creates a raw CAN socket
	// bound to the specified network interface and provides a higher-level pub/sub API.
	newBus := c.options.newBus
	if newBus == nil {
		newBus = can.NewBusForInterfaceWithName
	}
	bus, err := newBus(c.options.InterfaceName)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.bus = bus
	c.mu.Unlock()
	defer func() {
		c.mu.Lock()
		if c.bus == bus {
			c.bus = nil
			c.mu.Unlock()
			_ = bus.Disconnect()
			return
		}
		c.mu.Unlock()
	}()
	c.readyOnce.Do(func() { close(c.ready) })

	// Wrap the caller's handler so a nil handler becomes a no-op, then
	// subscribe it to receive all incoming CAN frames.
	var frameHandler can.HandlerFunc = func(frame can.Frame) {
		if handler != nil {
			handler(frame)
		}
	}
	bus.Subscribe(can.NewHandler(frameHandler))

	watchDone := make(chan struct{})
	defer close(watchDone)
	go func() {
		select {
		case <-ctx.Done():
			_ = c.Close()
		case <-watchDone:
		}
	}()

	c.log.Info("Opened SocketCAN and listening", "interfaceName", c.options.InterfaceName)

	// ConnectAndPublish opens the socket and enters a blocking read loop.
	// Each received CAN frame is published to all subscribed handlers.
	// This call blocks until the socket is closed or an error occurs.
	err = bus.ConnectAndPublish()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}

// Close shuts down the SocketCAN channel by disconnecting the underlying CAN
// socket. It is safe to call Close() if the bus was never opened (bus is nil).
func (c *socketCANChannel) Close() error {
	c.mu.Lock()
	bus := c.bus
	c.bus = nil
	c.mu.Unlock()
	if bus == nil {
		return nil
	}

	if err := bus.Disconnect(); err != nil {
		return errors.Wrap(err, "close underlying bus connection")
	}

	return nil
}

// WriteFrame sends a single CAN frame out on the SocketCAN bus.
// The brutella/can library handles encoding the frame into the Linux SocketCAN wire format
// and writing it to the raw CAN socket. Returns an error if the bus is not yet open.
func (c *socketCANChannel) WriteFrame(frame can.Frame) error {
	if frame.Length > 8 {
		return errors.Errorf("socketCAN: invalid classical CAN payload length %d", frame.Length)
	}
	c.mu.Lock()
	bus := c.bus
	c.mu.Unlock()
	if bus == nil {
		return errors.New("socketCAN: bus not open (interface not available or Run not called)")
	}
	return bus.Publish(frame)
}

// Ready closes once Run has opened the SocketCAN device for writes.
func (c *socketCANChannel) Ready() <-chan struct{} { return c.ready }

// NewSocketCAN creates a SocketCAN channel for the given Linux CAN interface name.
// The channel is not opened, and no handler is registered, until Run() is called.
func NewSocketCAN(log *slog.Logger, iface string) *socketCANChannel {
	return newSocketCANChannel(log, socketCANChannelOptions{
		InterfaceName: iface,
	})
}

// RunSocketCAN creates a SocketCAN channel for the given interface and runs it,
// calling handler for each received CAN frame. The interface must already be configured and up.
// Blocks until error or context done.
func RunSocketCAN(ctx context.Context, log *slog.Logger, iface string, handler func(can.Frame)) error {
	ch := newSocketCANChannel(log, socketCANChannelOptions{
		InterfaceName: iface,
	})

	return ch.Run(ctx, handler)
}
