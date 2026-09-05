package n2k

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/open-ships/n2k/internal/conformance"
	"github.com/stretchr/testify/require"
)

func TestConformanceRequirementIndex(t *testing.T) {
	data, err := os.ReadFile("conformance/requirements.json")
	require.NoError(t, err)
	var index conformance.Index
	require.NoError(t, json.Unmarshal(data, &index))
	require.NoError(t, conformance.Validate(index))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	catalog, err := conformance.Discover(ctx, ".", index)
	require.NoError(t, err)
	_, err = conformance.Resolve(index, catalog)
	require.NoError(t, err, "every evidence reference must match an executable test in its declared package")
}
