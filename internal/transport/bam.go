package transport

import (
	"github.com/brutella/can"
)

// handleBAMReceive processes an incoming CM_BAM announcement and sets up a receive session
// to reassemble the subsequent DT frames.
func (m *Manager) handleBAMReceive(frame can.Frame, source uint8) {
	totalSize := uint16(frame.Data[1]) | uint16(frame.Data[2])<<8
	numFrames := frame.Data[3]
	pgn := extractPGN(frame.Data)

	if totalSize == 0 || numFrames == 0 {
		m.logger.Warn("BAM with zero size or frames", "source", source, "pgn", pgn)
		return
	}

	key := sessionKey{
		source:      source,
		destination: BroadcastAddr,
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
		received:  0,
		data:      make([]byte, totalSize),
	}

	// Set a DT timeout for the first frame.
	sess.timer = m.afterFunc(DTTimeout, func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.logger.Warn("BAM DT timeout", "source", source, "pgn", pgn,
			"received", sess.received, "expected", numFrames)
		m.removeSession(key)
	})

	m.sessions[key] = sess
	m.logger.Debug("BAM receive session started",
		"source", source, "pgn", pgn, "totalSize", totalSize, "numFrames", numFrames)
}
