package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/api"
)

var quotasCmd = &cobra.Command{
	Use:   "quotas",
	Short: "Get API quotas",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := NewAPIClient()
		if err != nil {
			return err
		}

		url := api.URL("/api/v1/quotas/")
		result, err := client.Get(url)
		if err != nil {
			return err
		}

		fmt.Print(string(result.Raw))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(quotasCmd)
}
