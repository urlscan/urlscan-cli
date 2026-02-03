package scan

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var resultCmdExample = `  urlscan scan result <uuid>
  echo "<uuid>" | urlscan scan result -`

var resultCmd = &cobra.Command{
	Use:     "result <uuid>",
	Short:   "Get a result by UUID",
	Example: resultCmdExample,
	Annotations: map[string]string{
		"args": "exact1",
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		reader := utils.StringReaderFromCmdArgs(args)
		uuid, err := reader.ReadString()
		if err != nil {
			return err
		}

		err = utils.ValidateUUID(uuid)
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		result, err := client.NewRequest().Get(api.PrefixedPath(fmt.Sprintf("/result/%s/", uuid)))
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJSON())

		return nil
	},
}

func init() {
	RootCmd.AddCommand(resultCmd)
}
