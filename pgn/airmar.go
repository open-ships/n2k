package pgn

import (
	"fmt"
	"github.com/open-ships/n2k/units"
)

type AirmarBootStateAcknowledgment struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	BootState        BootStateConst        `json:"bootState"`
}

func (a *AirmarBootStateAcknowledgment) PGNNumber() uint32 { return 65285 }

func DecodeAirmarBootStateAcknowledgment(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarBootStateAcknowledgment{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarBootStateAcknowledgment-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarBootStateAcknowledgment-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarBootStateAcknowledgment-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarBootStateAcknowledgment-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarBootStateAcknowledgment-BootState: %w", err)
	} else {
		val.BootState = BootStateConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(45)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeAirmarBootStateAcknowledgment(val *AirmarBootStateAcknowledgment) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.BootState), 3)
	w.skipBits(45)
	return w.Bytes(), w.Err()
}

func encodeAirmarBootStateAcknowledgmentMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarBootStateAcknowledgment)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarBootStateAcknowledgment, got %T", v)
	}
	return EncodeAirmarBootStateAcknowledgment(val)
}

type AirmarBootStateRequest struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
}

func (a *AirmarBootStateRequest) PGNNumber() uint32 { return 65286 }

func DecodeAirmarBootStateRequest(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarBootStateRequest{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarBootStateRequest-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarBootStateRequest-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarBootStateRequest-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarBootStateRequest-IndustryCode: Expected %d != %d", 4, v)
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

func EncodeAirmarBootStateRequest(val *AirmarBootStateRequest) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(48)
	return w.Bytes(), w.Err()
}

func encodeAirmarBootStateRequestMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarBootStateRequest)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarBootStateRequest, got %T", v)
	}
	return EncodeAirmarBootStateRequest(val)
}

type AirmarAccessLevel struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	FormatCode       *uint8                `json:"formatCode"`
	AccessLevel      AccessLevelConst      `json:"accessLevel"`
	AccessSeedKey    *uint32               `json:"accessSeedKey"`
}

func (a *AirmarAccessLevel) PGNNumber() uint32 { return 65287 }

func DecodeAirmarAccessLevel(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarAccessLevel{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAccessLevel-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarAccessLevel-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarAccessLevel-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarAccessLevel-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAccessLevel-FormatCode: %w", err)
	} else {
		val.FormatCode = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAccessLevel-AccessLevel: %w", err)
	} else {
		val.AccessLevel = AccessLevelConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(5)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAccessLevel-AccessSeedKey: %w", err)
	} else {
		val.AccessSeedKey = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeAirmarAccessLevel(val *AirmarAccessLevel) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.FormatCode, 8)
	w.writeLookupField(uint64(val.AccessLevel), 3)
	w.skipBits(5)
	w.writeUInt32(val.AccessSeedKey, 32)
	return w.Bytes(), w.Err()
}

func encodeAirmarAccessLevelMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarAccessLevel)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarAccessLevel, got %T", v)
	}
	return EncodeAirmarAccessLevel(val)
}

type AirmarDepthQualityFactor struct {
	Info               MessageInfo                   `json:"info"`
	ManufacturerCode   ManufacturerCodeConst         `json:"manufacturerCode"`
	IndustryCode       IndustryCodeConst             `json:"industryCode"`
	Sid                *uint8                        `json:"sid"`
	DepthQualityFactor AirmarDepthQualityFactorConst `json:"depthQualityFactor"`
}

func (a *AirmarDepthQualityFactor) PGNNumber() uint32 { return 65408 }

func DecodeAirmarDepthQualityFactor(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarDepthQualityFactor{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarDepthQualityFactor-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarDepthQualityFactor-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarDepthQualityFactor-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarDepthQualityFactor-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarDepthQualityFactor-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarDepthQualityFactor-DepthQualityFactor: %w", err)
	} else {
		val.DepthQualityFactor = AirmarDepthQualityFactorConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(36)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeAirmarDepthQualityFactor(val *AirmarDepthQualityFactor) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Sid, 8)
	w.writeLookupField(uint64(val.DepthQualityFactor), 4)
	w.skipBits(36)
	return w.Bytes(), w.Err()
}

func encodeAirmarDepthQualityFactorMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarDepthQualityFactor)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarDepthQualityFactor, got %T", v)
	}
	return EncodeAirmarDepthQualityFactor(val)
}

