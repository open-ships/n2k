package gateway

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeBEMResponse(t *testing.T, connection net.Conn, command byte, data []byte) {
	t.Helper()
	payload := make([]byte, 12, 12+len(data))
	payload[0] = command
	payload = append(payload, data...)
	wire, err := actisense.EncodeDatagram(actisense.BSTBEMResponse, payload)
	require.NoError(t, err)
	_, err = connection.Write(wire)
	require.NoError(t, err)
}

func TestActisenseRawSerialBusReadinessFramesAndRestore(t *testing.T) {
	clientConnection, gatewayConnection := net.Pipe()
	defer func() { _ = gatewayConnection.Close() }()
	bus := NewActisenseRawSerialBus(testLogger(), "fake")
	bus.stream.open = func(context.Context) (actisenseConnection, error) { return clientConnection, nil }
	bus.stream.commandTimeout = time.Second

	datagrams := make(chan actisense.Datagram, 16)
	go func() {
		parser := actisense.NewParser()
		buf := make([]byte, 4096)
		for {
			n, err := gatewayConnection.Read(buf)
			if n > 0 {
				parser.Feed(buf[:n], func(datagram actisense.Datagram) { datagrams <- datagram }, nil)
			}
			if err != nil {
				return
			}
		}
	}()

	frames := make(chan can.Frame, 2)
	runErr := make(chan error, 1)
	go func() { runErr <- bus.Run(context.Background(), func(frame can.Frame) { frames <- frame }) }()

	getMode := <-datagrams
	require.Equal(t, byte(actisense.BSTBEMCommand), getMode.ID)
	require.Equal(t, []byte{actisense.BEMOperatingMode}, getMode.Payload)
	select {
	case <-bus.Ready():
		t.Fatal("raw Bus became ready before mode acknowledgement")
	case <-time.After(25 * time.Millisecond):
	}
	writeBEMResponse(t, gatewayConnection, actisense.BEMOperatingMode, actisense.OperatingModeSet(actisense.ModeTransferNormal))

	setMode := <-datagrams
	require.Len(t, setMode.Payload, 3)
	assert.Equal(t, uint16(actisense.ModeCANPacket), binary.LittleEndian.Uint16(setMode.Payload[1:3]))
	writeBEMResponse(t, gatewayConnection, actisense.BEMOperatingMode, actisense.OperatingModeSet(actisense.ModeCANPacket))
	select {
	case <-bus.Ready():
	case <-time.After(time.Second):
		t.Fatal("raw Bus did not become ready after acknowledgement")
	}

	inbound := can.Frame{ID: framer.BuildCANID(59904, 6, 0x20, 0x31), Length: 3, Data: [8]byte{0, 0xEE, 0}}
	wire, err := actisense.EncodeCANFrame(inbound, actisense.DirectionReceived, 10*time.Millisecond, 0)
	require.NoError(t, err)
	_, err = gatewayConnection.Write(wire)
	require.NoError(t, err)
	select {
	case frame := <-frames:
		assert.Equal(t, inbound, frame)
	case <-time.After(time.Second):
		t.Fatal("raw Bus did not deliver BST-95 frame")
	}

	outbound := can.Frame{ID: framer.BuildCANID(60928, 6, 0x2A, 255), Length: 8, Data: [8]byte{1, 2, 3, 4, 5, 6, 7, 8}}
	require.NoError(t, bus.WriteFrame(outbound))
	written := <-datagrams
	decoded, ok, err := actisense.DecodeCANFrame(written)
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, outbound, decoded.Frame)
	assert.Equal(t, actisense.DirectionTransmitted, decoded.Direction)

	closeErr := make(chan error, 1)
	go func() { closeErr <- bus.Close() }()
	restore := <-datagrams
	require.Equal(t, byte(actisense.BEMOperatingMode), restore.Payload[0])
	assert.Equal(t, uint16(actisense.ModeTransferNormal), binary.LittleEndian.Uint16(restore.Payload[1:3]))
	writeBEMResponse(t, gatewayConnection, actisense.BEMOperatingMode, actisense.OperatingModeSet(actisense.ModeTransferNormal))
	require.NoError(t, <-closeErr)
	select {
	case err := <-runErr:
		assert.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("Run did not stop after Close")
	}
}

