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
		opts, err := mapCmdToSubscriptionOptions(cmd)
		if err != nil {
			return cmd.Usage()
		}

		id, _ := cmd.Flags().GetString("subscription-id")
		if id == "" {
			id = uuid.New().String()
		}

		err = utils.ValidateUUID(id)
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		result, err := client.CreateSubscription(
			append([]api.SubscriptionOption{api.WithSubscriptionID(id)}, opts...)...,
		)
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJson())

		return nil
	},
}

func init() {
	setCreateOrUpdateFlags(createCmd)
	// optional flag, only for create command
	createCmd.Flags().StringP("subscription-id", "i", "", "Subscription ID (optional, if not provided a new id will be generated)")

	RootCmd.AddCommand(createCmd)
}
