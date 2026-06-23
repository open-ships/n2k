package pgn

import (
	"fmt"

	"github.com/open-ships/n2k/units"
)

type WindData struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	WindSpeed *units.Velocity `json:"windSpeed"`
	WindAngle *float32 `json:"windAngle"`
	Reference WindReferenceConst `json:"reference"`
}

func (x *WindData) PGNNumber() uint32  { return 130306 }

func DecodeWindData(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &WindData{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for WindData-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for WindData-WindSpeed: %w", err)
	} else {
		val.WindSpeed = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for WindData-WindAngle: %w", err)
	} else {
		val.WindAngle = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for WindData-Reference: %w", err)
	} else {
		val.Reference = WindReferenceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(21)
	if stream.isEOF() {
		return val, nil
		}	
	return val, nil
}

func EncodeWindData(val *WindData) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	var windSpeedRaw *float32
	if val.WindSpeed != nil {
		windSpeedRaw = &val.WindSpeed.Value
	}
	w.writeUnsignedResolution(windSpeedRaw, 16, 0.01)
	w.writeUnsignedResolution(val.WindAngle, 16, 0.0001)
	w.writeLookupField(uint64(val.Reference), 3)
	w.skipBits(21)
	return w.Bytes(), w.Err()
}

func encodeWindDataMsg(v Message) ([]byte, error) {
	val, ok := v.(*WindData)
	if !ok {
		return nil, fmt.Errorf("expected *WindData, got %T", v)
	}
	return EncodeWindData(val)
}

type EnvironmentalParametersObsolete struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	WaterTemperature *units.Temperature `json:"waterTemperature"`
	OutsideAmbientAirTemperature *units.Temperature `json:"outsideAmbientAirTemperature"`
	AtmosphericPressure *units.Pressure `json:"atmosphericPressure"`
}

func (x *EnvironmentalParametersObsolete) PGNNumber() uint32  { return 130310 }

func DecodeEnvironmentalParametersObsolete(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &EnvironmentalParametersObsolete{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for EnvironmentalParametersObsolete-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for EnvironmentalParametersObsolete-WaterTemperature: %w", err)
	} else {
		val.WaterTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for EnvironmentalParametersObsolete-OutsideAmbientAirTemperature: %w", err)
	} else {
		val.OutsideAmbientAirTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 100); err != nil {
		return nil, fmt.Errorf("parse failed for EnvironmentalParametersObsolete-AtmosphericPressure: %w", err)
	} else {
		val.AtmosphericPressure = nullableUnit(units.Pa, v, units.NewPressure)

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

func EncodeEnvironmentalParametersObsolete(val *EnvironmentalParametersObsolete) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	var waterTemperatureRaw *float32
	if val.WaterTemperature != nil {
		waterTemperatureRaw = &val.WaterTemperature.Value
	}
	w.writeUnsignedResolution(waterTemperatureRaw, 16, 0.01)
	var outsideAmbientAirTemperatureRaw *float32
	if val.OutsideAmbientAirTemperature != nil {
		outsideAmbientAirTemperatureRaw = &val.OutsideAmbientAirTemperature.Value
	}
	w.writeUnsignedResolution(outsideAmbientAirTemperatureRaw, 16, 0.01)
	var atmosphericPressureRaw *float32
	if val.AtmosphericPressure != nil {
		atmosphericPressureRaw = &val.AtmosphericPressure.Value
	}
	w.writeUnsignedResolution(atmosphericPressureRaw, 16, 100)
	w.skipBits(8)
	return w.Bytes(), w.Err()
}

func encodeEnvironmentalParametersObsoleteMsg(v Message) ([]byte, error) {
	val, ok := v.(*EnvironmentalParametersObsolete)
	if !ok {
		return nil, fmt.Errorf("expected *EnvironmentalParametersObsolete, got %T", v)
	}
	return EncodeEnvironmentalParametersObsolete(val)
}

type EnvironmentalParameters struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	TemperatureSource TemperatureSourceConst `json:"temperatureSource"`
	HumiditySource HumiditySourceConst `json:"humiditySource"`
	Temperature *units.Temperature `json:"temperature"`
	Humidity *float32 `json:"humidity"`
	AtmosphericPressure *units.Pressure `json:"atmosphericPressure"`
}

