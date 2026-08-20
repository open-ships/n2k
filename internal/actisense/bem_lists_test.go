package actisense

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPGNEnableRequestGoldenVectors(t *testing.T) {
	assert.Equal(t, []byte{0x10, 0xF0, 0x01, 0x00}, RxPGNEnableGet(126992))
	mask := uint32(0xAABBCCDD)
	assert.Equal(t, []byte{0x10, 0xF0, 0x01, 0x00, 0x01, 0xDD, 0xCC, 0xBB, 0xAA}, RxPGNEnableSet(126992, PGNEnabled, &mask))
	rate := uint32(1000)
	assert.Equal(t, []byte{0x10, 0xF0, 0x01, 0x00, 0x02, 0xE8, 0x03, 0x00, 0x00}, TxPGNEnableSetFull(126992, PGNRespondMode, &rate))
}

func TestRxF2RequiresStandardAndProprietaryForNGX(t *testing.T) {
	accumulator := &rxF2Accumulator{expectProprietary: true}
	proprietary := responseFor(BEMRxPGNEnableListF2, 2, ModelNGX1, []byte{
		7, 0x03, 0x11, 0x00, 0x00,
		1, 0x05,
		1, 0x02,
	})
	done, err := accumulator.feed(proprietary)
	require.NoError(t, err)
	assert.False(t, done, "proprietary-first must not truncate the standard list")
	standard := responseFor(BEMRxPGNEnableListF2, 1, ModelNGX1, []byte{
		7, 0x01, 0x11, 0x00, 0x00,
		2, 0, 2,
		3, 1,
		4, 0,
	})
	done, err = accumulator.feed(standard)
	require.NoError(t, err)
	assert.True(t, done)
	assert.Equal(t, []RxPGNListEntry{{Index: 3, Mask: 1}, {Index: 4, Mask: 0}}, accumulator.result.Entries)
	assert.Equal(t, []uint32{ProprietaryDP0Base, ProprietaryDP0Base + 2, ProprietaryDP1Base + 1}, accumulator.result.Proprietary.EnabledPGNs)
}

func TestTxF2WaitsForMissingStandardSlots(t *testing.T) {
	accumulator := &txF2Accumulator{expectProprietary: true}
	first := responseFor(BEMTxPGNEnableListF2, 1, ModelNGX1, []byte{
		8, 0x02, 0x11, 0x00, 0x00,
		2, 1, 1,
		9, 3, 0xE8, 0x03,
	})
	done, err := accumulator.feed(first)
	require.NoError(t, err)
	assert.False(t, done)
	proprietary := responseFor(BEMTxPGNEnableListF2, 2, ModelNGX1, []byte{
		8, 0x03, 0x11, 0x00, 0x00,
		0,
		0,
	})
	done, err = accumulator.feed(proprietary)
	require.NoError(t, err)
	assert.False(t, done, "a proprietary reply cannot hide standard slot zero")
	zeroth := responseFor(BEMTxPGNEnableListF2, 1, ModelNGX1, []byte{
		8, 0x02, 0x11, 0x00, 0x00,
		2, 0, 1,
		7, 6, 0xFF, 0xFF,
	})
	done, err = accumulator.feed(zeroth)
	require.NoError(t, err)
	assert.True(t, done)
	assert.Equal(t, uint16(0xFFFF), accumulator.result.Entries[0].RateMS)
}

func TestF2OlderModelCompletesWithoutProprietaryVariant(t *testing.T) {
	accumulator := &rxF2Accumulator{expectProprietary: false}
	done, err := accumulator.feed(responseFor(BEMRxPGNEnableListF2, 1, ModelNGT1, []byte{
		1, 0x01, 0x11, 0x00, 0x00,
		0, 0, 0,
	}))
	require.NoError(t, err)
	assert.True(t, done)
}

func TestSupportedPGNAccumulatorDoesNotCompleteAcrossAHole(t *testing.T) {
	accumulator := &supportedPGNAccumulator{}
	done, err := accumulator.feed(supportedPGNPart{
		transferID: 1, database: 2, total: 3, first: 1,
		entries: []SupportedPGN{{Index: 1, PGN: 126992}, {Index: 2, PGN: 127250}},
	})
	require.NoError(t, err)
	assert.False(t, done)
	done, err = accumulator.feed(supportedPGNPart{
		transferID: 1, database: 2, total: 3, first: 0,
		entries: []SupportedPGN{{Index: 0, PGN: 60928}},
	})
	require.NoError(t, err)
	assert.True(t, done)
}

func TestDecodePGNListParametersGolden(t *testing.T) {
	response := responseFor(BEMPGNListParameters, 1, ModelNGX1, []byte{
		10, 0, 2, 0, 3, 0,
		20, 0, 4, 0, 5, 0,
		0, 1,
	})
	parameters, err := DecodePGNListParameters(response)
	require.NoError(t, err)
	assert.Equal(t, uint16(10), parameters.RxMaximum)
	assert.Equal(t, uint16(5), parameters.TxActive)
	assert.True(t, parameters.RxSynchronized())
	assert.False(t, parameters.TxSynchronized())
}
