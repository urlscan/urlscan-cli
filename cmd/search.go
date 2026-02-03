package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var searchCmdExample = `  urlscan search <query>
  echo "<query>" | urlscan search -`

var searchCmd = &cobra.Command{
	Use:     "search <query>",
	Short:   "Search by a query",
	Example: searchCmdExample,
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
	flags.AddSizeFlag(searchCmd, 100) // non-pro user's max size is 100
	flags.AddLimitFlag(searchCmd)
	flags.AddAllFlag(searchCmd)
	searchCmd.Flags().String("search-after", "", "For retrieving the next batch of results, value of the sort attribute of the last (oldest) result you received (comma-separated)")
	searchCmd.Flags().StringP("datasource", "D", "scans", "Datasources to search: scans (urlscan.io), hostnames, incidents, notifications, certificates (urlscan Pro)")
	searchCmd.Flags().StringP("collapse", "c", "", "Field to collapse results on")

	RootCmd.AddCommand(searchCmd)
}
