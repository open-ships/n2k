package pgn

import (
	"fmt"

	"github.com/open-ships/n2k/units"
)

type ThrusterControlStatus struct {
	Info             MessageInfo                   `json:"info"`
	Sid              *uint8                        `json:"sid"`
	Identifier       *uint8                        `json:"identifier"`
	DirectionControl ThrusterDirectionControlConst `json:"directionControl"`
	PowerEnabled     OffOnConst                    `json:"powerEnabled"`
	RetractControl   ThrusterRetractControlConst   `json:"retractControl"`
	SpeedControl     *uint8                        `json:"speedControl"`
	ControlEvents    ThrusterControlEventsConst    `json:"controlEvents"`
	CommandTimeout   *float32                      `json:"commandTimeout"`
	AzimuthControl   *float32                      `json:"azimuthControl"`
}

func (x *ThrusterControlStatus) PGNNumber() uint32 { return 128006 }

func DecodeThrusterControlStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &ThrusterControlStatus{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterControlStatus-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterControlStatus-Identifier: %w", err)
	} else {
		val.Identifier = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterControlStatus-DirectionControl: %w", err)
	} else {
		val.DirectionControl = ThrusterDirectionControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterControlStatus-PowerEnabled: %w", err)
	} else {
		val.PowerEnabled = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterControlStatus-RetractControl: %w", err)
	} else {
		val.RetractControl = ThrusterRetractControlConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterControlStatus-SpeedControl: %w", err)
	} else {
		val.SpeedControl = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterControlStatus-ControlEvents: %w", err)
	} else {
		val.ControlEvents = ThrusterControlEventsConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(8, 0.005); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterControlStatus-CommandTimeout: %w", err)
	} else {
		val.CommandTimeout = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterControlStatus-AzimuthControl: %w", err)
	} else {
		val.AzimuthControl = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeThrusterControlStatus(val *ThrusterControlStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.Identifier, 8)
	w.writeLookupField(uint64(val.DirectionControl), 4)
	w.writeLookupField(uint64(val.PowerEnabled), 2)
	w.writeLookupField(uint64(val.RetractControl), 2)
	w.writeUInt8(val.SpeedControl, 8)
	w.writeLookupField(uint64(val.ControlEvents), 8)
	w.writeUnsignedResolution(val.CommandTimeout, 8, 0.005)
	w.writeUnsignedResolution(val.AzimuthControl, 16, 0.0001)
	return w.Bytes(), w.Err()
}

func encodeThrusterControlStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*ThrusterControlStatus)
	if !ok {
		return nil, fmt.Errorf("expected *ThrusterControlStatus, got %T", v)
	}
	return EncodeThrusterControlStatus(val)
}

type ThrusterInformation struct {
	Info                     MessageInfo            `json:"info"`
	Identifier               *uint8                 `json:"identifier"`
	MotorType                ThrusterMotorTypeConst `json:"motorType"`
	PowerRating              *uint16                `json:"powerRating"`
	MaximumTemperatureRating *units.Temperature     `json:"maximumTemperatureRating"`
	MaximumRotationalSpeed   *float32               `json:"maximumRotationalSpeed"`
}

func (x *ThrusterInformation) PGNNumber() uint32 { return 128007 }

func DecodeThrusterInformation(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &ThrusterInformation{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterInformation-Identifier: %w", err)
	} else {
		val.Identifier = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterInformation-MotorType: %w", err)
	} else {
		val.MotorType = ThrusterMotorTypeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterInformation-PowerRating: %w", err)
	} else {
		val.PowerRating = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterInformation-MaximumTemperatureRating: %w", err)
	} else {
		val.MaximumTemperatureRating = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.25); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterInformation-MaximumRotationalSpeed: %w", err)
	} else {
		val.MaximumRotationalSpeed = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeThrusterInformation(val *ThrusterInformation) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Identifier, 8)
	w.writeLookupField(uint64(val.MotorType), 4)
	w.writeReservedBits(4)
	w.writeUInt16(val.PowerRating, 16)
	var maximumTemperatureRatingRaw *float32
	if val.MaximumTemperatureRating != nil {
		maximumTemperatureRatingRaw = &val.MaximumTemperatureRating.Value
	}
	w.writeUnsignedResolution(maximumTemperatureRatingRaw, 16, 0.01)
	w.writeUnsignedResolution(val.MaximumRotationalSpeed, 16, 0.25)
	return w.Bytes(), w.Err()
}

func encodeThrusterInformationMsg(v Message) ([]byte, error) {
	val, ok := v.(*ThrusterInformation)
	if !ok {
		return nil, fmt.Errorf("expected *ThrusterInformation, got %T", v)
	}
	return EncodeThrusterInformation(val)
}

type ThrusterMotorStatus struct {
	Info          MessageInfo              `json:"info"`
	Sid           *uint8                   `json:"sid"`
	Identifier    *uint8                   `json:"identifier"`
	MotorEvents   ThrusterMotorEventsConst `json:"motorEvents"`
	Current       *uint8                   `json:"current"`
	Temperature   *units.Temperature       `json:"temperature"`
	OperatingTime *float32                 `json:"operatingTime"`
}

func (x *ThrusterMotorStatus) PGNNumber() uint32 { return 128008 }

func DecodeThrusterMotorStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &ThrusterMotorStatus{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterMotorStatus-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterMotorStatus-Identifier: %w", err)
	} else {
		val.Identifier = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterMotorStatus-MotorEvents: %w", err)
	} else {
		val.MotorEvents = ThrusterMotorEventsConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterMotorStatus-Current: %w", err)
	} else {
		val.Current = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterMotorStatus-Temperature: %w", err)
	} else {
		val.Temperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 60); err != nil {
		return nil, fmt.Errorf("parse failed for ThrusterMotorStatus-OperatingTime: %w", err)
	} else {
		val.OperatingTime = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeThrusterMotorStatus(val *ThrusterMotorStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.Identifier, 8)
	w.writeLookupField(uint64(val.MotorEvents), 8)
	w.writeUInt8(val.Current, 8)
	var temperatureRaw *float32
	if val.Temperature != nil {
		temperatureRaw = &val.Temperature.Value
	}
	w.writeUnsignedResolution(temperatureRaw, 16, 0.01)
	w.writeUnsignedResolution(val.OperatingTime, 16, 60)
	return w.Bytes(), w.Err()
}

func encodeThrusterMotorStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*ThrusterMotorStatus)
	if !ok {
		return nil, fmt.Errorf("expected *ThrusterMotorStatus, got %T", v)
	}
	return EncodeThrusterMotorStatus(val)
}
