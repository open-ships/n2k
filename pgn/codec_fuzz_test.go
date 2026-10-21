package pgn

import (
	"bytes"
	"sort"
	"testing"
)

func FuzzDecodeEncodeMessage(f *testing.F) {
	f.Add(uint32(127250), []byte{0, 0, 0, 0xFF, 0x7F, 0xFF, 0x7F, 0xFC})
	f.Add(uint32(126998), []byte{3, 1, 'a', 3, 1, 'b', 3, 1, 'c'})
	f.Add(uint32(65009), []byte{0, 0x5e, 0xd0, 0xb2, 0, 0x60, 0xfd, 0xff})
	f.Add(uint32(130944), []byte{0x87, 0x98, 0xfe, 2, 1, 0, 4, 2})
	f.Add(uint32(130820), []byte{0xa3, 0x99, 0x0b, 0x80, 1, 2, 0xe0, 0x51, 0x68, 6, 0xc7, 8, 'B', 'o', 'h', 'u', 's', 'l', 0x91, 'n'})
	var known []uint32
	for number := range PgnInfoLookup {
		known = append(known, number)
	}
	sort.Slice(known, func(i, j int) bool { return known[i] < known[j] })
	f.Fuzz(func(t *testing.T, selector uint32, payload []byte) {
		if len(payload) > 1785 {
			return
		}
		pgnNum := selector
		if _, exists := PgnInfoLookup[pgnNum]; !exists {
			pgnNum = known[int(selector%uint32(len(known)))]
		}
		decoded, err := DecodeMessage(MessageInfo{PGN: pgnNum}, payload)
		if err != nil {
			return
		}
		encoded, err := EncodeMessage(decoded)
		if err != nil || !bytes.Equal(encoded, payload) {
			t.Fatalf("successful PGN %d decode lost wire bytes: encode %x, %v; payload %x", pgnNum, encoded, err, payload)
		}
		owned, err := CloneMessage(decoded)
		if err != nil {
			t.Fatal(err)
		}
		copied, err := EncodeMessage(owned)
		if err != nil || !bytes.Equal(copied, payload) {
			t.Fatalf("PGN %d clone lost wire state: %x, %v", pgnNum, copied, err)
		}
	})
}
