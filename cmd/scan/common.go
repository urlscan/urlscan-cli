package scan

import (
	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
)

func addScanFlags(cmd *cobra.Command) {
	flags.AddTagsFlag(cmd)
	flags.AddCountryFlag(cmd)
	flags.AddCustomAgentFlag(cmd)
	flags.AddOverrideSafetyFlag(cmd)
	flags.AddRefererFlag(cmd)
	flags.AddVisibilityFlag(cmd)
	flags.AddWaitFlag(cmd)
	flags.AddMaxWaitFlag(cmd)

	flags.AddScreenshotFlag(cmd)
	flags.AddDOMFlag(cmd)
	flags.AddDownloadFlag(cmd)
}

func newScanOptions(cmd *cobra.Command) (opts []api.ScanOption) {
	country, _ := cmd.Flags().GetString("country")
	customAgent, _ := cmd.Flags().GetString("customagent")
	overrideSafety, _ := cmd.Flags().GetString("overrideSafety")
	referer, _ := cmd.Flags().GetString("referer")
	tags, _ := cmd.Flags().GetStringArray("tags")
	visibility, _ := cmd.Flags().GetString("visibility")

	return append(opts,
		api.WithScanCountry(country),
		api.WithScanCustomAgent(customAgent),
		api.WithScanOverrideSafety(overrideSafety),
		api.WithScanReferer(referer),
		api.WithScanTags(tags),
		api.WithScanVisibility(visibility),
	)
}

func newScreenshotFlag(cmd *cobra.Command) bool {
	download, _ := cmd.Flags().GetBool("download")
	screenshot, _ := cmd.Flags().GetBool("screenshot")
	return screenshot || download
}

func newDOMFlag(cmd *cobra.Command) bool {
	download, _ := cmd.Flags().GetBool("download")
	dom, _ := cmd.Flags().GetBool("dom")
	return dom || download
}

func newWaitFlag(cmd *cobra.Command) bool {
	wait, _ := cmd.Flags().GetBool("wait")
	screenshot := newScreenshotFlag(cmd)
	dom := newDOMFlag(cmd)
	return wait || screenshot || dom
}
