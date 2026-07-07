package search

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var countCmdExample = `  urlscan search count <query>
  echo "<query>" | urlscan search count -`

var countCmd = &cobra.Command{
	Use:     "count <query>",
	Short:   "Count the number of results matching a query",
	Example: countCmdExample,
	Annotations: map[string]string{
		"args": "exact1",
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		reader := utils.StringReaderFromCmdArgs(args)
		q, err := reader.ReadString()
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		it, err := client.Search(q, api.IteratorSize(0))
		if err != nil {
			return err
		}

		for _, err := range it.Iterate() {
			if err != nil {
				return err
			}
		}

		fmt.Println(it.Total)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(countCmd)
}
