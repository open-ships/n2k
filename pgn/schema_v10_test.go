package pgn

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestH5000TableCaptures(t *testing.T) {
	data, err := os.ReadFile("../testdata/h5000-tables.raw")
	require.NoError(t, err)
	// Values per table in the independent H5000 capture. The last table of
	// section 15 declares 22 values but the fast packet ends after 19.
	counts := map[byte][]int{11: {7, 7, 7, 7, 7}, 12: {7, 7, 7}, 13: {7, 7, 7}, 14: {7, 7}, 15: {22, 19}, 18: {7, 7, 7, 7, 7}}
	for line := range strings.SplitSeq(strings.TrimSpace(string(data)), "\n") {
		columns := strings.SplitN(line, ",", 7)
		require.Len(t, columns, 7)
		wire, err := hex.DecodeString(strings.ReplaceAll(columns[6], ",", ""))
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(wire), 14)
		section := wire[5]
		want, ok := counts[section]
		require.True(t, ok, "unexpected section %d", section)
		delete(counts, section)
		t.Run(fmt.Sprint(section), func(t *testing.T) {
			m, err := DecodePayload(MessageInfo{PGN: 130822}, wire)
			require.NoError(t, err)
			require.Equal(t, fmt.Sprintf("*pgn.NavicoUdbDatabaseTableDump%d", section), fmt.Sprintf("%T", m))
			encoded, err := m.EncodePayload()
			require.NoError(t, err)
			require.Equal(t, wire, encoded)
			raw, err := json.Marshal(m)
			require.NoError(t, err)
			var dump struct {
				Tables []struct {
					Count   int
					Records []struct {
						Length int
						Value  *float32
					} `json:"repeating2"`
				} `json:"repeating1"`
			}
			require.NoError(t, json.Unmarshal(raw, &dump))
			require.Len(t, dump.Tables, len(want))
			offset := 14
			for i, table := range dump.Tables {
				require.Equal(t, int(binary.LittleEndian.Uint16(wire[offset:])), table.Count)
				offset += 2
				require.Len(t, table.Records, want[i])
				for _, record := range table.Records {
					require.Equal(t, 4, record.Length)
					require.NotNil(t, record.Value)
					require.Equal(t, binary.LittleEndian.Uint32(wire[offset+1:]), math.Float32bits(*record.Value))
					offset += 5
				}
			}
			require.Equal(t, len(wire), offset)
		})
	}
	require.Empty(t, counts)
}

func TestNestedTableOwnershipAndCounts(t *testing.T) {
	wire, err := hex.DecodeString("1399ffc6000b002ffd0000005c00" +
		"0200040000803f04000000c0" + "0000" + "01000400004040")
	require.NoError(t, err)
	msg, err := DecodePayload(MessageInfo{PGN: 130822}, wire)
	require.NoError(t, err)
	table := msg.(*NavicoUdbDatabaseTableDump11)
	require.Len(t, table.Repeating1, 3)
	require.Empty(t, table.Repeating1[1].Repeating2)
	clone := table.Clone().(*NavicoUdbDatabaseTableDump11)
	*clone.Repeating1[0].Repeating2[0].Value = -4
	clone.Repeating1[0].Repeating2 = clone.Repeating1[0].Repeating2[:1]
	encoded, err := clone.EncodePayload()
	require.NoError(t, err)
	require.Equal(t, append(bytes.Clone(wire[:14]), []byte{1, 0, 4, 0, 0, 0x80, 0xc0, 0, 0, 1, 0, 4, 0, 0, 0x40, 0x40}...), encoded)
	original, err := table.EncodePayload()
	require.NoError(t, err)
	require.Equal(t, wire, original)
	require.Equal(t, float32(1), *table.Repeating1[0].Repeating2[0].Value)
	// A count ending at EOF still identifies an empty table.
	require.NoError(t, table.DecodePayload(wire[:28]))
	require.Len(t, table.Repeating1, 2)
	require.Empty(t, table.Repeating1[1].Repeating2)
	// Cut each byte of the final float off the wire. Keep the table header
	// and received records, and replay the partial tail verbatim.
	for n := 1; n <= 4; n++ {
		truncated := wire[:len(wire)-n]
		require.NoError(t, table.DecodePayload(truncated))
		require.Len(t, table.Repeating1, 3)
		require.Equal(t, uint64(1), *table.Repeating1[2].Count)
		require.Empty(t, table.Repeating1[2].Repeating2)
		replayed, err := table.EncodePayload()
		require.NoError(t, err)
		require.Equal(t, truncated, replayed)
	}
}

