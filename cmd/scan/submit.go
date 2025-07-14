package scan

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
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

func (s *scanner) doSingle(url string) error {
	scanResult, err := s.client.Scan(url, s.scanOpts...)
	if err != nil {
		return err
	}
	if !s.wait {
		fmt.Print(scanResult.PrettyJson())
		return nil
	}

	waitResult, err := s.client.WaitAndGetResult(s.ctx, scanResult.UUID, s.maxWait)
	if err != nil {
		return err
	}
	fmt.Print(waitResult.PrettyJson())
	return nil
}

func (s *scanner) doBatch(urls []string) error {
	tasks := make([]api.BatchTask[*json.RawMessage], len(urls))
	for i, url := range urls {
		if s.wait {
			tasks[i] = s.client.NewBatchScanWitWaitTask(url, s.maxWait, s.scanOpts...)
		} else {
			tasks[i] = s.client.NewBatchScanTask(url, s.scanOpts...)
		}
	}

	results, err := api.Batch(s.client.Client, tasks, s.batchOpts...)
	if err != nil {
		return err
	}

	pairs := utils.NewBatchJsonResultPairs(urls, results)

	b, err := json.MarshalIndent(pairs, "", "  ")
	if err != nil {
		return err
	}

	fmt.Print(string(b))

	return nil
}

func (s *scanner) do(urls []string) error {
	if len(urls) == 1 {
		return s.doSingle(urls[0])
	}
	return s.doBatch(urls)

}

func newScanner(cmd *cobra.Command) (*scanner, error) {
	country, _ := cmd.Flags().GetString("country")
	customAgent, _ := cmd.Flags().GetString("customagent")
	overrideSafety, _ := cmd.Flags().GetString("overrideSafety")
	referer, _ := cmd.Flags().GetString("referer")
	tags, _ := cmd.Flags().GetStringArray("tags")
	visibility, _ := cmd.Flags().GetString("visibility")
	wait, _ := cmd.Flags().GetBool("wait")
	maxWait, _ := cmd.Flags().GetInt("max-wait")

	maxConcurrency, _ := cmd.Flags().GetInt("max-concurrency")
	totalTimeout, _ := cmd.Flags().GetInt("total-timeout")

	client, err := utils.NewAPIClient()
	if err != nil {
		return nil, err
	}

	return &scanner{
		client: client,
		scanOpts: []api.ScanOption{
			api.WithScanCountry(country),
			api.WithScanCustomAgent(customAgent),
			api.WithScanOverrideSafety(overrideSafety),
			api.WithScanReferer(referer),
			api.WithScanVisibility(visibility),
			api.WithScanTags(tags),
		},
		batchOpts: []api.BatchOption{
			api.WithBatchMaxConcurrency(maxConcurrency),
			api.WithBatchTotalTimeout(totalTimeout),
		},
		wait:    wait,
		maxWait: maxWait,
		ctx:     cmd.Context(),
	}, nil
}

var submitCmdExample = `  urlscan scan submit <url>...
  echo "<url>" | urlscan scan submit -
  # submit with a file containing URLs per line, space, or tab
  urlscan scan submit list_of_urls.txt
  # combine the file input and the URL input
  urlscan scan submit list_of_urls.txt <url>`

var submitCmd = &cobra.Command{
	Use:     "submit <url>",
	Short:   "Submit a URL to scan",
	Example: submitCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Usage()
		}

		reader := utils.NewMappedStringReader(utils.StringReaderFromCmdArgs(args), utils.ResolveFileOrValue)
		urls, err := reader.ReadAll()
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
	submitCmd.Flags().StringArrayP("tags", "t", []string{}, "User-defined tags to annotate this scan")
	submitCmd.Flags().StringP("country", "c", "", " Specify which country the scan should be performed from (2-Letter ISO-3166-1 alpha-2 country")
	submitCmd.Flags().StringP("customagent", "a", "", "Override User-Agent for this scan")
	submitCmd.Flags().StringP("overrideSafety", "o", "", "If set to any value, this will disable reclassification of URLs with potential PII in them")
	submitCmd.Flags().StringP("referer", "r", "", "Override HTTP referer for this scan")
	submitCmd.Flags().StringP("visibility", "v", "", "One of public, unlisted, private")
	submitCmd.Flags().BoolP("wait", "w", false, "Wait for the scan to finish")
	submitCmd.Flags().IntP("max-wait", "m", 60, "Maximum wait time in seconds")

	flags.AddMaxConcurrencyFlag(submitCmd)
	flags.AddTotalTimeoutFlag(submitCmd)

	RootCmd.AddCommand(submitCmd)
}
