# Module Interfaces: Podcast Audio Downloader

**Feature**: Podcast Audio Downloader (CLI tool)
**Date**: 2026-02-08
**Purpose**: Define Go package interfaces and APIs for internal modules

---

## Overview

This document defines the contracts (interfaces) for each internal module in the podcast downloader CLI. These interfaces specify the public API that each module exposes, enabling parallel development and testing.

---

## Module: `internal/downloader`

### Purpose
Core download logic including HTTP requests, file writing, and progress tracking.

### Interface: `URLExtractor`

Extracts audio file URLs from podcast episode web pages.

```go
package downloader

import "context"

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
    //
    // Errors:
    //   ErrInvalidURL - If pageURL is malformed
    //   ErrPageNotFound - If page returns 404 or similar
    //   ErrAudioNotFound - If no audio URL found on page
    //   ErrAccessDenied - If page requires authentication
    ExtractURL(ctx context.Context, pageURL string) (audioURL string, title string, err error)
}
```

### Interface: `FileDownloader`

Downloads audio files from URLs to local filesystem with progress tracking.

```go
package downloader

import (
    "context"
    "io"
)

// FileDownloader defines the interface for downloading files with progress tracking.
type FileDownloader interface {
    // Download fetches the audio file and writes it to the local filesystem.
    //
    // Parameters:
    //   ctx - Context for cancellation and timeout
    //   audioURL - Direct URL to the audio file
    //   filePath - Local filesystem path to write the file
    //   progress - Optional progress writer for real-time feedback (can be nil)
    //
    // Returns:
    //   int64 - Number of bytes downloaded
    //   error - Err if download fails
    //
    // Errors:
    //   ErrNetworkTimeout - If request times out
    //   ErrConnectionRefused - If cannot connect to server
    //   ErrDiskFull - If insufficient disk space
    //   ErrPermissionDenied - If cannot write to filePath
    //   ErrInvalidAudio - If downloaded content is not a valid audio file
    Download(ctx context.Context, audioURL, filePath string, progress io.Writer) (bytesWritten int64, err error)

    // ValidateFile checks if the downloaded file is a valid audio file.
    //
    // Parameters:
    //   filePath - Local filesystem path to validate
    //
    // Returns:
    //   error - Err if file is not a valid audio file
    //
    // Errors:
    //   ErrInvalidAudio - If file magic bytes don't match .m4a format
    //   ErrFileNotFound - If file does not exist
    ValidateFile(filePath string) error
}
```

### Implementation Note

Both interfaces should use the research decisions:
- **URLExtractor**: Uses `goquery` for HTML parsing
- **FileDownloader**: Uses `net/http` with retry logic and `schollz/progressbar/v3` for progress

---

## Module: `internal/validator`

### Purpose
Validates user input (URLs, file paths, configuration) before processing.

### Interface: `URLValidator`

Validates podcast episode URLs.

```go
package validator

import "regexp"

// URLValidator defines the interface for validating URLs.
type URLValidator interface {
    // ValidateURL checks if a URL is valid and meets requirements.
    //
    // Parameters:
    //   url - The URL to validate
    //
    // Returns:
    //   bool - True if URL is valid
    //   string - Error message if invalid, empty if valid
    ValidateURL(url string) (isValid bool, errMsg string)
}

// XiaoyuzhouURLValidator validates Xiaoyuzhou FM episode URLs.
type XiaoyuzhouURLValidator struct {
    // Pattern matches: *.xiaoyuzhoufm.com/episode/{episode_id}
    pattern *regexp.Regexp
}

// ValidateURL implements URLValidator for Xiaoyuzhou FM URLs.
func (v *XiaoyuzhouURLValidator) ValidateURL(url string) (bool, string) {
    // Implementation checks:
    // 1. URL is well-formed
    // 2. Domain matches *.xiaoyuzhoufm.com
    // 3. Path contains /episode/
    // 4. Episode ID is present
}
```

### Interface: `FilePathValidator`

Validates file paths for writing downloaded files.

```go
package validator

import "os"

// FilePathValidator defines the interface for validating file paths.
type FilePathValidator interface {
    // ValidatePath checks if a file path is valid and writable.
    //
    // Parameters:
    //   path - The file path to validate
    //   createIfMissing - Create parent directories if they don't exist
    //
    // Returns:
    //   error - Err if path is invalid or not writable
    //
    // Errors:
    //   ErrInvalidPath - If path contains invalid characters
    //   ErrPermissionDenied - If parent directory is not writable
    //   ErrDiskFull - If insufficient disk space (optional check)
    ValidatePath(path string, createIfMissing bool) error
}
```

