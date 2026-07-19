package adapter

import (
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/decoder"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFastPacketOwnsFrameBytes(t *testing.T) {
	m := NewMultiBuilder()
	info := NewPacketInfo(&can.Frame{ID: framer.BuildCANID(130820, 1, 10, 255), Length: 8})
	frameZero := []byte{0x20, 8, 1, 2, 3, 4, 5, 6}
	p := decoder.NewPacket(info, frameZero)
	m.Add(p)

	for i := range frameZero {
		frameZero[i] = 0xEE
	}
	last := decoder.NewPacket(info, []byte{0x21, 7, 8})
	m.Add(last)
	require.True(t, last.Complete)
	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8}, last.Data)
}

func TestFastPacketSessionsAreIsolatedByNetwork(t *testing.T) {
	m := NewMultiBuilder()
	base := NewPacketInfo(&can.Frame{ID: framer.BuildCANID(130820, 1, 10, 255), Length: 8})
	infoA, infoB := base, base
	infoA.NetworkID = "can0"
	infoB.NetworkID = "can1"

	a0 := decoder.NewPacket(infoA, []byte{0x20, 8, 1, 2, 3, 4, 5, 6})
	b0 := decoder.NewPacket(infoB, []byte{0x20, 8, 11, 12, 13, 14, 15, 16})
	m.Add(a0)
	m.Add(b0)
	a1 := decoder.NewPacket(infoA, []byte{0x21, 7, 8})
	b1 := decoder.NewPacket(infoB, []byte{0x21, 17, 18})
	m.Add(a1)
	m.Add(b1)

	require.True(t, a1.Complete)
	require.True(t, b1.Complete)
	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8}, a1.Data)
	assert.Equal(t, []byte{11, 12, 13, 14, 15, 16, 17, 18}, b1.Data)
}

func TestMultiBuilderEvictsExpiredAndCapsState(t *testing.T) {
	m := NewMultiBuilder()
	now := time.Unix(100, 0)
	m.now = func() time.Time { return now }
	m.ttl = time.Second
	m.maxActive = 2
	addContinuation := func(source uint8) {
		info := NewPacketInfo(&can.Frame{ID: framer.BuildCANID(130820, 1, source, 255), Length: 8})
		m.Add(decoder.NewPacket(info, []byte{1, 1, 2, 3, 4, 5, 6, 7}))
	}

	addContinuation(1)
	addContinuation(2)
	addContinuation(3)
	assert.LessOrEqual(t, activeFastPacketSequences(m), 2)

	now = now.Add(time.Second)
	addContinuation(4)
	assert.Equal(t, 1, activeFastPacketSequences(m))
}

func activeFastPacketSequences(m *MultiBuilder) int {
	return len(m.sequences)
}

func FuzzCANAdapter(f *testing.F) {
	f.Add(uint32(0x09F20183), byte(8), []byte{0x60, 0x20})
	f.Add(uint32(0x09F20D00), byte(8), []byte{})
	f.Fuzz(func(t *testing.T, id uint32, length byte, input []byte) {
		a := NewCANAdapter()
		a.SetOutput(&mockHandler{})
		var data [8]byte
		copy(data[:], input)
		frame := can.Frame{ID: id, Length: length, Data: data}
		require.NotPanics(t, func() { a.HandleMessage(&frame) })
		assert.LessOrEqual(t, activeFastPacketSequences(a.multi), defaultMaxFastPacketState)
	})
}
