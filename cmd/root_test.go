package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := os.Setenv("URLSCAN_API_KEY", "dummy-api-key-for-testing")
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func listAnnotatedCommands(cmd *cobra.Command, key, value string) []*cobra.Command {
	var result []*cobra.Command

	if cmd.Annotations != nil {
		v, ok := cmd.Annotations[key]
		if ok && v == value {
			result = append(result, cmd)
		}
	}
	for _, sub := range cmd.Commands() {
		result = append(result, listAnnotatedCommands(sub, key, value)...)
	}

	return result
}

func TestExactArgs1CommandsShowUsage(t *testing.T) {
	cmds := listAnnotatedCommands(RootCmd, "args", "exact1")

	assert.NotEmpty(t, cmds)

	for _, cmd := range cmds {
		t.Run(cmd.CommandPath(), func(t *testing.T) {
			// Create a fresh root command for each test
			root := RootCmd
			root.SetArgs(commandPathToArgs(cmd))

			var out bytes.Buffer
			root.SetOut(&out)
			root.SetErr(&out)

			err := root.Execute()
			if err != nil {
				t.Errorf("command %s returned error: %v", cmd.CommandPath(), err)
			}

			output := out.String()
			assert.Contains(t, output, "Usage:", "command %s should show usage", cmd.CommandPath())
		})
	}
}

// commandPathToArgs converts a command path like "urlscan pro hostname" to args ["pro", "hostname"]
func commandPathToArgs(cmd *cobra.Command) []string {
	var args []string
	for c := cmd; c != nil && c.Parent() != nil; c = c.Parent() {
		args = append([]string{c.Name()}, args...)
	}
	return args
}
