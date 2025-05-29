# Development

## Requirements

- Go v1.24+
- [golangci-lint](https://github.com/golangci/golangci-lint) as a linter
- [Lefthook](https://github.com/evilmartians/lefthook) as a pre-commit hook manager
- [GoReleaser](https://github.com/goreleaser/goreleaser) as a release manager

## Setup

```bash
git clone https://github.com/urlscan/urlscan-cli
cd urlscan-cli

go mod tidy
```

## Project Layout

- `main.go`: the root command
- `api/`: main utilities for interacting with the urlscan.io API
- `cmd/`: individual urlscan commands
- `pkg/`: packages for implementing commands

## Test

```bash
go test ./...
```

## Lint

```bash
golangci-lint run --fix
```

## Build

```bash
# build the CLI for testing
go build main.go
```

```bash
# build binaries for releasing (without publishing)
goreleaser release --snapshot --clean
```

## Adding a New Command

Use [cobra-cli](https://github.com/spf13/cobra-cli).

```bash
cobra-cli add ...
```
