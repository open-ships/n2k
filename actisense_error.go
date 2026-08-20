package n2k

import (
	"errors"
	"fmt"

	"github.com/open-ships/n2k/internal/gateway"
)

// ErrActisenseGatewaySessionRequired reports an attempt to run a source-
// authoritative Client over gateway-owned BST-93/94 traffic. Use a public
// ActisenseGatewaySession for legacy message sends, or BST-95 raw mode for a
// Client with its own NMEA 2000 identity.
var ErrActisenseGatewaySessionRequired = errors.New("n2k: gateway-owned Actisense message formats cannot back Client; use NewActisenseTCPSession/NewActisenseSerialSession or a source-authoritative raw CAN format")

// ActisenseModeError reports that an Actisense gateway did not acknowledge a
// requested operating mode. RequestedMode is the Actisense wire value; raw
// CAN mode is 5. The wrapped error retains timeout, disconnect, negative
// acknowledgement, or device-error details.
type ActisenseModeError struct {
	RequestedMode uint16
	Err           error
}

func (e *ActisenseModeError) Error() string {
	return fmt.Sprintf("n2k: Actisense operating mode %d setup failed: %v", e.RequestedMode, e.Err)
}

// Unwrap returns the underlying gateway setup error.
func (e *ActisenseModeError) Unwrap() error { return e.Err }

func wrapActisenseModeError(err error) error {
	if err == nil {
		return nil
	}
	var publicErr *ActisenseModeError
	if errors.As(err, &publicErr) {
		return err
	}
	var setupErr *gateway.ActisenseModeSetupError
	if !errors.As(err, &setupErr) {
		return err
	}
	return &ActisenseModeError{RequestedMode: uint16(setupErr.RequestedMode), Err: err}
}
