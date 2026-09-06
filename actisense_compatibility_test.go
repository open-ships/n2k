package n2k

import (
	"context"
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/open-ships/n2k/internal/actisense"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type compatibilityPeer struct {
	connection net.Conn
	writeMu    sync.Mutex
	address    atomic.Uint32
	mode       ActisenseOperatingMode
	commands   chan actisense.Datagram
	wire       chan []byte
	onRemote   func(*compatibilityPeer, actisense.Message, actisense.Datagram)
	done       chan struct{}
}

func newCompatibilityPeer(t *testing.T, mode ActisenseOperatingMode, onRemote func(*compatibilityPeer, actisense.Message, actisense.Datagram)) (ActisenseByteStream, *compatibilityPeer) {
	t.Helper()
	host, device := net.Pipe()
	peer := &compatibilityPeer{connection: device, mode: mode, commands: make(chan actisense.Datagram, 64), wire: make(chan []byte, 64), onRemote: onRemote, done: make(chan struct{})}
	peer.address.Store(42)
	go func() {
		defer close(peer.done)
		parser := actisense.NewParser()
		buffer := make([]byte, 4096)
		for {
			n, err := device.Read(buffer)
			if n > 0 {
				peer.wire <- append([]byte(nil), buffer[:n]...)
				parser.Feed(buffer[:n], peer.handle, nil)
			}
			if err != nil {
				return
			}
		}
	}()
	t.Cleanup(func() {
		_ = host.Close()
		_ = device.Close()
		select {
		case <-peer.done:
		case <-time.After(time.Second):
			t.Error("compatibility peer did not stop")
		}
	})
	return host, peer
}

func (p *compatibilityPeer) send(id byte, payload []byte) {
	wire, err := actisense.EncodeDatagram(id, payload)
	if err != nil {
		panic(err)
	}
	p.writeMu.Lock()
	defer p.writeMu.Unlock()
	_, _ = p.connection.Write(wire)
}

func (p *compatibilityPeer) handle(datagram actisense.Datagram) {
	p.commands <- datagram
	if datagram.ID == actisense.BSTBEMCommand && len(datagram.Payload) > 0 {
		command, data := datagram.Payload[0], datagram.Payload[1:]
		response := make([]byte, 12)
		response[0], response[1], response[2] = command, 1, byte(ActisenseModelNGT1)
		switch command {
		case actisense.BEMOperatingMode:
			if len(data) == 2 {
				p.mode = ActisenseOperatingMode(binary.LittleEndian.Uint16(data))
			}
			response = binary.LittleEndian.AppendUint16(response, uint16(p.mode))
		case actisense.BEMEcho:
			response = append(response, data...)
		}
		p.send(actisense.BSTBEMResponse, response)
		return
	}
	message, ok, err := actisense.DecodeMessage(datagram)
	if !ok || err != nil || message.PGN != actisenseRemotePGN || len(message.Data) < 5 {
		return
	}
	inner, err := actisense.DecodeRaw(message.Data[2:])
	if err != nil || len(inner.Payload) == 0 {
		return
	}
	if inner.Payload[0] == actisense.BEMEcho {
		p.reply(message.Destination, uint8(p.address.Load()), actisense.BEMEcho, inner.Payload[1:], false)
	} else if p.onRemote != nil {
		p.onRemote(p, message, inner)
	}
}

func (p *compatibilityPeer) reply(source, destination, command byte, data []byte, d0 bool) {
	payload := remoteBEMPayload(command, 1, uint16(ActisenseModelNGX1), 1234, 0, data)
	if d0 {
		wire, err := actisense.EncodeMessageD0(actisense.Message{PGN: actisenseRemotePGN, Priority: 3, Source: source, Destination: destination, Data: payload, Direction: actisense.DirectionReceived})
		if err != nil {
			panic(err)
		}
		p.writeMu.Lock()
		defer p.writeMu.Unlock()
		_, _ = p.connection.Write(wire)
		return
	}
	data93 := []byte{3, 0, 0xEF, 1, destination, source, 0, 0, 0, 0, byte(len(payload))}
	p.send(actisense.BSTN2KReceive, append(data93, payload...))
}

func openCompatibilitySession(t *testing.T, host ActisenseByteStream, options ...ActisenseSessionOption) *ActisenseGatewaySession {
	t.Helper()
	options = append([]ActisenseSessionOption{WithActisenseCommandTimeout(time.Second)}, options...)
	session, err := NewActisenseGatewaySession(context.Background(), "compatibility", func(context.Context) (ActisenseByteStream, error) { return host, nil }, options...)
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })
	return session
}

