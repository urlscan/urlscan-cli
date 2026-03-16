package livescan

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var storeCmdExample = `  urlscan pro livescan store <scan-id> -S <scanner-id>
  echo <scan-id> | urlscan pro livescan store - -s <scanner-id>
  urlscan pro livescan store <scan-id> -s <scanner-id> --json '{"task":{"visibility":"private"}}'`

var storeCmd = &cobra.Command{
	Use:     "store",
	Short:   "Store the temporary scan as a permanent snapshot",
	Example: storeCmdExample,
	Annotations: map[string]string{
		"args": "exact1",
	},
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

		opts := []api.LiveScanStoreOption{}

		json, err := flags.GetJSON(cmd)
		if err != nil {
			return err
		}
		if json != nil {
			opts = append(opts, api.WithLiveScanStoreExtra(json))
		} else {
			visibility, _ := cmd.Flags().GetString("visibility")
			opts = append(opts, api.WithLiveScanStoreTaskVisibility(visibility))
		}

		result, err := client.StoreLiveScanResult(
			scannerId, scanId,
			opts...,
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
	flags.AddJSONFlag(storeCmd)

	RootCmd.AddCommand(storeCmd)
}
