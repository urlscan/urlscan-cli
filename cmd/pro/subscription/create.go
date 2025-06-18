package subscription

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var createCmdExample = `  urlscan pro subscription create -s <search-id-1> -s <search-id-2> -f <frequency> -e <email-address-1> -e <email-address-2> -n <name>`

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new subscription",
	Example: createCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
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
		id, _ := cmd.Flags().GetString("subscription-id")
		if id == "" {
			id = uuid.New().String()
		}

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
	setCreateOrUpdateFlags(createCmd)
	// optional flag, only for create command
	createCmd.Flags().StringP("subscription-id", "i", "", "Subscription ID (optional, if not provided a new id will be generated)")

	RootCmd.AddCommand(createCmd)
}
