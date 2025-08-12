package subscription

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var searchCmdExample = `  urlscan pro subscription search <subscription-id>
  echo "<subscription-id>" | urlscan pro subscription search`

var searchCmd = &cobra.Command{
	Use:     "search",
	Short:   "Get the search results for a specific subscription and datasource",
	Example: searchCmdExample,
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

		datasource, _ := cmd.Flags().GetString("datasource")

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		resp, err := client.NewRequest().Get(api.PrefixedPath(fmt.Sprintf("/user/subscriptions/%s/results/%s/", id, datasource)))
		if err != nil {
			return err
		}

		fmt.Print(resp.PrettyJSON())

		return nil
	},
}

func init() {
	searchCmd.Flags().StringP("datasource", "D", "scans", "datasource to search in (scans or hostnames)")

	RootCmd.AddCommand(searchCmd)
}
