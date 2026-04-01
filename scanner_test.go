package n2k

import (
	"context"
	"testing"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CAN frame for PGN 127501 (BinarySwitchBankStatus)
// CAN ID encodes: priority=2, PGN=0x1F20D (127501), source=0
var testFrame127501 = can.Frame{
	ID:     0x09F20D00,
	Length: 8,
	Data:   [8]uint8{0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
}

func TestScannerBasic(t *testing.T) {
	ctx := context.Background()
	s := NewScanner(ctx, Replay([]can.Frame{testFrame127501}), IncludeUnknown())

	got := false
	for s.Next() {
		got = true
		msg := s.Message()
		assert.NotNil(t, msg)
	}
	assert.True(t, got, "should have received at least one message")
	assert.NoError(t, s.Err())
}

func TestScannerNoSources(t *testing.T) {
	ctx := context.Background()
	s := NewScanner(ctx)

	assert.False(t, s.Next())
	assert.Error(t, s.Err())
	assert.Contains(t, s.Err().Error(), "at least one source")
}

func TestScannerContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	s := NewScanner(ctx, Replay([]can.Frame{testFrame127501}))
	for s.Next() {
	}
	// Should terminate
}

func TestScannerDropsUnknownByDefault(t *testing.T) {
	// Frame with PGN that won't decode
	unknownFrame := can.Frame{
		ID:     0x09000000, // PGN 0 or similar that won't have a decoder
		Length: 8,
		Data:   [8]uint8{0, 0, 0, 0, 0, 0, 0, 0},
	}

	ctx := context.Background()
	s := NewScanner(ctx, Replay([]can.Frame{unknownFrame}))

	count := 0
	for s.Next() {
		count++
	}
	assert.Equal(t, 0, count, "unknown PGNs should be dropped by default")
}

func TestScannerIncludeUnknown(t *testing.T) {
	// Frame with PGN that won't decode
	unknownFrame := can.Frame{
		ID:     0x09000000,
		Length: 8,
		Data:   [8]uint8{0, 0, 0, 0, 0, 0, 0, 0},
	}

	ctx := context.Background()
	s := NewScanner(ctx, Replay([]can.Frame{unknownFrame}), IncludeUnknown())

	count := 0
	for s.Next() {
		_, ok := s.Message().(*pgn.UnknownPGN)
		require.True(t, ok, "expected UnknownPGN")
		count++
	}
	assert.Equal(t, 1, count, "unknown PGN should be included when IncludeUnknown is set")
}

func TestScannerWithFilter(t *testing.T) {
	// Filter for PGN 127501 only
	ctx := context.Background()
	s := NewScanner(ctx,
		Replay([]can.Frame{testFrame127501}),
		Filter(`pgn == 127501`),
		IncludeUnknown(),
	)

	count := 0
	for s.Next() {
		count++
	}
	assert.NoError(t, s.Err())
	assert.GreaterOrEqual(t, count, 1, "should receive the matching PGN")
}

func TestScannerWithFilterNoMatch(t *testing.T) {
	ctx := context.Background()
	s := NewScanner(ctx,
		Replay([]can.Frame{testFrame127501}),
		Filter(`pgn == 0`),
		IncludeUnknown(),
	)

	count := 0
	for s.Next() {
		count++
	}
	assert.NoError(t, s.Err())
	assert.Equal(t, 0, count, "nothing should match pgn == 0")
}

func TestScannerWithInvalidFilter(t *testing.T) {
	ctx := context.Background()
	s := NewScanner(ctx,
		Replay([]can.Frame{testFrame127501}),
		Filter(`invalid !!!`),
	)

	assert.False(t, s.Next())
	assert.Error(t, s.Err())
}
