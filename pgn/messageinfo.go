package pgn

import (
	"time"
)

// MessageInfo carries the CAN bus header metadata that accompanies every NMEA 2000 message.
// It is extracted from the CAN frame's 29-bit identifier and timestamp before the payload
// is passed to a PGN decoder. Every PGN struct embeds a MessageInfo as its "info"
// field, with an exported MessageInfo() accessor.
type MessageInfo struct {
	Timestamp time.Time `json:"timestamp"`
	Priority  *uint8    `json:"priority"`
	PGN       uint32    `json:"pgn"`
	SourceId  uint8     `json:"sourceId"`
	TargetId  *uint8    `json:"targetId"`

	// rawPayload and rawCanonical are codec bookkeeping. An untouched decoded
	// message returns rawPayload exactly; once fields differ from rawCanonical,
	// encoding uses the current field values. Replacing MessageInfo clears both.
	rawPayload   []uint8
	rawCanonical []uint8
}

// Priority returns a pointer to v, for use in MessageInfo literal construction:
//
//	pgn.MessageInfo{Priority: pgn.Priority(2)}
func Priority(v uint8) *uint8 { return &v }

// Target returns a pointer to v, for use in MessageInfo literal construction:
//
//	pgn.MessageInfo{TargetId: pgn.Target(42)}
func Target(v uint8) *uint8 { return &v }
