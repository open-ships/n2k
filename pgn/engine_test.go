package pgn

import (
	"testing"

	"github.com/open-ships/n2k/units"
	"github.com/stretchr/testify/assert"
)

func TestRoundTrip_EncodeDecodeEngineParametersRapidUpdate(t *testing.T) {
	bp := units.NewPressure(units.Pa, 150000)
	original := &EngineParametersRapidUpdate{
		Instance:      EngineInstanceConst(0),
		Speed:         ptrFloat32(3000.0),
		BoostPressure: &bp,
		TiltTrim:      ptrInt8(5),
	}

	data, err := EncodeEngineParametersRapidUpdate(original)
	assert.NoError(t, err)

	decoded, err := DecodeEngineParametersRapidUpdate(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(*EngineParametersRapidUpdate)
	assert.Equal(t, EngineInstanceConst(0), result.Instance)
	assert.InDelta(t, float32(3000.0), *result.Speed, 0.25)
	assert.NotNil(t, result.BoostPressure)
	assert.InDelta(t, float32(150000), result.BoostPressure.Value, 100)
	assert.Equal(t, int8(5), *result.TiltTrim)
}

func TestRoundTrip_EncodeDecodeEngineParametersDynamic(t *testing.T) {
	oilPressure := units.NewPressure(units.Pa, 250000)
	oilTemp := units.NewTemperature(units.Kelvin, 370.0)
	temp := units.NewTemperature(units.Kelvin, 360.0)
	coolantPressure := units.NewPressure(units.Pa, 100000)
	fuelPressure := units.NewPressure(units.Pa, 300000)
	fuelRate := units.NewFlow(units.LitersPerHour, 15.5)

	original := &EngineParametersDynamic{
		Instance:            EngineInstanceConst(0),
		OilPressure:         &oilPressure,
		OilTemperature:      &oilTemp,
		Temperature:         &temp,
		AlternatorPotential: ptrFloat32(14.2),
		FuelRate:            &fuelRate,
		TotalEngineHours:    ptrUint32(36000),
		CoolantPressure:     &coolantPressure,
		FuelPressure:        &fuelPressure,
		DiscreteStatus1:     EngineStatus1Const(0),
		DiscreteStatus2:     EngineStatus2Const(0),
		EngineLoad:          ptrInt8(75),
		EngineTorque:        ptrInt8(80),
	}

	data, err := EncodeEngineParametersDynamic(original)
	assert.NoError(t, err)

	decoded, err := DecodeEngineParametersDynamic(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(*EngineParametersDynamic)
	assert.Equal(t, EngineInstanceConst(0), result.Instance)
	assert.NotNil(t, result.OilPressure)
	assert.InDelta(t, float32(250000), result.OilPressure.Value, 100)
	assert.NotNil(t, result.OilTemperature)
	assert.InDelta(t, float32(370.0), result.OilTemperature.Value, 0.1)
	assert.NotNil(t, result.Temperature)
	assert.InDelta(t, float32(360.0), result.Temperature.Value, 0.01)
	assert.InDelta(t, float32(14.2), *result.AlternatorPotential, 0.01)
	assert.NotNil(t, result.FuelRate)
	assert.InDelta(t, float32(15.5), result.FuelRate.Value, 0.1)
	assert.Equal(t, uint32(36000), *result.TotalEngineHours)
	assert.Equal(t, int8(75), *result.EngineLoad)
	assert.Equal(t, int8(80), *result.EngineTorque)
}
