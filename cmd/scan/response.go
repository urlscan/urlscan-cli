package scan

import (
	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/cmd/flags"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var responseCmdExample = `  urlscan scan response <file-hash>
  echo "<file-hash>" | urlscan scan response -`

var responseCmd = &cobra.Command{
	Use:     "response <flile-hash>",
	Short:   "Get a response by SHA256 file hash",
	Example: responseCmdExample,
	Annotations: map[string]string{
		"args": "exact1",
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		output, _ := cmd.Flags().GetString("output")
		force, _ := cmd.Flags().GetBool("force")
		directoryPrefix, _ := cmd.Flags().GetString("directory-prefix")

		reader := utils.StringReaderFromCmdArgs(args)
		fileHash, err := reader.ReadString()
		if err != nil {
			return err
		}

		err = utils.ValidateSHA256(fileHash)
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		if output == "" {
			output = fileHash
		}
		opts := utils.NewDownloadOptions(
			utils.WithDownloadClient(client),
			utils.WithDownloadResponse(fileHash),
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
	flags.AddOutputFlag(responseCmd, "<file-hash>")
	flags.AddForceFlag(responseCmd)
	flags.AddDirectoryPrefixFlag(responseCmd)

	RootCmd.AddCommand(responseCmd)
}
