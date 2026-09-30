package n2k

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/require"
)

// Keep startup pending until the test releases transport readiness. Frames
// still arrive, as on a busy gateway during its readiness/claim handshake.
type unstartedTestBus struct {
	*mockBus
	ready    chan struct{}
	running  chan struct{}
	runs     atomic.Int32
	claimErr error
}

func newUnstartedTestBus() *unstartedTestBus {
	return &unstartedTestBus{mockBus: newMockBus(), ready: make(chan struct{}), running: make(chan struct{})}
}

func (b *unstartedTestBus) Ready() <-chan struct{} { return b.ready }

func (b *unstartedTestBus) Run(ctx context.Context, handler func(can.Frame)) error {
	b.runs.Add(1)
	close(b.running)
	return b.mockBus.Run(ctx, handler)
}

func (b *unstartedTestBus) WriteFrame(frame can.Frame) error {
	if b.claimErr != nil {
		return b.claimErr
	}
	return b.mockBus.WriteFrame(frame)
}

func awaitStartup(t *testing.T, done <-chan error) error {
	t.Helper()
	select {
	case err := <-done:
		return err
	case <-time.After(5 * time.Second):
		t.Fatal("client lifecycle operation did not finish")
		return nil
	}
}

func TestUnstartedClientReceivesDuringStartup(t *testing.T) {
	bus := newUnstartedTestBus()
	client, err := NewUnstartedClient(context.Background(), WithBus(bus), WithReceiveBuffer(256),
		WithReadyTimeout(5*time.Second), WithClaimTimeout(20*time.Millisecond), WithHeartbeatInterval(0))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	require.Zero(t, bus.runs.Load())
	require.Empty(t, bus.getWritten())
	require.False(t, client.Status().Connected)
	require.False(t, client.Status().Ready)
	require.ErrorIs(t, client.Write(&pgn.VesselHeading{}).Wait(), ErrNotReady)

	scanner := client.Scanner()
	defer func() { _ = scanner.Close() }()
	require.Equal(t, 1, client.Status().ReceiveSubscribers)
	started := make(chan error, 1)
	go func() { started <- client.Start() }()

	// Acknowledge each consumed message so this tests startup accumulation,
	// independently of machine speed or sustained consumer overload. With the
	// blocking constructor no application reader can drain these 1024 messages.
	for i := range 1024 {
		frame := framer.FrameSingle(framer.BuildCANID(127250, 2, 12, 255),
			[]byte{byte(i % 250), 0x5C, 0x3D, 0xFF, 0x7F, 0xFF, 0x7F, 0xFC})
		bus.inbound <- frame
		require.True(t, scanner.Next(), "message %d: %v", i, scanner.Err())
		require.Equal(t, uint64(i%250), *scanner.Message().(*pgn.VesselHeading).Sid)
	}
	require.False(t, client.Status().Ready, "messages must be delivered before startup completes")
	select {
	case err := <-started:
		t.Fatalf("startup finished before readiness: %v", err)
	default:
	}
	close(bus.ready)
	require.NoError(t, awaitStartup(t, started))
	require.True(t, client.Status().Ready)
	require.Equal(t, int32(1), bus.runs.Load())
	require.Len(t, framesWithPGN(bus.getWritten(), framer.PGNISOAddressClaim), 1)
	require.NoError(t, client.Err())
	// Subsequent traffic uses the same subscription and connection.
	bus.inbound <- framer.FrameSingle(framer.BuildCANID(127250, 2, 12, 255),
		[]byte{7, 0x5C, 0x3D, 0xFF, 0x7F, 0xFF, 0x7F, 0xFC})
	require.True(t, scanner.Next())
	require.Equal(t, uint64(7), *scanner.Message().(*pgn.VesselHeading).Sid)
}

func TestUnstartedClientSlowSubscriberStillOverflows(t *testing.T) {
	bus := newUnstartedTestBus()
	client, err := NewUnstartedClient(context.Background(), WithBus(bus), WithReceiveBuffer(2),
		WithReadyTimeout(5*time.Second), WithClaimTimeout(20*time.Millisecond), WithHeartbeatInterval(0))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	scanner := client.Scanner()
	defer func() { _ = scanner.Close() }()
	started := make(chan error, 1)
	go func() { started <- client.Start() }()
	frame := framer.FrameSingle(framer.BuildCANID(127250, 2, 12, 255),
		[]byte{1, 0x5C, 0x3D, 0xFF, 0x7F, 0xFF, 0x7F, 0xFC})
	for range 3 {
		bus.inbound <- frame
	}
	require.Eventually(t, func() bool { return client.Status().ReceiveSubscribers == 0 }, time.Second, time.Millisecond)
	for scanner.Next() {
	}
	require.ErrorIs(t, scanner.Err(), ErrReceiveOverflow)
	require.NoError(t, client.Err(), "overflow must remain local to the subscriber")
	close(bus.ready)
	require.NoError(t, awaitStartup(t, started))
	require.True(t, client.Status().Ready)
}

