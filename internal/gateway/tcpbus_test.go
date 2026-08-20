package gateway

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

func respondToActisenseBEM(t *testing.T, conn net.Conn, parser *actisense.Parser, chunk []byte, mode *actisense.OperatingMode) {
	t.Helper()
	parser.Feed(chunk, func(datagram actisense.Datagram) {
		if datagram.ID != actisense.BSTBEMCommand || len(datagram.Payload) == 0 {
			return
		}
		command := datagram.Payload[0]
		var data []byte
		switch command {
		case actisense.BEMOperatingMode:
			if len(datagram.Payload) >= 3 {
				*mode = actisense.OperatingMode(binary.LittleEndian.Uint16(datagram.Payload[1:3]))
			}
			data = actisense.OperatingModeSet(*mode)
		case actisense.BEMTxPGNEnable:
			data = make([]byte, 14)
			if len(datagram.Payload) >= 5 {
				copy(data[:4], datagram.Payload[1:5])
			}
			data[4] = 1
		case actisense.BEMActivatePGNLists:
		default:
			return
		}
		response := make([]byte, 12, 12+len(data))
		response[0] = command
		response = append(response, data...)
		wire, err := actisense.EncodeDatagram(actisense.BSTBEMResponse, response)
		if err != nil {
			t.Errorf("encode fake gateway response: %v", err)
			return
		}
		if _, err = conn.Write(wire); err != nil {
			t.Errorf("write fake gateway response: %v", err)
		}
	}, func(err actisense.DecodeError) { t.Errorf("fake gateway decode: %v", err) })
}

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
	bus := NewYDRawTCPBus(testLogger(), addr, nil)
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
	bus := NewYDRawTCPBus(testLogger(), addr, nil)
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
	bus := NewYDRawTCPBus(testLogger(), "127.0.0.1:1", nil) // never dialed

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

// TestYDRawTCPBus_CancelReleasesWriterDuringOutage confirms a write parked
// during a reconnect outage is released (with an error) when the context is
// cancelled — await must not block forever on a reconnection that never comes.
func TestYDRawTCPBus_CancelReleasesWriterDuringOutage(t *testing.T) {
	// Nothing listens, so with reconnect enabled the writer parks in await
	// waiting for a connection that never arrives.
	bus := NewYDRawTCPBus(testLogger(), "127.0.0.1:1", &ReconnectPolicy{
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     20 * time.Millisecond,
	})
	defer func() { _ = bus.Close() }()

	ctx, cancel := context.WithCancel(context.Background())
	runErr := make(chan error, 1)
	go func() { runErr <- bus.Run(ctx, nil) }()

	writeErr := make(chan error, 1)
	go func() { writeErr <- bus.WriteFrame(can.Frame{ID: 1, Length: 1}) }()

	time.Sleep(50 * time.Millisecond) // let the writer park during the outage
	cancel()

	select {
	case err := <-writeErr:
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(2 * time.Second):
		t.Fatal("cancel did not release the writer parked during a reconnect outage")
	}
	select {
	case err := <-runErr:
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after cancellation")
	}
}