---

## Module: `internal/models`

### Purpose
Data structures representing Episode and DownloadSession entities.

### Struct: `Episode`

```go
package models

// Episode represents a podcast episode with metadata.
type Episode struct {
    // ID is the unique episode identifier
    ID string

    // Title is the episode title (may be empty if not found)
    Title string

    // AudioURL is the direct URL to the .m4a audio file
    AudioURL string

    // PageURL is the original episode page URL
    PageURL string

    // FileSize is the audio file size in bytes (-1 if unknown)
    FileSize int64

    // Duration is the episode duration in seconds (0 if unknown)
    Duration int
}

// SanitizedTitle returns a filesystem-safe version of the title.
// Falls back to episode ID if title is empty after sanitization.
func (e *Episode) SanitizedTitle() string

// GenerateFilename generates a filename in the format:
// {sanitized_title}_{episode_id}.m4a
func (e *Episode) GenerateFilename() string

// Validate checks if all required fields are populated and valid.
func (e *Episode) Validate() error
```

### Struct: `DownloadSession`

```go
package models

import "time"

// Status represents the current state of a download session.
type Status int

const (
    StatusPending    Status = iota
    StatusInProgress
    StatusCompleted
    StatusFailed
)

// DownloadSession represents a single download operation.
type DownloadSession struct {
    // EpisodeID references the episode being downloaded
    EpisodeID string

    // FilePath is the local filesystem path for the downloaded file
    FilePath string

    // Status is the current download state
    Status Status

    // Progress is the download progress (0.0 to 1.0)
    Progress float64

    // BytesDownloaded is the number of bytes downloaded
    BytesDownloaded int64

    // TotalBytes is the total file size (-1 if unknown)
    TotalBytes int64

    // StartTime is when the download started
    StartTime time.Time

    // EndTime is when the download completed (zero if not completed)
    EndTime time.Time

    // Error contains the error message if download failed
    Error error

    // RetryCount is the number of retry attempts
    RetryCount int

    // DownloadSpeed is the current download speed in bytes/sec
    DownloadSpeed float64
}

// UpdateProgress updates the download progress based on bytes downloaded.
func (s *DownloadSession) UpdateProgress(bytesDownloaded int64)

// Complete marks the download as completed.
func (s *DownloadSession) Complete()

// Fail marks the download as failed with an error.
func (s *DownloadSession) Fail(err error)

// CanRetry checks if the download can be retried.
func (s *DownloadSession) CanRetry(maxRetries int) bool

// IncrementRetry increments the retry count.
func (s *DownloadSession) IncrementRetry()
```

---

## Module: `internal/config`

### Purpose
Configuration management for the CLI tool.

### Struct: `Config`

```go
package config

import "time"

// Config holds the application configuration.
type Config struct {
    // OutputDirectory is where downloaded files are saved
    OutputDirectory string

    // OverwriteExisting controls whether to overwrite existing files
    OverwriteExisting bool

    // Timeout is the HTTP request timeout
    Timeout time.Duration

    // MaxRetries is the maximum number of retry attempts for failed downloads
    MaxRetries int

    // RetryDelay is the base delay between retries (exponential backoff)
    RetryDelay time.Duration

    // ShowProgress controls whether to display download progress
    ShowProgress bool

    // ValidateFiles controls whether to validate downloaded audio files
    ValidateFiles bool
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() *Config {
    return &Config{
        OutputDirectory:   "./downloads",
        OverwriteExisting: false,
        Timeout:           30 * time.Second,
        MaxRetries:        3,
        RetryDelay:        1 * time.Second,
        ShowProgress:      true,
        ValidateFiles:     true,
    }
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
    // Check OutputDirectory is valid path
    // Check Timeout is positive
    // Check MaxRetries is non-negative
}
```

---

## Module: `pkg/httpclient`

### Purpose
Reusable HTTP client with retry logic and timeout configuration (for future web service use).

### Interface: `Client`

