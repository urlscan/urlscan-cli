package flags

import (
	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/api"
)

func AddForceFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("force", "f", false, "Force overwrite an existing file.")
}

func AddNoLimitFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("no-limit", false, "Don't limit the number of results returned by the iterator, limit is ignored if it's set (default false)")
}

func AddLimitFlag(cmd *cobra.Command) {
	cmd.Flags().IntP("limit", "l", api.MaxTotal, "Maximum number of results that will be returned by the iterator")
}

func AddSizeFlag(cmd *cobra.Command) {
	cmd.Flags().IntP("size", "s", 100, "Number of results returned by the iterator in each batch")
}
