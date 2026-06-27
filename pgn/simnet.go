package pgn

import (
	"fmt"
	"github.com/open-ships/n2k/units"
)

type SimnetConfigureTemperatureSensor struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
}

func (s *SimnetConfigureTemperatureSensor) PGNNumber() uint32 { return 65287 }

func DecodeSimnetConfigureTemperatureSensor(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetConfigureTemperatureSensor{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetConfigureTemperatureSensor-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetConfigureTemperatureSensor-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetConfigureTemperatureSensor-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetConfigureTemperatureSensor-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(48)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeSimnetConfigureTemperatureSensor(val *SimnetConfigureTemperatureSensor) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeReservedBits(48)
	return w.Bytes(), w.Err()
}

func encodeSimnetConfigureTemperatureSensorMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetConfigureTemperatureSensor)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetConfigureTemperatureSensor, got %T", v)
	}
	return EncodeSimnetConfigureTemperatureSensor(val)
}

type SimnetTrimTabSensorCalibration struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
}

func (s *SimnetTrimTabSensorCalibration) PGNNumber() uint32 { return 65289 }

func DecodeSimnetTrimTabSensorCalibration(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetTrimTabSensorCalibration{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetTrimTabSensorCalibration-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetTrimTabSensorCalibration-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetTrimTabSensorCalibration-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetTrimTabSensorCalibration-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(48)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeSimnetTrimTabSensorCalibration(val *SimnetTrimTabSensorCalibration) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeReservedBits(48)
	return w.Bytes(), w.Err()
}

func encodeSimnetTrimTabSensorCalibrationMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetTrimTabSensorCalibration)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetTrimTabSensorCalibration, got %T", v)
	}
	return EncodeSimnetTrimTabSensorCalibration(val)
}

type SimnetPaddleWheelSpeedConfiguration struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
}

func (s *SimnetPaddleWheelSpeedConfiguration) PGNNumber() uint32 { return 65290 }

func DecodeSimnetPaddleWheelSpeedConfiguration(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetPaddleWheelSpeedConfiguration{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetPaddleWheelSpeedConfiguration-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetPaddleWheelSpeedConfiguration-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetPaddleWheelSpeedConfiguration-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetPaddleWheelSpeedConfiguration-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(48)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeSimnetPaddleWheelSpeedConfiguration(val *SimnetPaddleWheelSpeedConfiguration) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeReservedBits(48)
	return w.Bytes(), w.Err()
}

func encodeSimnetPaddleWheelSpeedConfigurationMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetPaddleWheelSpeedConfiguration)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetPaddleWheelSpeedConfiguration, got %T", v)
	}
	return EncodeSimnetPaddleWheelSpeedConfiguration(val)
}

type SimnetClearFluidLevelWarnings struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
}

func (s *SimnetClearFluidLevelWarnings) PGNNumber() uint32 { return 65292 }

func DecodeSimnetClearFluidLevelWarnings(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetClearFluidLevelWarnings{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetClearFluidLevelWarnings-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetClearFluidLevelWarnings-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetClearFluidLevelWarnings-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetClearFluidLevelWarnings-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(48)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeSimnetClearFluidLevelWarnings(val *SimnetClearFluidLevelWarnings) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeReservedBits(48)
	return w.Bytes(), w.Err()
}

func encodeSimnetClearFluidLevelWarningsMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetClearFluidLevelWarnings)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetClearFluidLevelWarnings, got %T", v)
	}
	return EncodeSimnetClearFluidLevelWarnings(val)
}

type SimnetLgc2000Configuration struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
}

func (s *SimnetLgc2000Configuration) PGNNumber() uint32 { return 65293 }

func DecodeSimnetLgc2000Configuration(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetLgc2000Configuration{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetLgc2000Configuration-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetLgc2000Configuration-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetLgc2000Configuration-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetLgc2000Configuration-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(48)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeSimnetLgc2000Configuration(val *SimnetLgc2000Configuration) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeReservedBits(48)
	return w.Bytes(), w.Err()
}

func encodeSimnetLgc2000ConfigurationMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetLgc2000Configuration)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetLgc2000Configuration, got %T", v)
	}
	return EncodeSimnetLgc2000Configuration(val)
}

type SimnetApUnknown1 struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	A                *uint8                `json:"a"`
	B                *uint8                `json:"b"`
	C                *uint16               `json:"c"`
	D                *uint8                `json:"d"`
}

func (s *SimnetApUnknown1) PGNNumber() uint32 { return 65302 }

func DecodeSimnetApUnknown1(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetApUnknown1{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown1-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetApUnknown1-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetApUnknown1-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetApUnknown1-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown1-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown1-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown1-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown1-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeSimnetApUnknown1(val *SimnetApUnknown1) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt16(val.C, 16)
	w.writeUInt8(val.D, 8)
	w.writeReservedBits(8)
	return w.Bytes(), w.Err()
}

func encodeSimnetApUnknown1Msg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetApUnknown1)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetApUnknown1, got %T", v)
	}
	return EncodeSimnetApUnknown1(val)
}

