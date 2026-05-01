package pgn

import (
	"time"
)

// MessageInfo carries the CAN bus header metadata that accompanies every NMEA 2000 message.
// It is extracted from the CAN frame's 29-bit identifier and timestamp before the payload
// is passed to a PGN decoder. Every PGN struct embeds a MessageInfo as its "info"
// field, with an exported Info() accessor.
type MessageInfo struct {
	Timestamp time.Time `json:"timestamp"`
	Priority  *uint8    `json:"priority"`
	PGN       uint32    `json:"pgn"`
	SourceId  uint8     `json:"sourceId"`
	TargetId  *uint8    `json:"targetId"`
}

func Priority(v uint8) *uint8 { return &v }

func Target(v uint8) *uint8 { return &v }
