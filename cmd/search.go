package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var searchCmdExample = `  urlscan search <query>
  echo "<query>" | urlscan search -
  urlscan search --params '{"size":"50","q":"..."}'`

var searchCmd = &cobra.Command{
	Use:     "search <query>",
	Short:   "Search by a query",
	Example: searchCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		// non query params
		limit, _ := cmd.Flags().GetInt("limit")
		all, _ := cmd.Flags().GetBool("all")

		opts := []api.IteratorOption{
			api.IteratorLimit(limit),
			api.IteratorAll(all),
		}

		params, err := flags.GetParams(cmd)
		if err != nil {
			return err
		}

		var q string
		if params != nil {
			opts = append(opts, api.IteratorExtra(params))
		} else {
			size, _ := cmd.Flags().GetInt("size")
			searchAfter, _ := cmd.Flags().GetString("search-after")
			datasource, _ := cmd.Flags().GetString("datasource")
			collapse, _ := cmd.Flags().GetString("collapse")

			reader := utils.StringReaderFromCmdArgs(args)
			q, err = reader.ReadString()
			// ignore reader's io.EOF error when params is given
			if err != nil && errors.Is(err, io.EOF) && len(params) == 0 {
				return fmt.Errorf("no query given from both STDIN & params: %w", err)
			}

			opts = append(opts,
				api.IteratorSize(size),
				api.IteratorSearchAfter(searchAfter),
				api.IteratorDatasource(datasource),
				api.IteratorCollapse(collapse),
			)
		}

		it, err := client.Search(q, opts...)
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
	flags.AddParamsFlag(searchCmd)

	searchCmd.Flags().String("search-after", "", "For retrieving the next batch of results, value of the sort attribute of the last (oldest) result you received (comma-separated)")
	searchCmd.Flags().StringP("datasource", "D", "scans", "Datasources to search: scans (urlscan.io), hostnames, incidents, notifications, certificates (urlscan Pro)")
	searchCmd.Flags().StringP("collapse", "c", "", "Field to collapse results on")

	RootCmd.AddCommand(searchCmd)
}
