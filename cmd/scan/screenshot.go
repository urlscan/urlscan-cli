package scan

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var screenshotCmdExample = `  urlscan screenshot <uuid>
  echo "<uuid>" | urlscan screenshot -`

var screenshotCmd = &cobra.Command{
	Use:     "screenshot",
	Short:   "Download a screenshot by UUID",
	Example: screenshotCmdExample,
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

		url := api.URL("%s", fmt.Sprintf("/screenshots/%s.png", uuid))
		if output == "" {
			output = fmt.Sprintf("%s.png", uuid)
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
	screenshotCmd.Flags().StringP("output", "o", "", "Output file name. Defaults to <uuid>.png.")
	screenshotCmd.Flags().BoolP("force", "f", false, "Enable to force overwriting an existing file.")

	RootCmd.AddCommand(screenshotCmd)
}
