package pgn

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchemaV92WireLayouts(t *testing.T) {
	tests := []struct {
		name   string
		number uint32
		wire   string
		check  func(*testing.T, PGN)
	}{
		{"heave delay", 127252, "0185fff401f2ffff", func(t *testing.T, p PGN) {
			m := p.(*Heave)
			require.Equal(t, int64(-123), *m.Heave)
			v, ok := m.DelayValue()
			require.True(t, ok)
			require.Equal(t, 5.0, v)
			require.Equal(t, uint64(DelaySourceFactoryDefault), *m.DelaySource)
			require.NoError(t, m.SetDelayValue(5))
		}},
		{"BEP brightness", 65294, "2799530102030405", func(t *testing.T, p PGN) {
			m := p.(*BepMarineProprietaryPgn65294)
			require.Equal(t, uint64(83), *m.Brightness)
			require.Equal(t, []byte{1, 2, 3, 4, 5}, m.Data)
		}},
		{"Furuno command", 61184, "3f9f0002ffffffff", func(t *testing.T, p PGN) {
			require.Equal(t, uint64(2), *p.(*FurunoCommand).Command)
		}},
		{"Furuno keyboard", 61184, "3f9f020405060708", func(t *testing.T, p PGN) {
			m := p.(*FurunoRemoteKeyboardReport)
			require.Equal(t, uint64(4), *m.Key1)
			require.Equal(t, uint64(8), *m.Key5)
		}},
		{"Furuno pointer", 61184, "3f9f03fe03ff0109", func(t *testing.T, p PGN) {
			m := p.(*FurunoRemotePointerReport)
			require.Equal(t, int64(-2), *m.Dx)
			require.Equal(t, int64(3), *m.Dy)
			require.Equal(t, int64(-1), *m.Wheel)
			require.Equal(t, uint64(1), *m.Button)
			require.Equal(t, uint64(1), *m.KeypadMode)
		}},
		{"Furuno announce", 65287, "3f9f01ffffffffff", func(t *testing.T, p PGN) {
			require.Equal(t, uint64(1), *p.(*FurunoRemoteAnnounce).Host)
		}},
		{"Furuno backlight", 65288, "3f9f012364ffffff", func(t *testing.T, p PGN) {
			m := p.(*FurunoRemoteBacklight)
			require.Equal(t, uint64(2), *m.Preset)
			require.Equal(t, uint64(100), *m.Brightness)
		}},
		{"Furuno info", 130816, "3f9ff0ff0102", func(t *testing.T, p PGN) {
			require.Equal(t, uint64(2), *p.(*FurunoRemoteInfo).B)
		}},
		{"Furuno host", 130849, "3f9f0102030405060708", func(t *testing.T, p PGN) {
			require.Equal(t, uint64(0x0807060504030201), *p.(*FurunoRemoteHost).Host)
		}},
		{"Simnet catch-all event", 130850, "419ffeff000211ffffffffff", func(t *testing.T, p PGN) {
			m := p.(*SimnetEvent)
			require.Equal(t, uint64(2), *m.CommandType)
			require.Equal(t, uint64(17), *m.Event)
		}},
		{"Furuno high-bit azimuth", 130845, "3f9f1111ff0101051027409cc4090000000001", func(t *testing.T, p PGN) {
			m := p.(*FurunoMultiSatsInViewExtended)
			require.Len(t, m.Repeating1, 1)
			require.Equal(t, uint64(40000), *m.Repeating1[0].Azimuth)
			v, ok := m.Repeating1[0].AzimuthValue()
			require.True(t, ok)
			require.Equal(t, 4.0, v)
			require.NoError(t, m.Repeating1[0].SetAzimuthValue(4))
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wire, err := hex.DecodeString(test.wire)
			require.NoError(t, err)
			m, err := DecodePayload(MessageInfo{PGN: test.number}, wire)
			require.NoError(t, err)
			test.check(t, m)
			// Bypass raw replay to independently check the generated layout.
			m.SetMessageInfo(MessageInfo{})
			encoded, err := m.EncodePayload()
			require.NoError(t, err)
			require.Equal(t, wire, encoded)
		})
	}
}

