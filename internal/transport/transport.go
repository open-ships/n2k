package transport

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
)

// ManagerConfig holds the dependencies for a transport protocol Manager.
type ManagerConfig struct {
	// WriteFrame sends a CAN frame onto the bus.
	WriteFrame func(can.Frame) error

	// OnComplete is called when a multi-frame message has been fully reassembled.
	// It receives the transported PGN, source address, destination address, and
	// the assembled payload.
	OnComplete func(pgn uint32, source uint8, destination uint8, data []byte)

	// Logger for transport protocol events. If nil, a no-op logger is used.
	Logger *slog.Logger
}

// sessionKey uniquely identifies an active transport protocol session.
type sessionKey struct {
	source      uint8
	destination uint8
	pgn         uint32
}

// sessionState tracks the phase of a transport protocol session.
type sessionState int

const (
	stateReceivingDT   sessionState = iota // accumulating DT frames
	stateWaitingForCTS                     // transmitter waiting for CTS
	stateSendingDT                         // transmitter sending DT frames
)

// session holds all state for a single transport protocol exchange.
type session struct {
	key       sessionKey
	state     sessionState
	totalSize uint16 // expected total payload bytes
	numFrames uint8  // expected total DT frame count
	received  int    // number of DT frames received so far
	data      []byte // reassembly buffer

	// RTS/CTS fields
	maxPerCTS uint8 // max DT frames per CTS cycle (from RTS byte 4)

	// Timeout management
	timer *time.Timer

	// Transmit-side fields for RTS/CTS
	txPayload []byte        // full payload to transmit
	txOffset  int           // current offset into txPayload
	txSeqNum  uint8         // next DT sequence number to send (1-based)
	txDone    chan struct{} // closed when transmit completes
	txErr     error         // set if transmit fails
}

// Manager orchestrates ISO 11783 transport protocol sessions, handling both
// BAM and RTS/CTS modes for receive and transmit.
type Manager struct {
	config   ManagerConfig
	sessions map[sessionKey]*session
	mu       sync.Mutex
	closed   bool
	logger   *slog.Logger
}

// NewManager creates a new transport protocol Manager.
func NewManager(cfg ManagerConfig) *Manager {
	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}
	return &Manager{
		config:   cfg,
		sessions: make(map[sessionKey]*session),
		logger:   logger,
	}
}

// HandleFrame routes an incoming CAN frame to the appropriate transport protocol
// handler based on PGN. Non-TP frames are silently ignored.
func (m *Manager) HandleFrame(frame can.Frame) {
	pgn, source, destination := parseCANID(frame.ID)

	switch pgn {
	case PGNCM:
		m.handleCM(frame, source, destination)
	case PGNDT:
		m.handleDT(frame, source, destination)
	}
}

// handleCM dispatches a Connection Management frame based on the control byte.
func (m *Manager) handleCM(frame can.Frame, source uint8, destination uint8) {
	if frame.Length < 8 {
		return
	}
	controlByte := frame.Data[0]

	switch controlByte {
	case ControlBAM:
		m.handleBAMReceive(frame, source)
	case ControlRTS:
		m.handleRTSReceive(frame, source, destination)
	case ControlCTS:
		m.handleCTSReceive(frame, source, destination)
	case ControlEndOfMsgAck:
		m.handleEndOfMsgAckReceive(frame, source, destination)
	case ControlAbort:
		m.handleAbortReceive(frame, source, destination)
	default:
		m.logger.Warn("unknown TP.CM control byte", "controlByte", controlByte, "source", source)
	}
}

