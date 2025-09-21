package channel

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var listCmdExample = `  urlscan pro channel list`

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get a list of channels for the current user",
	Example: listCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		resp, err := client.NewRequest().Get(api.PrefixedPath("/user/channels/"))
		if err != nil {
			return err
		}

		fmt.Print(resp.PrettyJSON())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
