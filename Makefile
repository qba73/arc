.phony: test testc

ROOT                := $(PWD)
GO_HTML_COV         := ./coverage.html
GO_TEST_OUTFILE     := ./c.out
GO_DOCKER_IMAGE     := golang:1.16
GO_DOCKER_CONTAINER := arct-container
CC_TEST_REPORTER_ID := ${CC_TEST_REPORTER_ID}
CC_PREFIX           := github.com/qba73/arct

SHELL     := /bin/bash
PROJECT   := arct
VCS_REF   := `git rev-parse HEAD`
ITERATION := $(shell date -u +%Y-%m-%dT%H-%M-%SZ)
BUILD_DATE := `date -u +"%Y-%m-%d-%H-%M-%SZ"`
GOARCH    := amd64
VERSION   := 0.1.0

# Let's parse make target comments prefixed with ## and generate help output for the user. 
define PRINT_HELP_PYSCRIPT
import re, sys

for line in sys.stdin:
	match = re.match(r'^([a-zA-Z_-]+):.*?## (.*)$$', line)
	if match:
		target, help = match.groups()
		print("%-20s %s" % (target, help))
endef
export PRINT_HELP_PYSCRIPT


default: help

help:
	@python -c "$$PRINT_HELP_PYSCRIPT" < $(MAKEFILE_LIST)

# ==============================================================================
# Running tests on the local machine

test: ## Run unit tests and staticcheck locally
	go test -race -count=1 -v
	staticcheck ./...

# ==============================================================================
# Modules support

deps-reset: ## Reset go mod dependencies
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy: ## Resolve dependencies
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list: ## List Go modules
	go list -mod=mod all

build: ## Build Go binaries
	go build -ldflags "-X main.Commit=${VCS_REF} -X main.Version=${VERSION} -X main.Date=${BUILD_DATE}" -o arct ./cmd/arc/main.go

# ==============================================================================
# Go releaser

snapshot:
	goreleaser build --snapshot --rm-dist

# ==============================================================================
## CodeClimate tests and coverage

clean: ## Remove docker container if exist
	docker rm -f ${GO_DOCKER_CONTAINER} || true

testc: ## Run unittests inside container
	docker run -w /app -v ${ROOT}:/app ${GO_DOCKER_IMAGE} go test -v ./... -coverprofile=${GO_TEST_OUTFILE}
	docker run -w /app -v ${ROOT}:/app ${GO_DOCKER_IMAGE} go tool cover -html=${GO_TEST_OUTFILE} -o ${GO_HTML_COV}

lint: ## Run linter inside container
	docker run --rm -v ${ROOT}:/data cytopia/golint .

# Custom logic for Code Climate (they have specific requirements)
_before-cc:
	# Download CC test report
	docker run -w /app -v ${ROOT}:/app ${GO_DOCKER_IMAGE} \
		/bin/bash -c \
		"curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter"
	
	# Make the test reporter executable
	docker run -w /app -v ${ROOT}:/app ${GO_DOCKER_IMAGE} chmod +x ./cc-test-reporter

	# Run before build
	docker run -w /app -v ${ROOT}:/app \
		-e CC_TEST_REPORTER_ID=${CC_TEST_REPORTER_ID} ${GO_DOCKER_IMAGE} \
		./cc-test-reporter before-build

_after-cc:
	# Handle custom prefix
	$(eval PREFIX=${CC_PREFIX})
ifdef prefix
	$(eval PREFIX=${prefix})
endif
	# Upload test coverage data file to Code Climate
	docker run -w /app -v ${ROOT}:/app \
		-e CC_TEST_REPORTER_ID=${CC_TEST_REPORTER_ID} \
		${GO_DOCKER_IMAGE} ./cc-test-reporter after-build --prefix ${PREFIX}

test-ci: _before-cc testc _after-cc

