package n2k

import (
	"context"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/open-ships/n2k/internal/actisense"
	internalebl "github.com/open-ships/n2k/internal/ebl"
	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func eblRecord(tag byte, payload []byte) []byte {
	out := []byte{0x1B, 0x01, tag}
	for _, value := range payload {
		out = append(out, value)
		if value == 0x1B {
			out = append(out, 0x1B)
		}
	}
	return append(out, 0x1B, 0x0A)
}

func eblVersionRecord(version uint32) []byte {
	payload := make([]byte, 4)
	binary.LittleEndian.PutUint32(payload, version)
	return eblRecord(internalebl.TagVersion, payload)
}

func eblBSTRaw(t *testing.T, wire []byte) []byte {
	t.Helper()
	var raw []byte
	actisense.NewParser().Feed(wire, func(datagram actisense.Datagram) {
		raw = append([]byte(nil), datagram.Raw...)
	}, func(err actisense.DecodeError) { t.Fatal(err) })
	require.NotEmpty(t, raw)
	return eblRecord(internalebl.TagBSTRawFrame, raw)
}

func TestEBLSourceDecodesBSTRawFrame(t *testing.T) {
	messagePayload := []byte{
		2, 0x12, 0xF1, 0x01, 255, 0x23,
		0, 0, 0, 0, 8,
		1, 0x5C, 0x3D, 0xFF, 0x7F, 0xFF, 0x7F, 0xFC,
	}
	wire, err := actisense.EncodeDatagram(actisense.BSTN2KReceive, messagePayload)
	require.NoError(t, err)
	capture := append(eblVersionRecord(internalebl.CurrentVersion), eblBSTRaw(t, wire)...)
	path := filepath.Join(t.TempDir(), "capture.ebl")
	require.NoError(t, os.WriteFile(path, capture, 0o600))

	var messages []pgn.Message
	for message, receiveErr := range Receive(context.Background(), EBL(path)) {
		require.NoError(t, receiveErr)
		messages = append(messages, message)
	}
	require.Len(t, messages, 1)
	_, ok := messages[0].(*pgn.VesselHeading)
	assert.True(t, ok)
}

func TestEBLSourceFeedsLargeD0AsAssembledMessage(t *testing.T) {
	data := make([]byte, 300)
	for i := range data {
		data[i] = byte(i)
	}
	wire, err := actisense.EncodeMessageD0(actisense.Message{
		Priority: 3, PGN: 130900, Destination: 255, Source: 0x31,
		Data: data, Type: actisense.MessageTransport,
	})
	require.NoError(t, err)
	capture := append(eblVersionRecord(internalebl.CurrentVersion), eblBSTRaw(t, wire)...)
	path := filepath.Join(t.TempDir(), "large-d0.ebl")
	require.NoError(t, os.WriteFile(path, capture, 0o600))

	var got *pgn.UnknownPGN
	for message, receiveErr := range Receive(context.Background(), EBL(path), IncludeUnknown()) {
		require.NoError(t, receiveErr)
		got, _ = message.(*pgn.UnknownPGN)
	}
	require.NotNil(t, got)
	assert.Equal(t, uint32(130900), got.Info.PGN)
	assert.Equal(t, data, got.Data)
}

func TestEBLSourceReportsNewerVersion(t *testing.T) {
	path := filepath.Join(t.TempDir(), "future.ebl")
	require.NoError(t, os.WriteFile(path, eblVersionRecord(internalebl.CurrentVersion+1), 0o600))
	for _, err := range Receive(context.Background(), EBL(path)) {
		require.Error(t, err)
		assert.ErrorIs(t, err, internalebl.ErrUnsupportedVersion)
		return
	}
	t.Fatal("expected unsupported EBL version error")
}

func TestEBLSourceSurfacesUnframedBytesAsOwnedEvidence(t *testing.T) {
	capture := append(eblVersionRecord(internalebl.CurrentVersion), []byte("boot banner\r\n")...)
	path := filepath.Join(t.TempDir(), "unframed.ebl")
	require.NoError(t, os.WriteFile(path, capture, 0o600))

	var observations []Observation
	for observation, err := range Observe(context.Background(), EBL(path)) {
		require.NoError(t, err)
		observations = append(observations, observation)
	}
	require.Len(t, observations, 1)
	assert.Equal(t, ObservationTransportError, observations[0].Kind)
	assert.Equal(t, "actisense-bdtp", observations[0].Protocol)
	assert.Equal(t, []byte("boot banner\r\n"), observations[0].Payload)
}

func TestPipelineHandlesOwnedMessageObservation(t *testing.T) {
	ctx := context.Background()
	data := []byte{1, 2, 3, 4}
	destination := uint8(255)
	var got []pgn.Message
	pipeline, err := newReadPipeline(ctx, config{includeUnknown: true}, func(message pgn.Message) { got = append(got, message) })
	require.NoError(t, err)
	pipeline.HandleObservation(Observation{
		Kind: ObservationMessage, PGN: 130900, Priority: 3, Source: 0x22,
		Destination: &destination, Payload: data,
	})
	data[0] = 0xFF
	require.Len(t, got, 1)
	unknown, ok := got[0].(*pgn.UnknownPGN)
	require.True(t, ok)
	assert.Equal(t, []byte{1, 2, 3, 4}, unknown.Data)
	assert.Equal(t, uint8(0x22), unknown.Info.SourceId)
}

func TestEBLOriginalTimingCanBeCancelled(t *testing.T) {
	version := eblVersionRecord(internalebl.CurrentVersion)
	baseTicks := uint64(11644473600+time.Now().Unix()) * 10_000_000
	firstTime := make([]byte, 8)
	secondTime := make([]byte, 8)
	binary.LittleEndian.PutUint64(firstTime, baseTicks)
	binary.LittleEndian.PutUint64(secondTime, baseTicks+uint64(5*time.Second/(100*time.Nanosecond)))
	wire, err := actisense.EncodeMessageD0(actisense.Message{PGN: 130900, Destination: 255, Source: 1, Data: []byte{1}})
	require.NoError(t, err)
	record := eblBSTRaw(t, wire)
	capture := append(version, eblRecord(internalebl.TagTimeUTC, firstTime)...)
	capture = append(capture, record...)
	capture = append(capture, eblRecord(internalebl.TagTimeUTC, secondTime)...)
	capture = append(capture, record...)
	path := filepath.Join(t.TempDir(), "paced.ebl")
	require.NoError(t, os.WriteFile(path, capture, 0o600))

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		for _, receiveErr := range Receive(ctx, EBL(path, OriginalTiming()), IncludeUnknown()) {
			if receiveErr != nil {
				done <- receiveErr
				return
			}
			cancel()
		}
		done <- nil
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("cancellation did not interrupt EBL timing")
	}
}
