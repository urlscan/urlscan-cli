package brand

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var availableCmdExample = `  urlscan pro brand available`

var availableCmd = &cobra.Command{
	Use:     "available",
	Short:   "API Endpoint to get a list of brands that are tracked as part of urlscan's brand and phishing detection",
	Example: availableCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		resp, err := client.NewRequest().Get(api.PrefixedPath(("/pro/availableBrands")))
		if err != nil {
			return err
		}

		fmt.Print(resp.PrettyJSON())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(availableCmd)
}
