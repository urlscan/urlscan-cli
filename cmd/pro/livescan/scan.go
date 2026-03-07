package livescan

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var scanCmdExample = `  urlscan pro livescan scan <url>
  echo <url> | urlscan pro livescan scan
  urlscan pro livescan scan -s <scanner-id> --json '{"task":{"url":"...","visibility":"private"},"scanner":{"pageTimeout":10000}}'`

var scanCmd = &cobra.Command{
	Use:     "scan",
	Short:   "Task a URL to be scanned",
	Example: scanCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		scannerId, _ := cmd.Flags().GetString("scanner-id")
		if scannerId == "" {
			return cmd.Usage()
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		opts := []api.LiveScanOption{}

		json, err := flags.GetJSON(cmd)
		if err != nil {
			return err
		}
		if json != nil {
			opts = append(opts, api.WithLiveScanExtra(json))
		} else {
			reader := utils.StringReaderFromCmdArgs(args)
			url, err := reader.ReadString()
			if err != nil {
				return err
			}

			visibility, _ := cmd.Flags().GetString("visibility")
			pageTimeout, _ := cmd.Flags().GetInt("page-timeout")
			captureDelay, _ := cmd.Flags().GetInt("capture-delay")
			extraHeaders, _ := cmd.Flags().GetStringToString("extra-headers")
			enableFeatures, _ := cmd.Flags().GetStringSlice("enable-features")
			disableFeatures, _ := cmd.Flags().GetStringSlice("disable-features")

			opts = append(opts,
				api.WithLiveScanTaskURL(url),
				api.WithLiveScanTaskVisibility(visibility),
				api.WithLiveScanScannerPageTimeout(pageTimeout),
				api.WithLiveScanScannerCaptureDelay(captureDelay),
				api.WithLiveScanScannerExtraHeaders(extraHeaders),
				api.WithLiveScanScannerEnableFeatures(enableFeatures),
				api.WithLiveScanScannerDisableFeatures(disableFeatures),
			)
		}

		blocking, _ := cmd.Flags().GetBool("blocking")
		var result *api.Response
		if blocking {
			result, err = client.TriggerLiveScan(scannerId, opts...)
			if err != nil {
				return err
			}
		} else {
			result, err = client.TriggerNonBlockingLiveScan(scannerId, opts...)
			if err != nil {
				return err
			}
		}

		fmt.Print(result.PrettyJSON())

		return nil
	},
}

func init() {
	scanCmd.Flags().IntP("page-timeout", "p", 10_000, "Time to wait for the whole scan process (in ms)")
	scanCmd.Flags().IntP("capture-delay", "c", 10_000, "Delay after page has finished loading before capturing page content (in ms)")
	scanCmd.Flags().StringToStringP("extra-headers", "H", map[string]string{}, "Extra headers to send with the request (e.g., User-Agent: urlscan-cli)")
	scanCmd.Flags().StringSliceP("enable-features", "e", []string{}, "Features to enable (bannerBypass, downloadWait, fullscreen)")
	scanCmd.Flags().StringSliceP("disable-features", "d", []string{}, "Features to disable (annotation, dom, downloads, hideheadless, pageInformation, responses, screenshot, screenshotCompression, stealth)")
	scanCmd.Flags().BoolP("blocking", "b", true, "Whether to do a blocking scan or not")

	addVisibilityFlag(scanCmd)
	addScannerIdFlag(scanCmd)
	flags.AddJSONFlag(scanCmd)

	RootCmd.AddCommand(scanCmd)
}
