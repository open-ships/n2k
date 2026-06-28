package pgn

import (
	"fmt"
	"github.com/open-ships/n2k/units"
)

type FurunoHeave struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	Heave            *units.Distance       `json:"heave"`
}

func (f *FurunoHeave) PGNNumber() uint32 { return 65280 }

func DecodeFurunoHeave(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &FurunoHeave{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoHeave-ManufacturerCode: %w", err)
	} else {
		if v != 1855 {
			return nil, fmt.Errorf("match failed for FurunoHeave-ManufacturerCode: Expected %d != %d", 1855, v)
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
		return nil, fmt.Errorf("parse failed for FurunoHeave-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FurunoHeave-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(32, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoHeave-Heave: %w", err)
	} else {
		val.Heave = nullableUnit(units.Meter, v, units.NewDistance)

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

func EncodeFurunoHeave(val *FurunoHeave) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	var heaveRaw *float32
	if val.Heave != nil {
		heaveRaw = &val.Heave.Value
	}
	w.writeSignedResolution(heaveRaw, 32, 0.001)
	w.writeReservedBits(16)
	return w.Bytes(), w.Err()
}

func encodeFurunoHeaveMsg(v Message) ([]byte, error) {
	val, ok := v.(*FurunoHeave)
	if !ok {
		return nil, fmt.Errorf("expected *FurunoHeave, got %T", v)
	}
	return EncodeFurunoHeave(val)
}

type FurunoUnknown130820 struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	A                *uint8                `json:"a"`
	B                *uint8                `json:"b"`
	C                *uint8                `json:"c"`
	D                *uint8                `json:"d"`
	E                *uint8                `json:"e"`
}

func (f *FurunoUnknown130820) PGNNumber() uint32 { return 130820 }

func DecodeFurunoUnknown130820(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &FurunoUnknown130820{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130820-ManufacturerCode: %w", err)
	} else {
		if v != 1855 {
			return nil, fmt.Errorf("match failed for FurunoUnknown130820-ManufacturerCode: Expected %d != %d", 1855, v)
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
		return nil, fmt.Errorf("parse failed for FurunoUnknown130820-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FurunoUnknown130820-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130820-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130820-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130820-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130820-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130820-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeFurunoUnknown130820(val *FurunoUnknown130820) ([]byte, error) {
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
	return w.Bytes(), w.Err()
}

func encodeFurunoUnknown130820Msg(v Message) ([]byte, error) {
	val, ok := v.(*FurunoUnknown130820)
	if !ok {
		return nil, fmt.Errorf("expected *FurunoUnknown130820, got %T", v)
	}
	return EncodeFurunoUnknown130820(val)
}

type FurunoUnknown130821 struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	Sid              *uint8                `json:"sid"`
	A                *uint8                `json:"a"`
	B                *uint8                `json:"b"`
	C                *uint8                `json:"c"`
	D                *uint8                `json:"d"`
	E                *uint8                `json:"e"`
	F                *uint8                `json:"f"`
	G                *uint8                `json:"g"`
	H                *uint8                `json:"h"`
	I                *uint8                `json:"i"`
}

func (f *FurunoUnknown130821) PGNNumber() uint32 { return 130821 }

func DecodeFurunoUnknown130821(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &FurunoUnknown130821{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-ManufacturerCode: %w", err)
	} else {
		if v != 1855 {
			return nil, fmt.Errorf("match failed for FurunoUnknown130821-ManufacturerCode: Expected %d != %d", 1855, v)
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
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FurunoUnknown130821-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-F: %w", err)
	} else {
		val.F = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-G: %w", err)
	} else {
		val.G = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-H: %w", err)
	} else {
		val.H = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-I: %w", err)
	} else {
		val.I = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeFurunoUnknown130821(val *FurunoUnknown130821) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.C, 8)
	w.writeUInt8(val.D, 8)
	w.writeUInt8(val.E, 8)
	w.writeUInt8(val.F, 8)
	w.writeUInt8(val.G, 8)
	w.writeUInt8(val.H, 8)
	w.writeUInt8(val.I, 8)
	return w.Bytes(), w.Err()
}

func encodeFurunoUnknown130821Msg(v Message) ([]byte, error) {
	val, ok := v.(*FurunoUnknown130821)
	if !ok {
		return nil, fmt.Errorf("expected *FurunoUnknown130821, got %T", v)
	}
	return EncodeFurunoUnknown130821(val)
}

type FurunoSixDegreesOfFreedomMovement struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	A                *int32                `json:"a"`
	B                *int32                `json:"b"`
	C                *int32                `json:"c"`
	D                *int8                 `json:"d"`
	E                *int32                `json:"e"`
	F                *int32                `json:"f"`
	G                *int16                `json:"g"`
	H                *int16                `json:"h"`
	I                *int16                `json:"i"`
}

func (f *FurunoSixDegreesOfFreedomMovement) PGNNumber() uint32 { return 130842 }

func DecodeFurunoSixDegreesOfFreedomMovement(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &FurunoSixDegreesOfFreedomMovement{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-ManufacturerCode: %w", err)
	} else {
		if v != 1855 {
			return nil, fmt.Errorf("match failed for FurunoSixDegreesOfFreedomMovement-ManufacturerCode: Expected %d != %d", 1855, v)
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
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FurunoSixDegreesOfFreedomMovement-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-F: %w", err)
	} else {
		val.F = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-G: %w", err)
	} else {
		val.G = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-H: %w", err)
	} else {
		val.H = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-I: %w", err)
	} else {
		val.I = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeFurunoSixDegreesOfFreedomMovement(val *FurunoSixDegreesOfFreedomMovement) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeInt32(val.A, 32)
	w.writeInt32(val.B, 32)
	w.writeInt32(val.C, 32)
	w.writeInt8(val.D, 8)
	w.writeInt32(val.E, 32)
	w.writeInt32(val.F, 32)
	w.writeInt16(val.G, 16)
	w.writeInt16(val.H, 16)
	w.writeInt16(val.I, 16)
	return w.Bytes(), w.Err()
}

func encodeFurunoSixDegreesOfFreedomMovementMsg(v Message) ([]byte, error) {
	val, ok := v.(*FurunoSixDegreesOfFreedomMovement)
	if !ok {
		return nil, fmt.Errorf("expected *FurunoSixDegreesOfFreedomMovement, got %T", v)
	}
	return EncodeFurunoSixDegreesOfFreedomMovement(val)
}

type FurunoHeelAngleRollInformation struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	A                *uint8                `json:"a"`
	B                *uint8                `json:"b"`
	Yaw              *float32              `json:"yaw"`
	Pitch            *float32              `json:"pitch"`
	Roll             *float32              `json:"roll"`
}

func (f *FurunoHeelAngleRollInformation) PGNNumber() uint32 { return 130843 }

func DecodeFurunoHeelAngleRollInformation(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &FurunoHeelAngleRollInformation{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoHeelAngleRollInformation-ManufacturerCode: %w", err)
	} else {
		if v != 1855 {
			return nil, fmt.Errorf("match failed for FurunoHeelAngleRollInformation-ManufacturerCode: Expected %d != %d", 1855, v)
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
		return nil, fmt.Errorf("parse failed for FurunoHeelAngleRollInformation-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FurunoHeelAngleRollInformation-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoHeelAngleRollInformation-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoHeelAngleRollInformation-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoHeelAngleRollInformation-Yaw: %w", err)
	} else {
		val.Yaw = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoHeelAngleRollInformation-Pitch: %w", err)
	} else {
		val.Pitch = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoHeelAngleRollInformation-Roll: %w", err)
	} else {
		val.Roll = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeFurunoHeelAngleRollInformation(val *FurunoHeelAngleRollInformation) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeSignedResolution(val.Yaw, 16, 0.0001)
	w.writeSignedResolution(val.Pitch, 16, 0.0001)
	w.writeSignedResolution(val.Roll, 16, 0.0001)
	return w.Bytes(), w.Err()
}

func encodeFurunoHeelAngleRollInformationMsg(v Message) ([]byte, error) {
	val, ok := v.(*FurunoHeelAngleRollInformation)
	if !ok {
		return nil, fmt.Errorf("expected *FurunoHeelAngleRollInformation, got %T", v)
	}
	return EncodeFurunoHeelAngleRollInformation(val)
}

type FurunoMultiSatsInViewExtended struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
}

func (f *FurunoMultiSatsInViewExtended) PGNNumber() uint32 { return 130845 }

func DecodeFurunoMultiSatsInViewExtended(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &FurunoMultiSatsInViewExtended{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoMultiSatsInViewExtended-ManufacturerCode: %w", err)
	} else {
		if v != 1855 {
			return nil, fmt.Errorf("match failed for FurunoMultiSatsInViewExtended-ManufacturerCode: Expected %d != %d", 1855, v)
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
		return nil, fmt.Errorf("parse failed for FurunoMultiSatsInViewExtended-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FurunoMultiSatsInViewExtended-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeFurunoMultiSatsInViewExtended(val *FurunoMultiSatsInViewExtended) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	return w.Bytes(), w.Err()
}

func encodeFurunoMultiSatsInViewExtendedMsg(v Message) ([]byte, error) {
	val, ok := v.(*FurunoMultiSatsInViewExtended)
	if !ok {
		return nil, fmt.Errorf("expected *FurunoMultiSatsInViewExtended, got %T", v)
	}
	return EncodeFurunoMultiSatsInViewExtended(val)
}

type FurunoMotionSensorStatusExtended struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
}

func (f *FurunoMotionSensorStatusExtended) PGNNumber() uint32 { return 130846 }

func DecodeFurunoMotionSensorStatusExtended(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &FurunoMotionSensorStatusExtended{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoMotionSensorStatusExtended-ManufacturerCode: %w", err)
	} else {
		if v != 1855 {
			return nil, fmt.Errorf("match failed for FurunoMotionSensorStatusExtended-ManufacturerCode: Expected %d != %d", 1855, v)
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
		return nil, fmt.Errorf("parse failed for FurunoMotionSensorStatusExtended-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FurunoMotionSensorStatusExtended-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeFurunoMotionSensorStatusExtended(val *FurunoMotionSensorStatusExtended) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	return w.Bytes(), w.Err()
}

func encodeFurunoMotionSensorStatusExtendedMsg(v Message) ([]byte, error) {
	val, ok := v.(*FurunoMotionSensorStatusExtended)
	if !ok {
		return nil, fmt.Errorf("expected *FurunoMotionSensorStatusExtended, got %T", v)
	}
	return EncodeFurunoMotionSensorStatusExtended(val)
}
