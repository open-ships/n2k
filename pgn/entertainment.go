package pgn

import (
	"fmt"
)

type FusionMediaControl struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId *uint8 `json:"proprietaryId"`
	Unknown *uint8 `json:"unknown"`
	SourceId *uint8 `json:"sourceId"`
	Command FusionCommandConst `json:"command"`
}
func DecodeFusionMediaControl(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionMediaControl
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMediaControl-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionMediaControl-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionMediaControl-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionMediaControl-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMediaControl-ProprietaryId: %w", err)
	} else {
		if v != nil && *v != 3 {
			return nil, fmt.Errorf("match failed for FusionMediaControl-ProprietaryId: Expected %d != %d", 3, *v)
		}
		val.ProprietaryId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMediaControl-Unknown: %w", err)
	} else {
		val.Unknown = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMediaControl-SourceId: %w", err)
	} else {
		val.SourceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMediaControl-Command: %w", err)
	} else {
		val.Command = FusionCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionSiriusControl struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId *uint8 `json:"proprietaryId"`
	Unknown *uint8 `json:"unknown"`
	SourceId *uint8 `json:"sourceId"`
	Command FusionSiriusCommandConst `json:"command"`
}
func DecodeFusionSiriusControl(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionSiriusControl
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSiriusControl-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionSiriusControl-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionSiriusControl-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionSiriusControl-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSiriusControl-ProprietaryId: %w", err)
	} else {
		if v != nil && *v != 30 {
			return nil, fmt.Errorf("match failed for FusionSiriusControl-ProprietaryId: Expected %d != %d", 30, *v)
		}
		val.ProprietaryId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSiriusControl-Unknown: %w", err)
	} else {
		val.Unknown = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSiriusControl-SourceId: %w", err)
	} else {
		val.SourceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSiriusControl-Command: %w", err)
	} else {
		val.Command = FusionSiriusCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionRequestStatus struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId FusionMessageIdConst `json:"proprietaryId"`
	Unknown *uint8 `json:"unknown"`
}
func DecodeFusionRequestStatus(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionRequestStatus
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionRequestStatus-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionRequestStatus-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionRequestStatus-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionRequestStatus-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionRequestStatus-ProprietaryId: %w", err)
	} else {
		if v != 1 {
			return nil, fmt.Errorf("match failed for FusionRequestStatus-ProprietaryId: Expected %d != %d", 1, v)
		}
		val.ProprietaryId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionRequestStatus-Unknown: %w", err)
	} else {
		val.Unknown = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionSetSource struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId FusionMessageIdConst `json:"proprietaryId"`
	Unknown *uint8 `json:"unknown"`
	SourceId *uint8 `json:"sourceId"`
}
func DecodeFusionSetSource(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionSetSource
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetSource-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionSetSource-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionSetSource-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionSetSource-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetSource-ProprietaryId: %w", err)
	} else {
		if v != 2 {
			return nil, fmt.Errorf("match failed for FusionSetSource-ProprietaryId: Expected %d != %d", 2, v)
		}
		val.ProprietaryId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetSource-Unknown: %w", err)
	} else {
		val.Unknown = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetSource-SourceId: %w", err)
	} else {
		val.SourceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionSetMute struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId FusionMessageIdConst `json:"proprietaryId"`
	Command FusionMuteCommandConst `json:"command"`
}
func DecodeFusionSetMute(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionSetMute
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetMute-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionSetMute-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionSetMute-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionSetMute-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetMute-ProprietaryId: %w", err)
	} else {
		if v != 23 {
			return nil, fmt.Errorf("match failed for FusionSetMute-ProprietaryId: Expected %d != %d", 23, v)
		}
		val.ProprietaryId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetMute-Command: %w", err)
	} else {
		val.Command = FusionMuteCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionSetZoneVolume struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId FusionMessageIdConst `json:"proprietaryId"`
	Unknown *uint8 `json:"unknown"`
	Zone *uint8 `json:"zone"`
	Volume *uint8 `json:"volume"`
}
func DecodeFusionSetZoneVolume(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionSetZoneVolume
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetZoneVolume-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionSetZoneVolume-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionSetZoneVolume-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionSetZoneVolume-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetZoneVolume-ProprietaryId: %w", err)
	} else {
		if v != 24 {
			return nil, fmt.Errorf("match failed for FusionSetZoneVolume-ProprietaryId: Expected %d != %d", 24, v)
		}
		val.ProprietaryId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetZoneVolume-Unknown: %w", err)
	} else {
		val.Unknown = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetZoneVolume-Zone: %w", err)
	} else {
		val.Zone = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetZoneVolume-Volume: %w", err)
	} else {
		val.Volume = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionSetAllVolumes struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId FusionMessageIdConst `json:"proprietaryId"`
	Unknown *uint8 `json:"unknown"`
	Zone1 *uint8 `json:"zone1"`
	Zone2 *uint8 `json:"zone2"`
	Zone3 *uint8 `json:"zone3"`
	Zone4 *uint8 `json:"zone4"`
}
func DecodeFusionSetAllVolumes(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionSetAllVolumes
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetAllVolumes-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionSetAllVolumes-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionSetAllVolumes-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionSetAllVolumes-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetAllVolumes-ProprietaryId: %w", err)
	} else {
		if v != 25 {
			return nil, fmt.Errorf("match failed for FusionSetAllVolumes-ProprietaryId: Expected %d != %d", 25, v)
		}
		val.ProprietaryId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetAllVolumes-Unknown: %w", err)
	} else {
		val.Unknown = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetAllVolumes-Zone1: %w", err)
	} else {
		val.Zone1 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetAllVolumes-Zone2: %w", err)
	} else {
		val.Zone2 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetAllVolumes-Zone3: %w", err)
	} else {
		val.Zone3 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSetAllVolumes-Zone4: %w", err)
	} else {
		val.Zone4 = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SonichubZoneInfo struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SonichubCommandConst `json:"proprietaryId"`
	Control SonichubControlConst `json:"control"`
	Zone *uint8 `json:"zone"`
}
func DecodeSonichubZoneInfo(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SonichubZoneInfo
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubZoneInfo-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubZoneInfo-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubZoneInfo-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubZoneInfo-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubZoneInfo-ProprietaryId: %w", err)
	} else {
		if v != 5 {
			return nil, fmt.Errorf("match failed for SonichubZoneInfo-ProprietaryId: Expected %d != %d", 5, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubZoneInfo-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubZoneInfo-Zone: %w", err)
	} else {
		val.Zone = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SonichubSource struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SonichubCommandConst `json:"proprietaryId"`
	Control SonichubControlConst `json:"control"`
	Source SonichubSourceConst `json:"source"`
}
func DecodeSonichubSource(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SonichubSource
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubSource-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubSource-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubSource-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubSource-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubSource-ProprietaryId: %w", err)
	} else {
		if v != 6 {
			return nil, fmt.Errorf("match failed for SonichubSource-ProprietaryId: Expected %d != %d", 6, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubSource-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubSource-Source: %w", err)
	} else {
		val.Source = SonichubSourceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SonichubSourceList struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SonichubCommandConst `json:"proprietaryId"`
	Control SonichubControlConst `json:"control"`
	SourceId *uint8 `json:"sourceId"`
	A *uint8 `json:"a"`
	Text string `json:"text"`
}
func DecodeSonichubSourceList(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SonichubSourceList
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubSourceList-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubSourceList-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubSourceList-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubSourceList-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubSourceList-ProprietaryId: %w", err)
	} else {
		if v != 8 {
			return nil, fmt.Errorf("match failed for SonichubSourceList-ProprietaryId: Expected %d != %d", 8, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubSourceList-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubSourceList-SourceId: %w", err)
	} else {
		val.SourceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubSourceList-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubSourceList-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SonichubControl struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SonichubCommandConst `json:"proprietaryId"`
	Control SonichubControlConst `json:"control"`
	Item FusionMuteCommandConst `json:"item"`
}
func DecodeSonichubControl(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SonichubControl
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubControl-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubControl-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubControl-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubControl-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubControl-ProprietaryId: %w", err)
	} else {
		if v != 9 {
			return nil, fmt.Errorf("match failed for SonichubControl-ProprietaryId: Expected %d != %d", 9, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubControl-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubControl-Item: %w", err)
	} else {
		val.Item = FusionMuteCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SonichubFmRadio struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SonichubCommandConst `json:"proprietaryId"`
	Control SonichubControlConst `json:"control"`
	Item SonichubTuningConst `json:"item"`
	Frequency *uint32 `json:"frequency"`
	NoiseLevel *uint8 `json:"noiseLevel"`
	SignalLevel *uint8 `json:"signalLevel"`
	Text string `json:"text"`
}
func DecodeSonichubFmRadio(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SonichubFmRadio
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubFmRadio-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubFmRadio-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-ProprietaryId: %w", err)
	} else {
		if v != 12 {
			return nil, fmt.Errorf("match failed for SonichubFmRadio-ProprietaryId: Expected %d != %d", 12, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-Item: %w", err)
	} else {
		val.Item = SonichubTuningConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-Frequency: %w", err)
	} else {
		val.Frequency = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(2); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-NoiseLevel: %w", err)
	} else {
		val.NoiseLevel = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-SignalLevel: %w", err)
	} else {
		val.SignalLevel = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubFmRadio-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SonichubPlaylist struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SonichubCommandConst `json:"proprietaryId"`
	Control SonichubControlConst `json:"control"`
	Item SonichubPlaylistConst `json:"item"`
	A *uint8 `json:"a"`
	CurrentTrack *uint32 `json:"currentTrack"`
	Tracks *uint32 `json:"tracks"`
	Length *float32 `json:"length"`
	PositionInTrack *float32 `json:"positionInTrack"`
}
func DecodeSonichubPlaylist(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SonichubPlaylist
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubPlaylist-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubPlaylist-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-ProprietaryId: %w", err)
	} else {
		if v != 13 {
			return nil, fmt.Errorf("match failed for SonichubPlaylist-ProprietaryId: Expected %d != %d", 13, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-Item: %w", err)
	} else {
		val.Item = SonichubPlaylistConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-CurrentTrack: %w", err)
	} else {
		val.CurrentTrack = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-Tracks: %w", err)
	} else {
		val.Tracks = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-Length: %w", err)
	} else {
		val.Length = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPlaylist-PositionInTrack: %w", err)
	} else {
		val.PositionInTrack = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SonichubTrack struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SonichubCommandConst `json:"proprietaryId"`
	Control SonichubControlConst `json:"control"`
	Item *uint32 `json:"item"`
	Text string `json:"text"`
}
func DecodeSonichubTrack(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SonichubTrack
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubTrack-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubTrack-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubTrack-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubTrack-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubTrack-ProprietaryId: %w", err)
	} else {
		if v != 14 {
			return nil, fmt.Errorf("match failed for SonichubTrack-ProprietaryId: Expected %d != %d", 14, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubTrack-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubTrack-Item: %w", err)
	} else {
		val.Item = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubTrack-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SonichubArtist struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SonichubCommandConst `json:"proprietaryId"`
	Control SonichubControlConst `json:"control"`
	Item *uint32 `json:"item"`
	Text string `json:"text"`
}
func DecodeSonichubArtist(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SonichubArtist
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubArtist-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubArtist-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubArtist-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubArtist-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubArtist-ProprietaryId: %w", err)
	} else {
		if v != 15 {
			return nil, fmt.Errorf("match failed for SonichubArtist-ProprietaryId: Expected %d != %d", 15, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubArtist-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubArtist-Item: %w", err)
	} else {
		val.Item = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubArtist-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SonichubAlbum struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SonichubCommandConst `json:"proprietaryId"`
	Control SonichubControlConst `json:"control"`
	Item *uint32 `json:"item"`
	Text string `json:"text"`
}
func DecodeSonichubAlbum(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SonichubAlbum
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubAlbum-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubAlbum-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubAlbum-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubAlbum-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubAlbum-ProprietaryId: %w", err)
	} else {
		if v != 16 {
			return nil, fmt.Errorf("match failed for SonichubAlbum-ProprietaryId: Expected %d != %d", 16, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubAlbum-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubAlbum-Item: %w", err)
	} else {
		val.Item = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubAlbum-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SonichubMenuItem struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SonichubCommandConst `json:"proprietaryId"`
	Control SonichubControlConst `json:"control"`
	Item *uint32 `json:"item"`
	C *uint8 `json:"c"`
	D *uint8 `json:"d"`
	E *uint8 `json:"e"`
	Text string `json:"text"`
}
func DecodeSonichubMenuItem(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SonichubMenuItem
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubMenuItem-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubMenuItem-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-ProprietaryId: %w", err)
	} else {
		if v != 19 {
			return nil, fmt.Errorf("match failed for SonichubMenuItem-ProprietaryId: Expected %d != %d", 19, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-Item: %w", err)
	} else {
		val.Item = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMenuItem-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SonichubZones struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SonichubCommandConst `json:"proprietaryId"`
	Control SonichubControlConst `json:"control"`
	Zones *uint8 `json:"zones"`
}
func DecodeSonichubZones(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SonichubZones
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubZones-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubZones-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubZones-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubZones-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubZones-ProprietaryId: %w", err)
	} else {
		if v != 20 {
			return nil, fmt.Errorf("match failed for SonichubZones-ProprietaryId: Expected %d != %d", 20, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubZones-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubZones-Zones: %w", err)
	} else {
		val.Zones = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SonichubMaxVolume struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SonichubCommandConst `json:"proprietaryId"`
	Control SonichubControlConst `json:"control"`
	Zone *uint8 `json:"zone"`
	Level *uint8 `json:"level"`
}
func DecodeSonichubMaxVolume(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SonichubMaxVolume
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMaxVolume-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubMaxVolume-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubMaxVolume-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubMaxVolume-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubMaxVolume-ProprietaryId: %w", err)
	} else {
		if v != 23 {
			return nil, fmt.Errorf("match failed for SonichubMaxVolume-ProprietaryId: Expected %d != %d", 23, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMaxVolume-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMaxVolume-Zone: %w", err)
	} else {
		val.Zone = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubMaxVolume-Level: %w", err)
	} else {
		val.Level = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SonichubVolume struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SonichubCommandConst `json:"proprietaryId"`
	Control SonichubControlConst `json:"control"`
	Zone *uint8 `json:"zone"`
	Level *uint8 `json:"level"`
}
func DecodeSonichubVolume(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SonichubVolume
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubVolume-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubVolume-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubVolume-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubVolume-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubVolume-ProprietaryId: %w", err)
	} else {
		if v != 24 {
			return nil, fmt.Errorf("match failed for SonichubVolume-ProprietaryId: Expected %d != %d", 24, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubVolume-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubVolume-Zone: %w", err)
	} else {
		val.Zone = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubVolume-Level: %w", err)
	} else {
		val.Level = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SonichubInit1 struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SonichubCommandConst `json:"proprietaryId"`
	Control SonichubControlConst `json:"control"`
}
func DecodeSonichubInit1(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SonichubInit1
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubInit1-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubInit1-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubInit1-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubInit1-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubInit1-ProprietaryId: %w", err)
	} else {
		if v != 25 {
			return nil, fmt.Errorf("match failed for SonichubInit1-ProprietaryId: Expected %d != %d", 25, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubInit1-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type SonichubPosition struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	ProprietaryId SonichubCommandConst `json:"proprietaryId"`
	Control SonichubControlConst `json:"control"`
	Position *float32 `json:"position"`
}
func DecodeSonichubPosition(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val SonichubPosition
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPosition-ManufacturerCode: %w", err)
	} else {
		if v != 275 {
			return nil, fmt.Errorf("match failed for SonichubPosition-ManufacturerCode: Expected %d != %d", 275, v)
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
		return nil, fmt.Errorf("parse failed for SonichubPosition-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for SonichubPosition-IndustryCode: Expected %d != %d", 4, v)
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
		return nil, fmt.Errorf("parse failed for SonichubPosition-ProprietaryId: %w", err)
	} else {
		if v != 48 {
			return nil, fmt.Errorf("match failed for SonichubPosition-ProprietaryId: Expected %d != %d", 48, v)
		}
		val.ProprietaryId = SonichubCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPosition-Control: %w", err)
	} else {
		val.Control = SonichubControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for SonichubPosition-Position: %w", err)
	} else {
		val.Position = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionSourceName struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId FusionMessageIdConst `json:"messageId"`
	A *uint8 `json:"a"`
	SourceId *uint8 `json:"sourceId"`
	CurrentSourceId *uint8 `json:"currentSourceId"`
	D *uint8 `json:"d"`
	E *uint8 `json:"e"`
	Source string `json:"source"`
}
func DecodeFusionSourceName(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionSourceName
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSourceName-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionSourceName-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionSourceName-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionSourceName-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSourceName-MessageId: %w", err)
	} else {
		if v != 2 {
			return nil, fmt.Errorf("match failed for FusionSourceName-MessageId: Expected %d != %d", 2, v)
		}
		val.MessageId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSourceName-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSourceName-SourceId: %w", err)
	} else {
		val.SourceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSourceName-CurrentSourceId: %w", err)
	} else {
		val.CurrentSourceId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSourceName-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSourceName-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSourceName-Source: %w", err)
	} else {
		val.Source = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionTrackInfo struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId FusionMessageIdConst `json:"messageId"`
	A *uint16 `json:"a"`
	Transport EntertainmentPlayStatusConst `json:"transport"`
	X *uint8 `json:"x"`
	B *uint8 `json:"b"`
	Track *uint16 `json:"track"`
	C *uint16 `json:"c"`
	TrackCount *uint16 `json:"trackCount"`
	E *uint16 `json:"e"`
	Length *float32 `json:"length"`
	PositionInTrack *float32 `json:"positionInTrack"`
	H *uint16 `json:"h"`
}
func DecodeFusionTrackInfo(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionTrackInfo
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrackInfo-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionTrackInfo-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionTrackInfo-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionTrackInfo-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrackInfo-MessageId: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionTrackInfo-MessageId: Expected %d != %d", 4, v)
		}
		val.MessageId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrackInfo-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrackInfo-Transport: %w", err)
	} else {
		val.Transport = EntertainmentPlayStatusConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrackInfo-X: %w", err)
	} else {
		val.X = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrackInfo-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrackInfo-Track: %w", err)
	} else {
		val.Track = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrackInfo-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrackInfo-TrackCount: %w", err)
	} else {
		val.TrackCount = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrackInfo-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(24, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrackInfo-Length: %w", err)
	} else {
		val.Length = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(24, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrackInfo-PositionInTrack: %w", err)
	} else {
		val.PositionInTrack = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrackInfo-H: %w", err)
	} else {
		val.H = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionTrack struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId FusionMessageIdConst `json:"messageId"`
	A *uint8 `json:"a"`
	B *uint64 `json:"b"`
	Track string `json:"track"`
}
func DecodeFusionTrack(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionTrack
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrack-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionTrack-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionTrack-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionTrack-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrack-MessageId: %w", err)
	} else {
		if v != 5 {
			return nil, fmt.Errorf("match failed for FusionTrack-MessageId: Expected %d != %d", 5, v)
		}
		val.MessageId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrack-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt64(40); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrack-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for FusionTrack-Track: %w", err)
	} else {
		val.Track = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionArtist struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId FusionMessageIdConst `json:"messageId"`
	A *uint8 `json:"a"`
	B *uint64 `json:"b"`
	Artist string `json:"artist"`
}
func DecodeFusionArtist(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionArtist
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionArtist-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionArtist-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionArtist-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionArtist-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionArtist-MessageId: %w", err)
	} else {
		if v != 6 {
			return nil, fmt.Errorf("match failed for FusionArtist-MessageId: Expected %d != %d", 6, v)
		}
		val.MessageId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionArtist-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt64(40); err != nil {
		return nil, fmt.Errorf("parse failed for FusionArtist-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for FusionArtist-Artist: %w", err)
	} else {
		val.Artist = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionAlbum struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId FusionMessageIdConst `json:"messageId"`
	A *uint8 `json:"a"`
	B *uint64 `json:"b"`
	Album string `json:"album"`
}
func DecodeFusionAlbum(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionAlbum
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionAlbum-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionAlbum-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionAlbum-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionAlbum-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionAlbum-MessageId: %w", err)
	} else {
		if v != 7 {
			return nil, fmt.Errorf("match failed for FusionAlbum-MessageId: Expected %d != %d", 7, v)
		}
		val.MessageId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionAlbum-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt64(40); err != nil {
		return nil, fmt.Errorf("parse failed for FusionAlbum-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for FusionAlbum-Album: %w", err)
	} else {
		val.Album = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionUnitName struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId FusionMessageIdConst `json:"messageId"`
	A *uint8 `json:"a"`
	Name string `json:"name"`
}
func DecodeFusionUnitName(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionUnitName
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionUnitName-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionUnitName-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionUnitName-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionUnitName-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionUnitName-MessageId: %w", err)
	} else {
		if v != 33 {
			return nil, fmt.Errorf("match failed for FusionUnitName-MessageId: Expected %d != %d", 33, v)
		}
		val.MessageId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionUnitName-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for FusionUnitName-Name: %w", err)
	} else {
		val.Name = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionZoneName struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId FusionMessageIdConst `json:"messageId"`
	A *uint8 `json:"a"`
	Number *uint8 `json:"number"`
	Name string `json:"name"`
}
func DecodeFusionZoneName(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionZoneName
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionZoneName-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionZoneName-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionZoneName-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionZoneName-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionZoneName-MessageId: %w", err)
	} else {
		if v != 45 {
			return nil, fmt.Errorf("match failed for FusionZoneName-MessageId: Expected %d != %d", 45, v)
		}
		val.MessageId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionZoneName-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionZoneName-Number: %w", err)
	} else {
		val.Number = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for FusionZoneName-Name: %w", err)
	} else {
		val.Name = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionPlayProgress struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId FusionMessageIdConst `json:"messageId"`
	A *uint8 `json:"a"`
	B *uint8 `json:"b"`
	Progress *float32 `json:"progress"`
}
func DecodeFusionPlayProgress(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionPlayProgress
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionPlayProgress-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionPlayProgress-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionPlayProgress-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionPlayProgress-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionPlayProgress-MessageId: %w", err)
	} else {
		if v != 9 {
			return nil, fmt.Errorf("match failed for FusionPlayProgress-MessageId: Expected %d != %d", 9, v)
		}
		val.MessageId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionPlayProgress-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionPlayProgress-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(24, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for FusionPlayProgress-Progress: %w", err)
	} else {
		val.Progress = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionAmFmStation struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId FusionMessageIdConst `json:"messageId"`
	A *uint8 `json:"a"`
	AmFm FusionRadioSourceConst `json:"amFm"`
	B *uint8 `json:"b"`
	Frequency *uint32 `json:"frequency"`
	C *uint8 `json:"c"`
	Track string `json:"track"`
}
func DecodeFusionAmFmStation(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionAmFmStation
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionAmFmStation-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionAmFmStation-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionAmFmStation-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionAmFmStation-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionAmFmStation-MessageId: %w", err)
	} else {
		if v != 11 {
			return nil, fmt.Errorf("match failed for FusionAmFmStation-MessageId: Expected %d != %d", 11, v)
		}
		val.MessageId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionAmFmStation-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionAmFmStation-AmFm: %w", err)
	} else {
		val.AmFm = FusionRadioSourceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionAmFmStation-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for FusionAmFmStation-Frequency: %w", err)
	} else {
		val.Frequency = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionAmFmStation-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for FusionAmFmStation-Track: %w", err)
	} else {
		val.Track = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionVhf struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId FusionMessageIdConst `json:"messageId"`
	A *uint8 `json:"a"`
	B *uint8 `json:"b"`
	Channel *uint8 `json:"channel"`
	D *uint32 `json:"d"`
}
func DecodeFusionVhf(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionVhf
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionVhf-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionVhf-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionVhf-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionVhf-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionVhf-MessageId: %w", err)
	} else {
		if v != 12 {
			return nil, fmt.Errorf("match failed for FusionVhf-MessageId: Expected %d != %d", 12, v)
		}
		val.MessageId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionVhf-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionVhf-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionVhf-Channel: %w", err)
	} else {
		val.Channel = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(24); err != nil {
		return nil, fmt.Errorf("parse failed for FusionVhf-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionSquelch struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId FusionMessageIdConst `json:"messageId"`
	A *uint8 `json:"a"`
	B *uint8 `json:"b"`
	Squelch *uint8 `json:"squelch"`
}
func DecodeFusionSquelch(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionSquelch
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSquelch-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionSquelch-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionSquelch-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionSquelch-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSquelch-MessageId: %w", err)
	} else {
		if v != 13 {
			return nil, fmt.Errorf("match failed for FusionSquelch-MessageId: Expected %d != %d", 13, v)
		}
		val.MessageId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSquelch-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSquelch-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSquelch-Squelch: %w", err)
	} else {
		val.Squelch = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionScan struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId FusionMessageIdConst `json:"messageId"`
	A *uint8 `json:"a"`
	B *uint8 `json:"b"`
	Scan YesNoConst `json:"scan"`
	C *uint8 `json:"c"`
}
func DecodeFusionScan(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionScan
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionScan-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionScan-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionScan-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionScan-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionScan-MessageId: %w", err)
	} else {
		if v != 14 {
			return nil, fmt.Errorf("match failed for FusionScan-MessageId: Expected %d != %d", 14, v)
		}
		val.MessageId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionScan-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionScan-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for FusionScan-Scan: %w", err)
	} else {
		val.Scan = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(6); err != nil {
		return nil, fmt.Errorf("parse failed for FusionScan-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionMenuItem struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId FusionMessageIdConst `json:"messageId"`
	A *uint8 `json:"a"`
	B *uint8 `json:"b"`
	Line *uint8 `json:"line"`
	E *uint8 `json:"e"`
	F *uint8 `json:"f"`
	G *uint8 `json:"g"`
	H *uint8 `json:"h"`
	I *uint8 `json:"i"`
	Text string `json:"text"`
}
func DecodeFusionMenuItem(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionMenuItem
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMenuItem-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionMenuItem-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionMenuItem-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionMenuItem-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMenuItem-MessageId: %w", err)
	} else {
		if v != 17 {
			return nil, fmt.Errorf("match failed for FusionMenuItem-MessageId: Expected %d != %d", 17, v)
		}
		val.MessageId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMenuItem-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMenuItem-B: %w", err)
	} else {
		val.B = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMenuItem-Line: %w", err)
	} else {
		val.Line = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMenuItem-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMenuItem-F: %w", err)
	} else {
		val.F = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMenuItem-G: %w", err)
	} else {
		val.G = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMenuItem-H: %w", err)
	} else {
		val.H = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMenuItem-I: %w", err)
	} else {
		val.I = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLength(); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMenuItem-Text: %w", err)
	} else {
		val.Text = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionReplay struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId FusionMessageIdConst `json:"messageId"`
	A *uint8 `json:"a"`
	Mode FusionReplayModeConst `json:"mode"`
	C *uint32 `json:"c"`
	D *uint8 `json:"d"`
	E *uint8 `json:"e"`
	Status FusionReplayStatusConst `json:"status"`
	H *uint8 `json:"h"`
	I *uint8 `json:"i"`
	J *uint8 `json:"j"`
}
func DecodeFusionReplay(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionReplay
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionReplay-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionReplay-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionReplay-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionReplay-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionReplay-MessageId: %w", err)
	} else {
		if v != 20 {
			return nil, fmt.Errorf("match failed for FusionReplay-MessageId: Expected %d != %d", 20, v)
		}
		val.MessageId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionReplay-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionReplay-Mode: %w", err)
	} else {
		val.Mode = FusionReplayModeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(24); err != nil {
		return nil, fmt.Errorf("parse failed for FusionReplay-C: %w", err)
	} else {
		val.C = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionReplay-D: %w", err)
	} else {
		val.D = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionReplay-E: %w", err)
	} else {
		val.E = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionReplay-Status: %w", err)
	} else {
		val.Status = FusionReplayStatusConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionReplay-H: %w", err)
	} else {
		val.H = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionReplay-I: %w", err)
	} else {
		val.I = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionReplay-J: %w", err)
	} else {
		val.J = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionMute struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId FusionMessageIdConst `json:"messageId"`
	A *uint8 `json:"a"`
	Mute FusionMuteCommandConst `json:"mute"`
}
func DecodeFusionMute(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionMute
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMute-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionMute-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionMute-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionMute-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMute-MessageId: %w", err)
	} else {
		if v != 23 {
			return nil, fmt.Errorf("match failed for FusionMute-MessageId: Expected %d != %d", 23, v)
		}
		val.MessageId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMute-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionMute-Mute: %w", err)
	} else {
		val.Mute = FusionMuteCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}
type FusionSubVolume struct {
	Info MessageInfo `json:"info"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	MessageId FusionMessageIdConst `json:"messageId"`
	A *uint8 `json:"a"`
	Zone1 *uint8 `json:"zone1"`
	Zone2 *uint8 `json:"zone2"`
	Zone3 *uint8 `json:"zone3"`
	Zone4 *uint8 `json:"zone4"`
}
func DecodeFusionSubVolume(Info MessageInfo, stream *PGNDataStream) (any, error) {
	var val FusionSubVolume
	val.Info = Info
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSubVolume-ManufacturerCode: %w", err)
	} else {
		if v != 419 {
			return nil, fmt.Errorf("match failed for FusionSubVolume-ManufacturerCode: Expected %d != %d", 419, v)
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
		return nil, fmt.Errorf("parse failed for FusionSubVolume-IndustryCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for FusionSubVolume-IndustryCode: Expected %d != %d", 4, v)
		}
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSubVolume-MessageId: %w", err)
	} else {
		if v != 26 {
			return nil, fmt.Errorf("match failed for FusionSubVolume-MessageId: Expected %d != %d", 26, v)
		}
		val.MessageId = FusionMessageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSubVolume-A: %w", err)
	} else {
		val.A = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSubVolume-Zone1: %w", err)
	} else {
		val.Zone1 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSubVolume-Zone2: %w", err)
	} else {
		val.Zone2 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSubVolume-Zone3: %w", err)
	} else {
		val.Zone3 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for FusionSubVolume-Zone4: %w", err)
	} else {
		val.Zone4 = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeFusionMediaControl(val *FusionMediaControl) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.ProprietaryId, 8)
	w.writeUInt8(val.Unknown, 8)
	w.writeUInt8(val.SourceId, 8)
	w.writeLookupField(uint64(val.Command), 8)
	return w.Bytes(), nil
}
func encodeFusionMediaControlAny(v any) ([]byte, error) {
	val, ok := v.(*FusionMediaControl)
	if !ok {
		return nil, fmt.Errorf("expected *FusionMediaControl, got %T", v)
	}
	return EncodeFusionMediaControl(val)
}

func EncodeFusionSiriusControl(val *FusionSiriusControl) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeUInt8(val.ProprietaryId, 8)
	w.writeUInt8(val.Unknown, 8)
	w.writeUInt8(val.SourceId, 8)
	w.writeLookupField(uint64(val.Command), 8)
	return w.Bytes(), nil
}
func encodeFusionSiriusControlAny(v any) ([]byte, error) {
	val, ok := v.(*FusionSiriusControl)
	if !ok {
		return nil, fmt.Errorf("expected *FusionSiriusControl, got %T", v)
	}
	return EncodeFusionSiriusControl(val)
}

func EncodeFusionRequestStatus(val *FusionRequestStatus) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt8(val.Unknown, 8)
	return w.Bytes(), nil
}
func encodeFusionRequestStatusAny(v any) ([]byte, error) {
	val, ok := v.(*FusionRequestStatus)
	if !ok {
		return nil, fmt.Errorf("expected *FusionRequestStatus, got %T", v)
	}
	return EncodeFusionRequestStatus(val)
}

func EncodeFusionSetSource(val *FusionSetSource) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt8(val.Unknown, 8)
	w.writeUInt8(val.SourceId, 8)
	return w.Bytes(), nil
}
func encodeFusionSetSourceAny(v any) ([]byte, error) {
	val, ok := v.(*FusionSetSource)
	if !ok {
		return nil, fmt.Errorf("expected *FusionSetSource, got %T", v)
	}
	return EncodeFusionSetSource(val)
}

func EncodeFusionSetMute(val *FusionSetMute) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Command), 8)
	return w.Bytes(), nil
}
func encodeFusionSetMuteAny(v any) ([]byte, error) {
	val, ok := v.(*FusionSetMute)
	if !ok {
		return nil, fmt.Errorf("expected *FusionSetMute, got %T", v)
	}
	return EncodeFusionSetMute(val)
}

func EncodeFusionSetZoneVolume(val *FusionSetZoneVolume) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt8(val.Unknown, 8)
	w.writeUInt8(val.Zone, 8)
	w.writeUInt8(val.Volume, 8)
	return w.Bytes(), nil
}
func encodeFusionSetZoneVolumeAny(v any) ([]byte, error) {
	val, ok := v.(*FusionSetZoneVolume)
	if !ok {
		return nil, fmt.Errorf("expected *FusionSetZoneVolume, got %T", v)
	}
	return EncodeFusionSetZoneVolume(val)
}

func EncodeFusionSetAllVolumes(val *FusionSetAllVolumes) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeUInt8(val.Unknown, 8)
	w.writeUInt8(val.Zone1, 8)
	w.writeUInt8(val.Zone2, 8)
	w.writeUInt8(val.Zone3, 8)
	w.writeUInt8(val.Zone4, 8)
	return w.Bytes(), nil
}
func encodeFusionSetAllVolumesAny(v any) ([]byte, error) {
	val, ok := v.(*FusionSetAllVolumes)
	if !ok {
		return nil, fmt.Errorf("expected *FusionSetAllVolumes, got %T", v)
	}
	return EncodeFusionSetAllVolumes(val)
}

func EncodeSonichubZoneInfo(val *SonichubZoneInfo) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt8(val.Zone, 8)
	return w.Bytes(), nil
}
func encodeSonichubZoneInfoAny(v any) ([]byte, error) {
	val, ok := v.(*SonichubZoneInfo)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubZoneInfo, got %T", v)
	}
	return EncodeSonichubZoneInfo(val)
}

func EncodeSonichubSource(val *SonichubSource) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeLookupField(uint64(val.Source), 8)
	return w.Bytes(), nil
}
func encodeSonichubSourceAny(v any) ([]byte, error) {
	val, ok := v.(*SonichubSource)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubSource, got %T", v)
	}
	return EncodeSonichubSource(val)
}

func EncodeSonichubSourceList(val *SonichubSourceList) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt8(val.SourceId, 8)
	w.writeUInt8(val.A, 8)
	w.writeStringWithLength(val.Text)
	return w.Bytes(), nil
}
func encodeSonichubSourceListAny(v any) ([]byte, error) {
	val, ok := v.(*SonichubSourceList)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubSourceList, got %T", v)
	}
	return EncodeSonichubSourceList(val)
}

func EncodeSonichubControl(val *SonichubControl) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeLookupField(uint64(val.Item), 8)
	return w.Bytes(), nil
}
func encodeSonichubControlAny(v any) ([]byte, error) {
	val, ok := v.(*SonichubControl)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubControl, got %T", v)
	}
	return EncodeSonichubControl(val)
}

func EncodeSonichubFmRadio(val *SonichubFmRadio) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeLookupField(uint64(val.Item), 8)
	w.writeUInt32(val.Frequency, 32)
	w.writeUInt8(val.NoiseLevel, 2)
	w.writeUInt8(val.SignalLevel, 4)
	w.skipBits(2)
	w.writeStringWithLength(val.Text)
	return w.Bytes(), nil
}
func encodeSonichubFmRadioAny(v any) ([]byte, error) {
	val, ok := v.(*SonichubFmRadio)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubFmRadio, got %T", v)
	}
	return EncodeSonichubFmRadio(val)
}

func EncodeSonichubPlaylist(val *SonichubPlaylist) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeLookupField(uint64(val.Item), 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt32(val.CurrentTrack, 32)
	w.writeUInt32(val.Tracks, 32)
	w.writeUnsignedResolution(val.Length, 32, 0.001)
	w.writeUnsignedResolution(val.PositionInTrack, 32, 0.001)
	return w.Bytes(), nil
}
func encodeSonichubPlaylistAny(v any) ([]byte, error) {
	val, ok := v.(*SonichubPlaylist)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubPlaylist, got %T", v)
	}
	return EncodeSonichubPlaylist(val)
}

func EncodeSonichubTrack(val *SonichubTrack) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt32(val.Item, 32)
	w.writeStringWithLength(val.Text)
	return w.Bytes(), nil
}
func encodeSonichubTrackAny(v any) ([]byte, error) {
	val, ok := v.(*SonichubTrack)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubTrack, got %T", v)
	}
	return EncodeSonichubTrack(val)
}

func EncodeSonichubArtist(val *SonichubArtist) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt32(val.Item, 32)
	w.writeStringWithLength(val.Text)
	return w.Bytes(), nil
}
func encodeSonichubArtistAny(v any) ([]byte, error) {
	val, ok := v.(*SonichubArtist)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubArtist, got %T", v)
	}
	return EncodeSonichubArtist(val)
}

func EncodeSonichubAlbum(val *SonichubAlbum) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt32(val.Item, 32)
	w.writeStringWithLength(val.Text)
	return w.Bytes(), nil
}
func encodeSonichubAlbumAny(v any) ([]byte, error) {
	val, ok := v.(*SonichubAlbum)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubAlbum, got %T", v)
	}
	return EncodeSonichubAlbum(val)
}

func EncodeSonichubMenuItem(val *SonichubMenuItem) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt32(val.Item, 32)
	w.writeUInt8(val.C, 8)
	w.writeUInt8(val.D, 8)
	w.writeUInt8(val.E, 8)
	w.writeStringWithLength(val.Text)
	return w.Bytes(), nil
}
func encodeSonichubMenuItemAny(v any) ([]byte, error) {
	val, ok := v.(*SonichubMenuItem)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubMenuItem, got %T", v)
	}
	return EncodeSonichubMenuItem(val)
}

func EncodeSonichubZones(val *SonichubZones) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt8(val.Zones, 8)
	return w.Bytes(), nil
}
func encodeSonichubZonesAny(v any) ([]byte, error) {
	val, ok := v.(*SonichubZones)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubZones, got %T", v)
	}
	return EncodeSonichubZones(val)
}

func EncodeSonichubMaxVolume(val *SonichubMaxVolume) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt8(val.Zone, 8)
	w.writeUInt8(val.Level, 8)
	return w.Bytes(), nil
}
func encodeSonichubMaxVolumeAny(v any) ([]byte, error) {
	val, ok := v.(*SonichubMaxVolume)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubMaxVolume, got %T", v)
	}
	return EncodeSonichubMaxVolume(val)
}

func EncodeSonichubVolume(val *SonichubVolume) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt8(val.Zone, 8)
	w.writeUInt8(val.Level, 8)
	return w.Bytes(), nil
}
func encodeSonichubVolumeAny(v any) ([]byte, error) {
	val, ok := v.(*SonichubVolume)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubVolume, got %T", v)
	}
	return EncodeSonichubVolume(val)
}

func EncodeSonichubInit1(val *SonichubInit1) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	return w.Bytes(), nil
}
func encodeSonichubInit1Any(v any) ([]byte, error) {
	val, ok := v.(*SonichubInit1)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubInit1, got %T", v)
	}
	return EncodeSonichubInit1(val)
}

func EncodeSonichubPosition(val *SonichubPosition) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(8)
	w.writeLookupField(uint64(val.ProprietaryId), 8)
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUnsignedResolution(val.Position, 32, 0.001)
	return w.Bytes(), nil
}
func encodeSonichubPositionAny(v any) ([]byte, error) {
	val, ok := v.(*SonichubPosition)
	if !ok {
		return nil, fmt.Errorf("expected *SonichubPosition, got %T", v)
	}
	return EncodeSonichubPosition(val)
}

func EncodeFusionSourceName(val *FusionSourceName) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.MessageId), 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.SourceId, 8)
	w.writeUInt8(val.CurrentSourceId, 8)
	w.writeUInt8(val.D, 8)
	w.writeUInt8(val.E, 8)
	w.writeStringWithLength(val.Source)
	return w.Bytes(), nil
}
func encodeFusionSourceNameAny(v any) ([]byte, error) {
	val, ok := v.(*FusionSourceName)
	if !ok {
		return nil, fmt.Errorf("expected *FusionSourceName, got %T", v)
	}
	return EncodeFusionSourceName(val)
}

func EncodeFusionTrackInfo(val *FusionTrackInfo) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.MessageId), 8)
	w.writeUInt16(val.A, 16)
	w.writeLookupField(uint64(val.Transport), 4)
	w.writeUInt8(val.X, 4)
	w.writeUInt8(val.B, 8)
	w.writeUInt16(val.Track, 16)
	w.writeUInt16(val.C, 16)
	w.writeUInt16(val.TrackCount, 16)
	w.writeUInt16(val.E, 16)
	w.writeUnsignedResolution(val.Length, 24, 0.001)
	w.writeUnsignedResolution(val.PositionInTrack, 24, 0.001)
	w.writeUInt16(val.H, 16)
	return w.Bytes(), nil
}
func encodeFusionTrackInfoAny(v any) ([]byte, error) {
	val, ok := v.(*FusionTrackInfo)
	if !ok {
		return nil, fmt.Errorf("expected *FusionTrackInfo, got %T", v)
	}
	return EncodeFusionTrackInfo(val)
}

func EncodeFusionTrack(val *FusionTrack) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.MessageId), 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt64(val.B, 40)
	w.writeStringWithLength(val.Track)
	return w.Bytes(), nil
}
func encodeFusionTrackAny(v any) ([]byte, error) {
	val, ok := v.(*FusionTrack)
	if !ok {
		return nil, fmt.Errorf("expected *FusionTrack, got %T", v)
	}
	return EncodeFusionTrack(val)
}

func EncodeFusionArtist(val *FusionArtist) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.MessageId), 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt64(val.B, 40)
	w.writeStringWithLength(val.Artist)
	return w.Bytes(), nil
}
func encodeFusionArtistAny(v any) ([]byte, error) {
	val, ok := v.(*FusionArtist)
	if !ok {
		return nil, fmt.Errorf("expected *FusionArtist, got %T", v)
	}
	return EncodeFusionArtist(val)
}

func EncodeFusionAlbum(val *FusionAlbum) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.MessageId), 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt64(val.B, 40)
	w.writeStringWithLength(val.Album)
	return w.Bytes(), nil
}
func encodeFusionAlbumAny(v any) ([]byte, error) {
	val, ok := v.(*FusionAlbum)
	if !ok {
		return nil, fmt.Errorf("expected *FusionAlbum, got %T", v)
	}
	return EncodeFusionAlbum(val)
}

func EncodeFusionUnitName(val *FusionUnitName) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.MessageId), 8)
	w.writeUInt8(val.A, 8)
	w.writeStringWithLength(val.Name)
	return w.Bytes(), nil
}
func encodeFusionUnitNameAny(v any) ([]byte, error) {
	val, ok := v.(*FusionUnitName)
	if !ok {
		return nil, fmt.Errorf("expected *FusionUnitName, got %T", v)
	}
	return EncodeFusionUnitName(val)
}

func EncodeFusionZoneName(val *FusionZoneName) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.MessageId), 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.Number, 8)
	w.writeStringWithLength(val.Name)
	return w.Bytes(), nil
}
func encodeFusionZoneNameAny(v any) ([]byte, error) {
	val, ok := v.(*FusionZoneName)
	if !ok {
		return nil, fmt.Errorf("expected *FusionZoneName, got %T", v)
	}
	return EncodeFusionZoneName(val)
}

func EncodeFusionPlayProgress(val *FusionPlayProgress) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.MessageId), 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUnsignedResolution(val.Progress, 24, 0.001)
	return w.Bytes(), nil
}
func encodeFusionPlayProgressAny(v any) ([]byte, error) {
	val, ok := v.(*FusionPlayProgress)
	if !ok {
		return nil, fmt.Errorf("expected *FusionPlayProgress, got %T", v)
	}
	return EncodeFusionPlayProgress(val)
}

func EncodeFusionAmFmStation(val *FusionAmFmStation) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.MessageId), 8)
	w.writeUInt8(val.A, 8)
	w.writeLookupField(uint64(val.AmFm), 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt32(val.Frequency, 32)
	w.writeUInt8(val.C, 8)
	w.writeStringWithLength(val.Track)
	return w.Bytes(), nil
}
func encodeFusionAmFmStationAny(v any) ([]byte, error) {
	val, ok := v.(*FusionAmFmStation)
	if !ok {
		return nil, fmt.Errorf("expected *FusionAmFmStation, got %T", v)
	}
	return EncodeFusionAmFmStation(val)
}

func EncodeFusionVhf(val *FusionVhf) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.MessageId), 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.Channel, 8)
	w.writeUInt32(val.D, 24)
	return w.Bytes(), nil
}
func encodeFusionVhfAny(v any) ([]byte, error) {
	val, ok := v.(*FusionVhf)
	if !ok {
		return nil, fmt.Errorf("expected *FusionVhf, got %T", v)
	}
	return EncodeFusionVhf(val)
}

func EncodeFusionSquelch(val *FusionSquelch) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.MessageId), 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.Squelch, 8)
	return w.Bytes(), nil
}
func encodeFusionSquelchAny(v any) ([]byte, error) {
	val, ok := v.(*FusionSquelch)
	if !ok {
		return nil, fmt.Errorf("expected *FusionSquelch, got %T", v)
	}
	return EncodeFusionSquelch(val)
}

func EncodeFusionScan(val *FusionScan) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.MessageId), 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeLookupField(uint64(val.Scan), 2)
	w.writeUInt8(val.C, 6)
	return w.Bytes(), nil
}
func encodeFusionScanAny(v any) ([]byte, error) {
	val, ok := v.(*FusionScan)
	if !ok {
		return nil, fmt.Errorf("expected *FusionScan, got %T", v)
	}
	return EncodeFusionScan(val)
}

func EncodeFusionMenuItem(val *FusionMenuItem) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.MessageId), 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.B, 8)
	w.writeUInt8(val.Line, 8)
	w.writeUInt8(val.E, 8)
	w.writeUInt8(val.F, 8)
	w.writeUInt8(val.G, 8)
	w.writeUInt8(val.H, 8)
	w.writeUInt8(val.I, 8)
	w.writeStringWithLength(val.Text)
	return w.Bytes(), nil
}
func encodeFusionMenuItemAny(v any) ([]byte, error) {
	val, ok := v.(*FusionMenuItem)
	if !ok {
		return nil, fmt.Errorf("expected *FusionMenuItem, got %T", v)
	}
	return EncodeFusionMenuItem(val)
}

func EncodeFusionReplay(val *FusionReplay) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.MessageId), 8)
	w.writeUInt8(val.A, 8)
	w.writeLookupField(uint64(val.Mode), 8)
	w.writeUInt32(val.C, 24)
	w.writeUInt8(val.D, 8)
	w.writeUInt8(val.E, 8)
	w.writeLookupField(uint64(val.Status), 8)
	w.writeUInt8(val.H, 8)
	w.writeUInt8(val.I, 8)
	w.writeUInt8(val.J, 8)
	return w.Bytes(), nil
}
func encodeFusionReplayAny(v any) ([]byte, error) {
	val, ok := v.(*FusionReplay)
	if !ok {
		return nil, fmt.Errorf("expected *FusionReplay, got %T", v)
	}
	return EncodeFusionReplay(val)
}

func EncodeFusionMute(val *FusionMute) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.MessageId), 8)
	w.writeUInt8(val.A, 8)
	w.writeLookupField(uint64(val.Mute), 8)
	return w.Bytes(), nil
}
func encodeFusionMuteAny(v any) ([]byte, error) {
	val, ok := v.(*FusionMute)
	if !ok {
		return nil, fmt.Errorf("expected *FusionMute, got %T", v)
	}
	return EncodeFusionMute(val)
}

func EncodeFusionSubVolume(val *FusionSubVolume) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.skipBits(2)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.writeLookupField(uint64(val.MessageId), 8)
	w.writeUInt8(val.A, 8)
	w.writeUInt8(val.Zone1, 8)
	w.writeUInt8(val.Zone2, 8)
	w.writeUInt8(val.Zone3, 8)
	w.writeUInt8(val.Zone4, 8)
	return w.Bytes(), nil
}
func encodeFusionSubVolumeAny(v any) ([]byte, error) {
	val, ok := v.(*FusionSubVolume)
	if !ok {
		return nil, fmt.Errorf("expected *FusionSubVolume, got %T", v)
	}
	return EncodeFusionSubVolume(val)
}
