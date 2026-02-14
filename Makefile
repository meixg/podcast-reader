.PHONY: all build test run-server run-server-watch run-server-dev run-cli clean fmt vet lint help

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

## run-server-watch: Run the API server with auto-reload on file changes
run-server-watch:
	@AIR_BIN="$(shell go env GOPATH)/bin/air"; \
	if [ ! -f "$$AIR_BIN" ]; then \
		echo "Installing air..."; \
		go install github.com/air-verse/air@latest; \
	fi; \
	"$$AIR_BIN"

## run-server-dev: Run the API server with simple auto-restart (no external tools required)
run-server-dev:
	@echo "Starting server with auto-restart..."
	@while true; do \
		echo "[$(shell date '+%H:%M:%S')] Building server..."; \
		$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_SERVER) ./cmd/server && \
		echo "[$(shell date '+%H:%M:%S')] Starting server..." && \
		./$(BUILD_DIR)/$(BINARY_SERVER) -verbose & \
		SERVER_PID=$$!; \
		inotifywait -qq -e modify -r --include '\.go$$' . 2>/dev/null || \
		(fstat -q . 2>/dev/null || sleep 5); \
		kill $$SERVER_PID 2>/dev/null; \
		echo "[$(shell date '+%H:%M:%S')] Server stopped, restarting..."; \
		sleep 1; \
	done

## run-cli: Run the CLI tool (requires URL argument)
run-cli: build-cli
	@echo "Usage: make run-cli URL=<url>"
	@echo "Example: make run-cli URL=https://www.xiaoyuzhoufm.com/episode/123"

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -rf tmp
	rm -f coverage.out coverage.html build-errors.log

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
