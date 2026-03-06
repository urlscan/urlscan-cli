package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/muesli/termenv"
	"golang.org/x/term"
)

func Notify(msg string, color termenv.Color) {
	width, _, err := term.GetSize(int(os.Stderr.Fd()))
	if err != nil || width <= 0 {
		width = 80
	}

	output := termenv.NewOutput(os.Stderr)
	styled := output.String(msg).Foreground(color)

	var formatted string
	if pad := width - len(msg); pad > 0 {
		formatted = fmt.Sprintf("%s%s", strings.Repeat(" ", pad), styled)
	} else {
		formatted = styled.String()
	}

	fmt.Fprintf(os.Stderr, "%s\n", formatted)
}
