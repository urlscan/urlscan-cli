package channel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/api"
)

func setCreateOrUpdateFlags(cmd *cobra.Command) {
	// required flags
	cmd.Flags().StringP("name", "n", "", "Channel name (required)")

	// optional flags
	cmd.Flags().String("type", "webhook", "Type of channel (webhook or email)")
	cmd.Flags().String("webhook-url", "", "Webhook URL (optional)")
	cmd.Flags().String("frequency", "", "Frequency of notifications (live, hourly or daily) (optional)")
	cmd.Flags().String("utc-time", "", "24 hour UTC time that daily emails are sent (optional)")
	cmd.Flags().Bool("is-active", true, "Whether the channel is active")
	cmd.Flags().Bool("is-default", false, "Whether the channel is the default channel (default false)")
	cmd.Flags().Bool("ignore-time", false, "Whether to ignore time constraints (default false)")
	cmd.Flags().StringSlice("email-addresses", []string{}, "Email addresses receiving the notifications")
	cmd.Flags().StringSlice("week-days", []string{}, "Days of the week alerts will be generated (Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday)")
	cmd.Flags().StringSlice("permissions", []string{}, "Permissions")
}

func mapCmdToChannelOptions(cmd *cobra.Command) (opts []api.ChannelOption, err error) {
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

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
