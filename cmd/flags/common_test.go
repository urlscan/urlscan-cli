package flags

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "test",
		RunE: func(cmd *cobra.Command, args []string) error { return nil },
	}
	AddJSONFlag(cmd)
	AddParamsFlag(cmd)
	AddJSONLFlag(cmd)
	return cmd
}

func TestGetJSON(t *testing.T) {
	cmd := newTestCommand()
	require.NoError(t, cmd.Flags().Set("json", `{"name":"hello"}`))

	json, err := GetJSON(cmd)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"name": "hello"}, json)
}

func TestGetParams(t *testing.T) {
	cmd := newTestCommand()
	require.NoError(t, cmd.Flags().Set("params", `{"size":"50","datasource":"hostnames"}`))

	params, err := GetParams(cmd)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"size": "50", "datasource": "hostnames"}, params)
}

func TestGetJSON_Empty(t *testing.T) {
	cmd := newTestCommand()

	jsonBody, err := GetJSON(cmd)
	require.NoError(t, err)
	assert.Empty(t, jsonBody)
}

func TestGetParams_Empty(t *testing.T) {
	cmd := newTestCommand()
	params, err := GetParams(cmd)
	require.NoError(t, err)
	assert.Nil(t, params)
}

func TestGetParams_InvalidParams(t *testing.T) {
	cmd := newTestCommand()
	require.NoError(t, cmd.Flags().Set("params", `{invalid}`))

	_, err := GetParams(cmd)
	assert.Error(t, err)
}

func TestGetJSONL(t *testing.T) {
	cmd := newTestCommand()
	require.NoError(t, cmd.Flags().Set("jsonl", "{\"url\":\"https://example.com\"}\n{\"url\":\"https://example.net\"}\n"))

	jsonl, err := GetJSONL(cmd)

	require.NoError(t, err)
	assert.Equal(t, []map[string]any{
		{"url": "https://example.com"},
		{"url": "https://example.net"},
	}, jsonl)
}

func TestGetJSONL_SkipsEmptyLines(t *testing.T) {
	cmd := newTestCommand()
	require.NoError(t, cmd.Flags().Set("jsonl", "{\"url\":\"https://example.com\"}\n\n{\"url\":\"https://example.net\"}\n\n"))

	jsonl, err := GetJSONL(cmd)
	require.NoError(t, err)
	assert.Len(t, jsonl, 2)
}

func TestGetJSONL_InvalidJSON(t *testing.T) {
	cmd := newTestCommand()
	require.NoError(t, cmd.Flags().Set("jsonl", "{\"url\":\"https://example.com\"}\nnot json\n"))

	_, err := GetJSONL(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}
