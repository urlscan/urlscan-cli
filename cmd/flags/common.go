package flags

import "github.com/spf13/cobra"

func AddForceFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("force", "f", false, "Enable to force overwriting an existing file.")
}
