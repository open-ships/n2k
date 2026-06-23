package pgn

import (
	"fmt"

	"github.com/open-ships/n2k/units"
)

type IsoAcknowledgement struct {
	Info MessageInfo `json:"info"`
	Control IsoControlConst `json:"control"`
	GroupFunction *uint8 `json:"groupFunction"`
	Pgn *uint32 `json:"pgn"`
}
func (i *IsoAcknowledgement) PGNNumber() uint32  { return 59392 }
func DecodeIsoAcknowledgement(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &IsoAcknowledgement{}
	val.Info = Info
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoAcknowledgement-Control: %w", err)
	} else {
		val.Control = IsoControlConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoAcknowledgement-GroupFunction: %w", err)
	} else {
		val.GroupFunction = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(24)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt32(24); err != nil {
		return nil, fmt.Errorf("parse failed for IsoAcknowledgement-Pgn: %w", err)
	} else {
		val.Pgn = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeIsoAcknowledgement(val *IsoAcknowledgement) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.Control), 8)
	w.writeUInt8(val.GroupFunction, 8)
	w.skipBits(24)
	w.writeUInt32(val.Pgn, 24)
	return w.Bytes(), w.Err()
}
func encodeIsoAcknowledgementMsg(v Message) ([]byte, error) {
	val, ok := v.(*IsoAcknowledgement)
	if !ok {
		return nil, fmt.Errorf("expected *IsoAcknowledgement, got %T", v)
	}
	return EncodeIsoAcknowledgement(val)
}
type IsoRequest struct {
	Info MessageInfo `json:"info"`
	Pgn *uint32 `json:"pgn"`
}
func (i *IsoRequest) PGNNumber() uint32  { return 59904 }
func DecodeIsoRequest(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &IsoRequest{}
	val.Info = Info
	if v, err := stream.readUInt32(24); err != nil {
		return nil, fmt.Errorf("parse failed for IsoRequest-Pgn: %w", err)
	} else {
		val.Pgn = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeIsoRequest(val *IsoRequest) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt32(val.Pgn, 24)
	return w.Bytes(), w.Err()
}
func encodeIsoRequestMsg(v Message) ([]byte, error) {
	val, ok := v.(*IsoRequest)
	if !ok {
		return nil, fmt.Errorf("expected *IsoRequest, got %T", v)
	}
	return EncodeIsoRequest(val)
}
type IsoTransportProtocolDataTransfer struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	Data []uint8 `json:"data"`
}
func (i *IsoTransportProtocolDataTransfer) PGNNumber() uint32  { return 60160 }
func DecodeIsoTransportProtocolDataTransfer(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &IsoTransportProtocolDataTransfer{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolDataTransfer-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readBinaryData(56); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolDataTransfer-Data: %w", err)
	} else {
		val.Data = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeIsoTransportProtocolDataTransfer(val *IsoTransportProtocolDataTransfer) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeBinaryData(val.Data, 56)
	return w.Bytes(), w.Err()
}
func encodeIsoTransportProtocolDataTransferMsg(v Message) ([]byte, error) {
	val, ok := v.(*IsoTransportProtocolDataTransfer)
	if !ok {
		return nil, fmt.Errorf("expected *IsoTransportProtocolDataTransfer, got %T", v)
	}
	return EncodeIsoTransportProtocolDataTransfer(val)
}
type IsoTransportProtocolConnectionManagementRequestToSend struct {
	Info MessageInfo `json:"info"`
	GroupFunctionCode IsoCommandConst `json:"groupFunctionCode"`
	MessageSize *uint16 `json:"messageSize"`
	Packets *uint8 `json:"packets"`
	PacketsReply *uint8 `json:"packetsReply"`
	Pgn *uint32 `json:"pgn"`
}
func (i *IsoTransportProtocolConnectionManagementRequestToSend) PGNNumber() uint32  { return 60416 }
func DecodeIsoTransportProtocolConnectionManagementRequestToSend(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &IsoTransportProtocolConnectionManagementRequestToSend{}
	val.Info = Info
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementRequestToSend-GroupFunctionCode: %w", err)
	} else {
		if v != 16 {
			return nil, fmt.Errorf("match failed for IsoTransportProtocolConnectionManagementRequestToSend-GroupFunctionCode: Expected %d != %d", 16, v)
		}
		val.GroupFunctionCode = IsoCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementRequestToSend-MessageSize: %w", err)
	} else {
		val.MessageSize = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementRequestToSend-Packets: %w", err)
	} else {
		val.Packets = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementRequestToSend-PacketsReply: %w", err)
	} else {
		val.PacketsReply = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(24); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementRequestToSend-Pgn: %w", err)
	} else {
		val.Pgn = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeIsoTransportProtocolConnectionManagementRequestToSend(val *IsoTransportProtocolConnectionManagementRequestToSend) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.GroupFunctionCode), 8)
	w.writeUInt16(val.MessageSize, 16)
	w.writeUInt8(val.Packets, 8)
	w.writeUInt8(val.PacketsReply, 8)
	w.writeUInt32(val.Pgn, 24)
	return w.Bytes(), w.Err()
}
func encodeIsoTransportProtocolConnectionManagementRequestToSendMsg(v Message) ([]byte, error) {
	val, ok := v.(*IsoTransportProtocolConnectionManagementRequestToSend)
	if !ok {
		return nil, fmt.Errorf("expected *IsoTransportProtocolConnectionManagementRequestToSend, got %T", v)
	}
	return EncodeIsoTransportProtocolConnectionManagementRequestToSend(val)
}
type IsoTransportProtocolConnectionManagementClearToSend struct {
	Info MessageInfo `json:"info"`
	GroupFunctionCode IsoCommandConst `json:"groupFunctionCode"`
	MaxPackets *uint8 `json:"maxPackets"`
	NextSid *uint8 `json:"nextSid"`
	Pgn *uint32 `json:"pgn"`
}
func (i *IsoTransportProtocolConnectionManagementClearToSend) PGNNumber() uint32  { return 60416 }
func DecodeIsoTransportProtocolConnectionManagementClearToSend(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &IsoTransportProtocolConnectionManagementClearToSend{}
	val.Info = Info
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementClearToSend-GroupFunctionCode: %w", err)
	} else {
		if v != 17 {
			return nil, fmt.Errorf("match failed for IsoTransportProtocolConnectionManagementClearToSend-GroupFunctionCode: Expected %d != %d", 17, v)
		}
		val.GroupFunctionCode = IsoCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementClearToSend-MaxPackets: %w", err)
	} else {
		val.MaxPackets = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementClearToSend-NextSid: %w", err)
	} else {
		val.NextSid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(16)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt32(24); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementClearToSend-Pgn: %w", err)
	} else {
		val.Pgn = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeIsoTransportProtocolConnectionManagementClearToSend(val *IsoTransportProtocolConnectionManagementClearToSend) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.GroupFunctionCode), 8)
	w.writeUInt8(val.MaxPackets, 8)
	w.writeUInt8(val.NextSid, 8)
	w.skipBits(16)
	w.writeUInt32(val.Pgn, 24)
	return w.Bytes(), w.Err()
}
func encodeIsoTransportProtocolConnectionManagementClearToSendMsg(v Message) ([]byte, error) {
	val, ok := v.(*IsoTransportProtocolConnectionManagementClearToSend)
	if !ok {
		return nil, fmt.Errorf("expected *IsoTransportProtocolConnectionManagementClearToSend, got %T", v)
	}
	return EncodeIsoTransportProtocolConnectionManagementClearToSend(val)
}
type IsoTransportProtocolConnectionManagementEndOfMessage struct {
	Info MessageInfo `json:"info"`
	GroupFunctionCode IsoCommandConst `json:"groupFunctionCode"`
	TotalMessageSize *uint16 `json:"totalMessageSize"`
	TotalNumberOfFramesReceived *uint8 `json:"totalNumberOfFramesReceived"`
	Pgn *uint32 `json:"pgn"`
}
func (i *IsoTransportProtocolConnectionManagementEndOfMessage) PGNNumber() uint32  { return 60416 }
func DecodeIsoTransportProtocolConnectionManagementEndOfMessage(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &IsoTransportProtocolConnectionManagementEndOfMessage{}
	val.Info = Info
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementEndOfMessage-GroupFunctionCode: %w", err)
	} else {
		if v != 19 {
			return nil, fmt.Errorf("match failed for IsoTransportProtocolConnectionManagementEndOfMessage-GroupFunctionCode: Expected %d != %d", 19, v)
		}
		val.GroupFunctionCode = IsoCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementEndOfMessage-TotalMessageSize: %w", err)
	} else {
		val.TotalMessageSize = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementEndOfMessage-TotalNumberOfFramesReceived: %w", err)
	} else {
		val.TotalNumberOfFramesReceived = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt32(24); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementEndOfMessage-Pgn: %w", err)
	} else {
		val.Pgn = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeIsoTransportProtocolConnectionManagementEndOfMessage(val *IsoTransportProtocolConnectionManagementEndOfMessage) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.GroupFunctionCode), 8)
	w.writeUInt16(val.TotalMessageSize, 16)
	w.writeUInt8(val.TotalNumberOfFramesReceived, 8)
	w.skipBits(8)
	w.writeUInt32(val.Pgn, 24)
	return w.Bytes(), w.Err()
}
func encodeIsoTransportProtocolConnectionManagementEndOfMessageMsg(v Message) ([]byte, error) {
	val, ok := v.(*IsoTransportProtocolConnectionManagementEndOfMessage)
	if !ok {
		return nil, fmt.Errorf("expected *IsoTransportProtocolConnectionManagementEndOfMessage, got %T", v)
	}
	return EncodeIsoTransportProtocolConnectionManagementEndOfMessage(val)
}
type IsoTransportProtocolConnectionManagementBroadcastAnnounce struct {
	Info MessageInfo `json:"info"`
	GroupFunctionCode IsoCommandConst `json:"groupFunctionCode"`
	MessageSize *uint16 `json:"messageSize"`
	Packets *uint8 `json:"packets"`
	Pgn *uint32 `json:"pgn"`
}
func (i *IsoTransportProtocolConnectionManagementBroadcastAnnounce) PGNNumber() uint32  { return 60416 }
func DecodeIsoTransportProtocolConnectionManagementBroadcastAnnounce(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &IsoTransportProtocolConnectionManagementBroadcastAnnounce{}
	val.Info = Info
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementBroadcastAnnounce-GroupFunctionCode: %w", err)
	} else {
		if v != 32 {
			return nil, fmt.Errorf("match failed for IsoTransportProtocolConnectionManagementBroadcastAnnounce-GroupFunctionCode: Expected %d != %d", 32, v)
		}
		val.GroupFunctionCode = IsoCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementBroadcastAnnounce-MessageSize: %w", err)
	} else {
		val.MessageSize = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementBroadcastAnnounce-Packets: %w", err)
	} else {
		val.Packets = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(8)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt32(24); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementBroadcastAnnounce-Pgn: %w", err)
	} else {
		val.Pgn = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeIsoTransportProtocolConnectionManagementBroadcastAnnounce(val *IsoTransportProtocolConnectionManagementBroadcastAnnounce) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.GroupFunctionCode), 8)
	w.writeUInt16(val.MessageSize, 16)
	w.writeUInt8(val.Packets, 8)
	w.skipBits(8)
	w.writeUInt32(val.Pgn, 24)
	return w.Bytes(), w.Err()
}
func encodeIsoTransportProtocolConnectionManagementBroadcastAnnounceMsg(v Message) ([]byte, error) {
	val, ok := v.(*IsoTransportProtocolConnectionManagementBroadcastAnnounce)
	if !ok {
		return nil, fmt.Errorf("expected *IsoTransportProtocolConnectionManagementBroadcastAnnounce, got %T", v)
	}
	return EncodeIsoTransportProtocolConnectionManagementBroadcastAnnounce(val)
}
type IsoTransportProtocolConnectionManagementAbort struct {
	Info MessageInfo `json:"info"`
	GroupFunctionCode IsoCommandConst `json:"groupFunctionCode"`
	Reason []uint8 `json:"reason"`
	Pgn *uint32 `json:"pgn"`
}
func (i *IsoTransportProtocolConnectionManagementAbort) PGNNumber() uint32  { return 60416 }
func DecodeIsoTransportProtocolConnectionManagementAbort(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &IsoTransportProtocolConnectionManagementAbort{}
	val.Info = Info
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementAbort-GroupFunctionCode: %w", err)
	} else {
		if v != 255 {
			return nil, fmt.Errorf("match failed for IsoTransportProtocolConnectionManagementAbort-GroupFunctionCode: Expected %d != %d", 255, v)
		}
		val.GroupFunctionCode = IsoCommandConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readBinaryData(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementAbort-Reason: %w", err)
	} else {
		val.Reason = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(24)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt32(24); err != nil {
		return nil, fmt.Errorf("parse failed for IsoTransportProtocolConnectionManagementAbort-Pgn: %w", err)
	} else {
		val.Pgn = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeIsoTransportProtocolConnectionManagementAbort(val *IsoTransportProtocolConnectionManagementAbort) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.GroupFunctionCode), 8)
	w.writeBinaryData(val.Reason, 8)
	w.skipBits(24)
	w.writeUInt32(val.Pgn, 24)
	return w.Bytes(), w.Err()
}
func encodeIsoTransportProtocolConnectionManagementAbortMsg(v Message) ([]byte, error) {
	val, ok := v.(*IsoTransportProtocolConnectionManagementAbort)
	if !ok {
		return nil, fmt.Errorf("expected *IsoTransportProtocolConnectionManagementAbort, got %T", v)
	}
	return EncodeIsoTransportProtocolConnectionManagementAbort(val)
}
type IsoAddressClaim struct {
	Info MessageInfo `json:"info"`
	UniqueNumber *uint32 `json:"uniqueNumber"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	DeviceInstanceLower *uint8 `json:"deviceInstanceLower"`
	DeviceInstanceUpper *uint8 `json:"deviceInstanceUpper"`
	DeviceFunction DeviceFunctionConst `json:"deviceFunction"`
	DeviceClass DeviceClassConst `json:"deviceClass"`
	SystemInstance *uint8 `json:"systemInstance"`
	IndustryGroup IndustryCodeConst `json:"industryGroup"`
	ArbitraryAddressCapable *uint8 `json:"arbitraryAddressCapable"`
}
func (i *IsoAddressClaim) PGNNumber() uint32  { return 60928 }
func DecodeIsoAddressClaim(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &IsoAddressClaim{}
	val.Info = Info
	if v, err := stream.readUInt32(21); err != nil {
		return nil, fmt.Errorf("parse failed for IsoAddressClaim-UniqueNumber: %w", err)
	} else {
		val.UniqueNumber = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for IsoAddressClaim-ManufacturerCode: %w", err)
	} else {
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(3); err != nil {
		return nil, fmt.Errorf("parse failed for IsoAddressClaim-DeviceInstanceLower: %w", err)
	} else {
		val.DeviceInstanceLower = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(5); err != nil {
		return nil, fmt.Errorf("parse failed for IsoAddressClaim-DeviceInstanceUpper: %w", err)
	} else {
		val.DeviceInstanceUpper = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoAddressClaim-DeviceFunction: %w", err)
	} else {
		val.DeviceFunction = DeviceFunctionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(1)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(7); err != nil {
		return nil, fmt.Errorf("parse failed for IsoAddressClaim-DeviceClass: %w", err)
	} else {
		val.DeviceClass = DeviceClassConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for IsoAddressClaim-SystemInstance: %w", err)
	} else {
		val.SystemInstance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for IsoAddressClaim-IndustryGroup: %w", err)
	} else {
		val.IndustryGroup = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(1); err != nil {
		return nil, fmt.Errorf("parse failed for IsoAddressClaim-ArbitraryAddressCapable: %w", err)
	} else {
		val.ArbitraryAddressCapable = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeIsoAddressClaim(val *IsoAddressClaim) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt32(val.UniqueNumber, 21)
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeUInt8(val.DeviceInstanceLower, 3)
	w.writeUInt8(val.DeviceInstanceUpper, 5)
	w.writeLookupField(uint64(val.DeviceFunction), 8)
	w.skipBits(1)
	w.writeLookupField(uint64(val.DeviceClass), 7)
	w.writeUInt8(val.SystemInstance, 4)
	w.writeLookupField(uint64(val.IndustryGroup), 3)
	w.writeUInt8(val.ArbitraryAddressCapable, 1)
	return w.Bytes(), w.Err()
}
func encodeIsoAddressClaimMsg(v Message) ([]byte, error) {
	val, ok := v.(*IsoAddressClaim)
	if !ok {
		return nil, fmt.Errorf("expected *IsoAddressClaim, got %T", v)
	}
	return EncodeIsoAddressClaim(val)
}
type IsoCommandedAddress struct {
	Info MessageInfo `json:"info"`
	UniqueNumber []uint8 `json:"uniqueNumber"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	DeviceInstanceLower *uint8 `json:"deviceInstanceLower"`
	DeviceInstanceUpper *uint8 `json:"deviceInstanceUpper"`
	DeviceFunction DeviceFunctionConst `json:"deviceFunction"`
	DeviceClass DeviceClassConst `json:"deviceClass"`
	SystemInstance *uint8 `json:"systemInstance"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	NewSourceAddress *uint8 `json:"newSourceAddress"`
}
func (i *IsoCommandedAddress) PGNNumber() uint32  { return 65240 }
func DecodeIsoCommandedAddress(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &IsoCommandedAddress{}
	val.Info = Info
	if v, err := stream.readBinaryData(21); err != nil {
		return nil, fmt.Errorf("parse failed for IsoCommandedAddress-UniqueNumber: %w", err)
	} else {
		val.UniqueNumber = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for IsoCommandedAddress-ManufacturerCode: %w", err)
	} else {
		val.ManufacturerCode = ManufacturerCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(3); err != nil {
		return nil, fmt.Errorf("parse failed for IsoCommandedAddress-DeviceInstanceLower: %w", err)
	} else {
		val.DeviceInstanceLower = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(5); err != nil {
		return nil, fmt.Errorf("parse failed for IsoCommandedAddress-DeviceInstanceUpper: %w", err)
	} else {
		val.DeviceInstanceUpper = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoCommandedAddress-DeviceFunction: %w", err)
	} else {
		val.DeviceFunction = DeviceFunctionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(1)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readLookupField(7); err != nil {
		return nil, fmt.Errorf("parse failed for IsoCommandedAddress-DeviceClass: %w", err)
	} else {
		val.DeviceClass = DeviceClassConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(4); err != nil {
		return nil, fmt.Errorf("parse failed for IsoCommandedAddress-SystemInstance: %w", err)
	} else {
		val.SystemInstance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for IsoCommandedAddress-IndustryCode: %w", err)
	} else {
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(1)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for IsoCommandedAddress-NewSourceAddress: %w", err)
	} else {
		val.NewSourceAddress = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeIsoCommandedAddress(val *IsoCommandedAddress) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeBinaryData(val.UniqueNumber, 21)
	w.writeLookupField(uint64(val.ManufacturerCode), 11)
	w.writeUInt8(val.DeviceInstanceLower, 3)
	w.writeUInt8(val.DeviceInstanceUpper, 5)
	w.writeLookupField(uint64(val.DeviceFunction), 8)
	w.skipBits(1)
	w.writeLookupField(uint64(val.DeviceClass), 7)
	w.writeUInt8(val.SystemInstance, 4)
	w.writeLookupField(uint64(val.IndustryCode), 3)
	w.skipBits(1)
	w.writeUInt8(val.NewSourceAddress, 8)
	return w.Bytes(), w.Err()
}
func encodeIsoCommandedAddressMsg(v Message) ([]byte, error) {
	val, ok := v.(*IsoCommandedAddress)
	if !ok {
		return nil, fmt.Errorf("expected *IsoCommandedAddress, got %T", v)
	}
	return EncodeIsoCommandedAddress(val)
}
type NmeaRequestGroupFunction struct {
	Info MessageInfo `json:"info"`
	FunctionCode GroupFunctionConst `json:"functionCode"`
	Pgn *uint32 `json:"pgn"`
	TransmissionInterval *float32 `json:"transmissionInterval"`
	TransmissionIntervalOffset *float32 `json:"transmissionIntervalOffset"`
	NumberOfParameters *uint8 `json:"numberOfParameters"`
	Repeating1 []NmeaRequestGroupFunctionRepeating1 `json:"repeating1"`
}
func (n *NmeaRequestGroupFunction) PGNNumber() uint32  { return 126208 }
type NmeaRequestGroupFunctionRepeating1 struct {
	Parameter *uint8 `json:"parameter"`
	Value []uint8 `json:"value"`
}
func DecodeNmeaRequestGroupFunction(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &NmeaRequestGroupFunction{}
	val.Info = Info
		var repeat1Count uint16 = 0
		var fieldIndex uint8
		var manufacturer ManufacturerCodeConst
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaRequestGroupFunction-FunctionCode: %w", err)
	} else {
		if v != 0 {
			return nil, fmt.Errorf("match failed for NmeaRequestGroupFunction-FunctionCode: Expected %d != %d", 0, v)
		}
		val.FunctionCode = GroupFunctionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(24); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaRequestGroupFunction-Pgn: %w", err)
	} else {
		val.Pgn = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaRequestGroupFunction-TransmissionInterval: %w", err)
	} else {
		val.TransmissionInterval = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaRequestGroupFunction-TransmissionIntervalOffset: %w", err)
	} else {
		val.TransmissionIntervalOffset = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaRequestGroupFunction-NumberOfParameters: %w", err)
	} else {
		val.NumberOfParameters = v
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
	val.Repeating1 = make([]NmeaRequestGroupFunctionRepeating1, 0)
	i := 0 
	for {
		var rep NmeaRequestGroupFunctionRepeating1
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaRequestGroupFunction-Parameter: %w", err)
		} else {
			rep.Parameter = v
			if v != nil {
				fieldIndex = *v
			}
		}
		if v, err := stream.readVariableData(*val.Pgn, manufacturer, fieldIndex); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaRequestGroupFunction-Value: %w", err)
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

func EncodeNmeaRequestGroupFunction(val *NmeaRequestGroupFunction) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.FunctionCode), 8)
	w.writeUInt32(val.Pgn, 24)
	w.writeUnsignedResolution(val.TransmissionInterval, 32, 0.001)
	w.writeUnsignedResolution(val.TransmissionIntervalOffset, 16, 0.01)
	w.writeUInt8(val.NumberOfParameters, 8)
	for _, rep := range val.Repeating1 {
		w.writeUInt8(rep.Parameter, 8)
		if val.Pgn != nil {
			w.writeVariableData(rep.Value, *val.Pgn, 0, func() uint8 {
				if rep.Parameter != nil {
					return *rep.Parameter
				}
				return 0
			}())
		}
	}
	return w.Bytes(), w.Err()
}
func encodeNmeaRequestGroupFunctionMsg(v Message) ([]byte, error) {
	val, ok := v.(*NmeaRequestGroupFunction)
	if !ok {
		return nil, fmt.Errorf("expected *NmeaRequestGroupFunction, got %T", v)
	}
	return EncodeNmeaRequestGroupFunction(val)
}
type NmeaCommandGroupFunction struct {
	Info MessageInfo `json:"info"`
	FunctionCode GroupFunctionConst `json:"functionCode"`
	Pgn *uint32 `json:"pgn"`
	Priority PriorityConst `json:"priority"`
	NumberOfParameters *uint8 `json:"numberOfParameters"`
	Repeating1 []NmeaCommandGroupFunctionRepeating1 `json:"repeating1"`
}
func (n *NmeaCommandGroupFunction) PGNNumber() uint32  { return 126208 }
type NmeaCommandGroupFunctionRepeating1 struct {
	Parameter *uint8 `json:"parameter"`
	Value []uint8 `json:"value"`
}
func DecodeNmeaCommandGroupFunction(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &NmeaCommandGroupFunction{}
	val.Info = Info
		var repeat1Count uint16 = 0
		var fieldIndex uint8
		var manufacturer ManufacturerCodeConst
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaCommandGroupFunction-FunctionCode: %w", err)
	} else {
		if v != 1 {
			return nil, fmt.Errorf("match failed for NmeaCommandGroupFunction-FunctionCode: Expected %d != %d", 1, v)
		}
		val.FunctionCode = GroupFunctionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(24); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaCommandGroupFunction-Pgn: %w", err)
	} else {
		val.Pgn = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaCommandGroupFunction-Priority: %w", err)
	} else {
		val.Priority = PriorityConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaCommandGroupFunction-NumberOfParameters: %w", err)
	} else {
		val.NumberOfParameters = v
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
	val.Repeating1 = make([]NmeaCommandGroupFunctionRepeating1, 0)
	i := 0 
	for {
		var rep NmeaCommandGroupFunctionRepeating1
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaCommandGroupFunction-Parameter: %w", err)
		} else {
			rep.Parameter = v
			if v != nil {
				fieldIndex = *v
			}
		}
		if v, err := stream.readVariableData(*val.Pgn, manufacturer, fieldIndex); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaCommandGroupFunction-Value: %w", err)
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

func EncodeNmeaCommandGroupFunction(val *NmeaCommandGroupFunction) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.FunctionCode), 8)
	w.writeUInt32(val.Pgn, 24)
	w.writeLookupField(uint64(val.Priority), 4)
	w.skipBits(4)
	w.writeUInt8(val.NumberOfParameters, 8)
	for _, rep := range val.Repeating1 {
		w.writeUInt8(rep.Parameter, 8)
		if val.Pgn != nil {
			w.writeVariableData(rep.Value, *val.Pgn, 0, func() uint8 {
				if rep.Parameter != nil {
					return *rep.Parameter
				}
				return 0
			}())
		}
	}
	return w.Bytes(), w.Err()
}
func encodeNmeaCommandGroupFunctionMsg(v Message) ([]byte, error) {
	val, ok := v.(*NmeaCommandGroupFunction)
	if !ok {
		return nil, fmt.Errorf("expected *NmeaCommandGroupFunction, got %T", v)
	}
	return EncodeNmeaCommandGroupFunction(val)
}
type NmeaAcknowledgeGroupFunction struct {
	Info MessageInfo `json:"info"`
	FunctionCode GroupFunctionConst `json:"functionCode"`
	Pgn *uint32 `json:"pgn"`
	PgnErrorCode PgnErrorCodeConst `json:"pgnErrorCode"`
	TransmissionIntervalPriorityErrorCode TransmissionIntervalConst `json:"transmissionIntervalPriorityErrorCode"`
	NumberOfParameters *uint8 `json:"numberOfParameters"`
	Repeating1 []NmeaAcknowledgeGroupFunctionRepeating1 `json:"repeating1"`
}
func (n *NmeaAcknowledgeGroupFunction) PGNNumber() uint32  { return 126208 }
type NmeaAcknowledgeGroupFunctionRepeating1 struct {
	Parameter ParameterFieldConst `json:"parameter"`
}
func DecodeNmeaAcknowledgeGroupFunction(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &NmeaAcknowledgeGroupFunction{}
	val.Info = Info
		var repeat1Count uint16 = 0
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaAcknowledgeGroupFunction-FunctionCode: %w", err)
	} else {
		if v != 2 {
			return nil, fmt.Errorf("match failed for NmeaAcknowledgeGroupFunction-FunctionCode: Expected %d != %d", 2, v)
		}
		val.FunctionCode = GroupFunctionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(24); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaAcknowledgeGroupFunction-Pgn: %w", err)
	} else {
		val.Pgn = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaAcknowledgeGroupFunction-PgnErrorCode: %w", err)
	} else {
		val.PgnErrorCode = PgnErrorCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaAcknowledgeGroupFunction-TransmissionIntervalPriorityErrorCode: %w", err)
	} else {
		val.TransmissionIntervalPriorityErrorCode = TransmissionIntervalConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaAcknowledgeGroupFunction-NumberOfParameters: %w", err)
	} else {
		val.NumberOfParameters = v
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
	val.Repeating1 = make([]NmeaAcknowledgeGroupFunctionRepeating1, 0)
	i := 0 
	for {
		var rep NmeaAcknowledgeGroupFunctionRepeating1
		if v, err := stream.readLookupField(4); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaAcknowledgeGroupFunction-Parameter: %w", err)
		} else {
			rep.Parameter = ParameterFieldConst(v)
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

func EncodeNmeaAcknowledgeGroupFunction(val *NmeaAcknowledgeGroupFunction) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.FunctionCode), 8)
	w.writeUInt32(val.Pgn, 24)
	w.writeLookupField(uint64(val.PgnErrorCode), 4)
	w.writeLookupField(uint64(val.TransmissionIntervalPriorityErrorCode), 4)
	w.writeUInt8(val.NumberOfParameters, 8)
	for _, rep := range val.Repeating1 {
		w.writeLookupField(uint64(rep.Parameter), 4)
	}
	return w.Bytes(), w.Err()
}
func encodeNmeaAcknowledgeGroupFunctionMsg(v Message) ([]byte, error) {
	val, ok := v.(*NmeaAcknowledgeGroupFunction)
	if !ok {
		return nil, fmt.Errorf("expected *NmeaAcknowledgeGroupFunction, got %T", v)
	}
	return EncodeNmeaAcknowledgeGroupFunction(val)
}
type NmeaReadFieldsGroupFunction struct {
	Info MessageInfo `json:"info"`
	FunctionCode GroupFunctionConst `json:"functionCode"`
	Pgn *uint32 `json:"pgn"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	UniqueId *uint8 `json:"uniqueId"`
	NumberOfSelectionPairs *uint8 `json:"numberOfSelectionPairs"`
	NumberOfParameters *uint8 `json:"numberOfParameters"`
	Repeating1 []NmeaReadFieldsGroupFunctionRepeating1 `json:"repeating1"`
	Repeating2 []NmeaReadFieldsGroupFunctionRepeating2 `json:"repeating2"`
}
func (n *NmeaReadFieldsGroupFunction) PGNNumber() uint32  { return 126208 }
type NmeaReadFieldsGroupFunctionRepeating1 struct {
	SelectionParameter *uint8 `json:"selectionParameter"`
	SelectionValue []uint8 `json:"selectionValue"`
}
type NmeaReadFieldsGroupFunctionRepeating2 struct {
	Parameter *uint8 `json:"parameter"`
}
func DecodeNmeaReadFieldsGroupFunction(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &NmeaReadFieldsGroupFunction{}
	val.Info = Info
		var repeat1Count uint16 = 0
		var repeat2Count uint16
		var fieldIndex uint8
		var manufacturer ManufacturerCodeConst
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaReadFieldsGroupFunction-FunctionCode: %w", err)
	} else {
		if v != 3 {
			return nil, fmt.Errorf("match failed for NmeaReadFieldsGroupFunction-FunctionCode: Expected %d != %d", 3, v)
		}
		val.FunctionCode = GroupFunctionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(24); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaReadFieldsGroupFunction-Pgn: %w", err)
	} else {
		val.Pgn = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if IsProprietaryPGN( *val.Pgn) {
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaReadFieldsGroupFunction-ManufacturerCode: %w", err)
	} else {
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
		return nil, fmt.Errorf("parse failed for NmeaReadFieldsGroupFunction-IndustryCode: %w", err)
	} else {
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaReadFieldsGroupFunction-UniqueId: %w", err)
	} else {
		val.UniqueId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaReadFieldsGroupFunction-NumberOfSelectionPairs: %w", err)
	} else {
		val.NumberOfSelectionPairs = v
		if v != nil {
			repeat1Count = uint16(*v)
		}

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaReadFieldsGroupFunction-NumberOfParameters: %w", err)
	} else {
		val.NumberOfParameters = v
		if v != nil {
			repeat2Count = uint16(*v)
		}

		if stream.isEOF() {
			return val, nil
		} 
	}
		if repeat1Count == 0 {
			return val, nil
		}
	val.Repeating1 = make([]NmeaReadFieldsGroupFunctionRepeating1, 0)
	i := 0 
	for {
		var rep NmeaReadFieldsGroupFunctionRepeating1
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaReadFieldsGroupFunction-SelectionParameter: %w", err)
		} else {
			rep.SelectionParameter = v
			if v != nil {
				fieldIndex = *v
			}
		}
		if v, err := stream.readVariableData(*val.Pgn, manufacturer, fieldIndex); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaReadFieldsGroupFunction-SelectionValue: %w", err)
		} else {
			rep.SelectionValue = v
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
		if repeat2Count == 0 {
			return val, nil
		}	
	val.Repeating2 = make([]NmeaReadFieldsGroupFunctionRepeating2, 0)
	i = 0
	for {
		var rep NmeaReadFieldsGroupFunctionRepeating2
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaReadFieldsGroupFunction-Parameter: %w", err)
		} else {
			rep.Parameter = v
		}
		val.Repeating2 = append(val.Repeating2, rep)
		if int(repeat2Count) == 0 {
			if stream.isEOF() {
				return val, nil
			}  
		} else {
			i++
			if i == int(repeat2Count) {
				break
			} 
		} 
	}	
	return val, nil
}

