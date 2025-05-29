package utils

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type StringReader interface {
	ReadString() (string, error)
}

type StringArrayReader struct {
	strings []string
	pos     int
}

func NewStringArrayReader(strings []string) *StringArrayReader {
	return &StringArrayReader{strings: strings}
}

func (sar *StringArrayReader) ReadString() (string, error) {
	if sar.pos == len(sar.strings) {
		return "", io.EOF
	}
	s := sar.strings[sar.pos]
	sar.pos++
	return s, nil
}

type StringIOReader struct {
	scanner *bufio.Scanner
}

func NewStringIOReader(r io.Reader) *StringIOReader {
	return &StringIOReader{scanner: bufio.NewScanner(r)}
}

func (sir *StringIOReader) ReadString() (string, error) {
	for sir.scanner.Scan() {
		s := strings.TrimSpace(sir.scanner.Text())
		if s != "" {
			return s, nil
		}
	}
	return "", io.EOF
}

func StringReaderFromCmdArgs(args []string) StringReader {
	if len(args) == 1 && args[0] == "-" {
		return NewStringIOReader(os.Stdin)
	}
	return NewStringArrayReader(args)
}
