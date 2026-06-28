package pgn

import (
	"fmt"
	"github.com/open-ships/n2k/units"
)

type MaretronProprietaryDcBreakerCurrent struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	BankInstance     *uint8                `json:"bankInstance"`
	IndicatorNumber  *uint8                `json:"indicatorNumber"`
	BreakerCurrent   *float32              `json:"breakerCurrent"`
}

func (m *MaretronProprietaryDcBreakerCurrent) PGNNumber() uint32 { return 65284 }

func DecodeMaretronProprietaryDcBreakerCurrent(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &MaretronProprietaryDcBreakerCurrent{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryDcBreakerCurrent-ManufacturerCode: %w", err)
	} else {
		if v != 137 {
			return nil, fmt.Errorf("match failed for MaretronProprietaryDcBreakerCurrent-ManufacturerCode: Expected %d != %d", 137, v)
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
		return nil, fmt.Errorf("parse failed for MaretronProprietaryDcBreakerCurrent-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for MaretronProprietaryDcBreakerCurrent-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryDcBreakerCurrent-BankInstance: %w", err)
	} else {
		val.BankInstance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryDcBreakerCurrent-IndicatorNumber: %w", err)
	} else {
		val.IndicatorNumber = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryDcBreakerCurrent-BreakerCurrent: %w", err)
	} else {
		val.BreakerCurrent = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(16)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeMaretronProprietaryDcBreakerCurrent(val *MaretronProprietaryDcBreakerCurrent) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.BankInstance, 8)
	w.writeUInt8(val.IndicatorNumber, 8)
	w.writeUnsignedResolution(val.BreakerCurrent, 16, 0.1)
	w.writeReservedBits(16)
	return w.Bytes(), w.Err()
}

func encodeMaretronProprietaryDcBreakerCurrentMsg(v Message) ([]byte, error) {
	val, ok := v.(*MaretronProprietaryDcBreakerCurrent)
	if !ok {
		return nil, fmt.Errorf("expected *MaretronProprietaryDcBreakerCurrent, got %T", v)
	}
	return EncodeMaretronProprietaryDcBreakerCurrent(val)
}

type MaretronSlaveResponse struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProductCode      *uint16               `json:"productCode"`
	SoftwareCode     *uint16               `json:"softwareCode"`
	Command          *uint8                `json:"command"`
	Status           *uint8                `json:"status"`
}

func (m *MaretronSlaveResponse) PGNNumber() uint32 { return 126720 }