type SimnetDeviceModeRequest struct {
	Info             MessageInfo             `json:"info"`
	ManufacturerCode ManufacturerCodeConst   `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst       `json:"industryCode"`
	Model            SimnetDeviceModelConst  `json:"model"`
	Report           SimnetDeviceReportConst `json:"report"`
}

func (s *SimnetDeviceModeRequest) PGNNumber() uint32 { return 65305 }

func DecodeSimnetDeviceModeRequest(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetDeviceModeRequest{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceModeRequest-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetDeviceModeRequest-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetDeviceModeRequest-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetDeviceModeRequest-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceModeRequest-Model: %w", err)
	} else {
		val.Model = SimnetDeviceModelConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceModeRequest-Report: %w", err)
	} else {
		if v != 11 {
			return nil, fmt.Errorf("match failed for SimnetDeviceModeRequest-Report: Expected %d != %d", 11, v)
		}
		val.Report = SimnetDeviceReportConst(v)

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

func EncodeSimnetDeviceModeRequest(val *SimnetDeviceModeRequest) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.Model), 8)
	w.writeLookupField(uint64(val.Report), 8)
	w.writeSpareBits(32)
	return w.Bytes(), w.Err()
}

func encodeSimnetDeviceModeRequestMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetDeviceModeRequest)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetDeviceModeRequest, got %T", v)
	}
	return EncodeSimnetDeviceModeRequest(val)
}

type SimnetDeviceStatus struct {
	Info             MessageInfo             `json:"info"`
	ManufacturerCode ManufacturerCodeConst   `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst       `json:"industryCode"`
	Model            SimnetDeviceModelConst  `json:"model"`
	Report           SimnetDeviceReportConst `json:"report"`
	Status           SimnetApStatusConst     `json:"status"`
}

func (s *SimnetDeviceStatus) PGNNumber() uint32 { return 65305 }

func DecodeSimnetDeviceStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetDeviceStatus{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatus-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetDeviceStatus-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatus-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetDeviceStatus-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatus-Model: %w", err)
	} else {
		val.Model = SimnetDeviceModelConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatus-Report: %w", err)
	} else {
		if v != 2 {
			return nil, fmt.Errorf("match failed for SimnetDeviceStatus-Report: Expected %d != %d", 2, v)
		}
		val.Report = SimnetDeviceReportConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatus-Status: %w", err)
	} else {
		val.Status = SimnetApStatusConst(v)

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

func EncodeSimnetDeviceStatus(val *SimnetDeviceStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.Model), 8)
	w.writeLookupField(uint64(val.Report), 8)
	w.writeLookupField(uint64(val.Status), 8)
	w.writeSpareBits(24)
	return w.Bytes(), w.Err()
}

func encodeSimnetDeviceStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetDeviceStatus)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetDeviceStatus, got %T", v)
	}
	return EncodeSimnetDeviceStatus(val)
}

type SimnetDeviceStatusRequest struct {
	Info             MessageInfo             `json:"info"`
	ManufacturerCode ManufacturerCodeConst   `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst       `json:"industryCode"`
	Model            SimnetDeviceModelConst  `json:"model"`
	Report           SimnetDeviceReportConst `json:"report"`
}

func (s *SimnetDeviceStatusRequest) PGNNumber() uint32 { return 65305 }

func DecodeSimnetDeviceStatusRequest(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetDeviceStatusRequest{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatusRequest-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetDeviceStatusRequest-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatusRequest-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetDeviceStatusRequest-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatusRequest-Model: %w", err)
	} else {
		val.Model = SimnetDeviceModelConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatusRequest-Report: %w", err)
	} else {
		if v != 3 {
			return nil, fmt.Errorf("match failed for SimnetDeviceStatusRequest-Report: Expected %d != %d", 3, v)
		}
		val.Report = SimnetDeviceReportConst(v)

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

func EncodeSimnetDeviceStatusRequest(val *SimnetDeviceStatusRequest) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.Model), 8)
	w.writeLookupField(uint64(val.Report), 8)
	w.writeSpareBits(32)
	return w.Bytes(), w.Err()
}

func encodeSimnetDeviceStatusRequestMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetDeviceStatusRequest)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetDeviceStatusRequest, got %T", v)
	}
	return EncodeSimnetDeviceStatusRequest(val)
}

type SimnetPilotMode struct {
	Info             MessageInfo               `json:"info"`
	ManufacturerCode ManufacturerCodeConst     `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst         `json:"industryCode"`
	Model            SimnetDeviceModelConst    `json:"model"`
	Report           SimnetDeviceReportConst   `json:"report"`
	Mode             SimnetApModeBitfieldConst `json:"mode"`
}

