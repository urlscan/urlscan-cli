package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var hostnameCmdExample = `  urlscan hostname <hostname>
  echo "<hostname>" | urlscan hostname -`

var hostnameCmd = &cobra.Command{
	Use:     "hostname",
	Short:   "Get the historical observations for a specific hostname in the hostname data source",
	Example: hostnameCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		result, err := client.GetHostname(hostname,
			api.WithHostnameLimit(limit),
			api.WithHostnamePageState(pageState),
		)
		if err != nil {
			return err
		}

		fmt.Print(string(result.Raw))

		return nil
	},
}

func init() {
	hostnameCmd.Flags().IntP("limit", "l", 1000, "Return at most this many results (Minimum 10, Maximum 10,000)")
	hostnameCmd.Flags().StringP("page-state", "p", "", "Continue return additional results starting from this page state from the previous API call")

	RootCmd.AddCommand(hostnameCmd)
}
