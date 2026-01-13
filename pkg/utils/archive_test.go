package utils

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExtractGzip(t *testing.T) {
	// setup fixture
	tempDir := t.TempDir()

	testContent := []byte("This is test content for gzip extracting")
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

	// test extracting
	err = Extract(gzFilePath, NewExtractOptions())
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	// verify the extracted file exists and has correct content
	extractedPath := filepath.Join(tempDir, "test")
	content, err := os.ReadFile(extractedPath)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}

	if string(content) != string(testContent) {
		t.Errorf("Content mismatch: got %q, want %q", string(content), string(testContent))
	}
}

func TestExtractTar(t *testing.T) {
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

	// test extracting
	err = Extract(tarFilePath, NewExtractOptions())
	if err != nil {
		t.Fatalf("Extract failed for tar: %v", err)
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

func TestExtractTarGzip(t *testing.T) {
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

	// test extracting
	err = Extract(tarGzFilePath, NewExtractOptions())
	if err != nil {
		t.Fatalf("Extract failed for tar.gz: %v", err)
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

func TestExtractTarWithZeroBlocks(t *testing.T) {
	// setup fixture - create a tar file with zero blocks inserted
	tempDir := t.TempDir()

	testFiles := map[string][]byte{
		"file1.txt": []byte("Content of file 1"),
		"file2.txt": []byte("Content of file 2"),
	}

	tarFilePath := filepath.Join(tempDir, "test_with_zeros.tar")

	tarFile, err := os.Create(tarFilePath)
	if err != nil {
		t.Fatalf("Failed to create test tar file: %v", err)
	}

	tarWriter := tar.NewWriter(tarFile)

	// write first file
	header1 := &tar.Header{
		Name:       "file1.txt",
		Mode:       0o644,
		Size:       int64(len(testFiles["file1.txt"])),
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
	if err := tarWriter.WriteHeader(header1); err != nil {
		t.Fatalf("Failed to write file1 header: %v", err)
	}
	if _, err := tarWriter.Write(testFiles["file1.txt"]); err != nil {
		t.Fatalf("Failed to write file1 content: %v", err)
	}

	// close tar writer to flush buffers
	if err := tarWriter.Close(); err != nil {
		t.Fatalf("Failed to close tar writer: %v", err)
	}

	// manually append zero blocks (this simulates corrupted or padded tar files)
	zeroBlock := make([]byte, tarBlockSize)
	for range 3 {
		_, err := tarFile.Write(zeroBlock)
		if err != nil {
			t.Fatalf("Failed to write zero block: %v", err)
		}
	}

	// now write second file by creating a new tar writer
	tarWriter2 := tar.NewWriter(tarFile)
	header2 := &tar.Header{
		Name:       "file2.txt",
		Mode:       0o644,
		Size:       int64(len(testFiles["file2.txt"])),
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
	err = tarWriter2.WriteHeader(header2)
	if err != nil {
		t.Fatalf("Failed to write file2 header: %v", err)
	}
	_, err = tarWriter2.Write(testFiles["file2.txt"])
	if err != nil {
		t.Fatalf("Failed to write file2 content: %v", err)
	}

	err = tarWriter2.Close()
	if err != nil {
		t.Fatalf("Failed to close second tar writer: %v", err)
	}

	err = tarFile.Close()
	if err != nil {
		t.Fatalf("Failed to close tar file: %v", err)
	}

	// test extracting - should succeed with zero block skipping
	err = Extract(tarFilePath, NewExtractOptions())
	if err != nil {
		t.Fatalf("Extract failed for tar with zero blocks: %v", err)
	}

	// verify both files were extracted with correct content
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
}

func TestExtractWithForceOption(t *testing.T) {
	// test that extraction fails when file exists and force is false
	tempDir := t.TempDir()

	testContent := []byte("Original content")
	gzFilePath := filepath.Join(tempDir, "test.gz")
	extractedPath := filepath.Join(tempDir, "test")

	// create gzip file
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

	// create existing file
	err = os.WriteFile(extractedPath, []byte("Existing content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	// test extraction with force=false should fail
	err = Extract(gzFilePath, NewExtractOptions(WithExtractForce(false)))
	if err == nil {
		t.Fatal("Expected error when extracting with existing file and force=false, got nil")
	}

	// verify the existing file was not overwritten
	content, err := os.ReadFile(extractedPath)
	if err != nil {
		t.Fatalf("Failed to read existing file: %v", err)
	}
	if string(content) != "Existing content" {
		t.Errorf("Existing file was modified: got %q, want %q", string(content), "Existing content")
	}

	// test extraction with force=true should succeed
	err = Extract(gzFilePath, NewExtractOptions(WithExtractForce(true)))
	if err != nil {
		t.Fatalf("Extract failed with force=true: %v", err)
	}

	// verify the file was overwritten with new content
	content, err = os.ReadFile(extractedPath)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}
	if string(content) != string(testContent) {
		t.Errorf("Content mismatch after forced extraction: got %q, want %q", string(content), string(testContent))
	}
}

func TestExtractTarGzipWithZeroBlocks(t *testing.T) {
	// setup fixture - create a tar.gz file with zero blocks inserted in the tar stream
	tempDir := t.TempDir()

	testFiles := map[string][]byte{
		"file1.txt": []byte("Content of file 1"),
		"file2.txt": []byte("Content of file 2"),
	}

	tarGzFilePath := filepath.Join(tempDir, "test_with_zeros.tar.gz")

	tarGzFile, err := os.Create(tarGzFilePath)
	if err != nil {
		t.Fatalf("Failed to create test tar.gz file: %v", err)
	}

	gzWriter := gzip.NewWriter(tarGzFile)
	tarWriter := tar.NewWriter(gzWriter)

	// write first file
	header1 := &tar.Header{
		Name:       "file1.txt",
		Mode:       0o644,
		Size:       int64(len(testFiles["file1.txt"])),
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
	if err := tarWriter.WriteHeader(header1); err != nil {
		t.Fatalf("Failed to write file1 header: %v", err)
	}
	if _, err := tarWriter.Write(testFiles["file1.txt"]); err != nil {
		t.Fatalf("Failed to write file1 content: %v", err)
	}

	// flush and close tar writer
	if err := tarWriter.Close(); err != nil {
		t.Fatalf("Failed to close tar writer: %v", err)
	}

	// manually append zero blocks through the gzip writer
	zeroBlock := make([]byte, tarBlockSize)
	for range 3 {
		if _, err := gzWriter.Write(zeroBlock); err != nil {
			t.Fatalf("Failed to write zero block: %v", err)
		}
	}

	// write second file with a new tar writer
	tarWriter2 := tar.NewWriter(gzWriter)
	header2 := &tar.Header{
		Name:       "file2.txt",
		Mode:       0o644,
		Size:       int64(len(testFiles["file2.txt"])),
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
	if err := tarWriter2.WriteHeader(header2); err != nil {
		t.Fatalf("Failed to write file2 header: %v", err)
	}
	if _, err := tarWriter2.Write(testFiles["file2.txt"]); err != nil {
		t.Fatalf("Failed to write file2 content: %v", err)
	}

	if err := tarWriter2.Close(); err != nil {
		t.Fatalf("Failed to close second tar writer: %v", err)
	}

	if err := gzWriter.Close(); err != nil {
		t.Fatalf("Failed to close gzip writer: %v", err)
	}

	if err := tarGzFile.Close(); err != nil {
		t.Fatalf("Failed to close tar.gz file: %v", err)
	}

	// test extracting - should succeed with zero block skipping
	err = Extract(tarGzFilePath, NewExtractOptions())
	if err != nil {
		t.Fatalf("Extract failed for tar.gz with zero blocks: %v", err)
	}

	// verify both files were extracted with correct content
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
}
