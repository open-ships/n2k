package main

import (
	"bytes"
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

	_, _, err := executeCommand(ctx, "sniff", "--file", "../../testdata/sample.log", "--timing")
	require.NoError(t, err)
}

func TestRecordContextWritesReplayableCandump(t *testing.T) {
	output := filepath.Join(t.TempDir(), "capture.log")
	_, _, err := executeCommand(
		context.Background(),
		"record", "--file", "../../testdata/sample.log", "--out", output,
	)
	require.NoError(t, err)
	data, err := os.ReadFile(output)
	require.NoError(t, err)
	text := string(data)
	require.Contains(t, text, "#")
	require.True(t, strings.HasPrefix(text, "("), text)
}

func TestRootHelpIsDiscoverable(t *testing.T) {
	stdout, _, err := executeCommand(context.Background(), "--help")
	require.NoError(t, err)
	require.Contains(t, stdout, "NMEA 2000 capture, diagnostics, and schema tools")
	require.Contains(t, stdout, "completion")
	require.Contains(t, stdout, "devices")
	require.Contains(t, stdout, "validate")
}

func TestCommandHelpUsesCanonicalLongFlags(t *testing.T) {
	stdout, _, err := executeCommand(context.Background(), "sniff", "--help")
	require.NoError(t, err)
	require.Contains(t, stdout, "--tcp")
	require.Contains(t, stdout, "--file")
	require.Contains(t, stdout, "-i, --interface")
	require.NotContains(t, stdout, " -tcp")
	require.NotContains(t, stdout, " -file")
}

func TestCompletionCommandGeneratesShellScripts(t *testing.T) {
	tests := map[string]string{
		"bash":       "__start_n2k",
		"fish":       "complete -c n2k",
		"powershell": "Register-ArgumentCompleter",
		"zsh":        "#compdef n2k",
	}
	for shell, marker := range tests {
		t.Run(shell, func(t *testing.T) {
			stdout, _, err := executeCommand(context.Background(), "completion", shell)
			require.NoError(t, err)
			require.Contains(t, stdout, marker)
		})
	}
}

func TestUnknownCommandSuggestsClosestCommand(t *testing.T) {
	_, _, err := executeCommand(context.Background(), "snif")
	require.Error(t, err)
	require.Contains(t, err.Error(), "sniff")
}

func TestPGNCompletionIncludesDescriptions(t *testing.T) {
	completions, directive := completePGNs(nil, nil, "12725")
	require.NotEmpty(t, completions)
	require.Contains(t, strings.Join(completions, "\n"), "127250\tVessel Heading")
	require.NotZero(t, directive)
}

func TestVersionCommandUsesConfiguredWriter(t *testing.T) {
	stdout, _, err := executeCommand(context.Background(), "version")
	require.NoError(t, err)
	require.Equal(t, "n2k dev (commit none, built unknown)\n", stdout)
}

func executeCommand(ctx context.Context, args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := newRootCommand(strings.NewReader(""), &stdout, &stderr)
	root.SetArgs(args)
	err := root.ExecuteContext(ctx)
	return stdout.String(), stderr.String(), err
}
