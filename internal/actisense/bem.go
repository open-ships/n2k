package actisense

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	BEMOperatingMode      = 0x11
	BEMTxPGNEnable        = 0x47
	BEMActivatePGNLists   = 0x4B
	BEMStartupStatus      = 0xF0
	BEMErrorReport        = 0xF1
	BEMSystemStatus       = 0xF2
	BEMNegativeAck        = 0xF4
	maxPendingBEMRequests = 16
)

// OperatingMode is an Actisense gateway operating-mode wire value.
type OperatingMode uint16

const (
	ModeUndefined          OperatingMode = 0
	ModeTransferNormal     OperatingMode = 1
	ModeTransferReceiveAll OperatingMode = 2
	ModeCANPacket          OperatingMode = 5
	ModeCANPacketASCII     OperatingMode = 6
)

// BEMResponse is the common local gateway response envelope. Data is owned.
type BEMResponse struct {
	BSTID        byte
	BEMID        byte
	Sequence     uint8
	ModelID      uint16
	SerialNumber uint32
	ErrorCode    int32
	Data         []byte
}

// DeviceError reports a command that reached the gateway but was rejected.
type DeviceError struct {
	Command byte
	Code    int32
}

func (e *DeviceError) Error() string {
	return fmt.Sprintf("actisense: BEM command 0x%02X failed with device error %d", e.Command, e.Code)
}

// DecodeBEMResponse decodes the 12-byte response header carried by local BEM
// response BST IDs. ok is false when d is not a BEM response datagram.
func DecodeBEMResponse(d Datagram) (response BEMResponse, ok bool, err error) {
	if d.ID != 0xA0 && d.ID != 0xA2 && d.ID != 0xA3 && d.ID != 0xA5 {
		return BEMResponse{}, false, nil
	}
	if len(d.Payload) < 12 {
		return BEMResponse{}, true, fmt.Errorf("actisense: BEM response payload is %d bytes; header requires 12", len(d.Payload))
	}
	return BEMResponse{
		BSTID:        d.ID,
		BEMID:        d.Payload[0],
		Sequence:     d.Payload[1],
		ModelID:      binary.LittleEndian.Uint16(d.Payload[2:4]),
		SerialNumber: binary.LittleEndian.Uint32(d.Payload[4:8]),
		ErrorCode:    int32(binary.LittleEndian.Uint32(d.Payload[8:12])),
		Data:         append([]byte(nil), d.Payload[12:]...),
	}, true, nil
}

// EncodeBEMRequest renders a local BEM command using BST-A1.
func EncodeBEMRequest(command byte, data []byte) ([]byte, error) {
	payload := make([]byte, 1, 1+len(data))
	payload[0] = command
	payload = append(payload, data...)
	return EncodeDatagram(BSTBEMCommand, payload)
}

func OperatingModeRequest() []byte { return nil }

func OperatingModeSet(mode OperatingMode) []byte {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, uint16(mode))
	return data
}

func DecodeOperatingMode(response BEMResponse) (OperatingMode, error) {
	if response.BEMID != BEMOperatingMode {
		return 0, fmt.Errorf("actisense: expected operating-mode response, got BEM 0x%02X", response.BEMID)
	}
	if len(response.Data) < 2 {
		return 0, errors.New("actisense: operating-mode response is shorter than two bytes")
	}
	return OperatingMode(binary.LittleEndian.Uint16(response.Data[:2])), nil
}

func TxPGNEnableGet(pgn uint32) []byte {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, pgn)
	return data
}

func TxPGNEnableSet(pgn uint32, enabled bool) []byte {
	data := make([]byte, 5)
	binary.LittleEndian.PutUint32(data, pgn)
	if enabled {
		data[4] = 1
	}
	return data
}

// TxPGNState is the complete 14-byte Tx PGN enable response.
type TxPGNState struct {
	PGN      uint32
	Enabled  uint8
	Rate     uint32
	Timeout  uint32
	Priority uint8
}

