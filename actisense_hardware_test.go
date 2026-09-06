package n2k

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const actisenseHardwareConfigEnvironment = "N2K_ACTISENSE_HARDWARE_CONFIG"

type actisenseHardwareMatrix struct {
	Devices           []actisenseHardwareDevice `json:"devices"`
	ArtifactDirectory string                    `json:"artifact_directory"`
}

type actisenseHardwareDevice struct {
	Name                string                  `json:"name"`
	Transport           string                  `json:"transport"`
	Endpoint            string                  `json:"endpoint"`
	Serial              ActisenseSerialConfig   `json:"serial"`
	ExpectedModelID     ActisenseModelID        `json:"expected_model_id"`
	SupportedPGNs       bool                    `json:"supported_pgn_list"`
	F2Lists             bool                    `json:"f2_lists"`
	PortInventory       bool                    `json:"port_inventory"`
	RawClient           bool                    `json:"raw_client"`
	ClientSource        uint8                   `json:"client_source"`
	RemoteSource        *uint8                  `json:"remote_source"`
	RemoteModelID       ActisenseModelID        `json:"remote_model_id"`
	CommandTimeoutMS    int                     `json:"command_timeout_ms"`
	ExpectedFirmware    string                  `json:"expected_firmware"`
	PreserveMode        bool                    `json:"preserve_mode"`
	ExpectedMode        *ActisenseOperatingMode `json:"expected_mode"`
	F1Lists             bool                    `json:"f1_lists"`
	PortDuplicateDelete bool                    `json:"port_duplicate_delete"`
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
	if matrix.ArtifactDirectory == "" {
		matrix.ArtifactDirectory = filepath.Join("conformance-artifacts", "actisense-hardware-"+time.Now().UTC().Format("20060102T150405.000000000Z"))
	}
	require.NoError(t, os.MkdirAll(matrix.ArtifactDirectory, 0o750))
	for index, device := range matrix.Devices {
		device := device
		t.Run(device.Name, func(t *testing.T) {
			runActisenseHardwareDevice(t, device, filepath.Join(matrix.ArtifactDirectory, fmt.Sprintf("device-%02d", index)))
		})
	}
}

func runActisenseHardwareDevice(t *testing.T, device actisenseHardwareDevice, artifactPrefix string) {
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

	traceOutput, err := os.OpenFile(artifactPrefix+".ebl", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o600)
	require.NoError(t, err)
	revision, revisionErr := exec.Command("git", "rev-parse", "HEAD").Output()
	worktree, worktreeErr := exec.Command("git", "status", "--porcelain").Output()
	var product ActisenseProductInfo
	var remoteProduct *ActisenseProductInfo
	var status ActisenseSessionStatus
	defer func() {
		defer func() { assert.NoError(t, traceOutput.Close()) }()
		_, seekErr := traceOutput.Seek(0, io.SeekStart)
		assert.NoError(t, seekErr)
		digest := sha256.New()
		_, copyErr := io.Copy(digest, traceOutput)
		assert.NoError(t, copyErr)
		result := "pass"
		assert.NoError(t, revisionErr)
		assert.NoError(t, worktreeErr)
		if t.Failed() {
			result = "fail"
		}
		evidence := struct {
			Device        actisenseHardwareDevice `json:"device"`
			Product       ActisenseProductInfo    `json:"product"`
			RemoteProduct *ActisenseProductInfo   `json:"remote_product,omitempty"`
			Status        ActisenseSessionStatus  `json:"status"`
			Result        string                  `json:"result"`
			CapturedAt    time.Time               `json:"captured_at"`
			CaptureSHA256 string                  `json:"capture_sha256"`
			SDKCommit     string                  `json:"sdk_commit"`
			RunnerCommit  string                  `json:"runner_commit"`
			RunnerDirty   bool                    `json:"runner_dirty"`
		}{device, product, remoteProduct, status, result, time.Now().UTC(), hex.EncodeToString(digest.Sum(nil)), "ed2268a6e8db0645f75e4ef17eed2e937d025040", strings.TrimSpace(string(revision)), len(worktree) != 0}
		encoded, encodeErr := json.MarshalIndent(evidence, "", "  ")
		assert.NoError(t, encodeErr)
		assert.NoError(t, os.WriteFile(artifactPrefix+".json", append(encoded, '\n'), 0o600))
		t.Logf("Actisense %s evidence: %s.json and .ebl", result, artifactPrefix)
	}()
	writer, err := NewEBLWriter(traceOutput, WithEBLDescription("n2k Actisense hardware conformance: "+device.Name))
	require.NoError(t, err)
	trace, err := NewActisenseEBLTrace(writer)
	require.NoError(t, err)
	options := []ActisenseSessionOption{
		WithActisenseCommandTimeout(timeout),
		WithActisenseSessionReadyTimeout(timeout),
		WithActisenseWireTrace(trace),
	}
	if device.PreserveMode {
		options = append(options, WithActisensePreserveOperatingMode())
	}
	session, err := openActisenseHardwareSession(ctx, device, options)
	require.NoError(t, err)
	defer func() {
		status = session.Status()
		assert.NoError(t, session.Close())
	}()

	status = session.Status()
	assert.True(t, status.Connected)
	assert.False(t, status.SourceAuthoritative)
	mode, err := session.GetOperatingMode(ctx)
	require.NoError(t, err)
	if device.ExpectedMode != nil {
		assert.Equal(t, *device.ExpectedMode, mode)
	} else if !device.PreserveMode {
		assert.Equal(t, ActisenseModeTransferReceiveAll, mode)
	}
	product, err = session.GetProductInfo(ctx)
	require.NoError(t, err)
	assert.Equal(t, device.ExpectedModelID, product.DeviceModelID)
	assert.NotEmpty(t, product.Model)
	if device.ExpectedFirmware != "" {
		assert.Equal(t, device.ExpectedFirmware, product.SoftwareVersion)
	}
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
	if device.F1Lists {
		_, err := session.GetRxPGNEnableListF1(ctx)
		require.NoError(t, err)
		_, err = session.GetTxPGNEnableListF1(ctx)
		require.NoError(t, err)
	}
	if device.PortDuplicateDelete {
		_, err := session.GetPortDuplicateDelete(ctx)
		require.NoError(t, err)
	}
	if device.RemoteSource != nil {
		remote, err := session.ActisenseRemoteDevice(*device.RemoteSource, WithActisenseRemoteTimeout(timeout))
		require.NoError(t, err)
		info, err := remote.GetProductInfo(ctx)
		require.NoError(t, err)
		remoteProduct = &info
		if device.RemoteModelID != ActisenseModelUnknown {
			assert.Equal(t, device.RemoteModelID, info.DeviceModelID)
		}
		echoed, err := remote.Echo(ctx, echo)
		require.NoError(t, err)
		assert.Equal(t, echo, echoed)
		require.NotNil(t, session.Status().GatewaySourceAddress)
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
