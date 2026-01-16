package flags

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/api"
)

func AddForceFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("force", "f", false, "Force overwrite an existing file")
}

func AddAllFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("all", false, "Return all results; limit is ignored if --all is specified (default false)")
}

func AddLimitFlag(cmd *cobra.Command) {
	cmd.Flags().IntP("limit", "l", api.MaxTotal, "Maximum number of results that will be returned by the iterator")
}

func AddSizeFlag(cmd *cobra.Command, value int) {
	cmd.Flags().IntP("size", "s", value, "Number of results returned by the iterator in each batch")
}

func AddOutputFlag(cmd *cobra.Command, defaultExample string) {
	cmd.Flags().StringP("output", "o", "", fmt.Sprintf("Output file name (default %s)", defaultExample))
}

func AddTagsFlag(cmd *cobra.Command) {
	cmd.Flags().StringArrayP("tags", "t", []string{}, "User-defined tags to annotate this scan")
}

func AddCountryFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("country", "c", "", "Specify which country the scan should be performed from (2-Letter ISO-3166-1 alpha-2 country")
}

func AddCustomAgentFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("customagent", "a", "", "Override User-Agent for this scan")
}

func AddOverrideSafetyFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("overrideSafety", "o", "", "If set to any value, this will disable reclassification of URLs with potential PII in them")
}

func AddRefererFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("referer", "r", "", "Override HTTP referer for this scan")
}

func AddVisibilityFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("visibility", "v", "", "One of public, unlisted, private")
}

func AddWaitFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("wait", "w", false, "Wait for the scan(s) to finish")
}

func AddMaxWaitFlag(cmd *cobra.Command) {
	cmd.Flags().IntP("max-wait", "m", 60, "Maximum wait time per scan in seconds")
}

func AddScreenshotFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("screenshot", false, "Download only the screenshot (overrides wait)")
}

func AddDOMFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("dom", false, "Download only the DOM contents (overrides wait)")
}

func AddDownloadFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("download", false, "Download screenshot and DOM contents (overrides wait/dom/screenshot)")
}

func AddDirectoryPrefixFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("directory-prefix", "P", ".", "Set directory prefix where file will be saved")
}
