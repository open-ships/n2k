package n2k

import "github.com/open-ships/n2k/internal/actisense"

// ActisenseOperatingMode selects the gateway operating mode; use the ActisenseMode constants.
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

// ActisenseBEMPath distinguishes local gateway commands from remote bus commands.
type ActisenseBEMPath = actisense.BEMPath

// ActisenseBEMOrigin identifies the local or remote source of a BEM response.
type ActisenseBEMOrigin = actisense.BEMOrigin

// ActisenseBEMResponse is an owned decoded BEM response with command, origin, device identity, and data.
type ActisenseBEMResponse = actisense.BEMResponse

// ActisenseDeviceError reports a device error code for a BEM command; use errors.As to inspect it.
type ActisenseDeviceError = actisense.DeviceError

// ActisenseNegativeAckError reports an asynchronous negative acknowledgment of a BEM command.
type ActisenseNegativeAckError = actisense.NegativeAckError

const (
	ActisenseBEMLocal  = actisense.BEMPathLocal
	ActisenseBEMRemote = actisense.BEMPathRemote
)

// ActisenseModelID identifies an Actisense hardware model from a BEM response.
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

// ActisenseDeviceCapabilities describes known protocol capabilities and model-specific restrictions.
type ActisenseDeviceCapabilities = actisense.DeviceCapabilities

// ActisenseProtocolMetrics is an owned snapshot of cumulative transport and BEM counters.
type ActisenseProtocolMetrics = actisense.SessionMetrics

// ActisenseSessionMetrics remains cumulative across reconnect epochs.
type ActisenseSessionMetrics struct {
	ConnectionEpochs uint64
	Reconnects       uint64
	GatewayResets    uint64
	Protocol         ActisenseProtocolMetrics
}

// ActisenseProductInfo contains the product identity returned by an Actisense device.
type ActisenseProductInfo = actisense.ProductInfo

// ActisenseHardwareProtocol identifies the protocol supported by a hardware port.
type ActisenseHardwareProtocol = actisense.HardwareProtocol

// ActisensePortMedia identifies a port medium, such as CAN, serial, or Ethernet.
type ActisensePortMedia = actisense.PortMedia

// ActisensePortInventoryEntry describes one hardware port in a device inventory.
type ActisensePortInventoryEntry = actisense.PortInventoryEntry

// ActisensePortInventory contains the complete decoded device port inventory.
type ActisensePortInventory = actisense.PortInventory

// ActisensePortBaudrate contains a port index and its baud-rate configuration.
type ActisensePortBaudrate = actisense.PortBaudrate

// ActisensePortPCode contains a port index and its proprietary-code setting.
type ActisensePortPCode = actisense.PortPCode

// ActisenseCANConfig contains stored CAN identity configuration; it does not prove the live claimed address.
type ActisenseCANConfig = actisense.CANConfig

// ActisenseCANInfoField selects a CAN installation-description or manufacturer-information field.
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

// ActisensePGNEnableFlag selects disabled, enabled, or request-response PGN behavior.
type ActisensePGNEnableFlag = actisense.PGNEnableFlag

// ActisenseRxPGNState contains one receive PGN entry and its effective flag and mask.
type ActisenseRxPGNState = actisense.RxPGNState

// ActisenseTxPGNState contains one transmit PGN entry, effective flag, rate, timeout, and priority.
type ActisenseTxPGNState = actisense.TxPGNState

// ActisensePGNListSelector selects the receive list, transmit list, or both lists.
type ActisensePGNListSelector = actisense.PGNListSelector

// ActisensePGNListParameters describes the gateway PGN-list configuration.
type ActisensePGNListParameters = actisense.PGNListParameters

// ActisenseSupportedPGN describes one PGN supported by the gateway.
type ActisenseSupportedPGN = actisense.SupportedPGN

// ActisenseSupportedPGNList contains an assembled list of gateway-supported PGNs.
type ActisenseSupportedPGNList = actisense.SupportedPGNList

// ActisenseRxPGNListEntry describes one entry in a decoded F2 receive list.
type ActisenseRxPGNListEntry = actisense.RxPGNListEntry

