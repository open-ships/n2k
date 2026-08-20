package actisense

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type commandCall struct {
	command byte
	data    []byte
	multi   bool
}

type scriptedRequester struct {
	calls          []commandCall
	response       BEMResponse
	responses      []BEMResponse
	err            error
	multiErr       error
	stopAfterReply int
}

func (r *scriptedRequester) Request(_ context.Context, command byte, data []byte) (BEMResponse, error) {
	r.calls = append(r.calls, commandCall{command: command, data: append([]byte(nil), data...)})
	response := r.response
	if response.BEMID == 0 {
		response = responseFor(command, 1, ModelNGX1, nil)
	}
	return response, r.err
}

func (r *scriptedRequester) RequestMulti(_ context.Context, command byte, data []byte, _ time.Duration, complete func([]BEMResponse) (bool, error)) ([]BEMResponse, error) {
	r.calls = append(r.calls, commandCall{command: command, data: append([]byte(nil), data...), multi: true})
	delivered := make([]BEMResponse, 0, len(r.responses))
	for index, response := range r.responses {
		if r.stopAfterReply > 0 && index >= r.stopAfterReply {
			break
		}
		delivered = append(delivered, response)
		done, err := complete(delivered)
		if err != nil || done {
			return delivered, err
		}
	}
	return delivered, r.multiErr
}

func TestCommandSetPreservesTypedProductInfoPartialOnTimeout(t *testing.T) {
	requester := &scriptedRequester{
		responses: []BEMResponse{
			responseFor(BEMProductInfo, 1, ModelNGT1, []byte{0x34, 0x08, 0x65, 0, 1, 2}),
			responseFor(BEMProductInfo, 2, ModelNGT1, padded32("NGT-1")),
		},
		multiErr: context.DeadlineExceeded,
	}
	commands := NewCommandSet(requester, CommandSetConfig{})
	info, err := commands.GetProductInfo(context.Background())
	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, "NGT-1", info.Model)
	assert.Equal(t, uint16(0x0834), info.NMEA2000Version)
	assert.Equal(t, ModelNGT1, commands.DeviceCapabilities().ModelID)
}

func TestCommandSetExposesExplicitPersistenceAndLifecycleVerbs(t *testing.T) {
	requester := &scriptedRequester{}
	commands := NewCommandSet(requester, CommandSetConfig{})
	require.NoError(t, commands.CommitEEPROM(context.Background()))
	require.NoError(t, commands.CommitFlash(context.Background()))
	require.NoError(t, commands.Reinitialize(context.Background()))
	require.NoError(t, commands.ActivatePGNLists(context.Background()))
	require.Len(t, requester.calls, 4)
	assert.Equal(t, []byte{BEMCommitEEPROM, BEMCommitFlash, BEMReinitialize, BEMActivatePGNLists}, []byte{
		requester.calls[0].command, requester.calls[1].command, requester.calls[2].command, requester.calls[3].command,
	})
	for _, call := range requester.calls {
		assert.Empty(t, call.data)
	}
}

func TestCommandSetPortAndCANSetPayloads(t *testing.T) {
	requester := &scriptedRequester{}
	commands := NewCommandSet(requester, CommandSetConfig{})
	requester.response = responseFor(BEMPortBaudrate, 1, ModelNGX1, []byte{
		2, 1, byte(HardwareSerialBST), 0x00, 0xC2, 0x01, 0x00, 0x00, 0xC2, 0x01, 0x00,
	})
	_, err := commands.SetPortBaudrate(context.Background(), 1, 115200, 115200)
	require.NoError(t, err)
	assert.Equal(t, []byte{1, 0x00, 0xC2, 0x01, 0x00, 0x00, 0xC2, 0x01, 0x00}, requester.calls[0].data)

	requester.response = responseFor(BEMCANConfig, 1, ModelNGX1, []byte{8, 7, 6, 5, 4, 3, 2, 1, 42})
	accepted, err := commands.SetCANConfig(context.Background(), CANConfig{NAME: 0x0102030405060708, SourceAddress: 42})
	require.NoError(t, err)
	assert.Equal(t, CANConfig{NAME: 0x0102030405060708, SourceAddress: 42}, accepted)
	assert.Equal(t, []byte{8, 7, 6, 5, 4, 3, 2, 1, 42}, requester.calls[1].data)
}

func TestCommandSetRemoteEchoCeilingAccountsForEnvelope(t *testing.T) {
	commands := NewCommandSet(&scriptedRequester{}, CommandSetConfig{Remote: true})
	_, err := commands.Echo(context.Background(), make([]byte, 207))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at most 206")
}

func TestCommandSetRejectsNilRequester(t *testing.T) {
	commands := NewCommandSet(nil, CommandSetConfig{})
	_, err := commands.RawRequest(context.Background(), BEMEcho, nil)
	assert.EqualError(t, err, "actisense: command handle is closed")
}
