package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/pro"
	"github.com/urlscan/urlscan-cli/cmd/scan"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

func addHostFlag(flags *pflag.FlagSet) {
	flags.String(
		"host", "urlscan.io",
		"API host name")
	flags.MarkHidden("host") //nolint:errcheck
}

func addProxyFlag(flags *pflag.FlagSet) {
	flags.String(
		"proxy", "",
		"HTTP proxy")
	flags.MarkHidden("proxy") //nolint:errcheck
}

var RootCmd = &cobra.Command{
	Use:          "urlscan",
	Short:        "A CLI tool for interacting with urlscan.io",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// bind flags
		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		host := viper.GetString("host")
		if host != "" {
			api.SetHost(host)
		}
		proxy := viper.GetString("proxy")
		if proxy != "" {
			os.Setenv("http_proxy", proxy) //nolint:errcheck
		}

		// check API key presence
		key, err := utils.GetKey()
		if err != nil || key == "" {
			return fmt.Errorf("API key not found, please set the URLSCAN_API_KEY environment variable or set it in keyring by `urlscan key set`")
		}

		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	addHostFlag(RootCmd.PersistentFlags())
	addProxyFlag(RootCmd.PersistentFlags())

	RootCmd.AddCommand(scan.RootCmd)
	RootCmd.AddCommand(pro.RootCmd)
}
