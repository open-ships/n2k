package pgn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeMessage_SelectsVariantForDuplicatePGN(t *testing.T) {
	msg := &GarminColorMode{
		ManufacturerCode: ptrUint64(229),
		IndustryCode:     ptrUint64(4),
		UnknownId1:       ptrUint64(222),
		UnknownId2:       ptrUint64(5),
		UnknownId3:       ptrUint64(5),
		UnknownId4:       ptrUint64(5),
		Mode:             ptrUint64(13),
		Color:            ptrUint64(1),
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
	require.Contains(t, err.Error(), "does not implement pgn.PGN")
}
