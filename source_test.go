package n2k

import (
	"context"
	"log/slog"
	"testing"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/raw"
	"github.com/stretchr/testify/assert"
)

func TestReplaySource(t *testing.T) {
	frames := []can.Frame{
		{ID: 0x09F10D01, Length: 8, Data: [8]uint8{1, 2, 3, 4, 5, 6, 7, 8}},
		{ID: 0x09F10D02, Length: 8, Data: [8]uint8{8, 7, 6, 5, 4, 3, 2, 1}},
	}

	src := replaySourceForTest(frames)
	received := make([]can.Frame, 0)
	ctx := context.Background()

	err := src.run(ctx, slog.Default(), func(observation raw.Observation) {
		received = append(received, *observation.Frame)
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
		replaySourceForTest(frames1),
		replaySourceForTest(frames2),
	}

	received := make([]can.Frame, 0)
	ctx := context.Background()

	err := runSources(ctx, slog.Default(), sources, func(observation raw.Observation) {
		received = append(received, *observation.Frame)
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
	src := replaySourceForTest(frames)

	count := 0
	cancel() // cancel before running
	err := src.run(ctx, slog.Default(), func(raw.Observation) {
		count++
	})

	assert.Error(t, err)
	assert.Less(t, count, 1000)
}

func replaySourceForTest(frames []can.Frame) *replaySource {
	observations := make([]raw.Observation, 0, len(frames))
	for _, frame := range frames {
		observations = append(observations, frameObservation(frame, "replay", "replay", raw.DirectionReceived))
	}
	return &replaySource{observations: observations}
}
