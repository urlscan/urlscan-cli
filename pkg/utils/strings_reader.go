package utils

import (
	"bufio"
	"io"
	"os"
)

type StringsReader interface {
	ReadStrings() (string, error)
	ReadAll() ([]string, error)
}

type MappedStringsReader struct {
	r     StringReader
	mapFn func(string) ([]string, error)
}

func NewMappedStringsReader(r StringReader, mapFn func(string) ([]string, error)) *MappedStringsReader {
	return &MappedStringsReader{r: r, mapFn: mapFn}
}

func (m *MappedStringsReader) ReadStrings() ([]string, error) {
	s, err := m.r.ReadString()
	if err != nil {
		return nil, err
	}

	mapped, err := m.mapFn(s)
	if err != nil {
		return nil, err
	}
	return mapped, nil
}

func (m *MappedStringsReader) ReadAll() ([]string, error) {
	var outputs []string

	for {
		mapped, err := m.ReadStrings()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		outputs = append(outputs, mapped...)
	}

	return outputs, nil
}

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
