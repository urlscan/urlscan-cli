package scan

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var resultCmdExample = `  urlscan result <uuid>
  echo "<uuid>" | urlscan result -`

var resultCmd = &cobra.Command{
	Use:     "result",
	Short:   "Get a result by UUID",
	Example: resultCmdExample,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := utils.StringReaderFromCmdArgs(args)
		uuid, err := reader.ReadString()
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		url := api.URL("%s", fmt.Sprintf("/api/v1/result/%s/", uuid))
		result, err := client.Get(url)
		if err != nil {
			return err
		}

		fmt.Print(string(result.Raw))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(resultCmd)
}
