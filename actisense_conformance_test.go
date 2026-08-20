package n2k

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/internal/ebl"
	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type publicActisenseGoldenCorpus struct {
	Remote []struct {
		ID      string `json:"id"`
		Command byte   `json:"command"`
		Data    string `json:"data_hex"`
		Payload string `json:"pgn_payload_hex"`
	} `json:"remote_bem_requests"`
	ASCII struct {
		CAN string `json:"can"`
		N2K string `json:"n2k"`
	} `json:"ascii_vectors"`
	EBL struct {
		Version string `json:"version_1002_hex"`
	} `json:"ebl_records"`
}

func loadPublicActisenseGoldenCorpus(t *testing.T) publicActisenseGoldenCorpus {
	t.Helper()
	contents, err := os.ReadFile("conformance/actisense-golden.json")
	require.NoError(t, err)
	var corpus publicActisenseGoldenCorpus
	require.NoError(t, json.Unmarshal(contents, &corpus))
	return corpus
}

func mustDecodeActisenseGoldenHex(t *testing.T, value string) []byte {
	t.Helper()
	decoded, err := hex.DecodeString(value)
	require.NoError(t, err)
	return decoded
}

func TestConformanceActisenseRemoteGoldenEnvelope(t *testing.T) {
	for _, vector := range loadPublicActisenseGoldenCorpus(t).Remote {
		t.Run(vector.ID, func(t *testing.T) {
			inner, err := encodeActisenseRemoteCommand(vector.Command, mustDecodeActisenseGoldenHex(t, vector.Data))
			require.NoError(t, err)
			payload := append([]byte{0x11, 0x99}, inner...)
			assert.Equal(t, mustDecodeActisenseGoldenHex(t, vector.Payload), payload)
		})
	}
}

func TestConformanceActisenseASCIIGoldenVectors(t *testing.T) {
	corpus := loadPublicActisenseGoldenCorpus(t)
	canObservation, err := ParseActisenseCANASCII(corpus.ASCII.CAN)
	require.NoError(t, err)
	require.NotNil(t, canObservation.Frame)
	canLine, err := EncodeActisenseCANASCII(*canObservation.Frame, canObservation.Direction, canObservation.TransportTimestamp)
	require.NoError(t, err)
	assert.Equal(t, corpus.ASCII.CAN, strings.TrimSpace(string(canLine)))

	n2kObservation, err := ParseActisenseN2KASCII(corpus.ASCII.N2K)
	require.NoError(t, err)
	require.NotNil(t, n2kObservation.Destination)
	n2kLine, err := EncodeActisenseN2KASCII(
		n2kObservation.PGN,
		n2kObservation.Priority,
		n2kObservation.Source,
		*n2kObservation.Destination,
		n2kObservation.Payload,
		n2kObservation.TransportTimestamp,
	)
	require.NoError(t, err)
	assert.Equal(t, corpus.ASCII.N2K, strings.TrimSpace(string(n2kLine)))
}

func TestConformanceActisenseEBLGoldenVersionRecord(t *testing.T) {
	corpus := loadPublicActisenseGoldenCorpus(t)
	var output bytes.Buffer
	_, err := ebl.NewWriter(&output, ebl.WriterConfig{StartTime: time.Unix(0, 0).UTC()})
	require.NoError(t, err)
	assert.True(t, bytes.Contains(output.Bytes(), mustDecodeActisenseGoldenHex(t, corpus.EBL.Version)))
	assert.Equal(t, byte(ebl.TagTimeUTC), output.Bytes()[2])
}

func TestConformanceActisenseRemotePayloadUsesVariableLengthPGN(t *testing.T) {
	destination := uint8(35)
	priority := uint8(3)
	message := &actisenseRemoteMessage{
		info:    actisenseMessageInfo(destination, priority),
		payload: []byte{0x11, 0x99, actisense.BSTBEMCommand, 0x01, actisense.BEMOperatingMode},
	}
	payload, err := message.EncodePayload()
	require.NoError(t, err)
	assert.Len(t, payload, 5)
}

func actisenseMessageInfo(destination, priority uint8) pgn.MessageInfo {
	return pgn.MessageInfo{PGN: actisenseRemotePGN, Priority: &priority, TargetId: &destination}
}
