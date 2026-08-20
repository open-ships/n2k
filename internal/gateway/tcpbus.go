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
	"github.com/open-ships/n2k/raw"
)

// errBusClosed is returned to writers and the read loop once the link is
// closed.
var errBusClosed = errors.New("gateway: bus closed")

// tcpLink manages the TCP connection to a gateway, shared between a bus's Run
// read loop and concurrent write calls. Run owns dialing; writers block in
// await until a connection is live (writes are issued as soon as a client
// starts, often before the connection is up).
//
// When a ReconnectPolicy is set, Run re-dials after a dropped connection with
// exponential backoff, and writers issued during an outage block until the
// link reconnects (or is closed) rather than failing. Without a policy, a
// dropped connection ends Run and pending writers observe the error.
type tcpLink struct {
	log       *slog.Logger
	addr      string
	reconnect *ReconnectPolicy // nil = no reconnection

	mu      sync.Mutex
	conn    net.Conn      // current live connection; nil while dialing or disconnected
	dialErr error         // last dial/read error, for the non-reconnect await path
	closed  bool          // link permanently closed by Close
	stopped bool          // run has returned; no further reconnect will occur
	stopErr error         // why run stopped, returned to writers still in await
	ready   chan struct{} // closed (and re-armed) on every state change await must observe
	done    chan struct{} // closed once by Close to interrupt dials, reads, and backoff
	epoch   uint64
	// observer runs before a new connection is published to writers and after
	// a dead connection is removed. It must not perform link I/O.
	observer func(connected bool, epoch uint64)

	// writeMu serializes writes so concurrent callers (address claimer,
	// heartbeat, user writes) cannot interleave partial frames.
	writeMu sync.Mutex
}

func newTCPLink(log *slog.Logger, addr string, reconnect *ReconnectPolicy) *tcpLink {
	return &tcpLink{log: log, addr: addr, reconnect: reconnect, ready: make(chan struct{}), done: make(chan struct{})}
}

func (l *tcpLink) setConnectionObserver(observer func(connected bool, epoch uint64)) {
	l.mu.Lock()
	l.observer = observer
	l.mu.Unlock()
}

// readFunc consumes one live connection until it drops (returning the read
// error, or nil on a clean EOF). A fresh readFunc invocation per connection
// gives protocols with reassembly state (Actisense) a clean slate after a
// reconnect.
type readFunc func(r io.Reader, handler func(raw.Observation)) error

// run dials, sends the optional greeting, and delivers incoming frames via
// read until ctx is cancelled or the link is closed. With a ReconnectPolicy
// it re-dials dropped connections with backoff; without one it returns on the
// first drop or dial failure.
func (l *tcpLink) run(ctx context.Context, greeting []byte, read readFunc, handler func(raw.Observation)) (rerr error) {
	// runCtx is cancelled by the caller's ctx or by Close, so an in-progress
	// dial, read, or backoff all unblock promptly on shutdown either way.
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Once run returns — cancelled or a non-reconnect terminal error — wake any
	// writer still parked in await so it fails instead of blocking forever for a
	// reconnection that will never happen.
	defer func() { l.stop(rerr) }()
	go func() {
		select {
		case <-l.done:
			cancel()
		case <-runCtx.Done():
		}
	}()

	var backoff *Backoff
	if l.reconnect != nil {
		backoff = NewBackoff(*l.reconnect)
	}

	// exit reports how run should return once a step ends: nil for a clean
	// Close, ctx.Err() for caller cancellation, and (nil, false) to continue.
	exit := func() (error, bool) {
		if l.isClosed() {
			return nil, true
		}
		if ctx.Err() != nil {
			return ctx.Err(), true
		}
		return nil, false
	}

	for {
		if err, done := exit(); done {
			return err
		}

		conn, err := l.dialOnce(runCtx, greeting)
		if err != nil {
			if e, done := exit(); done {
				return e
			}
			if backoff == nil {
				return err
			}
			l.log.Warn("gateway: dial failed, retrying", "addr", l.addr, "error", err)
			if !backoff.Wait(runCtx) {
				e, _ := exit()
				return e
			}
			continue
		}
		if backoff != nil {
			backoff.Reset()
		}

		// Close the connection when the read ends, however it ends: on
		// cancellation (ctx or Close) closing it unblocks the pending read, and
		// on a normal drop closing it here prevents leaking the socket (left in
		// CLOSE_WAIT) until finalization. readDone reports a self-terminating read.
		readDone := make(chan struct{})
		go func() {
			select {
			case <-runCtx.Done():
			case <-readDone:
			}
			_ = conn.Close()
		}()
		readErr := read(conn, handler)
		close(readDone)

		// Drop the dead connection so writers block for the next one instead
		// of writing into a closed socket during the backoff window.
		if backoff != nil {
			l.markDisconnected(readErr)
		}
		if e, done := exit(); done {
			return e
		}
		if backoff == nil {
			return readErr
		}

		if readErr != nil {
			l.log.Warn("gateway: connection lost, reconnecting", "addr", l.addr, "error", readErr)
		} else {
			l.log.Info("gateway: connection closed, reconnecting", "addr", l.addr)
		}
		if !backoff.Wait(runCtx) {
			e, _ := exit()
			return e
		}
	}
}

