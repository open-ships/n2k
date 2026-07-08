package n2k

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/brutella/can"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	Filter("pgn == 127250").apply(&cfg)
	assert.Equal(t, `pgn == 127250`, cfg.filterExpr)
}

func TestWithClaimTimeout(t *testing.T) {
	cfg := config{}
	WithClaimTimeout(2 * time.Second).apply(&cfg)
	require.NotNil(t, cfg.claimTimeout)
	assert.Equal(t, 2*time.Second, *cfg.claimTimeout)
}

func TestNewClient_InvalidDeviceNameRejected(t *testing.T) {
	_, err := NewClient(context.Background(),
		Replay(nil),
		WithName(DeviceName{ManufacturerCode: 5000}), // > 11 bits, must fail Validate
	)
	if err == nil {
		t.Fatal("expected invalid DeviceName error")
	}
}