func (s *SimnetPilotMode) PGNNumber() uint32 { return 65305 }

func DecodeSimnetPilotMode(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetPilotMode{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetPilotMode-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetPilotMode-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetPilotMode-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetPilotMode-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetPilotMode-Model: %w", err)
	} else {
		val.Model = SimnetDeviceModelConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetPilotMode-Report: %w", err)
	} else {
		if v != 10 {
			return nil, fmt.Errorf("match failed for SimnetPilotMode-Report: Expected %d != %d", 10, v)
		}
		val.Report = SimnetDeviceReportConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetPilotMode-Mode: %w", err)
	} else {
		val.Mode = SimnetApModeBitfieldConst(v)

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

func EncodeSimnetPilotMode(val *SimnetPilotMode) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.Model), 8)
	w.writeLookupField(uint64(val.Report), 8)
	w.writeLookupField(uint64(val.Mode), 16)
	w.writeSpareBits(16)
	return w.Bytes(), w.Err()
}

func encodeSimnetPilotModeMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetPilotMode)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetPilotMode, got %T", v)
	}
	return EncodeSimnetPilotMode(val)
}

type SimnetSailingProcessorStatus struct {
	Info             MessageInfo             `json:"info"`
	ManufacturerCode ManufacturerCodeConst   `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst       `json:"industryCode"`
	Model            SimnetDeviceModelConst  `json:"model"`
	Report           SimnetDeviceReportConst `json:"report"`
	Data             []uint8                 `json:"data"`
}

func (s *SimnetSailingProcessorStatus) PGNNumber() uint32 { return 65305 }

func DecodeSimnetSailingProcessorStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetSailingProcessorStatus{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetSailingProcessorStatus-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetSailingProcessorStatus-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetSailingProcessorStatus-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetSailingProcessorStatus-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetSailingProcessorStatus-Model: %w", err)
	} else {
		val.Model = SimnetDeviceModelConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetSailingProcessorStatus-Report: %w", err)
	} else {
		if v != 23 {
			return nil, fmt.Errorf("match failed for SimnetSailingProcessorStatus-Report: Expected %d != %d", 23, v)
		}
		val.Report = SimnetDeviceReportConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(32); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetSailingProcessorStatus-Data: %w", err)
	} else {
		val.Data = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSimnetSailingProcessorStatus(val *SimnetSailingProcessorStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.Model), 8)
	w.writeLookupField(uint64(val.Report), 8)
	w.writeBinaryData(val.Data, 32)
	return w.Bytes(), w.Err()
}

func encodeSimnetSailingProcessorStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetSailingProcessorStatus)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetSailingProcessorStatus, got %T", v)
	}
	return EncodeSimnetSailingProcessorStatus(val)
}

type SimnetApUnknown2 struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	A                *uint8                `json:"a"`
	B                *uint8                `json:"b"`
	C                *uint8                `json:"c"`
	D                *uint8                `json:"d"`
	E                *uint8                `json:"e"`
}

func (s *SimnetApUnknown2) PGNNumber() uint32 { return 65340 }

func DecodeSimnetApUnknown2(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetApUnknown2{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown2-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetApUnknown2-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetApUnknown2-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetApUnknown2-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown2-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown2-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown2-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown2-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown2-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeSimnetApUnknown2(val *SimnetApUnknown2) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.C, 8)
	w.writeUInt8(val.D, 8)
	w.writeUInt8(val.E, 8)
	w.writeReservedBits(8)
	return w.Bytes(), w.Err()
}

func encodeSimnetApUnknown2Msg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetApUnknown2)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetApUnknown2, got %T", v)
	}
	return EncodeSimnetApUnknown2(val)
}

type SimnetAutopilotAngle struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	Mode             SimnetApModeConst     `json:"mode"`
	Angle            *float32              `json:"angle"`
}

func (s *SimnetAutopilotAngle) PGNNumber() uint32 { return 65341 }

func DecodeSimnetAutopilotAngle(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetAutopilotAngle{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAutopilotAngle-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetAutopilotAngle-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetAutopilotAngle-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetAutopilotAngle-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(16)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAutopilotAngle-Mode: %w", err)
	} else {
		val.Mode = SimnetApModeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAutopilotAngle-Angle: %w", err)
	} else {
		val.Angle = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSimnetAutopilotAngle(val *SimnetAutopilotAngle) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeReservedBits(16)
	w.writeLookupField(uint64(val.Mode), 8)
	w.writeReservedBits(8)
	w.writeUnsignedResolution(val.Angle, 16, 0.0001)
	return w.Bytes(), w.Err()
}

func encodeSimnetAutopilotAngleMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetAutopilotAngle)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetAutopilotAngle, got %T", v)
	}
	return EncodeSimnetAutopilotAngle(val)
}

type SimnetMagneticField struct {
	Info MessageInfo `json:"info"`
	A    *float32    `json:"a"`
	B    *uint8      `json:"b"`
	C    *float32    `json:"c"`
	D    *float32    `json:"d"`
}

func (s *SimnetMagneticField) PGNNumber() uint32 { return 65350 }

func DecodeSimnetMagneticField(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetMagneticField{}
	val.Info = Info
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetMagneticField-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetMagneticField-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetMagneticField-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetMagneticField-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeSimnetMagneticField(val *SimnetMagneticField) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeSignedResolution(val.A, 16, 0.0001)
	w.writeUInt8(val.B, 8)
	w.writeSignedResolution(val.C, 16, 0.0001)
	w.writeSignedResolution(val.D, 16, 0.0001)
	w.writeReservedBits(8)
	return w.Bytes(), w.Err()
}

func encodeSimnetMagneticFieldMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetMagneticField)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetMagneticField, got %T", v)
	}
	return EncodeSimnetMagneticField(val)
}

type SimnetApUnknown3 struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	A                *uint8                `json:"a"`
	B                *uint8                `json:"b"`
	C                *uint8                `json:"c"`
	D                *uint8                `json:"d"`
	E                *uint8                `json:"e"`
}

func (s *SimnetApUnknown3) PGNNumber() uint32 { return 65420 }

func DecodeSimnetApUnknown3(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetApUnknown3{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown3-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetApUnknown3-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetApUnknown3-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetApUnknown3-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown3-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown3-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown3-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown3-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown3-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeSimnetApUnknown3(val *SimnetApUnknown3) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.C, 8)
	w.writeUInt8(val.D, 8)
	w.writeUInt8(val.E, 8)
	w.writeReservedBits(8)
	return w.Bytes(), w.Err()
}

func encodeSimnetApUnknown3Msg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetApUnknown3)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetApUnknown3, got %T", v)
	}
	return EncodeSimnetApUnknown3(val)
}

type SimnetAutopilotMode struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
}

func (s *SimnetAutopilotMode) PGNNumber() uint32 { return 65480 }

func DecodeSimnetAutopilotMode(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetAutopilotMode{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAutopilotMode-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetAutopilotMode-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetAutopilotMode-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetAutopilotMode-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(48)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeSimnetAutopilotMode(val *SimnetAutopilotMode) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeReservedBits(48)
	return w.Bytes(), w.Err()
}

func encodeSimnetAutopilotModeMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetAutopilotMode)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetAutopilotMode, got %T", v)
	}
	return EncodeSimnetAutopilotMode(val)
}

type SimnetReprogramData struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	Version          *uint16               `json:"version"`
	Sequence         *uint16               `json:"sequence"`
	Data             []uint8               `json:"data"`
}

func (s *SimnetReprogramData) PGNNumber() uint32 { return 130818 }

func DecodeSimnetReprogramData(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetReprogramData{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetReprogramData-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetReprogramData-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetReprogramData-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetReprogramData-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetReprogramData-Version: %w", err)
	} else {
		val.Version = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetReprogramData-Sequence: %w", err)
	} else {
		val.Sequence = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(1736); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetReprogramData-Data: %w", err)
	} else {
		val.Data = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSimnetReprogramData(val *SimnetReprogramData) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt16(val.Version, 16)
	w.writeUInt16(val.Sequence, 16)
	w.writeBinaryData(val.Data, 1736)
	return w.Bytes(), w.Err()
}

func encodeSimnetReprogramDataMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetReprogramData)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetReprogramData, got %T", v)
	}
	return EncodeSimnetReprogramData(val)
}

type SimnetFluidLevelSensorConfiguration struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	C                *uint8                `json:"c"`
	Device           *uint8                `json:"device"`
	Instance         *uint8                `json:"instance"`
	F                *uint8                `json:"f"`
	TankType         TankTypeConst         `json:"tankType"`
	Capacity         *units.Volume         `json:"capacity"`
	G                *uint8                `json:"g"`
	H                *int16                `json:"h"`
	I                *int8                 `json:"i"`
}

func (s *SimnetFluidLevelSensorConfiguration) PGNNumber() uint32 { return 130836 }

func DecodeSimnetFluidLevelSensorConfiguration(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetFluidLevelSensorConfiguration{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetFluidLevelSensorConfiguration-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetFluidLevelSensorConfiguration-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-Device: %w", err)
	} else {
		val.Device = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-F: %w", err)
	} else {
		val.F = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-TankType: %w", err)
	} else {
		val.TankType = TankTypeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(32, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-Capacity: %w", err)
	} else {
		val.Capacity = nullableUnit(units.Liter, v, units.NewVolume)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-G: %w", err)
	} else {
		val.G = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-H: %w", err)
	} else {
		val.H = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-I: %w", err)
	} else {
		val.I = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSimnetFluidLevelSensorConfiguration(val *SimnetFluidLevelSensorConfiguration) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.C, 8)
	w.writeUInt8(val.Device, 8)
	w.writeUInt8(val.Instance, 8)
	w.writeUInt8(val.F, 4)
	w.writeLookupField(uint64(val.TankType), 4)
	var capacityRaw *float32
	if val.Capacity != nil {
		capacityRaw = &val.Capacity.Value
	}
	w.writeUnsignedResolution(capacityRaw, 32, 0.1)
	w.writeUInt8(val.G, 8)
	w.writeInt16(val.H, 16)
	w.writeInt8(val.I, 8)
	return w.Bytes(), w.Err()
}

func encodeSimnetFluidLevelSensorConfigurationMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetFluidLevelSensorConfiguration)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetFluidLevelSensorConfiguration, got %T", v)
	}
	return EncodeSimnetFluidLevelSensorConfiguration(val)
}

type SimnetAisClassBStaticDataMsg24PartB struct {
	Info                           MessageInfo           `json:"info"`
	ManufacturerCode               ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode                   IndustryCodeConst     `json:"industryCode"`
	MessageId                      *uint8                `json:"messageId"`
	RepeatIndicator                RepeatIndicatorConst  `json:"repeatIndicator"`
	D                              *uint8                `json:"d"`
	E                              *uint8                `json:"e"`
	UserId                         *uint32               `json:"userId"`
	TypeOfShip                     ShipTypeConst         `json:"typeOfShip"`
	VendorId                       string                `json:"vendorId"`
	Callsign                       string                `json:"callsign"`
	Length                         *units.Distance       `json:"length"`
	Beam                           *units.Distance       `json:"beam"`
	PositionReferenceFromStarboard *units.Distance       `json:"positionReferenceFromStarboard"`
	PositionReferenceFromBow       *units.Distance       `json:"positionReferenceFromBow"`
	MothershipUserId               *uint32               `json:"mothershipUserId"`
}

func (s *SimnetAisClassBStaticDataMsg24PartB) PGNNumber() uint32 { return 130842 }

func DecodeSimnetAisClassBStaticDataMsg24PartB(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetAisClassBStaticDataMsg24PartB{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetAisClassBStaticDataMsg24PartB-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetAisClassBStaticDataMsg24PartB-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(6); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-MessageId: %w", err)
	} else {
		if v != nil && *v != 1 {
			return nil, fmt.Errorf("match failed for SimnetAisClassBStaticDataMsg24PartB-MessageId: Expected %d != %d", 1, *v)
		}
		val.MessageId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-UserId: %w", err)
	} else {
		val.UserId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-TypeOfShip: %w", err)
	} else {
		val.TypeOfShip = ShipTypeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(56); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-VendorId: %w", err)
	} else {
		val.VendorId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(56); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-Callsign: %w", err)
	} else {
		val.Callsign = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-Length: %w", err)
	} else {
		val.Length = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-Beam: %w", err)
	} else {
		val.Beam = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-PositionReferenceFromStarboard: %w", err)
	} else {
		val.PositionReferenceFromStarboard = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-PositionReferenceFromBow: %w", err)
	} else {
		val.PositionReferenceFromBow = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-MothershipUserId: %w", err)
	} else {
		val.MothershipUserId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeSimnetAisClassBStaticDataMsg24PartB(val *SimnetAisClassBStaticDataMsg24PartB) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.MessageId, 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt8(val.D, 8)
	w.writeUInt8(val.E, 8)
	w.writeUInt32(val.UserId, 32)
	w.writeLookupField(uint64(val.TypeOfShip), 8)
	w.writeFixedString(val.VendorId, 56)
	w.writeFixedString(val.Callsign, 56)
	var lengthRaw *float32
	if val.Length != nil {
		lengthRaw = &val.Length.Value
	}
	w.writeUnsignedResolution(lengthRaw, 16, 0.1)
	var beamRaw *float32
	if val.Beam != nil {
		beamRaw = &val.Beam.Value
	}
	w.writeUnsignedResolution(beamRaw, 16, 0.1)
	var positionReferenceFromStarboardRaw *float32
	if val.PositionReferenceFromStarboard != nil {
		positionReferenceFromStarboardRaw = &val.PositionReferenceFromStarboard.Value
	}
	w.writeUnsignedResolution(positionReferenceFromStarboardRaw, 16, 0.1)
	var positionReferenceFromBowRaw *float32
	if val.PositionReferenceFromBow != nil {
		positionReferenceFromBowRaw = &val.PositionReferenceFromBow.Value
	}
	w.writeUnsignedResolution(positionReferenceFromBowRaw, 16, 0.1)
	w.writeUInt32(val.MothershipUserId, 32)
	w.writeSpareBits(6)
	w.writeReservedBits(2)
	return w.Bytes(), w.Err()
}

func encodeSimnetAisClassBStaticDataMsg24PartBMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetAisClassBStaticDataMsg24PartB)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetAisClassBStaticDataMsg24PartB, got %T", v)
	}
	return EncodeSimnetAisClassBStaticDataMsg24PartB(val)
}

type SimnetKeyValue struct {
	Info             MessageInfo             `json:"info"`
	ManufacturerCode ManufacturerCodeConst   `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst       `json:"industryCode"`
	Address          *uint8                  `json:"address"`
	RepeatIndicator  RepeatIndicatorConst    `json:"repeatIndicator"`
	DisplayGroup     SimnetDisplayGroupConst `json:"displayGroup"`
	Key              SimnetKeyValueConst     `json:"key"`
	Minlength        *uint8                  `json:"minlength"`
	Value            []uint8                 `json:"value"`
}

