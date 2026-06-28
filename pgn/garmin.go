package pgn

import "fmt"

type GarminColorMode struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	UnknownId1       *uint8                `json:"unknownId1"`
	UnknownId2       *uint8                `json:"unknownId2"`
	UnknownId3       *uint8                `json:"unknownId3"`
	UnknownId4       *uint8                `json:"unknownId4"`
	Mode             GarminColorModeConst  `json:"mode"`
	Color            GarminColorConst      `json:"color"`
}

func (g *GarminColorMode) PGNNumber() uint32 { return 126720 }

func DecodeGarminColorMode(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GarminColorMode{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for GarminColorMode-ManufacturerCode: %w", err)
	} else {
		if v != 229 {
			return nil, fmt.Errorf("match failed for GarminColorMode-ManufacturerCode: Expected %d != %d", 229, v)
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
		return nil, fmt.Errorf("parse failed for GarminColorMode-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for GarminColorMode-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminColorMode-UnknownId1: %w", err)
	} else {
		if v != nil && *v != 222 {
			return nil, fmt.Errorf("match failed for GarminColorMode-UnknownId1: Expected %d != %d", 222, *v)
		}
		val.UnknownId1 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminColorMode-UnknownId2: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminColorMode-UnknownId2: Expected %d != %d", 5, *v)
		}
		val.UnknownId2 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminColorMode-UnknownId3: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminColorMode-UnknownId3: Expected %d != %d", 5, *v)
		}
		val.UnknownId3 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminColorMode-UnknownId4: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminColorMode-UnknownId4: Expected %d != %d", 5, *v)
		}
		val.UnknownId4 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(16)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminColorMode-Mode: %w", err)
	} else {
		if v != 13 {
			return nil, fmt.Errorf("match failed for GarminColorMode-Mode: Expected %d != %d", 13, v)
		}
		val.Mode = GarminColorModeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminColorMode-Color: %w", err)
	} else {
		val.Color = GarminColorConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeGarminColorMode(val *GarminColorMode) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.UnknownId1, 8)
	w.writeUInt8(val.UnknownId2, 8)
	w.writeUInt8(val.UnknownId3, 8)
	w.writeUInt8(val.UnknownId4, 8)
	w.writeSpareBits(16)
	w.writeLookupField(uint64(val.Mode), 8)
	w.writeSpareBits(8)
	w.writeLookupField(uint64(val.Color), 8)
	return w.Bytes(), w.Err()
}

func encodeGarminColorModeMsg(v Message) ([]byte, error) {
	val, ok := v.(*GarminColorMode)
	if !ok {
		return nil, fmt.Errorf("expected *GarminColorMode, got %T", v)
	}
	return EncodeGarminColorMode(val)
}

