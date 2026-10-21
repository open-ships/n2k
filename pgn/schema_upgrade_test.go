package pgn

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchemaAddressFieldsPreserveAllWireValues(t *testing.T) {
	// ADDRESS has no numeric null sentinel: 254 and 255 retain their network
	// meanings. This does not change the client's claimable-address policy.
	for _, definition := range sourceDefinitions {
		for _, field := range definition.Fields {
			if field.FieldType != "ADDRESS" {
				continue
			}
			t.Run(definition.StructName+"/"+field.SourceID, func(t *testing.T) {
				for _, address := range []uint64{0, 251, 252, 253, 254, 255} {
					message := structTypeRegistry[definition.StructName]()
					value := reflect.ValueOf(message).Elem()
					indexes, _, err := codecStructIndexes(value.Type())
					require.NoError(t, err)
					value.Field(indexes[field.Order]).Set(reflect.ValueOf(&address))
					wire, err := message.EncodePayload()
					require.NoError(t, err)
					require.Equal(t, byte(address), wire[*field.BitOffset/8])
					decoded := structTypeRegistry[definition.StructName]()
					require.NoError(t, decoded.DecodePayload(wire))
					got := reflect.ValueOf(decoded).Elem().Field(indexes[field.Order])
					require.False(t, got.IsNil(), "address %d became nil", address)
					require.Equal(t, address, got.Elem().Uint())
				}
			})
		}
	}
}

func TestSchemaUpgradeWireLayouts(t *testing.T) {
	// Literal payloads exercise corrected/new layouts independently of the
	// generator. Clearing MessageInfo also checks encoding without raw replay.
	tests := []struct {
		name    string
		number  uint32
		payload []byte
		check   func(*testing.T, PGN)
	}{
		{"Airmar master reset", 126720, []byte{0x87, 0x98, 1, 0xff, 0xff, 0xff}, func(t *testing.T, message PGN) {
			_, ok := message.(*AirmarMasterReset)
			require.True(t, ok, "%T", message)
		}},
		{"Airmar EEPROM reset", 126720, []byte{0x87, 0x98, 0x82, 0xf2, 0xff, 0xff}, func(t *testing.T, message PGN) {
			m, ok := message.(*AirmarResetEeprom)
			require.True(t, ok, "%T", message)
			require.Equal(t, uint64(2), *m.Options)
		}},
		{"Airmar POST results", 130944, []byte{0x87, 0x98, 0xfe, 2, 1, 0, 4, 2}, func(t *testing.T, message PGN) {
			m, ok := message.(*AirmarPost)
			require.True(t, ok, "%T", message)
			require.Len(t, m.Repeating1, 2)
			require.Equal(t, uint64(1), *m.Repeating1[0].TestId)
			require.Equal(t, uint64(0), *m.Repeating1[0].TestResult)
			require.Equal(t, uint64(4), *m.Repeating1[1].TestId)
			require.Equal(t, uint64(2), *m.Repeating1[1].TestResult)
		}},
		{"Airmar signed voltage minimum", 65410, []byte{0x87, 0x98, 1, 0x30, 0x75, 0, 0x80, 0xff}, func(t *testing.T, message PGN) {
			m, ok := message.(*AirmarDeviceInformation)
			require.True(t, ok, "%T", message)
			require.Equal(t, int64(-32768), *m.SupplyVoltage)
			voltage, valid := m.SupplyVoltageValue()
			require.True(t, valid)
			require.Equal(t, -327.68, voltage)
			require.NoError(t, m.SetSupplyVoltageValue(-327.68))
		}},
		{"Utility phase B reactive power", 65009, []byte{0, 0x5e, 0xd0, 0xb2, 0, 0x60, 0xfd, 0xff}, func(t *testing.T, message PGN) {
			m, ok := message.(*UtilityPhaseBAcReactivePower)
			require.True(t, ok, "%T", message)
			require.Equal(t, int64(3_000_000_000), *m.ReactivePower)
			power, valid := m.ReactivePowerValue()
			require.True(t, valid)
			require.Equal(t, float64(1_000_000_000), power)
			factor, valid := m.PowerFactorValue()
			require.True(t, valid)
			require.Equal(t, 0.5, factor)
			require.Equal(t, uint64(1), *m.PowerFactorLagging)
			require.NoError(t, m.SetPowerFactorValue(0.5))
		}},
		{"Utility phase C reactive power", 65006, []byte{0, 0x5e, 0xd0, 0xb2, 0, 0x60, 0xfe, 0xff}, func(t *testing.T, message PGN) {
			m, ok := message.(*UtilityPhaseCAcReactivePower)
			require.True(t, ok, "%T", message)
			require.Equal(t, int64(3_000_000_000), *m.ReactivePower)
			factor, valid := m.PowerFactorValue()
			require.True(t, valid)
			require.Equal(t, 0.5, factor)
			require.Equal(t, uint64(2), *m.PowerFactorLagging)
		}},
		{"Fusion menu lock", 130820, []byte{0xa3, 0x99, 0x12, 0x80, 0x78, 0x56, 0x34, 0x12, 7}, func(t *testing.T, message PGN) {
			m, ok := message.(*FusionMenuLockId)
			require.True(t, ok, "%T", message)
			require.Equal(t, uint64(0x12345678), *m.RequestToken)
			require.Equal(t, uint64(7), *m.LockId)
		}},
		{"Fusion system capabilities", 130820, []byte{0xa3, 0x99, 0x3d, 0x80, 1, 0, 0, 0, 0, 0, 0, 0x80}, func(t *testing.T, message PGN) {
			m, ok := message.(*FusionSystemCapabilities)
			require.True(t, ok, "%T", message)
			require.Equal(t, uint64(1)<<63|1, *m.Capabilities)
		}},
		{"Simnet ZC key", 65332, []byte{0x41, 0x9f, 0xff, 0x84, 0, 0, 1, 2}, func(t *testing.T, message PGN) {
			m, ok := message.(*SimnetZcKey)
			require.True(t, ok, "%T", message)
			require.Equal(t, uint64(255), *m.Address)
			require.Equal(t, uint64(1), *m.KeyEvent)
			require.Equal(t, uint64(2), *m.Key)
		}},
		{"Simnet ZC knob", 65332, []byte{0x41, 0x9f, 0xfe, 0x85, 0, 0, 0xfe, 0}, func(t *testing.T, message PGN) {
			m, ok := message.(*SimnetZcKnob)
			require.True(t, ok, "%T", message)
			require.Equal(t, uint64(254), *m.Address)
			require.Equal(t, int64(-2), *m.Ticks)
		}},
		{"Yanmar engine data", 65280, []byte{0xac, 0x98, 0x20, 0xf4, 0x41, 0xb8, 0x0b, 0x5a}, func(t *testing.T, message PGN) {
			m, ok := message.(*YanmarEngineDataA)
			require.True(t, ok, "%T", message)
			require.Equal(t, uint64(1), *m.EngineInstance)
			require.Equal(t, uint64(1), *m.TransmissionGear)
			throttle, valid := m.ThrottlePositionValue()
			require.True(t, valid)
			require.Equal(t, float64(50), throttle)
			require.Equal(t, uint64(3000), *m.EngineSpeed)
			require.Equal(t, []byte{0x5a}, m.UnknownData)
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			message, err := DecodePayload(MessageInfo{PGN: test.number}, test.payload)
			require.NoError(t, err)
			test.check(t, message)
			message.SetMessageInfo(MessageInfo{})
			wire, err := message.EncodePayload()
			require.NoError(t, err)
			require.Equal(t, test.payload, wire)
		})
	}
}

