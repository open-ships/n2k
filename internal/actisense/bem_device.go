package actisense

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"unicode/utf8"
)

const (
	productInfoVariant      = uint32(0x00000011)
	productInfoFormat2Seq   = uint8(6)
	portInventoryVariant    = uint32(0x00001104)
	portInventoryRecordSize = 22
	maxCANInfoFieldLength   = 70
	maxEchoPayloadLength    = 222
)

// ModelID is the Actisense hardware model identifier carried in every BEM
// response header.
type ModelID uint16

const (
	ModelUnknown ModelID = 0x0000
	ModelNGT1    ModelID = 0x000E
	ModelNGT1USB ModelID = 0x000F
	ModelNGW1    ModelID = 0x0010
	ModelEMU1    ModelID = 0x0011
	ModelPRONDC1 ModelID = 0x0020
	ModelWGX1    ModelID = 0x0030
	ModelNGX1    ModelID = 0x003B
)

// DeviceCapabilities records behavior that varies by known Actisense model.
// Unknown future models are deliberately conservative until verified.
type DeviceCapabilities struct {
	ModelID                       ModelID
	ModelName                     string
	ProprietaryPGNEnableListF2    bool
	ReceiveAllOmitsISOControlPGNs bool
	RewritesHostTransmitSID       bool
}

func CapabilitiesForModel(model ModelID) DeviceCapabilities {
	capabilities := DeviceCapabilities{ModelID: model, ModelName: fmt.Sprintf("Model-0x%04X", uint16(model))}
	switch model {
	case ModelNGT1:
		capabilities.ModelName = "NGT-1"
		capabilities.RewritesHostTransmitSID = true
	case ModelNGT1USB:
		capabilities.ModelName = "NGT-1 USB"
		capabilities.RewritesHostTransmitSID = true
	case ModelNGW1:
		capabilities.ModelName = "NGW-1"
	case ModelEMU1:
		capabilities.ModelName = "EMU-1"
	case ModelPRONDC1:
		capabilities.ModelName = "PRO-NDC-1-E2K"
	case ModelWGX1:
		capabilities.ModelName = "WGX"
	case ModelNGX1:
		capabilities.ModelName = "NGX-1"
		capabilities.ProprietaryPGNEnableListF2 = true
		capabilities.ReceiveAllOmitsISOControlPGNs = true
	case ModelUnknown:
		capabilities.ModelName = "Unknown"
	}
	return capabilities
}

// ProductInfo is the complete BEM 0x41 response. Legacy devices send five
// replies; current devices send one Format-2 reply.
type ProductInfo struct {
	StructureVariant   uint32
	NMEA2000Version    uint16
	ProductCode        uint16
	Model              string
	SoftwareVersion    string
	ModelVersion       string
	SerialCode         string
	Certification      uint8
	LoadEquivalency    uint8
	DeviceModelID      ModelID
	DeviceSerialNumber uint32
	Legacy             bool
}

func decodePaddedString(data []byte) string {
	if index := bytes.IndexAny(data, "\x00\xff"); index >= 0 {
		data = data[:index]
	}
	return string(data)
}

func DecodeProductInfo(response BEMResponse) (ProductInfo, error) {
	if response.BEMID != BEMProductInfo {
		return ProductInfo{}, fmt.Errorf("actisense: expected Product Info response, got BEM 0x%02X", response.BEMID)
	}
	if len(response.Data) < 138 {
		return ProductInfo{}, fmt.Errorf("actisense: Product Info Format-2 response is %d bytes; expected 138", len(response.Data))
	}
	variant := binary.LittleEndian.Uint32(response.Data[:4])
	if variant != productInfoVariant {
		return ProductInfo{}, fmt.Errorf("actisense: Product Info structure variant is 0x%08X; expected 0x%08X", variant, productInfoVariant)
	}
	return ProductInfo{
		StructureVariant:   variant,
		NMEA2000Version:    binary.LittleEndian.Uint16(response.Data[4:6]),
		ProductCode:        binary.LittleEndian.Uint16(response.Data[6:8]),
		Model:              decodePaddedString(response.Data[8:40]),
		SoftwareVersion:    decodePaddedString(response.Data[40:72]),
		ModelVersion:       decodePaddedString(response.Data[72:104]),
		SerialCode:         decodePaddedString(response.Data[104:136]),
		Certification:      response.Data[136],
		LoadEquivalency:    response.Data[137],
		DeviceModelID:      ModelID(response.ModelID),
		DeviceSerialNumber: response.SerialNumber,
	}, nil
}

