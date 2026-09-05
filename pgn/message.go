package pgn

import (
	"errors"
	"fmt"
	"reflect"
)

type Message interface {
	PGNNumber() uint32
}

type PGN interface {
	Message
	MessageInfo() MessageInfo
	SetMessageInfo(MessageInfo)
	DecodePayload([]uint8) error
	EncodePayload() ([]uint8, error)
}

// DecodeMessage decodes a raw PGN payload into the matching PGN struct.
func DecodeMessage(info MessageInfo, payload []uint8) (PGN, error) {
	return DecodePayload(info, payload)
}

// EncodeMessage serializes msg by calling its EncodePayload method.
func EncodeMessage(msg Message) ([]byte, error) {
	if isNilMessage(msg) {
		return nil, errors.New("nil PGN message")
	}

	pgnMsg, ok := msg.(PGN)
	if !ok {
		return nil, fmt.Errorf("%T does not implement pgn.PGN", msg)
	}
	return pgnMsg.EncodePayload()
}

func isNilMessage(msg Message) bool {
	if msg == nil {
		return true
	}
	value := reflect.ValueOf(msg)
	switch value.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		return value.IsNil()
	default:
		return false
	}
}
