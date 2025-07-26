package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var longCompletionCmd = `To load completions:

Bash:

    # for Linux
    $ urlscan completion bash > /etc/bash_completion.d/urlscan
    # for macOS
    $ urlscan completion bash > "$(brew --prefix)/etc/bash_completion.d/urlscan"

ZSH:

    $ urlscan completion zsh > "${fpath[1]}/_urlscan"
    # for oh-my-zsh
    $ mkdir -p "$ZSH/completions/"
    $ urlscan completion zsh > "$ZSH/completions/_urlscan"

Fish:

    $ urlscan completion fish > ~/.config/fish/completions/urlscan.fish`

var completionCmd = &cobra.Command{
	Use:   "completion <shell>",
	Short: "Output shell completion code for the specified shell (bash, zsh, fish)",
	Long:  longCompletionCmd,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}
		shell := args[0]

		var err error
		switch shell {
		case "bash":
			err = RootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			err = RootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			err = RootCmd.GenFishCompletion(os.Stdout, true)
		default:
			return fmt.Errorf("unsupported shell type %q", shell)
		}
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(completionCmd)
}
