package channel

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var updateCmdExample = `  urlscan pro channel update <channel-id> -n <name>
  echo <channel-id> | urlscan pro channel update - -n <name>
  urlscan pro channel update <channel-id> --json '{"channel":{"name":"..."}}'`

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a channel",
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

		err = utils.ValidateUUID(id)
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		opts := []api.ChannelOption{}

		json, err := flags.GetJSON(cmd)
		if err != nil {
			return err
		}
		if json != nil {
			opts = append(opts, api.WithChannelExtra(json))
		} else {
			mapped := mapCmdToChannelOptions(cmd)
			opts = append(opts, mapped...)
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
	flags.AddJSONFlag(updateCmd)

	RootCmd.AddCommand(updateCmd)
}