func TestActisensePreservedModeAndGenericSends(t *testing.T) {
	host, peer := newCompatibilityPeer(t, ActisenseModeCANPacketASCII, nil)
	session := openCompatibilitySession(t, host, WithActisensePreserveOperatingMode())
	assert.Equal(t, ActisenseModeCANPacketASCII, session.Status().OperatingMode)
	getMode := <-peer.commands
	assert.Equal(t, []byte{actisense.BEMOperatingMode}, getMode.Payload)
	assert.Empty(t, peer.commands)
	<-peer.wire

	// Independent Type-1 vector: unknown ID AF, payload 10 03, checksum 3C.
	require.NoError(t, session.SendBST(context.Background(), []byte{0xAF, 2, 0x10, 3}))
	assert.Equal(t, []byte{0x10, 2, 0xAF, 2, 0x10, 0x10, 3, 0x3C, 0x10, 3}, <-peer.wire)
	assert.Equal(t, byte(0xAF), (<-peer.commands).ID)
	require.NoError(t, session.SendRaw(context.Background(), []byte{0x21, 0x22, 0x23}))
	assert.Equal(t, []byte{0x21, 0x22, 0x23}, <-peer.wire)
	require.Error(t, session.SendBST(context.Background(), []byte{0xAF, 3, 1}))
	require.Error(t, session.SendRaw(context.Background(), make([]byte, ActisenseMaxRawWrite+1)))
	require.ErrorContains(t, session.SendRawPGN(context.Background(), 127250, 3, 255, []byte{1}), "mode 1 or 2")

	echo, err := session.Echo(context.Background(), []byte{1, 0x10, 0xFF})
	require.NoError(t, err)
	assert.Equal(t, []byte{1, 0x10, 0xFF}, echo)
	<-peer.commands
	require.NoError(t, session.SetOperatingMode(context.Background(), ActisenseModeTransferNormal))
	assert.Equal(t, ActisenseModeTransferNormal, session.Status().OperatingMode)
	<-peer.commands
	require.NoError(t, session.SendRawPGN(context.Background(), 127250, 3, 255, []byte{1}))
	assert.Equal(t, byte(actisense.BSTN2KTransmit), (<-peer.commands).ID)
	require.NoError(t, session.Close())
	assert.Empty(t, peer.commands, "preserving mode must not restore it on Close")
}

func TestActisenseGatewayRemoteUsesVerifiedAddressAndBothBinaryMessages(t *testing.T) {
	for _, d0 := range []bool{false, true} {
		t.Run(map[bool]string{false: "BST93", true: "BSTD0"}[d0], func(t *testing.T) {
			requests := make(chan actisense.Message, 2)
			host, _ := newCompatibilityPeer(t, ActisenseModeTransferNormal, func(peer *compatibilityPeer, message actisense.Message, inner actisense.Datagram) {
				requests <- message
				peer.reply(message.Destination, 41, inner.Payload[0], []byte{6, 0}, d0)
				peer.reply(message.Destination, 42, inner.Payload[0], []byte{5, 0}, d0)
			})
			session := openCompatibilitySession(t, host)
			remote, err := session.ActisenseRemoteDevice(35)
			require.NoError(t, err)
			mode, err := remote.GetOperatingMode(context.Background())
			require.NoError(t, err)
			assert.Equal(t, ActisenseModeCANPacket, mode)
			message := <-requests
			assert.Equal(t, uint8(3), message.Priority)
			assert.Equal(t, uint8(35), message.Destination)
			assert.Equal(t, []byte{0x11, 0x99, 0xA1, 1, 0x11}, message.Data)
			status := session.Status()
			require.NotNil(t, status.GatewaySourceAddress)
			assert.Equal(t, uint8(42), *status.GatewaySourceAddress)
			assert.False(t, status.SourceAuthoritative)
			assert.Equal(t, uint64(1), status.RemoteMetrics.BEMRequests)
		})
	}
}

func TestActisenseGatewayRemoteAddressChangeCancelsPendingTrain(t *testing.T) {
	waiting := make(chan struct{}, 1)
	host, peer := newCompatibilityPeer(t, ActisenseModeTransferNormal, func(peer *compatibilityPeer, message actisense.Message, inner actisense.Datagram) {
		if inner.Payload[0] == actisense.BEMProductInfo {
			waiting <- struct{}{}
			return
		}
		peer.reply(message.Destination, uint8(peer.address.Load()), inner.Payload[0], []byte{2, 0}, false)
	})
	session := openCompatibilitySession(t, host)
	remote, err := session.ActisenseRemoteDevice(35, WithActisenseRemoteMultiReplyInactivity(time.Second))
	require.NoError(t, err)
	result := make(chan error, 1)
	go func() { _, err := remote.GetProductInfo(context.Background()); result <- err }()
	select {
	case <-waiting:
	case <-time.After(time.Second):
		t.Fatal("product request did not reach gateway")
	}
	epoch := session.Status().IdentityEpoch
	peer.address.Store(43)
	_, err = remote.GetOperatingMode(context.Background())
	require.NoError(t, err)
	require.ErrorIs(t, <-result, ErrActisenseRemoteEpochChanged)
	assert.Greater(t, session.Status().IdentityEpoch, epoch)
	assert.Equal(t, uint8(43), *session.Status().GatewaySourceAddress)
}

