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
	"github.com/open-ships/n2k/raw"
)

const (
	actisenseCommandTimeout = 5 * time.Second
)

// ActisenseGatewayMetrics combines protocol-session counters across reconnect
// epochs with gateway lifecycle events.
type ActisenseGatewayMetrics struct {
	ConnectionEpochs uint64
	Reconnects       uint64
	GatewayResets    uint64
	Protocol         actisense.SessionMetrics
}

type actisenseConnection interface {
	io.Reader
	io.Writer
	io.Closer
}

// ActisenseConnection is the public package's byte-stream seam for a custom
// gateway transport.
type ActisenseConnection interface {
	io.Reader
	io.Writer
	io.Closer
}

// ActisenseOpen opens one connection epoch for a custom byte transport.
type ActisenseOpen func(context.Context) (ActisenseConnection, error)

type actisenseOpen func(context.Context) (actisenseConnection, error)

// ActisenseModeSetupError reports a failed acknowledged operating-mode setup.
// The public package translates it into ActisenseModeError.
type ActisenseModeSetupError struct {
	RequestedMode actisense.OperatingMode
	Err           error
}

func (e *ActisenseModeSetupError) Error() string {
	return fmt.Sprintf("actisense: acknowledged mode-%d setup failed: %v", e.RequestedMode, e.Err)
}

func (e *ActisenseModeSetupError) Unwrap() error { return e.Err }

