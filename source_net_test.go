package n2k

import (
	"bufio"
	"context"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/open-ships/n2k/internal/framer"
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

func TestTCPSource_Reconnects(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })

	// Each accepted connection delivers one heading line, then drops — so a
	// second and third message can only arrive if the source reconnects.
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			_, _ = conn.Write([]byte(ydHeadingLine))
			_ = conn.Close()
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count := 0
	for msg, err := range Receive(ctx,
		TCP(ln.Addr().String(), FormatYDRaw),
		WithReconnect(ReconnectPolicy{InitialBackoff: 5 * time.Millisecond, MaxBackoff: 20 * time.Millisecond}),
		WithLogger(testLogger()),
	) {
		require.NoError(t, err)
		_, ok := msg.(*pgn.VesselHeading)
		require.True(t, ok, "expected *pgn.VesselHeading, got %T", msg)
		count++
		if count == 3 {
			break
		}
	}
	assert.Equal(t, 3, count, "reconnect should keep delivering messages across drops")
}

func TestTCPClientReconnectReclaimsBeforeRestartingProtocolTraffic(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = listener.Close() })

	type observed struct {
		connection int
		pgn        uint32
	}
	observations := make(chan observed, 32)
	go func() {
		for connection := 1; connection <= 2; connection++ {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				return
			}
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				fields := strings.Fields(scanner.Text())
				if len(fields) == 0 {
					continue
				}
				id, parseErr := strconv.ParseUint(fields[0], 16, 32)
				if parseErr != nil {
					continue
				}
				pgnNum := framer.ParseCANID(uint32(id)).PGN
				observations <- observed{connection: connection, pgn: pgnNum}
				if connection == 1 && pgnNum == framer.PGNISOAddressClaim {
					// Let the initial contention window finish, then simulate a
					// gateway restart after NewClient has become active.
					time.Sleep(100 * time.Millisecond)
					_ = conn.Close()
					break
				}
			}
			if connection == 2 {
				_ = conn.Close()
			}
		}
	}()

	client, err := NewClient(context.Background(),
		TCP(listener.Addr().String(), FormatYDRaw),
		WithReconnect(ReconnectPolicy{InitialBackoff: 5 * time.Millisecond, MaxBackoff: 20 * time.Millisecond}),
		WithClaimTimeout(50*time.Millisecond),
		WithHeartbeatInterval(0),
		WithLogger(testLogger()),
	)
	require.NoError(t, err)
	defer func() { _ = client.Close() }()

	seenReclaim := false
	seenEnumeration := false
	deadline := time.After(3 * time.Second)
	for !seenReclaim || !seenEnumeration {
		select {
		case got := <-observations:
			if got.connection != 2 {
				continue
			}
			switch got.pgn {
			case framer.PGNISOAddressClaim:
				seenReclaim = true
			case framer.PGNISORequest:
				seenEnumeration = true
				require.True(t, seenReclaim, "protocol traffic resumed before reconnect claim")
			}
		case <-deadline:
			t.Fatalf("reconnect lifecycle incomplete: reclaim=%v enumeration=%v status=%+v", seenReclaim, seenEnumeration, client.Status())
		}
	}
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

// --- End-to-end write path: a real NewClient over a fake TCP gateway ---

// fakeYDGateway accepts one connection, forwards received RAW lines to the
// lines channel, and echoes each accepted transmit line back with a T
// direction (as a YDWG-02 does), which lets the client observe its own
// address claim.
func fakeYDGateway(t *testing.T, lines chan<- string) string {
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
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			line := scanner.Text()
			select {
			case lines <- line:
			default:
			}
			// Echo back with a timestamp and T direction.
			_, _ = conn.Write([]byte("00:00:00.000 T " + line + "\r\n"))
		}
	}()
	return ln.Addr().String()
}

