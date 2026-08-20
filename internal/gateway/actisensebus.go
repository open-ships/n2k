package gateway

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/raw"
)

const (
	actisenseCommandTimeout = 2 * time.Second
	actisenseRestoreTimeout = 500 * time.Millisecond
	maxRememberedTxPGNs     = 1024
)

type actisenseConnection interface {
	io.Reader
	io.Writer
	io.Closer
}

type actisenseOpen func(context.Context) (actisenseConnection, error)

type actisenseEpoch struct {
	connection  actisenseConnection
	session     *actisense.Session
	priorMode   actisense.OperatingMode
	changedMode bool
	txMu        sync.Mutex
	txPGNs      map[uint32]struct{}
	restoreOnce sync.Once
}

// actisenseStreamBus owns reconnect, handshake, and request state for both TCP
// and serial Adapters. A fresh protocol Session is created for every epoch.
type actisenseStreamBus struct {
	log            *slog.Logger
	endpoint       string
	adapterID      string
	open           actisenseOpen
	reconnect      *ReconnectPolicy
	mode           actisense.OperatingMode
	rawCAN         bool
	commandTimeout time.Duration

	mu        sync.Mutex
	epoch     *actisenseEpoch
	closed    bool
	stopped   bool
	stopErr   error
	changed   chan struct{}
	done      chan struct{}
	epochNum  uint64
	observer  func(bool, uint64)
	ready     chan struct{}
	readyOnce sync.Once
}

func newActisenseStreamBus(log *slog.Logger, endpoint, adapterID string, open actisenseOpen, reconnect *ReconnectPolicy, mode actisense.OperatingMode, rawCAN bool) *actisenseStreamBus {
	if log == nil {
		log = slog.Default()
	}
	return &actisenseStreamBus{
		log:            log,
		endpoint:       endpoint,
		adapterID:      adapterID,
		open:           open,
		reconnect:      reconnect,
		mode:           mode,
		rawCAN:         rawCAN,
		commandTimeout: actisenseCommandTimeout,
		changed:        make(chan struct{}),
		done:           make(chan struct{}),
		ready:          make(chan struct{}),
	}
}

func newActisenseTCPStreamBus(log *slog.Logger, addr string, reconnect *ReconnectPolicy, mode actisense.OperatingMode, rawCAN bool) *actisenseStreamBus {
	open := func(ctx context.Context) (actisenseConnection, error) {
		var dialer net.Dialer
		connection, err := dialer.DialContext(ctx, "tcp", addr)
		if err != nil {
			return nil, fmt.Errorf("gateway: dialing %s: %w", addr, err)
		}
		return connection, nil
	}
	return newActisenseStreamBus(log, addr, "tcp:"+addr, open, reconnect, mode, rawCAN)
}

func (b *actisenseStreamBus) Ready() <-chan struct{} { return b.ready }

func (b *actisenseStreamBus) SetConnectionObserver(observer func(bool, uint64)) {
	b.mu.Lock()
	b.observer = observer
	b.mu.Unlock()
}

func (b *actisenseStreamBus) Run(ctx context.Context, handler func(can.Frame)) error {
	return b.RunObservations(ctx, func(observation raw.Observation) {
		if handler != nil && observation.Frame != nil {
			handler(*observation.Frame)
		}
	})
}

