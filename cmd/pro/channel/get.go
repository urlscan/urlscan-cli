package channel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var getCmdExample = `  urlscan pro channel get <channel-id>
	echo <channel-id> | urlscan pro channel get -`

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "Get the search results for a specific notification channel",
	Example: getCmdExample,
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

		result, err := client.NewRequest().Get(api.PrefixedPath(fmt.Sprintf("/user/channels/%s", id)))
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJSON())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
