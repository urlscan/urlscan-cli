package utils

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/urlscan/urlscan-cli/api"
)

func newTestClient() *APIClient {
	c := api.NewClient("dummy")
	c.SetBaseURL(&url.URL{
		Scheme: "http",
		Host:   "testserver",
	})
	return &APIClient{Client: c}
}

func TestDownload(t *testing.T) {
	defer gock.Off()

	t.Run("downloads file when it doesn't exist", func(t *testing.T) {
		defer gock.Clean()

		gock.New("http://testserver").
			Get("/test").
			Reply(200).
			BodyString("test content")

		tmpDir := t.TempDir()
		outputPath := filepath.Join(tmpDir, "test.txt")

		opts := NewDownloadOptions(
			WithDownloadClient(newTestClient()),
			WithDownloadURL("/test"),
			WithDownloadOutput(outputPath),
			WithDownloadForce(false),
			WithDownloadSilent(true),
		)

		err := Download(opts)
		assert.NoError(t, err)
		assert.FileExists(t, outputPath)

		content, err := os.ReadFile(outputPath)
		assert.NoError(t, err)
		assert.Equal(t, "test content", string(content))

		assert.True(t, gock.IsDone())
	})

	t.Run("overwrites existing file when force is true", func(t *testing.T) {
		defer gock.Clean()

		gock.New("http://testserver").
			Get("/test").
			Reply(200).
			BodyString("test content")

		tmpDir := t.TempDir()
		outputPath := filepath.Join(tmpDir, "existing.txt")

		err := os.WriteFile(outputPath, []byte("old content"), 0o644)
		assert.NoError(t, err)

		opts := NewDownloadOptions(
			WithDownloadClient(newTestClient()),
			WithDownloadURL("/test"),
			WithDownloadOutput(outputPath),
			WithDownloadForce(true),
			WithDownloadSilent(true),
		)

		err = Download(opts)
		assert.NoError(t, err)
		assert.FileExists(t, outputPath)

		content, err := os.ReadFile(outputPath)
		assert.NoError(t, err)
		assert.Equal(t, "test content", string(content))

		assert.True(t, gock.IsDone())
	})

	t.Run("returns error when file exists and force is false", func(t *testing.T) {
		defer gock.Clean()

		gock.New("http://testserver").
			Get("/test").
			Reply(200).
			BodyString("test content")

		tmpDir := t.TempDir()
		outputPath := filepath.Join(tmpDir, "existing.txt")

		err := os.WriteFile(outputPath, []byte("existing content"), 0o644)
		assert.NoError(t, err)

		opts := NewDownloadOptions(
			WithDownloadClient(newTestClient()),
			WithDownloadURL("/test"),
			WithDownloadOutput(outputPath),
			WithDownloadForce(false),
			WithDownloadSilent(true),
		)

		err = Download(opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")

		content, err := os.ReadFile(outputPath)
		assert.NoError(t, err)
		assert.Equal(t, "existing content", string(content))

		// HTTP request should not have been made since file already exists
		assert.False(t, gock.IsDone())
	})
}