// ActisenseTxPGNListEntry describes one entry in a decoded F2 transmit list.
type ActisenseTxPGNListEntry = actisense.TxPGNListEntry

// ActisenseProprietaryPGNList contains an assembled proprietary PGN list.
type ActisenseProprietaryPGNList = actisense.ProprietaryPGNList

// ActisenseRxPGNEnableList contains a bounded, assembled F2 receive-enable list.
type ActisenseRxPGNEnableList = actisense.RxPGNEnableList

// ActisenseTxPGNEnableList contains a bounded, assembled F2 transmit-enable list.
type ActisenseTxPGNEnableList = actisense.TxPGNEnableList

// ActisenseRxPGNListF1Entry describes one entry in a legacy F1 receive-enable list.
type ActisenseRxPGNListF1Entry = actisense.RxPGNListF1Entry

// ActisenseTxPGNListF1Entry describes one entry in a legacy F1 transmit-enable list.
type ActisenseTxPGNListF1Entry = actisense.TxPGNListF1Entry

// ActisenseRxPGNEnableListF1 contains an assembled legacy F1 receive-enable list.
type ActisenseRxPGNEnableListF1 = actisense.RxPGNEnableListF1

// ActisenseTxPGNEnableListF1 contains an assembled legacy F1 transmit-enable list.
type ActisenseTxPGNEnableListF1 = actisense.TxPGNEnableListF1

// ActisensePortDuplicateDelete selects the port duplicate-suppression setting.
type ActisensePortDuplicateDelete = actisense.PortDuplicateDelete

const (
	ActisensePortDuplicateDeleteOff      = actisense.PortDuplicateDeleteOff
	ActisensePortDuplicateDeleteOn       = actisense.PortDuplicateDeleteOn
	ActisensePortDuplicateDeleteNoChange = actisense.PortDuplicateDeleteNoChange
)

const (
	ActisensePGNDisabled        = actisense.PGNDisabled
	ActisensePGNEnabled         = actisense.PGNEnabled
	ActisensePGNRespondMode     = actisense.PGNRespondMode
	ActisenseRxPGNMaskPGN       = actisense.RxPGNMaskPGN
	ActisenseRxPGNMaskPDUFormat = actisense.RxPGNMaskPDUFormat
	ActisenseRxPGNMaskPDUNibble = actisense.RxPGNMaskPDUNibble
	ActisenseRxPGNMaskDataPage  = actisense.RxPGNMaskDataPage
	ActisenseRxPGNMaskDefault   = actisense.RxPGNMaskDefault
	ActisenseRxPGNMaskNoChange  = actisense.RxPGNMaskNoChange
	ActisenseTxPGNRateNoChange  = actisense.TxPGNRateNoChange
	ActisenseTxPGNRateEvent     = actisense.TxPGNRateEvent
	ActisensePGNListRx          = actisense.PGNListRx
	ActisensePGNListTx          = actisense.PGNListTx
	ActisensePGNListBoth        = actisense.PGNListBoth
)

// Deprecated: the wire value means no change, not accept all. Use ActisenseRxPGNMaskNoChange.
const ActisenseRxPGNMaskAcceptAll = uint32(0xFFFFFFFF)

// Deprecated: the wire value leaves the rate unchanged. Use ActisenseTxPGNRateNoChange.
const ActisenseTxPGNRateDefault = uint32(0xFFFFFFFF)

// ActisenseDiagnosticKind identifies a startup, error, status, or negative-acknowledgment event.
type ActisenseDiagnosticKind = actisense.DiagnosticKind

// ActisenseDiagnostic is an owned BEM diagnostic with its response and origin.
type ActisenseDiagnostic = actisense.Diagnostic

// ActisenseStartupStatus contains decoded gateway startup information.
type ActisenseStartupStatus = actisense.StartupStatus

// ActisenseErrorReport contains a decoded gateway error report.
type ActisenseErrorReport = actisense.ErrorReport

// ActisenseSystemStatus contains a decoded gateway system-status report.
type ActisenseSystemStatus = actisense.SystemStatus

// ActisenseNegativeAck contains the rejected command identity and device error code.
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
