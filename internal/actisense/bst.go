package actisense

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/brutella/can"
	"github.com/open-ships/n2k/internal/framer"
)

const (
	BSTN2KReceive  = 0x93
	BSTN2KTransmit = 0x94
	BSTCANFrame    = 0x95
	BSTBEMResponse = 0xA0
	BSTBEMCommand  = 0xA1
	BSTN2KMessage  = 0xD0

	MaxNMEA2000Payload = 1785
)

// Direction is the direction encoded by BST-95 and BST-D0.
type Direction uint8

const (
	DirectionReceived Direction = iota
	DirectionTransmitted
)

// MessageType identifies the NMEA 2000 transfer form reported by BST-D0.
type MessageType uint8

const (
	MessageSingle MessageType = iota
	MessageFast
	MessageTransport
	MessageUnknown
)

// Message is one assembled NMEA 2000 payload carried in BST-93, BST-94, or
// BST-D0. Data is always owned.
type Message struct {
	Priority       uint8
	PGN            uint32
	Destination    uint8
	Source         uint8
	HasSource      bool
	Data           []byte
	Timestamp      time.Duration
	HasTimestamp   bool
	Direction      Direction
	Type           MessageType
	InternalSource bool
	FastSequence   uint8
}

// DecodeMessage decodes an NMEA 2000 BST datagram. ok is false when the BST ID
// carries another protocol record.
func DecodeMessage(d Datagram) (message Message, ok bool, err error) {
	switch d.ID {
	case BSTN2KReceive:
		message, err = decode93(d.Payload)
		return message, true, err
	case BSTN2KTransmit:
		message, err = decode94(d.Payload)
		return message, true, err
	case BSTN2KMessage:
		message, err = decodeD0(d.Payload)
		return message, true, err
	default:
		return Message{}, false, nil
	}
}

func decode93(payload []byte) (Message, error) {
	const headerLen = 11
	if len(payload) < headerLen {
		return Message{}, fmt.Errorf("actisense: BST-93 payload is %d bytes; header requires %d", len(payload), headerLen)
	}
	dataLen := int(payload[10])
	if len(payload) != headerLen+dataLen {
		return Message{}, fmt.Errorf("actisense: BST-93 data length is %d, actual length is %d", dataLen, len(payload)-headerLen)
	}
	return Message{
		Priority:     payload[0] & 0x07,
		PGN:          decodePGN(payload[1], payload[2], payload[3]),
		Destination:  payload[4],
		Source:       payload[5],
		HasSource:    true,
		Data:         append([]byte(nil), payload[headerLen:]...),
		Timestamp:    time.Duration(binary.LittleEndian.Uint32(payload[6:10])) * time.Millisecond,
		HasTimestamp: true,
		Direction:    DirectionReceived,
	}, nil
}

func decode94(payload []byte) (Message, error) {
	const headerLen = 6
	if len(payload) < headerLen {
		return Message{}, fmt.Errorf("actisense: BST-94 payload is %d bytes; header requires %d", len(payload), headerLen)
	}
	dataLen := int(payload[5])
	if len(payload) != headerLen+dataLen {
		return Message{}, fmt.Errorf("actisense: BST-94 data length is %d, actual length is %d", dataLen, len(payload)-headerLen)
	}
	return Message{
		Priority:    payload[0] & 0x07,
		PGN:         decodePGN(payload[1], payload[2], payload[3]),
		Destination: payload[4],
		Data:        append([]byte(nil), payload[headerLen:]...),
		Direction:   DirectionTransmitted,
	}, nil
}

func decodeD0(payload []byte) (Message, error) {
	const headerLen = 10
	if len(payload) < headerLen {
		return Message{}, fmt.Errorf("actisense: BST-D0 payload is %d bytes; header requires %d", len(payload), headerLen)
	}
	if len(payload)-headerLen > MaxNMEA2000Payload {
		return Message{}, fmt.Errorf("actisense: BST-D0 NMEA 2000 payload is %d bytes; maximum is %d", len(payload)-headerLen, MaxNMEA2000Payload)
	}
	pduSpecific, pduFormat, dpp, control := payload[2], payload[3], payload[4], payload[5]
	return Message{
		Destination:    payload[0],
		Source:         payload[1],
		HasSource:      true,
		Priority:       (dpp >> 2) & 0x07,
		PGN:            decodePGN(pduSpecific, pduFormat, dpp),
		Direction:      Direction((control >> 3) & 0x01),
		InternalSource: control&(1<<4) != 0,
		FastSequence:   (control >> 5) & 0x07,
		Type:           MessageType(control & 0x03),
		Timestamp:      time.Duration(binary.LittleEndian.Uint32(payload[6:10])) * time.Millisecond,
		HasTimestamp:   true,
		Data:           append([]byte(nil), payload[headerLen:]...),
	}, nil
}

// EncodeMessage94 renders one legacy gateway-owned assembled message.
func EncodeMessage94(message Message) ([]byte, error) {
	if len(message.Data) > 223 {
		return nil, fmt.Errorf("actisense: PGN %d payload is %d bytes; BST-94 maximum is 223", message.PGN, len(message.Data))
	}
	pduSpecific := byte(message.PGN)
	if byte(message.PGN>>8) < 240 {
		// BST-94 carries the destination separately. NGT firmware expects the
		// canonical zero PS byte here for PDU1 PGNs.
		pduSpecific = 0
	}
	payload := make([]byte, 0, 6+len(message.Data))
	payload = append(payload,
		message.Priority&0x07,
		pduSpecific, byte(message.PGN>>8), byte(message.PGN>>16)&0x03,
		message.Destination,
		byte(len(message.Data)),
	)
	payload = append(payload, message.Data...)
	return EncodeDatagram(BSTN2KTransmit, payload)
}

