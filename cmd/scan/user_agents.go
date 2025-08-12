package scan

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var userAgentsCmdExample = `  urlscan scan user-agents`

var userAgentsCmd = &cobra.Command{
	Use:     "user-agents",
	Short:   "Get grouped user agents to use with the Scan API",
	Example: userAgentsCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		resp, err := client.NewRequest().Get(api.PrefixedPath("/userAgents"))
		if err != nil {
			return err
		}

		fmt.Print(resp.PrettyJSON())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(userAgentsCmd)
}