// handleDT delivers a Data Transfer frame to the matching active session.
func (m *Manager) handleDT(frame can.Frame, source uint8, destination uint8) {
	if frame.Length < 8 {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Find the session. For BAM, destination is 255. For RTS/CTS, it's the actual target.
	sess := m.findSession(source, destination)
	if sess == nil {
		return
	}

	seqNum := frame.Data[0]
	expectedSeq := uint8(sess.received + 1)
	if seqNum != expectedSeq {
		m.logger.Warn("unexpected DT sequence number",
			"expected", expectedSeq, "got", seqNum,
			"source", source, "pgn", sess.key.pgn)
		m.removeSession(sess.key)
		return
	}

	// Copy data bytes into the reassembly buffer.
	offset := int(seqNum-1) * MaxDTDataBytes
	remaining := int(sess.totalSize) - offset
	n := MaxDTDataBytes
	if remaining < n {
		n = remaining
	}
	copy(sess.data[offset:offset+n], frame.Data[1:1+n])
	sess.received++

	// Reset the DT timeout timer.
	if sess.timer != nil {
		sess.timer.Stop()
	}

	if sess.received >= int(sess.numFrames) {
		// All frames received.
		key := sess.key
		data := make([]byte, sess.totalSize)
		copy(data, sess.data[:sess.totalSize])
		dst := sess.key.destination

		// For RTS/CTS receive: send EndOfMsgAck.
		if dst != BroadcastAddr {
			m.sendEndOfMsgAck(sess)
		}

		m.removeSession(key)
		// Deliver outside the lock would be nice, but we need to ensure ordering.
		// The callback should be fast.
		if m.config.OnComplete != nil {
			m.config.OnComplete(key.pgn, key.source, dst, data)
		}
	} else {
		// More frames expected; set a DT timeout.
		sess.timer = time.AfterFunc(DTTimeout, func() {
			m.mu.Lock()
			defer m.mu.Unlock()
			m.logger.Warn("DT timeout", "source", source, "pgn", sess.key.pgn,
				"received", sess.received, "expected", sess.numFrames)
			m.removeSession(sess.key)
		})
	}
}

// findSession looks up a session by source and destination. It tries an exact match
// first, then falls back to broadcast destination for BAM sessions.
func (m *Manager) findSession(source uint8, destination uint8) *session {
	// Try each PGN — we iterate sessions since PGN is part of the key.
	for k, s := range m.sessions {
		if k.source == source && k.destination == destination {
			return s
		}
	}
	// For DT frames with broadcast destination, try BAM sessions.
	if destination == BroadcastAddr {
		for k, s := range m.sessions {
			if k.source == source && k.destination == BroadcastAddr {
				return s
			}
		}
	}
	return nil
}

// removeSession stops timers and deletes a session from the map.
func (m *Manager) removeSession(key sessionKey) {
	if s, ok := m.sessions[key]; ok {
		if s.timer != nil {
			s.timer.Stop()
		}
		delete(m.sessions, key)
	}
}

// SendBAM transmits a multi-frame message using BAM (Broadcast Announce Message).
// This blocks the caller for the duration of the transmission due to the required
// inter-frame delays.
func (m *Manager) SendBAM(pgn uint32, source uint8, payload []byte) error {
	totalSize := uint16(len(payload))
	numFrames := uint8((len(payload) + MaxDTDataBytes - 1) / MaxDTDataBytes)

	// Build and send CM_BAM announcement.
	cmFrame := buildCMBAMFrame(totalSize, numFrames, pgn, source)
	if err := m.config.WriteFrame(cmFrame); err != nil {
		return fmt.Errorf("send CM_BAM: %w", err)
	}

	// Send DT frames with inter-frame delay.
	dtFrames := buildDTFrames(payload, source, BroadcastAddr)
	for i, f := range dtFrames {
		if i > 0 {
			time.Sleep(BAMInterFrameDelay)
		}
		if err := m.config.WriteFrame(f); err != nil {
			return fmt.Errorf("send DT frame %d: %w", i+1, err)
		}
	}

	return nil
}

// SendRTSCTS transmits a multi-frame message using RTS/CTS flow control.
// This blocks until the transfer completes or a timeout/error occurs.
func (m *Manager) SendRTSCTS(pgn uint32, source uint8, destination uint8, payload []byte) error {
	totalSize := uint16(len(payload))
	numFrames := uint8((len(payload) + MaxDTDataBytes - 1) / MaxDTDataBytes)

	key := sessionKey{source: source, destination: destination, pgn: pgn}

	sess := &session{
		key:       key,
		state:     stateWaitingForCTS,
		totalSize: totalSize,
		numFrames: numFrames,
		txPayload: payload,
		txOffset:  0,
		txSeqNum:  1,
		txDone:    make(chan struct{}),
	}

	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return fmt.Errorf("manager closed")
	}
	m.sessions[key] = sess

	// Set CTS timeout.
	sess.timer = time.AfterFunc(CTSTimeout, func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if _, ok := m.sessions[key]; ok {
			sess.txErr = fmt.Errorf("CTS timeout waiting for response")
			m.removeSession(key)
			close(sess.txDone)
		}
	})
	m.mu.Unlock()

	// Send RTS.
	rtsFrame := buildRTSFrame(totalSize, numFrames, pgn, source, destination)
	if err := m.config.WriteFrame(rtsFrame); err != nil {
		m.mu.Lock()
		m.removeSession(key)
		m.mu.Unlock()
		close(sess.txDone)
		return fmt.Errorf("send RTS: %w", err)
	}

	// Wait for completion.
	<-sess.txDone
	return sess.txErr
}

