//go:build linux || darwin || freebsd || openbsd

package gateway

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.bug.st/serial"
)

// A real pseudo-terminal can be supplied by the platform's PTY allocator.
// The peer must remain open without draining the terminal's output.
func TestActisenseSerialPseudoTerminalCancellation(t *testing.T) {
	path := os.Getenv("N2K_TEST_SERIAL_PTY")
	if path == "" {
		path = actisenseTestPTY(t)
	}
	connection, err := openActisenseSerialConnection(path, &serial.Mode{BaudRate: 115200, DataBits: 8})
	require.NoError(t, err)
	defer func() { _ = connection.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- writeActisenseUnitContext(ctx, connection, make([]byte, 1<<20)) }()
	select {
	case err := <-done:
		require.ErrorIs(t, err, context.DeadlineExceeded)
	case <-time.After(time.Second):
		t.Fatal("serial write on a real pseudo-terminal ignored its deadline")
	}
}
