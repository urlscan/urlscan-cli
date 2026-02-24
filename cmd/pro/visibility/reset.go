package visibility

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var resetCmdExample = `  urlscan pro visibility reset <scan-id>
  echo <scan-id> | urlscan pro visibility reset -`

var resetCmd = &cobra.Command{
	Use:     "reset",
	Short:   "Reset the visibility of a scan owned by you or your team to its original visibility",
	Example: resetCmdExample,
	Annotations: map[string]string{
		"args": "exact1",
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		reader := utils.StringReaderFromCmdArgs(args)
		scanId, err := reader.ReadString()
		if err != nil {
			return err
		}

		err = utils.ValidateUUID(scanId)
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		resp, err := client.NewRequest().Delete(api.PrefixedPath(fmt.Sprintf("/result/%s/visibility/", scanId)))
		if err != nil {
			return err
		}

		fmt.Print(resp.PrettyJSON())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(resetCmd)
}
