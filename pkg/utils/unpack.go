package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// NOTE:
// - Content-Type: application/gzip + .tar.gz => tar
// - Content-Type: binary/octet-stream + .gz => gzip
// - Content-Type: binary/octet-stream + .tar.gz => tar + gzip

func Unpack(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		closeErr := file.Close()

		if closeErr != nil {
			err = closeErr
		}
	}()

	isTarGzipped, err := isTarGzipFile(path)
	if err != nil {
		return fmt.Errorf("failed to check file format: %w", err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// unpack as tar.gz
	if isTarGzipped {
		outputDir := filepath.Dir(path)
		err = unpackTarGzip(file, outputDir)
		if err != nil {
			return fmt.Errorf("failed to extract tar.gz: %w", err)
		}
		return nil
	}

	isGzipped, err := isGzipFile(file)
	if err != nil {
		return fmt.Errorf("failed to check file format: %w", err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// unpack az gz
	if isGzipped {
		outputPath := strings.TrimSuffix(path, ".gz")
		err = unpackGzip(file, outputPath)
		if err != nil {
			return fmt.Errorf("failed to extract gz: %w", err)
		}
		return nil
	}

	// try to unpack as tar
	isTar, err := isTarFile(file)
	if err != nil {
		return fmt.Errorf("failed to check file format: %w", err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}

	if isTar {
		outputDir := filepath.Dir(path)
		err = unpackTar(file, outputDir)
		if err != nil {
			return fmt.Errorf("failed to extract tar: %w", err)
		}
		return nil
	}

	return fmt.Errorf("unsupported file format for unpacking")
}

func isGzipFile(file *os.File) (bool, error) {
	// read the first 2 bytes (gzip magic number)
	header := make([]byte, 2)
	n, err := file.Read(header)
	if err != nil && err != io.EOF {
		return false, err
	}
	if n < 2 {
		return false, nil
	}

	// check for gzip magic number (0x1f 0x8b)
	return header[0] == 0x1f && header[1] == 0x8b, nil
}

func isTarFile(file *os.File) (bool, error) {
	// read the first 512 bytes (tar header size)
	header := make([]byte, 512)
	n, err := file.Read(header)

	if err != nil && err != io.EOF {
		return false, err
	}
	if n < 512 {
		return false, nil
	}

	// check for tar magic number at offset 257
	// "ustar\x00" (POSIX tar) or "ustar  \x00" (GNU tar)
	magic := string(header[257:262])
	return magic == "ustar", nil
}

func isTarGzipFile(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		closeErr := file.Close()

		if closeErr != nil {
			err = closeErr
		}
	}()

	// first check if it's gzipped
	isGzipped, err := isGzipFile(file)
	if err != nil {
		return false, err
	}
	if !isGzipped {
		return false, nil
	}

	// reset file pointer
	_, err = file.Seek(0, 0)
	if err != nil {
		return false, err
	}

	// try to read the gzipped content and check if it's a tar archive
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return false, err
	}
	defer func() {
		closeErr := gzReader.Close()

		if closeErr != nil {
			err = closeErr
		}
	}()

	// read the first 512 bytes from the decompressed content
	header := make([]byte, 512)
	n, err := gzReader.Read(header)
	if err != nil && err != io.EOF {
		return false, err
	}
	if n < 512 {
		return false, nil
	}

	// check for tar magic number at offset 257
	magic := string(header[257:262])
	return magic == "ustar", nil
}

func unpackGzip(file *os.File, outputPath string) error {
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer func() {
		closeErr := gzReader.Close()

		if closeErr != nil {
			err = closeErr
		}
	}()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() {
		closeErr := outFile.Close()

		if closeErr != nil {
			err = closeErr
		}
	}()

	_, err = io.Copy(outFile, gzReader)
	if err != nil {
		return fmt.Errorf("failed to decompress file: %w", err)
	}

	return nil
}

func extractTar(tarReader *tar.Reader, outputDir string) error {
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		target := filepath.Join(outputDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			err := os.MkdirAll(target, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %w", target, err)
			}
		case tar.TypeReg:
			// ensure parent directory exists
			err := os.MkdirAll(filepath.Dir(target), 0o755)
			if err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", target, err)
			}

			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", target, err)
			}

			_, err = io.Copy(outFile, tarReader)
			closeErr := outFile.Close()
			if closeErr != nil && err == nil {
				err = closeErr
			}
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", target, err)
			}
		default:
			return fmt.Errorf("unsupported tar entry type: %c in file %s", header.Typeflag, header.Name)
		}
	}

	return nil
}

func unpackTar(file *os.File, outputDir string) error {
	tarReader := tar.NewReader(file)
	return extractTar(tarReader, outputDir)
}

func unpackTarGzip(file *os.File, outputDir string) error {
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer func() {
		closeErr := gzReader.Close()

		if closeErr != nil {
			err = closeErr
		}
	}()

	tarReader := tar.NewReader(gzReader)
	return extractTar(tarReader, outputDir)
}
