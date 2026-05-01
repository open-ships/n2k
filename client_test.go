package n2k

import (
	"context"
	"strings"
	"testing"

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

func TestClient_Write_EncodesAndFrames(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	heading := float32(1.5)
	msg := &pgn.VesselHeading{
		Info: pgn.MessageInfo{
			PGN:      127250,
			Priority: 2,
		},
		Heading: &heading,
	}

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

	heading := float32(0.5)
	msg := &pgn.VesselHeading{
		Info: pgn.MessageInfo{
			PGN:      127250,
			Priority: 3,
		},
		Heading: &heading,
	}

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

	heading := float32(0.5)
	msg := &pgn.VesselHeading{
		Info: pgn.MessageInfo{
			PGN: 127250,
			// Priority left at 0 => default to 6
		},
		Heading: &heading,
	}

	wr := c.Write(msg)
	require.NoError(t, wr.Wait())

	frames := c.WrittenFrames()
	require.Len(t, frames, 1)

	expectedID := framer.BuildCANID(127250, 6, 0, 255)
	assert.Equal(t, expectedID, frames[0].ID)
}

func TestClient_Write_UnknownType(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	// A plain struct with no Info field.
	type NotAPGN struct {
		Value int
	}
	wr := c.Write(&NotAPGN{Value: 42})
	err = wr.Wait()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no Info field")
}

func TestClient_Write_NonPointer(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	wr := c.Write("not a struct")
	err = wr.Wait()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected pointer to PGN struct")
}

func TestClient_Write_ZeroPGN(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	msg := &pgn.VesselHeading{
		Info: pgn.MessageInfo{
			PGN: 0, // zero PGN
		},
	}

	wr := c.Write(msg)
	err = wr.Wait()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "PGN number is 0")
}

func TestClient_Write_NoEncoder(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	// Use a PGN struct with a PGN number that has no encoder registered.
	// We use an unknown PGN and fake the Info.
	msg := &pgn.UnknownPGN{
		Info: pgn.MessageInfo{
			PGN: 999999, // unlikely to have an encoder
		},
	}

	wr := c.Write(msg)
	err = wr.Wait()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no encoder registered")
}

func TestClient_Write_AfterClose(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil))
	require.NoError(t, err)

	require.NoError(t, c.Close())

	heading := float32(1.0)
	msg := &pgn.VesselHeading{
		Info: pgn.MessageInfo{PGN: 127250},
		Heading: &heading,
	}
	wr := c.Write(msg)
	err = wr.Wait()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "closed")
}

func TestClient_Receive(t *testing.T) {
	// Build a pre-encoded VesselHeading frame for replay.
	heading := float32(1.5)
	msg := &pgn.VesselHeading{
		Heading: &heading,
	}
	encoder := pgn.EncoderLookup[127250]
	require.NotNil(t, encoder, "VesselHeading encoder must exist")
	payload, err := encoder(msg)
	require.NoError(t, err)

	canID := framer.BuildCANID(127250, 2, 42, 255)
	frame := framer.FrameSingle(canID, payload)

	c, err := NewClient(context.Background(), Replay([]can.Frame{frame}))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	var received []any
	for msg, err := range c.Receive() {
		require.NoError(t, err)
		received = append(received, msg)
	}
	require.Len(t, received, 1)

	vh, ok := received[0].(pgn.VesselHeading)
	require.True(t, ok, "expected VesselHeading, got %T", received[0])
	require.NotNil(t, vh.Heading)
	assert.InDelta(t, 1.5, float64(*vh.Heading), 0.001)
}

