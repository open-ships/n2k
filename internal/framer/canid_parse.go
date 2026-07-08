package framer

// Protocol PGNs shared across the read/write pipeline.
const (
	// PGNISOAddressClaim is the ISO 11783 Address Claim PGN.
	PGNISOAddressClaim uint32 = 60928
	// PGNISORequest is the ISO Request PGN.
	PGNISORequest uint32 = 59904
)

// BroadcastAddr is the NMEA 2000 broadcast destination address.
const BroadcastAddr uint8 = 255

// CANID is the decomposed form of a 29-bit extended CAN identifier.
// ParseCANID and BuildCANID are inverses for the fields they share.
type CANID struct {
	PGN         uint32
	Priority    uint8
	Source      uint8
	Destination uint8
	Addressed   bool
}

// ParseCANID extracts the NMEA 2000 fields from a 29-bit CAN ID. For PDU1
// (addressed) PGNs the low byte of the raw PGN field is the destination
// address; it is masked off the returned PGN. For PDU2 (broadcast) PGNs the
// destination is BroadcastAddr.
func ParseCANID(id uint32) CANID {
	c := CANID{
		Priority: uint8((id & 0x1C000000) >> 26),
		Source:   uint8(id & 0xFF),
		PGN:      (id & 0x3FFFF00) >> 8,
	}
	pduFormat := uint8((c.PGN & 0xFF00) >> 8)
	if pduFormat < 240 {
		c.Addressed = true
		c.Destination = uint8(c.PGN & 0xFF)
		c.PGN &= 0xFFF00
	} else {
		c.Destination = BroadcastAddr
	}
	return c
}

// FastPacketSeqFrame splits a fast-packet header byte into its 3-bit sequence
// ID (bits 7-5) and 5-bit frame number (bits 4-0).
func FastPacketSeqFrame(headerByte byte) (seqID uint8, frameNum uint8) {
	return (headerByte & 0xE0) >> 5, headerByte & 0x1F
}

// FastPacketHeader builds a fast-packet header byte from a sequence ID and
// frame number. Inverse of FastPacketSeqFrame.
func FastPacketHeader(seqID uint8, frameNum uint8) byte {
	return ((seqID & 0x07) << 5) | (frameNum & 0x1F)
}
