package utils

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	Orange = "\033[38;5;208m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	reset  = "\033[0m"
)

func Notify(msg string, color string) {
	width, _, err := term.GetSize(int(os.Stderr.Fd()))
	if err != nil || width <= 0 {
		width = 80
	}

	var formatted string
	if pad := width - len(msg); pad > 0 {
		formatted = fmt.Sprintf("%s%s%s%s\n", color, strings.Repeat(" ", pad), msg, reset)
	} else {
		formatted = fmt.Sprintf("%s%s%s\n", color, msg, reset)
	}

	fmt.Fprint(os.Stderr, formatted)
}
