package scan

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var submitCmdExample = `  urlscan scan submit <url>...
  echo "<url>" | urlscan scan submit -`

var submitCmd = &cobra.Command{
	Use:     "submit <url>",
	Short:   "Submit a URL to scan",
	Example: submitCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		scanOpts, err := mapCmdToScanOptions(cmd)
		if err != nil {
			return err
		}

		wait, _ := cmd.Flags().GetBool("wait")
		maxWait, _ := cmd.Flags().GetInt("max-wait")

		reader := utils.StringReaderFromCmdArgs(args)
		url, err := reader.ReadString()
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		scanResult, err := client.Scan(url, scanOpts...)
		if err != nil {
			return err
		}

		if !wait {
			fmt.Print(scanResult.PrettyJson())
			return nil
		}

		ctx := cmd.Context()
		waitResult, err := client.WaitAndGetResult(ctx, scanResult.UUID, maxWait)
		if err != nil {
			return err
		}

		fmt.Print(waitResult.PrettyJson())

		return nil
	},
}

func init() {
	addScanFlags(submitCmd)

	RootCmd.AddCommand(submitCmd)
}