func TestOffsetEncodedPowerRange(t *testing.T) {
	// All utility and generator power fields use 32 unsigned wire bits with
	// a -2e9 offset, despite signed Go storage. These boundaries must agree
	// across decode, encode, typed accessors, and runtime PhysicalValue.
	checked := 0
	for _, definition := range sourceDefinitions {
		for _, field := range definition.Fields {
			if !field.Signed || field.Offset == nil || *field.Offset != -2_000_000_000 {
				continue
			}
			checked++
			t.Run(definition.StructName+"/"+field.SourceID, func(t *testing.T) {
				for _, raw := range []uint32{0, 2_000_000_000, 0x7fffffff, 0x80000000, 3_000_000_000, 0xfffffffc, 0xfffffffd, 0xfffffffe, 0xffffffff} {
					payload := bytes.Repeat([]byte{0xff}, 8)
					binary.LittleEndian.PutUint32(payload[*field.BitOffset/8:], raw)
					message := structTypeRegistry[definition.StructName]()
					require.NoError(t, message.DecodePayload(payload))
					target := reflect.ValueOf(message).Elem()
					indexes, _, err := codecStructIndexes(target.Type())
					require.NoError(t, err)
					index := indexes[field.Order]
					value := target.Field(index)
					if raw == 0xffffffff {
						require.True(t, value.IsNil())
					} else {
						require.False(t, value.IsNil(), "raw %x became nil", raw)
						require.Equal(t, int64(raw), value.Elem().Int())
					}
					available := raw <= 0xfffffffc
					want := float64(raw) - 2_000_000_000
					physical, _, valid, err := PhysicalValue(message, field.Order)
					require.NoError(t, err)
					require.Equal(t, available, valid, "raw %x", raw)
					name := target.Type().Field(index).Name
					getter := reflect.ValueOf(message).MethodByName(name + "Value").Call(nil)
					require.Equal(t, available, getter[1].Bool())
					if available {
						require.Equal(t, want, physical)
						require.Equal(t, want, getter[0].Float())
						setter := reflect.ValueOf(message).MethodByName("Set" + name + "Value")
						require.True(t, setter.Call([]reflect.Value{reflect.ValueOf(want)})[0].IsNil())
						require.Equal(t, int64(raw), value.Elem().Int())
					}
					message.SetMessageInfo(MessageInfo{})
					encoded, err := message.EncodePayload()
					require.NoError(t, err)
					require.Equal(t, payload, encoded, "raw %x", raw)
				}
			})
		}
	}
	require.Equal(t, 24, checked)
}

