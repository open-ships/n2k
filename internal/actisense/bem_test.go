package actisense

import (
	"context"
	"encoding/binary"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func bemResponseDatagram(t *testing.T, command byte, errorCode int32, data []byte) []byte {
	return bemResponseDatagramWithBST(t, BSTBEMResponse, command, errorCode, data)
}

func bemResponseDatagramWithBST(t *testing.T, bstID, command byte, errorCode int32, data []byte) []byte {
	t.Helper()
	payload := make([]byte, 12, 12+len(data))
	payload[0] = command
	payload[1] = 7
	binary.LittleEndian.PutUint16(payload[2:4], 0x1234)
	binary.LittleEndian.PutUint32(payload[4:8], 0x55667788)
	binary.LittleEndian.PutUint32(payload[8:12], uint32(errorCode))
	payload = append(payload, data...)
	wire, err := EncodeDatagram(bstID, payload)
	require.NoError(t, err)
	return wire
}

func TestBEMRequestIndependentGoldenVectors(t *testing.T) {
	get, err := EncodeBEMRequest(BEMOperatingMode, nil)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x10, 0x02, 0xA1, 0x01, 0x11, 0x4D, 0x10, 0x03}, get)
	set, err := EncodeBEMRequest(BEMOperatingMode, []byte{0x05, 0x00})
	require.NoError(t, err)
	assert.Equal(t, []byte{0x10, 0x02, 0xA1, 0x03, 0x11, 0x05, 0x00, 0x46, 0x10, 0x03}, set)
}

func TestSessionRequiresExpectedResponseGroup(t *testing.T) {
	reader, writer := io.Pipe()
	defer func() { _ = reader.Close() }()
	session := NewSession(SessionConfig{Write: func(context.Context, []byte) error {
		if _, err := writer.Write(bemResponseDatagramWithBST(t, 0xA2, BEMOperatingMode, 0, []byte{1, 0})); err != nil {
			return err
		}
		_, err := writer.Write(bemResponseDatagram(t, BEMOperatingMode, 0, []byte{5, 0}))
		return err
	}})
	go func() { _ = session.Run(reader) }()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	mode, err := session.GetOperatingMode(ctx)
	require.NoError(t, err)
	assert.Equal(t, ModeCANPacket, mode)
	session.Close(nil)
	_ = writer.Close()
}

func TestSessionCollectsBoundedMultiResponseTrain(t *testing.T) {
	reader, writer := io.Pipe()
	defer func() { _ = reader.Close() }()
	session := NewSession(SessionConfig{Write: func(context.Context, []byte) error {
		for _, data := range [][]byte{{1}, {2}} {
			if _, err := writer.Write(bemResponseDatagram(t, BEMOperatingMode, 0, data)); err != nil {
				return err
			}
		}
		return nil
	}})
	go func() { _ = session.Run(reader) }()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	responses, err := session.RequestMulti(ctx, BEMOperatingMode, nil, 100*time.Millisecond, func(responses []BEMResponse) (bool, error) {
		return len(responses) == 2, nil
	})
	require.NoError(t, err)
	require.Len(t, responses, 2)
	assert.Equal(t, []byte{1}, responses[0].Data)
	assert.Equal(t, []byte{2}, responses[1].Data)
	session.Close(nil)
	_ = writer.Close()
}

