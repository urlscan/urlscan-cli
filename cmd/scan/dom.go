package scan

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/cmd/flags"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var domCmdExample = `  urlscan scan dom <uuid>
  echo "<uuid>" | urlscan scan dom -`

var domCmd = &cobra.Command{
	Use:     "dom <uuid>",
	Short:   "Download a dom by UUID",
	Example: domCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		output, _ := cmd.Flags().GetString("output")
		force, _ := cmd.Flags().GetBool("force")

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
			output = fmt.Sprintf("%s.html", uuid)
		}
		opts := utils.NewDownloadOptions(
			utils.WithDownloadClient(client),
			utils.WithDownloadDOM(uuid),
			utils.WithDownloadOutput(output),
			utils.WithDownloadForce(force),
		)
		err = utils.DownloadWithSpinner(opts)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	flags.AddOutputFlag(domCmd, "<uuid>.html")
	flags.AddForceFlag(domCmd)

	RootCmd.AddCommand(domCmd)
}
