package canbus

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/stretchr/testify/require"
)

func TestSocketCANWireFormat(t *testing.T) {
	local, peer := net.Pipe()
	t.Cleanup(func() { _ = peer.Close() })
	require.NoError(t, peer.SetDeadline(time.Now().Add(5*time.Second)))
	channel := NewSocketCAN(testLogger(), "test")
	channel.open = func(string) (io.ReadWriteCloser, error) { return local, nil }
	verifySocketCANPeer(t, channel, peer)
}

// The peer reads and writes Linux can_frame records independently of the Adapter.
func verifySocketCANPeer(t *testing.T, channel *socketCANChannel, peer io.ReadWriter) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	inbound := make(chan can.Frame, 8)
	go func() { done <- channel.Run(ctx, func(f can.Frame) { inbound <- f }) }()
	t.Cleanup(func() {
		cancel()
		require.NoError(t, channel.Close())
	})
	select {
	case <-channel.Ready():
	case err := <-done:
		t.Fatalf("open channel: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("channel failed to open")
	}

	const id = 0x19f11223
	frame := can.Frame{ID: id, Length: 3, Data: [8]byte{1, 2, 3}}
	written := make(chan error, 1)
	go func() { written <- channel.WriteFrame(frame) }()
	var record [16]byte
	_, err := io.ReadFull(peer, record[:])
	require.NoError(t, err)
	require.NoError(t, <-written)
	require.Equal(t, uint32(0x80000000|id), binary.NativeEndian.Uint32(record[:4]))
	require.Equal(t, []byte{3, 0, 0, 0, 1, 2, 3}, record[4:11])

	for _, wireID := range []uint32{0x123, 0xc0000000 | id, 0x80000000 | id} {
		binary.NativeEndian.PutUint32(record[:4], wireID)
		_, err := peer.Write(record[:])
		require.NoError(t, err)
	}
	select {
	case got := <-inbound:
		require.Equal(t, frame, got)
	case <-time.After(3 * time.Second):
		t.Fatal("no extended data frame received")
	}
	cancel()
	require.ErrorIs(t, <-done, context.Canceled)
}
