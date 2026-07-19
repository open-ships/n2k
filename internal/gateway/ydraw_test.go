package gateway

import (
	"github.com/brutella/can"
	"github.com/open-ships/n2k/raw"

	"testing"
	"time"

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

func TestParseYDRawObservation_PreservesTransportContext(t *testing.T) {
	observation, ok := ParseYDRawObservation("17:33:21.107 T 09F11201 01")
	require.True(t, ok)
	assert.Equal(t, raw.DirectionTransmitted, observation.Direction)
	assert.True(t, observation.HasTransportTimestamp)
	assert.Equal(t, 17*time.Hour+33*time.Minute+21*time.Second+107*time.Millisecond, observation.TransportTimestamp)
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
		"",                            // empty
		"\r\n",                        // blank
		"YDWG-02 firmware 1.60",       // service banner
		"17:33:21.107 R",              // no ID
		"17:33:21.107 R 09F11201",     // no data bytes
		"17:33:21.107 X 09F11201 01",  // bad direction
		"17:33:21.107 R ZZZZZZZZ 01",  // bad ID hex
		"17:33:21.107 R 209F11201 01", // ID > 29 bits
		"17:33:21.107 R 09F11201 GG",  // bad data hex
		"17:33:21.107 R 09F11201 01 02 03 04 05 06 07 08 09", // > 8 bytes
	}
	for _, line := range cases {
		_, ok := ParseYDRaw(line)
		assert.False(t, ok, "line should be rejected: %q", line)
	}
}

func TestFormatYDRawTX_FullFrame(t *testing.T) {
	frame := can.Frame{ID: 0x19F51323, Length: 8, Data: [8]uint8{0x01, 0x2F, 0x30, 0x70, 0x00, 0x2F, 0x30, 0x70}}
	assert.Equal(t, "19F51323 01 2F 30 70 00 2F 30 70\r\n", string(FormatYDRawTX(frame)))
}

func TestFormatYDRawTX_ShortFrame(t *testing.T) {
	frame := can.Frame{ID: 0x19F51323, Length: 2, Data: [8]uint8{0x01, 0x02}}
	assert.Equal(t, "19F51323 01 02\r\n", string(FormatYDRawTX(frame)))
}

func TestFormatYDRawTX_EchoRoundTrip(t *testing.T) {
	frame := can.Frame{ID: 0x09F80115, Length: 8, Data: [8]uint8{0xA0, 0x7D, 0xE6, 0x18, 0xC0, 0x05, 0xFB, 0xD5}}
	line := string(FormatYDRawTX(frame))
	// The gateway echoes accepted frames back with a timestamp and direction T.
	echoed, ok := ParseYDRaw("17:33:21.108 T " + line[:len(line)-2])
	require.True(t, ok)
	assert.Equal(t, frame, echoed)
}
