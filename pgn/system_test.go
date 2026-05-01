package pgn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundTrip_EncodeDecodeIsoAddressClaim(t *testing.T) {
	original := &IsoAddressClaim{
		UniqueNumber:            ptrUint32(123456),
		ManufacturerCode:        ManufacturerCodeConst(42),
		DeviceInstanceLower:     ptrUint8(3),
		DeviceInstanceUpper:     ptrUint8(5),
		DeviceFunction:          DeviceFunctionConst(130),
		DeviceClass:             DeviceClassConst(25),
		SystemInstance:          ptrUint8(0),
		IndustryGroup:           IndustryCodeConst(4),
		ArbitraryAddressCapable: ptrUint8(0),
	}

	data, err := EncodeIsoAddressClaim(original)
	assert.NoError(t, err)

	decoded, err := DecodeIsoAddressClaim(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(IsoAddressClaim)
	assert.Equal(t, *original.UniqueNumber, *result.UniqueNumber)
	assert.Equal(t, original.ManufacturerCode, result.ManufacturerCode)
	assert.Equal(t, *original.DeviceInstanceLower, *result.DeviceInstanceLower)
	assert.Equal(t, *original.DeviceInstanceUpper, *result.DeviceInstanceUpper)
	assert.Equal(t, original.DeviceFunction, result.DeviceFunction)
	assert.Equal(t, original.DeviceClass, result.DeviceClass)
	assert.Equal(t, *original.SystemInstance, *result.SystemInstance)
	assert.Equal(t, original.IndustryGroup, result.IndustryGroup)
	assert.Equal(t, uint8(0), *result.ArbitraryAddressCapable)
}

func TestRoundTrip_EncodeDecodeHeartbeat(t *testing.T) {
	original := &Heartbeat{
		DataTransmitOffset: ptrFloat32(1.5),
		SequenceCounter:    ptrUint8(42),
		Controller1State:   ControllerStateConst(1),
		Controller2State:   ControllerStateConst(0),
		EquipmentStatus:    EquipmentStatusConst(1),
	}

	data, err := EncodeHeartbeat(original)
	assert.NoError(t, err)

	decoded, err := DecodeHeartbeat(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(Heartbeat)
	assert.InDelta(t, float32(1.5), *result.DataTransmitOffset, 0.001)
	assert.Equal(t, uint8(42), *result.SequenceCounter)
	assert.Equal(t, ControllerStateConst(1), result.Controller1State)
	assert.Equal(t, ControllerStateConst(0), result.Controller2State)
	assert.Equal(t, EquipmentStatusConst(1), result.EquipmentStatus)
}

func TestRoundTrip_EncodeDecodeSystemTime(t *testing.T) {
	original := &SystemTime{
		Sid:    ptrUint8(7),
		Source: SystemTimeConst(1),
		Date:   ptrUint16(19724),
		Time:   ptrFloat32(43200.5),
	}

	data, err := EncodeSystemTime(original)
	assert.NoError(t, err)

	decoded, err := DecodeSystemTime(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(SystemTime)
	assert.Equal(t, uint8(7), *result.Sid)
	assert.Equal(t, SystemTimeConst(1), result.Source)
	assert.Equal(t, uint16(19724), *result.Date)
	assert.InDelta(t, float32(43200.5), *result.Time, 0.0001)
}

func TestRoundTrip_EncodeDecodeProductInformation(t *testing.T) {
	original := &ProductInformation{
		Nmea2000Version:     ptrFloat32(2.0),
		ProductCode:         ptrUint16(1234),
		ModelId:             "TestModel",
		SoftwareVersionCode: "1.0.0",
		ModelVersion:        "Rev A",
		ModelSerialCode:     "SN00001",
		CertificationLevel:  ptrUint8(1),
		LoadEquivalency:     ptrUint8(1),
	}

	data, err := EncodeProductInformation(original)
	assert.NoError(t, err)

	decoded, err := DecodeProductInformation(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(ProductInformation)
	assert.InDelta(t, float32(2.0), *result.Nmea2000Version, 0.001)
	assert.Equal(t, uint16(1234), *result.ProductCode)
	assert.Equal(t, "TestModel", result.ModelId)
	assert.Equal(t, "1.0.0", result.SoftwareVersionCode)
	assert.Equal(t, "Rev A", result.ModelVersion)
	assert.Equal(t, "SN00001", result.ModelSerialCode)
	assert.Equal(t, uint8(1), *result.CertificationLevel)
	assert.Equal(t, uint8(1), *result.LoadEquivalency)
}

func TestRoundTrip_EncodeDecodeIsoRequest(t *testing.T) {
	original := &IsoRequest{
		Pgn: ptrUint32(60928),
	}

	data, err := EncodeIsoRequest(original)
	assert.NoError(t, err)

	decoded, err := DecodeIsoRequest(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(IsoRequest)
	assert.Equal(t, uint32(60928), *result.Pgn)
}
