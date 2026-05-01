package pgn

import (
	"fmt"
)

type ChetcoDimmer struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Instance *uint8 `json:"instance"`
	Dimmer1 *uint8 `json:"dimmer1"`
	Dimmer2 *uint8 `json:"dimmer2"`
	Dimmer3 *uint8 `json:"dimmer3"`
	Dimmer4 *uint8 `json:"dimmer4"`
	Control *uint8 `json:"control"`
}
func (x *ChetcoDimmer) PGNNumber() uint32  { return 65286 }
func DecodeChetcoDimmer(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &ChetcoDimmer{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for ChetcoDimmer-ManufacturerCode: %w", err)
	} else {
		if v != 409 {
			return nil, fmt.Errorf("match failed for ChetcoDimmer-ManufacturerCode: Expected %d != %d", 409, v)
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
		return nil, fmt.Errorf("parse failed for ChetcoDimmer-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for ChetcoDimmer-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ChetcoDimmer-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ChetcoDimmer-Dimmer1: %w", err)
	} else {
		val.Dimmer1 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ChetcoDimmer-Dimmer2: %w", err)
	} else {
		val.Dimmer2 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ChetcoDimmer-Dimmer3: %w", err)
	} else {
		val.Dimmer3 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ChetcoDimmer-Dimmer4: %w", err)
	} else {
		val.Dimmer4 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ChetcoDimmer-Control: %w", err)
	} else {
		val.Control = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeChetcoDimmer(val *ChetcoDimmer) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Instance, 8)
	w.writeUInt8(val.Dimmer1, 8)
	w.writeUInt8(val.Dimmer2, 8)
	w.writeUInt8(val.Dimmer3, 8)
	w.writeUInt8(val.Dimmer4, 8)
	w.writeUInt8(val.Control, 8)
	return w.Bytes(), w.Err()
}
func encodeChetcoDimmerMsg(v Message) ([]byte, error) {
	val, ok := v.(*ChetcoDimmer)
	if !ok {
		return nil, fmt.Errorf("expected *ChetcoDimmer, got %T", v)
	}
	return EncodeChetcoDimmer(val)
}
