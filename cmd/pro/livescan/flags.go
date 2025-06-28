package livescan

import "github.com/spf13/cobra"

func addScannerIdFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("scanner-id", "s", "", "ID of the scanner (required)")
}

func addVisibilityFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("visibility", "v", "private", "Visibility of the scan (public, unlisted or private)")
}
