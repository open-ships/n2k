package adapter

import (
	"testing"

	"github.com/open-ships/n2k/internal/decoder"
	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
)

// TestPgn127501 verifies end-to-end decoding of PGN 127501 (Binary Switch Bank Status),
// a single-frame non-fast PGN. The test parses a raw CAN log line, creates a packet,
// filters candidates, and then decodes the payload to produce a PgnBinarySwitchBankStatus
// struct. This validates the full pipeline from raw CAN data to typed Go struct.
func TestPgn127501(t *testing.T) {
	raw := "2023-01-21T00:04:17Z,3,127501,224,0,8,00,03,c0,ff,ff,ff,ff,ff"
	f := CanFrameFromRaw(raw)
	pInfo := NewPacketInfo(&f)
	p := decoder.NewPacket(pInfo, f.Data[:])
	assert.NotEmpty(t, p.Candidates)
	p.FilterCandidates()
	assert.Equal(t, 1, len(p.Candidates))
	ret, err := pgn.DecodePayload(p.Info, p.Data)
	assert.Nil(t, err)
	assert.IsType(t, &pgn.PgnBinarySwitchBankStatus{}, ret)
}