// productInfoAccumulator accepts either the current one-message Format-2
// reply or the legacy five-part response train.
type productInfoAccumulator struct {
	result ProductInfo
	seen   uint8
	next   uint8
	legacy bool
	modern bool
}

func (a *productInfoAccumulator) feed(response BEMResponse) (bool, error) {
	if response.Sequence == productInfoFormat2Seq || (len(response.Data) >= 4 && binary.LittleEndian.Uint32(response.Data[:4]) == productInfoVariant) {
		if a.legacy {
			return false, errors.New("actisense: Product Info response mixed legacy and Format-2 records")
		}
		info, err := DecodeProductInfo(response)
		if err != nil {
			return false, err
		}
		a.result, a.modern = info, true
		return true, nil
	}
	if a.modern {
		return false, errors.New("actisense: Product Info response mixed Format-2 and legacy records")
	}
	a.legacy = true
	part := int(response.Sequence)
	if part < 1 || part > 5 {
		if a.next == 0 {
			a.next = 1
		}
		part = int(a.next)
	}
	if part < 1 || part > 5 {
		return false, fmt.Errorf("actisense: unexpected Product Info legacy part %d", part)
	}
	a.next = uint8(part + 1)
	a.result.DeviceModelID = ModelID(response.ModelID)
	a.result.DeviceSerialNumber = response.SerialNumber
	a.result.Legacy = true
	switch part {
	case 1:
		if len(response.Data) < 6 {
			return false, fmt.Errorf("actisense: Product Info legacy part 1 is %d bytes; expected 6", len(response.Data))
		}
		a.result.NMEA2000Version = binary.LittleEndian.Uint16(response.Data[:2])
		a.result.ProductCode = binary.LittleEndian.Uint16(response.Data[2:4])
		a.result.Certification = response.Data[4]
		a.result.LoadEquivalency = response.Data[5]
	case 2, 3, 4, 5:
		if len(response.Data) < 32 {
			return false, fmt.Errorf("actisense: Product Info legacy part %d is %d bytes; expected 32", part, len(response.Data))
		}
		value := decodePaddedString(response.Data[:32])
		switch part {
		case 2:
			a.result.Model = value
		case 3:
			a.result.SoftwareVersion = value
		case 4:
			a.result.ModelVersion = value
		case 5:
			a.result.SerialCode = value
		}
	}
	a.seen |= 1 << (part - 1)
	return a.seen == 0x1F, nil
}

// HardwareProtocol identifies the protocol carried over a gateway port.
type HardwareProtocol uint8

const (
	HardwareSerialNMEA0183   HardwareProtocol = 0
	HardwareSerialBST        HardwareProtocol = 1
	HardwareCANNMEA2000      HardwareProtocol = 32
	HardwareCANJ1939         HardwareProtocol = 33
	HardwareEthernetBST      HardwareProtocol = 64
	HardwareEthernetNMEA0183 HardwareProtocol = 65
	HardwareEthernetOneNet   HardwareProtocol = 66
)

// PortMedia identifies the physical medium independently from its protocol.
type PortMedia uint8

const (
	PortMediaCAN      PortMedia = 0
	PortMediaUART     PortMedia = 1
	PortMediaUSB      PortMedia = 2
	PortMediaBLE      PortMedia = 3
	PortMediaWiFi     PortMedia = 4
	PortMediaEthernet PortMedia = 5
	PortMediaIPStream PortMedia = 6
	PortMediaUnknown  PortMedia = 0xFF
	PortIndexNone               = uint8(0xFF)
	PortCanReceive              = uint8(0x01)
	PortCanTransmit             = uint8(0x02)
)

