package pgn

import (
	"fmt"

	"github.com/open-ships/n2k/units"
)

type Bus1PhaseCBasicAcQuantities struct {
	Info                    MessageInfo `json:"info"`
	LineLineAcRmsVoltage    *uint16     `json:"lineLineAcRmsVoltage"`
	LineNeutralAcRmsVoltage *uint16     `json:"lineNeutralAcRmsVoltage"`
	AcFrequency             *float32    `json:"acFrequency"`
}

func (b *Bus1PhaseCBasicAcQuantities) PGNNumber() uint32 { return 65001 }

func DecodeBus1PhaseCBasicAcQuantities(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Bus1PhaseCBasicAcQuantities{}
	val.Info = Info
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for Bus1PhaseCBasicAcQuantities-LineLineAcRmsVoltage: %w", err)
	} else {
		val.LineLineAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for Bus1PhaseCBasicAcQuantities-LineNeutralAcRmsVoltage: %w", err)
	} else {
		val.LineNeutralAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0078125); err != nil {
		return nil, fmt.Errorf("parse failed for Bus1PhaseCBasicAcQuantities-AcFrequency: %w", err)
	} else {
		val.AcFrequency = v

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

func EncodeBus1PhaseCBasicAcQuantities(val *Bus1PhaseCBasicAcQuantities) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.LineLineAcRmsVoltage, 16)
	w.writeUInt16(val.LineNeutralAcRmsVoltage, 16)
	w.writeUnsignedResolution(val.AcFrequency, 16, 0.0078125)
	w.writeReservedBits(16)
	return w.Bytes(), w.Err()
}

func encodeBus1PhaseCBasicAcQuantitiesMsg(v Message) ([]byte, error) {
	val, ok := v.(*Bus1PhaseCBasicAcQuantities)
	if !ok {
		return nil, fmt.Errorf("expected *Bus1PhaseCBasicAcQuantities, got %T", v)
	}
	return EncodeBus1PhaseCBasicAcQuantities(val)
}

type Bus1PhaseBBasicAcQuantities struct {
	Info                    MessageInfo `json:"info"`
	LineLineAcRmsVoltage    *uint16     `json:"lineLineAcRmsVoltage"`
	LineNeutralAcRmsVoltage *uint16     `json:"lineNeutralAcRmsVoltage"`
	AcFrequency             *float32    `json:"acFrequency"`
}

func (b *Bus1PhaseBBasicAcQuantities) PGNNumber() uint32 { return 65002 }

func DecodeBus1PhaseBBasicAcQuantities(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Bus1PhaseBBasicAcQuantities{}
	val.Info = Info
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for Bus1PhaseBBasicAcQuantities-LineLineAcRmsVoltage: %w", err)
	} else {
		val.LineLineAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for Bus1PhaseBBasicAcQuantities-LineNeutralAcRmsVoltage: %w", err)
	} else {
		val.LineNeutralAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0078125); err != nil {
		return nil, fmt.Errorf("parse failed for Bus1PhaseBBasicAcQuantities-AcFrequency: %w", err)
	} else {
		val.AcFrequency = v

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

func EncodeBus1PhaseBBasicAcQuantities(val *Bus1PhaseBBasicAcQuantities) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.LineLineAcRmsVoltage, 16)
	w.writeUInt16(val.LineNeutralAcRmsVoltage, 16)
	w.writeUnsignedResolution(val.AcFrequency, 16, 0.0078125)
	w.writeReservedBits(16)
	return w.Bytes(), w.Err()
}

func encodeBus1PhaseBBasicAcQuantitiesMsg(v Message) ([]byte, error) {
	val, ok := v.(*Bus1PhaseBBasicAcQuantities)
	if !ok {
		return nil, fmt.Errorf("expected *Bus1PhaseBBasicAcQuantities, got %T", v)
	}
	return EncodeBus1PhaseBBasicAcQuantities(val)
}

type Bus1PhaseABasicAcQuantities struct {
	Info                    MessageInfo `json:"info"`
	LineLineAcRmsVoltage    *uint16     `json:"lineLineAcRmsVoltage"`
	LineNeutralAcRmsVoltage *uint16     `json:"lineNeutralAcRmsVoltage"`
	AcFrequency             *float32    `json:"acFrequency"`
}

func (b *Bus1PhaseABasicAcQuantities) PGNNumber() uint32 { return 65003 }

func DecodeBus1PhaseABasicAcQuantities(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Bus1PhaseABasicAcQuantities{}
	val.Info = Info
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for Bus1PhaseABasicAcQuantities-LineLineAcRmsVoltage: %w", err)
	} else {
		val.LineLineAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for Bus1PhaseABasicAcQuantities-LineNeutralAcRmsVoltage: %w", err)
	} else {
		val.LineNeutralAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0078125); err != nil {
		return nil, fmt.Errorf("parse failed for Bus1PhaseABasicAcQuantities-AcFrequency: %w", err)
	} else {
		val.AcFrequency = v

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

func EncodeBus1PhaseABasicAcQuantities(val *Bus1PhaseABasicAcQuantities) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.LineLineAcRmsVoltage, 16)
	w.writeUInt16(val.LineNeutralAcRmsVoltage, 16)
	w.writeUnsignedResolution(val.AcFrequency, 16, 0.0078125)
	w.writeReservedBits(16)
	return w.Bytes(), w.Err()
}

func encodeBus1PhaseABasicAcQuantitiesMsg(v Message) ([]byte, error) {
	val, ok := v.(*Bus1PhaseABasicAcQuantities)
	if !ok {
		return nil, fmt.Errorf("expected *Bus1PhaseABasicAcQuantities, got %T", v)
	}
	return EncodeBus1PhaseABasicAcQuantities(val)
}

type Bus1AverageBasicAcQuantities struct {
	Info                    MessageInfo `json:"info"`
	LineLineAcRmsVoltage    *uint16     `json:"lineLineAcRmsVoltage"`
	LineNeutralAcRmsVoltage *uint16     `json:"lineNeutralAcRmsVoltage"`
	AcFrequency             *float32    `json:"acFrequency"`
}

func (b *Bus1AverageBasicAcQuantities) PGNNumber() uint32 { return 65004 }

func DecodeBus1AverageBasicAcQuantities(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Bus1AverageBasicAcQuantities{}
	val.Info = Info
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for Bus1AverageBasicAcQuantities-LineLineAcRmsVoltage: %w", err)
	} else {
		val.LineLineAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for Bus1AverageBasicAcQuantities-LineNeutralAcRmsVoltage: %w", err)
	} else {
		val.LineNeutralAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0078125); err != nil {
		return nil, fmt.Errorf("parse failed for Bus1AverageBasicAcQuantities-AcFrequency: %w", err)
	} else {
		val.AcFrequency = v

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

func EncodeBus1AverageBasicAcQuantities(val *Bus1AverageBasicAcQuantities) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.LineLineAcRmsVoltage, 16)
	w.writeUInt16(val.LineNeutralAcRmsVoltage, 16)
	w.writeUnsignedResolution(val.AcFrequency, 16, 0.0078125)
	w.writeReservedBits(16)
	return w.Bytes(), w.Err()
}

func encodeBus1AverageBasicAcQuantitiesMsg(v Message) ([]byte, error) {
	val, ok := v.(*Bus1AverageBasicAcQuantities)
	if !ok {
		return nil, fmt.Errorf("expected *Bus1AverageBasicAcQuantities, got %T", v)
	}
	return EncodeBus1AverageBasicAcQuantities(val)
}

type UtilityTotalAcEnergy struct {
	Info              MessageInfo `json:"info"`
	TotalEnergyExport *uint32     `json:"totalEnergyExport"`
	TotalEnergyImport *uint32     `json:"totalEnergyImport"`
}

func (u *UtilityTotalAcEnergy) PGNNumber() uint32 { return 65005 }

func DecodeUtilityTotalAcEnergy(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &UtilityTotalAcEnergy{}
	val.Info = Info
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityTotalAcEnergy-TotalEnergyExport: %w", err)
	} else {
		val.TotalEnergyExport = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityTotalAcEnergy-TotalEnergyImport: %w", err)
	} else {
		val.TotalEnergyImport = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeUtilityTotalAcEnergy(val *UtilityTotalAcEnergy) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt32(val.TotalEnergyExport, 32)
	w.writeUInt32(val.TotalEnergyImport, 32)
	return w.Bytes(), w.Err()
}

func encodeUtilityTotalAcEnergyMsg(v Message) ([]byte, error) {
	val, ok := v.(*UtilityTotalAcEnergy)
	if !ok {
		return nil, fmt.Errorf("expected *UtilityTotalAcEnergy, got %T", v)
	}
	return EncodeUtilityTotalAcEnergy(val)
}

type UtilityPhaseCAcReactivePower struct {
	Info               MessageInfo      `json:"info"`
	ReactivePower      *uint16          `json:"reactivePower"`
	PowerFactor        *float32         `json:"powerFactor"`
	PowerFactorLagging PowerFactorConst `json:"powerFactorLagging"`
}

func (u *UtilityPhaseCAcReactivePower) PGNNumber() uint32 { return 65006 }

func DecodeUtilityPhaseCAcReactivePower(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &UtilityPhaseCAcReactivePower{}
	val.Info = Info
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseCAcReactivePower-ReactivePower: %w", err)
	} else {
		val.ReactivePower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 6.10352e-05); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseCAcReactivePower-PowerFactor: %w", err)
	} else {
		val.PowerFactor = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseCAcReactivePower-PowerFactorLagging: %w", err)
	} else {
		val.PowerFactorLagging = PowerFactorConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(30)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeUtilityPhaseCAcReactivePower(val *UtilityPhaseCAcReactivePower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.ReactivePower, 16)
	w.writeUnsignedResolution(val.PowerFactor, 16, 6.10352e-05)
	w.writeLookupField(uint64(val.PowerFactorLagging), 2)
	w.writeReservedBits(30)
	return w.Bytes(), w.Err()
}

func encodeUtilityPhaseCAcReactivePowerMsg(v Message) ([]byte, error) {
	val, ok := v.(*UtilityPhaseCAcReactivePower)
	if !ok {
		return nil, fmt.Errorf("expected *UtilityPhaseCAcReactivePower, got %T", v)
	}
	return EncodeUtilityPhaseCAcReactivePower(val)
}

type UtilityPhaseCAcPower struct {
	Info          MessageInfo `json:"info"`
	RealPower     *int32      `json:"realPower"`
	ApparentPower *int32      `json:"apparentPower"`
}

func (u *UtilityPhaseCAcPower) PGNNumber() uint32 { return 65007 }

