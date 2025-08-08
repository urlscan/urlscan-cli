package livescan

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var scannersCmdExample = `  urlscan pro livescan scanners`

var scannersCmd = &cobra.Command{
	Use:     "scanners",
	Short:   "List available scanners along with their current metadata",
	Example: scannersCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		url := api.URL("/api/v1/livescan/scanners")
		result, err := client.Get(url)
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJSON())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(scannersCmd)
}
