package subscription

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var listCmdExample = `  urlscan pro subscription list`

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get a list of subscriptions for the current user",
	Example: listCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		url := api.URL("/api/v1/user/subscriptions/")
		result, err := client.Get(url)
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJson())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
