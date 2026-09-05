package transport

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
)

// ManagerConfig holds the dependencies for a transport protocol Manager.
type ManagerConfig struct {
	// WriteFrame sends a CAN frame onto the bus.
	WriteFrame func(can.Frame) error
	// WriteFrameContext, when provided, replaces WriteFrame and must respect
	// cancellation while waiting for write admission or transport I/O.
	WriteFrameContext func(context.Context, can.Frame) error
	// TransferTimeout bounds the entire transmit operation, including repeated
	// receiver holds. Zero selects DefaultTransferTimeout.
	TransferTimeout time.Duration
	// MaxSessions bounds combined receive and transmit state. Zero selects 512.
	MaxSessions int

	// LocalAddress returns the client's currently claimed address. When set,
	// addressed TP traffic for other nodes is ignored. A function is used
	// because address claiming may move the client at runtime.
	LocalAddress func() uint8

	// OnComplete is called when a multi-frame message has been fully reassembled.
	// It receives the transported PGN, source address, destination address, and
	// the assembled payload.
	OnComplete func(pgn uint32, source uint8, destination uint8, data []byte)
	// OnCompleteInfo is the source-aware completion callback. The MessageInfo
	// comes from the connection-management frame that began the transfer, with
	// PGN replaced by the transported PGN.
	OnCompleteInfo func(info pgn.MessageInfo, data []byte)

	// Logger for transport protocol events. If nil, a no-op logger is used.
	Logger *slog.Logger

	// AfterFunc schedules a callback after d; nil means time.AfterFunc.
	// Tests inject a fake to drive timeouts deterministically.
	AfterFunc func(d time.Duration, f func()) Timer

	// Sleep overrides BAM pacing for deterministic tests. Production callers
	// should leave it nil to use interruptible context-aware pacing.
	Sleep func(d time.Duration)
}

// Timer is the subset of *time.Timer the manager uses.
type Timer interface{ Stop() bool }

// sessionKey uniquely identifies an active transport protocol session.
type sessionKey struct {
	source      uint8
	destination uint8
	pgn         uint32
}

// sessionState tracks the phase of a transport protocol session.
type sessionState int

const (
	stateReceivingDT       sessionState = iota // accumulating DT frames
	stateWaitingForCTS                         // transmitter waiting for CTS
	stateSendingDT                             // transmitter sending DT frames
	stateSendingBAM                            // transmitter sending broadcast DT frames
	stateCompletingReceive                     // receiver sending its final acknowledgment
)

// session holds all state for a single transport protocol exchange.
type session struct {
	key       sessionKey
	state     sessionState
	totalSize uint16 // expected total payload bytes
	numFrames uint8  // expected total DT frame count
	received  int    // number of DT frames received so far
	data      []byte // reassembly buffer
	info      pgn.MessageInfo

	// RTS/CTS fields
	maxPerCTS uint8 // max DT frames per CTS cycle (from RTS byte 4)

	// Timeout management
	timer           Timer
	timerGeneration uint64
	deadlineTimer   Timer
	stopContext     func() bool
	ctx             context.Context
	cancel          context.CancelCauseFunc

	// Transmit-side fields for RTS/CTS
	txPayload   []byte        // full payload to transmit
	txSent      [256]bool     // successfully sent DT sequence numbers
	txSentCount int           // number of unique DT packets sent
	txDone      chan struct{} // closed when transmit completes
	txErr       error         // set if transmit fails
}

// Manager orchestrates ISO 11783 transport protocol sessions, handling both
// BAM and RTS/CTS modes for receive and transmit.
type Manager struct {
	config    ManagerConfig
	sessions  map[sessionKey]*session
	mu        sync.Mutex
	closed    bool
	logger    *slog.Logger
	afterFunc func(time.Duration, func()) Timer
	sleep     func(time.Duration)
}

