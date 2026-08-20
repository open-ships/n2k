package actisense

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func responseFor(command byte, sequence uint8, model ModelID, data []byte) BEMResponse {
	return BEMResponse{
		BSTID: BSTBEMResponse, BEMID: command, Origin: LocalBEMOrigin,
		Sequence: sequence, ModelID: uint16(model), SerialNumber: 0x55667788,
		Data: append([]byte(nil), data...),
	}
}

func padded32(text string) []byte {
	result := bytes.Repeat([]byte{0xFF}, 32)
	copy(result, text)
	return result
}

func TestDecodeProductInfoFormat2IndependentVector(t *testing.T) {
	data := make([]byte, 138)
	binary.LittleEndian.PutUint32(data[0:4], 0x00000011)
	binary.LittleEndian.PutUint16(data[4:6], 2101)
	binary.LittleEndian.PutUint16(data[6:8], 1234)
	copy(data[8:40], padded32("NGX-1"))
	copy(data[40:72], padded32("2.004"))
	copy(data[72:104], padded32("A"))
	copy(data[104:136], padded32("SN123"))
	data[136], data[137] = 2, 3

	info, err := DecodeProductInfo(responseFor(BEMProductInfo, 6, ModelNGX1, data))
	require.NoError(t, err)
	assert.Equal(t, uint16(2101), info.NMEA2000Version)
	assert.Equal(t, uint16(1234), info.ProductCode)
	assert.Equal(t, "NGX-1", info.Model)
	assert.Equal(t, "2.004", info.SoftwareVersion)
	assert.Equal(t, "A", info.ModelVersion)
	assert.Equal(t, "SN123", info.SerialCode)
	assert.Equal(t, ModelNGX1, info.DeviceModelID)
	assert.False(t, info.Legacy)
}

func TestProductInfoLegacyFiveMessageAssembly(t *testing.T) {
	main := []byte{0x34, 0x08, 0x65, 0x00, 0x01, 0x02}
	accumulator := &productInfoAccumulator{}
	parts := []BEMResponse{
		responseFor(BEMProductInfo, 1, ModelNGT1, main),
		responseFor(BEMProductInfo, 2, ModelNGT1, padded32("NGT-1")),
		responseFor(BEMProductInfo, 3, ModelNGT1, padded32("1.211")),
		responseFor(BEMProductInfo, 4, ModelNGT1, padded32("A")),
		responseFor(BEMProductInfo, 5, ModelNGT1, padded32("123456")),
	}
	for index, part := range parts {
		done, err := accumulator.feed(part)
		require.NoError(t, err)
		assert.Equal(t, index == len(parts)-1, done)
	}
	assert.Equal(t, uint16(0x0834), accumulator.result.NMEA2000Version)
	assert.Equal(t, uint16(0x0065), accumulator.result.ProductCode)
	assert.Equal(t, "NGT-1", accumulator.result.Model)
	assert.Equal(t, "123456", accumulator.result.SerialCode)
	assert.True(t, accumulator.result.Legacy)
}

func TestPortInventoryAggregatesAllSlots(t *testing.T) {
	makePart := func(first uint8, names ...string) BEMResponse {
		data := make([]byte, 8+22*len(names))
		data[0] = 9
		binary.LittleEndian.PutUint32(data[1:5], 0x00001104)
		data[5], data[6], data[7] = 3, first, uint8(len(names))
		for index, name := range names {
			record := data[8+22*index : 8+22*(index+1)]
			record[0] = first + uint8(index)
			record[1], record[2] = first+uint8(index), first+uint8(index)
			record[3], record[4], record[5] = byte(PortMediaCAN), byte(HardwareCANNMEA2000), PortCanReceive|PortCanTransmit
			binary.LittleEndian.PutUint32(record[6:10], 250000)
			binary.LittleEndian.PutUint32(record[10:14], 250000)
			copy(record[14:22], name)
		}
		return responseFor(BEMPortInventory, 1, ModelNGX1, data)
	}
	accumulator := &portInventoryAccumulator{}
	done, err := accumulator.feed(makePart(0, "CAN0"))
	require.NoError(t, err)
	assert.False(t, done)
	done, err = accumulator.feed(makePart(1, "CAN1", "CAN2"))
	require.NoError(t, err)
	assert.True(t, done)
	require.Len(t, accumulator.result.Ports, 3)
	assert.Equal(t, "CAN2", accumulator.result.Ports[2].Name)
	assert.True(t, accumulator.result.Ports[2].CanTransmit())
}

func TestDeviceCommandPayloadGoldenVectors(t *testing.T) {
	assert.Equal(t, []byte{7}, PortBaudrateGet(7))
	assert.Equal(t, []byte{7, 0x00, 0xC2, 0x01, 0x00, 0xFF, 0xFF, 0xFF, 0xFF}, PortBaudrateSet(7, 115200, BaudRateNoChange))
	assert.Equal(t, []byte{0x04, 0x01, 0x48, 0x69}, mustEncodeCANInfo(t, "Hi"))
	assert.Equal(t, []byte{3, 1, 2, 3}, mustEncodeEcho(t, []byte{1, 2, 3}))
	assert.Equal(t, []byte{0x78, 0x56, 0x34, 0x12, 0xEF, 0xCD, 0xAB, 0x90}, TotalTimeSet(0x12345678, 0x90ABCDEF))
}

func mustEncodeCANInfo(t *testing.T, text string) []byte {
	t.Helper()
	data, err := EncodeCANInfoField(text)
	require.NoError(t, err)
	return data
}

func mustEncodeEcho(t *testing.T, payload []byte) []byte {
	t.Helper()
	data, err := EncodeEcho(payload)
	require.NoError(t, err)
	return data
}

func TestCapabilitiesAreExplicitAndConservative(t *testing.T) {
	ngx := CapabilitiesForModel(ModelNGX1)
	assert.True(t, ngx.ProprietaryPGNEnableListF2)
	assert.True(t, ngx.ReceiveAllOmitsISOControlPGNs)
	assert.False(t, ngx.RewritesHostTransmitSID)
	ngt := CapabilitiesForModel(ModelNGT1)
	assert.True(t, ngt.RewritesHostTransmitSID)
	unknown := CapabilitiesForModel(0xCAFE)
	assert.False(t, unknown.ProprietaryPGNEnableListF2)
}
