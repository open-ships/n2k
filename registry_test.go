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

// claimFrame builds an ISO address claim (PGN 60928) broadcast from source.
func claimFrame(name DeviceName, source uint8) can.Frame {
	packed := name.Pack(true)
	payload := make([]byte, 8)
	for i := range payload {
		payload[i] = byte(packed >> (8 * i))
	}
	return framer.FrameSingle(framer.BuildCANID(60928, 6, source, 255), payload)
}

// infoRequestsTo counts ISO requests for requested written to dest.
func infoRequestsTo(mb *mockBus, requested uint32, dest uint8) int {
	count := 0
	for _, f := range mb.getWritten() {
		id := framer.ParseCANID(f.ID)
		if id.PGN != 59904 || id.Destination != dest || f.Length < 3 {
			continue
		}
		got := uint32(f.Data[0]) | uint32(f.Data[1])<<8 | uint32(f.Data[2])<<16
		if got == requested {
			count++
		}
	}
	return count
}

var testDeviceName = DeviceName{
	IdentityNumber:   424242,
	ManufacturerCode: 999,
	DeviceClass:      25,
	DeviceFunction:   130,
	IndustryGroup:    4,
}

func TestRegistry_ClaimAddsDevice(t *testing.T) {
	c, mb, _ := newCitizenClient(t)

	mb.inbound <- claimFrame(testDeviceName, 0x42)

	ok := waitFor(t, 2*time.Second, func() bool { return len(c.Devices()) == 1 })
	require.True(t, ok, "claim should register a device")

	d := c.Devices()[0]
	assert.Equal(t, uint8(0x42), d.Address)
	assert.Equal(t, testDeviceName.Pack(true), d.RawName)
	assert.Equal(t, uint16(999), d.Name.ManufacturerCode)
	assert.Equal(t, uint32(424242), d.Name.IdentityNumber)
	assert.False(t, d.LastSeen.IsZero())
	assert.Nil(t, d.ProductInfo)
}

func TestRegistry_RequestsInfoOncePerName(t *testing.T) {
	_, mb, _ := newCitizenClient(t)

	mb.inbound <- claimFrame(testDeviceName, 0x42)
	ok := waitFor(t, 2*time.Second, func() bool {
		return infoRequestsTo(mb, 126996, 0x42) == 1 && infoRequestsTo(mb, 126998, 0x42) == 1
	})
	require.True(t, ok, "first sight of a NAME should trigger product and config info requests")

	// A repeat claim from the same NAME must not re-request.
	mb.inbound <- claimFrame(testDeviceName, 0x42)
	time.Sleep(150 * time.Millisecond)
	assert.Equal(t, 1, infoRequestsTo(mb, 126996, 0x42))
	assert.Equal(t, 1, infoRequestsTo(mb, 126998, 0x42))
}

func TestRegistry_AttachesProductInfo(t *testing.T) {
	c, mb, _ := newCitizenClient(t)

	mb.inbound <- claimFrame(testDeviceName, 0x42)
	waitFor(t, 2*time.Second, func() bool { return len(c.Devices()) == 1 })

	for _, f := range fastPacketFrames(t, &pgn.ProductInformation{ModelId: "sonar"}, 6, 0x42, 255) {
		mb.inbound <- f
	}

	ok := waitFor(t, 2*time.Second, func() bool {
		devs := c.Devices()
		return len(devs) == 1 && devs[0].ProductInfo != nil
	})
	require.True(t, ok, "product information should attach to the claiming device")
	assert.Equal(t, "sonar", c.Devices()[0].ProductInfo.ModelId)
}

func TestRegistry_AddressMove(t *testing.T) {
	c, mb, _ := newCitizenClient(t)

	mb.inbound <- claimFrame(testDeviceName, 0x42)
	waitFor(t, 2*time.Second, func() bool { return len(c.Devices()) == 1 })

	// The same NAME moves to a new address.
	mb.inbound <- claimFrame(testDeviceName, 0x43)
	ok := waitFor(t, 2*time.Second, func() bool {
		_, found := c.DeviceAt(0x43)
		return found
	})
	require.True(t, ok)

	_, stale := c.DeviceAt(0x42)
	assert.False(t, stale, "old address binding must be evicted")
	require.Len(t, c.Devices(), 1, "an address move must not duplicate the device")
	assert.Equal(t, uint8(0x43), c.Devices()[0].Address)
}

func TestRegistry_DevicesSortedByAddress(t *testing.T) {
	c, mb, _ := newCitizenClient(t)

	second := testDeviceName
	second.IdentityNumber = 111111

	mb.inbound <- claimFrame(testDeviceName, 0x50)
	mb.inbound <- claimFrame(second, 0x30)

	ok := waitFor(t, 2*time.Second, func() bool { return len(c.Devices()) == 2 })
	require.True(t, ok)

	devs := c.Devices()
	assert.Equal(t, uint8(0x30), devs[0].Address)
	assert.Equal(t, uint8(0x50), devs[1].Address)
}

func TestRegistry_StartupEnumeration(t *testing.T) {
	_, mb, _ := newCitizenClient(t)

	ok := waitFor(t, 2*time.Second, func() bool {
		return infoRequestsTo(mb, 60928, 255) >= 1
	})
	assert.True(t, ok, "the client should enumerate the bus with a broadcast claim request at startup")
}

func TestRegistry_LastSeenAdvances(t *testing.T) {
	c, mb, _ := newCitizenClient(t)

	mb.inbound <- claimFrame(testDeviceName, 0x42)
	waitFor(t, 2*time.Second, func() bool { return len(c.Devices()) == 1 })
	first := c.Devices()[0].LastSeen

	time.Sleep(20 * time.Millisecond)
	h := uint64(15708)
	mb.inbound <- singleFrame(t, &pgn.VesselHeading{Heading: &h}, 2, 0x42, 255)

	ok := waitFor(t, 2*time.Second, func() bool {
		return c.Devices()[0].LastSeen.After(first)
	})
	assert.True(t, ok, "ordinary traffic from a known address should advance LastSeen")
}

func TestRegistry_DeviceAtUnknown(t *testing.T) {
	c, _, _ := newCitizenClient(t)
	_, found := c.DeviceAt(0x77)
	assert.False(t, found)
}

func TestRegistry_SnapshotsAreCopies(t *testing.T) {
	c, mb, _ := newCitizenClient(t)

	mb.inbound <- claimFrame(testDeviceName, 0x42)
	waitFor(t, 2*time.Second, func() bool { return len(c.Devices()) == 1 })
	for _, f := range fastPacketFrames(t, &pgn.ProductInformation{ModelId: "sonar"}, 6, 0x42, 255) {
		mb.inbound <- f
	}
	waitFor(t, 2*time.Second, func() bool { return c.Devices()[0].ProductInfo != nil })

	c.Devices()[0].ProductInfo.ModelId = "tampered"
	assert.Equal(t, "sonar", c.Devices()[0].ProductInfo.ModelId, "snapshot mutations must not leak into the registry")
}
