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

	// LocalAddress returns the client's currently claimed address. When set,
	// addressed TP traffic for other nodes is ignored. A function is used
	// because address claiming may move the client at runtime.
	LocalAddress func() uint8

	// OnComplete is called when a multi-frame message has been fully reassembled.
	// It receives the transported PGN, source address, destination address, and
	// the assembled payload.
	OnComplete func(pgn uint32, source uint8, destination uint8, data []byte)

	// Logger for transport protocol events. If nil, a no-op logger is used.
	Logger *slog.Logger

	// AfterFunc schedules a callback after d; nil means time.AfterFunc.
	// Tests inject a fake to drive timeouts deterministically.
	AfterFunc func(d time.Duration, f func()) Timer

	// Sleep pauses between BAM DT frames; nil means time.Sleep.
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
	timer Timer

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
	sleep := cfg.Sleep
	if sleep == nil {
		sleep = time.Sleep
	}
	return &Manager{
		config:    cfg,
		sessions:  make(map[sessionKey]*session),
		logger:    logger,
		afterFunc: afterFunc,
		sleep:     sleep,
	}
}

// HandleFrame routes an incoming CAN frame to the appropriate transport protocol
// handler based on PGN. Non-TP frames are silently ignored.
func (m *Manager) HandleFrame(frame can.Frame) {
	c := framer.ParseCANID(frame.ID)
	switch c.PGN {
	case PGNCM:
		m.handleCM(frame, c.Source, c.Destination)
	case PGNDT:
		m.handleDT(frame, c.Source, c.Destination)
	}
}

