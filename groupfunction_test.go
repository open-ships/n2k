package n2k

import (
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// requestGF builds the fast-packet frames of a request group function
// (PGN 126208, function code 0) for pgnNum. interval semantics follow the
// standard: nil = no change.
func requestGF(t *testing.T, pgnNum uint32, interval *uint32, source, dest uint8) []can.Frame {
	t.Helper()
	payload := []byte{0}                                                       // function code 0 = request
	payload = append(payload, byte(pgnNum), byte(pgnNum>>8), byte(pgnNum>>16)) // requested PGN
	iv := uint32(0xFFFFFFFF)
	if interval != nil {
		iv = *interval
	}
	payload = append(payload, byte(iv), byte(iv>>8), byte(iv>>16), byte(iv>>24))
	payload = append(payload, 0xFF, 0xFF) // interval offset: no change
	payload = append(payload, 0)          // no parameters
	canID := framer.BuildCANID(126208, 3, source, dest)
	return framer.FrameFastPacket(canID, payload, 0)
}

// requestGFWithParam is requestGF plus one parameter-selection pair.
func requestGFWithParam(t *testing.T, pgnNum uint32, source, dest uint8) []can.Frame {
	t.Helper()
	payload := []byte{0}
	payload = append(payload, byte(pgnNum), byte(pgnNum>>8), byte(pgnNum>>16))
	payload = append(payload, 0xFF, 0xFF, 0xFF, 0xFF)
	payload = append(payload, 0xFF, 0xFF)
	payload = append(payload, 1)          // one parameter pair
	payload = append(payload, 2)          // field index 2
	payload = append(payload, 0xD2, 0x04) // value 1234 (u16)
	canID := framer.BuildCANID(126208, 3, source, dest)
	return framer.FrameFastPacket(canID, payload, 0)
}

// commandGF builds the fast-packet frames of a command group function
// (function code 1) for pgnNum with no parameters.
func commandGF(t *testing.T, pgnNum uint32, source, dest uint8) []can.Frame {
	t.Helper()
	payload := []byte{1}
	payload = append(payload, byte(pgnNum), byte(pgnNum>>8), byte(pgnNum>>16))
	payload = append(payload, 0xFF) // priority: no change (4 bits) + reserved
	payload = append(payload, 0)    // no parameters
	canID := framer.BuildCANID(126208, 3, source, dest)
	return framer.FrameFastPacket(canID, payload, 0)
}

func inject(mb *mockBus, frames []can.Frame) {
	for _, f := range frames {
		mb.inbound <- f
	}
}

// decodeAck reassembles and decodes the first acknowledge group function the
// client wrote.
func decodeAck(t *testing.T, mb *mockBus) *pgn.NmeaAcknowledgeGroupFunction {
	t.Helper()
	frames := framesWithPGN(mb.getWritten(), 126208)
	require.NotEmpty(t, frames)
	payload := assembleFast(t, frames)
	msg, err := pgn.DecodeMessage(pgn.MessageInfo{PGN: 126208}, payload)
	require.NoError(t, err)
	ack, ok := msg.(*pgn.NmeaAcknowledgeGroupFunction)
	require.True(t, ok, "expected acknowledge group function, got %T", msg)
	return ack
}

func wroteAck(mb *mockBus) bool {
	return len(framesWithPGN(mb.getWritten(), 126208)) > 0
}

func u32(v uint32) *uint32 { return &v }

func TestGroupFunction_RequestSupportedPGNTransmitsOnce(t *testing.T) {
	_, mb, addr := newCitizenClient(t)

	inject(mb, requestGF(t, 126996, nil, 0x42, addr))

	ok := waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 126996)) > 0
	})
	assert.True(t, ok, "request GF for 126996 should transmit product info")
	assert.False(t, wroteAck(mb), "successful requests are not acknowledged")
}

func TestGroupFunction_RequestUnsupportedPGNAcked(t *testing.T) {
	_, mb, addr := newCitizenClient(t)

	inject(mb, requestGF(t, 130306, nil, 0x42, addr))

	ok := waitFor(t, 2*time.Second, func() bool { return wroteAck(mb) })
	require.True(t, ok, "expected an acknowledge group function")

	ack := decodeAck(t, mb)
	require.NotNil(t, ack.PgnErrorCode)
	assert.Equal(t, uint64(pgn.PgnErrorCodePGNNotSupported), *ack.PgnErrorCode)
	require.NotNil(t, ack.Pgn)
	assert.Equal(t, uint64(130306), *ack.Pgn)
}

