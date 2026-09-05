package pgn

import (
	"time"

	"github.com/open-ships/n2k/raw"
)

// MessageInfo carries the CAN bus header metadata that accompanies every NMEA 2000 message.
// It is extracted from the CAN frame's 29-bit identifier and timestamp before the payload
// is passed to a PGN decoder. Every PGN struct embeds a MessageInfo as its "info"
// field, with an exported MessageInfo() accessor.
type MessageInfo struct {
	Timestamp             time.Time     `json:"timestamp"`
	ReceivedAt            time.Time     `json:"receivedAt,omitempty"`
	TransportTimestamp    time.Duration `json:"transportTimestamp,omitempty"`
	HasTransportTimestamp bool          `json:"hasTransportTimestamp,omitempty"`
	AdapterID             string        `json:"adapterId,omitempty"`
	NetworkID             string        `json:"networkId,omitempty"`
	Direction             raw.Direction `json:"direction,omitempty"`
	ConnectionEpoch       uint64        `json:"connectionEpoch,omitempty"`
	ClaimEpoch            uint64        `json:"claimEpoch,omitempty"`
	// DecodeIssues makes partial decoding explicit while retained wire bytes
	// continue to support unchanged forwarding.
	DecodeIssues []string `json:"decodeIssues,omitempty"`
	Priority     *uint8   `json:"priority"`
	PGN          uint32   `json:"pgn"`
	SourceId     uint8    `json:"sourceId"`
	TargetId     *uint8   `json:"targetId"`

	// rawPayload and rawCanonical are codec bookkeeping. An untouched decoded
	// message returns rawPayload exactly; once fields differ from rawCanonical,
	// encoding uses the current field values. Replacing MessageInfo clears both.
	rawPayload   []uint8
	rawCanonical []uint8
}

// Clone returns independently owned metadata, including retained wire bytes.
func (info MessageInfo) Clone() MessageInfo {
	info.Priority = clonePointer(info.Priority)
	info.TargetId = clonePointer(info.TargetId)
	info.rawPayload = cloneSlice(info.rawPayload)
	info.rawCanonical = cloneSlice(info.rawCanonical)
	info.DecodeIssues = cloneSlice(info.DecodeIssues)
	return info
}

// Priority returns a pointer to v, for use in MessageInfo literal construction:
//
//	pgn.MessageInfo{Priority: pgn.Priority(2)}
func Priority(v uint8) *uint8 { return &v }

// Target returns a pointer to v, for use in MessageInfo literal construction:
//
//	pgn.MessageInfo{TargetId: pgn.Target(42)}
func Target(v uint8) *uint8 { return &v }
