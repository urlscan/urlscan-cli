package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/pkg/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("urlscan-cli %s\n", version.Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