func TestSessionResponse257TerminatesWithoutBlockingSoleReader(t *testing.T) {
	reader, writer := io.Pipe()
	defer func() { _ = reader.Close() }()
	session := NewSession(SessionConfig{Write: func(context.Context, []byte) error {
		for index := 0; index <= maxBEMResponseTrain; index++ {
			if _, err := writer.Write(bemResponseDatagram(t, BEMProductInfo, 0, []byte{byte(index)})); err != nil {
				return err
			}
		}
		return nil
	}})
	go func() { _ = session.Run(reader) }()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	responses, err := session.RequestMulti(ctx, BEMProductInfo, nil, time.Second, func([]BEMResponse) (bool, error) {
		return false, nil
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds 256")
	assert.Len(t, responses, 256)
	assert.Equal(t, uint64(1), session.Metrics().BEMResponseTrainOverflows)
	session.Close(nil)
	_ = writer.Close()
}

func TestSessionOperatingModeAcknowledged(t *testing.T) {
	reader, writer := io.Pipe()
	defer func() { _ = reader.Close() }()
	defer func() { _ = writer.Close() }()

	var writeMu sync.Mutex
	session := NewSession(SessionConfig{Write: func(_ context.Context, buf []byte) error {
		writeMu.Lock()
		defer writeMu.Unlock()
		_, err := writer.Write(bemResponseDatagram(t, BEMOperatingMode, 0, []byte{5, 0}))
		return err
	}})
	runErr := make(chan error, 1)
	go func() { runErr <- session.Run(reader) }()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	require.NoError(t, session.SetOperatingMode(ctx, ModeCANPacket))
	session.Close(nil)
	_ = writer.Close()
	<-runErr
}

func TestSessionRejectsSameVerbCollision(t *testing.T) {
	session := NewSession(SessionConfig{Write: func(context.Context, []byte) error { return nil }})
	ctx, cancel := context.WithCancel(context.Background())
	first := make(chan error, 1)
	go func() {
		_, err := session.Request(ctx, BEMOperatingMode, nil)
		first <- err
	}()
	require.Eventually(t, func() bool {
		session.mu.Lock()
		defer session.mu.Unlock()
		return len(session.pending) == 1
	}, time.Second, time.Millisecond)
	_, err := session.Request(context.Background(), BEMOperatingMode, nil)
	assert.ErrorIs(t, err, ErrRequestInFlight)
	cancel()
	assert.ErrorIs(t, <-first, context.Canceled)
}

func TestSessionSurfacesDeviceError(t *testing.T) {
	reader, writer := io.Pipe()
	defer func() { _ = reader.Close() }()
	session := NewSession(SessionConfig{Write: func(context.Context, []byte) error {
		_, err := writer.Write(bemResponseDatagram(t, BEMActivatePGNLists, -1159, nil))
		return err
	}})
	go func() { _ = session.Run(reader) }()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := session.ActivatePGNLists(ctx)
	var deviceErr *DeviceError
	require.ErrorAs(t, err, &deviceErr)
	assert.Equal(t, int32(-1159), deviceErr.Code)
	session.Close(nil)
	_ = writer.Close()
}

func TestDecodeDiagnostics(t *testing.T) {
	for _, test := range []struct {
		command byte
		data    []byte
		kind    DiagnosticKind
	}{
		{BEMStartupStatus, []byte{5, 0, 1, 0, 0, 0}, DiagnosticStartup},
		{BEMErrorReport, []byte{1, 0, 0, 0, 2, 0, 0, 0}, DiagnosticError},
		{BEMSystemStatus, []byte{1, 1, 2, 3, 4, 5, 6, 0, 0, 0, 0}, DiagnosticSystem},
		{BEMNegativeAck, []byte{1, 2, 3, 4}, DiagnosticNegativeAck},
	} {
		datagram := decodeOne(t, bemResponseDatagram(t, test.command, 0, test.data))
		response, ok, err := DecodeBEMResponse(datagram)
		require.NoError(t, err)
		require.True(t, ok)
		diagnostic, ok, err := DecodeDiagnostic(response)
		require.NoError(t, err)
		require.True(t, ok)
		assert.Equal(t, test.kind, diagnostic.Kind)
	}
}

func TestSessionNegativeAckFailsExactRequest(t *testing.T) {
	reader, writer := io.Pipe()
	defer func() { _ = reader.Close() }()
	session := NewSession(SessionConfig{Write: func(context.Context, []byte) error {
		data := make([]byte, 4)
		binary.LittleEndian.PutUint32(data, uint32(BEMOperatingMode))
		_, err := writer.Write(bemResponseDatagram(t, BEMNegativeAck, -1159, data))
		return err
	}})
	go func() { _ = session.Run(reader) }()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := session.Request(ctx, BEMOperatingMode, nil)
	var nack *NegativeAckError
	require.ErrorAs(t, err, &nack)
	assert.Equal(t, byte(BEMOperatingMode), nack.Command)
	assert.Equal(t, int32(-1159), nack.DeviceCode)
	session.Close(nil)
	_ = writer.Close()
}

func TestNegativeAckFallbackNeverCrossesResponseGroup(t *testing.T) {
	session := NewSession(SessionConfig{})
	local := &pendingRequest{results: make(chan responseResult, 1)}
	debug := &pendingRequest{results: make(chan responseResult, 1)}
	session.pending[responseKey{bstID: BSTBEMResponse, bemID: BEMEcho, origin: LocalBEMOrigin}] = local
	session.pending[responseKey{bstID: 0xA2, bemID: BEMOperatingMode, origin: LocalBEMOrigin}] = debug
	session.failNegativeAck(BEMResponse{BSTID: 0xA2, Origin: LocalBEMOrigin}, NegativeAck{UniqueCommandID: 0xDEADBEEF, ErrorCode: -1})

	select {
	case result := <-debug.results:
		var nack *NegativeAckError
		require.ErrorAs(t, result.err, &nack)
		assert.Equal(t, byte(BEMOperatingMode), nack.Command)
	default:
		t.Fatal("same-group sole fallback was not failed")
	}
	select {
	case <-local.results:
		t.Fatal("negative acknowledgement crossed from A2 into local A0 request")
	default:
	}
}

func TestSessionMetricsCountUnframedDecodeAndTimeout(t *testing.T) {
	session := NewSession(SessionConfig{Write: func(context.Context, []byte) error { return nil }})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	_, err := session.Request(ctx, BEMEcho, nil)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	metrics := session.Metrics()
	assert.Equal(t, uint64(1), metrics.BEMRequests)
	assert.Equal(t, uint64(1), metrics.BEMCompleted)
	assert.Equal(t, uint64(1), metrics.BEMTimeouts)
	assert.Zero(t, metrics.BEMInFlight)
}
