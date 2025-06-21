package search

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var getCmdExample = `  urlscan pro saved-search get <search-id>
  urlscan pro saved-search get --id <search-id>`

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

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		url := api.URL("/api/v1/user/searches/%s/results/", id)
		res, err := client.Get(url)
		if err != nil {
			return err
		}

		fmt.Print(string(res.Raw))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
