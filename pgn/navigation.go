package pgn

import (
	"fmt"

	"github.com/open-ships/n2k/units"
)

type HeadingTrackControl struct {
	Info                     MessageInfo             `json:"info"`
	RudderLimitExceeded      YesNoConst              `json:"rudderLimitExceeded"`
	OffHeadingLimitExceeded  YesNoConst              `json:"offHeadingLimitExceeded"`
	OffTrackLimitExceeded    YesNoConst              `json:"offTrackLimitExceeded"`
	Override                 YesNoConst              `json:"override"`
	SteeringMode             SteeringModeConst       `json:"steeringMode"`
	TurnMode                 TurnModeConst           `json:"turnMode"`
	HeadingReference         DirectionReferenceConst `json:"headingReference"`
	CommandedRudderDirection DirectionRudderConst    `json:"commandedRudderDirection"`
	CommandedRudderAngle     *float32                `json:"commandedRudderAngle"`
	HeadingToSteerCourse     *float32                `json:"headingToSteerCourse"`
	Track                    *float32                `json:"track"`
	RudderLimit              *float32                `json:"rudderLimit"`
	OffHeadingLimit          *float32                `json:"offHeadingLimit"`
	RadiusOfTurnOrder        *float32                `json:"radiusOfTurnOrder"`
	RateOfTurnOrder          *float32                `json:"rateOfTurnOrder"`
	OffTrackLimit            *units.Distance         `json:"offTrackLimit"`
	VesselHeading            *float32                `json:"vesselHeading"`
}

func (h *HeadingTrackControl) PGNNumber() uint32 { return 127237 }
func DecodeHeadingTrackControl(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &HeadingTrackControl{}
	val.Info = Info
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-RudderLimitExceeded: %w", err)
	} else {
		val.RudderLimitExceeded = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-OffHeadingLimitExceeded: %w", err)
	} else {
		val.OffHeadingLimitExceeded = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-OffTrackLimitExceeded: %w", err)
	} else {
		val.OffTrackLimitExceeded = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-Override: %w", err)
	} else {
		val.Override = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-SteeringMode: %w", err)
	} else {
		val.SteeringMode = SteeringModeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-TurnMode: %w", err)
	} else {
		val.TurnMode = TurnModeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-HeadingReference: %w", err)
	} else {
		val.HeadingReference = DirectionReferenceConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(5)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-CommandedRudderDirection: %w", err)
	} else {
		val.CommandedRudderDirection = DirectionRudderConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-CommandedRudderAngle: %w", err)
	} else {
		val.CommandedRudderAngle = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-HeadingToSteerCourse: %w", err)
	} else {
		val.HeadingToSteerCourse = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-Track: %w", err)
	} else {
		val.Track = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-RudderLimit: %w", err)
	} else {
		val.RudderLimit = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-OffHeadingLimit: %w", err)
	} else {
		val.OffHeadingLimit = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-RadiusOfTurnOrder: %w", err)
	} else {
		val.RadiusOfTurnOrder = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 3.125e-05); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-RateOfTurnOrder: %w", err)
	} else {
		val.RateOfTurnOrder = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-OffTrackLimit: %w", err)
	} else {
		val.OffTrackLimit = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for HeadingTrackControl-VesselHeading: %w", err)
	} else {
		val.VesselHeading = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeHeadingTrackControl(val *HeadingTrackControl) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.RudderLimitExceeded), 2)
	w.writeLookupField(uint64(val.OffHeadingLimitExceeded), 2)
	w.writeLookupField(uint64(val.OffTrackLimitExceeded), 2)
	w.writeLookupField(uint64(val.Override), 2)
	w.writeLookupField(uint64(val.SteeringMode), 3)
	w.writeLookupField(uint64(val.TurnMode), 3)
	w.writeLookupField(uint64(val.HeadingReference), 2)
	w.skipBits(5)
	w.writeLookupField(uint64(val.CommandedRudderDirection), 3)
	w.writeSignedResolution(val.CommandedRudderAngle, 16, 0.0001)
	w.writeUnsignedResolution(val.HeadingToSteerCourse, 16, 0.0001)
	w.writeUnsignedResolution(val.Track, 16, 0.0001)
	w.writeUnsignedResolution(val.RudderLimit, 16, 0.0001)
	w.writeUnsignedResolution(val.OffHeadingLimit, 16, 0.0001)
	w.writeSignedResolution(val.RadiusOfTurnOrder, 16, 0.0001)
	w.writeSignedResolution(val.RateOfTurnOrder, 16, 3.125e-05)
	var offTrackLimitRaw *int16
	if val.OffTrackLimit != nil {
		v := int16(val.OffTrackLimit.Value)
		offTrackLimitRaw = &v
	}
	w.writeInt16(offTrackLimitRaw, 16)
	w.writeUnsignedResolution(val.VesselHeading, 16, 0.0001)
	return w.Bytes(), w.Err()
}
func encodeHeadingTrackControlMsg(v Message) ([]byte, error) {
	val, ok := v.(*HeadingTrackControl)
	if !ok {
		return nil, fmt.Errorf("expected *HeadingTrackControl, got %T", v)
	}
	return EncodeHeadingTrackControl(val)
}

type Rudder struct {
	Info           MessageInfo          `json:"info"`
	Instance       *uint8               `json:"instance"`
	DirectionOrder DirectionRudderConst `json:"directionOrder"`
	AngleOrder     *float32             `json:"angleOrder"`
	Position       *float32             `json:"position"`
}

func (x *Rudder) PGNNumber() uint32 { return 127245 }

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

type VesselHeading struct {
	Info      MessageInfo             `json:"info"`
	Sid       *uint8                  `json:"sid"`
	Heading   *float32                `json:"heading"`
	Deviation *float32                `json:"deviation"`
	Variation *float32                `json:"variation"`
	Reference DirectionReferenceConst `json:"reference"`
}

