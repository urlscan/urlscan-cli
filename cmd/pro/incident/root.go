package incident

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "incident",
	Short: "Incident sub-commands",
}
