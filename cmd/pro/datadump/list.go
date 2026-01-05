package datadump

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var ListCmdExample = `  urlscan pro datadump list --time-window days --file-type api --date 20260101`

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get the list of data dump files",
	Example: ListCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		timeWindow, err := cmd.Flags().GetString("time-window")
		if err != nil {
			return err
		}

		fileType, err := cmd.Flags().GetString("file-type")
		if err != nil {
			return err
		}

		date, err := cmd.Flags().GetString("date")
		if err != nil {
			return err
		}

		resp, err := client.NewRequest().Get(api.PrefixedPath(fmt.Sprintf("/datadump/list/%s/%s/%s", timeWindow, fileType, date)))
		if err != nil {
			return err
		}

		fmt.Print(resp.PrettyJSON())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("time-window", "t", "", "Time window for the data dump (days, hours, minutes)")
	listCmd.Flags().StringP("file-type", "f", "", "File type for the data dump (api, search, screenshot, dom)")
	listCmd.Flags().StringP("date", "d", "", "Date for the data dump (e.g., 20260101)")
}
