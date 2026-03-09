package channel

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var createCmdExample = `  urlscan pro channel create -n <name>
  urlscan pro channel create --json '{"channel":{"name":"...","type":"webhook","webhookURL":"https://..."}}'`

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new channel",
	Example: createCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		result, err := client.CreateChannel(opts...)
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
