package n2k

import (
	"context"
	"encoding/binary"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/open-ships/n2k/internal/actisense"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sessionWireEvent struct {
	id      byte
	command byte
	data    []byte
}

func serveActisenseGateway(connection net.Conn, events chan<- sessionWireEvent) {
	defer func() { _ = connection.Close() }()
	parser := actisense.NewParser()
	mode := actisense.ModeTransferNormal
	buf := make([]byte, 4096)
	for {
		n, err := connection.Read(buf)
		if n > 0 {
			parser.Feed(buf[:n], func(datagram actisense.Datagram) {
				event := sessionWireEvent{id: datagram.ID}
				if datagram.ID == actisense.BSTBEMCommand && len(datagram.Payload) != 0 {
					event.command = datagram.Payload[0]
					event.data = append([]byte(nil), datagram.Payload[1:]...)
				}
				events <- event
				if event.id != actisense.BSTBEMCommand {
					return
				}
				responseData := []byte(nil)
				if event.command == actisense.BEMOperatingMode {
					if len(event.data) == 2 {
						mode = actisense.OperatingMode(binary.LittleEndian.Uint16(event.data))
					}
					responseData = make([]byte, 2)
					binary.LittleEndian.PutUint16(responseData, uint16(mode))
				}
				payload := make([]byte, 12, 12+len(responseData))
				payload[0], payload[1] = event.command, 1
				binary.LittleEndian.PutUint16(payload[2:4], uint16(actisense.ModelNGX1))
				binary.LittleEndian.PutUint32(payload[4:8], 1234)
				payload = append(payload, responseData...)
				wire, encodeErr := actisense.EncodeDatagram(actisense.BSTBEMResponse, payload)
				if encodeErr == nil {
					_, _ = connection.Write(wire)
				}
			}, nil)
		}
		if err != nil {
			return
		}
	}
}

func TestGatewaySessionDoesNotMutateTransmitListImplicitly(t *testing.T) {
	clientConnection, deviceConnection := net.Pipe()
	events := make(chan sessionWireEvent, 16)
	go serveActisenseGateway(deviceConnection, events)
	opened := false
	session, err := NewActisenseGatewaySession(context.Background(), "pipe", func(context.Context) (ActisenseByteStream, error) {
		if opened {
			return nil, errors.New("test stream reopened")
		}
		opened = true
		return clientConnection, nil
	})
	require.NoError(t, err)

	getMode := <-events
	setMode := <-events
	assert.Equal(t, byte(actisense.BEMOperatingMode), getMode.command)
	assert.Empty(t, getMode.data)
	assert.Equal(t, byte(actisense.BEMOperatingMode), setMode.command)
	assert.Equal(t, []byte{2, 0}, setMode.data)

	require.NoError(t, session.SendRawPGN(context.Background(), 127250, 2, 255, []byte{1, 2, 3}))
	send := <-events
	assert.Equal(t, byte(actisense.BSTN2KTransmit), send.id)
	assert.Zero(t, send.command)
	assert.Empty(t, events, "SendRawPGN must not issue hidden 0x47/0x4B list commands")

	status := session.Status()
	assert.True(t, status.Connected)
	assert.False(t, status.SourceAuthoritative)
	assert.True(t, status.ReceiveAll)
	assert.Equal(t, uint64(2), status.Metrics.Protocol.BEMRequests)
	require.NoError(t, session.Close())
	restore := <-events
	assert.Equal(t, byte(actisense.BEMOperatingMode), restore.command)
	assert.Equal(t, []byte{1, 0}, restore.data)
}

func TestGatewaySessionRejectsSourceAuthoritativeMode(t *testing.T) {
	_, err := NewActisenseGatewaySession(context.Background(), "unused", func(context.Context) (ActisenseByteStream, error) {
		return nil, nil
	}, WithActisenseSessionMode(ActisenseModeCANPacket), WithActisenseSessionReadyTimeout(10*time.Millisecond))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be 1 or 2")
}