type AirmarSpeedPulseCount struct {
	Info                   MessageInfo           `json:"info"`
	ManufacturerCode       ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode           IndustryCodeConst     `json:"industryCode"`
	Sid                    *uint8                `json:"sid"`
	DurationOfInterval     *float32              `json:"durationOfInterval"`
	NumberOfPulsesReceived *uint16               `json:"numberOfPulsesReceived"`
}

func (a *AirmarSpeedPulseCount) PGNNumber() uint32 { return 65409 }

func DecodeAirmarSpeedPulseCount(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarSpeedPulseCount{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedPulseCount-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarSpeedPulseCount-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarSpeedPulseCount-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarSpeedPulseCount-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedPulseCount-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedPulseCount-DurationOfInterval: %w", err)
	} else {
		val.DurationOfInterval = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedPulseCount-NumberOfPulsesReceived: %w", err)
	} else {
		val.NumberOfPulsesReceived = v

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

func EncodeAirmarSpeedPulseCount(val *AirmarSpeedPulseCount) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Sid, 8)
	w.writeUnsignedResolution(val.DurationOfInterval, 16, 0.001)
	w.writeUInt16(val.NumberOfPulsesReceived, 16)
	w.skipBits(8)
	return w.Bytes(), w.Err()
}

func encodeAirmarSpeedPulseCountMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarSpeedPulseCount)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarSpeedPulseCount, got %T", v)
	}
	return EncodeAirmarSpeedPulseCount(val)
}

type AirmarDeviceInformation struct {
	Info                      MessageInfo           `json:"info"`
	ManufacturerCode          ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode              IndustryCodeConst     `json:"industryCode"`
	Sid                       *uint8                `json:"sid"`
	InternalDeviceTemperature *units.Temperature    `json:"internalDeviceTemperature"`
	SupplyVoltage             *float32              `json:"supplyVoltage"`
}

func (a *AirmarDeviceInformation) PGNNumber() uint32 { return 65410 }

func DecodeAirmarDeviceInformation(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarDeviceInformation{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarDeviceInformation-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarDeviceInformation-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarDeviceInformation-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarDeviceInformation-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarDeviceInformation-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarDeviceInformation-InternalDeviceTemperature: %w", err)
	} else {
		val.InternalDeviceTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarDeviceInformation-SupplyVoltage: %w", err)
	} else {
		val.SupplyVoltage = v

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

func EncodeAirmarDeviceInformation(val *AirmarDeviceInformation) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Sid, 8)
	var internalDeviceTemperatureRaw *float32
	if val.InternalDeviceTemperature != nil {
		internalDeviceTemperatureRaw = &val.InternalDeviceTemperature.Value
	}
	w.writeUnsignedResolution(internalDeviceTemperatureRaw, 16, 0.01)
	w.writeUnsignedResolution(val.SupplyVoltage, 16, 0.01)
	w.skipBits(8)
	return w.Bytes(), w.Err()
}

func encodeAirmarDeviceInformationMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarDeviceInformation)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarDeviceInformation, got %T", v)
	}
	return EncodeAirmarDeviceInformation(val)
}

type AirmarAddressableMultiFrame struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    *uint8                `json:"proprietaryId"`
}

func (a *AirmarAddressableMultiFrame) PGNNumber() uint32 { return 126720 }

func DecodeAirmarAddressableMultiFrame(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarAddressableMultiFrame{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAddressableMultiFrame-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarAddressableMultiFrame-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarAddressableMultiFrame-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarAddressableMultiFrame-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAddressableMultiFrame-ProprietaryId: %w", err)
	} else {
		val.ProprietaryId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeAirmarAddressableMultiFrame(val *AirmarAddressableMultiFrame) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.ProprietaryId, 8)
	return w.Bytes(), w.Err()
}

func encodeAirmarAddressableMultiFrameMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarAddressableMultiFrame)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarAddressableMultiFrame, got %T", v)
	}
	return EncodeAirmarAddressableMultiFrame(val)
}

type AirmarAttitudeOffset struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    AirmarCommandConst    `json:"proprietaryId"`
	AzimuthOffset    *float32              `json:"azimuthOffset"`
	PitchOffset      *float32              `json:"pitchOffset"`
	RollOffset       *float32              `json:"rollOffset"`
}

func (a *AirmarAttitudeOffset) PGNNumber() uint32 { return 126720 }