// NewManager creates a new transport protocol Manager.
func NewManager(cfg ManagerConfig) *Manager {
	if cfg.WriteFrame == nil {
		cfg.WriteFrame = func(can.Frame) error {
			return fmt.Errorf("transport WriteFrame is nil")
		}
	}
	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}
	afterFunc := cfg.AfterFunc
	if afterFunc == nil {
		afterFunc = func(d time.Duration, f func()) Timer { return time.AfterFunc(d, f) }
	}
	if cfg.TransferTimeout <= 0 {
		cfg.TransferTimeout = DefaultTransferTimeout
	}
	if cfg.MaxSessions <= 0 {
		cfg.MaxSessions = 512
	}
	return &Manager{
		config:    cfg,
		sessions:  make(map[sessionKey]*session),
		logger:    logger,
		afterFunc: afterFunc,
		sleep:     cfg.Sleep,
	}
}

// HandleFrame routes an incoming CAN frame to the appropriate transport protocol
// handler based on PGN. Non-TP frames are silently ignored.
func (m *Manager) HandleFrame(frame can.Frame) {
	c := framer.ParseCANID(frame.ID)
	now := time.Now()
	priority := c.Priority
	info := pgn.MessageInfo{Timestamp: now, ReceivedAt: now, Priority: &priority, PGN: c.PGN, SourceId: c.Source}
	if c.Destination != BroadcastAddr {
		destination := c.Destination
		info.TargetId = &destination
	}
	m.HandleFrameWithInfo(frame, info)
}

// HandleFrameWithInfo routes a transport frame while retaining its Adapter,
// network, timing, and direction context through reassembly.
func (m *Manager) HandleFrameWithInfo(frame can.Frame, info pgn.MessageInfo) {
	c := framer.ParseCANID(frame.ID)
	switch c.PGN {
	case PGNCM:
		m.handleCM(frame, c.Source, c.Destination, info)
	case PGNDT:
		m.handleDT(frame, c.Source, c.Destination)
	}
}

// handleCM dispatches a Connection Management frame based on the control byte.
func (m *Manager) handleCM(frame can.Frame, source uint8, destination uint8, info pgn.MessageInfo) {
	if frame.Length != 8 {
		return
	}
	controlByte := frame.Data[0]

	switch controlByte {
	case ControlBAM:
		if destination != BroadcastAddr {
			return
		}
		m.handleBAMReceive(frame, source, info)
	case ControlRTS:
		if !m.isLocalDestination(destination) {
			return
		}
		m.handleRTSReceive(frame, source, destination, info)
	case ControlCTS:
		if !m.isLocalDestination(destination) {
			return
		}
		m.handleCTSReceive(frame, source, destination)
	case ControlEndOfMsgAck:
		if !m.isLocalDestination(destination) {
			return
		}
		m.handleEndOfMsgAckReceive(frame, source, destination)
	case ControlAbort:
		if !m.isLocalDestination(destination) {
			return
		}
		m.handleAbortReceive(frame, source, destination)
	default:
		m.logger.Warn("unknown TP.CM control byte", "controlByte", controlByte, "source", source)
	}
}

