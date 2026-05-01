package pgn

import (
	"testing"

	"github.com/open-ships/n2k/units"
	"github.com/stretchr/testify/assert"
)

func TestRoundTrip_EncodeDecodeWindData(t *testing.T) {
	ws := units.NewVelocity(units.MetersPerSecond, 12.5)
	original := &WindData{
		Sid:       ptrUint8(1),
		WindSpeed: &ws,
		WindAngle: ptrFloat32(1.5708),
		Reference: WindReferenceConst(0),
	}

	data, err := EncodeWindData(original)
	assert.NoError(t, err)

	decoded, err := DecodeWindData(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(*WindData)
	assert.Equal(t, uint8(1), *result.Sid)
	assert.NotNil(t, result.WindSpeed)
	assert.InDelta(t, float32(12.5), result.WindSpeed.Value, 0.01)
	assert.InDelta(t, float32(1.5708), *result.WindAngle, 0.0001)
	assert.Equal(t, WindReferenceConst(0), result.Reference)
}

func TestRoundTrip_EncodeDecodeTemperature(t *testing.T) {
	actualTemp := units.NewTemperature(units.Kelvin, 295.15)
	setTemp := units.NewTemperature(units.Kelvin, 293.15)
	original := &Temperature{
		Sid:               ptrUint8(5),
		Instance:          ptrUint8(0),
		Source:            TemperatureSourceConst(1),
		ActualTemperature: &actualTemp,
		SetTemperature:    &setTemp,
	}

	data, err := EncodeTemperature(original)
	assert.NoError(t, err)

	decoded, err := DecodeTemperature(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(*Temperature)
	assert.Equal(t, uint8(5), *result.Sid)
	assert.Equal(t, uint8(0), *result.Instance)
	assert.Equal(t, TemperatureSourceConst(1), result.Source)
	assert.NotNil(t, result.ActualTemperature)
	assert.InDelta(t, float32(295.15), result.ActualTemperature.Value, 0.01)
	assert.NotNil(t, result.SetTemperature)
	assert.InDelta(t, float32(293.15), result.SetTemperature.Value, 0.01)
}

func TestRoundTrip_EncodeDecodeHumidity(t *testing.T) {
	original := &Humidity{
		Sid:            ptrUint8(2),
		Instance:       ptrUint8(0),
		Source:         HumiditySourceConst(0),
		ActualHumidity: ptrFloat32(65.0),
		SetHumidity:    ptrFloat32(50.0),
	}

	data, err := EncodeHumidity(original)
	assert.NoError(t, err)

	decoded, err := DecodeHumidity(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(*Humidity)
	assert.Equal(t, uint8(2), *result.Sid)
	assert.Equal(t, uint8(0), *result.Instance)
	assert.Equal(t, HumiditySourceConst(0), result.Source)
	assert.InDelta(t, float32(65.0), *result.ActualHumidity, 0.004)
	assert.InDelta(t, float32(50.0), *result.SetHumidity, 0.004)
}

func TestRoundTrip_EncodeDecodeEnvironmentalParameters(t *testing.T) {
	temp := units.NewTemperature(units.Kelvin, 298.15)
	pressure := units.NewPressure(units.Pa, 101300)
	original := &EnvironmentalParameters{
		Sid:                 ptrUint8(3),
		TemperatureSource:   TemperatureSourceConst(1),
		HumiditySource:      HumiditySourceConst(0),
		Temperature:         &temp,
		Humidity:            ptrFloat32(55.0),
		AtmosphericPressure: &pressure,
	}

	data, err := EncodeEnvironmentalParameters(original)
	assert.NoError(t, err)

	decoded, err := DecodeEnvironmentalParameters(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(*EnvironmentalParameters)
	assert.Equal(t, uint8(3), *result.Sid)
	assert.Equal(t, TemperatureSourceConst(1), result.TemperatureSource)
	assert.NotNil(t, result.Temperature)
	assert.InDelta(t, float32(298.15), result.Temperature.Value, 0.01)
	assert.InDelta(t, float32(55.0), *result.Humidity, 0.004)
	assert.NotNil(t, result.AtmosphericPressure)
	assert.InDelta(t, float32(101300), result.AtmosphericPressure.Value, 100)
}
