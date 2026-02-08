# Quickstart Guide: Podcast API Server

**Feature**: 003-podcast-api-server
**Date**: 2026-02-08
**Phase**: Phase 1 - Design & Contracts

## Overview

This guide provides quick reference instructions for setting up, running, and testing the Podcast API Server during development.

## Prerequisites

### Required Software

- **Go**: Version 1.25.5 or later
  - Install from: https://go.dev/dl/
  - Verify: `go version`
- **Git**: For cloning the repository
- **curl** or **HTTPie**: For testing API endpoints (or any HTTP client)

### System Requirements

- **OS**: Linux, macOS, or Windows
- **Disk Space**: 100MB+ for code, additional space for downloaded podcasts
- **Network**: Internet connection for accessing Xiaoyuzhou FM
- **Permissions**: Write access to `downloads/` directory

## Development Setup

### 1. Clone Repository

```bash
git clone https://github.com/meixg/podcast-reader.git
cd podcast-reader
git checkout 003-podcast-api-server
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Verify Build

```bash
# Build existing CLI tool (should still work)
go build -o podcast-downloader cmd/podcast-downloader/main.go

# Build new server (will fail initially during development)
go build -o podcast-server cmd/podcast-server/main.go
```

## Running the Server

### Development Mode

```bash
# Run directly with Go
go run cmd/podcast-server/main.go

# Or build and run
go build -o podcast-server cmd/podcast-server/main.go
./podcast-server
```

### Configuration

The server supports the following environment variables and command-line flags:

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `-port` | `PODCAST_SERVER_PORT` | `8080` | HTTP server port |
| `-host` | `PODCAST_SERVER_HOST` | `localhost` | Server bind address |
| `-downloads` | `PODCAST_DOWNLOADS_DIR` | `./downloads` | Downloads directory path |
| `-verbose` | `PODCAST_VERBOSE` | `false` | Enable verbose logging |

**Examples**:

```bash
# Use default settings
./podcast-server

# Custom port
./podcast-server -port 3000

# Custom downloads directory
./podcast-server -downloads /path/to/podcasts

# Verbose logging
./podcast-server -verbose

# Using environment variables
PODCAST_SERVER_PORT=3000 PODCAST_VERBOSE=true ./podcast-server
```

### Expected Output

```
2026/02/08 10:30:00 Starting Podcast API Server v1.0.0
2026/02/08 10:30:00 Scanning downloads directory: ./downloads
2026/02/08 10:30:00 Found 15 downloaded podcasts
2026/02/08 10:30:00 Server listening on http://localhost:8080
```

## API Usage Examples

### 1. Submit Download Task

**Using Test URLs**:

The project includes valid Xiaoyuzhou FM URLs in `examples/xiaoyuzhou_urls` for testing:

```bash
# Use the first test URL
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3"}'

# Or read from file
TEST_URL=$(head -n 1 examples/xiaoyuzhou_urls)
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d "{\"url\": \"$TEST_URL\"}"
```

**Request**:

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.xiaoyuzhoufm.com/episode/12345678"}'
```

**Response** (202 Accepted):

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "url": "https://www.xiaoyuzhou.fm/episode/12345678",
  "status": "pending",
  "created_at": "2026-02-08T10:30:00Z",
  "started_at": null,
  "completed_at": null,
  "error": null,
  "progress": null,
  "podcast": null
}
```

### 2. Check Task Status

**Request**:

```bash
curl http://localhost:8080/tasks/550e8400-e29b-41d4-a716-446655440000
```

**Response** (In Progress):

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "url": "https://www.xiaoyuzhou.fm/episode/12345678",
  "status": "in_progress",
  "created_at": "2026-02-08T10:30:00Z",
  "started_at": "2026-02-08T10:30:01Z",
  "completed_at": null,
  "error": null,
  "progress": 45,
  "podcast": null
}
```

**Response** (Completed):

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "url": "https://www.xiauzhou.fm/episode/12345678",
  "status": "completed",
  "created_at": "2026-02-08T10:30:00Z",
  "started_at": "2026-02-08T10:30:01Z",
  "completed_at": "2026-02-08T10:32:15Z",
  "error": null,
  "progress": 100,
  "podcast": {
    "title": "Example Podcast Episode",
    "source_url": "https://www.xiaoyuzhou.fm/episode/12345678",
    "audio_path": "/home/user/downloads/Example Podcast Episode/podcast.m4a",
    "cover_path": "/home/user/downloads/Example Podcast Episode/cover.jpg",
    "cover_format": "jpg",
    "shownotes_path": "/home/user/downloads/Example Podcast Episode/shownotes.txt",
    "downloaded_at": "2026-02-08T10:32:15Z",
    "file_size_mb": 45.2
  }
}
```

### 3. List Downloaded Podcasts

**Request**:

```bash
curl "http://localhost:8080/podcasts?limit=10&offset=0"
```

**Response**:

```json
{
  "podcasts": [
    {
      "title": "Example Podcast Episode",
      "audio_url": "/home/user/downloads/Example Podcast Episode/podcast.m4a",
      "shownotes": "This is the show notes content...",
      "cover_image": "/home/user/downloads/Example Podcast Episode/cover.jpg"
    }
  ],
  "total": 15,
  "limit": 10,
  "offset": 0
}
```

## Testing

### Unit Tests

```bash
# Run all unit tests
go test ./internal/server/...
go test ./internal/taskmanager/...

