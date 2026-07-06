package search

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var rootCmdExample = `  urlscan search <query>
  echo "<query>" | urlscan search -`

var rootCmdLong = `Search by a query.

A query uses Elasticsearch query string syntax, For example, page.domain:example.com AND date:>now-1h.

To discover which fields you can use in a query, use "search fields".

To only count the number of results, use "search count <query>".`

var RootCmd = &cobra.Command{
	Use:     "search <query>",
	Short:   "Search by a query",
	Long:    rootCmdLong,
	Example: rootCmdExample,
	Annotations: map[string]string{
		"args": "exact1",
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		limit, _ := cmd.Flags().GetInt("limit")
		all, _ := cmd.Flags().GetBool("all")
		size, _ := cmd.Flags().GetInt("size")
		if limit == 0 {
			limit = size
		}
		searchAfter, _ := cmd.Flags().GetString("search-after")
		datasource, _ := cmd.Flags().GetString("datasource")
		collapse, _ := cmd.Flags().GetString("collapse")

		reader := utils.StringReaderFromCmdArgs(args)
		q, err := reader.ReadString()
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}
		it, err := client.Search(q,
			api.IteratorSize(size),
			api.IteratorLimit(limit),
			api.IteratorSearchAfter(searchAfter),
			api.IteratorAll(all),
			api.IteratorDatasource(datasource),
			api.IteratorCollapse(collapse),
		)
		if err != nil {
			return err
		}

		results := utils.NewSearchResults()
		for result, err := range it.Iterate() {
			if err != nil {
				return err
			}
			results.Results = append(results.Results, result.Raw)
		}

		results.HasMore = it.HasMore
		results.Total = it.Total

		b, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return err
		}

		fmt.Print(string(b))

		return nil
	},
}

func init() {
	flags.AddSizeFlag(RootCmd, 100) // non-pro user's max size is 100
	flags.AddLimitFlag(RootCmd)
	flags.AddAllFlag(RootCmd)
	RootCmd.Flags().String("search-after", "", "For retrieving the next batch of results, value of the sort attribute of the last (oldest) result you received (comma-separated)")
	RootCmd.Flags().StringP("datasource", "D", "scans", "Datasources to search: scans (urlscan.io), hostnames, incidents, notifications, certificates (urlscan Pro)")
	RootCmd.Flags().StringP("collapse", "c", "", "Field to collapse results on")
}
