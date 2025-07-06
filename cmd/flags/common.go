package flags

import (
	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/api"
)

func AddForceFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("force", "f", false, "Force overwrite an existing file.")
}

func AddAllFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("all", false, "Return all results; limit is ignored if --all is specified (default false)")
}

func AddLimitFlag(cmd *cobra.Command) {
	cmd.Flags().IntP("limit", "l", api.MaxTotal, "Maximum number of results that will be returned by the iterator")
}

func AddSizeFlag(cmd *cobra.Command, value int) {
	cmd.Flags().IntP("size", "s", value, "Number of results returned by the iterator in each batch")
}
