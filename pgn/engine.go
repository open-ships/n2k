package pgn

import (
	"fmt"

	"github.com/open-ships/n2k/units"
)

type Rudder struct {
	Info MessageInfo `json:"info"`
	Instance *uint8 `json:"instance"`
	DirectionOrder DirectionRudderConst `json:"directionOrder"`
	AngleOrder *float32 `json:"angleOrder"`
	Position *float32 `json:"position"`
}

func (x *Rudder) PGNNumber() uint32  { return 127245 }

func DecodeRudder(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Rudder{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Rudder-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for Rudder-DirectionOrder: %w", err)
	} else {
		val.DirectionOrder = DirectionRudderConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(5)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for Rudder-AngleOrder: %w", err)
	} else {
		val.AngleOrder = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for Rudder-Position: %w", err)
	} else {
		val.Position = v

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

func EncodeRudder(val *Rudder) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Instance, 8)
	w.writeLookupField(uint64(val.DirectionOrder), 3)
	w.skipBits(5)
	w.writeSignedResolution(val.AngleOrder, 16, 0.0001)
	w.writeSignedResolution(val.Position, 16, 0.0001)
	w.skipBits(16)
	return w.Bytes(), w.Err()
}

func encodeRudderMsg(v Message) ([]byte, error) {
	val, ok := v.(*Rudder)
	if !ok {
		return nil, fmt.Errorf("expected *Rudder, got %T", v)
	}
	return EncodeRudder(val)
}

type EngineParametersRapidUpdate struct {
	Info MessageInfo `json:"info"`
	Instance EngineInstanceConst `json:"instance"`
	Speed *float32 `json:"speed"`
	BoostPressure *units.Pressure `json:"boostPressure"`
	TiltTrim *int8 `json:"tiltTrim"`
}

func (x *EngineParametersRapidUpdate) PGNNumber() uint32  { return 127488 }

func DecodeEngineParametersRapidUpdate(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &EngineParametersRapidUpdate{}
	val.Info = Info
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersRapidUpdate-Instance: %w", err)
	} else {
		val.Instance = EngineInstanceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.25); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersRapidUpdate-Speed: %w", err)
	} else {
		val.Speed = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 100); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersRapidUpdate-BoostPressure: %w", err)
	} else {
		val.BoostPressure = nullableUnit(units.Pa, v, units.NewPressure)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersRapidUpdate-TiltTrim: %w", err)
	} else {
		val.TiltTrim = v

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

func EncodeEngineParametersRapidUpdate(val *EngineParametersRapidUpdate) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.Instance), 8)
	w.writeUnsignedResolution(val.Speed, 16, 0.25)
	var boostPressureRaw *float32
	if val.BoostPressure != nil {
		boostPressureRaw = &val.BoostPressure.Value
	}
	w.writeUnsignedResolution(boostPressureRaw, 16, 100)
	w.writeInt8(val.TiltTrim, 8)
	w.skipBits(16)
	return w.Bytes(), w.Err()
}

func encodeEngineParametersRapidUpdateMsg(v Message) ([]byte, error) {
	val, ok := v.(*EngineParametersRapidUpdate)
	if !ok {
		return nil, fmt.Errorf("expected *EngineParametersRapidUpdate, got %T", v)
	}
	return EncodeEngineParametersRapidUpdate(val)
}

type EngineParametersDynamic struct {
	Info MessageInfo `json:"info"`
	Instance EngineInstanceConst `json:"instance"`
	OilPressure *units.Pressure `json:"oilPressure"`
	OilTemperature *units.Temperature `json:"oilTemperature"`
	Temperature *units.Temperature `json:"temperature"`
	AlternatorPotential *float32 `json:"alternatorPotential"`
	FuelRate *units.Flow `json:"fuelRate"`
	TotalEngineHours *uint32 `json:"totalEngineHours"`
	CoolantPressure *units.Pressure `json:"coolantPressure"`
	FuelPressure *units.Pressure `json:"fuelPressure"`
	DiscreteStatus1 EngineStatus1Const `json:"discreteStatus1"`
	DiscreteStatus2 EngineStatus2Const `json:"discreteStatus2"`
	EngineLoad *int8 `json:"engineLoad"`
	EngineTorque *int8 `json:"engineTorque"`
}

func (x *EngineParametersDynamic) PGNNumber() uint32  { return 127489 }

