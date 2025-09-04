package pro

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
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
		HasMore:   false,
	}
}

var hostnameCmdExample = `  urlscan pro hostname <hostname>
  echo "<hostname>" | urlscan pro hostname -`

var hostnameLong = `To have the same idiom with the search command, this command has the following specs:

- Request:
  - limit: the maximum number of results that will be returned by the iterator.
  - size: the number of results returned by the iterator in each batch (equivalent to the API endpoint's "limit" query parameter).
- Response:
  - hasMore: indicates more results are available.`

var hostnameCmd = &cobra.Command{
	Use:     "hostname",
	Short:   "Get the historical observations for a specific hostname in the hostname data source",
	Long:    hostnameLong,
	Example: hostnameCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		size, _ := cmd.Flags().GetInt("size")
		limit, _ := cmd.Flags().GetInt("limit")
		all, _ := cmd.Flags().GetBool("all")
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
			api.HostnameIteratorAll(all),
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

		b, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return err
		}

		fmt.Print(string(b))

		return nil
	},
}

func init() {
	flags.AddSizeFlag(hostnameCmd, 1_000)
	flags.AddLimitFlag(hostnameCmd)
	flags.AddAllFlag(hostnameCmd)
	hostnameCmd.Flags().StringP("page-state", "p", "", "Returns additional results starting from this page state from the previous API call")

	RootCmd.AddCommand(hostnameCmd)
}