func DecodeAirmarAttitudeOffset(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarAttitudeOffset{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAttitudeOffset-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarAttitudeOffset-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarAttitudeOffset-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarAttitudeOffset-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAttitudeOffset-ProprietaryId: %w", err)
	} else {
		if v != 32 {
			return nil, fmt.Errorf("match failed for AirmarAttitudeOffset-ProprietaryId: Expected %d != %d", 32, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAttitudeOffset-AzimuthOffset: %w", err)
	} else {
		val.AzimuthOffset = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAttitudeOffset-PitchOffset: %w", err)
	} else {
		val.PitchOffset = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAttitudeOffset-RollOffset: %w", err)
	} else {
		val.RollOffset = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeAirmarAttitudeOffset(val *AirmarAttitudeOffset) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeSignedResolution(val.AzimuthOffset, 16, 0.0001)
	w.writeSignedResolution(val.PitchOffset, 16, 0.0001)
	w.writeSignedResolution(val.RollOffset, 16, 0.0001)
	return w.Bytes(), w.Err()
}

func encodeAirmarAttitudeOffsetMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarAttitudeOffset)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarAttitudeOffset, got %T", v)
	}
	return EncodeAirmarAttitudeOffset(val)
}

type AirmarCalibrateCompass struct {
	Info                   MessageInfo                  `json:"info"`
	ManufacturerCode       ManufacturerCodeConst        `json:"manufacturerCode"`
	IndustryCode           IndustryCodeConst            `json:"industryCode"`
	ProprietaryId          AirmarCommandConst           `json:"proprietaryId"`
	CalibrateFunction      AirmarCalibrateFunctionConst `json:"calibrateFunction"`
	CalibrationStatus      AirmarCalibrateStatusConst   `json:"calibrationStatus"`
	VerifyScore            *uint8                       `json:"verifyScore"`
	XAxisGainValue         *float32                     `json:"xAxisGainValue"`
	YAxisGainValue         *float32                     `json:"yAxisGainValue"`
	ZAxisGainValue         *float32                     `json:"zAxisGainValue"`
	XAxisLinearOffset      *float32                     `json:"xAxisLinearOffset"`
	YAxisLinearOffset      *float32                     `json:"yAxisLinearOffset"`
	ZAxisLinearOffset      *float32                     `json:"zAxisLinearOffset"`
	XAxisAngularOffset     *float32                     `json:"xAxisAngularOffset"`
	PitchAndRollDamping    *float32                     `json:"pitchAndRollDamping"`
	CompassRateGyroDamping *float32                     `json:"compassRateGyroDamping"`
}

func (a *AirmarCalibrateCompass) PGNNumber() uint32 { return 126720 }

func DecodeAirmarCalibrateCompass(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarCalibrateCompass{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateCompass-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateCompass-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-ProprietaryId: %w", err)
	} else {
		if v != 33 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateCompass-ProprietaryId: Expected %d != %d", 33, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-CalibrateFunction: %w", err)
	} else {
		val.CalibrateFunction = AirmarCalibrateFunctionConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-CalibrationStatus: %w", err)
	} else {
		val.CalibrationStatus = AirmarCalibrateStatusConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-VerifyScore: %w", err)
	} else {
		val.VerifyScore = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-XAxisGainValue: %w", err)
	} else {
		val.XAxisGainValue = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-YAxisGainValue: %w", err)
	} else {
		val.YAxisGainValue = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-ZAxisGainValue: %w", err)
	} else {
		val.ZAxisGainValue = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-XAxisLinearOffset: %w", err)
	} else {
		val.XAxisLinearOffset = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-YAxisLinearOffset: %w", err)
	} else {
		val.YAxisLinearOffset = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-ZAxisLinearOffset: %w", err)
	} else {
		val.ZAxisLinearOffset = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-XAxisAngularOffset: %w", err)
	} else {
		val.XAxisAngularOffset = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.05); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-PitchAndRollDamping: %w", err)
	} else {
		val.PitchAndRollDamping = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.05); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-CompassRateGyroDamping: %w", err)
	} else {
		val.CompassRateGyroDamping = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeAirmarCalibrateCompass(val *AirmarCalibrateCompass) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.CalibrateFunction), 8)
	w.writeLookupField(uint64(val.CalibrationStatus), 8)
	w.writeUInt8(val.VerifyScore, 8)
	w.writeSignedResolution(val.XAxisGainValue, 16, 0.01)
	w.writeSignedResolution(val.YAxisGainValue, 16, 0.01)
	w.writeSignedResolution(val.ZAxisGainValue, 16, 0.01)
	w.writeSignedResolution(val.XAxisLinearOffset, 16, 0.01)
	w.writeSignedResolution(val.YAxisLinearOffset, 16, 0.01)
	w.writeSignedResolution(val.ZAxisLinearOffset, 16, 0.01)
	w.writeSignedResolution(val.XAxisAngularOffset, 16, 0.1)
	w.writeSignedResolution(val.PitchAndRollDamping, 16, 0.05)
	w.writeSignedResolution(val.CompassRateGyroDamping, 16, 0.05)
	return w.Bytes(), w.Err()
}

func encodeAirmarCalibrateCompassMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarCalibrateCompass)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarCalibrateCompass, got %T", v)
	}
	return EncodeAirmarCalibrateCompass(val)
}

type AirmarCalibrateDepth struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    AirmarCommandConst    `json:"proprietaryId"`
	SpeedOfSoundMode *units.Velocity       `json:"speedOfSoundMode"`
}

func (a *AirmarCalibrateDepth) PGNNumber() uint32 { return 126720 }

func DecodeAirmarCalibrateDepth(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarCalibrateDepth{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateDepth-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateDepth-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarCalibrateDepth-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateDepth-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateDepth-ProprietaryId: %w", err)
	} else {
		if v != 40 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateDepth-ProprietaryId: Expected %d != %d", 40, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateDepth-SpeedOfSoundMode: %w", err)
	} else {
		val.SpeedOfSoundMode = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

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

func EncodeAirmarCalibrateDepth(val *AirmarCalibrateDepth) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	var speedOfSoundModeRaw *float32
	if val.SpeedOfSoundMode != nil {
		speedOfSoundModeRaw = &val.SpeedOfSoundMode.Value
	}
	w.writeUnsignedResolution(speedOfSoundModeRaw, 16, 0.1)
	w.skipBits(8)
	return w.Bytes(), w.Err()
}

func encodeAirmarCalibrateDepthMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarCalibrateDepth)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarCalibrateDepth, got %T", v)
	}
	return EncodeAirmarCalibrateDepth(val)
}

type AirmarCalibrateSpeed struct {
	Info                      MessageInfo                      `json:"info"`
	ManufacturerCode          ManufacturerCodeConst            `json:"manufacturerCode"`
	IndustryCode              IndustryCodeConst                `json:"industryCode"`
	ProprietaryId             AirmarCommandConst               `json:"proprietaryId"`
	NumberOfPairsOfDataPoints *uint8                           `json:"numberOfPairsOfDataPoints"`
	Repeating1                []AirmarCalibrateSpeedRepeating1 `json:"repeating1"`
}

type AirmarCalibrateSpeedRepeating1 struct {
	InputFrequency *float32        `json:"inputFrequency"`
	OutputSpeed    *units.Velocity `json:"outputSpeed"`
}

func (a *AirmarCalibrateSpeed) PGNNumber() uint32 { return 126720 }

func DecodeAirmarCalibrateSpeed(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarCalibrateSpeed{}
	val.Info = Info
	var repeat1Count uint16 = 0
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateSpeed-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateSpeed-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarCalibrateSpeed-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateSpeed-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateSpeed-ProprietaryId: %w", err)
	} else {
		if v != 41 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateSpeed-ProprietaryId: Expected %d != %d", 41, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateSpeed-NumberOfPairsOfDataPoints: %w", err)
	} else {
		val.NumberOfPairsOfDataPoints = v
		if v != nil {
			repeat1Count = uint16(*v)
		}

		if stream.isEOF() {
			return val, nil
		}
	}
	if repeat1Count == 0 {
		return val, nil
	}
	val.Repeating1 = make([]AirmarCalibrateSpeedRepeating1, 0)
	i := 0
	for {
		var rep AirmarCalibrateSpeedRepeating1
		if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
			return nil, fmt.Errorf("parse failed for AirmarCalibrateSpeed-InputFrequency: %w", err)
		} else {
			rep.InputFrequency = v
		}
		if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
			return nil, fmt.Errorf("parse failed for AirmarCalibrateSpeed-OutputSpeed: %w", err)
		} else {
			rep.OutputSpeed = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)
		}
		val.Repeating1 = append(val.Repeating1, rep)
		if int(repeat1Count) == 0 {
			if stream.isEOF() {
				return val, nil
			}
		} else {
			i++
			if i == int(repeat1Count) {
				break
			}
		}
	}
	return val, nil
}

