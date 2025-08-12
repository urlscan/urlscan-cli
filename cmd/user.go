package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

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

		resp, err := client.NewRequest().Get("/user/username")
		if err != nil {
			return err
		}

		fmt.Print(resp.PrettyJSON())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(userCmd)
}