func (s *SimnetKeyValue) PGNNumber() uint32 { return 130845 }

func DecodeSimnetKeyValue(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetKeyValue{}
	val.Info = Info
	var valueLength uint16

	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetKeyValue-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetKeyValue-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetKeyValue-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetKeyValue-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetKeyValue-Address: %w", err)
	} else {
		val.Address = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetKeyValue-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetKeyValue-DisplayGroup: %w", err)
	} else {
		val.DisplayGroup = SimnetDisplayGroupConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetKeyValue-Key: %w", err)
	} else {
		val.Key = SimnetKeyValueConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetKeyValue-Minlength: %w", err)
	} else {
		val.Minlength = v
		if v != nil {
			valueLength = uint16(*v) * 8
		}

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(valueLength); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetKeyValue-Value: %w", err)
	} else {
		val.Value = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSimnetKeyValue(val *SimnetKeyValue) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Address, 8)
	w.writeLookupField(uint64(val.RepeatIndicator), 8)
	w.writeLookupField(uint64(val.DisplayGroup), 8)
	w.writeReservedBits(8)
	w.writeLookupField(uint64(val.Key), 16)
	w.writeSpareBits(8)
	w.writeUInt8(val.Minlength, 8)
	return w.Bytes(), w.Err()
}

func encodeSimnetKeyValueMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetKeyValue)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetKeyValue, got %T", v)
	}
	return EncodeSimnetKeyValue(val)
}

