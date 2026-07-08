// Package transport implements ISO 11783 Transport Protocol for NMEA 2000 messages
// that exceed 8 bytes and cannot use fast-packet encoding. It supports two modes:
//
//   - BAM (Broadcast Announce Message): one-way broadcast without flow control
//   - RTS/CTS (Request To Send / Clear To Send): addressed point-to-point with flow control
//
// Both modes use TP.CM (PGN 60416) for connection management control frames and
// TP.DT (PGN 60160) for data transfer frames.
package transport

import (
	"time"

	"github.com/open-ships/n2k/internal/framer"
)

// TP PGN constants.
const (
	// PGNCM is the Connection Management PGN (TP.CM).
	PGNCM uint32 = 60416 // 0xEC00

	// PGNDT is the Data Transfer PGN (TP.DT).
	PGNDT uint32 = 60160 // 0xEB00
)

// CM control byte values identifying the type of connection management frame.
const (
	ControlRTS         uint8 = 16  // Request To Send
	ControlCTS         uint8 = 17  // Clear To Send
	ControlEndOfMsgAck uint8 = 19  // End of Message Acknowledgement
	ControlBAM         uint8 = 32  // Broadcast Announce Message
	ControlAbort       uint8 = 255 // Connection Abort
)

// Standard priority for transport protocol frames.
const TPPriority uint8 = 6

// BroadcastAddr is the broadcast destination address.
const BroadcastAddr = framer.BroadcastAddr

// Timeout durations per ISO 11783-3.
const (
	// DTTimeout is the maximum time to wait for the next DT frame.
	DTTimeout = 750 * time.Millisecond

	// CTSTimeout is the maximum time to wait for a CTS response after sending RTS or DT frames.
	CTSTimeout = 1250 * time.Millisecond
)

// BAM inter-frame delay (minimum 50ms per spec).
const BAMInterFrameDelay = 50 * time.Millisecond

// MaxDTDataBytes is the number of data bytes carried per DT frame (bytes 1-7).
const MaxDTDataBytes = 7