# Run with coverage
go test -cover ./internal/server/...
go test -cover ./internal/taskmanager/...

# Run with race detector
go test -race ./internal/server/...
```

### Integration Tests

```bash
# Run integration tests (starts test server)
go test ./tests/integration/...

# Run with verbose output
go test -v ./tests/integration/...
```

### Manual Testing with HTTPie

```bash
# Install HTTPie: pip install httpie

# Submit task using test URLs
http POST localhost:8080/tasks url="$(head -n 1 examples/xiaoyuzhou_urls)"

# Get status
http GET localhost:8080/tasks/550e8400-e29b-41d4-a716-446655440000

# List podcasts
http GET localhost:8080/podcasts

# List with pagination
http GET localhost:8080/podcasts limit==5 offset==10
```

### Batch Testing

Test multiple URLs from the examples file:

```bash
# Submit all test URLs
while IFS= read -r url; do
  echo "Submitting: $url"
  curl -X POST http://localhost:8080/tasks \
    -H "Content-Type: application/json" \
    -d "{\"url\": \"$url\"}"
  echo ""
done < examples/xiaoyuzhou_urls
```

## Project Structure

### Key Files to Modify During Development

```
cmd/podcast-server/main.go          # Server entry point
internal/server/server.go           # HTTP server setup
internal/server/handlers.go         # Request handlers
internal/taskmanager/manager.go     # Task lifecycle
internal/taskmanager/store.go       # In-memory storage
internal/taskmanager/scanner.go     # Directory scanner
examples/xiaoyuzhou_urls            # Test URLs for development
```

### Adding New Endpoints

1. Define handler in `internal/server/handlers.go`
2. Register route in `internal/server/server.go`
3. Add tests in `tests/unit/handler_test.go`
4. Update OpenAPI spec in `specs/003-podcast-api-server/contracts/openapi.yaml`

## Debugging

### Enable Verbose Logging

```bash
./podcast-server -verbose
```

### Check Logs

Logs are written to stdout with timestamps:

```
2026/02/08 10:30:00 [INFO] Starting server on :8080
2026/02/08 10:30:05 [DEBUG] POST /tasks - URL validated
2026/02/08 10:30:05 [INFO] Task created: 550e8400-e29b-41d4-a716-446655440000
2026/02/08 10:30:06 [DEBUG] Starting download for task 550e8400-e29b-41d4-a716-446655440000
```

### Common Issues

**Issue**: Port already in use
```bash
# Solution: Use different port
./podcast-server -port 3000
```

**Issue**: Permission denied writing to downloads directory
```bash
# Solution: Change ownership or permissions
chmod 755 downloads/
# Or use different directory
./podcast-server -downloads /tmp/podcasts
```

**Issue**: Task not found
```bash
# Solution: Check task ID is correct UUID format
# In-memory tasks are lost on server restart
```

## Performance Testing

### Load Testing with curl

```bash
# Submit 10 concurrent download tasks
for i in {1..10}; do
  curl -X POST http://localhost:8080/tasks \
    -H "Content-Type: application/json" \
    -d "{\"url\": \"https://www.xiaoyuzhou.fm/episode/123$i\"}" &
done
wait
```

### Benchmarking with Apache Bench

```bash
# Install: apt install apache2-utils (Linux)

# Benchmark task status endpoint
ab -n 1000 -c 10 http://localhost:8080/tasks/550e8400-e29b-41d4-a716-446655440000

# Benchmark list endpoint
ab -n 100 -c 5 http://localhost:8080/podcasts
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Test Podcast Server
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.25.5'
      - run: go test -v ./...
      - run: go build -o podcast-server cmd/podcast-server/main.go
```

## Deployment

### Development Deployment

```bash
# Build for current platform
go build -o podcast-server cmd/podcast-server/main.go

# Run in background
nohup ./podcast-server -port 8080 > server.log 2>&1 &

# Check running
ps aux | grep podcast-server

# Stop
pkill podcast-server
```

### Production Considerations

**Note**: This server is designed for local/private network use only (no authentication).

For production deployment:
1. Add reverse proxy (nginx/caddy) for HTTPS
2. Implement authentication (API keys, OAuth)
3. Add rate limiting
4. Configure systemd service for auto-restart
5. Set up log rotation
6. Monitor with health check endpoint

## Next Steps

1. **Implement Core Features**: Start with task submission endpoint
2. **Add Tests**: Write unit tests for each component
3. **Integration Testing**: Test full download flow
4. **Documentation**: Update API documentation as you implement
5. **Performance Testing**: Verify SC-001 to SC-008 targets are met

## Getting Help

- **API Documentation**: See `contracts/openapi.yaml`
- **Data Model**: See `data-model.md`
- **Technical Decisions**: See `research.md`
- **Feature Specification**: See `spec.md`

## Changelog

### Version 1.0.0 (Current)
- Initial release
- POST /tasks - Submit download tasks
- GET /tasks/{id} - Query task status
- GET /podcasts - List downloaded podcasts
- In-memory task management
- Directory scanning on startup
- Retry logic with exponential backoff
