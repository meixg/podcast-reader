# Podcast Reader Development Guidelines

Auto-generated from feature plans. Last updated: 2026-02-08

## Active Technologies

### Backend
- **Go 1.21+**: Primary language for CLI tools and web service backend
- **Standard Library**: `net/http`, `io`, `os`, `fmt`, `log`, `strings`, `regexp`, `time`, `math`

### CLI Dependencies
- **github.com/PuerkitoBio/goquery** (v1.8.1): HTML parsing for web scraping
- **github.com/urfave/cli/v2** (v2.27.1): CLI framework for command-line tools
- **github.com/schollz/progressbar/v3** (v3.14.1): Progress bar display for downloads

### Web Service (Future)
- **Frontend**: Vue 3 + TypeScript + Vite (planned for web service)
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
│   └── progress.go       # Progress tracking
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
│   ├── stores/          # Pinia stores
│   ├── types/           # TypeScript types
│   └── utils/           # Utilities
├── public/
├── package.json
├── vite.config.ts
└── tsconfig.json
```

## Commands

### Go Development
```bash
# Build CLI tool
go build -o podcast-downloader cmd/podcast-downloader/main.go

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

// URLExtractor defines the interface for extracting audio URLs from web pages.
type URLExtractor interface {
    // ExtractURL fetches the episode page and extracts the direct audio file URL.
    //
    // Parameters:
    //   ctx - Context for cancellation and timeout
    //   pageURL - The episode page URL to scrape
    //
    // Returns:
    //   string - The direct audio file URL (.m4a)
    //   string - The episode title for filename generation
    //   error - Err if page cannot be fetched or audio URL not found
    ExtractURL(ctx context.Context, pageURL string) (audioURL string, title string, err error)
}

// ExtractURL implements the URLExtractor interface using goquery for HTML parsing.
func (e *HTMLExtractor) ExtractURL(ctx context.Context, pageURL string) (string, string, error) {
    // Implementation...
}
```

**Testing**:
- Use table-driven tests for multiple test cases
- Test both success and failure paths
- Mock external dependencies (HTTP calls, file I/O)
- Use `t.Run()` for subtests

**Example**:
```go
func TestSanitizeFilename(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "removes invalid characters",
            input:    "file<>name",
            expected: "file__name",
        },
        {
            name:     "truncates long names",
            input:    string(make([]byte, 250)),
            expected: string(make([]byte, 200)),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := sanitizeFilename(tt.input)
            if result != tt.expected {
                t.Errorf("sanitizeFilename() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Recent Changes

### Feature 1: Podcast Audio Downloader (2026-02-08)
**What it added**:
- CLI tool for downloading podcast episodes from Xiaoyuzhou FM
- HTML parsing using goquery to extract audio URLs
- HTTP client with retry logic and exponential backoff
- Progress bar display for large file downloads
- File validation using magic bytes
- URL and file path validation

**Technologies introduced**:
- github.com/PuerkitoBio/goquery (HTML parsing)
- github.com/urfave/cli/v2 (CLI framework)
- github.com/schollz/progressbar/v3 (progress display)

**Architectural decisions**:
- CLI tool instead of web service (constitutional exception justified for utility tool)
- Sequential download processing (simpler for MVP)
- No database persistence (ephemeral in-memory data)
- File-based storage for downloaded audio

---

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
