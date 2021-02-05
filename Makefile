BIN := awsmfa
BUILD_LDFLAGS := "-s -w"
GOBIN ?= $(shell go env GOPATH)/bin
GODOWNLOADER_VERSION := 0.1.0
export GO111MODULE=on

repo_name := d-tsuji/awsmfa
current_dir := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: all
all: clean build

.PHONY: deps
deps:
	go mod tidy

.PHONY: devel-deps
devel-deps: deps
	sh -c '\
      tmpdir=$$(mktemp -d); \
      cd $$tmpdir; \
      go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.36.0; \
      rm -rf $$tmpdir'

.PHONY: build
build:
	go build -ldflags=$(BUILD_LDFLAGS) -o $(BIN) -trimpath ./cmd/awsmfa

.PHONY: test
test: deps
	go test -v ./...

.PHONY: test-cover
test-cover: deps
	go test -v ./... -cover -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html

.PHONY: lint
lint: devel-deps
	go vet ./...
	golangci-lint run --config .golangci.yml ./...

.PHONY: clean
clean:
	rm -rf $(BIN)
	go clean

.PHONY: installer
installer:
	sh -c '\
      tmpdir=$$(mktemp -d); \
      cd $$tmpdir; \
      # Build from source, because "go get github.com/goreleaser/godownloader" will result in an error; \
      git clone --quiet --depth 1 https://github.com/goreleaser/godownloader && cd godownloader; \
      go run . -f -r ${repo_name} -o ${current_dir}/install.sh; \
      rm -rf $$tmpdir'
