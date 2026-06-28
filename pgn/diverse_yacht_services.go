package pgn

import "fmt"

type DiverseYachtServicesLoadCell struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	Instance         *uint8                `json:"instance"`
	LoadCell         *uint32               `json:"loadCell"`
}

func (d *DiverseYachtServicesLoadCell) PGNNumber() uint32 { return 65293 }

func DecodeDiverseYachtServicesLoadCell(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &DiverseYachtServicesLoadCell{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for DiverseYachtServicesLoadCell-ManufacturerCode: %w", err)
	} else {
		if v != 641 {
			return nil, fmt.Errorf("match failed for DiverseYachtServicesLoadCell-ManufacturerCode: Expected %d != %d", 641, v)
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
		return nil, fmt.Errorf("parse failed for DiverseYachtServicesLoadCell-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for DiverseYachtServicesLoadCell-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for DiverseYachtServicesLoadCell-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for DiverseYachtServicesLoadCell-LoadCell: %w", err)
	} else {
		val.LoadCell = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeDiverseYachtServicesLoadCell(val *DiverseYachtServicesLoadCell) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Instance, 8)
	w.writeReservedBits(8)
	w.writeUInt32(val.LoadCell, 32)
	return w.Bytes(), w.Err()
}

func encodeDiverseYachtServicesLoadCellMsg(v Message) ([]byte, error) {
	val, ok := v.(*DiverseYachtServicesLoadCell)
	if !ok {
		return nil, fmt.Errorf("expected *DiverseYachtServicesLoadCell, got %T", v)
	}
	return EncodeDiverseYachtServicesLoadCell(val)
}
