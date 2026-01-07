package subscription

import (
	"fmt"

	"github.com/spf13/cobra"

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

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		resp, err := client.CreateSubscription(
			opts...,
		)
		if err != nil {
			return err
		}

		fmt.Print(resp.PrettyJSON())

		return nil
	},
}

func init() {
	setCreateOrUpdateFlags(createCmd)

	RootCmd.AddCommand(createCmd)
}
