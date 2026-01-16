package datadump

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

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
}
