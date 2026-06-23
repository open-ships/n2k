package pgn

import "fmt"

type SimradTextMessage struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    SimnetCommandConst    `json:"proprietaryId"`
	A                *uint8                `json:"a"`
	B                *uint8                `json:"b"`
	C                *uint8                `json:"c"`
	Sid              *uint8                `json:"sid"`
	Prio             *uint8                `json:"prio"`
	Text             string                `json:"text"`
}

func (s *SimradTextMessage) PGNNumber() uint32 { return 130816 }

func DecodeSimradTextMessage(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SimradTextMessage{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimradTextMessage-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimradTextMessage-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimradTextMessage-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimradTextMessage-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimradTextMessage-ProprietaryId: %w", err)
	} else {
		if v != 50 {
			return nil, fmt.Errorf("match failed for SimradTextMessage-ProprietaryId: Expected %d != %d", 50, v)
		}
		val.ProprietaryId = SimnetCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimradTextMessage-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimradTextMessage-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimradTextMessage-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimradTextMessage-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimradTextMessage-Prio: %w", err)
	} else {
		val.Prio = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for SimradTextMessage-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSimradTextMessage(val *SimradTextMessage) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.C, 8)
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.Prio, 8)
	w.writeFixedString(val.Text, 256)
	return w.Bytes(), w.Err()
}

func encodeSimradTextMessageMsg(v Message) ([]byte, error) {
	val, ok := v.(*SimradTextMessage)
	if !ok {
		return nil, fmt.Errorf("expected *SimradTextMessage, got %T", v)
	}
	return EncodeSimradTextMessage(val)
}
