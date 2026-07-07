package cmd

import (
	"github.com/njayp/ophis"
)

var config = &ophis.Config{
	CommandName: "mcp",
	Selectors: []ophis.Selector{
		{
			CmdSelector:           ophis.ExcludeCmdsContaining("urlscan key"),
			LocalFlagSelector:     nil,
			InheritedFlagSelector: nil,
			Middleware:            nil,
		},
	},
	DefaultEnv:     nil,
	ToolNamePrefix: "",
	SloggerOptions: nil,
	ServerOptions:  nil,
	Transport:      nil,
}

var mcpCmd = ophis.Command(config)

func init() {
	allow := map[string]bool{"start": true, "stream": true, "tools": true}

	for _, sub := range mcpCmd.Commands() {
		if !allow[sub.Name()] {
			mcpCmd.RemoveCommand(sub)
		}
	}

	RootCmd.AddCommand(mcpCmd)
}
