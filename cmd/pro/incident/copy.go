package incident

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var copyCmdExample = `  urlscan pro incident copy <incident-id>
  echo <incident-id> | urlscan pro incident copy -`

var copyCmd = &cobra.Command{
	Use:     "copy",
	Short:   "Copy an incident",
	Example: copyCmdExample,
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

		result, err := client.NewRequest().Post(api.PrefixedPath(fmt.Sprintf("/user/incidents/%s/copy", id)))
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJSON())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(copyCmd)
}
