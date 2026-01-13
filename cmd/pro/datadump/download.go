package datadump

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var DownloadCmdExample = `  urlscan pro datadump download days/api/20260101.gz
  urlscan pro datadump download hours/api/20260101/20260101-01.gz
  echo "<path>" | urlscan pro datadump download -`

var downloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "Download the data dump file",
	Example: DownloadCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		output, _ := cmd.Flags().GetString("output")
		force, _ := cmd.Flags().GetBool("force")
		extract, _ := cmd.Flags().GetBool("extract")

		reader := utils.StringReaderFromCmdArgs(args)
		path, err := reader.ReadString()
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}
		// disable auto gzip decompression for streamline extraction process
		client.SetDisableCompression(true)

		if output == "" {
			output = filepath.Base(path)
		}
		opts := utils.NewDownloadOptions(
			utils.WithDownloadClient(client),
			utils.WithDownloadOutput(output),
			utils.WithDownloadForce(force),
			utils.WithDownloadURL(api.PrefixedPath(fmt.Sprintf("/datadump/link/%s", path))),
		)
		err = utils.DownloadWithSpinner(opts)
		if err != nil {
			return err
		}

		if extract {
			err = utils.Extract(output)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	flags.AddOutputFlag(downloadCmd, "<path>.gz")
	flags.AddForceFlag(downloadCmd)

	downloadCmd.Flags().BoolP("extract", "e", false, "Extract the downloaded file")

	RootCmd.AddCommand(downloadCmd)
}
