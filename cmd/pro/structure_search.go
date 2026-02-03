package pro

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var structureSearchCmdExample = `  urlscan pro structure-search <uuid>
  echo "<uuid>" | urlscan pro structure-search -`

var structureSearchCmd = &cobra.Command{
	Use:     "structure-search <uuid>",
	Short:   "Get structurally similar results to a specific scan",
	Example: structureSearchCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		limit, _ := cmd.Flags().GetInt("limit")
		all, _ := cmd.Flags().GetBool("all")

		size, _ := cmd.Flags().GetInt("size")
		searchAfter, _ := cmd.Flags().GetString("search-after")
		q, _ := cmd.Flags().GetString("query")

		reader := utils.StringReaderFromCmdArgs(args)
		uuid, err := reader.ReadString()
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}
		it, err := client.StructureSearch(
			uuid,
			api.IteratorSize(size),
			api.IteratorLimit(limit),
			api.IteratorSearchAfter(searchAfter),
			api.IteratorQuery(q),
			api.IteratorAll(all),
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
	flags.AddSizeFlag(structureSearchCmd, 1_000)
	flags.AddLimitFlag(structureSearchCmd)
	flags.AddAllFlag(structureSearchCmd)
	structureSearchCmd.Flags().String("search-after", "", "For retrieving the next batch of results, value of the sort attribute of the last (oldest) result you received (comma-separated)")

	structureSearchCmd.Flags().StringP("query", "q", "", "Additional query filter")

	RootCmd.AddCommand(structureSearchCmd)
}
