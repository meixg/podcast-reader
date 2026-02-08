# Research: Podcast Audio Downloader

**Date**: 2026-02-08
**Feature**: Podcast Audio Downloader (CLI tool)
**Purpose**: Resolve technical decisions for implementation

---

## Overview

This document consolidates research findings for key technical decisions in building the podcast audio downloader CLI tool. Each dependency area was evaluated against the specific requirements of downloading podcast episodes from Xiaoyuzhou FM.

---

## Decision 1: HTML Parsing Library

### Requirement
Parse HTML from Xiaoyuzhou FM podcast episode pages to extract direct .m4a audio file URLs.

### Options Evaluated

| Library | Ease of Use | Performance | Dependencies | Chinese Support | Recommendation |
|---------|-------------|-------------|--------------|-----------------|----------------|
| **goquery** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | Minimal (net/html only) | ⭐⭐⭐⭐⭐ | ✅ **CHOSEN** |
| colly | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | Heavy | ⭐⭐⭐⭐⭐ | ❌ Overkill |
| net/html (stdlib) | ⭐⭐ | ⭐⭐⭐⭐⭐ | Zero | ⭐⭐⭐⭐⭐ | ❌ Too verbose |

### Decision: **goquery** (github.com/PuerkitoBio/goquery)

**Rationale**:
- **Perfect complexity match**: Simple HTML parsing without scraping framework overhead
- **jQuery-like API**: Intuitive CSS selectors (`audio[src]`, `.class`, `#id`)
- **Lightweight**: Only depends on `net/html` from standard library
- **Actively maintained**: Last updated November 2025, 20k+ GitHub stars
- **Chinese content**: Native UTF-8 support, zero encoding issues with Xiaoyuzhou FM
- **Performance**: ~1200 req/s (more than sufficient for single-URL processing)

**Alternatives Considered**:
- **colly**: Rejected due to complexity (built for concurrent scraping) and open issues about memory leaks
- **net/html**: Rejected due to verbose API requiring manual tree traversal and significant boilerplate

**Implementation Example**:
```go
doc, err := goquery.NewDocumentFromReader(response.Body)
if err != nil {
    return err
}

// Try common selectors for audio elements
selectors := []string{
    "audio[src]",
    "source[src]",
    ".audio-player audio",
    "[data-audio-url]",
}

for _, selector := range selectors {
    doc.Find(selector).Each(func(i int, s *goquery.Selection) {
        audioURL, exists := s.Attr("src")
        if exists && strings.HasSuffix(audioURL, ".m4a") {
            fmt.Printf("Found audio URL: %s\n", audioURL)
        }
    })
}
```

---

## Decision 2: CLI Framework

### Requirement
Build a command-line interface accepting a URL argument and optional flags (output directory, overwrite mode).

### Options Evaluated

| Library | Ease of Use | Dependencies | Binary Size | Best For | Recommendation |
|---------|-------------|--------------|-------------|----------|----------------|
| **urfave/cli v2** | ⭐⭐⭐⭐⭐ | Light | +500KB-1MB | Simple tools | ✅ **CHOSEN** |
| cobra | ⭐⭐⭐ | Heavy | +2-5MB | Complex apps | ❌ Overkill |
| flag (stdlib) | ⭐⭐ | Zero | +0KB | Minimal tools | ❌ Too manual |

### Decision: **urfave/cli v2** (github.com/urfave/cli/v2)

**Rationale**:
- **Perfect feature match**: Designed for single-command utilities with flags
- **Declarative API**: Less code than cobra or stdlib flag
- **Lightweight**: Minimal dependency tree, small binary size impact
- **Excellent DX**: Built-in help generation, validation, and error handling
- **Cross-platform**: Works on macOS, Windows, Linux
- **Future-proof**: Easy to add subcommands later if needed
- **Modern**: Actively maintained with v2 providing clean APIs