```go
package httpclient

import (
    "context"
    "net/http"
)

// Client defines the interface for HTTP requests with retry logic.
type Client interface {
    // Do executes an HTTP request with retry logic.
    //
    // Parameters:
    //   ctx - Context for cancellation and timeout
    //   req - HTTP request to execute
    //
    // Returns:
    //   *http.Response - HTTP response
    //   error - Err if all retry attempts fail
    Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// RetryableClient implements Client with exponential backoff retry.
type RetryableClient struct {
    // client is the underlying HTTP client
    client *http.Client

    // maxRetries is the maximum number of retry attempts
    maxRetries int

    // retryDelay is the base delay for exponential backoff
    retryDelay time.Duration
}

// NewRetryableClient creates a new client with retry logic.
func NewRetryableClient(timeout time.Duration, maxRetries int, retryDelay time.Duration) *RetryableClient
```

---

## Module: `cmd/podcast-downloader`

### Purpose
CLI application entry point integrating all modules.

### Function: `main`

```go
package main

import (
    "github.com/urfave/cli/v2"
)

// main is the application entry point.
func main() {
    app := &cli.App{
        Name:    "podcast-downloader",
        Usage:   "Download podcasts from Xiaoyuzhou FM",
        Version: "1.0.0",
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:    "output",
                Aliases: []string{"o"},
                Usage:   "Output directory for downloads",
                Value:   "./downloads",
                Action:  validateOutputDirectory,
            },
            &cli.BoolFlag{
                Name:    "overwrite",
                Aliases: []string{"f"},
                Usage:   "Overwrite existing files",
            },
            &cli.BoolFlag{
                Name:    "no-progress",
                Usage:   "Disable download progress bar",
            },
            &cli.IntFlag{
                Name:    "retry",
                Usage:   "Maximum number of retry attempts",
                Value:   3,
                Action:  validateRetryCount,
            },
            &cli.DurationFlag{
                Name:    "timeout",
                Usage:   "HTTP request timeout",
                Value:   30 * time.Second,
                Action:  validateTimeout,
            },
        },
        Action: downloadPodcast,
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

// downloadPodcast is the main action that orchestrates the download process.
func downloadPodcast(ctx *cli.Context) error {
    // 1. Validate URL argument
    // 2. Extract episode metadata (title, audio URL)
    // 3. Determine output file path
    // 4. Check for existing file
    // 5. Download with progress
    // 6. Validate downloaded file
    // 7. Report success/failure
}
```

---

## Integration Flow

```text
┌─────────────────┐
│  CLI Entry      │
│  (main.go)      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ URL Validator   │ ────[Valid]──┐
└─────────────────┘              │
                                 ▼
                        ┌─────────────────┐
                        │ URL Extractor   │ ────[Success]──┐
                        │ (goquery)       │                 │
                        └─────────────────┘                 │
                                                           ▼
                                                  ┌─────────────────┐
                                                  │ File Path       │ ────[Valid]──┐
                                                  │ Validator       │              │
                                                  └─────────────────┘              │
                                                                                  ▼
                                                                        ┌─────────────────┐
                                                                        │ File Downloader │
                                                                        │ (retry +        │
                                                                        │  progress)      │
                                                                        └────────┬────────┘
                                                                                 │
                                                                                 ▼
                                                                        ┌─────────────────┐
                                                                        │ File Validator  │
                                                                        │ (magic bytes)   │
                                                                        └────────┬────────┘
                                                                                 │
                                                                                 ▼
                                                                        ┌─────────────────┐
                                                                        │ Report Result   │
                                                                        │ (success/error) │
                                                                        └─────────────────┘
```

---

## Testing Contracts

Each module should include tests for its interface:

### `internal/downloader`
- Test URL extraction from mock HTML responses
- Test file download with progress tracking
- Test retry logic on network failures
- Test file validation

### `internal/validator`
- Test valid/invalid URL patterns
- Test file path validation
- Test directory creation

### `internal/models`
- Test episode title sanitization
- Test filename generation
- Test download session state transitions

### `internal/config`
- Test default configuration
- Test configuration validation
- Test flag parsing integration

---

## Next Steps

With module interfaces defined:
1. ✅ URLExtractor interface for HTML parsing
2. ✅ FileDownloader interface for HTTP downloads
3. ✅ Validator interfaces for input validation
4. ✅ Model structs for Episode and DownloadSession
5. ✅ Config struct for application settings
6. ✅ Main CLI entry point contract

Proceed to **quickstart.md** generation.
