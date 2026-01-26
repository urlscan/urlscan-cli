package datadump

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var ListCmdExample = `  urlscan pro datadump list days/api
  urlscan pro datadump list hours/api/20260101
  echo "<path>" | urlscan pro datadump list -

  NOTE: path format is <time-window>/<file-type>/<date>
        - time-window: days | hours | minutes. Required.
        - file-type: api | search | screenshot | dom. Required.
        - date: YYYYMMDD format date (optional if time-window is days). Optional.
        if date is not provided, all the available files (files within the last 7 days) will be listed`

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get the list of data dump files",
	Example: ListCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		reader := utils.StringReaderFromCmdArgs(args)
		path, err := reader.ReadString()
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		result, err := client.BulkGetDataDumpList(path)
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJSON())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