func (b *actisenseStreamBus) RunObservations(ctx context.Context, handler func(raw.Observation)) (runErr error) {
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer func() { b.stop(runErr) }()
	go func() {
		select {
		case <-b.done:
			cancel()
		case <-runCtx.Done():
		}
	}()

	var backoff *Backoff
	if b.reconnect != nil {
		backoff = NewBackoff(*b.reconnect)
	}

	for {
		if err, done := b.exitError(ctx); done {
			return err
		}
		connection, err := b.open(runCtx)
		if err != nil {
			if exitErr, done := b.exitError(ctx); done {
				return exitErr
			}
			if backoff == nil {
				return err
			}
			b.log.Warn("actisense: opening Adapter failed, retrying", "endpoint", b.endpoint, "error", err)
			if !backoff.Wait(runCtx) {
				exitErr, _ := b.exitError(ctx)
				return exitErr
			}
			continue
		}

		epochErr, connected := b.runEpoch(runCtx, connection, handler)
		if connected && backoff != nil {
			backoff.Reset()
		}
		if exitErr, done := b.exitError(ctx); done {
			return exitErr
		}
		if backoff == nil {
			return epochErr
		}
		if epochErr != nil {
			b.log.Warn("actisense: connection lost, retrying", "endpoint", b.endpoint, "error", epochErr)
		} else {
			b.log.Info("actisense: connection closed, retrying", "endpoint", b.endpoint)
		}
		if !backoff.Wait(runCtx) {
			exitErr, _ := b.exitError(ctx)
			return exitErr
		}
	}
}

func (b *actisenseStreamBus) runEpoch(ctx context.Context, connection actisenseConnection, handler func(raw.Observation)) (error, bool) {
	adapter := &actisenseAdapter{}
	var published atomic.Bool
	epoch := &actisenseEpoch{connection: connection, txPGNs: make(map[uint32]struct{})}
	epoch.session = actisense.NewSession(actisense.SessionConfig{
		Write: func(buf []byte) error { return writeActisenseUnit(connection, buf) },
		OnDatagram: func(datagram actisense.Datagram) {
			adapter.observeDatagram(datagram, func(observation raw.Observation) {
				observation.AdapterID = b.adapterID
				observation.NetworkID = b.endpoint
				if handler != nil {
					handler(observation)
				}
			})
		},
		OnDiagnostic: func(diagnostic actisense.Diagnostic) {
			if diagnostic.Kind == actisense.DiagnosticStartup && published.Load() {
				// A reset reverts volatile mode/list state. Closing this epoch forces
				// reconnect and a complete acknowledged handshake.
				_ = connection.Close()
			}
		},
		OnDecodeError: func(decodeErr actisense.DecodeError) {
			adapter.observeDecodeError(decodeErr, func(observation raw.Observation) {
				observation.AdapterID = b.adapterID
				observation.NetworkID = b.endpoint
				if handler != nil {
					handler(observation)
				}
			})
		},
	})
	readErr := make(chan error, 1)
	go func() { readErr <- epoch.session.Run(connection) }()

	handshakeCtx, cancel := context.WithTimeout(ctx, b.commandTimeout)
	priorMode, err := epoch.session.GetOperatingMode(handshakeCtx)
	if err == nil {
		epoch.priorMode = priorMode
		if priorMode != b.mode {
			err = epoch.session.SetOperatingMode(handshakeCtx, b.mode)
			epoch.changedMode = err == nil
		}
	}
	cancel()
	if err != nil {
		_ = connection.Close()
		epoch.session.Close(err)
		<-readErr
		return fmt.Errorf("actisense: acknowledged mode-%d setup failed: %w", b.mode, err), false
	}

	if !b.publishEpoch(epoch) {
		_ = connection.Close()
		epoch.session.Close(errBusClosed)
		<-readErr
		return errBusClosed, false
	}
	published.Store(true)

	var readResult error
	select {
	case readResult = <-readErr:
	case <-ctx.Done():
		b.restoreEpoch(epoch)
		_ = connection.Close()
		readResult = <-readErr
	}
	published.Store(false)
	b.clearEpoch(epoch)
	epoch.session.Close(readResult)
	_ = connection.Close()
	return readResult, true
}

func writeActisenseUnit(writer io.Writer, buf []byte) error {
	written := 0
	for written < len(buf) {
		n, err := writer.Write(buf[written:])
		written += n
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
	}
	return nil
}

