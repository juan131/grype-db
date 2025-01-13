package processors

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/anchore/grype-db/pkg/data"
	"github.com/anchore/grype-db/pkg/process/internal/tests"
	"github.com/anchore/grype-db/pkg/provider"
	"github.com/anchore/grype-db/pkg/provider/unmarshal"
)

func mockMatchExclusionProcessorTransform(vulnerability unmarshal.MatchExclusion) ([]data.Entry, error) {
	return []data.Entry{
		{
			DBSchemaVersion: 0,
			Data:            vulnerability,
		},
	}, nil
}

func TestMatchExclusionProcessor_Process(t *testing.T) {
	f, err := os.Open("test-fixtures/exclusions.json")
	require.NoError(t, err)
	defer tests.CloseFile(f)

	processor := NewMatchExclusionProcessor(mockMatchExclusionProcessorTransform)
	entries, err := processor.Process(f, provider.State{
		Provider: "match-exclusions",
	})

	require.NoError(t, err)
	assert.Len(t, entries, 3)
}

func TestMatchExclusionProcessor_IsSupported(t *testing.T) {
	tc := []struct {
		name      string
		schemaURL string
		expected  bool
	}{
		{
			name:      "valid schema URL with version 1.0.0",
			schemaURL: "https://example.com/vunnel/path/match-exclusion/schema-1.0.0.json",
			expected:  true,
		},
		{
			name:      "valid schema URL with version 1.3.4",
			schemaURL: "https://example.com/vunnel/path/match-exclusion/schema-1.3.4.json",
			expected:  true,
		},
		{
			name:      "invalid schema URL with unsupported version",
			schemaURL: "https://example.com/vunnel/path/match-exclusion/schema-2.0.0.json",
			expected:  false,
		},
		{
			name:      "invalid schema URL with missing version",
			schemaURL: "https://example.com/vunnel/path/match-exclusion/schema.json",
			expected:  false,
		},
		{
			name:      "completely invalid URL",
			schemaURL: "https://example.com/invalid/schema/url",
			expected:  false,
		},
	}

	p := matchExclusionProcessor{}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, p.IsSupported(tt.schemaURL))
		})
	}
}
