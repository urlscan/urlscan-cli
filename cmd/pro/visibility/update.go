package visibility

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var updateCmdExample = `  urlscan pro visibility update <scan-id> -v private
  echo <scan-id> | urlscan pro visibility update - -v public`

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Updated visibility of the scan result owned by you or your team",
	Example: updateCmdExample,
	Annotations: map[string]string{
		"args": "exact1",
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		reader := utils.StringReaderFromCmdArgs(args)
		scanId, err := reader.ReadString()
		if err != nil {
			return err
		}

		err = utils.ValidateUUID(scanId)
		if err != nil {
			return err
		}

		visibility, _ := cmd.Flags().GetString("visibility")
		if visibility == "" {
			return cmd.Usage()
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		resp, err := client.UpdateResultVisibility(scanId, api.WithResultVisibility(visibility))
		if err != nil {
			return err
		}

		fmt.Print(resp.PrettyJSON())

		return nil
	},
}

func init() {
	updateCmd.Flags().StringP("visibility", "v", "", "The new visibility of the scan result. public, unlisted, private or deleted.")

	RootCmd.AddCommand(updateCmd)
}