func (b *actisenseStreamBus) publishEpoch(epoch *actisenseEpoch) bool {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return false
	}
	b.epochNum++
	epochNumber := b.epochNum
	observer := b.observer
	b.mu.Unlock()

	if observer != nil {
		observer(true, epochNumber)
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return false
	}
	b.epoch = epoch
	b.broadcastLocked()
	b.readyOnce.Do(func() { close(b.ready) })
	return true
}

func (b *actisenseStreamBus) clearEpoch(epoch *actisenseEpoch) {
	b.mu.Lock()
	if b.epoch != epoch {
		b.mu.Unlock()
		return
	}
	b.epoch = nil
	observer := b.observer
	epochNumber := b.epochNum
	b.broadcastLocked()
	b.mu.Unlock()
	if observer != nil {
		observer(false, epochNumber)
	}
}

func (b *actisenseStreamBus) awaitEpoch() (*actisenseEpoch, error) {
	for {
		b.mu.Lock()
		switch {
		case b.closed:
			b.mu.Unlock()
			return nil, errBusClosed
		case b.epoch != nil:
			epoch := b.epoch
			b.mu.Unlock()
			return epoch, nil
		case b.stopped:
			err := b.stopErr
			b.mu.Unlock()
			if err == nil {
				err = errBusClosed
			}
			return nil, err
		}
		changed := b.changed
		b.mu.Unlock()
		<-changed
	}
}

func (b *actisenseStreamBus) writeFrame(frame can.Frame) error {
	if frame.Length > 8 {
		return fmt.Errorf("actisense: invalid CAN frame length %d", frame.Length)
	}
	epoch, err := b.awaitEpoch()
	if err != nil {
		return err
	}
	if b.rawCAN {
		buf, encodeErr := actisense.EncodeCANFrame(frame, actisense.DirectionTransmitted, 0, 0)
		if encodeErr != nil {
			return encodeErr
		}
		return epoch.session.Write(buf)
	}
	id := framer.ParseCANID(frame.ID)
	return b.writeMessage(id.PGN, id.Priority, id.Destination, frame.Data[:frame.Length])
}

func (b *actisenseStreamBus) writeMessage(pgnNumber uint32, priority, destination uint8, payload []byte) error {
	epoch, err := b.awaitEpoch()
	if err != nil {
		return err
	}
	if b.rawCAN {
		return errors.New("actisense: assembled-message writes are unavailable in raw CAN mode")
	}
	if err := b.ensureTxPGN(epoch, pgnNumber); err != nil {
		return err
	}
	buf, err := actisense.EncodeMessage94(actisense.Message{
		Priority: priority, PGN: pgnNumber, Destination: destination, Data: payload,
	})
	if err != nil {
		return err
	}
	return epoch.session.Write(buf)
}

func (b *actisenseStreamBus) ensureTxPGN(epoch *actisenseEpoch, pgnNumber uint32) error {
	epoch.txMu.Lock()
	defer epoch.txMu.Unlock()
	if _, ok := epoch.txPGNs[pgnNumber]; ok {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), b.commandTimeout)
	defer cancel()
	state, err := epoch.session.GetTxPGN(ctx, pgnNumber)
	if err != nil {
		return fmt.Errorf("actisense: query Tx PGN %d: %w", pgnNumber, err)
	}
	if state.Enabled != 1 {
		if err := epoch.session.SetTxPGN(ctx, pgnNumber, true); err != nil {
			return fmt.Errorf("actisense: enable Tx PGN %d: %w", pgnNumber, err)
		}
		if err := epoch.session.ActivatePGNLists(ctx); err != nil {
			return fmt.Errorf("actisense: activate Tx PGN %d: %w", pgnNumber, err)
		}
	}
	if len(epoch.txPGNs) < maxRememberedTxPGNs {
		epoch.txPGNs[pgnNumber] = struct{}{}
	}
	return nil
}

func (b *actisenseStreamBus) restoreEpoch(epoch *actisenseEpoch) {
	if epoch == nil || !epoch.changedMode {
		return
	}
	epoch.restoreOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), actisenseRestoreTimeout)
		defer cancel()
		if err := epoch.session.SetOperatingMode(ctx, epoch.priorMode); err != nil {
			b.log.Warn("actisense: could not restore prior operating mode", "endpoint", b.endpoint, "mode", epoch.priorMode, "error", err)
		}
	})
}