func DecodeEngineParametersDynamic(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &EngineParametersDynamic{}
	val.Info = Info
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersDynamic-Instance: %w", err)
	} else {
		val.Instance = EngineInstanceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 100); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersDynamic-OilPressure: %w", err)
	} else {
		val.OilPressure = nullableUnit(units.Pa, v, units.NewPressure)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersDynamic-OilTemperature: %w", err)
	} else {
		val.OilTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersDynamic-Temperature: %w", err)
	} else {
		val.Temperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersDynamic-AlternatorPotential: %w", err)
	} else {
		val.AlternatorPotential = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersDynamic-FuelRate: %w", err)
	} else {
		val.FuelRate = nullableUnit(units.LitersPerHour, v, units.NewFlow)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersDynamic-TotalEngineHours: %w", err)
	} else {
		val.TotalEngineHours = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 100); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersDynamic-CoolantPressure: %w", err)
	} else {
		val.CoolantPressure = nullableUnit(units.Pa, v, units.NewPressure)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 1000); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersDynamic-FuelPressure: %w", err)
	} else {
		val.FuelPressure = nullableUnit(units.Pa, v, units.NewPressure)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(16); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersDynamic-DiscreteStatus1: %w", err)
	} else {
		val.DiscreteStatus1 = EngineStatus1Const(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(16); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersDynamic-DiscreteStatus2: %w", err)
	} else {
		val.DiscreteStatus2 = EngineStatus2Const(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersDynamic-EngineLoad: %w", err)
	} else {
		val.EngineLoad = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersDynamic-EngineTorque: %w", err)
	} else {
		val.EngineTorque = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeEngineParametersDynamic(val *EngineParametersDynamic) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.Instance), 8)
	var oilPressureRaw *float32
	if val.OilPressure != nil {
		oilPressureRaw = &val.OilPressure.Value
	}
	w.writeUnsignedResolution(oilPressureRaw, 16, 100)
	var oilTemperatureRaw *float32
	if val.OilTemperature != nil {
		oilTemperatureRaw = &val.OilTemperature.Value
	}
	w.writeUnsignedResolution(oilTemperatureRaw, 16, 0.1)
	var temperatureRaw *float32
	if val.Temperature != nil {
		temperatureRaw = &val.Temperature.Value
	}
	w.writeUnsignedResolution(temperatureRaw, 16, 0.01)
	w.writeSignedResolution(val.AlternatorPotential, 16, 0.01)
	var fuelRateRaw *float32
	if val.FuelRate != nil {
		fuelRateRaw = &val.FuelRate.Value
	}
	w.writeSignedResolution(fuelRateRaw, 16, 0.1)
	w.writeUInt32(val.TotalEngineHours, 32)
	var coolantPressureRaw *float32
	if val.CoolantPressure != nil {
		coolantPressureRaw = &val.CoolantPressure.Value
	}
	w.writeUnsignedResolution(coolantPressureRaw, 16, 100)
	var fuelPressureRaw *float32
	if val.FuelPressure != nil {
		fuelPressureRaw = &val.FuelPressure.Value
	}
	w.writeUnsignedResolution(fuelPressureRaw, 16, 1000)
	w.skipBits(8)
	w.writeLookupField(uint64(val.DiscreteStatus1), 16)
	w.writeLookupField(uint64(val.DiscreteStatus2), 16)
	w.writeInt8(val.EngineLoad, 8)
	w.writeInt8(val.EngineTorque, 8)
	return w.Bytes(), w.Err()
}

func encodeEngineParametersDynamicMsg(v Message) ([]byte, error) {
	val, ok := v.(*EngineParametersDynamic)
	if !ok {
		return nil, fmt.Errorf("expected *EngineParametersDynamic, got %T", v)
	}
	return EncodeEngineParametersDynamic(val)
}

type TransmissionParametersDynamic struct {
	Info MessageInfo `json:"info"`
	Instance EngineInstanceConst `json:"instance"`
	TransmissionGear GearStatusConst `json:"transmissionGear"`
	OilPressure *units.Pressure `json:"oilPressure"`
	OilTemperature *units.Temperature `json:"oilTemperature"`
	DiscreteStatus1 *uint8 `json:"discreteStatus1"`
}

