package incident

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var closeCmdExample = `  urlscan pro incident close <incident-id>
  echo <incident-id> | urlscan pro incident close -`

var closeCmd = &cobra.Command{
	Use:     "close",
	Short:   "Close an incident",
	Example: closeCmdExample,
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

		url := api.URL("/api/v1/user/incidents/%s/close", id)
		res, err := client.Put(url, &api.Request{})
		if err != nil {
			return err
		}

		fmt.Print(string(res.Raw))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(closeCmd)
}
