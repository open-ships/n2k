package n2k

import (
	"context"
	"encoding/binary"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/internal/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func commandedAddressPayload(name uint64, address uint8) []byte {
	payload := make([]byte, 9)
	binary.LittleEndian.PutUint64(payload[:8], name)
	payload[8] = address
	return payload
}

func transportReceiveFrames(pgnNum uint32, source, destination uint8, payload []byte, bam bool) []can.Frame {
	packetCount := (len(payload) + transport.MaxDTDataBytes - 1) / transport.MaxDTDataBytes
	var announcement [8]byte
	if bam {
		announcement[0] = transport.ControlBAM
		announcement[4] = 0xFF
		destination = framer.BroadcastAddr
	} else {
		announcement[0] = transport.ControlRTS
		announcement[4] = uint8(packetCount)
	}
	binary.LittleEndian.PutUint16(announcement[1:3], uint16(len(payload)))
	announcement[3] = uint8(packetCount)
	announcement[5] = byte(pgnNum)
	announcement[6] = byte(pgnNum >> 8)
	announcement[7] = byte(pgnNum >> 16)

	frames := []can.Frame{{
		ID:     framer.BuildCANID(transport.PGNCM, transport.TPPriority, source, destination),
		Length: 8,
		Data:   announcement,
	}}
	for packet := 0; packet < packetCount; packet++ {
		data := [8]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
		data[0] = uint8(packet + 1)
		start := packet * transport.MaxDTDataBytes
		end := min(start+transport.MaxDTDataBytes, len(payload))
		copy(data[1:], payload[start:end])
		frames = append(frames, can.Frame{
			ID:     framer.BuildCANID(transport.PGNDT, transport.TPPriority, source, destination),
			Length: 8,
			Data:   data,
		})
	}
	return frames
}

func deliverFrames(c *Client, frames []can.Frame) {
	for _, frame := range frames {
		c.handleBusFrame(frame)
	}
}

type commandedAddressFailBus struct {
	*mockBus
	mu          sync.Mutex
	failAddress uint8
	fail        bool
	err         error
}

func (b *commandedAddressFailBus) failClaimAt(address uint8, err error) {
	b.mu.Lock()
	b.failAddress = address
	b.fail = true
	b.err = err
	b.mu.Unlock()
}

func (b *commandedAddressFailBus) WriteFrame(frame can.Frame) error {
	b.mu.Lock()
	fail := b.fail
	failAddress := b.failAddress
	err := b.err
	b.mu.Unlock()
	info := framer.ParseCANID(frame.ID)
	if fail && info.PGN == framer.PGNISOAddressClaim && info.Source == failAddress {
		return err
	}
	return b.mockBus.WriteFrame(frame)
}

func TestConformanceCommandedAddressMatchingBAMReclaimsAddress(t *testing.T) {
	c, bus, oldAddress := newCitizenClient(t)
	newAddress := uint8(42)
	require.NotEqual(t, oldAddress, newAddress)
	initialClaims := len(framesWithPGN(bus.getWritten(), framer.PGNISOAddressClaim))

	deliverFrames(c, transportReceiveFrames(
		pgnISOCommandedAddress,
		0x33,
		framer.BroadcastAddr,
		commandedAddressPayload(c.deviceName, newAddress),
		true,
	))

	status := c.Status()
	assert.Equal(t, newAddress, status.Address)
	assert.True(t, status.AddressClaimed)
	claims := framesWithPGN(bus.getWritten(), framer.PGNISOAddressClaim)
	require.Len(t, claims, initialClaims+1)
	lastClaim := claims[len(claims)-1]
	assert.Equal(t, newAddress, framer.ParseCANID(lastClaim.ID).Source)
	assert.Equal(t, c.deviceName, binary.LittleEndian.Uint64(lastClaim.Data[:8]))

	// Application traffic remains behind a fresh contention window, while the
	// protocol claim above is emitted immediately.
	c.mu.Lock()
	ready := c.txReady
	c.mu.Unlock()
	select {
	case <-ready:
		t.Fatal("application writes became ready before the commanded-address contention window")
	default:
	}
	require.Eventually(t, func() bool {
		select {
		case <-ready:
			return true
		default:
			return false
		}
	}, time.Second, 5*time.Millisecond)
}

func TestConformanceCommandedAddressRequiresExactTransferAndTarget(t *testing.T) {
	tests := []struct {
		name        string
		destination uint8
		payload     func(*Client, uint8) []byte
	}{
		{
			name:        "addressed transfer",
			destination: 251,
			payload: func(c *Client, address uint8) []byte {
				return commandedAddressPayload(c.deviceName, address)
			},
		},
		{
			name:        "short payload",
			destination: framer.BroadcastAddr,
			payload: func(c *Client, address uint8) []byte {
				return commandedAddressPayload(c.deviceName, address)[:8]
			},
		},
		{
			name:        "long payload",
			destination: framer.BroadcastAddr,
			payload: func(c *Client, address uint8) []byte {
				return append(commandedAddressPayload(c.deviceName, address), 0xFF)
			},
		},
		{
			name:        "NAME differs only in high bit",
			destination: framer.BroadcastAddr,
			payload: func(c *Client, address uint8) []byte {
				return commandedAddressPayload(c.deviceName^(uint64(1)<<63), address)
			},
		},
	}
	for _, address := range []uint8{252, 253, 254, 255} {
		commandedAddress := address
		tests = append(tests, struct {
			name        string
			destination uint8
			payload     func(*Client, uint8) []byte
		}{
			name:        "special address " + strconv.Itoa(int(commandedAddress)),
			destination: framer.BroadcastAddr,
			payload: func(c *Client, _ uint8) []byte {
				return commandedAddressPayload(c.deviceName, commandedAddress)
			},
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, bus, oldAddress := newCitizenClient(t)
			initialClaims := len(framesWithPGN(bus.getWritten(), framer.PGNISOAddressClaim))
			c.handleCommandedAddressTransfer(tt.destination, tt.payload(c, 42))

			assert.Equal(t, oldAddress, c.Status().Address)
			assert.Len(t, framesWithPGN(bus.getWritten(), framer.PGNISOAddressClaim), initialClaims)
		})
	}
}

func TestConformanceCommandedAddressRejectsFastPacketAndAddressedTP(t *testing.T) {
	t.Run("fast packet", func(t *testing.T) {
		c, bus, oldAddress := newCitizenClient(t)
		initialClaims := len(framesWithPGN(bus.getWritten(), framer.PGNISOAddressClaim))
		frames := framer.FrameFastPacket(
			framer.BuildCANID(pgnISOCommandedAddress, 6, 0x33, framer.BroadcastAddr),
			commandedAddressPayload(c.deviceName, 42),
			0,
		)
		deliverFrames(c, frames)

		assert.Equal(t, oldAddress, c.Status().Address)
		assert.Len(t, framesWithPGN(bus.getWritten(), framer.PGNISOAddressClaim), initialClaims)
	})

	t.Run("addressed transport", func(t *testing.T) {
		c, bus, oldAddress := newCitizenClient(t)
		initialClaims := len(framesWithPGN(bus.getWritten(), framer.PGNISOAddressClaim))
		deliverFrames(c, transportReceiveFrames(
			pgnISOCommandedAddress,
			0x33,
			oldAddress,
			commandedAddressPayload(c.deviceName, 42),
			false,
		))

		assert.Equal(t, oldAddress, c.Status().Address)
		assert.Len(t, framesWithPGN(bus.getWritten(), framer.PGNISOAddressClaim), initialClaims)
	})
}

func TestConformanceCommandedAddressClaimFailureTerminatesClient(t *testing.T) {
	bus := &commandedAddressFailBus{mockBus: newMockBus()}
	c, err := NewClient(context.Background(),
		WithBus(bus),
		WithClaimTimeout(50*time.Millisecond),
		WithHeartbeatInterval(0),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })

	writeErr := errors.New("injected commanded-address claim failure")
	bus.failClaimAt(42, writeErr)
	deliverFrames(c, transportReceiveFrames(
		pgnISOCommandedAddress,
		0x33,
		framer.BroadcastAddr,
		commandedAddressPayload(c.deviceName, 42),
		true,
	))

	require.ErrorIs(t, c.Err(), writeErr)
	assert.Equal(t, uint8(42), c.Status().Address)
}
