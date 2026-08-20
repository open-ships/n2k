package n2k

import "github.com/open-ships/n2k/internal/actisense"

type ActisenseOperatingMode = actisense.OperatingMode

const (
	ActisenseModeUndefined          = actisense.ModeUndefined
	ActisenseModeTransferNormal     = actisense.ModeTransferNormal
	ActisenseModeTransferReceiveAll = actisense.ModeTransferReceiveAll
	ActisenseModeTransferLegacyRaw  = actisense.ModeTransferLegacyRaw
	ActisenseModeConvertNormal      = actisense.ModeConvertNormal
	ActisenseModeCANPacket          = actisense.ModeCANPacket
	ActisenseModeCANPacketASCII     = actisense.ModeCANPacketASCII
	ActisenseModeBuffer1            = actisense.ModeBuffer1
	ActisenseModeBuffer2            = actisense.ModeBuffer2
	ActisenseModeBuffer3            = actisense.ModeBuffer3
	ActisenseModeAutoswitchDirect   = actisense.ModeAutoswitchDirect
	ActisenseModeAutoswitchSmart    = actisense.ModeAutoswitchSmart
	ActisenseModeCombineSlow        = actisense.ModeCombineSlow
	ActisenseModeCombineFast        = actisense.ModeCombineFast
	ActisenseModeTest1              = actisense.ModeTest1
	ActisenseModeNSI1               = actisense.ModeNSI1
	ActisenseModeNormal             = actisense.ModeNormal
	ActisenseModePredefined1        = actisense.ModePredefined1
	ActisenseModePredefined2        = actisense.ModePredefined2
	ActisenseModeUser1              = actisense.ModeUser1
	ActisenseModeUser2              = actisense.ModeUser2
	ActisenseModeUser3              = actisense.ModeUser3
	ActisenseModeUser4              = actisense.ModeUser4
	ActisenseModeUser5              = actisense.ModeUser5
	ActisenseModeNull               = actisense.ModeNull
)

type ActisenseBEMPath = actisense.BEMPath
type ActisenseBEMOrigin = actisense.BEMOrigin
type ActisenseBEMResponse = actisense.BEMResponse
type ActisenseDeviceError = actisense.DeviceError
type ActisenseNegativeAckError = actisense.NegativeAckError

const (
	ActisenseBEMLocal  = actisense.BEMPathLocal
	ActisenseBEMRemote = actisense.BEMPathRemote
)

type ActisenseModelID = actisense.ModelID

const (
	ActisenseModelUnknown = actisense.ModelUnknown
	ActisenseModelNGT1    = actisense.ModelNGT1
	ActisenseModelNGT1USB = actisense.ModelNGT1USB
	ActisenseModelNGW1    = actisense.ModelNGW1
	ActisenseModelEMU1    = actisense.ModelEMU1
	ActisenseModelPRONDC1 = actisense.ModelPRONDC1
	ActisenseModelWGX1    = actisense.ModelWGX1
	ActisenseModelNGX1    = actisense.ModelNGX1
)

type ActisenseDeviceCapabilities = actisense.DeviceCapabilities
type ActisenseProtocolMetrics = actisense.SessionMetrics

// ActisenseSessionMetrics remains cumulative across reconnect epochs.
type ActisenseSessionMetrics struct {
	ConnectionEpochs uint64
	Reconnects       uint64
	GatewayResets    uint64
	Protocol         ActisenseProtocolMetrics
}
type ActisenseProductInfo = actisense.ProductInfo
type ActisenseHardwareProtocol = actisense.HardwareProtocol
type ActisensePortMedia = actisense.PortMedia
type ActisensePortInventoryEntry = actisense.PortInventoryEntry
type ActisensePortInventory = actisense.PortInventory
type ActisensePortBaudrate = actisense.PortBaudrate
type ActisensePortPCode = actisense.PortPCode
type ActisenseCANConfig = actisense.CANConfig
type ActisenseCANInfoField = actisense.CANInfoField