func DecodeUtilityPhaseCAcPower(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &UtilityPhaseCAcPower{}
	val.Info = Info
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseCAcPower-RealPower: %w", err)
	} else {
		val.RealPower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseCAcPower-ApparentPower: %w", err)
	} else {
		val.ApparentPower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeUtilityPhaseCAcPower(val *UtilityPhaseCAcPower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeInt32(val.RealPower, 32)
	w.writeInt32(val.ApparentPower, 32)
	return w.Bytes(), w.Err()
}

func encodeUtilityPhaseCAcPowerMsg(v Message) ([]byte, error) {
	val, ok := v.(*UtilityPhaseCAcPower)
	if !ok {
		return nil, fmt.Errorf("expected *UtilityPhaseCAcPower, got %T", v)
	}
	return EncodeUtilityPhaseCAcPower(val)
}

type UtilityPhaseCBasicAcQuantities struct {
	Info                    MessageInfo `json:"info"`
	LineLineAcRmsVoltage    *uint16     `json:"lineLineAcRmsVoltage"`
	LineNeutralAcRmsVoltage *uint16     `json:"lineNeutralAcRmsVoltage"`
	AcFrequency             *float32    `json:"acFrequency"`
	AcRmsCurrent            *uint16     `json:"acRmsCurrent"`
}

func (u *UtilityPhaseCBasicAcQuantities) PGNNumber() uint32 { return 65008 }

func DecodeUtilityPhaseCBasicAcQuantities(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &UtilityPhaseCBasicAcQuantities{}
	val.Info = Info
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseCBasicAcQuantities-LineLineAcRmsVoltage: %w", err)
	} else {
		val.LineLineAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseCBasicAcQuantities-LineNeutralAcRmsVoltage: %w", err)
	} else {
		val.LineNeutralAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0078125); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseCBasicAcQuantities-AcFrequency: %w", err)
	} else {
		val.AcFrequency = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseCBasicAcQuantities-AcRmsCurrent: %w", err)
	} else {
		val.AcRmsCurrent = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeUtilityPhaseCBasicAcQuantities(val *UtilityPhaseCBasicAcQuantities) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.LineLineAcRmsVoltage, 16)
	w.writeUInt16(val.LineNeutralAcRmsVoltage, 16)
	w.writeUnsignedResolution(val.AcFrequency, 16, 0.0078125)
	w.writeUInt16(val.AcRmsCurrent, 16)
	return w.Bytes(), w.Err()
}

func encodeUtilityPhaseCBasicAcQuantitiesMsg(v Message) ([]byte, error) {
	val, ok := v.(*UtilityPhaseCBasicAcQuantities)
	if !ok {
		return nil, fmt.Errorf("expected *UtilityPhaseCBasicAcQuantities, got %T", v)
	}
	return EncodeUtilityPhaseCBasicAcQuantities(val)
}

type UtilityPhaseBAcReactivePower struct {
	Info               MessageInfo      `json:"info"`
	ReactivePower      *uint16          `json:"reactivePower"`
	PowerFactor        *float32         `json:"powerFactor"`
	PowerFactorLagging PowerFactorConst `json:"powerFactorLagging"`
}

func (u *UtilityPhaseBAcReactivePower) PGNNumber() uint32 { return 65009 }

func DecodeUtilityPhaseBAcReactivePower(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &UtilityPhaseBAcReactivePower{}
	val.Info = Info
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseBAcReactivePower-ReactivePower: %w", err)
	} else {
		val.ReactivePower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 6.10352e-05); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseBAcReactivePower-PowerFactor: %w", err)
	} else {
		val.PowerFactor = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseBAcReactivePower-PowerFactorLagging: %w", err)
	} else {
		val.PowerFactorLagging = PowerFactorConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(30)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeUtilityPhaseBAcReactivePower(val *UtilityPhaseBAcReactivePower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.ReactivePower, 16)
	w.writeUnsignedResolution(val.PowerFactor, 16, 6.10352e-05)
	w.writeLookupField(uint64(val.PowerFactorLagging), 2)
	w.writeReservedBits(30)
	return w.Bytes(), w.Err()
}

func encodeUtilityPhaseBAcReactivePowerMsg(v Message) ([]byte, error) {
	val, ok := v.(*UtilityPhaseBAcReactivePower)
	if !ok {
		return nil, fmt.Errorf("expected *UtilityPhaseBAcReactivePower, got %T", v)
	}
	return EncodeUtilityPhaseBAcReactivePower(val)
}

type UtilityPhaseBAcPower struct {
	Info          MessageInfo `json:"info"`
	RealPower     *int32      `json:"realPower"`
	ApparentPower *int32      `json:"apparentPower"`
}

func (u *UtilityPhaseBAcPower) PGNNumber() uint32 { return 65010 }

func DecodeUtilityPhaseBAcPower(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &UtilityPhaseBAcPower{}
	val.Info = Info
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseBAcPower-RealPower: %w", err)
	} else {
		val.RealPower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseBAcPower-ApparentPower: %w", err)
	} else {
		val.ApparentPower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeUtilityPhaseBAcPower(val *UtilityPhaseBAcPower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeInt32(val.RealPower, 32)
	w.writeInt32(val.ApparentPower, 32)
	return w.Bytes(), w.Err()
}

func encodeUtilityPhaseBAcPowerMsg(v Message) ([]byte, error) {
	val, ok := v.(*UtilityPhaseBAcPower)
	if !ok {
		return nil, fmt.Errorf("expected *UtilityPhaseBAcPower, got %T", v)
	}
	return EncodeUtilityPhaseBAcPower(val)
}

type UtilityPhaseBBasicAcQuantities struct {
	Info                    MessageInfo `json:"info"`
	LineLineAcRmsVoltage    *uint16     `json:"lineLineAcRmsVoltage"`
	LineNeutralAcRmsVoltage *uint16     `json:"lineNeutralAcRmsVoltage"`
	AcFrequency             *float32    `json:"acFrequency"`
	AcRmsCurrent            *uint16     `json:"acRmsCurrent"`
}

func (u *UtilityPhaseBBasicAcQuantities) PGNNumber() uint32 { return 65011 }

func DecodeUtilityPhaseBBasicAcQuantities(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &UtilityPhaseBBasicAcQuantities{}
	val.Info = Info
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseBBasicAcQuantities-LineLineAcRmsVoltage: %w", err)
	} else {
		val.LineLineAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseBBasicAcQuantities-LineNeutralAcRmsVoltage: %w", err)
	} else {
		val.LineNeutralAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0078125); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseBBasicAcQuantities-AcFrequency: %w", err)
	} else {
		val.AcFrequency = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseBBasicAcQuantities-AcRmsCurrent: %w", err)
	} else {
		val.AcRmsCurrent = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeUtilityPhaseBBasicAcQuantities(val *UtilityPhaseBBasicAcQuantities) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.LineLineAcRmsVoltage, 16)
	w.writeUInt16(val.LineNeutralAcRmsVoltage, 16)
	w.writeUnsignedResolution(val.AcFrequency, 16, 0.0078125)
	w.writeUInt16(val.AcRmsCurrent, 16)
	return w.Bytes(), w.Err()
}

func encodeUtilityPhaseBBasicAcQuantitiesMsg(v Message) ([]byte, error) {
	val, ok := v.(*UtilityPhaseBBasicAcQuantities)
	if !ok {
		return nil, fmt.Errorf("expected *UtilityPhaseBBasicAcQuantities, got %T", v)
	}
	return EncodeUtilityPhaseBBasicAcQuantities(val)
}

type UtilityPhaseAAcReactivePower struct {
	Info               MessageInfo      `json:"info"`
	ReactivePower      *int32           `json:"reactivePower"`
	PowerFactor        *float32         `json:"powerFactor"`
	PowerFactorLagging PowerFactorConst `json:"powerFactorLagging"`
}

func (u *UtilityPhaseAAcReactivePower) PGNNumber() uint32 { return 65012 }

func DecodeUtilityPhaseAAcReactivePower(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &UtilityPhaseAAcReactivePower{}
	val.Info = Info
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseAAcReactivePower-ReactivePower: %w", err)
	} else {
		val.ReactivePower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 6.10352e-05); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseAAcReactivePower-PowerFactor: %w", err)
	} else {
		val.PowerFactor = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseAAcReactivePower-PowerFactorLagging: %w", err)
	} else {
		val.PowerFactorLagging = PowerFactorConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(14)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeUtilityPhaseAAcReactivePower(val *UtilityPhaseAAcReactivePower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeInt32(val.ReactivePower, 32)
	w.writeUnsignedResolution(val.PowerFactor, 16, 6.10352e-05)
	w.writeLookupField(uint64(val.PowerFactorLagging), 2)
	w.writeReservedBits(14)
	return w.Bytes(), w.Err()
}

func encodeUtilityPhaseAAcReactivePowerMsg(v Message) ([]byte, error) {
	val, ok := v.(*UtilityPhaseAAcReactivePower)
	if !ok {
		return nil, fmt.Errorf("expected *UtilityPhaseAAcReactivePower, got %T", v)
	}
	return EncodeUtilityPhaseAAcReactivePower(val)
}

type UtilityPhaseAAcPower struct {
	Info          MessageInfo `json:"info"`
	RealPower     *int32      `json:"realPower"`
	ApparentPower *int32      `json:"apparentPower"`
}

func (u *UtilityPhaseAAcPower) PGNNumber() uint32 { return 65013 }

func DecodeUtilityPhaseAAcPower(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &UtilityPhaseAAcPower{}
	val.Info = Info
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseAAcPower-RealPower: %w", err)
	} else {
		val.RealPower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseAAcPower-ApparentPower: %w", err)
	} else {
		val.ApparentPower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeUtilityPhaseAAcPower(val *UtilityPhaseAAcPower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeInt32(val.RealPower, 32)
	w.writeInt32(val.ApparentPower, 32)
	return w.Bytes(), w.Err()
}

func encodeUtilityPhaseAAcPowerMsg(v Message) ([]byte, error) {
	val, ok := v.(*UtilityPhaseAAcPower)
	if !ok {
		return nil, fmt.Errorf("expected *UtilityPhaseAAcPower, got %T", v)
	}
	return EncodeUtilityPhaseAAcPower(val)
}

type UtilityPhaseABasicAcQuantities struct {
	Info                    MessageInfo `json:"info"`
	LineLineAcRmsVoltage    *uint16     `json:"lineLineAcRmsVoltage"`
	LineNeutralAcRmsVoltage *uint16     `json:"lineNeutralAcRmsVoltage"`
	AcFrequency             *float32    `json:"acFrequency"`
	AcRmsCurrent            *uint16     `json:"acRmsCurrent"`
}

func (u *UtilityPhaseABasicAcQuantities) PGNNumber() uint32 { return 65014 }

func DecodeUtilityPhaseABasicAcQuantities(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &UtilityPhaseABasicAcQuantities{}
	val.Info = Info
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseABasicAcQuantities-LineLineAcRmsVoltage: %w", err)
	} else {
		val.LineLineAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseABasicAcQuantities-LineNeutralAcRmsVoltage: %w", err)
	} else {
		val.LineNeutralAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0078125); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseABasicAcQuantities-AcFrequency: %w", err)
	} else {
		val.AcFrequency = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityPhaseABasicAcQuantities-AcRmsCurrent: %w", err)
	} else {
		val.AcRmsCurrent = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeUtilityPhaseABasicAcQuantities(val *UtilityPhaseABasicAcQuantities) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.LineLineAcRmsVoltage, 16)
	w.writeUInt16(val.LineNeutralAcRmsVoltage, 16)
	w.writeUnsignedResolution(val.AcFrequency, 16, 0.0078125)
	w.writeUInt16(val.AcRmsCurrent, 16)
	return w.Bytes(), w.Err()
}

func encodeUtilityPhaseABasicAcQuantitiesMsg(v Message) ([]byte, error) {
	val, ok := v.(*UtilityPhaseABasicAcQuantities)
	if !ok {
		return nil, fmt.Errorf("expected *UtilityPhaseABasicAcQuantities, got %T", v)
	}
	return EncodeUtilityPhaseABasicAcQuantities(val)
}

type UtilityTotalAcReactivePower struct {
	Info               MessageInfo      `json:"info"`
	ReactivePower      *int32           `json:"reactivePower"`
	PowerFactor        *float32         `json:"powerFactor"`
	PowerFactorLagging PowerFactorConst `json:"powerFactorLagging"`
}

func (u *UtilityTotalAcReactivePower) PGNNumber() uint32 { return 65015 }

