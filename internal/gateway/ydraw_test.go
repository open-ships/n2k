package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseYDRaw_ReceivedFrame(t *testing.T) {
	frame, ok := ParseYDRaw("17:33:21.107 R 09F11201 01 5C 3D FF 7F FF 7F FC")
	require.True(t, ok)
	assert.Equal(t, uint32(0x09F11201), frame.ID)
	assert.Equal(t, uint8(8), frame.Length)
	assert.Equal(t, [8]uint8{0x01, 0x5C, 0x3D, 0xFF, 0x7F, 0xFF, 0x7F, 0xFC}, frame.Data)
}

func TestParseYDRaw_TransmittedFrame(t *testing.T) {
	frame, ok := ParseYDRaw("17:33:21.108 T 09F11201 01 5C 3D FF 7F FF 7F FC")
	require.True(t, ok)
	assert.Equal(t, uint32(0x09F11201), frame.ID)
}

func TestParseYDRaw_ShortFrame(t *testing.T) {
	frame, ok := ParseYDRaw("00:00:00.000 R 18EEFF01 AA BB")
	require.True(t, ok)
	assert.Equal(t, uint8(2), frame.Length)
	assert.Equal(t, uint8(0xAA), frame.Data[0])
}

func TestParseYDRaw_CarriageReturnAndSpaces(t *testing.T) {
	frame, ok := ParseYDRaw("17:33:21.107  R  09F11201  01 5C 3D FF 7F FF 7F FC\r")
	require.True(t, ok)
	assert.Equal(t, uint8(8), frame.Length)
}

func TestParseYDRaw_Rejects(t *testing.T) {
	cases := []string{
		"",                              // empty
		"\r\n",                          // blank
		"YDWG-02 firmware 1.60",         // service banner
		"17:33:21.107 R",                // no ID
		"17:33:21.107 R 09F11201",       // no data bytes
		"17:33:21.107 X 09F11201 01",    // bad direction
		"17:33:21.107 R ZZZZZZZZ 01",    // bad ID hex
		"17:33:21.107 R 209F11201 01",   // ID > 29 bits
		"17:33:21.107 R 09F11201 GG",    // bad data hex
		"17:33:21.107 R 09F11201 01 02 03 04 05 06 07 08 09", // > 8 bytes
	}
	for _, line := range cases {
		_, ok := ParseYDRaw(line)
		assert.False(t, ok, "line should be rejected: %q", line)
	}
}
