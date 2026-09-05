package n2k

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func headingProvider(context.Context) pgn.Message {
	h := uint64(15708)
	return &pgn.VesselHeading{Heading: &h}
}

func TestBroadcast_PeriodicTransmission(t *testing.T) {
	c, mb, _ := newCitizenClient(t)

	stop, err := c.Broadcast(20*time.Millisecond, headingProvider)
	require.NoError(t, err)
	defer stop()

	ok := waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 127250)) >= 3
	})
	assert.True(t, ok, "expected at least 3 VesselHeading broadcasts")
}

func TestBroadcast_StopHaltsTransmission(t *testing.T) {
	c, mb, _ := newCitizenClient(t)

	stop, err := c.Broadcast(20*time.Millisecond, headingProvider)
	require.NoError(t, err)
	ok := waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 127250)) >= 1
	})
	require.True(t, ok)

	stop()
	stop() // idempotent
	count := len(framesWithPGN(mb.getWritten(), 127250))
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, count, len(framesWithPGN(mb.getWritten(), 127250)), "no transmissions after stop")
}

func TestBroadcast_NilMessageSkipsTickButKeepsGoing(t *testing.T) {
	c, mb, _ := newCitizenClient(t)

	var ready atomic.Bool
	stop, err := c.Broadcast(20*time.Millisecond, func(ctx context.Context) pgn.Message {
		if !ready.Load() {
			return nil
		}
		return headingProvider(ctx)
	})
	require.NoError(t, err)
	defer stop()

	time.Sleep(100 * time.Millisecond)
	assert.Empty(t, framesWithPGN(mb.getWritten(), 127250), "nil messages must not be transmitted")

	ready.Store(true)
	ok := waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 127250)) >= 1
	})
	assert.True(t, ok, "scheduler must keep ticking through nil messages")
}

func TestBroadcast_CloseStopsAll(t *testing.T) {
	mb := newMockBus()
	c, err := NewClient(context.Background(),
		WithBus(mb),
		WithClaimTimeout(50*time.Millisecond),
		WithHeartbeatInterval(0),
	)
	require.NoError(t, err)

	_, err = c.Broadcast(20*time.Millisecond, headingProvider)
	require.NoError(t, err)
	waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 127250)) >= 1
	})

	require.NoError(t, c.Close())
	count := len(framesWithPGN(mb.getWritten(), 127250))
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, count, len(framesWithPGN(mb.getWritten(), 127250)), "Close must stop broadcasters")
}

func TestBroadcast_SamePGNReplacesEarlier(t *testing.T) {
	c, _, _ := newCitizenClient(t)

	stop1, err := c.BroadcastPGN(127250, time.Hour, headingProvider)
	require.NoError(t, err)
	b1 := c.broadcasterFor(127250)
	require.NotNil(t, b1)

	stop2, err := c.BroadcastPGN(127250, time.Hour, headingProvider)
	require.NoError(t, err)
	b2 := c.broadcasterFor(127250)
	require.NotNil(t, b2)
	assert.NotSame(t, b1, b2, "re-registering a PGN should install a new broadcaster")

	stop1() // stopping the replaced broadcaster must not remove the new one
	assert.NotNil(t, c.broadcasterFor(127250))
	stop2()
	assert.Nil(t, c.broadcasterFor(127250))
}

func TestBroadcast_StopCancelsProviderAndWaitsForExit(t *testing.T) {
	c, _, _ := newCitizenClient(t)
	started, exited := make(chan struct{}), make(chan struct{})
	stop, err := c.BroadcastPGN(127250, time.Hour, func(ctx context.Context) pgn.Message {
		close(started)
		<-ctx.Done()
		close(exited)
		return nil
	})
	require.NoError(t, err)
	<-started
	stopDone := make(chan struct{})
	go func() { stop(); close(stopDone) }()
	select {
	case <-stopDone:
	case <-time.After(time.Second):
		t.Fatal("stop blocked on provider")
	}
	select {
	case <-exited:
	default:
		t.Fatal("stop returned while its provider was still running")
	}
	closeDone := make(chan error, 1)
	go func() { closeDone <- c.Close() }()
	select {
	case err := <-closeDone:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("Close blocked on provider")
	}
}

func TestBroadcast_ProviderPanicIsContained(t *testing.T) {
	c, _, _ := newCitizenClient(t)
	var calls atomic.Int32
	stop, err := c.BroadcastPGN(127250, 10*time.Millisecond, func(ctx context.Context) pgn.Message {
		if calls.Add(1) == 1 {
			panic("bad sensor callback")
		}
		return headingProvider(ctx)
	})
	require.NoError(t, err)
	defer stop()

	require.Eventually(t, func() bool { return calls.Load() >= 2 }, time.Second, time.Millisecond)
}

