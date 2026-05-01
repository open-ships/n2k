// Package framer builds CAN frames from encoded PGN payloads. It is the inverse of the
// read-side frame parsing in internal/adapter: where adapter.NewPacketInfo extracts NMEA 2000
// metadata from a CAN ID, this package constructs CAN IDs and frames for transmission.
package framer

// BuildCANID constructs a 29-bit extended CAN ID from NMEA 2000 parameters.
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
// For PDU2 (broadcast) PGNs where PF >= 240, the PGN encodes both PF and PS fields.
// For PDU1 (addressed) PGNs where PF < 240, the PGN's low byte is zero and the
// destination address is placed into the PS field (bits 15-8 of the CAN ID).
//
// Parameters:
//   - pgn: The Parameter Group Number (18-bit). For PDU1 PGNs, the low byte must be 0.
//   - priority: Message priority (0-7). Only the lower 3 bits are used.
//   - source: Source address of the transmitting device (0-253, 254-255 reserved).
//   - destination: Destination address for PDU1 (addressed) PGNs. Ignored for PDU2 (broadcast) PGNs.
//
// Returns the constructed 29-bit CAN ID.
func BuildCANID(pgn uint32, priority uint8, source uint8, destination uint8) uint32 {
	pduFormat := uint8((pgn >> 8) & 0xFF)

	var id uint32
	id |= uint32(priority&0x07) << 26
	id |= pgn << 8

	if pduFormat < 240 {
		// PDU1 (addressed): the PGN's low byte is 0, so we OR the destination
		// address into the PS field at bits 15-8 of the CAN ID.
		id |= uint32(destination) << 8
	}

	id |= uint32(source)

	return id
}