func DecodeUtilityTotalAcReactivePower(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &UtilityTotalAcReactivePower{}
	val.Info = Info
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityTotalAcReactivePower-ReactivePower: %w", err)
	} else {
		val.ReactivePower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 6.10352e-05); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityTotalAcReactivePower-PowerFactor: %w", err)
	} else {
		val.PowerFactor = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityTotalAcReactivePower-PowerFactorLagging: %w", err)
	} else {
		val.PowerFactorLagging = PowerFactorConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(14)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeUtilityTotalAcReactivePower(val *UtilityTotalAcReactivePower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeInt32(val.ReactivePower, 32)
	w.writeUnsignedResolution(val.PowerFactor, 16, 6.10352e-05)
	w.writeLookupField(uint64(val.PowerFactorLagging), 2)
	w.writeReservedBits(14)
	return w.Bytes(), w.Err()
}

func encodeUtilityTotalAcReactivePowerMsg(v Message) ([]byte, error) {
	val, ok := v.(*UtilityTotalAcReactivePower)
	if !ok {
		return nil, fmt.Errorf("expected *UtilityTotalAcReactivePower, got %T", v)
	}
	return EncodeUtilityTotalAcReactivePower(val)
}

type UtilityTotalAcPower struct {
	Info          MessageInfo `json:"info"`
	RealPower     *int32      `json:"realPower"`
	ApparentPower *int32      `json:"apparentPower"`
}

func (u *UtilityTotalAcPower) PGNNumber() uint32 { return 65016 }

func DecodeUtilityTotalAcPower(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &UtilityTotalAcPower{}
	val.Info = Info
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityTotalAcPower-RealPower: %w", err)
	} else {
		val.RealPower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityTotalAcPower-ApparentPower: %w", err)
	} else {
		val.ApparentPower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeUtilityTotalAcPower(val *UtilityTotalAcPower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeInt32(val.RealPower, 32)
	w.writeInt32(val.ApparentPower, 32)
	return w.Bytes(), w.Err()
}

func encodeUtilityTotalAcPowerMsg(v Message) ([]byte, error) {
	val, ok := v.(*UtilityTotalAcPower)
	if !ok {
		return nil, fmt.Errorf("expected *UtilityTotalAcPower, got %T", v)
	}
	return EncodeUtilityTotalAcPower(val)
}

type UtilityAverageBasicAcQuantities struct {
	Info                    MessageInfo `json:"info"`
	LineLineAcRmsVoltage    *uint16     `json:"lineLineAcRmsVoltage"`
	LineNeutralAcRmsVoltage *uint16     `json:"lineNeutralAcRmsVoltage"`
	AcFrequency             *float32    `json:"acFrequency"`
	AcRmsCurrent            *uint16     `json:"acRmsCurrent"`
}

func (u *UtilityAverageBasicAcQuantities) PGNNumber() uint32 { return 65017 }

func DecodeUtilityAverageBasicAcQuantities(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &UtilityAverageBasicAcQuantities{}
	val.Info = Info
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityAverageBasicAcQuantities-LineLineAcRmsVoltage: %w", err)
	} else {
		val.LineLineAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityAverageBasicAcQuantities-LineNeutralAcRmsVoltage: %w", err)
	} else {
		val.LineNeutralAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0078125); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityAverageBasicAcQuantities-AcFrequency: %w", err)
	} else {
		val.AcFrequency = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for UtilityAverageBasicAcQuantities-AcRmsCurrent: %w", err)
	} else {
		val.AcRmsCurrent = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeUtilityAverageBasicAcQuantities(val *UtilityAverageBasicAcQuantities) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.LineLineAcRmsVoltage, 16)
	w.writeUInt16(val.LineNeutralAcRmsVoltage, 16)
	w.writeUnsignedResolution(val.AcFrequency, 16, 0.0078125)
	w.writeUInt16(val.AcRmsCurrent, 16)
	return w.Bytes(), w.Err()
}

func encodeUtilityAverageBasicAcQuantitiesMsg(v Message) ([]byte, error) {
	val, ok := v.(*UtilityAverageBasicAcQuantities)
	if !ok {
		return nil, fmt.Errorf("expected *UtilityAverageBasicAcQuantities, got %T", v)
	}
	return EncodeUtilityAverageBasicAcQuantities(val)
}

type GeneratorTotalAcEnergy struct {
	Info              MessageInfo `json:"info"`
	TotalEnergyExport *uint32     `json:"totalEnergyExport"`
	TotalEnergyImport *uint32     `json:"totalEnergyImport"`
}

func (g *GeneratorTotalAcEnergy) PGNNumber() uint32 { return 65018 }

func DecodeGeneratorTotalAcEnergy(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GeneratorTotalAcEnergy{}
	val.Info = Info
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorTotalAcEnergy-TotalEnergyExport: %w", err)
	} else {
		val.TotalEnergyExport = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorTotalAcEnergy-TotalEnergyImport: %w", err)
	} else {
		val.TotalEnergyImport = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeGeneratorTotalAcEnergy(val *GeneratorTotalAcEnergy) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt32(val.TotalEnergyExport, 32)
	w.writeUInt32(val.TotalEnergyImport, 32)
	return w.Bytes(), w.Err()
}

func encodeGeneratorTotalAcEnergyMsg(v Message) ([]byte, error) {
	val, ok := v.(*GeneratorTotalAcEnergy)
	if !ok {
		return nil, fmt.Errorf("expected *GeneratorTotalAcEnergy, got %T", v)
	}
	return EncodeGeneratorTotalAcEnergy(val)
}

type GeneratorPhaseCAcReactivePower struct {
	Info               MessageInfo      `json:"info"`
	ReactivePower      *int32           `json:"reactivePower"`
	PowerFactor        *float32         `json:"powerFactor"`
	PowerFactorLagging PowerFactorConst `json:"powerFactorLagging"`
}

func (g *GeneratorPhaseCAcReactivePower) PGNNumber() uint32 { return 65019 }

func DecodeGeneratorPhaseCAcReactivePower(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GeneratorPhaseCAcReactivePower{}
	val.Info = Info
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseCAcReactivePower-ReactivePower: %w", err)
	} else {
		val.ReactivePower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 6.10352e-05); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseCAcReactivePower-PowerFactor: %w", err)
	} else {
		val.PowerFactor = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseCAcReactivePower-PowerFactorLagging: %w", err)
	} else {
		val.PowerFactorLagging = PowerFactorConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(14)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeGeneratorPhaseCAcReactivePower(val *GeneratorPhaseCAcReactivePower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeInt32(val.ReactivePower, 32)
	w.writeUnsignedResolution(val.PowerFactor, 16, 6.10352e-05)
	w.writeLookupField(uint64(val.PowerFactorLagging), 2)
	w.writeReservedBits(14)
	return w.Bytes(), w.Err()
}

func encodeGeneratorPhaseCAcReactivePowerMsg(v Message) ([]byte, error) {
	val, ok := v.(*GeneratorPhaseCAcReactivePower)
	if !ok {
		return nil, fmt.Errorf("expected *GeneratorPhaseCAcReactivePower, got %T", v)
	}
	return EncodeGeneratorPhaseCAcReactivePower(val)
}

type GeneratorPhaseCAcPower struct {
	Info          MessageInfo `json:"info"`
	RealPower     *int32      `json:"realPower"`
	ApparentPower *int32      `json:"apparentPower"`
}

func (g *GeneratorPhaseCAcPower) PGNNumber() uint32 { return 65020 }

func DecodeGeneratorPhaseCAcPower(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GeneratorPhaseCAcPower{}
	val.Info = Info
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseCAcPower-RealPower: %w", err)
	} else {
		val.RealPower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseCAcPower-ApparentPower: %w", err)
	} else {
		val.ApparentPower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeGeneratorPhaseCAcPower(val *GeneratorPhaseCAcPower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeInt32(val.RealPower, 32)
	w.writeInt32(val.ApparentPower, 32)
	return w.Bytes(), w.Err()
}

func encodeGeneratorPhaseCAcPowerMsg(v Message) ([]byte, error) {
	val, ok := v.(*GeneratorPhaseCAcPower)
	if !ok {
		return nil, fmt.Errorf("expected *GeneratorPhaseCAcPower, got %T", v)
	}
	return EncodeGeneratorPhaseCAcPower(val)
}

type GeneratorPhaseCBasicAcQuantities struct {
	Info                    MessageInfo `json:"info"`
	LineLineAcRmsVoltage    *uint16     `json:"lineLineAcRmsVoltage"`
	LineNeutralAcRmsVoltage *uint16     `json:"lineNeutralAcRmsVoltage"`
	AcFrequency             *float32    `json:"acFrequency"`
	AcRmsCurrent            *uint16     `json:"acRmsCurrent"`
}

func (g *GeneratorPhaseCBasicAcQuantities) PGNNumber() uint32 { return 65021 }

func DecodeGeneratorPhaseCBasicAcQuantities(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GeneratorPhaseCBasicAcQuantities{}
	val.Info = Info
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseCBasicAcQuantities-LineLineAcRmsVoltage: %w", err)
	} else {
		val.LineLineAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseCBasicAcQuantities-LineNeutralAcRmsVoltage: %w", err)
	} else {
		val.LineNeutralAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0078125); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseCBasicAcQuantities-AcFrequency: %w", err)
	} else {
		val.AcFrequency = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseCBasicAcQuantities-AcRmsCurrent: %w", err)
	} else {
		val.AcRmsCurrent = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeGeneratorPhaseCBasicAcQuantities(val *GeneratorPhaseCBasicAcQuantities) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.LineLineAcRmsVoltage, 16)
	w.writeUInt16(val.LineNeutralAcRmsVoltage, 16)
	w.writeUnsignedResolution(val.AcFrequency, 16, 0.0078125)
	w.writeUInt16(val.AcRmsCurrent, 16)
	return w.Bytes(), w.Err()
}

func encodeGeneratorPhaseCBasicAcQuantitiesMsg(v Message) ([]byte, error) {
	val, ok := v.(*GeneratorPhaseCBasicAcQuantities)
	if !ok {
		return nil, fmt.Errorf("expected *GeneratorPhaseCBasicAcQuantities, got %T", v)
	}
	return EncodeGeneratorPhaseCBasicAcQuantities(val)
}

type GeneratorPhaseBAcReactivePower struct {
	Info               MessageInfo      `json:"info"`
	ReactivePower      *int32           `json:"reactivePower"`
	PowerFactor        *float32         `json:"powerFactor"`
	PowerFactorLagging PowerFactorConst `json:"powerFactorLagging"`
}

func (g *GeneratorPhaseBAcReactivePower) PGNNumber() uint32 { return 65022 }