// handleDT delivers a Data Transfer frame to the matching active session.
func (m *Manager) handleDT(frame can.Frame, source uint8, destination uint8) {
	if frame.Length != 8 {
		return
	}
	if destination != BroadcastAddr && !m.isLocalDestination(destination) {
		return
	}

	m.mu.Lock()

	// Find the session. For BAM, destination is 255. For RTS/CTS, it's the actual target.
	sess := m.findSession(source, destination)
	if sess == nil {
		m.mu.Unlock()
		return
	}

	seqNum := frame.Data[0]
	expectedSeq := uint8(sess.received + 1)
	if seqNum != expectedSeq {
		m.logger.Warn("unexpected DT sequence number",
			"expected", expectedSeq, "got", seqNum,
			"source", source, "pgn", sess.key.pgn)
		m.finishSessionLocked(sess, errors.New("unexpected DT sequence number"))
		m.mu.Unlock()
		return
	}

	// Copy data bytes into the reassembly buffer.
	offset := int(seqNum-1) * MaxDTDataBytes
	remaining := int(sess.totalSize) - offset
	if remaining <= 0 || offset < 0 || offset >= len(sess.data) {
		m.logger.Warn("DT frame exceeds announced payload", "sequence", seqNum, "size", sess.totalSize)
		m.finishSessionLocked(sess, errors.New("DT frame exceeds announced payload"))
		m.mu.Unlock()
		return
	}
	n := MaxDTDataBytes
	if remaining < n {
		n = remaining
	}
	copy(sess.data[offset:offset+n], frame.Data[1:1+n])
	sess.received++

	// Reset the DT timeout timer.
	m.stopTimerLocked(sess)

	if sess.received < int(sess.numFrames) {
		// Addressed transfers are flow-controlled in blocks. Once the block
		// granted by the previous CTS has arrived, grant the next one. BAM has
		// no CTS exchanges and continues directly to the timeout re-arm below.
		var nextCTS *can.Frame
		if sess.key.destination != BroadcastAddr && sess.received%int(sess.maxPerCTS) == 0 {
			remainingFrames := int(sess.numFrames) - sess.received
			requested := int(sess.maxPerCTS)
			if remainingFrames < requested {
				requested = remainingFrames
			}
			f := buildCTSFrame(
				uint8(requested),
				uint8(sess.received+1),
				sess.key.pgn,
				sess.key.destination,
				sess.key.source,
			)
			nextCTS = &f
		}

		// More frames expected; set a DT timeout.
		m.armTimerLocked(sess, DTTimeout, errors.New("DT timeout"))
		key := sess.key
		m.mu.Unlock()
		if nextCTS != nil {
			if err := m.writeSessionFrame(sess, *nextCTS); err != nil {
				_ = m.finishSession(sess, err)
				m.logger.Warn("failed to send next CTS", "error", err, "pgn", key.pgn)
			}
		}
		return
	}

	// All frames received. Collect everything needed for delivery under the
	// lock, then release it so a slow OnComplete consumer cannot stall other
	// Manager calls.
	key := sess.key
	data := make([]byte, sess.totalSize)
	copy(data, sess.data[:sess.totalSize])
	dst := sess.key.destination
	info := sess.info

	// For RTS/CTS receive: build the EndOfMsgAck to send after unlocking.
	var ackFrame *can.Frame
	if dst != BroadcastAddr {
		f := buildEndOfMsgAckFrame(sess)
		ackFrame = &f
	}

	sess.state = stateCompletingReceive
	m.mu.Unlock()

	if ackFrame != nil {
		if err := m.writeSessionFrame(sess, *ackFrame); err != nil {
			_ = m.finishSession(sess, err)
			m.logger.Warn("failed to send EndOfMsgAck", "error", err)
			return
		}
	}
	m.mu.Lock()
	completed := m.finishSessionLocked(sess, nil)
	m.mu.Unlock()
	if !completed {
		return
	}
	if m.config.OnComplete != nil {
		m.config.OnComplete(key.pgn, key.source, dst, data)
	}
	if m.config.OnCompleteInfo != nil {
		m.config.OnCompleteInfo(info, data)
	}
}

// findSession looks up only receiving sessions by source and destination.
func (m *Manager) findSession(source uint8, destination uint8) *session {
	// Try each PGN — we iterate sessions since PGN is part of the key.
	for k, s := range m.sessions {
		if k.source == source && k.destination == destination && s.state == stateReceivingDT {
			return s
		}
	}
	return nil
}

func (m *Manager) isLocalDestination(destination uint8) bool {
	return m.config.LocalAddress == nil || destination == m.config.LocalAddress()
}

func validateAnnouncement(totalSize uint16, numFrames uint8) error {
	if totalSize == 0 || numFrames == 0 {
		return fmt.Errorf("zero size or frame count")
	}
	if totalSize > MaxPayloadBytes {
		return fmt.Errorf("payload size %d exceeds maximum %d", totalSize, MaxPayloadBytes)
	}
	want := (int(totalSize) + MaxDTDataBytes - 1) / MaxDTDataBytes
	if int(numFrames) != want {
		return fmt.Errorf("frame count %d does not match size %d (want %d)", numFrames, totalSize, want)
	}
	return nil
}

