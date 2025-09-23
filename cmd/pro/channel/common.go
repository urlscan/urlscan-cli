package channel

import (
	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/api"
)

func setCreateOrUpdateFlags(cmd *cobra.Command) {
	// required flags
	cmd.Flags().StringP("name", "n", "", "Channel name (required)")

	// optional flags
	cmd.Flags().String("type", "webhook", "Type of channel (webhook or email)")
	cmd.Flags().String("webhook-url", "", "Webhook URL (required for type: webhook)")
	cmd.Flags().String("frequency", "", "Frequency of notifications (live, hourly or daily) (optional)")
	cmd.Flags().String("utc-time", "", "24 hour UTC time that daily emails are sent at (optional)")
	cmd.Flags().Bool("is-active", true, "Set channel active")
	cmd.Flags().Bool("is-default", false, "Set channel as default (default false)")
	cmd.Flags().Bool("ignore-time", false, "Ignore time constraints (default false)")
	cmd.Flags().StringSlice("email-addresses", []string{}, "Email addresses receiving the notifications (required for type: email)")
	cmd.Flags().StringSlice("week-days", []string{}, "Days of the week alerts will be generated (Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday)")
	cmd.Flags().StringSlice("permissions", []string{}, "Permissions (optional; team:read, team:write)")
}

func mapCmdToChannelOptions(cmd *cobra.Command) (opts []api.ChannelOption, err error) {
	name, _ := cmd.Flags().GetString("name")
	channelType, _ := cmd.Flags().GetString("type")
	webhookURL, _ := cmd.Flags().GetString("webhook-url")
	frequency, _ := cmd.Flags().GetString("frequency")
	emailAddresses, _ := cmd.Flags().GetStringSlice("email-addresses")
	utcTime, _ := cmd.Flags().GetString("utc-time")
	isActive, _ := cmd.Flags().GetBool("is-active")
	isDefault, _ := cmd.Flags().GetBool("is-default")
	ignoreTime, _ := cmd.Flags().GetBool("ignore-time")
	weekDays, _ := cmd.Flags().GetStringSlice("week-days")
	permissions, _ := cmd.Flags().GetStringSlice("permissions")

	opts = append(opts,
		api.WithChannelType(channelType),
		api.WithChannelName(name),
		api.WithChannelWebhookURL(webhookURL),
		api.WithChannelFrequency(frequency),
		api.WithChannelEmailAddresses(emailAddresses),
		api.WithChannelUTCTime(utcTime),
		api.WithChannelIsActive(isActive),
		api.WithChannelIsDefault(isDefault),
		api.WithChannelIgnoreTime(ignoreTime),
		api.WithChannelWeekDays(weekDays),
		api.WithChannelPermissions(permissions),
	)
	return opts, nil
}