func EncodeNmeaReadFieldsGroupFunction(val *NmeaReadFieldsGroupFunction) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.FunctionCode), 8)
	w.writeUInt32(val.Pgn, 24)
	if val.Pgn != nil && IsProprietaryPGN(*val.Pgn) {
		w.writeLookupField(uint64(val.ManufacturerCode), 11)
		w.skipBits(2)
		w.writeLookupField(uint64(val.IndustryCode), 3)
	}
	w.writeUInt8(val.UniqueId, 8)
	w.writeUInt8(val.NumberOfSelectionPairs, 8)
	w.writeUInt8(val.NumberOfParameters, 8)
	for _, rep := range val.Repeating1 {
		w.writeUInt8(rep.SelectionParameter, 8)
		if val.Pgn != nil {
			w.writeVariableData(rep.SelectionValue, *val.Pgn, ManufacturerCodeConst(val.ManufacturerCode), func() uint8 {
				if rep.SelectionParameter != nil {
					return *rep.SelectionParameter
				}
				return 0
			}())
		}
	}
	for _, rep := range val.Repeating2 {
		w.writeUInt8(rep.Parameter, 8)
	}
	return w.Bytes(), w.Err()
}
func encodeNmeaReadFieldsGroupFunctionMsg(v Message) ([]byte, error) {
	val, ok := v.(*NmeaReadFieldsGroupFunction)
	if !ok {
		return nil, fmt.Errorf("expected *NmeaReadFieldsGroupFunction, got %T", v)
	}
	return EncodeNmeaReadFieldsGroupFunction(val)
}
type NmeaReadFieldsReplyGroupFunction struct {
	Info MessageInfo `json:"info"`
	FunctionCode GroupFunctionConst `json:"functionCode"`
	Pgn *uint32 `json:"pgn"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	UniqueId *uint8 `json:"uniqueId"`
	NumberOfSelectionPairs *uint8 `json:"numberOfSelectionPairs"`
	NumberOfParameters *uint8 `json:"numberOfParameters"`
	Repeating1 []NmeaReadFieldsReplyGroupFunctionRepeating1 `json:"repeating1"`
	Repeating2 []NmeaReadFieldsReplyGroupFunctionRepeating2 `json:"repeating2"`
}
func (n *NmeaReadFieldsReplyGroupFunction) PGNNumber() uint32  { return 126208 }
type NmeaReadFieldsReplyGroupFunctionRepeating1 struct {
	SelectionParameter *uint8 `json:"selectionParameter"`
	SelectionValue []uint8 `json:"selectionValue"`
}
type NmeaReadFieldsReplyGroupFunctionRepeating2 struct {
	Parameter *uint8 `json:"parameter"`
	Value []uint8 `json:"value"`
}
func DecodeNmeaReadFieldsReplyGroupFunction(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &NmeaReadFieldsReplyGroupFunction{}
	val.Info = Info
		var repeat1Count uint16 = 0
		var repeat2Count uint16
		var fieldIndex uint8
		var manufacturer ManufacturerCodeConst
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaReadFieldsReplyGroupFunction-FunctionCode: %w", err)
	} else {
		if v != 4 {
			return nil, fmt.Errorf("match failed for NmeaReadFieldsReplyGroupFunction-FunctionCode: Expected %d != %d", 4, v)
		}
		val.FunctionCode = GroupFunctionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(24); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaReadFieldsReplyGroupFunction-Pgn: %w", err)
	} else {
		val.Pgn = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if IsProprietaryPGN( *val.Pgn) {
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaReadFieldsReplyGroupFunction-ManufacturerCode: %w", err)
	} else {
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
		return nil, fmt.Errorf("parse failed for NmeaReadFieldsReplyGroupFunction-IndustryCode: %w", err)
	} else {
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaReadFieldsReplyGroupFunction-UniqueId: %w", err)
	} else {
		val.UniqueId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaReadFieldsReplyGroupFunction-NumberOfSelectionPairs: %w", err)
	} else {
		val.NumberOfSelectionPairs = v
		if v != nil {
			repeat1Count = uint16(*v)
		}

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaReadFieldsReplyGroupFunction-NumberOfParameters: %w", err)
	} else {
		val.NumberOfParameters = v
		if v != nil {
			repeat2Count = uint16(*v)
		}

		if stream.isEOF() {
			return val, nil
		} 
	}
		if repeat1Count == 0 {
			return val, nil
		}
	val.Repeating1 = make([]NmeaReadFieldsReplyGroupFunctionRepeating1, 0)
	i := 0 
	for {
		var rep NmeaReadFieldsReplyGroupFunctionRepeating1
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaReadFieldsReplyGroupFunction-SelectionParameter: %w", err)
		} else {
			rep.SelectionParameter = v
			if v != nil {
				fieldIndex = *v
			}
		}
		if v, err := stream.readVariableData(*val.Pgn, manufacturer, fieldIndex); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaReadFieldsReplyGroupFunction-SelectionValue: %w", err)
		} else {
			rep.SelectionValue = v
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
		if repeat2Count == 0 {
			return val, nil
		}	
	val.Repeating2 = make([]NmeaReadFieldsReplyGroupFunctionRepeating2, 0)
	i = 0
	for {
		var rep NmeaReadFieldsReplyGroupFunctionRepeating2
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaReadFieldsReplyGroupFunction-Parameter: %w", err)
		} else {
			rep.Parameter = v
		}
		if v, err := stream.readVariableData(*val.Pgn, manufacturer, fieldIndex); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaReadFieldsReplyGroupFunction-Value: %w", err)
		} else {
			rep.Value = v
		}
		val.Repeating2 = append(val.Repeating2, rep)
		if int(repeat2Count) == 0 {
			if stream.isEOF() {
				return val, nil
			}  
		} else {
			i++
			if i == int(repeat2Count) {
				break
			} 
		} 
	}	
	return val, nil
}

func EncodeNmeaReadFieldsReplyGroupFunction(val *NmeaReadFieldsReplyGroupFunction) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.FunctionCode), 8)
	w.writeUInt32(val.Pgn, 24)
	if val.Pgn != nil && IsProprietaryPGN(*val.Pgn) {
		w.writeLookupField(uint64(val.ManufacturerCode), 11)
		w.skipBits(2)
		w.writeLookupField(uint64(val.IndustryCode), 3)
	}
	w.writeUInt8(val.UniqueId, 8)
	w.writeUInt8(val.NumberOfSelectionPairs, 8)
	w.writeUInt8(val.NumberOfParameters, 8)
	for _, rep := range val.Repeating1 {
		w.writeUInt8(rep.SelectionParameter, 8)
		if val.Pgn != nil {
			w.writeVariableData(rep.SelectionValue, *val.Pgn, ManufacturerCodeConst(val.ManufacturerCode), func() uint8 {
				if rep.SelectionParameter != nil {
					return *rep.SelectionParameter
				}
				return 0
			}())
		}
	}
	for _, rep := range val.Repeating2 {
		w.writeUInt8(rep.Parameter, 8)
		if val.Pgn != nil {
			w.writeVariableData(rep.Value, *val.Pgn, ManufacturerCodeConst(val.ManufacturerCode), func() uint8 {
				if rep.Parameter != nil {
					return *rep.Parameter
				}
				return 0
			}())
		}
	}
	return w.Bytes(), w.Err()
}
func encodeNmeaReadFieldsReplyGroupFunctionMsg(v Message) ([]byte, error) {
	val, ok := v.(*NmeaReadFieldsReplyGroupFunction)
	if !ok {
		return nil, fmt.Errorf("expected *NmeaReadFieldsReplyGroupFunction, got %T", v)
	}
	return EncodeNmeaReadFieldsReplyGroupFunction(val)
}
type NmeaWriteFieldsGroupFunction struct {
	Info MessageInfo `json:"info"`
	FunctionCode GroupFunctionConst `json:"functionCode"`
	Pgn *uint32 `json:"pgn"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	UniqueId *uint8 `json:"uniqueId"`
	NumberOfSelectionPairs *uint8 `json:"numberOfSelectionPairs"`
	NumberOfParameters *uint8 `json:"numberOfParameters"`
	Repeating1 []NmeaWriteFieldsGroupFunctionRepeating1 `json:"repeating1"`
	Repeating2 []NmeaWriteFieldsGroupFunctionRepeating2 `json:"repeating2"`
}
func (n *NmeaWriteFieldsGroupFunction) PGNNumber() uint32  { return 126208 }
type NmeaWriteFieldsGroupFunctionRepeating1 struct {
	SelectionParameter *uint8 `json:"selectionParameter"`
	SelectionValue []uint8 `json:"selectionValue"`
}
type NmeaWriteFieldsGroupFunctionRepeating2 struct {
	Parameter *uint8 `json:"parameter"`
	Value []uint8 `json:"value"`
}
func DecodeNmeaWriteFieldsGroupFunction(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &NmeaWriteFieldsGroupFunction{}
	val.Info = Info
		var repeat1Count uint16 = 0
		var repeat2Count uint16
		var fieldIndex uint8
		var manufacturer ManufacturerCodeConst
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaWriteFieldsGroupFunction-FunctionCode: %w", err)
	} else {
		if v != 5 {
			return nil, fmt.Errorf("match failed for NmeaWriteFieldsGroupFunction-FunctionCode: Expected %d != %d", 5, v)
		}
		val.FunctionCode = GroupFunctionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(24); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaWriteFieldsGroupFunction-Pgn: %w", err)
	} else {
		val.Pgn = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if IsProprietaryPGN( *val.Pgn) {
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaWriteFieldsGroupFunction-ManufacturerCode: %w", err)
	} else {
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
		return nil, fmt.Errorf("parse failed for NmeaWriteFieldsGroupFunction-IndustryCode: %w", err)
	} else {
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaWriteFieldsGroupFunction-UniqueId: %w", err)
	} else {
		val.UniqueId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaWriteFieldsGroupFunction-NumberOfSelectionPairs: %w", err)
	} else {
		val.NumberOfSelectionPairs = v
		if v != nil {
			repeat1Count = uint16(*v)
		}

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaWriteFieldsGroupFunction-NumberOfParameters: %w", err)
	} else {
		val.NumberOfParameters = v
		if v != nil {
			repeat2Count = uint16(*v)
		}

		if stream.isEOF() {
			return val, nil
		} 
	}
		if repeat1Count == 0 {
			return val, nil
		}
	val.Repeating1 = make([]NmeaWriteFieldsGroupFunctionRepeating1, 0)
	i := 0 
	for {
		var rep NmeaWriteFieldsGroupFunctionRepeating1
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaWriteFieldsGroupFunction-SelectionParameter: %w", err)
		} else {
			rep.SelectionParameter = v
			if v != nil {
				fieldIndex = *v
			}
		}
		if v, err := stream.readVariableData(*val.Pgn, manufacturer, fieldIndex); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaWriteFieldsGroupFunction-SelectionValue: %w", err)
		} else {
			rep.SelectionValue = v
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
		if repeat2Count == 0 {
			return val, nil
		}	
	val.Repeating2 = make([]NmeaWriteFieldsGroupFunctionRepeating2, 0)
	i = 0
	for {
		var rep NmeaWriteFieldsGroupFunctionRepeating2
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaWriteFieldsGroupFunction-Parameter: %w", err)
		} else {
			rep.Parameter = v
		}
		if v, err := stream.readVariableData(*val.Pgn, manufacturer, fieldIndex); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaWriteFieldsGroupFunction-Value: %w", err)
		} else {
			rep.Value = v
		}
		val.Repeating2 = append(val.Repeating2, rep)
		if int(repeat2Count) == 0 {
			if stream.isEOF() {
				return val, nil
			}  
		} else {
			i++
			if i == int(repeat2Count) {
				break
			} 
		} 
	}	
	return val, nil
}

