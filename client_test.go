package n2k

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_NewClient_Validation(t *testing.T) {
	t.Run("no sources returns error", func(t *testing.T) {
		_, err := NewClient(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "at least one source")
	})

	t.Run("replay source succeeds", func(t *testing.T) {
		c, err := NewClient(context.Background(), Replay(nil))
		require.NoError(t, err)
		require.NotNil(t, c)
		assert.NoError(t, c.Close())
	})

	t.Run("with explicit source address", func(t *testing.T) {
		c, err := NewClient(context.Background(), Replay(nil), WithSourceAddress(42))
		require.NoError(t, err)
		assert.Equal(t, uint8(42), c.sourceAddr)
		assert.NoError(t, c.Close())
	})

	t.Run("with custom device name", func(t *testing.T) {
		name := DeviceName{
			IndustryGroup:    4,
			ManufacturerCode: 1234,
			DeviceClass:      25,
			DeviceFunction:   130,
			IdentityNumber:   99999,
		}
		c, err := NewClient(context.Background(), Replay(nil), WithName(name))
		require.NoError(t, err)
		assert.Equal(t, name.Pack(true), c.deviceName)
		assert.NoError(t, c.Close())
	})
}

func TestNewClient_Replay_BadFilterFailsEagerly(t *testing.T) {
	_, err := NewClient(context.Background(), Replay(nil), Filter("((("))
	if err == nil {
		t.Fatal("expected eager filter compile error for replay client")
	}
}

func TestClient_Write_EncodesAndFrames(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	heading := uint64(15000)
	msg := &pgn.VesselHeading{
		Heading: &heading,
	}
	msg.Info.Priority = ptrUint8(2)

	wr := c.Write(msg)
	err = wr.Wait()
	require.NoError(t, err)

	frames := c.WrittenFrames()
	require.Len(t, frames, 1, "VesselHeading is a single-frame PGN")

	// Verify the CAN ID encodes the correct PGN, priority, and source.
	expectedID := framer.BuildCANID(127250, 2, 0, 255)
	assert.Equal(t, expectedID, frames[0].ID)
}

func TestClient_Write_UsesPriorityFromInfo(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil), WithSourceAddress(10))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	heading := uint64(5000)
	msg := &pgn.VesselHeading{
		Heading: &heading,
	}
	msg.Info.Priority = ptrUint8(3)

	wr := c.Write(msg)
	require.NoError(t, wr.Wait())

	frames := c.WrittenFrames()
	require.Len(t, frames, 1)

	expectedID := framer.BuildCANID(127250, 3, 10, 255)
	assert.Equal(t, expectedID, frames[0].ID)
}

func TestClient_Write_DefaultPriority(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	heading := uint64(5000)
	msg := &pgn.VesselHeading{
		Heading: &heading,
		// Priority left nil => default to 6
	}

	wr := c.Write(msg)
	require.NoError(t, wr.Wait())

	frames := c.WrittenFrames()
	require.Len(t, frames, 1)

	expectedID := framer.BuildCANID(127250, 6, 0, 255)
	assert.Equal(t, expectedID, frames[0].ID)
}

func TestClient_Write_NoEncoder(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	// UnknownPGN is a decoded fallback, not a serializable PGN struct.
	msg := &pgn.UnknownPGN{}
	msg.Info.PGN = 999999

	wr := c.Write(msg)
	err = wr.Wait()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not implement pgn.PGN")
}

func TestClient_Write_AfterClose(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil))
	require.NoError(t, err)

	require.NoError(t, c.Close())

	heading := uint64(10000)
	msg := &pgn.VesselHeading{
		Heading: &heading,
	}
	wr := c.Write(msg)
	err = wr.Wait()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "closed")
}

func TestClient_Replay_Receive(t *testing.T) {
	// Build a pre-encoded VesselHeading frame for replay.
	heading := uint64(15000)
	msg := &pgn.VesselHeading{
		Heading: &heading,
	}
	payload, err := msg.EncodePayload()
	require.NoError(t, err)

	canID := framer.BuildCANID(127250, 2, 42, 255)
	frame := framer.FrameSingle(canID, payload)

	c, err := NewClient(context.Background(), Replay([]can.Frame{frame}))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	var received []pgn.Message
	for msg, err := range c.Receive() {
		require.NoError(t, err)
		received = append(received, msg)
	}
	require.Len(t, received, 1)

	vh, ok := received[0].(*pgn.VesselHeading)
	require.True(t, ok, "expected *pgn.VesselHeading, got %T", received[0])
	require.NotNil(t, vh.Heading)
	assert.Equal(t, uint64(15000), *vh.Heading)
}

