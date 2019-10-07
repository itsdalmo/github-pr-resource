mkfile := $(abspath $(lastword $(MAKEFILE_LIST)))
dir := $(dir $(ci_mkfile))

BUILD_DIR := build
DOCKER_REPO := teliaoss/github-pr-resource
BINARY := check
EXTENSION :=

GO11MODULE := on
export GO11MODULE

.PHONY: generate
generate:
	@go generate ./...

.PHONY: test
test: generate
	@gofmt -s -l -w .
	@go vet ./...
	@go test -race -v -cover ./...

.PHONY: e2e
e2e:
	@go test -race -v ./... -tags=e2e

.PHONY: docker
docker:
	@docker build -t $(DOCKER_REPO):dev .

.PHONY: build
build: test
	make go-build BINARY=check --no-print-directory
	make go-build BINARY=in --no-print-directory
	make go-build BINARY=out --no-print-directory

.PHONY: go-build
go-build:
	@CGO_ENABLED=0 GOOS=${OS} GOOS=${ARCH} go build -o $(BUILD_DIR)/$(BINARY)$(EXTENSION) -ldflags="-s -w" -v cmd/$(BINARY)/main.go

.PHONY: ci
ci: build
	@if [ -n "$(git status --porcelain)" ];then echo "Diff in generated files and/or formatting" && exit 1; fi