func (x *EnvironmentalParameters) PGNNumber() uint32  { return 130311 }

func DecodeEnvironmentalParameters(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &EnvironmentalParameters{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for EnvironmentalParameters-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for EnvironmentalParameters-TemperatureSource: %w", err)
	} else {
		val.TemperatureSource = TemperatureSourceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for EnvironmentalParameters-HumiditySource: %w", err)
	} else {
		val.HumiditySource = HumiditySourceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for EnvironmentalParameters-Temperature: %w", err)
	} else {
		val.Temperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.004); err != nil {
		return nil, fmt.Errorf("parse failed for EnvironmentalParameters-Humidity: %w", err)
	} else {
		val.Humidity = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 100); err != nil {
		return nil, fmt.Errorf("parse failed for EnvironmentalParameters-AtmosphericPressure: %w", err)
	} else {
		val.AtmosphericPressure = nullableUnit(units.Pa, v, units.NewPressure)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeEnvironmentalParameters(val *EnvironmentalParameters) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeLookupField(uint64(val.TemperatureSource), 6)
	w.writeLookupField(uint64(val.HumiditySource), 2)
	var temperatureRaw *float32
	if val.Temperature != nil {
		temperatureRaw = &val.Temperature.Value
	}
	w.writeUnsignedResolution(temperatureRaw, 16, 0.01)
	w.writeSignedResolution(val.Humidity, 16, 0.004)
	var atmosphericPressureRaw *float32
	if val.AtmosphericPressure != nil {
		atmosphericPressureRaw = &val.AtmosphericPressure.Value
	}
	w.writeUnsignedResolution(atmosphericPressureRaw, 16, 100)
	return w.Bytes(), w.Err()
}

func encodeEnvironmentalParametersMsg(v Message) ([]byte, error) {
	val, ok := v.(*EnvironmentalParameters)
	if !ok {
		return nil, fmt.Errorf("expected *EnvironmentalParameters, got %T", v)
	}
	return EncodeEnvironmentalParameters(val)
}

type Temperature struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	Instance *uint8 `json:"instance"`
	Source TemperatureSourceConst `json:"source"`
	ActualTemperature *units.Temperature `json:"actualTemperature"`
	SetTemperature *units.Temperature `json:"setTemperature"`
}

func (x *Temperature) PGNNumber() uint32  { return 130312 }

func DecodeTemperature(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Temperature{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Temperature-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Temperature-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for Temperature-Source: %w", err)
	} else {
		val.Source = TemperatureSourceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for Temperature-ActualTemperature: %w", err)
	} else {
		val.ActualTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for Temperature-SetTemperature: %w", err)
	} else {
		val.SetTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

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

func EncodeTemperature(val *Temperature) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.Instance, 8)
	w.writeLookupField(uint64(val.Source), 8)
	var actualTemperatureRaw *float32
	if val.ActualTemperature != nil {
		actualTemperatureRaw = &val.ActualTemperature.Value
	}
	w.writeUnsignedResolution(actualTemperatureRaw, 16, 0.01)
	var setTemperatureRaw *float32
	if val.SetTemperature != nil {
		setTemperatureRaw = &val.SetTemperature.Value
	}
	w.writeUnsignedResolution(setTemperatureRaw, 16, 0.01)
	w.skipBits(8)
	return w.Bytes(), w.Err()
}

func encodeTemperatureMsg(v Message) ([]byte, error) {
	val, ok := v.(*Temperature)
	if !ok {
		return nil, fmt.Errorf("expected *Temperature, got %T", v)
	}
	return EncodeTemperature(val)
}

