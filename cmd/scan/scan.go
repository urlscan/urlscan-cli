package scan

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var scanCmdExample = `  urlscan scan <url>
  echo "<url>" | urlscan scan -`

var scanCmd = &cobra.Command{
	Use:     "scan",
	Short:   "Scan a URL",
	Example: scanCmdExample,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		country, _ := cmd.Flags().GetString("country")
		customAgent, _ := cmd.Flags().GetString("customagent")
		overrideSafety, _ := cmd.Flags().GetString("overrideSafety")
		referer, _ := cmd.Flags().GetString("referer")
		tags, _ := cmd.Flags().GetStringArray("tags")
		visibility, _ := cmd.Flags().GetString("visibility")

		wait, _ := cmd.Flags().GetBool("wait")
		waitDeadline, _ := cmd.Flags().GetInt("wait-deadline")

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
		result, err := client.WaitAndGetResult(ctx, res.UUID, waitDeadline)
		if err != nil {
			return err
		}

		fmt.Print(string(result.Raw))

		return nil
	},
}

func init() {
	scanCmd.Flags().StringArrayP("tags", "t", []string{}, "User-defined tags to annotate this scan.")
	scanCmd.Flags().StringP("country", "c", "", " Specify which country the scan should be performed from (2-Letter ISO-3166-1 alpha-2 country")
	scanCmd.Flags().StringP("customagent", "a", "", "Override User-Agent for this scan")
	scanCmd.Flags().StringP("overrideSafety", "o", "", " If set to any value, this will disable reclassification of URLs with potential PII in them")
	scanCmd.Flags().StringP("referer", "r", "", "Override HTTP referer for this scan")
	scanCmd.Flags().StringP("visibility", "v", "", "One of public, unlisted, private")
	scanCmd.Flags().BoolP("wait", "w", false, "Wait for the scan to finish")
	scanCmd.Flags().IntP("wait-deadline", "d", 60, "Maximum waiting timeout in seconds")

	RootCmd.AddCommand(scanCmd)
}