func TestGroupFunction_BroadcastUnsupportedIgnored(t *testing.T) {
	_, mb, _ := newCitizenClient(t)

	inject(mb, requestGF(t, 130306, nil, 0x42, 255))

	time.Sleep(150 * time.Millisecond)
	assert.False(t, wroteAck(mb), "broadcast group functions for unsupported PGNs must not be NAKed")
}

func TestGroupFunction_AddressedToOtherNodeIgnored(t *testing.T) {
	_, mb, addr := newCitizenClient(t)

	inject(mb, requestGF(t, 126996, nil, 0x42, addr-1))

	time.Sleep(150 * time.Millisecond)
	assert.Empty(t, framesWithPGN(mb.getWritten(), 126996))
	assert.False(t, wroteAck(mb))
}

func TestGroupFunction_HeartbeatIntervalChange(t *testing.T) {
	c, mb, addr := newCitizenClient(t) // heartbeat disabled by helper

	inject(mb, requestGF(t, 126993, u32(2000), 0x42, addr))

	ok := waitFor(t, 2*time.Second, func() bool {
		return c.heartbeat.currentInterval() == 2*time.Second
	})
	assert.True(t, ok, "heartbeat interval should follow the group function")

	ok = waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 126993)) >= 1
	})
	assert.True(t, ok, "re-enabling the heartbeat should transmit immediately")
}

func TestGroupFunction_HeartbeatIntervalClamped(t *testing.T) {
	c, mb, addr := newCitizenClient(t)

	inject(mb, requestGF(t, 126993, u32(100), 0x42, addr))

	ok := waitFor(t, 2*time.Second, func() bool {
		return c.heartbeat.currentInterval() == time.Second
	})
	assert.True(t, ok, "sub-second heartbeat intervals must clamp to 1s")
}

func TestGroupFunction_HeartbeatStopAndRestoreDefault(t *testing.T) {
	c, mb, addr := newCitizenClient(t)

	inject(mb, requestGF(t, 126993, u32(0), 0x42, addr))
	ok := waitFor(t, 2*time.Second, func() bool {
		return c.heartbeat.currentInterval() == 0
	})
	require.True(t, ok, "interval 0 should stop the heartbeat")

	inject(mb, requestGF(t, 126993, u32(0xFFFFFFFE), 0x42, addr))
	ok = waitFor(t, 2*time.Second, func() bool {
		return c.heartbeat.currentInterval() == defaultHeartbeatInterval
	})
	assert.True(t, ok, "0xFFFFFFFE should restore the default interval")
}

func TestGroupFunction_RetimesBroadcaster(t *testing.T) {
	c, mb, addr := newCitizenClient(t)

	stop, err := c.Broadcast(time.Hour, headingProvider)
	require.NoError(t, err)
	defer stop()

	// Wait for the initial send, then retime to 20ms via group function.
	waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 127250)) >= 1
	})
	inject(mb, requestGF(t, 127250, u32(20), 0x42, addr))

	ok := waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 127250)) >= 3
	})
	assert.True(t, ok, "broadcaster should speed up to the commanded interval")
}

func TestGroupFunction_RequestWithParametersAcked(t *testing.T) {
	_, mb, addr := newCitizenClient(t)

	inject(mb, requestGFWithParam(t, 126996, 0x42, addr))

	ok := waitFor(t, 2*time.Second, func() bool { return wroteAck(mb) })
	require.True(t, ok, "parameterized requests are not supported and must be acknowledged")

	ack := decodeAck(t, mb)
	require.NotNil(t, ack.PgnErrorCode)
	assert.Equal(t, uint64(pgn.PgnErrorCodeNotSupported), *ack.PgnErrorCode)
	assert.Empty(t, framesWithPGN(mb.getWritten(), 126996), "refused requests must not transmit")
}

func TestGroupFunction_CommandAcked(t *testing.T) {
	_, mb, addr := newCitizenClient(t)

	inject(mb, commandGF(t, 126996, 0x42, addr))

	ok := waitFor(t, 2*time.Second, func() bool { return wroteAck(mb) })
	require.True(t, ok, "commands are not supported and must be acknowledged")

	ack := decodeAck(t, mb)
	require.NotNil(t, ack.PgnErrorCode)
	assert.Equal(t, uint64(pgn.PgnErrorCodeNotSupported), *ack.PgnErrorCode)

	require.NotNil(t, ack.FunctionCode)
	assert.Equal(t, uint64(2), *ack.FunctionCode, "acknowledgements carry function code 2")
}
