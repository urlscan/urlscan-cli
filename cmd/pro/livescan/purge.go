package livescan

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var purgeCmdExample = `  urlscan pro live-scan purge <scan-id> -s <scanner-id>
  echo <scan-id> | urlscan pro live-scan purge - -s <scanner-id>`

var purgeCmd = &cobra.Command{
	Use:     "purge",
	Short:   "Purge a temporary live scan",
	Example: purgeCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		reader := utils.StringReaderFromCmdArgs(args)
		scanId, err := reader.ReadString()
		if err != nil {
			return err
		}

		err = utils.ValidateUUID(scanId)
		if err != nil {
			return err
		}

		scannerId, _ := cmd.Flags().GetString("scanner-id")
		if scannerId == "" {
			return cmd.Usage()
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		result, err := client.Delete(api.URL("/api/v1/livescan/%s/%s/", scannerId, scanId))
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJson())

		return nil
	},
}

func init() {
	addScannerIdFlag(purgeCmd)

	RootCmd.AddCommand(purgeCmd)
}