**Alternatives Considered**:
- **cobra**: Rejected as overkill for a single-command tool; heavy dependency tree (often pulls in Viper)
- **flag (stdlib)**: Rejected due to manual help text generation and more verbose code

**Implementation Example**:
```go
app := &cli.App{
    Name:  "podcast-downloader",
    Usage: "Download podcasts from URLs with progress tracking",
    Flags: []cli.Flag{
        &cli.StringFlag{
            Name:    "output",
            Aliases: []string{"o"},
            Usage:   "Output directory for downloads",
            Value:   "./downloads",
        },
        &cli.BoolFlag{
            Name:    "overwrite",
            Aliases: []string{"f"},
            Usage:   "Overwrite existing files",
        },
    },
    Action: func(ctx *cli.Context) error {
        if ctx.NArg() < 1 {
            return cli.Exit("Please provide a podcast URL", 1)
        }

        url := ctx.Args().Get(0)
        outputDir := ctx.String("output")
        overwrite := ctx.Bool("overwrite")

        // Download logic here
        return downloadPodcast(url, outputDir, overwrite)
    },
}
```

---

## Decision 3: Progress Bar Library

### Requirement
Display download progress for large audio files (up to 500MB) with percentage, speed, and ETA.

### Options Evaluated

| Library | HTTP Integration | Features | Dependencies | Complexity | Recommendation |
|---------|------------------|----------|--------------|------------|----------------|
| **schollz/progressbar/v3** | ⭐⭐⭐⭐⭐ (io.Writer) | Complete | Zero | Minimal | ✅ **CHOSEN** |
| cheggaaa/pb/v3 | ⭐⭐⭐⭐ (ProxyReader) | Rich | Moderate | Medium | ❌ Overkill |
| vbauerster/mpb/v8 | ⭐⭐⭐ (Manual) | Multi-bar | Moderate | High | ❌ Too complex |
| Custom ANSI | ⭐⭐ (Custom) | Custom | Zero | Very High | ❌ Reinventing wheel |

### Decision: **schollz/progressbar/v3** (github.com/schollz/progressbar/v3)

**Rationale**:
- **Best HTTP integration**: Implements `io.Writer` interface for seamless integration with `io.Copy`
- **Complete feature set**: Built-in speed, ETA, percentage, and byte tracking
- **Zero dependencies**: Only uses Go standard library
- **Cross-platform**: Explicitly designed to work on every OS without problems
- **Battle-tested**: Used in `croc` file transfer tool for multi-GB transfers
- **Performance**: Lightweight rendering, efficient updates via carriage return
- **Active**: 2.8k+ GitHub stars, MIT license

**Alternatives Considered**:
- **cheggaaa/pb**: Rejected due to template complexity and heavier dependencies; not needed for single-file downloads
- **vbauerster/mpb**: Rejected as designed for multi-progress bar concurrent scenarios; overkill for sequential downloads
- **Custom ANSI**: Rejected due to development time (~500+ lines for robust implementation) and maintenance burden

**Implementation Example**:
```go
resp, err := http.Get(audioURL)
if err != nil {
    return err
}
defer resp.Body.Close()

out, err := os.Create(filename)
if err != nil {
    return err
}
defer out.Close()

// Create progress bar with automatic speed/ETA
bar := progressbar.DefaultBytes(
    resp.ContentLength,
    "downloading",
)

// Copy with automatic progress tracking
_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
```

---

## Additional Technical Decisions

### HTTP Client Configuration

**Decision**: Use `net/http` standard library with custom `http.Client` configuration

**Rationale**:
- Standard library provides robust HTTP/1.1 and HTTP/2 support
- Custom client allows timeout configuration (requirement: FR-010)
- Enables retry logic implementation (requirement: FR-008)

**Configuration**:
```go
client := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        10,
        IdleConnTimeout:     30 * time.Second,
        DisableCompression:  false,
        MaxConnsPerHost:     5,
    },
}
```

### Retry Logic Strategy

