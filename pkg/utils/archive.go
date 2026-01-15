package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ExtractOptions struct {
	force bool
}

type ExtractOption func(*ExtractOptions)

func WithExtractForce(force bool) ExtractOption {
	return func(opts *ExtractOptions) {
		opts.force = force
	}
}

func NewExtractOptions(opts ...ExtractOption) *ExtractOptions {
	var o ExtractOptions
	for _, fn := range opts {
		fn(&o)
	}
	return &o
}

func Extract(path string, opts *ExtractOptions) error {
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

	// check if gzipped and/or tared
	isGzipped := strings.HasSuffix(path, ".gz")
	isTared := strings.HasSuffix(path, ".tar") || strings.HasSuffix(path, ".tar.gz")

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

	if isTared {
		outputDir := filepath.Dir(path)
		err = extractTar(reader, outputDir, opts)
		if err != nil {
			return fmt.Errorf("failed to extract tar: %w", err)
		}
		return nil
	}

	if isGzipped {
		outputPath := strings.TrimSuffix(path, ".gz")

		// check if file exists and force is false
		if !opts.force {
			err := checkFileExists(outputPath)
			if err != nil {
				return err
			}
		}

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

		_, err = io.Copy(outFile, reader)
		if err != nil {
			return fmt.Errorf("failed to write extracted file: %w", err)
		}
		return nil
	}

	return fmt.Errorf("unsupported file format for extracting")
}

const tarBlockSize = 512

func isZeroBlock(b []byte) bool {
	return bytes.Count(b, []byte{0}) == len(b)
}

// zeroBlockSkippingReader wraps an io.Reader and skips over zero-filled 512-byte blocks
type zeroBlockSkippingReader struct {
	r      io.Reader
	buffer []byte // buffer for checking/skipping zero blocks
	offset int    // current offset in buffer
	valid  int    // valid bytes in buffer
}

func (z *zeroBlockSkippingReader) Read(p []byte) (n int, err error) {
	// if we have buffered data, return it first
	if z.offset < z.valid {
		n = copy(p, z.buffer[z.offset:z.valid])
		z.offset += n
		return n, nil
	}

	// read next block
	if z.buffer == nil {
		z.buffer = make([]byte, tarBlockSize)
	}
	z.valid = 0
	z.offset = 0

	// read full block
	for z.valid < tarBlockSize {
		rn, rerr := z.r.Read(z.buffer[z.valid:tarBlockSize])
		z.valid += rn
		if rerr != nil {
			if rerr == io.EOF && z.valid > 0 {
				break
			}
			return 0, rerr
		}
	}

	// if we read full block and it's all zeros, skip it and try again
	if z.valid == tarBlockSize && isZeroBlock(z.buffer) {
		z.valid = 0
		z.offset = 0
		return z.Read(p)
	}

	// return data from buffer
	n = copy(p, z.buffer[z.offset:z.valid])
	z.offset += n
	return n, nil
}

func extractTar(reader io.Reader, outputDir string, opts *ExtractOptions) error {
	zeroSkippingReader := &zeroBlockSkippingReader{r: reader, buffer: nil, offset: 0, valid: 0}
	tarReader := tar.NewReader(zeroSkippingReader)

	// clean and resolve the output directory to an absolute path
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
			// check if file exists and force is false
			if !opts.force {
				err := checkFileExists(target)
				if err != nil {
					return err
				}
			}

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
