package transport

import (
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/require"
)

func passiveFrames(destination uint8) []can.Frame {
	control := ControlBAM
	if destination != BroadcastAddr {
		control = ControlRTS
	}
	return []can.Frame{
		{ID: framer.BuildCANID(PGNCM, 6, 66, destination), Length: 8, Data: [8]byte{control, 9, 0, 2, 2, 0xD8, 0xFE, 0}},
		{ID: framer.BuildCANID(PGNDT, 6, 66, destination), Length: 8, Data: [8]byte{1, 1, 2, 3, 4, 5, 6, 7}},
		{ID: framer.BuildCANID(PGNDT, 6, 66, destination), Length: 8, Data: [8]byte{2, 8, 42, 255, 255, 255, 255, 255}},
	}
}
func TestPassiveReceiverAssemblesBAMAndAddressed(t *testing.T) {
	for _, dest := range []uint8{255, 42} {
		r := NewPassiveReceiver()
		stamp := time.Now()
		info := pgn.MessageInfo{SourceId: 66, Timestamp: stamp, NetworkID: "a", ConnectionEpoch: 2, ClaimEpoch: 3}
		var got []byte
		var meta pgn.MessageInfo
		for _, frame := range passiveFrames(dest) {
			var err error
			meta, got, err = r.Handle(frame, info)
			require.NoError(t, err)
		}
		require.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 42}, got)
		require.Equal(t, uint32(65240), meta.PGN)
		require.Equal(t, stamp, meta.Timestamp)
		if dest == 255 {
			require.Nil(t, meta.TargetId)
		} else {
			require.Equal(t, dest, *meta.TargetId)
		}
		require.Empty(t, r.sessions)
	}
}
func TestPassiveReceiverIsolatesNetworksEpochsAndDestinations(t *testing.T) {
	for _, kind := range []string{"network", "connection", "claim", "destination"} {
		t.Run(kind, func(t *testing.T) {
			r := NewPassiveReceiver()
			first := pgn.MessageInfo{NetworkID: "a", ConnectionEpoch: 1, ClaimEpoch: 1}
			second := first
			dest := uint8(42)
			switch kind {
			case "network":
				second.NetworkID = "b"
			case "connection":
				second.ConnectionEpoch = 2
			case "claim":
				second.ClaimEpoch = 2
			case "destination":
				dest = 43
			}
			a, b := passiveFrames(42), passiveFrames(dest)
			_, _, err := r.Handle(a[0], first)
			require.NoError(t, err)
			_, _, err = r.Handle(a[1], first)
			require.NoError(t, err)
			_, data, err := r.Handle(b[2], second)
			require.NoError(t, err)
			require.Nil(t, data)
			_, data, err = r.Handle(a[2], first)
			require.NoError(t, err)
			require.Len(t, data, 9)
		})
	}
}
func TestPassiveReceiverRejectsBrokenTransfers(t *testing.T) {
	for _, kind := range []string{"length", "count", "order", "capture timeout", "wall timeout", "abort", "reverse abort"} {
		t.Run(kind, func(t *testing.T) {
			r := NewPassiveReceiver()
			now := time.Now()
			r.now = func() time.Time { return now }
			frames := passiveFrames(42)
			info := pgn.MessageInfo{Timestamp: now}
			if kind == "count" {
				frames[0].Data[3] = 3
			}
			_, _, err := r.Handle(frames[0], info)
			if kind == "count" {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			switch kind {
			case "length":
				frames[1].Length = 3
			case "order":
				frames[1].Data[0] = 2
			case "capture timeout":
				info.Timestamp = now.Add(2 * time.Second)
			case "wall timeout":
				now = now.Add(2 * time.Second)
			case "abort", "reverse abort":
				abort := frames[0]
				abort.Data[0] = ControlAbort
				if kind == "reverse abort" {
					abort.ID = framer.BuildCANID(PGNCM, 6, 42, 66)
				}
				_, _, err = r.Handle(abort, info)
				require.NoError(t, err)
			}
			_, data, err := r.Handle(frames[1], info)
			if kind == "length" || kind == "order" {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Nil(t, data)
			_, data, err = r.Handle(frames[2], info)
			require.NoError(t, err)
			require.Nil(t, data)
			require.Empty(t, r.sessions)
		})
	}
}
func TestPassiveReceiverBoundsState(t *testing.T) {
	r := NewPassiveReceiver()
	frame := passiveFrames(255)[0]
	for i := 0; i < maxPassiveSessions; i++ {
		_, _, err := r.Handle(frame, pgn.MessageInfo{ConnectionEpoch: uint64(i)})
		require.NoError(t, err)
	}
	_, _, err := r.Handle(frame, pgn.MessageInfo{ConnectionEpoch: maxPassiveSessions})
	require.ErrorContains(t, err, "table exceeds")
	require.Len(t, r.sessions, maxPassiveSessions)
}

func TestPassiveReceiverCTSCannotKeepTransferAliveForever(t *testing.T) {
	for _, captureClock := range []bool{false, true} {
		r := NewPassiveReceiver()
		now := time.Now()
		r.now = func() time.Time { return now }
		info := pgn.MessageInfo{Timestamp: now}
		frames := passiveFrames(42)
		_, _, err := r.Handle(frames[0], info)
		require.NoError(t, err)
		cts := frames[0]
		cts.ID = framer.BuildCANID(PGNCM, 6, 42, 66)
		cts.Data[0] = ControlCTS
		for i := 0; i < 31; i++ {
			if captureClock {
				info.Timestamp = info.Timestamp.Add(time.Second)
			} else {
				now = now.Add(time.Second)
			}
			_, _, err := r.Handle(cts, info)
			require.NoError(t, err)
		}
		require.Empty(t, r.sessions)
	}
}

func FuzzPassiveReceiver(f *testing.F) {
	var seed []byte
	for _, frame := range passiveFrames(255) {
		seed = append(seed, byte(frame.ID), byte(frame.ID>>8), byte(frame.ID>>16), byte(frame.ID>>24), frame.Length)
		seed = append(seed, frame.Data[:]...)
	}
	f.Add(seed)
	f.Fuzz(func(t *testing.T, data []byte) {
		r := NewPassiveReceiver()
		for len(data) >= 13 {
			frame := can.Frame{ID: uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24, Length: data[4]}
			copy(frame.Data[:], data[5:13])
			_, payload, _ := r.Handle(frame, pgn.MessageInfo{})
			if len(payload) > 1785 || len(r.sessions) > maxPassiveSessions {
				t.Fatal("unbounded passive transfer")
			}
			data = data[13:]
		}
	})
}
