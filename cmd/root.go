package cmd

import (
	"os"

	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	api "github.com/urlscan/urlscan-cli/api"

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

func NewAPIClient() (*utils.APIClient, error) {
	return utils.NewAPIClient(fmt.Sprintf("urlscan-cli %s", Version))
}

var rootCmd = &cobra.Command{
	Use:   "urlscan",
	Short: "A CLI tool for interacting with urlscan.io",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	addHostFlag(rootCmd.PersistentFlags())
	addProxyFlag(rootCmd.PersistentFlags())
}
