package pgn

import "fmt"

type SeatalkWirelessKeypadControl struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	Pid              *uint8                `json:"pid"`
	Variant          *uint8                `json:"variant"`
	BeepControl      *uint8                `json:"beepControl"`
}

func (s *SeatalkWirelessKeypadControl) PGNNumber() uint32 { return 61184 }

func DecodeSeatalkWirelessKeypadControl(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SeatalkWirelessKeypadControl{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkWirelessKeypadControl-ManufacturerCode: %w", err)
	} else {
		if v != 1851 {
			return nil, fmt.Errorf("match failed for SeatalkWirelessKeypadControl-ManufacturerCode: Expected %d != %d", 1851, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkWirelessKeypadControl-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SeatalkWirelessKeypadControl-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkWirelessKeypadControl-Pid: %w", err)
	} else {
		val.Pid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkWirelessKeypadControl-Variant: %w", err)
	} else {
		val.Variant = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkWirelessKeypadControl-BeepControl: %w", err)
	} else {
		val.BeepControl = v

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

func EncodeSeatalkWirelessKeypadControl(val *SeatalkWirelessKeypadControl) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Pid, 8)
	w.writeUInt8(val.Variant, 8)
	w.writeUInt8(val.BeepControl, 8)
	w.skipBits(24)
	return w.Bytes(), w.Err()
}

func encodeSeatalkWirelessKeypadControlMsg(v Message) ([]byte, error) {
	val, ok := v.(*SeatalkWirelessKeypadControl)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkWirelessKeypadControl, got %T", v)
	}
	return EncodeSeatalkWirelessKeypadControl(val)
}

type SeatalkWirelessKeypadLightControl struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    *uint8                `json:"proprietaryId"`
	Variant          *uint8                `json:"variant"`
	WirelessSetting  *uint8                `json:"wirelessSetting"`
	WiredSetting     *uint8                `json:"wiredSetting"`
}

func (s *SeatalkWirelessKeypadLightControl) PGNNumber() uint32 { return 61184 }

