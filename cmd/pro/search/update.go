package search

import (
	"fmt"

	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var updateCmdExample = `  urlscan pro save-search update <search-id> -D scans -n <name> -q <query>
  echo "<search-id>" | urlscan pro save-search update - -D scans -n <name> -q <query> -`

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a saved search",
	Example: updateCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		reader := utils.StringReaderFromCmdArgs(args)
		id, err := reader.ReadString()
		if err != nil {
			return err
		}

		// required flags (show usage if any are missing)
		datasource, _ := cmd.Flags().GetString("datasource")
		if datasource == "" {
			return cmd.Usage()
		}
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return cmd.Usage()
		}
		query, _ := cmd.Flags().GetString("query")
		if query == "" {
			return cmd.Usage()
		}
		tlp, _ := cmd.Flags().GetString("tlp")
		if tlp == "" {
			return cmd.Usage()
		}
		pass, _ := cmd.Flags().GetInt("pass")
		if pass == 0 {
			return cmd.Usage()
		}

		// optional flags
		permissions, _ := cmd.Flags().GetStringSlice("permissions")
		description, _ := cmd.Flags().GetString("description")
		longDescription, _ := cmd.Flags().GetString("long-description")
		ownerDescription, _ := cmd.Flags().GetString("owner-description")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		userTags, _ := cmd.Flags().GetStringSlice("user-tags")

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		res, err := client.UpdateSavedSearch(
			api.WithSavedSearchID(id),
			api.WithSavedSearchDatasource(datasource),
			api.WithSavedSearchName(name),
			api.WithSavedSearchQuery(query),
			api.WithSavedSearchTLP(tlp),
			api.WithSavedSearchPass(pass),
			api.WithSavedSearchPermissions(permissions),
			api.WithSavedSearchDescription(description),
			api.WithSavedSearchLongDescription(longDescription),
			api.WithSavedSearchOwnerDescription(ownerDescription),
			api.WithSavedSearchTags(tags),
			api.WithSavedSearchUserTags(userTags),
		)
		if err != nil {
			return err
		}

		fmt.Print(string(res.Raw))

		return nil
	},
}

func init() {
	setCreateOrUpdateFlags(updateCmd)

	RootCmd.AddCommand(updateCmd)
}
