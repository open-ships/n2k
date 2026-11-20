package pgn

import (
	"errors"
	"fmt"
	"reflect"
)

// Message identifies a PGN by number. Receive also returns UnknownPGN values
// implementing only this minimal interface. Encoding and transmission require
// the stronger PGN interface; CloneMessage additionally requires Clone() Message.
type Message interface {
	PGNNumber() uint32
}

// PGN is a message that can encode and decode its payload and carry CAN header
// metadata. Generated implementations own decoded payload bytes. Decode commits
// on success; Encode returns owned bytes and preserves an unchanged decoded
// payload exactly. SetMessageInfo replaces metadata and clears retained wire
// bookkeeping. A message must not be mutated concurrently with encoding, decoding,
// or cloning; use CloneMessage to create independently owned copies.
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

// DecodePayload selects a variant by its wire selectors before constructing
// its typed message. A more general variant can decode when a specific one
// cannot interpret a payload.
func DecodePayload(info MessageInfo, payload []byte) (PGN, error) {
	if len(payload) > maxPayloadBytes {
		return nil, errPayloadTooLarge
	}
	variants := decodeVariants[info.PGN]
	if len(variants) == 0 {
		return nil, fmt.Errorf("no PGN struct for %d", info.PGN)
	}
	var errs []error
	for _, variant := range variants {
		if err := variant.matchPayload(payload); err != nil {
			if errors.Is(err, ErrUnexpectedPayloadEnd) {
				errs = append(errs, err)
			}
			continue
		}
		msg := structTypeRegistry[variant.Id]()
		msg.SetMessageInfo(info)
		if err := decodeSelectedFields(msg, payload); err != nil {
			errs = append(errs, err)
			continue
		}
		return msg, nil
	}
	if len(errs) == 0 {
		return nil, fmt.Errorf("no matching PGN struct for %d", info.PGN)
	}
	return nil, fmt.Errorf("no matching PGN struct for %d: %w", info.PGN, errors.Join(errs...))
}

// EncodeMessage serializes msg by calling its EncodePayload method. The value
// must implement PGN; nil, typed-nil, and read-only Message implementations
// (including UnknownPGN) return an error. The returned bytes are caller-owned.
func EncodeMessage(msg Message) ([]byte, error) {
	if isNilMessage(msg) {
		return nil, errors.New("nil PGN message")
	}

	pgnMsg, ok := msg.(PGN)
	if !ok {
		return nil, fmt.Errorf("%T does not implement pgn.PGN", msg)
	}
	payload, err := pgnMsg.EncodePayload()
	if err == nil && len(payload) > maxPayloadBytes {
		return nil, errPayloadTooLarge
	}
	return payload, err
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
