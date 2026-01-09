package utils

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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

	// check if gzipped
	isGzipped, err := isGzip(file)
	if err != nil {
		return fmt.Errorf("failed to check file format: %w", err)
	}
	// reset file pointer
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	// create reader chain based on compression
	var reader io.Reader = file
	var gzReader *gzip.Reader
	if isGzipped {
		gzReader, err = gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer func() {
			closeErr := gzReader.Close()
			if closeErr != nil && err == nil {
				err = closeErr
			}
		}()
		reader = gzReader
	}

	// use buffered reader to peek at content
	bufReader := bufio.NewReader(reader)
	isTared, err := isTarFromBuffered(bufReader)
	if err != nil {
		return fmt.Errorf("failed to check tar format: %w", err)
	}

	if isTared {
		outputDir := filepath.Dir(path)
		err = extractTar(bufReader, outputDir)
		if err != nil {
			return fmt.Errorf("failed to extract tar: %w", err)
		}
		return nil
	}

	if isGzipped {
		outputPath := strings.TrimSuffix(path, ".gz")
		outFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer func() {
			closeErr := outFile.Close()
			if closeErr != nil && err == nil {
				err = closeErr
			}
		}()

		_, err = io.Copy(outFile, bufReader)
		if err != nil {
			return fmt.Errorf("failed to write unpacked file: %w", err)
		}
		return nil
	}

	return fmt.Errorf("unsupported file format for unpacking")
}

func isGzip(reader io.Reader) (bool, error) {
	buf := make([]byte, 2)
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}
	if n < 2 {
		return false, nil
	}

	// check for gzip magic number (0x1f 0x8b)
	return buf[0] == 0x1f && buf[1] == 0x8b, nil
}

func isTarFromBuffered(reader *bufio.Reader) (bool, error) {
	// Peek at first 512 bytes without consuming them
	buf, err := reader.Peek(512)
	if err != nil && err != io.EOF {
		return false, err
	}
	if len(buf) < 512 {
		return false, nil
	}

	// check for tar magic number ("ustar" at offset 257)
	return string(buf[257:262]) == "ustar", nil
}

func extractTar(reader io.Reader, outputDir string) error {
	tarReader := tar.NewReader(reader)

	// Clean and resolve the output directory to an absolute path
	cleanOutputDir, err := filepath.Abs(filepath.Clean(outputDir))
	if err != nil {
		return fmt.Errorf("failed to resolve output directory: %w", err)
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		target := filepath.Join(cleanOutputDir, header.Name)

		// validate path doesn't escape outputDir (path traversal protection for just in case)
		cleanTarget := filepath.Clean(target)
		if !strings.HasPrefix(cleanTarget, cleanOutputDir+string(filepath.Separator)) && cleanTarget != cleanOutputDir {
			return fmt.Errorf("illegal file path: %s (attempted path traversal)", header.Name)
		}

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
