package main

import (
	"log"

	"github.com/spf13/cobra/doc"
	"github.com/urlscan/urlscan-cli/cmd"
)

func main() {
	cmd.RootCmd.DisableAutoGenTag = true

	err := doc.GenMarkdownTree(cmd.RootCmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