func TestClient_Scanner(t *testing.T) {
	heading := uint64(20000)
	msg := &pgn.VesselHeading{
		Heading: &heading,
	}
	payload, err := msg.EncodePayload()
	require.NoError(t, err)

	canID := framer.BuildCANID(127250, 2, 42, 255)
	frame := framer.FrameSingle(canID, payload)

	c, err := NewClient(context.Background(), Replay([]can.Frame{frame}))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	s := c.Scanner()
	require.True(t, s.Next())
	vh, ok := s.Message().(*pgn.VesselHeading)
	require.True(t, ok, "expected *pgn.VesselHeading, got %T", s.Message())
	require.NotNil(t, vh.Heading)
	assert.Equal(t, uint64(20000), *vh.Heading)
	assert.False(t, s.Next())
	assert.NoError(t, s.Err())
}

func TestClient_Close_Idempotent(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil))
	require.NoError(t, err)

	// Close multiple times, no panic.
	assert.NoError(t, c.Close())
	assert.NoError(t, c.Close())
	assert.NoError(t, c.Close())
}

func TestClient_WrittenFrames_Empty(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	frames := c.WrittenFrames()
	assert.Empty(t, frames)
}

func TestClient_MultipleWrites(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil), WithSourceAddress(5))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	heading := uint64(1000)
	for i := 0; i < 3; i++ {
		msg := &pgn.VesselHeading{
			Heading: &heading,
		}
		msg.Info.Priority = ptrUint8(2)
		wr := c.Write(msg)
		require.NoError(t, wr.Wait())
	}

	frames := c.WrittenFrames()
	assert.Len(t, frames, 3)
}

func TestClient_Write_FastPacket(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil), WithSourceAddress(10))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	// ProductInformation (PGN 126996) is a fast-packet PGN. Its string fields
	// produce a payload larger than 8 bytes, exercising the fast-packet framing path.
	ver := uint64(200)
	code := uint64(1234)
	cert := uint64(1)
	load := uint64(1)
	msg := &pgn.ProductInformation{
		Nmea2000Version:     &ver,
		ProductCode:         &code,
		ModelId:             "TestModel",
		SoftwareVersionCode: "1.0",
		ModelVersion:        "v1",
		ModelSerialCode:     "SN001",
		CertificationLevel:  &cert,
		LoadEquivalency:     &load,
	}
	msg.Info.Priority = ptrUint8(6)

	wr := c.Write(msg)
	require.NoError(t, wr.Wait())

	frames := c.WrittenFrames()
	require.Greater(t, len(frames), 1, "fast-packet PGN should produce multiple frames")

	// Verify frame 0 has the fast-packet header structure:
	// Data[0] = (seqId<<5 | frameNum), Data[1] = total length
	f0 := frames[0]
	frameNum := f0.Data[0] & 0x1F
	assert.Equal(t, uint8(0), frameNum, "first frame should have frame number 0")
	totalLen := f0.Data[1]
	assert.Greater(t, totalLen, uint8(8), "total payload length should exceed 8 bytes")

	// All frames should share the same CAN ID.
	for i, f := range frames {
		assert.Equal(t, f0.ID, f.ID, "frame %d should have same CAN ID as frame 0", i)
	}
}

func TestClient_CANSource_NoLongerStubbed(t *testing.T) {
	// CAN("can0") should attempt to construct a bus client, not error with
	// "not yet implemented". On a machine without can0, it will fail with a
	// hardware error — but NOT "not yet implemented".
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := NewClient(ctx, CAN("can_nonexistent_99"))
	require.Error(t, err)
	assert.NotContains(t, err.Error(), "not yet implemented",
		"should no longer return the old stub error")
}

// --- Pointer helpers for building test PGN structs ---

func ptrUint8(v uint8) *uint8    { return &v }
func ptrUint64(v uint64) *uint64 { return &v }

// --- End-to-end integration tests (write -> read round-trip) ---