func DecodeMaretronSlaveResponse(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &MaretronSlaveResponse{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSlaveResponse-ManufacturerCode: %w", err)
	} else {
		if v != 137 {
			return nil, fmt.Errorf("match failed for MaretronSlaveResponse-ManufacturerCode: Expected %d != %d", 137, v)
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
		return nil, fmt.Errorf("parse failed for MaretronSlaveResponse-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for MaretronSlaveResponse-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSlaveResponse-ProductCode: %w", err)
	} else {
		val.ProductCode = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSlaveResponse-SoftwareCode: %w", err)
	} else {
		val.SoftwareCode = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSlaveResponse-Command: %w", err)
	} else {
		val.Command = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSlaveResponse-Status: %w", err)
	} else {
		val.Status = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeMaretronSlaveResponse(val *MaretronSlaveResponse) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt16(val.ProductCode, 16)
	w.writeUInt16(val.SoftwareCode, 16)
	w.writeUInt8(val.Command, 8)
	w.writeUInt8(val.Status, 8)
	return w.Bytes(), w.Err()
}

func encodeMaretronSlaveResponseMsg(v Message) ([]byte, error) {
	val, ok := v.(*MaretronSlaveResponse)
	if !ok {
		return nil, fmt.Errorf("expected *MaretronSlaveResponse, got %T", v)
	}
	return EncodeMaretronSlaveResponse(val)
}

type MaretronProprietaryTemperatureHighRange struct {
	Info              MessageInfo            `json:"info"`
	ManufacturerCode  ManufacturerCodeConst  `json:"manufacturerCode"`
	IndustryCode      IndustryCodeConst      `json:"industryCode"`
	Sid               *uint8                 `json:"sid"`
	Instance          *uint8                 `json:"instance"`
	Source            TemperatureSourceConst `json:"source"`
	ActualTemperature *units.Temperature     `json:"actualTemperature"`
	SetTemperature    *units.Temperature     `json:"setTemperature"`
}

func (m *MaretronProprietaryTemperatureHighRange) PGNNumber() uint32 { return 130823 }

func DecodeMaretronProprietaryTemperatureHighRange(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &MaretronProprietaryTemperatureHighRange{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryTemperatureHighRange-ManufacturerCode: %w", err)
	} else {
		if v != 137 {
			return nil, fmt.Errorf("match failed for MaretronProprietaryTemperatureHighRange-ManufacturerCode: Expected %d != %d", 137, v)
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
		return nil, fmt.Errorf("parse failed for MaretronProprietaryTemperatureHighRange-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for MaretronProprietaryTemperatureHighRange-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryTemperatureHighRange-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryTemperatureHighRange-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryTemperatureHighRange-Source: %w", err)
	} else {
		val.Source = TemperatureSourceConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryTemperatureHighRange-ActualTemperature: %w", err)
	} else {
		val.ActualTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryTemperatureHighRange-SetTemperature: %w", err)
	} else {
		val.SetTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeMaretronProprietaryTemperatureHighRange(val *MaretronProprietaryTemperatureHighRange) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.Instance, 8)
	w.writeLookupField(uint64(val.Source), 8)
	var actualTemperatureRaw *float32
	if val.ActualTemperature != nil {
		actualTemperatureRaw = &val.ActualTemperature.Value
	}
	w.writeUnsignedResolution(actualTemperatureRaw, 16, 0.1)
	var setTemperatureRaw *float32
	if val.SetTemperature != nil {
		setTemperatureRaw = &val.SetTemperature.Value
	}
	w.writeUnsignedResolution(setTemperatureRaw, 16, 0.1)
	return w.Bytes(), w.Err()
}

func encodeMaretronProprietaryTemperatureHighRangeMsg(v Message) ([]byte, error) {
	val, ok := v.(*MaretronProprietaryTemperatureHighRange)
	if !ok {
		return nil, fmt.Errorf("expected *MaretronProprietaryTemperatureHighRange, got %T", v)
	}
	return EncodeMaretronProprietaryTemperatureHighRange(val)
}

type MaretronAnnunciator struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	Field4           *uint8                `json:"field4"`
	Field5           *uint8                `json:"field5"`
	Field6           *uint16               `json:"field6"`
	Field7           *uint8                `json:"field7"`
	Field8           *uint16               `json:"field8"`
}

func (m *MaretronAnnunciator) PGNNumber() uint32 { return 130824 }

func DecodeMaretronAnnunciator(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &MaretronAnnunciator{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronAnnunciator-ManufacturerCode: %w", err)
	} else {
		if v != 137 {
			return nil, fmt.Errorf("match failed for MaretronAnnunciator-ManufacturerCode: Expected %d != %d", 137, v)
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
		return nil, fmt.Errorf("parse failed for MaretronAnnunciator-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for MaretronAnnunciator-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronAnnunciator-Field4: %w", err)
	} else {
		val.Field4 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronAnnunciator-Field5: %w", err)
	} else {
		val.Field5 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronAnnunciator-Field6: %w", err)
	} else {
		val.Field6 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronAnnunciator-Field7: %w", err)
	} else {
		val.Field7 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronAnnunciator-Field8: %w", err)
	} else {
		val.Field8 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeMaretronAnnunciator(val *MaretronAnnunciator) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Field4, 8)
	w.writeUInt8(val.Field5, 8)
	w.writeUInt16(val.Field6, 16)
	w.writeUInt8(val.Field7, 8)
	w.writeUInt16(val.Field8, 16)
	return w.Bytes(), w.Err()
}

func encodeMaretronAnnunciatorMsg(v Message) ([]byte, error) {
	val, ok := v.(*MaretronAnnunciator)
	if !ok {
		return nil, fmt.Errorf("expected *MaretronAnnunciator, got %T", v)
	}
	return EncodeMaretronAnnunciator(val)
}

type MaretronSwitchStatusCounter struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	Instance         *uint8                `json:"instance"`
	IndicatorNumber  *uint8                `json:"indicatorNumber"`
	StartDate        *uint16               `json:"startDate"`
	StartTime        *float32              `json:"startTime"`
	OffCounter       *uint8                `json:"offCounter"`
	OnCounter        *uint8                `json:"onCounter"`
	ErrorCounter     *uint8                `json:"errorCounter"`
	SwitchStatus     OffOnConst            `json:"switchStatus"`
}

func (m *MaretronSwitchStatusCounter) PGNNumber() uint32 { return 130836 }

func DecodeMaretronSwitchStatusCounter(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &MaretronSwitchStatusCounter{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-ManufacturerCode: %w", err)
	} else {
		if v != 137 {
			return nil, fmt.Errorf("match failed for MaretronSwitchStatusCounter-ManufacturerCode: Expected %d != %d", 137, v)
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
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for MaretronSwitchStatusCounter-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-IndicatorNumber: %w", err)
	} else {
		val.IndicatorNumber = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-StartDate: %w", err)
	} else {
		val.StartDate = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(32, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-StartTime: %w", err)
	} else {
		val.StartTime = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-OffCounter: %w", err)
	} else {
		val.OffCounter = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-OnCounter: %w", err)
	} else {
		val.OnCounter = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-ErrorCounter: %w", err)
	} else {
		val.ErrorCounter = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-SwitchStatus: %w", err)
	} else {
		val.SwitchStatus = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeMaretronSwitchStatusCounter(val *MaretronSwitchStatusCounter) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Instance, 8)
	w.writeUInt8(val.IndicatorNumber, 8)
	w.writeUInt16(val.StartDate, 16)
	w.writeUnsignedResolution(val.StartTime, 32, 0.0001)
	w.writeUInt8(val.OffCounter, 8)
	w.writeUInt8(val.OnCounter, 8)
	w.writeUInt8(val.ErrorCounter, 8)
	w.writeLookupField(uint64(val.SwitchStatus), 2)
	w.writeReservedBits(6)
	return w.Bytes(), w.Err()
}

func encodeMaretronSwitchStatusCounterMsg(v Message) ([]byte, error) {
	val, ok := v.(*MaretronSwitchStatusCounter)
	if !ok {
		return nil, fmt.Errorf("expected *MaretronSwitchStatusCounter, got %T", v)
	}
	return EncodeMaretronSwitchStatusCounter(val)
}

type MaretronSwitchStatusTimer struct {
	Info                   MessageInfo           `json:"info"`
	ManufacturerCode       ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode           IndustryCodeConst     `json:"industryCode"`
	Instance               *uint8                `json:"instance"`
	IndicatorNumber        *uint8                `json:"indicatorNumber"`
	StartDate              *uint16               `json:"startDate"`
	StartTime              *float32              `json:"startTime"`
	AccumulatedOffPeriod   *uint32               `json:"accumulatedOffPeriod"`
	AccumulatedOnPeriod    *uint32               `json:"accumulatedOnPeriod"`
	AccumulatedErrorPeriod *uint32               `json:"accumulatedErrorPeriod"`
	SwitchStatus           OffOnConst            `json:"switchStatus"`
}

func (m *MaretronSwitchStatusTimer) PGNNumber() uint32 { return 130837 }

func DecodeMaretronSwitchStatusTimer(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &MaretronSwitchStatusTimer{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-ManufacturerCode: %w", err)
	} else {
		if v != 137 {
			return nil, fmt.Errorf("match failed for MaretronSwitchStatusTimer-ManufacturerCode: Expected %d != %d", 137, v)
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
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for MaretronSwitchStatusTimer-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-IndicatorNumber: %w", err)
	} else {
		val.IndicatorNumber = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-StartDate: %w", err)
	} else {
		val.StartDate = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(32, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-StartTime: %w", err)
	} else {
		val.StartTime = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-AccumulatedOffPeriod: %w", err)
	} else {
		val.AccumulatedOffPeriod = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-AccumulatedOnPeriod: %w", err)
	} else {
		val.AccumulatedOnPeriod = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-AccumulatedErrorPeriod: %w", err)
	} else {
		val.AccumulatedErrorPeriod = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-SwitchStatus: %w", err)
	} else {
		val.SwitchStatus = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeMaretronSwitchStatusTimer(val *MaretronSwitchStatusTimer) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Instance, 8)
	w.writeUInt8(val.IndicatorNumber, 8)
	w.writeUInt16(val.StartDate, 16)
	w.writeUnsignedResolution(val.StartTime, 32, 0.0001)
	w.writeUInt32(val.AccumulatedOffPeriod, 32)
	w.writeUInt32(val.AccumulatedOnPeriod, 32)
	w.writeUInt32(val.AccumulatedErrorPeriod, 32)
	w.writeLookupField(uint64(val.SwitchStatus), 2)
	w.writeReservedBits(6)
	return w.Bytes(), w.Err()
}

func encodeMaretronSwitchStatusTimerMsg(v Message) ([]byte, error) {
	val, ok := v.(*MaretronSwitchStatusTimer)
	if !ok {
		return nil, fmt.Errorf("expected *MaretronSwitchStatusTimer, got %T", v)
	}
	return EncodeMaretronSwitchStatusTimer(val)
}
