package gateway

import (
	"bytes"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/raw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActisenseCANASCIIIndependentSpecificationVector(t *testing.T) {
	observation, err := ParseActisenseCANASCIIObservation("17:33:21.107 R 19F51323 01 2F 30 70 00 2F 30 70")
	require.NoError(t, err)
	require.NotNil(t, observation.Frame)
	assert.Equal(t, uint32(0x19F51323), observation.Frame.ID)
	assert.Equal(t, uint8(8), observation.Frame.Length)
	assert.Equal(t, [8]byte{0x01, 0x2F, 0x30, 0x70, 0x00, 0x2F, 0x30, 0x70}, observation.Frame.Data)
	assert.Equal(t, raw.DirectionReceived, observation.Direction)
	assert.Equal(t, 17*time.Hour+33*time.Minute+21*time.Second+107*time.Millisecond, observation.TransportTimestamp)
	id := framer.ParseCANID(observation.Frame.ID)
	assert.Equal(t, uint32(128275), id.PGN)
	assert.Equal(t, uint8(6), id.Priority)
	assert.Equal(t, uint8(35), id.Source)
}

func TestActisenseCANASCIISecondsOnlyAndZeroLength(t *testing.T) {
	observation, err := ParseActisenseCANASCIIObservation("00:00:01 T 18EEFF2A")
	require.NoError(t, err)
	require.NotNil(t, observation.Frame)
	assert.Zero(t, observation.Frame.Length)
	assert.Equal(t, raw.DirectionTransmitted, observation.Direction)
	assert.Equal(t, time.Second, observation.TransportTimestamp)
}

func TestActisenseCANASCIIEncodeGolden(t *testing.T) {
	frame := can.Frame{ID: 0x19F51323, Length: 2, Data: [8]byte{0x01, 0x2F}}
	line, err := FormatActisenseCANASCII(frame, raw.DirectionTransmitted, 0)
	require.NoError(t, err)
	assert.Equal(t, "00:00:00.000 T 19F51323 01 2F\r\n", string(line))
}

func TestActisenseN2KASCIIIndependentSpecificationVector(t *testing.T) {
	observation, err := ParseActisenseN2KASCIIObservation("A173321.107 23FF7 1F513 012F3070002F30709F")
	require.NoError(t, err)
	assert.Equal(t, raw.KindMessage, observation.Kind)
	assert.Equal(t, uint32(128275), observation.PGN)
	assert.Equal(t, uint8(7), observation.Priority)
	assert.Equal(t, uint8(0x23), observation.Source)
	require.NotNil(t, observation.Destination)
	assert.Equal(t, uint8(0xFF), *observation.Destination)
	assert.Equal(t, []byte{0x01, 0x2F, 0x30, 0x70, 0x00, 0x2F, 0x30, 0x70, 0x9F}, observation.Payload)
	assert.Equal(t, 17*time.Hour+33*time.Minute+21*time.Second+107*time.Millisecond, observation.TransportTimestamp)
}

func TestActisenseN2KASCIIEncodeAddressedAndMaximumPayload(t *testing.T) {
	payload := bytes.Repeat([]byte{0xAB}, maxActisenseN2KPayload)
	line, err := FormatActisenseN2KASCII(126208, 3, 0x16, 0x23, payload, 2*time.Second)
	require.NoError(t, err)
	require.LessOrEqual(t, len(line), maxActisenseN2KASCIILine)
	observation, err := ParseActisenseN2KASCIIObservation(string(line))
	require.NoError(t, err)
	assert.Equal(t, uint32(126208), observation.PGN)
	assert.Equal(t, uint8(0x16), observation.Source)
	require.NotNil(t, observation.Destination)
	assert.Equal(t, uint8(0x23), *observation.Destination)
	assert.Equal(t, payload, observation.Payload)
}

func TestActisenseCANASCIILinesIgnoreBinaryAndRemainBounded(t *testing.T) {
	assembler := &actisenseCANASCIILines{}
	var observations []raw.Observation
	assembler.feed(bytes.Repeat([]byte{'x'}, maxActisenseASCIIFrameLine+10), func(value raw.Observation) {
		observations = append(observations, value)
	})
	assembler.feed([]byte("\n17:33:21 R 19F51323 01\r\n"), func(value raw.Observation) {
		observations = append(observations, value)
	})
	assert.LessOrEqual(t, len(assembler.buf), maxActisenseASCIIFrameLine)
	require.Len(t, observations, 1)
	assert.Equal(t, uint8(1), observations[0].Frame.Length)
}

func TestActisenseMode6DemuxDivertsBinaryBeforeLineBuffering(t *testing.T) {
	binaryFrame, err := actisense.EncodeDatagram(actisense.BSTBEMResponse, make([]byte, 12))
	require.NoError(t, err)
	asciiLine := []byte("17:33:21.107 R 19F51323 01\r\n")
	stream := append(append([]byte(nil), binaryFrame...), asciiLine...)

	lines := &actisenseCANASCIILines{}
	var (
		datagrams    []actisense.Datagram
		observations []raw.Observation
	)
	session := actisense.NewSession(actisense.SessionConfig{
		OnDatagram: func(datagram actisense.Datagram) { datagrams = append(datagrams, datagram) },
		OnUnframed: func(data []byte) {
			lines.feed(data, func(observation raw.Observation) { observations = append(observations, observation) })
		},
	})
	require.NoError(t, session.Run(bytes.NewReader(stream)))
	require.Len(t, datagrams, 1)
	require.Len(t, observations, 1)
	assert.Equal(t, uint32(0x19F51323), observations[0].Frame.ID)
	assert.Equal(t, uint64(len(asciiLine)), session.Metrics().UnframedBytes)
}
