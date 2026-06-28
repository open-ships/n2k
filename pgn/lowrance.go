package pgn

import (
	"fmt"
	"github.com/open-ships/n2k/units"
)

type LowranceTemperature struct {
	Info              MessageInfo            `json:"info"`
	ManufacturerCode  ManufacturerCodeConst  `json:"manufacturerCode"`
	IndustryCode      IndustryCodeConst      `json:"industryCode"`
	TemperatureSource TemperatureSourceConst `json:"temperatureSource"`
	ActualTemperature *units.Temperature     `json:"actualTemperature"`
}

func (l *LowranceTemperature) PGNNumber() uint32 { return 65285 }

func DecodeLowranceTemperature(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &LowranceTemperature{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceTemperature-ManufacturerCode: %w", err)
	} else {
		if v != 140 {
			return nil, fmt.Errorf("match failed for LowranceTemperature-ManufacturerCode: Expected %d != %d", 140, v)
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
		return nil, fmt.Errorf("parse failed for LowranceTemperature-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for LowranceTemperature-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceTemperature-TemperatureSource: %w", err)
	} else {
		val.TemperatureSource = TemperatureSourceConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceTemperature-ActualTemperature: %w", err)
	} else {
		val.ActualTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

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

func EncodeLowranceTemperature(val *LowranceTemperature) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.TemperatureSource), 8)
	var actualTemperatureRaw *float32
	if val.ActualTemperature != nil {
		actualTemperatureRaw = &val.ActualTemperature.Value
	}
	w.writeUnsignedResolution(actualTemperatureRaw, 16, 0.01)
	w.writeReservedBits(24)
	return w.Bytes(), w.Err()
}

func encodeLowranceTemperatureMsg(v Message) ([]byte, error) {
	val, ok := v.(*LowranceTemperature)
	if !ok {
		return nil, fmt.Errorf("expected *LowranceTemperature, got %T", v)
	}
	return EncodeLowranceTemperature(val)
}

type LowranceProductInformation struct {
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

func (l *LowranceProductInformation) PGNNumber() uint32 { return 130817 }

func DecodeLowranceProductInformation(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &LowranceProductInformation{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-ManufacturerCode: %w", err)
	} else {
		if v != 140 {
			return nil, fmt.Errorf("match failed for LowranceProductInformation-ManufacturerCode: Expected %d != %d", 140, v)
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
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for LowranceProductInformation-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-ProductCode: %w", err)
	} else {
		val.ProductCode = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-Model: %w", err)
	} else {
		val.Model = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(80); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-FirmwareVersion: %w", err)
	} else {
		val.FirmwareVersion = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-FirmwareDate: %w", err)
	} else {
		val.FirmwareDate = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-FirmwareTime: %w", err)
	} else {
		val.FirmwareTime = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeLowranceProductInformation(val *LowranceProductInformation) ([]byte, error) {
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

func encodeLowranceProductInformationMsg(v Message) ([]byte, error) {
	val, ok := v.(*LowranceProductInformation)
	if !ok {
		return nil, fmt.Errorf("expected *LowranceProductInformation, got %T", v)
	}
	return EncodeLowranceProductInformation(val)
}
