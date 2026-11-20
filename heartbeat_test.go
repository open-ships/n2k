package n2k

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// framesWithPGN filters frames whose CAN ID carries the given PGN.
func framesWithPGN(frames []can.Frame, pgnNum uint32) []can.Frame {
	var out []can.Frame
	for _, f := range frames {
		if framer.ParseCANID(f.ID).PGN == pgnNum {
			out = append(out, f)
		}
	}
	return out
}

func TestHeartbeat_PeriodicTransmission(t *testing.T) {
	mb := newMockBus()
	c, err := NewClient(context.Background(),
		WithBus(mb),
		WithClaimTimeout(50*time.Millisecond),
		WithHeartbeatInterval(30*time.Millisecond),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	ok := waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 126993)) >= 3
	})
	require.True(t, ok, "expected at least 3 heartbeats")

	beats := framesWithPGN(mb.getWritten(), 126993)
	// Priority 7 per the PGN's metadata default.
	assert.Equal(t, uint8(7), framer.ParseCANID(beats[0].ID).Priority)

	// Sequence counter advances monotonically across the first beats.
	var seqs []uint64
	for _, f := range beats[:3] {
		msg, err := pgn.DecodeMessage(pgn.MessageInfo{PGN: 126993}, f.Data[:f.Length])
		require.NoError(t, err)
		hb, ok := msg.(*pgn.Heartbeat)
		require.True(t, ok, "expected *pgn.Heartbeat, got %T", msg)
		require.NotNil(t, hb.SequenceCounter)
		seqs = append(seqs, *hb.SequenceCounter)
		require.NotNil(t, hb.DataTransmitOffset)
		assert.Equal(t, uint64(30), *hb.DataTransmitOffset, "offset should hold the interval in ms ticks")
	}
	assert.Equal(t, []uint64{0, 1, 2}, seqs)
}

func TestHeartbeat_DisabledWithZeroInterval(t *testing.T) {
	mb := newMockBus()
	c, err := NewClient(context.Background(),
		WithBus(mb),
		WithClaimTimeout(50*time.Millisecond),
		WithHeartbeatInterval(0),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	time.Sleep(150 * time.Millisecond)
	assert.Empty(t, framesWithPGN(mb.getWritten(), 126993))
}

func TestHeartbeat_SetIntervalReenables(t *testing.T) {
	mb := newMockBus()
	c, err := NewClient(context.Background(),
		WithBus(mb),
		WithClaimTimeout(50*time.Millisecond),
		WithHeartbeatInterval(0),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	c.heartbeat.setInterval(20 * time.Millisecond)
	ok := waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 126993)) >= 2
	})
	assert.True(t, ok, "setInterval should re-enable a disabled heartbeat")

	c.heartbeat.setInterval(0)
	time.Sleep(50 * time.Millisecond)
	before := len(framesWithPGN(mb.getWritten(), 126993))
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, before, len(framesWithPGN(mb.getWritten(), 126993)), "setInterval(0) should stop heartbeats")
}

func TestHeartbeat_SequenceCounterWraps(t *testing.T) {
	h := &heartbeater{interval: time.Second}
	h.seq = 252
	m1 := h.message()
	m2 := h.message()
	require.NotNil(t, m1.SequenceCounter)
	require.NotNil(t, m2.SequenceCounter)
	assert.Equal(t, uint64(252), *m1.SequenceCounter)
	assert.Equal(t, uint64(0), *m2.SequenceCounter, "sequence counter must wrap after 252")
}

func TestHeartbeat_ReplayClientHasNone(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	assert.Nil(t, c.heartbeat)
}

func TestHeartbeatFollowsConnectionAndClaimReadiness(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		bus := &startupEpochBus{mockBus: newMockBus()}
		const interval = 40 * time.Millisecond
		client, err := NewClient(t.Context(), WithBus(bus), WithName(DeviceName{IdentityNumber: 10}),
			WithClaimTimeout(20*time.Millisecond), WithHeartbeatInterval(interval))
		require.NoError(t, err)
		defer func() { _ = client.Close() }()
		synctest.Wait()
		beats := func() []can.Frame { return framesWithPGN(bus.getWritten(), 126993) }
		require.Len(t, beats(), 1)

		bus.observer(false, 1)
		synctest.Wait()
		admitted := client.Status().ProtocolWritesAccepted
		time.Sleep(3 * interval)
		synctest.Wait()
		require.Equal(t, admitted, client.Status().ProtocolWritesAccepted)
		require.Zero(t, client.Status().ProtocolRequiredQueueDepth)

		bus.observer(true, 2)
		synctest.Wait()
		time.Sleep(addressContentionWindow - time.Millisecond)
		require.Len(t, beats(), 1)
		time.Sleep(time.Millisecond)
		synctest.Wait()
		require.Len(t, beats(), 2)

		bus.inbound <- claimFrame(DeviceName{IdentityNumber: 1}, 251)
		synctest.Wait()
		time.Sleep(addressContentionWindow - time.Millisecond)
		require.Len(t, beats(), 2)
		require.Zero(t, client.Status().ProtocolRequiredQueueDepth)
		time.Sleep(time.Millisecond)
		synctest.Wait()
		require.Len(t, beats(), 3)
		require.Equal(t, uint8(250), framer.ParseCANID(beats()[2].ID).Source)
	})
}
