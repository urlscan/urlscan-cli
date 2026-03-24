package malicious

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "malicious",
	Short: "Malicious sub-commands",
}
