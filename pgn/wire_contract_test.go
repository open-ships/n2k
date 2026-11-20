package pgn

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCodecPayloadLimit(t *testing.T) {
	// PGN 129797 has an eight-byte header followed by the binary field.
	m := &AisBinaryBroadcastMessage{BinaryData: make([]byte, maxPayloadBytes-8)}
	wire, err := m.EncodePayload()
	require.NoError(t, err)
	require.Len(t, wire, maxPayloadBytes)
	_, err = DecodeMessage(MessageInfo{PGN: m.PGNNumber()}, wire)
	require.NoError(t, err)
	m.BinaryData = append(m.BinaryData, 1)
	_, err = m.EncodePayload()
	require.ErrorIs(t, err, errPayloadTooLarge)
	tooLong := append(wire, 0)
	_, err = DecodeMessage(MessageInfo{PGN: m.PGNNumber()}, tooLong)
	require.ErrorIs(t, err, errPayloadTooLarge)
	require.ErrorIs(t, m.DecodePayload(tooLong), errPayloadTooLarge)
}

func TestAISPreservesExplicitBitCountAfterEdit(t *testing.T) {
	wire := bytes.Clone(aisDgnssBroadcastBinaryMessagePayload)
	wire[15] = 13
	var m AisDgnssBroadcastBinaryMessage
	require.NoError(t, m.DecodePayload(wire))
	replay, err := m.EncodePayload()
	require.NoError(t, err)
	require.Equal(t, wire, replay)
	m.RepeatIndicator = new(uint64(2))
	edited, err := m.EncodePayload()
	require.NoError(t, err)
	require.Equal(t, uint16(13), binary.LittleEndian.Uint16(edited[15:17]))
	require.Equal(t, []byte{0xde, 0x0d}, edited[17:])
	m.BinaryData = append(m.BinaryData, 0x42)
	grown, err := m.EncodePayload()
	require.NoError(t, err)
	require.Equal(t, uint16(24), binary.LittleEndian.Uint16(grown[15:17]))
}

func TestSchemaNumericNulls(t *testing.T) {
	claim := &IsoAddressClaim{DeviceInstanceUpper: new(uint64(31))}
	wire, err := claim.EncodePayload()
	require.NoError(t, err)
	require.Equal(t, byte(31), wire[4]>>3)
	require.NoError(t, claim.DecodePayload(wire))
	require.Equal(t, uint64(31), *claim.DeviceInstanceUpper)

	angle := &GarminAutopilotTurnAngleOrder{TurnAngleOrder: new(uint64(65535))}
	_, ok := angle.TurnAngleOrderValue()
	require.True(t, ok)
	_, _, ok, err = PhysicalValue(angle, 10)
	require.NoError(t, err)
	require.True(t, ok)
	wire, err = angle.EncodePayload()
	require.NoError(t, err)
	require.NoError(t, angle.DecodePayload(wire))
	require.Equal(t, uint64(65535), *angle.TurnAngleOrder)

	gnss := &GnssPositionData{}
	wire, err = gnss.EncodePayload()
	require.NoError(t, err)
	require.Equal(t, uint64(math.MaxInt64), binary.LittleEndian.Uint64(wire[23:]))
	require.NoError(t, gnss.DecodePayload(wire))
	require.Nil(t, gnss.Latitude)
	require.Nil(t, gnss.Longitude)
	require.Nil(t, gnss.Altitude)
	gnss.Altitude = new(int64(math.MaxInt64))
	_, ok = gnss.AltitudeValue()
	require.False(t, ok)

	var datum UserDatum
	floatWire := make([]byte, 16)
	binary.LittleEndian.PutUint32(floatWire[12:], 0x7fc00001)
	require.NoError(t, datum.DecodePayload(floatWire))
	require.Nil(t, datum.RotationInX)
	_, err = json.Marshal(datum)
	require.NoError(t, err)
	wire, err = datum.EncodePayload()
	require.NoError(t, err)
	require.Equal(t, floatWire, wire)
}

