.PHONY: all gomod lint dev release cobra

target: all

gomod:
	go mod tidy

lint:
	golangci-lint help 2>/dev/null 1>&2 || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	golangci-lint run --fix

dev: gomod lint
	go build main.go

release:
	goreleaser release --snapshot --clean

all: dev
	cobra-cli help 2>/dev/null 1>&2 || go install github.com/spf13/cobra-cli@latest
	cobra-cli add ...
