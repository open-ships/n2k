package transport

import (
	"context"
	"errors"
	"github.com/brutella/can"
	"github.com/open-ships/n2k/pgn"
)

// handleBAMReceive processes an incoming CM_BAM announcement and sets up a receive session
// to reassemble the subsequent DT frames.
func (m *Manager) handleBAMReceive(frame can.Frame, source uint8, info pgn.MessageInfo) {
	totalSize := uint16(frame.Data[1]) | uint16(frame.Data[2])<<8
	numFrames := frame.Data[3]
	pgn := extractPGN(frame.Data)

	if err := validateAnnouncement(totalSize, numFrames); err != nil {
		m.logger.Warn("invalid BAM announcement", "source", source, "pgn", pgn, "error", err)
		return
	}
	if err := validateTransportPGN(pgn); err != nil {
		m.logger.Warn("invalid BAM PGN", "source", source, "pgn", pgn, "error", err)
		return
	}

	key := sessionKey{
		source:      source,
		destination: BroadcastAddr,
		pgn:         pgn,
	}
	info.PGN = pgn
	info.TargetId = nil

	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.admitReceiveLocked(key) {
		return
	}
	ctx, cancel := context.WithCancelCause(context.Background())

	sess := &session{
		key:       key,
		state:     stateReceivingDT,
		totalSize: totalSize,
		numFrames: numFrames,
		received:  0,
		data:      make([]byte, totalSize),
		info:      info,
		ctx:       ctx,
		cancel:    cancel,
	}

	m.sessions[key] = sess
	m.armTimerLocked(sess, DTTimeout, errors.New("BAM DT timeout"))
	m.logger.Debug("BAM receive session started",
		"source", source, "pgn", pgn, "totalSize", totalSize, "numFrames", numFrames)
}
