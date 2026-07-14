package gateway

import (
	"bufio"
	"context"
	"io"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

// startFakeGateway listens on loopback and hands the accepted connection to
// serve on a goroutine.
func startFakeGateway(t *testing.T, serve func(net.Conn)) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer func() { _ = conn.Close() }()
		serve(conn)
	}()
	return ln.Addr().String()
}

func TestYDRawTCPBus_ReadAndWrite(t *testing.T) {
	lines := make(chan string, 4)
	addr := startFakeGateway(t, func(conn net.Conn) {
		// Deliver one frame to the client, then record what it transmits.
		_, _ = conn.Write([]byte("17:33:21.141 R 09F80115 A0 7D E6 18 C0 05 FB D5\r\n"))
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bus := NewYDRawTCPBus(testLogger(), addr)
	defer func() { _ = bus.Close() }()

	frames := make(chan can.Frame, 4)
	runErr := make(chan error, 1)
	go func() { runErr <- bus.Run(ctx, func(f can.Frame) { frames <- f }) }()

	select {
	case f := <-frames:
		assert.Equal(t, uint32(0x09F80115), f.ID)
		assert.Equal(t, uint8(0xA0), f.Data[0])
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for received frame")
	}

	err := bus.WriteFrame(can.Frame{ID: 0x19F51323, Length: 2, Data: [8]uint8{0x01, 0x02}})
	require.NoError(t, err)

	select {
	case line := <-lines:
		assert.Equal(t, "19F51323 01 02", line)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for transmitted line")
	}

	cancel()
	select {
	case err := <-runErr:
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after cancellation")
	}
}

// TestYDRawTCPBus_WriteBeforeConnect exercises the await path: writes issued
// before Run has dialed must block until the connection is up, then succeed.
func TestYDRawTCPBus_WriteBeforeConnect(t *testing.T) {
	lines := make(chan string, 1)
	addr := startFakeGateway(t, func(conn net.Conn) {
		scanner := bufio.NewScanner(conn)
		if scanner.Scan() {
			lines <- scanner.Text()
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bus := NewYDRawTCPBus(testLogger(), addr)
	defer func() { _ = bus.Close() }()

	writeErr := make(chan error, 1)
	go func() {
		writeErr <- bus.WriteFrame(can.Frame{ID: 0x18EEFF01, Length: 8, Data: [8]uint8{1, 2, 3, 4, 5, 6, 7, 8}})
	}()

	// Give the writer a moment to block, then start Run.
	time.Sleep(50 * time.Millisecond)
	go func() { _ = bus.Run(ctx, nil) }()

	select {
	case err := <-writeErr:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("write did not complete after connect")
	}
	select {
	case line := <-lines:
		assert.Equal(t, "18EEFF01 01 02 03 04 05 06 07 08", line)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for transmitted line")
	}
}

func TestYDRawTCPBus_CloseUnblocksWriter(t *testing.T) {
	bus := NewYDRawTCPBus(testLogger(), "127.0.0.1:1") // never dialed

	writeErr := make(chan error, 1)
	go func() { writeErr <- bus.WriteFrame(can.Frame{ID: 1, Length: 1}) }()
	time.Sleep(50 * time.Millisecond)
	require.NoError(t, bus.Close())

	select {
	case err := <-writeErr:
		require.Error(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("write did not unblock after Close")
	}
}

func TestActisenseTCPBus_StartupReadAndWrite(t *testing.T) {
	received := make(chan []byte, 4)
	addr := startFakeGateway(t, func(conn net.Conn) {
		// Feed one assembled message to the client, then capture its writes.
		payload := n2kPayload(2, 127250, 255, 23, []byte{0xFF, 0x88, 0x4C, 0xFF, 0x7F, 0xFF, 0x7F, 0xFD})
		_, _ = conn.Write(wrapActisense(cmdN2KReceived, payload))
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n > 0 {
				chunk := make([]byte, n)
				copy(chunk, buf[:n])
				received <- chunk
			}
			if err != nil {
				return
			}
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bus := NewActisenseTCPBus(testLogger(), addr)
	defer func() { _ = bus.Close() }()

	frames := make(chan can.Frame, 4)
	go func() { _ = bus.Run(ctx, func(f can.Frame) { frames <- f }) }()

	// The re-framed single-frame PGN should arrive with the gateway's
	// source address in the CAN ID.
	select {
	case f := <-frames:
		assert.Equal(t, uint8(23), uint8(f.ID&0xFF), "source address")
		assert.Equal(t, uint8(0xFF), f.Data[0])
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for re-framed message")
	}

	// First bytes on the wire must be the startup command.
	collect := func() []byte {
		select {
		case chunk := <-received:
			return chunk
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for client bytes")
			return nil
		}
	}
	wire := collect()
	require.GreaterOrEqual(t, len(wire), len(EncodeStartup()))
	assert.Equal(t, EncodeStartup(), wire[:len(EncodeStartup())])
	wire = wire[len(EncodeStartup()):]

	// A fast-packet-sized payload goes out as one whole send command.
	data := make([]byte, 20)
	for i := range data {
		data[i] = byte(i)
	}
	require.NoError(t, bus.WriteMessage(126996, 6, 42, 255, data))
	want, err := EncodeSend(N2KMessage{Priority: 6, PGN: 126996, Destination: 255, Data: data})
	require.NoError(t, err)
	for len(wire) < len(want) {
		wire = append(wire, collect()...)
	}
	assert.Equal(t, want, wire)
}

func TestActisenseTCPBus_WriteFrameAsMessage(t *testing.T) {
	received := make(chan []byte, 4)
	addr := startFakeGateway(t, func(conn net.Conn) {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n > 0 {
				chunk := make([]byte, n)
				copy(chunk, buf[:n])
				received <- chunk
			}
			if err != nil {
				return
			}
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bus := NewActisenseTCPBus(testLogger(), addr)
	defer func() { _ = bus.Close() }()
	go func() { _ = bus.Run(ctx, nil) }()

	// An address-claim frame (PDU1, destination in the CAN ID) becomes a
	// whole message send.
	claimID := framer.BuildCANID(60928, 6, 42, 255)
	frame := can.Frame{ID: claimID, Length: 8, Data: [8]uint8{1, 2, 3, 4, 5, 6, 7, 8}}
	require.NoError(t, bus.WriteFrame(frame))

	want, err := EncodeSend(N2KMessage{Priority: 6, PGN: 60928, Destination: 255, Data: frame.Data[:]})
	require.NoError(t, err)

	var wire []byte
	deadline := time.After(2 * time.Second)
	for len(wire) < len(EncodeStartup())+len(want) {
		select {
		case chunk := <-received:
			wire = append(wire, chunk...)
		case <-deadline:
			t.Fatal("timed out waiting for client bytes")
		}
	}
	assert.Equal(t, want, wire[len(EncodeStartup()):])
}