func EncodeAirmarCalibrateSpeed(val *AirmarCalibrateSpeed) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt8(val.NumberOfPairsOfDataPoints, 8)
	for _, rep := range val.Repeating1 {
		w.writeUnsignedResolution(rep.InputFrequency, 16, 0.1)
		var outputSpeedRaw *float32
		if rep.OutputSpeed != nil {
			outputSpeedRaw = &rep.OutputSpeed.Value
		}
		w.writeUnsignedResolution(outputSpeedRaw, 16, 0.01)
	}
	return w.Bytes(), w.Err()
}

func encodeAirmarCalibrateSpeedMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarCalibrateSpeed)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarCalibrateSpeed, got %T", v)
	}
	return EncodeAirmarCalibrateSpeed(val)
}

type AirmarCalibrateTemperature struct {
	Info                MessageInfo                    `json:"info"`
	ManufacturerCode    ManufacturerCodeConst          `json:"manufacturerCode"`
	IndustryCode        IndustryCodeConst              `json:"industryCode"`
	ProprietaryId       AirmarCommandConst             `json:"proprietaryId"`
	TemperatureInstance AirmarTemperatureInstanceConst `json:"temperatureInstance"`
	TemperatureOffset   *units.Temperature             `json:"temperatureOffset"`
}

func (a *AirmarCalibrateTemperature) PGNNumber() uint32 { return 126720 }

func DecodeAirmarCalibrateTemperature(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarCalibrateTemperature{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateTemperature-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateTemperature-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarCalibrateTemperature-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateTemperature-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateTemperature-ProprietaryId: %w", err)
	} else {
		if v != 42 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateTemperature-ProprietaryId: Expected %d != %d", 42, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateTemperature-TemperatureInstance: %w", err)
	} else {
		val.TemperatureInstance = AirmarTemperatureInstanceConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readSignedResolution(16, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateTemperature-TemperatureOffset: %w", err)
	} else {
		val.TemperatureOffset = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeAirmarCalibrateTemperature(val *AirmarCalibrateTemperature) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.TemperatureInstance), 2)
	w.skipBits(6)
	var temperatureOffsetRaw *float32
	if val.TemperatureOffset != nil {
		temperatureOffsetRaw = &val.TemperatureOffset.Value
	}
	w.writeSignedResolution(temperatureOffsetRaw, 16, 0.001)
	return w.Bytes(), w.Err()
}

func encodeAirmarCalibrateTemperatureMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarCalibrateTemperature)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarCalibrateTemperature, got %T", v)
	}
	return EncodeAirmarCalibrateTemperature(val)
}

type AirmarNmea2000Options struct {
	Info                 MessageInfo                     `json:"info"`
	ManufacturerCode     ManufacturerCodeConst           `json:"manufacturerCode"`
	IndustryCode         IndustryCodeConst               `json:"industryCode"`
	ProprietaryId        AirmarCommandConst              `json:"proprietaryId"`
	TransmissionInterval AirmarTransmissionIntervalConst `json:"transmissionInterval"`
}

func (a *AirmarNmea2000Options) PGNNumber() uint32 { return 126720 }

func DecodeAirmarNmea2000Options(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarNmea2000Options{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarNmea2000Options-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarNmea2000Options-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarNmea2000Options-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarNmea2000Options-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarNmea2000Options-ProprietaryId: %w", err)
	} else {
		if v != 46 {
			return nil, fmt.Errorf("match failed for AirmarNmea2000Options-ProprietaryId: Expected %d != %d", 46, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarNmea2000Options-TransmissionInterval: %w", err)
	} else {
		val.TransmissionInterval = AirmarTransmissionIntervalConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(22)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeAirmarNmea2000Options(val *AirmarNmea2000Options) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.TransmissionInterval), 2)
	w.skipBits(22)
	return w.Bytes(), w.Err()
}

func encodeAirmarNmea2000OptionsMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarNmea2000Options)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarNmea2000Options, got %T", v)
	}
	return EncodeAirmarNmea2000Options(val)
}

type AirmarSimulateMode struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    AirmarCommandConst    `json:"proprietaryId"`
	SimulateMode     OffOnConst            `json:"simulateMode"`
}

