package pgn

import (
	"fmt"

	"github.com/open-ships/n2k/units"
)

type SeatalkWirelessKeypadLightControl struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId *uint8 `json:"proprietaryId"`
	Variant *uint8 `json:"variant"`
	WirelessSetting *uint8 `json:"wirelessSetting"`
	WiredSetting *uint8 `json:"wiredSetting"`
}
func DecodeSeatalkWirelessKeypadLightControl(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SeatalkWirelessKeypadLightControl
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
type SeatalkWirelessKeypadControl struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Pid *uint8 `json:"pid"`
	Variant *uint8 `json:"variant"`
	BeepControl *uint8 `json:"beepControl"`
}
func DecodeSeatalkWirelessKeypadControl(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SeatalkWirelessKeypadControl
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
type VictronBatteryRegister struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	RegisterId *uint16 `json:"registerId"`
	Payload *uint32 `json:"payload"`
}
func DecodeVictronBatteryRegister(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val VictronBatteryRegister
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for VictronBatteryRegister-ManufacturerCode: %w", err)
	} else {
		if v != 358 {
			return nil, fmt.Errorf("match failed for VictronBatteryRegister-ManufacturerCode: Expected %d != %d", 358, v)
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
		return nil, fmt.Errorf("parse failed for VictronBatteryRegister-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for VictronBatteryRegister-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for VictronBatteryRegister-RegisterId: %w", err)
	} else {
		val.RegisterId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for VictronBatteryRegister-Payload: %w", err)
	} else {
		val.Payload = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FurunoHeave struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Heave *units.Distance `json:"heave"`
}
func DecodeFurunoHeave(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FurunoHeave
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoHeave-ManufacturerCode: %w", err)
	} else {
		if v != 1855 {
			return nil, fmt.Errorf("match failed for FurunoHeave-ManufacturerCode: Expected %d != %d", 1855, v)
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
		return nil, fmt.Errorf("parse failed for FurunoHeave-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FurunoHeave-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(32, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoHeave-Heave: %w", err)
	} else {
		val.Heave = nullableUnit(units.Meter, v, units.NewDistance)

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
type MaretronProprietaryDcBreakerCurrent struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	BankInstance *uint8 `json:"bankInstance"`
	IndicatorNumber *uint8 `json:"indicatorNumber"`
	BreakerCurrent *float32 `json:"breakerCurrent"`
}
func DecodeMaretronProprietaryDcBreakerCurrent(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val MaretronProprietaryDcBreakerCurrent
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryDcBreakerCurrent-ManufacturerCode: %w", err)
	} else {
		if v != 137 {
			return nil, fmt.Errorf("match failed for MaretronProprietaryDcBreakerCurrent-ManufacturerCode: Expected %d != %d", 137, v)
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
		return nil, fmt.Errorf("parse failed for MaretronProprietaryDcBreakerCurrent-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for MaretronProprietaryDcBreakerCurrent-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryDcBreakerCurrent-BankInstance: %w", err)
	} else {
		val.BankInstance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryDcBreakerCurrent-IndicatorNumber: %w", err)
	} else {
		val.IndicatorNumber = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryDcBreakerCurrent-BreakerCurrent: %w", err)
	} else {
		val.BreakerCurrent = v

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
type AirmarBootStateAcknowledgment struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	BootState BootStateConst `json:"bootState"`
}
func DecodeAirmarBootStateAcknowledgment(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarBootStateAcknowledgment
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarBootStateAcknowledgment-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarBootStateAcknowledgment-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarBootStateAcknowledgment-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarBootStateAcknowledgment-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarBootStateAcknowledgment-BootState: %w", err)
	} else {
		val.BootState = BootStateConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(45)
	if stream.isEOF() {
		return val, nil
		}	
	return val, nil
}
type LowranceTemperature struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	TemperatureSource TemperatureSourceConst `json:"temperatureSource"`
	ActualTemperature *units.Temperature `json:"actualTemperature"`
}
func DecodeLowranceTemperature(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val LowranceTemperature
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceTemperature-ManufacturerCode: %w", err)
	} else {
		if v != 140 {
			return nil, fmt.Errorf("match failed for LowranceTemperature-ManufacturerCode: Expected %d != %d", 140, v)
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
		return nil, fmt.Errorf("parse failed for LowranceTemperature-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for LowranceTemperature-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceTemperature-TemperatureSource: %w", err)
	} else {
		val.TemperatureSource = TemperatureSourceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceTemperature-ActualTemperature: %w", err)
	} else {
		val.ActualTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

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
type AirmarBootStateRequest struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
}
func DecodeAirmarBootStateRequest(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarBootStateRequest
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarBootStateRequest-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarBootStateRequest-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarBootStateRequest-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarBootStateRequest-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

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
type AirmarAccessLevel struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	FormatCode *uint8 `json:"formatCode"`
	AccessLevel AccessLevelConst `json:"accessLevel"`
	AccessSeedKey *uint32 `json:"accessSeedKey"`
}
func DecodeAirmarAccessLevel(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarAccessLevel
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAccessLevel-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarAccessLevel-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarAccessLevel-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarAccessLevel-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAccessLevel-FormatCode: %w", err)
	} else {
		val.FormatCode = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAccessLevel-AccessLevel: %w", err)
	} else {
		val.AccessLevel = AccessLevelConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(5)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAccessLevel-AccessSeedKey: %w", err)
	} else {
		val.AccessSeedKey = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SimnetConfigureTemperatureSensor struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
}
func DecodeSimnetConfigureTemperatureSensor(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetConfigureTemperatureSensor
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetConfigureTemperatureSensor-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetConfigureTemperatureSensor-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetConfigureTemperatureSensor-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetConfigureTemperatureSensor-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

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
type SeatalkAlarm struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Sid []uint8 `json:"sid"`
	AlarmStatus SeatalkAlarmStatusConst `json:"alarmStatus"`
	AlarmId SeatalkAlarmIdConst `json:"alarmId"`
	AlarmGroup SeatalkAlarmGroupConst `json:"alarmGroup"`
	AlarmPriority []uint8 `json:"alarmPriority"`
}
func DecodeSeatalkAlarm(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SeatalkAlarm
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
type SimnetTrimTabSensorCalibration struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
}
func DecodeSimnetTrimTabSensorCalibration(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetTrimTabSensorCalibration
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetTrimTabSensorCalibration-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetTrimTabSensorCalibration-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetTrimTabSensorCalibration-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetTrimTabSensorCalibration-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

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
type SimnetPaddleWheelSpeedConfiguration struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
}
func DecodeSimnetPaddleWheelSpeedConfiguration(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetPaddleWheelSpeedConfiguration
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetPaddleWheelSpeedConfiguration-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetPaddleWheelSpeedConfiguration-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetPaddleWheelSpeedConfiguration-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetPaddleWheelSpeedConfiguration-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

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
type SimnetClearFluidLevelWarnings struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
}
func DecodeSimnetClearFluidLevelWarnings(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetClearFluidLevelWarnings
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetClearFluidLevelWarnings-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetClearFluidLevelWarnings-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetClearFluidLevelWarnings-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetClearFluidLevelWarnings-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

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
type SimnetLgc2000Configuration struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
}
func DecodeSimnetLgc2000Configuration(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetLgc2000Configuration
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetLgc2000Configuration-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetLgc2000Configuration-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetLgc2000Configuration-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetLgc2000Configuration-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

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
type DiverseYachtServicesLoadCell struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Instance *uint8 `json:"instance"`
	LoadCell *uint32 `json:"loadCell"`
}
func DecodeDiverseYachtServicesLoadCell(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val DiverseYachtServicesLoadCell
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for DiverseYachtServicesLoadCell-ManufacturerCode: %w", err)
	} else {
		if v != 641 {
			return nil, fmt.Errorf("match failed for DiverseYachtServicesLoadCell-ManufacturerCode: Expected %d != %d", 641, v)
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
		return nil, fmt.Errorf("parse failed for DiverseYachtServicesLoadCell-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for DiverseYachtServicesLoadCell-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for DiverseYachtServicesLoadCell-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for DiverseYachtServicesLoadCell-LoadCell: %w", err)
	} else {
		val.LoadCell = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SimnetApUnknown1 struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	A *uint8 `json:"a"`
	B *uint8 `json:"b"`
	C *uint16 `json:"c"`
	D *uint8 `json:"d"`
}
func DecodeSimnetApUnknown1(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetApUnknown1
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown1-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetApUnknown1-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetApUnknown1-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetApUnknown1-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown1-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown1-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown1-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown1-D: %w", err)
	} else {
		val.D = v

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
type SimnetDeviceStatus struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Model SimnetDeviceModelConst `json:"model"`
	Report SimnetDeviceReportConst `json:"report"`
	Status SimnetApStatusConst `json:"status"`
}
func DecodeSimnetDeviceStatus(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetDeviceStatus
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatus-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetDeviceStatus-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatus-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetDeviceStatus-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatus-Model: %w", err)
	} else {
		val.Model = SimnetDeviceModelConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatus-Report: %w", err)
	} else {
		if v != 2 {
			return nil, fmt.Errorf("match failed for SimnetDeviceStatus-Report: Expected %d != %d", 2, v)
		}
		val.Report = SimnetDeviceReportConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatus-Status: %w", err)
	} else {
		val.Status = SimnetApStatusConst(v)

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
type SimnetDeviceStatusRequest struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Model SimnetDeviceModelConst `json:"model"`
	Report SimnetDeviceReportConst `json:"report"`
}
func DecodeSimnetDeviceStatusRequest(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetDeviceStatusRequest
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatusRequest-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetDeviceStatusRequest-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatusRequest-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetDeviceStatusRequest-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatusRequest-Model: %w", err)
	} else {
		val.Model = SimnetDeviceModelConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceStatusRequest-Report: %w", err)
	} else {
		if v != 3 {
			return nil, fmt.Errorf("match failed for SimnetDeviceStatusRequest-Report: Expected %d != %d", 3, v)
		}
		val.Report = SimnetDeviceReportConst(v)

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
type SimnetPilotMode struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Model SimnetDeviceModelConst `json:"model"`
	Report SimnetDeviceReportConst `json:"report"`
	Mode SimnetApModeBitfieldConst `json:"mode"`
}
func DecodeSimnetPilotMode(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetPilotMode
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetPilotMode-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetPilotMode-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetPilotMode-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetPilotMode-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetPilotMode-Model: %w", err)
	} else {
		val.Model = SimnetDeviceModelConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetPilotMode-Report: %w", err)
	} else {
		if v != 10 {
			return nil, fmt.Errorf("match failed for SimnetPilotMode-Report: Expected %d != %d", 10, v)
		}
		val.Report = SimnetDeviceReportConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetPilotMode-Mode: %w", err)
	} else {
		val.Mode = SimnetApModeBitfieldConst(v)

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
type SimnetDeviceModeRequest struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Model SimnetDeviceModelConst `json:"model"`
	Report SimnetDeviceReportConst `json:"report"`
}
func DecodeSimnetDeviceModeRequest(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetDeviceModeRequest
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceModeRequest-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetDeviceModeRequest-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetDeviceModeRequest-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetDeviceModeRequest-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceModeRequest-Model: %w", err)
	} else {
		val.Model = SimnetDeviceModelConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetDeviceModeRequest-Report: %w", err)
	} else {
		if v != 11 {
			return nil, fmt.Errorf("match failed for SimnetDeviceModeRequest-Report: Expected %d != %d", 11, v)
		}
		val.Report = SimnetDeviceReportConst(v)

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
type SimnetSailingProcessorStatus struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Model SimnetDeviceModelConst `json:"model"`
	Report SimnetDeviceReportConst `json:"report"`
	Data []uint8 `json:"data"`
}
func DecodeSimnetSailingProcessorStatus(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetSailingProcessorStatus
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetSailingProcessorStatus-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetSailingProcessorStatus-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetSailingProcessorStatus-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetSailingProcessorStatus-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetSailingProcessorStatus-Model: %w", err)
	} else {
		val.Model = SimnetDeviceModelConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetSailingProcessorStatus-Report: %w", err)
	} else {
		if v != 23 {
			return nil, fmt.Errorf("match failed for SimnetSailingProcessorStatus-Report: Expected %d != %d", 23, v)
		}
		val.Report = SimnetDeviceReportConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readBinaryData(32); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetSailingProcessorStatus-Data: %w", err)
	} else {
		val.Data = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type NavicoWirelessBatteryStatus struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Status *uint8 `json:"status"`
	BatteryStatus *uint8 `json:"batteryStatus"`
	BatteryChargeStatus *uint8 `json:"batteryChargeStatus"`
}
func DecodeNavicoWirelessBatteryStatus(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val NavicoWirelessBatteryStatus
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoWirelessBatteryStatus-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for NavicoWirelessBatteryStatus-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for NavicoWirelessBatteryStatus-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for NavicoWirelessBatteryStatus-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoWirelessBatteryStatus-Status: %w", err)
	} else {
		val.Status = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoWirelessBatteryStatus-BatteryStatus: %w", err)
	} else {
		val.BatteryStatus = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoWirelessBatteryStatus-BatteryChargeStatus: %w", err)
	} else {
		val.BatteryChargeStatus = v

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
type NavicoWirelessSignalStatus struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Unknown *uint8 `json:"unknown"`
	SignalStrength *uint8 `json:"signalStrength"`
}
func DecodeNavicoWirelessSignalStatus(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val NavicoWirelessSignalStatus
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoWirelessSignalStatus-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for NavicoWirelessSignalStatus-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for NavicoWirelessSignalStatus-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for NavicoWirelessSignalStatus-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoWirelessSignalStatus-Unknown: %w", err)
	} else {
		val.Unknown = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoWirelessSignalStatus-SignalStrength: %w", err)
	} else {
		val.SignalStrength = v

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
type SimnetApUnknown2 struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	A *uint8 `json:"a"`
	B *uint8 `json:"b"`
	C *uint8 `json:"c"`
	D *uint8 `json:"d"`
	E *uint8 `json:"e"`
}
func DecodeSimnetApUnknown2(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetApUnknown2
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown2-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetApUnknown2-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetApUnknown2-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetApUnknown2-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown2-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown2-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown2-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown2-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown2-E: %w", err)
	} else {
		val.E = v

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
type SimnetAutopilotAngle struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Mode SimnetApModeConst `json:"mode"`
	Angle *float32 `json:"angle"`
}
func DecodeSimnetAutopilotAngle(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetAutopilotAngle
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAutopilotAngle-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetAutopilotAngle-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetAutopilotAngle-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetAutopilotAngle-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(16)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAutopilotAngle-Mode: %w", err)
	} else {
		val.Mode = SimnetApModeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAutopilotAngle-Angle: %w", err)
	} else {
		val.Angle = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SeatalkPilotWindDatum struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	WindDatum *float32 `json:"windDatum"`
	RollingAverageWindAngle *float32 `json:"rollingAverageWindAngle"`
}
func DecodeSeatalkPilotWindDatum(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SeatalkPilotWindDatum
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
type SimnetMagneticField struct {
	Info MessageInfo `json:"info"`
	A *float32 `json:"a"`
	B *uint8 `json:"b"`
	C *float32 `json:"c"`
	D *float32 `json:"d"`
}
func DecodeSimnetMagneticField(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetMagneticField
	val.Info = Info
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetMagneticField-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetMagneticField-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetMagneticField-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetMagneticField-D: %w", err)
	} else {
		val.D = v

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
type SeatalkPilotHeading struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Sid []uint8 `json:"sid"`
	HeadingTrue *float32 `json:"headingTrue"`
	HeadingMagnetic *float32 `json:"headingMagnetic"`
}
func DecodeSeatalkPilotHeading(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SeatalkPilotHeading
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
type SeatalkPilotLockedHeading struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Sid []uint8 `json:"sid"`
	TargetHeadingTrue *float32 `json:"targetHeadingTrue"`
	TargetHeadingMagnetic *float32 `json:"targetHeadingMagnetic"`
}
func DecodeSeatalkPilotLockedHeading(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SeatalkPilotLockedHeading
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
type SeatalkSilenceAlarm struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	AlarmId SeatalkAlarmIdConst `json:"alarmId"`
	AlarmGroup SeatalkAlarmGroupConst `json:"alarmGroup"`
}
func DecodeSeatalkSilenceAlarm(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SeatalkSilenceAlarm
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
type SeatalkKeypadMessage struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId *uint8 `json:"proprietaryId"`
	FirstKey *uint8 `json:"firstKey"`
	SecondKey *uint8 `json:"secondKey"`
	FirstKeyState *uint8 `json:"firstKeyState"`
	SecondKeyState *uint8 `json:"secondKeyState"`
	EncoderPosition *uint8 `json:"encoderPosition"`
}
func DecodeSeatalkKeypadMessage(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SeatalkKeypadMessage
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
type SeatalkKeypadHeartbeat struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId *uint8 `json:"proprietaryId"`
	Variant *uint8 `json:"variant"`
	Status *uint8 `json:"status"`
}
func DecodeSeatalkKeypadHeartbeat(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SeatalkKeypadHeartbeat
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
type SeatalkPilotMode struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	PilotMode SeatalkPilotMode16Const `json:"pilotMode"`
	SubMode []uint8 `json:"subMode"`
	PilotModeData []uint8 `json:"pilotModeData"`
}
func DecodeSeatalkPilotMode(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SeatalkPilotMode
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
type AirmarDepthQualityFactor struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Sid *uint8 `json:"sid"`
	DepthQualityFactor AirmarDepthQualityFactorConst `json:"depthQualityFactor"`
}
func DecodeAirmarDepthQualityFactor(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarDepthQualityFactor
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarDepthQualityFactor-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarDepthQualityFactor-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarDepthQualityFactor-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarDepthQualityFactor-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarDepthQualityFactor-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarDepthQualityFactor-DepthQualityFactor: %w", err)
	} else {
		val.DepthQualityFactor = AirmarDepthQualityFactorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(36)
	if stream.isEOF() {
		return val, nil
		}	
	return val, nil
}
type AirmarSpeedPulseCount struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Sid *uint8 `json:"sid"`
	DurationOfInterval *float32 `json:"durationOfInterval"`
	NumberOfPulsesReceived *uint16 `json:"numberOfPulsesReceived"`
}
func DecodeAirmarSpeedPulseCount(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarSpeedPulseCount
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedPulseCount-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarSpeedPulseCount-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarSpeedPulseCount-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarSpeedPulseCount-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedPulseCount-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedPulseCount-DurationOfInterval: %w", err)
	} else {
		val.DurationOfInterval = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedPulseCount-NumberOfPulsesReceived: %w", err)
	} else {
		val.NumberOfPulsesReceived = v

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
type AirmarDeviceInformation struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Sid *uint8 `json:"sid"`
	InternalDeviceTemperature *units.Temperature `json:"internalDeviceTemperature"`
	SupplyVoltage *float32 `json:"supplyVoltage"`
}
func DecodeAirmarDeviceInformation(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarDeviceInformation
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarDeviceInformation-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarDeviceInformation-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarDeviceInformation-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarDeviceInformation-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarDeviceInformation-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarDeviceInformation-InternalDeviceTemperature: %w", err)
	} else {
		val.InternalDeviceTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarDeviceInformation-SupplyVoltage: %w", err)
	} else {
		val.SupplyVoltage = v

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
type SimnetApUnknown3 struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	A *uint8 `json:"a"`
	B *uint8 `json:"b"`
	C *uint8 `json:"c"`
	D *uint8 `json:"d"`
	E *uint8 `json:"e"`
}
func DecodeSimnetApUnknown3(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetApUnknown3
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown3-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetApUnknown3-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetApUnknown3-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetApUnknown3-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown3-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown3-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown3-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown3-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown3-E: %w", err)
	} else {
		val.E = v

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
type SimnetAutopilotMode struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
}
func DecodeSimnetAutopilotMode(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetAutopilotMode
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAutopilotMode-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetAutopilotMode-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetAutopilotMode-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetAutopilotMode-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

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
type Seatalk1PilotMode struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId *uint16 `json:"proprietaryId"`
	Command *uint8 `json:"command"`
	Unknown1 []uint8 `json:"unknown1"`
	PilotMode SeatalkPilotModeConst `json:"pilotMode"`
	SubMode *uint8 `json:"subMode"`
	PilotModeData []uint8 `json:"pilotModeData"`
	Unknown2 []uint8 `json:"unknown2"`
}
func DecodeSeatalk1PilotMode(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val Seatalk1PilotMode
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
type Seatalk1Keystroke struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId *uint16 `json:"proprietaryId"`
	Command *uint8 `json:"command"`
	Device *uint8 `json:"device"`
	Key SeatalkKeystrokeConst `json:"key"`
	Keyinverted *uint8 `json:"keyinverted"`
	UnknownData []uint8 `json:"unknownData"`
}
func DecodeSeatalk1Keystroke(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val Seatalk1Keystroke
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
type Seatalk1DeviceIdentification struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId *uint16 `json:"proprietaryId"`
	Command *uint8 `json:"command"`
	Device SeatalkDeviceIdConst `json:"device"`
}
func DecodeSeatalk1DeviceIdentification(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val Seatalk1DeviceIdentification
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
type Seatalk1DisplayBrightness struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId *uint16 `json:"proprietaryId"`
	Group SeatalkNetworkGroupConst `json:"group"`
	Unknown1 []uint8 `json:"unknown1"`
	Command *uint8 `json:"command"`
	Brightness *uint8 `json:"brightness"`
	Unknown2 []uint8 `json:"unknown2"`
}
func DecodeSeatalk1DisplayBrightness(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val Seatalk1DisplayBrightness
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
type Seatalk1DisplayColor struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId *uint16 `json:"proprietaryId"`
	Group SeatalkNetworkGroupConst `json:"group"`
	Unknown1 []uint8 `json:"unknown1"`
	Command *uint8 `json:"command"`
	Color SeatalkDisplayColorConst `json:"color"`
	Unknown2 []uint8 `json:"unknown2"`
}
func DecodeSeatalk1DisplayColor(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val Seatalk1DisplayColor
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
type AirmarAttitudeOffset struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId AirmarCommandConst `json:"proprietaryId"`
	AzimuthOffset *float32 `json:"azimuthOffset"`
	PitchOffset *float32 `json:"pitchOffset"`
	RollOffset *float32 `json:"rollOffset"`
}
func DecodeAirmarAttitudeOffset(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarAttitudeOffset
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAttitudeOffset-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarAttitudeOffset-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarAttitudeOffset-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarAttitudeOffset-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAttitudeOffset-ProprietaryId: %w", err)
	} else {
		if v != 32 {
			return nil, fmt.Errorf("match failed for AirmarAttitudeOffset-ProprietaryId: Expected %d != %d", 32, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAttitudeOffset-AzimuthOffset: %w", err)
	} else {
		val.AzimuthOffset = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAttitudeOffset-PitchOffset: %w", err)
	} else {
		val.PitchOffset = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAttitudeOffset-RollOffset: %w", err)
	} else {
		val.RollOffset = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AirmarCalibrateCompass struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId AirmarCommandConst `json:"proprietaryId"`
	CalibrateFunction AirmarCalibrateFunctionConst `json:"calibrateFunction"`
	CalibrationStatus AirmarCalibrateStatusConst `json:"calibrationStatus"`
	VerifyScore *uint8 `json:"verifyScore"`
	XAxisGainValue *float32 `json:"xAxisGainValue"`
	YAxisGainValue *float32 `json:"yAxisGainValue"`
	ZAxisGainValue *float32 `json:"zAxisGainValue"`
	XAxisLinearOffset *float32 `json:"xAxisLinearOffset"`
	YAxisLinearOffset *float32 `json:"yAxisLinearOffset"`
	ZAxisLinearOffset *float32 `json:"zAxisLinearOffset"`
	XAxisAngularOffset *float32 `json:"xAxisAngularOffset"`
	PitchAndRollDamping *float32 `json:"pitchAndRollDamping"`
	CompassRateGyroDamping *float32 `json:"compassRateGyroDamping"`
}
func DecodeAirmarCalibrateCompass(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarCalibrateCompass
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateCompass-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateCompass-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-ProprietaryId: %w", err)
	} else {
		if v != 33 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateCompass-ProprietaryId: Expected %d != %d", 33, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-CalibrateFunction: %w", err)
	} else {
		val.CalibrateFunction = AirmarCalibrateFunctionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-CalibrationStatus: %w", err)
	} else {
		val.CalibrationStatus = AirmarCalibrateStatusConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-VerifyScore: %w", err)
	} else {
		val.VerifyScore = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-XAxisGainValue: %w", err)
	} else {
		val.XAxisGainValue = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-YAxisGainValue: %w", err)
	} else {
		val.YAxisGainValue = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-ZAxisGainValue: %w", err)
	} else {
		val.ZAxisGainValue = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-XAxisLinearOffset: %w", err)
	} else {
		val.XAxisLinearOffset = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-YAxisLinearOffset: %w", err)
	} else {
		val.YAxisLinearOffset = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-ZAxisLinearOffset: %w", err)
	} else {
		val.ZAxisLinearOffset = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-XAxisAngularOffset: %w", err)
	} else {
		val.XAxisAngularOffset = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.05); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-PitchAndRollDamping: %w", err)
	} else {
		val.PitchAndRollDamping = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.05); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateCompass-CompassRateGyroDamping: %w", err)
	} else {
		val.CompassRateGyroDamping = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AirmarTrueWindOptions struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId AirmarCommandConst `json:"proprietaryId"`
	CogSubstitutionForHdg YesNoConst `json:"cogSubstitutionForHdg"`
}
func DecodeAirmarTrueWindOptions(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarTrueWindOptions
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTrueWindOptions-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarTrueWindOptions-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarTrueWindOptions-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarTrueWindOptions-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTrueWindOptions-ProprietaryId: %w", err)
	} else {
		if v != 34 {
			return nil, fmt.Errorf("match failed for AirmarTrueWindOptions-ProprietaryId: Expected %d != %d", 34, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTrueWindOptions-CogSubstitutionForHdg: %w", err)
	} else {
		val.CogSubstitutionForHdg = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(22)
	if stream.isEOF() {
		return val, nil
		}	
	return val, nil
}
type AirmarSimulateMode struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId AirmarCommandConst `json:"proprietaryId"`
	SimulateMode OffOnConst `json:"simulateMode"`
}
func DecodeAirmarSimulateMode(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarSimulateMode
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSimulateMode-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarSimulateMode-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarSimulateMode-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarSimulateMode-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSimulateMode-ProprietaryId: %w", err)
	} else {
		if v != 35 {
			return nil, fmt.Errorf("match failed for AirmarSimulateMode-ProprietaryId: Expected %d != %d", 35, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSimulateMode-SimulateMode: %w", err)
	} else {
		val.SimulateMode = OffOnConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(22)
	if stream.isEOF() {
		return val, nil
		}	
	return val, nil
}
type AirmarCalibrateDepth struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId AirmarCommandConst `json:"proprietaryId"`
	SpeedOfSoundMode *units.Velocity `json:"speedOfSoundMode"`
}
func DecodeAirmarCalibrateDepth(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarCalibrateDepth
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateDepth-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateDepth-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarCalibrateDepth-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateDepth-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateDepth-ProprietaryId: %w", err)
	} else {
		if v != 40 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateDepth-ProprietaryId: Expected %d != %d", 40, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateDepth-SpeedOfSoundMode: %w", err)
	} else {
		val.SpeedOfSoundMode = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

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
type AirmarCalibrateSpeed struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId AirmarCommandConst `json:"proprietaryId"`
	NumberOfPairsOfDataPoints *uint8 `json:"numberOfPairsOfDataPoints"`
	Repeating1 []AirmarCalibrateSpeedRepeating1 `json:"repeating1"`
}
type AirmarCalibrateSpeedRepeating1 struct {
	InputFrequency *float32 `json:"inputFrequency"`
	OutputSpeed *units.Velocity `json:"outputSpeed"`
}
func DecodeAirmarCalibrateSpeed(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarCalibrateSpeed
	val.Info = Info
		var repeat1Count uint16 = 0
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateSpeed-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateSpeed-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarCalibrateSpeed-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateSpeed-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateSpeed-ProprietaryId: %w", err)
	} else {
		if v != 41 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateSpeed-ProprietaryId: Expected %d != %d", 41, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateSpeed-NumberOfPairsOfDataPoints: %w", err)
	} else {
		val.NumberOfPairsOfDataPoints = v
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
	val.Repeating1 = make([]AirmarCalibrateSpeedRepeating1, 0)
	i := 0 
	for {
		var rep AirmarCalibrateSpeedRepeating1
		if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
			return nil, fmt.Errorf("parse failed for AirmarCalibrateSpeed-InputFrequency: %w", err)
		} else {
			rep.InputFrequency = v
		}
		if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
			return nil, fmt.Errorf("parse failed for AirmarCalibrateSpeed-OutputSpeed: %w", err)
		} else {
			rep.OutputSpeed = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)
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
type AirmarCalibrateTemperature struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId AirmarCommandConst `json:"proprietaryId"`
	TemperatureInstance AirmarTemperatureInstanceConst `json:"temperatureInstance"`
	TemperatureOffset *units.Temperature `json:"temperatureOffset"`
}
func DecodeAirmarCalibrateTemperature(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarCalibrateTemperature
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateTemperature-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateTemperature-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarCalibrateTemperature-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateTemperature-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateTemperature-ProprietaryId: %w", err)
	} else {
		if v != 42 {
			return nil, fmt.Errorf("match failed for AirmarCalibrateTemperature-ProprietaryId: Expected %d != %d", 42, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateTemperature-TemperatureInstance: %w", err)
	} else {
		val.TemperatureInstance = AirmarTemperatureInstanceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readSignedResolution(16, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarCalibrateTemperature-TemperatureOffset: %w", err)
	} else {
		val.TemperatureOffset = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AirmarSpeedFilterNone struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId AirmarCommandConst `json:"proprietaryId"`
	FilterType *uint8 `json:"filterType"`
	SampleInterval *float32 `json:"sampleInterval"`
}
func DecodeAirmarSpeedFilterNone(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarSpeedFilterNone
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterNone-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarSpeedFilterNone-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterNone-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarSpeedFilterNone-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterNone-ProprietaryId: %w", err)
	} else {
		if v != 43 {
			return nil, fmt.Errorf("match failed for AirmarSpeedFilterNone-ProprietaryId: Expected %d != %d", 43, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterNone-FilterType: %w", err)
	} else {
		if v != nil && *v != 0 {
			return nil, fmt.Errorf("match failed for AirmarSpeedFilterNone-FilterType: Expected %d != %d", 0, *v)
		}
		val.FilterType = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterNone-SampleInterval: %w", err)
	} else {
		val.SampleInterval = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AirmarSpeedFilterIir struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId AirmarCommandConst `json:"proprietaryId"`
	FilterType *uint8 `json:"filterType"`
	SampleInterval *float32 `json:"sampleInterval"`
	FilterDuration *float32 `json:"filterDuration"`
}
func DecodeAirmarSpeedFilterIir(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarSpeedFilterIir
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterIir-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarSpeedFilterIir-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterIir-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarSpeedFilterIir-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterIir-ProprietaryId: %w", err)
	} else {
		if v != 43 {
			return nil, fmt.Errorf("match failed for AirmarSpeedFilterIir-ProprietaryId: Expected %d != %d", 43, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterIir-FilterType: %w", err)
	} else {
		if v != nil && *v != 1 {
			return nil, fmt.Errorf("match failed for AirmarSpeedFilterIir-FilterType: Expected %d != %d", 1, *v)
		}
		val.FilterType = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterIir-SampleInterval: %w", err)
	} else {
		val.SampleInterval = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarSpeedFilterIir-FilterDuration: %w", err)
	} else {
		val.FilterDuration = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AirmarTemperatureFilterNone struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId AirmarCommandConst `json:"proprietaryId"`
	FilterType *uint8 `json:"filterType"`
	SampleInterval *float32 `json:"sampleInterval"`
}
func DecodeAirmarTemperatureFilterNone(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarTemperatureFilterNone
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterNone-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarTemperatureFilterNone-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterNone-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarTemperatureFilterNone-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterNone-ProprietaryId: %w", err)
	} else {
		if v != 44 {
			return nil, fmt.Errorf("match failed for AirmarTemperatureFilterNone-ProprietaryId: Expected %d != %d", 44, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterNone-FilterType: %w", err)
	} else {
		if v != nil && *v != 0 {
			return nil, fmt.Errorf("match failed for AirmarTemperatureFilterNone-FilterType: Expected %d != %d", 0, *v)
		}
		val.FilterType = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterNone-SampleInterval: %w", err)
	} else {
		val.SampleInterval = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AirmarTemperatureFilterIir struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId AirmarCommandConst `json:"proprietaryId"`
	FilterType *uint8 `json:"filterType"`
	SampleInterval *float32 `json:"sampleInterval"`
	FilterDuration *float32 `json:"filterDuration"`
}
func DecodeAirmarTemperatureFilterIir(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarTemperatureFilterIir
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterIir-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarTemperatureFilterIir-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterIir-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarTemperatureFilterIir-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterIir-ProprietaryId: %w", err)
	} else {
		if v != 44 {
			return nil, fmt.Errorf("match failed for AirmarTemperatureFilterIir-ProprietaryId: Expected %d != %d", 44, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterIir-FilterType: %w", err)
	} else {
		if v != nil && *v != 1 {
			return nil, fmt.Errorf("match failed for AirmarTemperatureFilterIir-FilterType: Expected %d != %d", 1, *v)
		}
		val.FilterType = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterIir-SampleInterval: %w", err)
	} else {
		val.SampleInterval = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarTemperatureFilterIir-FilterDuration: %w", err)
	} else {
		val.FilterDuration = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type AirmarNmea2000Options struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId AirmarCommandConst `json:"proprietaryId"`
	TransmissionInterval AirmarTransmissionIntervalConst `json:"transmissionInterval"`
}
func DecodeAirmarNmea2000Options(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarNmea2000Options
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarNmea2000Options-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarNmea2000Options-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarNmea2000Options-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarNmea2000Options-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarNmea2000Options-ProprietaryId: %w", err)
	} else {
		if v != 46 {
			return nil, fmt.Errorf("match failed for AirmarNmea2000Options-ProprietaryId: Expected %d != %d", 46, v)
		}
		val.ProprietaryId = AirmarCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarNmea2000Options-TransmissionInterval: %w", err)
	} else {
		val.TransmissionInterval = AirmarTransmissionIntervalConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(22)
	if stream.isEOF() {
		return val, nil
		}	
	return val, nil
}
type AirmarAddressableMultiFrame struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId *uint8 `json:"proprietaryId"`
}
func DecodeAirmarAddressableMultiFrame(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val AirmarAddressableMultiFrame
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAddressableMultiFrame-ManufacturerCode: %w", err)
	} else {
		if v != 135 {
			return nil, fmt.Errorf("match failed for AirmarAddressableMultiFrame-ManufacturerCode: Expected %d != %d", 135, v)
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
		return nil, fmt.Errorf("parse failed for AirmarAddressableMultiFrame-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for AirmarAddressableMultiFrame-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AirmarAddressableMultiFrame-ProprietaryId: %w", err)
	} else {
		val.ProprietaryId = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type MaretronSlaveResponse struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProductCode *uint16 `json:"productCode"`
	SoftwareCode *uint16 `json:"softwareCode"`
	Command *uint8 `json:"command"`
	Status *uint8 `json:"status"`
}
func DecodeMaretronSlaveResponse(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val MaretronSlaveResponse
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSlaveResponse-ManufacturerCode: %w", err)
	} else {
		if v != 137 {
			return nil, fmt.Errorf("match failed for MaretronSlaveResponse-ManufacturerCode: Expected %d != %d", 137, v)
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
		return nil, fmt.Errorf("parse failed for MaretronSlaveResponse-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for MaretronSlaveResponse-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSlaveResponse-ProductCode: %w", err)
	} else {
		val.ProductCode = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSlaveResponse-SoftwareCode: %w", err)
	} else {
		val.SoftwareCode = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSlaveResponse-Command: %w", err)
	} else {
		val.Command = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSlaveResponse-Status: %w", err)
	} else {
		val.Status = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type GarminDayMode struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	UnknownId1 *uint8 `json:"unknownId1"`
	UnknownId2 *uint8 `json:"unknownId2"`
	UnknownId3 *uint8 `json:"unknownId3"`
	UnknownId4 *uint8 `json:"unknownId4"`
	Mode GarminColorModeConst `json:"mode"`
	Backlight GarminBacklightLevelConst `json:"backlight"`
}
func DecodeGarminDayMode(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val GarminDayMode
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for GarminDayMode-ManufacturerCode: %w", err)
	} else {
		if v != 229 {
			return nil, fmt.Errorf("match failed for GarminDayMode-ManufacturerCode: Expected %d != %d", 229, v)
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
		return nil, fmt.Errorf("parse failed for GarminDayMode-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for GarminDayMode-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminDayMode-UnknownId1: %w", err)
	} else {
		if v != nil && *v != 222 {
			return nil, fmt.Errorf("match failed for GarminDayMode-UnknownId1: Expected %d != %d", 222, *v)
		}
		val.UnknownId1 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminDayMode-UnknownId2: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminDayMode-UnknownId2: Expected %d != %d", 5, *v)
		}
		val.UnknownId2 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminDayMode-UnknownId3: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminDayMode-UnknownId3: Expected %d != %d", 5, *v)
		}
		val.UnknownId3 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminDayMode-UnknownId4: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminDayMode-UnknownId4: Expected %d != %d", 5, *v)
		}
		val.UnknownId4 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(16)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminDayMode-Mode: %w", err)
	} else {
		if v != 0 {
			return nil, fmt.Errorf("match failed for GarminDayMode-Mode: Expected %d != %d", 0, v)
		}
		val.Mode = GarminColorModeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminDayMode-Backlight: %w", err)
	} else {
		val.Backlight = GarminBacklightLevelConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type GarminNightMode struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	UnknownId1 *uint8 `json:"unknownId1"`
	UnknownId2 *uint8 `json:"unknownId2"`
	UnknownId3 *uint8 `json:"unknownId3"`
	UnknownId4 *uint8 `json:"unknownId4"`
	Mode GarminColorModeConst `json:"mode"`
	Backlight GarminBacklightLevelConst `json:"backlight"`
}
func DecodeGarminNightMode(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val GarminNightMode
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for GarminNightMode-ManufacturerCode: %w", err)
	} else {
		if v != 229 {
			return nil, fmt.Errorf("match failed for GarminNightMode-ManufacturerCode: Expected %d != %d", 229, v)
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
		return nil, fmt.Errorf("parse failed for GarminNightMode-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for GarminNightMode-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminNightMode-UnknownId1: %w", err)
	} else {
		if v != nil && *v != 222 {
			return nil, fmt.Errorf("match failed for GarminNightMode-UnknownId1: Expected %d != %d", 222, *v)
		}
		val.UnknownId1 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminNightMode-UnknownId2: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminNightMode-UnknownId2: Expected %d != %d", 5, *v)
		}
		val.UnknownId2 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminNightMode-UnknownId3: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminNightMode-UnknownId3: Expected %d != %d", 5, *v)
		}
		val.UnknownId3 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminNightMode-UnknownId4: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminNightMode-UnknownId4: Expected %d != %d", 5, *v)
		}
		val.UnknownId4 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(16)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminNightMode-Mode: %w", err)
	} else {
		if v != 1 {
			return nil, fmt.Errorf("match failed for GarminNightMode-Mode: Expected %d != %d", 1, v)
		}
		val.Mode = GarminColorModeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminNightMode-Backlight: %w", err)
	} else {
		val.Backlight = GarminBacklightLevelConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type GarminColorMode struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	UnknownId1 *uint8 `json:"unknownId1"`
	UnknownId2 *uint8 `json:"unknownId2"`
	UnknownId3 *uint8 `json:"unknownId3"`
	UnknownId4 *uint8 `json:"unknownId4"`
	Mode GarminColorModeConst `json:"mode"`
	Color GarminColorConst `json:"color"`
}
func DecodeGarminColorMode(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val GarminColorMode
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for GarminColorMode-ManufacturerCode: %w", err)
	} else {
		if v != 229 {
			return nil, fmt.Errorf("match failed for GarminColorMode-ManufacturerCode: Expected %d != %d", 229, v)
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
		return nil, fmt.Errorf("parse failed for GarminColorMode-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for GarminColorMode-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminColorMode-UnknownId1: %w", err)
	} else {
		if v != nil && *v != 222 {
			return nil, fmt.Errorf("match failed for GarminColorMode-UnknownId1: Expected %d != %d", 222, *v)
		}
		val.UnknownId1 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminColorMode-UnknownId2: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminColorMode-UnknownId2: Expected %d != %d", 5, *v)
		}
		val.UnknownId2 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminColorMode-UnknownId3: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminColorMode-UnknownId3: Expected %d != %d", 5, *v)
		}
		val.UnknownId3 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminColorMode-UnknownId4: %w", err)
	} else {
		if v != nil && *v != 5 {
			return nil, fmt.Errorf("match failed for GarminColorMode-UnknownId4: Expected %d != %d", 5, *v)
		}
		val.UnknownId4 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(16)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminColorMode-Mode: %w", err)
	} else {
		if v != 13 {
			return nil, fmt.Errorf("match failed for GarminColorMode-Mode: Expected %d != %d", 13, v)
		}
		val.Mode = GarminColorModeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for GarminColorMode-Color: %w", err)
	} else {
		val.Color = GarminColorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SimradTextMessage struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SimnetCommandConst `json:"proprietaryId"`
	A *uint8 `json:"a"`
	B *uint8 `json:"b"`
	C *uint8 `json:"c"`
	Sid *uint8 `json:"sid"`
	Prio *uint8 `json:"prio"`
	Text string `json:"text"`
}
func DecodeSimradTextMessage(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimradTextMessage
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimradTextMessage-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimradTextMessage-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimradTextMessage-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimradTextMessage-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimradTextMessage-ProprietaryId: %w", err)
	} else {
		if v != 50 {
			return nil, fmt.Errorf("match failed for SimradTextMessage-ProprietaryId: Expected %d != %d", 50, v)
		}
		val.ProprietaryId = SimnetCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimradTextMessage-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimradTextMessage-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimradTextMessage-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimradTextMessage-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimradTextMessage-Prio: %w", err)
	} else {
		val.Prio = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for SimradTextMessage-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type NavicoProductInformation struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProductCode *uint16 `json:"productCode"`
	Model string `json:"model"`
	A *uint8 `json:"a"`
	B *uint8 `json:"b"`
	C *uint8 `json:"c"`
	FirmwareVersion string `json:"firmwareVersion"`
	FirmwareDate string `json:"firmwareDate"`
	FirmwareTime string `json:"firmwareTime"`
}
func DecodeNavicoProductInformation(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val NavicoProductInformation
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for NavicoProductInformation-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for NavicoProductInformation-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-ProductCode: %w", err)
	} else {
		val.ProductCode = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-Model: %w", err)
	} else {
		val.Model = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(80); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-FirmwareVersion: %w", err)
	} else {
		val.FirmwareVersion = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-FirmwareDate: %w", err)
	} else {
		val.FirmwareDate = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoProductInformation-FirmwareTime: %w", err)
	} else {
		val.FirmwareTime = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type LowranceProductInformation struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProductCode *uint16 `json:"productCode"`
	Model string `json:"model"`
	A *uint8 `json:"a"`
	B *uint8 `json:"b"`
	C *uint8 `json:"c"`
	FirmwareVersion string `json:"firmwareVersion"`
	FirmwareDate string `json:"firmwareDate"`
	FirmwareTime string `json:"firmwareTime"`
}
func DecodeLowranceProductInformation(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val LowranceProductInformation
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-ManufacturerCode: %w", err)
	} else {
		if v != 140 {
			return nil, fmt.Errorf("match failed for LowranceProductInformation-ManufacturerCode: Expected %d != %d", 140, v)
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
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for LowranceProductInformation-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-ProductCode: %w", err)
	} else {
		val.ProductCode = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-Model: %w", err)
	} else {
		val.Model = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(80); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-FirmwareVersion: %w", err)
	} else {
		val.FirmwareVersion = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-FirmwareDate: %w", err)
	} else {
		val.FirmwareDate = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for LowranceProductInformation-FirmwareTime: %w", err)
	} else {
		val.FirmwareTime = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SimnetReprogramData struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Version *uint16 `json:"version"`
	Sequence *uint16 `json:"sequence"`
	Data []uint8 `json:"data"`
}
func DecodeSimnetReprogramData(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetReprogramData
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetReprogramData-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetReprogramData-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetReprogramData-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetReprogramData-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetReprogramData-Version: %w", err)
	} else {
		val.Version = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetReprogramData-Sequence: %w", err)
	} else {
		val.Sequence = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readBinaryData(1736); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetReprogramData-Data: %w", err)
	} else {
		val.Data = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FurunoUnknown130820 struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	A *uint8 `json:"a"`
	B *uint8 `json:"b"`
	C *uint8 `json:"c"`
	D *uint8 `json:"d"`
	E *uint8 `json:"e"`
}
func DecodeFurunoUnknown130820(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FurunoUnknown130820
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130820-ManufacturerCode: %w", err)
	} else {
		if v != 1855 {
			return nil, fmt.Errorf("match failed for FurunoUnknown130820-ManufacturerCode: Expected %d != %d", 1855, v)
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
		return nil, fmt.Errorf("parse failed for FurunoUnknown130820-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FurunoUnknown130820-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130820-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130820-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130820-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130820-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130820-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type NavicoAsciiData struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	A *uint8 `json:"a"`
	Message string `json:"message"`
}
func DecodeNavicoAsciiData(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val NavicoAsciiData
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoAsciiData-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for NavicoAsciiData-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for NavicoAsciiData-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for NavicoAsciiData-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoAsciiData-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(2048); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoAsciiData-Message: %w", err)
	} else {
		val.Message = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FurunoUnknown130821 struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Sid *uint8 `json:"sid"`
	A *uint8 `json:"a"`
	B *uint8 `json:"b"`
	C *uint8 `json:"c"`
	D *uint8 `json:"d"`
	E *uint8 `json:"e"`
	F *uint8 `json:"f"`
	G *uint8 `json:"g"`
	H *uint8 `json:"h"`
	I *uint8 `json:"i"`
}
func DecodeFurunoUnknown130821(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FurunoUnknown130821
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-ManufacturerCode: %w", err)
	} else {
		if v != 1855 {
			return nil, fmt.Errorf("match failed for FurunoUnknown130821-ManufacturerCode: Expected %d != %d", 1855, v)
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
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FurunoUnknown130821-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-F: %w", err)
	} else {
		val.F = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-G: %w", err)
	} else {
		val.G = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-H: %w", err)
	} else {
		val.H = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoUnknown130821-I: %w", err)
	} else {
		val.I = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type NavicoUnknown1 struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Data []uint8 `json:"data"`
}
func DecodeNavicoUnknown1(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val NavicoUnknown1
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoUnknown1-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for NavicoUnknown1-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for NavicoUnknown1-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for NavicoUnknown1-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readBinaryData(1848); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoUnknown1-Data: %w", err)
	} else {
		val.Data = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type MaretronProprietaryTemperatureHighRange struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Sid *uint8 `json:"sid"`
	Instance *uint8 `json:"instance"`
	Source TemperatureSourceConst `json:"source"`
	ActualTemperature *units.Temperature `json:"actualTemperature"`
	SetTemperature *units.Temperature `json:"setTemperature"`
}
func DecodeMaretronProprietaryTemperatureHighRange(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val MaretronProprietaryTemperatureHighRange
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryTemperatureHighRange-ManufacturerCode: %w", err)
	} else {
		if v != 137 {
			return nil, fmt.Errorf("match failed for MaretronProprietaryTemperatureHighRange-ManufacturerCode: Expected %d != %d", 137, v)
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
		return nil, fmt.Errorf("parse failed for MaretronProprietaryTemperatureHighRange-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for MaretronProprietaryTemperatureHighRange-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryTemperatureHighRange-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryTemperatureHighRange-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryTemperatureHighRange-Source: %w", err)
	} else {
		val.Source = TemperatureSourceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryTemperatureHighRange-ActualTemperature: %w", err)
	} else {
		val.ActualTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronProprietaryTemperatureHighRange-SetTemperature: %w", err)
	} else {
		val.SetTemperature = nullableUnit(units.Kelvin, v, units.NewTemperature)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type BGKeyValueData struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Repeating1 []BGKeyValueDataRepeating1 `json:"repeating1"`
}
type BGKeyValueDataRepeating1 struct {
	Key BandgKeyValueConst `json:"key"`
	Length *uint8 `json:"length"`
	Value []uint8 `json:"value"`
}
func DecodeBGKeyValueData(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val BGKeyValueData
	val.Info = Info
		var repeat1Count uint16 = 0
		var valueLength uint16
	
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for BGKeyValueData-ManufacturerCode: %w", err)
	} else {
		if v != 381 {
			return nil, fmt.Errorf("match failed for BGKeyValueData-ManufacturerCode: Expected %d != %d", 381, v)
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
		return nil, fmt.Errorf("parse failed for BGKeyValueData-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for BGKeyValueData-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	val.Repeating1 = make([]BGKeyValueDataRepeating1, 0)
	i := 0 
	for {
		var rep BGKeyValueDataRepeating1
		if v, err := stream.readLookupField(12); err != nil {
			return nil, fmt.Errorf("parse failed for BGKeyValueData-Key: %w", err)
		} else {
			rep.Key = BandgKeyValueConst(v)
		}
		if v, err := stream.readUInt8(4); err != nil {
			return nil, fmt.Errorf("parse failed for BGKeyValueData-Length: %w", err)
		} else {
			rep.Length = v
			if v != nil {
				valueLength = uint16(*v) * 8
			}
		
		}
		if v, err := stream.readBinaryData(valueLength); err != nil {
			return nil, fmt.Errorf("parse failed for BGKeyValueData-Value: %w", err)
		} else {
			rep.Value = v
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
type MaretronAnnunciator struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Field4 *uint8 `json:"field4"`
	Field5 *uint8 `json:"field5"`
	Field6 *uint16 `json:"field6"`
	Field7 *uint8 `json:"field7"`
	Field8 *uint16 `json:"field8"`
}
func DecodeMaretronAnnunciator(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val MaretronAnnunciator
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronAnnunciator-ManufacturerCode: %w", err)
	} else {
		if v != 137 {
			return nil, fmt.Errorf("match failed for MaretronAnnunciator-ManufacturerCode: Expected %d != %d", 137, v)
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
		return nil, fmt.Errorf("parse failed for MaretronAnnunciator-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for MaretronAnnunciator-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronAnnunciator-Field4: %w", err)
	} else {
		val.Field4 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronAnnunciator-Field5: %w", err)
	} else {
		val.Field5 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronAnnunciator-Field6: %w", err)
	} else {
		val.Field6 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronAnnunciator-Field7: %w", err)
	} else {
		val.Field7 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronAnnunciator-Field8: %w", err)
	} else {
		val.Field8 = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type NavicoUnknown2 struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Data []uint8 `json:"data"`
}
func DecodeNavicoUnknown2(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val NavicoUnknown2
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoUnknown2-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for NavicoUnknown2-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for NavicoUnknown2-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for NavicoUnknown2-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readBinaryData(80); err != nil {
		return nil, fmt.Errorf("parse failed for NavicoUnknown2-Data: %w", err)
	} else {
		val.Data = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type BGUserAndRemoteRename struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	DataType BandgKeyValueConst `json:"dataType"`
	Length *uint8 `json:"length"`
	Decimals BandgDecimalsConst `json:"decimals"`
	ShortName string `json:"shortName"`
	LongName string `json:"longName"`
}
func DecodeBGUserAndRemoteRename(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val BGUserAndRemoteRename
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for BGUserAndRemoteRename-ManufacturerCode: %w", err)
	} else {
		if v != 381 {
			return nil, fmt.Errorf("match failed for BGUserAndRemoteRename-ManufacturerCode: Expected %d != %d", 381, v)
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
		return nil, fmt.Errorf("parse failed for BGUserAndRemoteRename-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for BGUserAndRemoteRename-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(12); err != nil {
		return nil, fmt.Errorf("parse failed for BGUserAndRemoteRename-DataType: %w", err)
	} else {
		val.DataType = BandgKeyValueConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for BGUserAndRemoteRename-Length: %w", err)
	} else {
		val.Length = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for BGUserAndRemoteRename-Decimals: %w", err)
	} else {
		val.Decimals = BandgDecimalsConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(64); err != nil {
		return nil, fmt.Errorf("parse failed for BGUserAndRemoteRename-ShortName: %w", err)
	} else {
		val.ShortName = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(128); err != nil {
		return nil, fmt.Errorf("parse failed for BGUserAndRemoteRename-LongName: %w", err)
	} else {
		val.LongName = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SimnetFluidLevelSensorConfiguration struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	C *uint8 `json:"c"`
	Device *uint8 `json:"device"`
	Instance *uint8 `json:"instance"`
	F *uint8 `json:"f"`
	TankType TankTypeConst `json:"tankType"`
	Capacity *units.Volume `json:"capacity"`
	G *uint8 `json:"g"`
	H *int16 `json:"h"`
	I *int8 `json:"i"`
}
func DecodeSimnetFluidLevelSensorConfiguration(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetFluidLevelSensorConfiguration
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetFluidLevelSensorConfiguration-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetFluidLevelSensorConfiguration-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-Device: %w", err)
	} else {
		val.Device = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-F: %w", err)
	} else {
		val.F = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-TankType: %w", err)
	} else {
		val.TankType = TankTypeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-Capacity: %w", err)
	} else {
		val.Capacity = nullableUnit(units.Liter, v, units.NewVolume)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-G: %w", err)
	} else {
		val.G = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-H: %w", err)
	} else {
		val.H = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetFluidLevelSensorConfiguration-I: %w", err)
	} else {
		val.I = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type MaretronSwitchStatusCounter struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Instance *uint8 `json:"instance"`
	IndicatorNumber *uint8 `json:"indicatorNumber"`
	StartDate *uint16 `json:"startDate"`
	StartTime *float32 `json:"startTime"`
	OffCounter *uint8 `json:"offCounter"`
	OnCounter *uint8 `json:"onCounter"`
	ErrorCounter *uint8 `json:"errorCounter"`
	SwitchStatus OffOnConst `json:"switchStatus"`
}
func DecodeMaretronSwitchStatusCounter(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val MaretronSwitchStatusCounter
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-ManufacturerCode: %w", err)
	} else {
		if v != 137 {
			return nil, fmt.Errorf("match failed for MaretronSwitchStatusCounter-ManufacturerCode: Expected %d != %d", 137, v)
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
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for MaretronSwitchStatusCounter-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-IndicatorNumber: %w", err)
	} else {
		val.IndicatorNumber = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-StartDate: %w", err)
	} else {
		val.StartDate = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-StartTime: %w", err)
	} else {
		val.StartTime = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-OffCounter: %w", err)
	} else {
		val.OffCounter = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-OnCounter: %w", err)
	} else {
		val.OnCounter = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-ErrorCounter: %w", err)
	} else {
		val.ErrorCounter = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusCounter-SwitchStatus: %w", err)
	} else {
		val.SwitchStatus = OffOnConst(v)

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
type MaretronSwitchStatusTimer struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Instance *uint8 `json:"instance"`
	IndicatorNumber *uint8 `json:"indicatorNumber"`
	StartDate *uint16 `json:"startDate"`
	StartTime *float32 `json:"startTime"`
	AccumulatedOffPeriod *uint32 `json:"accumulatedOffPeriod"`
	AccumulatedOnPeriod *uint32 `json:"accumulatedOnPeriod"`
	AccumulatedErrorPeriod *uint32 `json:"accumulatedErrorPeriod"`
	SwitchStatus OffOnConst `json:"switchStatus"`
}
func DecodeMaretronSwitchStatusTimer(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val MaretronSwitchStatusTimer
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-ManufacturerCode: %w", err)
	} else {
		if v != 137 {
			return nil, fmt.Errorf("match failed for MaretronSwitchStatusTimer-ManufacturerCode: Expected %d != %d", 137, v)
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
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for MaretronSwitchStatusTimer-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-Instance: %w", err)
	} else {
		val.Instance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-IndicatorNumber: %w", err)
	} else {
		val.IndicatorNumber = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-StartDate: %w", err)
	} else {
		val.StartDate = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-StartTime: %w", err)
	} else {
		val.StartTime = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-AccumulatedOffPeriod: %w", err)
	} else {
		val.AccumulatedOffPeriod = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-AccumulatedOnPeriod: %w", err)
	} else {
		val.AccumulatedOnPeriod = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-AccumulatedErrorPeriod: %w", err)
	} else {
		val.AccumulatedErrorPeriod = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for MaretronSwitchStatusTimer-SwitchStatus: %w", err)
	} else {
		val.SwitchStatus = OffOnConst(v)

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
type FurunoSixDegreesOfFreedomMovement struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	A *int32 `json:"a"`
	B *int32 `json:"b"`
	C *int32 `json:"c"`
	D *int8 `json:"d"`
	E *int32 `json:"e"`
	F *int32 `json:"f"`
	G *int16 `json:"g"`
	H *int16 `json:"h"`
	I *int16 `json:"i"`
}
func DecodeFurunoSixDegreesOfFreedomMovement(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FurunoSixDegreesOfFreedomMovement
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-ManufacturerCode: %w", err)
	} else {
		if v != 1855 {
			return nil, fmt.Errorf("match failed for FurunoSixDegreesOfFreedomMovement-ManufacturerCode: Expected %d != %d", 1855, v)
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
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FurunoSixDegreesOfFreedomMovement-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-F: %w", err)
	} else {
		val.F = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-G: %w", err)
	} else {
		val.G = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-H: %w", err)
	} else {
		val.H = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoSixDegreesOfFreedomMovement-I: %w", err)
	} else {
		val.I = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SimnetAisClassBStaticDataMsg24PartB struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId *uint8 `json:"messageId"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	D *uint8 `json:"d"`
	E *uint8 `json:"e"`
	UserId *uint32 `json:"userId"`
	TypeOfShip ShipTypeConst `json:"typeOfShip"`
	VendorId string `json:"vendorId"`
	Callsign string `json:"callsign"`
	Length *units.Distance `json:"length"`
	Beam *units.Distance `json:"beam"`
	PositionReferenceFromStarboard *units.Distance `json:"positionReferenceFromStarboard"`
	PositionReferenceFromBow *units.Distance `json:"positionReferenceFromBow"`
	MothershipUserId *uint32 `json:"mothershipUserId"`
}
func DecodeSimnetAisClassBStaticDataMsg24PartB(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetAisClassBStaticDataMsg24PartB
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetAisClassBStaticDataMsg24PartB-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetAisClassBStaticDataMsg24PartB-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(6); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-MessageId: %w", err)
	} else {
		if v != nil && *v != 1 {
			return nil, fmt.Errorf("match failed for SimnetAisClassBStaticDataMsg24PartB-MessageId: Expected %d != %d", 1, *v)
		}
		val.MessageId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-UserId: %w", err)
	} else {
		val.UserId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-TypeOfShip: %w", err)
	} else {
		val.TypeOfShip = ShipTypeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(56); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-VendorId: %w", err)
	} else {
		val.VendorId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(56); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-Callsign: %w", err)
	} else {
		val.Callsign = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-Length: %w", err)
	} else {
		val.Length = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-Beam: %w", err)
	} else {
		val.Beam = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-PositionReferenceFromStarboard: %w", err)
	} else {
		val.PositionReferenceFromStarboard = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.1); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-PositionReferenceFromBow: %w", err)
	} else {
		val.PositionReferenceFromBow = nullableUnit(units.Meter, v, units.NewDistance)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAisClassBStaticDataMsg24PartB-MothershipUserId: %w", err)
	} else {
		val.MothershipUserId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(6)
	if stream.isEOF() {
		return val, nil
		}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}	
	return val, nil
}
type FurunoHeelAngleRollInformation struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	A *uint8 `json:"a"`
	B *uint8 `json:"b"`
	Yaw *float32 `json:"yaw"`
	Pitch *float32 `json:"pitch"`
	Roll *float32 `json:"roll"`
}
func DecodeFurunoHeelAngleRollInformation(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FurunoHeelAngleRollInformation
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoHeelAngleRollInformation-ManufacturerCode: %w", err)
	} else {
		if v != 1855 {
			return nil, fmt.Errorf("match failed for FurunoHeelAngleRollInformation-ManufacturerCode: Expected %d != %d", 1855, v)
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
		return nil, fmt.Errorf("parse failed for FurunoHeelAngleRollInformation-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FurunoHeelAngleRollInformation-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoHeelAngleRollInformation-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoHeelAngleRollInformation-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoHeelAngleRollInformation-Yaw: %w", err)
	} else {
		val.Yaw = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoHeelAngleRollInformation-Pitch: %w", err)
	} else {
		val.Pitch = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoHeelAngleRollInformation-Roll: %w", err)
	} else {
		val.Roll = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FurunoMultiSatsInViewExtended struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
}
func DecodeFurunoMultiSatsInViewExtended(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FurunoMultiSatsInViewExtended
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoMultiSatsInViewExtended-ManufacturerCode: %w", err)
	} else {
		if v != 1855 {
			return nil, fmt.Errorf("match failed for FurunoMultiSatsInViewExtended-ManufacturerCode: Expected %d != %d", 1855, v)
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
		return nil, fmt.Errorf("parse failed for FurunoMultiSatsInViewExtended-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FurunoMultiSatsInViewExtended-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SimnetKeyValue struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Address *uint8 `json:"address"`
	RepeatIndicator RepeatIndicatorConst `json:"repeatIndicator"`
	DisplayGroup SimnetDisplayGroupConst `json:"displayGroup"`
	Key SimnetKeyValueConst `json:"key"`
	Minlength *uint8 `json:"minlength"`
	Value []uint8 `json:"value"`
}
func DecodeSimnetKeyValue(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetKeyValue
	val.Info = Info
		var valueLength uint16
	
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetKeyValue-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetKeyValue-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetKeyValue-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetKeyValue-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetKeyValue-Address: %w", err)
	} else {
		val.Address = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetKeyValue-RepeatIndicator: %w", err)
	} else {
		val.RepeatIndicator = RepeatIndicatorConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetKeyValue-DisplayGroup: %w", err)
	} else {
		val.DisplayGroup = SimnetDisplayGroupConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetKeyValue-Key: %w", err)
	} else {
		val.Key = SimnetKeyValueConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetKeyValue-Minlength: %w", err)
	} else {
		val.Minlength = v
		if v != nil {
			valueLength = uint16(*v) * 8
		}
	

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readBinaryData(valueLength); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetKeyValue-Value: %w", err)
	} else {
		val.Value = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SimnetParameterSet struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Address *uint8 `json:"address"`
	B *uint8 `json:"b"`
	DisplayGroup SimnetDisplayGroupConst `json:"displayGroup"`
	D *uint16 `json:"d"`
	Key SimnetKeyValueConst `json:"key"`
	Length *uint8 `json:"length"`
	Value []uint8 `json:"value"`
}
func DecodeSimnetParameterSet(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetParameterSet
	val.Info = Info
		var valueLength uint16
	
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetParameterSet-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetParameterSet-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-Address: %w", err)
	} else {
		val.Address = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-DisplayGroup: %w", err)
	} else {
		val.DisplayGroup = SimnetDisplayGroupConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-Key: %w", err)
	} else {
		val.Key = SimnetKeyValueConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-Length: %w", err)
	} else {
		val.Length = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readBinaryData(valueLength); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetParameterSet-Value: %w", err)
	} else {
		val.Value = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FurunoMotionSensorStatusExtended struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
}
func DecodeFurunoMotionSensorStatusExtended(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FurunoMotionSensorStatusExtended
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FurunoMotionSensorStatusExtended-ManufacturerCode: %w", err)
	} else {
		if v != 1855 {
			return nil, fmt.Errorf("match failed for FurunoMotionSensorStatusExtended-ManufacturerCode: Expected %d != %d", 1855, v)
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
		return nil, fmt.Errorf("parse failed for FurunoMotionSensorStatusExtended-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FurunoMotionSensorStatusExtended-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SimnetApCommand struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Address *uint8 `json:"address"`
	ProprietaryId SimnetEventCommandConst `json:"proprietaryId"`
	ApStatus SimnetApStatusConst `json:"apStatus"`
	ApCommand SimnetApEventsConst `json:"apCommand"`
	Direction SimnetDirectionConst `json:"direction"`
	Angle *float32 `json:"angle"`
}
func DecodeSimnetApCommand(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetApCommand
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApCommand-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetApCommand-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetApCommand-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetApCommand-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApCommand-Address: %w", err)
	} else {
		val.Address = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApCommand-ProprietaryId: %w", err)
	} else {
		if v != 255 {
			return nil, fmt.Errorf("match failed for SimnetApCommand-ProprietaryId: Expected %d != %d", 255, v)
		}
		val.ProprietaryId = SimnetEventCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApCommand-ApStatus: %w", err)
	} else {
		val.ApStatus = SimnetApStatusConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApCommand-ApCommand: %w", err)
	} else {
		val.ApCommand = SimnetApEventsConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApCommand-Direction: %w", err)
	} else {
		val.Direction = SimnetDirectionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApCommand-Angle: %w", err)
	} else {
		val.Angle = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SimnetEventCommandApCommand struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SimnetEventCommandConst `json:"proprietaryId"`
	UnusedA *uint16 `json:"unusedA"`
	ControllingDevice *uint8 `json:"controllingDevice"`
	Event SimnetApEventsConst `json:"event"`
	UnusedB *uint8 `json:"unusedB"`
	Direction SimnetDirectionConst `json:"direction"`
	Angle *float32 `json:"angle"`
	UnusedC *uint8 `json:"unusedC"`
}
func DecodeSimnetEventCommandApCommand(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetEventCommandApCommand
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetEventCommandApCommand-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetEventCommandApCommand-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-ProprietaryId: %w", err)
	} else {
		if v != 2 {
			return nil, fmt.Errorf("match failed for SimnetEventCommandApCommand-ProprietaryId: Expected %d != %d", 2, v)
		}
		val.ProprietaryId = SimnetEventCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-UnusedA: %w", err)
	} else {
		val.UnusedA = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-ControllingDevice: %w", err)
	} else {
		val.ControllingDevice = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-Event: %w", err)
	} else {
		val.Event = SimnetApEventsConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-UnusedB: %w", err)
	} else {
		val.UnusedB = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-Direction: %w", err)
	} else {
		val.Direction = SimnetDirectionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-Angle: %w", err)
	} else {
		val.Angle = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventCommandApCommand-UnusedC: %w", err)
	} else {
		val.UnusedC = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SimnetAlarm struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	Address *uint8 `json:"address"`
	ProprietaryId SimnetEventCommandConst `json:"proprietaryId"`
	Alarm SimnetAlarmConst `json:"alarm"`
	MessageId *uint16 `json:"messageId"`
	F *uint8 `json:"f"`
	G *uint8 `json:"g"`
}
func DecodeSimnetAlarm(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetAlarm
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarm-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetAlarm-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetAlarm-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetAlarm-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarm-Address: %w", err)
	} else {
		val.Address = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarm-ProprietaryId: %w", err)
	} else {
		if v != 1 {
			return nil, fmt.Errorf("match failed for SimnetAlarm-ProprietaryId: Expected %d != %d", 1, v)
		}
		val.ProprietaryId = SimnetEventCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarm-Alarm: %w", err)
	} else {
		val.Alarm = SimnetAlarmConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarm-MessageId: %w", err)
	} else {
		val.MessageId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarm-F: %w", err)
	} else {
		val.F = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarm-G: %w", err)
	} else {
		val.G = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SimnetEventReplyApCommand struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SimnetEventCommandConst `json:"proprietaryId"`
	B *uint16 `json:"b"`
	Address *uint8 `json:"address"`
	Event SimnetApEventsConst `json:"event"`
	C *uint8 `json:"c"`
	Direction SimnetDirectionConst `json:"direction"`
	Angle *float32 `json:"angle"`
	G *uint8 `json:"g"`
}
func DecodeSimnetEventReplyApCommand(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetEventReplyApCommand
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetEventReplyApCommand-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetEventReplyApCommand-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-ProprietaryId: %w", err)
	} else {
		if v != 2 {
			return nil, fmt.Errorf("match failed for SimnetEventReplyApCommand-ProprietaryId: Expected %d != %d", 2, v)
		}
		val.ProprietaryId = SimnetEventCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-Address: %w", err)
	} else {
		val.Address = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-Event: %w", err)
	} else {
		val.Event = SimnetApEventsConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-Direction: %w", err)
	} else {
		val.Direction = SimnetDirectionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-Angle: %w", err)
	} else {
		val.Angle = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetEventReplyApCommand-G: %w", err)
	} else {
		val.G = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SimnetAlarmMessage struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId *uint16 `json:"messageId"`
	B *uint8 `json:"b"`
	C *uint8 `json:"c"`
	Text string `json:"text"`
}
func DecodeSimnetAlarmMessage(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetAlarmMessage
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarmMessage-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetAlarmMessage-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetAlarmMessage-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetAlarmMessage-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarmMessage-MessageId: %w", err)
	} else {
		val.MessageId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarmMessage-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarmMessage-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(1784); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetAlarmMessage-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SimnetApUnknown4 struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	A *uint8 `json:"a"`
	B *int32 `json:"b"`
	C *int32 `json:"c"`
	D *uint32 `json:"d"`
	E *int32 `json:"e"`
	F *uint32 `json:"f"`
}
func DecodeSimnetApUnknown4(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SimnetApUnknown4
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown4-ManufacturerCode: %w", err)
	} else {
		if v != 1857 {
			return nil, fmt.Errorf("match failed for SimnetApUnknown4-ManufacturerCode: Expected %d != %d", 1857, v)
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
		return nil, fmt.Errorf("parse failed for SimnetApUnknown4-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SimnetApUnknown4-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown4-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown4-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown4-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown4-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown4-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SimnetApUnknown4-F: %w", err)
	} else {
		val.F = v

		if stream.isEOF() {
			return val, nil
		} 
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
	return w.Bytes(), nil
}
func encodeSeatalkWirelessKeypadLightControlAny(v any) ([]byte, error) {
	val, ok := v.(*SeatalkWirelessKeypadLightControl)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkWirelessKeypadLightControl, got %T", v)
	}
	return EncodeSeatalkWirelessKeypadLightControl(val)
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
	return w.Bytes(), nil
}
func encodeSeatalkWirelessKeypadControlAny(v any) ([]byte, error) {
	val, ok := v.(*SeatalkWirelessKeypadControl)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkWirelessKeypadControl, got %T", v)
	}
	return EncodeSeatalkWirelessKeypadControl(val)
}

func EncodeVictronBatteryRegister(val *VictronBatteryRegister) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt16(val.RegisterId, 16)
	w.writeUInt32(val.Payload, 32)
	return w.Bytes(), nil
}
func encodeVictronBatteryRegisterAny(v any) ([]byte, error) {
	val, ok := v.(*VictronBatteryRegister)
	if !ok {
		return nil, fmt.Errorf("expected *VictronBatteryRegister, got %T", v)
	}
	return EncodeVictronBatteryRegister(val)
}

func EncodeFurunoHeave(val *FurunoHeave) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	var heaveRaw *float32
	if val.Heave != nil {
		heaveRaw = &val.Heave.Value
	}
	w.writeSignedResolution(heaveRaw, 32, 0.001)
	w.skipBits(16)
	return w.Bytes(), nil
}
func encodeFurunoHeaveAny(v any) ([]byte, error) {
	val, ok := v.(*FurunoHeave)
	if !ok {
		return nil, fmt.Errorf("expected *FurunoHeave, got %T", v)
	}
	return EncodeFurunoHeave(val)
}

func EncodeMaretronProprietaryDcBreakerCurrent(val *MaretronProprietaryDcBreakerCurrent) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.BankInstance, 8)
	w.writeUInt8(val.IndicatorNumber, 8)
	w.writeUnsignedResolution(val.BreakerCurrent, 16, 0.1)
	w.skipBits(16)
	return w.Bytes(), nil
}
func encodeMaretronProprietaryDcBreakerCurrentAny(v any) ([]byte, error) {
	val, ok := v.(*MaretronProprietaryDcBreakerCurrent)
	if !ok {
		return nil, fmt.Errorf("expected *MaretronProprietaryDcBreakerCurrent, got %T", v)
	}
	return EncodeMaretronProprietaryDcBreakerCurrent(val)
}

func EncodeAirmarBootStateAcknowledgment(val *AirmarBootStateAcknowledgment) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.BootState), 3)
	w.skipBits(45)
	return w.Bytes(), nil
}
func encodeAirmarBootStateAcknowledgmentAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarBootStateAcknowledgment)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarBootStateAcknowledgment, got %T", v)
	}
	return EncodeAirmarBootStateAcknowledgment(val)
}

func EncodeLowranceTemperature(val *LowranceTemperature) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.TemperatureSource), 8)
	var actualTemperatureRaw *float32
	if val.ActualTemperature != nil {
		actualTemperatureRaw = &val.ActualTemperature.Value
	}
	w.writeUnsignedResolution(actualTemperatureRaw, 16, 0.01)
	w.skipBits(24)
	return w.Bytes(), nil
}
func encodeLowranceTemperatureAny(v any) ([]byte, error) {
	val, ok := v.(*LowranceTemperature)
	if !ok {
		return nil, fmt.Errorf("expected *LowranceTemperature, got %T", v)
	}
	return EncodeLowranceTemperature(val)
}

func EncodeAirmarBootStateRequest(val *AirmarBootStateRequest) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(48)
	return w.Bytes(), nil
}
func encodeAirmarBootStateRequestAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarBootStateRequest)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarBootStateRequest, got %T", v)
	}
	return EncodeAirmarBootStateRequest(val)
}

func EncodeAirmarAccessLevel(val *AirmarAccessLevel) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.FormatCode, 8)
	w.writeLookupField(uint64(val.AccessLevel), 3)
	w.skipBits(5)
	w.writeUInt32(val.AccessSeedKey, 32)
	return w.Bytes(), nil
}
func encodeAirmarAccessLevelAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarAccessLevel)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarAccessLevel, got %T", v)
	}
	return EncodeAirmarAccessLevel(val)
}

func EncodeSimnetConfigureTemperatureSensor(val *SimnetConfigureTemperatureSensor) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(48)
	return w.Bytes(), nil
}
func encodeSimnetConfigureTemperatureSensorAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetConfigureTemperatureSensor)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetConfigureTemperatureSensor, got %T", v)
	}
	return EncodeSimnetConfigureTemperatureSensor(val)
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
	return w.Bytes(), nil
}
func encodeSeatalkAlarmAny(v any) ([]byte, error) {
	val, ok := v.(*SeatalkAlarm)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkAlarm, got %T", v)
	}
	return EncodeSeatalkAlarm(val)
}

func EncodeSimnetTrimTabSensorCalibration(val *SimnetTrimTabSensorCalibration) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(48)
	return w.Bytes(), nil
}
func encodeSimnetTrimTabSensorCalibrationAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetTrimTabSensorCalibration)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetTrimTabSensorCalibration, got %T", v)
	}
	return EncodeSimnetTrimTabSensorCalibration(val)
}

func EncodeSimnetPaddleWheelSpeedConfiguration(val *SimnetPaddleWheelSpeedConfiguration) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(48)
	return w.Bytes(), nil
}
func encodeSimnetPaddleWheelSpeedConfigurationAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetPaddleWheelSpeedConfiguration)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetPaddleWheelSpeedConfiguration, got %T", v)
	}
	return EncodeSimnetPaddleWheelSpeedConfiguration(val)
}

func EncodeSimnetClearFluidLevelWarnings(val *SimnetClearFluidLevelWarnings) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(48)
	return w.Bytes(), nil
}
func encodeSimnetClearFluidLevelWarningsAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetClearFluidLevelWarnings)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetClearFluidLevelWarnings, got %T", v)
	}
	return EncodeSimnetClearFluidLevelWarnings(val)
}

func EncodeSimnetLgc2000Configuration(val *SimnetLgc2000Configuration) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(48)
	return w.Bytes(), nil
}
func encodeSimnetLgc2000ConfigurationAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetLgc2000Configuration)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetLgc2000Configuration, got %T", v)
	}
	return EncodeSimnetLgc2000Configuration(val)
}

func EncodeDiverseYachtServicesLoadCell(val *DiverseYachtServicesLoadCell) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Instance, 8)
	w.skipBits(8)
	w.writeUInt32(val.LoadCell, 32)
	return w.Bytes(), nil
}
func encodeDiverseYachtServicesLoadCellAny(v any) ([]byte, error) {
	val, ok := v.(*DiverseYachtServicesLoadCell)
	if !ok {
		return nil, fmt.Errorf("expected *DiverseYachtServicesLoadCell, got %T", v)
	}
	return EncodeDiverseYachtServicesLoadCell(val)
}

func EncodeSimnetApUnknown1(val *SimnetApUnknown1) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt16(val.C, 16)
	w.writeUInt8(val.D, 8)
	w.skipBits(8)
	return w.Bytes(), nil
}
func encodeSimnetApUnknown1Any(v any) ([]byte, error) {
	val, ok := v.(*SimnetApUnknown1)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetApUnknown1, got %T", v)
	}
	return EncodeSimnetApUnknown1(val)
}

func EncodeSimnetDeviceStatus(val *SimnetDeviceStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.Model), 8)
	w.writeLookupField(uint64(val.Report), 8)
	w.writeLookupField(uint64(val.Status), 8)
	w.skipBits(24)
	return w.Bytes(), nil
}
func encodeSimnetDeviceStatusAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetDeviceStatus)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetDeviceStatus, got %T", v)
	}
	return EncodeSimnetDeviceStatus(val)
}

func EncodeSimnetDeviceStatusRequest(val *SimnetDeviceStatusRequest) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.Model), 8)
	w.writeLookupField(uint64(val.Report), 8)
	w.skipBits(32)
	return w.Bytes(), nil
}
func encodeSimnetDeviceStatusRequestAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetDeviceStatusRequest)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetDeviceStatusRequest, got %T", v)
	}
	return EncodeSimnetDeviceStatusRequest(val)
}

func EncodeSimnetPilotMode(val *SimnetPilotMode) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.Model), 8)
	w.writeLookupField(uint64(val.Report), 8)
	w.writeLookupField(uint64(val.Mode), 16)
	w.skipBits(16)
	return w.Bytes(), nil
}
func encodeSimnetPilotModeAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetPilotMode)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetPilotMode, got %T", v)
	}
	return EncodeSimnetPilotMode(val)
}

func EncodeSimnetDeviceModeRequest(val *SimnetDeviceModeRequest) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.Model), 8)
	w.writeLookupField(uint64(val.Report), 8)
	w.skipBits(32)
	return w.Bytes(), nil
}
func encodeSimnetDeviceModeRequestAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetDeviceModeRequest)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetDeviceModeRequest, got %T", v)
	}
	return EncodeSimnetDeviceModeRequest(val)
}

func EncodeSimnetSailingProcessorStatus(val *SimnetSailingProcessorStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.Model), 8)
	w.writeLookupField(uint64(val.Report), 8)
	w.writeBinaryData(val.Data, 32)
	return w.Bytes(), nil
}
func encodeSimnetSailingProcessorStatusAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetSailingProcessorStatus)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetSailingProcessorStatus, got %T", v)
	}
	return EncodeSimnetSailingProcessorStatus(val)
}

func EncodeNavicoWirelessBatteryStatus(val *NavicoWirelessBatteryStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Status, 8)
	w.writeUInt8(val.BatteryStatus, 8)
	w.writeUInt8(val.BatteryChargeStatus, 8)
	w.skipBits(24)
	return w.Bytes(), nil
}
func encodeNavicoWirelessBatteryStatusAny(v any) ([]byte, error) {
	val, ok := v.(*NavicoWirelessBatteryStatus)
	if !ok {
		return nil, fmt.Errorf("expected *NavicoWirelessBatteryStatus, got %T", v)
	}
	return EncodeNavicoWirelessBatteryStatus(val)
}

func EncodeNavicoWirelessSignalStatus(val *NavicoWirelessSignalStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Unknown, 8)
	w.writeUInt8(val.SignalStrength, 8)
	w.skipBits(32)
	return w.Bytes(), nil
}
func encodeNavicoWirelessSignalStatusAny(v any) ([]byte, error) {
	val, ok := v.(*NavicoWirelessSignalStatus)
	if !ok {
		return nil, fmt.Errorf("expected *NavicoWirelessSignalStatus, got %T", v)
	}
	return EncodeNavicoWirelessSignalStatus(val)
}

func EncodeSimnetApUnknown2(val *SimnetApUnknown2) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.C, 8)
	w.writeUInt8(val.D, 8)
	w.writeUInt8(val.E, 8)
	w.skipBits(8)
	return w.Bytes(), nil
}
func encodeSimnetApUnknown2Any(v any) ([]byte, error) {
	val, ok := v.(*SimnetApUnknown2)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetApUnknown2, got %T", v)
	}
	return EncodeSimnetApUnknown2(val)
}

func EncodeSimnetAutopilotAngle(val *SimnetAutopilotAngle) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(16)
	w.writeLookupField(uint64(val.Mode), 8)
	w.skipBits(8)
	w.writeUnsignedResolution(val.Angle, 16, 0.0001)
	return w.Bytes(), nil
}
func encodeSimnetAutopilotAngleAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetAutopilotAngle)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetAutopilotAngle, got %T", v)
	}
	return EncodeSimnetAutopilotAngle(val)
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
	return w.Bytes(), nil
}
func encodeSeatalkPilotWindDatumAny(v any) ([]byte, error) {
	val, ok := v.(*SeatalkPilotWindDatum)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkPilotWindDatum, got %T", v)
	}
	return EncodeSeatalkPilotWindDatum(val)
}

func EncodeSimnetMagneticField(val *SimnetMagneticField) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeSignedResolution(val.A, 16, 0.0001)
	w.writeUInt8(val.B, 8)
	w.writeSignedResolution(val.C, 16, 0.0001)
	w.writeSignedResolution(val.D, 16, 0.0001)
	w.skipBits(8)
	return w.Bytes(), nil
}
func encodeSimnetMagneticFieldAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetMagneticField)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetMagneticField, got %T", v)
	}
	return EncodeSimnetMagneticField(val)
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
	return w.Bytes(), nil
}
func encodeSeatalkPilotHeadingAny(v any) ([]byte, error) {
	val, ok := v.(*SeatalkPilotHeading)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkPilotHeading, got %T", v)
	}
	return EncodeSeatalkPilotHeading(val)
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
	return w.Bytes(), nil
}
func encodeSeatalkPilotLockedHeadingAny(v any) ([]byte, error) {
	val, ok := v.(*SeatalkPilotLockedHeading)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkPilotLockedHeading, got %T", v)
	}
	return EncodeSeatalkPilotLockedHeading(val)
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
	return w.Bytes(), nil
}
func encodeSeatalkSilenceAlarmAny(v any) ([]byte, error) {
	val, ok := v.(*SeatalkSilenceAlarm)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkSilenceAlarm, got %T", v)
	}
	return EncodeSeatalkSilenceAlarm(val)
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
	return w.Bytes(), nil
}
func encodeSeatalkKeypadMessageAny(v any) ([]byte, error) {
	val, ok := v.(*SeatalkKeypadMessage)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkKeypadMessage, got %T", v)
	}
	return EncodeSeatalkKeypadMessage(val)
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
	return w.Bytes(), nil
}
func encodeSeatalkKeypadHeartbeatAny(v any) ([]byte, error) {
	val, ok := v.(*SeatalkKeypadHeartbeat)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkKeypadHeartbeat, got %T", v)
	}
	return EncodeSeatalkKeypadHeartbeat(val)
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
	return w.Bytes(), nil
}
func encodeSeatalkPilotModeAny(v any) ([]byte, error) {
	val, ok := v.(*SeatalkPilotMode)
	if !ok {
		return nil, fmt.Errorf("expected *SeatalkPilotMode, got %T", v)
	}
	return EncodeSeatalkPilotMode(val)
}

func EncodeAirmarDepthQualityFactor(val *AirmarDepthQualityFactor) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Sid, 8)
	w.writeLookupField(uint64(val.DepthQualityFactor), 4)
	w.skipBits(36)
	return w.Bytes(), nil
}
func encodeAirmarDepthQualityFactorAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarDepthQualityFactor)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarDepthQualityFactor, got %T", v)
	}
	return EncodeAirmarDepthQualityFactor(val)
}

func EncodeAirmarSpeedPulseCount(val *AirmarSpeedPulseCount) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Sid, 8)
	w.writeUnsignedResolution(val.DurationOfInterval, 16, 0.001)
	w.writeUInt16(val.NumberOfPulsesReceived, 16)
	w.skipBits(8)
	return w.Bytes(), nil
}
func encodeAirmarSpeedPulseCountAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarSpeedPulseCount)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarSpeedPulseCount, got %T", v)
	}
	return EncodeAirmarSpeedPulseCount(val)
}

func EncodeAirmarDeviceInformation(val *AirmarDeviceInformation) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Sid, 8)
	var internalDeviceTemperatureRaw *float32
	if val.InternalDeviceTemperature != nil {
		internalDeviceTemperatureRaw = &val.InternalDeviceTemperature.Value
	}
	w.writeUnsignedResolution(internalDeviceTemperatureRaw, 16, 0.01)
	w.writeUnsignedResolution(val.SupplyVoltage, 16, 0.01)
	w.skipBits(8)
	return w.Bytes(), nil
}
func encodeAirmarDeviceInformationAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarDeviceInformation)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarDeviceInformation, got %T", v)
	}
	return EncodeAirmarDeviceInformation(val)
}

func EncodeSimnetApUnknown3(val *SimnetApUnknown3) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.C, 8)
	w.writeUInt8(val.D, 8)
	w.writeUInt8(val.E, 8)
	w.skipBits(8)
	return w.Bytes(), nil
}
func encodeSimnetApUnknown3Any(v any) ([]byte, error) {
	val, ok := v.(*SimnetApUnknown3)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetApUnknown3, got %T", v)
	}
	return EncodeSimnetApUnknown3(val)
}

func EncodeSimnetAutopilotMode(val *SimnetAutopilotMode) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(48)
	return w.Bytes(), nil
}
func encodeSimnetAutopilotModeAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetAutopilotMode)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetAutopilotMode, got %T", v)
	}
	return EncodeSimnetAutopilotMode(val)
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
	return w.Bytes(), nil
}
func encodeSeatalk1PilotModeAny(v any) ([]byte, error) {
	val, ok := v.(*Seatalk1PilotMode)
	if !ok {
		return nil, fmt.Errorf("expected *Seatalk1PilotMode, got %T", v)
	}
	return EncodeSeatalk1PilotMode(val)
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
	return w.Bytes(), nil
}
func encodeSeatalk1KeystrokeAny(v any) ([]byte, error) {
	val, ok := v.(*Seatalk1Keystroke)
	if !ok {
		return nil, fmt.Errorf("expected *Seatalk1Keystroke, got %T", v)
	}
	return EncodeSeatalk1Keystroke(val)
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
	return w.Bytes(), nil
}
func encodeSeatalk1DeviceIdentificationAny(v any) ([]byte, error) {
	val, ok := v.(*Seatalk1DeviceIdentification)
	if !ok {
		return nil, fmt.Errorf("expected *Seatalk1DeviceIdentification, got %T", v)
	}
	return EncodeSeatalk1DeviceIdentification(val)
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
	return w.Bytes(), nil
}
func encodeSeatalk1DisplayBrightnessAny(v any) ([]byte, error) {
	val, ok := v.(*Seatalk1DisplayBrightness)
	if !ok {
		return nil, fmt.Errorf("expected *Seatalk1DisplayBrightness, got %T", v)
	}
	return EncodeSeatalk1DisplayBrightness(val)
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
	return w.Bytes(), nil
}
func encodeSeatalk1DisplayColorAny(v any) ([]byte, error) {
	val, ok := v.(*Seatalk1DisplayColor)
	if !ok {
		return nil, fmt.Errorf("expected *Seatalk1DisplayColor, got %T", v)
	}
	return EncodeSeatalk1DisplayColor(val)
}

func EncodeAirmarAttitudeOffset(val *AirmarAttitudeOffset) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeSignedResolution(val.AzimuthOffset, 16, 0.0001)
	w.writeSignedResolution(val.PitchOffset, 16, 0.0001)
	w.writeSignedResolution(val.RollOffset, 16, 0.0001)
	return w.Bytes(), nil
}
func encodeAirmarAttitudeOffsetAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarAttitudeOffset)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarAttitudeOffset, got %T", v)
	}
	return EncodeAirmarAttitudeOffset(val)
}

func EncodeAirmarCalibrateCompass(val *AirmarCalibrateCompass) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.CalibrateFunction), 8)
	w.writeLookupField(uint64(val.CalibrationStatus), 8)
	w.writeUInt8(val.VerifyScore, 8)
	w.writeSignedResolution(val.XAxisGainValue, 16, 0.01)
	w.writeSignedResolution(val.YAxisGainValue, 16, 0.01)
	w.writeSignedResolution(val.ZAxisGainValue, 16, 0.01)
	w.writeSignedResolution(val.XAxisLinearOffset, 16, 0.01)
	w.writeSignedResolution(val.YAxisLinearOffset, 16, 0.01)
	w.writeSignedResolution(val.ZAxisLinearOffset, 16, 0.01)
	w.writeSignedResolution(val.XAxisAngularOffset, 16, 0.1)
	w.writeSignedResolution(val.PitchAndRollDamping, 16, 0.05)
	w.writeSignedResolution(val.CompassRateGyroDamping, 16, 0.05)
	return w.Bytes(), nil
}
func encodeAirmarCalibrateCompassAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarCalibrateCompass)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarCalibrateCompass, got %T", v)
	}
	return EncodeAirmarCalibrateCompass(val)
}

func EncodeAirmarTrueWindOptions(val *AirmarTrueWindOptions) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.CogSubstitutionForHdg), 2)
	w.skipBits(22)
	return w.Bytes(), nil
}
func encodeAirmarTrueWindOptionsAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarTrueWindOptions)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarTrueWindOptions, got %T", v)
	}
	return EncodeAirmarTrueWindOptions(val)
}

func EncodeAirmarSimulateMode(val *AirmarSimulateMode) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.SimulateMode), 2)
	w.skipBits(22)
	return w.Bytes(), nil
}
func encodeAirmarSimulateModeAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarSimulateMode)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarSimulateMode, got %T", v)
	}
	return EncodeAirmarSimulateMode(val)
}

func EncodeAirmarCalibrateDepth(val *AirmarCalibrateDepth) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	var speedOfSoundModeRaw *float32
	if val.SpeedOfSoundMode != nil {
		speedOfSoundModeRaw = &val.SpeedOfSoundMode.Value
	}
	w.writeUnsignedResolution(speedOfSoundModeRaw, 16, 0.1)
	w.skipBits(8)
	return w.Bytes(), nil
}
func encodeAirmarCalibrateDepthAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarCalibrateDepth)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarCalibrateDepth, got %T", v)
	}
	return EncodeAirmarCalibrateDepth(val)
}

func EncodeAirmarCalibrateSpeed(val *AirmarCalibrateSpeed) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt8(val.NumberOfPairsOfDataPoints, 8)
	for _, rep := range val.Repeating1 {
		w.writeUnsignedResolution(rep.InputFrequency, 16, 0.1)
		var outputSpeedRaw *float32
		if rep.OutputSpeed != nil {
			outputSpeedRaw = &rep.OutputSpeed.Value
		}
		w.writeUnsignedResolution(outputSpeedRaw, 16, 0.01)
	}
	return w.Bytes(), nil
}
func encodeAirmarCalibrateSpeedAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarCalibrateSpeed)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarCalibrateSpeed, got %T", v)
	}
	return EncodeAirmarCalibrateSpeed(val)
}

func EncodeAirmarCalibrateTemperature(val *AirmarCalibrateTemperature) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.TemperatureInstance), 2)
	w.skipBits(6)
	var temperatureOffsetRaw *float32
	if val.TemperatureOffset != nil {
		temperatureOffsetRaw = &val.TemperatureOffset.Value
	}
	w.writeSignedResolution(temperatureOffsetRaw, 16, 0.001)
	return w.Bytes(), nil
}
func encodeAirmarCalibrateTemperatureAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarCalibrateTemperature)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarCalibrateTemperature, got %T", v)
	}
	return EncodeAirmarCalibrateTemperature(val)
}

func EncodeAirmarSpeedFilterNone(val *AirmarSpeedFilterNone) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt8(val.FilterType, 4)
	w.skipBits(4)
	w.writeUnsignedResolution(val.SampleInterval, 16, 0.01)
	return w.Bytes(), nil
}
func encodeAirmarSpeedFilterNoneAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarSpeedFilterNone)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarSpeedFilterNone, got %T", v)
	}
	return EncodeAirmarSpeedFilterNone(val)
}

func EncodeAirmarSpeedFilterIir(val *AirmarSpeedFilterIir) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt8(val.FilterType, 4)
	w.skipBits(4)
	w.writeUnsignedResolution(val.SampleInterval, 16, 0.01)
	w.writeUnsignedResolution(val.FilterDuration, 16, 0.01)
	return w.Bytes(), nil
}
func encodeAirmarSpeedFilterIirAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarSpeedFilterIir)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarSpeedFilterIir, got %T", v)
	}
	return EncodeAirmarSpeedFilterIir(val)
}

func EncodeAirmarTemperatureFilterNone(val *AirmarTemperatureFilterNone) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt8(val.FilterType, 4)
	w.skipBits(4)
	w.writeUnsignedResolution(val.SampleInterval, 16, 0.01)
	return w.Bytes(), nil
}
func encodeAirmarTemperatureFilterNoneAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarTemperatureFilterNone)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarTemperatureFilterNone, got %T", v)
	}
	return EncodeAirmarTemperatureFilterNone(val)
}

func EncodeAirmarTemperatureFilterIir(val *AirmarTemperatureFilterIir) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt8(val.FilterType, 4)
	w.skipBits(4)
	w.writeUnsignedResolution(val.SampleInterval, 16, 0.01)
	w.writeUnsignedResolution(val.FilterDuration, 16, 0.01)
	return w.Bytes(), nil
}
func encodeAirmarTemperatureFilterIirAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarTemperatureFilterIir)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarTemperatureFilterIir, got %T", v)
	}
	return EncodeAirmarTemperatureFilterIir(val)
}

func EncodeAirmarNmea2000Options(val *AirmarNmea2000Options) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.TransmissionInterval), 2)
	w.skipBits(22)
	return w.Bytes(), nil
}
func encodeAirmarNmea2000OptionsAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarNmea2000Options)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarNmea2000Options, got %T", v)
	}
	return EncodeAirmarNmea2000Options(val)
}

func EncodeAirmarAddressableMultiFrame(val *AirmarAddressableMultiFrame) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.ProprietaryId, 8)
	return w.Bytes(), nil
}
func encodeAirmarAddressableMultiFrameAny(v any) ([]byte, error) {
	val, ok := v.(*AirmarAddressableMultiFrame)
	if !ok {
		return nil, fmt.Errorf("expected *AirmarAddressableMultiFrame, got %T", v)
	}
	return EncodeAirmarAddressableMultiFrame(val)
}

func EncodeMaretronSlaveResponse(val *MaretronSlaveResponse) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt16(val.ProductCode, 16)
	w.writeUInt16(val.SoftwareCode, 16)
	w.writeUInt8(val.Command, 8)
	w.writeUInt8(val.Status, 8)
	return w.Bytes(), nil
}
func encodeMaretronSlaveResponseAny(v any) ([]byte, error) {
	val, ok := v.(*MaretronSlaveResponse)
	if !ok {
		return nil, fmt.Errorf("expected *MaretronSlaveResponse, got %T", v)
	}
	return EncodeMaretronSlaveResponse(val)
}

func EncodeGarminDayMode(val *GarminDayMode) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.UnknownId1, 8)
	w.writeUInt8(val.UnknownId2, 8)
	w.writeUInt8(val.UnknownId3, 8)
	w.writeUInt8(val.UnknownId4, 8)
	w.skipBits(16)
	w.writeLookupField(uint64(val.Mode), 8)
	w.skipBits(8)
	w.writeLookupField(uint64(val.Backlight), 8)
	return w.Bytes(), nil
}
func encodeGarminDayModeAny(v any) ([]byte, error) {
	val, ok := v.(*GarminDayMode)
	if !ok {
		return nil, fmt.Errorf("expected *GarminDayMode, got %T", v)
	}
	return EncodeGarminDayMode(val)
}

func EncodeGarminNightMode(val *GarminNightMode) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.UnknownId1, 8)
	w.writeUInt8(val.UnknownId2, 8)
	w.writeUInt8(val.UnknownId3, 8)
	w.writeUInt8(val.UnknownId4, 8)
	w.skipBits(16)
	w.writeLookupField(uint64(val.Mode), 8)
	w.skipBits(8)
	w.writeLookupField(uint64(val.Backlight), 8)
	return w.Bytes(), nil
}
func encodeGarminNightModeAny(v any) ([]byte, error) {
	val, ok := v.(*GarminNightMode)
	if !ok {
		return nil, fmt.Errorf("expected *GarminNightMode, got %T", v)
	}
	return EncodeGarminNightMode(val)
}

func EncodeGarminColorMode(val *GarminColorMode) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.UnknownId1, 8)
	w.writeUInt8(val.UnknownId2, 8)
	w.writeUInt8(val.UnknownId3, 8)
	w.writeUInt8(val.UnknownId4, 8)
	w.skipBits(16)
	w.writeLookupField(uint64(val.Mode), 8)
	w.skipBits(8)
	w.writeLookupField(uint64(val.Color), 8)
	return w.Bytes(), nil
}
func encodeGarminColorModeAny(v any) ([]byte, error) {
	val, ok := v.(*GarminColorMode)
	if !ok {
		return nil, fmt.Errorf("expected *GarminColorMode, got %T", v)
	}
	return EncodeGarminColorMode(val)
}

func EncodeSimradTextMessage(val *SimradTextMessage) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.C, 8)
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.Prio, 8)
	w.writeFixedString(val.Text, 256)
	return w.Bytes(), nil
}
func encodeSimradTextMessageAny(v any) ([]byte, error) {
	val, ok := v.(*SimradTextMessage)
	if !ok {
		return nil, fmt.Errorf("expected *SimradTextMessage, got %T", v)
	}
	return EncodeSimradTextMessage(val)
}

func EncodeNavicoProductInformation(val *NavicoProductInformation) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt16(val.ProductCode, 16)
	w.writeFixedString(val.Model, 256)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.C, 8)
	w.writeFixedString(val.FirmwareVersion, 80)
	w.writeFixedString(val.FirmwareDate, 256)
	w.writeFixedString(val.FirmwareTime, 256)
	return w.Bytes(), nil
}
func encodeNavicoProductInformationAny(v any) ([]byte, error) {
	val, ok := v.(*NavicoProductInformation)
	if !ok {
		return nil, fmt.Errorf("expected *NavicoProductInformation, got %T", v)
	}
	return EncodeNavicoProductInformation(val)
}

func EncodeLowranceProductInformation(val *LowranceProductInformation) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt16(val.ProductCode, 16)
	w.writeFixedString(val.Model, 256)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.C, 8)
	w.writeFixedString(val.FirmwareVersion, 80)
	w.writeFixedString(val.FirmwareDate, 256)
	w.writeFixedString(val.FirmwareTime, 256)
	return w.Bytes(), nil
}
func encodeLowranceProductInformationAny(v any) ([]byte, error) {
	val, ok := v.(*LowranceProductInformation)
	if !ok {
		return nil, fmt.Errorf("expected *LowranceProductInformation, got %T", v)
	}
	return EncodeLowranceProductInformation(val)
}

func EncodeSimnetReprogramData(val *SimnetReprogramData) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt16(val.Version, 16)
	w.writeUInt16(val.Sequence, 16)
	w.writeBinaryData(val.Data, 1736)
	return w.Bytes(), nil
}
func encodeSimnetReprogramDataAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetReprogramData)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetReprogramData, got %T", v)
	}
	return EncodeSimnetReprogramData(val)
}

func EncodeFurunoUnknown130820(val *FurunoUnknown130820) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.C, 8)
	w.writeUInt8(val.D, 8)
	w.writeUInt8(val.E, 8)
	return w.Bytes(), nil
}
func encodeFurunoUnknown130820Any(v any) ([]byte, error) {
	val, ok := v.(*FurunoUnknown130820)
	if !ok {
		return nil, fmt.Errorf("expected *FurunoUnknown130820, got %T", v)
	}
	return EncodeFurunoUnknown130820(val)
}

func EncodeNavicoAsciiData(val *NavicoAsciiData) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.A, 8)
	w.writeFixedString(val.Message, 2048)
	return w.Bytes(), nil
}
func encodeNavicoAsciiDataAny(v any) ([]byte, error) {
	val, ok := v.(*NavicoAsciiData)
	if !ok {
		return nil, fmt.Errorf("expected *NavicoAsciiData, got %T", v)
	}
	return EncodeNavicoAsciiData(val)
}

func EncodeFurunoUnknown130821(val *FurunoUnknown130821) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.C, 8)
	w.writeUInt8(val.D, 8)
	w.writeUInt8(val.E, 8)
	w.writeUInt8(val.F, 8)
	w.writeUInt8(val.G, 8)
	w.writeUInt8(val.H, 8)
	w.writeUInt8(val.I, 8)
	return w.Bytes(), nil
}
func encodeFurunoUnknown130821Any(v any) ([]byte, error) {
	val, ok := v.(*FurunoUnknown130821)
	if !ok {
		return nil, fmt.Errorf("expected *FurunoUnknown130821, got %T", v)
	}
	return EncodeFurunoUnknown130821(val)
}

func EncodeNavicoUnknown1(val *NavicoUnknown1) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeBinaryData(val.Data, 1848)
	return w.Bytes(), nil
}
func encodeNavicoUnknown1Any(v any) ([]byte, error) {
	val, ok := v.(*NavicoUnknown1)
	if !ok {
		return nil, fmt.Errorf("expected *NavicoUnknown1, got %T", v)
	}
	return EncodeNavicoUnknown1(val)
}

func EncodeMaretronProprietaryTemperatureHighRange(val *MaretronProprietaryTemperatureHighRange) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Sid, 8)
	w.writeUInt8(val.Instance, 8)
	w.writeLookupField(uint64(val.Source), 8)
	var actualTemperatureRaw *float32
	if val.ActualTemperature != nil {
		actualTemperatureRaw = &val.ActualTemperature.Value
	}
	w.writeUnsignedResolution(actualTemperatureRaw, 16, 0.1)
	var setTemperatureRaw *float32
	if val.SetTemperature != nil {
		setTemperatureRaw = &val.SetTemperature.Value
	}
	w.writeUnsignedResolution(setTemperatureRaw, 16, 0.1)
	return w.Bytes(), nil
}
func encodeMaretronProprietaryTemperatureHighRangeAny(v any) ([]byte, error) {
	val, ok := v.(*MaretronProprietaryTemperatureHighRange)
	if !ok {
		return nil, fmt.Errorf("expected *MaretronProprietaryTemperatureHighRange, got %T", v)
	}
	return EncodeMaretronProprietaryTemperatureHighRange(val)
}

func EncodeBGKeyValueData(val *BGKeyValueData) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	for _, rep := range val.Repeating1 {
		w.writeLookupField(uint64(rep.Key), 12)
		w.writeUInt8(rep.Length, 4)
		var valueLength uint16
		if rep.Length != nil {
			valueLength = uint16(*rep.Length) * 8
		}
		w.writeBinaryData(rep.Value, valueLength)
	}
	return w.Bytes(), nil
}
func encodeBGKeyValueDataAny(v any) ([]byte, error) {
	val, ok := v.(*BGKeyValueData)
	if !ok {
		return nil, fmt.Errorf("expected *BGKeyValueData, got %T", v)
	}
	return EncodeBGKeyValueData(val)
}

func EncodeMaretronAnnunciator(val *MaretronAnnunciator) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Field4, 8)
	w.writeUInt8(val.Field5, 8)
	w.writeUInt16(val.Field6, 16)
	w.writeUInt8(val.Field7, 8)
	w.writeUInt16(val.Field8, 16)
	return w.Bytes(), nil
}
func encodeMaretronAnnunciatorAny(v any) ([]byte, error) {
	val, ok := v.(*MaretronAnnunciator)
	if !ok {
		return nil, fmt.Errorf("expected *MaretronAnnunciator, got %T", v)
	}
	return EncodeMaretronAnnunciator(val)
}

func EncodeNavicoUnknown2(val *NavicoUnknown2) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeBinaryData(val.Data, 80)
	return w.Bytes(), nil
}
func encodeNavicoUnknown2Any(v any) ([]byte, error) {
	val, ok := v.(*NavicoUnknown2)
	if !ok {
		return nil, fmt.Errorf("expected *NavicoUnknown2, got %T", v)
	}
	return EncodeNavicoUnknown2(val)
}

func EncodeBGUserAndRemoteRename(val *BGUserAndRemoteRename) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.DataType), 12)
	w.writeUInt8(val.Length, 4)
	w.skipBits(8)
	w.writeLookupField(uint64(val.Decimals), 8)
	w.writeFixedString(val.ShortName, 64)
	w.writeFixedString(val.LongName, 128)
	return w.Bytes(), nil
}
func encodeBGUserAndRemoteRenameAny(v any) ([]byte, error) {
	val, ok := v.(*BGUserAndRemoteRename)
	if !ok {
		return nil, fmt.Errorf("expected *BGUserAndRemoteRename, got %T", v)
	}
	return EncodeBGUserAndRemoteRename(val)
}

func EncodeSimnetFluidLevelSensorConfiguration(val *SimnetFluidLevelSensorConfiguration) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.C, 8)
	w.writeUInt8(val.Device, 8)
	w.writeUInt8(val.Instance, 8)
	w.writeUInt8(val.F, 4)
	w.writeLookupField(uint64(val.TankType), 4)
	var capacityRaw *float32
	if val.Capacity != nil {
		capacityRaw = &val.Capacity.Value
	}
	w.writeUnsignedResolution(capacityRaw, 32, 0.1)
	w.writeUInt8(val.G, 8)
	w.writeInt16(val.H, 16)
	w.writeInt8(val.I, 8)
	return w.Bytes(), nil
}
func encodeSimnetFluidLevelSensorConfigurationAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetFluidLevelSensorConfiguration)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetFluidLevelSensorConfiguration, got %T", v)
	}
	return EncodeSimnetFluidLevelSensorConfiguration(val)
}

func EncodeMaretronSwitchStatusCounter(val *MaretronSwitchStatusCounter) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Instance, 8)
	w.writeUInt8(val.IndicatorNumber, 8)
	w.writeUInt16(val.StartDate, 16)
	w.writeUnsignedResolution(val.StartTime, 32, 0.0001)
	w.writeUInt8(val.OffCounter, 8)
	w.writeUInt8(val.OnCounter, 8)
	w.writeUInt8(val.ErrorCounter, 8)
	w.writeLookupField(uint64(val.SwitchStatus), 2)
	w.skipBits(6)
	return w.Bytes(), nil
}
func encodeMaretronSwitchStatusCounterAny(v any) ([]byte, error) {
	val, ok := v.(*MaretronSwitchStatusCounter)
	if !ok {
		return nil, fmt.Errorf("expected *MaretronSwitchStatusCounter, got %T", v)
	}
	return EncodeMaretronSwitchStatusCounter(val)
}

func EncodeMaretronSwitchStatusTimer(val *MaretronSwitchStatusTimer) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Instance, 8)
	w.writeUInt8(val.IndicatorNumber, 8)
	w.writeUInt16(val.StartDate, 16)
	w.writeUnsignedResolution(val.StartTime, 32, 0.0001)
	w.writeUInt32(val.AccumulatedOffPeriod, 32)
	w.writeUInt32(val.AccumulatedOnPeriod, 32)
	w.writeUInt32(val.AccumulatedErrorPeriod, 32)
	w.writeLookupField(uint64(val.SwitchStatus), 2)
	w.skipBits(6)
	return w.Bytes(), nil
}
func encodeMaretronSwitchStatusTimerAny(v any) ([]byte, error) {
	val, ok := v.(*MaretronSwitchStatusTimer)
	if !ok {
		return nil, fmt.Errorf("expected *MaretronSwitchStatusTimer, got %T", v)
	}
	return EncodeMaretronSwitchStatusTimer(val)
}

func EncodeFurunoSixDegreesOfFreedomMovement(val *FurunoSixDegreesOfFreedomMovement) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeInt32(val.A, 32)
	w.writeInt32(val.B, 32)
	w.writeInt32(val.C, 32)
	w.writeInt8(val.D, 8)
	w.writeInt32(val.E, 32)
	w.writeInt32(val.F, 32)
	w.writeInt16(val.G, 16)
	w.writeInt16(val.H, 16)
	w.writeInt16(val.I, 16)
	return w.Bytes(), nil
}
func encodeFurunoSixDegreesOfFreedomMovementAny(v any) ([]byte, error) {
	val, ok := v.(*FurunoSixDegreesOfFreedomMovement)
	if !ok {
		return nil, fmt.Errorf("expected *FurunoSixDegreesOfFreedomMovement, got %T", v)
	}
	return EncodeFurunoSixDegreesOfFreedomMovement(val)
}

func EncodeSimnetAisClassBStaticDataMsg24PartB(val *SimnetAisClassBStaticDataMsg24PartB) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.MessageId, 6)
	w.writeLookupField(uint64(val.RepeatIndicator), 2)
	w.writeUInt8(val.D, 8)
	w.writeUInt8(val.E, 8)
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
	w.skipBits(6)
	w.skipBits(2)
	return w.Bytes(), nil
}
func encodeSimnetAisClassBStaticDataMsg24PartBAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetAisClassBStaticDataMsg24PartB)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetAisClassBStaticDataMsg24PartB, got %T", v)
	}
	return EncodeSimnetAisClassBStaticDataMsg24PartB(val)
}

func EncodeFurunoHeelAngleRollInformation(val *FurunoHeelAngleRollInformation) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeSignedResolution(val.Yaw, 16, 0.0001)
	w.writeSignedResolution(val.Pitch, 16, 0.0001)
	w.writeSignedResolution(val.Roll, 16, 0.0001)
	return w.Bytes(), nil
}
func encodeFurunoHeelAngleRollInformationAny(v any) ([]byte, error) {
	val, ok := v.(*FurunoHeelAngleRollInformation)
	if !ok {
		return nil, fmt.Errorf("expected *FurunoHeelAngleRollInformation, got %T", v)
	}
	return EncodeFurunoHeelAngleRollInformation(val)
}

func EncodeFurunoMultiSatsInViewExtended(val *FurunoMultiSatsInViewExtended) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	return w.Bytes(), nil
}
func encodeFurunoMultiSatsInViewExtendedAny(v any) ([]byte, error) {
	val, ok := v.(*FurunoMultiSatsInViewExtended)
	if !ok {
		return nil, fmt.Errorf("expected *FurunoMultiSatsInViewExtended, got %T", v)
	}
	return EncodeFurunoMultiSatsInViewExtended(val)
}

func EncodeSimnetKeyValue(val *SimnetKeyValue) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Address, 8)
	w.writeLookupField(uint64(val.RepeatIndicator), 8)
	w.writeLookupField(uint64(val.DisplayGroup), 8)
	w.skipBits(8)
	w.writeLookupField(uint64(val.Key), 16)
	w.skipBits(8)
	w.writeUInt8(val.Minlength, 8)
	return w.Bytes(), nil
}
func encodeSimnetKeyValueAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetKeyValue)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetKeyValue, got %T", v)
	}
	return EncodeSimnetKeyValue(val)
}

func EncodeSimnetParameterSet(val *SimnetParameterSet) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Address, 8)
	w.writeUInt8(val.B, 8)
	w.writeLookupField(uint64(val.DisplayGroup), 8)
	w.writeUInt16(val.D, 16)
	w.writeLookupField(uint64(val.Key), 16)
	w.skipBits(8)
	w.writeUInt8(val.Length, 8)
	return w.Bytes(), nil
}
func encodeSimnetParameterSetAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetParameterSet)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetParameterSet, got %T", v)
	}
	return EncodeSimnetParameterSet(val)
}

func EncodeFurunoMotionSensorStatusExtended(val *FurunoMotionSensorStatusExtended) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	return w.Bytes(), nil
}
func encodeFurunoMotionSensorStatusExtendedAny(v any) ([]byte, error) {
	val, ok := v.(*FurunoMotionSensorStatusExtended)
	if !ok {
		return nil, fmt.Errorf("expected *FurunoMotionSensorStatusExtended, got %T", v)
	}
	return EncodeFurunoMotionSensorStatusExtended(val)
}

func EncodeSimnetApCommand(val *SimnetApCommand) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Address, 8)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.ApStatus), 8)
	w.writeLookupField(uint64(val.ApCommand), 8)
	w.skipBits(8)
	w.writeLookupField(uint64(val.Direction), 8)
	w.writeUnsignedResolution(val.Angle, 16, 0.0001)
	return w.Bytes(), nil
}
func encodeSimnetApCommandAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetApCommand)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetApCommand, got %T", v)
	}
	return EncodeSimnetApCommand(val)
}

func EncodeSimnetEventCommandApCommand(val *SimnetEventCommandApCommand) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt16(val.UnusedA, 16)
	w.writeUInt8(val.ControllingDevice, 8)
	w.writeLookupField(uint64(val.Event), 8)
	w.writeUInt8(val.UnusedB, 8)
	w.writeLookupField(uint64(val.Direction), 8)
	w.writeUnsignedResolution(val.Angle, 16, 0.0001)
	w.writeUInt8(val.UnusedC, 8)
	return w.Bytes(), nil
}
func encodeSimnetEventCommandApCommandAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetEventCommandApCommand)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetEventCommandApCommand, got %T", v)
	}
	return EncodeSimnetEventCommandApCommand(val)
}

func EncodeSimnetAlarm(val *SimnetAlarm) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.Address, 8)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.skipBits(8)
	w.writeLookupField(uint64(val.Alarm), 16)
	w.writeUInt16(val.MessageId, 16)
	w.writeUInt8(val.F, 8)
	w.writeUInt8(val.G, 8)
	return w.Bytes(), nil
}
func encodeSimnetAlarmAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetAlarm)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetAlarm, got %T", v)
	}
	return EncodeSimnetAlarm(val)
}

func EncodeSimnetEventReplyApCommand(val *SimnetEventReplyApCommand) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt16(val.B, 16)
	w.writeUInt8(val.Address, 8)
	w.writeLookupField(uint64(val.Event), 8)
	w.writeUInt8(val.C, 8)
	w.writeLookupField(uint64(val.Direction), 8)
	w.writeUnsignedResolution(val.Angle, 16, 0.0001)
	w.writeUInt8(val.G, 8)
	return w.Bytes(), nil
}
func encodeSimnetEventReplyApCommandAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetEventReplyApCommand)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetEventReplyApCommand, got %T", v)
	}
	return EncodeSimnetEventReplyApCommand(val)
}

func EncodeSimnetAlarmMessage(val *SimnetAlarmMessage) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt16(val.MessageId, 16)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.C, 8)
	w.writeFixedString(val.Text, 1784)
	return w.Bytes(), nil
}
func encodeSimnetAlarmMessageAny(v any) ([]byte, error) {
	val, ok := v.(*SimnetAlarmMessage)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetAlarmMessage, got %T", v)
	}
	return EncodeSimnetAlarmMessage(val)
}

func EncodeSimnetApUnknown4(val *SimnetApUnknown4) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.A, 8)
	w.writeInt32(val.B, 32)
	w.writeInt32(val.C, 32)
	w.writeUInt32(val.D, 32)
	w.writeInt32(val.E, 32)
	w.writeUInt32(val.F, 32)
	return w.Bytes(), nil
}
func encodeSimnetApUnknown4Any(v any) ([]byte, error) {
	val, ok := v.(*SimnetApUnknown4)
	if !ok {
		return nil, fmt.Errorf("expected *SimnetApUnknown4, got %T", v)
	}
	return EncodeSimnetApUnknown4(val)
}
