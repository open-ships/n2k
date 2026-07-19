package n2k

import (
	"context"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/raw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplayObservationPreservesMessageInfo(t *testing.T) {
	frame := vesselHeadingFrame(t)
	timestamp := time.Unix(1_720_000_000, 123_000_000)
	observation := Observation{
		Kind:                  ObservationFrame,
		Timestamp:             timestamp,
		ReceivedAt:            timestamp.Add(time.Second),
		TransportTimestamp:    17*time.Hour + 33*time.Minute,
		HasTransportTimestamp: true,
		AdapterID:             "capture:test",
		NetworkID:             "can7",
		Direction:             DirectionReceived,
		Frame:                 &frame,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for message, err := range Receive(ctx, ReplayObservations([]Observation{observation})) {
		require.NoError(t, err)
		carrier, ok := message.(infoCarrier)
		require.True(t, ok)
		info := carrier.MessageInfo()
		assert.Equal(t, timestamp, info.Timestamp)
		assert.Equal(t, observation.ReceivedAt, info.ReceivedAt)
		assert.Equal(t, observation.TransportTimestamp, info.TransportTimestamp)
		assert.True(t, info.HasTransportTimestamp)
		assert.Equal(t, "capture:test", info.AdapterID)
		assert.Equal(t, "can7", info.NetworkID)
		assert.Equal(t, raw.DirectionReceived, info.Direction)
		return
	}
	t.Fatal("expected a decoded message")
}

func TestObservationHubClonesAndFailsOnlySlowSubscriber(t *testing.T) {
	hub := newObservationHub(1)
	fast := hub.subscribe()
	slow := hub.subscribe()
	frame := can.Frame{ID: 1, Length: 1, Data: [8]byte{7}}
	first := frameObservation(frame, "test", "can0", raw.DirectionReceived)
	hub.publish(first)
	got := <-fast.ch
	got.Frame.Data[0] = 99
	assert.Equal(t, byte(7), first.Frame.Data[0])

	hub.publish(first)
	hub.publish(first)
	assert.ErrorIs(t, slow.terminalError(), ErrObservationOverflow)
	assert.Equal(t, 0, hub.subscriberCount())
}
