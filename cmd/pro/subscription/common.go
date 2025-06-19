package subscription

import (
	"github.com/spf13/cobra"
)

// set common flags for create and update commands
func setCreateOrUpdateFlags(cmd *cobra.Command) {
	// required flags
	cmd.Flags().StringSliceP("search-ids", "s", []string{}, "Array of search IDs associated with this subscription (required)")
	cmd.Flags().StringP("frequency", "f", "", "Frequency of notifications (live, hourly or daily) (required) ")
	cmd.Flags().StringSliceP("email-addresses", "e", []string{}, "Email addresses to send notifications to (required)")
	cmd.Flags().StringP("name", "n", "", "Name of the subscription (required)")
	cmd.Flags().BoolP("is-active", "a", true, "Whether the subscription is active (required, defaults to true)")
	cmd.Flags().BoolP("ignore-time", "t", false, "Whether to ignore time constraints (required, defaults to false)")
	// optional flags
	cmd.Flags().StringP("description", "d", "", "Description of the subscription (optional)")
}
