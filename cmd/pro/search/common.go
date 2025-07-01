package search

import "github.com/spf13/cobra"

func setCreateOrUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("datasource", "D", "scans", "Which data this saved search operates on (hostnames or scans) (required")
	cmd.Flags().StringP("name", "n", "", "Name of the saved search (required)")
	cmd.Flags().StringP("tlp", "t", "red", "TLP (Traffic Light Protocol) of the saved search (required)")
	cmd.Flags().StringP("query", "q", "", "Search query of the saved search (required)")
	cmd.Flags().IntP("pass", "P", 2, "2 for inline-matching, 10 for bookmark-only (required)")

	cmd.Flags().StringSliceP("permissions", "p", []string{}, "Permissions of the saved search (optional)")
	cmd.Flags().StringP("description", "d", "", "Short description of the saved search (optional)")
	cmd.Flags().StringP("long-description", "l", "", "Long description of the saved search (optional)")
	cmd.Flags().StringP("owner-description", "o", "", "Owner description of the saved search (optional)")
	cmd.Flags().StringSliceP("tags", "T", []string{}, "Tags of the saved search (optional)")
	cmd.Flags().StringSliceP("user-tags", "u", []string{}, "User tags of the saved search (optional)")
}
