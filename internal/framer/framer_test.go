package framer_test

import (
	"testing"

	"github.com/open-ships/n2k/internal/framer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrameSingle_FullPayload(t *testing.T) {
	canID := framer.BuildCANID(127250, 2, 42, 255)
	payload := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	frame := framer.FrameSingle(canID, payload)

	assert.Equal(t, canID, frame.ID)
	assert.Equal(t, uint8(8), frame.Length)
	assert.Equal(t, [8]uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, frame.Data)
}

func TestFrameSingle_PaddedPayload(t *testing.T) {
	canID := framer.BuildCANID(127250, 2, 42, 255)
	payload := []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE}

	frame := framer.FrameSingle(canID, payload)

	assert.Equal(t, canID, frame.ID)
	assert.Equal(t, uint8(8), frame.Length)
	// First 5 bytes are payload, remaining 3 are padded with 0xFF.
	assert.Equal(t, [8]uint8{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0xFF, 0xFF}, frame.Data)
}

func TestFrameSingle_EmptyPayload(t *testing.T) {
	canID := framer.BuildCANID(127250, 2, 42, 255)
	payload := []byte{}

	frame := framer.FrameSingle(canID, payload)

	assert.Equal(t, canID, frame.ID)
	// All 8 bytes should be 0xFF padding.
	assert.Equal(t, [8]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, frame.Data)
}

func TestFrameFastPacket_SmallPayload(t *testing.T) {
	// 10 bytes: frame 0 carries 6 data bytes, frame 1 carries 4 data bytes.
	canID := framer.BuildCANID(127250, 2, 42, 255)
	payload := make([]byte, 10)
	for i := range payload {
		payload[i] = uint8(i + 1) // 1..10
	}
	seqId := uint8(3)

	frames := framer.FrameFastPacket(canID, payload, seqId)

	require.Len(t, frames, 2, "10-byte payload should produce 2 frames")

	// Frame 0
	f0 := frames[0]
	assert.Equal(t, canID, f0.ID)
	assert.Equal(t, uint8(3<<5), f0.Data[0], "frame 0 header: seqId=3, frameNum=0")
	assert.Equal(t, uint8(10), f0.Data[1], "frame 0 length byte")
	assert.Equal(t, payload[:6], f0.Data[2:8], "frame 0 should carry first 6 payload bytes")

	// Frame 1
	f1 := frames[1]
	assert.Equal(t, canID, f1.ID)
	assert.Equal(t, uint8((3<<5)|1), f1.Data[0], "frame 1 header: seqId=3, frameNum=1")
	assert.Equal(t, payload[6:10], f1.Data[1:5], "frame 1 should carry remaining 4 payload bytes")
	// Remaining bytes should be padded with 0xFF.
	assert.Equal(t, uint8(0xFF), f1.Data[5])
	assert.Equal(t, uint8(0xFF), f1.Data[6])
	assert.Equal(t, uint8(0xFF), f1.Data[7])
}

func TestFrameFastPacket_LargerPayload(t *testing.T) {
	// 24 bytes: frame 0 carries 6, frame 1 carries 7, frame 2 carries 7, frame 3 carries 4.
	canID := framer.BuildCANID(130306, 2, 100, 255)
	payload := make([]byte, 24)
	for i := range payload {
		payload[i] = uint8(i)
	}
	seqId := uint8(0)

	frames := framer.FrameFastPacket(canID, payload, seqId)

	// 24 bytes: 6 in frame 0, then 18 remaining => ceil(18/7) = 3 continuation frames
	require.Len(t, frames, 4, "24-byte payload should produce 4 frames")

	// Frame 0: 6 data bytes
	assert.Equal(t, uint8(0), frames[0].Data[0]&0x1F, "frame 0 number should be 0")
	assert.Equal(t, uint8(24), frames[0].Data[1], "total length should be 24")
	assert.Equal(t, payload[:6], frames[0].Data[2:8])

	// Frame 1: 7 data bytes (payload[6:13])
	assert.Equal(t, uint8(1), frames[1].Data[0]&0x1F, "frame 1 number")
	assert.Equal(t, payload[6:13], frames[1].Data[1:8])

	// Frame 2: 7 data bytes (payload[13:20])
	assert.Equal(t, uint8(2), frames[2].Data[0]&0x1F, "frame 2 number")
	assert.Equal(t, payload[13:20], frames[2].Data[1:8])

	// Frame 3: 4 data bytes (payload[20:24])
	assert.Equal(t, uint8(3), frames[3].Data[0]&0x1F, "frame 3 number")
	assert.Equal(t, payload[20:24], frames[3].Data[1:5])
	assert.Equal(t, uint8(0xFF), frames[3].Data[5], "padding byte")
	assert.Equal(t, uint8(0xFF), frames[3].Data[6], "padding byte")
	assert.Equal(t, uint8(0xFF), frames[3].Data[7], "padding byte")
}

