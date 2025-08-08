package search

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var getCmdExample = `  urlscan pro saved-search get <search-id>
  echo <search-id> | urlscan pro saved-search get -`

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "Get search results for a specified saved search",
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

		resp, err := client.NewRequest().Get(api.PrefixedPath(fmt.Sprintf("/user/searches/%s/results/", id)))
		if err != nil {
			return err
		}

		fmt.Print(resp.PrettyJSON())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