func (x *TransmissionParametersDynamic) PGNNumber() uint32  { return 127493 }

func DecodeTransmissionParametersDynamic(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &TransmissionParametersDynamic{}
	val.Info = Info
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for TransmissionParametersDynamic-Instance: %w", err)
	} else {
		val.Instance = EngineInstanceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for TransmissionParametersDynamic-TransmissionGear: %w", err)
	} else {
		val.TransmissionGear = GearStatusConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUnsignedResolution(16, 100); err != nil {
		return nil, fmt.Errorf("parse failed for TransmissionParametersDynamic-OilPressure: %w", err)
	} else {
		val.OilPressure = nullableUnit(units.Pa, v, units.NewPressure)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for TransmissionParametersDynamic-OilTemperature: %w", err)
	} else {
		val.OilTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for TransmissionParametersDynamic-DiscreteStatus1: %w", err)
	} else {
		val.DiscreteStatus1 = v

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

func EncodeTransmissionParametersDynamic(val *TransmissionParametersDynamic) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.Instance), 8)
	w.writeLookupField(uint64(val.TransmissionGear), 2)
	w.skipBits(6)
	var oilPressureRaw *float32
	if val.OilPressure != nil {
		oilPressureRaw = &val.OilPressure.Value
	}
	w.writeUnsignedResolution(oilPressureRaw, 16, 100)
	var oilTemperatureRaw *float32
	if val.OilTemperature != nil {
		oilTemperatureRaw = &val.OilTemperature.Value
	}
	w.writeUnsignedResolution(oilTemperatureRaw, 16, 0.1)
	w.writeUInt8(val.DiscreteStatus1, 8)
	w.skipBits(8)
	return w.Bytes(), w.Err()
}

func encodeTransmissionParametersDynamicMsg(v Message) ([]byte, error) {
	val, ok := v.(*TransmissionParametersDynamic)
	if !ok {
		return nil, fmt.Errorf("expected *TransmissionParametersDynamic, got %T", v)
	}
	return EncodeTransmissionParametersDynamic(val)
}

type TripParametersVessel struct {
	Info MessageInfo `json:"info"`
	TimeToEmpty *float32 `json:"timeToEmpty"`
	DistanceToEmpty *units.Distance `json:"distanceToEmpty"`
	EstimatedFuelRemaining *units.Volume `json:"estimatedFuelRemaining"`
	TripRunTime *float32 `json:"tripRunTime"`
}

func (x *TripParametersVessel) PGNNumber() uint32  { return 127496 }