func TestUnstartedClientStartFailureReachesReader(t *testing.T) {
	bus := newUnstartedTestBus()
	bus.claimErr = errors.New("claim write failed")
	close(bus.ready)
	client, err := NewUnstartedClient(context.Background(), WithBus(bus), WithClaimTimeout(time.Second))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	scanner := client.Scanner()
	defer func() { _ = scanner.Close() }()
	err = client.Start()
	require.ErrorIs(t, err, bus.claimErr)
	require.False(t, scanner.Next())
	require.ErrorIs(t, scanner.Err(), bus.claimErr)
	require.ErrorIs(t, client.Err(), bus.claimErr)
	require.True(t, client.Status().Closed)
	require.ErrorIs(t, client.Start(), ErrClientClosed)
}

func TestUnstartedClientCloseBeforeStart(t *testing.T) {
	bus := newUnstartedTestBus()
	client, err := NewUnstartedClient(context.Background(), WithBus(bus))
	require.NoError(t, err)
	scanner := client.Scanner()
	require.NoError(t, client.Close())
	require.False(t, scanner.Next())
	require.NoError(t, scanner.Err())
	require.ErrorIs(t, client.Start(), ErrClientClosed)
	require.Zero(t, bus.runs.Load())
	require.Empty(t, bus.getWritten())
	require.NoError(t, client.Close())
}

func TestUnstartedClientCloseAndCancelDuringStart(t *testing.T) {
	for _, mode := range []string{"close", "cancel"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			bus := newUnstartedTestBus()
			client, err := NewUnstartedClient(ctx, WithBus(bus), WithReadyTimeout(time.Minute))
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, client.Close()) })
			scanner := client.Scanner()
			defer func() { _ = scanner.Close() }()
			started := make(chan error, 1)
			go func() { started <- client.Start() }()
			select {
			case <-bus.running:
			case <-time.After(time.Second):
				t.Fatal("bus did not start")
			}
			if mode == "close" {
				closed := make(chan error, 1)
				go func() { closed <- client.Close() }()
				require.NoError(t, awaitStartup(t, closed))
			} else {
				cancel()
			}
			require.Error(t, awaitStartup(t, started))
			require.False(t, scanner.Next())
			require.True(t, client.Status().Closed)
			require.Empty(t, bus.getWritten())
		})
	}
}

func TestUnstartedClientConcurrentStart(t *testing.T) {
	bus := newUnstartedTestBus()
	close(bus.ready)
	client, err := NewUnstartedClient(context.Background(), WithBus(bus),
		WithClaimTimeout(20*time.Millisecond), WithHeartbeatInterval(0))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	started := make(chan error, 8)
	for range cap(started) {
		go func() { started <- client.Start() }()
	}
	for range cap(started) {
		require.NoError(t, awaitStartup(t, started))
	}
	require.NoError(t, client.Start())
	require.Equal(t, int32(1), bus.runs.Load())
	require.Len(t, framesWithPGN(bus.getWritten(), framer.PGNISOAddressClaim), 1)
}

func TestUnstartedClientCanceledBeforeStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	bus := newUnstartedTestBus()
	client, err := NewUnstartedClient(ctx, WithBus(bus))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	scanner := client.Scanner()
	cancel()
	require.ErrorIs(t, client.Start(), context.Canceled)
	require.False(t, scanner.Next())
	require.ErrorIs(t, scanner.Err(), context.Canceled)
	require.Zero(t, bus.runs.Load())
}

func TestUnstartedClientReplayRemainsReady(t *testing.T) {
	client, err := NewUnstartedClient(context.Background(), Replay(nil))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	require.True(t, client.Status().Ready)
	require.NoError(t, client.Write(&pgn.VesselHeading{}).Wait())
	require.NoError(t, client.Start())
	require.NoError(t, client.Start())
	require.Len(t, client.WrittenFrames(), 1)
}

func TestUnstartedClientReceivesDuringAddressClaim(t *testing.T) {
	bus := newMockBus()
	// Hold the claim window open independently of machine speed. Close ends
	// the test once traffic has been consumed; no wall-clock sleep is needed.
	client, err := NewUnstartedClient(context.Background(), WithBus(bus),
		WithReceiveBuffer(256), WithClaimTimeout(time.Minute), WithHeartbeatInterval(0))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	scanner := client.Scanner()
	defer func() { _ = scanner.Close() }()
	started := make(chan error, 1)
	go func() { started <- client.Start() }()
	require.Eventually(t, func() bool {
		return len(framesWithPGN(bus.getWritten(), framer.PGNISOAddressClaim)) == 1
	}, time.Second, time.Millisecond)

	for i := range 1024 {
		bus.inbound <- framer.FrameSingle(framer.BuildCANID(127250, 2, 12, 255),
			[]byte{byte(i % 250), 0x5C, 0x3D, 0xFF, 0x7F, 0xFF, 0x7F, 0xFC})
		require.True(t, scanner.Next(), "message %d: %v", i, scanner.Err())
		require.Equal(t, uint64(i%250), *scanner.Message().(*pgn.VesselHeading).Sid)
	}
	require.False(t, client.Status().Ready)
	require.NoError(t, client.Err())
	require.ErrorIs(t, client.Write(&pgn.VesselHeading{}).Wait(), ErrNotReady)
	closed := make(chan error, 1)
	go func() { closed <- client.Close() }()
	require.NoError(t, awaitStartup(t, closed))
	require.Error(t, awaitStartup(t, started))
}
