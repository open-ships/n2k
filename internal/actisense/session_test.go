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
	session := NewSession(SessionConfig{Write: func(ctx context.Context, _ []byte) error {
		writes.Add(1)
		close(started)
		<-ctx.Done()
		return ctx.Err()
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
	session := NewSession(SessionConfig{Write: func(ctx context.Context, _ []byte) error {
		close(started)
		<-ctx.Done()
		return ctx.Err()
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