func DecodeTripParametersVessel(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &TripParametersVessel{}
	val.Info = Info
	if v, err := stream.readUnsignedResolution(32, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for TripParametersVessel-TimeToEmpty: %w", err)
	} else {
		val.TimeToEmpty = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for TripParametersVessel-DistanceToEmpty: %w", err)
	} else {
		val.DistanceToEmpty = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for TripParametersVessel-EstimatedFuelRemaining: %w", err)
	} else {
		val.EstimatedFuelRemaining = nullableUnit(units.Liter, v, units.NewVolume)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for TripParametersVessel-TripRunTime: %w", err)
	} else {
		val.TripRunTime = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeTripParametersVessel(val *TripParametersVessel) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUnsignedResolution(val.TimeToEmpty, 32, 0.001)
	var distanceToEmptyRaw *float32
	if val.DistanceToEmpty != nil {
		distanceToEmptyRaw = &val.DistanceToEmpty.Value
	}
	w.writeUnsignedResolution(distanceToEmptyRaw, 32, 0.01)
	var estimatedFuelRemainingRaw *uint16
	if val.EstimatedFuelRemaining != nil {
		v := uint16(val.EstimatedFuelRemaining.Value)
		estimatedFuelRemainingRaw = &v
	}
	w.writeUInt16(estimatedFuelRemainingRaw, 16)
	w.writeUnsignedResolution(val.TripRunTime, 32, 0.001)
	return w.Bytes(), w.Err()
}

func encodeTripParametersVesselMsg(v Message) ([]byte, error) {
	val, ok := v.(*TripParametersVessel)
	if !ok {
		return nil, fmt.Errorf("expected *TripParametersVessel, got %T", v)
	}
	return EncodeTripParametersVessel(val)
}

type TripParametersEngine struct {
	Info MessageInfo `json:"info"`
	Instance EngineInstanceConst `json:"instance"`
	TripFuelUsed *units.Volume `json:"tripFuelUsed"`
	FuelRateAverage *units.Flow `json:"fuelRateAverage"`
	FuelRateEconomy *units.Flow `json:"fuelRateEconomy"`
	InstantaneousFuelEconomy *units.Flow `json:"instantaneousFuelEconomy"`
}

func (x *TripParametersEngine) PGNNumber() uint32  { return 127497 }

func DecodeTripParametersEngine(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &TripParametersEngine{}
	val.Info = Info
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for TripParametersEngine-Instance: %w", err)
	} else {
		val.Instance = EngineInstanceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for TripParametersEngine-TripFuelUsed: %w", err)
	} else {
		val.TripFuelUsed = nullableUnit(units.Liter, v, units.NewVolume)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for TripParametersEngine-FuelRateAverage: %w", err)
	} else {
		val.FuelRateAverage = nullableUnit(units.LitersPerHour, v, units.NewFlow)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for TripParametersEngine-FuelRateEconomy: %w", err)
	} else {
		val.FuelRateEconomy = nullableUnit(units.LitersPerHour, v, units.NewFlow)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for TripParametersEngine-InstantaneousFuelEconomy: %w", err)
	} else {
		val.InstantaneousFuelEconomy = nullableUnit(units.LitersPerHour, v, units.NewFlow)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeTripParametersEngine(val *TripParametersEngine) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.Instance), 8)
	var tripFuelUsedRaw *uint16
	if val.TripFuelUsed != nil {
		v := uint16(val.TripFuelUsed.Value)
		tripFuelUsedRaw = &v
	}
	w.writeUInt16(tripFuelUsedRaw, 16)
	var fuelRateAverageRaw *float32
	if val.FuelRateAverage != nil {
		fuelRateAverageRaw = &val.FuelRateAverage.Value
	}
	w.writeSignedResolution(fuelRateAverageRaw, 16, 0.1)
	var fuelRateEconomyRaw *float32
	if val.FuelRateEconomy != nil {
		fuelRateEconomyRaw = &val.FuelRateEconomy.Value
	}
	w.writeSignedResolution(fuelRateEconomyRaw, 16, 0.1)
	var instantaneousFuelEconomyRaw *float32
	if val.InstantaneousFuelEconomy != nil {
		instantaneousFuelEconomyRaw = &val.InstantaneousFuelEconomy.Value
	}
	w.writeSignedResolution(instantaneousFuelEconomyRaw, 16, 0.1)
	return w.Bytes(), w.Err()
}

func encodeTripParametersEngineMsg(v Message) ([]byte, error) {
	val, ok := v.(*TripParametersEngine)
	if !ok {
		return nil, fmt.Errorf("expected *TripParametersEngine, got %T", v)
	}
	return EncodeTripParametersEngine(val)
}

type EngineParametersStatic struct {
	Info MessageInfo `json:"info"`
	Instance EngineInstanceConst `json:"instance"`
	RatedEngineSpeed *float32 `json:"ratedEngineSpeed"`
	Vin string `json:"vin"`
	SoftwareId string `json:"softwareId"`
}

func (x *EngineParametersStatic) PGNNumber() uint32  { return 127498 }

func DecodeEngineParametersStatic(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &EngineParametersStatic{}
	val.Info = Info
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersStatic-Instance: %w", err)
	} else {
		val.Instance = EngineInstanceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.25); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersStatic-RatedEngineSpeed: %w", err)
	} else {
		val.RatedEngineSpeed = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(136); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersStatic-Vin: %w", err)
	} else {
		val.Vin = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for EngineParametersStatic-SoftwareId: %w", err)
	} else {
		val.SoftwareId = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeEngineParametersStatic(val *EngineParametersStatic) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.Instance), 8)
	w.writeUnsignedResolution(val.RatedEngineSpeed, 16, 0.25)
	w.writeFixedString(val.Vin, 136)
	w.writeFixedString(val.SoftwareId, 256)
	return w.Bytes(), w.Err()
}

func encodeEngineParametersStaticMsg(v Message) ([]byte, error) {
	val, ok := v.(*EngineParametersStatic)
	if !ok {
		return nil, fmt.Errorf("expected *EngineParametersStatic, got %T", v)
	}
	return EncodeEngineParametersStatic(val)
}
