package n2k

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Actisense test framing helpers (mirror internal/gateway wire format) ---

const (
	testDLE            = 0x10
	testSTX            = 0x02
	testETX            = 0x03
	testCmdN2KReceived = 0x93
)

// wrapActisense builds a DLE/STX-framed Actisense message around cmd+payload.
func wrapActisense(cmd byte, payload []byte) []byte {
	body := append([]byte{cmd, byte(len(payload))}, payload...)
	var sum byte
	for _, b := range body {
		sum += b
	}
	body = append(body, -sum)

	out := []byte{testDLE, testSTX}
	for _, b := range body {
		if b == testDLE {
			out = append(out, testDLE, testDLE)
		} else {
			out = append(out, b)
		}
	}
	return append(out, testDLE, testETX)
}

// actisenseN2K wraps an assembled PGN payload as an Actisense N2K message.
func actisenseN2K(prio uint8, pgnNum uint32, dst, src uint8, data []byte) []byte {
	p := []byte{
		prio,
		byte(pgnNum), byte(pgnNum >> 8), byte(pgnNum >> 16),
		dst, src,
		0, 0, 0, 0,
		byte(len(data)),
	}
	return wrapActisense(testCmdN2KReceived, append(p, data...))
}

const ydHeadingLine = "17:33:21.107 R 09F11201 01 5C 3D FF 7F FF 7F FC\r\n"

func collectMessages(t *testing.T, ctx context.Context, opt Option, want int) []pgn.Message {
	t.Helper()
	var msgs []pgn.Message
	for msg, err := range Receive(ctx, opt) {
		require.NoError(t, err)
		msgs = append(msgs, msg)
		if len(msgs) == want {
			break
		}
	}
	return msgs
}

func TestTCPSource_YDRaw(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() { _ = ln.Close() }()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		_, _ = conn.Write([]byte("YDWG-02 service line\r\n" + ydHeadingLine))
		_ = conn.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msgs := collectMessages(t, ctx, TCP(ln.Addr().String(), FormatYDRaw), 1)
	require.Len(t, msgs, 1)
	vh, ok := msgs[0].(*pgn.VesselHeading)
	require.True(t, ok, "expected *pgn.VesselHeading, got %T", msgs[0])
	require.NotNil(t, vh.Heading)
	assert.Equal(t, uint64(15708), *vh.Heading)
}

func TestTCPSource_ActisenseFastPacket(t *testing.T) {
	// A fast-packet ProductInformation crosses re-framing and reassembly.
	product := &pgn.ProductInformation{
		ProductCode: u64(1234),
		ModelId:     "test gateway",
	}
	payload, err := pgn.EncodeMessage(product)
	require.NoError(t, err)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() { _ = ln.Close() }()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		_, _ = conn.Write(actisenseN2K(6, 126996, 255, 0x42, payload))
		_ = conn.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msgs := collectMessages(t, ctx, TCP(ln.Addr().String(), FormatActisense), 1)
	require.Len(t, msgs, 1)
	pi, ok := msgs[0].(*pgn.ProductInformation)
	require.True(t, ok, "expected *pgn.ProductInformation, got %T", msgs[0])
	require.NotNil(t, pi.ProductCode)
	assert.Equal(t, uint64(1234), *pi.ProductCode)
	assert.Equal(t, "test gateway", pi.ModelId)
	assert.Equal(t, uint8(0x42), pi.Info.SourceId)
}

func TestTCPSource_DialFailure(t *testing.T) {
	// A port nobody listens on: grab one, then close it.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := ln.Addr().String()
	require.NoError(t, ln.Close())

	var firstErr error
	for _, err := range Receive(context.Background(), TCP(addr, FormatYDRaw)) {
		if err != nil {
			firstErr = err
			break
		}
	}
	require.Error(t, firstErr)
}

func TestTCPSource_ContextCancel(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() { _ = ln.Close() }()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		// Hold the connection open, sending nothing.
		time.Sleep(10 * time.Second)
		_ = conn.Close()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	for range Receive(ctx, TCP(ln.Addr().String(), FormatYDRaw)) {
	}
	assert.Less(t, time.Since(start), 5*time.Second, "cancel should terminate a blocked read promptly")
}

func TestUDPSource_YDRaw(t *testing.T) {
	// Reserve a local UDP address for the source to listen on.
	probe, err := net.ListenPacket("udp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := probe.LocalAddr().String()
	require.NoError(t, probe.Close())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		conn, err := net.Dial("udp", addr)
		if err != nil {
			return
		}
		defer func() { _ = conn.Close() }()
		// Retry until the listener is up (or the test times out).
		for range 50 {
			_, _ = conn.Write([]byte(ydHeadingLine))
			time.Sleep(20 * time.Millisecond)
		}
	}()

	msgs := collectMessages(t, ctx, UDP(addr, FormatYDRaw), 1)
	require.Len(t, msgs, 1)
	_, ok := msgs[0].(*pgn.VesselHeading)
	require.True(t, ok, "expected *pgn.VesselHeading, got %T", msgs[0])
}

func u64(v uint64) *uint64 { return &v }
