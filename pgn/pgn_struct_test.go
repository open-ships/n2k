package pgn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPGNRoundTripVesselHeading(t *testing.T) {
	msg := &VesselHeading{
		Info:      MessageInfo{PGN: 127250},
		Sid:       ptrUint64(5),
		Heading:   ptrUint64(15708),
		Deviation: ptrInt64(-523),
		Variation: ptrInt64(349),
		Reference: ptrUint64(1),
	}

	payload, err := msg.EncodePayload()
	require.NoError(t, err)

	decoded, err := DecodePayload(MessageInfo{PGN: 127250, SourceId: 4}, payload)
	require.NoError(t, err)

	result, ok := decoded.(*VesselHeading)
	require.True(t, ok, "expected *pgn.VesselHeading, got %T", decoded)
	require.Equal(t, uint32(127250), result.PGNNumber())
	require.Equal(t, uint8(4), result.Info.SourceId)
	require.Equal(t, uint64(5), *result.Sid)
	require.Equal(t, uint64(15708), *result.Heading)
	require.Equal(t, int64(-523), *result.Deviation)
	require.Equal(t, int64(349), *result.Variation)
	require.Equal(t, uint64(1), *result.Reference)
}

func TestPGNDecodeUsesVariantMatches(t *testing.T) {
	msg := &BGKeyValueData{
		Info:             MessageInfo{PGN: 130824},
		ManufacturerCode: ptrUint64(381),
		IndustryCode:     ptrUint64(4),
		Key:              ptrUint64(1),
		Length:           ptrUint64(2),
		Value:            []uint8{0x10, 0x20},
	}

	payload, err := msg.EncodePayload()
	require.NoError(t, err)

	decoded, err := DecodePayload(MessageInfo{PGN: 130824}, payload)
	require.NoError(t, err)
	_, ok := decoded.(*BGKeyValueData)
	require.True(t, ok, "expected *pgn.BGKeyValueData, got %T", decoded)
}

func TestPGNMetadata(t *testing.T) {
	infos := PgnInfoLookup[127250]
	require.Len(t, infos, 1)
	require.Equal(t, "VesselHeading", infos[0].Id)
	require.Equal(t, "Vessel Heading", infos[0].Description)
	require.False(t, infos[0].Fast)
	require.NotEmpty(t, infos[0].Fields)

	fd, err := GetFieldDescriptor(127250, 0, 2)
	require.NoError(t, err)
	require.Equal(t, "Heading", fd.Name)
	require.Equal(t, "*uint64", fd.GolangType)
}

func TestDebugDumpPGNStruct(t *testing.T) {
	msg := &VesselHeading{
		Info:    MessageInfo{PGN: 127250, SourceId: 7, Priority: ptrUint8(3)},
		Sid:     ptrUint64(3),
		Heading: ptrUint64(15000),
	}

	result := DebugDumpPGN(msg)
	require.Contains(t, result, "VesselHeading:")
	require.Contains(t, result, "PGN=")
	require.Contains(t, result, "127250")
	require.Contains(t, result, "SourceId=")
	require.Contains(t, result, "Sid=")
	require.Contains(t, result, "Heading=")
}
