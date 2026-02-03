package scan

import (
	"fmt"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var openCmdExample = `  urlscan scan open <uuid>
  echo "<uuid>" | urlscan scan open -`

var openCmd = &cobra.Command{
	Use:     "open <uuid>",
	Short:   "Open a scan result in your browser by UUID",
	Example: openCmdExample,
	Annotations: map[string]string{
		"args": "exact1",
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		reader := utils.StringReaderFromCmdArgs(args)
		uuid, err := reader.ReadString()
		if err != nil {
			return err
		}

		err = utils.ValidateUUID(uuid)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("https://urlscan.io/result/%s/", uuid)
		err = browser.OpenURL(url)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(openCmd)
}
