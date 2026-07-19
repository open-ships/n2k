package n2k

import (
	"context"
	"testing"

	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func protocolQueueTestClient() (*Client, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		ctx:        ctx,
		cancel:     cancel,
		log:        testLogger(),
		writeCh:    make(chan writeJob, 1),
		protocolTx: newProtocolTransmitter(testLogger()),
	}
	client.protocolTx.required = make(chan writeJob, 1)
	client.protocolTx.advisory = make(chan writeJob, 1)
	return client, cancel
}

func TestProtocolTransmissionHasDedicatedAdmission(t *testing.T) {
	client, cancel := protocolQueueTestClient()
	defer cancel()

	client.writeCh <- writeJob{msg: &pgn.VesselHeading{}, result: newWriteResult()}
	result := client.writeProtocol("heartbeat", protocolRequired, &pgn.Heartbeat{})

	assert.Len(t, client.writeCh, 1)
	assert.Len(t, client.protocolTx.required, 1)
	select {
	case <-result.Done():
		t.Fatal("admitted protocol write completed before a writer handled it")
	default:
	}
}

func TestRequiredProtocolQueueOverflowIsTerminal(t *testing.T) {
	client, cancel := protocolQueueTestClient()
	defer cancel()

	client.protocolTx.required <- writeJob{result: newWriteResult(), protocol: true}
	err := client.writeProtocol("ISO response", protocolRequired, &pgn.IsoAcknowledgement{}).Wait()

	require.ErrorIs(t, err, ErrProtocolQueueFull)
	require.Error(t, client.Err())
	assert.Contains(t, client.Err().Error(), "required protocol transmission")
}

func TestAdvisoryProtocolQueueOverflowIsObservableButNotTerminal(t *testing.T) {
	client, cancel := protocolQueueTestClient()
	defer cancel()

	client.protocolTx.advisory <- writeJob{result: newWriteResult(), protocol: true}
	err := client.writeProtocol("device probe", protocolAdvisory, &pgn.IsoRequest{}).Wait()

	require.ErrorIs(t, err, ErrProtocolQueueFull)
	assert.NoError(t, client.Err())
	assert.Equal(t, uint64(1), client.protocolTx.rejected.Load())
}
