package n2k

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
	"github.com/open-ships/n2k/pgn"
)

// testLogger returns a slog.Logger that discards all output, for use in
// tests that need to satisfy config.logger without polluting test output.
func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func vesselHeadingFrame(t *testing.T) can.Frame {
	t.Helper()
	h := uint64(15708)
	payload, err := pgn.EncodeMessage(&pgn.VesselHeading{Heading: &h})
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	id := framer.BuildCANID(127250, 2, 42, 255)
	return framer.FrameSingle(id, payload)
}

func TestReadPipeline_DecodesFrame(t *testing.T) {
	ctx := context.Background()
	out := make(chan pgn.Message, 1)
	p, err := newReadPipeline(ctx, config{logger: testLogger()}, out)
	if err != nil {
		t.Fatalf("newReadPipeline: %v", err)
	}
	p.HandleFrame(vesselHeadingFrame(t))
	select {
	case msg := <-out:
		if _, ok := msg.(*pgn.VesselHeading); !ok {
			t.Fatalf("got %T, want *pgn.VesselHeading", msg)
		}
	case <-time.After(time.Second):
		t.Fatal("no message delivered")
	}
}

func TestReadPipeline_FilterCompileErrorIsEager(t *testing.T) {
	out := make(chan pgn.Message, 1)
	_, err := newReadPipeline(context.Background(), config{logger: testLogger(), filterExpr: "((("}, out)
	if err == nil {
		t.Fatal("expected filter compile error")
	}
}

func TestReadPipeline_DropsUnknownByDefault(t *testing.T) {
	ctx := context.Background()
	out := make(chan pgn.Message, 1)
	p, err := newReadPipeline(ctx, config{logger: testLogger()}, out)
	if err != nil {
		t.Fatalf("newReadPipeline: %v", err)
	}
	// PGN 1 is not in the metadata: becomes UnknownPGN, dropped by default.
	p.HandleFrame(can.Frame{ID: framer.BuildCANID(1, 6, 9, 255), Length: 8})
	select {
	case msg := <-out:
		t.Fatalf("expected drop, got %T", msg)
	case <-time.After(50 * time.Millisecond):
	}
}

func TestReadPipeline_InjectAssembledAppliesPreFilter(t *testing.T) {
	ctx := context.Background()
	out := make(chan pgn.Message, 2)
	p, err := newReadPipeline(ctx, config{logger: testLogger(), filterExpr: "pgn == 127250"}, out)
	if err != nil {
		t.Fatalf("newReadPipeline: %v", err)
	}

	// Non-matching PGN must be dropped by the pre-filter before decode.
	h := uint64(15708)
	payload, err := pgn.EncodeMessage(&pgn.VesselHeading{Heading: &h})
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	priority := uint8(6)
	infoWrongPGN := pgn.MessageInfo{Timestamp: time.Now(), PGN: 130306, SourceId: 7, Priority: &priority}
	p.InjectAssembled(infoWrongPGN, payload)

	select {
	case msg := <-out:
		t.Fatalf("expected pre-filter drop, got %T", msg)
	case <-time.After(50 * time.Millisecond):
	}

	// Matching PGN must pass the pre-filter and decode.
	infoMatch := pgn.MessageInfo{Timestamp: time.Now(), PGN: 127250, SourceId: 7, Priority: &priority}
	p.InjectAssembled(infoMatch, payload)

	select {
	case msg := <-out:
		if _, ok := msg.(*pgn.VesselHeading); !ok {
			t.Fatalf("got %T, want *pgn.VesselHeading", msg)
		}
	case <-time.After(time.Second):
		t.Fatal("no message delivered for matching PGN")
	}
}
