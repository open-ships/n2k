package actisense

import (
	"context"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRxPGNMaskAcknowledgements(t *testing.T) {
	for _, mask := range []uint32{RxPGNMaskPGN, RxPGNMaskPDUFormat, RxPGNMaskPDUNibble, RxPGNMaskDataPage, RxPGNMaskDefault, RxPGNMaskNoChange} {
		t.Run(fmt.Sprintf("%08X", mask), func(t *testing.T) {
			effective := mask
			if mask >= RxPGNMaskDefault {
				effective = RxPGNMaskPGN
			}
			requester := &scriptedRequester{response: responseFor(BEMRxPGNEnable, 1, ModelNGT1, RxPGNEnableSet(127250, PGNEnabled, &effective))}
			state, err := NewCommandSet(requester, CommandSetConfig{}).SetRxPGN(context.Background(), 127250, PGNEnabled, &mask)
			require.NoError(t, err)
			assert.Equal(t, effective, state.Mask)
			assert.Equal(t, mask, binary.LittleEndian.Uint32(requester.calls[0].data[5:]))
		})
	}
	for _, mask := range []uint32{0, 0xAABBCCDD, 0x00FF0000} {
		requester := &scriptedRequester{}
		_, err := NewCommandSet(requester, CommandSetConfig{}).SetRxPGN(context.Background(), 127250, PGNEnabled, &mask)
		require.ErrorContains(t, err, "unsupported Rx PGN mask")
		assert.Empty(t, requester.calls)
	}
	mask := RxPGNMaskPGN
	wrong := RxPGNMaskDataPage
	requester := &scriptedRequester{response: responseFor(BEMRxPGNEnable, 1, ModelNGT1, RxPGNEnableSet(127250, PGNEnabled, &wrong))}
	_, err := NewCommandSet(requester, CommandSetConfig{}).SetRxPGN(context.Background(), 127250, PGNEnabled, &mask)
	require.ErrorContains(t, err, "acknowledged Rx PGN state")
}

func TestTxPGNRateAcknowledgements(t *testing.T) {
	for _, rate := range []uint32{0, 1, 65534, 65535, 65536, 0xFFFFFFFE, 0xFFFFFFFF} {
		t.Run(fmt.Sprint(rate), func(t *testing.T) {
			effective := rate
			if rate >= 65535 {
				effective = 1000
			}
			data := append(TxPGNEnableSetFull(127250, PGNEnabled, &effective), 0, 0, 0, 0, 3)
			requester := &scriptedRequester{response: responseFor(BEMTxPGNEnable, 1, ModelNGT1, data)}
			state, err := NewCommandSet(requester, CommandSetConfig{}).SetTxPGN(context.Background(), 127250, PGNEnabled, &rate)
			require.NoError(t, err)
			assert.Equal(t, effective, state.Rate)
			assert.Equal(t, rate, binary.LittleEndian.Uint32(requester.calls[0].data[5:]))
		})
	}
	rate, wrong := uint32(1000), uint32(500)
	requester := &scriptedRequester{response: responseFor(BEMTxPGNEnable, 1, ModelNGT1, append(TxPGNEnableSetFull(127250, PGNEnabled, &wrong), 0, 0, 0, 0, 3))}
	_, err := NewCommandSet(requester, CommandSetConfig{}).SetTxPGN(context.Background(), 127250, PGNEnabled, &rate)
	require.ErrorContains(t, err, "acknowledged Tx PGN state")
}
