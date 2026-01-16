package utils

import (
	"path/filepath"

	"github.com/adrg/xdg"
	"go.etcd.io/bbolt"
)

const (
	namespace          = "urlscan"
	databaseFilename   = "state.db"
	dataDumpBucketName = "datadump"
)

type Database struct {
	db *bbolt.DB
}

func getDatabaseFile() (string, error) {
	return xdg.DataFile(filepath.Join(namespace, databaseFilename))
}

func NewDatabase() (d *Database, err error) {
	path, err := getDatabaseFile()
	if err != nil {
		return nil, err
	}

	db, err := bbolt.Open(path, 0o600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(dataDumpBucketName))
		return err
	})
	if err != nil {
		defer func() {
			closeErr := db.Close()
			if closeErr != nil && err == nil {
				err = closeErr
			}
		}()
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetDataDump(path string) (localPath string, exists bool, err error) {
	err = d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(dataDumpBucketName))
		v := b.Get([]byte(path))
		if v != nil {
			localPath = string(v)
			exists = true
		}
		return nil
	})

	return localPath, exists, err
}

func (d *Database) SetDataDump(path, localPath string) error {
	return d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(dataDumpBucketName))
		return b.Put([]byte(path), []byte(localPath))
	})
}

func (d *Database) DeleteDataDump(path string) error {
	return d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(dataDumpBucketName))
		return b.Delete([]byte(path))
	})
}

func (d *Database) IsDataDumpDownloaded(path string) (bool, error) {
	localPath, exists, err := d.GetDataDump(path)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}

	if !fileExists(localPath) {
		// if file is deleted, remove it from the database
		err := d.DeleteDataDump(path)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}