const (
	ActisenseHardwareSerialNMEA0183   = actisense.HardwareSerialNMEA0183
	ActisenseHardwareSerialBST        = actisense.HardwareSerialBST
	ActisenseHardwareCANNMEA2000      = actisense.HardwareCANNMEA2000
	ActisenseHardwareCANJ1939         = actisense.HardwareCANJ1939
	ActisenseHardwareEthernetBST      = actisense.HardwareEthernetBST
	ActisenseHardwareEthernetNMEA0183 = actisense.HardwareEthernetNMEA0183
	ActisenseHardwareEthernetOneNet   = actisense.HardwareEthernetOneNet

	ActisensePortMediaCAN      = actisense.PortMediaCAN
	ActisensePortMediaUART     = actisense.PortMediaUART
	ActisensePortMediaUSB      = actisense.PortMediaUSB
	ActisensePortMediaBLE      = actisense.PortMediaBLE
	ActisensePortMediaWiFi     = actisense.PortMediaWiFi
	ActisensePortMediaEthernet = actisense.PortMediaEthernet
	ActisensePortMediaIPStream = actisense.PortMediaIPStream
	ActisensePortMediaUnknown  = actisense.PortMediaUnknown
	ActisensePortIndexNone     = actisense.PortIndexNone

	ActisenseBaudRateNoChange       = actisense.BaudRateNoChange
	ActisenseBaudRateDefault        = actisense.BaudRateDefault
	ActisenseBaudRateAdoptAlternate = actisense.BaudRateAdoptAlternate

	ActisensePortPCodeOff      = actisense.PortPCodeOff
	ActisensePortPCodeOn       = actisense.PortPCodeOn
	ActisensePortPCodeNoChange = actisense.PortPCodeNoChange

	ActisenseCANInfoInstallationDescription1 = actisense.CANInfoInstallationDescription1
	ActisenseCANInfoInstallationDescription2 = actisense.CANInfoInstallationDescription2
	ActisenseCANInfoManufacturerInformation  = actisense.CANInfoManufacturerInformation
)

type ActisensePGNEnableFlag = actisense.PGNEnableFlag
type ActisenseRxPGNState = actisense.RxPGNState
type ActisenseTxPGNState = actisense.TxPGNState
type ActisensePGNListSelector = actisense.PGNListSelector
type ActisensePGNListParameters = actisense.PGNListParameters
type ActisenseSupportedPGN = actisense.SupportedPGN
type ActisenseSupportedPGNList = actisense.SupportedPGNList
type ActisenseRxPGNListEntry = actisense.RxPGNListEntry
type ActisenseTxPGNListEntry = actisense.TxPGNListEntry
type ActisenseProprietaryPGNList = actisense.ProprietaryPGNList
type ActisenseRxPGNEnableList = actisense.RxPGNEnableList
type ActisenseTxPGNEnableList = actisense.TxPGNEnableList

const (
	ActisensePGNDisabled        = actisense.PGNDisabled
	ActisensePGNEnabled         = actisense.PGNEnabled
	ActisensePGNRespondMode     = actisense.PGNRespondMode
	ActisenseRxPGNMaskAcceptAll = actisense.RxPGNMaskAcceptAll
	ActisenseTxPGNRateDefault   = actisense.TxPGNRateDefault
	ActisenseTxPGNRateEvent     = actisense.TxPGNRateEvent
	ActisensePGNListRx          = actisense.PGNListRx
	ActisensePGNListTx          = actisense.PGNListTx
	ActisensePGNListBoth        = actisense.PGNListBoth
)

type ActisenseDiagnosticKind = actisense.DiagnosticKind
type ActisenseDiagnostic = actisense.Diagnostic
type ActisenseStartupStatus = actisense.StartupStatus
type ActisenseErrorReport = actisense.ErrorReport
type ActisenseSystemStatus = actisense.SystemStatus
type ActisenseNegativeAck = actisense.NegativeAck

const (
	ActisenseDiagnosticStartup     = actisense.DiagnosticStartup
	ActisenseDiagnosticError       = actisense.DiagnosticError
	ActisenseDiagnosticSystem      = actisense.DiagnosticSystem
	ActisenseDiagnosticNegativeAck = actisense.DiagnosticNegativeAck
)

// ActisenseDevice exposes the same typed command Interface for a locally
// attached gateway and a remote Actisense device.
type ActisenseDevice struct {
	*actisense.CommandSet
}
