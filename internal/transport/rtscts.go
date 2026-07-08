package transport

import (
	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
)

// handleRTSReceive processes an incoming RTS frame and sets up a receive session.
// It responds with a CTS frame to begin the data transfer.
func (m *Manager) handleRTSReceive(frame can.Frame, source uint8, destination uint8) {
	totalSize := uint16(frame.Data[1]) | uint16(frame.Data[2])<<8
	numFrames := frame.Data[3]
	maxPerCTS := frame.Data[4]
	pgn := extractPGN(frame.Data)

	if totalSize == 0 || numFrames == 0 {
		m.logger.Warn("RTS with zero size or frames", "source", source, "pgn", pgn)
		return
	}

	key := sessionKey{
		source:      source,
		destination: destination,
		pgn:         pgn,
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// If there's already a session for this key, clean it up first.
	m.removeSession(key)

	sess := &session{
		key:       key,
		state:     stateReceivingDT,
		totalSize: totalSize,
		numFrames: numFrames,
		maxPerCTS: maxPerCTS,
		received:  0,
		data:      make([]byte, totalSize),
	}

	// Set a DT timeout.
	sess.timer = m.afterFunc(DTTimeout, func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.logger.Warn("RTS/CTS DT timeout", "source", source, "pgn", pgn,
			"received", sess.received, "expected", numFrames)
		m.removeSession(key)
	})

	m.sessions[key] = sess

	m.logger.Debug("RTS/CTS receive session started",
		"source", source, "destination", destination, "pgn", pgn,
		"totalSize", totalSize, "numFrames", numFrames)

	// Send CTS: request all frames starting from 1.
	m.sendCTS(numFrames, 1, pgn, destination, source)
}

// sendCTS sends a CTS (Clear To Send) frame.
func (m *Manager) sendCTS(numFrames uint8, nextSeqNum uint8, pgn uint32, source uint8, destination uint8) {
	var data [8]uint8
	data[0] = ControlCTS
	data[1] = numFrames
	data[2] = nextSeqNum
	data[3] = 0xFF
	data[4] = 0xFF
	encodePGN(data[5:8], pgn)

	canID := framer.BuildCANID(PGNCM, TPPriority, source, destination)
	f := can.Frame{
		ID:     canID,
		Length: 8,
		Data:   data,
	}
	if err := m.config.WriteFrame(f); err != nil {
		m.logger.Warn("failed to send CTS", "error", err, "pgn", pgn)
	}
}