func DecodeGeneratorPhaseBAcReactivePower(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GeneratorPhaseBAcReactivePower{}
	val.Info = Info
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseBAcReactivePower-ReactivePower: %w", err)
	} else {
		val.ReactivePower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 6.10352e-05); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseBAcReactivePower-PowerFactor: %w", err)
	} else {
		val.PowerFactor = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseBAcReactivePower-PowerFactorLagging: %w", err)
	} else {
		val.PowerFactorLagging = PowerFactorConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(14)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeGeneratorPhaseBAcReactivePower(val *GeneratorPhaseBAcReactivePower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeInt32(val.ReactivePower, 32)
	w.writeUnsignedResolution(val.PowerFactor, 16, 6.10352e-05)
	w.writeLookupField(uint64(val.PowerFactorLagging), 2)
	w.writeReservedBits(14)
	return w.Bytes(), w.Err()
}

func encodeGeneratorPhaseBAcReactivePowerMsg(v Message) ([]byte, error) {
	val, ok := v.(*GeneratorPhaseBAcReactivePower)
	if !ok {
		return nil, fmt.Errorf("expected *GeneratorPhaseBAcReactivePower, got %T", v)
	}
	return EncodeGeneratorPhaseBAcReactivePower(val)
}

type GeneratorPhaseBAcPower struct {
	Info          MessageInfo `json:"info"`
	RealPower     *int32      `json:"realPower"`
	ApparentPower *int32      `json:"apparentPower"`
}

func (g *GeneratorPhaseBAcPower) PGNNumber() uint32 { return 65023 }

func DecodeGeneratorPhaseBAcPower(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GeneratorPhaseBAcPower{}
	val.Info = Info
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseBAcPower-RealPower: %w", err)
	} else {
		val.RealPower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseBAcPower-ApparentPower: %w", err)
	} else {
		val.ApparentPower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeGeneratorPhaseBAcPower(val *GeneratorPhaseBAcPower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeInt32(val.RealPower, 32)
	w.writeInt32(val.ApparentPower, 32)
	return w.Bytes(), w.Err()
}

func encodeGeneratorPhaseBAcPowerMsg(v Message) ([]byte, error) {
	val, ok := v.(*GeneratorPhaseBAcPower)
	if !ok {
		return nil, fmt.Errorf("expected *GeneratorPhaseBAcPower, got %T", v)
	}
	return EncodeGeneratorPhaseBAcPower(val)
}

type GeneratorPhaseBBasicAcQuantities struct {
	Info                    MessageInfo `json:"info"`
	LineLineAcRmsVoltage    *uint16     `json:"lineLineAcRmsVoltage"`
	LineNeutralAcRmsVoltage *uint16     `json:"lineNeutralAcRmsVoltage"`
	AcFrequency             *float32    `json:"acFrequency"`
	AcRmsCurrent            *uint16     `json:"acRmsCurrent"`
}

func (g *GeneratorPhaseBBasicAcQuantities) PGNNumber() uint32 { return 65024 }

func DecodeGeneratorPhaseBBasicAcQuantities(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GeneratorPhaseBBasicAcQuantities{}
	val.Info = Info
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseBBasicAcQuantities-LineLineAcRmsVoltage: %w", err)
	} else {
		val.LineLineAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseBBasicAcQuantities-LineNeutralAcRmsVoltage: %w", err)
	} else {
		val.LineNeutralAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0078125); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseBBasicAcQuantities-AcFrequency: %w", err)
	} else {
		val.AcFrequency = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseBBasicAcQuantities-AcRmsCurrent: %w", err)
	} else {
		val.AcRmsCurrent = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeGeneratorPhaseBBasicAcQuantities(val *GeneratorPhaseBBasicAcQuantities) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.LineLineAcRmsVoltage, 16)
	w.writeUInt16(val.LineNeutralAcRmsVoltage, 16)
	w.writeUnsignedResolution(val.AcFrequency, 16, 0.0078125)
	w.writeUInt16(val.AcRmsCurrent, 16)
	return w.Bytes(), w.Err()
}

func encodeGeneratorPhaseBBasicAcQuantitiesMsg(v Message) ([]byte, error) {
	val, ok := v.(*GeneratorPhaseBBasicAcQuantities)
	if !ok {
		return nil, fmt.Errorf("expected *GeneratorPhaseBBasicAcQuantities, got %T", v)
	}
	return EncodeGeneratorPhaseBBasicAcQuantities(val)
}

type GeneratorPhaseAAcReactivePower struct {
	Info               MessageInfo      `json:"info"`
	ReactivePower      *int32           `json:"reactivePower"`
	PowerFactor        *float32         `json:"powerFactor"`
	PowerFactorLagging PowerFactorConst `json:"powerFactorLagging"`
}

func (g *GeneratorPhaseAAcReactivePower) PGNNumber() uint32 { return 65025 }

func DecodeGeneratorPhaseAAcReactivePower(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GeneratorPhaseAAcReactivePower{}
	val.Info = Info
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseAAcReactivePower-ReactivePower: %w", err)
	} else {
		val.ReactivePower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 6.10352e-05); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseAAcReactivePower-PowerFactor: %w", err)
	} else {
		val.PowerFactor = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseAAcReactivePower-PowerFactorLagging: %w", err)
	} else {
		val.PowerFactorLagging = PowerFactorConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(14)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeGeneratorPhaseAAcReactivePower(val *GeneratorPhaseAAcReactivePower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeInt32(val.ReactivePower, 32)
	w.writeUnsignedResolution(val.PowerFactor, 16, 6.10352e-05)
	w.writeLookupField(uint64(val.PowerFactorLagging), 2)
	w.writeReservedBits(14)
	return w.Bytes(), w.Err()
}

func encodeGeneratorPhaseAAcReactivePowerMsg(v Message) ([]byte, error) {
	val, ok := v.(*GeneratorPhaseAAcReactivePower)
	if !ok {
		return nil, fmt.Errorf("expected *GeneratorPhaseAAcReactivePower, got %T", v)
	}
	return EncodeGeneratorPhaseAAcReactivePower(val)
}

type GeneratorPhaseAAcPower struct {
	Info          MessageInfo `json:"info"`
	RealPower     *int32      `json:"realPower"`
	ApparentPower *int32      `json:"apparentPower"`
}

func (g *GeneratorPhaseAAcPower) PGNNumber() uint32 { return 65026 }

func DecodeGeneratorPhaseAAcPower(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GeneratorPhaseAAcPower{}
	val.Info = Info
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseAAcPower-RealPower: %w", err)
	} else {
		val.RealPower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseAAcPower-ApparentPower: %w", err)
	} else {
		val.ApparentPower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeGeneratorPhaseAAcPower(val *GeneratorPhaseAAcPower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeInt32(val.RealPower, 32)
	w.writeInt32(val.ApparentPower, 32)
	return w.Bytes(), w.Err()
}

func encodeGeneratorPhaseAAcPowerMsg(v Message) ([]byte, error) {
	val, ok := v.(*GeneratorPhaseAAcPower)
	if !ok {
		return nil, fmt.Errorf("expected *GeneratorPhaseAAcPower, got %T", v)
	}
	return EncodeGeneratorPhaseAAcPower(val)
}

type GeneratorPhaseABasicAcQuantities struct {
	Info                    MessageInfo `json:"info"`
	LineLineAcRmsVoltage    *uint16     `json:"lineLineAcRmsVoltage"`
	LineNeutralAcRmsVoltage *uint16     `json:"lineNeutralAcRmsVoltage"`
	AcFrequency             *float32    `json:"acFrequency"`
	AcRmsCurrent            *uint16     `json:"acRmsCurrent"`
}

func (g *GeneratorPhaseABasicAcQuantities) PGNNumber() uint32 { return 65027 }

func DecodeGeneratorPhaseABasicAcQuantities(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GeneratorPhaseABasicAcQuantities{}
	val.Info = Info
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseABasicAcQuantities-LineLineAcRmsVoltage: %w", err)
	} else {
		val.LineLineAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseABasicAcQuantities-LineNeutralAcRmsVoltage: %w", err)
	} else {
		val.LineNeutralAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0078125); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseABasicAcQuantities-AcFrequency: %w", err)
	} else {
		val.AcFrequency = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorPhaseABasicAcQuantities-AcRmsCurrent: %w", err)
	} else {
		val.AcRmsCurrent = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeGeneratorPhaseABasicAcQuantities(val *GeneratorPhaseABasicAcQuantities) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.LineLineAcRmsVoltage, 16)
	w.writeUInt16(val.LineNeutralAcRmsVoltage, 16)
	w.writeUnsignedResolution(val.AcFrequency, 16, 0.0078125)
	w.writeUInt16(val.AcRmsCurrent, 16)
	return w.Bytes(), w.Err()
}

func encodeGeneratorPhaseABasicAcQuantitiesMsg(v Message) ([]byte, error) {
	val, ok := v.(*GeneratorPhaseABasicAcQuantities)
	if !ok {
		return nil, fmt.Errorf("expected *GeneratorPhaseABasicAcQuantities, got %T", v)
	}
	return EncodeGeneratorPhaseABasicAcQuantities(val)
}

type GeneratorTotalAcReactivePower struct {
	Info               MessageInfo      `json:"info"`
	ReactivePower      *int32           `json:"reactivePower"`
	PowerFactor        *float32         `json:"powerFactor"`
	PowerFactorLagging PowerFactorConst `json:"powerFactorLagging"`
}

func (g *GeneratorTotalAcReactivePower) PGNNumber() uint32 { return 65028 }

func DecodeGeneratorTotalAcReactivePower(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GeneratorTotalAcReactivePower{}
	val.Info = Info
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorTotalAcReactivePower-ReactivePower: %w", err)
	} else {
		val.ReactivePower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 6.10352e-05); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorTotalAcReactivePower-PowerFactor: %w", err)
	} else {
		val.PowerFactor = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorTotalAcReactivePower-PowerFactorLagging: %w", err)
	} else {
		val.PowerFactorLagging = PowerFactorConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(14)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeGeneratorTotalAcReactivePower(val *GeneratorTotalAcReactivePower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeInt32(val.ReactivePower, 32)
	w.writeUnsignedResolution(val.PowerFactor, 16, 6.10352e-05)
	w.writeLookupField(uint64(val.PowerFactorLagging), 2)
	w.writeReservedBits(14)
	return w.Bytes(), w.Err()
}

func encodeGeneratorTotalAcReactivePowerMsg(v Message) ([]byte, error) {
	val, ok := v.(*GeneratorTotalAcReactivePower)
	if !ok {
		return nil, fmt.Errorf("expected *GeneratorTotalAcReactivePower, got %T", v)
	}
	return EncodeGeneratorTotalAcReactivePower(val)
}

type GeneratorTotalAcPower struct {
	Info          MessageInfo `json:"info"`
	RealPower     *int32      `json:"realPower"`
	ApparentPower *int32      `json:"apparentPower"`
}

func (g *GeneratorTotalAcPower) PGNNumber() uint32 { return 65029 }

func DecodeGeneratorTotalAcPower(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GeneratorTotalAcPower{}
	val.Info = Info
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorTotalAcPower-RealPower: %w", err)
	} else {
		val.RealPower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorTotalAcPower-ApparentPower: %w", err)
	} else {
		val.ApparentPower = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeGeneratorTotalAcPower(val *GeneratorTotalAcPower) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeInt32(val.RealPower, 32)
	w.writeInt32(val.ApparentPower, 32)
	return w.Bytes(), w.Err()
}

func encodeGeneratorTotalAcPowerMsg(v Message) ([]byte, error) {
	val, ok := v.(*GeneratorTotalAcPower)
	if !ok {
		return nil, fmt.Errorf("expected *GeneratorTotalAcPower, got %T", v)
	}
	return EncodeGeneratorTotalAcPower(val)
}

type GeneratorAverageBasicAcQuantities struct {
	Info                    MessageInfo `json:"info"`
	LineLineAcRmsVoltage    *uint16     `json:"lineLineAcRmsVoltage"`
	LineNeutralAcRmsVoltage *uint16     `json:"lineNeutralAcRmsVoltage"`
	AcFrequency             *float32    `json:"acFrequency"`
	AcRmsCurrent            *uint16     `json:"acRmsCurrent"`
}

func (g *GeneratorAverageBasicAcQuantities) PGNNumber() uint32 { return 65030 }

func DecodeGeneratorAverageBasicAcQuantities(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GeneratorAverageBasicAcQuantities{}
	val.Info = Info
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorAverageBasicAcQuantities-LineLineAcRmsVoltage: %w", err)
	} else {
		val.LineLineAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorAverageBasicAcQuantities-LineNeutralAcRmsVoltage: %w", err)
	} else {
		val.LineNeutralAcRmsVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0078125); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorAverageBasicAcQuantities-AcFrequency: %w", err)
	} else {
		val.AcFrequency = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for GeneratorAverageBasicAcQuantities-AcRmsCurrent: %w", err)
	} else {
		val.AcRmsCurrent = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeGeneratorAverageBasicAcQuantities(val *GeneratorAverageBasicAcQuantities) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.LineLineAcRmsVoltage, 16)
	w.writeUInt16(val.LineNeutralAcRmsVoltage, 16)
	w.writeUnsignedResolution(val.AcFrequency, 16, 0.0078125)
	w.writeUInt16(val.AcRmsCurrent, 16)
	return w.Bytes(), w.Err()
}

