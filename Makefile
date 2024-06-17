# See: https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
# For a list of valid GOOS and GOARCH values
# Note: these can be overriden on the command line e.g. `make PLATFORM=<platform> ARCH=<arch>`
PLATFORM=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)

GOTESTSUM=go run gotest.tools/gotestsum@v1.10.0

.DEFAULT_GOAL := help
.PHONY: server frontend frontend-dev

##@ Building

server: ## Build the server
	@echo "Building the Server API ..."
	@CGO_ENABLED=0 GOOS=$(PLATFORM) GOARCH=$(ARCH) go build -trimpath --installsuffix cgo --ldflags "-s" -o dist/$(server) main.go

frontend: 
	cd frontend && yarn build

frontend-dev:
	yarn dev

##@ Dependencies

tidy: ## Tidy up the go.mod file
	@go mod tidy

##@ Testing
.PHONY: test-server

test:	## Run server tests
	$(GOTESTSUM) --format pkgname-and-test-fails --format-hide-empty-pkg --hide-summary skipped -- -cover  ./...

##@ Cleanup

clean: ## Remove all build artifacts
	@echo "Clearing the dist directory..."
	@rm -rf dist/*

##@ Helpers

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
