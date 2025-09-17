## urlscan completion

Output shell completion code for the specified shell (bash, zsh, fish)

### Synopsis

To load completions:

Bash (Linux or macOS):

    # for Linux (make sure you have bash-completion package)
    $ urlscan completion bash > /etc/bash_completion.d/urlscan
    # for macOS
    $ urlscan completion bash > "$(brew --prefix)/etc/bash_completion.d/urlscan"

ZSH (Linux or macOS):

    $ urlscan completion zsh > "${fpath[1]}/_urlscan"
    # for oh-my-zsh
    $ mkdir -p "$ZSH/completions/"
    $ urlscan completion zsh > "$ZSH/completions/_urlscan"

Fish (Linux or macOS):

    $ urlscan completion fish > ~/.config/fish/completions/urlscan.fish

```
urlscan completion <shell> [flags]
```

### Options

```
  -h, --help   help for completion
```

### SEE ALSO

* [urlscan](urlscan.md)	 - A CLI tool for interacting with urlscan.io

