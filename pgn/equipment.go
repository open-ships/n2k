package pgn

import (
	"fmt"
	"github.com/open-ships/n2k/units"
)

type FluidLevel struct {
	Info     MessageInfo   `json:"info"`
	Instance *uint8        `json:"instance"`
	Type     TankTypeConst `json:"type"`
	Level    *float32      `json:"level"`
	Capacity *units.Volume `json:"capacity"`
}

func (f *FluidLevel) PGNNumber() uint32 { return 127505 }

func DecodeFluidLevel(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &FluidLevel{}
	val.Info = Info
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for FluidLevel-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for FluidLevel-Type: %w", err)
	} else {
		val.Type = TankTypeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.004); err != nil {
		return nil, fmt.Errorf("parse failed for FluidLevel-Level: %w", err)
	} else {
		val.Level = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(32, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for FluidLevel-Capacity: %w", err)
	} else {
		val.Capacity = nullableUnit(units.Liter, v, units.NewVolume)

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

func EncodeFluidLevel(val *FluidLevel) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Instance, 4)
	w.writeLookupField(uint64(val.Type), 4)
	w.writeSignedResolution(val.Level, 16, 0.004)
	var capacityRaw *float32
	if val.Capacity != nil {
		capacityRaw = &val.Capacity.Value
	}
	w.writeUnsignedResolution(capacityRaw, 32, 0.1)
	w.writeReservedBits(8)
	return w.Bytes(), w.Err()
}

func encodeFluidLevelMsg(v Message) ([]byte, error) {
	val, ok := v.(*FluidLevel)
	if !ok {
		return nil, fmt.Errorf("expected *FluidLevel, got %T", v)
	}
	return EncodeFluidLevel(val)
}

type WindlassControlStatus struct {
	Info                     MessageInfo            `json:"info"`
	Sid                      *uint8                 `json:"sid"`
	WindlassId               *uint8                 `json:"windlassId"`
	WindlassDirectionControl WindlassDirectionConst `json:"windlassDirectionControl"`
	AnchorDockingControl     OffOnConst             `json:"anchorDockingControl"`
	SpeedControlType         SpeedTypeConst         `json:"speedControlType"`
	SpeedControl             []uint8                `json:"speedControl"`
	PowerEnable              OffOnConst             `json:"powerEnable"`
	MechanicalLock           OffOnConst             `json:"mechanicalLock"`
	DeckAndAnchorWash        OffOnConst             `json:"deckAndAnchorWash"`
	AnchorLight              OffOnConst             `json:"anchorLight"`
	CommandTimeout           *float32               `json:"commandTimeout"`
	WindlassControlEvents    WindlassControlConst   `json:"windlassControlEvents"`
}

func (x *WindlassControlStatus) PGNNumber() uint32 { return 128776 }

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
	w.writeReservedBits(2)
	w.writeBinaryData(val.SpeedControl, 8)
	w.writeLookupField(uint64(val.PowerEnable), 2)
	w.writeLookupField(uint64(val.MechanicalLock), 2)
	w.writeLookupField(uint64(val.DeckAndAnchorWash), 2)
	w.writeLookupField(uint64(val.AnchorLight), 2)
	w.writeUnsignedResolution(val.CommandTimeout, 8, 0.005)
	w.writeLookupField(uint64(val.WindlassControlEvents), 4)
	w.writeReservedBits(12)
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
	Info                     MessageInfo            `json:"info"`
	Sid                      *uint8                 `json:"sid"`
	WindlassId               *uint8                 `json:"windlassId"`
	WindlassDirectionControl WindlassDirectionConst `json:"windlassDirectionControl"`
	WindlassMotionStatus     WindlassMotionConst    `json:"windlassMotionStatus"`
	RodeTypeStatus           RodeTypeConst          `json:"rodeTypeStatus"`
	RodeCounterValue         *units.Distance        `json:"rodeCounterValue"`
	WindlassLineSpeed        *units.Velocity        `json:"windlassLineSpeed"`
	AnchorDockingStatus      DockingStatusConst     `json:"anchorDockingStatus"`
	WindlassOperatingEvents  WindlassOperationConst `json:"windlassOperatingEvents"`
}

