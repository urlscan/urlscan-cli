package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
