package gateway

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
)

// tcpLink manages one TCP connection to a gateway, shared between a bus's
// Run read loop and concurrent write calls. Run dials; writers block in
// await until the dial completes (writes are issued as soon as a client
// starts, often before the connection is up).
type tcpLink struct {
	log  *slog.Logger
	addr string

	mu      sync.Mutex
	conn    net.Conn
	dialErr error
	closed  bool

	readyOnce sync.Once
	ready     chan struct{}

	// writeMu serializes writes so concurrent callers (address claimer,
	// heartbeat, user writes) cannot interleave partial frames.
	writeMu sync.Mutex
}

func newTCPLink(log *slog.Logger, addr string) *tcpLink {
	return &tcpLink{log: log, addr: addr, ready: make(chan struct{})}
}

// dial opens the connection, transmits the optional greeting, and only then
// releases any writers blocked in await — so a protocol's initialization
// bytes always reach the wire before the first frame. The returned
// connection is closed when ctx is cancelled or the link is closed.
func (l *tcpLink) dial(ctx context.Context, greeting []byte) (net.Conn, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", l.addr)
	if err == nil && len(greeting) > 0 {
		if _, werr := conn.Write(greeting); werr != nil {
			_ = conn.Close()
			conn, err = nil, fmt.Errorf("sending greeting: %w", werr)
		}
	}

	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		if conn != nil {
			_ = conn.Close()
		}
		return nil, errors.New("gateway: bus closed")
	}
	l.conn, l.dialErr = conn, err
	l.mu.Unlock()
	l.signalReady()

	if err != nil {
		return nil, fmt.Errorf("gateway: dialing %s: %w", l.addr, err)
	}

	// Closing the connection on cancellation unblocks any pending Read.
	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()
	return conn, nil
}

// await blocks until dial has run (or the link was closed), then returns the
// live connection.
func (l *tcpLink) await() (net.Conn, error) {
	<-l.ready
	l.mu.Lock()
	defer l.mu.Unlock()
	switch {
	case l.closed:
		return nil, errors.New("gateway: bus closed")
	case l.dialErr != nil:
		return nil, fmt.Errorf("gateway: dialing %s: %w", l.addr, l.dialErr)
	case l.conn == nil:
		return nil, errors.New("gateway: bus not connected")
	}
	return l.conn, nil
}

// write sends one wire-encoded unit (a RAW line or a framed command) whole.
func (l *tcpLink) write(buf []byte) error {
	conn, err := l.await()
	if err != nil {
		return err
	}
	l.writeMu.Lock()
	defer l.writeMu.Unlock()
	n, err := conn.Write(buf)
	if err != nil {
		return err
	}
	if n != len(buf) {
		return fmt.Errorf("gateway: wrote %d of %d bytes", n, len(buf))
	}
	return nil
}

// Close releases the connection and unblocks pending writers. Safe to call
// before, during, or after Run.
func (l *tcpLink) Close() error {
	l.mu.Lock()
	l.closed = true
	conn := l.conn
	l.conn = nil
	l.mu.Unlock()
	l.signalReady()
	if conn != nil {
		return conn.Close()
	}
	return nil
}

func (l *tcpLink) signalReady() {
	l.readyOnce.Do(func() { close(l.ready) })
}

// runResult normalizes a read-loop exit: context cancellation wins over the
// read error it provoked by closing the connection.
func (l *tcpLink) runResult(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}

// YDRawTCPBus is a read/write bus over a TCP gateway speaking the Yacht
// Devices RAW line protocol (YDWG-02 RAW server mode and compatibles). RAW
// mode is frame-level in both directions, so the full client stack —
// address claiming included — works exactly as it does on CAN hardware.
// The gateway echoes transmitted frames back with direction T, so the
// client also observes its own traffic.
type YDRawTCPBus struct {
	link *tcpLink
}

// NewYDRawTCPBus returns an unconnected bus; Run dials.
func NewYDRawTCPBus(log *slog.Logger, addr string) *YDRawTCPBus {
	return &YDRawTCPBus{link: newTCPLink(log, addr)}
}

// Run connects and delivers incoming frames to handler until ctx is
// cancelled or the connection fails.
func (b *YDRawTCPBus) Run(ctx context.Context, handler func(can.Frame)) error {
	conn, err := b.link.dial(ctx, nil)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if frame, ok := ParseYDRaw(scanner.Text()); ok && handler != nil {
			handler(frame)
		}
	}
	return b.link.runResult(ctx, scanner.Err())
}

// WriteFrame transmits one CAN frame as a RAW line.
func (b *YDRawTCPBus) WriteFrame(frame can.Frame) error {
	return b.link.write(FormatYDRawTX(frame))
}

// Close releases the connection.
func (b *YDRawTCPBus) Close() error { return b.link.Close() }

// ActisenseTCPBus is a read/write bus over a TCP connection speaking the
// Actisense binary stream protocol (W2K-1 gateways, or an NGT-1 behind a
// TCP bridge). The protocol is message-level: incoming messages arrive
// assembled and are re-framed for the frame-based decode pipeline, and
// outgoing messages are sent whole via WriteMessage — the gateway performs
// fast-packet fragmentation and stamps its own claimed source address, so
// the client's source address is not authoritative on the wire.
type ActisenseTCPBus struct {
	link *tcpLink
}

// NewActisenseTCPBus returns an unconnected bus; Run dials and sends the
// gateway initialization command.
func NewActisenseTCPBus(log *slog.Logger, addr string) *ActisenseTCPBus {
	return &ActisenseTCPBus{link: newTCPLink(log, addr)}
}

// Run connects, initializes the gateway, and delivers incoming messages to
// handler as re-framed CAN frames until ctx is cancelled or the connection
// fails.
func (b *ActisenseTCPBus) Run(ctx context.Context, handler func(can.Frame)) error {
	conn, err := b.link.dial(ctx, EncodeStartup())
	if err != nil {
		return err
	}

	reader := NewActisenseReader()
	emit := ReframeEmitter(func(f can.Frame) {
		if handler != nil {
			handler(f)
		}
	})
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			reader.Feed(buf[:n], emit)
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				return b.link.runResult(ctx, nil)
			}
			return b.link.runResult(ctx, err)
		}
	}
}

// WriteFrame transmits one CAN frame as a whole message. This is correct
// for single-frame PGNs (address claims, ISO requests, transport-protocol
// frames); fast-packet payloads must go through WriteMessage instead, which
// the client's write path does automatically.
func (b *ActisenseTCPBus) WriteFrame(frame can.Frame) error {
	id := framer.ParseCANID(frame.ID)
	return b.WriteMessage(id.PGN, id.Priority, id.Source, id.Destination, frame.Data[:frame.Length])
}

// WriteMessage transmits an assembled PGN payload (up to 223 bytes). The
// source address is accepted for interface symmetry but not transmitted:
// the gateway substitutes its own claimed address.
func (b *ActisenseTCPBus) WriteMessage(pgnNum uint32, priority, _, destination uint8, payload []byte) error {
	buf, err := EncodeSend(N2KMessage{
		Priority:    priority,
		PGN:         pgnNum,
		Destination: destination,
		Data:        payload,
	})
	if err != nil {
		return err
	}
	return b.link.write(buf)
}

// Close releases the connection.
func (b *ActisenseTCPBus) Close() error { return b.link.Close() }
