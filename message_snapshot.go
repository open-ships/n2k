package n2k

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/open-ships/n2k/pgn"
)

// messageSnapshot is an owned wire payload admitted to a write queue. Encoding
// happens while the caller still owns its message, before Write returns.
type messageSnapshot struct {
	info    pgn.MessageInfo
	number  uint32
	payload []byte
}

func snapshotMessage(msg pgn.Message) (snapshot *messageSnapshot, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("n2k: panic encoding message: %v", r)
		}
	}()
	if msg == nil {
		return nil, errors.New("n2k: cannot write a nil message")
	}
	value := reflect.ValueOf(msg)
	if value.Kind() == reflect.Pointer && value.IsNil() {
		return nil, fmt.Errorf("n2k: cannot write a nil %T", msg)
	}
	typed, ok := msg.(pgn.PGN)
	if !ok {
		return nil, fmt.Errorf("n2k: %T does not implement pgn.PGN", msg)
	}
	payload, err := typed.EncodePayload()
	if err != nil {
		return nil, fmt.Errorf("n2k: encode PGN %d: %w", msg.PGNNumber(), err)
	}
	if len(payload) > 1785 {
		return nil, errors.New("n2k: payload exceeds 1785 bytes")
	}
	return &messageSnapshot{number: msg.PGNNumber(), info: typed.MessageInfo().Clone(), payload: append([]byte(nil), payload...)}, nil
}

func (m *messageSnapshot) PGNNumber() uint32                   { return m.number }
func (m *messageSnapshot) MessageInfo() pgn.MessageInfo        { return m.info.Clone() }
func (m *messageSnapshot) SetMessageInfo(info pgn.MessageInfo) { m.info = info.Clone() }
func (m *messageSnapshot) EncodePayload() ([]byte, error) {
	return append([]byte(nil), m.payload...), nil
}
func (m *messageSnapshot) DecodePayload([]byte) error {
	return errors.New("n2k: write snapshots cannot decode")
}
