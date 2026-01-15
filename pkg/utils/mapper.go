package utils

import (
	"bufio"
	"fmt"
	"os"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func checkFileExists(path string) error {
	if fileExists(path) {
		return fmt.Errorf("%s already exists, use --force to overwrite", path)
	}
	return nil
}

func resolveFile(s string) (outputs []string, err error) {
	file, err := os.Open(s)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		if word != "" {
			outputs = append(outputs, word)
		}
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return outputs, nil
}

func ResolveFileOrValue(value string) ([]string, error) {
	if fileExists(value) {
		return resolveFile(value)
	}
	return []string{value}, nil
}
