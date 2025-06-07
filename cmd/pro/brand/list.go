package brand

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var ListCmdExample = `  urlscan pro brand list`

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get the list of brands that we are able to detect phishing pages for, the total number of detected pages, and the latest hit for each brand",
	Example: ListCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		url := api.URL("/api/v1/pro/brands")
		result, err := client.Get(url)
		if err != nil {
			return err
		}

		fmt.Print(string(result.Raw))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
