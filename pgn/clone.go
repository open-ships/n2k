package pgn

import (
	"fmt"
)

// CloneMessage copies a message at an ownership boundary. Generated messages
// own all fields and retained wire bytes in the result. Custom message types
// can participate by implementing Clone() Message with the same ownership
// contract; types without that contract return an error.
func CloneMessage(msg Message) (Message, error) {
	if isNilMessage(msg) {
		return nil, fmt.Errorf("cannot clone nil PGN message")
	}
	cloner, ok := msg.(interface{ Clone() Message })
	if !ok {
		return nil, fmt.Errorf("%T does not implement Clone() pgn.Message", msg)
	}
	clone := cloner.Clone()
	if isNilMessage(clone) {
		return nil, fmt.Errorf("%T returned a nil PGN clone", msg)
	}
	return clone, nil
}

func clonePointer[T any](value *T) *T {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneSlice[T any](value []T) []T {
	if value == nil {
		return nil
	}
	return append(make([]T, 0, len(value)), value...)
}
