package scan

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var domCmdExample = `  urlscan scan dom <uuid>
  echo "<uuid>" | urlscan scan dom -`

var domCmd = &cobra.Command{
	Use:     "dom <uuid>",
	Short:   "Download a dom by UUID",
	Example: domCmdExample,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		force, _ := cmd.Flags().GetBool("force")

		reader := utils.StringReaderFromCmdArgs(args)
		uuid, err := reader.ReadString()
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		url := api.URL("%s", fmt.Sprintf("/dom/%s/", uuid))
		if output == "" {
			output = fmt.Sprintf("%s.html", uuid)
		}
		options := utils.NewDownloadOptions(client, url, output, force)
		err = utils.Download(options)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(domCmd)
	domCmd.Flags().StringP("output", "o", "", "Output file name. Defaults to <uuid>.html.")
}
