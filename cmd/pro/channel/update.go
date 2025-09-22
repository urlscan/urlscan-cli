package channel

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var updateCmdExample = `  urlscan pro channel update <channel-id> -n <name>
  echo <channel-id> | urlscan pro channel update - -n <name>`

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a channel",
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

		err = utils.ValidateUUID(id)
		if err != nil {
			return err
		}

		opts, err := mapCmdToChannelOptions(cmd)
		if err != nil {
			return cmd.Usage()
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		result, err := client.UpdateChannel(id, opts...)
		if err != nil {
			return err
		}

		fmt.Print(result.PrettyJSON())

		return nil
	},
}

func init() {
	setCreateOrUpdateFlags(updateCmd)

	RootCmd.AddCommand(updateCmd)
}
