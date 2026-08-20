package ebl

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func record(tag byte, payload []byte) []byte {
	out := []byte{escape, start, tag}
	for _, b := range payload {
		out = append(out, b)
		if b == escape {
			out = append(out, escape)
		}
	}
	return append(out, escape, end)
}

func TestReaderWorkedWireTraceVector(t *testing.T) {
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, CurrentVersion)
	ticks := uint64(133_485_408_001_000_000)
	timePayload := make([]byte, 8)
	binary.LittleEndian.PutUint64(timePayload, ticks)
	bst := []byte{0x93, 0x01, 0x42}

	stream := append(record(TagTimeUTC, timePayload), record(TagVersion, version)...)
	stream = append(stream, record(TagDirectionMarker, []byte{0})...)
	stream = append(stream, record(TagBSTRawFrame, bst)...)

	var events []Event
	var warnings []Warning
	err := NewReader().Read(bytes.NewReader(stream), func(event Event) { events = append(events, event) }, func(warning Warning) { warnings = append(warnings, warning) })
	require.NoError(t, err)
	assert.Empty(t, warnings)
	require.Len(t, events, 1)
	assert.True(t, events[0].RawBST)
	assert.Equal(t, DirectionReceived, events[0].Direction)
	assert.Equal(t, bst, events[0].Payload)
	assert.False(t, events[0].Timestamp.IsZero())
}

func TestReaderRawStreamUnstuffsEscape(t *testing.T) {
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, CurrentVersion)
	stream := append(record(TagVersion, version), 1, escape, escape, 2)
	var got Event
	require.NoError(t, NewReader().Read(bytes.NewReader(stream), func(event Event) { got = event }, nil))
	assert.Equal(t, []byte{1, escape, 2}, got.Payload)
}

func TestReaderRejectsNewerVersion(t *testing.T) {
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, CurrentVersion+1)
	err := NewReader().Read(bytes.NewReader(record(TagVersion, version)), nil, nil)
	assert.ErrorIs(t, err, ErrUnsupportedVersion)
}

func TestReaderOversizedTagRecovers(t *testing.T) {
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, CurrentVersion)
	oversized := record(TagDescription, bytes.Repeat([]byte{'x'}, MaxTagPayload+1))
	stream := append(record(TagVersion, version), oversized...)
	stream = append(stream, record(TagBSTRawFrame, []byte{0x93, 0})...)
	var events []Event
	var warnings []Warning
	require.NoError(t, NewReader().Read(bytes.NewReader(stream), func(event Event) { events = append(events, event) }, func(warning Warning) { warnings = append(warnings, warning) }))
	require.Len(t, warnings, 1)
	require.Len(t, events, 1)
}

func TestFiletimeConversion(t *testing.T) {
	const unixEpochTicks = uint64(filetimeEpochOffsetSeconds) * 10_000_000
	assert.Equal(t, time.Unix(0, 0).UTC(), filetimeToTime(unixEpochTicks))
}
