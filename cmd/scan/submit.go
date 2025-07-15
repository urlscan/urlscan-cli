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

		scanOpts, err := mapCmdToScanOptions(cmd)
		if err != nil {
			return err
		}

		wait, _ := cmd.Flags().GetBool("wait")
		maxWait, _ := cmd.Flags().GetInt("max-wait")

		withBoth, _ := cmd.Flags().GetBool("with-both")
		withScreenshot, _ := cmd.Flags().GetBool("with-screenshot")
		withScreenshot = withScreenshot || withBoth
		withDOM, _ := cmd.Flags().GetBool("with-dom")
		withDOM = withDOM || withBoth
		force, _ := cmd.Flags().GetBool("force")

		// override wait if any of with flag is set
		wait = wait || withScreenshot || withDOM

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
			fmt.Print(scanResult.PrettyJson())
			return nil
		}

		ctx := cmd.Context()
		waitResult, err := client.WaitAndGetResult(ctx, scanResult.UUID, maxWait)
		if err != nil {
			return err
		}

		fmt.Print(waitResult.PrettyJson())

		if withScreenshot {
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

		if withDOM {
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

	submitCmd.Flags().Bool("with-screenshot", false, "Download a screenshot (this flag overrides wait flag to true)")
	submitCmd.Flags().Bool("with-dom", false, "Download a DOM (this flag overrides wait flag to true)")
	submitCmd.Flags().Bool("with-both", false, "Download both a screenshot and a DOM (this flag overrides wait, with-screen and with-both flags to true)")

	RootCmd.AddCommand(submitCmd)
}
