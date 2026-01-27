package datadump

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/urlscan/urlscan-cli/api"
	"github.com/urlscan/urlscan-cli/cmd/flags"
	"github.com/urlscan/urlscan-cli/pkg/utils"
)

var DownloadCmdExample = `  urlscan pro datadump download days/api/20260101.gz
  urlscan pro datadump download hours/api/20260101/20260101-01.gz
  echo "<path>" | urlscan pro datadump download -

  # use --follow option to download all files from a datadump path
  # for example, the following commands download all the files listed by 'urlscan pro datadump list hours/dom/20260101/'
  # note: --follow memoizes downloaded files in a local database to avoid re-downloading, so it's safe to run it periodically
  urlscan pro datadump download hours/dom/20260101/ --follow
  # if date is not provided, all the available files (files within the last 7 days) will be downloaded
  urlscan pro datadump download hours/api/ --follow`

var downloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "Download the data dump file",
	Example: DownloadCmdExample,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 1 {
			return cmd.Usage()
		}

		output, _ := cmd.Flags().GetString("output")
		directoryPrefix, _ := cmd.Flags().GetString("directory-prefix")
		force, _ := cmd.Flags().GetBool("force")
		extract, _ := cmd.Flags().GetBool("extract")
		follow, _ := cmd.Flags().GetBool("follow")

		reader := utils.StringReaderFromCmdArgs(args)
		path, err := reader.ReadString()
		if err != nil {
			return err
		}

		client, err := utils.NewAPIClient()
		if err != nil {
			return err
		}
		// disable auto gzip decompression to streamline extraction process
		client.SetDisableCompression(true)

		// open the database
		db, err := utils.NewDatabase()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer func() {
			closeErr := db.Close()
			if closeErr != nil && err == nil {
				err = closeErr
			}
		}()

		isSingleFile := strings.HasSuffix(path, ".gz")
		if follow && isSingleFile {
			return fmt.Errorf("--follow option can only be used on entire directories, not individual files")
		}
		if !follow && !isSingleFile {
			return fmt.Errorf("please use --follow option to download entire directories")
		}

		// explode paths to download
		paths := []string{path}
		if follow {
			missingPaths, err := findMissingPaths(db, client, path, force)
			if err != nil {
				return err
			}
			paths = missingPaths
		}

		for _, path := range paths {
			if err := download(client, db, path, output, directoryPrefix, force, extract); err != nil {
				return err
			}
		}

		return nil
	},
}

func download(client *utils.APIClient, db *utils.Database, path, output, directoryPrefix string, force, extract bool) error {
	if output == "" {
		output = filepath.Base(path)
	}

	err := utils.DownloadWithSpinner(
		utils.NewDownloadOptions(
			utils.WithDownloadClient(client),
			utils.WithDownloadOutput(output),
			utils.WithDownloadDirectoryPrefix(directoryPrefix),
			utils.WithDownloadForce(force),
			utils.WithDownloadURL(api.PrefixedPath(fmt.Sprintf("/datadump/link/%s", path))),
		))
	if err != nil {
		return err
	}

	// update the database after successful download
	err = db.SetDataDump(path, filepath.Join(directoryPrefix, output))
	if err != nil {
		return fmt.Errorf("failed to update the database: %w", err)
	}

	// extract if requested
	if extract {
		err = utils.Extract(output, utils.NewExtractOptions(utils.WithExtractForce(force), utils.WithExtractDirectoryPrefix(directoryPrefix)))
		if err != nil {
			return err
		}
	}

	return nil
}

func findMissingPaths(db *utils.Database, client *utils.APIClient, path string, force bool) ([]string, error) {
	list, err := client.BulkGetDataDumpList(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get datadump list: %w", err)
	}

	paths := []string{}
	for _, file := range list.Files {
		// if force is set, re-download all files
		if force {
			paths = append(paths, file.Path)
			continue
		}

		downloaded, err := db.HasDataDumpBeenDownloaded(file.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to check download status for %s: %w", file.Path, err)
		}
		if !downloaded {
			paths = append(paths, file.Path)
		}
	}

	return paths, nil
}

func init() {
	flags.AddOutputFlag(downloadCmd, "<path>.gz")
	flags.AddForceFlag(downloadCmd)
	flags.AddDirectoryPrefixFlag(downloadCmd)

	downloadCmd.Flags().BoolP("extract", "x", false, "Extract the downloaded file")
	downloadCmd.Flags().BoolP("follow", "F", false, "Download missing files from the datadump path")

	RootCmd.AddCommand(downloadCmd)
}
