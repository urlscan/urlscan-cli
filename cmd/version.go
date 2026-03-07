package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/pkg/utils"
	"github.com/urlscan/urlscan-cli/pkg/version"
)

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
				msg := fmt.Sprintf("Update available: %s", r.version)
				utils.Notify(msg, utils.Orange)
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
