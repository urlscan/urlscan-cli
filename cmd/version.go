package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// This variable is initialized using the -X linker flag when building the binary.
var Version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("urlscan-cli %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
