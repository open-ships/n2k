package pgn

import (
	"testing"

	"github.com/open-ships/n2k/units"
	"github.com/stretchr/testify/assert"
)

func TestRoundTrip_EncodeDecodeVesselHeading(t *testing.T) {
	original := &VesselHeading{
		Sid:       ptrUint8(5),
		Heading:   ptrFloat32(1.5708),
		Deviation: ptrFloat32(-0.0523),
		Variation: ptrFloat32(0.0349),
		Reference: DirectionReferenceConst(1),
	}

	data, err := EncodeVesselHeading(original)
	assert.NoError(t, err)

	decoded, err := DecodeVesselHeading(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(VesselHeading)
	assert.Equal(t, uint8(5), *result.Sid)
	assert.InDelta(t, float32(1.5708), *result.Heading, 0.0002)
	assert.InDelta(t, float32(-0.0523), *result.Deviation, 0.0002)
	assert.InDelta(t, float32(0.0349), *result.Variation, 0.0002)
	assert.Equal(t, DirectionReferenceConst(1), result.Reference)
}

func TestRoundTrip_EncodeDecodeCogSogRapidUpdate(t *testing.T) {
	sog := units.NewVelocity(units.MetersPerSecond, 5.25)
	original := &CogSogRapidUpdate{
		Sid:          ptrUint8(3),
		CogReference: DirectionReferenceConst(0),
		Cog:          ptrFloat32(3.1416),
		Sog:          &sog,
	}

	data, err := EncodeCogSogRapidUpdate(original)
	assert.NoError(t, err)

	decoded, err := DecodeCogSogRapidUpdate(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(CogSogRapidUpdate)
	assert.Equal(t, uint8(3), *result.Sid)
	assert.Equal(t, DirectionReferenceConst(0), result.CogReference)
	assert.InDelta(t, float32(3.1416), *result.Cog, 0.0002)
	assert.NotNil(t, result.Sog)
	assert.InDelta(t, float32(5.25), result.Sog.Value, 0.02)
}

func TestRoundTrip_EncodeDecodePositionRapidUpdate(t *testing.T) {
	original := &PositionRapidUpdate{
		Latitude:  ptrFloat64(37.7749),
		Longitude: ptrFloat64(-122.4194),
	}

	data, err := EncodePositionRapidUpdate(original)
	assert.NoError(t, err)

	decoded, err := DecodePositionRapidUpdate(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(PositionRapidUpdate)
	assert.InDelta(t, float64(37.7749), *result.Latitude, 1e-06)
	assert.InDelta(t, float64(-122.4194), *result.Longitude, 1e-06)
}

func TestRoundTrip_EncodeDecodeGnssPositionData(t *testing.T) {
	alt := units.NewDistance(units.Meter, 100.5)
	geoSep := units.NewDistance(units.Meter, -32.5)
	refStations := uint8(1)
	original := &GnssPositionData{
		Sid:               ptrUint8(7),
		Date:              ptrUint16(19724),
		Time:              ptrFloat32(43200.5),
		Latitude:          ptrFloat64(37.7749),
		Longitude:         ptrFloat64(-122.4194),
		Altitude:          &alt,
		GnssType:          GnsConst(2),
		Method:            GnsMethodConst(1),
		Integrity:         GnsIntegrityConst(0),
		NumberOfSvs:       ptrUint8(12),
		Hdop:              ptrFloat32(1.2),
		Pdop:              ptrFloat32(2.1),
		GeoidalSeparation: &geoSep,
		ReferenceStations: &refStations,
		Repeating1: []GnssPositionDataRepeating1{
			{
				ReferenceStationType:   GnsConst(3),
				ReferenceStationId:     ptrUint16(42),
				AgeOfDgnssCorrections:  ptrFloat32(1.5),
			},
		},
	}

	data, err := EncodeGnssPositionData(original)
	assert.NoError(t, err)

	decoded, err := DecodeGnssPositionData(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(GnssPositionData)
	assert.Equal(t, uint8(7), *result.Sid)
	assert.Equal(t, uint16(19724), *result.Date)
	assert.InDelta(t, float32(43200.5), *result.Time, 0.001)
	assert.InDelta(t, float64(37.7749), *result.Latitude, 1e-06)
	assert.InDelta(t, float64(-122.4194), *result.Longitude, 1e-06)
	assert.NotNil(t, result.Altitude)
	assert.InDelta(t, float32(100.5), result.Altitude.Value, 0.01)
	assert.Equal(t, GnsConst(2), result.GnssType)
	assert.Equal(t, GnsMethodConst(1), result.Method)
	assert.Equal(t, GnsIntegrityConst(0), result.Integrity)
	assert.Equal(t, uint8(12), *result.NumberOfSvs)
	assert.InDelta(t, float32(1.2), *result.Hdop, 0.02)
	assert.InDelta(t, float32(2.1), *result.Pdop, 0.02)
	assert.NotNil(t, result.GeoidalSeparation)
	assert.InDelta(t, float32(-32.5), result.GeoidalSeparation.Value, 0.02)
	assert.Equal(t, uint8(1), *result.ReferenceStations)
	assert.Len(t, result.Repeating1, 1)
	assert.Equal(t, GnsConst(3), result.Repeating1[0].ReferenceStationType)
	assert.Equal(t, uint16(42), *result.Repeating1[0].ReferenceStationId)
	assert.InDelta(t, float32(1.5), *result.Repeating1[0].AgeOfDgnssCorrections, 0.02)
}

func TestRoundTrip_EncodeDecodeAttitude(t *testing.T) {
	original := &Attitude{
		Sid:   ptrUint8(1),
		Yaw:   ptrFloat32(0.5236),
		Pitch: ptrFloat32(-0.1745),
		Roll:  ptrFloat32(0.0873),
	}

	data, err := EncodeAttitude(original)
	assert.NoError(t, err)

	decoded, err := DecodeAttitude(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(Attitude)
	assert.Equal(t, uint8(1), *result.Sid)
	assert.InDelta(t, float32(0.5236), *result.Yaw, 0.0002)
	assert.InDelta(t, float32(-0.1745), *result.Pitch, 0.0002)
	assert.InDelta(t, float32(0.0873), *result.Roll, 0.0002)
}

func TestRoundTrip_EncodeDecodeSmallCraftStatus(t *testing.T) {
	original := &SmallCraftStatus{
		PortTrimTab:      ptrInt8(-50),
		StarboardTrimTab: ptrInt8(75),
	}

	data, err := EncodeSmallCraftStatus(original)
	assert.NoError(t, err)

	decoded, err := DecodeSmallCraftStatus(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(SmallCraftStatus)
	assert.Equal(t, int8(-50), *result.PortTrimTab)
	assert.Equal(t, int8(75), *result.StarboardTrimTab)
}
