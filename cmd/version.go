package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/urlscan/urlscan-cli/pkg/version"
)

const (
	orange = "\033[38;5;208m"
	reset  = "\033[0m"
)

func checkAndPrintUpdate() error {
	latest, err := version.CheckLatest()
	if err != nil {
		return err
	}
	if version.IsNewer(version.Version, latest) {
		printUpdateNotice(latest)
	}
	return nil
}

func printUpdateNotice(latest string) {
	msg := fmt.Sprintf("Update available: %s", latest)

	width, _, err := term.GetSize(int(os.Stderr.Fd()))
	if err != nil || width <= 0 {
		width = 80
	}

	if pad := width - len(msg); pad > 0 {
		fmt.Fprintf(os.Stderr, "%s%s%s%s\n", orange, strings.Repeat(" ", pad), msg, reset)
	} else {
		fmt.Fprintf(os.Stderr, "%s%s%s\n", orange, msg, reset)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return cmd.Usage()
		}

		fmt.Printf("urlscan-cli %s\n", version.Version)
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return checkAndPrintUpdate()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
