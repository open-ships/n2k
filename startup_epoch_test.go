package n2k

import (
	"context"
	"sync"
	"testing"
	"testing/synctest"
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

type replacementClaimBus struct {
	*mockBus
	hold    uint8
	release chan struct{}
}

func (b *replacementClaimBus) WriteFrameContext(ctx context.Context, frame can.Frame) error {
	id := framer.ParseCANID(frame.ID)
	if id.PGN == framer.PGNISOAddressClaim && id.Source == b.hold {
		select {
		case <-b.release:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return b.WriteFrame(frame)
}

func TestStartupContentionFollowsSuccessfulReplacementClaim(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		bus := &replacementClaimBus{mockBus: newMockBus(), hold: 250, release: make(chan struct{})}
		client, err := NewUnstartedClient(t.Context(), WithBus(bus), WithName(DeviceName{IdentityNumber: 10}),
			WithClaimTimeout(time.Second), WithHeartbeatInterval(0))
		require.NoError(t, err)
		defer func() { _ = client.Close() }()
		started := make(chan error, 1)
		go func() { started <- client.Start() }()
		synctest.Wait()
		firstClaims := framesWithPGN(bus.getWritten(), framer.PGNISOAddressClaim)
		require.Len(t, firstClaims, 1)
		time.Sleep(150 * time.Millisecond)
		bus.inbound <- claimFrame(DeviceName{IdentityNumber: 1}, 251)
		synctest.Wait()
		require.Equal(t, uint8(250), client.Status().Address)
		require.ErrorIs(t, client.writeAddressClaim(firstClaims[0]), ErrEpochChanged)
		time.Sleep(addressContentionWindow)
		require.False(t, client.Status().Ready, "the replacement claim is still in the writer")

		close(bus.release)
		synctest.Wait()
		time.Sleep(addressContentionWindow - time.Millisecond)
		require.False(t, client.Status().Ready)
		// A second change invalidates the first replacement's pending timer.
		bus.inbound <- claimFrame(DeviceName{IdentityNumber: 2}, 250)
		synctest.Wait()
		require.Equal(t, uint8(249), client.Status().Address)
		time.Sleep(addressContentionWindow - time.Millisecond)
		require.False(t, client.Status().Ready)
		time.Sleep(time.Millisecond)
		synctest.Wait()
		require.NoError(t, <-started)
		require.True(t, client.Status().Ready)
	})
}

func TestSlowInitialClaimObservesPostWriteContention(t *testing.T) {
	for _, tc := range []struct {
		name        string
		deadline    time.Duration
		writeDelay  time.Duration
		observation time.Duration
		wantErr     error
	}{
		{"full contention misses deadline", time.Second, 900 * time.Millisecond, 100 * time.Millisecond, context.DeadlineExceeded},
		{"short override starts at write", 100 * time.Millisecond, 80 * time.Millisecond, 100 * time.Millisecond, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				bus := &replacementClaimBus{mockBus: newMockBus(), hold: 251, release: make(chan struct{})}
				client, err := NewUnstartedClient(t.Context(), WithBus(bus), WithClaimTimeout(tc.deadline),
					WithHeartbeatInterval(0), WithWriteTimeout(2*time.Second))
				require.NoError(t, err)
				defer func() { _ = client.Close() }()
				started := make(chan error, 1)
				go func() { started <- client.Start() }()
				synctest.Wait()
				time.Sleep(tc.writeDelay)
				close(bus.release)
				synctest.Wait()
				time.Sleep(tc.observation - time.Millisecond)
				require.False(t, client.Status().Ready)
				time.Sleep(time.Millisecond)
				require.ErrorIs(t, <-started, tc.wantErr)
				require.Equal(t, tc.wantErr == nil, client.Status().Ready)
			})
		})
	}
}

func TestStartupDeadlineRejectsUnsettledAddress(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		bus := newMockBus()
		client, err := NewUnstartedClient(t.Context(), WithBus(bus), WithName(DeviceName{IdentityNumber: 10}),
			WithClaimTimeout(time.Second), WithHeartbeatInterval(0))
		require.NoError(t, err)
		defer func() { _ = client.Close() }()
		started := make(chan error, 1)
		go func() { started <- client.Start() }()
		synctest.Wait()
		time.Sleep(900 * time.Millisecond)
		bus.inbound <- claimFrame(DeviceName{IdentityNumber: 1}, 251)
		synctest.Wait()
		require.Equal(t, uint8(250), client.Status().Address)
		time.Sleep(100 * time.Millisecond)
		require.ErrorIs(t, <-started, context.DeadlineExceeded)
		require.False(t, client.Status().Ready)
	})
}
