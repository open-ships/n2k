// Package n2k is a standalone Go toolkit for NMEA 2000 marine networks. It
// reads and writes messages from CAN hardware, USB-CAN adapters, network
// gateways, and capture/replay sources, decoding PGNs into strongly typed Go
// structs from package pgn.
package n2k

import (
	"context"
	"iter"

	"github.com/open-ships/n2k/pgn"
)

// Receive returns an iterator of decoded NMEA 2000 messages from the configured sources.
// Each yielded value is a pointer to a PGN struct (e.g., *pgn.VesselHeading)
// or *pgn.UnknownPGN if IncludeUnknown() is set.
func Receive(ctx context.Context, opts ...Option) iter.Seq2[pgn.Message, error] {
	return func(yield func(pgn.Message, error) bool) {
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
