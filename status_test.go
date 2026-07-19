package n2k

import (
	"context"
	"testing"

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
