package adapter

import (
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
)

// This data structure is copied from
// https://github.com/brutella/can/blob/master/frame.go
// licensed under the MIT License, following

/*
The MIT License (MIT)

Copyright (c) 2016 Matthias Hochgatterer

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

// NewPacketInfo extracts NMEA 2000 message metadata from a raw CAN bus frame's 29-bit
// extended identifier. The CAN ID encodes several fields using the NMEA 2000 / ISO 11783
// bit layout:
//
// CAN ID bit layout (29 bits total, extended frame format):
//
//	Bits 28-26: Priority (3 bits, 0-7, lower = higher priority)
//	Bit  25:    Reserved
//	Bit  24:    Data Page (DP)
//	Bits 23-16: PDU Format (PF) - determines if message is broadcast or addressed
//	Bits 15-8:  PDU Specific (PS) - destination address (if PF < 240) or group extension (if PF >= 240)
//	Bits 7-0:   Source Address (SA) - the sender's address on the bus
//
// The PGN (Parameter Group Number) is extracted from bits 25-8 of the CAN ID. For
// addressed messages (PDU Format < 240), the PS field contains the destination address
// rather than being part of the PGN, so the lower 8 bits are masked off and stored as
// TargetId instead.
//
// Parameters:
//   - message: A pointer to a can.Frame containing the raw CAN bus frame with its 29-bit ID.
//
// Returns a pgn.MessageInfo populated with the extracted PGN, source, priority, target,
// and current timestamp.
func NewPacketInfo(message *can.Frame) pgn.MessageInfo {
	return NewPacketInfoAt(message, time.Now())
}

// NewPacketInfoAt extracts CAN identifier metadata using timestamp as the
// observation time. It is used when a capture or gateway supplied a more
// faithful time than the host's current clock.
func NewPacketInfoAt(message *can.Frame, timestamp time.Time) pgn.MessageInfo {
	c := framer.ParseCANID(message.ID)
	p := pgn.MessageInfo{
		Timestamp: timestamp,
		SourceId:  c.Source,
		PGN:       c.PGN,
		Priority:  &c.Priority,
	}
	if c.Addressed {
		targetId := c.Destination
		p.TargetId = &targetId
	}
	return p
}
