package n2k

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
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

func TestGatewaySessionConstructorBoundsStalledHandshake(t *testing.T) {
	for _, test := range []struct {
		name           string
		commandTimeout time.Duration
		readyTimeout   time.Duration
	}{
		{"command deadline", 20 * time.Millisecond, time.Second},
		{"readiness deadline", time.Second, 20 * time.Millisecond},
	} {
		t.Run(test.name, func(t *testing.T) {
			host, peer := net.Pipe()
			defer func() { _ = peer.Close() }()
			result := make(chan error, 1)
			go func() {
				session, err := NewActisenseGatewaySession(context.Background(), "stalled", func(context.Context) (ActisenseByteStream, error) {
					return host, nil
				}, WithActisenseCommandTimeout(test.commandTimeout), WithActisenseSessionReadyTimeout(test.readyTimeout))
				if session != nil {
					_ = session.Close()
				}
				result <- err
			}()
			select {
			case err := <-result:
				require.Error(t, err)
			case <-time.After(500 * time.Millisecond):
				t.Fatal("constructor ignored its deadline while the peer never read")
			}
		})
	}
}

// stalledGatewaySession acknowledges exactly the setup exchange, then leaves
// the peer open without reading any more bytes. net.Pipe makes the ensuing
// physical write stall deterministic, without relying on a socket buffer size.
func stalledGatewaySession(t *testing.T, changeMode bool) *ActisenseGatewaySession {
	t.Helper()
	host, peer := net.Pipe()
	t.Cleanup(func() { _ = peer.Close() })
	deviceErr := make(chan error, 1)
	go func() {
		buf := make([]byte, 4096)
		for _, mode := range []actisense.OperatingMode{actisense.ModeTransferNormal, actisense.ModeTransferReceiveAll} {
			if _, err := peer.Read(buf); err != nil {
				deviceErr <- err
				return
			}
			payload := make([]byte, 14)
			payload[0] = actisense.BEMOperatingMode
			binary.LittleEndian.PutUint16(payload[12:], uint16(mode))
			wire, err := actisense.EncodeDatagram(actisense.BSTBEMResponse, payload)
			if err == nil {
				_, err = peer.Write(wire)
			}
			if err != nil {
				deviceErr <- err
				return
			}
			if !changeMode {
				break
			}
		}
		deviceErr <- nil
	}()
	mode := ActisenseModeTransferNormal
	if changeMode {
		mode = ActisenseModeTransferReceiveAll
	}
	session, err := NewActisenseGatewaySession(context.Background(), "stall-after-ready", func(context.Context) (ActisenseByteStream, error) {
		return host, nil
	}, WithActisenseSessionMode(mode), WithActisenseCommandTimeout(40*time.Millisecond))
	require.NoError(t, err)
	require.NoError(t, <-deviceErr)
	t.Cleanup(func() { _ = session.Close() })
	return session
}

func TestGatewaySessionEstablishedWritesHonorDeadlines(t *testing.T) {
	for _, operation := range []string{"message", "BEM", "BEM command timeout"} {
		t.Run(operation, func(t *testing.T) {
			session := stalledGatewaySession(t, false)
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
			defer cancel()
			if operation == "BEM command timeout" {
				ctx = context.Background()
			}
			result := make(chan error, 1)
			go func() {
				if operation == "message" {
					result <- session.SendRawPGN(ctx, 127250, 2, 255, []byte{1, 2, 3})
				} else {
					_, err := session.Echo(ctx, []byte{1, 2, 3})
					result <- err
				}
			}()
			select {
			case err := <-result:
				require.ErrorIs(t, err, context.DeadlineExceeded)
			case <-time.After(500 * time.Millisecond):
				t.Fatal("established physical write ignored the caller deadline")
			}
			select {
			case <-session.done:
			case <-time.After(500 * time.Millisecond):
				t.Fatal("canceled physical write left its connection running")
			}
		})
	}
}

func TestGatewaySessionCanceledWriteLeavesReadyConnectionIntact(t *testing.T) {
	session := stalledGatewaySession(t, false)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := session.SendRawPGN(ctx, 127250, 2, 255, []byte{1, 2, 3})
	require.ErrorIs(t, err, context.Canceled)
	status := session.Status()
	require.True(t, status.Connected)
	require.Equal(t, uint64(1), status.Metrics.Protocol.TransportWriteCalls)
}

func TestGatewaySessionCloseBoundsStalledModeRestoration(t *testing.T) {
	session := stalledGatewaySession(t, true)
	result := make(chan error, 1)
	go func() { result <- session.Close() }()
	select {
	case err := <-result:
		// A restoration deadline closes the pipe to interrupt the write.
		if err != nil {
			require.ErrorIs(t, err, io.ErrClosedPipe)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Close waited indefinitely for mode restoration")
	}
	select {
	case <-session.done:
	default:
		t.Fatal("Close returned before the session reader stopped")
	}
}

func TestGatewaySessionCloseReleasesStalledWriterAndRestoration(t *testing.T) {
	session := stalledGatewaySession(t, true)
	writeResult := make(chan error, 1)
	go func() {
		writeResult <- session.SendRawPGN(context.Background(), 127250, 2, 255, []byte{1, 2, 3})
	}()
	require.Eventually(t, func() bool {
		return session.Status().Metrics.Protocol.TransportWriteCalls == 3
	}, time.Second, time.Millisecond)
	closeResult := make(chan error, 1)
	go func() { closeResult <- session.Close() }()
	select {
	case <-closeResult:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Close could not pass a stalled writer to finish bounded restoration")
	}
	select {
	case err := <-writeResult:
		require.Error(t, err)
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Close left the active physical write running")
	}
}

func TestGatewaySessionCloseBoundsStalledTransmitListRestoration(t *testing.T) {
	session := stalledGatewaySession(t, false)
	session.mu.Lock()
	session.txOriginal[127250] = ActisenseTxPGNState{PGN: 127250, Enabled: 1}
	session.txEpoch = session.epoch
	session.mu.Unlock()
	result := make(chan error, 1)
	go func() { result <- session.Close() }()
	select {
	case <-result:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Close waited indefinitely for Tx-list restoration")
	}
}
