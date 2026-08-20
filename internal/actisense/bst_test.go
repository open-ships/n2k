package actisense

import (
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func decodeOne(t *testing.T, wire []byte) Datagram {
	t.Helper()
	var got Datagram
	NewParser().Feed(wire, func(d Datagram) { got = d }, func(err DecodeError) { t.Fatal(err) })
	require.NotZero(t, got.ID)
	return got
}

func TestBST93GoldenVector(t *testing.T) {
	payload := []byte{
		2, 0x12, 0xF1, 0x01, 0xFF, 0x23,
		0x78, 0x56, 0x34, 0x12, 3,
		1, 2, 3,
	}
	wire, err := EncodeDatagram(BSTN2KReceive, payload)
	require.NoError(t, err)
	message, ok, err := DecodeMessage(decodeOne(t, wire))
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, uint32(127250), message.PGN)
	assert.Equal(t, uint8(0x23), message.Source)
	assert.True(t, message.HasSource)
	assert.Equal(t, time.Duration(0x12345678)*time.Millisecond, message.Timestamp)
	assert.Equal(t, []byte{1, 2, 3}, message.Data)
}

func TestBST93PDU1UsesCanonicalPGNAndExplicitDestination(t *testing.T) {
	payload := []byte{
		6, 0x31, 0xEA, 0, 0x31, 0x23,
		0, 0, 0, 0, 3,
		0, 0xEE, 0,
	}
	wire, err := EncodeDatagram(BSTN2KReceive, payload)
	require.NoError(t, err)
	message, ok, err := DecodeMessage(decodeOne(t, wire))
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, uint32(59904), message.PGN)
	assert.Equal(t, uint8(0x31), message.Destination)
}

func TestBST94RoundTrip(t *testing.T) {
	wire, err := EncodeMessage94(Message{Priority: 3, PGN: 129029, Destination: 255, Data: []byte{DLE, 2, 3}})
	require.NoError(t, err)
	message, ok, err := DecodeMessage(decodeOne(t, wire))
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, DirectionTransmitted, message.Direction)
	assert.False(t, message.HasSource)
	assert.Equal(t, []byte{DLE, 2, 3}, message.Data)
}

func TestBST94PDU1EncodesCanonicalPS(t *testing.T) {
	wire, err := EncodeMessage94(Message{Priority: 6, PGN: 59904, Destination: 0x31, Data: []byte{0, 0xEE, 0}})
	require.NoError(t, err)
	datagram := decodeOne(t, wire)
	require.Len(t, datagram.Payload, 9)
	assert.Equal(t, byte(0), datagram.Payload[1])
	message, ok, err := DecodeMessage(datagram)
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, uint32(59904), message.PGN)
	assert.Equal(t, uint8(0x31), message.Destination)
}

func TestBST95PDU1RoundTripPreservesSourceAndDestination(t *testing.T) {
	want := can.Frame{
		ID:     framer.BuildCANID(59904, 6, 0x31, 0x22),
		Length: 3,
		Data:   [8]byte{0x00, 0xEE, 0x00},
	}
	wire, err := EncodeCANFrame(want, DirectionTransmitted, 1234*time.Microsecond, 3)
	require.NoError(t, err)
	decoded, ok, err := DecodeCANFrame(decodeOne(t, wire))
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, want, decoded.Frame)
	assert.Equal(t, DirectionTransmitted, decoded.Direction)
	assert.Equal(t, 1234*time.Microsecond, decoded.Timestamp)
}

func TestBST95PublishedGoldenVector(t *testing.T) {
	raw := []byte{0x95, 0x0E, 0x01, 0x20, 0x30, 0x02, 0xF8, 0x09, 0xFF, 0xFC, 0x37, 0x0A, 0x00, 0x10, 0xFF, 0xFF}
	datagram, err := DecodeRaw(raw)
	require.NoError(t, err)
	decoded, ok, err := DecodeCANFrame(datagram)
	require.NoError(t, err)
	require.True(t, ok)
	id := framer.ParseCANID(decoded.Frame.ID)
	assert.Equal(t, uint32(129026), id.PGN)
	assert.Equal(t, uint8(48), id.Source)
	assert.Equal(t, uint8(2), id.Priority)
	assert.Equal(t, 8193*time.Millisecond, decoded.Timestamp)
}

func TestBSTD0MaximumPayloadAndOwnership(t *testing.T) {
	data := make([]byte, MaxNMEA2000Payload)
	for i := range data {
		data[i] = byte(i)
	}
	wire, err := EncodeMessageD0(Message{
		Priority:     3,
		PGN:          129029,
		Destination:  255,
		Source:       0x42,
		Data:         data,
		Timestamp:    42 * time.Second,
		HasTimestamp: true,
		Direction:    DirectionReceived,
		Type:         MessageTransport,
	})
	require.NoError(t, err)
	message, ok, err := DecodeMessage(decodeOne(t, wire))
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, uint32(129029), message.PGN)
	assert.Equal(t, uint8(0x42), message.Source)
	assert.Equal(t, MessageTransport, message.Type)
	assert.Equal(t, 42*time.Second, message.Timestamp)
	require.Equal(t, data, message.Data)
	data[0] ^= 0xFF
	assert.NotEqual(t, data[0], message.Data[0])
}