func TestBroadcastPGN_SkipsMismatchedProviderMessage(t *testing.T) {
	c, mb, _ := newCitizenClient(t)
	var calls atomic.Int32
	stop, err := c.BroadcastPGN(127250, 10*time.Millisecond, func(context.Context) pgn.Message {
		calls.Add(1)
		return &pgn.WaterDepth{}
	})
	require.NoError(t, err)
	defer stop()

	require.Eventually(t, func() bool { return calls.Load() >= 2 }, time.Second, time.Millisecond)
	assert.Empty(t, framesWithPGN(mb.getWritten(), 128267), "a provider must not silently change its declared PGN")
}

func TestBroadcast_CloseOwnsProviderBeforePGNIsKnown(t *testing.T) {
	c, _, _ := newCitizenClient(t)
	started, exited := make(chan struct{}), make(chan struct{})
	_, err := c.Broadcast(time.Hour, func(ctx context.Context) pgn.Message {
		close(started)
		<-ctx.Done()
		close(exited)
		return nil
	})
	require.NoError(t, err)
	<-started
	require.NoError(t, c.Close())
	select {
	case <-exited:
	default:
		t.Fatal("Close lost an unknown-PGN provider")
	}
}

func TestBroadcast_UnknownSchedulesHaveBoundedAdmission(t *testing.T) {
	c, _, _ := newCitizenClient(t)
	for range maxBroadcastWorkers {
		_, err := c.Broadcast(0, headingProvider)
		require.NoError(t, err)
	}
	stop, err := c.Broadcast(0, headingProvider)
	require.ErrorIs(t, err, ErrBroadcastLimit)
	require.Nil(t, stop)
	require.NoError(t, c.Close())
	c.bMu.Lock()
	defer c.bMu.Unlock()
	require.Empty(t, c.broadcastWorkers)
	_, err = c.Broadcast(0, headingProvider)
	require.ErrorIs(t, err, ErrClientClosed)
}

func TestBroadcast_RequiredTriggerSurvivesBusyPeriodicProvider(t *testing.T) {
	c, mb, _ := newCitizenClient(t)
	started, release := make(chan struct{}), make(chan struct{})
	var calls atomic.Int32
	stop, err := c.BroadcastPGN(127250, time.Hour, func(ctx context.Context) pgn.Message {
		if calls.Add(1) == 1 {
			close(started)
			select {
			case <-release:
			case <-ctx.Done():
				return nil
			}
		}
		return headingProvider(ctx)
	})
	require.NoError(t, err)
	defer stop()
	<-started
	require.NoError(t, c.broadcasterFor(127250).trigger(c.ctx))
	close(release)
	require.Eventually(t, func() bool {
		return len(framesWithPGN(mb.getWritten(), 127250)) == 2
	}, time.Second, time.Millisecond)
	require.Equal(t, int32(2), calls.Load())
}

func TestBroadcast_RequiredQueueOverflowIsTerminalAndBounded(t *testing.T) {
	c, _, _ := newCitizenClient(t)
	started := make(chan struct{})
	var calls atomic.Int32
	_, err := c.BroadcastPGN(127250, time.Hour, func(ctx context.Context) pgn.Message {
		calls.Add(1)
		close(started)
		<-ctx.Done()
		return nil
	})
	require.NoError(t, err)
	<-started
	c.applyRequestedInterval(127250, nil)
	c.applyRequestedInterval(127250, nil)
	require.ErrorIs(t, c.Err(), ErrBroadcastQueueFull)
	require.NoError(t, c.Close())
	require.Equal(t, int32(1), calls.Load(), "triggers must not start concurrent providers")
}

func TestBroadcast_AdmissionRacingCloseLeavesNoWorkers(t *testing.T) {
	c, _, _ := newCitizenClient(t)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for range 32 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			stop, err := c.Broadcast(0, headingProvider)
			if err == nil {
				stop()
			} else if !errors.Is(err, ErrClientClosed) {
				t.Errorf("unexpected broadcast admission error: %v", err)
			}
		}()
	}
	close(start)
	require.NoError(t, c.Close())
	wg.Wait()
	c.bMu.Lock()
	defer c.bMu.Unlock()
	require.Empty(t, c.broadcastWorkers)
}

func TestBroadcast_RequiredResponseCannotCrossEpoch(t *testing.T) {
	c, bus, _ := newCitizenClient(t)
	started, canceled := make(chan struct{}), make(chan struct{})
	stop, err := c.BroadcastPGN(127250, 0, func(ctx context.Context) pgn.Message {
		close(started)
		<-ctx.Done()
		close(canceled)
		// Even a provider returning a final value while cancellation unwinds
		// must not transmit an old request's response in the next epoch.
		return headingProvider(ctx)
	})
	require.NoError(t, err)
	defer stop()
	c.applyRequestedInterval(127250, nil)
	<-started
	c.handleConnectionChange(false, 1)
	select {
	case <-canceled:
	case <-time.After(time.Second):
		t.Fatal("required provider retained a departed request epoch")
	}
	c.handleConnectionChange(true, 2)
	require.Eventually(t, func() bool { return !c.Status().Rejoining }, time.Second, time.Millisecond)
	stop()
	require.Empty(t, framesWithPGN(bus.getWritten(), 127250))
	require.NoError(t, c.Err())
}
