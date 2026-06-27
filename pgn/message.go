package pgn

import (
	"errors"
	"fmt"
)

type Message interface {
	PGNNumber() uint32
}

// EncodeMessage serializes msg using the encoder registered for its concrete PGN
// variant. Multiple CANboat variants can share a numeric PGN, so callers must not
// select encoders by PGN number alone.
func EncodeMessage(msg Message) ([]byte, error) {
	if msg == nil {
		return nil, errors.New("nil PGN message")
	}

	pgnNum := msg.PGNNumber()
	candidates := PgnInfoLookup[pgnNum]
	var errs []error
	tried := false

	for _, candidate := range candidates {
		if candidate.Encoder == nil {
			continue
		}
		tried = true
		payload, err := candidate.Encoder(msg)
		if err == nil {
			return payload, nil
		}
		errs = append(errs, err)
	}

	// Preserve compatibility for callers/tests that populate EncoderLookup directly,
	// while keeping PgnInfoLookup as the authoritative variant-aware registry.
	if !tried {
		if encoder, ok := EncoderLookup[pgnNum]; ok {
			return encoder(msg)
		}
		return nil, fmt.Errorf("no encoder registered for PGN %d", pgnNum)
	}

	return nil, fmt.Errorf("no encoder matched %T for PGN %d: %w", msg, pgnNum, errors.Join(errs...))
}