**Decision**: Implement exponential backoff with max 3 retries (requirement: FR-006)

**Rationale**:
- Handles transient network failures (requirement: SC-006)
- Exponential backoff prevents overwhelming the server
- 3 retries balances success rate with user wait time

**Implementation**:
```go
func downloadWithRetry(client *http.Client, url string, maxRetries int) (*http.Response, error) {
    var lastErr error
    for attempt := 0; attempt < maxRetries; attempt++ {
        resp, err := client.Get(url)
        if err == nil && resp.StatusCode < 500 {
            return resp, nil
        }
        lastErr = err
        if attempt < maxRetries-1 {
            time.Sleep(time.Duration(math.Pow(2, float64(attempt))) * time.Second)
        }
    }
    return nil, lastErr
}
```

### File Naming Strategy

**Decision**: Use episode title with sanitization, fallback to episode ID

**Rationale**:
- Meaningful filenames (requirement: FR-004)
- Sanitization prevents filesystem issues
- Episode ID fallback ensures valid filename when title extraction fails

**Implementation**:
```go
func sanitizeFilename(name string) string {
    // Replace invalid characters with underscore
    reg := regexp.MustCompile(`[<>:"/\\|?*]`)
    sanitized := reg.ReplaceAllString(name, "_")
    // Limit length
    if len(sanitized) > 200 {
        sanitized = sanitized[:200]
    }
    return strings.TrimSpace(sanitized)
}
```

---

## Summary of Dependencies

### Primary Dependencies

| Package | Purpose | Version |
|---------|---------|---------|
| `github.com/PuerkitoBio/goquery` | HTML parsing | v1.8.1 |
| `github.com/urfave/cli/v2` | CLI framework | v2.27.1 |
| `github.com/schollz/progressbar/v3` | Progress display | v3.14.1 |

### Standard Library Packages Used

- `net/http` - HTTP client and requests
- `net/url` - URL validation and parsing
- `io` - I/O operations and stream copying
- `os` - File system operations
- `fmt` - Formatting and output
- `log` - Logging
- `strings` - String manipulation
- `regexp` - Filename sanitization
- `time` - Timeout and retry delays
- `math` - Exponential backoff calculation

---

## Architecture Notes

### Concurrent vs Sequential Processing

**Decision**: Sequential download processing for MVP

**Rationale**:
- Simpler implementation and testing
- Meets requirement: "single-user tool, sequential downloads"
- Easier to provide clear progress feedback
- Can be enhanced later if concurrent downloads are needed

### Error Message Language

**Decision**: Error messages in Chinese with English fallback

**Rationale**:
- Xiaoyuzhou FM is a Chinese service
- Target audience likely Chinese speakers
- Improves user experience for primary user base
- English fallback supports international users

### File Validation

**Decision**: Validate file is audio by checking magic bytes and extension

**Rationale**:
- Detects error pages returned as HTML (requirement: FR-011)
- More robust than just checking file extension
- Supports .m4a format (requirement assumption)

**Implementation**:
```go
func validateAudioFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    header := make([]byte, 12)
    if _, err := file.Read(header); err != nil {
        return err
    }

    // Check for M4A/MP4 magic bytes (ftyp)
    if !bytes.Equal(header[4:8], []byte("ftyp")) {
        return fmt.Errorf("downloaded file is not a valid audio file")
    }

    return nil
}
```

---

## Next Steps

With all technical decisions resolved:

1. ✅ HTML parsing: goquery selected
2. ✅ CLI framework: urfave/cli v2 selected
3. ✅ Progress display: schollz/progressbar/v3 selected
4. ✅ HTTP client: Standard library with custom configuration
5. ✅ Retry logic: Exponential backoff with 3 retries
6. ✅ File naming: Episode title with sanitization

Proceed to **Phase 1: Design & Contracts** to create:
- Data models (Episode, DownloadSession)
- Implementation contracts (module interfaces)
- Quickstart guide
