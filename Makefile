DOCKER_REPO  = teliaoss/github-pr-resource
TARGET      ?= darwin
ARCH        ?= amd64
SRC          = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

export GO111MODULE=on

default: test

generate:
	@echo "== Go Generate =="
	go generate ./...

build: test
	@echo "== Build =="
	CGO_ENABLED=0 GOOS=$(TARGET) GOARCH=$(ARCH) go build -o check -v cmd/check/main.go
	CGO_ENABLED=0 GOOS=$(TARGET) GOARCH=$(ARCH) go build -o in -v cmd/in/main.go
	CGO_ENABLED=0 GOOS=$(TARGET) GOARCH=$(ARCH) go build -o out -v cmd/out/main.go

test:
	@echo "== Test =="
	gofmt -s -l -w $(SRC)
	go vet -v ./...
	go test -race -v ./...

e2e: test
	@echo "== Integration =="
	go test -race -v ./... -tags=e2e

clean:
	@echo "== Cleaning =="
	rm check
	rm in
	rm out

lint:
	@echo "== Lint =="
	golint cmd
	golint src
	golint e2e

docker:
	@echo "== Docker build =="
	docker build -t $(DOCKER_REPO):dev .

.PHONY: default generate build test docker e2e clean lint
