package search

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	api "github.com/urlscan/urlscan-cli/api"

	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var createCmdExample = `  urlscan pro saved-search create -D scans -n <name> -q <query>`

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new saved search",
	Example: createCmdExample,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		id, _ := cmd.Flags().GetString("search-id")
		if id == "" {
			id = uuid.New().String()
		}
		err := utils.ValidateUUID(id)
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}

		result, err := client.CreateSavedSearch(
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

		fmt.Print(result.PrettyJSON())

		return nil
	},
}

func init() {
	setCreateOrUpdateFlags(createCmd)
	// optional flag, only for create command
	createCmd.Flags().StringP("search-id", "i", "", "Search ID (optional, if not provided a new ID will be generated)")

	RootCmd.AddCommand(createCmd)
}
