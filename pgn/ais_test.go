package pgn

import (
	"testing"

	"github.com/open-ships/n2k/units"
	"github.com/stretchr/testify/assert"
)

func TestRoundTrip_EncodeDecodeAisClassAPositionReport(t *testing.T) {
	sog := units.NewVelocity(units.MetersPerSecond, 5.14)
	original := &AisClassAPositionReport{
		MessageId:                 AisMessageIdConst(1),
		RepeatIndicator:           RepeatIndicatorConst(0),
		UserId:                    ptrUint32(211234567),
		Longitude:                 ptrFloat64(9.9936),
		Latitude:                  ptrFloat64(53.5511),
		PositionAccuracy:          PositionAccuracyConst(1),
		Raim:                      RaimFlagConst(0),
		TimeStamp:                 TimeStampConst(30),
		Cog:                       ptrFloat32(2.3562),
		Sog:                       &sog,
		CommunicationState:        []uint8{0x00, 0x00, 0x00},
		AisTransceiverInformation: AisTransceiverConst(0),
		Heading:                   ptrFloat32(2.2689),
		RateOfTurn:                ptrFloat32(0.0),
		NavStatus:                 NavStatusConst(0),
		SpecialManeuverIndicator:  AisSpecialManeuverConst(0),
		SequenceId:                ptrUint8(0),
	}

	data, err := EncodeAisClassAPositionReport(original)
	assert.NoError(t, err)

	decoded, err := DecodeAisClassAPositionReport(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(AisClassAPositionReport)
	assert.Equal(t, AisMessageIdConst(1), result.MessageId)
	assert.Equal(t, RepeatIndicatorConst(0), result.RepeatIndicator)
	assert.Equal(t, uint32(211234567), *result.UserId)
	assert.InDelta(t, 9.9936, *result.Longitude, 0.0000001)
	assert.InDelta(t, 53.5511, *result.Latitude, 0.0000001)
	assert.Equal(t, PositionAccuracyConst(1), result.PositionAccuracy)
	assert.Equal(t, RaimFlagConst(0), result.Raim)
	assert.Equal(t, TimeStampConst(30), result.TimeStamp)
	assert.InDelta(t, float32(2.3562), *result.Cog, 0.0001)
	assert.NotNil(t, result.Sog)
	assert.InDelta(t, float32(5.14), result.Sog.Value, 0.01)
	assert.Equal(t, AisTransceiverConst(0), result.AisTransceiverInformation)
	assert.InDelta(t, float32(2.2689), *result.Heading, 0.0001)
	assert.InDelta(t, float32(0.0), *result.RateOfTurn, 0.0001)
	assert.Equal(t, NavStatusConst(0), result.NavStatus)
}

func TestRoundTrip_EncodeDecodeAisClassBPositionReport(t *testing.T) {
	sog := units.NewVelocity(units.MetersPerSecond, 3.0)
	original := &AisClassBPositionReport{
		MessageId:                 AisMessageIdConst(18),
		RepeatIndicator:           RepeatIndicatorConst(0),
		UserId:                    ptrUint32(338123456),
		Longitude:                 ptrFloat64(-122.4194),
		Latitude:                  ptrFloat64(37.7749),
		PositionAccuracy:          PositionAccuracyConst(0),
		Raim:                      RaimFlagConst(0),
		TimeStamp:                 TimeStampConst(15),
		Cog:                       ptrFloat32(3.1416),
		Sog:                       &sog,
		AisTransceiverInformation: AisTransceiverConst(0),
		Heading:                   ptrFloat32(3.0543),
		UnitType:                  AisTypeConst(0),
		IntegratedDisplay:         YesNoConst(0),
		Dsc:                       YesNoConst(0),
		Band:                      AisBandConst(0),
		CanHandleMsg22:            YesNoConst(0),
		AisMode:                   AisModeConst(0),
		AisCommunicationState:     AisCommunicationStateConst(0),
	}

	data, err := EncodeAisClassBPositionReport(original)
	assert.NoError(t, err)

	decoded, err := DecodeAisClassBPositionReport(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(AisClassBPositionReport)
	assert.Equal(t, AisMessageIdConst(18), result.MessageId)
	assert.Equal(t, uint32(338123456), *result.UserId)
	assert.InDelta(t, -122.4194, *result.Longitude, 0.0000001)
	assert.InDelta(t, 37.7749, *result.Latitude, 0.0000001)
	assert.InDelta(t, float32(3.1416), *result.Cog, 0.0001)
	assert.NotNil(t, result.Sog)
	assert.InDelta(t, float32(3.0), result.Sog.Value, 0.01)
}