func validatePayload(payload []byte) error {
	if len(payload) == 0 {
		return fmt.Errorf("empty transport payload")
	}
	if len(payload) > MaxPayloadBytes {
		return fmt.Errorf("payload size %d exceeds maximum %d", len(payload), MaxPayloadBytes)
	}
	return nil
}

func validateTransportPGN(pgn uint32) error {
	if pgn > MaxPGN {
		return fmt.Errorf("PGN %d exceeds the 18-bit range", pgn)
	}
	if pgn&0x20000 != 0 {
		return fmt.Errorf("PGN %d sets the reserved CAN-ID bit", pgn)
	}
	if ((pgn>>8)&0xFF) < 240 && pgn&0xFF != 0 {
		return fmt.Errorf("PDU1 PGN %d must have a zero group-extension byte", pgn)
	}
	return nil
}

func (m *Manager) admitReceiveLocked(key sessionKey) bool {
	if m.closed {
		return false
	}
	for activeKey, sess := range m.sessions {
		if activeKey.source == key.source && activeKey.destination == key.destination && sess.txDone != nil {
			// A transport echo must not replace an outgoing transfer.
			return false
		}
	}
	m.removeReceiveSessions(key.source, key.destination)
	if len(m.sessions) >= m.config.MaxSessions {
		m.logger.Warn("transport receive session table is full", "source", key.source, "destination", key.destination)
		return false
	}
	return true
}

// removeReceiveSessions replaces the one receive transfer for these addresses.
func (m *Manager) removeReceiveSessions(source, destination uint8) {
	for key, sess := range m.sessions {
		if key.source == source && key.destination == destination && sess.txDone == nil {
			m.finishSessionLocked(sess, errors.New("receive transfer superseded"))
		}
	}
}

// stopTimerLocked invalidates callbacks which have already started and cannot
// be stopped by Timer.Stop. The caller holds m.mu.
func (m *Manager) stopTimerLocked(sess *session) {
	sess.timerGeneration++
	if sess.timer != nil {
		sess.timer.Stop()
		sess.timer = nil
	}
}

func (m *Manager) armTimerLocked(sess *session, duration time.Duration, err error) {
	m.stopTimerLocked(sess)
	generation := sess.timerGeneration
	sess.timer = m.afterFunc(duration, func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if m.sessions[sess.key] != sess || sess.timerGeneration != generation {
			return
		}
		m.finishSessionLocked(sess, err)
	})
}

// finishSessionLocked owns every terminal transition and completion signal.
// Pointer identity protects a replacement transfer with an identical wire key.
func (m *Manager) finishSessionLocked(sess *session, err error) bool {
	if m.sessions[sess.key] != sess {
		return false
	}
	delete(m.sessions, sess.key)
	m.stopTimerLocked(sess)
	if sess.deadlineTimer != nil {
		sess.deadlineTimer.Stop()
	}
	if sess.stopContext != nil {
		sess.stopContext()
	}
	sess.txErr = err
	if sess.cancel != nil {
		sess.cancel(err)
	}
	if sess.txDone != nil {
		close(sess.txDone)
	}
	return true
}

func (m *Manager) finishSession(sess *session, err error) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.finishSessionLocked(sess, err)
	return sess.txErr
}

// writeSessionFrame rejects obsolete work both before and after transport I/O.
// The context-aware writer is responsible for cancellation during the I/O.
func (m *Manager) writeSessionFrame(sess *session, frame can.Frame) error {
	m.mu.Lock()
	active := m.sessions[sess.key] == sess
	err := sess.txErr
	m.mu.Unlock()
	if !active {
		if err == nil {
			err = errors.New("transport session is no longer active")
		}
		return err
	}
	if err := context.Cause(sess.ctx); err != nil {
		return err
	}
	var writeErr error
	if m.config.WriteFrameContext != nil {
		writeErr = m.config.WriteFrameContext(sess.ctx, frame)
	} else {
		writeErr = m.config.WriteFrame(frame)
	}
	m.mu.Lock()
	active = m.sessions[sess.key] == sess
	err = sess.txErr
	m.mu.Unlock()
	if !active {
		return err
	}
	if err := context.Cause(sess.ctx); err != nil {
		return err
	}
	return writeErr
}

