package n2k

import (
	"bytes"
	"testing"
	"time"

	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/internal/ebl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActisenseEBLTraceClassifiesWireEvidence(t *testing.T) {
	var output bytes.Buffer
	writer, err := NewEBLWriter(&output, WithEBLStartTime(time.Unix(0, 0)))
	require.NoError(t, err)
	trace, err := NewActisenseEBLTrace(writer)
	require.NoError(t, err)

	valid, err := actisense.EncodeDatagram(actisense.BSTN2KReceive, []byte{0x42})
	require.NoError(t, err)
	received := append([]byte("boot\r\n"), valid...)
	trace.trace(actisense.WireReceived, time.Unix(1, 0), received)
	trace.trace(actisense.WireTransmitted, time.Unix(2, 0), valid)
	require.NoError(t, trace.Flush())

	var events []ebl.Event
	require.NoError(t, ebl.NewReader().Read(bytes.NewReader(output.Bytes()), func(event ebl.Event) {
		events = append(events, event)
	}, func(warning ebl.Warning) { t.Fatal(warning) }))
	require.Len(t, events, 3)
	assert.Equal(t, []byte("boot\r\n"), events[0].Payload)
	assert.False(t, events[0].RawBST)
	assert.True(t, events[1].RawBST)
	assert.Equal(t, []byte{actisense.BSTN2KReceive, 1, 0x42}, events[1].Payload)
	assert.Equal(t, valid, events[2].Payload)
	assert.Equal(t, ebl.DirectionTransmitted, events[2].Direction)
}

func TestActisenseEBLTraceRetainsInvalidFrameExactly(t *testing.T) {
	var output bytes.Buffer
	writer, err := NewEBLWriter(&output, WithEBLStartTime(time.Unix(0, 0)))
	require.NoError(t, err)
	trace, err := NewActisenseEBLTrace(writer)
	require.NoError(t, err)
	wire, err := actisense.EncodeDatagram(actisense.BSTN2KReceive, []byte{1, 2})
	require.NoError(t, err)
	wire[len(wire)-3]++
	trace.trace(actisense.WireReceived, time.Unix(1, 0), wire)
	require.NoError(t, trace.Flush())

	var events []ebl.Event
	require.NoError(t, ebl.NewReader().Read(bytes.NewReader(output.Bytes()), func(event ebl.Event) {
		events = append(events, event)
	}, nil))
	require.Len(t, events, 1)
	assert.False(t, events[0].RawBST)
	assert.Equal(t, wire, events[0].Payload)
}
