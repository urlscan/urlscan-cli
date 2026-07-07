package search

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var fieldsCmdExample = `  urlscan search fields`

var fieldsCmd = &cobra.Command{
	Use:     "fields",
	Short:   "List the queryable fields available for search",
	Example: fieldsCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		resp, err := client.NewRequest().Get("/user/username")
		if err != nil {
			return err
		}

		var user struct {
			Limits struct {
				QueryableFields []string `json:"queryableFields"`
			} `json:"limits"`
		}
		if err := resp.Unmarshal(&user); err != nil {
			return err
		}

		b, err := json.MarshalIndent(user.Limits.QueryableFields, "", "  ")
		if err != nil {
			return err
		}

		fmt.Print(string(b))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(fieldsCmd)
}