// TestYDRawTCPBus_CloseReleasesWriterDuringOutage confirms Close releases a
// write parked during a reconnect outage while Run is active.
func TestYDRawTCPBus_CloseReleasesWriterDuringOutage(t *testing.T) {
	bus := NewYDRawTCPBus(testLogger(), "127.0.0.1:1", &ReconnectPolicy{
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     20 * time.Millisecond,
	})
	runErr := make(chan error, 1)
	go func() { runErr <- bus.Run(context.Background(), nil) }()

	writeErr := make(chan error, 1)
	go func() { writeErr <- bus.WriteFrame(can.Frame{ID: 1, Length: 1}) }()

	time.Sleep(50 * time.Millisecond) // let the writer park during the outage
	require.NoError(t, bus.Close())

	select {
	case err := <-writeErr:
		require.Error(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Close did not release the writer parked during a reconnect outage")
	}
	select {
	case err := <-runErr:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after Close")
	}
}

// TestYDRawTCPBus_Reconnects drops the connection after each frame; with a
// reconnect policy the bus re-dials and keeps delivering frames from the
// fresh connection.
func TestYDRawTCPBus_Reconnects(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })

	frameLines := []string{
		"17:33:21.141 R 09F80115 A0 7D E6 18 C0 05 FB D5\r\n",
		"17:33:22.141 R 09F80116 B0 7D E6 18 C0 05 FB D5\r\n",
	}
	go func() {
		for i := 0; ; i++ {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			if i < len(frameLines) {
				_, _ = conn.Write([]byte(frameLines[i]))
			}
			// Drop the connection to force a reconnect.
			_ = conn.Close()
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bus := NewYDRawTCPBus(testLogger(), ln.Addr().String(), &ReconnectPolicy{
		InitialBackoff: 5 * time.Millisecond,
		MaxBackoff:     20 * time.Millisecond,
	})
	defer func() { _ = bus.Close() }()

	frames := make(chan can.Frame, 4)
	runErr := make(chan error, 1)
	go func() { runErr <- bus.Run(ctx, func(f can.Frame) { frames <- f }) }()

	// The first connection delivers frame A.
	select {
	case f := <-frames:
		assert.Equal(t, uint32(0x09F80115), f.ID)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for first frame")
	}
	// After the drop, the bus reconnects and the second connection delivers B.
	select {
	case f := <-frames:
		assert.Equal(t, uint32(0x09F80116), f.ID)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for frame after reconnect")
	}

	cancel()
	select {
	case err := <-runErr:
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after cancellation")
	}
}

// TestYDRawTCPBus_CloseStopsReconnect confirms Close interrupts a pending
// reconnect backoff instead of waiting for it to elapse, even when the
// caller's context is never cancelled.
func TestYDRawTCPBus_CloseStopsReconnect(t *testing.T) {
	// Nothing listens here, so the first dial fails and the bus enters a long
	// backoff that only Close should be able to cut short.
	bus := NewYDRawTCPBus(testLogger(), "127.0.0.1:1", &ReconnectPolicy{
		InitialBackoff: 30 * time.Second,
		MaxBackoff:     30 * time.Second,
	})

	runErr := make(chan error, 1)
	go func() { runErr <- bus.Run(context.Background(), nil) }()

	time.Sleep(50 * time.Millisecond) // let the dial fail and enter backoff
	require.NoError(t, bus.Close())

	select {
	case err := <-runErr:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Close did not interrupt the reconnect backoff")
	}
}

// TestYDRawTCPBus_NoReconnectEndsOnDrop confirms the default (nil policy)
// still ends Run when the connection drops.
func TestYDRawTCPBus_NoReconnectEndsOnDrop(t *testing.T) {
	addr := startFakeGateway(t, func(conn net.Conn) {
		_, _ = conn.Write([]byte("17:33:21.141 R 09F80115 A0 7D E6 18 C0 05 FB D5\r\n"))
		// serve returns, closing the connection and dropping the client.
	})

	bus := NewYDRawTCPBus(testLogger(), addr, nil)
	defer func() { _ = bus.Close() }()

	runErr := make(chan error, 1)
	go func() { runErr <- bus.Run(context.Background(), nil) }()

	select {
	case err := <-runErr:
		assert.NoError(t, err) // clean EOF ends Run without reconnecting
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after the connection dropped")
	}
}

// scriptedConn fakes just enough of net.Conn for the write-retry tests: only
// Write is exercised, so the embedded nil net.Conn is never dereferenced.
type failingConn struct{ net.Conn }

func (*failingConn) Write([]byte) (int, error) {
	return 0, errors.New("gateway_test: connection reset")
}

type partialConn struct{ net.Conn }

func (*partialConn) Write(p []byte) (int, error) {
	return 1, errors.New("gateway_test: connection reset")
}

type recordingConn struct {
	net.Conn
	mu  sync.Mutex
	buf []byte
}

func (c *recordingConn) Write(p []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.buf = append(c.buf, p...)
	return len(p), nil
}

func (c *recordingConn) written() []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]byte(nil), c.buf...)
}

// TestTCPLink_WriteRetriesFreshConnAfterDrop confirms a write whose connection
// dropped before any byte was sent (n == 0) waits for the reconnect and
// re-sends the whole frame on the fresh connection.
func TestTCPLink_WriteRetriesFreshConnAfterDrop(t *testing.T) {
	l := newTCPLink(testLogger(), "127.0.0.1:1", &ReconnectPolicy{
		InitialBackoff: time.Millisecond,
		MaxBackoff:     time.Millisecond,
	})
	dead := &failingConn{}
	good := &recordingConn{}

	// Publish a connection whose first write fails with nothing sent.
	l.mu.Lock()
	l.conn = dead
	l.mu.Unlock()

	writeErr := make(chan error, 1)
	go func() { writeErr <- l.write([]byte("hello")) }()

	// The write fails on dead (n == 0) and waits for a different connection.
	time.Sleep(20 * time.Millisecond)
	require.True(t, l.markConnected(good))

	select {
	case err := <-writeErr:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("write did not complete after reconnect")
	}
	assert.Equal(t, []byte("hello"), good.written())
}

