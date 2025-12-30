.DEFAULT_GOAL := help

COMMIT=$(shell git rev-parse HEAD)
BUILD_VERSION := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u "+20%y-%m-%dT%H:%MZ")

GO_SRC_FOLDER=pkg
GO_UNIT_TESTS=./pkg/...
GO=go
GOPATH_BIN=$(subst :,/bin:,$(GOPATH))/bin
export PATH := $(PATH):$(GOPATH_BIN)

.PHONY: test lint

#gogenerate: @ Generate mock and other test artifacts
gogenerate:
	@echo "[MockGeneration]=============================="
	@test -x $(GOPATH)/bin/mockgen || $(GO) get github.com/golang/mock/mockgen
	$(GO) generate -v ./...

#test: @ convenient task to run tests
test: lint
	@echo ====[Running Unit Tests]=================================================================
	$(GO) vet ./$(GO_SRC_FOLDER)/...
	@mkdir -p coverage
	$(GO) run gotest.tools/gotestsum@latest --junitfile coverage/coverage.xml --format pkgname \
		-- -v $(GO_UNIT_TESTS) -coverprofile coverage/coverage.fmt fmt
	$(GO) tool cover -func=coverage/coverage.fmt
	$(GO) tool cover -html=coverage/coverage.fmt -o coverage/coverage.html

#clean: @ remove generated files
clean:
	rm -rf pkg/logging/http/mocks

#lint: @ run golint to check code style
lint:
	@echo ====[GOLINT]=============================================================================
	@test -x $(GOPATH)/bin/staticcheck || (echo ====[Installing staticcheck@latest]====================================================== && $(GO) install honnef.co/go/tools/cmd/staticcheck@latest)
	@echo ====[Performing static checks]===========================================================
	@staticcheck ./...

#help: @ List available tasks on this project
help:
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#'  \
	| awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
