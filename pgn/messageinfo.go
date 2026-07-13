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

	// rawPayload is the exact wire payload the owning message was decoded
	// from, stashed by decodeFields and returned verbatim by encodeFields --
	// a decoded message always re-encodes to the bytes it came from,
	// regardless of any field changes made in between. Unexported: this is
	// codec bookkeeping, not part of MessageInfo's public contract, and is
	// cleared whenever a caller replaces Info wholesale (e.g. via
	// SetMessageInfo with a fresh value), which conveniently also
	// invalidates it whenever a caller is about to repurpose the message
	// for something else (the only supported way back to the normal,
	// build-from-scratch encode path for a message that was once decoded).
	rawPayload []uint8
}

// Priority returns a pointer to v, for use in MessageInfo literal construction:
//
//	pgn.MessageInfo{Priority: pgn.Priority(2)}
func Priority(v uint8) *uint8 { return &v }

// Target returns a pointer to v, for use in MessageInfo literal construction:
//
//	pgn.MessageInfo{TargetId: pgn.Target(42)}
func Target(v uint8) *uint8 { return &v }