// PortInventoryEntry maps the independent inventory, System Status, and Port
// Baudrate index spaces without guessing.
type PortInventoryEntry struct {
	PortIndex          uint8
	SystemStatusIndex  uint8
	BaudratePortNumber uint8
	Media              PortMedia
	Protocol           HardwareProtocol
	Capabilities       uint8
	SessionBaud        uint32
	StoreBaud          uint32
	Name               string
}

func (e PortInventoryEntry) CanReceive() bool         { return e.Capabilities&PortCanReceive != 0 }
func (e PortInventoryEntry) CanTransmit() bool        { return e.Capabilities&PortCanTransmit != 0 }
func (e PortInventoryEntry) HasSystemStatus() bool    { return e.SystemStatusIndex != PortIndexNone }
func (e PortInventoryEntry) HasBaudrateControl() bool { return e.BaudratePortNumber != PortIndexNone }
func (e PortInventoryEntry) HasSessionOverride() bool { return e.SessionBaud != e.StoreBaud }

type PortInventory struct {
	TransferID uint8
	Ports      []PortInventoryEntry
}

type portInventoryPart struct {
	transferID uint8
	total      uint8
	first      uint8
	entries    []PortInventoryEntry
}

func decodePortInventoryPart(response BEMResponse) (portInventoryPart, error) {
	if response.BEMID != BEMPortInventory {
		return portInventoryPart{}, fmt.Errorf("actisense: expected Port Inventory response, got BEM 0x%02X", response.BEMID)
	}
	data := response.Data
	if len(data) < 8 {
		return portInventoryPart{}, fmt.Errorf("actisense: Port Inventory response is %d bytes; expected at least 8", len(data))
	}
	if variant := binary.LittleEndian.Uint32(data[1:5]); variant != portInventoryVariant {
		return portInventoryPart{}, fmt.Errorf("actisense: Port Inventory structure variant is 0x%08X; expected 0x%08X", variant, portInventoryVariant)
	}
	part := portInventoryPart{transferID: data[0], total: data[5], first: data[6]}
	count := int(data[7])
	needed := 8 + count*portInventoryRecordSize
	if len(data) < needed {
		return portInventoryPart{}, fmt.Errorf("actisense: Port Inventory response needs %d bytes for %d records; got %d", needed, count, len(data))
	}
	part.entries = make([]PortInventoryEntry, 0, count)
	for index := 0; index < count; index++ {
		record := data[8+index*portInventoryRecordSize : 8+(index+1)*portInventoryRecordSize]
		part.entries = append(part.entries, PortInventoryEntry{
			PortIndex:          record[0],
			SystemStatusIndex:  record[1],
			BaudratePortNumber: record[2],
			Media:              PortMedia(record[3]),
			Protocol:           HardwareProtocol(record[4]),
			Capabilities:       record[5],
			SessionBaud:        binary.LittleEndian.Uint32(record[6:10]),
			StoreBaud:          binary.LittleEndian.Uint32(record[10:14]),
			Name:               decodePaddedString(record[14:22]),
		})
	}
	return part, nil
}

type portInventoryAccumulator struct {
	result      PortInventory
	seen        []bool
	initialized bool
	received    int
}

func (a *portInventoryAccumulator) feed(response BEMResponse) (bool, error) {
	part, err := decodePortInventoryPart(response)
	if err != nil {
		return false, err
	}
	if !a.initialized {
		a.result.TransferID = part.transferID
		a.result.Ports = make([]PortInventoryEntry, int(part.total))
		a.seen = make([]bool, int(part.total))
		a.initialized = true
	} else if part.transferID != a.result.TransferID {
		return false, fmt.Errorf("actisense: Port Inventory transfer ID changed from %d to %d", a.result.TransferID, part.transferID)
	} else if int(part.total) != len(a.result.Ports) {
		return false, fmt.Errorf("actisense: Port Inventory total changed from %d to %d", len(a.result.Ports), part.total)
	}
	end := int(part.first) + len(part.entries)
	if end > len(a.result.Ports) {
		return false, fmt.Errorf("actisense: Port Inventory range %d:%d exceeds total %d", part.first, end, len(a.result.Ports))
	}
	for index, entry := range part.entries {
		slot := int(part.first) + index
		a.result.Ports[slot] = entry
		if !a.seen[slot] {
			a.seen[slot] = true
			a.received++
		}
	}
	return a.received == len(a.result.Ports), nil
}

