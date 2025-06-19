package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	api "github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var searchCmdExample = `  urlscan search <query>
  echo "<query>" | urlscan search -`

var searchCmd = &cobra.Command{
	Use:     "search <query>",
	Short:   "Search by a query",
	Example: searchCmdExample,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		size, _ := cmd.Flags().GetInt("size")
		searchAfter, _ := cmd.Flags().GetString("search-after")

		reader := utils.StringReaderFromCmdArgs(args)
		q, err := reader.ReadString()
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}
		it, err := client.Search(q, api.IteratorSize(size), api.IteratorLimit(limit), api.IteratorSearchAfter(searchAfter))
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

		b, err := json.Marshal(results)
		if err != nil {
			return err
		}

		fmt.Print(string(b))

		return nil
	},
}

func init() {
	searchCmd.Flags().IntP("limit", "l", api.MaxTotal, "Maximum number of results that will be returned by the iterator")
	searchCmd.Flags().IntP("size", "s", 100, "Number of results returned by the iterator in each batch")
	searchCmd.Flags().String("search-after", "", "For retrieving the next batch of results, value of the sort attribute of the last (oldest) result you received (comma-separated)")

	RootCmd.AddCommand(searchCmd)
}
