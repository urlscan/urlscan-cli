package incident

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var createCmdExample = `  urlscan pro incident create -o <observable>
  urlscan pro incident create --json '{"incident":{"observable":"..."}}'`

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new incident",
	Example: createCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		result, err := client.CreateIncident(opts...)
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJSON())

		return nil
	},
}

func init() {
	setCreateOrUpdateFlags(createCmd)
	flags.AddJSONFlag(createCmd)

	RootCmd.AddCommand(createCmd)
}
