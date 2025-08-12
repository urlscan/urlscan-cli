package incident

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var createCmdExample = `  urlscan pro incident create -o <observable>`

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new incident",
	Example: createCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts, err := mapCmdToIncidentOptions(cmd)
		if err != nil {
			return cmd.Usage()
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		resp, err := client.CreateIncident(opts...)
		if err != nil {
			return err
		}

		fmt.Print(resp.PrettyJSON())

		return nil
	},
}

func init() {
	setCreateOrUpdateFlags(createCmd)

	RootCmd.AddCommand(createCmd)
}
