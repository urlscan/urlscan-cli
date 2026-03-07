package subscription

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var updateCmdExample = `  urlscan pro subscription update <subscription-id> -s <search-id-1> -s <search-id-2> -f <frequency> -e <email-address-1> -e <email-address-2> -n <name>
  urlscan pro subscription update <subscription-id> --json '{"subscription":{"searchIds":["..."],"frequency":"live","emailAddresses":["..."],"name":"my-sub"}}'`

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a subscription",
	Example: updateCmdExample,
	Annotations: map[string]string{
		"args": "exact1",
	},
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

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		opts := []api.SubscriptionOption{}

		json, err := flags.GetJSON(cmd)
		if err != nil {
			return err
		}
		if json != nil {
			opts = append(opts, api.WithSubscriptionExtra(json))
		} else {
			mapped, err := mapCmdToSubscriptionOptions(cmd)
			if err != nil {
				return err
			}
			opts = append(opts, mapped...)
		}

		result, err := client.UpdateSubscription(
			id,
			opts...,
		)
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJSON())

		return nil
	},
}

func init() {
	setCreateOrUpdateFlags(updateCmd)
	flags.AddJSONFlag(updateCmd)

	RootCmd.AddCommand(updateCmd)
}
