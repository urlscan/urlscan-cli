package utils

import (
	"bufio"
	"os"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func resolveFile(s string) ([]string, error) {
	var outputs []string

	file, err := os.Open(s)
	if err != nil {
		return nil, err
	}
	defer file.Close() // nolint:errcheck

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