func decodePGN(pduSpecific, pduFormat, dataPage byte) uint32 {
	pgn := uint32(dataPage&0x03)<<16 | uint32(pduFormat)<<8
	if pduFormat >= 240 {
		pgn |= uint32(pduSpecific)
	}
	return pgn
}

// EncodeMessageD0 renders an assembled NMEA 2000 message using BST Type 2.
func EncodeMessageD0(message Message) ([]byte, error) {
	if len(message.Data) > MaxNMEA2000Payload {
		return nil, fmt.Errorf("actisense: PGN %d payload is %d bytes; maximum is %d", message.PGN, len(message.Data), MaxNMEA2000Payload)
	}
	pduFormat := byte(message.PGN >> 8)
	pduSpecific := byte(message.PGN)
	if pduFormat < 240 {
		pduSpecific = message.Destination
	}
	dpp := byte(message.PGN>>16)&0x03 | (message.Priority&0x07)<<2
	control := byte(message.Type)&0x03 | byte(message.Direction&0x01)<<3 | (message.FastSequence&0x07)<<5
	if message.InternalSource {
		control |= 1 << 4
	}
	payload := make([]byte, 10, 10+len(message.Data))
	payload[0] = message.Destination
	payload[1] = message.Source
	payload[2] = pduSpecific
	payload[3] = pduFormat
	payload[4] = dpp
	payload[5] = control
	if message.HasTimestamp {
		binary.LittleEndian.PutUint32(payload[6:10], uint32(message.Timestamp/time.Millisecond))
	}
	payload = append(payload, message.Data...)
	return EncodeDatagram(BSTN2KMessage, payload)
}

// CANFrame is one raw CAN Frame plus BST-95 transport metadata.
type CANFrame struct {
	Frame        can.Frame
	Timestamp    time.Duration
	HasTimestamp bool
	Direction    Direction
	Resolution   uint8
}

// DecodeCANFrame decodes BST-95. ok is false for other BST IDs.
func DecodeCANFrame(d Datagram) (decoded CANFrame, ok bool, err error) {
	if d.ID != BSTCANFrame {
		return CANFrame{}, false, nil
	}
	if len(d.Payload) < 6 || len(d.Payload) > 14 {
		return CANFrame{}, true, fmt.Errorf("actisense: BST-95 payload is %d bytes; expected 6 through 14", len(d.Payload))
	}
	dppc := d.Payload[5]
	pduSpecific, pduFormat := d.Payload[3], d.Payload[4]
	pgn := uint32(dppc&0x03)<<16 | uint32(pduFormat)<<8
	destination := uint8(255)
	if pduFormat < 240 {
		destination = pduSpecific
	} else {
		pgn |= uint32(pduSpecific)
	}
	data := d.Payload[6:]
	frame := can.Frame{
		ID:     framer.BuildCANID(pgn, (dppc>>2)&0x07, d.Payload[2], destination),
		Length: uint8(len(data)),
	}
	copy(frame.Data[:], data)
	resolution := (dppc >> 5) & 0x03
	units := [...]time.Duration{time.Millisecond, 100 * time.Microsecond, 10 * time.Microsecond, time.Microsecond}
	return CANFrame{
		Frame:        frame,
		Timestamp:    time.Duration(binary.LittleEndian.Uint16(d.Payload[:2])) * units[resolution],
		HasTimestamp: true,
		Direction:    Direction((dppc >> 7) & 0x01),
		Resolution:   resolution,
	}, true, nil
}

// EncodeCANFrame renders a source-authoritative raw CAN frame as BST-95.
func EncodeCANFrame(frame can.Frame, direction Direction, timestamp time.Duration, resolution uint8) ([]byte, error) {
	if frame.Length > 8 {
		return nil, errors.New("actisense: classical CAN frame length exceeds 8")
	}
	if resolution > 3 {
		return nil, fmt.Errorf("actisense: timestamp resolution %d is outside 0-3", resolution)
	}
	id := framer.ParseCANID(frame.ID)
	pduFormat := byte(id.PGN >> 8)
	pduSpecific := byte(id.PGN)
	if pduFormat < 240 {
		pduSpecific = id.Destination
	}
	dppc := byte(id.PGN>>16)&0x03 | (id.Priority&0x07)<<2 | (resolution&0x03)<<5 | byte(direction&0x01)<<7
	units := [...]time.Duration{time.Millisecond, 100 * time.Microsecond, 10 * time.Microsecond, time.Microsecond}
	rawTimestamp := uint16(timestamp / units[resolution])
	payload := make([]byte, 6, 6+frame.Length)
	binary.LittleEndian.PutUint16(payload[:2], rawTimestamp)
	payload[2] = id.Source
	payload[3] = pduSpecific
	payload[4] = pduFormat
	payload[5] = dppc
	payload = append(payload, frame.Data[:frame.Length]...)
	return EncodeDatagram(BSTCANFrame, payload)
}
