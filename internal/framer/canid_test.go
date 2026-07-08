package framer_test

import (
	"testing"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/adapter"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/stretchr/testify/assert"
)

func TestBuildCANID_PDU2Broadcast(t *testing.T) {
	// VesselHeading PGN 127250 = 0x1F112, PDU Format = 0xF1 (241 >= 240 -> broadcast)
	canID := framer.BuildCANID(127250, 2, 42, 255)

	frame := can.Frame{ID: canID}
	info := adapter.NewPacketInfo(&frame)

	assert.Equal(t, uint32(127250), info.PGN, "PGN should be 127250")
	assert.Equal(t, uint8(42), info.SourceId, "Source should be 42")
	assert.Equal(t, uint8(2), *info.Priority, "Priority should be 2")
}

func TestBuildCANID_PDU1Addressed(t *testing.T) {
	// IsoRequest PGN 59904 = 0xEA00, PDU Format = 0xEA (234 < 240 -> addressed)
	canID := framer.BuildCANID(59904, 6, 0, 255)

	frame := can.Frame{ID: canID}
	info := adapter.NewPacketInfo(&frame)

	assert.Equal(t, uint32(59904), info.PGN, "PGN should be 59904")
	assert.Equal(t, uint8(0), info.SourceId, "Source should be 0")
	assert.Equal(t, uint8(6), *info.Priority, "Priority should be 6")
	assert.Equal(t, uint8(255), *info.TargetId, "TargetId should be 255")
}

func TestBuildCANID_PDU1AddressedSpecificTarget(t *testing.T) {
	// AddressClaimPGN 60928 = 0xEE00, PDU Format = 0xEE (238 < 240 -> addressed)
	canID := framer.BuildCANID(60928, 3, 10, 20)

	frame := can.Frame{ID: canID}
	info := adapter.NewPacketInfo(&frame)

	assert.Equal(t, uint32(60928), info.PGN, "PGN should be 60928")
	assert.Equal(t, uint8(10), info.SourceId, "Source should be 10")
	assert.Equal(t, uint8(3), *info.Priority, "Priority should be 3")
	assert.Equal(t, uint8(20), *info.TargetId, "TargetId should be 20")
}

func TestBuildCANID_PriorityRange(t *testing.T) {
	for priority := uint8(0); priority <= 7; priority++ {
		canID := framer.BuildCANID(127250, priority, 1, 255)

		frame := can.Frame{ID: canID}
		info := adapter.NewPacketInfo(&frame)

		assert.Equal(t, priority, *info.Priority, "Priority %d should round-trip", priority)
		assert.Equal(t, uint32(127250), info.PGN, "PGN should be preserved for priority %d", priority)
		assert.Equal(t, uint8(1), info.SourceId, "Source should be preserved for priority %d", priority)
	}
}

func TestBuildCANID_PDU2VariousPGNs(t *testing.T) {
	// Test several well-known broadcast PGNs
	tests := []struct {
		name     string
		pgn      uint32
		priority uint8
		source   uint8
	}{
		{"WindData", 130306, 2, 100},
		{"SystemTime", 126992, 3, 0},
		{"WaterDepth", 128267, 2, 55},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			canID := framer.BuildCANID(tt.pgn, tt.priority, tt.source, 255)

			frame := can.Frame{ID: canID}
			info := adapter.NewPacketInfo(&frame)

			assert.Equal(t, tt.pgn, info.PGN)
			assert.Equal(t, tt.source, info.SourceId)
			assert.Equal(t, tt.priority, *info.Priority)
		})
	}
}

func TestBuildCANID_PriorityClampedTo3Bits(t *testing.T) {
	// Priorities above 7 should be masked to 3 bits
	canID := framer.BuildCANID(127250, 0xFF, 1, 255)

	frame := can.Frame{ID: canID}
	info := adapter.NewPacketInfo(&frame)

	assert.Equal(t, uint8(7), *info.Priority, "Priority 0xFF should be masked to 7")
}

func TestParseCANIDInverseOfBuild(t *testing.T) {
	cases := []struct {
		name               string
		pgn                uint32
		priority, src, dst uint8
		wantAddressed      bool
	}{
		{"pdu2 broadcast", 127250, 2, 42, 255, false},
		{"pdu1 addressed", 59904, 6, 10, 42, true},
		{"pdu1 broadcast dest", 60928, 6, 253, 255, true},
		{"high pgn", 130824, 7, 1, 255, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			id := framer.BuildCANID(tc.pgn, tc.priority, tc.src, tc.dst)
			got := framer.ParseCANID(id)
			if got.PGN != tc.pgn || got.Priority != tc.priority || got.Source != tc.src {
				t.Fatalf("ParseCANID(BuildCANID) = %+v, want pgn=%d pri=%d src=%d", got, tc.pgn, tc.priority, tc.src)
			}
			if got.Addressed != tc.wantAddressed {
				t.Fatalf("Addressed = %v, want %v", got.Addressed, tc.wantAddressed)
			}
			if tc.wantAddressed && got.Destination != tc.dst {
				t.Fatalf("Destination = %d, want %d", got.Destination, tc.dst)
			}
			if !tc.wantAddressed && got.Destination != framer.BroadcastAddr {
				t.Fatalf("Destination = %d, want broadcast", got.Destination)
			}
		})
	}
}

func TestFastPacketHeaderInverse(t *testing.T) {
	for seq := uint8(0); seq < 8; seq++ {
		for frame := uint8(0); frame < 32; frame++ {
			b := framer.FastPacketHeader(seq, frame)
			gotSeq, gotFrame := framer.FastPacketSeqFrame(b)
			if gotSeq != seq || gotFrame != frame {
				t.Fatalf("round trip (%d,%d) -> %d -> (%d,%d)", seq, frame, b, gotSeq, gotFrame)
			}
		}
	}
}