func TestClient_Scanner(t *testing.T) {
	heading := float32(2.0)
	msg := &pgn.VesselHeading{
		Heading: &heading,
	}
	encoder := pgn.EncoderLookup[127250]
	require.NotNil(t, encoder)
	payload, err := encoder(msg)
	require.NoError(t, err)

	canID := framer.BuildCANID(127250, 2, 42, 255)
	frame := framer.FrameSingle(canID, payload)

	c, err := NewClient(context.Background(), Replay([]can.Frame{frame}))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	s := c.Scanner()
	require.True(t, s.Next())
	vh, ok := s.Message().(pgn.VesselHeading)
	require.True(t, ok, "expected VesselHeading, got %T", s.Message())
	require.NotNil(t, vh.Heading)
	assert.InDelta(t, 2.0, float64(*vh.Heading), 0.001)
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

	heading := float32(0.1)
	for i := 0; i < 3; i++ {
		msg := &pgn.VesselHeading{
			Info: pgn.MessageInfo{PGN: 127250, Priority: 2},
			Heading: &heading,
		}
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
	ver := float32(2.0)
	code := uint16(1234)
	cert := uint8(1)
	load := uint8(1)
	msg := &pgn.ProductInformation{
		Info:                pgn.MessageInfo{PGN: 126996, Priority: 6},
		Nmea2000Version:     &ver,
		ProductCode:         &code,
		ModelId:             "TestModel",
		SoftwareVersionCode: "1.0",
		ModelVersion:        "v1",
		ModelSerialCode:     "SN001",
		CertificationLevel:  &cert,
		LoadEquivalency:     &load,
	}

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

// --- Pointer helpers for building test PGN structs ---

func ptrFloat32(v float32) *float32 { return &v }
func ptrUint8(v uint8) *uint8       { return &v }
func ptrUint16(v uint16) *uint16    { return &v }

// --- End-to-end integration tests (write -> read round-trip) ---

func TestClient_RoundTrip_SingleFrame(t *testing.T) {
	ctx := context.Background()

	// 1. Build a VesselHeading with known values.
	original := &pgn.VesselHeading{
		Info: pgn.MessageInfo{
			PGN:      127250,
			Priority: 2,
		},
		Sid:     ptrUint8(0),
		Heading: ptrFloat32(1.5708),
	}

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
	var decoded []any
	for msg, err := range Receive(ctx, Replay(frames)) {
		require.NoError(t, err)
		decoded = append(decoded, msg)
	}
	require.Len(t, decoded, 1, "should decode exactly 1 message from the captured frames")

	// 5. Assert the decoded VesselHeading matches.
	vh, ok := decoded[0].(pgn.VesselHeading)
	require.True(t, ok, "expected pgn.VesselHeading, got %T", decoded[0])

	require.NotNil(t, vh.Heading, "Heading should not be nil")
	assert.InDelta(t, 1.5708, float64(*vh.Heading), 0.001, "Heading should round-trip within tolerance")

	require.NotNil(t, vh.Sid, "Sid should not be nil")
	assert.Equal(t, uint8(0), *vh.Sid, "Sid should match")
}

func TestClient_RoundTrip_FastPacket(t *testing.T) {
	ctx := context.Background()

	// 1. Build a ProductInformation with known values.
	original := &pgn.ProductInformation{
		Info: pgn.MessageInfo{
			PGN:      126996,
			Priority: 6,
		},
		Nmea2000Version:     ptrFloat32(2.0),
		ProductCode:         ptrUint16(42),
		ModelId:             "TestModel",
		SoftwareVersionCode: "1.0.0",
		ModelVersion:        "v2",
		ModelSerialCode:     "SN12345",
		CertificationLevel:  ptrUint8(1),
		LoadEquivalency:     ptrUint8(1),
	}

	// 2. Write via client — should produce multiple frames (fast-packet).
	writer, err := NewClient(ctx, Replay(nil))
	require.NoError(t, err)
	defer func() { _ = writer.Close() }()

	wr := writer.Write(original)
	require.NoError(t, wr.Wait())

	// 3. Get captured frames, verify more than 1 frame.
	frames := writer.WrittenFrames()
	require.Greater(t, len(frames), 1, "ProductInformation is fast-packet and must produce multiple frames")

	// 4. Feed frames into Receive to decode.
	var decoded []any
	for msg, err := range Receive(ctx, Replay(frames)) {
		require.NoError(t, err)
		decoded = append(decoded, msg)
	}
	require.Len(t, decoded, 1, "should decode exactly 1 ProductInformation from the captured frames")

	// 5. Assert decoded ProductInformation matches.
	pi, ok := decoded[0].(pgn.ProductInformation)
	require.True(t, ok, "expected pgn.ProductInformation, got %T", decoded[0])

	require.NotNil(t, pi.ProductCode, "ProductCode should not be nil")
	assert.Equal(t, uint16(42), *pi.ProductCode, "ProductCode should match")

	require.NotNil(t, pi.Nmea2000Version, "Nmea2000Version should not be nil")
	assert.InDelta(t, 2.0, float64(*pi.Nmea2000Version), 0.01, "Nmea2000Version should round-trip")

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
	assert.Equal(t, uint8(1), *pi.CertificationLevel, "CertificationLevel should match")

	require.NotNil(t, pi.LoadEquivalency, "LoadEquivalency should not be nil")
	assert.Equal(t, uint8(1), *pi.LoadEquivalency, "LoadEquivalency should match")
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
			Info: pgn.MessageInfo{
				PGN:      127250,
				Priority: 2,
			},
			Sid:     ptrUint8(uint8(i)),
			Heading: ptrFloat32(float32(i) * 0.1),
		}
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
	var decoded []pgn.VesselHeading
	for msg, err := range Receive(ctx, Replay(frames)) {
		require.NoError(t, err)
		vh, ok := msg.(pgn.VesselHeading)
		require.True(t, ok, "expected pgn.VesselHeading, got %T", msg)
		decoded = append(decoded, vh)
	}
	require.Len(t, decoded, count, "should decode exactly %d messages", count)

	// Verify ordering: each message's Sid should match its sequence position.
	for i, vh := range decoded {
		require.NotNil(t, vh.Sid, "message %d Sid should not be nil", i)
		assert.Equal(t, uint8(i), *vh.Sid, "message %d Sid should reflect write order", i)

		require.NotNil(t, vh.Heading, "message %d Heading should not be nil", i)
		assert.InDelta(t, float64(i)*0.1, float64(*vh.Heading), 0.001,
			"message %d Heading should match write order", i)
	}
}