// dialOnce opens one connection, transmits the optional greeting, and
// publishes it so blocked writers unblock — the greeting always reaches the
// wire before the first frame.
func (l *tcpLink) dialOnce(ctx context.Context, greeting []byte) (net.Conn, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", l.addr)
	if err == nil && len(greeting) > 0 {
		if _, werr := conn.Write(greeting); werr != nil {
			_ = conn.Close()
			conn, err = nil, fmt.Errorf("sending greeting: %w", werr)
		}
	}
	if err != nil {
		l.markDisconnected(err)
		return nil, fmt.Errorf("gateway: dialing %s: %w", l.addr, err)
	}
	if !l.markConnected(conn) {
		_ = conn.Close()
		return nil, errBusClosed
	}
	return conn, nil
}

// markConnected publishes conn as the live connection and wakes blocked
// writers. It returns false (leaving conn for the caller to close) if the
// link was closed while dialing.
func (l *tcpLink) markConnected(conn net.Conn) bool {
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		return false
	}
	l.epoch++
	epoch := l.epoch
	observer := l.observer
	l.mu.Unlock()

	// Install the client's reconnect gate before making conn available to any
	// writer parked in await.
	if observer != nil {
		observer(true, epoch)
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed {
		return false
	}
	l.conn = conn
	l.dialErr = nil
	l.broadcast()
	return true
}

// markDisconnected clears the live connection after a drop or dial failure.
// With reconnection enabled, writers keep blocking for the next connection;
// without it, they are woken to observe the error.
func (l *tcpLink) markDisconnected(err error) {
	l.mu.Lock()
	wasConnected := l.conn != nil
	l.conn = nil
	l.dialErr = err
	observer := l.observer
	epoch := l.epoch
	if l.reconnect == nil {
		l.broadcast()
	}
	l.mu.Unlock()
	if wasConnected && observer != nil {
		observer(false, epoch)
	}
}

// stop wakes writers parked in await once run has returned for good, so a
// pending write fails with the reason run ended rather than blocking for a
// reconnection that will never come. A later Close still supersedes this, since
// await checks closed first. No-op once the link is closed or already stopped.
func (l *tcpLink) stop(err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed || l.stopped {
		return
	}
	l.stopped = true
	l.stopErr = err
	l.conn = nil
	l.broadcast()
}

// broadcast wakes every goroutine parked in await and re-arms the gate for the
// next state change. The caller must hold l.mu.
func (l *tcpLink) broadcast() {
	close(l.ready)
	l.ready = make(chan struct{})
}

// maxWriteRetries bounds how many freshly-reconnected connections a single
// write will try before giving up, so a persistently flapping gateway cannot
// stall a writer indefinitely.
const maxWriteRetries = 4

// await blocks until a live connection is available and returns it. A non-nil
// stale is a connection the caller has just found dead: await then waits for a
// different one, so a retry does not spin on the same dead conn before the read
// loop has observed the drop. It returns an error if the link is closed, run
// has stopped, or (without reconnection) the dial failed.
func (l *tcpLink) await(stale net.Conn) (net.Conn, error) {
	for {
		l.mu.Lock()
		switch {
		case l.closed:
			l.mu.Unlock()
			return nil, errBusClosed
		case l.conn != nil && l.conn != stale:
			conn := l.conn
			l.mu.Unlock()
			return conn, nil
		case l.reconnect == nil && l.dialErr != nil:
			err := l.dialErr
			l.mu.Unlock()
			return nil, fmt.Errorf("gateway: dialing %s: %w", l.addr, err)
		case l.stopped:
			err := l.stopErr
			l.mu.Unlock()
			if err == nil {
				err = errBusClosed
			}
			return nil, err
		}
		ready := l.ready
		l.mu.Unlock()
		<-ready
	}
}