type SimnetParameterSet struct {
	Info             MessageInfo             `json:"info"`
	ManufacturerCode ManufacturerCodeConst   `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst       `json:"industryCode"`
	Address          *uint8                  `json:"address"`
	B                *uint8                  `json:"b"`
	DisplayGroup     SimnetDisplayGroupConst `json:"displayGroup"`
	D                *uint16                 `json:"d"`
	Key              SimnetKeyValueConst     `json:"key"`
	Length           *uint8                  `json:"length"`
	Value            []uint8                 `json:"value"`
}

func (s *SimnetParameterSet) PGNNumber() uint32 { return 130846 }

func DecodeSimnetParameterSet(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetParameterSet{}
	val.Info = Info
	var valueLength uint16

	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetParameterSet-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetParameterSet-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-Address: %w", err)
	} else {
		val.Address = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-DisplayGroup: %w", err)
	} else {
		val.DisplayGroup = SimnetDisplayGroupConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-Key: %w", err)
	} else {
		val.Key = SimnetKeyValueConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-Length: %w", err)
	} else {
		val.Length = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(valueLength); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-Value: %w", err)
	} else {
		val.Value = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSimnetParameterSet(val *SimnetParameterSet) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Address, 8)
	w.writeUInt8(val.B, 8)
	w.writeLookupField(uint64(val.DisplayGroup), 8)
	w.writeUInt16(val.D, 16)
	w.writeLookupField(uint64(val.Key), 16)
	w.writeSpareBits(8)
	w.writeUInt8(val.Length, 8)
	return w.Bytes(), w.Err()
}

func encodeSimnetParameterSetMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetParameterSet)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetParameterSet, got %T", v)
	}
	return EncodeSimnetParameterSet(val)
}

type SimnetAlarm struct {
	Info             MessageInfo             `json:"info"`
	ManufacturerCode ManufacturerCodeConst   `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst       `json:"industryCode"`
	Address          *uint8                  `json:"address"`
	ProprietaryId    SimnetEventCommandConst `json:"proprietaryId"`
	Alarm            SimnetAlarmConst        `json:"alarm"`
	MessageId        *uint16                 `json:"messageId"`
	F                *uint8                  `json:"f"`
	G                *uint8                  `json:"g"`
}

