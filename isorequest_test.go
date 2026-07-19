package n2k

import (
	"context"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// isoRequestFrame builds a single-frame ISO request (PGN 59904) for
// requestedPGN, from source to dest.
func isoRequestFrame(requested uint32, source, dest uint8) can.Frame {
	payload := []byte{byte(requested), byte(requested >> 8), byte(requested >> 16)}
	return framer.FrameSingle(framer.BuildCANID(59904, 6, source, dest), payload)
}

// assembleFast reassembles fast-packet frames into a single payload.
func assembleFast(t *testing.T, frames []can.Frame) []byte {
	t.Helper()
	require.NotEmpty(t, frames)
	total := int(frames[0].Data[1])
	payload := append([]byte{}, frames[0].Data[2:8]...)
	for _, f := range frames[1:] {
		payload = append(payload, f.Data[1:8]...)
	}
	require.GreaterOrEqual(t, len(payload), total)
	return payload[:total]
}

// newCitizenClient builds a mockBus client with fast claim and no heartbeat
// noise, returning the client, bus, and the client's claimed address.
func newCitizenClient(t *testing.T, opts ...Option) (*Client, *mockBus, uint8) {
	t.Helper()
	mb := newMockBus()
	base := []Option{
		WithBus(mb),
		WithClaimTimeout(50 * time.Millisecond),
		WithHeartbeatInterval(0),
	}
	c, err := NewClient(context.Background(), append(base, opts...)...)
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })

	c.mu.Lock()
	addr := c.sourceAddr
	c.mu.Unlock()
	return c, mb, addr
}

func TestISORequest_ProductInfoBroadcast(t *testing.T) {
	_, mb, _ := newCitizenClient(t, WithProductInfo(ProductInfo{
		ProductCode: 1234,
		ModelID:     "gateway mk2",
	}))

	mb.inbound <- isoRequestFrame(126996, 0x42, 255)

	ok := waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 126996)) > 0
	})
	require.True(t, ok, "expected a ProductInformation response")

	payload := assembleFast(t, framesWithPGN(mb.getWritten(), 126996))
	msg, err := pgn.DecodeMessage(pgn.MessageInfo{PGN: 126996}, payload)
	require.NoError(t, err)
	pi, isPI := msg.(*pgn.ProductInformation)
	require.True(t, isPI, "expected *pgn.ProductInformation, got %T", msg)
	assert.Equal(t, "gateway mk2", pi.ModelId)
	require.NotNil(t, pi.ProductCode)
	assert.Equal(t, uint64(1234), *pi.ProductCode)
}

func TestISORequest_ProductInfoDefaults(t *testing.T) {
	_, mb, _ := newCitizenClient(t)

	mb.inbound <- isoRequestFrame(126996, 0x42, 255)

	ok := waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 126996)) > 0
	})
	require.True(t, ok, "product info must be answered even when WithProductInfo is not set")
}

func TestISORequest_ConfigInfo(t *testing.T) {
	_, mb, addr := newCitizenClient(t, WithConfigInfo(ConfigInfo{
		InstallationDescription1: "port engine room",
	}))

	mb.inbound <- isoRequestFrame(126998, 0x42, addr)

	ok := waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 126998)) > 0
	})
	require.True(t, ok, "expected a ConfigurationInformation response")

	payload := assembleFast(t, framesWithPGN(mb.getWritten(), 126998))
	msg, err := pgn.DecodeMessage(pgn.MessageInfo{PGN: 126998}, payload)
	require.NoError(t, err)
	ci, isCI := msg.(*pgn.ConfigurationInformation)
	require.True(t, isCI, "expected *pgn.ConfigurationInformation, got %T", msg)
	assert.Equal(t, "port engine room", ci.InstallationDescription1)
}

func TestISORequest_HeartbeatOnDemand(t *testing.T) {
	_, mb, _ := newCitizenClient(t) // heartbeats disabled in helper

	mb.inbound <- isoRequestFrame(126993, 0x42, 255)

	ok := waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 126993)) == 1
	})
	assert.True(t, ok, "a request for 126993 should trigger exactly one heartbeat")
}

func TestISORequest_AddressedToOtherNodeIgnored(t *testing.T) {
	_, mb, addr := newCitizenClient(t)
	otherNode := addr - 1

	mb.inbound <- isoRequestFrame(126996, 0x42, otherNode)

	time.Sleep(150 * time.Millisecond)
	assert.Empty(t, framesWithPGN(mb.getWritten(), 126996))
}

func TestISORequest_UnsupportedAddressedGetsNAK(t *testing.T) {
	c, mb, addr := newCitizenClient(t)

	mb.inbound <- isoRequestFrame(130306, 0x42, addr)

	ok := waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 59392)) > 0
	})
	require.True(t, ok, "expected an IsoAcknowledgement NAK; client error: %v", c.Err())

	nakFrame := framesWithPGN(mb.getWritten(), 59392)[0]
	// Wire layout: control (1B), group function (1B), reserved (3B),
	// refused PGN (3B little-endian).
	assert.Equal(t, uint8(pgn.NAK), nakFrame.Data[0], "control must be NAK")
	assert.Equal(t, uint8(0xFF), nakFrame.Data[1], "group function must be the null sentinel")
	refused := uint32(nakFrame.Data[5]) | uint32(nakFrame.Data[6])<<8 | uint32(nakFrame.Data[7])<<16
	assert.Equal(t, uint32(130306), refused)
}

func TestISORequestNAKEncodes(t *testing.T) {
	_, err := pgn.EncodeMessage(nakFor(130306))
	require.NoError(t, err)
}

func TestISORequest_UnsupportedBroadcastIgnored(t *testing.T) {
	_, mb, _ := newCitizenClient(t)

	mb.inbound <- isoRequestFrame(130306, 0x42, 255)

	time.Sleep(150 * time.Millisecond)
	assert.Empty(t, framesWithPGN(mb.getWritten(), 59392), "broadcast requests for unsupported PGNs must not be NAKed")
}

func TestISORequest_ClaimRequestForOtherNodeIgnored(t *testing.T) {
	_, mb, addr := newCitizenClient(t)
	initialClaims := len(framesWithPGN(mb.getWritten(), 60928))

	// Addressed to another node: our claimer must stay silent.
	mb.inbound <- isoRequestFrame(60928, 0x42, addr-1)
	time.Sleep(150 * time.Millisecond)
	assert.Equal(t, initialClaims, len(framesWithPGN(mb.getWritten(), 60928)))

	// Broadcast: we must answer with a fresh claim.
	mb.inbound <- isoRequestFrame(60928, 0x42, 255)
	ok := waitFor(t, 2*time.Second, func() bool {
		return len(framesWithPGN(mb.getWritten(), 60928)) > initialClaims
	})
	assert.True(t, ok, "broadcast claim requests must still be answered")
}
