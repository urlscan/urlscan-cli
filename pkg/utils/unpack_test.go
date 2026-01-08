package utils

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestUnpackGzip(t *testing.T) {
	// setup fixture
	tempDir := t.TempDir()

	testContent := []byte("This is test content for gzip unpacking")
	gzFilePath := filepath.Join(tempDir, "test.gz")

	gzFile, err := os.Create(gzFilePath)
	if err != nil {
		t.Fatalf("Failed to create test gzip file: %v", err)
	}

	gzWriter := gzip.NewWriter(gzFile)
	_, err = gzWriter.Write(testContent)
	if err != nil {
		t.Fatalf("Failed to write test content: %v", err)
	}
	err = gzWriter.Close()
	if err != nil {
		t.Fatalf("Failed to close gzip writer: %v", err)
	}
	err = gzFile.Close()
	if err != nil {
		t.Fatalf("Failed to close gzip file: %v", err)
	}

	// test unpacking
	err = Unpack(gzFilePath)
	if err != nil {
		t.Fatalf("Unpack failed: %v", err)
	}

	// verify the unpacked file exists and has correct content
	unpackedPath := filepath.Join(tempDir, "test")
	content, err := os.ReadFile(unpackedPath)
	if err != nil {
		t.Fatalf("Failed to read unpacked file: %v", err)
	}

	if string(content) != string(testContent) {
		t.Errorf("Content mismatch: got %q, want %q", string(content), string(testContent))
	}
}

func TestUnpackTar(t *testing.T) {
	// setup fixture
	tempDir := t.TempDir()

	testFiles := map[string][]byte{
		"file1.txt":     []byte("Content of file 1"),
		"file2.txt":     []byte("Content of file 2"),
		"dir/file3.txt": []byte("Content of file 3 in subdirectory"),
	}

	tarFilePath := filepath.Join(tempDir, "test.tar")

	tarFile, err := os.Create(tarFilePath)
	if err != nil {
		t.Fatalf("Failed to create test tar file: %v", err)
	}

	tarWriter := tar.NewWriter(tarFile)

	for filename, content := range testFiles {
		if dir := filepath.Dir(filename); dir != "." {
			header := &tar.Header{
				Name:       dir + "/",
				Mode:       0o755,
				Typeflag:   tar.TypeDir,
				Format:     tar.FormatPAX,
				Linkname:   "",
				Size:       0,
				Uid:        0,
				Gid:        0,
				Uname:      "",
				Gname:      "",
				ModTime:    time.Time{},
				AccessTime: time.Time{},
				ChangeTime: time.Time{},
				Devmajor:   0,
				Devminor:   0,
				PAXRecords: nil,
				Xattrs:     nil,
			}
			err := tarWriter.WriteHeader(header)
			if err != nil {
				t.Fatalf("Failed to write dir header: %v", err)
			}
		}

		header := &tar.Header{
			Name:       filename,
			Mode:       0o644,
			Size:       int64(len(content)),
			Typeflag:   tar.TypeReg,
			Format:     tar.FormatPAX,
			Linkname:   "",
			Uid:        0,
			Gid:        0,
			Uname:      "",
			Gname:      "",
			ModTime:    time.Time{},
			AccessTime: time.Time{},
			ChangeTime: time.Time{},
			Devmajor:   0,
			Devminor:   0,
			PAXRecords: nil,
			Xattrs:     nil,
		}
		err = tarWriter.WriteHeader(header)
		if err != nil {
			t.Fatalf("Failed to write file header: %v", err)
		}
		_, err = tarWriter.Write(content)
		if err != nil {
			t.Fatalf("Failed to write file content: %v", err)
		}
	}

	err = tarWriter.Close()
	if err != nil {
		t.Fatalf("Failed to close tar writer: %v", err)
	}
	err = tarFile.Close()
	if err != nil {
		t.Fatalf("Failed to close tar file: %v", err)
	}

	// test unpacking
	err = Unpack(tarFilePath)
	if err != nil {
		t.Fatalf("Unpack failed for tar: %v", err)
	}

	// verify all files were extracted with correct content
	for filename, expectedContent := range testFiles {
		extractedPath := filepath.Join(tempDir, filename)
		content, err := os.ReadFile(extractedPath)
		if err != nil {
			t.Errorf("Failed to read extracted file %s: %v", filename, err)
			continue
		}

		if string(content) != string(expectedContent) {
			t.Errorf("Content mismatch for %s: got %q, want %q", filename, string(content), string(expectedContent))
		}
	}

	// verify directory was created
	dirPath := filepath.Join(tempDir, "dir")
	info, err := os.Stat(dirPath)
	if err != nil {
		t.Errorf("Directory 'dir' was not created: %v", err)
	} else if !info.IsDir() {
		t.Error("'dir' is not a directory")
	}
}

