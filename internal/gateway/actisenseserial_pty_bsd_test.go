//go:build darwin || freebsd || openbsd

package gateway

import (
	"testing"

	"github.com/creack/pty"
	"github.com/stretchr/testify/require"
)

func actisenseTestPTY(t *testing.T) string {
	t.Helper()
	master, slave, err := pty.Open()
	require.NoError(t, err)
	t.Cleanup(func() { _ = master.Close() })
	t.Cleanup(func() { _ = slave.Close() })
	return slave.Name()
}
