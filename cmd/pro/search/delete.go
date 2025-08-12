package search

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var deleteCmdExample = `  urlscan pro saved-search delete <search-id>
  echo "<search-id>" | urlscan pro saved-search delete -`

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a saved search",
	Example: deleteCmdExample,
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

		resp, err := client.NewRequest().Delete(api.PrefixedPath(fmt.Sprintf("/user/searches/%s/", id)))
		if err != nil {
			return err
		}

		fmt.Print(resp.PrettyJSON())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
