package pgn

import "fmt"

type BGKeyValueData struct {
	Info             MessageInfo                `json:"info"`
	ManufacturerCode ManufacturerCodeConst      `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst          `json:"industryCode"`
	Repeating1       []BGKeyValueDataRepeating1 `json:"repeating1"`
}

type BGKeyValueDataRepeating1 struct {
	Key    BandgKeyValueConst `json:"key"`
	Length *uint8             `json:"length"`
	Value  []uint8            `json:"value"`
}

func (b *BGKeyValueData) PGNNumber() uint32 { return 130824 }

func DecodeBGKeyValueData(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &BGKeyValueData{}
	val.Info = Info
	var repeat1Count uint16 = 0
	var valueLength uint16

	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for BGKeyValueData-ManufacturerCode: %w", err)
	} else {
		if v != 381 {
			return nil, fmt.Errorf("match failed for BGKeyValueData-ManufacturerCode: Expected %d != %d", 381, v)
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
		return nil, fmt.Errorf("parse failed for BGKeyValueData-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for BGKeyValueData-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	val.Repeating1 = make([]BGKeyValueDataRepeating1, 0)
	i := 0
	for {
		var rep BGKeyValueDataRepeating1
		if v, err := stream.readLookupField(12); err != nil {
			return nil, fmt.Errorf("parse failed for BGKeyValueData-Key: %w", err)
		} else {
			rep.Key = BandgKeyValueConst(v)
		}
		if v, err := stream.readUInt8(4); err != nil {
			return nil, fmt.Errorf("parse failed for BGKeyValueData-Length: %w", err)
		} else {
			rep.Length = v
			if v != nil {
				valueLength = uint16(*v) * 8
			}

		}
		if v, err := stream.readBinaryData(valueLength); err != nil {
			return nil, fmt.Errorf("parse failed for BGKeyValueData-Value: %w", err)
		} else {
			rep.Value = v
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

func EncodeBGKeyValueData(val *BGKeyValueData) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	for _, rep := range val.Repeating1 {
		w.writeLookupField(uint64(rep.Key), 12)
		w.writeUInt8(rep.Length, 4)
		var valueLength uint16
		if rep.Length != nil {
			valueLength = uint16(*rep.Length) * 8
		}
		w.writeBinaryData(rep.Value, valueLength)
	}
	return w.Bytes(), w.Err()
}

func encodeBGKeyValueDataMsg(v Message) ([]byte, error) {
	val, ok := v.(*BGKeyValueData)
	if !ok {
		return nil, fmt.Errorf("expected *BGKeyValueData, got %T", v)
	}
	return EncodeBGKeyValueData(val)
}

type BGUserAndRemoteRename struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	DataType         BandgKeyValueConst    `json:"dataType"`
	Length           *uint8                `json:"length"`
	Decimals         BandgDecimalsConst    `json:"decimals"`
	ShortName        string                `json:"shortName"`
	LongName         string                `json:"longName"`
}

func (b *BGUserAndRemoteRename) PGNNumber() uint32 { return 130833 }

func DecodeBGUserAndRemoteRename(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &BGUserAndRemoteRename{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for BGUserAndRemoteRename-ManufacturerCode: %w", err)
	} else {
		if v != 381 {
			return nil, fmt.Errorf("match failed for BGUserAndRemoteRename-ManufacturerCode: Expected %d != %d", 381, v)
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
		return nil, fmt.Errorf("parse failed for BGUserAndRemoteRename-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for BGUserAndRemoteRename-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(12); err != nil {
		return nil, fmt.Errorf("parse failed for BGUserAndRemoteRename-DataType: %w", err)
	} else {
		val.DataType = BandgKeyValueConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for BGUserAndRemoteRename-Length: %w", err)
	} else {
		val.Length = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for BGUserAndRemoteRename-Decimals: %w", err)
	} else {
		val.Decimals = BandgDecimalsConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(64); err != nil {
		return nil, fmt.Errorf("parse failed for BGUserAndRemoteRename-ShortName: %w", err)
	} else {
		val.ShortName = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(128); err != nil {
		return nil, fmt.Errorf("parse failed for BGUserAndRemoteRename-LongName: %w", err)
	} else {
		val.LongName = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeBGUserAndRemoteRename(val *BGUserAndRemoteRename) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.DataType), 12)
	w.writeUInt8(val.Length, 4)
	w.skipBits(8)
	w.writeLookupField(uint64(val.Decimals), 8)
	w.writeFixedString(val.ShortName, 64)
	w.writeFixedString(val.LongName, 128)
	return w.Bytes(), w.Err()
}

func encodeBGUserAndRemoteRenameMsg(v Message) ([]byte, error) {
	val, ok := v.(*BGUserAndRemoteRename)
	if !ok {
		return nil, fmt.Errorf("expected *BGUserAndRemoteRename, got %T", v)
	}
	return EncodeBGUserAndRemoteRename(val)
}
