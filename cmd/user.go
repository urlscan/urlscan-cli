package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	api "github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var userCmdExample = `  urlscan user`

var userCmd = &cobra.Command{
	Use:     "user",
	Short:   "Get information about the current user or API key making the request",
	Example: userCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		url := api.URL("/user/username")
		result, err := client.Get(url)
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJson())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(userCmd)
}