func encodeGeneratorAverageBasicAcQuantitiesMsg(v Message) ([]byte, error) {
	val, ok := v.(*GeneratorAverageBasicAcQuantities)
	if !ok {
		return nil, fmt.Errorf("expected *GeneratorAverageBasicAcQuantities, got %T", v)
	}
	return EncodeGeneratorAverageBasicAcQuantities(val)
}

type LoadControllerConnectionStateControl struct {
	Info                     MessageInfo `json:"info"`
	SequenceId               *uint8      `json:"sequenceId"`
	ConnectionId             *uint8      `json:"connectionId"`
	State                    *uint8      `json:"state"`
	Status                   *uint8      `json:"status"`
	OperationalStatusControl *uint8      `json:"operationalStatusControl"`
	PwmDutyCycle             *uint8      `json:"pwmDutyCycle"`
	Timeon                   *uint8      `json:"timeon"`
	Timeoff                  *uint8      `json:"timeoff"`
}

func (l *LoadControllerConnectionStateControl) PGNNumber() uint32 { return 127500 }

func DecodeLoadControllerConnectionStateControl(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &LoadControllerConnectionStateControl{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for LoadControllerConnectionStateControl-SequenceId: %w", err)
	} else {
		val.SequenceId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for LoadControllerConnectionStateControl-ConnectionId: %w", err)
	} else {
		val.ConnectionId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for LoadControllerConnectionStateControl-State: %w", err)
	} else {
		val.State = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for LoadControllerConnectionStateControl-Status: %w", err)
	} else {
		val.Status = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for LoadControllerConnectionStateControl-OperationalStatusControl: %w", err)
	} else {
		val.OperationalStatusControl = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for LoadControllerConnectionStateControl-PwmDutyCycle: %w", err)
	} else {
		val.PwmDutyCycle = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for LoadControllerConnectionStateControl-Timeon: %w", err)
	} else {
		val.Timeon = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for LoadControllerConnectionStateControl-Timeoff: %w", err)
	} else {
		val.Timeoff = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeLoadControllerConnectionStateControl(val *LoadControllerConnectionStateControl) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.SequenceId, 8)
	w.writeUInt8(val.ConnectionId, 8)
	w.writeUInt8(val.State, 8)
	w.writeUInt8(val.Status, 8)
	w.writeUInt8(val.OperationalStatusControl, 8)
	w.writeUInt8(val.PwmDutyCycle, 8)
	w.writeUInt8(val.Timeon, 8)
	w.writeUInt8(val.Timeoff, 8)
	return w.Bytes(), w.Err()
}

func encodeLoadControllerConnectionStateControlMsg(v Message) ([]byte, error) {
	val, ok := v.(*LoadControllerConnectionStateControl)
	if !ok {
		return nil, fmt.Errorf("expected *LoadControllerConnectionStateControl, got %T", v)
	}
	return EncodeLoadControllerConnectionStateControl(val)
}

type BinarySwitchBankStatus struct {
	Info        MessageInfo `json:"info"`
	Instance    *uint8      `json:"instance"`
	Indicator1  OffOnConst  `json:"indicator1"`
	Indicator2  OffOnConst  `json:"indicator2"`
	Indicator3  OffOnConst  `json:"indicator3"`
	Indicator4  OffOnConst  `json:"indicator4"`
	Indicator5  OffOnConst  `json:"indicator5"`
	Indicator6  OffOnConst  `json:"indicator6"`
	Indicator7  OffOnConst  `json:"indicator7"`
	Indicator8  OffOnConst  `json:"indicator8"`
	Indicator9  OffOnConst  `json:"indicator9"`
	Indicator10 OffOnConst  `json:"indicator10"`
	Indicator11 OffOnConst  `json:"indicator11"`
	Indicator12 OffOnConst  `json:"indicator12"`
	Indicator13 OffOnConst  `json:"indicator13"`
	Indicator14 OffOnConst  `json:"indicator14"`
	Indicator15 OffOnConst  `json:"indicator15"`
	Indicator16 OffOnConst  `json:"indicator16"`
	Indicator17 OffOnConst  `json:"indicator17"`
	Indicator18 OffOnConst  `json:"indicator18"`
	Indicator19 OffOnConst  `json:"indicator19"`
	Indicator20 OffOnConst  `json:"indicator20"`
	Indicator21 OffOnConst  `json:"indicator21"`
	Indicator22 OffOnConst  `json:"indicator22"`
	Indicator23 OffOnConst  `json:"indicator23"`
	Indicator24 OffOnConst  `json:"indicator24"`
	Indicator25 OffOnConst  `json:"indicator25"`
	Indicator26 OffOnConst  `json:"indicator26"`
	Indicator27 OffOnConst  `json:"indicator27"`
	Indicator28 OffOnConst  `json:"indicator28"`
}

func (b *BinarySwitchBankStatus) PGNNumber() uint32 { return 127501 }

func DecodeBinarySwitchBankStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &BinarySwitchBankStatus{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator1: %w", err)
	} else {
		val.Indicator1 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator2: %w", err)
	} else {
		val.Indicator2 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator3: %w", err)
	} else {
		val.Indicator3 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator4: %w", err)
	} else {
		val.Indicator4 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator5: %w", err)
	} else {
		val.Indicator5 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator6: %w", err)
	} else {
		val.Indicator6 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator7: %w", err)
	} else {
		val.Indicator7 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator8: %w", err)
	} else {
		val.Indicator8 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator9: %w", err)
	} else {
		val.Indicator9 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator10: %w", err)
	} else {
		val.Indicator10 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator11: %w", err)
	} else {
		val.Indicator11 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator12: %w", err)
	} else {
		val.Indicator12 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator13: %w", err)
	} else {
		val.Indicator13 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator14: %w", err)
	} else {
		val.Indicator14 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator15: %w", err)
	} else {
		val.Indicator15 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator16: %w", err)
	} else {
		val.Indicator16 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator17: %w", err)
	} else {
		val.Indicator17 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator18: %w", err)
	} else {
		val.Indicator18 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator19: %w", err)
	} else {
		val.Indicator19 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator20: %w", err)
	} else {
		val.Indicator20 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator21: %w", err)
	} else {
		val.Indicator21 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator22: %w", err)
	} else {
		val.Indicator22 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator23: %w", err)
	} else {
		val.Indicator23 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator24: %w", err)
	} else {
		val.Indicator24 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator25: %w", err)
	} else {
		val.Indicator25 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator26: %w", err)
	} else {
		val.Indicator26 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator27: %w", err)
	} else {
		val.Indicator27 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BinarySwitchBankStatus-Indicator28: %w", err)
	} else {
		val.Indicator28 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeBinarySwitchBankStatus(val *BinarySwitchBankStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Instance, 8)
	w.writeLookupField(uint64(val.Indicator1), 2)
	w.writeLookupField(uint64(val.Indicator2), 2)
	w.writeLookupField(uint64(val.Indicator3), 2)
	w.writeLookupField(uint64(val.Indicator4), 2)
	w.writeLookupField(uint64(val.Indicator5), 2)
	w.writeLookupField(uint64(val.Indicator6), 2)
	w.writeLookupField(uint64(val.Indicator7), 2)
	w.writeLookupField(uint64(val.Indicator8), 2)
	w.writeLookupField(uint64(val.Indicator9), 2)
	w.writeLookupField(uint64(val.Indicator10), 2)
	w.writeLookupField(uint64(val.Indicator11), 2)
	w.writeLookupField(uint64(val.Indicator12), 2)
	w.writeLookupField(uint64(val.Indicator13), 2)
	w.writeLookupField(uint64(val.Indicator14), 2)
	w.writeLookupField(uint64(val.Indicator15), 2)
	w.writeLookupField(uint64(val.Indicator16), 2)
	w.writeLookupField(uint64(val.Indicator17), 2)
	w.writeLookupField(uint64(val.Indicator18), 2)
	w.writeLookupField(uint64(val.Indicator19), 2)
	w.writeLookupField(uint64(val.Indicator20), 2)
	w.writeLookupField(uint64(val.Indicator21), 2)
	w.writeLookupField(uint64(val.Indicator22), 2)
	w.writeLookupField(uint64(val.Indicator23), 2)
	w.writeLookupField(uint64(val.Indicator24), 2)
	w.writeLookupField(uint64(val.Indicator25), 2)
	w.writeLookupField(uint64(val.Indicator26), 2)
	w.writeLookupField(uint64(val.Indicator27), 2)
	w.writeLookupField(uint64(val.Indicator28), 2)
	return w.Bytes(), w.Err()
}

func encodeBinarySwitchBankStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*BinarySwitchBankStatus)
	if !ok {
		return nil, fmt.Errorf("expected *BinarySwitchBankStatus, got %T", v)
	}
	return EncodeBinarySwitchBankStatus(val)
}

type SwitchBankControl struct {
	Info     MessageInfo `json:"info"`
	Instance *uint8      `json:"instance"`
	Switch1  OffOnConst  `json:"switch1"`
	Switch2  OffOnConst  `json:"switch2"`
	Switch3  OffOnConst  `json:"switch3"`
	Switch4  OffOnConst  `json:"switch4"`
	Switch5  OffOnConst  `json:"switch5"`
	Switch6  OffOnConst  `json:"switch6"`
	Switch7  OffOnConst  `json:"switch7"`
	Switch8  OffOnConst  `json:"switch8"`
	Switch9  OffOnConst  `json:"switch9"`
	Switch10 OffOnConst  `json:"switch10"`
	Switch11 OffOnConst  `json:"switch11"`
	Switch12 OffOnConst  `json:"switch12"`
	Switch13 OffOnConst  `json:"switch13"`
	Switch14 OffOnConst  `json:"switch14"`
	Switch15 OffOnConst  `json:"switch15"`
	Switch16 OffOnConst  `json:"switch16"`
	Switch17 OffOnConst  `json:"switch17"`
	Switch18 OffOnConst  `json:"switch18"`
	Switch19 OffOnConst  `json:"switch19"`
	Switch20 OffOnConst  `json:"switch20"`
	Switch21 OffOnConst  `json:"switch21"`
	Switch22 OffOnConst  `json:"switch22"`
	Switch23 OffOnConst  `json:"switch23"`
	Switch24 OffOnConst  `json:"switch24"`
	Switch25 OffOnConst  `json:"switch25"`
	Switch26 OffOnConst  `json:"switch26"`
	Switch27 OffOnConst  `json:"switch27"`
	Switch28 OffOnConst  `json:"switch28"`
}

func (s *SwitchBankControl) PGNNumber() uint32 { return 127502 }

