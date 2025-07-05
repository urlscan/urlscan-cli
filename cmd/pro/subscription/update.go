package subscription

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var updateCmdExample = `  urlscan pro subscription update <subscription-id> -s <search-id-1> -s <search-id-2> -f <frequency> -e <email-address-1> -e <email-address-2> -n <name>`

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a subscription",
	Example: updateCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		reader := utils.StringReaderFromCmdArgs(args)
		id, err := reader.ReadString()
		if err != nil {
			return err
		}
		err = utils.ValidateUUID(id)
		if err != nil {
			return err
		}

		opts, err := mapCmdToSubscriptionOptions(cmd)
		if err != nil {
			return cmd.Usage()
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		result, err := client.UpdateSubscription(
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
	setCreateOrUpdateFlags(updateCmd)

	RootCmd.AddCommand(updateCmd)
}
