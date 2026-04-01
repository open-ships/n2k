// Package n2k decodes NMEA 2000 marine network messages from CAN bus hardware
// into strongly-typed Go structs.
package n2k

import (
	"context"
	"iter"
)

// Receive returns an iterator of decoded NMEA 2000 messages from the configured sources.
// Each yielded value is a pointer to a typed PGN struct (e.g., *pgn.VesselHeading)
// or *pgn.UnknownPGN if IncludeUnknown() is set.
func Receive(ctx context.Context, opts ...Option) iter.Seq2[any, error] {
	return func(yield func(any, error) bool) {
		cfg := config{}
		for _, o := range opts {
			o.apply(&cfg)
		}
		if err := cfg.validate(); err != nil {
			yield(nil, err)
			return
		}

		s := NewScanner(ctx, opts...)
		for s.Next() {
			if !yield(s.Message(), nil) {
				return
			}
		}
		if s.Err() != nil {
			yield(nil, s.Err())
		}
	}
}
