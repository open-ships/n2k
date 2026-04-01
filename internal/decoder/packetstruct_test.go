package decoder

import (
	"fmt"
	"testing"

	"github.com/open-ships/n2k/pgn"
	"github.com/stretchr/testify/assert"
)

// mockHandler is a test double for the Handler interface that captures all
// structs passed to HandleStruct for later assertion. This allows tests to verify what
// Decoder produces without needing a real downstream consumer.
type mockHandler struct {
	// received accumulates all structs passed to HandleStruct in order.
	received []any
}

// HandleStruct appends the received struct to the captured list for test verification.
func (m *mockHandler) HandleStruct(v any) {
	m.received = append(m.received, v)
}

// TestNew verifies that New returns a non-nil instance.
func TestNew(t *testing.T) {
	ps := New()
	assert.NotNil(t, ps)
}

// TestSetOutput verifies that SetOutput correctly registers a handler on Decoder,
// making it available for forwarding decoded structs.
func TestSetOutput(t *testing.T) {
	ps := New()
	handler := &mockHandler{}
	ps.SetOutput(handler)
	assert.NotNil(t, ps.handler)
}

// TestDecode_ValidDecoder verifies the happy path: a packet with a single working
// decoder produces the expected typed struct (VesselHeading) and forwards it to the
// handler. This confirms the full decode-and-forward pipeline works end to end.
func TestDecode_ValidDecoder(t *testing.T) {
	ps := New()
	handler := &mockHandler{}
	ps.SetOutput(handler)

	info := pgn.MessageInfo{PGN: 127250, SourceId: 1}
	// A mock decoder that always succeeds and returns a VesselHeading struct.
	dec := func(mi pgn.MessageInfo, s *pgn.PGNDataStream) (any, error) {
		return pgn.VesselHeading{Info: mi}, nil
	}

	pkt := Packet{
		Info:     info,
		Data:     []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		Decoders: []func(pgn.MessageInfo, *pgn.PGNDataStream) (any, error){dec},
	}

	ps.Decode(pkt)

	// Verify exactly one struct was forwarded and it's the correct type with correct PGN.
	assert.Equal(t, 1, len(handler.received))
	vh, ok := handler.received[0].(pgn.VesselHeading)
	assert.True(t, ok)
	assert.Equal(t, uint32(127250), vh.Info.PGN)
}

// TestDecode_DecoderFails_FallsToUnknown verifies that when the only decoder fails
// with an error, Decoder falls back to producing an UnknownPGN. This tests the error
// accumulation and fallback behavior.
func TestDecode_DecoderFails_FallsToUnknown(t *testing.T) {
	ps := New()
	handler := &mockHandler{}
	ps.SetOutput(handler)

	info := pgn.MessageInfo{PGN: 127250, SourceId: 1}
	// A mock decoder that always fails.
	failDecoder := func(mi pgn.MessageInfo, s *pgn.PGNDataStream) (any, error) {
		return nil, fmt.Errorf("decode error")
	}

	pkt := Packet{
		Info:     info,
		Data:     []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		Decoders: []func(pgn.MessageInfo, *pgn.PGNDataStream) (any, error){failDecoder},
	}

	ps.Decode(pkt)

	// Should still produce output, but as UnknownPGN instead of VesselHeading.
	assert.Equal(t, 1, len(handler.received))
	u, ok := handler.received[0].(pgn.UnknownPGN)
	assert.True(t, ok)
	assert.Equal(t, uint32(127250), u.Info.PGN)
}

