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
}

func mapCmdToScanOptions(cmd *cobra.Command) (opts []api.ScanOption, err error) {
	country, _ := cmd.Flags().GetString("country")
	customAgent, _ := cmd.Flags().GetString("customagent")
	overrideSafety, _ := cmd.Flags().GetString("overrideSafety")
	referer, _ := cmd.Flags().GetString("referer")
	tags, _ := cmd.Flags().GetStringArray("tags")
	visibility, _ := cmd.Flags().GetString("visibility")

	opts = append(opts,
		api.WithScanCountry(country),
		api.WithScanCustomAgent(customAgent),
		api.WithScanOverrideSafety(overrideSafety),
		api.WithScanReferer(referer),
		api.WithScanTags(tags),
		api.WithScanVisibility(visibility),
	)

	return opts, nil
}