// handleCTSReceive sends a validated receiver-granted packet window.
func (m *Manager) handleCTSReceive(frame can.Frame, source uint8, destination uint8) {
	numFrames := frame.Data[1]
	nextSeqNum := frame.Data[2]
	key := sessionKey{source: destination, destination: source, pgn: extractPGN(frame.Data)}

	m.mu.Lock()
	sess := m.sessions[key]
	if sess == nil || sess.state != stateWaitingForCTS {
		m.mu.Unlock()
		return
	}
	if numFrames == 0 {
		m.armTimerLocked(sess, CTSTimeout, errors.New("CTS timeout while receiver paused transfer"))
		m.mu.Unlock()
		return
	}
	lastSeqNum := int(nextSeqNum) + int(numFrames) - 1
	if nextSeqNum == 0 || int(nextSeqNum) > int(sess.numFrames) || lastSeqNum > int(sess.numFrames) {
		m.finishSessionLocked(sess, fmt.Errorf("invalid CTS range: first=%d count=%d total=%d", nextSeqNum, numFrames, sess.numFrames))
		m.mu.Unlock()
		return
	}
	m.stopTimerLocked(sess)
	sess.state = stateSendingDT
	payload := sess.txPayload
	m.mu.Unlock()

	for _, dtFrame := range buildDTFrameRange(payload, nextSeqNum, int(numFrames), destination, source) {
		if err := m.writeSessionFrame(sess, dtFrame); err != nil {
			_ = m.finishSession(sess, fmt.Errorf("send DT frame %d: %w", dtFrame.Data[0], err))
			return
		}
		m.mu.Lock()
		if m.sessions[key] != sess {
			m.mu.Unlock()
			return
		}
		seq := dtFrame.Data[0]
		if !sess.txSent[seq] {
			sess.txSent[seq] = true
			sess.txSentCount++
		}
		m.mu.Unlock()
	}

	m.mu.Lock()
	if m.sessions[key] == sess {
		sess.state = stateWaitingForCTS
		m.armTimerLocked(sess, CTSTimeout, errors.New("CTS timeout after sending DT frames"))
	}
	m.mu.Unlock()
}

// handleEndOfMsgAckReceive validates and completes a transmit-side session.
func (m *Manager) handleEndOfMsgAckReceive(frame can.Frame, source uint8, destination uint8) {
	key := sessionKey{source: destination, destination: source, pgn: extractPGN(frame.Data)}
	m.mu.Lock()
	defer m.mu.Unlock()
	sess := m.sessions[key]
	if sess == nil || sess.state != stateWaitingForCTS {
		return
	}
	totalSize := uint16(frame.Data[1]) | uint16(frame.Data[2])<<8
	numFrames := frame.Data[3]
	if totalSize != sess.totalSize || numFrames != sess.numFrames || sess.txSentCount != int(sess.numFrames) {
		m.finishSessionLocked(sess, fmt.Errorf(
			"invalid EndOfMsgAck: size=%d frames=%d packets sent=%d; want size=%d frames=%d",
			totalSize, numFrames, sess.txSentCount, sess.totalSize, sess.numFrames))
		return
	}
	m.finishSessionLocked(sess, nil)
}