type Humidity struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	Instance *uint8 `json:"instance"`
	Source HumiditySourceConst `json:"source"`
	ActualHumidity *float32 `json:"actualHumidity"`
	SetHumidity *float32 `json:"setHumidity"`
}

func (x *Humidity) PGNNumber() uint32  { return 130313 }

func DecodeHumidity(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Humidity{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Humidity-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Humidity-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for Humidity-Source: %w", err)
	} else {
		val.Source = HumiditySourceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.004); err != nil {
		return nil, fmt.Errorf("parse failed for Humidity-ActualHumidity: %w", err)
	} else {
		val.ActualHumidity = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.004); err != nil {
		return nil, fmt.Errorf("parse failed for Humidity-SetHumidity: %w", err)
	} else {
		val.SetHumidity = v

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

func EncodeHumidity(val *Humidity) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.Instance, 8)
	w.writeLookupField(uint64(val.Source), 8)
	w.writeSignedResolution(val.ActualHumidity, 16, 0.004)
	w.writeSignedResolution(val.SetHumidity, 16, 0.004)
	w.skipBits(8)
	return w.Bytes(), w.Err()
}

func encodeHumidityMsg(v Message) ([]byte, error) {
	val, ok := v.(*Humidity)
	if !ok {
		return nil, fmt.Errorf("expected *Humidity, got %T", v)
	}
	return EncodeHumidity(val)
}

type ActualPressure struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	Instance *uint8 `json:"instance"`
	Source PressureSourceConst `json:"source"`
	Pressure *units.Pressure `json:"pressure"`
}

func (x *ActualPressure) PGNNumber() uint32  { return 130314 }

func DecodeActualPressure(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &ActualPressure{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ActualPressure-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ActualPressure-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for ActualPressure-Source: %w", err)
	} else {
		val.Source = PressureSourceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(32, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for ActualPressure-Pressure: %w", err)
	} else {
		val.Pressure = nullableUnit(units.Pa, v, units.NewPressure)

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

func EncodeActualPressure(val *ActualPressure) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.Instance, 8)
	w.writeLookupField(uint64(val.Source), 8)
	var pressureRaw *float32
	if val.Pressure != nil {
		pressureRaw = &val.Pressure.Value
	}
	w.writeSignedResolution(pressureRaw, 32, 0.1)
	w.skipBits(8)
	return w.Bytes(), w.Err()
}

func encodeActualPressureMsg(v Message) ([]byte, error) {
	val, ok := v.(*ActualPressure)
	if !ok {
		return nil, fmt.Errorf("expected *ActualPressure, got %T", v)
	}
	return EncodeActualPressure(val)
}

type SetPressure struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	Instance *uint8 `json:"instance"`
	Source PressureSourceConst `json:"source"`
	Pressure *units.Pressure `json:"pressure"`
}

func (x *SetPressure) PGNNumber() uint32  { return 130315 }

func DecodeSetPressure(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SetPressure{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SetPressure-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SetPressure-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SetPressure-Source: %w", err)
	} else {
		val.Source = PressureSourceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for SetPressure-Pressure: %w", err)
	} else {
		val.Pressure = nullableUnit(units.Pa, v, units.NewPressure)

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

func EncodeSetPressure(val *SetPressure) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.Instance, 8)
	w.writeLookupField(uint64(val.Source), 8)
	var pressureRaw *float32
	if val.Pressure != nil {
		pressureRaw = &val.Pressure.Value
	}
	w.writeUnsignedResolution(pressureRaw, 32, 0.1)
	w.skipBits(8)
	return w.Bytes(), w.Err()
}

func encodeSetPressureMsg(v Message) ([]byte, error) {
	val, ok := v.(*SetPressure)
	if !ok {
		return nil, fmt.Errorf("expected *SetPressure, got %T", v)
	}
	return EncodeSetPressure(val)
}

