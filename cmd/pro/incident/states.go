package incident

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var statesCmdExample = `  urlscan pro incident states <incident-id>
  echo <incident-id> | urlscan pro incident states -`

var statesCmd = &cobra.Command{
	Use:     "states",
	Short:   "Get states of an incident",
	Example: statesCmdExample,
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

		result, err := client.NewRequest().Get(api.PrefixedPath(fmt.Sprintf("/user/incidentstates/%s/", id)))
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJSON())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(statesCmd)
}
