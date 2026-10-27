BINARY=./dist/urlscan

# setup the -ldflags option for go build
LDFLAGS=-ldflags "-X github.com/urlscan/urlscan-cli/pkg/version.Version=${VERSION}"

.PHONY: update lint build docs clean release all

update:
	go get -u
	go mod tidy

lint:
	go tool golangci-lint run --fix

build: lint
	go build ${LDFLAGS} -o ${BINARY} main.go

docs:
	rm docs/urlscan*.md
	go run docutil/main.go

clean:
	rm -rf ./dist

release:
	go tool goreleaser release --snapshot --clean

all: build docs
