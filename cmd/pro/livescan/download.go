package livescan

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var downloadCmdExample = `  urlscan pro livescan download <file-hash> -s <scanner-id>
  echo <file-hash> | urlscan pro livescan download - -s <scanner-id>`

var downloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "Download a resource of a live scan by SHA256 file hash",
	Example: downloadCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		reader := utils.StringReaderFromCmdArgs(args)
		fileHash, err := reader.ReadString()
		if err != nil {
			return err
		}
		err = utils.ValidateSHA256(fileHash)
		if err != nil {
			return err
		}

		scannerId, _ := cmd.Flags().GetString("scanner-id")
		if scannerId == "" {
			return cmd.Usage()
		}

		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = fileHash
		}
		force, _ := cmd.Flags().GetBool("force")

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		url := api.PrefixedPath(fmt.Sprintf("/livescan/%s/download/%s/", scannerId, fileHash))
		opts := utils.NewDownloadOptions(
			utils.WithDownloadClient(client),
			utils.WithDownloadURL(url),
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
	addScannerIdFlag(downloadCmd)
	flags.AddForceFlag(downloadCmd)
	flags.AddOutputFlag(downloadCmd, "<file-hash>")

	RootCmd.AddCommand(downloadCmd)
}