func TestControlTextDecoding(t *testing.T) {
	for _, tt := range []struct {
		name string
		wire []byte
		text string
	}{
		{"UTF16LE", []byte{10, 0, 0x41, 0, 0x3d, 0xd8, 0xe0, 0xde, 0, 0}, "A🛠"},
		{"UTF16BE BOM", []byte{8, 0, 0xfe, 0xff, 0, 0x41, 0xff, 0xff}, "A"},
		{"Latin1", []byte{6, 1, 'B', 0xe6, 'k', 0}, "Bæk"},
		{"UTF8", []byte{5, 1, 'B', 0xc3, 0xa6}, "Bæ"},
		{"Long Latin1", append([]byte{255, 1}, bytes.Repeat([]byte{0xe6}, 253)...), strings.Repeat("æ", 253)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var msg ConfigurationInformation
			require.NoError(t, msg.DecodePayload(tt.wire))
			require.Equal(t, tt.text, msg.InstallationDescription1)
			wire, err := msg.EncodePayload()
			require.NoError(t, err)
			require.Equal(t, tt.wire, wire)
			msg.ManufacturerInformation = "updated"
			wire, err = msg.EncodePayload()
			require.NoError(t, err)
			require.NoError(t, msg.DecodePayload(wire))
			require.Equal(t, tt.text, msg.InstallationDescription1)
			require.Equal(t, "updated", msg.ManufacturerInformation)
		})
	}
}

func TestSimnetCapturedCompactAlarm(t *testing.T) {
	// Reassembled 1DFF280A frames in testdata/sample.log.
	wire := append([]byte{0x41, 0x9f, 0x6d, 0x00, 0x28, 0x01}, []byte("No rudder data. Steering is unreliable")...)
	msg, err := DecodeMessage(MessageInfo{PGN: 130856}, wire)
	require.NoError(t, err)
	alarm := msg.(*SimnetAlarmMessage)
	require.Equal(t, "No rudder data. Steering is unreliable", alarm.Text)
	require.Equal(t, uint64(109), *alarm.AlarmId)
	require.Equal(t, uint64(40), *alarm.B)
	require.Equal(t, uint64(1), *alarm.C)
	replay, err := alarm.EncodePayload()
	require.NoError(t, err)
	require.Equal(t, wire, replay)
	alarm.Text = "updated"
	edited, err := alarm.EncodePayload()
	require.NoError(t, err)
	require.Equal(t, append([]byte{0x41, 0x9f, 0x6d, 0x00, 9, 1}, []byte("updated")...), edited)
}

func TestMetadataFieldLookupUsesSpecificAndUnambiguousVariants(t *testing.T) {
	field, err := GetFieldDescriptor(59392, 0, 2)
	require.NoError(t, err)
	require.Equal(t, "Group Function", field.Name)
	_, err = GetFieldDescriptor(126208, 0, 5)
	require.ErrorContains(t, err, "ambiguous")
	_, err = GetFieldDescriptor(61184, ManufacturerCodeFuruno, 5)
	require.ErrorContains(t, err, "ambiguous")
	field, err = GetFieldDescriptor(61184, ManufacturerCodeFuruno, 1)
	require.NoError(t, err)
	require.Equal(t, "Manufacturer Code", field.Name)
}

func TestProprietaryDispatchSelectsOnlyWireMatch(t *testing.T) {
	data := bytes.Clone(bgKeyValuePayload)
	msg, err := DecodeMessage(MessageInfo{PGN: 130824}, data)
	require.NoError(t, err)
	require.IsType(t, &BGKeyValueData{}, msg)
	data[0]--
	_, err = DecodeMessage(MessageInfo{PGN: 130824}, data)
	require.ErrorContains(t, err, "no matching PGN struct")
}

func TestNearCapacityTextSurvivesReplayAndEdits(t *testing.T) {
	// Seventeen zone names (each 100 bytes of Latin-1) expand
	// beyond the transport ceiling if converted wholesale to UTF-8.
	wire := []byte{1, 17, 17}
	for id := 1; id <= 17; id++ {
		wire = append(wire, byte(id), 102, 1)
		wire = append(wire, bytes.Repeat([]byte{0xe6}, 100)...)
	}
	var zones SupportedZoneData
	require.NoError(t, zones.DecodePayload(wire))
	got, err := zones.EncodePayload()
	require.NoError(t, err)
	require.Equal(t, wire, got)
	zones.FirstZoneId = new(uint64(2))
	got, err = zones.EncodePayload()
	require.NoError(t, err)
	wire[0] = 2
	require.Equal(t, wire, got)
}
