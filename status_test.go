package n2k

import (
	"context"
	"testing"

	"github.com/open-ships/n2k/internal/gateway"
	"github.com/open-ships/n2k/raw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientStatusTracksLifecycle(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil))
	require.NoError(t, err)
	status := c.Status()
	assert.True(t, status.AddressClaimed)
	assert.False(t, status.Closed)
	assert.NoError(t, status.TerminalError)
	assert.Equal(t, defaultWriteQueue, status.WriteQueueCapacity)

	require.NoError(t, c.Close())
	assert.True(t, c.Status().Closed)
	assert.NoError(t, c.Err())
}

func TestNilClientStatus(t *testing.T) {
	var c *Client
	assert.True(t, c.Status().Closed)
	assert.ErrorIs(t, c.Err(), ErrClientClosed)
}

func TestClientStatusTracksGatewayAndTransportObservations(t *testing.T) {
	c := &Client{observationHub: newObservationHub(2)}

	c.publishObservation(raw.Observation{Kind: raw.KindGateway})
	c.publishObservation(raw.Observation{Kind: raw.KindTransportError})
	status := c.Status()
	assert.Equal(t, uint64(1), status.GatewayEventsObserved)
	assert.Equal(t, uint64(1), status.TransportErrorsObserved)
}

func TestClientStatusExposesActisenseBusMetrics(t *testing.T) {
	bus := gateway.NewActisenseRawTCPBus(nil, "127.0.0.1:1", nil)
	c := &Client{bus: bus}

	status := c.Status()
	require.NotNil(t, status.Actisense)
	assert.Zero(t, status.Actisense.ConnectionEpochs)
	assert.NotNil(t, status.Actisense.Protocol.BSTFrames)
}