func EncodeNmeaWriteFieldsGroupFunction(val *NmeaWriteFieldsGroupFunction) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.FunctionCode), 8)
	w.writeUInt32(val.Pgn, 24)
	if val.Pgn != nil && IsProprietaryPGN(*val.Pgn) {
		w.writeLookupField(uint64(val.ManufacturerCode), 11)
		w.skipBits(2)
		w.writeLookupField(uint64(val.IndustryCode), 3)
	}
	w.writeUInt8(val.UniqueId, 8)
	w.writeUInt8(val.NumberOfSelectionPairs, 8)
	w.writeUInt8(val.NumberOfParameters, 8)
	for _, rep := range val.Repeating1 {
		w.writeUInt8(rep.SelectionParameter, 8)
		if val.Pgn != nil {
			w.writeVariableData(rep.SelectionValue, *val.Pgn, ManufacturerCodeConst(val.ManufacturerCode), func() uint8 {
				if rep.SelectionParameter != nil {
					return *rep.SelectionParameter
				}
				return 0
			}())
		}
	}
	for _, rep := range val.Repeating2 {
		w.writeUInt8(rep.Parameter, 8)
		if val.Pgn != nil {
			w.writeVariableData(rep.Value, *val.Pgn, ManufacturerCodeConst(val.ManufacturerCode), func() uint8 {
				if rep.Parameter != nil {
					return *rep.Parameter
				}
				return 0
			}())
		}
	}
	return w.Bytes(), w.Err()
}
func encodeNmeaWriteFieldsGroupFunctionMsg(v Message) ([]byte, error) {
	val, ok := v.(*NmeaWriteFieldsGroupFunction)
	if !ok {
		return nil, fmt.Errorf("expected *NmeaWriteFieldsGroupFunction, got %T", v)
	}
	return EncodeNmeaWriteFieldsGroupFunction(val)
}
type NmeaWriteFieldsReplyGroupFunction struct {
	Info MessageInfo `json:"info"`
	FunctionCode GroupFunctionConst `json:"functionCode"`
	Pgn *uint32 `json:"pgn"`
	ManufacturerCode ManufacturerCodeConst `json:"manufacturerCode"`
	IndustryCode IndustryCodeConst `json:"industryCode"`
	UniqueId *uint8 `json:"uniqueId"`
	NumberOfSelectionPairs *uint8 `json:"numberOfSelectionPairs"`
	NumberOfParameters *uint8 `json:"numberOfParameters"`
	Repeating1 []NmeaWriteFieldsReplyGroupFunctionRepeating1 `json:"repeating1"`
	Repeating2 []NmeaWriteFieldsReplyGroupFunctionRepeating2 `json:"repeating2"`
}
func (n *NmeaWriteFieldsReplyGroupFunction) PGNNumber() uint32  { return 126208 }
type NmeaWriteFieldsReplyGroupFunctionRepeating1 struct {
	SelectionParameter *uint8 `json:"selectionParameter"`
	SelectionValue []uint8 `json:"selectionValue"`
}
type NmeaWriteFieldsReplyGroupFunctionRepeating2 struct {
	Parameter *uint8 `json:"parameter"`
	Value []uint8 `json:"value"`
}
func DecodeNmeaWriteFieldsReplyGroupFunction(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &NmeaWriteFieldsReplyGroupFunction{}
	val.Info = Info
		var repeat1Count uint16 = 0
		var repeat2Count uint16
		var fieldIndex uint8
		var manufacturer ManufacturerCodeConst
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaWriteFieldsReplyGroupFunction-FunctionCode: %w", err)
	} else {
		if v != 6 {
			return nil, fmt.Errorf("match failed for NmeaWriteFieldsReplyGroupFunction-FunctionCode: Expected %d != %d", 6, v)
		}
		val.FunctionCode = GroupFunctionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(24); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaWriteFieldsReplyGroupFunction-Pgn: %w", err)
	} else {
		val.Pgn = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if IsProprietaryPGN( *val.Pgn) {
	if v, err := stream.readLookupField(11); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaWriteFieldsReplyGroupFunction-ManufacturerCode: %w", err)
	} else {
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
		return nil, fmt.Errorf("parse failed for NmeaWriteFieldsReplyGroupFunction-IndustryCode: %w", err)
	} else {
		val.IndustryCode = IndustryCodeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaWriteFieldsReplyGroupFunction-UniqueId: %w", err)
	} else {
		val.UniqueId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaWriteFieldsReplyGroupFunction-NumberOfSelectionPairs: %w", err)
	} else {
		val.NumberOfSelectionPairs = v
		if v != nil {
			repeat1Count = uint16(*v)
		}

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for NmeaWriteFieldsReplyGroupFunction-NumberOfParameters: %w", err)
	} else {
		val.NumberOfParameters = v
		if v != nil {
			repeat2Count = uint16(*v)
		}

		if stream.isEOF() {
			return val, nil
		} 
	}
		if repeat1Count == 0 {
			return val, nil
		}
	val.Repeating1 = make([]NmeaWriteFieldsReplyGroupFunctionRepeating1, 0)
	i := 0 
	for {
		var rep NmeaWriteFieldsReplyGroupFunctionRepeating1
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaWriteFieldsReplyGroupFunction-SelectionParameter: %w", err)
		} else {
			rep.SelectionParameter = v
			if v != nil {
				fieldIndex = *v
			}
		}
		if v, err := stream.readVariableData(*val.Pgn, manufacturer, fieldIndex); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaWriteFieldsReplyGroupFunction-SelectionValue: %w", err)
		} else {
			rep.SelectionValue = v
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
		if repeat2Count == 0 {
			return val, nil
		}	
	val.Repeating2 = make([]NmeaWriteFieldsReplyGroupFunctionRepeating2, 0)
	i = 0
	for {
		var rep NmeaWriteFieldsReplyGroupFunctionRepeating2
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaWriteFieldsReplyGroupFunction-Parameter: %w", err)
		} else {
			rep.Parameter = v
		}
		if v, err := stream.readVariableData(*val.Pgn, manufacturer, fieldIndex); err != nil {
			return nil, fmt.Errorf("parse failed for NmeaWriteFieldsReplyGroupFunction-Value: %w", err)
		} else {
			rep.Value = v
		}
		val.Repeating2 = append(val.Repeating2, rep)
		if int(repeat2Count) == 0 {
			if stream.isEOF() {
				return val, nil
			}  
		} else {
			i++
			if i == int(repeat2Count) {
				break
			} 
		} 
	}	
	return val, nil
}