func DecodeTxPGNState(response BEMResponse) (TxPGNState, error) {
	if response.BEMID != BEMTxPGNEnable {
		return TxPGNState{}, fmt.Errorf("actisense: expected Tx PGN response, got BEM 0x%02X", response.BEMID)
	}
	if len(response.Data) < 14 {
		return TxPGNState{}, fmt.Errorf("actisense: Tx PGN response is %d bytes; expected 14", len(response.Data))
	}
	return TxPGNState{
		PGN:      binary.LittleEndian.Uint32(response.Data[:4]),
		Enabled:  response.Data[4],
		Rate:     binary.LittleEndian.Uint32(response.Data[5:9]),
		Timeout:  binary.LittleEndian.Uint32(response.Data[9:13]),
		Priority: response.Data[13],
	}, nil
}

// DiagnosticKind identifies an unsolicited local BEM event.
type DiagnosticKind string

const (
	DiagnosticStartup     DiagnosticKind = "startup_status"
	DiagnosticError       DiagnosticKind = "error_report"
	DiagnosticSystem      DiagnosticKind = "system_status"
	DiagnosticNegativeAck DiagnosticKind = "negative_ack"
)

type StartupStatus struct {
	Legacy      bool
	StartupMode uint16
	ErrorCode   uint32
}

type ErrorReport struct {
	Variant     uint32
	ErrorCode   uint32
	Timestamp   *uint32
	ContextData []byte
}

type IndividualBufferStatus struct {
	ReceiveBandwidth  uint8
	ReceiveLoading    uint8
	ReceiveFiltered   uint8
	ReceiveDropped    uint8
	TransmitBandwidth uint8
	TransmitLoading   uint8
}

type UnifiedBufferStatus struct {
	Bandwidth      uint8
	Deleted        uint8
	Loading        uint8
	PointerLoading uint8
}

type CANStatus struct {
	ReceiveErrors  uint8
	TransmitErrors uint8
	Flags          uint8
}

type SystemStatus struct {
	Individual    []IndividualBufferStatus
	Unified       []UnifiedBufferStatus
	CAN           *CANStatus
	OperatingMode *OperatingMode
}

type NegativeAck struct {
	UniqueCommandID uint32
	ErrorCode       int32
}

// Diagnostic retains the common response origin alongside one typed payload.
type Diagnostic struct {
	Kind        DiagnosticKind
	Response    BEMResponse
	Startup     *StartupStatus
	ErrorReport *ErrorReport
	System      *SystemStatus
	NegativeAck *NegativeAck
}

// DecodeDiagnostic decodes F0/F1/F2/F4. ok is false for ordinary solicited
// responses.
func DecodeDiagnostic(response BEMResponse) (diagnostic Diagnostic, ok bool, err error) {
	diagnostic.Response = response
	switch response.BEMID {
	case BEMStartupStatus:
		status, decodeErr := decodeStartupStatus(response.Data)
		diagnostic.Kind, diagnostic.Startup = DiagnosticStartup, &status
		return diagnostic, true, decodeErr
	case BEMErrorReport:
		report, decodeErr := decodeErrorReport(response.Data)
		diagnostic.Kind, diagnostic.ErrorReport = DiagnosticError, &report
		return diagnostic, true, decodeErr
	case BEMSystemStatus:
		status, decodeErr := decodeSystemStatus(response.Data)
		diagnostic.Kind, diagnostic.System = DiagnosticSystem, &status
		return diagnostic, true, decodeErr
	case BEMNegativeAck:
		if len(response.Data) < 4 {
			return Diagnostic{}, true, errors.New("actisense: negative acknowledgement is shorter than four bytes")
		}
		nack := NegativeAck{UniqueCommandID: binary.LittleEndian.Uint32(response.Data[:4]), ErrorCode: response.ErrorCode}
		diagnostic.Kind, diagnostic.NegativeAck = DiagnosticNegativeAck, &nack
		return diagnostic, true, nil
	default:
		return Diagnostic{}, false, nil
	}
}

