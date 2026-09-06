package n2k

import (
	"context"
	"fmt"

	"github.com/open-ships/n2k/internal/actisense"
	"github.com/open-ships/n2k/internal/gateway"
)

// ErrActisenseNotReady means no acknowledged gateway connection is available.
// The operation was not queued for a future connection.
var ErrActisenseNotReady = gateway.ErrActisenseNotReady

const ActisenseMaxRawWrite = 65536

// SendBST sends checksum-free BST bytes (ID, length, payload), adding the
// checksum, DLE escaping, and BDTP framing. Unknown BST IDs are supported.
// The input must contain exactly one complete datagram, at most 1800 bytes.
// This is an explicit protocol escape hatch; the caller selects a record
// appropriate for the connected device's operating mode.
func (s *ActisenseGatewaySession) SendBST(ctx context.Context, bst []byte) error {
	datagram, err := actisense.DecodeRaw(bst)
	if err != nil {
		return fmt.Errorf("n2k: Actisense BST send: %w", err)
	}
	wire, err := actisense.EncodeDatagram(datagram.ID, datagram.Payload)
	if err != nil {
		return err
	}
	return s.SendRaw(ctx, wire)
}

// SendRaw sends an owned snapshot of verbatim wire bytes through the same
// serialized writer as BEM and PGN operations. It adds no checksum or framing.
// Each call is limited to ActisenseMaxRawWrite bytes and the command timeout.
// Cancellation or partial failure closes that connection; writes are never
// retried on a reconnect. It does not implement any additional wire protocol.
func (s *ActisenseGatewaySession) SendRaw(ctx context.Context, wire []byte) error {
	if s == nil {
		return actisense.ErrSessionClosed
	}
	if len(wire) == 0 || len(wire) > ActisenseMaxRawWrite {
		return fmt.Errorf("n2k: Actisense raw write must contain 1-%d bytes; got %d", ActisenseMaxRawWrite, len(wire))
	}
	s.mu.Lock()
	closed := s.closed
	s.mu.Unlock()
	if closed {
		return actisense.ErrSessionClosed
	}
	writeCtx, cancel := s.writeContext(ctx)
	defer cancel()
	return s.transport.WriteContext(writeCtx, append([]byte(nil), wire...))
}
