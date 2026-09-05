package transport

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// SendBAM transmits a broadcast transfer with the manager's total deadline.
func (m *Manager) SendBAM(pgn uint32, source uint8, payload []byte) error {
	return m.SendBAMContext(context.Background(), pgn, source, payload)
}

// SendBAMContext transmits one broadcast transfer. Cancellation, Reset, and
// Close interrupt pacing and prevent subsequent packets from being sent.
func (m *Manager) SendBAMContext(ctx context.Context, pgn uint32, source uint8, payload []byte) error {
	sess, err := m.beginTransmit(ctx, sessionKey{source: source, destination: BroadcastAddr, pgn: pgn}, payload, stateSendingBAM)
	if err != nil {
		return fmt.Errorf("send BAM: %w", err)
	}
	if err := m.writeSessionFrame(sess, buildCMBAMFrame(sess.totalSize, sess.numFrames, pgn, source)); err != nil {
		return m.finishSession(sess, fmt.Errorf("send CM_BAM: %w", err))
	}
	for i, frame := range buildDTFrameRange(sess.txPayload, 1, int(sess.numFrames), source, BroadcastAddr) {
		if i > 0 {
			if err := m.waitBAM(sess); err != nil {
				return m.finishSession(sess, err)
			}
		}
		if err := m.writeSessionFrame(sess, frame); err != nil {
			return m.finishSession(sess, fmt.Errorf("send DT frame %d: %w", i+1, err))
		}
	}
	return m.finishSession(sess, nil)
}

func (m *Manager) waitBAM(sess *session) error {
	if m.sleep != nil {
		m.sleep(BAMInterFrameDelay)
		return context.Cause(sess.ctx)
	}
	timer := time.NewTimer(BAMInterFrameDelay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return context.Cause(sess.ctx)
	case <-sess.ctx.Done():
		return context.Cause(sess.ctx)
	}
}

// SendRTSCTS transmits an addressed transfer with the manager's total deadline.
func (m *Manager) SendRTSCTS(pgn uint32, source, destination uint8, payload []byte) error {
	return m.SendRTSCTSContext(context.Background(), pgn, source, destination, payload)
}

// SendRTSCTSContext waits for acknowledgment, caller cancellation, transport
// reset, inactivity, or the absolute transfer deadline. Receiver holds never
// extend that absolute deadline.
func (m *Manager) SendRTSCTSContext(ctx context.Context, pgn uint32, source, destination uint8, payload []byte) error {
	if destination == BroadcastAddr {
		return errors.New("send RTS/CTS: broadcast destination is invalid")
	}
	sess, err := m.beginTransmit(ctx, sessionKey{source: source, destination: destination, pgn: pgn}, payload, stateWaitingForCTS)
	if err != nil {
		return fmt.Errorf("send RTS/CTS: %w", err)
	}
	if err := m.writeSessionFrame(sess, buildRTSFrame(sess.totalSize, sess.numFrames, pgn, source, destination)); err != nil {
		return m.finishSession(sess, fmt.Errorf("send RTS: %w", err))
	}
	<-sess.txDone
	return sess.txErr
}

func (m *Manager) beginTransmit(ctx context.Context, key sessionKey, payload []byte, state sessionState) (*session, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := validateTransportPGN(key.pgn); err != nil {
		return nil, err
	}
	if err := validatePayload(payload); err != nil {
		return nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := context.Cause(ctx); err != nil {
		return nil, err
	}
	if m.closed {
		return nil, errors.New("manager closed")
	}
	for activeKey := range m.sessions {
		if activeKey.source == key.source && activeKey.destination == key.destination {
			return nil, fmt.Errorf("transport session already active for source %d destination %d", key.source, key.destination)
		}
	}
	if len(m.sessions) >= m.config.MaxSessions {
		return nil, errors.New("transport session table is full")
	}
	sessionCtx, cancel := context.WithCancelCause(ctx)
	sess := &session{
		key: key, state: state, totalSize: uint16(len(payload)),
		numFrames: uint8((len(payload) + MaxDTDataBytes - 1) / MaxDTDataBytes),
		txPayload: append([]byte(nil), payload...), txDone: make(chan struct{}),
		ctx: sessionCtx, cancel: cancel,
	}
	m.sessions[key] = sess
	sess.deadlineTimer = m.afterFunc(m.config.TransferTimeout, func() {
		_ = m.finishSession(sess, fmt.Errorf("transport transfer deadline: %w", context.DeadlineExceeded))
	})
	sess.stopContext = context.AfterFunc(ctx, func() {
		_ = m.finishSession(sess, context.Cause(ctx))
	})
	if state == stateWaitingForCTS {
		m.armTimerLocked(sess, CTSTimeout, errors.New("CTS timeout waiting for response"))
	}
	return sess, nil
}
