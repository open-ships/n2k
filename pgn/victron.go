package pgn

import "fmt"

type VictronBatteryRegister struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	RegisterId       *uint16               `json:"registerId"`
	Payload          *uint32               `json:"payload"`
}

func (v *VictronBatteryRegister) PGNNumber() uint32 { return 61184 }

func DecodeVictronBatteryRegister(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &VictronBatteryRegister{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for VictronBatteryRegister-ManufacturerCode: %w", err)
	} else {
		if v != 358 {
			return nil, fmt.Errorf("match failed for VictronBatteryRegister-ManufacturerCode: Expected %d != %d", 358, v)
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
		return nil, fmt.Errorf("parse failed for VictronBatteryRegister-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for VictronBatteryRegister-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for VictronBatteryRegister-RegisterId: %w", err)
	} else {
		val.RegisterId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for VictronBatteryRegister-Payload: %w", err)
	} else {
		val.Payload = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeVictronBatteryRegister(val *VictronBatteryRegister) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt16(val.RegisterId, 16)
	w.writeUInt32(val.Payload, 32)
	return w.Bytes(), w.Err()
}

func encodeVictronBatteryRegisterMsg(v Message) ([]byte, error) {
	val, ok := v.(*VictronBatteryRegister)
	if !ok {
		return nil, fmt.Errorf("expected *VictronBatteryRegister, got %T", v)
	}
	return EncodeVictronBatteryRegister(val)
}
