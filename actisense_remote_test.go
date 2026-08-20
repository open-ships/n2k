package n2k

import (
	"encoding/binary"
	"testing"

	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/pgn"
	"github.com/open-ships/n2k/raw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func remoteBEMPayload(command byte, sequence uint8, model uint16, serial uint32, errorCode int32, data []byte) []byte {
	response := make([]byte, 14, 14+len(data))
	response[0] = actisense.BSTBEMResponse
	response[1] = byte(12 + len(data))
	response[2] = command
	response[3] = sequence
	binary.LittleEndian.PutUint16(response[4:6], model)
	binary.LittleEndian.PutUint32(response[6:10], serial)
	binary.LittleEndian.PutUint32(response[10:14], uint32(errorCode))
	return append([]byte{0x11, 0x99}, append(response, data...)...)
}

func TestActisenseRemoteIndependentEnvelopeGoldenVector(t *testing.T) {
	inner, err := encodeActisenseRemoteCommand(actisense.BEMOperatingMode, nil)
	require.NoError(t, err)
	assert.Equal(t, []byte{0xA1, 0x01, 0x11}, inner)
	destination := uint8(35)
	priority := uint8(3)
	message := &actisenseRemoteMessage{
		info:    pgn.MessageInfo{PGN: actisenseRemotePGN, Priority: &priority, TargetId: &destination},
		payload: append([]byte{0x11, 0x99}, inner...),
	}
	payload, err := pgn.EncodeMessage(message)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x11, 0x99, 0xA1, 0x01, 0x11}, payload)
}

func TestActisenseRemoteCorrelationRequiresSourceDestinationAndEpoch(t *testing.T) {
	client := &Client{sourceAddr: 42, connectionEpoch: 3, claimEpoch: 4}
	manager := newActisenseRemoteManager(client)
	key := actisenseRemoteKey{
		source: 35, destination: 42, connectionEpoch: 3, claimEpoch: 4,
		bstID: actisense.BSTBEMResponse, bemID: actisense.BEMOperatingMode,
	}
	pending := &actisenseRemotePending{results: make(chan actisenseRemoteResult, 1)}
	manager.pending[key] = pending
	wrongDestination := uint8(41)
	manager.observe(raw.Observation{
		Kind: raw.KindMessage, PGN: actisenseRemotePGN, Source: 35, Destination: &wrongDestination,
		Payload: remoteBEMPayload(actisense.BEMOperatingMode, 1, uint16(actisense.ModelNGX1), 99, 0, []byte{5, 0}),
	})
	select {
	case <-pending.results:
		t.Fatal("promiscuously observed response with wrong local destination correlated")
	default:
	}

	destination := uint8(42)
	manager.observe(raw.Observation{
		Kind: raw.KindMessage, PGN: actisenseRemotePGN, Source: 35, Destination: &destination,
		Payload: remoteBEMPayload(actisense.BEMOperatingMode, 1, uint16(actisense.ModelNGX1), 99, 0, []byte{5, 0}),
	})
	result := <-pending.results
	require.NoError(t, result.err)
	assert.Equal(t, uint8(35), result.response.Origin.Source)
	assert.Equal(t, actisense.BEMPathRemote, result.response.Origin.Path)
	assert.Equal(t, []byte{5, 0}, result.response.Data)
}

func TestActisenseRemoteEpochInvalidationCompletesEveryPendingRequest(t *testing.T) {
	manager := newActisenseRemoteManager(&Client{})
	first := &actisenseRemotePending{results: make(chan actisenseRemoteResult, 1)}
	second := &actisenseRemotePending{results: make(chan actisenseRemoteResult, 1)}
	manager.pending[actisenseRemoteKey{source: 1, bemID: 1}] = first
	manager.pending[actisenseRemoteKey{source: 2, bemID: 2}] = second
	manager.invalidate(ErrActisenseRemoteEpochChanged)
	assert.ErrorIs(t, (<-first.results).err, ErrActisenseRemoteEpochChanged)
	assert.ErrorIs(t, (<-second.results).err, ErrActisenseRemoteEpochChanged)
	assert.Empty(t, manager.pending)
}

func TestActisenseRemoteDeviceErrorRetainsResponseIdentity(t *testing.T) {
	client := &Client{sourceAddr: 42, connectionEpoch: 1, claimEpoch: 1}
	manager := newActisenseRemoteManager(client)
	key := actisenseRemoteKey{
		source: 7, destination: 42, connectionEpoch: 1, claimEpoch: 1,
		bstID: actisense.BSTBEMResponse, bemID: actisense.BEMEcho,
	}
	pending := &actisenseRemotePending{results: make(chan actisenseRemoteResult, 1)}
	manager.pending[key] = pending
	destination := uint8(42)
	manager.observe(raw.Observation{
		Kind: raw.KindMessage, PGN: actisenseRemotePGN, Source: 7, Destination: &destination,
		Payload: remoteBEMPayload(actisense.BEMEcho, 1, uint16(actisense.ModelNGT1), 123, -1159, nil),
	})
	result := <-pending.results
	var deviceErr *actisense.DeviceError
	require.ErrorAs(t, result.err, &deviceErr)
	assert.Equal(t, byte(actisense.BEMEcho), result.response.BEMID)
	assert.Equal(t, uint16(actisense.ModelNGT1), result.response.ModelID)
}