func TestFrameFastPacket_SequenceIdEncoding(t *testing.T) {
	payload := make([]byte, 10)
	canID := framer.BuildCANID(127250, 2, 1, 255)

	for seqId := uint8(0); seqId < 8; seqId++ {
		frames := framer.FrameFastPacket(canID, payload, seqId)

		extractedSeqId := (frames[0].Data[0] & 0xE0) >> 5
		assert.Equal(t, seqId, extractedSeqId, "sequence ID %d should be encoded in upper 3 bits", seqId)

		extractedFrameNum := frames[0].Data[0] & 0x1F
		assert.Equal(t, uint8(0), extractedFrameNum, "frame 0 should have frame number 0")
	}
}

func TestFrameFastPacket_RoundTrip(t *testing.T) {
	// Verify that framing a payload and then manually reassembling it using the
	// same logic as the adapter's sequence.complete() produces the original payload.
	payload := make([]byte, 50)
	for i := range payload {
		payload[i] = uint8(i * 3)
	}

	canID := framer.BuildCANID(127250, 2, 42, 255)
	frames := framer.FrameFastPacket(canID, payload, 5)

	// Reassemble using the same logic as adapter/sequence.go:
	// Frame 0: expected length = Data[1], data = Data[2:8] (6 bytes)
	// Continuation frames: data = Data[1:8] (7 bytes)
	expectedLen := frames[0].Data[1]
	assert.Equal(t, uint8(50), expectedLen)

	var assembled []byte
	for i, f := range frames {
		if i == 0 {
			assembled = append(assembled, f.Data[2:]...)
		} else {
			assembled = append(assembled, f.Data[1:]...)
		}
		if len(assembled) >= int(expectedLen) {
			break
		}
	}
	assembled = assembled[:expectedLen]

	assert.Equal(t, payload, assembled, "reassembled payload should match original")
}

func TestFrameFastPacket_ExactlyThirteenBytes(t *testing.T) {
	// 13 bytes: frame 0 carries 6, frame 1 carries 7. Exactly fills 2 frames.
	payload := make([]byte, 13)
	for i := range payload {
		payload[i] = uint8(i + 0x10)
	}

	frames := framer.FrameFastPacket(framer.BuildCANID(127250, 2, 1, 255), payload, 0)
	require.Len(t, frames, 2, "13-byte payload should produce exactly 2 frames")

	// Frame 1 should be fully utilized (no padding except Data[0] header).
	assert.Equal(t, payload[6:13], frames[1].Data[1:8])
}

func TestFrameFastPacket_MaxPayload(t *testing.T) {
	// 223 bytes is the maximum fast-packet payload.
	payload := make([]byte, 223)
	for i := range payload {
		payload[i] = uint8(i)
	}

	frames := framer.FrameFastPacket(framer.BuildCANID(127250, 2, 1, 255), payload, 7)

	// 223 bytes: 6 in frame 0, then 217 remaining => ceil(217/7) = 31 continuation frames
	// Total: 32 frames (frame numbers 0-31).
	require.Len(t, frames, 32, "223-byte payload should produce 32 frames")

	// Verify last frame number is 31.
	lastFrameNum := frames[31].Data[0] & 0x1F
	assert.Equal(t, uint8(31), lastFrameNum)

	// Round-trip check.
	var assembled []byte
	for i, f := range frames {
		if i == 0 {
			assembled = append(assembled, f.Data[2:]...)
		} else {
			assembled = append(assembled, f.Data[1:]...)
		}
	}
	assembled = assembled[:223]
	assert.Equal(t, payload, assembled)
}