func TestNewClient_TCPYDRaw_ClaimsAndWrites(t *testing.T) {
	lines := make(chan string, 64)
	addr := fakeYDGateway(t, lines)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := NewClient(ctx,
		TCP(addr, FormatYDRaw),
		WithClaimTimeout(500*time.Millisecond),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	// The client must have transmitted an address claim (PGN 60928) as a
	// RAW line during NewClient.
	sawClaim := false
	deadline := time.After(2 * time.Second)
drainClaim:
	for {
		select {
		case line := <-lines:
			if frame, ok := gatewayParseForTest(line); ok {
				rawPGN := (frame & 0x3FFFF00) >> 8
				if uint8((rawPGN>>8)&0xFF) < 240 {
					rawPGN &= 0xFFF00
				}
				if rawPGN == 60928 {
					sawClaim = true
					break drainClaim
				}
			}
		case <-deadline:
			break drainClaim
		}
	}
	require.True(t, sawClaim, "client should have transmitted an address claim over the gateway")

	// A user write reaches the gateway as a RAW line for its PGN.
	heading := &pgn.VesselHeading{}
	heading.SetHeadingValue(1.5708)
	require.NoError(t, c.Write(heading).Wait())

	sawHeading := false
	deadline = time.After(2 * time.Second)
drainHeading:
	for {
		select {
		case line := <-lines:
			if frame, ok := gatewayParseForTest(line); ok {
				if (frame&0x3FFFF00)>>8 == 127250 {
					sawHeading = true
					break drainHeading
				}
			}
		case <-deadline:
			break drainHeading
		}
	}
	require.True(t, sawHeading, "user write should reach the gateway as a RAW line")
}

// gatewayParseForTest extracts the CAN ID from an application-to-gateway RAW
// line (no time/direction prefix).
func gatewayParseForTest(line string) (uint32, bool) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0, false
	}
	id, err := strconv.ParseUint(fields[0], 16, 32)
	if err != nil {
		return 0, false
	}
	return uint32(id), true
}

func TestNewClient_TCPActisense_Writes(t *testing.T) {
	sends := make(chan []byte, 16)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer func() { _ = conn.Close() }()
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n > 0 {
				chunk := make([]byte, n)
				copy(chunk, buf[:n])
				select {
				case sends <- chunk:
				default:
				}
			}
			if err != nil {
				return
			}
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Actisense sends carry no source address, so address claiming cannot
	// arbitrate; use an explicit address so NewClient returns promptly.
	c, err := NewClient(ctx,
		TCP(ln.Addr().String(), FormatActisense),
		WithSourceAddress(42),
		WithClaimTimeout(500*time.Millisecond),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	heading := &pgn.VesselHeading{}
	heading.SetHeadingValue(1.5708)
	require.NoError(t, c.Write(heading).Wait())

	// Collect stream bytes until a 0x94 send command for PGN 127250 appears.
	var wire []byte
	deadline := time.After(2 * time.Second)
	found := false
	for !found {
		select {
		case chunk := <-sends:
			wire = append(wire, chunk...)
			found = containsActisenseSend(wire, 127250)
		case <-deadline:
			t.Fatalf("no Actisense send command for PGN 127250 seen; got %d bytes", len(wire))
		}
	}
}

// containsActisenseSend reports whether buf contains a framed 0x94 send
// command whose payload targets pgnNum. It is a lenient scan sufficient for
// the test: it looks for the DLE STX 0x94 header and decodes the PGN.
func containsActisenseSend(buf []byte, pgnNum uint32) bool {
	for i := 0; i+6 < len(buf); i++ {
		if buf[i] == 0x10 && buf[i+1] == 0x02 && buf[i+2] == 0x94 {
			// payload begins after cmd(1)+len(1): prio(1) then PGN(3 LE).
			p := i + 5
			if p+3 >= len(buf) {
				return false
			}
			got := uint32(buf[p]) | uint32(buf[p+1])<<8 | uint32(buf[p+2])<<16
			if got == pgnNum {
				return true
			}
		}
	}
	return false
}
