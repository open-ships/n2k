package transport

import (
	"context"
	"errors"
	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
)

// handleRTSReceive processes an incoming RTS frame and sets up a receive session.
// It responds with a CTS frame to begin the data transfer.
func (m *Manager) handleRTSReceive(frame can.Frame, source uint8, destination uint8, info pgn.MessageInfo) {
	totalSize := uint16(frame.Data[1]) | uint16(frame.Data[2])<<8
	numFrames := frame.Data[3]
	maxPerCTS := frame.Data[4]
	pgn := extractPGN(frame.Data)

	if err := validateAnnouncement(totalSize, numFrames); err != nil {
		m.logger.Warn("invalid RTS announcement", "source", source, "pgn", pgn, "error", err)
		return
	}
	if err := validateTransportPGN(pgn); err != nil {
		m.logger.Warn("invalid RTS PGN", "source", source, "pgn", pgn, "error", err)
		return
	}
	if maxPerCTS == 0 {
		m.logger.Warn("invalid RTS maximum packets per CTS", "source", source, "pgn", pgn)
		return
	}

	key := sessionKey{
		source:      source,
		destination: destination,
		pgn:         pgn,
	}
	info.PGN = pgn
	info.TargetId = pgnTarget(destination)

	m.mu.Lock()

	if !m.admitReceiveLocked(key) {
		m.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancelCause(context.Background())

	sess := &session{
		key:       key,
		state:     stateReceivingDT,
		totalSize: totalSize,
		numFrames: numFrames,
		maxPerCTS: maxPerCTS,
		received:  0,
		data:      make([]byte, totalSize),
		info:      info,
		ctx:       ctx,
		cancel:    cancel,
	}

	m.sessions[key] = sess
	m.armTimerLocked(sess, DTTimeout, errors.New("RTS/CTS DT timeout"))
	m.mu.Unlock()

	m.logger.Debug("RTS/CTS receive session started",
		"source", source, "destination", destination, "pgn", pgn,
		"totalSize", totalSize, "numFrames", numFrames)

	// Send CTS: request all frames starting from 1.
	requested := numFrames
	if maxPerCTS < requested {
		requested = maxPerCTS
	}
	if err := m.writeSessionFrame(sess, buildCTSFrame(requested, 1, pgn, destination, source)); err != nil {
		_ = m.finishSession(sess, err)
		m.logger.Warn("failed to send CTS", "error", err, "pgn", pgn)
	}
}

func pgnTarget(destination uint8) *uint8 { return &destination }

// buildCTSFrame constructs a CTS (Clear To Send) frame.
func buildCTSFrame(numFrames uint8, nextSeqNum uint8, pgn uint32, source uint8, destination uint8) can.Frame {
	var data [8]uint8
	data[0] = ControlCTS
	data[1] = numFrames
	data[2] = nextSeqNum
	data[3] = 0xFF
	data[4] = 0xFF
	encodePGN(data[5:8], pgn)

	canID := framer.BuildCANID(PGNCM, TPPriority, source, destination)
	return can.Frame{
		ID:     canID,
		Length: 8,
		Data:   data,
	}
}