func (v *VesselHeading) PGNNumber() uint32 { return 127250 }
func DecodeVesselHeading(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &VesselHeading{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for VesselHeading-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for VesselHeading-Heading: %w", err)
	} else {
		val.Heading = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for VesselHeading-Deviation: %w", err)
	} else {
		val.Deviation = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for VesselHeading-Variation: %w", err)
	} else {
		val.Variation = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for VesselHeading-Reference: %w", err)
	} else {
		val.Reference = DirectionReferenceConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeVesselHeading(val *VesselHeading) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUnsignedResolution(val.Heading, 16, 0.0001)
	w.writeSignedResolution(val.Deviation, 16, 0.0001)
	w.writeSignedResolution(val.Variation, 16, 0.0001)
	w.writeLookupField(uint64(val.Reference), 2)
	w.skipBits(6)
	return w.Bytes(), w.Err()
}
func encodeVesselHeadingMsg(v Message) ([]byte, error) {
	val, ok := v.(*VesselHeading)
	if !ok {
		return nil, fmt.Errorf("expected *VesselHeading, got %T", v)
	}
	return EncodeVesselHeading(val)
}

type RateOfTurn struct {
	Info MessageInfo `json:"info"`
	Sid  *uint8      `json:"sid"`
	Rate *float64    `json:"rate"`
}

func (r *RateOfTurn) PGNNumber() uint32 { return 127251 }
func DecodeRateOfTurn(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &RateOfTurn{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for RateOfTurn-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution64Override(32, 3.125e-08); err != nil {
		return nil, fmt.Errorf("parse failed for RateOfTurn-Rate: %w", err)
	} else {
		val.Rate = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(24)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeRateOfTurn(val *RateOfTurn) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeSignedResolution64Override(val.Rate, 32, 3.125e-08)
	w.skipBits(24)
	return w.Bytes(), w.Err()
}
func encodeRateOfTurnMsg(v Message) ([]byte, error) {
	val, ok := v.(*RateOfTurn)
	if !ok {
		return nil, fmt.Errorf("expected *RateOfTurn, got %T", v)
	}
	return EncodeRateOfTurn(val)
}

type Heave struct {
	Info  MessageInfo     `json:"info"`
	Sid   *uint8          `json:"sid"`
	Heave *units.Distance `json:"heave"`
}

func (h *Heave) PGNNumber() uint32 { return 127252 }
func DecodeHeave(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Heave{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Heave-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for Heave-Heave: %w", err)
	} else {
		val.Heave = nullableUnit(units.Meter, v, units.NewDistance)

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

func EncodeHeave(val *Heave) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	var heaveRaw *float32
	if val.Heave != nil {
		heaveRaw = &val.Heave.Value
	}
	w.writeSignedResolution(heaveRaw, 16, 0.01)
	w.skipBits(40)
	return w.Bytes(), w.Err()
}
func encodeHeaveMsg(v Message) ([]byte, error) {
	val, ok := v.(*Heave)
	if !ok {
		return nil, fmt.Errorf("expected *Heave, got %T", v)
	}
	return EncodeHeave(val)
}

type Attitude struct {
	Info  MessageInfo `json:"info"`
	Sid   *uint8      `json:"sid"`
	Yaw   *float32    `json:"yaw"`
	Pitch *float32    `json:"pitch"`
	Roll  *float32    `json:"roll"`
}

func (a *Attitude) PGNNumber() uint32 { return 127257 }
func DecodeAttitude(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Attitude{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Attitude-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for Attitude-Yaw: %w", err)
	} else {
		val.Yaw = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for Attitude-Pitch: %w", err)
	} else {
		val.Pitch = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for Attitude-Roll: %w", err)
	} else {
		val.Roll = v

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

func EncodeAttitude(val *Attitude) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeSignedResolution(val.Yaw, 16, 0.0001)
	w.writeSignedResolution(val.Pitch, 16, 0.0001)
	w.writeSignedResolution(val.Roll, 16, 0.0001)
	w.skipBits(8)
	return w.Bytes(), w.Err()
}
func encodeAttitudeMsg(v Message) ([]byte, error) {
	val, ok := v.(*Attitude)
	if !ok {
		return nil, fmt.Errorf("expected *Attitude, got %T", v)
	}
	return EncodeAttitude(val)
}

type MagneticVariation struct {
	Info         MessageInfo            `json:"info"`
	Sid          *uint8                 `json:"sid"`
	Source       MagneticVariationConst `json:"source"`
	AgeOfService *uint16                `json:"ageOfService"`
	Variation    *float32               `json:"variation"`
}

func (m *MagneticVariation) PGNNumber() uint32 { return 127258 }
func DecodeMagneticVariation(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &MagneticVariation{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MagneticVariation-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for MagneticVariation-Source: %w", err)
	} else {
		val.Source = MagneticVariationConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for MagneticVariation-AgeOfService: %w", err)
	} else {
		val.AgeOfService = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for MagneticVariation-Variation: %w", err)
	} else {
		val.Variation = v

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

func EncodeMagneticVariation(val *MagneticVariation) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeLookupField(uint64(val.Source), 4)
	w.skipBits(4)
	w.writeUInt16(val.AgeOfService, 16)
	w.writeSignedResolution(val.Variation, 16, 0.0001)
	w.skipBits(16)
	return w.Bytes(), w.Err()
}
func encodeMagneticVariationMsg(v Message) ([]byte, error) {
	val, ok := v.(*MagneticVariation)
	if !ok {
		return nil, fmt.Errorf("expected *MagneticVariation, got %T", v)
	}
	return EncodeMagneticVariation(val)
}

type LeewayAngle struct {
	Info        MessageInfo `json:"info"`
	Sid         *uint8      `json:"sid"`
	LeewayAngle *float32    `json:"leewayAngle"`
}

func (x *LeewayAngle) PGNNumber() uint32 { return 128000 }

func DecodeLeewayAngle(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &LeewayAngle{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for LeewayAngle-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for LeewayAngle-LeewayAngle: %w", err)
	} else {
		val.LeewayAngle = v

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

func EncodeLeewayAngle(val *LeewayAngle) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeSignedResolution(val.LeewayAngle, 16, 0.0001)
	w.skipBits(40)
	return w.Bytes(), w.Err()
}

func encodeLeewayAngleMsg(v Message) ([]byte, error) {
	val, ok := v.(*LeewayAngle)
	if !ok {
		return nil, fmt.Errorf("expected *LeewayAngle, got %T", v)
	}
	return EncodeLeewayAngle(val)
}

type Speed struct {
	Info                     MessageInfo         `json:"info"`
	Sid                      *uint8              `json:"sid"`
	SpeedWaterReferenced     *units.Velocity     `json:"speedWaterReferenced"`
	SpeedGroundReferenced    *units.Velocity     `json:"speedGroundReferenced"`
	SpeedWaterReferencedType WaterReferenceConst `json:"speedWaterReferencedType"`
	SpeedDirection           *uint8              `json:"speedDirection"`
}

func (x *Speed) PGNNumber() uint32 { return 128259 }

func DecodeSpeed(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Speed{}
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
	return w.Bytes(), w.Err()
}

func encodeSpeedMsg(v Message) ([]byte, error) {
	val, ok := v.(*Speed)
	if !ok {
		return nil, fmt.Errorf("expected *Speed, got %T", v)
	}
	return EncodeSpeed(val)
}

type WaterDepth struct {
	Info   MessageInfo     `json:"info"`
	Sid    *uint8          `json:"sid"`
	Depth  *units.Distance `json:"depth"`
	Offset *units.Distance `json:"offset"`
	Range  *units.Distance `json:"range"`
}

func (x *WaterDepth) PGNNumber() uint32 { return 128267 }

func DecodeWaterDepth(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &WaterDepth{}
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
	return w.Bytes(), w.Err()
}

func encodeWaterDepthMsg(v Message) ([]byte, error) {
	val, ok := v.(*WaterDepth)
	if !ok {
		return nil, fmt.Errorf("expected *WaterDepth, got %T", v)
	}
	return EncodeWaterDepth(val)
}

type DistanceLog struct {
	Info    MessageInfo     `json:"info"`
	Date    *uint16         `json:"date"`
	Time    *float32        `json:"time"`
	Log     *units.Distance `json:"log"`
	TripLog *units.Distance `json:"tripLog"`
}

func (x *DistanceLog) PGNNumber() uint32 { return 128275 }

func DecodeDistanceLog(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &DistanceLog{}
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
	return w.Bytes(), w.Err()
}

func encodeDistanceLogMsg(v Message) ([]byte, error) {
	val, ok := v.(*DistanceLog)
	if !ok {
		return nil, fmt.Errorf("expected *DistanceLog, got %T", v)
	}
	return EncodeDistanceLog(val)
}

type PositionRapidUpdate struct {
	Info      MessageInfo `json:"info"`
	Latitude  *float64    `json:"latitude"`
	Longitude *float64    `json:"longitude"`
}

func (p *PositionRapidUpdate) PGNNumber() uint32 { return 129025 }
func DecodePositionRapidUpdate(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &PositionRapidUpdate{}
	val.Info = Info
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for PositionRapidUpdate-Latitude: %w", err)
	} else {
		val.Latitude = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for PositionRapidUpdate-Longitude: %w", err)
	} else {
		val.Longitude = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodePositionRapidUpdate(val *PositionRapidUpdate) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeSignedResolution64Override(val.Latitude, 32, 1e-07)
	w.writeSignedResolution64Override(val.Longitude, 32, 1e-07)
	return w.Bytes(), w.Err()
}
func encodePositionRapidUpdateMsg(v Message) ([]byte, error) {
	val, ok := v.(*PositionRapidUpdate)
	if !ok {
		return nil, fmt.Errorf("expected *PositionRapidUpdate, got %T", v)
	}
	return EncodePositionRapidUpdate(val)
}

type CogSogRapidUpdate struct {
	Info         MessageInfo             `json:"info"`
	Sid          *uint8                  `json:"sid"`
	CogReference DirectionReferenceConst `json:"cogReference"`
	Cog          *float32                `json:"cog"`
	Sog          *units.Velocity         `json:"sog"`
}

func (c *CogSogRapidUpdate) PGNNumber() uint32 { return 129026 }
func DecodeCogSogRapidUpdate(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &CogSogRapidUpdate{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for CogSogRapidUpdate-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for CogSogRapidUpdate-CogReference: %w", err)
	} else {
		val.CogReference = DirectionReferenceConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for CogSogRapidUpdate-Cog: %w", err)
	} else {
		val.Cog = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for CogSogRapidUpdate-Sog: %w", err)
	} else {
		val.Sog = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

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

func EncodeCogSogRapidUpdate(val *CogSogRapidUpdate) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeLookupField(uint64(val.CogReference), 2)
	w.skipBits(6)
	w.writeUnsignedResolution(val.Cog, 16, 0.0001)
	var sogRaw *float32
	if val.Sog != nil {
		sogRaw = &val.Sog.Value
	}
	w.writeUnsignedResolution(sogRaw, 16, 0.01)
	w.skipBits(16)
	return w.Bytes(), w.Err()
}
func encodeCogSogRapidUpdateMsg(v Message) ([]byte, error) {
	val, ok := v.(*CogSogRapidUpdate)
	if !ok {
		return nil, fmt.Errorf("expected *CogSogRapidUpdate, got %T", v)
	}
	return EncodeCogSogRapidUpdate(val)
}

type PositionDeltaRapidUpdate struct {
	Info           MessageInfo `json:"info"`
	Sid            *uint8      `json:"sid"`
	TimeDelta      *uint16     `json:"timeDelta"`
	LatitudeDelta  *int16      `json:"latitudeDelta"`
	LongitudeDelta *int16      `json:"longitudeDelta"`
}

func (p *PositionDeltaRapidUpdate) PGNNumber() uint32 { return 129027 }
func DecodePositionDeltaRapidUpdate(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &PositionDeltaRapidUpdate{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for PositionDeltaRapidUpdate-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for PositionDeltaRapidUpdate-TimeDelta: %w", err)
	} else {
		val.TimeDelta = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for PositionDeltaRapidUpdate-LatitudeDelta: %w", err)
	} else {
		val.LatitudeDelta = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for PositionDeltaRapidUpdate-LongitudeDelta: %w", err)
	} else {
		val.LongitudeDelta = v

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

func EncodePositionDeltaRapidUpdate(val *PositionDeltaRapidUpdate) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt16(val.TimeDelta, 16)
	w.writeInt16(val.LatitudeDelta, 16)
	w.writeInt16(val.LongitudeDelta, 16)
	w.skipBits(8)
	return w.Bytes(), w.Err()
}
func encodePositionDeltaRapidUpdateMsg(v Message) ([]byte, error) {
	val, ok := v.(*PositionDeltaRapidUpdate)
	if !ok {
		return nil, fmt.Errorf("expected *PositionDeltaRapidUpdate, got %T", v)
	}
	return EncodePositionDeltaRapidUpdate(val)
}

type AltitudeDeltaRapidUpdate struct {
	Info          MessageInfo `json:"info"`
	Sid           *uint8      `json:"sid"`
	TimeDelta     *int16      `json:"timeDelta"`
	GnssQuality   *uint8      `json:"gnssQuality"`
	Direction     *uint8      `json:"direction"`
	Cog           *float32    `json:"cog"`
	AltitudeDelta *int16      `json:"altitudeDelta"`
}

func (a *AltitudeDeltaRapidUpdate) PGNNumber() uint32 { return 129028 }
func DecodeAltitudeDeltaRapidUpdate(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AltitudeDeltaRapidUpdate{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AltitudeDeltaRapidUpdate-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for AltitudeDeltaRapidUpdate-TimeDelta: %w", err)
	} else {
		val.TimeDelta = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(2); err != nil {
		return nil, fmt.Errorf("parse failed for AltitudeDeltaRapidUpdate-GnssQuality: %w", err)
	} else {
		val.GnssQuality = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(2); err != nil {
		return nil, fmt.Errorf("parse failed for AltitudeDeltaRapidUpdate-Direction: %w", err)
	} else {
		val.Direction = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for AltitudeDeltaRapidUpdate-Cog: %w", err)
	} else {
		val.Cog = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for AltitudeDeltaRapidUpdate-AltitudeDelta: %w", err)
	} else {
		val.AltitudeDelta = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeAltitudeDeltaRapidUpdate(val *AltitudeDeltaRapidUpdate) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeInt16(val.TimeDelta, 16)
	w.writeUInt8(val.GnssQuality, 2)
	w.writeUInt8(val.Direction, 2)
	w.skipBits(4)
	w.writeUnsignedResolution(val.Cog, 16, 0.0001)
	w.writeInt16(val.AltitudeDelta, 16)
	return w.Bytes(), w.Err()
}
func encodeAltitudeDeltaRapidUpdateMsg(v Message) ([]byte, error) {
	val, ok := v.(*AltitudeDeltaRapidUpdate)
	if !ok {
		return nil, fmt.Errorf("expected *AltitudeDeltaRapidUpdate, got %T", v)
	}
	return EncodeAltitudeDeltaRapidUpdate(val)
}

type GnssPositionData struct {
	Info              MessageInfo                  `json:"info"`
	Sid               *uint8                       `json:"sid"`
	Date              *uint16                      `json:"date"`
	Time              *float32                     `json:"time"`
	Latitude          *float64                     `json:"latitude"`
	Longitude         *float64                     `json:"longitude"`
	Altitude          *units.Distance              `json:"altitude"`
	GnssType          GnsConst                     `json:"gnssType"`
	Method            GnsMethodConst               `json:"method"`
	Integrity         GnsIntegrityConst            `json:"integrity"`
	NumberOfSvs       *uint8                       `json:"numberOfSvs"`
	Hdop              *float32                     `json:"hdop"`
	Pdop              *float32                     `json:"pdop"`
	GeoidalSeparation *units.Distance              `json:"geoidalSeparation"`
	ReferenceStations *uint8                       `json:"referenceStations"`
	Repeating1        []GnssPositionDataRepeating1 `json:"repeating1"`
}
type GnssPositionDataRepeating1 struct {
	ReferenceStationType  GnsConst `json:"referenceStationType"`
	ReferenceStationId    *uint16  `json:"referenceStationId"`
	AgeOfDgnssCorrections *float32 `json:"ageOfDgnssCorrections"`
}

func (g *GnssPositionData) PGNNumber() uint32 { return 129029 }
func DecodeGnssPositionData(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GnssPositionData{}
	val.Info = Info
	var repeat1Count uint16 = 0
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GnssPositionData-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for GnssPositionData-Date: %w", err)
	} else {
		val.Date = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(32, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for GnssPositionData-Time: %w", err)
	} else {
		val.Time = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution64Override(64, 1e-16); err != nil {
		return nil, fmt.Errorf("parse failed for GnssPositionData-Latitude: %w", err)
	} else {
		val.Latitude = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution64Override(64, 1e-16); err != nil {
		return nil, fmt.Errorf("parse failed for GnssPositionData-Longitude: %w", err)
	} else {
		val.Longitude = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(64, 1e-06); err != nil {
		return nil, fmt.Errorf("parse failed for GnssPositionData-Altitude: %w", err)
	} else {
		val.Altitude = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for GnssPositionData-GnssType: %w", err)
	} else {
		val.GnssType = GnsConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for GnssPositionData-Method: %w", err)
	} else {
		val.Method = GnsMethodConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for GnssPositionData-Integrity: %w", err)
	} else {
		val.Integrity = GnsIntegrityConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GnssPositionData-NumberOfSvs: %w", err)
	} else {
		val.NumberOfSvs = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for GnssPositionData-Hdop: %w", err)
	} else {
		val.Hdop = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for GnssPositionData-Pdop: %w", err)
	} else {
		val.Pdop = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(32, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for GnssPositionData-GeoidalSeparation: %w", err)
	} else {
		val.GeoidalSeparation = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GnssPositionData-ReferenceStations: %w", err)
	} else {
		val.ReferenceStations = v
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
	val.Repeating1 = make([]GnssPositionDataRepeating1, 0)
	i := 0
	for {
		var rep GnssPositionDataRepeating1
		if v, err := stream.readLookupField(4); err != nil {
			return nil, fmt.Errorf("parse failed for GnssPositionData-ReferenceStationType: %w", err)
		} else {
			rep.ReferenceStationType = GnsConst(v)
		}
		if v, err := stream.readUInt16(12); err != nil {
			return nil, fmt.Errorf("parse failed for GnssPositionData-ReferenceStationId: %w", err)
		} else {
			rep.ReferenceStationId = v
		}
		if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
			return nil, fmt.Errorf("parse failed for GnssPositionData-AgeOfDgnssCorrections: %w", err)
		} else {
			rep.AgeOfDgnssCorrections = v
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

func EncodeGnssPositionData(val *GnssPositionData) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt16(val.Date, 16)
	w.writeUnsignedResolution(val.Time, 32, 0.0001)
	w.writeSignedResolution64Override(val.Latitude, 64, 1e-16)
	w.writeSignedResolution64Override(val.Longitude, 64, 1e-16)
	var altitudeRaw *float32
	if val.Altitude != nil {
		altitudeRaw = &val.Altitude.Value
	}
	w.writeSignedResolution(altitudeRaw, 64, 1e-06)
	w.writeLookupField(uint64(val.GnssType), 4)
	w.writeLookupField(uint64(val.Method), 4)
	w.writeLookupField(uint64(val.Integrity), 2)
	w.skipBits(6)
	w.writeUInt8(val.NumberOfSvs, 8)
	w.writeSignedResolution(val.Hdop, 16, 0.01)
	w.writeSignedResolution(val.Pdop, 16, 0.01)
	var geoidalSeparationRaw *float32
	if val.GeoidalSeparation != nil {
		geoidalSeparationRaw = &val.GeoidalSeparation.Value
	}
	w.writeSignedResolution(geoidalSeparationRaw, 32, 0.01)
	w.writeUInt8(val.ReferenceStations, 8)
	for _, rep := range val.Repeating1 {
		w.writeLookupField(uint64(rep.ReferenceStationType), 4)
		w.writeUInt16(rep.ReferenceStationId, 12)
		w.writeUnsignedResolution(rep.AgeOfDgnssCorrections, 16, 0.01)
	}
	return w.Bytes(), w.Err()
}
func encodeGnssPositionDataMsg(v Message) ([]byte, error) {
	val, ok := v.(*GnssPositionData)
	if !ok {
		return nil, fmt.Errorf("expected *GnssPositionData, got %T", v)
	}
	return EncodeGnssPositionData(val)
}

type TimeDate struct {
	Info        MessageInfo `json:"info"`
	Date        *uint16     `json:"date"`
	Time        *float32    `json:"time"`
	LocalOffset *float32    `json:"localOffset"`
}

func (t *TimeDate) PGNNumber() uint32 { return 129033 }
func DecodeTimeDate(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &TimeDate{}
	val.Info = Info
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for TimeDate-Date: %w", err)
	} else {
		val.Date = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(32, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for TimeDate-Time: %w", err)
	} else {
		val.Time = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 60); err != nil {
		return nil, fmt.Errorf("parse failed for TimeDate-LocalOffset: %w", err)
	} else {
		val.LocalOffset = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeTimeDate(val *TimeDate) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.Date, 16)
	w.writeUnsignedResolution(val.Time, 32, 0.0001)
	w.writeSignedResolution(val.LocalOffset, 16, 60)
	return w.Bytes(), w.Err()
}
func encodeTimeDateMsg(v Message) ([]byte, error) {
	val, ok := v.(*TimeDate)
	if !ok {
		return nil, fmt.Errorf("expected *TimeDate, got %T", v)
	}
	return EncodeTimeDate(val)
}

type Datum struct {
	Info           MessageInfo     `json:"info"`
	LocalDatum     string          `json:"localDatum"`
	DeltaLatitude  *float64        `json:"deltaLatitude"`
	DeltaLongitude *float64        `json:"deltaLongitude"`
	DeltaAltitude  *units.Distance `json:"deltaAltitude"`
	ReferenceDatum string          `json:"referenceDatum"`
}

func (d *Datum) PGNNumber() uint32 { return 129044 }
func DecodeDatum(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Datum{}
	val.Info = Info
	if v, err := stream.readFixedString(32); err != nil {
		return nil, fmt.Errorf("parse failed for Datum-LocalDatum: %w", err)
	} else {
		val.LocalDatum = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for Datum-DeltaLatitude: %w", err)
	} else {
		val.DeltaLatitude = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for Datum-DeltaLongitude: %w", err)
	} else {
		val.DeltaLongitude = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(32, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for Datum-DeltaAltitude: %w", err)
	} else {
		val.DeltaAltitude = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(32); err != nil {
		return nil, fmt.Errorf("parse failed for Datum-ReferenceDatum: %w", err)
	} else {
		val.ReferenceDatum = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeDatum(val *Datum) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeFixedString(val.LocalDatum, 32)
	w.writeSignedResolution64Override(val.DeltaLatitude, 32, 1e-07)
	w.writeSignedResolution64Override(val.DeltaLongitude, 32, 1e-07)
	var deltaAltitudeRaw *float32
	if val.DeltaAltitude != nil {
		deltaAltitudeRaw = &val.DeltaAltitude.Value
	}
	w.writeSignedResolution(deltaAltitudeRaw, 32, 0.01)
	w.writeFixedString(val.ReferenceDatum, 32)
	return w.Bytes(), w.Err()
}
func encodeDatumMsg(v Message) ([]byte, error) {
	val, ok := v.(*Datum)
	if !ok {
		return nil, fmt.Errorf("expected *Datum, got %T", v)
	}
	return EncodeDatum(val)
}

type UserDatum struct {
	Info                       MessageInfo     `json:"info"`
	DeltaX                     *units.Distance `json:"deltaX"`
	DeltaY                     *units.Distance `json:"deltaY"`
	DeltaZ                     *units.Distance `json:"deltaZ"`
	RotationInX                *float32        `json:"rotationInX"`
	RotationInY                *float32        `json:"rotationInY"`
	RotationInZ                *float32        `json:"rotationInZ"`
	Scale                      *float32        `json:"scale"`
	EllipsoidSemiMajorAxis     *units.Distance `json:"ellipsoidSemiMajorAxis"`
	EllipsoidFlatteningInverse *float32        `json:"ellipsoidFlatteningInverse"`
	DatumName                  string          `json:"datumName"`
}

func (u *UserDatum) PGNNumber() uint32 { return 129045 }
func DecodeUserDatum(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &UserDatum{}
	val.Info = Info
	if v, err := stream.readSignedResolution(32, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for UserDatum-DeltaX: %w", err)
	} else {
		val.DeltaX = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(32, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for UserDatum-DeltaY: %w", err)
	} else {
		val.DeltaY = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(32, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for UserDatum-DeltaZ: %w", err)
	} else {
		val.DeltaZ = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFloat32(); err != nil {
		return nil, fmt.Errorf("parse failed for UserDatum-RotationInX: %w", err)
	} else {
		val.RotationInX = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFloat32(); err != nil {
		return nil, fmt.Errorf("parse failed for UserDatum-RotationInY: %w", err)
	} else {
		val.RotationInY = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFloat32(); err != nil {
		return nil, fmt.Errorf("parse failed for UserDatum-RotationInZ: %w", err)
	} else {
		val.RotationInZ = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFloat32(); err != nil {
		return nil, fmt.Errorf("parse failed for UserDatum-Scale: %w", err)
	} else {
		val.Scale = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(32, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for UserDatum-EllipsoidSemiMajorAxis: %w", err)
	} else {
		val.EllipsoidSemiMajorAxis = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFloat32(); err != nil {
		return nil, fmt.Errorf("parse failed for UserDatum-EllipsoidFlatteningInverse: %w", err)
	} else {
		val.EllipsoidFlatteningInverse = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readFixedString(32); err != nil {
		return nil, fmt.Errorf("parse failed for UserDatum-DatumName: %w", err)
	} else {
		val.DatumName = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeUserDatum(val *UserDatum) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	var deltaXRaw *float32
	if val.DeltaX != nil {
		deltaXRaw = &val.DeltaX.Value
	}
	w.writeSignedResolution(deltaXRaw, 32, 0.01)
	var deltaYRaw *float32
	if val.DeltaY != nil {
		deltaYRaw = &val.DeltaY.Value
	}
	w.writeSignedResolution(deltaYRaw, 32, 0.01)
	var deltaZRaw *float32
	if val.DeltaZ != nil {
		deltaZRaw = &val.DeltaZ.Value
	}
	w.writeSignedResolution(deltaZRaw, 32, 0.01)
	w.writeFloat32(val.RotationInX)
	w.writeFloat32(val.RotationInY)
	w.writeFloat32(val.RotationInZ)
	w.writeFloat32(val.Scale)
	var ellipsoidSemiMajorAxisRaw *float32
	if val.EllipsoidSemiMajorAxis != nil {
		ellipsoidSemiMajorAxisRaw = &val.EllipsoidSemiMajorAxis.Value
	}
	w.writeSignedResolution(ellipsoidSemiMajorAxisRaw, 32, 0.01)
	w.writeFloat32(val.EllipsoidFlatteningInverse)
	w.writeFixedString(val.DatumName, 32)
	return w.Bytes(), w.Err()
}
func encodeUserDatumMsg(v Message) ([]byte, error) {
	val, ok := v.(*UserDatum)
	if !ok {
		return nil, fmt.Errorf("expected *UserDatum, got %T", v)
	}
	return EncodeUserDatum(val)
}

type CrossTrackError struct {
	Info                 MessageInfo       `json:"info"`
	Sid                  *uint8            `json:"sid"`
	XteMode              ResidualModeConst `json:"xteMode"`
	NavigationTerminated YesNoConst        `json:"navigationTerminated"`
	Xte                  *units.Distance   `json:"xte"`
}

func (c *CrossTrackError) PGNNumber() uint32 { return 129283 }
func DecodeCrossTrackError(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &CrossTrackError{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for CrossTrackError-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for CrossTrackError-XteMode: %w", err)
	} else {
		val.XteMode = ResidualModeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for CrossTrackError-NavigationTerminated: %w", err)
	} else {
		val.NavigationTerminated = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(32, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for CrossTrackError-Xte: %w", err)
	} else {
		val.Xte = nullableUnit(units.Meter, v, units.NewDistance)

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

func EncodeCrossTrackError(val *CrossTrackError) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeLookupField(uint64(val.XteMode), 4)
	w.skipBits(2)
	w.writeLookupField(uint64(val.NavigationTerminated), 2)
	var xteRaw *float32
	if val.Xte != nil {
		xteRaw = &val.Xte.Value
	}
	w.writeSignedResolution(xteRaw, 32, 0.01)
	w.skipBits(16)
	return w.Bytes(), w.Err()
}
func encodeCrossTrackErrorMsg(v Message) ([]byte, error) {
	val, ok := v.(*CrossTrackError)
	if !ok {
		return nil, fmt.Errorf("expected *CrossTrackError, got %T", v)
	}
	return EncodeCrossTrackError(val)
}

type NavigationData struct {
	Info                                 MessageInfo             `json:"info"`
	Sid                                  *uint8                  `json:"sid"`
	DistanceToWaypoint                   *units.Distance         `json:"distanceToWaypoint"`
	CourseBearingReference               DirectionReferenceConst `json:"courseBearingReference"`
	PerpendicularCrossed                 YesNoConst              `json:"perpendicularCrossed"`
	ArrivalCircleEntered                 YesNoConst              `json:"arrivalCircleEntered"`
	CalculationType                      BearingModeConst        `json:"calculationType"`
	EtaTime                              *float32                `json:"etaTime"`
	EtaDate                              *uint16                 `json:"etaDate"`
	BearingOriginToDestinationWaypoint   *float32                `json:"bearingOriginToDestinationWaypoint"`
	BearingPositionToDestinationWaypoint *float32                `json:"bearingPositionToDestinationWaypoint"`
	OriginWaypointNumber                 *uint32                 `json:"originWaypointNumber"`
	DestinationWaypointNumber            *uint32                 `json:"destinationWaypointNumber"`
	DestinationLatitude                  *float64                `json:"destinationLatitude"`
	DestinationLongitude                 *float64                `json:"destinationLongitude"`
	WaypointClosingVelocity              *units.Velocity         `json:"waypointClosingVelocity"`
}

func (n *NavigationData) PGNNumber() uint32 { return 129284 }
func DecodeNavigationData(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &NavigationData{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationData-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(32, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationData-DistanceToWaypoint: %w", err)
	} else {
		val.DistanceToWaypoint = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationData-CourseBearingReference: %w", err)
	} else {
		val.CourseBearingReference = DirectionReferenceConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationData-PerpendicularCrossed: %w", err)
	} else {
		val.PerpendicularCrossed = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationData-ArrivalCircleEntered: %w", err)
	} else {
		val.ArrivalCircleEntered = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationData-CalculationType: %w", err)
	} else {
		val.CalculationType = BearingModeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(32, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationData-EtaTime: %w", err)
	} else {
		val.EtaTime = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationData-EtaDate: %w", err)
	} else {
		val.EtaDate = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationData-BearingOriginToDestinationWaypoint: %w", err)
	} else {
		val.BearingOriginToDestinationWaypoint = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationData-BearingPositionToDestinationWaypoint: %w", err)
	} else {
		val.BearingPositionToDestinationWaypoint = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationData-OriginWaypointNumber: %w", err)
	} else {
		val.OriginWaypointNumber = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationData-DestinationWaypointNumber: %w", err)
	} else {
		val.DestinationWaypointNumber = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationData-DestinationLatitude: %w", err)
	} else {
		val.DestinationLatitude = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationData-DestinationLongitude: %w", err)
	} else {
		val.DestinationLongitude = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationData-WaypointClosingVelocity: %w", err)
	} else {
		val.WaypointClosingVelocity = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeNavigationData(val *NavigationData) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	var distanceToWaypointRaw *float32
	if val.DistanceToWaypoint != nil {
		distanceToWaypointRaw = &val.DistanceToWaypoint.Value
	}
	w.writeUnsignedResolution(distanceToWaypointRaw, 32, 0.01)
	w.writeLookupField(uint64(val.CourseBearingReference), 2)
	w.writeLookupField(uint64(val.PerpendicularCrossed), 2)
	w.writeLookupField(uint64(val.ArrivalCircleEntered), 2)
	w.writeLookupField(uint64(val.CalculationType), 2)
	w.writeUnsignedResolution(val.EtaTime, 32, 0.0001)
	w.writeUInt16(val.EtaDate, 16)
	w.writeUnsignedResolution(val.BearingOriginToDestinationWaypoint, 16, 0.0001)
	w.writeUnsignedResolution(val.BearingPositionToDestinationWaypoint, 16, 0.0001)
	w.writeUInt32(val.OriginWaypointNumber, 32)
	w.writeUInt32(val.DestinationWaypointNumber, 32)
	w.writeSignedResolution64Override(val.DestinationLatitude, 32, 1e-07)
	w.writeSignedResolution64Override(val.DestinationLongitude, 32, 1e-07)
	var waypointClosingVelocityRaw *float32
	if val.WaypointClosingVelocity != nil {
		waypointClosingVelocityRaw = &val.WaypointClosingVelocity.Value
	}
	w.writeSignedResolution(waypointClosingVelocityRaw, 16, 0.01)
	return w.Bytes(), w.Err()
}
func encodeNavigationDataMsg(v Message) ([]byte, error) {
	val, ok := v.(*NavigationData)
	if !ok {
		return nil, fmt.Errorf("expected *NavigationData, got %T", v)
	}
	return EncodeNavigationData(val)
}

type NavigationRouteWpInformation struct {
	Info                              MessageInfo                              `json:"info"`
	StartRps                          *uint16                                  `json:"startRps"`
	Nitems                            *uint16                                  `json:"nitems"`
	DatabaseId                        *uint16                                  `json:"databaseId"`
	RouteId                           *uint16                                  `json:"routeId"`
	NavigationDirectionInRoute        DirectionConst                           `json:"navigationDirectionInRoute"`
	SupplementaryRouteWpDataAvailable OffOnConst                               `json:"supplementaryRouteWpDataAvailable"`
	RouteName                         string                                   `json:"routeName"`
	Repeating1                        []NavigationRouteWpInformationRepeating1 `json:"repeating1"`
}
type NavigationRouteWpInformationRepeating1 struct {
	WpId        *uint16  `json:"wpId"`
	WpName      string   `json:"wpName"`
	WpLatitude  *float64 `json:"wpLatitude"`
	WpLongitude *float64 `json:"wpLongitude"`
}

func (n *NavigationRouteWpInformation) PGNNumber() uint32 { return 129285 }
func DecodeNavigationRouteWpInformation(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &NavigationRouteWpInformation{}
	val.Info = Info
	var repeat1Count uint16 = 0
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationRouteWpInformation-StartRps: %w", err)
	} else {
		val.StartRps = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationRouteWpInformation-Nitems: %w", err)
	} else {
		val.Nitems = v
		if v != nil {
			repeat1Count = uint16(*v)
		}

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationRouteWpInformation-DatabaseId: %w", err)
	} else {
		val.DatabaseId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationRouteWpInformation-RouteId: %w", err)
	} else {
		val.RouteId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationRouteWpInformation-NavigationDirectionInRoute: %w", err)
	} else {
		val.NavigationDirectionInRoute = DirectionConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationRouteWpInformation-SupplementaryRouteWpDataAvailable: %w", err)
	} else {
		val.SupplementaryRouteWpDataAvailable = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(3)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readStringWithLengthAndControl(); err != nil {
		return nil, fmt.Errorf("parse failed for NavigationRouteWpInformation-RouteName: %w", err)
	} else {
		val.RouteName = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	if repeat1Count == 0 {
		return val, nil
	}
	val.Repeating1 = make([]NavigationRouteWpInformationRepeating1, 0)
	i := 0
	for {
		var rep NavigationRouteWpInformationRepeating1
		if v, err := stream.readUInt16(16); err != nil {
			return nil, fmt.Errorf("parse failed for NavigationRouteWpInformation-WpId: %w", err)
		} else {
			rep.WpId = v
		}
		if v, err := stream.readStringWithLengthAndControl(); err != nil {
			return nil, fmt.Errorf("parse failed for NavigationRouteWpInformation-WpName: %w", err)
		} else {
			rep.WpName = v
		}
		if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
			return nil, fmt.Errorf("parse failed for NavigationRouteWpInformation-WpLatitude: %w", err)
		} else {
			rep.WpLatitude = v
		}
		if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
			return nil, fmt.Errorf("parse failed for NavigationRouteWpInformation-WpLongitude: %w", err)
		} else {
			rep.WpLongitude = v
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

func EncodeNavigationRouteWpInformation(val *NavigationRouteWpInformation) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt16(val.StartRps, 16)
	w.writeUInt16(val.Nitems, 16)
	w.writeUInt16(val.DatabaseId, 16)
	w.writeUInt16(val.RouteId, 16)
	w.writeLookupField(uint64(val.NavigationDirectionInRoute), 3)
	w.writeLookupField(uint64(val.SupplementaryRouteWpDataAvailable), 2)
	w.skipBits(3)
	w.writeStringWithLengthAndControl(val.RouteName)
	w.skipBits(8)
	for _, rep := range val.Repeating1 {
		w.writeUInt16(rep.WpId, 16)
		w.writeStringWithLengthAndControl(rep.WpName)
		w.writeSignedResolution64Override(rep.WpLatitude, 32, 1e-07)
		w.writeSignedResolution64Override(rep.WpLongitude, 32, 1e-07)
	}
	return w.Bytes(), w.Err()
}
func encodeNavigationRouteWpInformationMsg(v Message) ([]byte, error) {
	val, ok := v.(*NavigationRouteWpInformation)
	if !ok {
		return nil, fmt.Errorf("expected *NavigationRouteWpInformation, got %T", v)
	}
	return EncodeNavigationRouteWpInformation(val)
}

type SetDriftRapidUpdate struct {
	Info         MessageInfo             `json:"info"`
	Sid          *uint8                  `json:"sid"`
	SetReference DirectionReferenceConst `json:"setReference"`
	Set          *float32                `json:"set"`
	Drift        *units.Velocity         `json:"drift"`
}

func (s *SetDriftRapidUpdate) PGNNumber() uint32 { return 129291 }
func DecodeSetDriftRapidUpdate(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SetDriftRapidUpdate{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SetDriftRapidUpdate-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SetDriftRapidUpdate-SetReference: %w", err)
	} else {
		val.SetReference = DirectionReferenceConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SetDriftRapidUpdate-Set: %w", err)
	} else {
		val.Set = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for SetDriftRapidUpdate-Drift: %w", err)
	} else {
		val.Drift = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

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

func EncodeSetDriftRapidUpdate(val *SetDriftRapidUpdate) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeLookupField(uint64(val.SetReference), 2)
	w.skipBits(6)
	w.writeUnsignedResolution(val.Set, 16, 0.0001)
	var driftRaw *float32
	if val.Drift != nil {
		driftRaw = &val.Drift.Value
	}
	w.writeUnsignedResolution(driftRaw, 16, 0.01)
	w.skipBits(16)
	return w.Bytes(), w.Err()
}
func encodeSetDriftRapidUpdateMsg(v Message) ([]byte, error) {
	val, ok := v.(*SetDriftRapidUpdate)
	if !ok {
		return nil, fmt.Errorf("expected *SetDriftRapidUpdate, got %T", v)
	}
	return EncodeSetDriftRapidUpdate(val)
}

type GnssDops struct {
	Info        MessageInfo   `json:"info"`
	Sid         *uint8        `json:"sid"`
	DesiredMode GnssModeConst `json:"desiredMode"`
	ActualMode  GnssModeConst `json:"actualMode"`
	Hdop        *float32      `json:"hdop"`
	Vdop        *float32      `json:"vdop"`
	Tdop        *float32      `json:"tdop"`
}

func (g *GnssDops) PGNNumber() uint32 { return 129539 }
func DecodeGnssDops(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GnssDops{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GnssDops-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for GnssDops-DesiredMode: %w", err)
	} else {
		val.DesiredMode = GnssModeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for GnssDops-ActualMode: %w", err)
	} else {
		val.ActualMode = GnssModeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for GnssDops-Hdop: %w", err)
	} else {
		val.Hdop = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for GnssDops-Vdop: %w", err)
	} else {
		val.Vdop = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for GnssDops-Tdop: %w", err)
	} else {
		val.Tdop = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeGnssDops(val *GnssDops) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeLookupField(uint64(val.DesiredMode), 3)
	w.writeLookupField(uint64(val.ActualMode), 3)
	w.skipBits(2)
	w.writeSignedResolution(val.Hdop, 16, 0.01)
	w.writeSignedResolution(val.Vdop, 16, 0.01)
	w.writeSignedResolution(val.Tdop, 16, 0.01)
	return w.Bytes(), w.Err()
}
func encodeGnssDopsMsg(v Message) ([]byte, error) {
	val, ok := v.(*GnssDops)
	if !ok {
		return nil, fmt.Errorf("expected *GnssDops, got %T", v)
	}
	return EncodeGnssDops(val)
}

type GnssSatsInView struct {
	Info              MessageInfo                `json:"info"`
	Sid               *uint8                     `json:"sid"`
	RangeResidualMode RangeResidualModeConst     `json:"rangeResidualMode"`
	SatsInView        *uint8                     `json:"satsInView"`
	Repeating1        []GnssSatsInViewRepeating1 `json:"repeating1"`
}
type GnssSatsInViewRepeating1 struct {
	Prn            *uint8               `json:"prn"`
	Elevation      *float32             `json:"elevation"`
	Azimuth        *float32             `json:"azimuth"`
	Snr            *float32             `json:"snr"`
	RangeResiduals *int32               `json:"rangeResiduals"`
	Status         SatelliteStatusConst `json:"status"`
}

func (g *GnssSatsInView) PGNNumber() uint32 { return 129540 }
func DecodeGnssSatsInView(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GnssSatsInView{}
	val.Info = Info
	var repeat1Count uint16 = 0
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GnssSatsInView-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for GnssSatsInView-RangeResidualMode: %w", err)
	} else {
		val.RangeResidualMode = RangeResidualModeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GnssSatsInView-SatsInView: %w", err)
	} else {
		val.SatsInView = v
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
	val.Repeating1 = make([]GnssSatsInViewRepeating1, 0)
	i := 0
	for {
		var rep GnssSatsInViewRepeating1
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for GnssSatsInView-Prn: %w", err)
		} else {
			rep.Prn = v
		}
		if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
			return nil, fmt.Errorf("parse failed for GnssSatsInView-Elevation: %w", err)
		} else {
			rep.Elevation = v
		}
		if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
			return nil, fmt.Errorf("parse failed for GnssSatsInView-Azimuth: %w", err)
		} else {
			rep.Azimuth = v
		}
		if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
			return nil, fmt.Errorf("parse failed for GnssSatsInView-Snr: %w", err)
		} else {
			rep.Snr = v
		}
		if v, err := stream.readInt32(32); err != nil {
			return nil, fmt.Errorf("parse failed for GnssSatsInView-RangeResiduals: %w", err)
		} else {
			rep.RangeResiduals = v
		}
		if v, err := stream.readLookupField(4); err != nil {
			return nil, fmt.Errorf("parse failed for GnssSatsInView-Status: %w", err)
		} else {
			rep.Status = SatelliteStatusConst(v)
		}
		stream.skipBits(4)
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

func EncodeGnssSatsInView(val *GnssSatsInView) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeLookupField(uint64(val.RangeResidualMode), 2)
	w.skipBits(6)
	w.writeUInt8(val.SatsInView, 8)
	for _, rep := range val.Repeating1 {
		w.writeUInt8(rep.Prn, 8)
		w.writeSignedResolution(rep.Elevation, 16, 0.0001)
		w.writeUnsignedResolution(rep.Azimuth, 16, 0.0001)
		w.writeUnsignedResolution(rep.Snr, 16, 0.01)
		w.writeInt32(rep.RangeResiduals, 32)
		w.writeLookupField(uint64(rep.Status), 4)
		w.skipBits(4)
	}
	return w.Bytes(), w.Err()
}
func encodeGnssSatsInViewMsg(v Message) ([]byte, error) {
	val, ok := v.(*GnssSatsInView)
	if !ok {
		return nil, fmt.Errorf("expected *GnssSatsInView, got %T", v)
	}
	return EncodeGnssSatsInView(val)
}

type GpsAlmanacData struct {
	Info                     MessageInfo `json:"info"`
	Prn                      *uint8      `json:"prn"`
	GpsWeekNumber            *uint16     `json:"gpsWeekNumber"`
	SvHealthBits             []uint8     `json:"svHealthBits"`
	Eccentricity             *float32    `json:"eccentricity"`
	AlmanacReferenceTime     *float32    `json:"almanacReferenceTime"`
	InclinationAngle         *float32    `json:"inclinationAngle"`
	RateOfRightAscension     *float64    `json:"rateOfRightAscension"`
	RootOfSemiMajorAxis      *float32    `json:"rootOfSemiMajorAxis"`
	ArgumentOfPerigee        *float32    `json:"argumentOfPerigee"`
	LongitudeOfAscensionNode *float32    `json:"longitudeOfAscensionNode"`
	MeanAnomaly              *float32    `json:"meanAnomaly"`
	ClockParameter1          *float32    `json:"clockParameter1"`
	ClockParameter2          *float64    `json:"clockParameter2"`
}

func (g *GpsAlmanacData) PGNNumber() uint32 { return 129541 }
func DecodeGpsAlmanacData(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &GpsAlmanacData{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GpsAlmanacData-Prn: %w", err)
	} else {
		val.Prn = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for GpsAlmanacData-GpsWeekNumber: %w", err)
	} else {
		val.GpsWeekNumber = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(8); err != nil {
		return nil, fmt.Errorf("parse failed for GpsAlmanacData-SvHealthBits: %w", err)
	} else {
		val.SvHealthBits = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 4.76837e-07); err != nil {
		return nil, fmt.Errorf("parse failed for GpsAlmanacData-Eccentricity: %w", err)
	} else {
		val.Eccentricity = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(8, 4096); err != nil {
		return nil, fmt.Errorf("parse failed for GpsAlmanacData-AlmanacReferenceTime: %w", err)
	} else {
		val.AlmanacReferenceTime = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 1.90735e-06); err != nil {
		return nil, fmt.Errorf("parse failed for GpsAlmanacData-InclinationAngle: %w", err)
	} else {
		val.InclinationAngle = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution64Override(16, 3.63798e-12); err != nil {
		return nil, fmt.Errorf("parse failed for GpsAlmanacData-RateOfRightAscension: %w", err)
	} else {
		val.RateOfRightAscension = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(24, 0.000488281); err != nil {
		return nil, fmt.Errorf("parse failed for GpsAlmanacData-RootOfSemiMajorAxis: %w", err)
	} else {
		val.RootOfSemiMajorAxis = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(24, 1.19209e-07); err != nil {
		return nil, fmt.Errorf("parse failed for GpsAlmanacData-ArgumentOfPerigee: %w", err)
	} else {
		val.ArgumentOfPerigee = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(24, 1.19209e-07); err != nil {
		return nil, fmt.Errorf("parse failed for GpsAlmanacData-LongitudeOfAscensionNode: %w", err)
	} else {
		val.LongitudeOfAscensionNode = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(24, 1.19209e-07); err != nil {
		return nil, fmt.Errorf("parse failed for GpsAlmanacData-MeanAnomaly: %w", err)
	} else {
		val.MeanAnomaly = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(11, 9.53674e-07); err != nil {
		return nil, fmt.Errorf("parse failed for GpsAlmanacData-ClockParameter1: %w", err)
	} else {
		val.ClockParameter1 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution64Override(11, 3.63798e-12); err != nil {
		return nil, fmt.Errorf("parse failed for GpsAlmanacData-ClockParameter2: %w", err)
	} else {
		val.ClockParameter2 = v

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

func EncodeGpsAlmanacData(val *GpsAlmanacData) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Prn, 8)
	w.writeUInt16(val.GpsWeekNumber, 16)
	w.writeBinaryData(val.SvHealthBits, 8)
	w.writeUnsignedResolution(val.Eccentricity, 16, 4.76837e-07)
	w.writeUnsignedResolution(val.AlmanacReferenceTime, 8, 4096)
	w.writeSignedResolution(val.InclinationAngle, 16, 1.90735e-06)
	w.writeSignedResolution64Override(val.RateOfRightAscension, 16, 3.63798e-12)
	w.writeUnsignedResolution(val.RootOfSemiMajorAxis, 24, 0.000488281)
	w.writeSignedResolution(val.ArgumentOfPerigee, 24, 1.19209e-07)
	w.writeSignedResolution(val.LongitudeOfAscensionNode, 24, 1.19209e-07)
	w.writeSignedResolution(val.MeanAnomaly, 24, 1.19209e-07)
	w.writeSignedResolution(val.ClockParameter1, 11, 9.53674e-07)
	w.writeSignedResolution64Override(val.ClockParameter2, 11, 3.63798e-12)
	w.skipBits(2)
	return w.Bytes(), w.Err()
}
func encodeGpsAlmanacDataMsg(v Message) ([]byte, error) {
	val, ok := v.(*GpsAlmanacData)
	if !ok {
		return nil, fmt.Errorf("expected *GpsAlmanacData, got %T", v)
	}
	return EncodeGpsAlmanacData(val)
}

type SmallCraftStatus struct {
	Info             MessageInfo `json:"info"`
	PortTrimTab      *int8       `json:"portTrimTab"`
	StarboardTrimTab *int8       `json:"starboardTrimTab"`
}

func (s *SmallCraftStatus) PGNNumber() uint32 { return 130576 }
func DecodeSmallCraftStatus(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SmallCraftStatus{}
	val.Info = Info
	if v, err := stream.readInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SmallCraftStatus-PortTrimTab: %w", err)
	} else {
		val.PortTrimTab = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SmallCraftStatus-StarboardTrimTab: %w", err)
	} else {
		val.StarboardTrimTab = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(48)
	if stream.isEOF() {
		return val, nil
	}
	return val, nil
}

func EncodeSmallCraftStatus(val *SmallCraftStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeInt8(val.PortTrimTab, 8)
	w.writeInt8(val.StarboardTrimTab, 8)
	w.skipBits(48)
	return w.Bytes(), w.Err()
}
func encodeSmallCraftStatusMsg(v Message) ([]byte, error) {
	val, ok := v.(*SmallCraftStatus)
	if !ok {
		return nil, fmt.Errorf("expected *SmallCraftStatus, got %T", v)
	}
	return EncodeSmallCraftStatus(val)
}

type VesselSpeedComponents struct {
	Info                              MessageInfo     `json:"info"`
	LongitudinalSpeedWaterReferenced  *units.Velocity `json:"longitudinalSpeedWaterReferenced"`
	TransverseSpeedWaterReferenced    *units.Velocity `json:"transverseSpeedWaterReferenced"`
	LongitudinalSpeedGroundReferenced *units.Velocity `json:"longitudinalSpeedGroundReferenced"`
	TransverseSpeedGroundReferenced   *units.Velocity `json:"transverseSpeedGroundReferenced"`
	SternSpeedWaterReferenced         *units.Velocity `json:"sternSpeedWaterReferenced"`
	SternSpeedGroundReferenced        *units.Velocity `json:"sternSpeedGroundReferenced"`
}

func (v *VesselSpeedComponents) PGNNumber() uint32 { return 130578 }
func DecodeVesselSpeedComponents(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &VesselSpeedComponents{}
	val.Info = Info
	if v, err := stream.readSignedResolution(16, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for VesselSpeedComponents-LongitudinalSpeedWaterReferenced: %w", err)
	} else {
		val.LongitudinalSpeedWaterReferenced = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for VesselSpeedComponents-TransverseSpeedWaterReferenced: %w", err)
	} else {
		val.TransverseSpeedWaterReferenced = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for VesselSpeedComponents-LongitudinalSpeedGroundReferenced: %w", err)
	} else {
		val.LongitudinalSpeedGroundReferenced = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for VesselSpeedComponents-TransverseSpeedGroundReferenced: %w", err)
	} else {
		val.TransverseSpeedGroundReferenced = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for VesselSpeedComponents-SternSpeedWaterReferenced: %w", err)
	} else {
		val.SternSpeedWaterReferenced = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readSignedResolution(16, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for VesselSpeedComponents-SternSpeedGroundReferenced: %w", err)
	} else {
		val.SternSpeedGroundReferenced = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeVesselSpeedComponents(val *VesselSpeedComponents) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	var longitudinalSpeedWaterReferencedRaw *float32
	if val.LongitudinalSpeedWaterReferenced != nil {
		longitudinalSpeedWaterReferencedRaw = &val.LongitudinalSpeedWaterReferenced.Value
	}
	w.writeSignedResolution(longitudinalSpeedWaterReferencedRaw, 16, 0.001)
	var transverseSpeedWaterReferencedRaw *float32
	if val.TransverseSpeedWaterReferenced != nil {
		transverseSpeedWaterReferencedRaw = &val.TransverseSpeedWaterReferenced.Value
	}
	w.writeSignedResolution(transverseSpeedWaterReferencedRaw, 16, 0.001)
	var longitudinalSpeedGroundReferencedRaw *float32
	if val.LongitudinalSpeedGroundReferenced != nil {
		longitudinalSpeedGroundReferencedRaw = &val.LongitudinalSpeedGroundReferenced.Value
	}
	w.writeSignedResolution(longitudinalSpeedGroundReferencedRaw, 16, 0.001)
	var transverseSpeedGroundReferencedRaw *float32
	if val.TransverseSpeedGroundReferenced != nil {
		transverseSpeedGroundReferencedRaw = &val.TransverseSpeedGroundReferenced.Value
	}
	w.writeSignedResolution(transverseSpeedGroundReferencedRaw, 16, 0.001)
	var sternSpeedWaterReferencedRaw *float32
	if val.SternSpeedWaterReferenced != nil {
		sternSpeedWaterReferencedRaw = &val.SternSpeedWaterReferenced.Value
	}
	w.writeSignedResolution(sternSpeedWaterReferencedRaw, 16, 0.001)
	var sternSpeedGroundReferencedRaw *float32
	if val.SternSpeedGroundReferenced != nil {
		sternSpeedGroundReferencedRaw = &val.SternSpeedGroundReferenced.Value
	}
	w.writeSignedResolution(sternSpeedGroundReferencedRaw, 16, 0.001)
	return w.Bytes(), w.Err()
}
func encodeVesselSpeedComponentsMsg(v Message) ([]byte, error) {
	val, ok := v.(*VesselSpeedComponents)
	if !ok {
		return nil, fmt.Errorf("expected *VesselSpeedComponents, got %T", v)
	}
	return EncodeVesselSpeedComponents(val)
}