func (a *AirmarSimulateMode) PGNNumber() uint32 { return 126720 }

func DecodeAirmarSimulateMode(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarSimulateMode{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSimulateMode-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarSimulateMode-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarSimulateMode-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarSimulateMode-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSimulateMode-ProprietaryId: %w", err)
	} else {
		if v != 35 {
			return nil, fmt.Errorf("match failed for AirmarSimulateMode-ProprietaryId: Expected %d != %d", 35, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSimulateMode-SimulateMode: %w", err)
	} else {
		val.SimulateMode = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(22)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeAirmarSimulateMode(val *AirmarSimulateMode) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.SimulateMode), 2)
	w.skipBits(22)
	return w.Bytes(), w.Err()
}

func encodeAirmarSimulateModeMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarSimulateMode)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarSimulateMode, got %T", v)
	}
	return EncodeAirmarSimulateMode(val)
}

type AirmarSpeedFilterIir struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    AirmarCommandConst    `json:"proprietaryId"`
	FilterType       *uint8                `json:"filterType"`
	SampleInterval   *float32              `json:"sampleInterval"`
	FilterDuration   *float32              `json:"filterDuration"`
}

func (a *AirmarSpeedFilterIir) PGNNumber() uint32 { return 126720 }

func DecodeAirmarSpeedFilterIir(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarSpeedFilterIir{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterIir-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarSpeedFilterIir-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterIir-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarSpeedFilterIir-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterIir-ProprietaryId: %w", err)
	} else {
		if v != 43 {
			return nil, fmt.Errorf("match failed for AirmarSpeedFilterIir-ProprietaryId: Expected %d != %d", 43, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterIir-FilterType: %w", err)
	} else {
		if v != nil && *v != 1 {
			return nil, fmt.Errorf("match failed for AirmarSpeedFilterIir-FilterType: Expected %d != %d", 1, *v)
		}
		val.FilterType = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterIir-SampleInterval: %w", err)
	} else {
		val.SampleInterval = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterIir-FilterDuration: %w", err)
	} else {
		val.FilterDuration = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeAirmarSpeedFilterIir(val *AirmarSpeedFilterIir) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt8(val.FilterType, 4)
	w.skipBits(4)
	w.writeUnsignedResolution(val.SampleInterval, 16, 0.01)
	w.writeUnsignedResolution(val.FilterDuration, 16, 0.01)
	return w.Bytes(), w.Err()
}

func encodeAirmarSpeedFilterIirMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarSpeedFilterIir)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarSpeedFilterIir, got %T", v)
	}
	return EncodeAirmarSpeedFilterIir(val)
}

type AirmarSpeedFilterNone struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    AirmarCommandConst    `json:"proprietaryId"`
	FilterType       *uint8                `json:"filterType"`
	SampleInterval   *float32              `json:"sampleInterval"`
}

func (a *AirmarSpeedFilterNone) PGNNumber() uint32 { return 126720 }

func DecodeAirmarSpeedFilterNone(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarSpeedFilterNone{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterNone-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarSpeedFilterNone-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterNone-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarSpeedFilterNone-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterNone-ProprietaryId: %w", err)
	} else {
		if v != 43 {
			return nil, fmt.Errorf("match failed for AirmarSpeedFilterNone-ProprietaryId: Expected %d != %d", 43, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterNone-FilterType: %w", err)
	} else {
		if v != nil && *v != 0 {
			return nil, fmt.Errorf("match failed for AirmarSpeedFilterNone-FilterType: Expected %d != %d", 0, *v)
		}
		val.FilterType = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterNone-SampleInterval: %w", err)
	} else {
		val.SampleInterval = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeAirmarSpeedFilterNone(val *AirmarSpeedFilterNone) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt8(val.FilterType, 4)
	w.skipBits(4)
	w.writeUnsignedResolution(val.SampleInterval, 16, 0.01)
	return w.Bytes(), w.Err()
}

func encodeAirmarSpeedFilterNoneMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarSpeedFilterNone)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarSpeedFilterNone, got %T", v)
	}
	return EncodeAirmarSpeedFilterNone(val)
}

type AirmarTemperatureFilterIir struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    AirmarCommandConst    `json:"proprietaryId"`
	FilterType       *uint8                `json:"filterType"`
	SampleInterval   *float32              `json:"sampleInterval"`
	FilterDuration   *float32              `json:"filterDuration"`
}

func (a *AirmarTemperatureFilterIir) PGNNumber() uint32 { return 126720 }