// write sends one wire-encoded unit (a RAW line or a framed command) whole.
// If the connection drops before any byte is written and reconnection is
// enabled, write waits for the reconnect and retries — only ever when nothing
// was sent (n == 0), so a retry can never duplicate a partially-written frame.
func (l *tcpLink) write(buf []byte) error {
	l.writeMu.Lock()
	defer l.writeMu.Unlock()
	var stale net.Conn
	for attempt := 0; ; attempt++ {
		conn, err := l.await(stale)
		if err != nil {
			return err
		}
		n, werr := conn.Write(buf)
		if werr == nil {
			if n != len(buf) {
				return fmt.Errorf("gateway: wrote %d of %d bytes", n, len(buf))
			}
			return nil
		}
		// A partial write, a link without reconnection, or an exhausted retry
		// budget is terminal; otherwise wait for a connection other than this
		// dead one and try the whole frame again.
		if n != 0 || l.reconnect == nil || attempt >= maxWriteRetries {
			return werr
		}
		stale = conn
	}
}

// isClosed reports whether Close has been called.
func (l *tcpLink) isClosed() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.closed
}

// Close releases the connection and unblocks pending writers and any reconnect
// backoff. Safe to call before, during, or after Run, and more than once.
func (l *tcpLink) Close() error {
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		return nil
	}
	l.closed = true
	conn := l.conn
	l.conn = nil
	close(l.done)
	l.broadcast()
	l.mu.Unlock()
	if conn != nil {
		return conn.Close()
	}
	return nil
}

// ReadYDRaw delivers Yacht Devices RAW lines from r as CAN frames until EOF or
// a read error. It is the single reader for the RAW protocol, shared by the TCP
// bus and the read-only network sources.
func ReadYDRaw(r io.Reader, handler func(can.Frame)) error {
	return ReadYDRawObservations(r, func(observation raw.Observation) {
		if handler != nil && observation.Frame != nil {
			handler(*observation.Frame)
		}
	})
}

// ReadYDRawObservations preserves direction and gateway-relative timestamps.
func ReadYDRawObservations(r io.Reader, handler func(raw.Observation)) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if observation, ok := ParseYDRawObservation(scanner.Text()); ok && handler != nil {
			handler(observation)
		}
	}
	return scanner.Err()
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

// NewYDRawTCPBus returns an unconnected bus; Run dials. A non-nil reconnect
// policy makes Run re-dial dropped connections.
func NewYDRawTCPBus(log *slog.Logger, addr string, reconnect *ReconnectPolicy) *YDRawTCPBus {
	return &YDRawTCPBus{link: newTCPLink(log, addr, reconnect)}
}

// Run connects and delivers incoming frames to handler until ctx is
// cancelled or (without a reconnect policy) the connection fails.
func (b *YDRawTCPBus) Run(ctx context.Context, handler func(can.Frame)) error {
	return b.RunObservations(ctx, func(observation raw.Observation) {
		if handler != nil && observation.Frame != nil {
			handler(*observation.Frame)
		}
	})
}

// RunObservations preserves gateway transport context.
func (b *YDRawTCPBus) RunObservations(ctx context.Context, handler func(raw.Observation)) error {
	return b.link.run(ctx, nil, ReadYDRawObservations, func(observation raw.Observation) {
		observation.AdapterID = "tcp:" + b.link.addr
		observation.NetworkID = b.link.addr
		if handler != nil {
			handler(observation)
		}
	})
}

// WriteFrame transmits one CAN frame as a RAW line.
func (b *YDRawTCPBus) WriteFrame(frame can.Frame) error {
	return b.link.write(FormatYDRawTX(frame))
}

// SetConnectionObserver installs reconnect lifecycle observation for Client.
func (b *YDRawTCPBus) SetConnectionObserver(observer func(bool, uint64)) {
	b.link.setConnectionObserver(observer)
}

// Close releases the connection.
func (b *YDRawTCPBus) Close() error { return b.link.Close() }
