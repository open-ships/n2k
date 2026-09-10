package transport

import (
	"fmt"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
)

const maxPassiveSessions = 512

type passiveKey struct {
	network             string
	connection, claim   uint64
	source, destination uint8
}

type passiveSession struct {
	info                        pgn.MessageInfo
	data                        []byte
	received                    int
	started, updated, timestamp time.Time
}

// PassiveReceiver reconstructs observed BAM and addressed transfers without
// writing CTS or acknowledgments. It requires an announcement and every DT
// frame in order, but does not require that the capture contain the peer's CTS.
// State is bounded across all networks and expires on capture or receipt time.
// The owner must serialize Handle calls. There are no background workers.
type PassiveReceiver struct {
	sessions map[passiveKey]*passiveSession
	now      func() time.Time
}

func NewPassiveReceiver() *PassiveReceiver {
	return &PassiveReceiver{sessions: make(map[passiveKey]*passiveSession), now: time.Now}
}

// Handle returns an owned payload only on completion. Malformed or excess
// transfers return an error. Expired sessions are discarded; unrelated frames
// and DT frames without a current announcement return no payload or error.
func (r *PassiveReceiver) Handle(frame can.Frame, info pgn.MessageInfo) (pgn.MessageInfo, []byte, error) {
	var empty pgn.MessageInfo
	id := framer.ParseCANID(frame.ID)
	if id.PGN != PGNCM && id.PGN != PGNDT {
		return empty, nil, nil
	}
	network := info.NetworkID
	if network == "" {
		network = info.AdapterID
	}
	key := passiveKey{network, info.ConnectionEpoch, info.ClaimEpoch, id.Source, id.Destination}
	now := r.now()
	for k, session := range r.sessions {
		captureExpired := k.network == key.network && k.connection == key.connection && k.claim == key.claim &&
			!info.Timestamp.IsZero() && !session.timestamp.IsZero() &&
			(info.Timestamp.Sub(session.timestamp) > CTSTimeout || info.Timestamp.Sub(session.info.Timestamp) > DefaultTransferTimeout)
		if now.Sub(session.updated) > CTSTimeout || now.Sub(session.started) > DefaultTransferTimeout || captureExpired {
			delete(r.sessions, k)
		}
	}
	if frame.Length != 8 {
		delete(r.sessions, key)
		return empty, nil, fmt.Errorf("ISO transport frame has %d bytes; expected 8", frame.Length)
	}
	if id.PGN == PGNCM {
		switch frame.Data[0] {
		case ControlBAM, ControlRTS:
			delete(r.sessions, key)
			bam := frame.Data[0] == ControlBAM
			if bam != (id.Destination == BroadcastAddr) || (!bam && frame.Data[4] == 0) {
				return empty, nil, fmt.Errorf("invalid passive ISO transport announcement")
			}
			size := uint16(frame.Data[1]) | uint16(frame.Data[2])<<8
			if err := validateAnnouncement(size, frame.Data[3]); err != nil {
				return empty, nil, err
			}
			number := extractPGN(frame.Data)
			if err := validateTransportPGN(number); err != nil {
				return empty, nil, err
			}
			if len(r.sessions) >= maxPassiveSessions {
				return empty, nil, fmt.Errorf("passive ISO transport table exceeds %d sessions", maxPassiveSessions)
			}
			info = info.Clone()
			info.PGN = number
			info.TargetId = nil
			if !bam {
				info.TargetId = pgn.Target(id.Destination)
			}
			r.sessions[key] = &passiveSession{info: info, data: make([]byte, int(size)), started: now, updated: now, timestamp: info.Timestamp}
		case ControlAbort:
			reverse := key
			reverse.source, reverse.destination = key.destination, key.source
			for _, k := range []passiveKey{key, reverse} {
				if s := r.sessions[k]; s != nil && s.info.PGN == extractPGN(frame.Data) {
					delete(r.sessions, k)
				}
			}
		case ControlCTS:
			reverse := key
			reverse.source, reverse.destination = key.destination, key.source
			if s := r.sessions[reverse]; s != nil && s.info.PGN == extractPGN(frame.Data) {
				s.updated, s.timestamp = now, info.Timestamp
			}
		}
		return empty, nil, nil
	}
	s := r.sessions[key]
	if s == nil {
		return empty, nil, nil
	}
	if int(frame.Data[0]) != s.received+1 {
		delete(r.sessions, key)
		return empty, nil, fmt.Errorf("passive ISO transport expected DT %d, got %d", s.received+1, frame.Data[0])
	}
	offset := s.received * MaxDTDataBytes
	copy(s.data[offset:], frame.Data[1:])
	s.received++
	s.updated, s.timestamp = now, info.Timestamp
	if s.received*MaxDTDataBytes < len(s.data) {
		return empty, nil, nil
	}
	delete(r.sessions, key)
	return s.info, s.data, nil
}
