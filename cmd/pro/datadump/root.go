package datadump

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "datadump",
	Short: "Data dump sub-commands",
}
