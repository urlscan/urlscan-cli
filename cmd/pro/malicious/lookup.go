package malicious

import (
	"fmt"
	"net/url"
	"slices"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var validTypes = []string{"ip", "hostname", "domain", "url"}

var lookupCmdExample = `  urlscan pro malicious lookup ip 192.0.2.1
  urlscan pro malicious lookup hostname www.example.com
  urlscan pro malicious lookup domain example.com
  urlscan pro malicious lookup url "https://example.com/path"
  echo "192.0.2.1" | urlscan pro malicious lookup ip -`

var lookupCmd = &cobra.Command{
	Use:     "lookup <type> <value>",
	Short:   "Look up how often an observable has been seen in malicious scan results",
	Long:    "Look up how often an observable has been seen in malicious scan results, along with first and last seen timestamps. Type must be one of: ip, hostname, domain, url.",
	Example: lookupCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return cmd.Usage()
		}

		observableType := args[0]
		if !slices.Contains(validTypes, observableType) {
			return fmt.Errorf("invalid type %q: must be one of %v", observableType, validTypes)
		}

		refang, _ := cmd.Flags().GetBool("refang")

		reader := utils.StringReaderFromCmdArgs(args[1:])
		value, err := reader.ReadString()
		if err != nil {
			return err
		}

		if refang {
			value = utils.Refang(value)
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		resp, err := client.NewRequest().Get(api.PrefixedPath(fmt.Sprintf("/malicious/%s/%s", observableType, url.PathEscape(value))))
		if err != nil {
			return err
		}

		fmt.Print(resp.PrettyJSON())

		return nil
	},
}

func init() {
	flags.AddRefangFlag(lookupCmd)

	RootCmd.AddCommand(lookupCmd)
}