type TemperatureExtendedRange struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	Instance *uint8 `json:"instance"`
	Source TemperatureSourceConst `json:"source"`
	Temperature *units.Temperature `json:"temperature"`
	SetTemperature *units.Temperature `json:"setTemperature"`
}

func (x *TemperatureExtendedRange) PGNNumber() uint32  { return 130316 }

func DecodeTemperatureExtendedRange(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &TemperatureExtendedRange{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for TemperatureExtendedRange-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for TemperatureExtendedRange-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for TemperatureExtendedRange-Source: %w", err)
	} else {
		val.Source = TemperatureSourceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(24, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for TemperatureExtendedRange-Temperature: %w", err)
	} else {
		val.Temperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for TemperatureExtendedRange-SetTemperature: %w", err)
	} else {
		val.SetTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeTemperatureExtendedRange(val *TemperatureExtendedRange) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.Instance, 8)
	w.writeLookupField(uint64(val.Source), 8)
	var temperatureRaw *float32
	if val.Temperature != nil {
		temperatureRaw = &val.Temperature.Value
	}
	w.writeUnsignedResolution(temperatureRaw, 24, 0.001)
	var setTemperatureRaw *float32
	if val.SetTemperature != nil {
		setTemperatureRaw = &val.SetTemperature.Value
	}
	w.writeUnsignedResolution(setTemperatureRaw, 16, 0.1)
	return w.Bytes(), w.Err()
}

func encodeTemperatureExtendedRangeMsg(v Message) ([]byte, error) {
	val, ok := v.(*TemperatureExtendedRange)
	if !ok {
		return nil, fmt.Errorf("expected *TemperatureExtendedRange, got %T", v)
	}
	return EncodeTemperatureExtendedRange(val)
}

type TideStationData struct {
	Info MessageInfo `json:"info"`
	Mode ResidualModeConst `json:"mode"`
	TideTendency TideConst `json:"tideTendency"`
	MeasurementDate *uint16 `json:"measurementDate"`
	MeasurementTime *float32 `json:"measurementTime"`
	StationLatitude *float64 `json:"stationLatitude"`
	StationLongitude *float64 `json:"stationLongitude"`
	TideLevel *units.Distance `json:"tideLevel"`
	TideLevelStandardDeviation *units.Distance `json:"tideLevelStandardDeviation"`
	StationId string `json:"stationId"`
	StationName string `json:"stationName"`
}

func (x *TideStationData) PGNNumber() uint32  { return 130320 }

func DecodeTideStationData(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &TideStationData{}
	val.Info = Info
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for TideStationData-Mode: %w", err)
	} else {
		val.Mode = ResidualModeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for TideStationData-TideTendency: %w", err)
	} else {
		val.TideTendency = TideConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for TideStationData-MeasurementDate: %w", err)
	} else {
		val.MeasurementDate = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for TideStationData-MeasurementTime: %w", err)
	} else {
		val.MeasurementTime = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for TideStationData-StationLatitude: %w", err)
	} else {
		val.StationLatitude = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for TideStationData-StationLongitude: %w", err)
	} else {
		val.StationLongitude = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for TideStationData-TideLevel: %w", err)
	} else {
		val.TideLevel = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for TideStationData-TideLevelStandardDeviation: %w", err)
	} else {
		val.TideLevelStandardDeviation = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLengthAndControl(); err != nil {
		return nil, fmt.Errorf("parse failed for TideStationData-StationId: %w", err)
	} else {
		val.StationId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLengthAndControl(); err != nil {
		return nil, fmt.Errorf("parse failed for TideStationData-StationName: %w", err)
	} else {
		val.StationName = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeTideStationData(val *TideStationData) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.Mode), 4)
	w.writeLookupField(uint64(val.TideTendency), 2)
	w.skipBits(2)
	w.writeUInt16(val.MeasurementDate, 16)
	w.writeUnsignedResolution(val.MeasurementTime, 32, 0.0001)
	w.writeSignedResolution64Override(val.StationLatitude, 32, 1e-07)
	w.writeSignedResolution64Override(val.StationLongitude, 32, 1e-07)
	var tideLevelRaw *float32
	if val.TideLevel != nil {
		tideLevelRaw = &val.TideLevel.Value
	}
	w.writeSignedResolution(tideLevelRaw, 16, 0.001)
	var tideLevelStandardDeviationRaw *float32
	if val.TideLevelStandardDeviation != nil {
		tideLevelStandardDeviationRaw = &val.TideLevelStandardDeviation.Value
	}
	w.writeUnsignedResolution(tideLevelStandardDeviationRaw, 16, 0.01)
	w.writeStringWithLengthAndControl(val.StationId)
	w.writeStringWithLengthAndControl(val.StationName)
	return w.Bytes(), w.Err()
}

func encodeTideStationDataMsg(v Message) ([]byte, error) {
	val, ok := v.(*TideStationData)
	if !ok {
		return nil, fmt.Errorf("expected *TideStationData, got %T", v)
	}
	return EncodeTideStationData(val)
}

type SalinityStationData struct {
	Info MessageInfo `json:"info"`
	Mode ResidualModeConst `json:"mode"`
	MeasurementDate *uint16 `json:"measurementDate"`
	MeasurementTime *float32 `json:"measurementTime"`
	StationLatitude *float64 `json:"stationLatitude"`
	StationLongitude *float64 `json:"stationLongitude"`
	Salinity *float32 `json:"salinity"`
	WaterTemperature *units.Temperature `json:"waterTemperature"`
	StationId string `json:"stationId"`
	StationName string `json:"stationName"`
}

func (x *SalinityStationData) PGNNumber() uint32  { return 130321 }

func DecodeSalinityStationData(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SalinityStationData{}
	val.Info = Info
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for SalinityStationData-Mode: %w", err)
	} else {
		val.Mode = ResidualModeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SalinityStationData-MeasurementDate: %w", err)
	} else {
		val.MeasurementDate = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SalinityStationData-MeasurementTime: %w", err)
	} else {
		val.MeasurementTime = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for SalinityStationData-StationLatitude: %w", err)
	} else {
		val.StationLatitude = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for SalinityStationData-StationLongitude: %w", err)
	} else {
		val.StationLongitude = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFloat32(); err != nil {
		return nil, fmt.Errorf("parse failed for SalinityStationData-Salinity: %w", err)
	} else {
		val.Salinity = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for SalinityStationData-WaterTemperature: %w", err)
	} else {
		val.WaterTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLengthAndControl(); err != nil {
		return nil, fmt.Errorf("parse failed for SalinityStationData-StationId: %w", err)
	} else {
		val.StationId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLengthAndControl(); err != nil {
		return nil, fmt.Errorf("parse failed for SalinityStationData-StationName: %w", err)
	} else {
		val.StationName = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeSalinityStationData(val *SalinityStationData) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.Mode), 4)
	w.skipBits(4)
	w.writeUInt16(val.MeasurementDate, 16)
	w.writeUnsignedResolution(val.MeasurementTime, 32, 0.0001)
	w.writeSignedResolution64Override(val.StationLatitude, 32, 1e-07)
	w.writeSignedResolution64Override(val.StationLongitude, 32, 1e-07)
	w.writeFloat32(val.Salinity)
	var waterTemperatureRaw *float32
	if val.WaterTemperature != nil {
		waterTemperatureRaw = &val.WaterTemperature.Value
	}
	w.writeUnsignedResolution(waterTemperatureRaw, 16, 0.01)
	w.writeStringWithLengthAndControl(val.StationId)
	w.writeStringWithLengthAndControl(val.StationName)
	return w.Bytes(), w.Err()
}

func encodeSalinityStationDataMsg(v Message) ([]byte, error) {
	val, ok := v.(*SalinityStationData)
	if !ok {
		return nil, fmt.Errorf("expected *SalinityStationData, got %T", v)
	}
	return EncodeSalinityStationData(val)
}
