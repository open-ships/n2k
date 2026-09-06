package actisense

import (
	"context"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func legacyPart(command byte, values ...uint32) BEMResponse {
	data := make([]byte, 1+4*len(values))
	data[0] = byte(len(values))
	for i, value := range values {
		binary.LittleEndian.PutUint32(data[1+4*i:], value)
	}
	return responseFor(command, 0, ModelNGT1, data)
}

func TestLegacyF1ResponseTrains(t *testing.T) {
	r := &scriptedRequester{responses: []BEMResponse{
		legacyPart(BEMRxPGNEnableListF1, 127250, 129025),
		legacyPart(BEMRxPGNEnableListF1, RxPGNMaskPGN, RxPGNMaskDataPage),
	}}
	commands := NewCommandSet(r, CommandSetConfig{})
	rx, err := commands.GetRxPGNEnableListF1(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 2, rx.PartsReceived)
	assert.Equal(t, []RxPGNListF1Entry{{127250, RxPGNMaskPGN}, {129025, RxPGNMaskDataPage}}, rx.Entries)
	assert.Equal(t, byte(0x48), r.calls[0].command)
	assert.Empty(t, r.calls[0].data)
	r.responses = []BEMResponse{
		legacyPart(BEMTxPGNEnableListF1, 127250, 129025),
		legacyPart(BEMTxPGNEnableListF1, 1000, 500),
		legacyPart(BEMTxPGNEnableListF1, 2000, 0),
		responseFor(BEMTxPGNEnableListF1, 0, ModelNGT1, []byte{2, 3, 6}),
	}
	tx, err := commands.GetTxPGNEnableListF1(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 4, tx.PartsReceived)
	assert.Equal(t, []TxPGNListF1Entry{{127250, 1000, 2000, 3}, {129025, 500, 0, 6}}, tx.Entries)
	r.stopAfterReply, r.multiErr = 2, context.DeadlineExceeded
	tx, err = commands.GetTxPGNEnableListF1(context.Background())
	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, 2, tx.PartsReceived)
	assert.Equal(t, uint32(1000), tx.Entries[0].RateMS)
	assert.Zero(t, tx.Entries[0].TimeoutMS)
}

func TestLegacyF1RejectsMalformedOrMixedTrains(t *testing.T) {
	for _, second := range []BEMResponse{
		legacyPart(BEMRxPGNEnableListF1),
		legacyPart(BEMTxPGNEnableListF1, RxPGNMaskPGN),
		responseFor(BEMRxPGNEnableListF1, 0, ModelNGX1, []byte{1, 0, 0xFF, 0xFF, 3}),
		responseFor(BEMRxPGNEnableListF1, 0, ModelNGT1, []byte{1, 0}),
		legacyPart(BEMRxPGNEnableListF1, 0xFFFFFFFF),
		legacyPart(BEMRxPGNEnableListF1, make([]uint32, 51)...),
	} {
		r := &scriptedRequester{responses: []BEMResponse{legacyPart(BEMRxPGNEnableListF1, 127250), second}}
		result, err := NewCommandSet(r, CommandSetConfig{}).GetRxPGNEnableListF1(context.Background())
		require.Error(t, err)
		assert.Equal(t, 1, result.PartsReceived)
		assert.Equal(t, uint32(127250), result.Entries[0].PGN)
	}
	r := &scriptedRequester{responses: []BEMResponse{legacyPart(BEMRxPGNEnableListF1), legacyPart(BEMRxPGNEnableListF1)}}
	result, err := NewCommandSet(r, CommandSetConfig{}).GetRxPGNEnableListF1(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 2, result.PartsReceived)
	assert.Empty(t, result.Entries)
}

func TestPortDuplicateDeleteWireContract(t *testing.T) {
	r := &scriptedRequester{response: responseFor(BEMPortDuplicateDelete, 1, ModelNGT1, []byte{3, 0, 1, 0})}
	commands := NewCommandSet(r, CommandSetConfig{})
	settings, err := commands.GetPortDuplicateDelete(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []PortDuplicateDelete{0, 1, 0}, settings)
	_, err = commands.SetPortDuplicateDelete(context.Background(), []PortDuplicateDelete{0, 1, 255})
	require.NoError(t, err)
	assert.Equal(t, []byte{0, 1, 255}, r.calls[1].data)
	assert.Len(t, r.calls, 2, "no implicit EEPROM command")
	_, err = commands.SetPortDuplicateDelete(context.Background(), []PortDuplicateDelete{1, 1, 255})
	require.Error(t, err)
	for _, invalid := range [][]PortDuplicateDelete{nil, {2}, make([]PortDuplicateDelete, 224)} {
		_, err = commands.SetPortDuplicateDelete(context.Background(), invalid)
		require.Error(t, err)
	}
	assert.Len(t, r.calls, 3)
}
