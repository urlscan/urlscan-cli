package subscription

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var createCmdExample = `  urlscan pro subscription create -s <search-id-1> -s <search-id-2> -f <frequency> -e <email-address-1> -e <email-address-2> -n <name>
  urlscan pro subscription create --json '{"subscription":{"searchIds":["..."],"frequency":"live","emailAddresses":["..."],"name":"my-sub"}}'`

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new subscription",
	Example: createCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		result, err := client.CreateSubscription(opts...)
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJSON())

		return nil
	},
}

func init() {
	setCreateOrUpdateFlags(createCmd)
	flags.AddJSONFlag(createCmd)

	RootCmd.AddCommand(createCmd)
}
