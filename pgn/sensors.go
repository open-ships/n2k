package pgn

import (
	"fmt"

	"github.com/open-ships/n2k/units"
)

type TrackedTargetData struct {
	Info              MessageInfo             `json:"info"`
	Sid               *uint8                  `json:"sid"`
	TargetId          *uint8                  `json:"targetId"`
	TrackStatus       TrackingConst           `json:"trackStatus"`
	ReportedTarget    YesNoConst              `json:"reportedTarget"`
	TargetAcquisition TargetAcquisitionConst  `json:"targetAcquisition"`
	BearingReference  DirectionReferenceConst `json:"bearingReference"`
	Bearing           *float32                `json:"bearing"`
	Distance          *units.Distance         `json:"distance"`
	Course            *float32                `json:"course"`
	Speed             *units.Velocity         `json:"speed"`
	Cpa               *units.Distance         `json:"cpa"`
	Tcpa              *float32                `json:"tcpa"`
	UtcOfFix          *float32                `json:"utcOfFix"`
	Name              string                  `json:"name"`
}

func (x *TrackedTargetData) PGNNumber() uint32 { return 128520 }

func DecodeTrackedTargetData(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &TrackedTargetData{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for TrackedTargetData-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for TrackedTargetData-TargetId: %w", err)
	} else {
		val.TargetId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for TrackedTargetData-TrackStatus: %w", err)
	} else {
		val.TrackStatus = TrackingConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for TrackedTargetData-ReportedTarget: %w", err)
	} else {
		val.ReportedTarget = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for TrackedTargetData-TargetAcquisition: %w", err)
	} else {
		val.TargetAcquisition = TargetAcquisitionConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for TrackedTargetData-BearingReference: %w", err)
	} else {
		val.BearingReference = DirectionReferenceConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for TrackedTargetData-Bearing: %w", err)
	} else {
		val.Bearing = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(32, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for TrackedTargetData-Distance: %w", err)
	} else {
		val.Distance = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for TrackedTargetData-Course: %w", err)
	} else {
		val.Course = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for TrackedTargetData-Speed: %w", err)
	} else {
		val.Speed = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(32, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for TrackedTargetData-Cpa: %w", err)
	} else {
		val.Cpa = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(32, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for TrackedTargetData-Tcpa: %w", err)
	} else {
		val.Tcpa = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(32, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for TrackedTargetData-UtcOfFix: %w", err)
	} else {
		val.UtcOfFix = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(1784); err != nil {
		return nil, fmt.Errorf("parse failed for TrackedTargetData-Name: %w", err)
	} else {
		val.Name = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeTrackedTargetData(val *TrackedTargetData) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.TargetId, 8)
	w.writeLookupField(uint64(val.TrackStatus), 2)
	w.writeLookupField(uint64(val.ReportedTarget), 1)
	w.writeLookupField(uint64(val.TargetAcquisition), 1)
	w.writeLookupField(uint64(val.BearingReference), 2)
	w.writeReservedBits(2)
	w.writeUnsignedResolution(val.Bearing, 16, 0.0001)
	var distanceRaw *float32
	if val.Distance != nil {
		distanceRaw = &val.Distance.Value
	}
	w.writeUnsignedResolution(distanceRaw, 32, 0.001)
	w.writeUnsignedResolution(val.Course, 16, 0.0001)
	var speedRaw *float32
	if val.Speed != nil {
		speedRaw = &val.Speed.Value
	}
	w.writeUnsignedResolution(speedRaw, 16, 0.01)
	var cpaRaw *float32
	if val.Cpa != nil {
		cpaRaw = &val.Cpa.Value
	}
	w.writeUnsignedResolution(cpaRaw, 32, 0.01)
	w.writeSignedResolution(val.Tcpa, 32, 0.001)
	w.writeUnsignedResolution(val.UtcOfFix, 32, 0.0001)
	w.writeFixedString(val.Name, 1784)
	return w.Bytes(), w.Err()
}

func encodeTrackedTargetDataMsg(v Message) ([]byte, error) {
	val, ok := v.(*TrackedTargetData)
	if !ok {
		return nil, fmt.Errorf("expected *TrackedTargetData, got %T", v)
	}
	return EncodeTrackedTargetData(val)
}
