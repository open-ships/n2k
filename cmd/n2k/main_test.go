package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSniffContextCancellationIsClean(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := sniffContext(ctx, []string{"-file", "../../testdata/sample.log", "-timing"})
	require.NoError(t, err)
}
