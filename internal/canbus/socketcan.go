package canbus

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/brutella/can"
)

// A classical Linux struct can_frame contains a native-endian can_id, a
// length byte, three padding/extension bytes, and eight data bytes.
const socketCANFrameSize = 16

type socketCANChannel struct {
	iface string
	log   *slog.Logger
	open  func(string) (io.ReadWriteCloser, error)

	mu     sync.Mutex
	conn   io.ReadWriteCloser
	closed bool
	ready  chan struct{}
}

// Run reads extended classical data frames from an already configured interface.
func (c *socketCANChannel) Run(ctx context.Context, handler func(can.Frame)) error {
	conn, err := c.open(c.iface)
	if err != nil {
		return err
	}
	c.mu.Lock()
	if c.closed || c.conn != nil {
		c.mu.Unlock()
		_ = conn.Close()
		return errors.New("socketCAN: channel already started or closed")
	}
	c.conn = conn
	close(c.ready)
	c.mu.Unlock()
	defer func() { _ = c.Close() }()
	stop := context.AfterFunc(ctx, func() { _ = c.Close() })
	defer stop()

	c.log.Info("Opened SocketCAN and listening", "interfaceName", c.iface)
	var record [socketCANFrameSize]byte
	for {
		n, err := conn.Read(record[:])
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err != nil {
			return err
		}
		if n != len(record) {
			return fmt.Errorf("socketCAN: got %d bytes in a classical CAN record", n)
		}
		id := binary.NativeEndian.Uint32(record[:4])
		if id & ^uint32(can.MaskIDEff) != can.MaskEff || record[4] > 8 {
			continue
		}
		if handler != nil {
			handler(can.Frame{ID: id & can.MaskIDEff, Length: record[4], Data: [8]byte(record[8:])})
		}
	}
}

// Close interrupts I/O and prevents a concurrent open from publishing a socket.
func (c *socketCANChannel) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	return err
}

func (c *socketCANChannel) WriteFrame(frame can.Frame) error {
	if frame.Length > 8 || frame.ID > can.MaskIDEff {
		return fmt.Errorf("socketCAN: invalid classical CAN frame ID %#x length %d", frame.ID, frame.Length)
	}
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()
	if conn == nil {
		return errors.New("socketCAN: bus not open")
	}

	var record [socketCANFrameSize]byte
	binary.NativeEndian.PutUint32(record[:4], frame.ID|can.MaskEff)
	record[4] = frame.Length
	copy(record[8:], frame.Data[:frame.Length])
	n, err := conn.Write(record[:])
	if err == nil && n != len(record) {
		return io.ErrShortWrite
	}
	return err
}

// Ready closes once Run has opened the SocketCAN device for writes.
func (c *socketCANChannel) Ready() <-chan struct{} { return c.ready }

// NewSocketCAN creates a channel for the named Linux CAN interface.
func NewSocketCAN(log *slog.Logger, iface string) *socketCANChannel {
	return &socketCANChannel{iface: iface, log: log, open: openSocketCAN, ready: make(chan struct{})}
}

// RunSocketCAN reads frames until ctx ends or the interface fails.
func RunSocketCAN(ctx context.Context, log *slog.Logger, iface string, handler func(can.Frame)) error {
	return NewSocketCAN(log, iface).Run(ctx, handler)
}
