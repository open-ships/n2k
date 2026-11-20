package actisense

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSessionQueuedWriteHonorsContextWithoutInterruptingWriter(t *testing.T) {
	started := make(chan struct{})
	var writes atomic.Int32
	session := NewSession(SessionConfig{Write: func(ctx context.Context, _ []byte) (int, error) {
		writes.Add(1)
		close(started)
		<-ctx.Done()
		return 0, ctx.Err()
	}})
	t.Cleanup(func() { session.Close(nil) })
	firstCtx, cancelFirst := context.WithCancel(context.Background())
	defer cancelFirst()
	first := make(chan error, 1)
	go func() { first <- session.WriteContext(firstCtx, []byte{1}) }()
	<-started

	queuedCtx, cancelQueued := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancelQueued()
	_, err := session.Request(queuedCtx, BEMEcho, []byte{2})
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Equal(t, int32(1), writes.Load())
	require.Zero(t, session.Metrics().BEMInFlight)
	select {
	case <-first:
		t.Fatal("a queued request canceled another writer")
	default:
	}
	cancelFirst()
	select {
	case err := <-first:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("active writer did not stop after cancellation")
	}
}

func TestSessionCloseCancelsActiveAndQueuedWrites(t *testing.T) {
	started := make(chan struct{})
	session := NewSession(SessionConfig{Write: func(ctx context.Context, _ []byte) (int, error) {
		close(started)
		<-ctx.Done()
		return 0, ctx.Err()
	}})
	results := make(chan error, 2)
	go func() { results <- session.Write([]byte{1}) }()
	<-started
	go func() { results <- session.Write([]byte{2}) }()
	terminal := errors.New("transport failed")
	session.Close(terminal)
	for range 2 {
		select {
		case err := <-results:
			require.ErrorIs(t, err, terminal)
		case <-time.After(time.Second):
			t.Fatal("Close did not release a writer")
		}
	}
}

func TestSessionTracesOnlyAcceptedWriteBytes(t *testing.T) {
	failure := errors.New("write stopped")
	for _, accepted := range []int{0, 2} {
		var trace []byte
		session := NewSession(SessionConfig{
			Write: func(context.Context, []byte) (int, error) { return accepted, failure },
			OnWireBytes: func(direction WireDirection, _ time.Time, data []byte) {
				require.Equal(t, WireTransmitted, direction)
				require.NotEmpty(t, data)
				trace = append(trace, data...)
			},
		})
		payload := []byte{1, 2, 3, 4}
		require.ErrorIs(t, session.Write(payload), failure)
		require.Equal(t, payload[:accepted], append([]byte{}, trace...))
		require.Equal(t, uint64(accepted), session.Metrics().TransportWriteBytes)
		require.Equal(t, uint64(1), session.Metrics().TransportWriteErrors)
		session.Close(nil)
	}
}