func TestUnpackTarGzip(t *testing.T) {
	// setup fixture
	tempDir := t.TempDir()

	testFiles := map[string][]byte{
		"file1.txt":     []byte("Content of file 1"),
		"file2.txt":     []byte("Content of file 2"),
		"dir/file3.txt": []byte("Content of file 3 in subdirectory"),
	}

	tarGzFilePath := filepath.Join(tempDir, "test.tar.gz")

	tarGzFile, err := os.Create(tarGzFilePath)
	if err != nil {
		t.Fatalf("Failed to create test tar.gz file: %v", err)
	}

	gzWriter := gzip.NewWriter(tarGzFile)
	tarWriter := tar.NewWriter(gzWriter)

	for filename, content := range testFiles {
		if dir := filepath.Dir(filename); dir != "." {
			header := &tar.Header{
				Name:       dir + "/",
				Mode:       0o755,
				Typeflag:   tar.TypeDir,
				Format:     tar.FormatPAX,
				Linkname:   "",
				Size:       0,
				Uid:        0,
				Gid:        0,
				Uname:      "",
				Gname:      "",
				ModTime:    time.Time{},
				AccessTime: time.Time{},
				ChangeTime: time.Time{},
				Devmajor:   0,
				Devminor:   0,
				PAXRecords: nil,
				Xattrs:     nil,
			}
			if err := tarWriter.WriteHeader(header); err != nil {
				t.Fatalf("Failed to write dir header: %v", err)
			}
		}

		header := &tar.Header{
			Name:       filename,
			Mode:       0o644,
			Size:       int64(len(content)),
			Typeflag:   tar.TypeReg,
			Format:     tar.FormatPAX,
			Linkname:   "",
			Uid:        0,
			Gid:        0,
			Uname:      "",
			Gname:      "",
			ModTime:    time.Time{},
			AccessTime: time.Time{},
			ChangeTime: time.Time{},
			Devmajor:   0,
			Devminor:   0,
			PAXRecords: nil,
			Xattrs:     nil,
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			t.Fatalf("Failed to write file header: %v", err)
		}
		if _, err := tarWriter.Write(content); err != nil {
			t.Fatalf("Failed to write file content: %v", err)
		}
	}

	err = tarWriter.Close()
	if err != nil {
		t.Fatalf("Failed to close tar writer: %v", err)
	}
	err = gzWriter.Close()
	if err != nil {
		t.Fatalf("Failed to close gzip writer: %v", err)
	}
	err = tarGzFile.Close()
	if err != nil {
		t.Fatalf("Failed to close tar.gz file: %v", err)
	}

	// test unpacking
	err = Unpack(tarGzFilePath)
	if err != nil {
		t.Fatalf("Unpack failed for tar.gz: %v", err)
	}

	// verify all files were extracted with correct content
	for filename, expectedContent := range testFiles {
		extractedPath := filepath.Join(tempDir, filename)
		content, err := os.ReadFile(extractedPath)
		if err != nil {
			t.Errorf("Failed to read extracted file %s: %v", filename, err)
			continue
		}

		if string(content) != string(expectedContent) {
			t.Errorf("Content mismatch for %s: got %q, want %q", filename, string(content), string(expectedContent))
		}
	}

	// verify directory was created
	dirPath := filepath.Join(tempDir, "dir")
	info, err := os.Stat(dirPath)
	if err != nil {
		t.Errorf("Directory 'dir' was not created: %v", err)
	} else if !info.IsDir() {
		t.Error("'dir' is not a directory")
	}
}
