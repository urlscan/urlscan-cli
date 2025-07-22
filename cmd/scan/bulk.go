package scan

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

type scanner struct {
	client    *utils.APIClient
	scanOpts  []api.ScanOption
	batchOpts []api.BatchOption
	wait      bool
	maxWait   int
	ctx       context.Context
}

func (s *scanner) do(urls []string) error {
	tasks := make([]api.BatchTask[*json.RawMessage], len(urls))
	for i, url := range urls {
		if s.wait {
			tasks[i] = s.client.NewBatchScanWithWaitTask(url, s.maxWait, s.scanOpts...)
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
	scanOpts, err := mapCmdToScanOptions(cmd)
	if err != nil {
		return nil, err
	}

	wait, _ := cmd.Flags().GetBool("wait")
	maxWait, _ := cmd.Flags().GetInt("max-wait")

	maxConcurrency, _ := cmd.Flags().GetInt("max-concurrency")
	timeout, _ := cmd.Flags().GetInt("timeout")

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
		wait:    wait,
		maxWait: maxWait,
		ctx:     cmd.Context(),
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

		reader := utils.NewFilteredStringReader(utils.NewMappedStringReader(utils.StringReaderFromCmdArgs(args), utils.ResolveFileOrValue), utils.ValidateURL)
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

	bulkSubmitCmd.Flags().Int("max-concurrency", 5, "Maximum number of concurrent requests for batch operation")
	bulkSubmitCmd.Flags().Int("timeout", 60*30, "Timeout for the batch operation in seconds, 0 means no timeout")

	RootCmd.AddCommand(bulkSubmitCmd)
}
