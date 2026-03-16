package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrefixedPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "without leading /",
			input:    "foo",
			expected: "/api/v1/foo",
		},
		{
			name:     "with leading /",
			input:    "/foo",
			expected: "/api/v1/foo",
		},
		{
			name:     "without leading /, with trailing /",
			input:    "foo/",
			expected: "/api/v1/foo/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PrefixedPath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMarshalJSONWithExtra(t *testing.T) {
	type inner struct {
		Name string `json:"name"`
	}
	type opts struct {
		Inner inner          `json:"inner"`
		Extra map[string]any `json:"-"`
	}

	t.Run("without extra", func(t *testing.T) {
		v := opts{Inner: inner{Name: "test"}, Extra: nil}
		b, err := marshalJSONWithExtra(v, v.Extra)
		require.NoError(t, err)

		var m map[string]any
		require.NoError(t, json.Unmarshal(b, &m))
		assert.Len(t, m, 1)
		assert.Contains(t, m, "inner")
	})

	t.Run("with extra", func(t *testing.T) {
		v := opts{
			Inner: inner{Name: "test"},
			Extra: map[string]any{"custom": "value"},
		}
		b, err := marshalJSONWithExtra(v, v.Extra)
		require.NoError(t, err)

		var m map[string]any
		require.NoError(t, json.Unmarshal(b, &m))
		assert.Len(t, m, 2)
		assert.Contains(t, m, "inner")
		assert.Equal(t, "value", m["custom"])
	})

	t.Run("extra overrides struct field", func(t *testing.T) {
		v := opts{
			Inner: inner{Name: "test"},
			Extra: map[string]any{"inner": "overridden"},
		}
		b, err := marshalJSONWithExtra(v, v.Extra)
		require.NoError(t, err)

		var m map[string]any
		require.NoError(t, json.Unmarshal(b, &m))
		assert.Equal(t, "overridden", m["inner"])
	})

	t.Run("empty extra map", func(t *testing.T) {
		v := opts{
			Inner: inner{Name: "test"},
			Extra: map[string]any{},
		}
		b, err := marshalJSONWithExtra(v, v.Extra)
		require.NoError(t, err)

		var m map[string]any
		require.NoError(t, json.Unmarshal(b, &m))
		assert.Len(t, m, 1)
	})
}
