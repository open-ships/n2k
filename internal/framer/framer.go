package framer

import (
	"github.com/brutella/can"
)

// maxFastPacketPayload is the maximum payload size for fast-packet messages:
// 6 bytes in frame 0 + 31 continuation frames * 7 bytes = 223 bytes.
const maxFastPacketPayload = 223

// FrameSingle creates a single CAN frame for payloads of 8 bytes or fewer.
// Payloads shorter than 8 bytes are padded with 0xFF, which is the NMEA 2000
// convention for unused bytes.
//
// Parameters:
//   - canID: The 29-bit CAN ID constructed by BuildCANID.
//   - payload: The raw PGN payload bytes (must be <= 8 bytes).
//
// Returns a can.Frame ready for transmission.
func FrameSingle(canID uint32, payload []byte) can.Frame {
	f := can.Frame{
		ID:     canID,
		Length: 8,
	}

	// Initialize all data bytes to 0xFF (NMEA 2000 padding value).
	for i := range f.Data {
		f.Data[i] = 0xFF
	}

	// Copy payload into the frame, up to 8 bytes.
	n := len(payload)
	if n > 8 {
		n = 8
	}
	copy(f.Data[:n], payload)

	return f
}

// FrameFastPacket creates multiple CAN frames for fast-packet PGNs whose payload
// exceeds 8 bytes. The fast-packet protocol splits the payload across multiple CAN
// frames:
//
// Frame 0:
//   - Data[0] = (sequenceId << 5) | 0   (frame counter = 0)
//   - Data[1] = total payload length
//   - Data[2:8] = first 6 bytes of payload
//
// Continuation frames (1..N):
//   - Data[0] = (sequenceId << 5) | frameNumber
//   - Data[1:8] = next 7 bytes of payload
//
// The sequence ID is a 3-bit counter (0-7) managed by the caller to distinguish
// concurrent transmissions of the same PGN. Unused bytes in the last frame are
// padded with 0xFF.
//
// Parameters:
//   - canID: The 29-bit CAN ID constructed by BuildCANID.
//   - payload: The raw PGN payload bytes (up to 223 bytes).
//   - sequenceId: A 3-bit sequence counter (0-7) provided by the caller.
//
// Returns a slice of can.Frame values ready for sequential transmission.
func FrameFastPacket(canID uint32, payload []byte, sequenceId uint8) []can.Frame {
	totalLen := len(payload)
	if totalLen > maxFastPacketPayload {
		totalLen = maxFastPacketPayload
		payload = payload[:totalLen]
	}

	seqBits := (sequenceId & 0x07) << 5

	var frames []can.Frame
	offset := 0
	frameNum := 0

	// Frame 0: header byte, length byte, then up to 6 data bytes.
	{
		f := can.Frame{
			ID:     canID,
			Length: 8,
		}
		for i := range f.Data {
			f.Data[i] = 0xFF
		}

		f.Data[0] = seqBits | uint8(frameNum)
		f.Data[1] = uint8(totalLen)

		n := totalLen
		if n > 6 {
			n = 6
		}
		copy(f.Data[2:2+n], payload[:n])
		offset = n

		frames = append(frames, f)
		frameNum++
	}

	// Continuation frames: header byte, then up to 7 data bytes.
	for offset < totalLen {
		f := can.Frame{
			ID:     canID,
			Length: 8,
		}
		for i := range f.Data {
			f.Data[i] = 0xFF
		}

		f.Data[0] = seqBits | uint8(frameNum)

		n := totalLen - offset
		if n > 7 {
			n = 7
		}
		copy(f.Data[1:1+n], payload[offset:offset+n])
		offset += n

		frames = append(frames, f)
		frameNum++
	}

	return frames
}
