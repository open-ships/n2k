package pgn

import (
	"fmt"

	"github.com/open-ships/n2k/units"
)

type ThrusterControlStatus struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	Identifier *uint8 `json:"identifier"`
	DirectionControl ThrusterDirectionControlConst `json:"directionControl"`
	PowerEnabled OffOnConst `json:"powerEnabled"`
	RetractControl ThrusterRetractControlConst `json:"retractControl"`
	SpeedControl *uint8 `json:"speedControl"`
	ControlEvents ThrusterControlEventsConst `json:"controlEvents"`
	CommandTimeout *float32 `json:"commandTimeout"`
	AzimuthControl *float32 `json:"azimuthControl"`
}
func DecodeThrusterControlStatus(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val ThrusterControlStatus
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
type ThrusterInformation struct {
	Info MessageInfo `json:"info"`
	Identifier *uint8 `json:"identifier"`
	MotorType ThrusterMotorTypeConst `json:"motorType"`
	PowerRating *uint16 `json:"powerRating"`
	MaximumTemperatureRating *units.Temperature `json:"maximumTemperatureRating"`
	MaximumRotationalSpeed *float32 `json:"maximumRotationalSpeed"`
}
func DecodeThrusterInformation(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val ThrusterInformation
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
type ThrusterMotorStatus struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	Identifier *uint8 `json:"identifier"`
	MotorEvents ThrusterMotorEventsConst `json:"motorEvents"`
	Current *uint8 `json:"current"`
	Temperature *units.Temperature `json:"temperature"`
	OperatingTime *float32 `json:"operatingTime"`
}
func DecodeThrusterMotorStatus(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val ThrusterMotorStatus
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
type Speed struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	SpeedWaterReferenced *units.Velocity `json:"speedWaterReferenced"`
	SpeedGroundReferenced *units.Velocity `json:"speedGroundReferenced"`
	SpeedWaterReferencedType WaterReferenceConst `json:"speedWaterReferencedType"`
	SpeedDirection *uint8 `json:"speedDirection"`
}
func DecodeSpeed(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val Speed
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Speed-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for Speed-SpeedWaterReferenced: %w", err)
	} else {
		val.SpeedWaterReferenced = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for Speed-SpeedGroundReferenced: %w", err)
	} else {
		val.SpeedGroundReferenced = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for Speed-SpeedWaterReferencedType: %w", err)
	} else {
		val.SpeedWaterReferencedType = WaterReferenceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for Speed-SpeedDirection: %w", err)
	} else {
		val.SpeedDirection = v

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
type WaterDepth struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	Depth *units.Distance `json:"depth"`
	Offset *units.Distance `json:"offset"`
	Range *units.Distance `json:"range"`
}
func DecodeWaterDepth(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val WaterDepth
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for WaterDepth-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for WaterDepth-Depth: %w", err)
	} else {
		val.Depth = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for WaterDepth-Offset: %w", err)
	} else {
		val.Offset = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(8, 10); err != nil {
		return nil, fmt.Errorf("parse failed for WaterDepth-Range: %w", err)
	} else {
		val.Range = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type DistanceLog struct {
	Info MessageInfo `json:"info"`
	Date *uint16 `json:"date"`
	Time *float32 `json:"time"`
	Log *units.Distance `json:"log"`
	TripLog *units.Distance `json:"tripLog"`
}
func DecodeDistanceLog(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val DistanceLog
	val.Info = Info
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for DistanceLog-Date: %w", err)
	} else {
		val.Date = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for DistanceLog-Time: %w", err)
	} else {
		val.Time = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for DistanceLog-Log: %w", err)
	} else {
		val.Log = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for DistanceLog-TripLog: %w", err)
	} else {
		val.TripLog = nullableUnit(units.Meter, v, units.NewDistance)

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
	return w.Bytes(), nil
}
func encodeThrusterControlStatusAny(v any) ([]byte, error) {
	val, ok := v.(*ThrusterControlStatus)
	if !ok {
		return nil, fmt.Errorf("expected *ThrusterControlStatus, got %T", v)
	}
	return EncodeThrusterControlStatus(val)
}

func EncodeThrusterInformation(val *ThrusterInformation) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Identifier, 8)
	w.writeLookupField(uint64(val.MotorType), 4)
	w.skipBits(4)
	w.writeUInt16(val.PowerRating, 16)
	var maximumTemperatureRatingRaw *float32
	if val.MaximumTemperatureRating != nil {
		maximumTemperatureRatingRaw = &val.MaximumTemperatureRating.Value
	}
	w.writeUnsignedResolution(maximumTemperatureRatingRaw, 16, 0.01)
	w.writeUnsignedResolution(val.MaximumRotationalSpeed, 16, 0.25)
	return w.Bytes(), nil
}
func encodeThrusterInformationAny(v any) ([]byte, error) {
	val, ok := v.(*ThrusterInformation)
	if !ok {
		return nil, fmt.Errorf("expected *ThrusterInformation, got %T", v)
	}
	return EncodeThrusterInformation(val)
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
	return w.Bytes(), nil
}
func encodeThrusterMotorStatusAny(v any) ([]byte, error) {
	val, ok := v.(*ThrusterMotorStatus)
	if !ok {
		return nil, fmt.Errorf("expected *ThrusterMotorStatus, got %T", v)
	}
	return EncodeThrusterMotorStatus(val)
}

func EncodeSpeed(val *Speed) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	var speedWaterReferencedRaw *float32
	if val.SpeedWaterReferenced != nil {
		speedWaterReferencedRaw = &val.SpeedWaterReferenced.Value
	}
	w.writeUnsignedResolution(speedWaterReferencedRaw, 16, 0.01)
	var speedGroundReferencedRaw *float32
	if val.SpeedGroundReferenced != nil {
		speedGroundReferencedRaw = &val.SpeedGroundReferenced.Value
	}
	w.writeUnsignedResolution(speedGroundReferencedRaw, 16, 0.01)
	w.writeLookupField(uint64(val.SpeedWaterReferencedType), 8)
	w.writeUInt8(val.SpeedDirection, 4)
	w.skipBits(12)
	return w.Bytes(), nil
}
func encodeSpeedAny(v any) ([]byte, error) {
	val, ok := v.(*Speed)
	if !ok {
		return nil, fmt.Errorf("expected *Speed, got %T", v)
	}
	return EncodeSpeed(val)
}

