package n2k

import (
	"log/slog"
	"testing"

	"github.com/brutella/can"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	var cfg config
	opts := []Option{
		CAN("can0"),
		CAN("can1"),
		USB("/dev/ttyUSB0"),
		Replay([]can.Frame{{ID: 1}}),
		IncludeUnknown(),
		WithLogger(slog.Default()),
	}

	for _, o := range opts {
		o.apply(&cfg)
	}

	assert.Equal(t, 4, len(cfg.sources))
	assert.True(t, cfg.includeUnknown)
	assert.NotNil(t, cfg.logger)
}

func TestNoSourcesError(t *testing.T) {
	cfg := config{}
	err := cfg.validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one source")
}

func TestFilterOption(t *testing.T) {
	var cfg config
	Filter(`pgn == 127250`).apply(&cfg)
	assert.Equal(t, `pgn == 127250`, cfg.filterExpr)
}
