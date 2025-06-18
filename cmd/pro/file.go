package pro

import (
	"fmt"

	"github.com/spf13/cobra"

	api "github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var fileCmdExample = `  urlscan pro file <file-hash>
  echo "<file-hash>" | urlscan pro file -`

var fileCmd = &cobra.Command{
	Use:     "file",
	Short:   "Download a file",
	Example: fileCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		reader := utils.StringReaderFromCmdArgs(args)
		hash, err := reader.ReadString()
		if err != nil {
			return err
		}

		filename, _ := cmd.Flags().GetString("filename")
		if filename == "" {
			filename = fmt.Sprintf("%s.zip", hash)
		}
		password, _ := cmd.Flags().GetString("password")
		force, _ := cmd.Flags().GetBool("force")

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		url := api.URL("/downloads/%s", hash)
		query := url.Query()
		query.Add("filename", filename)
		query.Add("password", password)
		url.RawQuery = query.Encode()

		options := utils.NewDownloadOptions(client, url, filename, force)
		err = utils.Download(options)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	fileCmd.Flags().StringP("filename", "F", "", "Specify the name of the ZIP file that should be downloaded (defaults to <hash>.zip)")
	fileCmd.Flags().StringP("password", "p", "urlscan!", "The password to use to encrypt the ZIP file")
	fileCmd.Flags().BoolP("force", "f", false, "Enable to force overwriting an existing file")

	RootCmd.AddCommand(fileCmd)
}
