SHELL := /bin/bash

export PROJECT = arct
VCS_REF=`git rev-parse HEAD`
ITERATION=$(shell date -u +%Y-%m-%dT%H-%M-%SZ)
GOARCH=amd64

.phony: test

# ==============================================================================
# Running tests within the local computer

test:
	go test ./... -v -count=1
	staticcheck ./...

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

## Building binaries
build:
	go build -ldflags "-X main.build=${VCS_REF}" ./cmd/arct/main.go

## Go releaser

snapshot:
	goreleaser build --snapshot --rm-dist
