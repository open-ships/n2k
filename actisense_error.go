package n2k

import (
	"errors"
	"fmt"

	"github.com/open-ships/n2k/internal/gateway"
)

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
