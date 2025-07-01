package flags

import "github.com/spf13/cobra"

func AddForceFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("force", "f", false, "Force overwrite an existing file.")
}
