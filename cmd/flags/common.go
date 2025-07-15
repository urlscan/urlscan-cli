package flags

import (
	"fmt"

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

func AddOutputFlag(cmd *cobra.Command, defaultExample string) {
	cmd.Flags().StringP("output", "o", "", fmt.Sprintf("Output file name (default %s)", defaultExample))
}

func AddMaxConcurrencyFlag(cmd *cobra.Command) {
	cmd.Flags().Int("max-concurrency", 5, "Maximum number of concurrent requests for batch operation")
}

func AddTotalTimeoutFlag(cmd *cobra.Command) {
	cmd.Flags().Int("total-timeout", 60*30, "Total timeout for the batch operation in seconds, 0 means no timeout")
}