func TestClient_RoundTrip_SingleFrame(t *testing.T) {
	ctx := context.Background()

	// 1. Build a VesselHeading with known values.
	original := &pgn.VesselHeading{
		Sid:     ptrUint64(0),
		Heading: ptrUint64(15708),
	}
	original.Info.Priority = ptrUint8(2)

	// 2. Write via client.
	writer, err := NewClient(ctx, Replay(nil))
	require.NoError(t, err)
	defer func() { _ = writer.Close() }()

	wr := writer.Write(original)
	require.NoError(t, wr.Wait())

	// 3. Get captured frames.
	frames := writer.WrittenFrames()
	require.Len(t, frames, 1, "VesselHeading should produce exactly 1 frame")

	// 4. Feed frames into a new Receive to decode them.
	var decoded []pgn.Message
	for msg, err := range Receive(ctx, Replay(frames)) {
		require.NoError(t, err)
		decoded = append(decoded, msg)
	}
	require.Len(t, decoded, 1, "should decode exactly 1 message from the captured frames")

	// 5. Assert the decoded VesselHeading matches.
	vh, ok := decoded[0].(*pgn.VesselHeading)
	require.True(t, ok, "expected *pgn.VesselHeading, got %T", decoded[0])

	require.NotNil(t, vh.Heading, "Heading should not be nil")
	assert.Equal(t, uint64(15708), *vh.Heading, "Heading should round-trip")

	require.NotNil(t, vh.Sid, "Sid should not be nil")
	assert.Equal(t, uint64(0), *vh.Sid, "Sid should match")
}

func TestClient_RoundTrip_FastPacket(t *testing.T) {
	ctx := context.Background()

	// 1. Build a ProductInformation with known values.
	original := &pgn.ProductInformation{
		Nmea2000Version:     ptrUint64(200),
		ProductCode:         ptrUint64(42),
		ModelId:             "TestModel",
		SoftwareVersionCode: "1.0.0",
		ModelVersion:        "v2",
		ModelSerialCode:     "SN12345",
		CertificationLevel:  ptrUint64(1),
		LoadEquivalency:     ptrUint64(1),
	}
	original.Info.Priority = ptrUint8(6)

	// 2. Write via client -- should produce multiple frames (fast-packet).
	writer, err := NewClient(ctx, Replay(nil))
	require.NoError(t, err)
	defer func() { _ = writer.Close() }()

	wr := writer.Write(original)
	require.NoError(t, wr.Wait())

	// 3. Get captured frames, verify more than 1 frame.
	frames := writer.WrittenFrames()
	require.Greater(t, len(frames), 1, "ProductInformation is fast-packet and must produce multiple frames")

	// 4. Feed frames into Receive to decode.
	var decoded []pgn.Message
	for msg, err := range Receive(ctx, Replay(frames)) {
		require.NoError(t, err)
		decoded = append(decoded, msg)
	}
	require.Len(t, decoded, 1, "should decode exactly 1 ProductInformation from the captured frames")

	// 5. Assert decoded ProductInformation matches.
	pi, ok := decoded[0].(*pgn.ProductInformation)
	require.True(t, ok, "expected *pgn.ProductInformation, got %T", decoded[0])

	require.NotNil(t, pi.ProductCode, "ProductCode should not be nil")
	assert.Equal(t, uint64(42), *pi.ProductCode, "ProductCode should match")

	require.NotNil(t, pi.Nmea2000Version, "Nmea2000Version should not be nil")
	assert.Equal(t, uint64(200), *pi.Nmea2000Version, "Nmea2000Version should round-trip")

	// String fields are fixed-width on the wire, so decoded strings may have trailing padding.
	trimmed := func(s string) string {
		return strings.TrimRightFunc(s, func(r rune) bool {
			return r == 0x00 || r == 0xFF || r == '@' || r == ' ' || r == 0xFFFD
		})
	}

	assert.Equal(t, "TestModel", trimmed(pi.ModelId), "ModelId should match after trimming padding")
	assert.Equal(t, "1.0.0", trimmed(pi.SoftwareVersionCode), "SoftwareVersionCode should match")
	assert.Equal(t, "v2", trimmed(pi.ModelVersion), "ModelVersion should match")
	assert.Equal(t, "SN12345", trimmed(pi.ModelSerialCode), "ModelSerialCode should match")

	require.NotNil(t, pi.CertificationLevel, "CertificationLevel should not be nil")
	assert.Equal(t, uint64(1), *pi.CertificationLevel, "CertificationLevel should match")

	require.NotNil(t, pi.LoadEquivalency, "LoadEquivalency should not be nil")
	assert.Equal(t, uint64(1), *pi.LoadEquivalency, "LoadEquivalency should match")
}

