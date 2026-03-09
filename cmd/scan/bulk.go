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

type scanRequest = struct {
	url  string
	opts []api.ScanOption
}

type scanner struct {
	client          *utils.APIClient
	batchOpts       []api.BatchOption
	wait            bool
	maxWait         int
	force           bool
	directoryPrefix string
	screenshot      bool
	dom             bool
	ctx             context.Context
}

func (s *scanner) newBatchScanWithDownloadTask(scanRequest scanRequest) api.BatchTask[*api.Response] {
	return func(c *api.Client, ctx context.Context) mo.Result[*api.Response] {
		req := c.NewScanRequest(scanRequest.url, scanRequest.opts...)
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
				utils.WithDownloadDirectoryPrefix(s.directoryPrefix),
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
				utils.WithDownloadDirectoryPrefix(s.directoryPrefix),
			)
			downloadErr := utils.Download(downloadOpts)
			if downloadErr != nil {
				fmt.Fprint(os.Stderr, "Error downloading DOM: ", downloadErr)
			}
		}

		return mo.Ok(resp)
	}
}

func (s *scanner) do(scanRequests []scanRequest) error {
	urls := make([]string, len(scanRequests))
	tasks := make([]api.BatchTask[*api.Response], len(scanRequests))

	for i, req := range scanRequests {
		if s.wait {
			tasks[i] = s.newBatchScanWithDownloadTask(req)
		} else {
			tasks[i] = s.client.NewBatchScanTask(req.url, req.opts...)
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
	maxConcurrency, _ := cmd.Flags().GetInt("max-concurrency")
	timeout, _ := cmd.Flags().GetInt("timeout")

	wait := newWaitFlag(cmd)
	maxWait, _ := cmd.Flags().GetInt("max-wait")

	screenshot := newScreenshotFlag(cmd)
	dom := newDOMFlag(cmd)
	force, _ := cmd.Flags().GetBool("force")
	directoryPrefix, _ := cmd.Flags().GetString("directory-prefix")

	// override wait if dom or screenshot flag is set
	wait = wait || screenshot || dom

	client, err := utils.NewAPIClient()
	if err != nil {
		return nil, err
	}

	return &scanner{
		client: client,
		batchOpts: []api.BatchOption{
			api.WithBatchMaxConcurrency(maxConcurrency),
			api.WithBatchTimeout(timeout),
		},
		wait:            wait,
		maxWait:         maxWait,
		dom:             dom,
		screenshot:      screenshot,
		force:           force,
		directoryPrefix: directoryPrefix,
		ctx:             cmd.Context(),
	}, nil
}

func getURL(json map[string]any) (string, error) {
	url, ok := json["url"].(string)
	if !ok {
		return "", fmt.Errorf("url field is missing or not a string in JSON: %v", json)
	}
	return url, nil
}

var bulkSubmitCmdExample = `  urlscan scan bulk-submit <url>...
  # submit with a file containing URLs per line, space, or tab
  urlscan scan bulk-submit list_of_urls.txt
  # combine the file input and the URL input
  urlscan scan bulk-submit list_of_urls.txt <url>
  # submit with JSONL file where each line is a JSON payload
  urlscan scan bulk-submit --jsonl payloads.jsonl
  # read JSONL from stdin
  echo '{"url":"...","visibility":"public"}' | urlscan scan bulk-submit --jsonl -`

var bulkSubmitCmdLong = `Submit multiple URLs to scan in bulk.

This command allows you to submit a list of URLs for scanning in bulk. You can provide URLs via command line arguments or through a file.
Note that the URLs will be validated before submission, and only valid URLs will be processed.`

var bulkSubmitCmd = &cobra.Command{
	Use:     "bulk-submit <url>...",
	Short:   "Bulk submit URLs to scan",
	Long:    bulkSubmitCmdLong,
	Example: bulkSubmitCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonl, err := flags.GetJSONL(cmd)
		if err != nil {
			return err
		}

		scanRequests := []scanRequest{}
		if jsonl != nil {
			for _, json := range jsonl {
				url, err := getURL(json)
				if err != nil {
					return err
				}
				scanRequests = append(scanRequests, scanRequest{
					url:  url,
					opts: []api.ScanOption{api.WithScanExtra(json)},
				})
			}
		} else {
			reader := utils.NewFilteredStringReader(utils.NewMappedStringReader(utils.StringReaderFromCmdArgs(args), utils.ResolveFileOrValue), utils.ValidateNetworkIndicator)
			urls, err := utils.ReadAllFromReader(reader)
			if err != nil {
				return err
			}
			opts := newScanOptions(cmd)
			for _, url := range urls {
				scanRequests = append(scanRequests, scanRequest{
					url:  url,
					opts: opts,
				})
			}
		}

		if len(scanRequests) == 0 {
			return cmd.Usage()
		}

		scanner, err := newScanner(cmd)
		if err != nil {
			return err
		}

		return scanner.do(scanRequests)
	},
}

func init() {
	addScanFlags(bulkSubmitCmd)
	flags.AddForceFlag(bulkSubmitCmd)
	flags.AddDirectoryPrefixFlag(bulkSubmitCmd)
	flags.AddJSONLFlag(bulkSubmitCmd)

	bulkSubmitCmd.Flags().Int("max-concurrency", 5, "Maximum number of concurrent requests for batch operation")
	bulkSubmitCmd.Flags().Int("timeout", 60*30, "Timeout for the batch operation in seconds, 0 means no timeout")

	RootCmd.AddCommand(bulkSubmitCmd)
}