// handleCTSReceive processes an incoming CTS frame for a transmit-side RTS/CTS session.
func (m *Manager) handleCTSReceive(frame can.Frame, source uint8, destination uint8) {
	numFrames := frame.Data[1]
	nextSeqNum := frame.Data[2]
	pgnBytes := extractPGN(frame.Data)

	// For a transmit session, the CTS comes FROM the receiver (source of CTS) TO us (destination of CTS).
	// Our session key has our address as source and the CTS sender as destination.
	key := sessionKey{source: destination, destination: source, pgn: pgnBytes}

	m.mu.Lock()
	sess, ok := m.sessions[key]
	if !ok {
		m.mu.Unlock()
		return
	}

	// Stop the CTS timeout timer.
	if sess.timer != nil {
		sess.timer.Stop()
	}

	sess.state = stateSendingDT
	sess.txSeqNum = nextSeqNum
	m.mu.Unlock()

	// Send the requested DT frames.
	for i := uint8(0); i < numFrames; i++ {
		seqNum := nextSeqNum + i
		offset := int(seqNum-1) * MaxDTDataBytes
		if offset >= len(sess.txPayload) {
			break
		}

		end := offset + MaxDTDataBytes
		if end > len(sess.txPayload) {
			end = len(sess.txPayload)
		}

		var data [7]byte
		for j := range data {
			data[j] = 0xFF
		}
		copy(data[:], sess.txPayload[offset:end])

		dtFrame := makeDTFrame(seqNum, data, destination, source)
		if err := m.config.WriteFrame(dtFrame); err != nil {
			m.mu.Lock()
			sess.txErr = fmt.Errorf("send DT frame %d: %w", seqNum, err)
			m.removeSession(key)
			m.mu.Unlock()
			close(sess.txDone)
			return
		}
	}

	// After sending DT frames, set a CTS timeout for the next CTS or EndOfMsgAck.
	m.mu.Lock()
	if _, ok := m.sessions[key]; ok {
		sess.state = stateWaitingForCTS
		sess.timer = time.AfterFunc(CTSTimeout, func() {
			m.mu.Lock()
			defer m.mu.Unlock()
			if _, ok := m.sessions[key]; ok {
				sess.txErr = fmt.Errorf("CTS timeout after sending DT frames")
				m.removeSession(key)
				close(sess.txDone)
			}
		})
	}
	m.mu.Unlock()
}

