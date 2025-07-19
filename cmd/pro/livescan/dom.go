package livescan

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var domCmdExample = `  urlscan pro live-scan dom <scan-id> -s <scanner-id>
  echo <scan-id> | urlscan pro live-scan dom - -s <scanner-id>`

var domCmd = &cobra.Command{
	Use:     "dom",
	Short:   "Get dom of a live scan",
	Example: domCmdExample,
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

		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = fmt.Sprintf("%s.html", scanId)
		}
		force, _ := cmd.Flags().GetBool("force")

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		url := api.URL("/api/v1/livescan/%s/dom/%s", scannerId, scanId)
		opts := utils.NewDownloadOptions(
			utils.WithDownloadClient(client),
			utils.WithDownloadURL(url),
			utils.WithDownloadOutput(output),
			utils.WithDownloadForce(force),
		)
		err = utils.DownloadWithSpinner(opts)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	addScannerIdFlag(domCmd)
	flags.AddForceFlag(domCmd)
	flags.AddOutputFlag(domCmd, "<uuid>.html")

	RootCmd.AddCommand(domCmd)
}
