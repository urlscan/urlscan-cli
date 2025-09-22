package incident

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/api"
)

func setCreateOrUpdateFlags(cmd *cobra.Command) {
	// required flags
	cmd.Flags().StringP("observable", "o", "", "Observable (hostname, domain, IP or URL) (required)")

	// defaulted flags
	cmd.Flags().Int("countries-per-interval", 1, "Countries per interval")
	cmd.Flags().Int("expire-after", 0, "Expire after in seconds (default 0)")
	cmd.Flags().Int("scan-interval-after-malicious", 0, "Scan interval after malicious in seconds (default 0)")
	cmd.Flags().Int("scan-interval-after-suspended", 0, "Scan interval after suspended in seconds (default 0)")
	cmd.Flags().Int("scan-interval", 0, "Scan interval in seconds (default 0)")
	cmd.Flags().Int("stop-delay-inactive", 0, "Stop delay inactive in seconds (default 0)")
	cmd.Flags().Int("stop-delay-malicious", 0, "Stop delay malicious in seconds (default 0)")
	cmd.Flags().Int("stop-delay-suspended", 0, "Stop delay suspended in seconds (default 0)")
	cmd.Flags().Int("user-agents-per-interval", 1, "User agents per interval")
	cmd.Flags().String("scan-interval-mode", "automatic", "Scan interval mode (automatic or manual)")
	cmd.Flags().String("visibility", "private", "Visibility (unlisted or private)")

	// optional flags
	cmd.Flags().String("incident-profile", "", "Incident profile (optional)")
	cmd.Flags().String("expire-at", "", "Expire at (optional)")
	cmd.Flags().StringSlice("channels", []string{}, "Channels")
	cmd.Flags().StringSlice("countries", []string{}, "Countries")
	cmd.Flags().StringSlice("user-agents", []string{}, "User agents")
	cmd.Flags().StringSlice("watched-attributes", []string{}, "Watched attributes")
}

func mapCmdToIncidentOptions(cmd *cobra.Command) (opts []api.IncidentOption, err error) {
	observable, _ := cmd.Flags().GetString("observable")
	if observable == "" {
		return nil, fmt.Errorf("observable is required")
	}

	channels, _ := cmd.Flags().GetStringSlice("channels")
	countries, _ := cmd.Flags().GetStringSlice("countries")
	countriesPerInterval, _ := cmd.Flags().GetInt("countries-per-interval")
	expireAfter, _ := cmd.Flags().GetInt("expire-after")
	expireAt, _ := cmd.Flags().GetString("expire-at")
	incidentProfile, _ := cmd.Flags().GetString("incident-profile")
	scanInterval, _ := cmd.Flags().GetInt("scan-interval")
	scanIntervalAfterMalicious, _ := cmd.Flags().GetInt("scan-interval-after-malicious")
	scanIntervalAfterSuspended, _ := cmd.Flags().GetInt("scan-interval-after-suspended")
	scanIntervalMode, _ := cmd.Flags().GetString("scan-interval-mode")
	stopDelayInactive, _ := cmd.Flags().GetInt("stop-delay-inactive")
	stopDelayMalicious, _ := cmd.Flags().GetInt("stop-delay-malicious")
	stopDelaySuspended, _ := cmd.Flags().GetInt("stop-delay-suspended")
	userAgents, _ := cmd.Flags().GetStringSlice("user-agents")
	userAgentsPerInterval, _ := cmd.Flags().GetInt("user-agents-per-interval")
	visibility, _ := cmd.Flags().GetString("visibility")
	watchedAttributes, _ := cmd.Flags().GetStringSlice("watched-attributes")

	opts = append(opts,
		api.WithIncidentChannels(channels),
		api.WithIncidentExpireAfter(expireAfter),
		api.WithIncidentScanInterval(scanInterval),
		api.WithIncidentScanIntervalMode(scanIntervalMode),
		api.WithIncidentWatchedAttributes(watchedAttributes),
		api.WithIncidentUserAgents(userAgents),
		api.WithIncidentUserAgentsPerInterval(userAgentsPerInterval),
		api.WithIncidentCountries(countries),
		api.WithIncidentCountriesPerInterval(countriesPerInterval),
		api.WithIncidentStopDelaySuspended(stopDelaySuspended),
		api.WithIncidentStopDelayInactive(stopDelayInactive),
		api.WithIncidentStopDelayMalicious(stopDelayMalicious),
		api.WithIncidentScanIntervalAfterSuspended(scanIntervalAfterSuspended),
		api.WithIncidentScanIntervalAfterMalicious(scanIntervalAfterMalicious),
		api.WithIncidentVisibility(visibility),
		api.WithIncidentExpireAt(expireAt),
		api.WithIncidentProfile(incidentProfile),
		api.WithIncidentObservable(observable),
	)

	return opts, nil
}