func TestNavicoSettingsDumpLengths(t *testing.T) {
	wire, err := hex.DecodeString("1399ffc6000a012ffd0000005c00" + "01ff0107" + "01ff040000803f" + "01ff00")
	require.NoError(t, err)
	msg, err := DecodePayload(MessageInfo{PGN: 130822}, wire)
	require.NoError(t, err)
	settings := msg.(*NavicoUdbDatabaseSettingsDump)
	require.Len(t, settings.Repeating1, 3)
	require.Equal(t, []byte{7}, settings.Repeating1[0].Value)
	require.Equal(t, []byte{0, 0, 0x80, 0x3f}, settings.Repeating1[1].Value)
	require.Empty(t, settings.Repeating1[2].Value)
	settings.Repeating1[1].Value = []byte{0xfe}
	encoded, err := settings.EncodePayload()
	require.NoError(t, err)
	want := append(bytes.Clone(wire[:18]), []byte{1, 0xff, 1, 0xfe, 1, 0xff, 0}...)
	require.Equal(t, want, encoded)
}

func TestSchemaV10SignedReadingsAndMasks(t *testing.T) {
	var pilot FurunoNavpilotStatus
	require.NoError(t, pilot.DecodePayload([]byte{0x3f, 0x9f, 1, 0x78, 0xec}))
	value, valid := pilot.RudderAngleValue()
	require.True(t, valid)
	require.Equal(t, -0.5, value)
	require.NoError(t, pilot.DecodePayload([]byte{0x3f, 0x9f, 1, 0xff, 0x7f}))
	require.Nil(t, pilot.RudderAngle)

	var seatalk Seatalk1PilotMode
	require.NoError(t, seatalk.DecodePayload([]byte{0x3b, 0x9f, 0xf0, 0x81, 0x84, 0, 0, 0, 0, 0, 0xf6}))
	position, valid := seatalk.PilotModeDataValue()
	require.True(t, valid)
	require.Equal(t, -10.0, position)

	var mercury MercuryEngineStatus
	require.NoError(t, mercury.DecodePayload([]byte{0x90, 0x98, 0, 0xf3, 3}))
	require.Nil(t, mercury.Helm)
	require.Equal(t, uint64(3), *mercury.SubHelm)
	require.Equal(t, MercuryEngineCapabilityBeacon|MercuryEngineCapabilityCruiseControl, MercuryEngineCapabilityConst(*mercury.Capabilities))
	var czone BepMarineCzoneCircuitStatus
	require.NoError(t, czone.DecodePayload([]byte{0x27, 0x99, 0, 0, 1, 0, 0, 0x80}))
	require.Equal(t, uint64(0x80000001), *czone.Bitmap)
}

func TestSchemaV10TrailingSpares(t *testing.T) {
	for _, tt := range []struct {
		msg               PGN
		previous, current int
	}{
		{&FurunoSensorSetup{}, 28, 30},
		{&FusionVersions{}, 11, 12},
		{&FusionTrackPosition{}, 8, 9},
	} {
		wire, err := tt.msg.EncodePayload()
		require.NoError(t, err)
		require.Len(t, wire, tt.current)
		require.Equal(t, make([]byte, tt.current-tt.previous), wire[tt.previous:])
		require.NoError(t, tt.msg.DecodePayload(wire[:tt.previous]))
		replayed, err := tt.msg.EncodePayload()
		require.NoError(t, err)
		require.Equal(t, wire[:tt.previous], replayed)
	}
	var setup FurunoSensorSetup
	for _, pa := range []float64{-327680, 327640} {
		require.NoError(t, setup.SetAirPressureOffsetValue(pa))
		got, ok := setup.AirPressureOffsetValue()
		require.True(t, ok)
		require.Equal(t, pa, got)
	}
}
