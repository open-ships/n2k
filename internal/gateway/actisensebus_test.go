package gateway

import (
	"bytes"
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
	"github.com/open-ships/n2k/raw"
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

func TestActisenseHandshakeTimeoutAndCloseInterruptInitialWrite(t *testing.T) {
	for _, closeDuringHandshake := range []bool{false, true} {
		name := "deadline"
		if closeDuringHandshake {
			name = "close"
		}
		t.Run(name, func(t *testing.T) {
			host, peer := net.Pipe()
			defer func() { _ = peer.Close() }()
			opened := make(chan struct{})
			bus := newActisenseStreamBus(testLogger(), "stalled", "stalled", func(context.Context) (actisenseConnection, error) {
				close(opened)
				return host, nil
			}, nil, actisense.ModeCANPacket, true)
			t.Cleanup(func() { _ = bus.Close() })
			timeout := 20 * time.Millisecond
			if closeDuringHandshake {
				timeout = time.Second
			}
			require.NoError(t, bus.SetCommandTimeout(timeout))
			result := make(chan error, 1)
			go func() { result <- bus.Run(context.Background(), nil) }()
			<-opened
			if closeDuringHandshake {
				require.NoError(t, bus.Close())
			}
			select {
			case err := <-result:
				if closeDuringHandshake {
					require.NoError(t, err)
				} else {
					require.ErrorIs(t, err, context.DeadlineExceeded)
				}
			case <-time.After(500 * time.Millisecond):
				t.Fatal("initial write ignored handshake timeout or Close")
			}
			select {
			case <-bus.Ready():
				t.Fatal("unacknowledged connection became ready")
			default:
			}
		})
	}
}

func TestActisenseEpochRequesterNeverFollowsReconnect(t *testing.T) {
	gateway := NewActisenseCustomGatewaySession(testLogger(), "epoch", "epoch", nil, nil, actisense.ModeTransferNormal)
	first := actisense.NewSession(actisense.SessionConfig{})
	gateway.stream.epoch = &actisenseEpoch{number: 1, session: first}
	requester, err := gateway.EpochRequester(1)
	require.NoError(t, err)
	first.Close(actisense.ErrSessionClosed)

	second := actisense.NewSession(actisense.SessionConfig{Write: func(_ context.Context, data []byte) (int, error) {
		t.Error("old epoch cleanup reached the new connection")
		return len(data), nil
	}})
	defer second.Close(nil)
	gateway.stream.epoch = &actisenseEpoch{number: 2, session: second}
	_, err = requester.Request(context.Background(), actisense.BEMTxPGNEnable, nil)
	require.ErrorIs(t, err, actisense.ErrSessionClosed)
	_, err = gateway.EpochRequester(1)
	require.Error(t, err)
	require.Zero(t, second.Metrics().BEMRequests)
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
	frames := make(chan raw.Observation, 2)
	runErr := make(chan error, 1)
	go func() {
		runErr <- bus.RunObservations(context.Background(), func(observation raw.Observation) {
			if observation.Kind == raw.KindFrame {
				frames <- observation
			}
		})
	}()

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
		case observation := <-frames:
			assert.Equal(t, uint64(epoch), observation.ConnectionEpoch)
			id := framer.ParseCANID(observation.Frame.ID)
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

func TestActisenseRawBusRequiresAcknowledgedFormatWhileCapturePreservesRecords(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mode   actisense.OperatingMode
		rawBus bool
		want   []uint32
	}{
		{"binary bus", actisense.ModeCANPacket, true, []uint32{2}},
		{"binary capture", actisense.ModeCANPacket, false, []uint32{1, 127250, 130900, 2}},
		{"ASCII bus", actisense.ModeCANPacketASCII, true, []uint32{3}},
		{"ASCII capture", actisense.ModeCANPacketASCII, false, []uint32{1, 127250, 130900, 2, 3}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			host, peer := net.Pipe()
			t.Cleanup(func() { _ = peer.Close(); _ = host.Close() })
			require.NoError(t, peer.SetDeadline(time.Now().Add(3*time.Second)))
			bus := newActisenseStreamBus(testLogger(), "test", "test", nil, nil, tc.mode, tc.rawBus)
			var observations []raw.Observation
			result := make(chan error, 1)
			go func() {
				err, _ := bus.runEpoch(context.Background(), host, 42, func(o raw.Observation) {
					observations = append(observations, o)
				})
				result <- err
			}()
			get, err := actisense.EncodeBEMRequest(actisense.BEMOperatingMode, nil)
			require.NoError(t, err)
			got := make([]byte, len(get))
			_, err = io.ReadFull(peer, got)
			require.NoError(t, err)
			require.Equal(t, get, got)
			before, err := actisense.EncodeCANFrame(can.Frame{ID: 1}, actisense.DirectionReceived, 0, 0)
			require.NoError(t, err)
			after, err := actisense.EncodeCANFrame(can.Frame{ID: 2}, actisense.DirectionReceived, 0, 0)
			require.NoError(t, err)
			modeResponse := make([]byte, 12, 14)
			modeResponse[0] = actisense.BEMOperatingMode
			modeResponse = append(modeResponse, actisense.OperatingModeSet(tc.mode)...)
			ack, err := actisense.EncodeDatagram(actisense.BSTBEMResponse, modeResponse)
			require.NoError(t, err)
			assembled93 := wrapActisense(cmdN2KReceived, n2kPayload(2, 127250, 255, 42, []byte{1, 2}))
			assembledD0, err := actisense.EncodeMessageD0(actisense.Message{PGN: 130900, Destination: 255, Source: 42})
			require.NoError(t, err)
			asciiFrame, err := FormatActisenseCANASCII(can.Frame{ID: 3}, raw.DirectionReceived, 0)
			require.NoError(t, err)
			_, err = peer.Write(bytes.Join([][]byte{before, ack, assembled93, assembledD0, after, asciiFrame}, nil))
			require.NoError(t, err)
			select {
			case <-bus.Ready():
			case err := <-result:
				t.Fatalf("epoch stopped during handshake: %v", err)
			case <-time.After(3 * time.Second):
				t.Fatal("epoch did not publish acknowledged mode")
			}
			require.NoError(t, peer.Close())
			require.NoError(t, <-result)
			var networkRecords []uint32
			for _, o := range observations {
				require.Equal(t, uint64(42), o.ConnectionEpoch)
				if o.Kind == raw.KindMessage {
					networkRecords = append(networkRecords, o.PGN)
				} else if o.Frame != nil {
					networkRecords = append(networkRecords, o.Frame.ID)
				}
			}
			assert.Equal(t, tc.want, networkRecords)
		})
	}
}

func TestActisenseWireTraceMatchesPartialPhysicalWrite(t *testing.T) {
	host, peer := net.Pipe()
	defer func() { _ = host.Close() }()
	defer func() { _ = peer.Close() }()
	peerRead := make(chan []byte, 1)
	go func() {
		prefix := make([]byte, 2)
		n, _ := io.ReadFull(peer, prefix)
		_ = peer.Close()
		peerRead <- prefix[:n]
	}()
	var traced []byte
	session := actisense.NewSession(actisense.SessionConfig{
		Write: func(ctx context.Context, data []byte) (int, error) {
			return writeActisenseUnitContext(ctx, host, data)
		},
		OnWireBytes: func(direction actisense.WireDirection, _ time.Time, data []byte) {
			require.Equal(t, actisense.WireTransmitted, direction)
			traced = append(traced, data...)
		},
	})
	defer session.Close(nil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	require.ErrorIs(t, session.WriteContext(ctx, []byte{0x10, 0x02, 0xA1, 0x01}), io.ErrClosedPipe)
	require.Equal(t, []byte{0x10, 0x02}, traced)
	require.Equal(t, <-peerRead, traced)
	require.Equal(t, uint64(2), session.Metrics().TransportWriteBytes)
}
