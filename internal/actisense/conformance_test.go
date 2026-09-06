package actisense

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type actisenseGoldenVector struct {
	ID      string `json:"id"`
	Command byte   `json:"command"`
	Data    string `json:"data_hex"`
	BDTP    string `json:"bdtp_hex"`
}

type actisenseGoldenCorpus struct {
	Baseline struct {
		SDKCommit string `json:"sdk_commit"`
	} `json:"baseline"`
	Requests            []actisenseGoldenVector `json:"local_bem_requests"`
	Responses           []actisenseGoldenVector `json:"local_bem_responses"`
	DocumentedRequests  []actisenseGoldenVector `json:"documented_bem_requests"`
	DocumentedResponses []actisenseGoldenVector `json:"documented_bem_responses"`
}

func TestActisenseGoldenDocumentedCommands(t *testing.T) {
	corpus := loadActisenseGoldenCorpus(t)
	for _, vector := range corpus.DocumentedRequests {
		wire, err := EncodeBEMRequest(vector.Command, decodeGoldenHex(t, vector.Data))
		require.NoError(t, err)
		assert.Equal(t, decodeGoldenHex(t, vector.BDTP), wire, vector.ID)
	}
	responses := make(map[byte][]BEMResponse)
	for _, vector := range corpus.DocumentedResponses {
		count := 0
		NewParser().Feed(decodeGoldenHex(t, vector.BDTP), func(datagram Datagram) {
			response, ok, err := DecodeBEMResponse(datagram)
			require.NoError(t, err)
			require.True(t, ok)
			assert.Equal(t, vector.Command, response.BEMID)
			assert.Equal(t, decodeGoldenHex(t, vector.Data), response.Data)
			responses[response.BEMID] = append(responses[response.BEMID], response)
			count++
		}, func(err DecodeError) { t.Error(err) })
		require.Equal(t, 1, count, vector.ID)
	}
	require.Len(t, responses[BEMPortDuplicateDelete], 1)
	r := &scriptedRequester{response: responses[BEMPortDuplicateDelete][0]}
	commands := NewCommandSet(r, CommandSetConfig{})
	ports, err := commands.GetPortDuplicateDelete(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []PortDuplicateDelete{0, 1, 0}, ports)
	r.responses = responses[BEMRxPGNEnableListF1]
	rx, err := commands.GetRxPGNEnableListF1(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 2, rx.PartsReceived)
	assert.Equal(t, []RxPGNListF1Entry{{127250, RxPGNMaskPGN}, {129025, RxPGNMaskDataPage}}, rx.Entries)
	r.responses = responses[BEMTxPGNEnableListF1]
	tx, err := commands.GetTxPGNEnableListF1(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 4, tx.PartsReceived)
	assert.Equal(t, []TxPGNListF1Entry{{127250, 1000, 2000, 3}, {129025, 500, 0, 6}}, tx.Entries)
}

func loadActisenseGoldenCorpus(t *testing.T) actisenseGoldenCorpus {
	t.Helper()
	contents, err := os.ReadFile("../../conformance/actisense-golden.json")
	require.NoError(t, err)
	var corpus actisenseGoldenCorpus
	require.NoError(t, json.Unmarshal(contents, &corpus))
	assert.Equal(t, "ed2268a6e8db0645f75e4ef17eed2e937d025040", corpus.Baseline.SDKCommit)
	return corpus
}

func decodeGoldenHex(t *testing.T, value string) []byte {
	t.Helper()
	decoded, err := hex.DecodeString(value)
	require.NoError(t, err)
	return decoded
}

func compiledSolicitedBEMCommands() map[byte]struct{} {
	return map[byte]struct{}{
		BEMReinitialize: {}, BEMCommitEEPROM: {}, BEMCommitFlash: {}, BEMOperatingMode: {},
		BEMPortPCode: {}, BEMTotalTime: {}, BEMPortBaudrate: {}, BEMEcho: {}, BEMPortInventory: {},
		BEMSupportedPGNList: {}, BEMProductInfo: {}, BEMCANConfig: {}, BEMCANInfoField1: {},
		BEMCANInfoField2: {}, BEMCANInfoField3: {}, BEMRxPGNEnable: {}, BEMTxPGNEnable: {},
		BEMDeletePGNLists: {}, BEMActivatePGNLists: {}, BEMDefaultPGNLists: {},
		BEMPGNListParameters: {}, BEMRxPGNEnableListF2: {}, BEMTxPGNEnableListF2: {},
	}
}

func TestActisenseGoldenCorpusCoversEveryCompiledSolicitedCommand(t *testing.T) {
	corpus := loadActisenseGoldenCorpus(t)
	requestCoverage := make(map[byte]struct{})
	for _, vector := range corpus.Requests {
		requestCoverage[vector.Command] = struct{}{}
	}
	responseCoverage := make(map[byte]struct{})
	for _, vector := range corpus.Responses {
		responseCoverage[vector.Command] = struct{}{}
	}
	assert.Equal(t, compiledSolicitedBEMCommands(), requestCoverage)
	assert.Equal(t, compiledSolicitedBEMCommands(), responseCoverage)
}

func TestActisenseGoldenLocalBEMRequests(t *testing.T) {
	for _, vector := range loadActisenseGoldenCorpus(t).Requests {
		t.Run(vector.ID, func(t *testing.T) {
			wire, err := EncodeBEMRequest(vector.Command, decodeGoldenHex(t, vector.Data))
			require.NoError(t, err)
			assert.Equal(t, decodeGoldenHex(t, vector.BDTP), wire)
		})
	}
}

