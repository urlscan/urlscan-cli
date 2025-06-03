package scan

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var submitCmdExample = `  urlscan scan submit <url>
  echo "<url>" | urlscan scan submit -`

var submitCmd = &cobra.Command{
	Use:     "submit",
	Short:   "Submit a URL to scan",
	Example: submitCmdExample,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		country, _ := cmd.Flags().GetString("country")
		customAgent, _ := cmd.Flags().GetString("customagent")
		overrideSafety, _ := cmd.Flags().GetString("overrideSafety")
		referer, _ := cmd.Flags().GetString("referer")
		tags, _ := cmd.Flags().GetStringArray("tags")
		visibility, _ := cmd.Flags().GetString("visibility")

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

		res, err := client.Scan(url,
			api.WithScanCountry(country),
			api.WithScanCustomAgent(customAgent),
			api.WithScanOverrideSafety(overrideSafety),
			api.WithScanReferer(referer),
			api.WithScanTags(tags),
			api.WithScanVisibility(visibility),
		)
		if err != nil {
			return err
		}

		if !wait {
			fmt.Print(string(res.Raw))
			return nil
		}

		ctx := cmd.Context()
		result, err := client.WaitAndGetResult(ctx, res.UUID, maxWait)
		if err != nil {
			return err
		}

		fmt.Print(string(result.Raw))

		return nil
	},
}

func init() {
	submitCmd.Flags().StringArrayP("tags", "t", []string{}, "User-defined tags to annotate this scan")
	submitCmd.Flags().StringP("country", "c", "", " Specify which country the scan should be performed from (2-Letter ISO-3166-1 alpha-2 country")
	submitCmd.Flags().StringP("customagent", "a", "", "Override User-Agent for this scan")
	submitCmd.Flags().StringP("overrideSafety", "o", "", "If set to any value, this will disable reclassification of URLs with potential PII in them")
	submitCmd.Flags().StringP("referer", "r", "", "Override HTTP referer for this scan")
	submitCmd.Flags().StringP("visibility", "v", "", "One of public, unlisted, private")
	submitCmd.Flags().BoolP("wait", "w", false, "Wait for the scan to finish")
	submitCmd.Flags().IntP("max-wait", "m", 60, "Maximum wait time in seconds")

	RootCmd.AddCommand(submitCmd)
}
