package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/urlscan/urlscan-cli/pkg/version"
)

const (
	orange = "\033[38;5;208m"
	reset  = "\033[0m"
)

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
	PersistentPostRunE: func(cmd *cobra.Command, args []string) (err error) {
		type result struct {
			version   string
			hasUpdate bool
		}
		ch := make(chan result, 1)

		timeout := 5

		go func() {
			latest, err := version.CheckLatest(timeout)
			hasUpdate := false
			if err == nil {
				hasUpdate = version.IsNewer(version.Version, latest)
			}
			ch <- result{latest, hasUpdate}
		}()

		select {
		case r := <-ch:
			if r.hasUpdate {
				printUpdateNotice(r.version)
			}
		case <-time.After(time.Duration(timeout) * time.Second): // don't hang if GitHub is slow/unreachable
			// silently skip
		}

		return err
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