// handleEndOfMsgAckReceive processes an EndOfMsgAck for a transmit-side session.
func (m *Manager) handleEndOfMsgAckReceive(frame can.Frame, source uint8, destination uint8) {
	pgnBytes := extractPGN(frame.Data)

	// Our session key: we are the transmitter (destination of EndOfMsgAck), receiver is source.
	key := sessionKey{source: destination, destination: source, pgn: pgnBytes}

	m.mu.Lock()
	sess, ok := m.sessions[key]
	if !ok {
		m.mu.Unlock()
		return
	}

	if sess.timer != nil {
		sess.timer.Stop()
	}
	m.removeSession(key)
	m.mu.Unlock()

	// Signal completion.
	close(sess.txDone)
}

// handleAbortReceive processes an Abort frame.
func (m *Manager) handleAbortReceive(frame can.Frame, source uint8, destination uint8) {
	pgnBytes := extractPGN(frame.Data)
	reason := frame.Data[1]

	m.logger.Warn("TP abort received", "source", source, "pgn", pgnBytes, "reason", reason)

	// Try both session key orientations.
	key1 := sessionKey{source: destination, destination: source, pgn: pgnBytes}
	key2 := sessionKey{source: source, destination: destination, pgn: pgnBytes}

	m.mu.Lock()
	defer m.mu.Unlock()

	if sess, ok := m.sessions[key1]; ok {
		if sess.txDone != nil {
			sess.txErr = fmt.Errorf("remote abort, reason=%d", reason)
			m.removeSession(key1)
			close(sess.txDone)
		} else {
			m.removeSession(key1)
		}
	}
	if sess, ok := m.sessions[key2]; ok {
		if sess.txDone != nil {
			sess.txErr = fmt.Errorf("remote abort, reason=%d", reason)
			m.removeSession(key2)
			close(sess.txDone)
		} else {
			m.removeSession(key2)
		}
	}
}

// sendEndOfMsgAck sends an EndOfMsgAck frame for a completed receive session.
func (m *Manager) sendEndOfMsgAck(sess *session) {
	var data [8]uint8
	data[0] = ControlEndOfMsgAck
	data[1] = uint8(sess.totalSize)
	data[2] = uint8(sess.totalSize >> 8)
	data[3] = sess.numFrames
	data[4] = 0xFF
	encodePGN(data[5:8], sess.key.pgn)

	canID := framer.BuildCANID(PGNCM, TPPriority, sess.key.destination, sess.key.source)
	f := can.Frame{
		ID:     canID,
		Length: 8,
		Data:   data,
	}
	if err := m.config.WriteFrame(f); err != nil {
		m.logger.Warn("failed to send EndOfMsgAck", "error", err)
	}
}

// Close stops all active sessions and prevents new ones from being created.
func (m *Manager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.closed = true
	for key, sess := range m.sessions {
		if sess.timer != nil {
			sess.timer.Stop()
		}
		if sess.txDone != nil {
			sess.txErr = fmt.Errorf("manager closed")
			select {
			case <-sess.txDone:
			default:
				close(sess.txDone)
			}
		}
		delete(m.sessions, key)
	}
}

// parseCANID extracts the PGN, source, and destination from a 29-bit CAN ID.
func parseCANID(id uint32) (pgn uint32, source uint8, destination uint8) {
	source = uint8(id & 0xFF)
	pgn = (id & 0x3FFFF00) >> 8

	pduFormat := uint8((pgn >> 8) & 0xFF)
	if pduFormat < 240 {
		destination = uint8(pgn & 0xFF)
		pgn &= 0xFFF00
	} else {
		destination = BroadcastAddr
	}
	return
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

// buildDTFrames creates all DT frames for a payload.
func buildDTFrames(payload []byte, source uint8, destination uint8) []can.Frame {
	numFrames := (len(payload) + MaxDTDataBytes - 1) / MaxDTDataBytes
	frames := make([]can.Frame, 0, numFrames)

	for i := 0; i < numFrames; i++ {
		seqNum := uint8(i + 1)
		offset := i * MaxDTDataBytes
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
