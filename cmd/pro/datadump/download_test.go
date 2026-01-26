package datadump

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
	"go.etcd.io/bbolt"
)

func newTestClient() *utils.APIClient {
	c := api.NewClient("dummy")
	c.SetBaseURL(&url.URL{
		Scheme: "http",
		Host:   "testserver",
	})
	return &utils.APIClient{Client: c}
}

func newTestDatabase(t *testing.T) *utils.Database {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := bbolt.Open(dbPath, 0o600, nil)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("datadump"))
		return err
	})
	if err != nil {
		db.Close() // nolint:errcheck
		t.Fatalf("failed to create bucket: %v", err)
	}

	return &utils.Database{DB: db}
}

func TestDownload(t *testing.T) {
	defer gock.Off()

	t.Run("downloads file with default output name", func(t *testing.T) {
		defer gock.Clean()

		gock.New("http://testserver").
			Get("/api/v1/datadump/link/hours/api/20260101/20260101-01.gz").
			Reply(200).
			BodyString("test content")

		db := newTestDatabase(t)
		defer db.Close() // nolint:errcheck

		client := newTestClient()
		tmpDir := t.TempDir()

		err := download(client, db, "hours/api/20260101/20260101-01.gz", "", tmpDir, false, false)
		assert.NoError(t, err)

		// verify file was downloaded
		outputPath := filepath.Join(tmpDir, "20260101-01.gz")
		assert.FileExists(t, outputPath)

		content, err := os.ReadFile(outputPath)
		assert.NoError(t, err)
		assert.Equal(t, "test content", string(content))

		// verify database was updated
		downloaded, err := db.HasDataDumpBeenDownloaded("hours/api/20260101/20260101-01.gz")
		assert.NoError(t, err)
		assert.True(t, downloaded)

		assert.True(t, gock.IsDone())
	})

	t.Run("downloads file with custom output name", func(t *testing.T) {
		defer gock.Clean()

		gock.New("http://testserver").
			Get("/api/v1/datadump/link/hours/api/20260101/20260101-01.gz").
			Reply(200).
			BodyString("custom output content")

		db := newTestDatabase(t)
		defer db.Close() // nolint:errcheck

		client := newTestClient()
		tmpDir := t.TempDir()

		err := download(client, db, "hours/api/20260101/20260101-01.gz", "custom.gz", tmpDir, false, false)
		assert.NoError(t, err)

		// verify file was downloaded with custom name
		outputPath := filepath.Join(tmpDir, "custom.gz")
		assert.FileExists(t, outputPath)

		content, err := os.ReadFile(outputPath)
		assert.NoError(t, err)
		assert.Equal(t, "custom output content", string(content))

		assert.True(t, gock.IsDone())
	})

	t.Run("returns error when file exists and force is false", func(t *testing.T) {
		defer gock.Clean()

		db := newTestDatabase(t)
		defer db.Close() // nolint:errcheck

		client := newTestClient()
		tmpDir := t.TempDir()

		// create existing file
		existingFile := filepath.Join(tmpDir, "20260101-01.gz")
		err := os.WriteFile(existingFile, []byte("existing content"), 0o644)
		assert.NoError(t, err)

		err = download(client, db, "hours/api/20260101/20260101-01.gz", "", tmpDir, false, false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")

		// verify original content is unchanged
		content, err := os.ReadFile(existingFile)
		assert.NoError(t, err)
		assert.Equal(t, "existing content", string(content))
	})

	t.Run("overwrites file when force is true", func(t *testing.T) {
		defer gock.Clean()

		gock.New("http://testserver").
			Get("/api/v1/datadump/link/hours/api/20260101/20260101-01.gz").
			Reply(200).
			BodyString("new content")

		db := newTestDatabase(t)
		defer db.Close() // nolint:errcheck

		client := newTestClient()
		tmpDir := t.TempDir()

		// create existing file
		existingFile := filepath.Join(tmpDir, "20260101-01.gz")
		err := os.WriteFile(existingFile, []byte("old content"), 0o644)
		assert.NoError(t, err)

		err = download(client, db, "hours/api/20260101/20260101-01.gz", "", tmpDir, true, false)
		assert.NoError(t, err)

		// verify file was overwritten
		content, err := os.ReadFile(existingFile)
		assert.NoError(t, err)
		assert.Equal(t, "new content", string(content))

		assert.True(t, gock.IsDone())
	})
}

func TestDownloadWithFollow(t *testing.T) {
	defer gock.Off()

	t.Run("downloads multiple files with follow option", func(t *testing.T) {
		defer gock.Clean()

		// mock the list endpoint
		gock.New("http://testserver").
			Get("/api/v1/datadump/list/hours/api/20260101/").
			Reply(200).
			JSON(map[string]any{
				"files": []map[string]any{
					{"path": "hours/api/20260101/20260101-01.gz", "size": 100},
					{"path": "hours/api/20260101/20260101-02.gz", "size": 200},
				},
			})

		// mock the download endpoints
		gock.New("http://testserver").
			Get("/api/v1/datadump/link/hours/api/20260101/20260101-01.gz").
			Reply(200).
			BodyString("content 01")

		gock.New("http://testserver").
			Get("/api/v1/datadump/link/hours/api/20260101/20260101-02.gz").
			Reply(200).
			BodyString("content 02")

		db := newTestDatabase(t)
		defer db.Close() // nolint:errcheck

		client := newTestClient()
		tmpDir := t.TempDir()

		// find missing paths
		paths, err := findMissingPaths(db, client, "hours/api/20260101/", false)
		assert.NoError(t, err)
		assert.Len(t, paths, 2)

		// download each file
		for _, path := range paths {
			err := download(client, db, path, "", tmpDir, false, false)
			assert.NoError(t, err)
		}

		// verify files were downloaded
		content1, err := os.ReadFile(filepath.Join(tmpDir, "20260101-01.gz"))
		assert.NoError(t, err)
		assert.Equal(t, "content 01", string(content1))

		content2, err := os.ReadFile(filepath.Join(tmpDir, "20260101-02.gz"))
		assert.NoError(t, err)
		assert.Equal(t, "content 02", string(content2))

		// verify database was updated
		downloaded1, err := db.HasDataDumpBeenDownloaded("hours/api/20260101/20260101-01.gz")
		assert.NoError(t, err)
		assert.True(t, downloaded1)

		downloaded2, err := db.HasDataDumpBeenDownloaded("hours/api/20260101/20260101-02.gz")
		assert.NoError(t, err)
		assert.True(t, downloaded2)

		assert.True(t, gock.IsDone())
	})

	t.Run("skips already downloaded files with follow option", func(t *testing.T) {
		defer gock.Clean()

		// mock the list endpoint
		gock.New("http://testserver").
			Get("/api/v1/datadump/list/hours/api/20260101/").
			Reply(200).
			JSON(map[string]any{
				"files": []map[string]any{
					{"path": "hours/api/20260101/20260101-01.gz", "size": 100},
					{"path": "hours/api/20260101/20260101-02.gz", "size": 200},
				},
			})

		// only mock download for the second file (first is already downloaded)
		gock.New("http://testserver").
			Get("/api/v1/datadump/link/hours/api/20260101/20260101-02.gz").
			Reply(200).
			BodyString("content 02")

		db := newTestDatabase(t)
		defer db.Close() // nolint:errcheck

		client := newTestClient()
		tmpDir := t.TempDir()

		// simulate first file already downloaded
		existingFile := filepath.Join(tmpDir, "20260101-01.gz")
		err := os.WriteFile(existingFile, []byte("existing content"), 0o644)
		assert.NoError(t, err)
		err = db.SetDataDump("hours/api/20260101/20260101-01.gz", existingFile)
		assert.NoError(t, err)

		// find missing paths (should only return second file)
		paths, err := findMissingPaths(db, client, "hours/api/20260101/", false)
		assert.NoError(t, err)
		assert.Len(t, paths, 1)
		assert.Equal(t, "hours/api/20260101/20260101-02.gz", paths[0])

		// download missing files
		for _, path := range paths {
			err := download(client, db, path, "", tmpDir, false, false)
			assert.NoError(t, err)
		}

		// verify first file was not changed
		content1, err := os.ReadFile(existingFile)
		assert.NoError(t, err)
		assert.Equal(t, "existing content", string(content1))

		// verify second file was downloaded
		content2, err := os.ReadFile(filepath.Join(tmpDir, "20260101-02.gz"))
		assert.NoError(t, err)
		assert.Equal(t, "content 02", string(content2))

		assert.True(t, gock.IsDone())
	})

	t.Run("re-downloads all files with follow and force options", func(t *testing.T) {
		defer gock.Clean()

		// mock the list endpoint
		gock.New("http://testserver").
			Get("/api/v1/datadump/list/hours/api/20260101/").
			Reply(200).
			JSON(map[string]any{
				"files": []map[string]any{
					{"path": "hours/api/20260101/20260101-01.gz", "size": 100},
					{"path": "hours/api/20260101/20260101-02.gz", "size": 200},
				},
			})

		// mock download for both files
		gock.New("http://testserver").
			Get("/api/v1/datadump/link/hours/api/20260101/20260101-01.gz").
			Reply(200).
			BodyString("new content 01")

		gock.New("http://testserver").
			Get("/api/v1/datadump/link/hours/api/20260101/20260101-02.gz").
			Reply(200).
			BodyString("new content 02")

		db := newTestDatabase(t)
		defer db.Close() // nolint:errcheck

		client := newTestClient()
		tmpDir := t.TempDir()

		// simulate first file already downloaded
		existingFile := filepath.Join(tmpDir, "20260101-01.gz")
		err := os.WriteFile(existingFile, []byte("old content"), 0o644)
		assert.NoError(t, err)
		err = db.SetDataDump("hours/api/20260101/20260101-01.gz", existingFile)
		assert.NoError(t, err)

		// find paths with force=true (should return both files)
		paths, err := findMissingPaths(db, client, "hours/api/20260101/", true)
		assert.NoError(t, err)
		assert.Len(t, paths, 2)

		// download all files with force
		for _, path := range paths {
			err := download(client, db, path, "", tmpDir, true, false)
			assert.NoError(t, err)
		}

		// verify first file was overwritten
		content1, err := os.ReadFile(existingFile)
		assert.NoError(t, err)
		assert.Equal(t, "new content 01", string(content1))

		// verify second file was downloaded
		content2, err := os.ReadFile(filepath.Join(tmpDir, "20260101-02.gz"))
		assert.NoError(t, err)
		assert.Equal(t, "new content 02", string(content2))

		assert.True(t, gock.IsDone())
	})

	t.Run("returns no paths when all files are downloaded", func(t *testing.T) {
		defer gock.Clean()

		// mock the list endpoint
		gock.New("http://testserver").
			Get("/api/v1/datadump/list/hours/api/20260101/").
			Reply(200).
			JSON(map[string]any{
				"files": []map[string]any{
					{"path": "hours/api/20260101/20260101-01.gz", "size": 100},
				},
			})

		db := newTestDatabase(t)
		defer db.Close() // nolint:errcheck

		client := newTestClient()
		tmpDir := t.TempDir()

		// simulate file already downloaded
		existingFile := filepath.Join(tmpDir, "20260101-01.gz")
		err := os.WriteFile(existingFile, []byte("content"), 0o644)
		assert.NoError(t, err)
		err = db.SetDataDump("hours/api/20260101/20260101-01.gz", existingFile)
		assert.NoError(t, err)

		// find missing paths (should return empty)
		paths, err := findMissingPaths(db, client, "hours/api/20260101/", false)
		assert.NoError(t, err)
		assert.Len(t, paths, 0)

		assert.True(t, gock.IsDone())
	})
}

func TestFindMissingPaths(t *testing.T) {
	defer gock.Off()

	t.Run("returns all paths when database is empty", func(t *testing.T) {
		defer gock.Clean()

		gock.New("http://testserver").
			Get("/api/v1/datadump/list/hours/api/20260101/").
			Reply(200).
			JSON(map[string]any{
				"files": []map[string]any{
					{"path": "hours/api/20260101/20260101-01.gz", "size": 100},
					{"path": "hours/api/20260101/20260101-02.gz", "size": 200},
					{"path": "hours/api/20260101/20260101-03.gz", "size": 300},
				},
			})

		db := newTestDatabase(t)
		defer db.Close() // nolint:errcheck

		client := newTestClient()

		paths, err := findMissingPaths(db, client, "hours/api/20260101/", false)
		assert.NoError(t, err)
		assert.Len(t, paths, 3)
		assert.Contains(t, paths, "hours/api/20260101/20260101-01.gz")
		assert.Contains(t, paths, "hours/api/20260101/20260101-02.gz")
		assert.Contains(t, paths, "hours/api/20260101/20260101-03.gz")

		assert.True(t, gock.IsDone())
	})

	t.Run("excludes already downloaded paths", func(t *testing.T) {
		defer gock.Clean()

		gock.New("http://testserver").
			Get("/api/v1/datadump/list/hours/api/20260101/").
			Reply(200).
			JSON(map[string]any{
				"files": []map[string]any{
					{"path": "hours/api/20260101/20260101-01.gz", "size": 100},
					{"path": "hours/api/20260101/20260101-02.gz", "size": 200},
					{"path": "hours/api/20260101/20260101-03.gz", "size": 300},
				},
			})

		db := newTestDatabase(t)
		defer db.Close() // nolint:errcheck

		// create a temporary file to simulate a downloaded file
		tmpDir := t.TempDir()
		downloadedFile := filepath.Join(tmpDir, "20260101-02.gz")
		err := os.WriteFile(downloadedFile, []byte("test"), 0o644)
		assert.NoError(t, err)

		// mark the second path as downloaded with the actual file path
		err = db.SetDataDump("hours/api/20260101/20260101-02.gz", downloadedFile)
		assert.NoError(t, err)

		client := newTestClient()

		paths, err := findMissingPaths(db, client, "hours/api/20260101/", false)
		assert.NoError(t, err)
		assert.Len(t, paths, 2)
		assert.Contains(t, paths, "hours/api/20260101/20260101-01.gz")
		assert.Contains(t, paths, "hours/api/20260101/20260101-03.gz")
		assert.NotContains(t, paths, "hours/api/20260101/20260101-02.gz")

		assert.True(t, gock.IsDone())
	})

	t.Run("returns all paths when force is true", func(t *testing.T) {
		defer gock.Clean()

		gock.New("http://testserver").
			Get("/api/v1/datadump/list/hours/api/20260101/").
			Reply(200).
			JSON(map[string]any{
				"files": []map[string]any{
					{"path": "hours/api/20260101/20260101-01.gz", "size": 100},
					{"path": "hours/api/20260101/20260101-02.gz", "size": 200},
				},
			})

		db := newTestDatabase(t)
		defer db.Close() // nolint:errcheck

		tmpDir := t.TempDir()
		downloadedFile := filepath.Join(tmpDir, "20260101-01.gz")
		err := os.WriteFile(downloadedFile, []byte("test"), 0o644)
		assert.NoError(t, err)

		err = db.SetDataDump("hours/api/20260101/20260101-01.gz", downloadedFile)
		assert.NoError(t, err)

		client := newTestClient()

		paths, err := findMissingPaths(db, client, "hours/api/20260101/", true)
		assert.NoError(t, err)
		assert.Len(t, paths, 2)
		assert.Contains(t, paths, "hours/api/20260101/20260101-01.gz")
		assert.Contains(t, paths, "hours/api/20260101/20260101-02.gz")

		assert.True(t, gock.IsDone())
	})

	t.Run("returns empty slice when all files are downloaded", func(t *testing.T) {
		defer gock.Clean()

		gock.New("http://testserver").
			Get("/api/v1/datadump/list/hours/api/20260101/").
			Reply(200).
			JSON(map[string]any{
				"files": []map[string]any{
					{"path": "hours/api/20260101/20260101-01.gz", "size": 100},
				},
			})

		db := newTestDatabase(t)
		defer db.Close() // nolint:errcheck

		tmpDir := t.TempDir()
		downloadedFile := filepath.Join(tmpDir, "20260101-01.gz")
		err := os.WriteFile(downloadedFile, []byte("test"), 0o644)
		assert.NoError(t, err)

		err = db.SetDataDump("hours/api/20260101/20260101-01.gz", downloadedFile)
		assert.NoError(t, err)

		client := newTestClient()

		paths, err := findMissingPaths(db, client, "hours/api/20260101/", false)
		assert.NoError(t, err)
		assert.Len(t, paths, 0)

		assert.True(t, gock.IsDone())
	})

	t.Run("enumerates last 7 days when path has no date", func(t *testing.T) {
		defer gock.Clean()

		today := time.Now().UTC()
		days := api.GetLast7Days(today)

		for _, day := range days {
			gock.New("http://testserver").
				Get("/api/v1/datadump/list/hours/api/" + day + "/").
				Reply(200).
				JSON(map[string]any{
					"files": []map[string]any{
						{"path": "hours/api/" + day + "/" + day + ".gz", "size": 100},
					},
				})
		}

		db := newTestDatabase(t)
		defer db.Close() // nolint:errcheck

		client := newTestClient()

		paths, err := findMissingPaths(db, client, "hours/api/", false)
		assert.NoError(t, err)
		assert.Len(t, paths, 7)

		assert.True(t, gock.IsDone())
	})
}
