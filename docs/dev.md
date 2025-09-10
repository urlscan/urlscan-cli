# Development

## Requirements

- Go v1.24+
- [golangci-lint](https://github.com/golangci/golangci-lint) as a linter
- [Lefthook](https://github.com/evilmartians/lefthook) as a pre-commit hook manager
- [GoReleaser](https://github.com/goreleaser/goreleaser) as a release manager
- [Bats](https://github.com/bats-core/bats-core) as a integration test framework

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
- `docutil`: a package for generating Cobra command docs
- `pkg/`: packages for implementing commands

## Test

### Unit Test

```bash
go test ./...
```

### Integration Test

> [!NOTE]
> Make sure to install [bats-core](https://github.com/bats-core/bats-core), [bats-assert](https://github.com/bats-core/bats-assert) and [bats-support](https://github.com/bats-core/bats-support) (by `npm install -g bats bats-assert bats-support`).

```bash
bats test
```

## Lint

```bash
golangci-lint run --fix
```

## Build

```bash
make build
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
