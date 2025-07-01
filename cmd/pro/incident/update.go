package incident

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var updateCmdExample = `  urlscan pro incident update <incident-id> -o <observable>
  echo <incident-id> | urlscan pro incident update - -o <observable>`

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update an incident",
	Example: updateCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		reader := utils.StringReaderFromCmdArgs(args)
		id, err := reader.ReadString()
		if err != nil {
			return err
		}

		opts, err := mapCmdToIncidentOptions(cmd)
		if err != nil {
			return cmd.Usage()
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		res, err := client.UpdateIncident(id, opts...)
		if err != nil {
			return err
		}

		fmt.Print(string(res.Raw))

		return nil
	},
}

func init() {
	setCreateOrUpdateFlags(updateCmd)

	RootCmd.AddCommand(updateCmd)
}
