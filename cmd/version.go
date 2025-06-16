package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