func TestClient_Write_FIFO_Ordering(t *testing.T) {
	ctx := context.Background()

	// 1. Create client with source address 5.
	writer, err := NewClient(ctx, Replay(nil), WithSourceAddress(5))
	require.NoError(t, err)
	defer func() { _ = writer.Close() }()

	// 2. Write 10 VesselHeading messages sequentially.
	const count = 10
	for i := 0; i < count; i++ {
		msg := &pgn.VesselHeading{
			Sid:     ptrUint64(uint64(i)),
			Heading: ptrUint64(uint64(i) * 1000),
		}
		msg.Info.Priority = ptrUint8(2)
		wr := writer.Write(msg)
		require.NoError(t, wr.Wait(), "write %d should succeed", i)
	}

	// 3. Get all captured frames.
	frames := writer.WrittenFrames()
	require.Len(t, frames, count, "each single-frame VesselHeading should produce 1 frame")

	// 4. Verify frames have consistent CAN IDs (all identical since same PGN/priority/source).
	expectedID := framer.BuildCANID(127250, 2, 5, 255)
	for i, f := range frames {
		assert.Equal(t, expectedID, f.ID, "frame %d should have expected CAN ID", i)
	}

	// 5. Feed all frames into Receive, verify 10 messages come out in order.
	var decoded []*pgn.VesselHeading
	for msg, err := range Receive(ctx, Replay(frames)) {
		require.NoError(t, err)
		vh, ok := msg.(*pgn.VesselHeading)
		require.True(t, ok, "expected *pgn.VesselHeading, got %T", msg)
		decoded = append(decoded, vh)
	}
	require.Len(t, decoded, count, "should decode exactly %d messages", count)

	// Verify ordering: each message's Sid should match its sequence position.
	for i, vh := range decoded {
		require.NotNil(t, vh.Sid, "message %d Sid should not be nil", i)
		assert.Equal(t, uint64(i), *vh.Sid, "message %d Sid should reflect write order", i)

		require.NotNil(t, vh.Heading, "message %d Heading should not be nil", i)
		assert.Equal(t, uint64(i)*1000, *vh.Heading, "message %d Heading should match write order", i)
	}
}

// --- Mock bus for bus path tests ---

type mockBus struct {
	inbound chan can.Frame
	written []can.Frame
	mu      sync.Mutex
	handler func(can.Frame)
}

func newMockBus() *mockBus {
	return &mockBus{inbound: make(chan can.Frame, 64)}
}

func (m *mockBus) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case f, ok := <-m.inbound:
			if !ok {
				return nil
			}
			if m.handler != nil {
				m.handler(f)
			}
		}
	}
}

func (m *mockBus) Close() error { return nil }

func (m *mockBus) WriteFrame(frame can.Frame) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.written = append(m.written, frame)
	return nil
}

func (m *mockBus) SetHandler(h func(can.Frame)) {
	m.handler = h
}

func (m *mockBus) getWritten() []can.Frame {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]can.Frame, len(m.written))
	copy(out, m.written)
	return out
}

