package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

type HostnameResults struct {
	Results   []json.RawMessage `json:"results"`
	PageState string            `json:"pageState"`
	HasMore   bool              `json:"has_more"`
}

func newHostnameResults() HostnameResults {
	return HostnameResults{
		Results:   make([]json.RawMessage, 0),
		PageState: "",
	}
}

var hostnameCmdExample = `  urlscan hostname <hostname>
  echo "<hostname>" | urlscan hostname -`

var hostnameLong = `To have the same idiom with the search command, this command has the following specs:

- Request:
  - limit: is not exactly same as the API endpoint's "limit" query parameter. It is the maximum number of results that will be returned by the iterator.
	- size: is equivalent to the API endpoint's "limit" query parameter, which is the number of results returned by the iterator in each batch.
- Response:
  - hasMore: is an additional field (not included in the API endpoint response) indicates if there are more results available.
`

var hostnameCmd = &cobra.Command{
	Use:     "hostname",
	Short:   "Get the historical observations for a specific hostname in the hostname data source",
	Long:    hostnameLong,
	Example: hostnameCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		size, _ := cmd.Flags().GetInt("size")
		limit, _ := cmd.Flags().GetInt("limit")
		pageState, _ := cmd.Flags().GetString("page-state")

		reader := utils.StringReaderFromCmdArgs(args)
		hostname, err := reader.ReadString()
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		it, err := client.IterateHostname(hostname,
			api.HostnameIteratorLimit(limit),
			api.HostnameIteratorSize(size),
			api.HostnameIteratorPageState(pageState),
		)
		if err != nil {
			return err
		}

		results := newHostnameResults()
		for result, err := range it.Iterate() {
			if err != nil {
				return err
			}
			results.Results = append(results.Results, *result)
		}

		results.HasMore = it.HasMore

		b, err := json.Marshal(results)
		if err != nil {
			return err
		}

		fmt.Print(string(b))

		return nil
	},
}

func init() {
	hostnameCmd.Flags().IntP("limit", "l", 10000, "Maximum number of results that will be returned by the iterator")
	hostnameCmd.Flags().IntP("size", "s", 1000, "Number of results returned by the iterator in each batch")
	hostnameCmd.Flags().StringP("page-state", "p", "", "Continue return additional results starting from this page state from the previous API call")

	RootCmd.AddCommand(hostnameCmd)
}
