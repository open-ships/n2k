package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSniffContextCancellationIsClean(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := sniffContext(ctx, []string{"-file", "../../testdata/sample.log", "-timing"})
	require.NoError(t, err)
}

func TestRecordContextWritesReplayableCandump(t *testing.T) {
	output := filepath.Join(t.TempDir(), "capture.log")
	err := recordContext(context.Background(), []string{
		"-file", "../../testdata/sample.log",
		"-out", output,
	})
	require.NoError(t, err)
	data, err := os.ReadFile(output)
	require.NoError(t, err)
	text := string(data)
	require.Contains(t, text, "#")
	require.True(t, strings.HasPrefix(text, "("), text)
}
