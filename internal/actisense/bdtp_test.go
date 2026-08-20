package actisense

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParserType1SplitAndEscapedDLE(t *testing.T) {
	wire, err := EncodeDatagram(BSTN2KReceive, []byte{1, DLE, 2})
	require.NoError(t, err)

	var got []Datagram
	parser := NewParser()
	for _, chunk := range [][]byte{wire[:1], wire[1:4], wire[4:7], wire[7:]} {
		parser.Feed(chunk, func(d Datagram) { got = append(got, d) }, func(err DecodeError) { t.Fatal(err) })
	}
	require.Len(t, got, 1)
	assert.Equal(t, byte(BSTN2KReceive), got[0].ID)
	assert.Equal(t, []byte{1, DLE, 2}, got[0].Payload)
}

func TestParserType2MaximumPayload(t *testing.T) {
	payload := bytes.Repeat([]byte{0xA5}, MaxDatagramLength-3)
	wire, err := EncodeDatagram(BSTN2KMessage, payload)
	require.NoError(t, err)

	var got Datagram
	NewParser().Feed(wire, func(d Datagram) { got = d }, func(err DecodeError) { t.Fatal(err) })
	assert.Equal(t, payload, got.Payload)
	assert.Len(t, got.Raw, MaxDatagramLength)
}

func TestParserReportsChecksumAndRecovers(t *testing.T) {
	bad, err := EncodeDatagram(0x22, []byte{1, 2, 3})
	require.NoError(t, err)
	bad[len(bad)-3]++
	good, err := EncodeDatagram(0x23, []byte{4})
	require.NoError(t, err)

	var (
		got    []Datagram
		errors []DecodeError
	)
	parser := NewParser()
	parser.Feed(append(bad, good...), func(d Datagram) { got = append(got, d) }, func(err DecodeError) { errors = append(errors, err) })
	require.Len(t, errors, 1)
	assert.Equal(t, DecodeChecksum, errors[0].Kind)
	require.Len(t, got, 1)
	assert.Equal(t, byte(0x23), got[0].ID)
}

func TestDecodeRawRejectsLengthMismatch(t *testing.T) {
	_, err := DecodeRaw([]byte{0x93, 2, 1})
	require.Error(t, err)
	_, err = DecodeRaw([]byte{0xD0, 4, 0, 1, 2})
	require.Error(t, err)
}

func TestParserOversizeIsBoundedAndRecovers(t *testing.T) {
	oversized := append([]byte{DLE, STX}, bytes.Repeat([]byte{0x55}, MaxDatagramLength+2)...)
	oversized = append(oversized, DLE, ETX)
	good, err := EncodeDatagram(0x20, []byte{1})
	require.NoError(t, err)

	var (
		got    []Datagram
		errors []DecodeError
	)
	parser := NewParser()
	parser.Feed(append(oversized, good...), func(d Datagram) { got = append(got, d) }, func(err DecodeError) { errors = append(errors, err) })
	assert.LessOrEqual(t, parser.BufferedBytes(), MaxDatagramLength+1)
	require.Len(t, errors, 1)
	assert.Equal(t, DecodeOversize, errors[0].Kind)
	require.Len(t, got, 1)
}

func TestParserReportsExactUnframedBytesAcrossChunks(t *testing.T) {
	wire, err := EncodeDatagram(BSTN2KReceive, []byte{1})
	require.NoError(t, err)
	input := append([]byte("boot\r\n"), DLE, DLE)
	input = append(input, wire[1:]...)
	input = append(input, 'x', DLE)

	parser := NewParser()
	var (
		unframed []byte
		dgrams   []Datagram
	)
	for _, chunk := range [][]byte{input[:3], input[3:8], input[8:]} {
		parser.FeedWithUnframed(chunk, func(datagram Datagram) { dgrams = append(dgrams, datagram) }, func(err DecodeError) { t.Fatal(err) }, func(data []byte) {
			unframed = append(unframed, data...)
		})
	}
	parser.EndWithUnframed(func(err DecodeError) { t.Fatal(err) }, func(data []byte) {
		unframed = append(unframed, data...)
	})

	require.Len(t, dgrams, 1)
	assert.Equal(t, append([]byte("boot\r\n"), DLE, 'x', DLE), unframed)
	assert.Equal(t, uint64(len(unframed)), parser.UnframedBytes())
}

func FuzzParser(f *testing.F) {
	seed, err := EncodeDatagram(BSTN2KReceive, []byte{1, 2, 3})
	require.NoError(f, err)
	f.Add(seed)
	f.Add([]byte{DLE, STX, 0xD0, 0xFF, 0xFF})
	f.Fuzz(func(t *testing.T, input []byte) {
		parser := NewParser()
		emitted := 0
		for start := 0; start < len(input); start += 7 {
			end := min(start+7, len(input))
			parser.Feed(input[start:end], func(d Datagram) {
				emitted++
				assert.LessOrEqual(t, len(d.Raw), MaxDatagramLength)
			}, func(DecodeError) {})
		}
		assert.LessOrEqual(t, parser.BufferedBytes(), MaxDatagramLength+1)
		assert.LessOrEqual(t, emitted, len(input)/3+1)
	})
}
