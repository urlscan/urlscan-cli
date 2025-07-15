package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveFileOrValue(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"foo", []string{"foo"}},
		// with a file that exists
		{"testdata/input.txt", []string{"foo", "bar", "baz", "qux"}},
		// with a file that does not exist
		{"testdata/404.txt", []string{"testdata/404.txt"}},
	}

	for _, test := range tests {
		result, err := ResolveFileOrValue(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, result)
	}
}
