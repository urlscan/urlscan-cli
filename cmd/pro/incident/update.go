package incident

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var updateCmdExample = `  urlscan pro incident update <incident-id> -o <observable>
  echo <incident-id> | urlscan pro incident update - -o <observable>
  urlscan pro incident update <incident-id> --json '{"incident":{"observable":"..."}}'`

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update an incident",
	Example: updateCmdExample,
	Annotations: map[string]string{
		"args": "exact1",
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		reader := utils.StringReaderFromCmdArgs(args)
		id, err := reader.ReadString()
		if err != nil {
			return err
		}

		err = utils.ValidateULID(id)
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		opts := []api.IncidentOption{}

		json, err := flags.GetJSON(cmd)
		if err != nil {
			return err
		}
		if json != nil {
			opts = append(opts, api.WithIncidentExtra(json))
		} else {
			mapped, err := mapCmdToIncidentOptions(cmd)
			if err != nil {
				return err
			}
			opts = append(opts, mapped...)
		}

		result, err := client.UpdateIncident(id, opts...)
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJSON())

		return nil
	},
}

func init() {
	setCreateOrUpdateFlags(updateCmd)
	flags.AddJSONFlag(updateCmd)

	RootCmd.AddCommand(updateCmd)
}
