package n2k

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fastPacketFrames encodes msg and frames it as fast-packet CAN frames from
// the given source address.
func fastPacketFrames(t *testing.T, msg pgn.Message, prio, source, dest uint8) []can.Frame {
	t.Helper()
	payload, err := pgn.EncodeMessage(msg)
	require.NoError(t, err)
	canID := framer.BuildCANID(msg.PGNNumber(), prio, source, dest)
	return framer.FrameFastPacket(canID, payload, 0)
}

// singleFrame encodes msg into one CAN frame from the given source address.
func singleFrame(t *testing.T, msg pgn.Message, prio, source, dest uint8) can.Frame {
	t.Helper()
	payload, err := pgn.EncodeMessage(msg)
	require.NoError(t, err)
	canID := framer.BuildCANID(msg.PGNNumber(), prio, source, dest)
	return framer.FrameSingle(canID, payload)
}

// collectSystem registers a capture handler on the client's system router.
func collectSystem(c *Client) (func() []pgn.Message, func(pgn.Message) bool) {
	var mu sync.Mutex
	var got []pgn.Message
	c.system.addHandler(func(m pgn.Message) {
		mu.Lock()
		got = append(got, m)
		mu.Unlock()
	})
	snapshot := func() []pgn.Message {
		mu.Lock()
		defer mu.Unlock()
		return append([]pgn.Message(nil), got...)
	}
	return snapshot, nil
}

// waitFor polls until cond returns true or the timeout elapses.
func waitFor(t *testing.T, timeout time.Duration, cond func() bool) bool {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return true
		}
		time.Sleep(5 * time.Millisecond)
	}
	return cond()
}

func TestSystemRouter_DecodesProtocolPGNsDespiteUserFilter(t *testing.T) {
	mb := newMockBus()
	c, err := NewClient(context.Background(),
		WithBus(mb),
		WithClaimTimeout(50*time.Millisecond),
		Filter("pgn == 999999"), // user filter matches nothing
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	snapshot, _ := collectSystem(c)

	product := &pgn.ProductInformation{ModelId: "radar", ProductCode: u64(77)}
	for _, f := range fastPacketFrames(t, product, 6, 0x42, 255) {
		mb.inbound <- f
	}

	ok := waitFor(t, 2*time.Second, func() bool {
		for _, m := range snapshot() {
			if pi, ok := m.(*pgn.ProductInformation); ok && pi.ModelId == "radar" {
				return true
			}
		}
		return false
	})
	assert.True(t, ok, "system router should decode ProductInformation even with a non-matching user filter")
}

func TestSystemRouter_IgnoresUnlistedPGNs(t *testing.T) {
	mb := newMockBus()
	c, err := NewClient(context.Background(),
		WithBus(mb),
		WithClaimTimeout(50*time.Millisecond),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	snapshot, _ := collectSystem(c)

	h := uint64(15708)
	heading := &pgn.VesselHeading{Heading: &h}
	mb.inbound <- singleFrame(t, heading, 2, 0x42, 255)

	// Give the router a moment; VesselHeading is not a system PGN.
	time.Sleep(100 * time.Millisecond)
	for _, m := range snapshot() {
		_, isHeading := m.(*pgn.VesselHeading)
		assert.False(t, isHeading, "VesselHeading must not pass the system gate")
	}
}

func TestSystemRouter_DynamicPGNRefcount(t *testing.T) {
	mb := newMockBus()
	c, err := NewClient(context.Background(),
		WithBus(mb),
		WithClaimTimeout(50*time.Millisecond),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	snapshot, _ := collectSystem(c)

	h := uint64(15708)
	heading := &pgn.VesselHeading{Heading: &h}

	c.system.addPGN(127250)
	c.system.addPGN(127250)
	c.system.removePGN(127250) // still one reference held
	mb.inbound <- singleFrame(t, heading, 2, 0x42, 255)

	ok := waitFor(t, 2*time.Second, func() bool {
		for _, m := range snapshot() {
			if _, ok := m.(*pgn.VesselHeading); ok {
				return true
			}
		}
		return false
	})
	assert.True(t, ok, "PGN with one remaining reference should pass the gate")

	c.system.removePGN(127250) // last reference gone
	before := len(snapshot())
	mb.inbound <- singleFrame(t, heading, 2, 0x42, 255)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, before, len(snapshot()), "PGN with no references should be gated out")
}