func (s *SimnetAlarm) PGNNumber() uint32 { return 130850 }

func DecodeSimnetAlarm(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetAlarm{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarm-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetAlarm-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetAlarm-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetAlarm-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarm-Address: %w", err)
	} else {
		val.Address = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarm-ProprietaryId: %w", err)
	} else {
		if v != 1 {
			return nil, fmt.Errorf("match failed for SimnetAlarm-ProprietaryId: Expected %d != %d", 1, v)
		}
		val.ProprietaryId = SimnetEventCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarm-Alarm: %w", err)
	} else {
		val.Alarm = SimnetAlarmConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarm-MessageId: %w", err)
	} else {
		val.MessageId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarm-F: %w", err)
	} else {
		val.F = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarm-G: %w", err)
	} else {
		val.G = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSimnetAlarm(val *SimnetAlarm) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Address, 8)
	w.writeReservedBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeReservedBits(8)
	w.writeLookupField(uint64(val.Alarm), 16)
	w.writeUInt16(val.MessageId, 16)
	w.writeUInt8(val.F, 8)
	w.writeUInt8(val.G, 8)
	return w.Bytes(), w.Err()
}

func encodeSimnetAlarmMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetAlarm)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetAlarm, got %T", v)
	}
	return EncodeSimnetAlarm(val)
}