func TestActisenseGatewayRemoteDisconnectCancelsWithoutRetry(t *testing.T) {
	seen := make(chan struct{}, 1)
	host, peer := newCompatibilityPeer(t, ActisenseModeTransferNormal, func(_ *compatibilityPeer, _ actisense.Message, _ actisense.Datagram) { seen <- struct{}{} })
	session := openCompatibilitySession(t, host)
	remote, err := session.ActisenseRemoteDevice(35)
	require.NoError(t, err)
	result := make(chan error, 1)
	go func() { _, err := remote.GetOperatingMode(context.Background()); result <- err }()
	select {
	case <-seen:
	case <-time.After(time.Second):
		t.Fatal("request did not reach gateway")
	}
	require.NoError(t, peer.connection.Close())
	select {
	case err := <-result:
		require.True(t, errors.Is(err, ErrActisenseRemoteEpochChanged) || errors.Is(err, actisense.ErrSessionClosed), "%v", err)
	case <-time.After(time.Second):
		t.Fatal("disconnect did not cancel request")
	}
	assert.Nil(t, session.Status().GatewaySourceAddress)
	require.ErrorIs(t, session.SendRaw(context.Background(), []byte{1}), ErrActisenseNotReady)
}

func TestActisenseGatewayRemoteReconnectPreservesEachMode(t *testing.T) {
	seen := make(chan struct{}, 1)
	firstHost, first := newCompatibilityPeer(t, ActisenseModeTransferNormal, func(_ *compatibilityPeer, _ actisense.Message, _ actisense.Datagram) { seen <- struct{}{} })
	secondHost, second := newCompatibilityPeer(t, ActisenseModeTransferReceiveAll, func(peer *compatibilityPeer, message actisense.Message, inner actisense.Datagram) {
		peer.reply(message.Destination, uint8(peer.address.Load()), inner.Payload[0], []byte{5, 0}, false)
	})
	second.address.Store(43)
	var opens atomic.Uint32
	session, err := NewActisenseGatewaySession(context.Background(), "reconnect", func(context.Context) (ActisenseByteStream, error) {
		switch opens.Add(1) {
		case 1:
			return firstHost, nil
		case 2:
			return secondHost, nil
		default:
			return nil, errors.New("unexpected reconnect")
		}
	}, WithActisensePreserveOperatingMode(), WithActisenseSessionReconnect(ReconnectPolicy{InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond}))
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })
	remote, err := session.ActisenseRemoteDevice(35)
	require.NoError(t, err)
	result := make(chan error, 1)
	go func() { _, err := remote.GetOperatingMode(context.Background()); result <- err }()
	select {
	case <-seen:
	case <-time.After(time.Second):
		t.Fatal("first request did not reach gateway")
	}
	require.NoError(t, first.connection.Close())
	require.ErrorIs(t, <-result, ErrActisenseRemoteEpochChanged)
	require.Eventually(t, func() bool {
		_, err := session.transport.EpochRequester(2)
		return err == nil
	}, time.Second, time.Millisecond)
	assert.Equal(t, ActisenseModeTransferReceiveAll, session.Status().OperatingMode)
	assert.Nil(t, session.Status().GatewaySourceAddress)
	mode, err := remote.GetOperatingMode(context.Background())
	require.NoError(t, err)
	assert.Equal(t, ActisenseModeCANPacket, mode)
	assert.Equal(t, uint8(43), *session.Status().GatewaySourceAddress)
	assert.Equal(t, uint64(2), session.Status().ConnectionEpoch)
	require.NoError(t, session.Close())
	for _, peer := range []*compatibilityPeer{first, second} {
		for len(peer.commands) > 0 {
			command := <-peer.commands
			if command.ID == actisense.BSTBEMCommand {
				assert.Equal(t, []byte{actisense.BEMOperatingMode}, command.Payload, "no implicit mode setters")
			}
		}
	}
}

func TestActisenseGatewayRemoteProbeRejectsUnrelatedEcho(t *testing.T) {
	probe := &actisenseGatewayProbe{source: 35, epoch: 2, result: make(chan actisenseGatewayProbeResult, 1)}
	copy(probe.challenge[:], "unique challenge")
	session := &ActisenseGatewaySession{epoch: 2, remoteProbe: probe}
	echo, err := actisense.EncodeEcho(probe.challenge[:])
	require.NoError(t, err)
	message := actisense.Message{PGN: actisenseRemotePGN, Source: 35, HasSource: true, Destination: 42,
		Data: remoteBEMPayload(actisense.BEMEcho, 1, uint16(ActisenseModelNGT1), 1234, 0, echo)}
	for _, change := range []func(*actisense.Message){
		func(m *actisense.Message) { m.Source = 34 },
		func(m *actisense.Message) { m.Destination = 255 },
		func(m *actisense.Message) { m.Data[2] = 0xA3 },
		func(m *actisense.Message) { m.Data[len(m.Data)-1] ^= 1 },
	} {
		wrong := message
		wrong.Data = append([]byte(nil), message.Data...)
		change(&wrong)
		assert.False(t, session.handleRemoteProbe(wrong))
		assert.Empty(t, probe.result)
	}
	assert.True(t, session.handleRemoteProbe(message))
	assert.Equal(t, uint8(42), (<-probe.result).address)
}
