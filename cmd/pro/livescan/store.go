package livescan

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var storeCmdExample = `  urlscan pro livescan store <scan-id> -S <scanner-id>
  echo <scan-id> | urlscan pro livescan store - -s <scanner-id>`

var storeCmd = &cobra.Command{
	Use:     "store",
	Short:   "Store the temporary scan as a permanent snapshot",
	Example: storeCmdExample,
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

		visibility, _ := cmd.Flags().GetString("visibility")

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		result, err := client.StoreLiveScanResult(
			scannerId, scanId,
			api.WithLiveScanStoreTaskVisibility(visibility),
		)
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJSON())

		return nil
	},
}

func init() {
	addScannerIdFlag(storeCmd)
	addVisibilityFlag(storeCmd)

	RootCmd.AddCommand(storeCmd)
}
