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
	require.Same(t, first, <-fast.ch)

	hub.publish(second)
	require.Same(t, second, <-fast.ch)
	require.Same(t, first, <-slow.ch)
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
	require.Same(t, last, <-sub.ch)
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