func decodeStartupStatus(data []byte) (StartupStatus, error) {
	if len(data) >= 6 {
		return StartupStatus{
			StartupMode: binary.LittleEndian.Uint16(data[:2]),
			ErrorCode:   binary.LittleEndian.Uint32(data[2:6]),
		}, nil
	}
	if len(data) >= 3 {
		return StartupStatus{
			Legacy:      true,
			StartupMode: uint16(data[0]),
			ErrorCode:   uint32(binary.LittleEndian.Uint16(data[1:3])),
		}, nil
	}
	return StartupStatus{}, errors.New("actisense: startup status is shorter than three bytes")
}

func decodeErrorReport(data []byte) (ErrorReport, error) {
	if len(data) < 4 {
		return ErrorReport{}, errors.New("actisense: error report is shorter than four bytes")
	}
	report := ErrorReport{Variant: binary.LittleEndian.Uint32(data[:4])}
	if len(data) >= 8 {
		report.ErrorCode = binary.LittleEndian.Uint32(data[4:8])
	}
	switch report.Variant {
	case 3:
		if len(data) >= 12 {
			timestamp := binary.LittleEndian.Uint32(data[8:12])
			report.Timestamp = &timestamp
		}
	case 2:
		if len(data) > 8 {
			report.ContextData = append([]byte(nil), data[8:]...)
		}
	default:
		if len(data) > 8 {
			report.ContextData = append([]byte(nil), data[8:]...)
		}
	}
	return report, nil
}

func decodeSystemStatus(data []byte) (SystemStatus, error) {
	if len(data) < 1 {
		return SystemStatus{}, errors.New("actisense: system status has no buffer count")
	}
	offset := 1
	individualCount := int(data[0])
	if individualCount < 1 || individualCount > 16 {
		return SystemStatus{}, fmt.Errorf("actisense: system status individual buffer count %d is outside 1-16", individualCount)
	}
	if offset+individualCount*6 > len(data) {
		return SystemStatus{}, errors.New("actisense: system status is truncated in individual buffers")
	}
	status := SystemStatus{Individual: make([]IndividualBufferStatus, 0, individualCount)}
	for range individualCount {
		status.Individual = append(status.Individual, IndividualBufferStatus{
			ReceiveBandwidth:  data[offset],
			ReceiveLoading:    data[offset+1],
			ReceiveFiltered:   data[offset+2],
			ReceiveDropped:    data[offset+3],
			TransmitBandwidth: data[offset+4],
			TransmitLoading:   data[offset+5],
		})
		offset += 6
	}
	if offset == len(data) {
		return status, nil
	}
	unifiedCount := int(data[offset])
	offset++
	if unifiedCount > 8 {
		return SystemStatus{}, fmt.Errorf("actisense: system status unified buffer count %d exceeds 8", unifiedCount)
	}
	if offset+unifiedCount*4 > len(data) {
		return SystemStatus{}, errors.New("actisense: system status is truncated in unified buffers")
	}
	status.Unified = make([]UnifiedBufferStatus, 0, unifiedCount)
	for range unifiedCount {
		status.Unified = append(status.Unified, UnifiedBufferStatus{
			Bandwidth:      data[offset],
			Deleted:        data[offset+1],
			Loading:        data[offset+2],
			PointerLoading: data[offset+3],
		})
		offset += 4
	}
	if len(data)-offset >= 3 {
		status.CAN = &CANStatus{ReceiveErrors: data[offset], TransmitErrors: data[offset+1], Flags: data[offset+2]}
		offset += 3
	}
	if len(data)-offset >= 2 {
		mode := OperatingMode(binary.LittleEndian.Uint16(data[offset : offset+2]))
		status.OperatingMode = &mode
	}
	return status, nil
}