func TestActisenseRawBusRejectsUnacknowledgedMode(t *testing.T) {
	clientConnection, gatewayConnection := net.Pipe()
	defer func() { _ = gatewayConnection.Close() }()
	bus := NewActisenseRawSerialBus(testLogger(), "fake")
	bus.stream.open = func(context.Context) (actisenseConnection, error) { return clientConnection, nil }
	bus.stream.commandTimeout = 25 * time.Millisecond
	go func() { _, _ = io.Copy(io.Discard, gatewayConnection) }()
	err := bus.Run(context.Background(), nil)
	require.Error(t, err)
	var modeErr *ActisenseModeSetupError
	require.ErrorAs(t, err, &modeErr)
	assert.Equal(t, actisense.ModeCANPacket, modeErr.RequestedMode)
	assert.Contains(t, err.Error(), "acknowledged mode-5 setup failed")
	select {
	case <-bus.Ready():
		t.Fatal("unacknowledged raw mode must not become ready")
	default:
	}
}

func TestActisenseRawBusReconnectsWithFreshHandshake(t *testing.T) {
	handshakes := make(chan int, 2)
	serveErrors := make(chan error, 2)
	opens := 0
	open := func(ctx context.Context) (actisenseConnection, error) {
		opens++
		epoch := opens
		if epoch > 2 {
			<-ctx.Done()
			return nil, ctx.Err()
		}
		clientConnection, gatewayConnection := net.Pipe()
		go func() {
			defer func() { _ = gatewayConnection.Close() }()
			parser := actisense.NewParser()
			buf := make([]byte, 4096)
			finished := false
			for !finished {
				n, err := gatewayConnection.Read(buf)
				if n > 0 {
					parser.Feed(buf[:n], func(datagram actisense.Datagram) {
						if datagram.ID != actisense.BSTBEMCommand || len(datagram.Payload) == 0 || datagram.Payload[0] != actisense.BEMOperatingMode {
							return
						}
						mode := actisense.ModeTransferNormal
						if len(datagram.Payload) == 1 {
							handshakes <- epoch
						} else if len(datagram.Payload) >= 3 {
							mode = actisense.OperatingMode(binary.LittleEndian.Uint16(datagram.Payload[1:3]))
						}
						payload := make([]byte, 12, 14)
						payload[0] = actisense.BEMOperatingMode
						payload = append(payload, actisense.OperatingModeSet(mode)...)
						wire, encodeErr := actisense.EncodeDatagram(actisense.BSTBEMResponse, payload)
						if encodeErr != nil {
							serveErrors <- encodeErr
							finished = true
							return
						}
						if _, writeErr := gatewayConnection.Write(wire); writeErr != nil {
							serveErrors <- writeErr
							finished = true
							return
						}
						if mode == actisense.ModeCANPacket {
							frame := can.Frame{
								ID:     framer.BuildCANID(127250, 2, uint8(0x20+epoch), 255),
								Length: 8,
							}
							frameWire, frameErr := actisense.EncodeCANFrame(frame, actisense.DirectionReceived, 0, 0)
							if frameErr != nil {
								serveErrors <- frameErr
							} else if _, frameErr = gatewayConnection.Write(frameWire); frameErr != nil {
								serveErrors <- frameErr
							}
							finished = true
						}
					}, func(decodeErr actisense.DecodeError) {
						serveErrors <- decodeErr
						finished = true
					})
				}
				if err != nil && !errors.Is(err, io.EOF) {
					serveErrors <- err
					return
				}
			}
		}()
		return clientConnection, nil
	}

	stream := newActisenseStreamBus(testLogger(), "fake", "fake", open, &ReconnectPolicy{
		InitialBackoff: time.Millisecond,
		MaxBackoff:     5 * time.Millisecond,
	}, actisense.ModeCANPacket, true)
	bus := &ActisenseRawTCPBus{stream: stream}
	frames := make(chan can.Frame, 2)
	runErr := make(chan error, 1)
	go func() { runErr <- bus.Run(context.Background(), func(frame can.Frame) { frames <- frame }) }()

	for epoch := 1; epoch <= 2; epoch++ {
		select {
		case gotEpoch := <-handshakes:
			assert.Equal(t, epoch, gotEpoch)
		case err := <-serveErrors:
			t.Fatalf("fake gateway failed: %v", err)
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for Actisense epoch handshake")
		}
		select {
		case frame := <-frames:
			id := framer.ParseCANID(frame.ID)
			assert.Equal(t, uint8(0x20+epoch), id.Source)
		case err := <-serveErrors:
			t.Fatalf("fake gateway failed: %v", err)
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for Actisense epoch frame")
		}
	}

	require.NoError(t, bus.Close())
	select {
	case err := <-runErr:
		assert.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("Actisense reconnect loop did not stop")
	}
}
