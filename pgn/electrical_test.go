package pgn

import (
	"testing"

	"github.com/open-ships/n2k/units"
	"github.com/stretchr/testify/assert"
)

func TestRoundTrip_EncodeDecodeBatteryStatus(t *testing.T) {
	temp := units.NewTemperature(units.Kelvin, 298.15)
	original := &BatteryStatus{
		Instance:    ptrUint8(1),
		Voltage:     ptrFloat32(12.8),
		Current:     ptrFloat32(-5.5),
		Temperature: &temp,
		Sid:         ptrUint8(7),
	}

	data, err := EncodeBatteryStatus(original)
	assert.NoError(t, err)

	decoded, err := DecodeBatteryStatus(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(*BatteryStatus)
	assert.Equal(t, uint8(1), *result.Instance)
	assert.InDelta(t, float32(12.8), *result.Voltage, 0.01)
	assert.InDelta(t, float32(-5.5), *result.Current, 0.1)
	assert.NotNil(t, result.Temperature)
	assert.InDelta(t, float32(298.15), result.Temperature.Value, 0.01)
	assert.Equal(t, uint8(7), *result.Sid)
}

func TestRoundTrip_EncodeDecodeBinarySwitchBankStatus(t *testing.T) {
	original := &BinarySwitchBankStatus{
		Instance:    ptrUint8(1),
		Indicator1:  OffOnConst(1),
		Indicator2:  OffOnConst(0),
		Indicator3:  OffOnConst(1),
		Indicator4:  OffOnConst(0),
		Indicator5:  OffOnConst(0),
		Indicator6:  OffOnConst(0),
		Indicator7:  OffOnConst(0),
		Indicator8:  OffOnConst(0),
		Indicator9:  OffOnConst(0),
		Indicator10: OffOnConst(0),
		Indicator11: OffOnConst(0),
		Indicator12: OffOnConst(0),
		Indicator13: OffOnConst(0),
		Indicator14: OffOnConst(0),
		Indicator15: OffOnConst(0),
		Indicator16: OffOnConst(0),
		Indicator17: OffOnConst(0),
		Indicator18: OffOnConst(0),
		Indicator19: OffOnConst(0),
		Indicator20: OffOnConst(0),
		Indicator21: OffOnConst(0),
		Indicator22: OffOnConst(0),
		Indicator23: OffOnConst(0),
		Indicator24: OffOnConst(0),
		Indicator25: OffOnConst(0),
		Indicator26: OffOnConst(0),
		Indicator27: OffOnConst(0),
		Indicator28: OffOnConst(0),
	}

	data, err := EncodeBinarySwitchBankStatus(original)
	assert.NoError(t, err)

	decoded, err := DecodeBinarySwitchBankStatus(MessageInfo{}, NewPgnDataStream(data))
	assert.NoError(t, err)

	result := decoded.(*BinarySwitchBankStatus)
	assert.Equal(t, uint8(1), *result.Instance)
	assert.Equal(t, OffOnConst(1), result.Indicator1)
	assert.Equal(t, OffOnConst(0), result.Indicator2)
	assert.Equal(t, OffOnConst(1), result.Indicator3)
}
