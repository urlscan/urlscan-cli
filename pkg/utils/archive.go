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

func Extract(path string, opts *ExtractOptions) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil && err == nil {
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

const (
	tarBlockSize        = 512
	tarEndZeroBlockSize = 2
)

func isZeroBlock(b []byte) bool {
	return bytes.Count(b, []byte{0}) == len(b)
}

func scanConcatTar(r io.Reader, handle func(*tar.Reader) error) error {
	var buf bytes.Buffer
	zeroBlockCount := 0

	block := make([]byte, tarBlockSize)

	for {
		_, err := io.ReadFull(r, block)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if isZeroBlock(block) {
			zeroBlockCount++
		} else {
			zeroBlockCount = 0
		}

		buf.Write(block)

		if zeroBlockCount >= tarEndZeroBlockSize {
			tr := tar.NewReader(bytes.NewReader(buf.Bytes()))
			err := handle(tr)
			if err != nil {
				return err
			}

			buf.Reset()
			zeroBlockCount = 0
		}
	}
}

func extractTarFile(reader *tar.Reader, outputDir string, opts *ExtractOptions) error {
	header, err := reader.Next()
	if err == io.EOF {
		return err
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

		_, err = io.Copy(outFile, reader)
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

	return nil
}

func extractTar(reader io.Reader, outputDir string, opts *ExtractOptions) error {
	// clean and resolve the output directory to an absolute path
	cleanOutputDir, err := filepath.Abs(filepath.Clean(outputDir))
	if err != nil {
		return fmt.Errorf("failed to resolve output directory: %w", err)
	}

	return scanConcatTar(reader, func(tr *tar.Reader) error {
		for {
			err := extractTarFile(tr, cleanOutputDir, opts)
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
		}
	})
}
