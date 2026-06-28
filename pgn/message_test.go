package pgn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeMessage_SelectsVariantForDuplicatePGN(t *testing.T) {
	msg := &GarminColorMode{
		ManufacturerCode: 229,
		IndustryCode:     4,
		UnknownId1:       ptrUint8(222),
		UnknownId2:       ptrUint8(5),
		UnknownId3:       ptrUint8(5),
		UnknownId4:       ptrUint8(5),
		Mode:             GarminColorModeConst(13),
		Color:            GarminColorConst(1),
	}

	payload, err := EncodeMessage(msg)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
}

func TestEncodeMessage_NoEncoder(t *testing.T) {
	msg := &UnknownPGN{}
	msg.Info.PGN = 999999

	_, err := EncodeMessage(msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no encoder registered")
}
