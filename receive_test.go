package n2k

import (
	"context"
	"testing"

	"github.com/brutella/can"
	"github.com/stretchr/testify/assert"
)

func TestReceiveBasic(t *testing.T) {
	ctx := context.Background()
	count := 0
	for msg, err := range Receive(ctx, Replay([]can.Frame{testFrame127501}), IncludeUnknown()) {
		assert.NoError(t, err)
		assert.NotNil(t, msg)
		count++
	}
	assert.GreaterOrEqual(t, count, 1)
}

func TestReceiveNoSources(t *testing.T) {
	ctx := context.Background()
	for _, err := range Receive(ctx) {
		assert.Error(t, err)
		break
	}
}

func TestReceiveWithFilter(t *testing.T) {
	ctx := context.Background()
	count := 0
	for msg, err := range Receive(ctx, Replay([]can.Frame{testFrame127501}), Filter(`pgn == 0`), IncludeUnknown()) {
		assert.NoError(t, err)
		_ = msg
		count++
	}
	assert.Equal(t, 0, count)
}
