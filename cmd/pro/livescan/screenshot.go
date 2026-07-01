package livescan

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var screenshotCmdExample = `  urlscan pro livescan screenshot <scan-id> -s <scanner-id>
  echo <scan-id> | urlscan pro livescan screenshot - -s <scanner-id>`

var screenshotCmd = &cobra.Command{
	Use:     "screenshot",
	Short:   "Get a screenshot of a live scan",
	Example: screenshotCmdExample,
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
			output = fmt.Sprintf("%s.png", scanId)
		}
		force, _ := cmd.Flags().GetBool("force")

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		url := api.PrefixedPath(fmt.Sprintf("/livescan/%s/screenshot/%s/", scannerId, scanId))
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
	addScannerIdFlag(screenshotCmd)
	flags.AddForceFlag(screenshotCmd)
	flags.AddOutputFlag(screenshotCmd, "<uuid>.png")

	RootCmd.AddCommand(screenshotCmd)
}