func TestOffsetEncodedPowerValidationAndJSON(t *testing.T) {
	m := &UtilityPhaseBAcReactivePower{}
	require.NoError(t, m.SetReactivePowerValue(1_000_000_000))
	object := marshalJSONObject(t, m)
	require.Equal(t, "3000000000", string(object["reactivePower"]))
	checkJSONMeasurement(t, object, "reactivePower", 1_000_000_000, "VAR")
	previous := m.ReactivePower
	for _, invalid := range []float64{-2_000_000_001, 2_294_967_293, 2_294_967_294, 2_294_967_295, math.NaN(), math.Inf(1)} {
		require.ErrorIs(t, m.SetReactivePowerValue(invalid), ErrInvalidPhysicalValue)
		require.Same(t, previous, m.ReactivePower)
	}
	for _, invalid := range []int64{-1, 0x100000000} {
		m.ReactivePower = &invalid
		_, err := m.EncodePayload()
		require.Error(t, err)
	}
}

func TestAirmarPostRepeatingCount(t *testing.T) {
	m := &AirmarPost{}
	require.NoError(t, m.DecodePayload([]byte{0x87, 0x98, 0xfe, 2, 1, 0, 4, 2}))
	m.Repeating1 = append(m.Repeating1, m.Repeating1[0])
	encoded, err := m.EncodePayload()
	require.NoError(t, err)
	require.Equal(t, []byte{0x87, 0x98, 0xfe, 3, 1, 0, 4, 2, 1, 0}, encoded)
	m.Repeating1 = nil
	encoded, err = m.EncodePayload()
	require.NoError(t, err)
	require.Equal(t, []byte{0x87, 0x98, 0xfe, 0}, encoded)
	require.NoError(t, m.DecodePayload(encoded))
	require.Empty(t, m.Repeating1)
}

func TestFusionSignedGainAndBalance(t *testing.T) {
	for _, message := range []PGN{&FusionAuxGain{}, &FusionBalance{}} {
		t.Run(fmt.Sprintf("%T", message), func(t *testing.T) {
			negative := int64(-12)
			switch m := message.(type) {
			case *FusionAuxGain:
				m.Gain = &negative
			case *FusionBalance:
				m.Value = &negative
			}
			wire, err := message.EncodePayload()
			require.NoError(t, err)
			require.Equal(t, byte(0xf4), wire[5])
			require.NoError(t, message.DecodePayload(wire))
			switch m := message.(type) {
			case *FusionAuxGain:
				require.Equal(t, negative, *m.Gain)
			case *FusionBalance:
				require.Equal(t, negative, *m.Value)
			}
		})
	}
}

func TestFusionTunerRDSText(t *testing.T) {
	// PGN 130820, message 32779: FM 107.5 MHz, text in the Basic RDS G0
	// repertoire. The character at 0x91 is ä, not a Latin-1 control code.
	payload := []byte{0xa3, 0x99, 0x0b, 0x80, 1, 2, 0xe0, 0x51, 0x68, 6, 0xc7, 8, 'B', 'o', 'h', 'u', 's', 'l', 0x91, 'n'}
	decoded, err := DecodePayload(MessageInfo{PGN: 130820}, payload)
	require.NoError(t, err)
	m, ok := decoded.(*FusionTuner)
	require.True(t, ok, "%T", decoded)
	require.Equal(t, "Bohuslän", m.Track)
	require.Equal(t, "RDS_G0", structInfoLookup["FusionTuner"].Fields[9].Encoding)
	encoded, err := m.EncodePayload()
	require.NoError(t, err)
	require.Equal(t, payload, encoded)
	data, err := json.Marshal(m)
	require.NoError(t, err)
	require.Contains(t, string(data), `"track":"Bohuslän"`)

	m.Track = "Malmö"
	encoded, err = m.EncodePayload()
	require.NoError(t, err)
	require.Equal(t, []byte{5, 'M', 'a', 'l', 'm', 0x97}, encoded[11:])
	require.NoError(t, m.DecodePayload(encoded))
	require.Equal(t, "Malmö", m.Track)

	for _, value := range []string{"Bohuslän", "¤―║¯", "Ã€", "日本", strings.Repeat("ä", 254)} {
		m.Track = value
		encoded, err = m.EncodePayload()
		require.NoError(t, err)
		require.NoError(t, m.DecodePayload(encoded))
		require.Equal(t, value, m.Track)
	}
	m.Track = strings.Repeat("ä", 255)
	_, err = m.EncodePayload()
	require.Error(t, err)
	m.Track = "\xff"
	_, err = m.EncodePayload()
	require.Error(t, err)
}

func TestRDSTextRepertoire(t *testing.T) {
	// Include a non-UTF-8 byte to force G0 interpretation, including its four
	// differences in the ASCII range and positions changed by later DAB tables.
	require.Equal(t, "ä¤―║¯Ĳĳπ€Ŧŧÿ", decodeRDSText("\x91$^`~\x8f\x9f\xa8\xa9\xee\xfe\xff"))
	require.Equal(t, "Bohuslän $^`~", decodeRDSText("Bohuslän $^`~"))
}
