package utils

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringArrayReader(t *testing.T) {
	args := []string{"one", "two", "three"}
	reader := NewStringArrayReader(args)

	for _, expected := range args {
		got, err := reader.ReadString()
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	}

	_, err := reader.ReadString()
	assert.Equal(t, io.EOF, err)
}

func TestStringIOReader(t *testing.T) {
	input := "foo\nbar\nbaz"
	reader := NewStringIOReader(strings.NewReader(input))

	expectedLines := []string{"foo", "bar", "baz"}
	for _, expected := range expectedLines {
		got, err := reader.ReadString()
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	}

	_, err := reader.ReadString()
	assert.Equal(t, io.EOF, err)
}

func TestMappedStringReader(t *testing.T) {
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
		reader := NewMappedStringReader(StringReaderFromCmdArgs(test.input), ResolveFileOrValue)
		got, err := reader.ReadAll()
		assert.NoError(t, err)
		assert.Equal(t, test.expected, got)
	}
}
