package subscription

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var createCmdExample = `  urlscan pro subscription create <subscription-id> -s <search-id-1> -s <search-id-2> -f <frequency> -e <email-address-1> -e <email-address-2> -n <name>`

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new subscription",
	Example: createCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		reader := utils.StringReaderFromCmdArgs(args)
		id, err := reader.ReadString()
		if err != nil {
			return err
		}

		// required flags (show usage if any are missing)
		searchIds, _ := cmd.Flags().GetStringSlice("search-ids")
		if len(searchIds) == 0 {
			return cmd.Usage()
		}
		frequency, _ := cmd.Flags().GetString("frequency")
		if frequency == "" {
			return cmd.Usage()
		}
		emailAddresses, _ := cmd.Flags().GetStringSlice("email-addresses")
		if len(emailAddresses) == 0 {
			return cmd.Usage()
		}
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return cmd.Usage()
		}
		isActive, _ := cmd.Flags().GetBool("is-active")
		ignoreTime, _ := cmd.Flags().GetBool("ignore-time")

		// optional flags
		description, _ := cmd.Flags().GetString("description")

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		res, err := client.CreateSubscription(
			api.WithSubscriptionID(id),
			api.WithSubscriptionSearchIds(searchIds),
			api.WithSubscriptionFrequency(frequency),
			api.WithSubscriptionEmailAddresses(emailAddresses),
			api.WithSubscriptionName(name),
			api.WithSubscriptionDescription(description),
			api.WithSubscriptionIsActive(isActive),
			api.WithSubscriptionIgnoreTime(ignoreTime),
		)
		if err != nil {
			return err
		}

		fmt.Print(string(res.Raw))

		return nil
	},
}

func init() {
	// required flags
	createCmd.Flags().StringSliceP("search-ids", "s", []string{}, "Array of search IDs associated with this subscription (required)")
	createCmd.Flags().StringP("frequency", "f", "", "Frequency of notifications (required/live, hourly or daily)")
	createCmd.Flags().StringSliceP("email-addresses", "e", []string{}, "Email addresses to send notifications to (required)")
	createCmd.Flags().StringP("name", "n", "", "Name of the subscription (required)")
	createCmd.Flags().BoolP("is-active", "a", true, "Whether the subscription is active (required/defaults to true)")
	createCmd.Flags().BoolP("ignore-time", "t", false, "Whether to ignore time constraints (required/defaults to false)")
	// optional flags
	createCmd.Flags().StringP("description", "d", "", "Description of the subscription (optional)")

	RootCmd.AddCommand(createCmd)
}
