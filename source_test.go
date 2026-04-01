package n2k

import (
	"context"
	"log/slog"
	"testing"

	"github.com/brutella/can"
	"github.com/stretchr/testify/assert"
)

func TestReplaySource(t *testing.T) {
	frames := []can.Frame{
		{ID: 0x09F10D01, Length: 8, Data: [8]uint8{1, 2, 3, 4, 5, 6, 7, 8}},
		{ID: 0x09F10D02, Length: 8, Data: [8]uint8{8, 7, 6, 5, 4, 3, 2, 1}},
	}

	src := &replaySource{frames: frames}
	received := make([]can.Frame, 0)
	ctx := context.Background()

	err := src.run(ctx, slog.Default(), func(f can.Frame) {
		received = append(received, f)
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, len(received))
	assert.Equal(t, frames[0].ID, received[0].ID)
	assert.Equal(t, frames[1].ID, received[1].ID)
}

func TestFanIn(t *testing.T) {
	frames1 := []can.Frame{
		{ID: 0x09F10D01, Length: 8, Data: [8]uint8{1, 2, 3, 4, 5, 6, 7, 8}},
	}
	frames2 := []can.Frame{
		{ID: 0x09F10D02, Length: 8, Data: [8]uint8{8, 7, 6, 5, 4, 3, 2, 1}},
	}

	sources := []source{
		&replaySource{frames: frames1},
		&replaySource{frames: frames2},
	}

	received := make([]can.Frame, 0)
	ctx := context.Background()

	err := runSources(ctx, slog.Default(), sources, func(f can.Frame) {
		received = append(received, f)
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, len(received))
}

func TestReplaySourceContextCancel(t *testing.T) {
	frames := make([]can.Frame, 1000)
	for i := range frames {
		frames[i] = can.Frame{ID: uint32(i)}
	}

	ctx, cancel := context.WithCancel(context.Background())
	src := &replaySource{frames: frames}

	count := 0
	cancel() // cancel before running
	err := src.run(ctx, slog.Default(), func(f can.Frame) {
		count++
	})

	assert.Error(t, err)
	assert.Less(t, count, 1000)
}
