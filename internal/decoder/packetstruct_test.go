package decoder

import (
	"testing"

	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHandler struct {
	received []pgn.Message
}

func (m *mockHandler) HandleStruct(v pgn.Message) {
	m.received = append(m.received, v)
}

func TestNew(t *testing.T) {
	ps := New()
	assert.NotNil(t, ps)
}

func TestSetOutput(t *testing.T) {
	ps := New()
	handler := &mockHandler{}
	ps.SetOutput(handler)
	assert.NotNil(t, ps.handler)
}

func TestDecode_ValidPGN(t *testing.T) {
	ps := New()
	handler := &mockHandler{}
	ps.SetOutput(handler)

	payload, err := (&pgn.PgnVesselHeading{
		Sid:     ptrUint64(1),
		Heading: ptrUint64(15000),
	}).EncodePayload()
	require.NoError(t, err)

	ps.Decode(Packet{
		Info: pgn.MessageInfo{PGN: 127250, SourceId: 1},
		Data: payload,
	})

	require.Len(t, handler.received, 1)
	vh, ok := handler.received[0].(*pgn.PgnVesselHeading)
	require.True(t, ok, "expected *pgn.PgnVesselHeading, got %T", handler.received[0])
	assert.Equal(t, uint32(127250), vh.Info.PGN)
	assert.Equal(t, uint8(1), vh.Info.SourceId)
	require.NotNil(t, vh.Heading)
	assert.Equal(t, uint64(15000), *vh.Heading)
}

func TestDecode_ErrorFallsToUnknown(t *testing.T) {
	ps := New()
	handler := &mockHandler{}
	ps.SetOutput(handler)

	ps.Decode(Packet{
		Info: pgn.MessageInfo{PGN: 999999, SourceId: 1},
		Data: []uint8{0x01, 0x02, 0x03},
	})

	require.Len(t, handler.received, 1)
	u, ok := handler.received[0].(*pgn.UnknownPGN)
	require.True(t, ok, "expected *pgn.UnknownPGN, got %T", handler.received[0])
	assert.Equal(t, uint32(999999), u.Info.PGN)
	require.NotNil(t, u.Reason)
}

func TestDecode_PgnStructFromMetadata(t *testing.T) {
	ps := New()
	handler := &mockHandler{}
	ps.SetOutput(handler)

	payload, err := (&pgn.PgnElectricDriveStatusDynamic{
		InverterMotorIdentifier: ptrUint64(1),
		OperatingMode:           ptrUint64(2),
		MotorTemperature:        ptrUint64(3000),
		InverterTemperature:     ptrUint64(3001),
		CoolantTemperature:      ptrUint64(3002),
		GearTemperature:         ptrUint64(3003),
		ShaftTorque:             ptrUint64(500),
	}).EncodePayload()
	require.NoError(t, err)

	ps.Decode(Packet{
		Info: pgn.MessageInfo{PGN: 127490, SourceId: 1},
		Data: payload,
	})

	require.Len(t, handler.received, 1)
	msg, ok := handler.received[0].(*pgn.PgnElectricDriveStatusDynamic)
	require.True(t, ok, "expected *pgn.PgnElectricDriveStatusDynamic, got %T", handler.received[0])
	assert.Equal(t, uint32(127490), msg.Info.PGN)
	require.NotNil(t, msg.InverterMotorIdentifier)
	assert.Equal(t, uint64(1), *msg.InverterMotorIdentifier)
}

func TestDecode_NoHandler(t *testing.T) {
	ps := New()
	payload, err := (&pgn.PgnVesselHeading{
		Heading: ptrUint64(15000),
	}).EncodePayload()
	require.NoError(t, err)

	assert.NotPanics(t, func() {
		ps.Decode(Packet{
			Info: pgn.MessageInfo{PGN: 127250, SourceId: 1},
			Data: payload,
		})
	})
}

func ptrUint64(v uint64) *uint64 {
	return &v
}
