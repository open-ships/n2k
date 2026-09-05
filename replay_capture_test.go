package n2k

import (
	"context"
	"testing"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplayCaptureBoundedRecentWindowAndStatus(t *testing.T) {
	client, err := NewClient(context.Background(), Replay(nil))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	initial := client.Status()
	assert.Equal(t, ReplayFrameCapacity, initial.ReplayFrameCapacity)
	assert.Zero(t, initial.ReplayFramesRetained)
	assert.Zero(t, initial.ReplayFramesDropped)

	// More than two rotations distinguishes a chronological ring snapshot
	// from its physical storage order and bounds the allocated storage too.
	const written = 2*ReplayFrameCapacity + 17
	for i := range written {
		require.NoError(t, client.writeFrame(can.Frame{ID: uint32(i), Length: 1, Data: [8]byte{byte(i)}}))
	}
	frames := client.WrittenFrames()
	require.Len(t, frames, ReplayFrameCapacity)
	for i, frame := range frames {
		expected := written - ReplayFrameCapacity + i
		require.Equal(t, uint32(expected), frame.ID)
		require.Equal(t, byte(expected), frame.Data[0])
	}
	status := client.Status()
	assert.Equal(t, ReplayFrameCapacity, status.ReplayFramesRetained)
	assert.Equal(t, ReplayFrameCapacity, status.ReplayFrameCapacity)
	assert.Equal(t, uint64(written-ReplayFrameCapacity), status.ReplayFramesDropped)
	client.mu.Lock()
	capacity := cap(client.writtenFrames)
	client.mu.Unlock()
	assert.Equal(t, ReplayFrameCapacity, capacity)

	frames[0].ID = 0
	frames[0].Data[0] = 0
	assert.NotEqual(t, frames[0], client.WrittenFrames()[0], "returned frames must not alias capture storage")

	// Full capture retains diagnostics by evicting an old entry; it must not
	// cause a successful replay transmission to report an artificial failure.
	require.NoError(t, client.Write(&pgn.VesselHeading{}).Wait())
	assert.Equal(t, uint64(written-ReplayFrameCapacity+1), client.Status().ReplayFramesDropped)
	assert.Equal(t, uint64(1), client.Status().ApplicationWritesCompleted)
}
