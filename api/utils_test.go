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
			input:    "foo",
			expected: "/api/v1/foo",
		},
		{
			input:    "/foo",
			expected: "/api/v1/foo",
		},
		{
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
