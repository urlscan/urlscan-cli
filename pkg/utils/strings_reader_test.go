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

func TestMappedStringsReader(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{[]string{"foo"}, []string{"foo"}},
		{[]string{"testdata/input.txt"}, []string{"foo", "bar", "baz", "qux"}},
		{[]string{"testdata/input.txt", "a", "b", "c"}, []string{"foo", "bar", "baz", "qux", "a", "b", "c"}},
		{[]string{"testdata/404.txt", "a", "b", "c"}, []string{"testdata/404.txt", "a", "b", "c"}},
	}

	for _, test := range tests {
		reader := NewMappedStringsReader(StringReaderFromCmdArgs(test.input), ResolveFileOrValue)
		got, err := reader.ReadAll()
		assert.NoError(t, err)
		assert.Equal(t, test.expected, got)
	}
}
