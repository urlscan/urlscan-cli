BINARY=./dist/urlscan

# setup the -ldflags option for go build
LDFLAGS=-ldflags "-X github.com/urlscan/urlscan-cli/pkg/version.Version=${VERSION}"

.PHONY: install update lint build docs clean release

install:
	go mod tidy
	# install dev tools
	cobra-cli help 2>/dev/null 1>&2 || go install github.com/spf13/cobra-cli@latest
	golangci-lint help 2>/dev/null 1>&2 || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	goreleaser help 2>/dev/null 1>&2 || go install github.com/goreleaser/goreleaser/v2@latest

update:
	go get -u ./...
	go mod tidy

lint: install
	golangci-lint run --fix

build: install lint
	go build ${LDFLAGS} -o ${BINARY} main.go

docs: install
	rm docs/urlscan*.md
	go run docutil/main.go

clean:
	rm -rf ./dist

release:
	goreleaser release --snapshot --clean

all: install build docs
