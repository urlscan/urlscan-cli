package flags

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/api"
)

func AddJSONFlag(cmd *cobra.Command) {
	cmd.Flags().String("json", "", "JSON payload to send as request body")
}

func AddParamsFlag(cmd *cobra.Command) {
	cmd.Flags().String("params", "", "Query string parameters as JSON (e.g. '{\"key\":\"value\"}')")
}

func AddJSONLFlag(cmd *cobra.Command) {
	cmd.Flags().String("jsonl", "", "JSONL payload to send as request bodies (one JSON payload per line)")
}

func JSONToMap(s string) (map[string]any, error) {
	var m map[string]any
	err := json.Unmarshal([]byte(s), &m)
	return m, err
}

func GetJSONL(cmd *cobra.Command) ([]map[string]any, error) {
	s, _ := cmd.Flags().GetString("jsonl")
	if s == "" {
		return nil, nil
	}
	lines := strings.Split(s, "\n")
	var result []map[string]any
	for _, line := range lines {
		if line == "" {
			continue
		}
		m, err := JSONToMap(line)
		if err != nil {
			return nil, fmt.Errorf("invalid JSONL line: %w", err)
		}
		result = append(result, m)
	}
	return result, nil
}

func GetParams(cmd *cobra.Command) (map[string]any, error) {
	s, _ := cmd.Flags().GetString("params")
	if s != "" {
		return JSONToMap(s)
	}
	return nil, nil
}

func GetJSON(cmd *cobra.Command) (map[string]any, error) {
	s, _ := cmd.Flags().GetString("json")
	if s != "" {
		return JSONToMap(s)
	}
	return nil, nil
}

func AddForceFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("force", "f", false, "Force overwrite an existing file")
}

func AddAllFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("all", false, "Return all results; limit is ignored if --all is specified (default false)")
}

func AddLimitFlag(cmd *cobra.Command) {
	cmd.Flags().IntP("limit", "l", api.MaxTotal, "Maximum number of results that will be returned by the iterator")
}

func AddSizeFlag(cmd *cobra.Command, value int) {
	cmd.Flags().IntP("size", "s", value, "Number of results returned by the iterator in each batch")
}

func AddOutputFlag(cmd *cobra.Command, defaultExample string) {
	cmd.Flags().StringP("output", "o", "", fmt.Sprintf("Output file name (default %s)", defaultExample))
}

func AddTagsFlag(cmd *cobra.Command) {
	cmd.Flags().StringArrayP("tags", "t", []string{}, "User-defined tags to annotate this scan")
}

func AddCountryFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("country", "c", "", "Specify which country the scan should be performed from (2-Letter ISO-3166-1 alpha-2 country")
}

func AddCustomAgentFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("customagent", "a", "", "Override User-Agent for this scan")
}

func AddOverrideSafetyFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("overrideSafety", "o", "", "If set to any value, this will disable reclassification of URLs with potential PII in them")
}

func AddRefererFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("referer", "r", "", "Override HTTP referer for this scan")
}

func AddVisibilityFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("visibility", "v", "", "One of public, unlisted, private")
}

func AddWaitFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("wait", "w", false, "Wait for the scan(s) to finish")
}

func AddMaxWaitFlag(cmd *cobra.Command) {
	cmd.Flags().IntP("max-wait", "m", 60, "Maximum wait time per scan in seconds")
}

func AddScreenshotFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("screenshot", false, "Download only the screenshot (overrides wait)")
}

func AddDOMFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("dom", false, "Download only the DOM contents (overrides wait)")
}

func AddDownloadFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("download", false, "Download screenshot and DOM contents (overrides wait/dom/screenshot)")
}

func AddDirectoryPrefixFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("directory-prefix", "P", ".", "Set directory prefix where file will be saved")
}
