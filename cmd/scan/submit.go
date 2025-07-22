package scan

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var submitCmdExample = `  urlscan scan submit <url>...
  echo "<url>" | urlscan scan submit -`

var submitCmd = &cobra.Command{
	Use:     "submit <url>",
	Short:   "Submit a URL to scan",
	Example: submitCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		scanOpts := newScanOptions(cmd)

		wait := newWaitFlag(cmd)
		maxWait, _ := cmd.Flags().GetInt("max-wait")

		screenshot := newScreenshotFlag(cmd)
		dom := newDOMFlag(cmd)
		force, _ := cmd.Flags().GetBool("force")

		reader := utils.StringReaderFromCmdArgs(args)
		url, err := reader.ReadString()
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		scanResult, err := client.Scan(url, scanOpts...)
		if err != nil {
			return err
		}

		if !wait {
			fmt.Print(scanResult.PrettyJSON())
			return nil
		}

		ctx := cmd.Context()
		waitResult, err := client.WaitAndGetResult(ctx, scanResult.UUID, maxWait)
		if err != nil {
			return err
		}

		fmt.Print(waitResult.PrettyJSON())

		if screenshot {
			downloadOpts := utils.NewDownloadOptions(
				utils.WithDownloadClient(client),
				utils.WithDownloadScreenshot(scanResult.UUID),
				utils.WithDownloadOutput(fmt.Sprintf("%s.png", scanResult.UUID)),
				utils.WithDownloadForce(force),
				utils.WithDownloadSilent(true),
			)
			err := utils.DownloadWithSpinner(downloadOpts)
			if err != nil {
				return err
			}
		}

		if dom {
			downloadOpts := utils.NewDownloadOptions(
				utils.WithDownloadClient(client),
				utils.WithDownloadDOM(scanResult.UUID),
				utils.WithDownloadOutput(fmt.Sprintf("%s.html", scanResult.UUID)),
				utils.WithDownloadForce(force),
				utils.WithDownloadSilent(true),
			)
			err := utils.DownloadWithSpinner(downloadOpts)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	addScanFlags(submitCmd)
	flags.AddForceFlag(submitCmd)

	RootCmd.AddCommand(submitCmd)
}