func TestSchemaV92CzoneOutputs(t *testing.T) {
	for _, number := range []uint32{130817, 130822, 130825} {
		t.Run(fmt.Sprint(number), func(t *testing.T) {
			wire := []byte{0x27, 0x99, 1, 2}
			if number == 130825 {
				wire = append(wire, 0xa5)
				data, err := hex.DecodeString("0a40df02d2c700f535607d0e60df03da0701f745e07d")
				require.NoError(t, err)
				wire = append(wire, data...)
			} else {
				wire = append(wire, bytes.Repeat([]byte{5, 0xf4, 0xa9}, 8)...)
			}
			m, err := DecodePayload(MessageInfo{PGN: number}, wire)
			require.NoError(t, err)
			value := reflect.ValueOf(m).Elem()
			for i := 0; i < 8; i++ {
				if number == 130825 {
					require.Equal(t, uint64(10+i), value.FieldByName(fmt.Sprintf("OutputCurrentRaw%d", i)).Elem().Uint())
					require.Equal(t, uint64(1000+i), value.FieldByName(fmt.Sprintf("OutputLevelRaw%d", i)).Elem().Uint())
					continue
				}
				require.Equal(t, uint64(5), value.FieldByName(fmt.Sprintf("OutputCurrent%d", i)).Elem().Uint())
				require.Equal(t, uint64(500), value.FieldByName(fmt.Sprintf("OutputLevel%d", i)).Elem().Uint())
				require.Equal(t, []byte{42}, value.FieldByName(fmt.Sprintf("OutputUnknownFlags%d", i)).Bytes())
				current, unit, valid, err := PhysicalValue(m, 6+3*i)
				require.NoError(t, err)
				require.True(t, valid)
				require.Equal(t, "A", unit)
				want := 0.5
				if number == 130817 {
					want = 1
				}
				require.Equal(t, want, current)
			}
			m.SetMessageInfo(MessageInfo{})
			encoded, err := m.EncodePayload()
			require.NoError(t, err)
			require.Equal(t, wire, encoded)
		})
	}
}

func TestSchemaV92FastPacketLimits(t *testing.T) {
	for _, m := range []PGN{&LumishoreProprietary{}, &MaretronProprietaryConfiguration{}, &NavicoBoatSpeedPolarTable{}, &Message0x1ff000x1ffffManufacturerSpecificFastPacketNonAddressed{}} {
		t.Run(reflect.TypeOf(m).Elem().Name(), func(t *testing.T) {
			wire, err := m.EncodePayload()
			require.NoError(t, err)
			require.Len(t, wire, 223)
			require.NoError(t, m.DecodePayload(wire))
			m.SetMessageInfo(MessageInfo{})
			again, err := m.EncodePayload()
			require.NoError(t, err)
			require.Equal(t, wire, again)
		})
	}
}

func TestSchemaV92DscSelectionAndDecimals(t *testing.T) {
	for _, test := range []struct {
		format, category byte
		name, address    string
	}{
		{112, 100, "DscDistressCallInformation", "003456789"},
		{112, 112, "DscDistressCallInformation", "003456789"},
		{120, 112, "DscDistressAcknowledgementOrRelay", "003456789"},
		{120, 100, "DscCallInformation", "003456789"},
		{102, 100, "DscCallInformation", "0034567890"},
	} {
		t.Run(fmt.Sprintf("%d/%d", test.format, test.category), func(t *testing.T) {
			wire := []byte{test.format, test.category, 0, 34, 56, 78, 90}
			m, err := DecodePayload(MessageInfo{PGN: 129808}, wire)
			require.NoError(t, err)
			require.Equal(t, test.name, reflect.TypeOf(m).Elem().Name())
			v := reflect.ValueOf(m).Elem()
			require.Equal(t, test.address, v.FieldByName("DscMessageAddress").String())
			raw, err := json.Marshal(m)
			require.NoError(t, err)
			require.Contains(t, string(raw), `"dscMessageAddress":"`+test.address+`"`)
			encoded, err := m.EncodePayload()
			require.NoError(t, err)
			require.Equal(t, wire, encoded)
			// Mutation must invalidate replay and encode decimal pairs, retaining
			// leading zeros and appending the MMSI padding digit exactly once.
			v.FieldByName("MmsiOfShipInDistress").SetString("012345678")
			encoded, err = m.EncodePayload()
			require.NoError(t, err)
			require.Equal(t, wire, encoded[:7])
			require.Equal(t, []byte{1, 23, 45, 67, 80}, encoded[35:40])
			decoded := structTypeRegistry[test.name]()
			require.NoError(t, decoded.DecodePayload(encoded))
			require.Equal(t, "012345678", reflect.ValueOf(decoded).Elem().FieldByName("MmsiOfShipInDistress").String())
		})
	}
}

func TestSchemaV92DscUnavailableAndInvalidAddresses(t *testing.T) {
	for _, data := range [][]byte{{0xff, 0xff, 0xff, 0xff, 0xff}, {0, 34, 100, 78, 90}, {0, 34, 56, 78}} {
		wire := append([]byte{120, 100}, data...)
		m := &DscCallInformation{}
		require.NoError(t, m.DecodePayload(wire))
		require.Empty(t, m.DscMessageAddress)
		encoded, err := m.EncodePayload()
		require.NoError(t, err)
		require.Equal(t, wire, encoded)
	}
	for _, address := range []string{"123", "12345678901", "12345678x", "0034567890"} {
		m := &DscCallInformation{DscMessageAddress: address}
		_, err := m.EncodePayload()
		require.Error(t, err)
	}
	format := uint64(102)
	_, err := (&DscCallInformation{DscFormatSymbol: &format, DscMessageAddress: "003456789"}).EncodePayload()
	require.Error(t, err)
}
