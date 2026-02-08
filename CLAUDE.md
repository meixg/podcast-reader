# Podcast Reader Development Guidelines

Auto-generated from feature plans. Last updated: 2026-02-08 (Feature 2: Save Cover Images and Show Notes)

## Active Technologies

### Backend
- **Go 1.21+**: Primary language for CLI tools and web service backend
- **Standard Library**: `net/http`, `io`, `os`, `fmt`, `log`, `strings`, `regexp`, `time`, `math`

### CLI Dependencies
- **github.com/PuerkitoBio/goquery** (v1.8.1): HTML parsing for web scraping
- **github.com/urfave/cli/v2** (v2.27.1): CLI framework for command-line tools
- **github.com/schollz/progressbar/v3** (v3.14.1): Progress bar display for downloads
- **github.com/fatih/color** (v1.18.0): Colored console output

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

### Feature 2: Save Cover Images and Show Notes (2026-02-08)
**What it added**:
- Automatic cover image download alongside podcast audio
- Show notes extraction and plain text file generation
- Multi-fallback HTML extraction strategy for robustness
- UTF-8 with BOM encoding for international character support
- Graceful degradation (audio download continues even if cover/show notes fail)
- Podcast subdirectories with simplified filenames (podcast.m4a, cover.jpg, shownotes.txt)

**Technologies introduced**:
- github.com/fatih/color (v1.18.0) for colored console output
- No other new external dependencies (uses existing goquery)
- UTF-8 BOM encoding for cross-platform text file compatibility
- Magic byte detection for image format validation

**Architectural decisions**:
- Extended URLExtractor interface to return EpisodeMetadata struct (breaking change)
- New services: ImageDownloader, ShowNotesSaver
- Multi-fallback strategy for HTML element selection (aria-label → semantic selectors)
- HTML to plain text conversion with structure preservation (links, lists, headers)
- Image format preservation (JPEG, PNG, WebP) with format detection via magic bytes
- File organization: Podcast title subdirectories with simplified asset filenames

**Key implementation details**:
- Cover image extraction: `.avater-container` first image (Xiaoyuzhou FM specific)
- Show notes extraction: `<section aria-label="节目show notes">` → aria-label containing "show notes" → semantic selectors
- Image validation: Magic byte detection (JPEG: FF D8 FF, PNG: 89 50 4E 47, WebP: RIFF....WEBP)
- Text formatting: Links → "text (URL: url)", lists → bullets/numbers, headers → uppercase with underline
- Error handling: Detailed warning messages with specific reasons, no failure for non-critical assets
- HTTP timeouts: 1 hour for audio, 2 minutes for images

---

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
