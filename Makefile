.PHONY: all build test run-server run-cli clean fmt vet lint help

# Variables
BINARY_SERVER=podcast-server
BINARY_CLI=podcast-downloader
BUILD_DIR=build
GO=go
GOFLAGS=-v

all: build

## build: Build both CLI and server binaries
build: build-cli build-server

## build-cli: Build CLI tool
build-cli:
	@echo "Building CLI tool..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_CLI) ./cmd/downloader

## build-server: Build API server
build-server:
	@echo "Building API server..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_SERVER) ./cmd/server

## test: Run all tests
test:
	$(GO) test -v -race -cover ./...

## test-cover: Run tests with coverage report
test-cover:
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

## run-server: Run the API server
run-server: build-server
	./$(BUILD_DIR)/$(BINARY_SERVER) -verbose

## run-cli: Run the CLI tool (requires URL argument)
run-cli: build-cli
	@echo "Usage: make run-cli URL=<url>"
	@echo "Example: make run-cli URL=https://www.xiaoyuzhoufm.com/episode/123"

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

## fmt: Format code
fmt:
	$(GO) fmt ./...

## vet: Run go vet
vet:
	$(GO) vet ./...

## lint: Run golangci-lint (requires golangci-lint to be installed)
lint:
	golangci-lint run ./...

## mod-tidy: Tidy go modules
mod-tidy:
	$(GO) mod tidy

## mod-download: Download go modules
mod-download:
	$(GO) mod download

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
