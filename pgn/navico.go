package pgn

import "fmt"

type NavicoWirelessBatteryStatus struct {
	Info                MessageInfo           `json:"info"`
	ManufacturerCode    ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode        IndustryCodeConst     `json:"industryCode"`
	Status              *uint8                `json:"status"`
	BatteryStatus       *uint8                `json:"batteryStatus"`
	BatteryChargeStatus *uint8                `json:"batteryChargeStatus"`
}

func (n *NavicoWirelessBatteryStatus) PGNNumber() uint32 { return 65309 }

func DecodeNavicoWirelessBatteryStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &NavicoWirelessBatteryStatus{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoWirelessBatteryStatus-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for NavicoWirelessBatteryStatus-ManufacturerCode: Expected %d != %d", 275, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoWirelessBatteryStatus-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for NavicoWirelessBatteryStatus-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoWirelessBatteryStatus-Status: %w", err)
	} else {
		val.Status = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoWirelessBatteryStatus-BatteryStatus: %w", err)
	} else {
		val.BatteryStatus = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoWirelessBatteryStatus-BatteryChargeStatus: %w", err)
	} else {
		val.BatteryChargeStatus = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(24)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeNavicoWirelessBatteryStatus(val *NavicoWirelessBatteryStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Status, 8)
	w.writeUInt8(val.BatteryStatus, 8)
	w.writeUInt8(val.BatteryChargeStatus, 8)
	w.writeReservedBits(24)
	return w.Bytes(), w.Err()
}

func encodeNavicoWirelessBatteryStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*NavicoWirelessBatteryStatus)
	if !ok {
		return nil, fmt.Errorf("expected *NavicoWirelessBatteryStatus, got %T", v)
	}
	return EncodeNavicoWirelessBatteryStatus(val)
}

type NavicoWirelessSignalStatus struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	Unknown          *uint8                `json:"unknown"`
	SignalStrength   *uint8                `json:"signalStrength"`
}

func (n *NavicoWirelessSignalStatus) PGNNumber() uint32 { return 65312 }

func DecodeNavicoWirelessSignalStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &NavicoWirelessSignalStatus{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoWirelessSignalStatus-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for NavicoWirelessSignalStatus-ManufacturerCode: Expected %d != %d", 275, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoWirelessSignalStatus-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for NavicoWirelessSignalStatus-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoWirelessSignalStatus-Unknown: %w", err)
	} else {
		val.Unknown = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoWirelessSignalStatus-SignalStrength: %w", err)
	} else {
		val.SignalStrength = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(32)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeNavicoWirelessSignalStatus(val *NavicoWirelessSignalStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Unknown, 8)
	w.writeUInt8(val.SignalStrength, 8)
	w.writeReservedBits(32)
	return w.Bytes(), w.Err()
}

func encodeNavicoWirelessSignalStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*NavicoWirelessSignalStatus)
	if !ok {
		return nil, fmt.Errorf("expected *NavicoWirelessSignalStatus, got %T", v)
	}
	return EncodeNavicoWirelessSignalStatus(val)
}

type NavicoProductInformation struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProductCode      *uint16               `json:"productCode"`
	Model            string                `json:"model"`
	A                *uint8                `json:"a"`
	B                *uint8                `json:"b"`
	C                *uint8                `json:"c"`
	FirmwareVersion  string                `json:"firmwareVersion"`
	FirmwareDate     string                `json:"firmwareDate"`
	FirmwareTime     string                `json:"firmwareTime"`
}

func (n *NavicoProductInformation) PGNNumber() uint32 { return 130817 }

func DecodeNavicoProductInformation(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &NavicoProductInformation{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for NavicoProductInformation-ManufacturerCode: Expected %d != %d", 275, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for NavicoProductInformation-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-ProductCode: %w", err)
	} else {
		val.ProductCode = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-Model: %w", err)
	} else {
		val.Model = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(80); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-FirmwareVersion: %w", err)
	} else {
		val.FirmwareVersion = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-FirmwareDate: %w", err)
	} else {
		val.FirmwareDate = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-FirmwareTime: %w", err)
	} else {
		val.FirmwareTime = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeNavicoProductInformation(val *NavicoProductInformation) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt16(val.ProductCode, 16)
	w.writeFixedString(val.Model, 256)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.C, 8)
	w.writeFixedString(val.FirmwareVersion, 80)
	w.writeFixedString(val.FirmwareDate, 256)
	w.writeFixedString(val.FirmwareTime, 256)
	return w.Bytes(), w.Err()
}

func encodeNavicoProductInformationMsg(v Message) ([]byte, error) {
	val, ok := v.(*NavicoProductInformation)
	if !ok {
		return nil, fmt.Errorf("expected *NavicoProductInformation, got %T", v)
	}
	return EncodeNavicoProductInformation(val)
}

type NavicoAsciiData struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	A                *uint8                `json:"a"`
	Message          string                `json:"message"`
}

func (n *NavicoAsciiData) PGNNumber() uint32 { return 130821 }

func DecodeNavicoAsciiData(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &NavicoAsciiData{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoAsciiData-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for NavicoAsciiData-ManufacturerCode: Expected %d != %d", 275, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoAsciiData-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for NavicoAsciiData-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoAsciiData-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(2048); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoAsciiData-Message: %w", err)
	} else {
		val.Message = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeNavicoAsciiData(val *NavicoAsciiData) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.A, 8)
	w.writeFixedString(val.Message, 2048)
	return w.Bytes(), w.Err()
}

func encodeNavicoAsciiDataMsg(v Message) ([]byte, error) {
	val, ok := v.(*NavicoAsciiData)
	if !ok {
		return nil, fmt.Errorf("expected *NavicoAsciiData, got %T", v)
	}
	return EncodeNavicoAsciiData(val)
}

type NavicoUnknown1 struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	Data             []uint8               `json:"data"`
}

func (n *NavicoUnknown1) PGNNumber() uint32 { return 130822 }

func DecodeNavicoUnknown1(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &NavicoUnknown1{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoUnknown1-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for NavicoUnknown1-ManufacturerCode: Expected %d != %d", 275, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoUnknown1-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for NavicoUnknown1-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(1848); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoUnknown1-Data: %w", err)
	} else {
		val.Data = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeNavicoUnknown1(val *NavicoUnknown1) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeBinaryData(val.Data, 1848)
	return w.Bytes(), w.Err()
}

func encodeNavicoUnknown1Msg(v Message) ([]byte, error) {
	val, ok := v.(*NavicoUnknown1)
	if !ok {
		return nil, fmt.Errorf("expected *NavicoUnknown1, got %T", v)
	}
	return EncodeNavicoUnknown1(val)
}

type NavicoUnknown2 struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	Data             []uint8               `json:"data"`
}

func (n *NavicoUnknown2) PGNNumber() uint32 { return 130825 }

func DecodeNavicoUnknown2(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &NavicoUnknown2{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoUnknown2-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for NavicoUnknown2-ManufacturerCode: Expected %d != %d", 275, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoUnknown2-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for NavicoUnknown2-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(80); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoUnknown2-Data: %w", err)
	} else {
		val.Data = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeNavicoUnknown2(val *NavicoUnknown2) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeBinaryData(val.Data, 80)
	return w.Bytes(), w.Err()
}

func encodeNavicoUnknown2Msg(v Message) ([]byte, error) {
	val, ok := v.(*NavicoUnknown2)
	if !ok {
		return nil, fmt.Errorf("expected *NavicoUnknown2, got %T", v)
	}
	return EncodeNavicoUnknown2(val)
}
