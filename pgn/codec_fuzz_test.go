package pgn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func FuzzDecodeEncodeMessage(f *testing.F) {
	f.Add(uint32(127250), []byte{0, 0, 0, 0xFF, 0x7F, 0xFF, 0x7F, 0xFC})
	f.Add(uint32(126998), []byte{2, 1, 'a', 2, 1, 'b', 2, 1, 'c'})
	known := []uint32{59904, 60928, 126208, 126996, 126998, 127250, 129029, 130306}
	f.Fuzz(func(t *testing.T, selector uint32, payload []byte) {
		pgnNum := known[int(selector%uint32(len(known)))]
		var decoded PGN
		require.NotPanics(t, func() {
			var err error
			decoded, err = DecodeMessage(MessageInfo{PGN: pgnNum}, payload)
			if err == nil {
				_, _ = EncodeMessage(decoded)
			}
		})
	})
}
