package livescan

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var scannersCmdExample = `  urlscan pro live-scan scanners`

var scannersCmd = &cobra.Command{
	Use:     "scanners",
	Short:   "List available scanners along with their current metadata",
	Example: scannersCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		resp, err := client.NewRequest().Get(api.PrefixedPath("/livescan/scanners"))
		if err != nil {
			return err
		}

		fmt.Print(resp.PrettyJSON())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(scannersCmd)
}