func (x *AnchorWindlassOperatingStatus) PGNNumber() uint32 { return 128777 }

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
	w.writeReservedBits(2)
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
	Info                     MessageInfo             `json:"info"`
	Sid                      *uint8                  `json:"sid"`
	WindlassId               *uint8                  `json:"windlassId"`
	WindlassMonitoringEvents WindlassMonitoringConst `json:"windlassMonitoringEvents"`
	ControllerVoltage        *float32                `json:"controllerVoltage"`
	MotorCurrent             *uint8                  `json:"motorCurrent"`
	TotalMotorTime           *float32                `json:"totalMotorTime"`
}

func (x *AnchorWindlassMonitoringStatus) PGNNumber() uint32 { return 128778 }

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
	w.writeReservedBits(8)
	return w.Bytes(), w.Err()
}

func encodeAnchorWindlassMonitoringStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*AnchorWindlassMonitoringStatus)
	if !ok {
		return nil, fmt.Errorf("expected *AnchorWindlassMonitoringStatus, got %T", v)
	}
	return EncodeAnchorWindlassMonitoringStatus(val)
}

type WatermakerInputSettingAndStatus struct {
	Info                       MessageInfo          `json:"info"`
	WatermakerOperatingState   WatermakerStateConst `json:"watermakerOperatingState"`
	ProductionStartStop        YesNoConst           `json:"productionStartStop"`
	RinseStartStop             YesNoConst           `json:"rinseStartStop"`
	LowPressurePumpStatus      YesNoConst           `json:"lowPressurePumpStatus"`
	HighPressurePumpStatus     YesNoConst           `json:"highPressurePumpStatus"`
	EmergencyStop              YesNoConst           `json:"emergencyStop"`
	ProductSolenoidValveStatus OkWarningConst       `json:"productSolenoidValveStatus"`
	FlushModeStatus            YesNoConst           `json:"flushModeStatus"`
	SalinityStatus             OkWarningConst       `json:"salinityStatus"`
	SensorStatus               OkWarningConst       `json:"sensorStatus"`
	OilChangeIndicatorStatus   OkWarningConst       `json:"oilChangeIndicatorStatus"`
	FilterStatus               OkWarningConst       `json:"filterStatus"`
	SystemStatus               OkWarningConst       `json:"systemStatus"`
	Salinity                   *uint16              `json:"salinity"`
	ProductWaterTemperature    *units.Temperature   `json:"productWaterTemperature"`
	PreFilterPressure          *units.Pressure      `json:"preFilterPressure"`
	PostFilterPressure         *units.Pressure      `json:"postFilterPressure"`
	FeedPressure               *units.Pressure      `json:"feedPressure"`
	SystemHighPressure         *units.Pressure      `json:"systemHighPressure"`
	ProductWaterFlow           *units.Flow          `json:"productWaterFlow"`
	BrineWaterFlow             *units.Flow          `json:"brineWaterFlow"`
	RunTime                    *uint32              `json:"runTime"`
}

func (x *WatermakerInputSettingAndStatus) PGNNumber() uint32 { return 130567 }

func DecodeWatermakerInputSettingAndStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &WatermakerInputSettingAndStatus{}
	val.Info = Info
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-WatermakerOperatingState: %w", err)
	} else {
		val.WatermakerOperatingState = WatermakerStateConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-ProductionStartStop: %w", err)
	} else {
		val.ProductionStartStop = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-RinseStartStop: %w", err)
	} else {
		val.RinseStartStop = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-LowPressurePumpStatus: %w", err)
	} else {
		val.LowPressurePumpStatus = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-HighPressurePumpStatus: %w", err)
	} else {
		val.HighPressurePumpStatus = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-EmergencyStop: %w", err)
	} else {
		val.EmergencyStop = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-ProductSolenoidValveStatus: %w", err)
	} else {
		val.ProductSolenoidValveStatus = OkWarningConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-FlushModeStatus: %w", err)
	} else {
		val.FlushModeStatus = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-SalinityStatus: %w", err)
	} else {
		val.SalinityStatus = OkWarningConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-SensorStatus: %w", err)
	} else {
		val.SensorStatus = OkWarningConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-OilChangeIndicatorStatus: %w", err)
	} else {
		val.OilChangeIndicatorStatus = OkWarningConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-FilterStatus: %w", err)
	} else {
		val.FilterStatus = OkWarningConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-SystemStatus: %w", err)
	} else {
		val.SystemStatus = OkWarningConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-Salinity: %w", err)
	} else {
		val.Salinity = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-ProductWaterTemperature: %w", err)
	} else {
		val.ProductWaterTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 100); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-PreFilterPressure: %w", err)
	} else {
		val.PreFilterPressure = nullableUnit(units.Pa, v, units.NewPressure)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 100); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-PostFilterPressure: %w", err)
	} else {
		val.PostFilterPressure = nullableUnit(units.Pa, v, units.NewPressure)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 1000); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-FeedPressure: %w", err)
	} else {
		val.FeedPressure = nullableUnit(units.Pa, v, units.NewPressure)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 1000); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-SystemHighPressure: %w", err)
	} else {
		val.SystemHighPressure = nullableUnit(units.Pa, v, units.NewPressure)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-ProductWaterFlow: %w", err)
	} else {
		val.ProductWaterFlow = nullableUnit(units.LitersPerHour, v, units.NewFlow)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-BrineWaterFlow: %w", err)
	} else {
		val.BrineWaterFlow = nullableUnit(units.LitersPerHour, v, units.NewFlow)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for WatermakerInputSettingAndStatus-RunTime: %w", err)
	} else {
		val.RunTime = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeWatermakerInputSettingAndStatus(val *WatermakerInputSettingAndStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.WatermakerOperatingState), 6)
	w.writeLookupField(uint64(val.ProductionStartStop), 2)
	w.writeLookupField(uint64(val.RinseStartStop), 2)
	w.writeLookupField(uint64(val.LowPressurePumpStatus), 2)
	w.writeLookupField(uint64(val.HighPressurePumpStatus), 2)
	w.writeLookupField(uint64(val.EmergencyStop), 2)
	w.writeLookupField(uint64(val.ProductSolenoidValveStatus), 2)
	w.writeLookupField(uint64(val.FlushModeStatus), 2)
	w.writeLookupField(uint64(val.SalinityStatus), 2)
	w.writeLookupField(uint64(val.SensorStatus), 2)
	w.writeLookupField(uint64(val.OilChangeIndicatorStatus), 2)
	w.writeLookupField(uint64(val.FilterStatus), 2)
	w.writeLookupField(uint64(val.SystemStatus), 2)
	w.writeReservedBits(2)
	w.writeUInt16(val.Salinity, 16)
	var productWaterTemperatureRaw *float32
	if val.ProductWaterTemperature != nil {
		productWaterTemperatureRaw = &val.ProductWaterTemperature.Value
	}
	w.writeUnsignedResolution(productWaterTemperatureRaw, 16, 0.01)
	var preFilterPressureRaw *float32
	if val.PreFilterPressure != nil {
		preFilterPressureRaw = &val.PreFilterPressure.Value
	}
	w.writeUnsignedResolution(preFilterPressureRaw, 16, 100)
	var postFilterPressureRaw *float32
	if val.PostFilterPressure != nil {
		postFilterPressureRaw = &val.PostFilterPressure.Value
	}
	w.writeUnsignedResolution(postFilterPressureRaw, 16, 100)
	var feedPressureRaw *float32
	if val.FeedPressure != nil {
		feedPressureRaw = &val.FeedPressure.Value
	}
	w.writeSignedResolution(feedPressureRaw, 16, 1000)
	var systemHighPressureRaw *float32
	if val.SystemHighPressure != nil {
		systemHighPressureRaw = &val.SystemHighPressure.Value
	}
	w.writeUnsignedResolution(systemHighPressureRaw, 16, 1000)
	var productWaterFlowRaw *float32
	if val.ProductWaterFlow != nil {
		productWaterFlowRaw = &val.ProductWaterFlow.Value
	}
	w.writeSignedResolution(productWaterFlowRaw, 16, 0.1)
	var brineWaterFlowRaw *float32
	if val.BrineWaterFlow != nil {
		brineWaterFlowRaw = &val.BrineWaterFlow.Value
	}
	w.writeSignedResolution(brineWaterFlowRaw, 16, 0.1)
	w.writeUInt32(val.RunTime, 32)
	return w.Bytes(), w.Err()
}

func encodeWatermakerInputSettingAndStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*WatermakerInputSettingAndStatus)
	if !ok {
		return nil, fmt.Errorf("expected *WatermakerInputSettingAndStatus, got %T", v)
	}
	return EncodeWatermakerInputSettingAndStatus(val)
}