const (
	BaudRateNoChange       uint32 = 0xFFFFFFFF
	BaudRateDefault        uint32 = 0xFFFFFFFE
	BaudRateAdoptAlternate uint32 = 0xFFFFFFFC
)

type PortBaudrate struct {
	TotalPorts  uint8
	PortNumber  uint8
	Protocol    HardwareProtocol
	SessionBaud uint32
	StoreBaud   uint32
}

func PortBaudrateGet(port uint8) []byte { return []byte{port} }

func PortBaudrateSet(port uint8, sessionBaud, storeBaud uint32) []byte {
	data := make([]byte, 9)
	data[0] = port
	binary.LittleEndian.PutUint32(data[1:5], sessionBaud)
	binary.LittleEndian.PutUint32(data[5:9], storeBaud)
	return data
}

func DecodePortBaudrate(response BEMResponse) (PortBaudrate, error) {
	if response.BEMID != BEMPortBaudrate {
		return PortBaudrate{}, fmt.Errorf("actisense: expected Port Baudrate response, got BEM 0x%02X", response.BEMID)
	}
	if len(response.Data) < 11 {
		return PortBaudrate{}, fmt.Errorf("actisense: Port Baudrate response is %d bytes; expected 11", len(response.Data))
	}
	return PortBaudrate{
		TotalPorts: response.Data[0], PortNumber: response.Data[1], Protocol: HardwareProtocol(response.Data[2]),
		SessionBaud: binary.LittleEndian.Uint32(response.Data[3:7]), StoreBaud: binary.LittleEndian.Uint32(response.Data[7:11]),
	}, nil
}

type PortPCode uint8

const (
	PortPCodeOff      PortPCode = 0
	PortPCodeOn       PortPCode = 1
	PortPCodeNoChange PortPCode = 0xFF
)

func DecodePortPCodes(response BEMResponse) ([]PortPCode, error) {
	if response.BEMID != BEMPortPCode {
		return nil, fmt.Errorf("actisense: expected Port P-Code response, got BEM 0x%02X", response.BEMID)
	}
	if len(response.Data) == 0 {
		return nil, errors.New("actisense: Port P-Code response has no count byte")
	}
	count := int(response.Data[0])
	if len(response.Data) < 1+count {
		return nil, fmt.Errorf("actisense: Port P-Code response reports %d values; only %d present", count, len(response.Data)-1)
	}
	result := make([]PortPCode, count)
	for index := range count {
		result[index] = PortPCode(response.Data[index+1])
	}
	return result, nil
}

func EncodePortPCodes(codes []PortPCode) ([]byte, error) {
	if len(codes) > 254 {
		return nil, errors.New("actisense: at most 254 Port P-Code values fit in one BEM command")
	}
	data := make([]byte, len(codes))
	for index, code := range codes {
		data[index] = byte(code)
	}
	return data, nil
}

type CANConfig struct {
	NAME          uint64
	SourceAddress uint8
}

func CANConfigSet(config CANConfig) []byte {
	data := make([]byte, 9)
	binary.LittleEndian.PutUint64(data[:8], config.NAME)
	data[8] = config.SourceAddress
	return data
}

func DecodeCANConfig(response BEMResponse) (CANConfig, error) {
	if response.BEMID != BEMCANConfig {
		return CANConfig{}, fmt.Errorf("actisense: expected CAN Config response, got BEM 0x%02X", response.BEMID)
	}
	if len(response.Data) < 9 {
		return CANConfig{}, fmt.Errorf("actisense: CAN Config response is %d bytes; expected 9", len(response.Data))
	}
	return CANConfig{NAME: binary.LittleEndian.Uint64(response.Data[:8]), SourceAddress: response.Data[8]}, nil
}

