package subscription

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/api"
)

// set common flags for create and update commands
func setCreateOrUpdateFlags(cmd *cobra.Command) {
	// required flags
	cmd.Flags().StringSliceP("search-ids", "s", []string{}, "Array of search IDs associated with this subscription (required)")
	cmd.Flags().StringP("frequency", "f", "", "Frequency of notifications (live, hourly or daily) (required)")
	cmd.Flags().StringSliceP("email-addresses", "e", []string{}, "Email addresses to send notifications to (required)")
	cmd.Flags().StringP("name", "n", "", "Name of the subscription (required)")

	// defaulted flags
	cmd.Flags().BoolP("is-active", "a", true, "Whether the subscription is active")
	cmd.Flags().BoolP("ignore-time", "t", false, "Whether to ignore time constraints (default false)")

	// optional flags
	cmd.Flags().StringP("description", "d", "", "Description of the subscription (optional)")
}

func mapCmdToSubscriptionOptions(cmd *cobra.Command) (opts []api.SubscriptionOption, err error) {
	searchIds, _ := cmd.Flags().GetStringSlice("search-ids")
	if len(searchIds) == 0 {
		return nil, fmt.Errorf("search-ids is required")
	}
	frequency, _ := cmd.Flags().GetString("frequency")
	if frequency == "" {
		return nil, fmt.Errorf("frequency is required")
	}
	emailAddresses, _ := cmd.Flags().GetStringSlice("email-addresses")
	if len(emailAddresses) == 0 {
		return nil, fmt.Errorf("email-addresses is required")
	}
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	// optional flags
	isActive, _ := cmd.Flags().GetBool("is-active")
	ignoreTime, _ := cmd.Flags().GetBool("ignore-time")
	description, _ := cmd.Flags().GetString("description")

	return []api.SubscriptionOption{
		api.WithSubscriptionSearchIds(searchIds),
		api.WithSubscriptionFrequency(frequency),
		api.WithSubscriptionEmailAddresses(emailAddresses),
		api.WithSubscriptionName(name),
		api.WithSubscriptionDescription(description),
		api.WithSubscriptionIsActive(isActive),
		api.WithSubscriptionIgnoreTime(ignoreTime),
	}, nil

}
