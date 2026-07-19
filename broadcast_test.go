package n2k

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func headingProvider() pgn.Message {
	h := uint64(15708)
	return &pgn.VesselHeading{Heading: &h}
}

func TestBroadcast_PeriodicTransmission(t *testing.T) {
	c, mb, _ := newCitizenClient(t)

	stop := c.Broadcast(20*time.Millisecond, headingProvider)
	defer stop()

	ok := waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 127250)) >= 3
	})
	assert.True(t, ok, "expected at least 3 VesselHeading broadcasts")
}

func TestBroadcast_StopHaltsTransmission(t *testing.T) {
	c, mb, _ := newCitizenClient(t)

	stop := c.Broadcast(20*time.Millisecond, headingProvider)
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
	stop := c.Broadcast(20*time.Millisecond, func() pgn.Message {
		if !ready.Load() {
			return nil
		}
		return headingProvider()
	})
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

	_ = c.Broadcast(20*time.Millisecond, headingProvider)
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

	stop1 := c.BroadcastPGN(127250, time.Hour, headingProvider)
	b1 := c.broadcasterFor(127250)
	require.NotNil(t, b1)

	stop2 := c.BroadcastPGN(127250, time.Hour, headingProvider)
	b2 := c.broadcasterFor(127250)
	require.NotNil(t, b2)
	assert.NotSame(t, b1, b2, "re-registering a PGN should install a new broadcaster")

	stop1() // stopping the replaced broadcaster must not remove the new one
	assert.NotNil(t, c.broadcasterFor(127250))
	stop2()
	assert.Nil(t, c.broadcasterFor(127250))
}

func TestBroadcast_BlockingProviderDoesNotBlockStopOrClose(t *testing.T) {
	c, _, _ := newCitizenClient(t)
	blocked := make(chan struct{})
	stop := c.BroadcastPGN(127250, time.Hour, func() pgn.Message {
		<-blocked
		return headingProvider()
	})

	require.Eventually(t, func() bool {
		return c.broadcasterFor(127250) != nil
	}, time.Second, time.Millisecond)
	stopDone := make(chan struct{})
	go func() { stop(); close(stopDone) }()
	select {
	case <-stopDone:
	case <-time.After(time.Second):
		t.Fatal("stop blocked on provider")
	}
	closeDone := make(chan error, 1)
	go func() { closeDone <- c.Close() }()
	select {
	case err := <-closeDone:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("Close blocked on provider")
	}
	close(blocked)
}

func TestBroadcast_ProviderPanicIsContained(t *testing.T) {
	c, _, _ := newCitizenClient(t)
	var calls atomic.Int32
	stop := c.BroadcastPGN(127250, 10*time.Millisecond, func() pgn.Message {
		if calls.Add(1) == 1 {
			panic("bad sensor callback")
		}
		return headingProvider()
	})
	defer stop()

	require.Eventually(t, func() bool { return calls.Load() >= 2 }, time.Second, time.Millisecond)
}

func TestBroadcastPGN_SkipsMismatchedProviderMessage(t *testing.T) {
	c, mb, _ := newCitizenClient(t)
	var calls atomic.Int32
	stop := c.BroadcastPGN(127250, 10*time.Millisecond, func() pgn.Message {
		calls.Add(1)
		return &pgn.WaterDepth{}
	})
	defer stop()

	require.Eventually(t, func() bool { return calls.Load() >= 2 }, time.Second, time.Millisecond)
	assert.Empty(t, framesWithPGN(mb.getWritten(), 128267), "a provider must not silently change its declared PGN")
}
