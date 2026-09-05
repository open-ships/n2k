package n2k

import (
	"errors"
	"testing"

	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageHubSubscriptionsAreIndependent(t *testing.T) {
	hub := newMessageHub(1)
	fast := hub.subscribe()
	slow := hub.subscribe()
	t.Cleanup(fast.unsubscribe)
	t.Cleanup(slow.unsubscribe)

	first := &pgn.VesselHeading{}
	second := &pgn.RateOfTurn{}
	hub.publish(first)
	require.Equal(t, first, <-fast.ch)

	hub.publish(second)
	require.Equal(t, second, <-fast.ch)
	require.Equal(t, first, <-slow.ch)
	_, open := <-slow.ch
	assert.False(t, open)
	assert.ErrorIs(t, slow.terminalError(), ErrReceiveOverflow)
	assert.NoError(t, fast.terminalError())
}

func TestMessageHubReportsBacklogLoss(t *testing.T) {
	hub := newMessageHub(1)
	hub.publish(&pgn.VesselHeading{})
	last := &pgn.RateOfTurn{}
	hub.publish(last)

	sub := hub.subscribe()
	t.Cleanup(sub.unsubscribe)
	require.Equal(t, last, <-sub.ch)
	_, open := <-sub.ch
	assert.False(t, open)
	assert.ErrorIs(t, sub.terminalError(), ErrReceiveOverflow)
}

func TestMessageHubPropagatesTerminalError(t *testing.T) {
	hub := newMessageHub(1)
	sub := hub.subscribe()
	want := errors.New("bus failed")
	hub.close(want)

	_, open := <-sub.ch
	assert.False(t, open)
	assert.ErrorIs(t, sub.terminalError(), want)

	late := hub.subscribe()
	_, open = <-late.ch
	assert.False(t, open)
	assert.ErrorIs(t, late.terminalError(), want)
}

func TestMessageHubSubscribersOwnMeasurementsAndMetadata(t *testing.T) {
	hub := newMessageHub(2)
	first := hub.subscribe()
	second := hub.subscribe()
	t.Cleanup(first.unsubscribe)
	t.Cleanup(second.unsubscribe)
	heading := uint64(15708)
	message := &pgn.VesselHeading{
		Info:    pgn.MessageInfo{Priority: pgn.Priority(2), DecodeIssues: []string{"diagnostic"}},
		Heading: &heading,
	}
	hub.publish(message)
	// Ownership starts when publish returns, before either subscriber reads.
	*message.Heading = 123
	message.Info.DecodeIssues[0] = "publisher changed"
	ownedFirst := (<-first.ch).(*pgn.VesselHeading)
	ownedSecond := (<-second.ch).(*pgn.VesselHeading)
	require.NotSame(t, ownedFirst, ownedSecond)
	require.Equal(t, uint64(15708), *ownedFirst.Heading)
	require.Equal(t, uint64(15708), *ownedSecond.Heading)
	*ownedFirst.Heading = 456
	*ownedFirst.Info.Priority = 7
	ownedFirst.Info.DecodeIssues[0] = "subscriber changed"
	assert.Equal(t, uint64(15708), *ownedSecond.Heading)
	assert.Equal(t, uint8(2), *ownedSecond.Info.Priority)
	assert.Equal(t, []string{"diagnostic"}, ownedSecond.Info.DecodeIssues)
}

func TestMessageHubBacklogOwnsPublisherData(t *testing.T) {
	hub := newMessageHub(1)
	message := &pgn.UnknownPGN{Data: []byte{1, 2, 3}}
	hub.publish(message)
	message.Data[0] = 9
	sub := hub.subscribe()
	t.Cleanup(sub.unsubscribe)
	owned := (<-sub.ch).(*pgn.UnknownPGN)
	assert.Equal(t, []byte{1, 2, 3}, owned.Data)
}
