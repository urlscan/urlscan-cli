package scan

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/samber/mo"
	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

type scanner struct {
	client     *utils.APIClient
	scanOpts   []api.ScanOption
	batchOpts  []api.BatchOption
	wait       bool
	maxWait    int
	force      bool
	screenshot bool
	dom        bool
	ctx        context.Context
}

func (s *scanner) newBatchScanWithDownloadTask(url string) api.BatchTask[*api.Response] {
	return func(c *api.Client, ctx context.Context) mo.Result[*api.Response] {
		req := c.NewScanRequest(url, s.scanOpts...)
		resp, err := req.Do()
		if err != nil {
			return mo.Err[*api.Response](err)
		}

		scanResult := &api.ScanResult{} // nolint: exhaustruct
		err = resp.Unmarshal(scanResult)
		if err != nil {
			return mo.Err[*api.Response](err)
		}
		_, err = c.WaitAndGetResult(ctx, scanResult.UUID, s.maxWait)
		if err != nil {
			return mo.Err[*api.Response](err)
		}

		if s.screenshot {
			downloadOpts := utils.NewDownloadOptions(
				utils.WithDownloadClient(s.client),
				utils.WithDownloadScreenshot(scanResult.UUID),
				utils.WithDownloadOutput(fmt.Sprintf("%s.png", scanResult.UUID)),
				utils.WithDownloadForce(s.force),
				utils.WithDownloadSilent(true),
			)
			downloadErr := utils.Download(downloadOpts)
			if downloadErr != nil {
				fmt.Fprint(os.Stderr, "Error downloading screenshot: ", downloadErr)
			}
		}

		if s.dom {
			downloadOpts := utils.NewDownloadOptions(
				utils.WithDownloadClient(s.client),
				utils.WithDownloadDOM(scanResult.UUID),
				utils.WithDownloadOutput(fmt.Sprintf("%s.html", scanResult.UUID)),
				utils.WithDownloadForce(s.force),
				utils.WithDownloadSilent(true),
			)
			downloadErr := utils.Download(downloadOpts)
			if downloadErr != nil {
				fmt.Fprint(os.Stderr, "Error downloading DOM: ", downloadErr)
			}
		}

		return mo.Ok(resp)
	}
}

func (s *scanner) do(urls []string) error {
	tasks := make([]api.BatchTask[*api.Response], len(urls))
	for i, url := range urls {
		if s.wait {
			tasks[i] = s.newBatchScanWithDownloadTask(url)
		} else {
			tasks[i] = s.client.NewBatchScanTask(url, s.scanOpts...)
		}
	}

	results, err := api.Batch(s.client.Client, tasks, s.batchOpts...)
	if err != nil {
		return err
	}

	pairs := utils.NewBatchJSONResultPairs(urls, results)

	b, err := json.MarshalIndent(pairs, "", "  ")
	if err != nil {
		return err
	}

	fmt.Print(string(b))

	return nil
}

func newScanner(cmd *cobra.Command) (*scanner, error) {
	scanOpts := newScanOptions(cmd)

	maxConcurrency, _ := cmd.Flags().GetInt("max-concurrency")
	timeout, _ := cmd.Flags().GetInt("timeout")

	wait := newWaitFlag(cmd)
	maxWait, _ := cmd.Flags().GetInt("max-wait")

	screenshot := newScreenshotFlag(cmd)
	dom := newDOMFlag(cmd)
	force, _ := cmd.Flags().GetBool("force")

	// override wait if dom or screenshot flag is set
	wait = wait || screenshot || dom

	client, err := utils.NewAPIClient()
	if err != nil {
		return nil, err
	}

	return &scanner{
		client:   client,
		scanOpts: scanOpts,
		batchOpts: []api.BatchOption{
			api.WithBatchMaxConcurrency(maxConcurrency),
			api.WithBatchTimeout(timeout),
		},
		wait:       wait,
		maxWait:    maxWait,
		dom:        dom,
		screenshot: screenshot,
		force:      force,
		ctx:        cmd.Context(),
	}, nil
}

var bulkSubmitCmdExample = `  urlscan scan bulk-submit <url>...
  # submit with a file containing URLs per line, space, or tab
  urlscan scan bulk-submit list_of_urls.txt
  # combine the file input and the URL input
  urlscan scan bulk-submit list_of_urls.txt <url>`

var bulkSubmitCmdLong = `Submit multiple URLs to scan in bulk.

This command allows you to submit a list of URLs for scanning in bulk. You can provide URLs via command line arguments or through a file.
Note that the URLs will be validated before submission, and only valid URLs will be processed.`

var bulkSubmitCmd = &cobra.Command{
	Use:     "bulk-submit <url>...",
	Short:   "Bulk submit URLs to scan",
	Long:    bulkSubmitCmdLong,
	Example: bulkSubmitCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Usage()
		}

		reader := utils.NewFilteredStringReader(utils.NewMappedStringReader(utils.StringReaderFromCmdArgs(args), utils.ResolveFileOrValue), utils.ValidateNetworkIndicator)
		urls, err := utils.ReadAllFromReader(reader)
		if err != nil {
			return err
		}

		scanner, err := newScanner(cmd)
		if err != nil {
			return err
		}

		return scanner.do(urls)
	},
}

func init() {
	addScanFlags(bulkSubmitCmd)
	flags.AddForceFlag(bulkSubmitCmd)

	bulkSubmitCmd.Flags().Int("max-concurrency", 5, "Maximum number of concurrent requests for batch operation")
	bulkSubmitCmd.Flags().Int("timeout", 60*30, "Timeout for the batch operation in seconds, 0 means no timeout")

	RootCmd.AddCommand(bulkSubmitCmd)
}