// TestDecode_NoDecoders_SendsUnknown verifies that a packet with no decoders at all
// (nil Decoders slice) produces an UnknownPGN with a "no matching decoder" error. This
// happens when the PGN is unknown or all candidates were filtered out.
func TestDecode_NoDecoders_SendsUnknown(t *testing.T) {
	ps := New()
	handler := &mockHandler{}
	ps.SetOutput(handler)

	info := pgn.MessageInfo{PGN: 127250, SourceId: 1}
	pkt := Packet{
		Info:     info,
		Data:     []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		Decoders: nil, // No decoders available.
	}

	ps.Decode(pkt)

	assert.Equal(t, 1, len(handler.received))
	u, ok := handler.received[0].(pgn.UnknownPGN)
	assert.True(t, ok)
	assert.Contains(t, u.Reason.Error(), "no matching decoder")
}

// TestDecode_MultipleDecoders_FirstFails_SecondSucceeds verifies the decoder
// fallthrough behavior: when the first decoder fails, Decoder tries the next one.
// This simulates the real scenario where multiple proprietary variants exist for the
// same PGN and only one matches the actual data.
func TestDecode_MultipleDecoders_FirstFails_SecondSucceeds(t *testing.T) {
	ps := New()
	handler := &mockHandler{}
	ps.SetOutput(handler)

	info := pgn.MessageInfo{PGN: 127250, SourceId: 1}
	failDecoder := func(mi pgn.MessageInfo, s *pgn.PGNDataStream) (any, error) {
		return nil, fmt.Errorf("first decoder failed")
	}
	successDecoder := func(mi pgn.MessageInfo, s *pgn.PGNDataStream) (any, error) {
		return pgn.VesselHeading{Info: mi}, nil
	}

	pkt := Packet{
		Info: info,
		Data: []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		Decoders: []func(pgn.MessageInfo, *pgn.PGNDataStream) (any, error){
			failDecoder,    // This one fails...
			successDecoder, // ...so this one should be tried and succeed.
		},
	}

	ps.Decode(pkt)

	// The second decoder should have produced a VesselHeading.
	assert.Equal(t, 1, len(handler.received))
	_, ok := handler.received[0].(pgn.VesselHeading)
	assert.True(t, ok)
}

// TestDecode_NoHandler verifies that Decode does not panic when no
// Handler has been registered. This is a safety check ensuring the nil-guard
// in pgnReady works correctly.
func TestDecode_NoHandler(t *testing.T) {
	ps := New()
	// Intentionally no handler set -- should not panic.
	info := pgn.MessageInfo{PGN: 127250, SourceId: 1}
	dec := func(mi pgn.MessageInfo, s *pgn.PGNDataStream) (any, error) {
		return pgn.VesselHeading{Info: mi}, nil
	}

	pkt := Packet{
		Info:     info,
		Data:     []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		Decoders: []func(pgn.MessageInfo, *pgn.PGNDataStream) (any, error){dec},
	}

	assert.NotPanics(t, func() {
		ps.Decode(pkt)
	})
}

// TestDecode_AllDecodersFail verifies that when multiple decoders are present and
// ALL of them fail, the packet falls through to an UnknownPGN. This tests the exhaustive
// decoder-attempt loop and the final fallback at the end of the loop.
func TestDecode_AllDecodersFail(t *testing.T) {
	ps := New()
	handler := &mockHandler{}
	ps.SetOutput(handler)

	info := pgn.MessageInfo{PGN: 127250, SourceId: 1}
	fail1 := func(mi pgn.MessageInfo, s *pgn.PGNDataStream) (any, error) {
		return nil, fmt.Errorf("fail1")
	}
	fail2 := func(mi pgn.MessageInfo, s *pgn.PGNDataStream) (any, error) {
		return nil, fmt.Errorf("fail2")
	}

	pkt := Packet{
		Info:     info,
		Data:     []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		Decoders: []func(pgn.MessageInfo, *pgn.PGNDataStream) (any, error){fail1, fail2},
	}

	ps.Decode(pkt)

	// Should produce an UnknownPGN with a non-nil Reason containing the error details.
	assert.Equal(t, 1, len(handler.received))
	u, ok := handler.received[0].(pgn.UnknownPGN)
	assert.True(t, ok)
	assert.NotNil(t, u.Reason)
}
