.PHONY: build run restart clean test

# Build the server binary
build:
	go build -o podcast-server ./cmd/server

# Run the server
run:
	go run ./cmd/server/main.go

# Restart the server: kill existing, rebuild, and start
restart:
	@echo "Restarting podcast server..."
	@pkill -9 podcast-server 2>/dev/null || true
	@sleep 1
	@go build -o podcast-server ./cmd/server
	@./podcast-server & \
	echo "Server restarted on port 8080"
	@echo "Logs: tail -f output/server.log"

# Clean build artifacts
clean:
	rm -f podcast-server
	rm -rf output/

# Run tests
test:
	go test ./...

# Format code
fmt:
	gofmt -w .

# Run linter
lint:
	go vet ./...
