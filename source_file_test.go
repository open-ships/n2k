package n2k

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// headingLine is a candump -L line carrying a VesselHeading (PGN 127250)
// frame: SID=1, Heading=15708 (1.5708 rad), deviation/variation null,
// reference=True.
const headingLine = "(1720000000.000000) can0 09F11201#015C3DFF7FFF7FFC"

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "capture.log")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestFileSource_DecodesCandumpLog(t *testing.T) {
	path := writeTempLog(t, headingLine+"\n# a comment\n\n")

	var msgs []pgn.Message
	for msg, err := range Receive(context.Background(), File(path)) {
		require.NoError(t, err)
		msgs = append(msgs, msg)
	}

	require.Len(t, msgs, 1)
	vh, ok := msgs[0].(*pgn.VesselHeading)
	require.True(t, ok, "expected *pgn.VesselHeading, got %T", msgs[0])
	require.NotNil(t, vh.Heading)
	assert.Equal(t, uint64(15708), *vh.Heading)
	assert.Equal(t, uint8(0x01), vh.Info.SourceId)
}

func TestFileSource_OriginalTimingPacesFrames(t *testing.T) {
	log := "(1720000000.000000) can0 09F11201#015C3DFF7FFF7FFC\n" +
		"(1720000000.080000) can0 09F11201#025C3DFF7FFF7FFC\n"
	path := writeTempLog(t, log)

	start := time.Now()
	count := 0
	for _, err := range Receive(context.Background(), File(path, OriginalTiming())) {
		require.NoError(t, err)
		count++
	}
	elapsed := time.Since(start)

	require.Equal(t, 2, count)
	assert.GreaterOrEqual(t, elapsed, 60*time.Millisecond,
		"OriginalTiming should sleep ~80ms between the two frames")
}

func TestFileSource_MissingFile(t *testing.T) {
	var firstErr error
	for _, err := range Receive(context.Background(), File(filepath.Join(t.TempDir(), "nope.log"))) {
		if err != nil {
			firstErr = err
			break
		}
	}
	require.Error(t, firstErr)
}

func TestFileSource_ContextCancelStopsPacing(t *testing.T) {
	// Two frames 10 minutes apart; cancellation must interrupt the sleep.
	log := "(1720000000.000000) can0 09F11201#015C3DFF7FFF7FFC\n" +
		"(1720000600.000000) can0 09F11201#025C3DFF7FFF7FFC\n"
	path := writeTempLog(t, log)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	for range Receive(ctx, File(path, OriginalTiming())) {
	}
	assert.Less(t, time.Since(start), 5*time.Second)
}

func TestNewClient_RejectsReadOnlySource(t *testing.T) {
	path := writeTempLog(t, headingLine+"\n")
	_, err := NewClient(context.Background(), File(path))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read-only")
}
