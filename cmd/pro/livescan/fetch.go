package livescan

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var fetchCmdExample = `  urlscan pro livescan fetch <file-hash> -s <scanner-id>
  echo <file-hash> | urlscan pro livescan fetch - -s <scanner-id>`

var fetchCmd = &cobra.Command{
	Use:     "fetch",
	Short:   "Fetch a response or a download of a live scan by SHA256 file hash",
	Example: fetchCmdExample,
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

		responseURL := api.PrefixedPath(fmt.Sprintf("/livescan/%s/response/%s/", scannerId, fileHash))
		downloadURL := api.PrefixedPath(fmt.Sprintf("/livescan/%s/download/%s/", scannerId, fileHash))

		resp, err := client.NewRequest().Head(downloadURL)
		var url string
		switch {
		case resp != nil && resp.Response != nil && resp.IsSuccess():
			url = downloadURL
		case resp != nil && resp.Response != nil && resp.StatusCode == http.StatusNotFound:
			url = responseURL
		default:
			if err != nil {
				return err
			}
			return fmt.Errorf("unexpected error fetching %q", downloadURL)
		}

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
	addScannerIdFlag(fetchCmd)
	flags.AddForceFlag(fetchCmd)
	flags.AddOutputFlag(fetchCmd, "<file-hash>")

	RootCmd.AddCommand(fetchCmd)
}