func (b *actisenseStreamBus) exitError(ctx context.Context) (error, bool) {
	b.mu.Lock()
	closed := b.closed
	b.mu.Unlock()
	if closed {
		return nil, true
	}
	if ctx.Err() != nil {
		return ctx.Err(), true
	}
	return nil, false
}

func (b *actisenseStreamBus) broadcastLocked() {
	close(b.changed)
	b.changed = make(chan struct{})
}

func (b *actisenseStreamBus) stop(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed || b.stopped {
		return
	}
	b.stopped = true
	b.stopErr = err
	b.broadcastLocked()
}

func (b *actisenseStreamBus) Close() error {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return nil
	}
	b.closed = true
	epoch := b.epoch
	close(b.done)
	b.broadcastLocked()
	b.mu.Unlock()
	if epoch == nil {
		return nil
	}
	b.restoreEpoch(epoch)
	return epoch.connection.Close()
}

// ActisenseTCPBus is the v1-compatible, gateway-owned BST-93/94 message
// session. It is retained for compatibility; raw CAN mode is the authoritative
// Client Bus.
type ActisenseTCPBus struct{ stream *actisenseStreamBus }

func NewActisenseTCPBus(log *slog.Logger, addr string, reconnect *ReconnectPolicy) *ActisenseTCPBus {
	return &ActisenseTCPBus{stream: newActisenseTCPStreamBus(log, addr, reconnect, actisense.ModeTransferReceiveAll, false)}
}

func (b *ActisenseTCPBus) Run(ctx context.Context, handler func(can.Frame)) error {
	return b.stream.Run(ctx, handler)
}
func (b *ActisenseTCPBus) RunObservations(ctx context.Context, handler func(raw.Observation)) error {
	return b.stream.RunObservations(ctx, handler)
}
func (b *ActisenseTCPBus) WriteFrame(frame can.Frame) error { return b.stream.writeFrame(frame) }
func (b *ActisenseTCPBus) WriteMessage(pgnNumber uint32, priority, _ uint8, destination uint8, payload []byte) error {
	return b.stream.writeMessage(pgnNumber, priority, destination, payload)
}
func (b *ActisenseTCPBus) SetConnectionObserver(observer func(bool, uint64)) {
	b.stream.SetConnectionObserver(observer)
}
func (b *ActisenseTCPBus) Ready() <-chan struct{} { return b.stream.Ready() }
func (b *ActisenseTCPBus) Close() error           { return b.stream.Close() }

// ActisenseRawTCPBus is a source-authoritative BST-95 raw CAN Bus.
type ActisenseRawTCPBus struct{ stream *actisenseStreamBus }

func NewActisenseRawTCPBus(log *slog.Logger, addr string, reconnect *ReconnectPolicy) *ActisenseRawTCPBus {
	return &ActisenseRawTCPBus{stream: newActisenseTCPStreamBus(log, addr, reconnect, actisense.ModeCANPacket, true)}
}

func (b *ActisenseRawTCPBus) Run(ctx context.Context, handler func(can.Frame)) error {
	return b.stream.Run(ctx, handler)
}
func (b *ActisenseRawTCPBus) RunObservations(ctx context.Context, handler func(raw.Observation)) error {
	return b.stream.RunObservations(ctx, handler)
}
func (b *ActisenseRawTCPBus) WriteFrame(frame can.Frame) error { return b.stream.writeFrame(frame) }
func (b *ActisenseRawTCPBus) SetConnectionObserver(observer func(bool, uint64)) {
	b.stream.SetConnectionObserver(observer)
}
func (b *ActisenseRawTCPBus) Ready() <-chan struct{} { return b.stream.Ready() }
func (b *ActisenseRawTCPBus) Close() error           { return b.stream.Close() }