func DecodeSeatalkWirelessKeypadLightControl(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SeatalkWirelessKeypadLightControl{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkWirelessKeypadLightControl-ManufacturerCode: %w", err)
	} else {
		if v != 1851 {
			return nil, fmt.Errorf("match failed for SeatalkWirelessKeypadLightControl-ManufacturerCode: Expected %d != %d", 1851, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkWirelessKeypadLightControl-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SeatalkWirelessKeypadLightControl-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkWirelessKeypadLightControl-ProprietaryId: %w", err)
	} else {
		if v != nil && *v != 1 {
			return nil, fmt.Errorf("match failed for SeatalkWirelessKeypadLightControl-ProprietaryId: Expected %d != %d", 1, *v)
		}
		val.ProprietaryId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkWirelessKeypadLightControl-Variant: %w", err)
	} else {
		val.Variant = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkWirelessKeypadLightControl-WirelessSetting: %w", err)
	} else {
		val.WirelessSetting = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkWirelessKeypadLightControl-WiredSetting: %w", err)
	} else {
		val.WiredSetting = v

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

func EncodeSeatalkWirelessKeypadLightControl(val *SeatalkWirelessKeypadLightControl) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.ProprietaryId, 8)
	w.writeUInt8(val.Variant, 8)
	w.writeUInt8(val.WirelessSetting, 8)
	w.writeUInt8(val.WiredSetting, 8)
	w.skipBits(16)
	return w.Bytes(), w.Err()
}

func encodeSeatalkWirelessKeypadLightControlMsg(v Message) ([]byte, error) {
	val, ok := v.(*SeatalkWirelessKeypadLightControl)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkWirelessKeypadLightControl, got %T", v)
	}
	return EncodeSeatalkWirelessKeypadLightControl(val)
}

type SeatalkAlarm struct {
	Info             MessageInfo             `json:"info"`
	ManufacturerCode ManufacturerCodeConst   `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst       `json:"industryCode"`
	Sid              []uint8                 `json:"sid"`
	AlarmStatus      SeatalkAlarmStatusConst `json:"alarmStatus"`
	AlarmId          SeatalkAlarmIdConst     `json:"alarmId"`
	AlarmGroup       SeatalkAlarmGroupConst  `json:"alarmGroup"`
	AlarmPriority    []uint8                 `json:"alarmPriority"`
}

func (s *SeatalkAlarm) PGNNumber() uint32 { return 65288 }

func DecodeSeatalkAlarm(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SeatalkAlarm{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkAlarm-ManufacturerCode: %w", err)
	} else {
		if v != 1851 {
			return nil, fmt.Errorf("match failed for SeatalkAlarm-ManufacturerCode: Expected %d != %d", 1851, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkAlarm-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SeatalkAlarm-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkAlarm-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkAlarm-AlarmStatus: %w", err)
	} else {
		val.AlarmStatus = SeatalkAlarmStatusConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkAlarm-AlarmId: %w", err)
	} else {
		val.AlarmId = SeatalkAlarmIdConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkAlarm-AlarmGroup: %w", err)
	} else {
		val.AlarmGroup = SeatalkAlarmGroupConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(16); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkAlarm-AlarmPriority: %w", err)
	} else {
		val.AlarmPriority = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSeatalkAlarm(val *SeatalkAlarm) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeBinaryData(val.Sid, 8)
	w.writeLookupField(uint64(val.AlarmStatus), 8)
	w.writeLookupField(uint64(val.AlarmId), 8)
	w.writeLookupField(uint64(val.AlarmGroup), 8)
	w.writeBinaryData(val.AlarmPriority, 16)
	return w.Bytes(), w.Err()
}

func encodeSeatalkAlarmMsg(v Message) ([]byte, error) {
	val, ok := v.(*SeatalkAlarm)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkAlarm, got %T", v)
	}
	return EncodeSeatalkAlarm(val)
}

type SeatalkPilotWindDatum struct {
	Info                    MessageInfo           `json:"info"`
	ManufacturerCode        ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode            IndustryCodeConst     `json:"industryCode"`
	WindDatum               *float32              `json:"windDatum"`
	RollingAverageWindAngle *float32              `json:"rollingAverageWindAngle"`
}

func (s *SeatalkPilotWindDatum) PGNNumber() uint32 { return 65345 }

func DecodeSeatalkPilotWindDatum(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SeatalkPilotWindDatum{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotWindDatum-ManufacturerCode: %w", err)
	} else {
		if v != 1851 {
			return nil, fmt.Errorf("match failed for SeatalkPilotWindDatum-ManufacturerCode: Expected %d != %d", 1851, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotWindDatum-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SeatalkPilotWindDatum-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotWindDatum-WindDatum: %w", err)
	} else {
		val.WindDatum = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotWindDatum-RollingAverageWindAngle: %w", err)
	} else {
		val.RollingAverageWindAngle = v

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

func EncodeSeatalkPilotWindDatum(val *SeatalkPilotWindDatum) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUnsignedResolution(val.WindDatum, 16, 0.0001)
	w.writeUnsignedResolution(val.RollingAverageWindAngle, 16, 0.0001)
	w.skipBits(16)
	return w.Bytes(), w.Err()
}

func encodeSeatalkPilotWindDatumMsg(v Message) ([]byte, error) {
	val, ok := v.(*SeatalkPilotWindDatum)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkPilotWindDatum, got %T", v)
	}
	return EncodeSeatalkPilotWindDatum(val)
}

type SeatalkPilotHeading struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	Sid              []uint8               `json:"sid"`
	HeadingTrue      *float32              `json:"headingTrue"`
	HeadingMagnetic  *float32              `json:"headingMagnetic"`
}

func (s *SeatalkPilotHeading) PGNNumber() uint32 { return 65359 }

func DecodeSeatalkPilotHeading(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SeatalkPilotHeading{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotHeading-ManufacturerCode: %w", err)
	} else {
		if v != 1851 {
			return nil, fmt.Errorf("match failed for SeatalkPilotHeading-ManufacturerCode: Expected %d != %d", 1851, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotHeading-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SeatalkPilotHeading-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotHeading-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotHeading-HeadingTrue: %w", err)
	} else {
		val.HeadingTrue = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotHeading-HeadingMagnetic: %w", err)
	} else {
		val.HeadingMagnetic = v

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

func EncodeSeatalkPilotHeading(val *SeatalkPilotHeading) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeBinaryData(val.Sid, 8)
	w.writeUnsignedResolution(val.HeadingTrue, 16, 0.0001)
	w.writeUnsignedResolution(val.HeadingMagnetic, 16, 0.0001)
	w.skipBits(8)
	return w.Bytes(), w.Err()
}

func encodeSeatalkPilotHeadingMsg(v Message) ([]byte, error) {
	val, ok := v.(*SeatalkPilotHeading)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkPilotHeading, got %T", v)
	}
	return EncodeSeatalkPilotHeading(val)
}

type SeatalkPilotLockedHeading struct {
	Info                  MessageInfo           `json:"info"`
	ManufacturerCode      ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode          IndustryCodeConst     `json:"industryCode"`
	Sid                   []uint8               `json:"sid"`
	TargetHeadingTrue     *float32              `json:"targetHeadingTrue"`
	TargetHeadingMagnetic *float32              `json:"targetHeadingMagnetic"`
}

func (s *SeatalkPilotLockedHeading) PGNNumber() uint32 { return 65360 }

func DecodeSeatalkPilotLockedHeading(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SeatalkPilotLockedHeading{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotLockedHeading-ManufacturerCode: %w", err)
	} else {
		if v != 1851 {
			return nil, fmt.Errorf("match failed for SeatalkPilotLockedHeading-ManufacturerCode: Expected %d != %d", 1851, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotLockedHeading-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SeatalkPilotLockedHeading-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotLockedHeading-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotLockedHeading-TargetHeadingTrue: %w", err)
	} else {
		val.TargetHeadingTrue = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotLockedHeading-TargetHeadingMagnetic: %w", err)
	} else {
		val.TargetHeadingMagnetic = v

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

func EncodeSeatalkPilotLockedHeading(val *SeatalkPilotLockedHeading) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeBinaryData(val.Sid, 8)
	w.writeUnsignedResolution(val.TargetHeadingTrue, 16, 0.0001)
	w.writeUnsignedResolution(val.TargetHeadingMagnetic, 16, 0.0001)
	w.skipBits(8)
	return w.Bytes(), w.Err()
}

func encodeSeatalkPilotLockedHeadingMsg(v Message) ([]byte, error) {
	val, ok := v.(*SeatalkPilotLockedHeading)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkPilotLockedHeading, got %T", v)
	}
	return EncodeSeatalkPilotLockedHeading(val)
}

type SeatalkSilenceAlarm struct {
	Info             MessageInfo            `json:"info"`
	ManufacturerCode ManufacturerCodeConst  `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst      `json:"industryCode"`
	AlarmId          SeatalkAlarmIdConst    `json:"alarmId"`
	AlarmGroup       SeatalkAlarmGroupConst `json:"alarmGroup"`
}

func (s *SeatalkSilenceAlarm) PGNNumber() uint32 { return 65361 }

func DecodeSeatalkSilenceAlarm(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SeatalkSilenceAlarm{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkSilenceAlarm-ManufacturerCode: %w", err)
	} else {
		if v != 1851 {
			return nil, fmt.Errorf("match failed for SeatalkSilenceAlarm-ManufacturerCode: Expected %d != %d", 1851, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkSilenceAlarm-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SeatalkSilenceAlarm-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkSilenceAlarm-AlarmId: %w", err)
	} else {
		val.AlarmId = SeatalkAlarmIdConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkSilenceAlarm-AlarmGroup: %w", err)
	} else {
		val.AlarmGroup = SeatalkAlarmGroupConst(v)

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

func EncodeSeatalkSilenceAlarm(val *SeatalkSilenceAlarm) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.AlarmId), 8)
	w.writeLookupField(uint64(val.AlarmGroup), 8)
	w.skipBits(32)
	return w.Bytes(), w.Err()
}

func encodeSeatalkSilenceAlarmMsg(v Message) ([]byte, error) {
	val, ok := v.(*SeatalkSilenceAlarm)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkSilenceAlarm, got %T", v)
	}
	return EncodeSeatalkSilenceAlarm(val)
}

type SeatalkKeypadMessage struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    *uint8                `json:"proprietaryId"`
	FirstKey         *uint8                `json:"firstKey"`
	SecondKey        *uint8                `json:"secondKey"`
	FirstKeyState    *uint8                `json:"firstKeyState"`
	SecondKeyState   *uint8                `json:"secondKeyState"`
	EncoderPosition  *uint8                `json:"encoderPosition"`
}

