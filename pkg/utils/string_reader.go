package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type StringReader interface {
	Next() bool
	ReadString() (string, error)
	Value() string
	Err() error
}

type StringArrayReader struct {
	strings []string
	pos     int
	value   string
	err     error
}

func NewStringArrayReader(strings []string) *StringArrayReader {
	return &StringArrayReader{strings: strings}
}

func (sar *StringArrayReader) Next() bool {
	if sar.pos < len(sar.strings) {
		sar.value = sar.strings[sar.pos]
		sar.pos++
		return true
	}
	sar.err = io.EOF
	return false
}

func (sar *StringArrayReader) ReadString() (string, error) {
	if sar.Next() {
		return sar.value, nil
	}
	return "", sar.err
}

func (sar *StringArrayReader) Value() string {
	return sar.value
}

func (sar *StringArrayReader) Err() error {
	return sar.err
}

type StringIOReader struct {
	scanner *bufio.Scanner
	value   string
	err     error
}

func NewStringIOReader(r io.Reader) *StringIOReader {
	return &StringIOReader{scanner: bufio.NewScanner(r)}
}

func (sir *StringIOReader) Next() bool {
	if sir.scanner.Scan() {
		sir.value = strings.TrimSpace(sir.scanner.Text())
		if sir.value != "" {
			return true
		}
	}
	sir.err = sir.scanner.Err()
	if sir.err == nil {
		sir.err = io.EOF
	}
	return false
}

func (sir *StringIOReader) ReadString() (string, error) {
	if sir.Next() {
		return sir.value, nil
	}
	return "", sir.err
}

func (sir *StringIOReader) Value() string {
	return sir.value
}

func (sir *StringIOReader) Err() error {
	return sir.err
}

type MappedStringReader struct {
	q     []string
	r     StringReader
	mapFn func(string) ([]string, error)
	value string
	err   error
}

func NewMappedStringReader(r StringReader, mapFn func(string) ([]string, error)) *MappedStringReader {
	return &MappedStringReader{r: r, mapFn: mapFn}
}

func (msr *MappedStringReader) Pop() (string, error) {
	if len(msr.q) == 0 {
		return "", fmt.Errorf("no strings left in the queue")
	}
	value := msr.value
	msr.q = msr.q[1:]
	return value, nil
}

func (msr *MappedStringReader) setError(err error) {
	msr.err = err
	msr.value = ""
	msr.q = nil
}

func (msr *MappedStringReader) Next() bool {
	// return from the queue if available
	if len(msr.q) > 0 {
		msr.value = msr.q[0]
		msr.q = msr.q[1:]
		return true
	}

	// read the next string and apply the mapping function
	got, err := msr.r.ReadString()
	if err != nil {
		msr.setError(err)
		return false
	}

	mapped, err := msr.mapFn(got)
	if err != nil {
		msr.err = err
		return false
	}
	// set the firs value as the current value, queue the rest
	if len(mapped) > 0 {
		msr.value = mapped[0]
		msr.q = mapped[1:]
		return true
	}

	// if no mapped values, set error and return false
	msr.setError(fmt.Errorf("no mapped values for: %s", got))
	return false
}

func (msr *MappedStringReader) Value() string {
	return msr.value
}

func (msr *MappedStringReader) Err() error {
	return msr.err
}

func (msr *MappedStringReader) ReadString() (string, error) {
	if msr.Next() {
		return msr.value, nil
	}
	return "", msr.err
}

type FilteredStringReader struct {
	r          StringReader
	validateFn func(string) error
}

func NewFilteredStringReader(r StringReader, validateFn func(string) error) *FilteredStringReader {
	return &FilteredStringReader{r: r, validateFn: validateFn}
}

func (f *FilteredStringReader) ReadString() (s string, err error) {
	if f.Next() {
		return f.Value(), nil
	}
	return "", f.r.Err()
}

func (f *FilteredStringReader) Next() bool {
	for next := f.r.Next(); next; next = f.r.Next() {
		value := f.r.Value()
		if f.validateFn(value) == nil {
			return true
		}
	}
	return false
}

func (f *FilteredStringReader) Value() string {
	return f.r.Value()
}

func (f *FilteredStringReader) Err() error {
	return f.r.Err()
}

func ReadAllFromReader(r StringReader) ([]string, error) {
	var results []string
	for r.Next() {
		results = append(results, r.Value())
	}
	if r.Err() != nil && r.Err() != io.EOF {
		return nil, r.Err()
	}
	return results, nil
}

func StringReaderFromCmdArgs(args []string) StringReader {
	if len(args) == 1 && args[0] == "-" {
		return NewStringIOReader(os.Stdin)
	}
	return NewStringArrayReader(args)
}