func EncodeWaterDepth(val *WaterDepth) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	var depthRaw *float32
	if val.Depth != nil {
		depthRaw = &val.Depth.Value
	}
	w.writeUnsignedResolution(depthRaw, 32, 0.01)
	var offsetRaw *float32
	if val.Offset != nil {
		offsetRaw = &val.Offset.Value
	}
	w.writeSignedResolution(offsetRaw, 16, 0.001)
	var rangeRaw *float32
	if val.Range != nil {
		rangeRaw = &val.Range.Value
	}
	w.writeUnsignedResolution(rangeRaw, 8, 10)
	return w.Bytes(), nil
}
func encodeWaterDepthAny(v any) ([]byte, error) {
	val, ok := v.(*WaterDepth)
	if !ok {
		return nil, fmt.Errorf("expected *WaterDepth, got %T", v)
	}
	return EncodeWaterDepth(val)
}

func EncodeDistanceLog(val *DistanceLog) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.Date, 16)
	w.writeUnsignedResolution(val.Time, 32, 0.0001)
	var logRaw *uint32
	if val.Log != nil {
		v := uint32(val.Log.Value)
		logRaw = &v
	}
	w.writeUInt32(logRaw, 32)
	var tripLogRaw *uint32
	if val.TripLog != nil {
		v := uint32(val.TripLog.Value)
		tripLogRaw = &v
	}
	w.writeUInt32(tripLogRaw, 32)
	return w.Bytes(), nil
}
func encodeDistanceLogAny(v any) ([]byte, error) {
	val, ok := v.(*DistanceLog)
	if !ok {
		return nil, fmt.Errorf("expected *DistanceLog, got %T", v)
	}
	return EncodeDistanceLog(val)
}
