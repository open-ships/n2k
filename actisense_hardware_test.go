package n2k

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const actisenseHardwareConfigEnvironment = "N2K_ACTISENSE_HARDWARE_CONFIG"

type actisenseHardwareMatrix struct {
	Devices []actisenseHardwareDevice `json:"devices"`
}

type actisenseHardwareDevice struct {
	Name             string                `json:"name"`
	Transport        string                `json:"transport"`
	Endpoint         string                `json:"endpoint"`
	Serial           ActisenseSerialConfig `json:"serial"`
	ExpectedModelID  ActisenseModelID      `json:"expected_model_id"`
	SupportedPGNs    bool                  `json:"supported_pgn_list"`
	F2Lists          bool                  `json:"f2_lists"`
	PortInventory    bool                  `json:"port_inventory"`
	RawClient        bool                  `json:"raw_client"`
	ClientSource     uint8                 `json:"client_source"`
	RemoteSource     *uint8                `json:"remote_source"`
	RemoteModelID    ActisenseModelID      `json:"remote_model_id"`
	CommandTimeoutMS int                   `json:"command_timeout_ms"`
}

func TestActisenseHardwareMatrix(t *testing.T) {
	configPath := os.Getenv(actisenseHardwareConfigEnvironment)
	if configPath == "" {
		t.Skipf("set %s to an Actisense lab matrix; see conformance/actisense-hardware.example.json", actisenseHardwareConfigEnvironment)
	}
	contents, err := os.ReadFile(configPath)
	require.NoError(t, err)
	var matrix actisenseHardwareMatrix
	require.NoError(t, json.Unmarshal(contents, &matrix))
	require.NotEmpty(t, matrix.Devices)
	for _, device := range matrix.Devices {
		device := device
		t.Run(device.Name, func(t *testing.T) {
			runActisenseHardwareDevice(t, device)
		})
	}
}

func runActisenseHardwareDevice(t *testing.T, device actisenseHardwareDevice) {
	t.Helper()
	require.NotEmpty(t, device.Name)
	require.NotEmpty(t, device.Endpoint)
	require.NotZero(t, device.ExpectedModelID)
	timeout := 30 * time.Second
	if device.CommandTimeoutMS > 0 {
		timeout = time.Duration(device.CommandTimeoutMS) * time.Millisecond
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var traceOutput bytes.Buffer
	writer, err := NewEBLWriter(&traceOutput, WithEBLDescription("n2k Actisense hardware conformance: "+device.Name))
	require.NoError(t, err)
	trace, err := NewActisenseEBLTrace(writer)
	require.NoError(t, err)
	options := []ActisenseSessionOption{
		WithActisenseCommandTimeout(timeout),
		WithActisenseSessionReadyTimeout(timeout),
		WithActisenseWireTrace(trace),
	}
	session, err := openActisenseHardwareSession(ctx, device, options)
	require.NoError(t, err)
	defer func() { require.NoError(t, session.Close()) }()

	status := session.Status()
	assert.True(t, status.Connected)
	assert.False(t, status.SourceAuthoritative)
	mode, err := session.GetOperatingMode(ctx)
	require.NoError(t, err)
	assert.Equal(t, ActisenseModeTransferReceiveAll, mode)
	product, err := session.GetProductInfo(ctx)
	require.NoError(t, err)
	assert.Equal(t, device.ExpectedModelID, product.DeviceModelID)
	assert.NotEmpty(t, product.Model)
	echo := []byte{0x00, 0x10, 0x1B, 0xFF}
	echoed, err := session.Echo(ctx, echo)
	require.NoError(t, err)
	assert.Equal(t, echo, echoed)
	if device.SupportedPGNs {
		supported, err := session.GetSupportedPGNs(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, supported.Entries)
	}
	if device.F2Lists {
		rxList, err := session.GetRxPGNEnableList(ctx)
		require.NoError(t, err)
		assert.NotNil(t, rxList.Entries)
		txList, err := session.GetTxPGNEnableList(ctx)
		require.NoError(t, err)
		assert.NotNil(t, txList.Entries)
	}
	if device.PortInventory {
		inventory, err := session.GetPortInventory(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, inventory.Ports)
	}
	metrics := session.Status().Metrics
	assert.GreaterOrEqual(t, metrics.Protocol.BEMRequests, uint64(3))
	assert.GreaterOrEqual(t, writer.Metrics().Records, uint64(3))

	if !device.RawClient {
		return
	}
	require.NoError(t, session.Close())
	runActisenseRawHardwareChecks(t, ctx, device, timeout)
}

func openActisenseHardwareSession(ctx context.Context, device actisenseHardwareDevice, options []ActisenseSessionOption) (*ActisenseGatewaySession, error) {
	switch device.Transport {
	case "tcp":
		return NewActisenseTCPSession(ctx, device.Endpoint, options...)
	case "serial":
		return NewActisenseSerialSession(ctx, device.Endpoint, device.Serial, options...)
	default:
		return nil, fmt.Errorf("unknown Actisense hardware transport %q", device.Transport)
	}
}

func runActisenseRawHardwareChecks(t *testing.T, ctx context.Context, device actisenseHardwareDevice, timeout time.Duration) {
	t.Helper()
	var source Option
	switch device.Transport {
	case "tcp":
		source = ActisenseTCP(device.Endpoint)
	case "serial":
		source = ActisenseSerial(device.Endpoint, WithActisenseSerialConfig(device.Serial))
	default:
		t.Fatalf("unknown Actisense hardware transport %q", device.Transport)
	}
	client, err := NewClient(
		ctx,
		source,
		WithSourceAddress(device.ClientSource),
		WithReadyTimeout(timeout),
		WithClaimTimeout(timeout),
		WithHeartbeatInterval(0),
	)
	require.NoError(t, err)
	defer func() { require.NoError(t, client.Close()) }()
	status := client.Status()
	assert.True(t, status.AddressClaimed)
	require.NotNil(t, status.Actisense)
	assert.GreaterOrEqual(t, status.Actisense.ConnectionEpochs, uint64(1))
	if device.RemoteSource == nil {
		return
	}
	remote, err := client.ActisenseRemoteDevice(*device.RemoteSource, WithActisenseRemoteTimeout(timeout))
	require.NoError(t, err)
	product, err := remote.GetProductInfo(ctx)
	require.NoError(t, err)
	if device.RemoteModelID != ActisenseModelUnknown {
		assert.Equal(t, device.RemoteModelID, product.DeviceModelID)
	}
	echo := []byte{0x10, 0x00, 0xFF}
	echoed, err := remote.Echo(ctx, echo)
	require.NoError(t, err)
	assert.Equal(t, echo, echoed)
	assert.GreaterOrEqual(t, remote.Metrics().BEMRequests, uint64(2))
}
