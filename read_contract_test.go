package n2k

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
	"github.com/open-ships/n2k/raw"
	"github.com/stretchr/testify/require"
)

func TestReadOnlyTCPReconnectIsolatesFastPackets(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() { _ = listener.Close() }()
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	id := framer.BuildCANID(126998, 6, 66, 255)
	old := framer.FrameFastPacket(id, []byte{3, 1, 'a', 3, 1, 'b', 3, 1, 'c'}, 0)
	current := framer.FrameFastPacket(id, []byte{3, 1, 'x', 3, 1, 'y', 3, 1, 'z'}, 0)
	served := make(chan error, 1)
	go func() {
		for _, frames := range [][]can.Frame{{old[0]}, {current[1], current[0], current[1]}} {
			conn, err := listener.Accept()
			if err != nil {
				served <- err
				return
			}
			for _, frame := range frames {
				line := fmt.Sprintf("12:00:00.000 R %08X", frame.ID)
				for _, b := range frame.Data[:frame.Length] {
					line += fmt.Sprintf(" %02X", b)
				}
				if _, err := fmt.Fprintln(conn, line); err != nil {
					_ = conn.Close()
					served <- err
					return
				}
			}
			_ = conn.Close()
		}
		served <- nil
	}()
	scanner := NewScanner(ctx, YachtDevicesTCP(listener.Addr().String()), WithReconnect(ReconnectPolicy{InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond}))
	defer func() { _ = scanner.Close() }()
	require.True(t, scanner.Next(), "scanner error: %v", scanner.Err())
	msg := scanner.Message().(*pgn.ConfigurationInformation)
	require.Equal(t, "xyz", msg.InstallationDescription1+msg.InstallationDescription2+msg.ManufacturerInformation)
	require.Equal(t, uint64(2), msg.Info.ConnectionEpoch)
	require.NoError(t, <-served)
}

func TestReadOnlyBAMWithTransportedPGNFilter(t *testing.T) {
	frames := []can.Frame{
		{ID: framer.BuildCANID(60416, 6, 66, 255), Length: 8, Data: [8]byte{32, 9, 0, 2, 255, 0xD8, 0xFE, 0}},
		{ID: framer.BuildCANID(60160, 6, 66, 255), Length: 8, Data: [8]byte{1, 1, 2, 3, 4, 5, 6, 7}},
		{ID: framer.BuildCANID(60160, 6, 66, 255), Length: 8, Data: [8]byte{2, 8, 42, 255, 255, 255, 255, 255}},
	}
	var got []pgn.Message
	for msg, err := range Receive(t.Context(), Replay(frames), Filter("pgn == 65240")) {
		require.NoError(t, err)
		got = append(got, msg)
	}
	require.Len(t, got, 1)
	payload, err := pgn.EncodeMessage(got[0])
	require.NoError(t, err)
	require.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 42}, payload)
}

func TestClientWriteBAMReplaysThroughReadAPI(t *testing.T) {
	writer, err := NewClient(t.Context(), Replay(nil), WithSourceAddress(66))
	require.NoError(t, err)
	defer func() { _ = writer.Close() }()
	wire := []byte{1, 2, 3, 4, 5, 6, 7, 8, 42}
	message, err := pgn.DecodeMessage(pgn.MessageInfo{PGN: 65240}, wire)
	require.NoError(t, err)
	require.NoError(t, writer.Write(message).WaitContext(t.Context()))
	frames := writer.WrittenFrames()
	require.Len(t, frames, 3)
	require.Equal(t, byte(32), frames[0].Data[0], "writer must use BAM for Commanded Address")
	scanner := NewScanner(t.Context(), Replay(frames), Filter("pgn == 65240"))
	defer func() { _ = scanner.Close() }()
	require.True(t, scanner.Next())
	got, err := pgn.EncodeMessage(scanner.Message())
	require.NoError(t, err)
	require.Equal(t, wire, got)
	require.Equal(t, uint8(66), scanner.Message().(pgn.PGN).MessageInfo().SourceId)
	require.False(t, scanner.Next())
	require.NoError(t, scanner.Err())
}

type cleanupSource struct{ closing, release, done chan struct{} }

func (s *cleanupSource) run(ctx context.Context, _ *slog.Logger, handler func(raw.Observation)) error {
	handler(raw.Observation{Kind: raw.KindMessage, PGN: 127250, Payload: []byte{0, 0x5c, 0x3d, 0xff, 0x7f, 0xff, 0x7f, 0xfc}})
	<-ctx.Done()
	close(s.closing)
	<-s.release
	close(s.done)
	return nil
}
func TestScannerCloseJoinsSource(t *testing.T) {
	src := &cleanupSource{make(chan struct{}), make(chan struct{}), make(chan struct{})}
	scanner := NewScanner(t.Context(), optionFunc(func(c *config) { c.sources = append(c.sources, src) }))
	require.True(t, scanner.Next())
	returned := make(chan struct{})
	go func() { _ = scanner.Close(); close(returned) }()
	<-src.closing
	select {
	case <-returned:
		t.Error("Close returned before source cleanup")
	case <-time.After(20 * time.Millisecond):
	}
	close(src.release)
	<-returned
	select {
	case <-src.done:
	default:
		t.Error("source still running after Close")
	}
	require.NoError(t, scanner.Close())
}
func TestScannerCloseBeforeNext(t *testing.T) {
	src := &cleanupSource{}
	scanner := NewScanner(t.Context(), optionFunc(func(c *config) { c.sources = append(c.sources, src) }))
	require.NoError(t, scanner.Close())
	require.False(t, scanner.Next())
	require.NoError(t, scanner.Err())
}

func TestObserveEarlyExitJoinsSource(t *testing.T) {
	src := &cleanupSource{make(chan struct{}), make(chan struct{}), make(chan struct{})}
	returned := make(chan struct{})
	go func() {
		defer close(returned)
		for range Observe(t.Context(), optionFunc(func(c *config) { c.sources = append(c.sources, src) })) {
			break
		}
	}()
	<-src.closing
	select {
	case <-returned:
		t.Error("iterator returned before source cleanup")
	case <-time.After(20 * time.Millisecond):
	}
	close(src.release)
	<-returned
	select {
	case <-src.done:
	default:
		t.Error("source still running after iteration ended")
	}
}

func TestReadAPIsRejectWithBus(t *testing.T) {
	// Any non-nil Bus is sufficient; rejection must happen before opening it.
	bus := &unusedReadBus{}
	scanner := NewScanner(t.Context(), WithBus(bus))
	defer func() { _ = scanner.Close() }()
	require.False(t, scanner.Next())
	require.ErrorContains(t, scanner.Err(), "WithBus requires NewClient")
	for _, iterator := range []func(func(Observation, error) bool){Observe(t.Context(), WithBus(bus))} {
		count := 0
		iterator(func(_ Observation, err error) bool {
			count++
			require.ErrorContains(t, err, "WithBus requires NewClient")
			return true
		})
		require.Equal(t, 1, count)
	}
	for _, err := range Receive(t.Context(), WithBus(bus)) {
		require.ErrorContains(t, err, "WithBus requires NewClient")
	}
}

type unusedReadBus struct{}

func (*unusedReadBus) Run(context.Context, func(can.Frame)) error { panic("unexpected Bus.Run") }
func (*unusedReadBus) WriteFrame(can.Frame) error                 { panic("unexpected Bus.WriteFrame") }
func (*unusedReadBus) Close() error                               { panic("unexpected Bus.Close") }