func (s *SeatalkKeypadMessage) PGNNumber() uint32 { return 65371 }

func DecodeSeatalkKeypadMessage(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SeatalkKeypadMessage{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkKeypadMessage-ManufacturerCode: %w", err)
	} else {
		if v != 1851 {
			return nil, fmt.Errorf("match failed for SeatalkKeypadMessage-ManufacturerCode: Expected %d != %d", 1851, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkKeypadMessage-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SeatalkKeypadMessage-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkKeypadMessage-ProprietaryId: %w", err)
	} else {
		val.ProprietaryId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkKeypadMessage-FirstKey: %w", err)
	} else {
		val.FirstKey = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkKeypadMessage-SecondKey: %w", err)
	} else {
		val.SecondKey = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(2); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkKeypadMessage-FirstKeyState: %w", err)
	} else {
		val.FirstKeyState = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(2); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkKeypadMessage-SecondKeyState: %w", err)
	} else {
		val.SecondKeyState = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkKeypadMessage-EncoderPosition: %w", err)
	} else {
		val.EncoderPosition = v

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

func EncodeSeatalkKeypadMessage(val *SeatalkKeypadMessage) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.ProprietaryId, 8)
	w.writeUInt8(val.FirstKey, 8)
	w.writeUInt8(val.SecondKey, 8)
	w.writeUInt8(val.FirstKeyState, 2)
	w.writeUInt8(val.SecondKeyState, 2)
	w.skipBits(4)
	w.writeUInt8(val.EncoderPosition, 8)
	w.skipBits(8)
	return w.Bytes(), w.Err()
}

func encodeSeatalkKeypadMessageMsg(v Message) ([]byte, error) {
	val, ok := v.(*SeatalkKeypadMessage)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkKeypadMessage, got %T", v)
	}
	return EncodeSeatalkKeypadMessage(val)
}

type SeatalkKeypadHeartbeat struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    *uint8                `json:"proprietaryId"`
	Variant          *uint8                `json:"variant"`
	Status           *uint8                `json:"status"`
}