type actisenseEpoch struct {
	connection  actisenseConnection
	session     *actisense.Session
	priorMode   actisense.OperatingMode
	changedMode bool
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

	mu                 sync.Mutex
	epoch              *actisenseEpoch
	closed             bool
	stopped            bool
	stopErr            error
	changed            chan struct{}
	done               chan struct{}
	epochNum           uint64
	observer           func(bool, uint64)
	diagnosticObserver func(actisense.Diagnostic)
	wireObserver       func(actisense.WireDirection, time.Time, []byte)
	ready              chan struct{}
	readyOnce          sync.Once
	completedMetrics   actisense.SessionMetrics
	gatewayResets      uint64
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

func (b *actisenseStreamBus) SetDiagnosticObserver(observer func(actisense.Diagnostic)) {
	b.mu.Lock()
	b.diagnosticObserver = observer
	b.mu.Unlock()
}

func (b *actisenseStreamBus) SetWireObserver(observer func(actisense.WireDirection, time.Time, []byte)) {
	b.mu.Lock()
	b.wireObserver = observer
	b.mu.Unlock()
}

func (b *actisenseStreamBus) SetCommandTimeout(timeout time.Duration) error {
	if timeout <= 0 {
		return errors.New("actisense: command timeout must be positive")
	}
	b.mu.Lock()
	b.commandTimeout = timeout
	b.mu.Unlock()
	return nil
}

func (b *actisenseStreamBus) Metrics() ActisenseGatewayMetrics {
	b.mu.Lock()
	epochs := b.epochNum
	resets := b.gatewayResets
	protocol := cloneActisenseSessionMetrics(b.completedMetrics)
	if b.epoch != nil {
		mergeActisenseSessionMetrics(&protocol, b.epoch.session.Metrics())
	}
	b.mu.Unlock()
	reconnects := uint64(0)
	if epochs > 1 {
		reconnects = epochs - 1
	}
	return ActisenseGatewayMetrics{ConnectionEpochs: epochs, Reconnects: reconnects, GatewayResets: resets, Protocol: protocol}
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
	asciiLines := &actisenseCANASCIILines{}
	var published atomic.Bool
	epoch := &actisenseEpoch{connection: connection}
	emitObservation := func(observation raw.Observation) {
		observation.AdapterID = b.adapterID
		observation.NetworkID = b.endpoint
		if handler != nil {
			handler(observation)
		}
	}
	epoch.session = actisense.NewSession(actisense.SessionConfig{
		Write: func(buf []byte) error { return writeActisenseUnit(connection, buf) },
		OnDatagram: func(datagram actisense.Datagram) {
			adapter.observeDatagram(datagram, emitObservation)
		},
		OnDiagnostic: func(diagnostic actisense.Diagnostic) {
			b.mu.Lock()
			observer := b.diagnosticObserver
			b.mu.Unlock()
			if observer != nil {
				observer(diagnostic)
			}
			if diagnostic.Kind == actisense.DiagnosticStartup && published.Load() {
				b.mu.Lock()
				b.gatewayResets++
				b.mu.Unlock()
				// A reset reverts volatile mode/list state. Closing this epoch forces
				// reconnect and a complete acknowledged handshake.
				_ = connection.Close()
			}
		},
		OnDecodeError: func(decodeErr actisense.DecodeError) {
			adapter.observeDecodeError(decodeErr, emitObservation)
		},
		OnWireBytes: func(direction actisense.WireDirection, timestamp time.Time, data []byte) {
			b.mu.Lock()
			observer := b.wireObserver
			b.mu.Unlock()
			if observer != nil {
				observer(direction, timestamp, data)
			}
		},
		OnUnframed: func(data []byte) {
			if b.mode == actisense.ModeCANPacketASCII {
				asciiLines.feed(data, emitObservation)
			}
		},
	})
	readErr := make(chan error, 1)
	go func() { readErr <- epoch.session.Run(connection) }()

	b.mu.Lock()
	commandTimeout := b.commandTimeout
	b.mu.Unlock()
	handshakeCtx, cancel := context.WithTimeout(ctx, commandTimeout)
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
		return &ActisenseModeSetupError{RequestedMode: b.mode, Err: err}, false
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
	b.recordEpochMetrics(epoch.session.Metrics())
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
	return b.awaitEpochContext(context.Background())
}

func (b *actisenseStreamBus) awaitEpochContext(ctx context.Context) (*actisenseEpoch, error) {
	if ctx == nil {
		ctx = context.Background()
	}
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
		select {
		case <-changed:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (b *actisenseStreamBus) writeFrame(frame can.Frame) error {
	if frame.Length > 8 {
		return fmt.Errorf("actisense: invalid CAN frame length %d", frame.Length)
	}
	if !b.rawCAN {
		return errors.New("actisense: a gateway-owned message session cannot write source-authoritative CAN frames")
	}
	epoch, err := b.awaitEpoch()
	if err != nil {
		return err
	}
	if b.mode == actisense.ModeCANPacketASCII {
		buf, encodeErr := FormatActisenseCANASCII(frame, raw.DirectionTransmitted, 0)
		if encodeErr != nil {
			return encodeErr
		}
		return epoch.session.Write(buf)
	}
	buf, encodeErr := actisense.EncodeCANFrame(frame, actisense.DirectionTransmitted, 0, 0)
	if encodeErr != nil {
		return encodeErr
	}
	return epoch.session.Write(buf)
}

func (b *actisenseStreamBus) writeMessage(pgnNumber uint32, priority, destination uint8, payload []byte) error {
	return b.writeMessageContext(context.Background(), pgnNumber, priority, destination, payload)
}

func (b *actisenseStreamBus) writeMessageContext(ctx context.Context, pgnNumber uint32, priority, destination uint8, payload []byte) error {
	epoch, err := b.awaitEpochContext(ctx)
	if err != nil {
		return err
	}
	if b.rawCAN {
		return errors.New("actisense: assembled-message writes are unavailable in raw CAN mode")
	}
	buf, err := actisense.EncodeMessage94(actisense.Message{
		Priority: priority, PGN: pgnNumber, Destination: destination, Data: payload,
	})
	if err != nil {
		return err
	}
	return epoch.session.Write(buf)
}

func (b *actisenseStreamBus) Request(ctx context.Context, command byte, data []byte) (actisense.BEMResponse, error) {
	epoch, err := b.awaitEpochContext(ctx)
	if err != nil {
		return actisense.BEMResponse{}, err
	}
	return epoch.session.Request(ctx, command, data)
}

func (b *actisenseStreamBus) RequestMulti(ctx context.Context, command byte, data []byte, inactivity time.Duration, complete func([]actisense.BEMResponse) (bool, error)) ([]actisense.BEMResponse, error) {
	epoch, err := b.awaitEpochContext(ctx)
	if err != nil {
		return nil, err
	}
	return epoch.session.RequestMulti(ctx, command, data, inactivity, complete)
}

func (b *actisenseStreamBus) restoreEpoch(epoch *actisenseEpoch) {
	if epoch == nil || !epoch.changedMode {
		return
	}
	epoch.restoreOnce.Do(func() {
		b.mu.Lock()
		commandTimeout := b.commandTimeout
		b.mu.Unlock()
		ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
		defer cancel()
		if err := epoch.session.SetOperatingMode(ctx, epoch.priorMode); err != nil {
			b.log.Warn("actisense: could not restore prior operating mode", "endpoint", b.endpoint, "mode", epoch.priorMode, "error", err)
		}
	})
}

func (b *actisenseStreamBus) recordEpochMetrics(metrics actisense.SessionMetrics) {
	b.mu.Lock()
	mergeActisenseSessionMetrics(&b.completedMetrics, metrics)
	b.mu.Unlock()
}

func cloneActisenseSessionMetrics(metrics actisense.SessionMetrics) actisense.SessionMetrics {
	copy := metrics
	copy.BSTFrames = make(map[byte]uint64, len(metrics.BSTFrames))
	for id, count := range metrics.BSTFrames {
		copy.BSTFrames[id] = count
	}
	return copy
}

func mergeActisenseSessionMetrics(target *actisense.SessionMetrics, value actisense.SessionMetrics) {
	priorCompleted := target.BEMCompleted
	priorTotal := uint64(target.BEMLatencyAverage) * priorCompleted
	valueTotal := uint64(value.BEMLatencyAverage) * value.BEMCompleted
	target.TransportReadCalls += value.TransportReadCalls
	target.TransportReadBytes += value.TransportReadBytes
	target.TransportReadErrors += value.TransportReadErrors
	target.TransportWriteCalls += value.TransportWriteCalls
	target.TransportWriteBytes += value.TransportWriteBytes
	target.TransportWriteErrors += value.TransportWriteErrors
	target.Datagrams += value.Datagrams
	target.UnframedBytes += value.UnframedBytes
	target.FramingErrors += value.FramingErrors
	target.ChecksumErrors += value.ChecksumErrors
	target.LengthErrors += value.LengthErrors
	target.OversizeErrors += value.OversizeErrors
	target.BEMRequests += value.BEMRequests
	target.BEMResponses += value.BEMResponses
	target.BEMCompleted += value.BEMCompleted
	target.BEMCorrelationMisses += value.BEMCorrelationMisses
	target.BEMDuplicateRequests += value.BEMDuplicateRequests
	target.BEMTimeouts += value.BEMTimeouts
	target.BEMDeviceErrors += value.BEMDeviceErrors
	target.BEMNegativeAcks += value.BEMNegativeAcks
	target.BEMResponseTrainOverflows += value.BEMResponseTrainOverflows
	target.BEMInFlight = value.BEMInFlight
	target.BEMMaxInFlight = max(target.BEMMaxInFlight, value.BEMMaxInFlight)
	if target.BEMLatencyMinimum == 0 || value.BEMLatencyMinimum != 0 && value.BEMLatencyMinimum < target.BEMLatencyMinimum {
		target.BEMLatencyMinimum = value.BEMLatencyMinimum
	}
	target.BEMLatencyMaximum = max(target.BEMLatencyMaximum, value.BEMLatencyMaximum)
	if target.BEMCompleted != 0 {
		target.BEMLatencyAverage = time.Duration((priorTotal + valueTotal) / target.BEMCompleted)
	}
	if target.BSTFrames == nil {
		target.BSTFrames = make(map[byte]uint64)
	}
	for id, count := range value.BSTFrames {
		target.BSTFrames[id] += count
	}
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
func (b *ActisenseRawTCPBus) Metrics() ActisenseGatewayMetrics {
	return b.stream.Metrics()
}

// ActisenseCANASCIITCPBus is a source-authoritative mode-6 CAN ASCII Bus.
// Binary BEM control replies are demultiplexed from ASCII frame lines.
type ActisenseCANASCIITCPBus struct{ stream *actisenseStreamBus }

func NewActisenseCANASCIITCPBus(log *slog.Logger, addr string, reconnect *ReconnectPolicy) *ActisenseCANASCIITCPBus {
	return &ActisenseCANASCIITCPBus{stream: newActisenseTCPStreamBus(log, addr, reconnect, actisense.ModeCANPacketASCII, true)}
}

func (b *ActisenseCANASCIITCPBus) Run(ctx context.Context, handler func(can.Frame)) error {
	return b.stream.Run(ctx, handler)
}
func (b *ActisenseCANASCIITCPBus) RunObservations(ctx context.Context, handler func(raw.Observation)) error {
	return b.stream.RunObservations(ctx, handler)
}
func (b *ActisenseCANASCIITCPBus) WriteFrame(frame can.Frame) error {
	return b.stream.writeFrame(frame)
}
func (b *ActisenseCANASCIITCPBus) SetConnectionObserver(observer func(bool, uint64)) {
	b.stream.SetConnectionObserver(observer)
}
func (b *ActisenseCANASCIITCPBus) Ready() <-chan struct{} { return b.stream.Ready() }
func (b *ActisenseCANASCIITCPBus) Close() error           { return b.stream.Close() }
func (b *ActisenseCANASCIITCPBus) Metrics() ActisenseGatewayMetrics {
	return b.stream.Metrics()
}
