package subscription

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/api"
)

// set common flags for create and update commands
func setCreateOrUpdateFlags(cmd *cobra.Command) {
	// required flags
	cmd.Flags().StringSliceP("search-ids", "s", []string{}, "Array of search IDs associated with this subscription (required)")
	cmd.Flags().StringP("frequency", "f", "", "Frequency of notifications (live, hourly or daily) (required)")
	cmd.Flags().StringSliceP("email-addresses", "e", []string{}, "Email addresses to send notifications to (required)")
	cmd.Flags().StringP("name", "n", "", "Name of the subscription (required)")

	// defaulted flags
	cmd.Flags().BoolP("is-active", "a", true, "Whether the subscription is active")
	cmd.Flags().BoolP("ignore-time", "t", false, "Whether to ignore time constraints (default false)")

	// optional flags
	cmd.Flags().StringP("description", "d", "", "Description of the subscription (optional)")
	cmd.Flags().StringSlice("week-days", []string{}, "Days of the week alerts will be generated (Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday)")
	cmd.Flags().StringSlice("permissions", []string{}, "Permissions (team:read, team:write)")
	cmd.Flags().StringSlice("channel-ids", []string{}, "Array of channel IDs associated with this subscription")
	cmd.Flags().StringSlice("incident-channel-ids", []string{}, "Array of incident channel IDs associated with this subscription")
	cmd.Flags().String("incident-profile-id", "", "Incident Profile ID associated with this subscription")
	cmd.Flags().String("incident-visibility", "", "Incident visibility (unlisted, private)")
	cmd.Flags().String("incident-creation-mode", "", "Incident creation rule (none, default, always, ignore-if-exists)")
	cmd.Flags().String("incident-watch-keys", "", "Source/key to watch in the incident (scans/page.url, scans/page.domain, scans/page.ip, scans/page.apexDomain, hostnames/hostname, hostnames/domain)")
}

func mapCmdToSubscriptionOptions(cmd *cobra.Command) (opts []api.SubscriptionOption, err error) {
	searchIds, _ := cmd.Flags().GetStringSlice("search-ids")
	if len(searchIds) == 0 {
		return nil, fmt.Errorf("search-ids is required")
	}
	frequency, _ := cmd.Flags().GetString("frequency")
	if frequency == "" {
		return nil, fmt.Errorf("frequency is required")
	}
	emailAddresses, _ := cmd.Flags().GetStringSlice("email-addresses")
	if len(emailAddresses) == 0 {
		return nil, fmt.Errorf("email-addresses is required")
	}
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	// optional flags
	isActive, _ := cmd.Flags().GetBool("is-active")
	ignoreTime, _ := cmd.Flags().GetBool("ignore-time")
	description, _ := cmd.Flags().GetString("description")
	weekDays, _ := cmd.Flags().GetStringSlice("week-days")
	permissions, _ := cmd.Flags().GetStringSlice("permissions")
	channelIds, _ := cmd.Flags().GetStringSlice("channel-ids")
	incidentChannelIds, _ := cmd.Flags().GetStringSlice("incident-channel-ids")
	incidentProfileId, _ := cmd.Flags().GetString("incident-profile-id")
	incidentVisibility, _ := cmd.Flags().GetString("incident-visibility")
	incidentCreationMode, _ := cmd.Flags().GetString("incident-creation-mode")
	incidentWatchKeys, _ := cmd.Flags().GetString("incident-watch-keys")

	return []api.SubscriptionOption{
		api.WithSubscriptionSearchIds(searchIds),
		api.WithSubscriptionFrequency(frequency),
		api.WithSubscriptionEmailAddresses(emailAddresses),
		api.WithSubscriptionName(name),
		api.WithSubscriptionDescription(description),
		api.WithSubscriptionIsActive(isActive),
		api.WithSubscriptionIgnoreTime(ignoreTime),
		api.WithSubscriptionWeekDays(weekDays),
		api.WithSubscriptionPermissions(permissions),
		api.WithSubscriptionChannelIds(channelIds),
		api.WithSubscriptionIncidentChannelIds(incidentChannelIds),
		api.WithSubscriptionIncidentProfileId(incidentProfileId),
		api.WithSubscriptionIncidentVisibility(incidentVisibility),
		api.WithSubscriptionIncidentCreationMode(incidentCreationMode),
		api.WithSubscriptionIncidentWatchKeys(incidentWatchKeys),
	}, nil
}
