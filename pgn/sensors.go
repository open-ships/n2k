package pgn

import (
	"fmt"

	"github.com/open-ships/n2k/units"
)

type TrackedTargetData struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	TargetId *uint8 `json:"targetId"`
	TrackStatus TrackingConst `json:"trackStatus"`
	ReportedTarget YesNoConst `json:"reportedTarget"`
	TargetAcquisition TargetAcquisitionConst `json:"targetAcquisition"`
	BearingReference DirectionReferenceConst `json:"bearingReference"`
	Bearing *float32 `json:"bearing"`
	Distance *units.Distance `json:"distance"`
	Course *float32 `json:"course"`
	Speed *units.Velocity `json:"speed"`
	Cpa *units.Distance `json:"cpa"`
	Tcpa *float32 `json:"tcpa"`
	UtcOfFix *float32 `json:"utcOfFix"`
	Name string `json:"name"`
}

func (x *TrackedTargetData) PGNNumber() uint32  { return 128520 }

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
	w.skipBits(2)
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

type WindlassControlStatus struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	WindlassId *uint8 `json:"windlassId"`
	WindlassDirectionControl WindlassDirectionConst `json:"windlassDirectionControl"`
	AnchorDockingControl OffOnConst `json:"anchorDockingControl"`
	SpeedControlType SpeedTypeConst `json:"speedControlType"`
	SpeedControl []uint8 `json:"speedControl"`
	PowerEnable OffOnConst `json:"powerEnable"`
	MechanicalLock OffOnConst `json:"mechanicalLock"`
	DeckAndAnchorWash OffOnConst `json:"deckAndAnchorWash"`
	AnchorLight OffOnConst `json:"anchorLight"`
	CommandTimeout *float32 `json:"commandTimeout"`
	WindlassControlEvents WindlassControlConst `json:"windlassControlEvents"`
}

func (x *WindlassControlStatus) PGNNumber() uint32  { return 128776 }

func DecodeWindlassControlStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &WindlassControlStatus{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for WindlassControlStatus-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for WindlassControlStatus-WindlassId: %w", err)
	} else {
		val.WindlassId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WindlassControlStatus-WindlassDirectionControl: %w", err)
	} else {
		val.WindlassDirectionControl = WindlassDirectionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WindlassControlStatus-AnchorDockingControl: %w", err)
	} else {
		val.AnchorDockingControl = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WindlassControlStatus-SpeedControlType: %w", err)
	} else {
		val.SpeedControlType = SpeedTypeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readBinaryData(8); err != nil {
		return nil, fmt.Errorf("parse failed for WindlassControlStatus-SpeedControl: %w", err)
	} else {
		val.SpeedControl = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WindlassControlStatus-PowerEnable: %w", err)
	} else {
		val.PowerEnable = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WindlassControlStatus-MechanicalLock: %w", err)
	} else {
		val.MechanicalLock = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WindlassControlStatus-DeckAndAnchorWash: %w", err)
	} else {
		val.DeckAndAnchorWash = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WindlassControlStatus-AnchorLight: %w", err)
	} else {
		val.AnchorLight = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(8, 0.005); err != nil {
		return nil, fmt.Errorf("parse failed for WindlassControlStatus-CommandTimeout: %w", err)
	} else {
		val.CommandTimeout = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for WindlassControlStatus-WindlassControlEvents: %w", err)
	} else {
		val.WindlassControlEvents = WindlassControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(12)
	if stream.isEOF() {
		return val, nil
		}	
	return val, nil
}

func EncodeWindlassControlStatus(val *WindlassControlStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.WindlassId, 8)
	w.writeLookupField(uint64(val.WindlassDirectionControl), 2)
	w.writeLookupField(uint64(val.AnchorDockingControl), 2)
	w.writeLookupField(uint64(val.SpeedControlType), 2)
	w.skipBits(2)
	w.writeBinaryData(val.SpeedControl, 8)
	w.writeLookupField(uint64(val.PowerEnable), 2)
	w.writeLookupField(uint64(val.MechanicalLock), 2)
	w.writeLookupField(uint64(val.DeckAndAnchorWash), 2)
	w.writeLookupField(uint64(val.AnchorLight), 2)
	w.writeUnsignedResolution(val.CommandTimeout, 8, 0.005)
	w.writeLookupField(uint64(val.WindlassControlEvents), 4)
	w.skipBits(12)
	return w.Bytes(), w.Err()
}

func encodeWindlassControlStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*WindlassControlStatus)
	if !ok {
		return nil, fmt.Errorf("expected *WindlassControlStatus, got %T", v)
	}
	return EncodeWindlassControlStatus(val)
}

type AnchorWindlassOperatingStatus struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	WindlassId *uint8 `json:"windlassId"`
	WindlassDirectionControl WindlassDirectionConst `json:"windlassDirectionControl"`
	WindlassMotionStatus WindlassMotionConst `json:"windlassMotionStatus"`
	RodeTypeStatus RodeTypeConst `json:"rodeTypeStatus"`
	RodeCounterValue *units.Distance `json:"rodeCounterValue"`
	WindlassLineSpeed *units.Velocity `json:"windlassLineSpeed"`
	AnchorDockingStatus DockingStatusConst `json:"anchorDockingStatus"`
	WindlassOperatingEvents WindlassOperationConst `json:"windlassOperatingEvents"`
}

