package n2k

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/require"
)

type unwindingWriteBus struct {
	*mockBus
	entered  chan struct{}
	canceled chan struct{}
	release  chan struct{}
	once     sync.Once
}

func (b *unwindingWriteBus) WriteFrameContext(ctx context.Context, frame can.Frame) error {
	if framer.ParseCANID(frame.ID).PGN == 127250 {
		close(b.entered)
		<-ctx.Done()
		close(b.canceled)
		// Model a physical record completing while cancellation interrupts I/O.
		// Its acceptance must be reflected before the public result completes.
		<-b.release
	}
	return b.WriteFrame(frame)
}

func (b *unwindingWriteBus) Close() error {
	b.once.Do(func() { close(b.release) })
	return b.mockBus.Close()
}

func TestCanceledWriteJoinsPhysicalCompletion(t *testing.T) {
	bus := &unwindingWriteBus{mockBus: newMockBus(), entered: make(chan struct{}), canceled: make(chan struct{}), release: make(chan struct{})}
	client, err := NewClient(context.Background(), WithBus(bus), WithClaimTimeout(20*time.Millisecond), WithHeartbeatInterval(0))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := client.WriteContext(ctx, &pgn.VesselHeading{})
	<-bus.entered
	cancel()
	<-bus.canceled
	select {
	case <-result.Done():
		t.Fatal("write completed before its owned physical write returned")
	case <-time.After(20 * time.Millisecond):
	}
	bus.once.Do(func() { close(bus.release) })
	err = result.Wait()
	require.ErrorIs(t, err, context.Canceled)
	var partial *WriteError
	require.ErrorAs(t, err, &partial)
	require.True(t, partial.TransmissionUncertain)
	require.Equal(t, uint64(1), partial.CompletedRecords)
	// A canceled context-aware write is not itself a transport failure.
	require.NoError(t, client.wire.send(context.Background(), func(context.Context) error { return nil }))
	require.NoError(t, client.Err())
}

func TestWireAdmissionCannotOutliveWorker(t *testing.T) {
	client, _, _ := newCitizenClient(t)
	require.NoError(t, client.Close())
	ctx := context.WithValue(context.Background(), writeStampKey{}, writeStamp{})
	done := make(chan error, 1)
	go func() {
		done <- client.wire.send(ctx, func(context.Context) error { return errors.New("must not run") })
	}()
	select {
	case err := <-done:
		require.ErrorIs(t, err, ErrClientClosed)
	case <-time.After(time.Second):
		t.Fatal("wire work was admitted after its worker exited")
	}
}

func TestWireCallerDeadlineVersusPhysicalTimeout(t *testing.T) {
	for _, callerDeadline := range []bool{false, true} {
		name := "physical timeout is terminal"
		if callerDeadline {
			name = "caller deadline is not terminal"
		}
		t.Run(name, func(t *testing.T) {
			bus := &reliabilityBus{mockBus: newMockBus(), entered: make(chan struct{}), release: make(chan struct{})}
			writeTimeout := 20 * time.Millisecond
			if callerDeadline {
				writeTimeout = time.Second
			}
			client, err := NewClient(context.Background(), WithBus(bus), WithClaimTimeout(20*time.Millisecond), WithWriteTimeout(writeTimeout), WithHeartbeatInterval(0))
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, client.Close()) })
			ctx := context.Background()
			if callerDeadline {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, 20*time.Millisecond)
				defer cancel()
			}
			require.ErrorIs(t, client.WriteContext(ctx, &pgn.VesselHeading{}).Wait(), context.DeadlineExceeded)
			if callerDeadline {
				require.NoError(t, client.wire.send(context.Background(), func(context.Context) error { return nil }))
				require.NoError(t, client.Err())
			} else {
				require.Eventually(t, func() bool { return errors.Is(client.Err(), context.DeadlineExceeded) }, time.Second, time.Millisecond)
			}
		})
	}
}

func TestWriteContextPreservesInheritedEpochStamp(t *testing.T) {
	client, bus, _ := newCitizenClient(t)
	client.mu.Lock()
	oldContext, stop := client.writeContextLocked(context.Background(), true)
	// Hold the old context active to model an epoch cancellation callback
	// that has not been scheduled yet. Identity checks cannot rely on it.
	client.claimEpoch++
	client.mu.Unlock()
	defer stop()
	require.NoError(t, oldContext.Err())
	require.ErrorIs(t, client.WriteContext(oldContext, &pgn.VesselHeading{}).Wait(), ErrEpochChanged)
	require.ErrorIs(t, client.writeProtocolContext(oldContext, "stale fixture", protocolRequired, &pgn.VesselHeading{}).Wait(), ErrEpochChanged)
	require.Empty(t, framesWithPGN(bus.getWritten(), 127250))
	require.NoError(t, client.Err())
}
