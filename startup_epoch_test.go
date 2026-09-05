package n2k

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/stretchr/testify/require"
)

// Startup must never publish an identity claimed on a departed connection.
type startupEpochBus struct {
	*mockBus
	observer  func(bool, uint64)
	reconnect bool
	flip      bool
	connect   <-chan struct{}
	claimOnce sync.Once
}

func (b *startupEpochBus) SetConnectionObserver(observer func(bool, uint64)) {
	b.observer = observer
}

func (b *startupEpochBus) Run(ctx context.Context, handler func(can.Frame)) error {
	if b.connect != nil {
		select {
		case <-b.connect:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	b.observer(true, 1)
	return b.mockBus.Run(ctx, handler)
}

func (b *startupEpochBus) WriteFrame(frame can.Frame) error {
	if err := b.mockBus.WriteFrame(frame); err != nil {
		return err
	}
	if b.flip && framer.ParseCANID(frame.ID).PGN == framer.PGNISOAddressClaim {
		b.claimOnce.Do(func() {
			b.observer(false, 1)
			if b.reconnect {
				b.observer(true, 2)
			}
		})
	}
	return nil
}

func TestStartupClaimCannotCrossConnectionEpoch(t *testing.T) {
	for _, reconnect := range []bool{false, true} {
		name := "disconnect"
		if reconnect {
			name = "reconnect"
		}
		t.Run(name, func(t *testing.T) {
			bus := &startupEpochBus{mockBus: newMockBus(), reconnect: reconnect, flip: true}
			client, err := NewClient(context.Background(), WithBus(bus), WithClaimTimeout(30*time.Millisecond), WithHeartbeatInterval(0))
			if client != nil {
				require.NoError(t, client.Close())
			}
			require.ErrorIs(t, err, ErrEpochChanged)
			require.Len(t, framesWithPGN(bus.getWritten(), framer.PGNISOAddressClaim), 1)
		})
	}
}

func TestStartupWaitsForKnownConnectionEpoch(t *testing.T) {
	connect := make(chan struct{})
	bus := &startupEpochBus{mockBus: newMockBus(), connect: connect}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	type outcome struct {
		client *Client
		err    error
	}
	done := make(chan outcome, 1)
	go func() {
		client, err := NewClient(ctx, WithBus(bus), WithReadyTimeout(time.Second), WithClaimTimeout(20*time.Millisecond), WithHeartbeatInterval(0))
		done <- outcome{client: client, err: err}
	}()
	select {
	case result := <-done:
		if result.client != nil {
			require.NoError(t, result.client.Close())
		}
		t.Fatalf("construction ended before connection notification: %v", result.err)
	case <-time.After(40 * time.Millisecond):
	}
	require.Empty(t, bus.getWritten(), "no claim may be stamped with unknown connection epoch")
	close(connect)
	select {
	case result := <-done:
		require.NoError(t, result.err)
		require.NoError(t, result.client.Close())
		require.Len(t, framesWithPGN(bus.getWritten(), framer.PGNISOAddressClaim), 1)
	case <-time.After(time.Second):
		t.Fatal("construction did not resume after connection notification")
	}
}

func TestStartupConnectionWaitIsBounded(t *testing.T) {
	bus := &startupEpochBus{mockBus: newMockBus(), connect: make(chan struct{})}
	started := time.Now()
	client, err := NewClient(context.Background(), WithBus(bus), WithReadyTimeout(20*time.Millisecond), WithClaimTimeout(20*time.Millisecond), WithHeartbeatInterval(0))
	require.Nil(t, client)
	require.ErrorContains(t, err, "ready timeout")
	require.Less(t, time.Since(started), time.Second)
	require.Empty(t, bus.getWritten())
}