func (x *AnchorWindlassOperatingStatus) PGNNumber() uint32  { return 128777 }

func DecodeAnchorWindlassOperatingStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AnchorWindlassOperatingStatus{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AnchorWindlassOperatingStatus-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AnchorWindlassOperatingStatus-WindlassId: %w", err)
	} else {
		val.WindlassId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AnchorWindlassOperatingStatus-WindlassDirectionControl: %w", err)
	} else {
		val.WindlassDirectionControl = WindlassDirectionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AnchorWindlassOperatingStatus-WindlassMotionStatus: %w", err)
	} else {
		val.WindlassMotionStatus = WindlassMotionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AnchorWindlassOperatingStatus-RodeTypeStatus: %w", err)
	} else {
		val.RodeTypeStatus = RodeTypeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AnchorWindlassOperatingStatus-RodeCounterValue: %w", err)
	} else {
		val.RodeCounterValue = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AnchorWindlassOperatingStatus-WindlassLineSpeed: %w", err)
	} else {
		val.WindlassLineSpeed = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AnchorWindlassOperatingStatus-AnchorDockingStatus: %w", err)
	} else {
		val.AnchorDockingStatus = DockingStatusConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AnchorWindlassOperatingStatus-WindlassOperatingEvents: %w", err)
	} else {
		val.WindlassOperatingEvents = WindlassOperationConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeAnchorWindlassOperatingStatus(val *AnchorWindlassOperatingStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.WindlassId, 8)
	w.writeLookupField(uint64(val.WindlassDirectionControl), 2)
	w.writeLookupField(uint64(val.WindlassMotionStatus), 2)
	w.writeLookupField(uint64(val.RodeTypeStatus), 2)
	w.skipBits(2)
	var rodeCounterValueRaw *float32
	if val.RodeCounterValue != nil {
		rodeCounterValueRaw = &val.RodeCounterValue.Value
	}
	w.writeUnsignedResolution(rodeCounterValueRaw, 16, 0.1)
	var windlassLineSpeedRaw *float32
	if val.WindlassLineSpeed != nil {
		windlassLineSpeedRaw = &val.WindlassLineSpeed.Value
	}
	w.writeUnsignedResolution(windlassLineSpeedRaw, 16, 0.01)
	w.writeLookupField(uint64(val.AnchorDockingStatus), 2)
	w.writeLookupField(uint64(val.WindlassOperatingEvents), 6)
	return w.Bytes(), w.Err()
}

func encodeAnchorWindlassOperatingStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*AnchorWindlassOperatingStatus)
	if !ok {
		return nil, fmt.Errorf("expected *AnchorWindlassOperatingStatus, got %T", v)
	}
	return EncodeAnchorWindlassOperatingStatus(val)
}

type AnchorWindlassMonitoringStatus struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	WindlassId *uint8 `json:"windlassId"`
	WindlassMonitoringEvents WindlassMonitoringConst `json:"windlassMonitoringEvents"`
	ControllerVoltage *float32 `json:"controllerVoltage"`
	MotorCurrent *uint8 `json:"motorCurrent"`
	TotalMotorTime *float32 `json:"totalMotorTime"`
}

func (x *AnchorWindlassMonitoringStatus) PGNNumber() uint32  { return 128778 }

func DecodeAnchorWindlassMonitoringStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AnchorWindlassMonitoringStatus{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AnchorWindlassMonitoringStatus-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AnchorWindlassMonitoringStatus-WindlassId: %w", err)
	} else {
		val.WindlassId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AnchorWindlassMonitoringStatus-WindlassMonitoringEvents: %w", err)
	} else {
		val.WindlassMonitoringEvents = WindlassMonitoringConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(8, 0.2); err != nil {
		return nil, fmt.Errorf("parse failed for AnchorWindlassMonitoringStatus-ControllerVoltage: %w", err)
	} else {
		val.ControllerVoltage = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AnchorWindlassMonitoringStatus-MotorCurrent: %w", err)
	} else {
		val.MotorCurrent = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 60); err != nil {
		return nil, fmt.Errorf("parse failed for AnchorWindlassMonitoringStatus-TotalMotorTime: %w", err)
	} else {
		val.TotalMotorTime = v

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

func EncodeAnchorWindlassMonitoringStatus(val *AnchorWindlassMonitoringStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.WindlassId, 8)
	w.writeLookupField(uint64(val.WindlassMonitoringEvents), 8)
	w.writeUnsignedResolution(val.ControllerVoltage, 8, 0.2)
	w.writeUInt8(val.MotorCurrent, 8)
	w.writeUnsignedResolution(val.TotalMotorTime, 16, 60)
	w.skipBits(8)
	return w.Bytes(), w.Err()
}

func encodeAnchorWindlassMonitoringStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*AnchorWindlassMonitoringStatus)
	if !ok {
		return nil, fmt.Errorf("expected *AnchorWindlassMonitoringStatus, got %T", v)
	}
	return EncodeAnchorWindlassMonitoringStatus(val)
}
