package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandPath(t *testing.T) {
	t.Run("error when path has no file type", func(t *testing.T) {
		paths, err := expandPath("hours")
		assert.Nil(t, paths)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "path must include file type")
	})

	t.Run("no expansion when path has date", func(t *testing.T) {
		paths, err := expandPath("hours/api/20260101")
		assert.NoError(t, err)
		assert.Len(t, paths, 1)
		assert.Equal(t, "hours/api/20260101", paths[0])
	})

	t.Run("no expansion when path has days window", func(t *testing.T) {
		paths, err := expandPath("days/api")
		assert.NoError(t, err)
		assert.Len(t, paths, 1)
		assert.Equal(t, "days/api", paths[0])
	})

	t.Run("expands to 7 days when no date or days window", func(t *testing.T) {
		paths, err := expandPath("hours/api")
		assert.NoError(t, err)
		assert.Len(t, paths, 7)

		// verify each path has the expected format
		for _, p := range paths {
			assert.Regexp(t, `^hours/api/\d{8}/$`, p)
		}
	})

	t.Run("expands with leading slash", func(t *testing.T) {
		paths, err := expandPath("/hours/screenshot")
		assert.NoError(t, err)
		assert.Len(t, paths, 7)

		for _, p := range paths {
			assert.Regexp(t, `^/hours/screenshot/\d{8}/$`, p)
		}
	})

	t.Run("expands with trailing slash", func(t *testing.T) {
		paths, err := expandPath("hours/dom/")
		assert.NoError(t, err)
		assert.Len(t, paths, 7)

		for _, p := range paths {
			assert.Regexp(t, `^hours/dom/\d{8}/$`, p)
		}
	})
}