func TestActisenseGoldenLocalBEMResponses(t *testing.T) {
	for _, vector := range loadActisenseGoldenCorpus(t).Responses {
		t.Run(vector.ID, func(t *testing.T) {
			var datagrams []Datagram
			var decodeErrors []DecodeError
			parser := NewParser()
			parser.Feed(decodeGoldenHex(t, vector.BDTP), func(datagram Datagram) {
				datagrams = append(datagrams, datagram)
			}, func(err DecodeError) {
				decodeErrors = append(decodeErrors, err)
			})
			require.Empty(t, decodeErrors)
			require.Len(t, datagrams, 1)
			response, ok, err := DecodeBEMResponse(datagrams[0])
			require.NoError(t, err)
			require.True(t, ok)
			assert.Equal(t, vector.Command, response.BEMID)
			assert.Equal(t, uint8(1), response.Sequence)
			assert.Equal(t, uint16(ModelNGT1), response.ModelID)
			assert.Equal(t, uint32(0x11223344), response.SerialNumber)
			assert.Zero(t, response.ErrorCode)
			assert.True(t, bytes.Equal(decodeGoldenHex(t, vector.Data), response.Data))
			assertGoldenTypedResponse(t, response)
		})
	}
}

func assertGoldenTypedResponse(t *testing.T, response BEMResponse) {
	t.Helper()
	switch response.BEMID {
	case BEMReinitialize, BEMCommitEEPROM, BEMCommitFlash, BEMDeletePGNLists, BEMActivatePGNLists, BEMDefaultPGNLists:
		assert.Empty(t, response.Data)
	case BEMOperatingMode:
		mode, err := DecodeOperatingMode(response)
		require.NoError(t, err)
		assert.Equal(t, ModeCANPacket, mode)
	case BEMPortPCode:
		codes, err := DecodePortPCodes(response)
		require.NoError(t, err)
		assert.Equal(t, []PortPCode{PortPCodeOff, PortPCodeOn, PortPCodeNoChange}, codes)
	case BEMTotalTime:
		seconds, err := DecodeTotalTime(response)
		require.NoError(t, err)
		assert.Equal(t, uint32(0x12345678), seconds)
	case BEMPortBaudrate:
		state, err := DecodePortBaudrate(response)
		require.NoError(t, err)
		assert.Equal(t, uint32(115200), state.SessionBaud)
	case BEMEcho:
		payload, err := DecodeEcho(response)
		require.NoError(t, err)
		assert.Equal(t, []byte{1, 2, 3}, payload)
	case BEMPortInventory:
		accumulator := &portInventoryAccumulator{}
		done, err := accumulator.feed(response)
		require.NoError(t, err)
		assert.True(t, done)
		require.Len(t, accumulator.result.Ports, 1)
		assert.Equal(t, "CAN0", accumulator.result.Ports[0].Name)
	case BEMSupportedPGNList:
		part, err := decodeSupportedPGNPart(response)
		require.NoError(t, err)
		done, err := (&supportedPGNAccumulator{}).feed(part)
		require.NoError(t, err)
		assert.True(t, done)
		require.Len(t, part.entries, 1)
		assert.Equal(t, uint32(126992), part.entries[0].PGN)
	case BEMProductInfo:
		accumulator := &productInfoAccumulator{}
		done, err := accumulator.feed(response)
		require.NoError(t, err)
		assert.False(t, done)
		assert.Equal(t, uint16(101), accumulator.result.ProductCode)
	case BEMCANConfig:
		config, err := DecodeCANConfig(response)
		require.NoError(t, err)
		assert.Equal(t, uint64(0x0102030405060708), config.NAME)
		assert.Equal(t, uint8(42), config.SourceAddress)
	case BEMCANInfoField1, BEMCANInfoField2, BEMCANInfoField3:
		field := CANInfoField(response.BEMID - BEMCANInfoField1 + 1)
		value, err := DecodeCANInfoField(response, field)
		require.NoError(t, err)
		assert.Equal(t, "Hi", value)
	case BEMRxPGNEnable:
		state, err := DecodeRxPGNState(response)
		require.NoError(t, err)
		assert.Equal(t, uint32(126992), state.PGN)
		assert.Equal(t, uint32(0x03FFFF00), state.Mask)
	case BEMTxPGNEnable:
		state, err := DecodeTxPGNState(response)
		require.NoError(t, err)
		assert.Equal(t, uint32(126992), state.PGN)
		assert.Equal(t, uint32(1000), state.Rate)
		assert.Equal(t, uint8(3), state.Priority)
	case BEMPGNListParameters:
		parameters, err := DecodePGNListParameters(response)
		require.NoError(t, err)
		assert.Equal(t, uint16(10), parameters.RxMaximum)
		assert.Equal(t, uint16(5), parameters.TxActive)
	case BEMRxPGNEnableListF2:
		part, err := decodeRxF2Part(response)
		require.NoError(t, err)
		require.Len(t, part.entries, 1)
		assert.Equal(t, RxPGNListEntry{Index: 3, Mask: 1}, part.entries[0])
	case BEMTxPGNEnableListF2:
		part, err := decodeTxF2Part(response)
		require.NoError(t, err)
		require.Len(t, part.entries, 1)
		assert.Equal(t, TxPGNListEntry{Index: 7, Priority: 6, RateMS: 1000}, part.entries[0])
	default:
		t.Fatalf("unclassified golden response BEM 0x%02X", response.BEMID)
	}
}