func TestClient_AddressClaim(t *testing.T) {
	mb := newMockBus()

	c, err := NewClient(context.Background(),
		WithBus(mb),
		WithClaimTimeout(250*time.Millisecond),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	// The first written frame should be an address claim (PGN 60928).
	written := mb.getWritten()
	require.NotEmpty(t, written, "client should have written at least one frame (address claim)")

	// Parse the CAN ID to extract the PGN properly (accounting for addressed PGNs).
	firstFrame := written[0]
	rawPGN := (firstFrame.ID & 0x3FFFF00) >> 8
	pduFormat := uint8((rawPGN >> 8) & 0xFF)
	if pduFormat < 240 {
		rawPGN &= 0xFFF00 // mask off destination byte for addressed PGNs
	}
	assert.Equal(t, uint32(60928), rawPGN, "first frame should be address claim PGN 60928")
}

func TestClient_Write(t *testing.T) {
	mb := newMockBus()

	c, err := NewClient(context.Background(),
		WithBus(mb),
		WithClaimTimeout(250*time.Millisecond),
	)
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	heading := uint64(15000)
	msg := &pgn.VesselHeading{
		Heading: &heading,
	}
	msg.Info.Priority = ptrUint8(2)

	wr := c.Write(msg)
	require.NoError(t, wr.Wait())

	// Find the VesselHeading frame in the written frames (skip address claim frames).
	written := mb.getWritten()
	var found bool
	for _, f := range written {
		framePGN := (f.ID & 0x3FFFF00) >> 8
		if framePGN == 127250 {
			found = true
			break
		}
	}
	assert.True(t, found, "should find a VesselHeading (PGN 127250) frame in written frames")
}

func TestClient_Receive(t *testing.T) {
	mb := newMockBus()

	c, err := NewClient(context.Background(),
		WithBus(mb),
		WithClaimTimeout(250*time.Millisecond),
	)
	require.NoError(t, err)

	// Build and inject a VesselHeading frame into the mock bus.
	heading := uint64(15000)
	msg := &pgn.VesselHeading{
		Heading: &heading,
	}
	payload, err := msg.EncodePayload()
	require.NoError(t, err)

	canID := framer.BuildCANID(127250, 2, 42, 255)
	frame := framer.FrameSingle(canID, payload)

	// Inject the frame into the mock bus inbound channel.
	mb.inbound <- frame

	// Read messages via msgCh until we get a VesselHeading (skip address claim
	// echoes that may arrive first from the claimer's own WriteFrame calls).
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var received *pgn.VesselHeading
	done := make(chan struct{})
	go func() {
		defer close(done)
		for m, err := range c.Receive() {
			if err != nil {
				return
			}
			if vh, ok := m.(*pgn.VesselHeading); ok {
				received = vh
				return
			}
		}
	}()

	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("timed out waiting for VesselHeading message")
	}

	require.NotNil(t, received, "should have received a decoded VesselHeading")
	require.NotNil(t, received.Heading)
	assert.Equal(t, uint64(15000), *received.Heading)

	_ = c.Close()
}

func TestClient_FilterPGN(t *testing.T) {
	mb := newMockBus()

	// Filter: only accept PGN 127250 (VesselHeading).
	c, err := NewClient(context.Background(),
		WithBus(mb),
		WithClaimTimeout(250*time.Millisecond),
		Filter("pgn == 127250"),
	)
	require.NoError(t, err)

	// Build a VesselHeading frame (should pass filter).
	heading := uint64(15000)
	payload, err := (&pgn.VesselHeading{Heading: &heading}).EncodePayload()
	require.NoError(t, err)
	headingFrame := framer.FrameSingle(framer.BuildCANID(127250, 2, 42, 255), payload)

	// Build a SystemTime frame PGN 126992 (should be filtered out).
	sysPayload, err := (&pgn.SystemTime{}).EncodePayload()
	require.NoError(t, err)
	sysFrame := framer.FrameSingle(framer.BuildCANID(126992, 3, 42, 255), sysPayload)

	// Send the SystemTime first, then VesselHeading.
	mb.inbound <- sysFrame
	mb.inbound <- headingFrame

	// We should only receive VesselHeading.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var received pgn.Message
	done := make(chan struct{})
	go func() {
		defer close(done)
		for m, err := range c.Receive() {
			if err != nil {
				return
			}
			received = m
			return
		}
	}()

	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("timed out waiting for message")
	}

	require.NotNil(t, received)
	_, ok := received.(*pgn.VesselHeading)
	assert.True(t, ok, "expected VesselHeading, got %T", received)

	_ = c.Close()
}

// nonPGNMessage implements pgn.Message but not pgn.PGN, to prove the write
// path rejects it without panicking (the old reflection panicked on missing
// Info fields).
type nonPGNMessage struct{}

func (nonPGNMessage) PGNNumber() uint32 { return 127250 }

func TestClient_Write_NonPGNMessageErrors(t *testing.T) {
	ctx := context.Background()
	c, err := NewClient(ctx, Replay(nil))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer c.Close()

	res := c.Write(nonPGNMessage{})
	if err := res.Wait(); err == nil {
		t.Fatal("expected error writing a non-PGN message, got nil")
	}
}
