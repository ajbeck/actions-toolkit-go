GOFILES := $(shell find . -name '*.go' -type f)
TEST_PKGS ?= ./...
TEST_ARGS ?=

.PHONY: fmt vet build test

fmt: $(GOFILES) go.mod
	go fmt ./...

vet: $(GOFILES) go.mod
	go vet -json ./...

build: $(GOFILES) go.mod
	go build -v ./...

test: $(GOFILES) go.mod
	GOCACHE=$(CURDIR)/.gocache/build GOMODCACHE=$(CURDIR)/.gocache/mod go test $(TEST_PKGS) $(TEST_ARGS) -v -cover -coverprofile=coverage.out
	go tool cover -func=coverage.out