func EncodeNmeaWriteFieldsReplyGroupFunction(val *NmeaWriteFieldsReplyGroupFunction) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.FunctionCode), 8)
	w.writeUInt32(val.Pgn, 24)
	if val.Pgn != nil && IsProprietaryPGN(*val.Pgn) {
		w.writeLookupField(uint64(val.ManufacturerCode), 11)
		w.skipBits(2)
		w.writeLookupField(uint64(val.IndustryCode), 3)
	}
	w.writeUInt8(val.UniqueId, 8)
	w.writeUInt8(val.NumberOfSelectionPairs, 8)
	w.writeUInt8(val.NumberOfParameters, 8)
	for _, rep := range val.Repeating1 {
		w.writeUInt8(rep.SelectionParameter, 8)
		if val.Pgn != nil {
			w.writeVariableData(rep.SelectionValue, *val.Pgn, ManufacturerCodeConst(val.ManufacturerCode), func() uint8 {
				if rep.SelectionParameter != nil {
					return *rep.SelectionParameter
				}
				return 0
			}())
		}
	}
	for _, rep := range val.Repeating2 {
		w.writeUInt8(rep.Parameter, 8)
		if val.Pgn != nil {
			w.writeVariableData(rep.Value, *val.Pgn, ManufacturerCodeConst(val.ManufacturerCode), func() uint8 {
				if rep.Parameter != nil {
					return *rep.Parameter
				}
				return 0
			}())
		}
	}
	return w.Bytes(), w.Err()
}
func encodeNmeaWriteFieldsReplyGroupFunctionMsg(v Message) ([]byte, error) {
	val, ok := v.(*NmeaWriteFieldsReplyGroupFunction)
	if !ok {
		return nil, fmt.Errorf("expected *NmeaWriteFieldsReplyGroupFunction, got %T", v)
	}
	return EncodeNmeaWriteFieldsReplyGroupFunction(val)
}
type PgnListTransmitAndReceive struct {
	Info MessageInfo `json:"info"`
	FunctionCode PgnListFunctionConst `json:"functionCode"`
	Repeating1 []PgnListTransmitAndReceiveRepeating1 `json:"repeating1"`
}
func (p *PgnListTransmitAndReceive) PGNNumber() uint32  { return 126464 }
type PgnListTransmitAndReceiveRepeating1 struct {
	Pgn *uint32 `json:"pgn"`
}
func DecodePgnListTransmitAndReceive(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &PgnListTransmitAndReceive{}
	val.Info = Info
		var repeat1Count uint16 = 0
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for PgnListTransmitAndReceive-FunctionCode: %w", err)
	} else {
		val.FunctionCode = PgnListFunctionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	val.Repeating1 = make([]PgnListTransmitAndReceiveRepeating1, 0)
	i := 0 
	for {
		var rep PgnListTransmitAndReceiveRepeating1
		if v, err := stream.readUInt32(24); err != nil {
			return nil, fmt.Errorf("parse failed for PgnListTransmitAndReceive-Pgn: %w", err)
		} else {
			rep.Pgn = v
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

func EncodePgnListTransmitAndReceive(val *PgnListTransmitAndReceive) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.FunctionCode), 8)
	for _, rep := range val.Repeating1 {
		w.writeUInt32(rep.Pgn, 24)
	}
	return w.Bytes(), w.Err()
}
func encodePgnListTransmitAndReceiveMsg(v Message) ([]byte, error) {
	val, ok := v.(*PgnListTransmitAndReceive)
	if !ok {
		return nil, fmt.Errorf("expected *PgnListTransmitAndReceive, got %T", v)
	}
	return EncodePgnListTransmitAndReceive(val)
}
type Alert struct {
	Info MessageInfo `json:"info"`
	AlertType AlertTypeConst `json:"alertType"`
	AlertCategory AlertCategoryConst `json:"alertCategory"`
	AlertSystem *uint8 `json:"alertSystem"`
	AlertSubSystem *uint8 `json:"alertSubSystem"`
	AlertId *uint16 `json:"alertId"`
	DataSourceNetworkIdName *uint64 `json:"dataSourceNetworkIdName"`
	DataSourceInstance *uint8 `json:"dataSourceInstance"`
	DataSourceIndexSource *uint8 `json:"dataSourceIndexSource"`
	AlertOccurrenceNumber *uint8 `json:"alertOccurrenceNumber"`
	TemporarySilenceStatus YesNoConst `json:"temporarySilenceStatus"`
	AcknowledgeStatus YesNoConst `json:"acknowledgeStatus"`
	EscalationStatus YesNoConst `json:"escalationStatus"`
	TemporarySilenceSupport YesNoConst `json:"temporarySilenceSupport"`
	AcknowledgeSupport YesNoConst `json:"acknowledgeSupport"`
	EscalationSupport YesNoConst `json:"escalationSupport"`
	AcknowledgeSourceNetworkIdName *uint64 `json:"acknowledgeSourceNetworkIdName"`
	TriggerCondition AlertTriggerConditionConst `json:"triggerCondition"`
	ThresholdStatus AlertThresholdStatusConst `json:"thresholdStatus"`
	AlertPriority *uint8 `json:"alertPriority"`
	AlertState AlertStateConst `json:"alertState"`
}
func (a *Alert) PGNNumber() uint32  { return 126983 }
func DecodeAlert(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Alert{}
	val.Info = Info
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-AlertType: %w", err)
	} else {
		val.AlertType = AlertTypeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-AlertCategory: %w", err)
	} else {
		val.AlertCategory = AlertCategoryConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-AlertSystem: %w", err)
	} else {
		val.AlertSystem = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-AlertSubSystem: %w", err)
	} else {
		val.AlertSubSystem = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-AlertId: %w", err)
	} else {
		val.AlertId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt64(64); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-DataSourceNetworkIdName: %w", err)
	} else {
		val.DataSourceNetworkIdName = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-DataSourceInstance: %w", err)
	} else {
		val.DataSourceInstance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-DataSourceIndexSource: %w", err)
	} else {
		val.DataSourceIndexSource = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-AlertOccurrenceNumber: %w", err)
	} else {
		val.AlertOccurrenceNumber = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-TemporarySilenceStatus: %w", err)
	} else {
		val.TemporarySilenceStatus = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-AcknowledgeStatus: %w", err)
	} else {
		val.AcknowledgeStatus = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-EscalationStatus: %w", err)
	} else {
		val.EscalationStatus = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-TemporarySilenceSupport: %w", err)
	} else {
		val.TemporarySilenceSupport = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-AcknowledgeSupport: %w", err)
	} else {
		val.AcknowledgeSupport = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(1); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-EscalationSupport: %w", err)
	} else {
		val.EscalationSupport = YesNoConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(2)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt64(64); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-AcknowledgeSourceNetworkIdName: %w", err)
	} else {
		val.AcknowledgeSourceNetworkIdName = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-TriggerCondition: %w", err)
	} else {
		val.TriggerCondition = AlertTriggerConditionConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-ThresholdStatus: %w", err)
	} else {
		val.ThresholdStatus = AlertThresholdStatusConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-AlertPriority: %w", err)
	} else {
		val.AlertPriority = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for Alert-AlertState: %w", err)
	} else {
		val.AlertState = AlertStateConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeAlert(val *Alert) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.AlertType), 4)
	w.writeLookupField(uint64(val.AlertCategory), 4)
	w.writeUInt8(val.AlertSystem, 8)
	w.writeUInt8(val.AlertSubSystem, 8)
	w.writeUInt16(val.AlertId, 16)
	w.writeUInt64(val.DataSourceNetworkIdName, 64)
	w.writeUInt8(val.DataSourceInstance, 8)
	w.writeUInt8(val.DataSourceIndexSource, 8)
	w.writeUInt8(val.AlertOccurrenceNumber, 8)
	w.writeLookupField(uint64(val.TemporarySilenceStatus), 1)
	w.writeLookupField(uint64(val.AcknowledgeStatus), 1)
	w.writeLookupField(uint64(val.EscalationStatus), 1)
	w.writeLookupField(uint64(val.TemporarySilenceSupport), 1)
	w.writeLookupField(uint64(val.AcknowledgeSupport), 1)
	w.writeLookupField(uint64(val.EscalationSupport), 1)
	w.skipBits(2)
	w.writeUInt64(val.AcknowledgeSourceNetworkIdName, 64)
	w.writeLookupField(uint64(val.TriggerCondition), 4)
	w.writeLookupField(uint64(val.ThresholdStatus), 4)
	w.writeUInt8(val.AlertPriority, 8)
	w.writeLookupField(uint64(val.AlertState), 8)
	return w.Bytes(), w.Err()
}
func encodeAlertMsg(v Message) ([]byte, error) {
	val, ok := v.(*Alert)
	if !ok {
		return nil, fmt.Errorf("expected *Alert, got %T", v)
	}
	return EncodeAlert(val)
}
type AlertResponse struct {
	Info MessageInfo `json:"info"`
	AlertType AlertTypeConst `json:"alertType"`
	AlertCategory AlertCategoryConst `json:"alertCategory"`
	AlertSystem *uint8 `json:"alertSystem"`
	AlertSubSystem *uint8 `json:"alertSubSystem"`
	AlertId *uint16 `json:"alertId"`
	DataSourceNetworkIdName *uint64 `json:"dataSourceNetworkIdName"`
	DataSourceInstance *uint8 `json:"dataSourceInstance"`
	DataSourceIndexSource *uint8 `json:"dataSourceIndexSource"`
	AlertOccurrenceNumber *uint8 `json:"alertOccurrenceNumber"`
	AcknowledgeSourceNetworkIdName *uint64 `json:"acknowledgeSourceNetworkIdName"`
	ResponseCommand AlertResponseCommandConst `json:"responseCommand"`
}
func (a *AlertResponse) PGNNumber() uint32  { return 126984 }
func DecodeAlertResponse(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AlertResponse{}
	val.Info = Info
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AlertResponse-AlertType: %w", err)
	} else {
		val.AlertType = AlertTypeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AlertResponse-AlertCategory: %w", err)
	} else {
		val.AlertCategory = AlertCategoryConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertResponse-AlertSystem: %w", err)
	} else {
		val.AlertSystem = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertResponse-AlertSubSystem: %w", err)
	} else {
		val.AlertSubSystem = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for AlertResponse-AlertId: %w", err)
	} else {
		val.AlertId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt64(64); err != nil {
		return nil, fmt.Errorf("parse failed for AlertResponse-DataSourceNetworkIdName: %w", err)
	} else {
		val.DataSourceNetworkIdName = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertResponse-DataSourceInstance: %w", err)
	} else {
		val.DataSourceInstance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertResponse-DataSourceIndexSource: %w", err)
	} else {
		val.DataSourceIndexSource = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertResponse-AlertOccurrenceNumber: %w", err)
	} else {
		val.AlertOccurrenceNumber = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt64(64); err != nil {
		return nil, fmt.Errorf("parse failed for AlertResponse-AcknowledgeSourceNetworkIdName: %w", err)
	} else {
		val.AcknowledgeSourceNetworkIdName = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for AlertResponse-ResponseCommand: %w", err)
	} else {
		val.ResponseCommand = AlertResponseCommandConst(v)

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

func EncodeAlertResponse(val *AlertResponse) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.AlertType), 4)
	w.writeLookupField(uint64(val.AlertCategory), 4)
	w.writeUInt8(val.AlertSystem, 8)
	w.writeUInt8(val.AlertSubSystem, 8)
	w.writeUInt16(val.AlertId, 16)
	w.writeUInt64(val.DataSourceNetworkIdName, 64)
	w.writeUInt8(val.DataSourceInstance, 8)
	w.writeUInt8(val.DataSourceIndexSource, 8)
	w.writeUInt8(val.AlertOccurrenceNumber, 8)
	w.writeUInt64(val.AcknowledgeSourceNetworkIdName, 64)
	w.writeLookupField(uint64(val.ResponseCommand), 2)
	w.skipBits(6)
	return w.Bytes(), w.Err()
}
func encodeAlertResponseMsg(v Message) ([]byte, error) {
	val, ok := v.(*AlertResponse)
	if !ok {
		return nil, fmt.Errorf("expected *AlertResponse, got %T", v)
	}
	return EncodeAlertResponse(val)
}
type AlertText struct {
	Info MessageInfo `json:"info"`
	AlertType AlertTypeConst `json:"alertType"`
	AlertCategory AlertCategoryConst `json:"alertCategory"`
	AlertSystem *uint8 `json:"alertSystem"`
	AlertSubSystem *uint8 `json:"alertSubSystem"`
	AlertId *uint16 `json:"alertId"`
	DataSourceNetworkIdName *uint64 `json:"dataSourceNetworkIdName"`
	DataSourceInstance *uint8 `json:"dataSourceInstance"`
	DataSourceIndexSource *uint8 `json:"dataSourceIndexSource"`
	AlertOccurrenceNumber *uint8 `json:"alertOccurrenceNumber"`
	LanguageId AlertLanguageIdConst `json:"languageId"`
	AlertTextDescription string `json:"alertTextDescription"`
	AlertLocationTextDescription string `json:"alertLocationTextDescription"`
}
func (a *AlertText) PGNNumber() uint32  { return 126985 }
func DecodeAlertText(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AlertText{}
	val.Info = Info
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AlertText-AlertType: %w", err)
	} else {
		val.AlertType = AlertTypeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AlertText-AlertCategory: %w", err)
	} else {
		val.AlertCategory = AlertCategoryConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertText-AlertSystem: %w", err)
	} else {
		val.AlertSystem = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertText-AlertSubSystem: %w", err)
	} else {
		val.AlertSubSystem = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for AlertText-AlertId: %w", err)
	} else {
		val.AlertId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt64(64); err != nil {
		return nil, fmt.Errorf("parse failed for AlertText-DataSourceNetworkIdName: %w", err)
	} else {
		val.DataSourceNetworkIdName = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertText-DataSourceInstance: %w", err)
	} else {
		val.DataSourceInstance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertText-DataSourceIndexSource: %w", err)
	} else {
		val.DataSourceIndexSource = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertText-AlertOccurrenceNumber: %w", err)
	} else {
		val.AlertOccurrenceNumber = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertText-LanguageId: %w", err)
	} else {
		val.LanguageId = AlertLanguageIdConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLengthAndControl(); err != nil {
		return nil, fmt.Errorf("parse failed for AlertText-AlertTextDescription: %w", err)
	} else {
		val.AlertTextDescription = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLengthAndControl(); err != nil {
		return nil, fmt.Errorf("parse failed for AlertText-AlertLocationTextDescription: %w", err)
	} else {
		val.AlertLocationTextDescription = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeAlertText(val *AlertText) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.AlertType), 4)
	w.writeLookupField(uint64(val.AlertCategory), 4)
	w.writeUInt8(val.AlertSystem, 8)
	w.writeUInt8(val.AlertSubSystem, 8)
	w.writeUInt16(val.AlertId, 16)
	w.writeUInt64(val.DataSourceNetworkIdName, 64)
	w.writeUInt8(val.DataSourceInstance, 8)
	w.writeUInt8(val.DataSourceIndexSource, 8)
	w.writeUInt8(val.AlertOccurrenceNumber, 8)
	w.writeLookupField(uint64(val.LanguageId), 8)
	w.writeStringWithLengthAndControl(val.AlertTextDescription)
	w.writeStringWithLengthAndControl(val.AlertLocationTextDescription)
	return w.Bytes(), w.Err()
}
func encodeAlertTextMsg(v Message) ([]byte, error) {
	val, ok := v.(*AlertText)
	if !ok {
		return nil, fmt.Errorf("expected *AlertText, got %T", v)
	}
	return EncodeAlertText(val)
}
type AlertConfiguration struct {
	Info MessageInfo `json:"info"`
	AlertType AlertTypeConst `json:"alertType"`
	AlertCategory AlertCategoryConst `json:"alertCategory"`
	AlertSystem *uint8 `json:"alertSystem"`
	AlertSubSystem *uint8 `json:"alertSubSystem"`
	AlertId *uint16 `json:"alertId"`
	DataSourceNetworkIdName *uint64 `json:"dataSourceNetworkIdName"`
	DataSourceInstance *uint8 `json:"dataSourceInstance"`
	DataSourceIndexSource *uint8 `json:"dataSourceIndexSource"`
	AlertOccurrenceNumber *uint8 `json:"alertOccurrenceNumber"`
	AlertControl *uint8 `json:"alertControl"`
	UserDefinedAlertAssignment *uint8 `json:"userDefinedAlertAssignment"`
	ReactivationPeriod *uint8 `json:"reactivationPeriod"`
	TemporarySilencePeriod *uint8 `json:"temporarySilencePeriod"`
	EscalationPeriod *uint8 `json:"escalationPeriod"`
}
func (a *AlertConfiguration) PGNNumber() uint32  { return 126986 }
func DecodeAlertConfiguration(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AlertConfiguration{}
	val.Info = Info
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AlertConfiguration-AlertType: %w", err)
	} else {
		val.AlertType = AlertTypeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AlertConfiguration-AlertCategory: %w", err)
	} else {
		val.AlertCategory = AlertCategoryConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertConfiguration-AlertSystem: %w", err)
	} else {
		val.AlertSystem = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertConfiguration-AlertSubSystem: %w", err)
	} else {
		val.AlertSubSystem = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for AlertConfiguration-AlertId: %w", err)
	} else {
		val.AlertId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt64(64); err != nil {
		return nil, fmt.Errorf("parse failed for AlertConfiguration-DataSourceNetworkIdName: %w", err)
	} else {
		val.DataSourceNetworkIdName = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertConfiguration-DataSourceInstance: %w", err)
	} else {
		val.DataSourceInstance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertConfiguration-DataSourceIndexSource: %w", err)
	} else {
		val.DataSourceIndexSource = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertConfiguration-AlertOccurrenceNumber: %w", err)
	} else {
		val.AlertOccurrenceNumber = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(2); err != nil {
		return nil, fmt.Errorf("parse failed for AlertConfiguration-AlertControl: %w", err)
	} else {
		val.AlertControl = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(2); err != nil {
		return nil, fmt.Errorf("parse failed for AlertConfiguration-UserDefinedAlertAssignment: %w", err)
	} else {
		val.UserDefinedAlertAssignment = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertConfiguration-ReactivationPeriod: %w", err)
	} else {
		val.ReactivationPeriod = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertConfiguration-TemporarySilencePeriod: %w", err)
	} else {
		val.TemporarySilencePeriod = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertConfiguration-EscalationPeriod: %w", err)
	} else {
		val.EscalationPeriod = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeAlertConfiguration(val *AlertConfiguration) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.AlertType), 4)
	w.writeLookupField(uint64(val.AlertCategory), 4)
	w.writeUInt8(val.AlertSystem, 8)
	w.writeUInt8(val.AlertSubSystem, 8)
	w.writeUInt16(val.AlertId, 16)
	w.writeUInt64(val.DataSourceNetworkIdName, 64)
	w.writeUInt8(val.DataSourceInstance, 8)
	w.writeUInt8(val.DataSourceIndexSource, 8)
	w.writeUInt8(val.AlertOccurrenceNumber, 8)
	w.writeUInt8(val.AlertControl, 2)
	w.writeUInt8(val.UserDefinedAlertAssignment, 2)
	w.skipBits(4)
	w.writeUInt8(val.ReactivationPeriod, 8)
	w.writeUInt8(val.TemporarySilencePeriod, 8)
	w.writeUInt8(val.EscalationPeriod, 8)
	return w.Bytes(), w.Err()
}
func encodeAlertConfigurationMsg(v Message) ([]byte, error) {
	val, ok := v.(*AlertConfiguration)
	if !ok {
		return nil, fmt.Errorf("expected *AlertConfiguration, got %T", v)
	}
	return EncodeAlertConfiguration(val)
}
type AlertThreshold struct {
	Info MessageInfo `json:"info"`
	AlertType AlertTypeConst `json:"alertType"`
	AlertCategory AlertCategoryConst `json:"alertCategory"`
	AlertSystem *uint8 `json:"alertSystem"`
	AlertSubSystem *uint8 `json:"alertSubSystem"`
	AlertId *uint16 `json:"alertId"`
	DataSourceNetworkIdName *uint64 `json:"dataSourceNetworkIdName"`
	DataSourceInstance *uint8 `json:"dataSourceInstance"`
	DataSourceIndexSource *uint8 `json:"dataSourceIndexSource"`
	AlertOccurrenceNumber *uint8 `json:"alertOccurrenceNumber"`
	NumberOfParameters *uint8 `json:"numberOfParameters"`
	Repeating1 []AlertThresholdRepeating1 `json:"repeating1"`
}
func (a *AlertThreshold) PGNNumber() uint32  { return 126987 }
type AlertThresholdRepeating1 struct {
	ParameterNumber *uint8 `json:"parameterNumber"`
	TriggerMethod *uint8 `json:"triggerMethod"`
	ThresholdDataFormat *uint8 `json:"thresholdDataFormat"`
	ThresholdLevel *uint64 `json:"thresholdLevel"`
}
func DecodeAlertThreshold(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AlertThreshold{}
	val.Info = Info
		var repeat1Count uint16 = 0
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AlertThreshold-AlertType: %w", err)
	} else {
		val.AlertType = AlertTypeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AlertThreshold-AlertCategory: %w", err)
	} else {
		val.AlertCategory = AlertCategoryConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertThreshold-AlertSystem: %w", err)
	} else {
		val.AlertSystem = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertThreshold-AlertSubSystem: %w", err)
	} else {
		val.AlertSubSystem = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for AlertThreshold-AlertId: %w", err)
	} else {
		val.AlertId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt64(64); err != nil {
		return nil, fmt.Errorf("parse failed for AlertThreshold-DataSourceNetworkIdName: %w", err)
	} else {
		val.DataSourceNetworkIdName = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertThreshold-DataSourceInstance: %w", err)
	} else {
		val.DataSourceInstance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertThreshold-DataSourceIndexSource: %w", err)
	} else {
		val.DataSourceIndexSource = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertThreshold-AlertOccurrenceNumber: %w", err)
	} else {
		val.AlertOccurrenceNumber = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertThreshold-NumberOfParameters: %w", err)
	} else {
		val.NumberOfParameters = v
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
	val.Repeating1 = make([]AlertThresholdRepeating1, 0)
	i := 0 
	for {
		var rep AlertThresholdRepeating1
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for AlertThreshold-ParameterNumber: %w", err)
		} else {
			rep.ParameterNumber = v
		}
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for AlertThreshold-TriggerMethod: %w", err)
		} else {
			rep.TriggerMethod = v
		}
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for AlertThreshold-ThresholdDataFormat: %w", err)
		} else {
			rep.ThresholdDataFormat = v
		}
		if v, err := stream.readUInt64(64); err != nil {
			return nil, fmt.Errorf("parse failed for AlertThreshold-ThresholdLevel: %w", err)
		} else {
			rep.ThresholdLevel = v
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

func EncodeAlertThreshold(val *AlertThreshold) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.AlertType), 4)
	w.writeLookupField(uint64(val.AlertCategory), 4)
	w.writeUInt8(val.AlertSystem, 8)
	w.writeUInt8(val.AlertSubSystem, 8)
	w.writeUInt16(val.AlertId, 16)
	w.writeUInt64(val.DataSourceNetworkIdName, 64)
	w.writeUInt8(val.DataSourceInstance, 8)
	w.writeUInt8(val.DataSourceIndexSource, 8)
	w.writeUInt8(val.AlertOccurrenceNumber, 8)
	w.writeUInt8(val.NumberOfParameters, 8)
	for _, rep := range val.Repeating1 {
		w.writeUInt8(rep.ParameterNumber, 8)
		w.writeUInt8(rep.TriggerMethod, 8)
		w.writeUInt8(rep.ThresholdDataFormat, 8)
		w.writeUInt64(rep.ThresholdLevel, 64)
	}
	return w.Bytes(), w.Err()
}
func encodeAlertThresholdMsg(v Message) ([]byte, error) {
	val, ok := v.(*AlertThreshold)
	if !ok {
		return nil, fmt.Errorf("expected *AlertThreshold, got %T", v)
	}
	return EncodeAlertThreshold(val)
}
type AlertValue struct {
	Info MessageInfo `json:"info"`
	AlertType AlertTypeConst `json:"alertType"`
	AlertCategory AlertCategoryConst `json:"alertCategory"`
	AlertSystem *uint8 `json:"alertSystem"`
	AlertSubSystem *uint8 `json:"alertSubSystem"`
	AlertId *uint16 `json:"alertId"`
	DataSourceNetworkIdName *uint64 `json:"dataSourceNetworkIdName"`
	DataSourceInstance *uint8 `json:"dataSourceInstance"`
	DataSourceIndexSource *uint8 `json:"dataSourceIndexSource"`
	AlertOccurrenceNumber *uint8 `json:"alertOccurrenceNumber"`
	NumberOfParameters *uint8 `json:"numberOfParameters"`
	Repeating1 []AlertValueRepeating1 `json:"repeating1"`
}
func (a *AlertValue) PGNNumber() uint32  { return 126988 }
type AlertValueRepeating1 struct {
	ValueParameterNumber *uint8 `json:"valueParameterNumber"`
	ValueDataFormat *uint8 `json:"valueDataFormat"`
	ValueData *uint64 `json:"valueData"`
}
func DecodeAlertValue(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &AlertValue{}
	val.Info = Info
		var repeat1Count uint16 = 0
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AlertValue-AlertType: %w", err)
	} else {
		val.AlertType = AlertTypeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for AlertValue-AlertCategory: %w", err)
	} else {
		val.AlertCategory = AlertCategoryConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertValue-AlertSystem: %w", err)
	} else {
		val.AlertSystem = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertValue-AlertSubSystem: %w", err)
	} else {
		val.AlertSubSystem = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for AlertValue-AlertId: %w", err)
	} else {
		val.AlertId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt64(64); err != nil {
		return nil, fmt.Errorf("parse failed for AlertValue-DataSourceNetworkIdName: %w", err)
	} else {
		val.DataSourceNetworkIdName = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertValue-DataSourceInstance: %w", err)
	} else {
		val.DataSourceInstance = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertValue-DataSourceIndexSource: %w", err)
	} else {
		val.DataSourceIndexSource = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertValue-AlertOccurrenceNumber: %w", err)
	} else {
		val.AlertOccurrenceNumber = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for AlertValue-NumberOfParameters: %w", err)
	} else {
		val.NumberOfParameters = v
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
	val.Repeating1 = make([]AlertValueRepeating1, 0)
	i := 0 
	for {
		var rep AlertValueRepeating1
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for AlertValue-ValueParameterNumber: %w", err)
		} else {
			rep.ValueParameterNumber = v
		}
		if v, err := stream.readUInt8(8); err != nil {
			return nil, fmt.Errorf("parse failed for AlertValue-ValueDataFormat: %w", err)
		} else {
			rep.ValueDataFormat = v
		}
		if v, err := stream.readUInt64(64); err != nil {
			return nil, fmt.Errorf("parse failed for AlertValue-ValueData: %w", err)
		} else {
			rep.ValueData = v
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

func EncodeAlertValue(val *AlertValue) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeLookupField(uint64(val.AlertType), 4)
	w.writeLookupField(uint64(val.AlertCategory), 4)
	w.writeUInt8(val.AlertSystem, 8)
	w.writeUInt8(val.AlertSubSystem, 8)
	w.writeUInt16(val.AlertId, 16)
	w.writeUInt64(val.DataSourceNetworkIdName, 64)
	w.writeUInt8(val.DataSourceInstance, 8)
	w.writeUInt8(val.DataSourceIndexSource, 8)
	w.writeUInt8(val.AlertOccurrenceNumber, 8)
	w.writeUInt8(val.NumberOfParameters, 8)
	for _, rep := range val.Repeating1 {
		w.writeUInt8(rep.ValueParameterNumber, 8)
		w.writeUInt8(rep.ValueDataFormat, 8)
		w.writeUInt64(rep.ValueData, 64)
	}
	return w.Bytes(), w.Err()
}
func encodeAlertValueMsg(v Message) ([]byte, error) {
	val, ok := v.(*AlertValue)
	if !ok {
		return nil, fmt.Errorf("expected *AlertValue, got %T", v)
	}
	return EncodeAlertValue(val)
}
type SystemTime struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	Source SystemTimeConst `json:"source"`
	Date *uint16 `json:"date"`
	Time *float32 `json:"time"`
}
func (s *SystemTime) PGNNumber() uint32  { return 126992 }
func DecodeSystemTime(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &SystemTime{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for SystemTime-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(4); err != nil {
		return nil, fmt.Errorf("parse failed for SystemTime-Source: %w", err)
	} else {
		val.Source = SystemTimeConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(4)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for SystemTime-Date: %w", err)
	} else {
		val.Date = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for SystemTime-Time: %w", err)
	} else {
		val.Time = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeSystemTime(val *SystemTime) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeLookupField(uint64(val.Source), 4)
	w.skipBits(4)
	w.writeUInt16(val.Date, 16)
	w.writeUnsignedResolution(val.Time, 32, 0.0001)
	return w.Bytes(), w.Err()
}
func encodeSystemTimeMsg(v Message) ([]byte, error) {
	val, ok := v.(*SystemTime)
	if !ok {
		return nil, fmt.Errorf("expected *SystemTime, got %T", v)
	}
	return EncodeSystemTime(val)
}
type Heartbeat struct {
	Info MessageInfo `json:"info"`
	DataTransmitOffset *float32 `json:"dataTransmitOffset"`
	SequenceCounter *uint8 `json:"sequenceCounter"`
	Controller1State ControllerStateConst `json:"controller1State"`
	Controller2State ControllerStateConst `json:"controller2State"`
	EquipmentStatus EquipmentStatusConst `json:"equipmentStatus"`
}
func (h *Heartbeat) PGNNumber() uint32  { return 126993 }
func DecodeHeartbeat(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &Heartbeat{}
	val.Info = Info
	if v, err := stream.readUnsignedResolution(16, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for Heartbeat-DataTransmitOffset: %w", err)
	} else {
		val.DataTransmitOffset = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for Heartbeat-SequenceCounter: %w", err)
	} else {
		val.SequenceCounter = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for Heartbeat-Controller1State: %w", err)
	} else {
		val.Controller1State = ControllerStateConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for Heartbeat-Controller2State: %w", err)
	} else {
		val.Controller2State = ControllerStateConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for Heartbeat-EquipmentStatus: %w", err)
	} else {
		val.EquipmentStatus = EquipmentStatusConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(34)
	if stream.isEOF() {
		return val, nil
		}	
	return val, nil
}

func EncodeHeartbeat(val *Heartbeat) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUnsignedResolution(val.DataTransmitOffset, 16, 0.001)
	w.writeUInt8(val.SequenceCounter, 8)
	w.writeLookupField(uint64(val.Controller1State), 2)
	w.writeLookupField(uint64(val.Controller2State), 2)
	w.writeLookupField(uint64(val.EquipmentStatus), 2)
	w.skipBits(34)
	return w.Bytes(), w.Err()
}
func encodeHeartbeatMsg(v Message) ([]byte, error) {
	val, ok := v.(*Heartbeat)
	if !ok {
		return nil, fmt.Errorf("expected *Heartbeat, got %T", v)
	}
	return EncodeHeartbeat(val)
}
type ProductInformation struct {
	Info MessageInfo `json:"info"`
	Nmea2000Version *float32 `json:"nmea2000Version"`
	ProductCode *uint16 `json:"productCode"`
	ModelId string `json:"modelId"`
	SoftwareVersionCode string `json:"softwareVersionCode"`
	ModelVersion string `json:"modelVersion"`
	ModelSerialCode string `json:"modelSerialCode"`
	CertificationLevel *uint8 `json:"certificationLevel"`
	LoadEquivalency *uint8 `json:"loadEquivalency"`
}
func (p *ProductInformation) PGNNumber() uint32  { return 126996 }
func DecodeProductInformation(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &ProductInformation{}
	val.Info = Info
	if v, err := stream.readUnsignedResolution(16, 0.001); err != nil {
		return nil, fmt.Errorf("parse failed for ProductInformation-Nmea2000Version: %w", err)
	} else {
		val.Nmea2000Version = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for ProductInformation-ProductCode: %w", err)
	} else {
		val.ProductCode = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for ProductInformation-ModelId: %w", err)
	} else {
		val.ModelId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for ProductInformation-SoftwareVersionCode: %w", err)
	} else {
		val.SoftwareVersionCode = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for ProductInformation-ModelVersion: %w", err)
	} else {
		val.ModelVersion = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readFixedString(256); err != nil {
		return nil, fmt.Errorf("parse failed for ProductInformation-ModelSerialCode: %w", err)
	} else {
		val.ModelSerialCode = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ProductInformation-CertificationLevel: %w", err)
	} else {
		val.CertificationLevel = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ProductInformation-LoadEquivalency: %w", err)
	} else {
		val.LoadEquivalency = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeProductInformation(val *ProductInformation) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUnsignedResolution(val.Nmea2000Version, 16, 0.001)
	w.writeUInt16(val.ProductCode, 16)
	w.writeFixedString(val.ModelId, 256)
	w.writeFixedString(val.SoftwareVersionCode, 256)
	w.writeFixedString(val.ModelVersion, 256)
	w.writeFixedString(val.ModelSerialCode, 256)
	w.writeUInt8(val.CertificationLevel, 8)
	w.writeUInt8(val.LoadEquivalency, 8)
	return w.Bytes(), w.Err()
}
func encodeProductInformationMsg(v Message) ([]byte, error) {
	val, ok := v.(*ProductInformation)
	if !ok {
		return nil, fmt.Errorf("expected *ProductInformation, got %T", v)
	}
	return EncodeProductInformation(val)
}
type ConfigurationInformation struct {
	Info MessageInfo `json:"info"`
	InstallationDescription1 string `json:"installationDescription1"`
	InstallationDescription2 string `json:"installationDescription2"`
	ManufacturerInformation string `json:"manufacturerInformation"`
}
func (c *ConfigurationInformation) PGNNumber() uint32  { return 126998 }
func DecodeConfigurationInformation(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &ConfigurationInformation{}
	val.Info = Info
	if v, err := stream.readStringWithLengthAndControl(); err != nil {
		return nil, fmt.Errorf("parse failed for ConfigurationInformation-InstallationDescription1: %w", err)
	} else {
		val.InstallationDescription1 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLengthAndControl(); err != nil {
		return nil, fmt.Errorf("parse failed for ConfigurationInformation-InstallationDescription2: %w", err)
	} else {
		val.InstallationDescription2 = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readStringWithLengthAndControl(); err != nil {
		return nil, fmt.Errorf("parse failed for ConfigurationInformation-ManufacturerInformation: %w", err)
	} else {
		val.ManufacturerInformation = v

		if stream.isEOF() {
			return val, nil
		} 
	}	
	return val, nil
}

func EncodeConfigurationInformation(val *ConfigurationInformation) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeStringWithLengthAndControl(val.InstallationDescription1)
	w.writeStringWithLengthAndControl(val.InstallationDescription2)
	w.writeStringWithLengthAndControl(val.ManufacturerInformation)
	return w.Bytes(), w.Err()
}
func encodeConfigurationInformationMsg(v Message) ([]byte, error) {
	val, ok := v.(*ConfigurationInformation)
	if !ok {
		return nil, fmt.Errorf("expected *ConfigurationInformation, got %T", v)
	}
	return EncodeConfigurationInformation(val)
}
type ManOverboardNotification struct {
	Info MessageInfo `json:"info"`
	Sid *uint8 `json:"sid"`
	MobEmitterId *uint32 `json:"mobEmitterId"`
	ManOverboardStatus MobStatusConst `json:"manOverboardStatus"`
	ActivationTime *float32 `json:"activationTime"`
	PositionSource MobPositionSourceConst `json:"positionSource"`
	PositionDate *uint16 `json:"positionDate"`
	PositionTime *float32 `json:"positionTime"`
	Latitude *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	CogReference DirectionReferenceConst `json:"cogReference"`
	Cog *float32 `json:"cog"`
	Sog *units.Velocity `json:"sog"`
	MmsiOfVesselOfOrigin *uint32 `json:"mmsiOfVesselOfOrigin"`
	MobEmitterBatteryLowStatus LowBatteryConst `json:"mobEmitterBatteryLowStatus"`
}
func (m *ManOverboardNotification) PGNNumber() uint32  { return 127233 }
func DecodeManOverboardNotification(Info MessageInfo, stream *PGNDataStream) (Message, error) {
	val := &ManOverboardNotification{}
	val.Info = Info
	if v, err := stream.readUInt8(8); err != nil {
		return nil, fmt.Errorf("parse failed for ManOverboardNotification-Sid: %w", err)
	} else {
		val.Sid = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for ManOverboardNotification-MobEmitterId: %w", err)
	} else {
		val.MobEmitterId = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for ManOverboardNotification-ManOverboardStatus: %w", err)
	} else {
		val.ManOverboardStatus = MobStatusConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(5)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUnsignedResolution(32, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for ManOverboardNotification-ActivationTime: %w", err)
	} else {
		val.ActivationTime = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for ManOverboardNotification-PositionSource: %w", err)
	} else {
		val.PositionSource = MobPositionSourceConst(v)

		if stream.isEOF() {
			return val, nil
		} 
	}
	stream.skipBits(5)
	if stream.isEOF() {
		return val, nil
		}
	if v, err := stream.readUInt16(16); err != nil {
		return nil, fmt.Errorf("parse failed for ManOverboardNotification-PositionDate: %w", err)
	} else {
		val.PositionDate = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(32, 0.0001); err != nil {
		return nil, fmt.Errorf("parse failed for ManOverboardNotification-PositionTime: %w", err)
	} else {
		val.PositionTime = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for ManOverboardNotification-Latitude: %w", err)
	} else {
		val.Latitude = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readSignedResolution64Override(32, 1e-07); err != nil {
		return nil, fmt.Errorf("parse failed for ManOverboardNotification-Longitude: %w", err)
	} else {
		val.Longitude = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(2); err != nil {
		return nil, fmt.Errorf("parse failed for ManOverboardNotification-CogReference: %w", err)
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
		return nil, fmt.Errorf("parse failed for ManOverboardNotification-Cog: %w", err)
	} else {
		val.Cog = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUnsignedResolution(16, 0.01); err != nil {
		return nil, fmt.Errorf("parse failed for ManOverboardNotification-Sog: %w", err)
	} else {
		val.Sog = nullableUnit(units.MetersPerSecond, v, units.NewVelocity)

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readUInt32(32); err != nil {
		return nil, fmt.Errorf("parse failed for ManOverboardNotification-MmsiOfVesselOfOrigin: %w", err)
	} else {
		val.MmsiOfVesselOfOrigin = v

		if stream.isEOF() {
			return val, nil
		} 
	}
	if v, err := stream.readLookupField(3); err != nil {
		return nil, fmt.Errorf("parse failed for ManOverboardNotification-MobEmitterBatteryLowStatus: %w", err)
	} else {
		val.MobEmitterBatteryLowStatus = LowBatteryConst(v)

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

func EncodeManOverboardNotification(val *ManOverboardNotification) ([]byte, error) {
	w := NewPGNDataStreamWriter()
	// TODO: cross-field validation not yet implemented
	w.writeUInt8(val.Sid, 8)
	w.writeUInt32(val.MobEmitterId, 32)
	w.writeLookupField(uint64(val.ManOverboardStatus), 3)
	w.skipBits(5)
	w.writeUnsignedResolution(val.ActivationTime, 32, 0.0001)
	w.writeLookupField(uint64(val.PositionSource), 3)
	w.skipBits(5)
	w.writeUInt16(val.PositionDate, 16)
	w.writeUnsignedResolution(val.PositionTime, 32, 0.0001)
	w.writeSignedResolution64Override(val.Latitude, 32, 1e-07)
	w.writeSignedResolution64Override(val.Longitude, 32, 1e-07)
	w.writeLookupField(uint64(val.CogReference), 2)
	w.skipBits(6)
	w.writeUnsignedResolution(val.Cog, 16, 0.0001)
	var sogRaw *float32
	if val.Sog != nil {
		sogRaw = &val.Sog.Value
	}
	w.writeUnsignedResolution(sogRaw, 16, 0.01)
	w.writeUInt32(val.MmsiOfVesselOfOrigin, 32)
	w.writeLookupField(uint64(val.MobEmitterBatteryLowStatus), 3)
	w.skipBits(5)
	return w.Bytes(), w.Err()
}
func encodeManOverboardNotificationMsg(v Message) ([]byte, error) {
	val, ok := v.(*ManOverboardNotification)
	if !ok {
		return nil, fmt.Errorf("expected *ManOverboardNotification, got %T", v)
	}
	return EncodeManOverboardNotification(val)
}