func DecodeAirmarTemperatureFilterIir(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarTemperatureFilterIir{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterIir-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarTemperatureFilterIir-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterIir-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarTemperatureFilterIir-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterIir-ProprietaryId: %w", err)
	} else {
		if v != 44 {
			return nil, fmt.Errorf("match failed for AirmarTemperatureFilterIir-ProprietaryId: Expected %d != %d", 44, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterIir-FilterType: %w", err)
	} else {
		if v != nil && *v != 1 {
			return nil, fmt.Errorf("match failed for AirmarTemperatureFilterIir-FilterType: Expected %d != %d", 1, *v)
		}
		val.FilterType = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterIir-SampleInterval: %w", err)
	} else {
		val.SampleInterval = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterIir-FilterDuration: %w", err)
	} else {
		val.FilterDuration = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeAirmarTemperatureFilterIir(val *AirmarTemperatureFilterIir) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt8(val.FilterType, 4)
	w.skipBits(4)
	w.writeUnsignedResolution(val.SampleInterval, 16, 0.01)
	w.writeUnsignedResolution(val.FilterDuration, 16, 0.01)
	return w.Bytes(), w.Err()
}

func encodeAirmarTemperatureFilterIirMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarTemperatureFilterIir)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarTemperatureFilterIir, got %T", v)
	}
	return EncodeAirmarTemperatureFilterIir(val)
}

type AirmarTemperatureFilterNone struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    AirmarCommandConst    `json:"proprietaryId"`
	FilterType       *uint8                `json:"filterType"`
	SampleInterval   *float32              `json:"sampleInterval"`
}

func (a *AirmarTemperatureFilterNone) PGNNumber() uint32 { return 126720 }

func DecodeAirmarTemperatureFilterNone(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarTemperatureFilterNone{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterNone-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarTemperatureFilterNone-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterNone-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarTemperatureFilterNone-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterNone-ProprietaryId: %w", err)
	} else {
		if v != 44 {
			return nil, fmt.Errorf("match failed for AirmarTemperatureFilterNone-ProprietaryId: Expected %d != %d", 44, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterNone-FilterType: %w", err)
	} else {
		if v != nil && *v != 0 {
			return nil, fmt.Errorf("match failed for AirmarTemperatureFilterNone-FilterType: Expected %d != %d", 0, *v)
		}
		val.FilterType = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterNone-SampleInterval: %w", err)
	} else {
		val.SampleInterval = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeAirmarTemperatureFilterNone(val *AirmarTemperatureFilterNone) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt8(val.FilterType, 4)
	w.skipBits(4)
	w.writeUnsignedResolution(val.SampleInterval, 16, 0.01)
	return w.Bytes(), w.Err()
}

func encodeAirmarTemperatureFilterNoneMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarTemperatureFilterNone)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarTemperatureFilterNone, got %T", v)
	}
	return EncodeAirmarTemperatureFilterNone(val)
}

type AirmarTrueWindOptions struct {
	Info                  MessageInfo           `json:"info"`
	ManufacturerCode      ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode          IndustryCodeConst     `json:"industryCode"`
	ProprietaryId         AirmarCommandConst    `json:"proprietaryId"`
	CogSubstitutionForHdg YesNoConst            `json:"cogSubstitutionForHdg"`
}

func (a *AirmarTrueWindOptions) PGNNumber() uint32 { return 126720 }

func DecodeAirmarTrueWindOptions(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AirmarTrueWindOptions{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTrueWindOptions-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarTrueWindOptions-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarTrueWindOptions-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarTrueWindOptions-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTrueWindOptions-ProprietaryId: %w", err)
	} else {
		if v != 34 {
			return nil, fmt.Errorf("match failed for AirmarTrueWindOptions-ProprietaryId: Expected %d != %d", 34, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTrueWindOptions-CogSubstitutionForHdg: %w", err)
	} else {
		val.CogSubstitutionForHdg = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(22)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeAirmarTrueWindOptions(val *AirmarTrueWindOptions) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.CogSubstitutionForHdg), 2)
	w.skipBits(22)
	return w.Bytes(), w.Err()
}

func encodeAirmarTrueWindOptionsMsg(v Message) ([]byte, error) {
	val, ok := v.(*AirmarTrueWindOptions)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarTrueWindOptions, got %T", v)
	}
	return EncodeAirmarTrueWindOptions(val)
}
