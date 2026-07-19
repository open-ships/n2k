package gateway

import (
	"testing"
	"time"

	"github.com/open-ships/n2k/internal/framer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// wrapActisense builds a DLE/STX-framed Actisense message around cmd+payload,
// applying DLE escaping and the mod-256 checksum.
func wrapActisense(cmd byte, payload []byte) []byte {
	body := append([]byte{cmd, byte(len(payload))}, payload...)
	var sum byte
	for _, b := range body {
		sum += b
	}
	body = append(body, -sum) // checksum byte: total ≡ 0 (mod 256)

	out := []byte{dle, stx}
	for _, b := range body {
		if b == dle {
			out = append(out, dle, dle)
		} else {
			out = append(out, b)
		}
	}
	return append(out, dle, etx)
}

// n2kPayload builds the payload of an Actisense N2K message (command 0x93).
func n2kPayload(prio uint8, pgnNum uint32, dst, src uint8, data []byte) []byte {
	p := []byte{
		prio,
		byte(pgnNum), byte(pgnNum >> 8), byte(pgnNum >> 16),
		dst, src,
		0, 0, 0, 0, // timestamp (ignored)
		byte(len(data)),
	}
	return append(p, data...)
}

func TestActisenseReader_SingleMessage(t *testing.T) {
	data := []byte{0x01, 0x5C, 0x3D, 0xFF, 0x7F, 0xFF, 0x7F, 0xFC}
	stream := wrapActisense(cmdN2KReceived, n2kPayload(2, 127250, 255, 0x23, data))

	var got []N2KMessage
	r := NewActisenseReader()
	r.Feed(stream, func(m N2KMessage) { got = append(got, m) })

	require.Len(t, got, 1)
	assert.Equal(t, uint8(2), got[0].Priority)
	assert.Equal(t, uint32(127250), got[0].PGN)
	assert.Equal(t, uint8(255), got[0].Destination)
	assert.Equal(t, uint8(0x23), got[0].Source)
	assert.Equal(t, data, got[0].Data)
}

func TestActisenseReader_PreservesTransportTimestamp(t *testing.T) {
	payload := n2kPayload(2, 127250, 255, 0x23, []byte{1})
	payload[6], payload[7], payload[8], payload[9] = 0x78, 0x56, 0x34, 0x12
	stream := wrapActisense(cmdN2KReceived, payload)

	var got N2KMessage
	NewActisenseReader().Feed(stream, func(message N2KMessage) { got = message })
	assert.True(t, got.HasTimestamp)
	assert.Equal(t, time.Duration(0x12345678)*time.Millisecond, got.Timestamp)
}

func TestActisenseReader_SplitAcrossFeeds(t *testing.T) {
	data := []byte{0xAA, 0xBB, 0xCC}
	stream := wrapActisense(cmdN2KReceived, n2kPayload(6, 60928, 255, 0x01, data))

	var got []N2KMessage
	r := NewActisenseReader()
	emit := func(m N2KMessage) { got = append(got, m) }
	r.Feed(stream[:3], emit)
	r.Feed(stream[3:7], emit)
	r.Feed(stream[7:], emit)

	require.Len(t, got, 1)
	assert.Equal(t, uint32(60928), got[0].PGN)
	assert.Equal(t, data, got[0].Data)
}

func TestActisenseReader_EscapedDLEInPayload(t *testing.T) {
	// 0x10 (DLE) inside the CAN data must round-trip through escaping.
	data := []byte{0x10, 0x10, 0x02}
	stream := wrapActisense(cmdN2KReceived, n2kPayload(3, 127250, 255, 0x10, data))

	var got []N2KMessage
	r := NewActisenseReader()
	r.Feed(stream, func(m N2KMessage) { got = append(got, m) })

	require.Len(t, got, 1)
	assert.Equal(t, uint8(0x10), got[0].Source)
	assert.Equal(t, data, got[0].Data)
}

func TestActisenseReader_BadChecksumDropped(t *testing.T) {
	stream := wrapActisense(cmdN2KReceived, n2kPayload(2, 127250, 255, 0x23, []byte{0x01}))
	stream[len(stream)-3] ^= 0xFF // corrupt the checksum byte

	count := 0
	r := NewActisenseReader()
	r.Feed(stream, func(N2KMessage) { count++ })
	assert.Zero(t, count)
}

func TestActisenseReader_OtherCommandIgnored(t *testing.T) {
	stream := wrapActisense(0xA0, []byte{0x01, 0x02})

	count := 0
	r := NewActisenseReader()
	r.Feed(stream, func(N2KMessage) { count++ })
	assert.Zero(t, count)
}

func TestActisenseReader_GarbageBetweenFrames(t *testing.T) {
	m1 := wrapActisense(cmdN2KReceived, n2kPayload(2, 127250, 255, 0x23, []byte{0x01}))
	m2 := wrapActisense(cmdN2KReceived, n2kPayload(2, 127250, 255, 0x24, []byte{0x02}))
	stream := append([]byte{0xDE, 0xAD, 0xBE, 0xEF}, m1...)
	stream = append(stream, 0x55, 0x66)
	stream = append(stream, m2...)

	var got []N2KMessage
	r := NewActisenseReader()
	r.Feed(stream, func(m N2KMessage) { got = append(got, m) })

	require.Len(t, got, 2)
	assert.Equal(t, uint8(0x23), got[0].Source)
	assert.Equal(t, uint8(0x24), got[1].Source)
}

func TestActisenseReader_TruncatedPayloadDropped(t *testing.T) {
	// dataLen claims 8 bytes but only 2 are present.
	p := n2kPayload(2, 127250, 255, 0x23, []byte{0x01, 0x02})
	p[10] = 8
	stream := wrapActisense(cmdN2KReceived, p)

	count := 0
	r := NewActisenseReader()
	r.Feed(stream, func(N2KMessage) { count++ })
	assert.Zero(t, count)
}

func TestActisenseReader_OversizedBodyRecovers(t *testing.T) {
	oversized := append([]byte{dle, stx}, make([]byte, maxActisenseBody+1)...)
	valid := wrapActisense(cmdN2KReceived, n2kPayload(2, 127250, 255, 0x23, []byte{1}))
	stream := append(oversized, valid...)

	var got []N2KMessage
	r := NewActisenseReader()
	r.Feed(stream, func(m N2KMessage) { got = append(got, m) })
	require.Len(t, got, 1)
	assert.Equal(t, []byte{1}, got[0].Data)
}

func FuzzActisenseReader(f *testing.F) {
	f.Add(wrapActisense(cmdN2KReceived, n2kPayload(2, 127250, 255, 0x23, []byte{1, 2, 3})))
	f.Add([]byte{dle, stx, 0x93, 0xFF})
	f.Fuzz(func(t *testing.T, input []byte) {
		r := NewActisenseReader()
		emitted := 0
		for start := 0; start < len(input); {
			end := start + 7
			if end > len(input) {
				end = len(input)
			}
			r.Feed(input[start:end], func(m N2KMessage) {
				emitted++
				assert.LessOrEqual(t, len(m.Data), 244)
			})
			start = end
		}
		assert.LessOrEqual(t, len(r.body), maxActisenseBody)
		assert.LessOrEqual(t, emitted, len(input)/3+1)
	})
}

func TestReframe_SingleFrame(t *testing.T) {
	// PGN 127250 (VesselHeading) is a single-frame PGN.
	data := []byte{0x01, 0x5C, 0x3D, 0xFF, 0x7F, 0xFF, 0x7F, 0xFC}
	frames, err := Reframe(N2KMessage{Priority: 2, PGN: 127250, Destination: 255, Source: 0x23, Data: data}, 0)
	require.NoError(t, err)
	require.Len(t, frames, 1)
	assert.Equal(t, framer.BuildCANID(127250, 2, 0x23, 255), frames[0].ID)
	assert.Equal(t, uint8(8), frames[0].Length)
	assert.Equal(t, data, frames[0].Data[:8])
}

func TestReframe_FastPacket(t *testing.T) {
	// PGN 129029 (GNSS Position Data) is a fast-packet PGN, nominal 43 bytes.
	data := make([]byte, 43)
	for i := range data {
		data[i] = byte(i)
	}
	frames, err := Reframe(N2KMessage{Priority: 3, PGN: 129029, Destination: 255, Source: 0x42, Data: data}, 5)
	require.NoError(t, err)
	// 6 bytes in frame 0, 7 per continuation: 1 + ceil(37/7) = 7 frames.
	require.Len(t, frames, 7)
	seq, frameNum := framer.FastPacketSeqFrame(frames[0].Data[0])
	assert.Equal(t, uint8(5), seq)
	assert.Equal(t, uint8(0), frameNum)
	assert.Equal(t, uint8(43), frames[0].Data[1], "frame 0 carries total length")
	assert.Equal(t, data[:6], frames[0].Data[2:8])
}

func TestReframe_UnknownPGNBySize(t *testing.T) {
	// Unknown PGN (60000 has no metadata entry), small payload → single frame.
	frames, err := Reframe(N2KMessage{Priority: 6, PGN: 60000, Destination: 255, Source: 1, Data: []byte{1, 2, 3, 4}}, 0)
	require.NoError(t, err)
	require.Len(t, frames, 1)

	// Unknown PGN, large payload → fast-packet.
	big := make([]byte, 20)
	frames, err = Reframe(N2KMessage{Priority: 6, PGN: 60000, Destination: 255, Source: 1, Data: big}, 0)
	require.NoError(t, err)
	assert.Greater(t, len(frames), 1)
}

func TestReframe_TooLarge(t *testing.T) {
	_, err := Reframe(N2KMessage{PGN: 129029, Data: make([]byte, 224)}, 0)
	require.Error(t, err)
}

func TestEncodeSend_Layout(t *testing.T) {
	msg := N2KMessage{Priority: 2, PGN: 127250, Destination: 255, Data: []byte{0xFF, 0x88, 0x4C, 0xFF, 0x7F, 0xFF, 0x7F, 0xFD}}
	buf, err := EncodeSend(msg)
	require.NoError(t, err)
	// Send payload: prio, PGN little-endian, dst, len, data -- no source or
	// timestamp bytes.
	want := wrapActisense(cmdN2KSend, append([]byte{2, 0x12, 0xF1, 0x01, 255, 8}, msg.Data...))
	assert.Equal(t, want, buf)
}

func TestEncodeSend_EscapesDLE(t *testing.T) {
	msg := N2KMessage{Priority: 6, PGN: 60928, Destination: 255, Data: []byte{0x10, 0x10, 0x03}}
	buf, err := EncodeSend(msg)
	require.NoError(t, err)
	want := wrapActisense(cmdN2KSend, append([]byte{6, 0x00, 0xEE, 0x00, 255, 3}, msg.Data...))
	assert.Equal(t, want, buf)
}

func TestEncodeSend_TooLarge(t *testing.T) {
	_, err := EncodeSend(N2KMessage{PGN: 126996, Data: make([]byte, 224)})
	require.Error(t, err)
}

func TestEncodeStartup(t *testing.T) {
	assert.Equal(t, []byte{0x10, 0x02, 0xA1, 0x03, 0x11, 0x02, 0x00, 0x49, 0x10, 0x03}, EncodeStartup())
}
