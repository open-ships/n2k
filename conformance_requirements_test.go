package n2k

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type requirementIndex struct {
	SchemaVersion int `json:"schemaVersion"`
	Requirements  []struct {
		ID       string   `json:"id"`
		Category string   `json:"category"`
		Behavior string   `json:"behavior"`
		Evidence []string `json:"evidence"`
	} `json:"requirements"`
}

func TestConformanceRequirementIndex(t *testing.T) {
	data, err := os.ReadFile("conformance/requirements.json")
	require.NoError(t, err)
	var index requirementIndex
	require.NoError(t, json.Unmarshal(data, &index))
	require.Equal(t, 1, index.SchemaVersion)
	require.NotEmpty(t, index.Requirements)
	seen := make(map[string]struct{}, len(index.Requirements))
	for _, requirement := range index.Requirements {
		require.NotEmpty(t, requirement.ID)
		require.NotEmpty(t, requirement.Category, requirement.ID)
		require.NotEmpty(t, requirement.Behavior, requirement.ID)
		require.NotEmpty(t, requirement.Evidence, requirement.ID)
		_, duplicate := seen[requirement.ID]
		require.False(t, duplicate, "duplicate requirement ID %s", requirement.ID)
		seen[requirement.ID] = struct{}{}
	}
}
