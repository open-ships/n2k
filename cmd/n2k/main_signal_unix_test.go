//go:build aix || android || darwin || dragonfly || freebsd || illumos || ios || linux || netbsd || openbsd || solaris

package main

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSniffTerminatesCleanlyOnSIGTERM(t *testing.T) {
	if os.Getenv("N2K_SIGTERM_HELPER") == "1" {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
		defer stop()
		root := newRootCommand(os.Stdin, os.Stdout, os.Stderr)
		root.SetArgs([]string{"sniff", "--file", "../../testdata/sample.log", "--timing", "--unknown"})
		err := root.ExecuteContext(ctx)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestSniffTerminatesCleanlyOnSIGTERM$")
	cmd.Env = append(os.Environ(), "N2K_SIGTERM_HELPER=1")
	cmd.Stderr = io.Discard
	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)
	require.NoError(t, cmd.Start())

	ready := make(chan error, 1)
	go func() {
		_, err := bufio.NewReader(stdout).ReadBytes('\n')
		ready <- err
	}()
	select {
	case err := <-ready:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		t.Fatal("sniff subprocess did not become ready")
	}

	require.NoError(t, cmd.Process.Signal(syscall.SIGTERM))
	require.NoError(t, cmd.Wait())
}
