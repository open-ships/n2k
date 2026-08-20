package ebl

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriterPreambleUsesSDKOrderAndVersionVector(t *testing.T) {
	var output bytes.Buffer
	_, err := NewWriter(&output, WriterConfig{StartTime: time.Unix(0, 0).UTC()})
	require.NoError(t, err)

	versionVector := []byte{0x1B, 0x01, 0x01, 0xEA, 0x03, 0x00, 0x00, 0x1B, 0x0A}
	assert.True(t, bytes.Contains(output.Bytes(), versionVector))
	assert.Equal(t, byte(TagTimeUTC), output.Bytes()[2], "TimeUTC must precede Version")
}

func TestWriterReaderRoundTripRawBSTAndUnframedBytes(t *testing.T) {
	var output bytes.Buffer
	writer, err := NewWriter(&output, WriterConfig{StartTime: time.Unix(1, 200).UTC(), Description: "wire"})
	require.NoError(t, err)
	bst := []byte{0x93, 0x02, 0x10, 0x42}
	rawBytes := []byte{0x10, 0x02, 0x93, 0x00, 0x6D, 0x10, 0x03}
	require.NoError(t, writer.WriteRawBST(time.Unix(2, 300).UTC(), DirectionReceived, bst))
	require.NoError(t, writer.WriteRawBytes(time.Unix(3, 400).UTC(), DirectionTransmitted, rawBytes))

	var events []Event
	require.NoError(t, NewReader().Read(bytes.NewReader(output.Bytes()), func(event Event) {
		events = append(events, event)
	}, func(warning Warning) { t.Fatal(warning) }))
	require.Len(t, events, 2)
	assert.True(t, events[0].RawBST)
	assert.Equal(t, DirectionReceived, events[0].Direction)
	assert.Equal(t, bst, events[0].Payload)
	assert.False(t, events[1].RawBST)
	assert.Equal(t, DirectionTransmitted, events[1].Direction)
	assert.Equal(t, rawBytes, events[1].Payload)
	assert.Equal(t, uint64(8), writer.Metrics().Records)
}

type switchErrorWriter struct {
	output bytes.Buffer
	fail   bool
}

func (w *switchErrorWriter) Write(data []byte) (int, error) {
	if w.fail {
		return 0, errors.New("capture disk full")
	}
	return w.output.Write(data)
}

func TestWriterRetainsFirstOutputError(t *testing.T) {
	output := &switchErrorWriter{}
	writer, err := NewWriter(output, WriterConfig{StartTime: time.Unix(0, 0)})
	require.NoError(t, err)
	output.fail = true
	err = writer.WriteDescription("fails")
	require.EqualError(t, err, "capture disk full")
	assert.EqualError(t, writer.Err(), "capture disk full")
	assert.EqualError(t, writer.WriteDescription("still fails"), "capture disk full")
	assert.Equal(t, uint64(1), writer.Metrics().Errors)
}

func TestWriterRejectsTimestampBeforeFILETIMEEpoch(t *testing.T) {
	_, err := NewWriter(&bytes.Buffer{}, WriterConfig{StartTime: time.Date(1600, 1, 1, 0, 0, 0, 0, time.UTC)})
	assert.Error(t, err)
}