func (m *Manager) handleAbortReceive(frame can.Frame, source uint8, destination uint8) {
	pgnNumber := extractPGN(frame.Data)
	err := fmt.Errorf("remote abort, reason=%d", frame.Data[1])
	keys := [2]sessionKey{
		{source: destination, destination: source, pgn: pgnNumber},
		{source: source, destination: destination, pgn: pgnNumber},
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, key := range keys {
		if sess := m.sessions[key]; sess != nil {
			m.finishSessionLocked(sess, err)
		}
	}
}

// buildEndOfMsgAckFrame constructs an EndOfMsgAck frame for a completed
// receive session. The caller writes it to the bus outside the manager lock.
func buildEndOfMsgAckFrame(sess *session) can.Frame {
	var data [8]uint8
	data[0] = ControlEndOfMsgAck
	data[1] = uint8(sess.totalSize)
	data[2] = uint8(sess.totalSize >> 8)
	data[3] = sess.numFrames
	data[4] = 0xFF
	encodePGN(data[5:8], sess.key.pgn)

	canID := framer.BuildCANID(PGNCM, TPPriority, sess.key.destination, sess.key.source)
	return can.Frame{
		ID:     canID,
		Length: 8,
		Data:   data,
	}
}

// Reset invalidates every active transfer without closing the manager. The
// supplied cause reaches blocked transmitters and context-aware frame writes.
// Call this whenever the network connection or local claimed address changes.
func (m *Manager) Reset(err error) {
	if err == nil {
		err = errors.New("transport manager reset")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, sess := range m.sessions {
		m.finishSessionLocked(sess, err)
	}
}

// Close stops active transfers and permanently prevents new sessions.
func (m *Manager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	for _, sess := range m.sessions {
		m.finishSessionLocked(sess, errors.New("manager closed"))
	}
}

// extractPGN reads a 3-byte little-endian PGN from CM frame bytes 5-7.
func extractPGN(data [8]uint8) uint32 {
	return uint32(data[5]) | uint32(data[6])<<8 | uint32(data[7])<<16
}

// encodePGN writes a PGN as 3 little-endian bytes.
func encodePGN(dst []byte, pgn uint32) {
	dst[0] = uint8(pgn)
	dst[1] = uint8(pgn >> 8)
	dst[2] = uint8(pgn >> 16)
}

// buildCMBAMFrame constructs a CM_BAM announcement frame.
func buildCMBAMFrame(totalSize uint16, numFrames uint8, pgn uint32, source uint8) can.Frame {
	var data [8]uint8
	data[0] = ControlBAM
	data[1] = uint8(totalSize)
	data[2] = uint8(totalSize >> 8)
	data[3] = numFrames
	data[4] = 0xFF
	encodePGN(data[5:8], pgn)

	canID := framer.BuildCANID(PGNCM, TPPriority, source, BroadcastAddr)
	return can.Frame{
		ID:     canID,
		Length: 8,
		Data:   data,
	}
}

// buildRTSFrame constructs an RTS frame for connection-managed transfers.
func buildRTSFrame(totalSize uint16, numFrames uint8, pgn uint32, source uint8, destination uint8) can.Frame {
	var data [8]uint8
	data[0] = ControlRTS
	data[1] = uint8(totalSize)
	data[2] = uint8(totalSize >> 8)
	data[3] = numFrames
	data[4] = 0xFF // unlimited frames per CTS
	encodePGN(data[5:8], pgn)

	canID := framer.BuildCANID(PGNCM, TPPriority, source, destination)
	return can.Frame{
		ID:     canID,
		Length: 8,
		Data:   data,
	}
}

// buildDTFrameRange builds DT frames [startSeq, startSeq+count) for payload.
// Frames past the end of payload are omitted.
func buildDTFrameRange(payload []byte, startSeq uint8, count int, source uint8, destination uint8) []can.Frame {
	frames := make([]can.Frame, 0, count)

	for i := 0; i < count; i++ {
		seqNum := startSeq + uint8(i)
		offset := int(seqNum-1) * MaxDTDataBytes
		if offset >= len(payload) {
			break
		}

		end := offset + MaxDTDataBytes
		if end > len(payload) {
			end = len(payload)
		}

		var data [7]byte
		for j := range data {
			data[j] = 0xFF
		}
		copy(data[:], payload[offset:end])

		frames = append(frames, makeDTFrame(seqNum, data, source, destination))
	}

	return frames
}

// makeDTFrame constructs a single DT frame.
func makeDTFrame(seqNum uint8, data [7]byte, source uint8, destination uint8) can.Frame {
	var frameData [8]uint8
	frameData[0] = seqNum
	copy(frameData[1:], data[:])

	canID := framer.BuildCANID(PGNDT, TPPriority, source, destination)
	return can.Frame{
		ID:     canID,
		Length: 8,
		Data:   frameData,
	}
}