type GarminDayMode struct {
	Info             MessageInfo               `json:"info"`
	ManufacturerCode ManufacturerCodeConst     `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst         `json:"industryCode"`
	UnknownId1       *uint8                    `json:"unknownId1"`
	UnknownId2       *uint8                    `json:"unknownId2"`
	UnknownId3       *uint8                    `json:"unknownId3"`
	UnknownId4       *uint8                    `json:"unknownId4"`
	Mode             GarminColorModeConst      `json:"mode"`
	Backlight        GarminBacklightLevelConst `json:"backlight"`
}

func (g *GarminDayMode) PGNNumber() uint32 { return 126720 }

func DecodeGarminDayMode(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GarminDayMode{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for GarminDayMode-ManufacturerCode: %w", err)
	} else {
		if v != 229 {
			return nil, fmt.Errorf("match failed for GarminDayMode-ManufacturerCode: Expected %d != %d", 229, v)
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
		return nil, fmt.Errorf("parse failed for GarminDayMode-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for GarminDayMode-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminDayMode-UnknownId1: %w", err)
	} else {
		if v != nil && *v != 222 {
			return nil, fmt.Errorf("match failed for GarminDayMode-UnknownId1: Expected %d != %d", 222, *v)
		}
		val.UnknownId1 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminDayMode-UnknownId2: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminDayMode-UnknownId2: Expected %d != %d", 5, *v)
		}
		val.UnknownId2 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminDayMode-UnknownId3: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminDayMode-UnknownId3: Expected %d != %d", 5, *v)
		}
		val.UnknownId3 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminDayMode-UnknownId4: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminDayMode-UnknownId4: Expected %d != %d", 5, *v)
		}
		val.UnknownId4 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(16)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminDayMode-Mode: %w", err)
	} else {
		if v != 0 {
			return nil, fmt.Errorf("match failed for GarminDayMode-Mode: Expected %d != %d", 0, v)
		}
		val.Mode = GarminColorModeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminDayMode-Backlight: %w", err)
	} else {
		val.Backlight = GarminBacklightLevelConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeGarminDayMode(val *GarminDayMode) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.UnknownId1, 8)
	w.writeUInt8(val.UnknownId2, 8)
	w.writeUInt8(val.UnknownId3, 8)
	w.writeUInt8(val.UnknownId4, 8)
	w.writeSpareBits(16)
	w.writeLookupField(uint64(val.Mode), 8)
	w.writeSpareBits(8)
	w.writeLookupField(uint64(val.Backlight), 8)
	return w.Bytes(), w.Err()
}

func encodeGarminDayModeMsg(v Message) ([]byte, error) {
	val, ok := v.(*GarminDayMode)
	if !ok {
		return nil, fmt.Errorf("expected *GarminDayMode, got %T", v)
	}
	return EncodeGarminDayMode(val)
}

type GarminNightMode struct {
	Info             MessageInfo               `json:"info"`
	ManufacturerCode ManufacturerCodeConst     `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst         `json:"industryCode"`
	UnknownId1       *uint8                    `json:"unknownId1"`
	UnknownId2       *uint8                    `json:"unknownId2"`
	UnknownId3       *uint8                    `json:"unknownId3"`
	UnknownId4       *uint8                    `json:"unknownId4"`
	Mode             GarminColorModeConst      `json:"mode"`
	Backlight        GarminBacklightLevelConst `json:"backlight"`
}

func (g *GarminNightMode) PGNNumber() uint32 { return 126720 }

func DecodeGarminNightMode(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GarminNightMode{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for GarminNightMode-ManufacturerCode: %w", err)
	} else {
		if v != 229 {
			return nil, fmt.Errorf("match failed for GarminNightMode-ManufacturerCode: Expected %d != %d", 229, v)
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
		return nil, fmt.Errorf("parse failed for GarminNightMode-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for GarminNightMode-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminNightMode-UnknownId1: %w", err)
	} else {
		if v != nil && *v != 222 {
			return nil, fmt.Errorf("match failed for GarminNightMode-UnknownId1: Expected %d != %d", 222, *v)
		}
		val.UnknownId1 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminNightMode-UnknownId2: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminNightMode-UnknownId2: Expected %d != %d", 5, *v)
		}
		val.UnknownId2 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminNightMode-UnknownId3: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminNightMode-UnknownId3: Expected %d != %d", 5, *v)
		}
		val.UnknownId3 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminNightMode-UnknownId4: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminNightMode-UnknownId4: Expected %d != %d", 5, *v)
		}
		val.UnknownId4 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(16)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminNightMode-Mode: %w", err)
	} else {
		if v != 1 {
			return nil, fmt.Errorf("match failed for GarminNightMode-Mode: Expected %d != %d", 1, v)
		}
		val.Mode = GarminColorModeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminNightMode-Backlight: %w", err)
	} else {
		val.Backlight = GarminBacklightLevelConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeGarminNightMode(val *GarminNightMode) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.UnknownId1, 8)
	w.writeUInt8(val.UnknownId2, 8)
	w.writeUInt8(val.UnknownId3, 8)
	w.writeUInt8(val.UnknownId4, 8)
	w.writeSpareBits(16)
	w.writeLookupField(uint64(val.Mode), 8)
	w.writeSpareBits(8)
	w.writeLookupField(uint64(val.Backlight), 8)
	return w.Bytes(), w.Err()
}

func encodeGarminNightModeMsg(v Message) ([]byte, error) {
	val, ok := v.(*GarminNightMode)
	if !ok {
		return nil, fmt.Errorf("expected *GarminNightMode, got %T", v)
	}
	return EncodeGarminNightMode(val)
}
