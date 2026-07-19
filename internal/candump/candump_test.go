package candump

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_LineWithTimestamp(t *testing.T) {
	frame, ts, ok := Parse("(1720000000.500000) can0 09F50E7F#00FFFFFFFFFFFFFF")
	require.True(t, ok)
	assert.Equal(t, uint32(0x09F50E7F), frame.ID)
	assert.Equal(t, uint8(8), frame.Length)
	assert.Equal(t, [8]uint8{0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, frame.Data)
	assert.Equal(t, time.Unix(1720000000, 500000000).UTC(), ts.UTC())
}

func TestParseRecord_PreservesNetwork(t *testing.T) {
	record, ok := ParseRecord("(1720000000.500000) can7 09F50E7F#00FFFFFFFFFFFFFF")
	require.True(t, ok)
	assert.Equal(t, "can7", record.Network)
	assert.Equal(t, time.Unix(1720000000, 500000000), record.Timestamp)
}

func TestParse_LineWithoutTimestamp(t *testing.T) {
	frame, ts, ok := Parse("09F10D01#0102030405060708")
	require.True(t, ok)
	assert.Equal(t, uint32(0x09F10D01), frame.ID)
	assert.Equal(t, uint8(8), frame.Length)
	assert.True(t, ts.IsZero(), "missing timestamp should yield zero time")
}

func TestParse_ShortPayload(t *testing.T) {
	frame, _, ok := Parse("(1.000000) can0 18EEFF01#AABB")
	require.True(t, ok)
	assert.Equal(t, uint8(2), frame.Length)
	assert.Equal(t, uint8(0xAA), frame.Data[0])
	assert.Equal(t, uint8(0xBB), frame.Data[1])
}

func TestParse_Rejects(t *testing.T) {
	cases := []string{
		"",                                  // empty
		"# comment line",                    // comment
		"(1.0) can0 123#R",                  // RTR frame
		"(1.0) can0 123##1AABB",             // CAN FD frame
		"(1.0) can0 ZZZ#AABB",               // bad ID
		"(1.0) can0 123#GGHH",               // bad hex payload
		"(1.0) can0 123#010203040506070809", // > 8 bytes
		"no frame data here",                // no # field
	}
	for _, line := range cases {
		_, _, ok := Parse(line)
		assert.False(t, ok, "line should be rejected: %q", line)
	}
}

func TestParse_BadTimestampStillParsesFrame(t *testing.T) {
	frame, ts, ok := Parse("(notatime) can0 09F10D01#01")
	require.True(t, ok)
	assert.Equal(t, uint8(1), frame.Length)
	assert.True(t, ts.IsZero())
}
