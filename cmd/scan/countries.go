package scan

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var countriesCmdExample = `  urlscan scan countries`

var countriesCmd = &cobra.Command{
	Use:     "countries",
	Short:   "Retrieve countries available for scanning using the Scan API",
	Example: countriesCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		result, err := client.NewRequest().Get(api.PrefixedPath("/availableCountries"))
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJSON())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(countriesCmd)
}
