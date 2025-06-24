package pro

import (
	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/cmd/pro/brand"
	"github.com/urlscan/urlscan-cli/cmd/pro/incident"
	"github.com/urlscan/urlscan-cli/cmd/pro/subscription"
)

var RootCmd = &cobra.Command{
	Use:   "pro",
	Short: "Pro sub-commands",
}

func init() {
	RootCmd.AddCommand(brand.RootCmd)
	RootCmd.AddCommand(subscription.RootCmd)
	RootCmd.AddCommand(incident.RootCmd)
}