func DecodeSwitchBankControl(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SwitchBankControl{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch1: %w", err)
	} else {
		val.Switch1 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch2: %w", err)
	} else {
		val.Switch2 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch3: %w", err)
	} else {
		val.Switch3 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch4: %w", err)
	} else {
		val.Switch4 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch5: %w", err)
	} else {
		val.Switch5 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch6: %w", err)
	} else {
		val.Switch6 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch7: %w", err)
	} else {
		val.Switch7 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch8: %w", err)
	} else {
		val.Switch8 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch9: %w", err)
	} else {
		val.Switch9 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch10: %w", err)
	} else {
		val.Switch10 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch11: %w", err)
	} else {
		val.Switch11 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch12: %w", err)
	} else {
		val.Switch12 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch13: %w", err)
	} else {
		val.Switch13 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch14: %w", err)
	} else {
		val.Switch14 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch15: %w", err)
	} else {
		val.Switch15 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch16: %w", err)
	} else {
		val.Switch16 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch17: %w", err)
	} else {
		val.Switch17 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch18: %w", err)
	} else {
		val.Switch18 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch19: %w", err)
	} else {
		val.Switch19 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch20: %w", err)
	} else {
		val.Switch20 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch21: %w", err)
	} else {
		val.Switch21 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch22: %w", err)
	} else {
		val.Switch22 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch23: %w", err)
	} else {
		val.Switch23 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch24: %w", err)
	} else {
		val.Switch24 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch25: %w", err)
	} else {
		val.Switch25 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch26: %w", err)
	} else {
		val.Switch26 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch27: %w", err)
	} else {
		val.Switch27 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SwitchBankControl-Switch28: %w", err)
	} else {
		val.Switch28 = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSwitchBankControl(val *SwitchBankControl) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Instance, 8)
	w.writeLookupField(uint64(val.Switch1), 2)
	w.writeLookupField(uint64(val.Switch2), 2)
	w.writeLookupField(uint64(val.Switch3), 2)
	w.writeLookupField(uint64(val.Switch4), 2)
	w.writeLookupField(uint64(val.Switch5), 2)
	w.writeLookupField(uint64(val.Switch6), 2)
	w.writeLookupField(uint64(val.Switch7), 2)
	w.writeLookupField(uint64(val.Switch8), 2)
	w.writeLookupField(uint64(val.Switch9), 2)
	w.writeLookupField(uint64(val.Switch10), 2)
	w.writeLookupField(uint64(val.Switch11), 2)
	w.writeLookupField(uint64(val.Switch12), 2)
	w.writeLookupField(uint64(val.Switch13), 2)
	w.writeLookupField(uint64(val.Switch14), 2)
	w.writeLookupField(uint64(val.Switch15), 2)
	w.writeLookupField(uint64(val.Switch16), 2)
	w.writeLookupField(uint64(val.Switch17), 2)
	w.writeLookupField(uint64(val.Switch18), 2)
	w.writeLookupField(uint64(val.Switch19), 2)
	w.writeLookupField(uint64(val.Switch20), 2)
	w.writeLookupField(uint64(val.Switch21), 2)
	w.writeLookupField(uint64(val.Switch22), 2)
	w.writeLookupField(uint64(val.Switch23), 2)
	w.writeLookupField(uint64(val.Switch24), 2)
	w.writeLookupField(uint64(val.Switch25), 2)
	w.writeLookupField(uint64(val.Switch26), 2)
	w.writeLookupField(uint64(val.Switch27), 2)
	w.writeLookupField(uint64(val.Switch28), 2)
	return w.Bytes(), w.Err()
}

func encodeSwitchBankControlMsg(v Message) ([]byte, error) {
	val, ok := v.(*SwitchBankControl)
	if !ok {
		return nil, fmt.Errorf("expected *SwitchBankControl, got %T", v)
	}
	return EncodeSwitchBankControl(val)
}

type AcInputStatus struct {
	Info          MessageInfo               `json:"info"`
	Instance      *uint8                    `json:"instance"`
	NumberOfLines *uint8                    `json:"numberOfLines"`
	Repeating1    []AcInputStatusRepeating1 `json:"repeating1"`
}

type AcInputStatusRepeating1 struct {
	Line          *uint8             `json:"line"`
	Acceptability AcceptabilityConst `json:"acceptability"`
	Voltage       *float32           `json:"voltage"`
	Current       *float32           `json:"current"`
	Frequency     *float32           `json:"frequency"`
	BreakerSize   *float32           `json:"breakerSize"`
	RealPower     *uint32            `json:"realPower"`
	ReactivePower *uint32            `json:"reactivePower"`
	PowerFactor   *float32           `json:"powerFactor"`
}

func (a *AcInputStatus) PGNNumber() uint32 { return 127503 }

func DecodeAcInputStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AcInputStatus{}
	val.Info = Info
	var repeat1Count uint16 = 0
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AcInputStatus-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AcInputStatus-NumberOfLines: %w", err)
	} else {
		val.NumberOfLines = v
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
	val.Repeating1 = make([]AcInputStatusRepeating1, 0)
	i := 0
	for {
		var rep AcInputStatusRepeating1
		if v, err := stream.readUInt8(2); err != nil {
			return nil, fmt.Errorf("parse failed for AcInputStatus-Line: %w", err)
		} else {
			rep.Line = v
		}
		if v, err := stream.readLookupField(2); err != nil {
			return nil, fmt.Errorf("parse failed for AcInputStatus-Acceptability: %w", err)
		} else {
			rep.Acceptability = AcceptabilityConst(v)
		}
		stream.skipBits(4)
		if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
			return nil, fmt.Errorf("parse failed for AcInputStatus-Voltage: %w", err)
		} else {
			rep.Voltage = v
		}
		if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
			return nil, fmt.Errorf("parse failed for AcInputStatus-Current: %w", err)
		} else {
			rep.Current = v
		}
		if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
			return nil, fmt.Errorf("parse failed for AcInputStatus-Frequency: %w", err)
		} else {
			rep.Frequency = v
		}
		if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
			return nil, fmt.Errorf("parse failed for AcInputStatus-BreakerSize: %w", err)
		} else {
			rep.BreakerSize = v
		}
		if v, err := stream.readUInt32(32); err != nil {
			return nil, fmt.Errorf("parse failed for AcInputStatus-RealPower: %w", err)
		} else {
			rep.RealPower = v
		}
		if v, err := stream.readUInt32(32); err != nil {
			return nil, fmt.Errorf("parse failed for AcInputStatus-ReactivePower: %w", err)
		} else {
			rep.ReactivePower = v
		}
		if v, err := stream.readUnsignedResolution(8, 0.01); err != nil {
			return nil, fmt.Errorf("parse failed for AcInputStatus-PowerFactor: %w", err)
		} else {
			rep.PowerFactor = v
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

func EncodeAcInputStatus(val *AcInputStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Instance, 8)
	w.writeUInt8(val.NumberOfLines, 8)
	for _, rep := range val.Repeating1 {
		w.writeUInt8(rep.Line, 2)
		w.writeLookupField(uint64(rep.Acceptability), 2)
		w.writeReservedBits(4)
		w.writeUnsignedResolution(rep.Voltage, 16, 0.01)
		w.writeUnsignedResolution(rep.Current, 16, 0.1)
		w.writeUnsignedResolution(rep.Frequency, 16, 0.01)
		w.writeUnsignedResolution(rep.BreakerSize, 16, 0.1)
		w.writeUInt32(rep.RealPower, 32)
		w.writeUInt32(rep.ReactivePower, 32)
		w.writeUnsignedResolution(rep.PowerFactor, 8, 0.01)
	}
	w.skipBits(4)
	return w.Bytes(), w.Err()
}

func encodeAcInputStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*AcInputStatus)
	if !ok {
		return nil, fmt.Errorf("expected *AcInputStatus, got %T", v)
	}
	return EncodeAcInputStatus(val)
}

type AcOutputStatus struct {
	Info          MessageInfo                `json:"info"`
	Instance      *uint8                     `json:"instance"`
	NumberOfLines *uint8                     `json:"numberOfLines"`
	Repeating1    []AcOutputStatusRepeating1 `json:"repeating1"`
}

type AcOutputStatusRepeating1 struct {
	Line          LineConst     `json:"line"`
	Waveform      WaveformConst `json:"waveform"`
	Voltage       *float32      `json:"voltage"`
	Current       *float32      `json:"current"`
	Frequency     *float32      `json:"frequency"`
	BreakerSize   *float32      `json:"breakerSize"`
	RealPower     *uint32       `json:"realPower"`
	ReactivePower *uint32       `json:"reactivePower"`
	PowerFactor   *float32      `json:"powerFactor"`
}

func (a *AcOutputStatus) PGNNumber() uint32 { return 127504 }

func DecodeAcOutputStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AcOutputStatus{}
	val.Info = Info
	var repeat1Count uint16 = 0
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AcOutputStatus-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AcOutputStatus-NumberOfLines: %w", err)
	} else {
		val.NumberOfLines = v
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
	val.Repeating1 = make([]AcOutputStatusRepeating1, 0)
	i := 0
	for {
		var rep AcOutputStatusRepeating1
		if v, err := stream.readLookupField(2); err != nil {
			return nil, fmt.Errorf("parse failed for AcOutputStatus-Line: %w", err)
		} else {
			rep.Line = LineConst(v)
		}
		if v, err := stream.readLookupField(3); err != nil {
			return nil, fmt.Errorf("parse failed for AcOutputStatus-Waveform: %w", err)
		} else {
			rep.Waveform = WaveformConst(v)
		}
		stream.skipBits(3)
		if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
			return nil, fmt.Errorf("parse failed for AcOutputStatus-Voltage: %w", err)
		} else {
			rep.Voltage = v
		}
		if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
			return nil, fmt.Errorf("parse failed for AcOutputStatus-Current: %w", err)
		} else {
			rep.Current = v
		}
		if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
			return nil, fmt.Errorf("parse failed for AcOutputStatus-Frequency: %w", err)
		} else {
			rep.Frequency = v
		}
		if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
			return nil, fmt.Errorf("parse failed for AcOutputStatus-BreakerSize: %w", err)
		} else {
			rep.BreakerSize = v
		}
		if v, err := stream.readUInt32(32); err != nil {
			return nil, fmt.Errorf("parse failed for AcOutputStatus-RealPower: %w", err)
		} else {
			rep.RealPower = v
		}
		if v, err := stream.readUInt32(32); err != nil {
			return nil, fmt.Errorf("parse failed for AcOutputStatus-ReactivePower: %w", err)
		} else {
			rep.ReactivePower = v
		}
		if v, err := stream.readUnsignedResolution(8, 0.01); err != nil {
			return nil, fmt.Errorf("parse failed for AcOutputStatus-PowerFactor: %w", err)
		} else {
			rep.PowerFactor = v
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

func EncodeAcOutputStatus(val *AcOutputStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Instance, 8)
	w.writeUInt8(val.NumberOfLines, 8)
	for _, rep := range val.Repeating1 {
		w.writeLookupField(uint64(rep.Line), 2)
		w.writeLookupField(uint64(rep.Waveform), 3)
		w.writeReservedBits(3)
		w.writeUnsignedResolution(rep.Voltage, 16, 0.01)
		w.writeUnsignedResolution(rep.Current, 16, 0.1)
		w.writeUnsignedResolution(rep.Frequency, 16, 0.01)
		w.writeUnsignedResolution(rep.BreakerSize, 16, 0.1)
		w.writeUInt32(rep.RealPower, 32)
		w.writeUInt32(rep.ReactivePower, 32)
		w.writeUnsignedResolution(rep.PowerFactor, 8, 0.01)
	}
	w.skipBits(3)
	return w.Bytes(), w.Err()
}

func encodeAcOutputStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*AcOutputStatus)
	if !ok {
		return nil, fmt.Errorf("expected *AcOutputStatus, got %T", v)
	}
	return EncodeAcOutputStatus(val)
}

type DcDetailedStatus struct {
	Info              MessageInfo   `json:"info"`
	Sid               *uint8        `json:"sid"`
	Instance          *uint8        `json:"instance"`
	DcType            DcSourceConst `json:"dcType"`
	StateOfCharge     *uint8        `json:"stateOfCharge"`
	StateOfHealth     *uint8        `json:"stateOfHealth"`
	TimeRemaining     *float32      `json:"timeRemaining"`
	RippleVoltage     *float32      `json:"rippleVoltage"`
	RemainingCapacity *uint16       `json:"remainingCapacity"`
}

func (d *DcDetailedStatus) PGNNumber() uint32 { return 127506 }

func DecodeDcDetailedStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &DcDetailedStatus{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for DcDetailedStatus-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for DcDetailedStatus-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for DcDetailedStatus-DcType: %w", err)
	} else {
		val.DcType = DcSourceConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for DcDetailedStatus-StateOfCharge: %w", err)
	} else {
		val.StateOfCharge = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for DcDetailedStatus-StateOfHealth: %w", err)
	} else {
		val.StateOfHealth = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 60); err != nil {
		return nil, fmt.Errorf("parse failed for DcDetailedStatus-TimeRemaining: %w", err)
	} else {
		val.TimeRemaining = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for DcDetailedStatus-RippleVoltage: %w", err)
	} else {
		val.RippleVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for DcDetailedStatus-RemainingCapacity: %w", err)
	} else {
		val.RemainingCapacity = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeDcDetailedStatus(val *DcDetailedStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.Instance, 8)
	w.writeLookupField(uint64(val.DcType), 8)
	w.writeUInt8(val.StateOfCharge, 8)
	w.writeUInt8(val.StateOfHealth, 8)
	w.writeUnsignedResolution(val.TimeRemaining, 16, 60)
	w.writeUnsignedResolution(val.RippleVoltage, 16, 0.01)
	w.writeUInt16(val.RemainingCapacity, 16)
	return w.Bytes(), w.Err()
}

func encodeDcDetailedStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*DcDetailedStatus)
	if !ok {
		return nil, fmt.Errorf("expected *DcDetailedStatus, got %T", v)
	}
	return EncodeDcDetailedStatus(val)
}

type ChargerStatus struct {
	Info                      MessageInfo       `json:"info"`
	Instance                  *uint8            `json:"instance"`
	BatteryInstance           *uint8            `json:"batteryInstance"`
	OperatingState            ChargerStateConst `json:"operatingState"`
	ChargeMode                ChargerModeConst  `json:"chargeMode"`
	Enabled                   OffOnConst        `json:"enabled"`
	EqualizationPending       OffOnConst        `json:"equalizationPending"`
	EqualizationTimeRemaining *float32          `json:"equalizationTimeRemaining"`
}

func (c *ChargerStatus) PGNNumber() uint32 { return 127507 }

func DecodeChargerStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &ChargerStatus{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ChargerStatus-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ChargerStatus-BatteryInstance: %w", err)
	} else {
		val.BatteryInstance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for ChargerStatus-OperatingState: %w", err)
	} else {
		val.OperatingState = ChargerStateConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for ChargerStatus-ChargeMode: %w", err)
	} else {
		val.ChargeMode = ChargerModeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for ChargerStatus-Enabled: %w", err)
	} else {
		val.Enabled = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for ChargerStatus-EqualizationPending: %w", err)
	} else {
		val.EqualizationPending = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUnsignedResolution(16, 60); err != nil {
		return nil, fmt.Errorf("parse failed for ChargerStatus-EqualizationTimeRemaining: %w", err)
	} else {
		val.EqualizationTimeRemaining = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeChargerStatus(val *ChargerStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Instance, 8)
	w.writeUInt8(val.BatteryInstance, 8)
	w.writeLookupField(uint64(val.OperatingState), 4)
	w.writeLookupField(uint64(val.ChargeMode), 4)
	w.writeLookupField(uint64(val.Enabled), 2)
	w.writeLookupField(uint64(val.EqualizationPending), 2)
	w.writeReservedBits(4)
	w.writeUnsignedResolution(val.EqualizationTimeRemaining, 16, 60)
	return w.Bytes(), w.Err()
}

func encodeChargerStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*ChargerStatus)
	if !ok {
		return nil, fmt.Errorf("expected *ChargerStatus, got %T", v)
	}
	return EncodeChargerStatus(val)
}

type BatteryStatus struct {
	Info        MessageInfo        `json:"info"`
	Instance    *uint8             `json:"instance"`
	Voltage     *float32           `json:"voltage"`
	Current     *float32           `json:"current"`
	Temperature *units.Temperature `json:"temperature"`
	Sid         *uint8             `json:"sid"`
}

func (b *BatteryStatus) PGNNumber() uint32 { return 127508 }

func DecodeBatteryStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &BatteryStatus{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for BatteryStatus-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for BatteryStatus-Voltage: %w", err)
	} else {
		val.Voltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for BatteryStatus-Current: %w", err)
	} else {
		val.Current = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for BatteryStatus-Temperature: %w", err)
	} else {
		val.Temperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for BatteryStatus-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeBatteryStatus(val *BatteryStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Instance, 8)
	w.writeUnsignedResolution(val.Voltage, 16, 0.01)
	w.writeSignedResolution(val.Current, 16, 0.1)
	var temperatureRaw *float32
	if val.Temperature != nil {
		temperatureRaw = &val.Temperature.Value
	}
	w.writeUnsignedResolution(temperatureRaw, 16, 0.01)
	w.writeUInt8(val.Sid, 8)
	return w.Bytes(), w.Err()
}

func encodeBatteryStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*BatteryStatus)
	if !ok {
		return nil, fmt.Errorf("expected *BatteryStatus, got %T", v)
	}
	return EncodeBatteryStatus(val)
}

type InverterStatus struct {
	Info           MessageInfo        `json:"info"`
	Instance       *uint8             `json:"instance"`
	AcInstance     *uint8             `json:"acInstance"`
	DcInstance     *uint8             `json:"dcInstance"`
	OperatingState InverterStateConst `json:"operatingState"`
	InverterEnable OffOnConst         `json:"inverterEnable"`
}

func (i *InverterStatus) PGNNumber() uint32 { return 127509 }

func DecodeInverterStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &InverterStatus{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for InverterStatus-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for InverterStatus-AcInstance: %w", err)
	} else {
		val.AcInstance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for InverterStatus-DcInstance: %w", err)
	} else {
		val.DcInstance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for InverterStatus-OperatingState: %w", err)
	} else {
		val.OperatingState = InverterStateConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for InverterStatus-InverterEnable: %w", err)
	} else {
		val.InverterEnable = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeInverterStatus(val *InverterStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Instance, 8)
	w.writeUInt8(val.AcInstance, 8)
	w.writeUInt8(val.DcInstance, 8)
	w.writeLookupField(uint64(val.OperatingState), 4)
	w.writeLookupField(uint64(val.InverterEnable), 2)
	w.writeReservedBits(2)
	return w.Bytes(), w.Err()
}

func encodeInverterStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*InverterStatus)
	if !ok {
		return nil, fmt.Errorf("expected *InverterStatus, got %T", v)
	}
	return EncodeInverterStatus(val)
}

type InverterConfigurationStatus struct {
	Info                    MessageInfo `json:"info"`
	Instance                *uint8      `json:"instance"`
	AcInstance              *uint8      `json:"acInstance"`
	DcInstance              *uint8      `json:"dcInstance"`
	InverterEnableDisable   *uint8      `json:"inverterEnableDisable"`
	InverterMode            *uint8      `json:"inverterMode"`
	LoadSenseEnableDisable  *uint8      `json:"loadSenseEnableDisable"`
	LoadSensePowerThreshold *uint8      `json:"loadSensePowerThreshold"`
	LoadSenseInterval       *uint8      `json:"loadSenseInterval"`
}

func (i *InverterConfigurationStatus) PGNNumber() uint32 { return 127511 }

func DecodeInverterConfigurationStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &InverterConfigurationStatus{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for InverterConfigurationStatus-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for InverterConfigurationStatus-AcInstance: %w", err)
	} else {
		val.AcInstance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for InverterConfigurationStatus-DcInstance: %w", err)
	} else {
		val.DcInstance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(2); err != nil {
		return nil, fmt.Errorf("parse failed for InverterConfigurationStatus-InverterEnableDisable: %w", err)
	} else {
		val.InverterEnableDisable = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for InverterConfigurationStatus-InverterMode: %w", err)
	} else {
		val.InverterMode = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for InverterConfigurationStatus-LoadSenseEnableDisable: %w", err)
	} else {
		val.LoadSenseEnableDisable = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for InverterConfigurationStatus-LoadSensePowerThreshold: %w", err)
	} else {
		val.LoadSensePowerThreshold = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for InverterConfigurationStatus-LoadSenseInterval: %w", err)
	} else {
		val.LoadSenseInterval = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeInverterConfigurationStatus(val *InverterConfigurationStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Instance, 8)
	w.writeUInt8(val.AcInstance, 8)
	w.writeUInt8(val.DcInstance, 8)
	w.writeUInt8(val.InverterEnableDisable, 2)
	w.writeReservedBits(6)
	w.writeUInt8(val.InverterMode, 8)
	w.writeUInt8(val.LoadSenseEnableDisable, 8)
	w.writeUInt8(val.LoadSensePowerThreshold, 8)
	w.writeUInt8(val.LoadSenseInterval, 8)
	return w.Bytes(), w.Err()
}

func encodeInverterConfigurationStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*InverterConfigurationStatus)
	if !ok {
		return nil, fmt.Errorf("expected *InverterConfigurationStatus, got %T", v)
	}
	return EncodeInverterConfigurationStatus(val)
}

type AgsConfigurationStatus struct {
	Info              MessageInfo `json:"info"`
	Instance          *uint8      `json:"instance"`
	GeneratorInstance *uint8      `json:"generatorInstance"`
	AgsMode           *uint8      `json:"agsMode"`
}

func (a *AgsConfigurationStatus) PGNNumber() uint32 { return 127512 }

func DecodeAgsConfigurationStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AgsConfigurationStatus{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AgsConfigurationStatus-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AgsConfigurationStatus-GeneratorInstance: %w", err)
	} else {
		val.GeneratorInstance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AgsConfigurationStatus-AgsMode: %w", err)
	} else {
		val.AgsMode = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(40)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeAgsConfigurationStatus(val *AgsConfigurationStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Instance, 8)
	w.writeUInt8(val.GeneratorInstance, 8)
	w.writeUInt8(val.AgsMode, 8)
	w.writeReservedBits(40)
	return w.Bytes(), w.Err()
}

func encodeAgsConfigurationStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*AgsConfigurationStatus)
	if !ok {
		return nil, fmt.Errorf("expected *AgsConfigurationStatus, got %T", v)
	}
	return EncodeAgsConfigurationStatus(val)
}

type BatteryConfigurationStatus struct {
	Info                   MessageInfo           `json:"info"`
	Instance               *uint8                `json:"instance"`
	BatteryType            BatteryTypeConst      `json:"batteryType"`
	SupportsEqualization   YesNoConst            `json:"supportsEqualization"`
	NominalVoltage         BatteryVoltageConst   `json:"nominalVoltage"`
	Chemistry              BatteryChemistryConst `json:"chemistry"`
	Capacity               *uint16               `json:"capacity"`
	TemperatureCoefficient *int8                 `json:"temperatureCoefficient"`
	PeukertExponent        *float32              `json:"peukertExponent"`
	ChargeEfficiencyFactor *int8                 `json:"chargeEfficiencyFactor"`
}

func (b *BatteryConfigurationStatus) PGNNumber() uint32 { return 127513 }

func DecodeBatteryConfigurationStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &BatteryConfigurationStatus{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for BatteryConfigurationStatus-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for BatteryConfigurationStatus-BatteryType: %w", err)
	} else {
		val.BatteryType = BatteryTypeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for BatteryConfigurationStatus-SupportsEqualization: %w", err)
	} else {
		val.SupportsEqualization = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for BatteryConfigurationStatus-NominalVoltage: %w", err)
	} else {
		val.NominalVoltage = BatteryVoltageConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for BatteryConfigurationStatus-Chemistry: %w", err)
	} else {
		val.Chemistry = BatteryChemistryConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for BatteryConfigurationStatus-Capacity: %w", err)
	} else {
		val.Capacity = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for BatteryConfigurationStatus-TemperatureCoefficient: %w", err)
	} else {
		val.TemperatureCoefficient = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(8, 0.002); err != nil {
		return nil, fmt.Errorf("parse failed for BatteryConfigurationStatus-PeukertExponent: %w", err)
	} else {
		val.PeukertExponent = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for BatteryConfigurationStatus-ChargeEfficiencyFactor: %w", err)
	} else {
		val.ChargeEfficiencyFactor = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeBatteryConfigurationStatus(val *BatteryConfigurationStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Instance, 8)
	w.writeLookupField(uint64(val.BatteryType), 4)
	w.writeLookupField(uint64(val.SupportsEqualization), 2)
	w.writeReservedBits(2)
	w.writeLookupField(uint64(val.NominalVoltage), 4)
	w.writeLookupField(uint64(val.Chemistry), 4)
	w.writeUInt16(val.Capacity, 16)
	w.writeInt8(val.TemperatureCoefficient, 8)
	w.writeUnsignedResolution(val.PeukertExponent, 8, 0.002)
	w.writeInt8(val.ChargeEfficiencyFactor, 8)
	return w.Bytes(), w.Err()
}

func encodeBatteryConfigurationStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*BatteryConfigurationStatus)
	if !ok {
		return nil, fmt.Errorf("expected *BatteryConfigurationStatus, got %T", v)
	}
	return EncodeBatteryConfigurationStatus(val)
}

type AgsStatus struct {
	Info               MessageInfo `json:"info"`
	Instance           *uint8      `json:"instance"`
	GeneratorInstance  *uint8      `json:"generatorInstance"`
	AgsOperatingState  *uint8      `json:"agsOperatingState"`
	GeneratorState     *uint8      `json:"generatorState"`
	GeneratorOnReason  *uint8      `json:"generatorOnReason"`
	GeneratorOffReason *uint8      `json:"generatorOffReason"`
}

func (a *AgsStatus) PGNNumber() uint32 { return 127514 }

func DecodeAgsStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AgsStatus{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AgsStatus-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AgsStatus-GeneratorInstance: %w", err)
	} else {
		val.GeneratorInstance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AgsStatus-AgsOperatingState: %w", err)
	} else {
		val.AgsOperatingState = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AgsStatus-GeneratorState: %w", err)
	} else {
		val.GeneratorState = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AgsStatus-GeneratorOnReason: %w", err)
	} else {
		val.GeneratorOnReason = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AgsStatus-GeneratorOffReason: %w", err)
	} else {
		val.GeneratorOffReason = v

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

func EncodeAgsStatus(val *AgsStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Instance, 8)
	w.writeUInt8(val.GeneratorInstance, 8)
	w.writeUInt8(val.AgsOperatingState, 8)
	w.writeUInt8(val.GeneratorState, 8)
	w.writeUInt8(val.GeneratorOnReason, 8)
	w.writeUInt8(val.GeneratorOffReason, 8)
	w.writeReservedBits(16)
	return w.Bytes(), w.Err()
}

func encodeAgsStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*AgsStatus)
	if !ok {
		return nil, fmt.Errorf("expected *AgsStatus, got %T", v)
	}
	return EncodeAgsStatus(val)
}

type AcPowerCurrentPhaseA struct {
	Info             MessageInfo `json:"info"`
	Sid              *uint8      `json:"sid"`
	ConnectionNumber *uint8      `json:"connectionNumber"`
	AcRmsCurrent     *float32    `json:"acRmsCurrent"`
	Power            *int32      `json:"power"`
}

func (a *AcPowerCurrentPhaseA) PGNNumber() uint32 { return 127744 }

func DecodeAcPowerCurrentPhaseA(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AcPowerCurrentPhaseA{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AcPowerCurrentPhaseA-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AcPowerCurrentPhaseA-ConnectionNumber: %w", err)
	} else {
		val.ConnectionNumber = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AcPowerCurrentPhaseA-AcRmsCurrent: %w", err)
	} else {
		val.AcRmsCurrent = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AcPowerCurrentPhaseA-Power: %w", err)
	} else {
		val.Power = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeAcPowerCurrentPhaseA(val *AcPowerCurrentPhaseA) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.ConnectionNumber, 8)
	w.writeUnsignedResolution(val.AcRmsCurrent, 16, 0.1)
	w.writeInt32(val.Power, 32)
	return w.Bytes(), w.Err()
}

func encodeAcPowerCurrentPhaseAMsg(v Message) ([]byte, error) {
	val, ok := v.(*AcPowerCurrentPhaseA)
	if !ok {
		return nil, fmt.Errorf("expected *AcPowerCurrentPhaseA, got %T", v)
	}
	return EncodeAcPowerCurrentPhaseA(val)
}

type AcPowerCurrentPhaseB struct {
	Info             MessageInfo `json:"info"`
	Sid              *uint8      `json:"sid"`
	ConnectionNumber *uint8      `json:"connectionNumber"`
	AcRmsCurrent     *float32    `json:"acRmsCurrent"`
	Power            *int32      `json:"power"`
}

func (a *AcPowerCurrentPhaseB) PGNNumber() uint32 { return 127745 }

func DecodeAcPowerCurrentPhaseB(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AcPowerCurrentPhaseB{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AcPowerCurrentPhaseB-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AcPowerCurrentPhaseB-ConnectionNumber: %w", err)
	} else {
		val.ConnectionNumber = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AcPowerCurrentPhaseB-AcRmsCurrent: %w", err)
	} else {
		val.AcRmsCurrent = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AcPowerCurrentPhaseB-Power: %w", err)
	} else {
		val.Power = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeAcPowerCurrentPhaseB(val *AcPowerCurrentPhaseB) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.ConnectionNumber, 8)
	w.writeUnsignedResolution(val.AcRmsCurrent, 16, 0.1)
	w.writeInt32(val.Power, 32)
	return w.Bytes(), w.Err()
}

func encodeAcPowerCurrentPhaseBMsg(v Message) ([]byte, error) {
	val, ok := v.(*AcPowerCurrentPhaseB)
	if !ok {
		return nil, fmt.Errorf("expected *AcPowerCurrentPhaseB, got %T", v)
	}
	return EncodeAcPowerCurrentPhaseB(val)
}

type AcPowerCurrentPhaseC struct {
	Info             MessageInfo `json:"info"`
	Sid              *uint8      `json:"sid"`
	ConnectionNumber *uint8      `json:"connectionNumber"`
	AcRmsCurrent     *float32    `json:"acRmsCurrent"`
	Power            *int32      `json:"power"`
}

func (a *AcPowerCurrentPhaseC) PGNNumber() uint32 { return 127746 }

func DecodeAcPowerCurrentPhaseC(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AcPowerCurrentPhaseC{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AcPowerCurrentPhaseC-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AcPowerCurrentPhaseC-ConnectionNumber: %w", err)
	} else {
		val.ConnectionNumber = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AcPowerCurrentPhaseC-AcRmsCurrent: %w", err)
	} else {
		val.AcRmsCurrent = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AcPowerCurrentPhaseC-Power: %w", err)
	} else {
		val.Power = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeAcPowerCurrentPhaseC(val *AcPowerCurrentPhaseC) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.ConnectionNumber, 8)
	w.writeUnsignedResolution(val.AcRmsCurrent, 16, 0.1)
	w.writeInt32(val.Power, 32)
	return w.Bytes(), w.Err()
}

func encodeAcPowerCurrentPhaseCMsg(v Message) ([]byte, error) {
	val, ok := v.(*AcPowerCurrentPhaseC)
	if !ok {
		return nil, fmt.Errorf("expected *AcPowerCurrentPhaseC, got %T", v)
	}
	return EncodeAcPowerCurrentPhaseC(val)
}

type ConverterStatus struct {
	Info              MessageInfo           `json:"info"`
	Sid               []uint8               `json:"sid"`
	ConnectionNumber  *uint8                `json:"connectionNumber"`
	OperatingState    ConverterStateConst   `json:"operatingState"`
	TemperatureState  GoodWarningErrorConst `json:"temperatureState"`
	OverloadState     GoodWarningErrorConst `json:"overloadState"`
	LowDcVoltageState GoodWarningErrorConst `json:"lowDcVoltageState"`
	RippleState       GoodWarningErrorConst `json:"rippleState"`
}

func (c *ConverterStatus) PGNNumber() uint32 { return 127750 }

func DecodeConverterStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &ConverterStatus{}
	val.Info = Info
	if v, err := stream.readBinaryData(8); err != nil {
		return nil, fmt.Errorf("parse failed for ConverterStatus-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ConverterStatus-ConnectionNumber: %w", err)
	} else {
		val.ConnectionNumber = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for ConverterStatus-OperatingState: %w", err)
	} else {
		val.OperatingState = ConverterStateConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for ConverterStatus-TemperatureState: %w", err)
	} else {
		val.TemperatureState = GoodWarningErrorConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for ConverterStatus-OverloadState: %w", err)
	} else {
		val.OverloadState = GoodWarningErrorConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for ConverterStatus-LowDcVoltageState: %w", err)
	} else {
		val.LowDcVoltageState = GoodWarningErrorConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for ConverterStatus-RippleState: %w", err)
	} else {
		val.RippleState = GoodWarningErrorConst(v)

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

func EncodeConverterStatus(val *ConverterStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeBinaryData(val.Sid, 8)
	w.writeUInt8(val.ConnectionNumber, 8)
	w.writeLookupField(uint64(val.OperatingState), 8)
	w.writeLookupField(uint64(val.TemperatureState), 2)
	w.writeLookupField(uint64(val.OverloadState), 2)
	w.writeLookupField(uint64(val.LowDcVoltageState), 2)
	w.writeLookupField(uint64(val.RippleState), 2)
	w.writeReservedBits(32)
	return w.Bytes(), w.Err()
}

func encodeConverterStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*ConverterStatus)
	if !ok {
		return nil, fmt.Errorf("expected *ConverterStatus, got %T", v)
	}
	return EncodeConverterStatus(val)
}

type DcVoltageCurrent struct {
	Info             MessageInfo `json:"info"`
	Sid              []uint8     `json:"sid"`
	ConnectionNumber *uint8      `json:"connectionNumber"`
	DcVoltage        *float32    `json:"dcVoltage"`
	DcCurrent        *float32    `json:"dcCurrent"`
}

func (d *DcVoltageCurrent) PGNNumber() uint32 { return 127751 }

func DecodeDcVoltageCurrent(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &DcVoltageCurrent{}
	val.Info = Info
	if v, err := stream.readBinaryData(8); err != nil {
		return nil, fmt.Errorf("parse failed for DcVoltageCurrent-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for DcVoltageCurrent-ConnectionNumber: %w", err)
	} else {
		val.ConnectionNumber = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for DcVoltageCurrent-DcVoltage: %w", err)
	} else {
		val.DcVoltage = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(24, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for DcVoltageCurrent-DcCurrent: %w", err)
	} else {
		val.DcCurrent = v

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

func EncodeDcVoltageCurrent(val *DcVoltageCurrent) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeBinaryData(val.Sid, 8)
	w.writeUInt8(val.ConnectionNumber, 8)
	w.writeUnsignedResolution(val.DcVoltage, 16, 0.1)
	w.writeSignedResolution(val.DcCurrent, 24, 0.01)
	w.writeReservedBits(8)
	return w.Bytes(), w.Err()
}

func encodeDcVoltageCurrentMsg(v Message) ([]byte, error) {
	val, ok := v.(*DcVoltageCurrent)
	if !ok {
		return nil, fmt.Errorf("expected *DcVoltageCurrent, got %T", v)
	}
	return EncodeDcVoltageCurrent(val)
}
