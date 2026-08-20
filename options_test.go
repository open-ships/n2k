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

func TestConfigValidationRejectsInvalidOperationalValues(t *testing.T) {
	for _, opt := range []Option{
		WithSourceAddress(252),
		WithSourceAddress(254),
		WithPreferredAddress(252),
		WithClaimTimeout(0),
		WithHeartbeatInterval(-time.Second),
		WithReceiveBuffer(0),
		WithWriteQueue(0),
	} {
		_, err := NewClient(context.Background(), Replay(nil), opt)
		assert.Error(t, err)
	}
}

func TestAddressOptions(t *testing.T) {
	for _, opt := range []Option{WithSourceAddress(0), WithSourceAddress(251), WithPreferredAddress(0), WithPreferredAddress(251)} {
		c, err := NewClient(context.Background(), Replay(nil), opt)
		require.NoError(t, err)
		require.NoError(t, c.Close())
	}

	_, err := NewClient(context.Background(), Replay(nil), WithSourceAddress(10), WithPreferredAddress(10))
	require.Error(t, err)
}

func TestBufferOptions(t *testing.T) {
	c, err := NewClient(context.Background(), Replay(nil), WithReceiveBuffer(3), WithWriteQueue(5))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()
	status := c.Status()
	assert.Equal(t, 5, status.WriteQueueCapacity)
	assert.True(t, status.AddressClaimed)
}

func TestNewClientRejectsAmbiguousSources(t *testing.T) {
	_, err := NewClient(context.Background(), Replay(nil), Replay(nil))
	assert.NoError(t, err, "multiple replay sources are valid")

	_, err = NewClient(context.Background(), CAN("can0"), Replay(nil))
	assert.Error(t, err)
}

func TestConfigValidationRejectsInvalidCaptureAndSerialSources(t *testing.T) {
	for _, opt := range []Option{
		Serial("", FormatActisense),
		Serial("/dev/ttyUSB0", FormatYDRaw),
		EBL(""),
	} {
		cfg := config{}
		opt.apply(&cfg)
		assert.Error(t, cfg.validate())
	}
}

func TestConfigValidationRequiresTCPForReconnect(t *testing.T) {
	cfg := config{}
	Serial("/dev/ttyUSB0", FormatActisenseRaw).apply(&cfg)
	WithReconnect(ReconnectPolicy{}).apply(&cfg)
	assert.ErrorContains(t, cfg.validate(), "requires a TCP source")
}
