package datadump

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var ListCmdExample = `  urlscan pro datadump list days/api
  urlscan pro datadump list hours/api/20260101
  echo "<path>" | urlscan pro datadump list -

  NOTE: path format is <time-window>/<file-type>/<date>
        - time-window: days | hours | minutes
        - file-type: api | search | screenshot | dom
        - date: YYYYMMDD format date (optional if time-window is days)`

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get the list of data dump files",
	Example: ListCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := utils.StringReaderFromCmdArgs(args)
		path, err := reader.ReadString()
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		resp, err := client.NewRequest().Get(api.PrefixedPath(fmt.Sprintf("/datadump/list/%s", path)))
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
