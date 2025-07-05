package livescan

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var resultCmdExample = `  urlscan pro live-scan result <scan-id> -s <scanner-id>
  echo <scan-id> | urlscan pro live-scan result - -s <scanner-id>`

var resultCmd = &cobra.Command{
	Use:     "result",
	Short:   "Get a result of a live scan",
	Example: resultCmdExample,
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

		result, err := client.Get(api.URL("/api/v1/livescan/%s/result/%s", scannerId, scanId))
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJson())

		return nil
	},
}

func init() {
	addScannerIdFlag(resultCmd)

	RootCmd.AddCommand(resultCmd)
}
