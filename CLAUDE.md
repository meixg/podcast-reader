# Podcast Reader Development Guidelines

Auto-generated from feature plans. Last updated: 2026-02-08 (Feature 2: Save Cover Images and Show Notes)

## Active Technologies
- Go 1.25.5 + `net/http` (standard library), `github.com/PuerkitoBio/goquery` (existing), `github.com/fatih/color` (existing), existing downloader packages (3-podcast-api-server)
- In-memory for active tasks (lost on restart), filesystem for downloaded podcasts (scanned on startup), no database (3-podcast-api-server)
- File-based (existing downloads directory, in-memory task queue) (004-frontend-web-app)
- Go 1.21+ (existing project standard) (005-metadata-extraction)
- File-based (existing downloads directory structure) (005-metadata-extraction)
- Go 1.21+ (existing project standard) + Docker, GitHub Actions, docker/build-push-action (006-docker-packaging)
- File-based (downloads directory mounted as volume) (006-docker-packaging)

### Backend
- **Go 1.21+**: Primary language for CLI tools and web service backend
- **Standard Library**: `net/http`, `io`, `os`, `fmt`, `log`, `strings`, `regexp`, `time`, `math`

### CLI Dependencies
- **github.com/PuerkitoBio/goquery** (v1.8.1): HTML parsing for web scraping
- **github.com/urfave/cli/v2** (v2.27.1): CLI framework for command-line tools
- **github.com/schollz/progressbar/v3** (v3.14.1): Progress bar display for downloads
- **github.com/fatih/color** (v1.18.0): Colored console output

### Web Service (Future)
- **Frontend**: Vue 3 + TypeScript + Vite + Tailwind CSS (planned for web service)
  - Vue 3 Composition API for component logic
  - TypeScript strict mode for type safety
  - Tailwind CSS for utility-first styling
  - No external state management (use Vue 3 reactive state)
- **Backend**: RESTful APIs with Go (planned for web service)

## Project Structure

```text
# CLI Application (Current)
cmd/podcast-downloader/
├── main.go               # Application entry point
└── root.go               # CLI command definitions

internal/
├── downloader/           # Download logic
│   ├── downloader.go     # Main download service
│   ├── url_extractor.go  # HTML parsing and URL extraction
│   ├── metadata.go       # Episode metadata structures
│   ├── image_downloader.go   # Cover image download service
│   └── shownotes_saver.go    # Show notes file writer
├── models/               # Data structures
│   ├── episode.go        # Episode metadata
│   └── download_session.go
├── validator/            # Input validation
│   └── url_validator.go
└── config/               # Configuration
    └── config.go

pkg/
└── httpclient/           # Reusable HTTP client with retry logic
    └── client.go

downloads/                # Default download directory (configurable)
├── Podcast Title/        # Subdirectory per podcast
│   ├── podcast.m4a       # Downloaded audio files
│   ├── cover.jpg         # Downloaded cover images
│   └── shownotes.txt     # Saved show notes (UTF-8 with BOM)

# Web Service (Future - Planned Architecture)
backend/
├── cmd/server/           # Application entry point
├── internal/
│   ├── handlers/         # HTTP handlers
│   ├── services/         # Business logic
│   ├── models/           # Data structures
│   └── config/           # Configuration
├── pkg/                  # Public packages
├── storage/              # File-based storage
│   ├── downloads/
│   ├── transcripts/
│   └── briefings/
├── go.mod
└── go.sum

frontend/
├── src/
│   ├── components/       # Vue components
│   ├── views/           # Page components
│   ├── services/        # API calls
│   ├── composables/     # Vue 3 composables for shared logic
│   ├── types/           # TypeScript types
│   └── utils/           # Utilities
├── public/
├── package.json
├── vite.config.ts
├── tailwind.config.js
└── tsconfig.json
```

## Commands

### Go Development
```bash
# Using Make (recommended for server operations)
cd backend
make build              # Build server binary
make run                 # Run server directly
make restart            # Restart server (kill, rebuild, and start in background)
make clean              # Clean build artifacts and logs
make test               # Run tests
make fmt                # Format code
make lint               # Run linter

# Manual build commands (if Makefile not available)
go build -o podcast-server cmd/podcast-downloader/main.go

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run with race detector
go run -race cmd/podcast-downloader/main.go

# Format code
gofmt -w .

# Lint code
golint ./...
go vet ./...

# Install dependencies
go mod tidy
```

### Server Management
```bash
# IMPORTANT: Always use Makefile to restart server
cd backend
make restart

# View server logs
tail -f backend/output/server.log

# Check if server is running
ps aux | grep podcast-server | grep -v grep
```

### CLI Tool Usage
```bash
# Basic usage
./podcast-downloader "https://www.xiaoyuzhoufm.com/episode/..."

# With options
./podcast-downloader --output ~/music --overwrite --retry 5 "https://..."

# Show help
./podcast-downloader --help

# Show version
./podcast-downloader --version
```

## Code Style

### Go Code Standards

**Formatting**:
- Use `gofmt` for all code (automatic formatting)
- Follow Go conventions from [Effective Go](https://golang.org/doc/effective_go)
- Maximum line length: 100 characters (soft limit)

**Naming Conventions**:
- Package names: lowercase, single word, no underscores
- Exported functions/variables: PascalCase (e.g., `DownloadFile`)
- Private functions/variables: camelCase (e.g., `sanitizeFilename`)
- Constants: PascalCase or UPPER_CASE
- Interfaces: PascalCase with "er" suffix (e.g., `URLExtractor`)

**Error Handling**:
- Always handle errors explicitly (never ignore)
- Use wrapped errors with context: `fmt.Errorf("download failed: %w", err)`
- Return errors from functions, don't panic in production code
- Use custom error types for domain-specific errors

**Comments**:
- Exported functions must have doc comments
- Use godoc format: `// FunctionName does...`
- Add context for non-obvious code
- Comment complex algorithms and business logic

**Example**:
```go
// Package downloader provides functionality for downloading podcast episodes
// from Xiaoyuzhou FM with progress tracking and retry logic.
package downloader

// URLExtractor defines the interface for extracting metadata from podcast pages.
type URLExtractor interface {
	// ExtractURL fetches the episode page and extracts metadata.
	//
	// Parameters:
	//   ctx - Context for cancellation and timeout
	//   pageURL - The episode page URL to scrape
	//
	// Returns:
	//   *EpisodeMetadata - Contains audio URL, cover URL, show notes, title, etc.
	//   error - Err if page cannot be fetched or required data not found
	ExtractURL(ctx context.Context, pageURL string) (*EpisodeMetadata, error)
}
```

**Testing**:
- Use table-driven tests for multiple test cases
- Test both success and failure paths
- Mock external dependencies (HTTP calls, file I/O)
- Use `t.Run()` for subtests

## Recent Changes
- 006-docker-packaging: Added Go 1.21+ (existing project standard) + Docker, GitHub Actions, docker/build-push-action
- 005-metadata-extraction: Added Go 1.21+ (existing project standard)
- 004-frontend-web-app: Added File-based (existing downloads directory, in-memory task queue)

### Feature 2: Save Cover Images and Show Notes (2026-02-08)
**What it added**:

**Technologies introduced**:

**Architectural decisions**:

**Key implementation details**:

---

### Feature 1: Podcast Audio Downloader (2026-02-08)
**What it added**:

**Technologies introduced**:

**Architectural decisions**:

---

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