type CANInfoField uint8

const (
	CANInfoInstallationDescription1 CANInfoField = 1
	CANInfoInstallationDescription2 CANInfoField = 2
	CANInfoManufacturerInformation  CANInfoField = 3
)

func (f CANInfoField) command() (byte, error) {
	switch f {
	case CANInfoInstallationDescription1:
		return BEMCANInfoField1, nil
	case CANInfoInstallationDescription2:
		return BEMCANInfoField2, nil
	case CANInfoManufacturerInformation:
		return BEMCANInfoField3, nil
	default:
		return 0, fmt.Errorf("actisense: unknown CAN Info field %d", f)
	}
}

func EncodeCANInfoField(text string) ([]byte, error) {
	if !utf8.ValidString(text) {
		return nil, errors.New("actisense: CAN Info field is not valid UTF-8")
	}
	if len(text) > maxCANInfoFieldLength {
		return nil, fmt.Errorf("actisense: CAN Info field is %d bytes; maximum is %d", len(text), maxCANInfoFieldLength)
	}
	for _, value := range []byte(text) {
		if value > 0x7F {
			return nil, errors.New("actisense: CAN Info fields support ASCII only")
		}
	}
	data := make([]byte, 2, 2+len(text))
	data[0] = byte(2 + len(text))
	data[1] = 1
	return append(data, text...), nil
}

func DecodeCANInfoField(response BEMResponse, field CANInfoField) (string, error) {
	command, err := field.command()
	if err != nil {
		return "", err
	}
	if response.BEMID != command {
		return "", fmt.Errorf("actisense: expected CAN Info field %d response, got BEM 0x%02X", field, response.BEMID)
	}
	if len(response.Data) < 2 {
		return "", errors.New("actisense: CAN Info field response is shorter than two bytes")
	}
	length := int(response.Data[0])
	if length < 2 || length > 2+maxCANInfoFieldLength || length > len(response.Data) {
		return "", fmt.Errorf("actisense: CAN Info field length %d is invalid for %d response bytes", length, len(response.Data))
	}
	if response.Data[1] != 1 {
		return "", fmt.Errorf("actisense: CAN Info field encoding %d is unsupported", response.Data[1])
	}
	return string(response.Data[2:length]), nil
}

func EncodeEcho(payload []byte) ([]byte, error) {
	if len(payload) > maxEchoPayloadLength {
		return nil, fmt.Errorf("actisense: Echo payload is %d bytes; maximum is %d", len(payload), maxEchoPayloadLength)
	}
	data := make([]byte, 1, 1+len(payload))
	data[0] = byte(len(payload))
	return append(data, payload...), nil
}

func DecodeEcho(response BEMResponse) ([]byte, error) {
	if response.BEMID != BEMEcho {
		return nil, fmt.Errorf("actisense: expected Echo response, got BEM 0x%02X", response.BEMID)
	}
	if len(response.Data) == 0 {
		return nil, nil
	}
	length := int(response.Data[0])
	if len(response.Data) < 1+length {
		return nil, fmt.Errorf("actisense: Echo response reports %d bytes; only %d present", length, len(response.Data)-1)
	}
	return append([]byte(nil), response.Data[1:1+length]...), nil
}

func TotalTimeSet(seconds, passkey uint32) []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint32(data[:4], seconds)
	binary.LittleEndian.PutUint32(data[4:8], passkey)
	return data
}

func DecodeTotalTime(response BEMResponse) (uint32, error) {
	if response.BEMID != BEMTotalTime {
		return 0, fmt.Errorf("actisense: expected Total Time response, got BEM 0x%02X", response.BEMID)
	}
	if len(response.Data) < 4 {
		return 0, fmt.Errorf("actisense: Total Time response is %d bytes; expected 4", len(response.Data))
	}
	return binary.LittleEndian.Uint32(response.Data[:4]), nil
}