func (s *SeatalkKeypadHeartbeat) PGNNumber() uint32 { return 65374 }

func DecodeSeatalkKeypadHeartbeat(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SeatalkKeypadHeartbeat{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkKeypadHeartbeat-ManufacturerCode: %w", err)
	} else {
		if v != 1851 {
			return nil, fmt.Errorf("match failed for SeatalkKeypadHeartbeat-ManufacturerCode: Expected %d != %d", 1851, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkKeypadHeartbeat-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SeatalkKeypadHeartbeat-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkKeypadHeartbeat-ProprietaryId: %w", err)
	} else {
		val.ProprietaryId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkKeypadHeartbeat-Variant: %w", err)
	} else {
		val.Variant = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkKeypadHeartbeat-Status: %w", err)
	} else {
		val.Status = v

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

func EncodeSeatalkKeypadHeartbeat(val *SeatalkKeypadHeartbeat) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.ProprietaryId, 8)
	w.writeUInt8(val.Variant, 8)
	w.writeUInt8(val.Status, 8)
	w.skipBits(24)
	return w.Bytes(), w.Err()
}

func encodeSeatalkKeypadHeartbeatMsg(v Message) ([]byte, error) {
	val, ok := v.(*SeatalkKeypadHeartbeat)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkKeypadHeartbeat, got %T", v)
	}
	return EncodeSeatalkKeypadHeartbeat(val)
}

type SeatalkPilotMode struct {
	Info             MessageInfo             `json:"info"`
	ManufacturerCode ManufacturerCodeConst   `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst       `json:"industryCode"`
	PilotMode        SeatalkPilotMode16Const `json:"pilotMode"`
	SubMode          []uint8                 `json:"subMode"`
	PilotModeData    []uint8                 `json:"pilotModeData"`
}

func (s *SeatalkPilotMode) PGNNumber() uint32 { return 65379 }

func DecodeSeatalkPilotMode(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SeatalkPilotMode{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotMode-ManufacturerCode: %w", err)
	} else {
		if v != 1851 {
			return nil, fmt.Errorf("match failed for SeatalkPilotMode-ManufacturerCode: Expected %d != %d", 1851, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotMode-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SeatalkPilotMode-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(16); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotMode-PilotMode: %w", err)
	} else {
		val.PilotMode = SeatalkPilotMode16Const(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(16); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotMode-SubMode: %w", err)
	} else {
		val.SubMode = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(8); err != nil {
		return nil, fmt.Errorf("parse failed for SeatalkPilotMode-PilotModeData: %w", err)
	} else {
		val.PilotModeData = v

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

func EncodeSeatalkPilotMode(val *SeatalkPilotMode) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.PilotMode), 16)
	w.writeBinaryData(val.SubMode, 16)
	w.writeBinaryData(val.PilotModeData, 8)
	w.skipBits(8)
	return w.Bytes(), w.Err()
}

func encodeSeatalkPilotModeMsg(v Message) ([]byte, error) {
	val, ok := v.(*SeatalkPilotMode)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkPilotMode, got %T", v)
	}
	return EncodeSeatalkPilotMode(val)
}

type Seatalk1DeviceIdentification struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    *uint16               `json:"proprietaryId"`
	Command          *uint8                `json:"command"`
	Device           SeatalkDeviceIdConst  `json:"device"`
}

func (s *Seatalk1DeviceIdentification) PGNNumber() uint32 { return 126720 }

func DecodeSeatalk1DeviceIdentification(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Seatalk1DeviceIdentification{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DeviceIdentification-ManufacturerCode: %w", err)
	} else {
		if v != 1851 {
			return nil, fmt.Errorf("match failed for Seatalk1DeviceIdentification-ManufacturerCode: Expected %d != %d", 1851, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DeviceIdentification-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for Seatalk1DeviceIdentification-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DeviceIdentification-ProprietaryId: %w", err)
	} else {
		if v != nil && *v != 33264 {
			return nil, fmt.Errorf("match failed for Seatalk1DeviceIdentification-ProprietaryId: Expected %d != %d", 33264, *v)
		}
		val.ProprietaryId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DeviceIdentification-Command: %w", err)
	} else {
		if v != nil && *v != 144 {
			return nil, fmt.Errorf("match failed for Seatalk1DeviceIdentification-Command: Expected %d != %d", 144, *v)
		}
		val.Command = v

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DeviceIdentification-Device: %w", err)
	} else {
		val.Device = SeatalkDeviceIdConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSeatalk1DeviceIdentification(val *Seatalk1DeviceIdentification) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt16(val.ProprietaryId, 16)
	w.writeUInt8(val.Command, 8)
	w.skipBits(8)
	w.writeLookupField(uint64(val.Device), 8)
	return w.Bytes(), w.Err()
}

func encodeSeatalk1DeviceIdentificationMsg(v Message) ([]byte, error) {
	val, ok := v.(*Seatalk1DeviceIdentification)
	if !ok {
		return nil, fmt.Errorf("expected *Seatalk1DeviceIdentification, got %T", v)
	}
	return EncodeSeatalk1DeviceIdentification(val)
}

type Seatalk1DisplayBrightness struct {
	Info             MessageInfo              `json:"info"`
	ManufacturerCode ManufacturerCodeConst    `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst        `json:"industryCode"`
	ProprietaryId    *uint16                  `json:"proprietaryId"`
	Group            SeatalkNetworkGroupConst `json:"group"`
	Unknown1         []uint8                  `json:"unknown1"`
	Command          *uint8                   `json:"command"`
	Brightness       *uint8                   `json:"brightness"`
	Unknown2         []uint8                  `json:"unknown2"`
}

func (s *Seatalk1DisplayBrightness) PGNNumber() uint32 { return 126720 }

func DecodeSeatalk1DisplayBrightness(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Seatalk1DisplayBrightness{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DisplayBrightness-ManufacturerCode: %w", err)
	} else {
		if v != 1851 {
			return nil, fmt.Errorf("match failed for Seatalk1DisplayBrightness-ManufacturerCode: Expected %d != %d", 1851, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DisplayBrightness-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for Seatalk1DisplayBrightness-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DisplayBrightness-ProprietaryId: %w", err)
	} else {
		if v != nil && *v != 3212 {
			return nil, fmt.Errorf("match failed for Seatalk1DisplayBrightness-ProprietaryId: Expected %d != %d", 3212, *v)
		}
		val.ProprietaryId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DisplayBrightness-Group: %w", err)
	} else {
		val.Group = SeatalkNetworkGroupConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DisplayBrightness-Unknown1: %w", err)
	} else {
		val.Unknown1 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DisplayBrightness-Command: %w", err)
	} else {
		if v != nil && *v != 0 {
			return nil, fmt.Errorf("match failed for Seatalk1DisplayBrightness-Command: Expected %d != %d", 0, *v)
		}
		val.Command = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DisplayBrightness-Brightness: %w", err)
	} else {
		val.Brightness = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DisplayBrightness-Unknown2: %w", err)
	} else {
		val.Unknown2 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSeatalk1DisplayBrightness(val *Seatalk1DisplayBrightness) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt16(val.ProprietaryId, 16)
	w.writeLookupField(uint64(val.Group), 8)
	w.writeBinaryData(val.Unknown1, 8)
	w.writeUInt8(val.Command, 8)
	w.writeUInt8(val.Brightness, 8)
	w.writeBinaryData(val.Unknown2, 8)
	return w.Bytes(), w.Err()
}

func encodeSeatalk1DisplayBrightnessMsg(v Message) ([]byte, error) {
	val, ok := v.(*Seatalk1DisplayBrightness)
	if !ok {
		return nil, fmt.Errorf("expected *Seatalk1DisplayBrightness, got %T", v)
	}
	return EncodeSeatalk1DisplayBrightness(val)
}

type Seatalk1DisplayColor struct {
	Info             MessageInfo              `json:"info"`
	ManufacturerCode ManufacturerCodeConst    `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst        `json:"industryCode"`
	ProprietaryId    *uint16                  `json:"proprietaryId"`
	Group            SeatalkNetworkGroupConst `json:"group"`
	Unknown1         []uint8                  `json:"unknown1"`
	Command          *uint8                   `json:"command"`
	Color            SeatalkDisplayColorConst `json:"color"`
	Unknown2         []uint8                  `json:"unknown2"`
}

func (s *Seatalk1DisplayColor) PGNNumber() uint32 { return 126720 }

func DecodeSeatalk1DisplayColor(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Seatalk1DisplayColor{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DisplayColor-ManufacturerCode: %w", err)
	} else {
		if v != 1851 {
			return nil, fmt.Errorf("match failed for Seatalk1DisplayColor-ManufacturerCode: Expected %d != %d", 1851, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DisplayColor-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for Seatalk1DisplayColor-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DisplayColor-ProprietaryId: %w", err)
	} else {
		if v != nil && *v != 3212 {
			return nil, fmt.Errorf("match failed for Seatalk1DisplayColor-ProprietaryId: Expected %d != %d", 3212, *v)
		}
		val.ProprietaryId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DisplayColor-Group: %w", err)
	} else {
		val.Group = SeatalkNetworkGroupConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DisplayColor-Unknown1: %w", err)
	} else {
		val.Unknown1 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DisplayColor-Command: %w", err)
	} else {
		if v != nil && *v != 1 {
			return nil, fmt.Errorf("match failed for Seatalk1DisplayColor-Command: Expected %d != %d", 1, *v)
		}
		val.Command = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DisplayColor-Color: %w", err)
	} else {
		val.Color = SeatalkDisplayColorConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1DisplayColor-Unknown2: %w", err)
	} else {
		val.Unknown2 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSeatalk1DisplayColor(val *Seatalk1DisplayColor) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt16(val.ProprietaryId, 16)
	w.writeLookupField(uint64(val.Group), 8)
	w.writeBinaryData(val.Unknown1, 8)
	w.writeUInt8(val.Command, 8)
	w.writeLookupField(uint64(val.Color), 8)
	w.writeBinaryData(val.Unknown2, 8)
	return w.Bytes(), w.Err()
}

func encodeSeatalk1DisplayColorMsg(v Message) ([]byte, error) {
	val, ok := v.(*Seatalk1DisplayColor)
	if !ok {
		return nil, fmt.Errorf("expected *Seatalk1DisplayColor, got %T", v)
	}
	return EncodeSeatalk1DisplayColor(val)
}

type Seatalk1Keystroke struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    *uint16               `json:"proprietaryId"`
	Command          *uint8                `json:"command"`
	Device           *uint8                `json:"device"`
	Key              SeatalkKeystrokeConst `json:"key"`
	Keyinverted      *uint8                `json:"keyinverted"`
	UnknownData      []uint8               `json:"unknownData"`
}

func (s *Seatalk1Keystroke) PGNNumber() uint32 { return 126720 }

func DecodeSeatalk1Keystroke(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Seatalk1Keystroke{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1Keystroke-ManufacturerCode: %w", err)
	} else {
		if v != 1851 {
			return nil, fmt.Errorf("match failed for Seatalk1Keystroke-ManufacturerCode: Expected %d != %d", 1851, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1Keystroke-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for Seatalk1Keystroke-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1Keystroke-ProprietaryId: %w", err)
	} else {
		if v != nil && *v != 33264 {
			return nil, fmt.Errorf("match failed for Seatalk1Keystroke-ProprietaryId: Expected %d != %d", 33264, *v)
		}
		val.ProprietaryId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1Keystroke-Command: %w", err)
	} else {
		if v != nil && *v != 134 {
			return nil, fmt.Errorf("match failed for Seatalk1Keystroke-Command: Expected %d != %d", 134, *v)
		}
		val.Command = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1Keystroke-Device: %w", err)
	} else {
		val.Device = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1Keystroke-Key: %w", err)
	} else {
		val.Key = SeatalkKeystrokeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1Keystroke-Keyinverted: %w", err)
	} else {
		val.Keyinverted = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(112); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1Keystroke-UnknownData: %w", err)
	} else {
		val.UnknownData = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSeatalk1Keystroke(val *Seatalk1Keystroke) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt16(val.ProprietaryId, 16)
	w.writeUInt8(val.Command, 8)
	w.writeUInt8(val.Device, 8)
	w.writeLookupField(uint64(val.Key), 8)
	w.writeUInt8(val.Keyinverted, 8)
	w.writeBinaryData(val.UnknownData, 112)
	return w.Bytes(), w.Err()
}

func encodeSeatalk1KeystrokeMsg(v Message) ([]byte, error) {
	val, ok := v.(*Seatalk1Keystroke)
	if !ok {
		return nil, fmt.Errorf("expected *Seatalk1Keystroke, got %T", v)
	}
	return EncodeSeatalk1Keystroke(val)
}

type Seatalk1PilotMode struct {
	Info             MessageInfo           `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode     IndustryCodeConst     `json:"industryCode"`
	ProprietaryId    *uint16               `json:"proprietaryId"`
	Command          *uint8                `json:"command"`
	Unknown1         []uint8               `json:"unknown1"`
	PilotMode        SeatalkPilotModeConst `json:"pilotMode"`
	SubMode          *uint8                `json:"subMode"`
	PilotModeData    []uint8               `json:"pilotModeData"`
	Unknown2         []uint8               `json:"unknown2"`
}

func (s *Seatalk1PilotMode) PGNNumber() uint32 { return 126720 }

func DecodeSeatalk1PilotMode(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Seatalk1PilotMode{}
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1PilotMode-ManufacturerCode: %w", err)
	} else {
		if v != 1851 {
			return nil, fmt.Errorf("match failed for Seatalk1PilotMode-ManufacturerCode: Expected %d != %d", 1851, v)
		}
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1PilotMode-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for Seatalk1PilotMode-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1PilotMode-ProprietaryId: %w", err)
	} else {
		if v != nil && *v != 33264 {
			return nil, fmt.Errorf("match failed for Seatalk1PilotMode-ProprietaryId: Expected %d != %d", 33264, *v)
		}
		val.ProprietaryId = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1PilotMode-Command: %w", err)
	} else {
		if v != nil && *v != 132 {
			return nil, fmt.Errorf("match failed for Seatalk1PilotMode-Command: Expected %d != %d", 132, *v)
		}
		val.Command = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(24); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1PilotMode-Unknown1: %w", err)
	} else {
		val.Unknown1 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1PilotMode-PilotMode: %w", err)
	} else {
		val.PilotMode = SeatalkPilotModeConst(v)

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1PilotMode-SubMode: %w", err)
	} else {
		val.SubMode = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(8); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1PilotMode-PilotModeData: %w", err)
	} else {
		val.PilotModeData = v

		if stream.isEOF() {
			return val, nil
		}
	}
	if v, err := stream.readBinaryData(80); err != nil {
		return nil, fmt.Errorf("parse failed for Seatalk1PilotMode-Unknown2: %w", err)
	} else {
		val.Unknown2 = v

		if stream.isEOF() {
			return val, nil
		}
	}
	return val, nil
}

func EncodeSeatalk1PilotMode(val *Seatalk1PilotMode) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt16(val.ProprietaryId, 16)
	w.writeUInt8(val.Command, 8)
	w.writeBinaryData(val.Unknown1, 24)
	w.writeLookupField(uint64(val.PilotMode), 8)
	w.writeUInt8(val.SubMode, 8)
	w.writeBinaryData(val.PilotModeData, 8)
	w.writeBinaryData(val.Unknown2, 80)
	return w.Bytes(), w.Err()
}

func encodeSeatalk1PilotModeMsg(v Message) ([]byte, error) {
	val, ok := v.(*Seatalk1PilotMode)
	if !ok {
		return nil, fmt.Errorf("expected *Seatalk1PilotMode, got %T", v)
	}
	return EncodeSeatalk1PilotMode(val)
}