// handleCM dispatches a Connection Management frame based on the control byte.
func (m *Manager) handleCM(frame can.Frame, source uint8, destination uint8) {
	if frame.Length != 8 {
		return
	}
	controlByte := frame.Data[0]

	switch controlByte {
	case ControlBAM:
		if destination != BroadcastAddr {
			return
		}
		m.handleBAMReceive(frame, source)
	case ControlRTS:
		if !m.isLocalDestination(destination) {
			return
		}
		m.handleRTSReceive(frame, source, destination)
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
		m.removeSession(sess.key)
		m.mu.Unlock()
		return
	}

	// Copy data bytes into the reassembly buffer.
	offset := int(seqNum-1) * MaxDTDataBytes
	remaining := int(sess.totalSize) - offset
	if remaining <= 0 || offset < 0 || offset >= len(sess.data) {
		m.logger.Warn("DT frame exceeds announced payload", "sequence", seqNum, "size", sess.totalSize)
		m.removeSession(sess.key)
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
	if sess.timer != nil {
		sess.timer.Stop()
	}

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
		sess.timer = m.afterFunc(DTTimeout, func() {
			m.mu.Lock()
			defer m.mu.Unlock()
			m.logger.Warn("DT timeout", "source", source, "pgn", sess.key.pgn,
				"received", sess.received, "expected", sess.numFrames)
			m.removeSession(sess.key)
		})
		key := sess.key
		m.mu.Unlock()
		if nextCTS != nil {
			if err := m.config.WriteFrame(*nextCTS); err != nil {
				m.mu.Lock()
				m.removeSession(key)
				m.mu.Unlock()
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

	// For RTS/CTS receive: build the EndOfMsgAck to send after unlocking.
	var ackFrame *can.Frame
	if dst != BroadcastAddr {
		f := buildEndOfMsgAckFrame(sess)
		ackFrame = &f
	}

	m.removeSession(key)
	m.mu.Unlock()

	if ackFrame != nil {
		if err := m.config.WriteFrame(*ackFrame); err != nil {
			m.logger.Warn("failed to send EndOfMsgAck", "error", err)
		}
	}
	if m.config.OnComplete != nil {
		m.config.OnComplete(key.pgn, key.source, dst, data)
	}
}

// findSession looks up a session by source and destination. It tries an exact match
// first, then falls back to broadcast destination for BAM sessions.
func (m *Manager) findSession(source uint8, destination uint8) *session {
	// Try each PGN — we iterate sessions since PGN is part of the key.
	for k, s := range m.sessions {
		if k.source == source && k.destination == destination && s.state == stateReceivingDT {
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

// removeReceiveSessions removes the one active receive transfer for a source
// and destination. DT packets carry no PGN, so allowing multiple such
// sessions would make packet ownership ambiguous.
func (m *Manager) removeReceiveSessions(source, destination uint8) {
	for key, sess := range m.sessions {
		if key.source == source && key.destination == destination && sess.state == stateReceivingDT {
			m.removeSession(key)
		}
	}
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
	if err := validateTransportPGN(pgn); err != nil {
		return fmt.Errorf("send BAM: %w", err)
	}
	if err := validatePayload(payload); err != nil {
		return fmt.Errorf("send BAM: %w", err)
	}
	m.mu.Lock()
	closed := m.closed
	m.mu.Unlock()
	if closed {
		return fmt.Errorf("send BAM: manager closed")
	}
	payload = append([]byte(nil), payload...)
	totalSize := uint16(len(payload))
	numFrames := uint8((len(payload) + MaxDTDataBytes - 1) / MaxDTDataBytes)

	// Build and send CM_BAM announcement.
	cmFrame := buildCMBAMFrame(totalSize, numFrames, pgn, source)
	if err := m.config.WriteFrame(cmFrame); err != nil {
		return fmt.Errorf("send CM_BAM: %w", err)
	}

	// Send DT frames with inter-frame delay.
	dtFrames := buildDTFrameRange(payload, 1, int(numFrames), source, BroadcastAddr)
	for i, f := range dtFrames {
		if i > 0 {
			m.sleep(BAMInterFrameDelay)
		}
		m.mu.Lock()
		closed := m.closed
		m.mu.Unlock()
		if closed {
			return fmt.Errorf("send BAM: manager closed")
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
	if destination == BroadcastAddr {
		return fmt.Errorf("send RTS/CTS: broadcast destination is invalid")
	}
	if err := validatePayload(payload); err != nil {
		return fmt.Errorf("send RTS/CTS: %w", err)
	}
	if err := validateTransportPGN(pgn); err != nil {
		return fmt.Errorf("send RTS/CTS: %w", err)
	}
	totalSize := uint16(len(payload))
	numFrames := uint8((len(payload) + MaxDTDataBytes - 1) / MaxDTDataBytes)

	key := sessionKey{source: source, destination: destination, pgn: pgn}

	sess := &session{
		key:       key,
		state:     stateWaitingForCTS,
		totalSize: totalSize,
		numFrames: numFrames,
		txPayload: append([]byte(nil), payload...),
		txDone:    make(chan struct{}),
	}

	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return fmt.Errorf("manager closed")
	}
	for activeKey, active := range m.sessions {
		if active.txDone != nil && activeKey.source == source && activeKey.destination == destination {
			m.mu.Unlock()
			return fmt.Errorf("transport session already active for source %d destination %d", source, destination)
		}
	}
	m.sessions[key] = sess

	// Set CTS timeout.
	sess.timer = m.afterFunc(CTSTimeout, func() {
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
		// The CTS timeout may have fired (removing the session and closing
		// txDone) while WriteFrame blocked; only close txDone if we removed
		// the session ourselves.
		m.mu.Lock()
		_, active := m.sessions[key]
		if active {
			m.removeSession(key)
		}
		m.mu.Unlock()
		if active {
			close(sess.txDone)
		}
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
	if !ok || sess.state != stateWaitingForCTS {
		m.mu.Unlock()
		return
	}
	if numFrames == 0 {
		// A zero-packet CTS asks the transmitter to keep waiting. Re-arm the
		// timeout so a peer can apply backpressure without corrupting state.
		if sess.timer != nil {
			sess.timer.Stop()
		}
		sess.timer = m.afterFunc(CTSTimeout, func() {
			m.mu.Lock()
			defer m.mu.Unlock()
			if _, active := m.sessions[key]; active {
				sess.txErr = fmt.Errorf("CTS timeout while receiver paused transfer")
				m.removeSession(key)
				close(sess.txDone)
			}
		})
		m.mu.Unlock()
		return
	}
	lastSeqNum := int(nextSeqNum) + int(numFrames) - 1
	if nextSeqNum == 0 || int(nextSeqNum) > int(sess.numFrames) || lastSeqNum > int(sess.numFrames) {
		sess.txErr = fmt.Errorf(
			"invalid CTS range: first=%d count=%d total=%d",
			nextSeqNum,
			numFrames,
			sess.numFrames,
		)
		m.removeSession(key)
		m.mu.Unlock()
		close(sess.txDone)
		return
	}

	// Stop the CTS timeout timer.
	if sess.timer != nil {
		sess.timer.Stop()
	}

	sess.state = stateSendingDT
	payload := sess.txPayload // immutable after session creation; safe to read unlocked
	m.mu.Unlock()

	// Send the requested DT frames without holding the manager lock.
	for _, dtFrame := range buildDTFrameRange(payload, nextSeqNum, int(numFrames), destination, source) {
		if err := m.config.WriteFrame(dtFrame); err != nil {
			// The session may have been removed (abort/close) while writing;
			// only mutate and signal completion if we still own it.
			m.mu.Lock()
			_, active := m.sessions[key]
			if active {
				sess.txErr = fmt.Errorf("send DT frame %d: %w", dtFrame.Data[0], err)
				m.removeSession(key)
			}
			m.mu.Unlock()
			if active {
				close(sess.txDone)
			}
			return
		}
		m.mu.Lock()
		seq := dtFrame.Data[0]
		if !sess.txSent[seq] {
			sess.txSent[seq] = true
			sess.txSentCount++
		}
		m.mu.Unlock()
	}

	// After sending DT frames, set a CTS timeout for the next CTS or EndOfMsgAck.
	// The session may have been removed (abort/close) while frames were being
	// written; only re-arm if it still exists.
	m.mu.Lock()
	if _, ok := m.sessions[key]; ok {
		sess.state = stateWaitingForCTS
		sess.timer = m.afterFunc(CTSTimeout, func() {
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
	if !ok || sess.state != stateWaitingForCTS {
		m.mu.Unlock()
		return
	}
	totalSize := uint16(frame.Data[1]) | uint16(frame.Data[2])<<8
	numFrames := frame.Data[3]
	if totalSize != sess.totalSize || numFrames != sess.numFrames || sess.txSentCount != int(sess.numFrames) {
		sess.txErr = fmt.Errorf(
			"invalid EndOfMsgAck: size=%d frames=%d packets sent=%d; want size=%d frames=%d",
			totalSize,
			numFrames,
			sess.txSentCount,
			sess.totalSize,
			sess.numFrames,
		)
		m.removeSession(key)
		m.mu.Unlock()
		close(sess.txDone)
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