// TestTCPLink_WriteNoRetryOnPartialWrite confirms a partial write is never
// retried — a retry could duplicate the bytes already on the wire — so write
// fails promptly instead of awaiting a reconnect.
func TestTCPLink_WriteNoRetryOnPartialWrite(t *testing.T) {
	l := newTCPLink(testLogger(), "127.0.0.1:1", &ReconnectPolicy{
		InitialBackoff: time.Millisecond,
		MaxBackoff:     time.Millisecond,
	})
	l.mu.Lock()
	l.conn = &partialConn{}
	l.mu.Unlock()

	writeErr := make(chan error, 1)
	go func() { writeErr <- l.write([]byte("hello")) }()
	select {
	case err := <-writeErr:
		require.Error(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("write blocked retrying a partial write instead of failing")
	}
}

// TestTCPLink_WriteNoRetryWithoutReconnect confirms that without a reconnect
// policy a failed write returns the error rather than retrying.
func TestTCPLink_WriteNoRetryWithoutReconnect(t *testing.T) {
	l := newTCPLink(testLogger(), "127.0.0.1:1", nil) // no reconnect
	l.mu.Lock()
	l.conn = &failingConn{}
	l.mu.Unlock()

	require.Error(t, l.write([]byte("hello")))
}

func TestTCPLinkConnectionObserverGatesPublicationAndTracksEpochs(t *testing.T) {
	link := newTCPLink(testLogger(), "127.0.0.1:1", &ReconnectPolicy{})
	type event struct {
		connected bool
		epoch     uint64
		published bool
	}
	var events []event
	link.setConnectionObserver(func(connected bool, epoch uint64) {
		link.mu.Lock()
		published := link.conn != nil
		link.mu.Unlock()
		events = append(events, event{connected: connected, epoch: epoch, published: published})
	})

	first := &recordingConn{}
	second := &recordingConn{}
	require.True(t, link.markConnected(first))
	link.markDisconnected(errors.New("drop"))
	require.True(t, link.markConnected(second))

	require.Equal(t, []event{
		{connected: true, epoch: 1, published: false},
		{connected: false, epoch: 1, published: false},
		{connected: true, epoch: 2, published: false},
	}, events)
}

func TestActisenseTCPBus_StartupReadAndWrite(t *testing.T) {
	received := make(chan []byte, 4)
	addr := startFakeGateway(t, func(conn net.Conn) {
		// Feed one assembled message to the client, then capture its writes.
		payload := n2kPayload(2, 127250, 255, 23, []byte{0xFF, 0x88, 0x4C, 0xFF, 0x7F, 0xFF, 0x7F, 0xFD})
		_, _ = conn.Write(wrapActisense(cmdN2KReceived, payload))
		parser := actisense.NewParser()
		mode := actisense.ModeTransferNormal
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n > 0 {
				chunk := make([]byte, n)
				copy(chunk, buf[:n])
				received <- chunk
				respondToActisenseBEM(t, conn, parser, chunk, &mode)
			}
			if err != nil {
				return
			}
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bus := NewActisenseTCPBus(testLogger(), addr, nil)
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

	// First bytes on the wire query and then acknowledge the requested mode.
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
	for !containsBEMCommand(wire, actisense.BEMOperatingMode, true) {
		wire = append(wire, collect()...)
	}

	// A fast-packet-sized payload goes out as one whole send command.
	data := make([]byte, 20)
	for i := range data {
		data[i] = byte(i)
	}
	require.NoError(t, bus.WriteMessage(126996, 6, 42, 255, data))
	want, err := EncodeSend(N2KMessage{Priority: 6, PGN: 126996, Destination: 255, Data: data})
	require.NoError(t, err)
	for !bytes.Contains(wire, want) {
		wire = append(wire, collect()...)
	}
	assert.True(t, bytes.Contains(wire, want))
}

func TestActisenseTCPBus_WriteFrameAsMessage(t *testing.T) {
	received := make(chan []byte, 4)
	addr := startFakeGateway(t, func(conn net.Conn) {
		parser := actisense.NewParser()
		mode := actisense.ModeTransferNormal
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n > 0 {
				chunk := make([]byte, n)
				copy(chunk, buf[:n])
				received <- chunk
				respondToActisenseBEM(t, conn, parser, chunk, &mode)
			}
			if err != nil {
				return
			}
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bus := NewActisenseTCPBus(testLogger(), addr, nil)
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
	for !bytes.Contains(wire, want) {
		select {
		case chunk := <-received:
			wire = append(wire, chunk...)
		case <-deadline:
			t.Fatal("timed out waiting for client bytes")
		}
	}
	assert.True(t, bytes.Contains(wire, want))
}

func containsBEMCommand(wire []byte, command byte, withData bool) bool {
	found := false
	actisense.NewParser().Feed(wire, func(datagram actisense.Datagram) {
		if datagram.ID == actisense.BSTBEMCommand && len(datagram.Payload) > 0 && datagram.Payload[0] == command {
			found = !withData || len(datagram.Payload) > 1
		}
	}, nil)
	return found
}
