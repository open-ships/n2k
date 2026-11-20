package decoder

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/open-ships/n2k/pgn"
)

// TestValid verifies that a completely empty Packet fails validation with exactly two errors:
// one for PGN == 0 and one for empty data. This tests the dual-check logic in Valid().
func TestValid(t *testing.T) {
	p := &Packet{}
	p.Valid()
	assert.Equal(t, len(p.ParseErrors), 2)
}

// TestGetSeqFrame verifies extraction of the sequence ID and frame number from the first
// byte of fast-packet data. The test byte 0x62 decodes as:
//   - Bits 7-5 = 0b011 = 3 (sequence ID)
//   - Bits 4-0 = 0b00010 = 2 (frame number)
//
// This confirms the bit masking and shifting logic in GetSeqFrame.
func TestGetSeqFrame(t *testing.T) {
	p := &Packet{
		Data: []uint8{0x62, 0x55, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0x7F},
	}
	p.GetSeqFrame()
	assert.Equal(t, uint8(2), p.FrameNum)
	assert.Equal(t, uint8(3), p.SeqId)
}

// TestNewPacket verifies that NewPacket correctly initializes all fields for a known
// proprietary PGN (130824). It checks:
//   - MessageInfo fields are preserved (PGN, SourceId, Priority)
//   - No parse errors for a valid, known PGN
//   - Fast flag is set (130824 is classified as a fast-packet PGN)
func TestNewPacket(t *testing.T) {
	pri, tgt := uint8(1), uint8(0)
	pInfo := pgn.MessageInfo{
		PGN:      130824,
		SourceId: 7,
		Priority: &pri,
		TargetId: &tgt,
	}
	p := NewPacket(pInfo, []uint8{(381 & 0xFF), (381 >> 8) | (4 << 5), 3, 4, 5, 0xFF, 0xFF, 0xFF})
	assert.Equal(t, uint32(130824), p.Info.PGN)
	assert.Equal(t, uint8(7), p.Info.SourceId)
	assert.Equal(t, uint8(1), *p.Info.Priority)
	assert.Equal(t, 0, len(p.ParseErrors))
	assert.True(t, p.Fast)
}

// TestBroadcast verifies that NewPacket handles broadcast PGNs (TargetId = 255) correctly.
// PGN 0x1ef00 is a broadcast PGN. The test confirms that the PGN number and broadcast
// target address (255) are preserved in the packet's MessageInfo.
func TestBroadcast(t *testing.T) {
	pri, tgt := uint8(1), uint8(255)
	pInfo := pgn.MessageInfo{
		PGN:      0x1ef00,
		SourceId: 7,
		Priority: &pri,
		TargetId: &tgt,
	}
	p := NewPacket(pInfo, []uint8{160, 9, (137 & 0xFF), (137 >> 8) | (4 << 5), 1, 2, 3, 4})
	assert.Equal(t, uint32(0x1ef00), p.Info.PGN)
	assert.Equal(t, uint8(255), *p.Info.TargetId)
}
