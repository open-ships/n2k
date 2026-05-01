package pgn

import (
	"fmt"

	"github.com/open-ships/n2k/units"
)

type AisClassAPositionReport struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	UserId *uint32 `json:"userId"`
	Longitude *float64 `json:"longitude"`
	Latitude *float64 `json:"latitude"`
	PositionAccuracy PositionAccuracyConst `json:"positionAccuracy"`
	Raim RaimFlagConst `json:"raim"`
	TimeStamp TimeStampConst `json:"timeStamp"`
	Cog *float32 `json:"cog"`
	Sog *units.Velocity `json:"sog"`
	CommunicationState []uint8 `json:"communicationState"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
	Heading *float32 `json:"heading"`
	RateOfTurn *float32 `json:"rateOfTurn"`
	NavStatus NavStatusConst `json:"navStatus"`
	SpecialManeuverIndicator AisSpecialManeuverConst `json:"specialManeuverIndicator"`
	SequenceId *uint8 `json:"sequenceId"`
}
func DecodeAisClassAPositionReport(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisClassAPositionReport
	val.Info = Info
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-UserId: %w", err)
	} else {
		val.UserId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-Longitude: %w", err)
	} else {
		val.Longitude = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-Latitude: %w", err)
	} else {
		val.Latitude = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-PositionAccuracy: %w", err)
	} else {
		val.PositionAccuracy = PositionAccuracyConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-Raim: %w", err)
	} else {
		val.Raim = RaimFlagConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-TimeStamp: %w", err)
	} else {
		val.TimeStamp = TimeStampConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-Cog: %w", err)
	} else {
		val.Cog = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-Sog: %w", err)
	} else {
		val.Sog = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readBinaryData(19); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-CommunicationState: %w", err)
	} else {
		val.CommunicationState = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-Heading: %w", err)
	} else {
		val.Heading = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 3.125e-05); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-RateOfTurn: %w", err)
	} else {
		val.RateOfTurn = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-NavStatus: %w", err)
	} else {
		val.NavStatus = NavStatusConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-SpecialManeuverIndicator: %w", err)
	} else {
		val.SpecialManeuverIndicator = AisSpecialManeuverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}
	stream.skipBits(3)
	if stream.isEOF() {
		return val, nil
		}
	stream.skipBits(5)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAPositionReport-SequenceId: %w", err)
	} else {
		val.SequenceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AisClassBPositionReport struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	UserId *uint32 `json:"userId"`
	Longitude *float64 `json:"longitude"`
	Latitude *float64 `json:"latitude"`
	PositionAccuracy PositionAccuracyConst `json:"positionAccuracy"`
	Raim RaimFlagConst `json:"raim"`
	TimeStamp TimeStampConst `json:"timeStamp"`
	Cog *float32 `json:"cog"`
	Sog *units.Velocity `json:"sog"`
	CommunicationState []uint8 `json:"communicationState"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
	Heading *float32 `json:"heading"`
	UnitType AisTypeConst `json:"unitType"`
	IntegratedDisplay YesNoConst `json:"integratedDisplay"`
	Dsc YesNoConst `json:"dsc"`
	Band AisBandConst `json:"band"`
	CanHandleMsg22 YesNoConst `json:"canHandleMsg22"`
	AisMode AisModeConst `json:"aisMode"`
	AisCommunicationState AisCommunicationStateConst `json:"aisCommunicationState"`
}
func DecodeAisClassBPositionReport(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisClassBPositionReport
	val.Info = Info
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-UserId: %w", err)
	} else {
		val.UserId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-Longitude: %w", err)
	} else {
		val.Longitude = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-Latitude: %w", err)
	} else {
		val.Latitude = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-PositionAccuracy: %w", err)
	} else {
		val.PositionAccuracy = PositionAccuracyConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-Raim: %w", err)
	} else {
		val.Raim = RaimFlagConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-TimeStamp: %w", err)
	} else {
		val.TimeStamp = TimeStampConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-Cog: %w", err)
	} else {
		val.Cog = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-Sog: %w", err)
	} else {
		val.Sog = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readBinaryData(19); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-CommunicationState: %w", err)
	} else {
		val.CommunicationState = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-Heading: %w", err)
	} else {
		val.Heading = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-UnitType: %w", err)
	} else {
		val.UnitType = AisTypeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-IntegratedDisplay: %w", err)
	} else {
		val.IntegratedDisplay = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-Dsc: %w", err)
	} else {
		val.Dsc = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-Band: %w", err)
	} else {
		val.Band = AisBandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-CanHandleMsg22: %w", err)
	} else {
		val.CanHandleMsg22 = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-AisMode: %w", err)
	} else {
		val.AisMode = AisModeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBPositionReport-AisCommunicationState: %w", err)
	} else {
		val.AisCommunicationState = AisCommunicationStateConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(15)
	if stream.isEOF() {
		return val, nil
		}	
	return val, nil
}
type AisClassBExtendedPositionReport struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	UserId *uint32 `json:"userId"`
	Longitude *float64 `json:"longitude"`
	Latitude *float64 `json:"latitude"`
	PositionAccuracy PositionAccuracyConst `json:"positionAccuracy"`
	Raim RaimFlagConst `json:"raim"`
	TimeStamp TimeStampConst `json:"timeStamp"`
	Cog *float32 `json:"cog"`
	Sog *units.Velocity `json:"sog"`
	TypeOfShip ShipTypeConst `json:"typeOfShip"`
	TrueHeading *float32 `json:"trueHeading"`
	GnssType PositionFixDeviceConst `json:"gnssType"`
	Length *units.Distance `json:"length"`
	Beam *units.Distance `json:"beam"`
	PositionReferenceFromStarboard *units.Distance `json:"positionReferenceFromStarboard"`
	PositionReferenceFromBow *units.Distance `json:"positionReferenceFromBow"`
	Name string `json:"name"`
	Dte AvailableConst `json:"dte"`
	AisMode AisModeConst `json:"aisMode"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
}
func DecodeAisClassBExtendedPositionReport(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisClassBExtendedPositionReport
	val.Info = Info
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-UserId: %w", err)
	} else {
		val.UserId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-Longitude: %w", err)
	} else {
		val.Longitude = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-Latitude: %w", err)
	} else {
		val.Latitude = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-PositionAccuracy: %w", err)
	} else {
		val.PositionAccuracy = PositionAccuracyConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-Raim: %w", err)
	} else {
		val.Raim = RaimFlagConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-TimeStamp: %w", err)
	} else {
		val.TimeStamp = TimeStampConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-Cog: %w", err)
	} else {
		val.Cog = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-Sog: %w", err)
	} else {
		val.Sog = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
		}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-TypeOfShip: %w", err)
	} else {
		val.TypeOfShip = ShipTypeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-TrueHeading: %w", err)
	} else {
		val.TrueHeading = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-GnssType: %w", err)
	} else {
		val.GnssType = PositionFixDeviceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-Length: %w", err)
	} else {
		val.Length = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-Beam: %w", err)
	} else {
		val.Beam = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-PositionReferenceFromStarboard: %w", err)
	} else {
		val.PositionReferenceFromStarboard = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-PositionReferenceFromBow: %w", err)
	} else {
		val.PositionReferenceFromBow = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(160); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-Name: %w", err)
	} else {
		val.Name = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-Dte: %w", err)
	} else {
		val.Dte = AvailableConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-AisMode: %w", err)
	} else {
		val.AisMode = AisModeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBExtendedPositionReport-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(5)
	if stream.isEOF() {
		return val, nil
		}	
	return val, nil
}
type AisAidsToNavigationAtonReport struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	UserId *uint32 `json:"userId"`
	Longitude *float64 `json:"longitude"`
	Latitude *float64 `json:"latitude"`
	PositionAccuracy PositionAccuracyConst `json:"positionAccuracy"`
	Raim RaimFlagConst `json:"raim"`
	TimeStamp TimeStampConst `json:"timeStamp"`
	LengthDiameter *units.Distance `json:"lengthDiameter"`
	BeamDiameter *units.Distance `json:"beamDiameter"`
	PositionReferenceFromStarboardEdge *units.Distance `json:"positionReferenceFromStarboardEdge"`
	PositionReferenceFromTrueNorthFacingEdge *units.Distance `json:"positionReferenceFromTrueNorthFacingEdge"`
	AtonType AtonTypeConst `json:"atonType"`
	OffPositionIndicator YesNoConst `json:"offPositionIndicator"`
	VirtualAtonFlag YesNoConst `json:"virtualAtonFlag"`
	AssignedModeFlag AisAssignedModeConst `json:"assignedModeFlag"`
	PositionFixingDeviceType PositionFixDeviceConst `json:"positionFixingDeviceType"`
	AtonStatus []uint8 `json:"atonStatus"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
	AtonName string `json:"atonName"`
}
func DecodeAisAidsToNavigationAtonReport(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisAidsToNavigationAtonReport
	val.Info = Info
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-UserId: %w", err)
	} else {
		val.UserId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-Longitude: %w", err)
	} else {
		val.Longitude = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-Latitude: %w", err)
	} else {
		val.Latitude = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-PositionAccuracy: %w", err)
	} else {
		val.PositionAccuracy = PositionAccuracyConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-Raim: %w", err)
	} else {
		val.Raim = RaimFlagConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-TimeStamp: %w", err)
	} else {
		val.TimeStamp = TimeStampConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-LengthDiameter: %w", err)
	} else {
		val.LengthDiameter = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-BeamDiameter: %w", err)
	} else {
		val.BeamDiameter = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-PositionReferenceFromStarboardEdge: %w", err)
	} else {
		val.PositionReferenceFromStarboardEdge = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-PositionReferenceFromTrueNorthFacingEdge: %w", err)
	} else {
		val.PositionReferenceFromTrueNorthFacingEdge = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-AtonType: %w", err)
	} else {
		val.AtonType = AtonTypeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-OffPositionIndicator: %w", err)
	} else {
		val.OffPositionIndicator = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-VirtualAtonFlag: %w", err)
	} else {
		val.VirtualAtonFlag = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-AssignedModeFlag: %w", err)
	} else {
		val.AssignedModeFlag = AisAssignedModeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(1)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-PositionFixingDeviceType: %w", err)
	} else {
		val.PositionFixingDeviceType = PositionFixDeviceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(3)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readBinaryData(8); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-AtonStatus: %w", err)
	} else {
		val.AtonStatus = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(3)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readStringWithLengthAndControl(); err != nil {
		return nil, fmt.Errorf("parse failed for AisAidsToNavigationAtonReport-AtonName: %w", err)
	} else {
		val.AtonName = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AisUtcAndDateReport struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	UserId *uint32 `json:"userId"`
	Longitude *float64 `json:"longitude"`
	Latitude *float64 `json:"latitude"`
	PositionAccuracy PositionAccuracyConst `json:"positionAccuracy"`
	Raim RaimFlagConst `json:"raim"`
	PositionTime *float32 `json:"positionTime"`
	CommunicationState []uint8 `json:"communicationState"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
	PositionDate *uint16 `json:"positionDate"`
	GnssType PositionFixDeviceConst `json:"gnssType"`
}
func DecodeAisUtcAndDateReport(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisUtcAndDateReport
	val.Info = Info
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcAndDateReport-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcAndDateReport-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcAndDateReport-UserId: %w", err)
	} else {
		val.UserId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcAndDateReport-Longitude: %w", err)
	} else {
		val.Longitude = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcAndDateReport-Latitude: %w", err)
	} else {
		val.Latitude = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcAndDateReport-PositionAccuracy: %w", err)
	} else {
		val.PositionAccuracy = PositionAccuracyConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcAndDateReport-Raim: %w", err)
	} else {
		val.Raim = RaimFlagConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUnsignedResolution(32, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcAndDateReport-PositionTime: %w", err)
	} else {
		val.PositionTime = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readBinaryData(19); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcAndDateReport-CommunicationState: %w", err)
	} else {
		val.CommunicationState = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcAndDateReport-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcAndDateReport-PositionDate: %w", err)
	} else {
		val.PositionDate = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcAndDateReport-GnssType: %w", err)
	} else {
		val.GnssType = PositionFixDeviceConst(v)

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
type AisClassAStaticAndVoyageRelatedData struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	UserId *uint32 `json:"userId"`
	ImoNumber *uint32 `json:"imoNumber"`
	Callsign string `json:"callsign"`
	Name string `json:"name"`
	TypeOfShip ShipTypeConst `json:"typeOfShip"`
	Length *units.Distance `json:"length"`
	Beam *units.Distance `json:"beam"`
	PositionReferenceFromStarboard *units.Distance `json:"positionReferenceFromStarboard"`
	PositionReferenceFromBow *units.Distance `json:"positionReferenceFromBow"`
	EtaDate *uint16 `json:"etaDate"`
	EtaTime *float32 `json:"etaTime"`
	Draft *units.Distance `json:"draft"`
	Destination string `json:"destination"`
	AisVersionIndicator AisVersionConst `json:"aisVersionIndicator"`
	GnssType PositionFixDeviceConst `json:"gnssType"`
	Dte AvailableConst `json:"dte"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
}
func DecodeAisClassAStaticAndVoyageRelatedData(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisClassAStaticAndVoyageRelatedData
	val.Info = Info
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-UserId: %w", err)
	} else {
		val.UserId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-ImoNumber: %w", err)
	} else {
		val.ImoNumber = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(56); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-Callsign: %w", err)
	} else {
		val.Callsign = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(160); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-Name: %w", err)
	} else {
		val.Name = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-TypeOfShip: %w", err)
	} else {
		val.TypeOfShip = ShipTypeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-Length: %w", err)
	} else {
		val.Length = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-Beam: %w", err)
	} else {
		val.Beam = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-PositionReferenceFromStarboard: %w", err)
	} else {
		val.PositionReferenceFromStarboard = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-PositionReferenceFromBow: %w", err)
	} else {
		val.PositionReferenceFromBow = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-EtaDate: %w", err)
	} else {
		val.EtaDate = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-EtaTime: %w", err)
	} else {
		val.EtaTime = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-Draft: %w", err)
	} else {
		val.Draft = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(160); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-Destination: %w", err)
	} else {
		val.Destination = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-AisVersionIndicator: %w", err)
	} else {
		val.AisVersionIndicator = AisVersionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-GnssType: %w", err)
	} else {
		val.GnssType = PositionFixDeviceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-Dte: %w", err)
	} else {
		val.Dte = AvailableConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(1)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassAStaticAndVoyageRelatedData-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(3)
	if stream.isEOF() {
		return val, nil
		}	
	return val, nil
}
type AisAddressedBinaryMessage struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	SourceId *uint32 `json:"sourceId"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
	SequenceNumber *uint8 `json:"sequenceNumber"`
	DestinationId *uint32 `json:"destinationId"`
	RetransmitFlag *uint8 `json:"retransmitFlag"`
	NumberOfBitsInBinaryDataField *uint16 `json:"numberOfBitsInBinaryDataField"`
	BinaryData []uint8 `json:"binaryData"`
}
func DecodeAisAddressedBinaryMessage(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisAddressedBinaryMessage
	val.Info = Info
		var binaryLength uint16 = 0
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedBinaryMessage-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedBinaryMessage-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedBinaryMessage-SourceId: %w", err)
	} else {
		val.SourceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(1)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedBinaryMessage-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedBinaryMessage-SequenceNumber: %w", err)
	} else {
		val.SequenceNumber = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedBinaryMessage-DestinationId: %w", err)
	} else {
		val.DestinationId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt8(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedBinaryMessage-RetransmitFlag: %w", err)
	} else {
		val.RetransmitFlag = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(1)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedBinaryMessage-NumberOfBitsInBinaryDataField: %w", err)
	} else {
		val.NumberOfBitsInBinaryDataField = v
		if v != nil {
			binaryLength = uint16(*v)
		}

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readBinaryData(binaryLength); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedBinaryMessage-BinaryData: %w", err)
	} else {
		val.BinaryData = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AisAcknowledge struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	SourceId *uint32 `json:"sourceId"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
	DestinationId1 *uint32 `json:"destinationId1"`
	SequenceNumberForId1 []uint8 `json:"sequenceNumberForId1"`
	SequenceNumberForIdN []uint8 `json:"sequenceNumberForIdN"`
}
func DecodeAisAcknowledge(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisAcknowledge
	val.Info = Info
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisAcknowledge-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisAcknowledge-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisAcknowledge-SourceId: %w", err)
	} else {
		val.SourceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(1)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisAcknowledge-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisAcknowledge-DestinationId1: %w", err)
	} else {
		val.DestinationId1 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readBinaryData(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisAcknowledge-SequenceNumberForId1: %w", err)
	} else {
		val.SequenceNumberForId1 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readBinaryData(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisAcknowledge-SequenceNumberForIdN: %w", err)
	} else {
		val.SequenceNumberForIdN = v

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
type AisBinaryBroadcastMessage struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	SourceId *uint32 `json:"sourceId"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
	NumberOfBitsInBinaryDataField *uint16 `json:"numberOfBitsInBinaryDataField"`
	BinaryData []uint8 `json:"binaryData"`
}
func DecodeAisBinaryBroadcastMessage(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisBinaryBroadcastMessage
	val.Info = Info
		var binaryLength uint16 = 0
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisBinaryBroadcastMessage-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisBinaryBroadcastMessage-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisBinaryBroadcastMessage-SourceId: %w", err)
	} else {
		val.SourceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(1)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisBinaryBroadcastMessage-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for AisBinaryBroadcastMessage-NumberOfBitsInBinaryDataField: %w", err)
	} else {
		val.NumberOfBitsInBinaryDataField = v
		if v != nil {
			binaryLength = uint16(*v)
		}

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readBinaryData(binaryLength); err != nil {
		return nil, fmt.Errorf("parse failed for AisBinaryBroadcastMessage-BinaryData: %w", err)
	} else {
		val.BinaryData = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AisUtcDateInquiry struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	SourceId *uint32 `json:"sourceId"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
	DestinationId *uint32 `json:"destinationId"`
}
func DecodeAisUtcDateInquiry(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisUtcDateInquiry
	val.Info = Info
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcDateInquiry-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcDateInquiry-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcDateInquiry-SourceId: %w", err)
	} else {
		val.SourceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcDateInquiry-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(3)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisUtcDateInquiry-DestinationId: %w", err)
	} else {
		val.DestinationId = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AisAddressedSafetyRelatedMessage struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	SourceId *uint32 `json:"sourceId"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
	SequenceNumber *uint8 `json:"sequenceNumber"`
	DestinationId *uint32 `json:"destinationId"`
	RetransmitFlag *uint8 `json:"retransmitFlag"`
	SafetyRelatedText string `json:"safetyRelatedText"`
}
func DecodeAisAddressedSafetyRelatedMessage(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisAddressedSafetyRelatedMessage
	val.Info = Info
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedSafetyRelatedMessage-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedSafetyRelatedMessage-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedSafetyRelatedMessage-SourceId: %w", err)
	} else {
		val.SourceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedSafetyRelatedMessage-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedSafetyRelatedMessage-SequenceNumber: %w", err)
	} else {
		val.SequenceNumber = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(1)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedSafetyRelatedMessage-DestinationId: %w", err)
	} else {
		val.DestinationId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(1); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedSafetyRelatedMessage-RetransmitFlag: %w", err)
	} else {
		val.RetransmitFlag = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(7)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readFixedString(936); err != nil {
		return nil, fmt.Errorf("parse failed for AisAddressedSafetyRelatedMessage-SafetyRelatedText: %w", err)
	} else {
		val.SafetyRelatedText = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AisSafetyRelatedBroadcastMessage struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	SourceId *uint32 `json:"sourceId"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
	SafetyRelatedText string `json:"safetyRelatedText"`
}
func DecodeAisSafetyRelatedBroadcastMessage(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisSafetyRelatedBroadcastMessage
	val.Info = Info
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisSafetyRelatedBroadcastMessage-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisSafetyRelatedBroadcastMessage-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisSafetyRelatedBroadcastMessage-SourceId: %w", err)
	} else {
		val.SourceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisSafetyRelatedBroadcastMessage-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(3)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readFixedString(1296); err != nil {
		return nil, fmt.Errorf("parse failed for AisSafetyRelatedBroadcastMessage-SafetyRelatedText: %w", err)
	} else {
		val.SafetyRelatedText = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AisInterrogation struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	SourceId *uint32 `json:"sourceId"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
	DestinationId1 *uint32 `json:"destinationId1"`
	MessageId11 AisMessageIdConst `json:"messageId11"`
	SlotOffset11 *uint16 `json:"slotOffset11"`
	MessageId12 AisMessageIdConst `json:"messageId12"`
	SlotOffset12 *uint16 `json:"slotOffset12"`
	DestinationId2 *uint32 `json:"destinationId2"`
	MessageId21 AisMessageIdConst `json:"messageId21"`
	SlotOffset21 *uint16 `json:"slotOffset21"`
	Sid *uint8 `json:"sid"`
}
func DecodeAisInterrogation(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisInterrogation
	val.Info = Info
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisInterrogation-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisInterrogation-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisInterrogation-SourceId: %w", err)
	} else {
		val.SourceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(1)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisInterrogation-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisInterrogation-DestinationId1: %w", err)
	} else {
		val.DestinationId1 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisInterrogation-MessageId11: %w", err)
	} else {
		val.MessageId11 = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(12); err != nil {
		return nil, fmt.Errorf("parse failed for AisInterrogation-SlotOffset11: %w", err)
	} else {
		val.SlotOffset11 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisInterrogation-MessageId12: %w", err)
	} else {
		val.MessageId12 = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(12); err != nil {
		return nil, fmt.Errorf("parse failed for AisInterrogation-SlotOffset12: %w", err)
	} else {
		val.SlotOffset12 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisInterrogation-DestinationId2: %w", err)
	} else {
		val.DestinationId2 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisInterrogation-MessageId21: %w", err)
	} else {
		val.MessageId21 = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(12); err != nil {
		return nil, fmt.Errorf("parse failed for AisInterrogation-SlotOffset21: %w", err)
	} else {
		val.SlotOffset21 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AisInterrogation-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AisDataLinkManagementMessage struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	SourceId *uint32 `json:"sourceId"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
	Repeating1 []AisDataLinkManagementMessageRepeating1 `json:"repeating1"`
}
type AisDataLinkManagementMessageRepeating1 struct {
	Offset *uint16 `json:"offset"`
	NumberOfSlots *uint8 `json:"numberOfSlots"`
	Timeout *uint8 `json:"timeout"`
	Increment *uint16 `json:"increment"`
}
func DecodeAisDataLinkManagementMessage(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisDataLinkManagementMessage
	val.Info = Info
		var repeat1Count uint16 = 0
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisDataLinkManagementMessage-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisDataLinkManagementMessage-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisDataLinkManagementMessage-SourceId: %w", err)
	} else {
		val.SourceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisDataLinkManagementMessage-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(3)
	if stream.isEOF() {
		return val, nil
		}
	val.Repeating1 = make([]AisDataLinkManagementMessageRepeating1, 0)
	i := 0 
	for {
		var rep AisDataLinkManagementMessageRepeating1
		if v, err := stream.readUInt16(16); err != nil {
			return nil, fmt.Errorf("parse failed for AisDataLinkManagementMessage-Offset: %w", err)
		} else {
			rep.Offset = v
		}
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for AisDataLinkManagementMessage-NumberOfSlots: %w", err)
		} else {
			rep.NumberOfSlots = v
		}
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for AisDataLinkManagementMessage-Timeout: %w", err)
		} else {
			rep.Timeout = v
		}
		if v, err := stream.readUInt16(16); err != nil {
			return nil, fmt.Errorf("parse failed for AisDataLinkManagementMessage-Increment: %w", err)
		} else {
			rep.Increment = v
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
type AisChannelManagement struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	SourceId *uint32 `json:"sourceId"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
	ChannelA *uint8 `json:"channelA"`
	ChannelB *uint8 `json:"channelB"`
	Power *uint8 `json:"power"`
	TxRxMode *uint8 `json:"txRxMode"`
	NorthEastLongitudeCorner1 *float64 `json:"northEastLongitudeCorner1"`
	NorthEastLatitudeCorner1 *float64 `json:"northEastLatitudeCorner1"`
	SouthWestLongitudeCorner1 *float64 `json:"southWestLongitudeCorner1"`
	SouthWestLatitudeCorner2 *float64 `json:"southWestLatitudeCorner2"`
	AddressedOrBroadcastMessageIndicator *uint8 `json:"addressedOrBroadcastMessageIndicator"`
	ChannelABandwidth *uint8 `json:"channelABandwidth"`
	ChannelBBandwidth *uint8 `json:"channelBBandwidth"`
	TransitionalZoneSize *uint8 `json:"transitionalZoneSize"`
}
func DecodeAisChannelManagement(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisChannelManagement
	val.Info = Info
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisChannelManagement-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisChannelManagement-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisChannelManagement-SourceId: %w", err)
	} else {
		val.SourceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisChannelManagement-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(3)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt8(7); err != nil {
		return nil, fmt.Errorf("parse failed for AisChannelManagement-ChannelA: %w", err)
	} else {
		val.ChannelA = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(7); err != nil {
		return nil, fmt.Errorf("parse failed for AisChannelManagement-ChannelB: %w", err)
	} else {
		val.ChannelB = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AisChannelManagement-Power: %w", err)
	} else {
		val.Power = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AisChannelManagement-TxRxMode: %w", err)
	} else {
		val.TxRxMode = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for AisChannelManagement-NorthEastLongitudeCorner1: %w", err)
	} else {
		val.NorthEastLongitudeCorner1 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for AisChannelManagement-NorthEastLatitudeCorner1: %w", err)
	} else {
		val.NorthEastLatitudeCorner1 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for AisChannelManagement-SouthWestLongitudeCorner1: %w", err)
	} else {
		val.SouthWestLongitudeCorner1 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for AisChannelManagement-SouthWestLatitudeCorner2: %w", err)
	} else {
		val.SouthWestLatitudeCorner2 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt8(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisChannelManagement-AddressedOrBroadcastMessageIndicator: %w", err)
	} else {
		val.AddressedOrBroadcastMessageIndicator = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(7); err != nil {
		return nil, fmt.Errorf("parse failed for AisChannelManagement-ChannelABandwidth: %w", err)
	} else {
		val.ChannelABandwidth = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(7); err != nil {
		return nil, fmt.Errorf("parse failed for AisChannelManagement-ChannelBBandwidth: %w", err)
	} else {
		val.ChannelBBandwidth = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AisChannelManagement-TransitionalZoneSize: %w", err)
	} else {
		val.TransitionalZoneSize = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AisClassBStaticDataMsg24PartA struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	UserId *uint32 `json:"userId"`
	Name string `json:"name"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
	SequenceId *uint8 `json:"sequenceId"`
}
func DecodeAisClassBStaticDataMsg24PartA(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisClassBStaticDataMsg24PartA
	val.Info = Info
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartA-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartA-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartA-UserId: %w", err)
	} else {
		val.UserId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(160); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartA-Name: %w", err)
	} else {
		val.Name = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartA-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(3)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartA-SequenceId: %w", err)
	} else {
		val.SequenceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AisClassBStaticDataMsg24PartB struct {
	Info MessageInfo `json:"info"`
	MessageId AisMessageIdConst `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	UserId *uint32 `json:"userId"`
	TypeOfShip ShipTypeConst `json:"typeOfShip"`
	VendorId string `json:"vendorId"`
	Callsign string `json:"callsign"`
	Length *units.Distance `json:"length"`
	Beam *units.Distance `json:"beam"`
	PositionReferenceFromStarboard *units.Distance `json:"positionReferenceFromStarboard"`
	PositionReferenceFromBow *units.Distance `json:"positionReferenceFromBow"`
	MothershipUserId *uint32 `json:"mothershipUserId"`
	AisTransceiverInformation AisTransceiverConst `json:"aisTransceiverInformation"`
	SequenceId *uint8 `json:"sequenceId"`
}
func DecodeAisClassBStaticDataMsg24PartB(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AisClassBStaticDataMsg24PartB
	val.Info = Info
	if v, err := stream.readLookupField(6); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartB-MessageId: %w", err)
	} else {
		val.MessageId = AisMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartB-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartB-UserId: %w", err)
	} else {
		val.UserId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartB-TypeOfShip: %w", err)
	} else {
		val.TypeOfShip = ShipTypeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(56); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartB-VendorId: %w", err)
	} else {
		val.VendorId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(56); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartB-Callsign: %w", err)
	} else {
		val.Callsign = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartB-Length: %w", err)
	} else {
		val.Length = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartB-Beam: %w", err)
	} else {
		val.Beam = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartB-PositionReferenceFromStarboard: %w", err)
	} else {
		val.PositionReferenceFromStarboard = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartB-PositionReferenceFromBow: %w", err)
	} else {
		val.PositionReferenceFromBow = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartB-MothershipUserId: %w", err)
	} else {
		val.MothershipUserId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(5); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartB-AisTransceiverInformation: %w", err)
	} else {
		val.AisTransceiverInformation = AisTransceiverConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(3)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AisClassBStaticDataMsg24PartB-SequenceId: %w", err)
	} else {
		val.SequenceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeAisClassAPositionReport(val *AisClassAPositionReport) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.UserId, 32)
	w.writeSignedResolution64Override(val.Longitude, 32, 1e-07)
	w.writeSignedResolution64Override(val.Latitude, 32, 1e-07)
	w.writeLookupField(uint64(val.PositionAccuracy), 1)
	w.writeLookupField(uint64(val.Raim), 1)
	w.writeLookupField(uint64(val.TimeStamp), 6)
	w.writeUnsignedResolution(val.Cog, 16, 0.0001)
	var sogRaw *float32
	if val.Sog != nil {
		sogRaw = &val.Sog.Value
	}
	w.writeUnsignedResolution(sogRaw, 16, 0.01)
	w.writeBinaryData(val.CommunicationState, 19)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.writeUnsignedResolution(val.Heading, 16, 0.0001)
	w.writeSignedResolution(val.RateOfTurn, 16, 3.125e-05)
	w.writeLookupField(uint64(val.NavStatus), 4)
	w.writeLookupField(uint64(val.SpecialManeuverIndicator), 2)
	w.skipBits(2)
	w.skipBits(3)
	w.skipBits(5)
	w.writeUInt8(val.SequenceId, 8)
	return w.Bytes(), nil
}
func encodeAisClassAPositionReportAny(v any) ([]byte, error) {
	val, ok := v.(*AisClassAPositionReport)
	if !ok {
		return nil, fmt.Errorf("expected *AisClassAPositionReport, got %T", v)
	}
	return EncodeAisClassAPositionReport(val)
}

func EncodeAisClassBPositionReport(val *AisClassBPositionReport) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.UserId, 32)
	w.writeSignedResolution64Override(val.Longitude, 32, 1e-07)
	w.writeSignedResolution64Override(val.Latitude, 32, 1e-07)
	w.writeLookupField(uint64(val.PositionAccuracy), 1)
	w.writeLookupField(uint64(val.Raim), 1)
	w.writeLookupField(uint64(val.TimeStamp), 6)
	w.writeUnsignedResolution(val.Cog, 16, 0.0001)
	var sogRaw *float32
	if val.Sog != nil {
		sogRaw = &val.Sog.Value
	}
	w.writeUnsignedResolution(sogRaw, 16, 0.01)
	w.writeBinaryData(val.CommunicationState, 19)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.writeUnsignedResolution(val.Heading, 16, 0.0001)
	w.skipBits(8)
	w.skipBits(2)
	w.writeLookupField(uint64(val.UnitType), 1)
	w.writeLookupField(uint64(val.IntegratedDisplay), 1)
	w.writeLookupField(uint64(val.Dsc), 1)
	w.writeLookupField(uint64(val.Band), 1)
	w.writeLookupField(uint64(val.CanHandleMsg22), 1)
	w.writeLookupField(uint64(val.AisMode), 1)
	w.writeLookupField(uint64(val.AisCommunicationState), 1)
	w.skipBits(15)
	return w.Bytes(), nil
}
func encodeAisClassBPositionReportAny(v any) ([]byte, error) {
	val, ok := v.(*AisClassBPositionReport)
	if !ok {
		return nil, fmt.Errorf("expected *AisClassBPositionReport, got %T", v)
	}
	return EncodeAisClassBPositionReport(val)
}

func EncodeAisClassBExtendedPositionReport(val *AisClassBExtendedPositionReport) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.UserId, 32)
	w.writeSignedResolution64Override(val.Longitude, 32, 1e-07)
	w.writeSignedResolution64Override(val.Latitude, 32, 1e-07)
	w.writeLookupField(uint64(val.PositionAccuracy), 1)
	w.writeLookupField(uint64(val.Raim), 1)
	w.writeLookupField(uint64(val.TimeStamp), 6)
	w.writeUnsignedResolution(val.Cog, 16, 0.0001)
	var sogRaw *float32
	if val.Sog != nil {
		sogRaw = &val.Sog.Value
	}
	w.writeUnsignedResolution(sogRaw, 16, 0.01)
	w.skipBits(8)
	w.skipBits(4)
	w.skipBits(4)
	w.writeLookupField(uint64(val.TypeOfShip), 8)
	w.writeUnsignedResolution(val.TrueHeading, 16, 0.0001)
	w.skipBits(4)
	w.writeLookupField(uint64(val.GnssType), 4)
	var lengthRaw *float32
	if val.Length != nil {
		lengthRaw = &val.Length.Value
	}
	w.writeUnsignedResolution(lengthRaw, 16, 0.1)
	var beamRaw *float32
	if val.Beam != nil {
		beamRaw = &val.Beam.Value
	}
	w.writeUnsignedResolution(beamRaw, 16, 0.1)
	var positionReferenceFromStarboardRaw *float32
	if val.PositionReferenceFromStarboard != nil {
		positionReferenceFromStarboardRaw = &val.PositionReferenceFromStarboard.Value
	}
	w.writeUnsignedResolution(positionReferenceFromStarboardRaw, 16, 0.1)
	var positionReferenceFromBowRaw *float32
	if val.PositionReferenceFromBow != nil {
		positionReferenceFromBowRaw = &val.PositionReferenceFromBow.Value
	}
	w.writeUnsignedResolution(positionReferenceFromBowRaw, 16, 0.1)
	w.writeFixedString(val.Name, 160)
	w.writeLookupField(uint64(val.Dte), 1)
	w.writeLookupField(uint64(val.AisMode), 1)
	w.skipBits(4)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.skipBits(5)
	return w.Bytes(), nil
}
func encodeAisClassBExtendedPositionReportAny(v any) ([]byte, error) {
	val, ok := v.(*AisClassBExtendedPositionReport)
	if !ok {
		return nil, fmt.Errorf("expected *AisClassBExtendedPositionReport, got %T", v)
	}
	return EncodeAisClassBExtendedPositionReport(val)
}

func EncodeAisAidsToNavigationAtonReport(val *AisAidsToNavigationAtonReport) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.UserId, 32)
	w.writeSignedResolution64Override(val.Longitude, 32, 1e-07)
	w.writeSignedResolution64Override(val.Latitude, 32, 1e-07)
	w.writeLookupField(uint64(val.PositionAccuracy), 1)
	w.writeLookupField(uint64(val.Raim), 1)
	w.writeLookupField(uint64(val.TimeStamp), 6)
	var lengthDiameterRaw *float32
	if val.LengthDiameter != nil {
		lengthDiameterRaw = &val.LengthDiameter.Value
	}
	w.writeUnsignedResolution(lengthDiameterRaw, 16, 0.1)
	var beamDiameterRaw *float32
	if val.BeamDiameter != nil {
		beamDiameterRaw = &val.BeamDiameter.Value
	}
	w.writeUnsignedResolution(beamDiameterRaw, 16, 0.1)
	var positionReferenceFromStarboardEdgeRaw *float32
	if val.PositionReferenceFromStarboardEdge != nil {
		positionReferenceFromStarboardEdgeRaw = &val.PositionReferenceFromStarboardEdge.Value
	}
	w.writeUnsignedResolution(positionReferenceFromStarboardEdgeRaw, 16, 0.1)
	var positionReferenceFromTrueNorthFacingEdgeRaw *float32
	if val.PositionReferenceFromTrueNorthFacingEdge != nil {
		positionReferenceFromTrueNorthFacingEdgeRaw = &val.PositionReferenceFromTrueNorthFacingEdge.Value
	}
	w.writeUnsignedResolution(positionReferenceFromTrueNorthFacingEdgeRaw, 16, 0.1)
	w.writeLookupField(uint64(val.AtonType), 5)
	w.writeLookupField(uint64(val.OffPositionIndicator), 1)
	w.writeLookupField(uint64(val.VirtualAtonFlag), 1)
	w.writeLookupField(uint64(val.AssignedModeFlag), 1)
	w.skipBits(1)
	w.writeLookupField(uint64(val.PositionFixingDeviceType), 4)
	w.skipBits(3)
	w.writeBinaryData(val.AtonStatus, 8)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.skipBits(3)
	w.writeStringWithLengthAndControl(val.AtonName)
	return w.Bytes(), nil
}
func encodeAisAidsToNavigationAtonReportAny(v any) ([]byte, error) {
	val, ok := v.(*AisAidsToNavigationAtonReport)
	if !ok {
		return nil, fmt.Errorf("expected *AisAidsToNavigationAtonReport, got %T", v)
	}
	return EncodeAisAidsToNavigationAtonReport(val)
}

func EncodeAisUtcAndDateReport(val *AisUtcAndDateReport) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.UserId, 32)
	w.writeSignedResolution64Override(val.Longitude, 32, 1e-07)
	w.writeSignedResolution64Override(val.Latitude, 32, 1e-07)
	w.writeLookupField(uint64(val.PositionAccuracy), 1)
	w.writeLookupField(uint64(val.Raim), 1)
	w.skipBits(6)
	w.writeUnsignedResolution(val.PositionTime, 32, 0.0001)
	w.writeBinaryData(val.CommunicationState, 19)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.writeUInt16(val.PositionDate, 16)
	w.skipBits(4)
	w.writeLookupField(uint64(val.GnssType), 4)
	w.skipBits(8)
	return w.Bytes(), nil
}
func encodeAisUtcAndDateReportAny(v any) ([]byte, error) {
	val, ok := v.(*AisUtcAndDateReport)
	if !ok {
		return nil, fmt.Errorf("expected *AisUtcAndDateReport, got %T", v)
	}
	return EncodeAisUtcAndDateReport(val)
}

func EncodeAisClassAStaticAndVoyageRelatedData(val *AisClassAStaticAndVoyageRelatedData) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.UserId, 32)
	w.writeUInt32(val.ImoNumber, 32)
	w.writeFixedString(val.Callsign, 56)
	w.writeFixedString(val.Name, 160)
	w.writeLookupField(uint64(val.TypeOfShip), 8)
	var lengthRaw *float32
	if val.Length != nil {
		lengthRaw = &val.Length.Value
	}
	w.writeUnsignedResolution(lengthRaw, 16, 0.1)
	var beamRaw *float32
	if val.Beam != nil {
		beamRaw = &val.Beam.Value
	}
	w.writeUnsignedResolution(beamRaw, 16, 0.1)
	var positionReferenceFromStarboardRaw *float32
	if val.PositionReferenceFromStarboard != nil {
		positionReferenceFromStarboardRaw = &val.PositionReferenceFromStarboard.Value
	}
	w.writeUnsignedResolution(positionReferenceFromStarboardRaw, 16, 0.1)
	var positionReferenceFromBowRaw *float32
	if val.PositionReferenceFromBow != nil {
		positionReferenceFromBowRaw = &val.PositionReferenceFromBow.Value
	}
	w.writeUnsignedResolution(positionReferenceFromBowRaw, 16, 0.1)
	w.writeUInt16(val.EtaDate, 16)
	w.writeUnsignedResolution(val.EtaTime, 32, 0.0001)
	var draftRaw *float32
	if val.Draft != nil {
		draftRaw = &val.Draft.Value
	}
	w.writeUnsignedResolution(draftRaw, 16, 0.01)
	w.writeFixedString(val.Destination, 160)
	w.writeLookupField(uint64(val.AisVersionIndicator), 2)
	w.writeLookupField(uint64(val.GnssType), 4)
	w.writeLookupField(uint64(val.Dte), 1)
	w.skipBits(1)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.skipBits(3)
	return w.Bytes(), nil
}
func encodeAisClassAStaticAndVoyageRelatedDataAny(v any) ([]byte, error) {
	val, ok := v.(*AisClassAStaticAndVoyageRelatedData)
	if !ok {
		return nil, fmt.Errorf("expected *AisClassAStaticAndVoyageRelatedData, got %T", v)
	}
	return EncodeAisClassAStaticAndVoyageRelatedData(val)
}

func EncodeAisAddressedBinaryMessage(val *AisAddressedBinaryMessage) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.SourceId, 32)
	w.skipBits(1)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.writeUInt8(val.SequenceNumber, 2)
	w.writeUInt32(val.DestinationId, 32)
	w.skipBits(6)
	w.writeUInt8(val.RetransmitFlag, 1)
	w.skipBits(1)
	w.writeUInt16(val.NumberOfBitsInBinaryDataField, 16)
	return w.Bytes(), nil
}
func encodeAisAddressedBinaryMessageAny(v any) ([]byte, error) {
	val, ok := v.(*AisAddressedBinaryMessage)
	if !ok {
		return nil, fmt.Errorf("expected *AisAddressedBinaryMessage, got %T", v)
	}
	return EncodeAisAddressedBinaryMessage(val)
}

func EncodeAisAcknowledge(val *AisAcknowledge) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.SourceId, 32)
	w.skipBits(1)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.skipBits(2)
	w.writeUInt32(val.DestinationId1, 32)
	w.writeBinaryData(val.SequenceNumberForId1, 2)
	w.skipBits(6)
	w.writeBinaryData(val.SequenceNumberForIdN, 2)
	w.skipBits(6)
	return w.Bytes(), nil
}
func encodeAisAcknowledgeAny(v any) ([]byte, error) {
	val, ok := v.(*AisAcknowledge)
	if !ok {
		return nil, fmt.Errorf("expected *AisAcknowledge, got %T", v)
	}
	return EncodeAisAcknowledge(val)
}

func EncodeAisBinaryBroadcastMessage(val *AisBinaryBroadcastMessage) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.SourceId, 32)
	w.skipBits(1)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.skipBits(2)
	w.writeUInt16(val.NumberOfBitsInBinaryDataField, 16)
	return w.Bytes(), nil
}
func encodeAisBinaryBroadcastMessageAny(v any) ([]byte, error) {
	val, ok := v.(*AisBinaryBroadcastMessage)
	if !ok {
		return nil, fmt.Errorf("expected *AisBinaryBroadcastMessage, got %T", v)
	}
	return EncodeAisBinaryBroadcastMessage(val)
}

func EncodeAisUtcDateInquiry(val *AisUtcDateInquiry) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.SourceId, 32)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.skipBits(3)
	w.writeUInt32(val.DestinationId, 32)
	return w.Bytes(), nil
}
func encodeAisUtcDateInquiryAny(v any) ([]byte, error) {
	val, ok := v.(*AisUtcDateInquiry)
	if !ok {
		return nil, fmt.Errorf("expected *AisUtcDateInquiry, got %T", v)
	}
	return EncodeAisUtcDateInquiry(val)
}

func EncodeAisAddressedSafetyRelatedMessage(val *AisAddressedSafetyRelatedMessage) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.SourceId, 32)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.writeUInt8(val.SequenceNumber, 2)
	w.skipBits(1)
	w.writeUInt32(val.DestinationId, 32)
	w.writeUInt8(val.RetransmitFlag, 1)
	w.skipBits(7)
	w.writeFixedString(val.SafetyRelatedText, 936)
	return w.Bytes(), nil
}
func encodeAisAddressedSafetyRelatedMessageAny(v any) ([]byte, error) {
	val, ok := v.(*AisAddressedSafetyRelatedMessage)
	if !ok {
		return nil, fmt.Errorf("expected *AisAddressedSafetyRelatedMessage, got %T", v)
	}
	return EncodeAisAddressedSafetyRelatedMessage(val)
}

func EncodeAisSafetyRelatedBroadcastMessage(val *AisSafetyRelatedBroadcastMessage) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.SourceId, 32)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.skipBits(3)
	w.writeFixedString(val.SafetyRelatedText, 1296)
	return w.Bytes(), nil
}
func encodeAisSafetyRelatedBroadcastMessageAny(v any) ([]byte, error) {
	val, ok := v.(*AisSafetyRelatedBroadcastMessage)
	if !ok {
		return nil, fmt.Errorf("expected *AisSafetyRelatedBroadcastMessage, got %T", v)
	}
	return EncodeAisSafetyRelatedBroadcastMessage(val)
}

func EncodeAisInterrogation(val *AisInterrogation) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.SourceId, 32)
	w.skipBits(1)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.skipBits(2)
	w.writeUInt32(val.DestinationId1, 32)
	w.writeLookupField(uint64(val.MessageId11), 6)
	w.writeUInt16(val.SlotOffset11, 12)
	w.skipBits(2)
	w.writeLookupField(uint64(val.MessageId12), 6)
	w.writeUInt16(val.SlotOffset12, 12)
	w.skipBits(2)
	w.writeUInt32(val.DestinationId2, 32)
	w.writeLookupField(uint64(val.MessageId21), 6)
	w.writeUInt16(val.SlotOffset21, 12)
	w.skipBits(2)
	w.skipBits(4)
	w.writeUInt8(val.Sid, 8)
	return w.Bytes(), nil
}
func encodeAisInterrogationAny(v any) ([]byte, error) {
	val, ok := v.(*AisInterrogation)
	if !ok {
		return nil, fmt.Errorf("expected *AisInterrogation, got %T", v)
	}
	return EncodeAisInterrogation(val)
}

func EncodeAisDataLinkManagementMessage(val *AisDataLinkManagementMessage) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.SourceId, 32)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.skipBits(3)
	for _, rep := range val.Repeating1 {
		w.writeUInt16(rep.Offset, 16)
		w.writeUInt8(rep.NumberOfSlots, 8)
		w.writeUInt8(rep.Timeout, 8)
		w.writeUInt16(rep.Increment, 16)
	}
	return w.Bytes(), nil
}
func encodeAisDataLinkManagementMessageAny(v any) ([]byte, error) {
	val, ok := v.(*AisDataLinkManagementMessage)
	if !ok {
		return nil, fmt.Errorf("expected *AisDataLinkManagementMessage, got %T", v)
	}
	return EncodeAisDataLinkManagementMessage(val)
}

func EncodeAisChannelManagement(val *AisChannelManagement) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.SourceId, 32)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.skipBits(3)
	w.writeUInt8(val.ChannelA, 7)
	w.writeUInt8(val.ChannelB, 7)
	w.skipBits(2)
	w.writeUInt8(val.Power, 8)
	w.writeUInt8(val.TxRxMode, 8)
	w.writeSignedResolution64Override(val.NorthEastLongitudeCorner1, 32, 1e-07)
	w.writeSignedResolution64Override(val.NorthEastLatitudeCorner1, 32, 1e-07)
	w.writeSignedResolution64Override(val.SouthWestLongitudeCorner1, 32, 1e-07)
	w.writeSignedResolution64Override(val.SouthWestLatitudeCorner2, 32, 1e-07)
	w.skipBits(6)
	w.writeUInt8(val.AddressedOrBroadcastMessageIndicator, 2)
	w.writeUInt8(val.ChannelABandwidth, 7)
	w.writeUInt8(val.ChannelBBandwidth, 7)
	w.skipBits(2)
	w.writeUInt8(val.TransitionalZoneSize, 8)
	return w.Bytes(), nil
}
func encodeAisChannelManagementAny(v any) ([]byte, error) {
	val, ok := v.(*AisChannelManagement)
	if !ok {
		return nil, fmt.Errorf("expected *AisChannelManagement, got %T", v)
	}
	return EncodeAisChannelManagement(val)
}

func EncodeAisClassBStaticDataMsg24PartA(val *AisClassBStaticDataMsg24PartA) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.UserId, 32)
	w.writeFixedString(val.Name, 160)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.skipBits(3)
	w.writeUInt8(val.SequenceId, 8)
	return w.Bytes(), nil
}
func encodeAisClassBStaticDataMsg24PartAAny(v any) ([]byte, error) {
	val, ok := v.(*AisClassBStaticDataMsg24PartA)
	if !ok {
		return nil, fmt.Errorf("expected *AisClassBStaticDataMsg24PartA, got %T", v)
	}
	return EncodeAisClassBStaticDataMsg24PartA(val)
}

func EncodeAisClassBStaticDataMsg24PartB(val *AisClassBStaticDataMsg24PartB) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.MessageId), 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt32(val.UserId, 32)
	w.writeLookupField(uint64(val.TypeOfShip), 8)
	w.writeFixedString(val.VendorId, 56)
	w.writeFixedString(val.Callsign, 56)
	var lengthRaw *float32
	if val.Length != nil {
		lengthRaw = &val.Length.Value
	}
	w.writeUnsignedResolution(lengthRaw, 16, 0.1)
	var beamRaw *float32
	if val.Beam != nil {
		beamRaw = &val.Beam.Value
	}
	w.writeUnsignedResolution(beamRaw, 16, 0.1)
	var positionReferenceFromStarboardRaw *float32
	if val.PositionReferenceFromStarboard != nil {
		positionReferenceFromStarboardRaw = &val.PositionReferenceFromStarboard.Value
	}
	w.writeUnsignedResolution(positionReferenceFromStarboardRaw, 16, 0.1)
	var positionReferenceFromBowRaw *float32
	if val.PositionReferenceFromBow != nil {
		positionReferenceFromBowRaw = &val.PositionReferenceFromBow.Value
	}
	w.writeUnsignedResolution(positionReferenceFromBowRaw, 16, 0.1)
	w.writeUInt32(val.MothershipUserId, 32)
	w.skipBits(2)
	w.skipBits(6)
	w.writeLookupField(uint64(val.AisTransceiverInformation), 5)
	w.skipBits(3)
	w.writeUInt8(val.SequenceId, 8)
	return w.Bytes(), nil
}
func encodeAisClassBStaticDataMsg24PartBAny(v any) ([]byte, error) {
	val, ok := v.(*AisClassBStaticDataMsg24PartB)
	if !ok {
		return nil, fmt.Errorf("expected *AisClassBStaticDataMsg24PartB, got %T", v)
	}
	return EncodeAisClassBStaticDataMsg24PartB(val)
}