type SimnetApCommand struct {
	Info             MessageInfo             `json:"info"`
	ManufacturerCode ManufacturerCodeConst   `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst       `json:"industryCode"`
	Address          *uint8                  `json:"address"`
	ProprietaryId    SimnetEventCommandConst `json:"proprietaryId"`
	ApStatus         SimnetApStatusConst     `json:"apStatus"`
	ApCommand        SimnetApEventsConst     `json:"apCommand"`
	Direction        SimnetDirectionConst    `json:"direction"`
	Angle            *float32                `json:"angle"`
}

func (s *SimnetApCommand) PGNNumber() uint32 { return 130850 }

func DecodeSimnetApCommand(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetApCommand{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApCommand-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetApCommand-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetApCommand-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetApCommand-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApCommand-Address: %w", err)
	} else {
		val.Address = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApCommand-ProprietaryId: %w", err)
	} else {
		if v != 255 {
			return nil, fmt.Errorf("match failed for SimnetApCommand-ProprietaryId: Expected %d != %d", 255, v)
		}
		val.ProprietaryId = SimnetEventCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApCommand-ApStatus: %w", err)
	} else {
		val.ApStatus = SimnetApStatusConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApCommand-ApCommand: %w", err)
	} else {
		val.ApCommand = SimnetApEventsConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApCommand-Direction: %w", err)
	} else {
		val.Direction = SimnetDirectionConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApCommand-Angle: %w", err)
	} else {
		val.Angle = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSimnetApCommand(val *SimnetApCommand) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Address, 8)
	w.writeReservedBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.ApStatus), 8)
	w.writeLookupField(uint64(val.ApCommand), 8)
	w.writeSpareBits(8)
	w.writeLookupField(uint64(val.Direction), 8)
	w.writeUnsignedResolution(val.Angle, 16, 0.0001)
	return w.Bytes(), w.Err()
}

func encodeSimnetApCommandMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetApCommand)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetApCommand, got %T", v)
	}
	return EncodeSimnetApCommand(val)
}

type SimnetEventCommandApCommand struct {
	Info              MessageInfo             `json:"info"`
	ManufacturerCode  ManufacturerCodeConst   `json:"manufacturerCode"`
	IndustryCode      IndustryCodeConst       `json:"industryCode"`
	ProprietaryId     SimnetEventCommandConst `json:"proprietaryId"`
	UnusedA           *uint16                 `json:"unusedA"`
	ControllingDevice *uint8                  `json:"controllingDevice"`
	Event             SimnetApEventsConst     `json:"event"`
	UnusedB           *uint8                  `json:"unusedB"`
	Direction         SimnetDirectionConst    `json:"direction"`
	Angle             *float32                `json:"angle"`
	UnusedC           *uint8                  `json:"unusedC"`
}

func (s *SimnetEventCommandApCommand) PGNNumber() uint32 { return 130850 }

func DecodeSimnetEventCommandApCommand(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetEventCommandApCommand{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetEventCommandApCommand-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetEventCommandApCommand-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-ProprietaryId: %w", err)
	} else {
		if v != 2 {
			return nil, fmt.Errorf("match failed for SimnetEventCommandApCommand-ProprietaryId: Expected %d != %d", 2, v)
		}
		val.ProprietaryId = SimnetEventCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-UnusedA: %w", err)
	} else {
		val.UnusedA = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-ControllingDevice: %w", err)
	} else {
		val.ControllingDevice = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-Event: %w", err)
	} else {
		val.Event = SimnetApEventsConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-UnusedB: %w", err)
	} else {
		val.UnusedB = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-Direction: %w", err)
	} else {
		val.Direction = SimnetDirectionConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-Angle: %w", err)
	} else {
		val.Angle = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-UnusedC: %w", err)
	} else {
		val.UnusedC = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSimnetEventCommandApCommand(val *SimnetEventCommandApCommand) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt16(val.UnusedA, 16)
	w.writeUInt8(val.ControllingDevice, 8)
	w.writeLookupField(uint64(val.Event), 8)
	w.writeUInt8(val.UnusedB, 8)
	w.writeLookupField(uint64(val.Direction), 8)
	w.writeUnsignedResolution(val.Angle, 16, 0.0001)
	w.writeUInt8(val.UnusedC, 8)
	return w.Bytes(), w.Err()
}

func encodeSimnetEventCommandApCommandMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetEventCommandApCommand)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetEventCommandApCommand, got %T", v)
	}
	return EncodeSimnetEventCommandApCommand(val)
}

type SimnetEventReplyApCommand struct {
	Info             MessageInfo             `json:"info"`
	ManufacturerCode ManufacturerCodeConst   `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst       `json:"industryCode"`
	ProprietaryId    SimnetEventCommandConst `json:"proprietaryId"`
	B                *uint16                 `json:"b"`
	Address          *uint8                  `json:"address"`
	Event            SimnetApEventsConst     `json:"event"`
	C                *uint8                  `json:"c"`
	Direction        SimnetDirectionConst    `json:"direction"`
	Angle            *float32                `json:"angle"`
	G                *uint8                  `json:"g"`
}

