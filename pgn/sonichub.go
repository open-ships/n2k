package pgn

import "fmt"

type SonichubAlbum struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    SonichubCommandConst  `json:"proprietaryId"`
	Control          SonichubControlConst  `json:"control"`
	Item             *uint32               `json:"item"`
	Text             string                `json:"text"`
}

func (x *SonichubAlbum) PGNNumber() uint32 { return 130816 }

func DecodeSonichubAlbum(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SonichubAlbum{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubAlbum-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubAlbum-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubAlbum-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubAlbum-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubAlbum-ProprietaryId: %w", err)
	} else {
		if v != 16 {
			return nil, fmt.Errorf("match failed for SonichubAlbum-ProprietaryId: Expected %d != %d", 16, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubAlbum-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubAlbum-Item: %w", err)
	} else {
		val.Item = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubAlbum-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSonichubAlbum(val *SonichubAlbum) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt32(val.Item, 32)
	w.writeStringWithLength(val.Text)
	return w.Bytes(), w.Err()
}

func encodeSonichubAlbumMsg(v Message) ([]byte, error) {
	val, ok := v.(*SonichubAlbum)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubAlbum, got %T", v)
	}
	return EncodeSonichubAlbum(val)
}

type SonichubArtist struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    SonichubCommandConst  `json:"proprietaryId"`
	Control          SonichubControlConst  `json:"control"`
	Item             *uint32               `json:"item"`
	Text             string                `json:"text"`
}

func (x *SonichubArtist) PGNNumber() uint32 { return 130816 }

func DecodeSonichubArtist(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SonichubArtist{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubArtist-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubArtist-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubArtist-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubArtist-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubArtist-ProprietaryId: %w", err)
	} else {
		if v != 15 {
			return nil, fmt.Errorf("match failed for SonichubArtist-ProprietaryId: Expected %d != %d", 15, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubArtist-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubArtist-Item: %w", err)
	} else {
		val.Item = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubArtist-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSonichubArtist(val *SonichubArtist) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt32(val.Item, 32)
	w.writeStringWithLength(val.Text)
	return w.Bytes(), w.Err()
}

func encodeSonichubArtistMsg(v Message) ([]byte, error) {
	val, ok := v.(*SonichubArtist)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubArtist, got %T", v)
	}
	return EncodeSonichubArtist(val)
}

type SonichubControl struct {
	Info             MessageInfo            `json:"info"`
	ManufacturerCode ManufacturerCodeConst  `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst      `json:"industryCode"`
	ProprietaryId    SonichubCommandConst   `json:"proprietaryId"`
	Control          SonichubControlConst   `json:"control"`
	Item             FusionMuteCommandConst `json:"item"`
}

func (x *SonichubControl) PGNNumber() uint32 { return 130816 }

func DecodeSonichubControl(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SonichubControl{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubControl-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubControl-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubControl-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubControl-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubControl-ProprietaryId: %w", err)
	} else {
		if v != 9 {
			return nil, fmt.Errorf("match failed for SonichubControl-ProprietaryId: Expected %d != %d", 9, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubControl-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubControl-Item: %w", err)
	} else {
		val.Item = FusionMuteCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSonichubControl(val *SonichubControl) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeLookupField(uint64(val.Item), 8)
	return w.Bytes(), w.Err()
}

func encodeSonichubControlMsg(v Message) ([]byte, error) {
	val, ok := v.(*SonichubControl)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubControl, got %T", v)
	}
	return EncodeSonichubControl(val)
}

type SonichubFmRadio struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    SonichubCommandConst  `json:"proprietaryId"`
	Control          SonichubControlConst  `json:"control"`
	Item             SonichubTuningConst   `json:"item"`
	Frequency        *uint32               `json:"frequency"`
	NoiseLevel       *uint8                `json:"noiseLevel"`
	SignalLevel      *uint8                `json:"signalLevel"`
	Text             string                `json:"text"`
}

func (x *SonichubFmRadio) PGNNumber() uint32 { return 130816 }

func DecodeSonichubFmRadio(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SonichubFmRadio{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubFmRadio-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubFmRadio-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-ProprietaryId: %w", err)
	} else {
		if v != 12 {
			return nil, fmt.Errorf("match failed for SonichubFmRadio-ProprietaryId: Expected %d != %d", 12, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-Item: %w", err)
	} else {
		val.Item = SonichubTuningConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-Frequency: %w", err)
	} else {
		val.Frequency = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(2); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-NoiseLevel: %w", err)
	} else {
		val.NoiseLevel = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-SignalLevel: %w", err)
	} else {
		val.SignalLevel = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSonichubFmRadio(val *SonichubFmRadio) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeLookupField(uint64(val.Item), 8)
	w.writeUInt32(val.Frequency, 32)
	w.writeUInt8(val.NoiseLevel, 2)
	w.writeUInt8(val.SignalLevel, 4)
	w.skipBits(2)
	w.writeStringWithLength(val.Text)
	return w.Bytes(), w.Err()
}

func encodeSonichubFmRadioMsg(v Message) ([]byte, error) {
	val, ok := v.(*SonichubFmRadio)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubFmRadio, got %T", v)
	}
	return EncodeSonichubFmRadio(val)
}

type SonichubInit1 struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    SonichubCommandConst  `json:"proprietaryId"`
	Control          SonichubControlConst  `json:"control"`
}

func (x *SonichubInit1) PGNNumber() uint32 { return 130816 }

func DecodeSonichubInit1(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SonichubInit1{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubInit1-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubInit1-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubInit1-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubInit1-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubInit1-ProprietaryId: %w", err)
	} else {
		if v != 25 {
			return nil, fmt.Errorf("match failed for SonichubInit1-ProprietaryId: Expected %d != %d", 25, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubInit1-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSonichubInit1(val *SonichubInit1) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	return w.Bytes(), w.Err()
}

func encodeSonichubInit1Msg(v Message) ([]byte, error) {
	val, ok := v.(*SonichubInit1)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubInit1, got %T", v)
	}
	return EncodeSonichubInit1(val)
}

type SonichubMaxVolume struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    SonichubCommandConst  `json:"proprietaryId"`
	Control          SonichubControlConst  `json:"control"`
	Zone             *uint8                `json:"zone"`
	Level            *uint8                `json:"level"`
}

func (x *SonichubMaxVolume) PGNNumber() uint32 { return 130816 }

func DecodeSonichubMaxVolume(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SonichubMaxVolume{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMaxVolume-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubMaxVolume-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubMaxVolume-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubMaxVolume-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubMaxVolume-ProprietaryId: %w", err)
	} else {
		if v != 23 {
			return nil, fmt.Errorf("match failed for SonichubMaxVolume-ProprietaryId: Expected %d != %d", 23, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMaxVolume-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMaxVolume-Zone: %w", err)
	} else {
		val.Zone = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMaxVolume-Level: %w", err)
	} else {
		val.Level = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSonichubMaxVolume(val *SonichubMaxVolume) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt8(val.Zone, 8)
	w.writeUInt8(val.Level, 8)
	return w.Bytes(), w.Err()
}

func encodeSonichubMaxVolumeMsg(v Message) ([]byte, error) {
	val, ok := v.(*SonichubMaxVolume)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubMaxVolume, got %T", v)
	}
	return EncodeSonichubMaxVolume(val)
}

type SonichubMenuItem struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    SonichubCommandConst  `json:"proprietaryId"`
	Control          SonichubControlConst  `json:"control"`
	Item             *uint32               `json:"item"`
	C                *uint8                `json:"c"`
	D                *uint8                `json:"d"`
	E                *uint8                `json:"e"`
	Text             string                `json:"text"`
}

func (x *SonichubMenuItem) PGNNumber() uint32 { return 130816 }

func DecodeSonichubMenuItem(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SonichubMenuItem{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubMenuItem-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubMenuItem-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-ProprietaryId: %w", err)
	} else {
		if v != 19 {
			return nil, fmt.Errorf("match failed for SonichubMenuItem-ProprietaryId: Expected %d != %d", 19, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-Item: %w", err)
	} else {
		val.Item = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSonichubMenuItem(val *SonichubMenuItem) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt32(val.Item, 32)
	w.writeUInt8(val.C, 8)
	w.writeUInt8(val.D, 8)
	w.writeUInt8(val.E, 8)
	w.writeStringWithLength(val.Text)
	return w.Bytes(), w.Err()
}

func encodeSonichubMenuItemMsg(v Message) ([]byte, error) {
	val, ok := v.(*SonichubMenuItem)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubMenuItem, got %T", v)
	}
	return EncodeSonichubMenuItem(val)
}

type SonichubPlaylist struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    SonichubCommandConst  `json:"proprietaryId"`
	Control          SonichubControlConst  `json:"control"`
	Item             SonichubPlaylistConst `json:"item"`
	A                *uint8                `json:"a"`
	CurrentTrack     *uint32               `json:"currentTrack"`
	Tracks           *uint32               `json:"tracks"`
	Length           *float32              `json:"length"`
	PositionInTrack  *float32              `json:"positionInTrack"`
}

func (x *SonichubPlaylist) PGNNumber() uint32 { return 130816 }

func DecodeSonichubPlaylist(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SonichubPlaylist{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubPlaylist-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubPlaylist-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-ProprietaryId: %w", err)
	} else {
		if v != 13 {
			return nil, fmt.Errorf("match failed for SonichubPlaylist-ProprietaryId: Expected %d != %d", 13, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-Item: %w", err)
	} else {
		val.Item = SonichubPlaylistConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-CurrentTrack: %w", err)
	} else {
		val.CurrentTrack = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-Tracks: %w", err)
	} else {
		val.Tracks = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(32, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-Length: %w", err)
	} else {
		val.Length = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(32, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-PositionInTrack: %w", err)
	} else {
		val.PositionInTrack = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSonichubPlaylist(val *SonichubPlaylist) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeLookupField(uint64(val.Item), 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt32(val.CurrentTrack, 32)
	w.writeUInt32(val.Tracks, 32)
	w.writeUnsignedResolution(val.Length, 32, 0.001)
	w.writeUnsignedResolution(val.PositionInTrack, 32, 0.001)
	return w.Bytes(), w.Err()
}

func encodeSonichubPlaylistMsg(v Message) ([]byte, error) {
	val, ok := v.(*SonichubPlaylist)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubPlaylist, got %T", v)
	}
	return EncodeSonichubPlaylist(val)
}

type SonichubPosition struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    SonichubCommandConst  `json:"proprietaryId"`
	Control          SonichubControlConst  `json:"control"`
	Position         *float32              `json:"position"`
}

func (x *SonichubPosition) PGNNumber() uint32 { return 130816 }

func DecodeSonichubPosition(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SonichubPosition{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPosition-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubPosition-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubPosition-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubPosition-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubPosition-ProprietaryId: %w", err)
	} else {
		if v != 48 {
			return nil, fmt.Errorf("match failed for SonichubPosition-ProprietaryId: Expected %d != %d", 48, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPosition-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(32, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPosition-Position: %w", err)
	} else {
		val.Position = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSonichubPosition(val *SonichubPosition) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUnsignedResolution(val.Position, 32, 0.001)
	return w.Bytes(), w.Err()
}

func encodeSonichubPositionMsg(v Message) ([]byte, error) {
	val, ok := v.(*SonichubPosition)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubPosition, got %T", v)
	}
	return EncodeSonichubPosition(val)
}

type SonichubSource struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    SonichubCommandConst  `json:"proprietaryId"`
	Control          SonichubControlConst  `json:"control"`
	Source           SonichubSourceConst   `json:"source"`
}

func (x *SonichubSource) PGNNumber() uint32 { return 130816 }

func DecodeSonichubSource(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SonichubSource{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubSource-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubSource-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubSource-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubSource-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubSource-ProprietaryId: %w", err)
	} else {
		if v != 6 {
			return nil, fmt.Errorf("match failed for SonichubSource-ProprietaryId: Expected %d != %d", 6, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubSource-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubSource-Source: %w", err)
	} else {
		val.Source = SonichubSourceConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSonichubSource(val *SonichubSource) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeLookupField(uint64(val.Source), 8)
	return w.Bytes(), w.Err()
}

func encodeSonichubSourceMsg(v Message) ([]byte, error) {
	val, ok := v.(*SonichubSource)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubSource, got %T", v)
	}
	return EncodeSonichubSource(val)
}

type SonichubSourceList struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    SonichubCommandConst  `json:"proprietaryId"`
	Control          SonichubControlConst  `json:"control"`
	SourceId         *uint8                `json:"sourceId"`
	A                *uint8                `json:"a"`
	Text             string                `json:"text"`
}

func (x *SonichubSourceList) PGNNumber() uint32 { return 130816 }

func DecodeSonichubSourceList(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SonichubSourceList{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubSourceList-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubSourceList-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubSourceList-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubSourceList-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubSourceList-ProprietaryId: %w", err)
	} else {
		if v != 8 {
			return nil, fmt.Errorf("match failed for SonichubSourceList-ProprietaryId: Expected %d != %d", 8, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubSourceList-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubSourceList-SourceId: %w", err)
	} else {
		val.SourceId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubSourceList-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubSourceList-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSonichubSourceList(val *SonichubSourceList) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt8(val.SourceId, 8)
	w.writeUInt8(val.A, 8)
	w.writeStringWithLength(val.Text)
	return w.Bytes(), w.Err()
}

func encodeSonichubSourceListMsg(v Message) ([]byte, error) {
	val, ok := v.(*SonichubSourceList)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubSourceList, got %T", v)
	}
	return EncodeSonichubSourceList(val)
}

type SonichubTrack struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    SonichubCommandConst  `json:"proprietaryId"`
	Control          SonichubControlConst  `json:"control"`
	Item             *uint32               `json:"item"`
	Text             string                `json:"text"`
}

func (x *SonichubTrack) PGNNumber() uint32 { return 130816 }

func DecodeSonichubTrack(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SonichubTrack{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubTrack-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubTrack-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubTrack-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubTrack-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubTrack-ProprietaryId: %w", err)
	} else {
		if v != 14 {
			return nil, fmt.Errorf("match failed for SonichubTrack-ProprietaryId: Expected %d != %d", 14, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubTrack-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubTrack-Item: %w", err)
	} else {
		val.Item = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubTrack-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSonichubTrack(val *SonichubTrack) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt32(val.Item, 32)
	w.writeStringWithLength(val.Text)
	return w.Bytes(), w.Err()
}

func encodeSonichubTrackMsg(v Message) ([]byte, error) {
	val, ok := v.(*SonichubTrack)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubTrack, got %T", v)
	}
	return EncodeSonichubTrack(val)
}

type SonichubVolume struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    SonichubCommandConst  `json:"proprietaryId"`
	Control          SonichubControlConst  `json:"control"`
	Zone             *uint8                `json:"zone"`
	Level            *uint8                `json:"level"`
}

func (x *SonichubVolume) PGNNumber() uint32 { return 130816 }

func DecodeSonichubVolume(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SonichubVolume{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubVolume-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubVolume-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubVolume-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubVolume-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubVolume-ProprietaryId: %w", err)
	} else {
		if v != 24 {
			return nil, fmt.Errorf("match failed for SonichubVolume-ProprietaryId: Expected %d != %d", 24, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubVolume-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubVolume-Zone: %w", err)
	} else {
		val.Zone = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubVolume-Level: %w", err)
	} else {
		val.Level = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSonichubVolume(val *SonichubVolume) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt8(val.Zone, 8)
	w.writeUInt8(val.Level, 8)
	return w.Bytes(), w.Err()
}

func encodeSonichubVolumeMsg(v Message) ([]byte, error) {
	val, ok := v.(*SonichubVolume)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubVolume, got %T", v)
	}
	return EncodeSonichubVolume(val)
}

type SonichubZoneInfo struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    SonichubCommandConst  `json:"proprietaryId"`
	Control          SonichubControlConst  `json:"control"`
	Zone             *uint8                `json:"zone"`
}

func (x *SonichubZoneInfo) PGNNumber() uint32 { return 130816 }

func DecodeSonichubZoneInfo(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SonichubZoneInfo{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubZoneInfo-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubZoneInfo-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubZoneInfo-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubZoneInfo-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubZoneInfo-ProprietaryId: %w", err)
	} else {
		if v != 5 {
			return nil, fmt.Errorf("match failed for SonichubZoneInfo-ProprietaryId: Expected %d != %d", 5, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubZoneInfo-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubZoneInfo-Zone: %w", err)
	} else {
		val.Zone = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSonichubZoneInfo(val *SonichubZoneInfo) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt8(val.Zone, 8)
	return w.Bytes(), w.Err()
}

func encodeSonichubZoneInfoMsg(v Message) ([]byte, error) {
	val, ok := v.(*SonichubZoneInfo)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubZoneInfo, got %T", v)
	}
	return EncodeSonichubZoneInfo(val)
}

type SonichubZones struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    SonichubCommandConst  `json:"proprietaryId"`
	Control          SonichubControlConst  `json:"control"`
	Zones            *uint8                `json:"zones"`
}

func (x *SonichubZones) PGNNumber() uint32 { return 130816 }

func DecodeSonichubZones(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SonichubZones{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubZones-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubZones-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubZones-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubZones-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubZones-ProprietaryId: %w", err)
	} else {
		if v != 20 {
			return nil, fmt.Errorf("match failed for SonichubZones-ProprietaryId: Expected %d != %d", 20, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubZones-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubZones-Zones: %w", err)
	} else {
		val.Zones = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSonichubZones(val *SonichubZones) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt8(val.Zones, 8)
	return w.Bytes(), w.Err()
}

func encodeSonichubZonesMsg(v Message) ([]byte, error) {
	val, ok := v.(*SonichubZones)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubZones, got %T", v)
	}
	return EncodeSonichubZones(val)
}
