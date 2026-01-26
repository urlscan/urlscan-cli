package scan

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var screenshotCmdExample = `  urlscan scan screenshot <uuid>
  echo "<uuid>" | urlscan scan screenshot -`

var screenshotCmd = &cobra.Command{
	Use:     "screenshot <uuid>",
	Short:   "Download a screenshot by UUID",
	Example: screenshotCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		output, _ := cmd.Flags().GetString("output")
		force, _ := cmd.Flags().GetBool("force")
		directoryPrefix, _ := cmd.Flags().GetString("directory-prefix")

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

		if output == "" {
			output = fmt.Sprintf("%s.png", uuid)
		}
		opts := utils.NewDownloadOptions(
			utils.WithDownloadClient(client),
			utils.WithDownloadScreenshot(uuid),
			utils.WithDownloadOutput(output),
			utils.WithDownloadForce(force),
			utils.WithDownloadDirectoryPrefix(directoryPrefix),
		)
		err = utils.DownloadWithSpinner(opts)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	flags.AddOutputFlag(screenshotCmd, "<uuid>.png")
	flags.AddForceFlag(screenshotCmd)
	flags.AddDirectoryPrefixFlag(screenshotCmd)

	RootCmd.AddCommand(screenshotCmd)
}