func (s *SimnetEventReplyApCommand) PGNNumber() uint32 { return 130851 }

func DecodeSimnetEventReplyApCommand(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetEventReplyApCommand{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetEventReplyApCommand-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetEventReplyApCommand-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-ProprietaryId: %w", err)
	} else {
		if v != 2 {
			return nil, fmt.Errorf("match failed for SimnetEventReplyApCommand-ProprietaryId: Expected %d != %d", 2, v)
		}
		val.ProprietaryId = SimnetEventCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-Address: %w", err)
	} else {
		val.Address = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-Event: %w", err)
	} else {
		val.Event = SimnetApEventsConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-Direction: %w", err)
	} else {
		val.Direction = SimnetDirectionConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-Angle: %w", err)
	} else {
		val.Angle = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-G: %w", err)
	} else {
		val.G = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSimnetEventReplyApCommand(val *SimnetEventReplyApCommand) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt16(val.B, 16)
	w.writeUInt8(val.Address, 8)
	w.writeLookupField(uint64(val.Event), 8)
	w.writeUInt8(val.C, 8)
	w.writeLookupField(uint64(val.Direction), 8)
	w.writeUnsignedResolution(val.Angle, 16, 0.0001)
	w.writeUInt8(val.G, 8)
	return w.Bytes(), w.Err()
}

func encodeSimnetEventReplyApCommandMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetEventReplyApCommand)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetEventReplyApCommand, got %T", v)
	}
	return EncodeSimnetEventReplyApCommand(val)
}

type SimnetAlarmMessage struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	MessageId        *uint16               `json:"messageId"`
	B                *uint8                `json:"b"`
	C                *uint8                `json:"c"`
	Text             string                `json:"text"`
}

func (s *SimnetAlarmMessage) PGNNumber() uint32 { return 130856 }

func DecodeSimnetAlarmMessage(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetAlarmMessage{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarmMessage-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetAlarmMessage-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetAlarmMessage-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetAlarmMessage-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarmMessage-MessageId: %w", err)
	} else {
		val.MessageId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarmMessage-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarmMessage-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(1784); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarmMessage-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSimnetAlarmMessage(val *SimnetAlarmMessage) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt16(val.MessageId, 16)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.C, 8)
	w.writeFixedString(val.Text, 1784)
	return w.Bytes(), w.Err()
}

func encodeSimnetAlarmMessageMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetAlarmMessage)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetAlarmMessage, got %T", v)
	}
	return EncodeSimnetAlarmMessage(val)
}

type SimnetApUnknown4 struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	A                *uint8                `json:"a"`
	B                *int32                `json:"b"`
	C                *int32                `json:"c"`
	D                *uint32               `json:"d"`
	E                *int32                `json:"e"`
	F                *uint32               `json:"f"`
}

func (s *SimnetApUnknown4) PGNNumber() uint32 { return 130860 }

func DecodeSimnetApUnknown4(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimnetApUnknown4{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown4-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetApUnknown4-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetApUnknown4-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetApUnknown4-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown4-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown4-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown4-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown4-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown4-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown4-F: %w", err)
	} else {
		val.F = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSimnetApUnknown4(val *SimnetApUnknown4) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.A, 8)
	w.writeInt32(val.B, 32)
	w.writeInt32(val.C, 32)
	w.writeUInt32(val.D, 32)
	w.writeInt32(val.E, 32)
	w.writeUInt32(val.F, 32)
	return w.Bytes(), w.Err()
}

func encodeSimnetApUnknown4Msg(v Message) ([]byte, error) {
	val, ok := v.(*SimnetApUnknown4)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetApUnknown4, got %T", v)
	}
	return EncodeSimnetApUnknown4(val)
}
