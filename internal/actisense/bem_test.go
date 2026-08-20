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
	t.Helper()
	payload := make([]byte, 12, 12+len(data))
	payload[0] = command
	payload[1] = 7
	binary.LittleEndian.PutUint16(payload[2:4], 0x1234)
	binary.LittleEndian.PutUint32(payload[4:8], 0x55667788)
	binary.LittleEndian.PutUint32(payload[8:12], uint32(errorCode))
	payload = append(payload, data...)
	wire, err := EncodeDatagram(BSTBEMResponse, payload)
	require.NoError(t, err)
	return wire
}

func TestSessionOperatingModeAcknowledged(t *testing.T) {
	reader, writer := io.Pipe()
	defer func() { _ = reader.Close() }()
	defer func() { _ = writer.Close() }()

	var writeMu sync.Mutex
	session := NewSession(SessionConfig{Write: func(buf []byte) error {
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
	session := NewSession(SessionConfig{Write: func([]byte) error { return nil }})
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
	session := NewSession(SessionConfig{Write: func([]byte) error {
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
	session := NewSession(SessionConfig{Write: func([]byte) error {
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
