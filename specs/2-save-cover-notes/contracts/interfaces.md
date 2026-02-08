# Interface Contracts: Save Cover Images and Show Notes

**Feature**: 2-save-cover-notes
**Date**: 2026-02-08
**Type**: Go Interface Definitions

## Overview

This document defines the Go interfaces for the enhanced podcast downloader. Since this is a CLI tool (not a web service), these are internal service interfaces rather than REST API contracts.

## Core Interfaces

### URLExtractor (Enhanced)

Extracts all metadata (audio, cover, show notes) from a podcast episode page.

```go
package downloader

import "context"

// EpisodeMetadata contains all extractable metadata from a podcast episode page.
type EpisodeMetadata struct {
    // AudioURL is the direct URL to the audio file (.m4a, .mp3, etc.)
    AudioURL string

    // CoverURL is the URL to the cover image (JPEG, PNG, WEBP). May be empty.
    CoverURL string

    // ShowNotes is the plain text content of the show notes. May be empty.
    ShowNotes string

    // Title is the episode title used for filename generation.
    Title string

    // EpisodeNumber is the episode number if available. May be empty.
    EpisodeNumber string

    // PodcastName is the podcast/series name. May be empty.
    PodcastName string

    // PublicationDate is when the episode was published. Zero if not available.
    PublicationDate time.Time
}

// URLExtractor defines the interface for extracting metadata from podcast pages.
type URLExtractor interface {
    // ExtractURL fetches the episode page and extracts all metadata.
    //
    // Parameters:
    //   ctx - Context for cancellation and timeout
    //   pageURL - The episode page URL to scrape
    //
    // Returns:
    //   *EpisodeMetadata - All extracted metadata (audio, cover, show notes, title, etc.)
    //   error - Error if page cannot be fetched or essential data (audio URL) not found.
    //           Non-essential data (cover, show notes) may be missing without error.
    ExtractURL(ctx context.Context, pageURL string) (*EpisodeMetadata, error)
}
```

**Error Behavior**:
- Returns error only if page fetch fails or audio URL is not found
- Cover image and show notes may be empty (no error) - this is graceful degradation
- Wraps errors with context: `fmt.Errorf("extract metadata: %w", err)`

**Thread Safety**: Not thread-safe (create new instance per goroutine if needed)

---

### ImageDownloader

Downloads cover images with progress tracking and validation.

```go
package downloader

import (
    "context"
    "io"
)

// ImageDownloader defines the interface for downloading cover images.
type ImageDownloader interface {
    // Download fetches the cover image and writes it to the local filesystem.
    //
    // Parameters:
    //   ctx - Context for cancellation and timeout
    //   imageURL - The direct image URL to download
    //   filePath - Local file path where the image should be saved
    //   progress - Optional writer for progress updates (can be nil)
    //
    // Returns:
    //   int64 - Number of bytes written
    //   error - Error if download fails, validation fails, or context is cancelled
    Download(ctx context.Context, imageURL, filePath string, progress io.Writer) (int64, error)

    // ValidateImage checks if the downloaded file is a valid image.
    //
    // Parameters:
    //   filePath - Local file path to validate
    //
    // Returns:
    //   string - Detected image format ("jpg", "png", "webp", "gif")
    //   error - Error if file is invalid, corrupted, or not an image
    ValidateImage(filePath string) (string, error)
}
```

**Error Behavior**:
- Returns error on network failures (timeout, connection refused, 404, etc.)
- Returns error if downloaded file is not a valid image (magic byte validation)
- Implements retry logic with exponential backoff for transient failures
- Wraps errors with context: `fmt.Errorf("download cover: %w", err)`

**Thread Safety**: Not thread-safe (create new instance per goroutine if needed)

---

### ShowNotesSaver

Formats and saves show notes content to text files.

```go
package downloader

// ShowNotesSaver defines the interface for saving show notes.
type ShowNotesSaver interface {
    // Save formats HTML content and saves it as a UTF-8-BOM text file.
    //
    // Parameters:
    //   content - Raw HTML content of the show notes section
    //   filePath - Local file path where the text file should be saved
    //
    // Returns:
    //   int64 - Number of bytes written
    //   error - Error if formatting fails or file cannot be written
    Save(content, filePath string) (int64, error)

    // FormatHTMLToText converts HTML content to readable plain text.
    //
    // This method:
    // - Converts links to "text (URL: url)" format
    // - Converts lists to bullet/numbered lists
    // - Converts headers to uppercase with underlines
    // - Preserves blockquotes with "> " prefix
    // - Strips excessive whitespace
    //
    // Parameters:
    //   htmlContent - Raw HTML content to format
    //
    // Returns:
    //   string - Formatted plain text
    FormatHTMLToText(htmlContent string) string
}
```

**Error Behavior**:
- Returns error if file cannot be created (permission denied, disk full, etc.)
- Returns error if content contains invalid UTF-8 sequences
- Never returns error for empty content (creates empty file)
- Wraps errors with context: `fmt.Errorf("save show notes: %w", err)`

**Thread Safety**: Not thread-safe (create new instance per goroutine if needed)

---

### FileDownloader (Existing, Unchanged)

Downloads audio files with progress tracking and validation.

```go
package downloader

import (
    "context"
    "io"
)

// FileDownloader defines the interface for downloading audio files.
type FileDownloader interface {
    Download(ctx context.Context, audioURL, filePath string, progress io.Writer) (int64, error)
    ValidateFile(filePath string) error
}
```

**Note**: This interface already exists and is not modified. The new `ImageDownloader` interface follows the same pattern for consistency.

---

## Implementation Contracts

### HTMLExtractor Implementation

```go
package downloader

import (
    "context"
    "github.com/PuerkitoBio/goquery"
)

// HTMLExtractor implements URLExtractor using goquery for HTML parsing.
type HTMLExtractor struct {
    client Doer
}

// Doer is the interface for HTTP GET requests (goquery.Document).
type Doer interface {
    Get(url string) (*goquery.Document, error)
}

// NewHTMLExtractor creates a new HTML extractor.
func NewHTMLExtractor(client Doer) *HTMLExtractor {
    return &HTMLExtractor{client: client}
}

// ExtractURL implements URLExtractor interface.
func (e *HTMLExtractor) ExtractURL(ctx context.Context, pageURL string) (*EpisodeMetadata, error) {
    doc, err := e.client.Get(pageURL)
    if err != nil {
        return nil, fmt.Errorf("fetch page: %w", err)
    }

    metadata := &EpisodeMetadata{}

    // Extract audio URL (required)
    metadata.AudioURL, err = e.extractAudioURL(doc)
    if err != nil {
        return nil, err // Audio is required
    }

    // Extract title (required)
    metadata.Title = e.extractTitle(doc)
    if metadata.Title == "" {
        metadata.Title = "episode" // Fallback
    }

    // Extract cover URL (optional)
    metadata.CoverURL, _ = e.extractCoverURL(doc)
    // No error if cover not found - graceful degradation

    // Extract show notes (optional)
    metadata.ShowNotes, _ = e.extractShowNotes(doc)
    // No error if show notes not found - graceful degradation

    // Extract optional metadata
    metadata.EpisodeNumber = e.extractEpisodeNumber(doc)
    metadata.PodcastName = e.extractPodcastName(doc)
    metadata.PublicationDate = e.extractPublicationDate(doc)

    return metadata, nil
}

// extractCoverURL extracts cover image URL from avatar-container.
// Cover images are guaranteed to exist on all podcast episodes.
func (e *HTMLExtractor) extractCoverURL(doc *goquery.Document) (string, error) {
    // avatar-container class selector (Xiaoyuzhou FM specific)
    // IMPORTANT: This container has TWO images - first is episode cover, second is podcast account cover
    // We only want the first image (episode cover)
    if selection := doc.Find(".avatar-container img").First(); selection.Length() > 0 {
        if src, exists := selection.Attr("src"); exists && src != "" {
            return src, nil
        }
    }

    return "", ErrCoverNotFound
}

// extractShowNotes extracts show notes using multi-fallback strategy.
func (e *HTMLExtractor) extractShowNotes(doc *goquery.Document) (string, error) {
    // Try 1: Exact aria-label match
    if selection := doc.Find("section[aria-label='节目show notes']").First(); selection.Length() > 0 {
        return selection.Html(), nil
    }

    // Try 2: Partial aria-label match (case-insensitive)
    // ... (implementation details)

    // Try 3: Semantic selectors
    // ... (implementation details)

    return "", ErrShowNotesNotFound
}
```

**Contract Guarantees**:
- Always returns non-nil `EpisodeMetadata` if no error
- `AudioURL` is always set and valid if no error
- `Title` is always non-empty if no error (falls back to "episode")
- `CoverURL` and `ShowNotes` may be empty (not an error)
- All methods are idempotent (same input → same output)

---

### HTTPImageDownloader Implementation

```go
package downloader

import (
    "context"
    "io"
    "net/http"
)

// HTTPImageDownloader implements ImageDownloader with HTTP client.
type HTTPImageDownloader struct {
    client       *http.Client
    maxImageSize int64 // Maximum image size in bytes (default: 10MB)
}

// NewHTTPImageDownloader creates a new image downloader.
func NewHTTPImageDownloader(client *http.Client, maxSize int64) *HTTPImageDownloader {
    return &HTTPImageDownloader{
        client:       client,
        maxImageSize: maxSize,
    }
}

// Download implements ImageDownloader interface.
func (d *HTTPImageDownloader) Download(ctx context.Context, imageURL, filePath string, progress io.Writer) (int64, error) {
    // Implementation similar to HTTPDownloader.Download()
    // - Create HTTP request with context
    // - Execute with retry logic
    // - Check content length <= maxImageSize
    // - Stream to file with progress tracking
    // - Validate after download
    // ... (implementation details)
    return 0, nil
}

// ValidateImage implements ImageDownloader interface.
func (d *HTTPImageDownloader) ValidateImage(filePath string) (string, error) {
    // Read first 12 bytes for magic byte detection
    file, err := os.Open(filePath)
    if err != nil {
        return "", fmt.Errorf("open file: %w", err)
    }
    defer file.Close()

    header := make([]byte, 12)
    if _, err := io.ReadFull(file, header); err != nil {
        return "", fmt.Errorf("read header: %w", err)
    }

    format := d.detectFormat(header)
    if format == "" {
        return "", ErrInvalidImage
    }

    // Verify file size is reasonable
    stat, err := file.Stat()
    if err != nil {
        return "", fmt.Errorf("stat file: %w", err)
    }

    if stat.Size() > d.maxImageSize {
        return "", fmt.Errorf("%w: size %d exceeds limit %d", ErrImageTooLarge, stat.Size(), d.maxImageSize)
    }

    return format, nil
}

// detectFormat identifies image format from magic bytes.
func (d *HTTPImageDownloader) detectFormat(data []byte) string {
    if len(data) < 4 {
        return ""
    }

    // JPEG: FF D8 FF
    if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
        return "jpg"
    }

    // PNG: 89 50 4E 47
    if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
        return "png"
    }

    // WebP: RIFF....WEBP
    if len(data) >= 12 &&
        data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 &&
        data[8] == 0x57 && data[9] == 0x45 && data[10] == 0x42 && data[11] == 0x50 {
        return "webp"
    }

    return ""
}
```

**Contract Guarantees**:
- Downloads complete images or returns error (no partial files on error)
- Validates image format using magic bytes (not just file extension)
- Enforces maximum file size limit
- Implements retry logic for transient network failures
- Respects context cancellation

---

### PlainTextShowNotesSaver Implementation

```go
package downloader

import (
    "bytes"
    "unicode/utf8"
)

// PlainTextShowNotesSaver implements ShowNotesSaver.
type PlainTextShowNotesSaver struct{}

// NewPlainTextShowNotesSaver creates a new show notes saver.
func NewPlainTextShowNotesSaver() *PlainTextShowNotesSaver {
    return &PlainTextShowNotesSaver{}
}

// Save implements ShowNotesSaver interface.
func (s *PlainTextShowNotesSaver) Save(content, filePath string) (int64, error) {
    // Format HTML to plain text
    textContent := s.FormatHTMLToText(content)

    // Validate UTF-8 encoding
    if !utf8.ValidString(textContent) {
        return 0, fmt.Errorf("%w: show notes contains invalid UTF-8", ErrInvalidEncoding)
    }

    // Create file with UTF-8 BOM
    file, err := os.Create(filePath)
    if err != nil {
        return 0, fmt.Errorf("create file: %w", err)
    }
    defer file.Close()

    // Write UTF-8 BOM (0xEF 0xBB 0xBF)
    bom := []byte{0xEF, 0xBB, 0xBF}
    if _, err := file.Write(bom); err != nil {
        return 0, fmt.Errorf("write BOM: %w", err)
    }

    // Write content
    bytesWritten, err := file.WriteString(textContent)
    if err != nil {
        return 0, fmt.Errorf("write content: %w", err)
    }

    return int64(bytesWritten + len(bom)), nil
}

// FormatHTMLToText implements ShowNotesSaver interface.
func (s *PlainTextShowNotesSaver) FormatHTMLToText(htmlContent string) string {
    doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(htmlContent)))
    if err != nil {
        return htmlContent // Fallback to original HTML
    }

    // Convert links
    doc.Find("a").Each(func(i int, s *goquery.Selection) {
        text := s.Text()
        if href, exists := s.Attr("href"); exists && href != "" {
            replacement := fmt.Sprintf("%s (URL: %s)", text, href)
            s.ReplaceWithHtml(replacement)
        }
    })

    // Convert lists
    doc.Find("ul").Each(func(i int, s *goquery.Selection) {
        s.Find("li").Each(func(j int, li *goquery.Selection) {
            text := li.Text()
            li.ReplaceWithHtml(fmt.Sprintf("• %s", text))
        })
        s.AfterHtml("\n")
    })

    doc.Find("ol").Each(func(i int, s *goquery.Selection) {
        s.Find("li").Each(func(j int, li *goquery.Selection) {
            text := li.Text()
            li.ReplaceWithHtml(fmt.Sprintf("%d. %s", j+1, text))
        })
        s.AfterHtml("\n")
    })

    // Convert headers
    for i := 1; i <= 6; i++ {
        selector := fmt.Sprintf("h%d", i)
        doc.Find(selector).Each(func(j int, s *goquery.Selection) {
            text := s.Text()
            underline := strings.Repeat("=", len(text)) // h1 and h2
            if i >= 3 {
                underline = "" // h3-h6 no underline
            }
            if underline != "" {
                s.ReplaceWithHtml(fmt.Sprintf("%s\n%s\n\n", text, underline))
            } else {
                s.ReplaceWithHtml(fmt.Sprintf("%s\n\n", text))
            }
        })
    }

    // Get text and clean up whitespace
    text := doc.Text()

    // Remove excessive whitespace
    text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
    text = strings.TrimSpace(text)

    return text
}
```

**Contract Guarantees**:
- Always writes UTF-8-BOM encoded files
- Validates UTF-8 encoding before writing
- Preserves links, lists, and headers in plain text format
- Handles malformed HTML gracefully (fallback to original)
- Creates parent directories if needed

---

## Integration Contract

### Main Download Workflow

```go
package main

func downloadEpisode(ctx context.Context, episodeURL, outputDir string) error {
    // 1. Extract metadata
    extractor := downloader.NewHTMLExtractor(httpClient)
    metadata, err := extractor.ExtractURL(ctx, episodeURL)
    if err != nil {
        return fmt.Errorf("extract metadata: %w", err)
    }

    // 2. Generate filenames
    baseFilename := sanitizer.SanitizeFilename(metadata.Title)
    audioPath := filepath.Join(outputDir, baseFilename+".m4a")
    coverPath := filepath.Join(outputDir, baseFilename+".jpg")
    showNotesPath := filepath.Join(outputDir, baseFilename+".txt")

    // 3. Download audio (required)
    audioDownloader := downloader.NewHTTPDownloader(httpClient, true)
    progress := newProgressBar("Downloading audio...")
    if _, err := audioDownloader.Download(ctx, metadata.AudioURL, audioPath, progress); err != nil {
        return fmt.Errorf("download audio: %w", err)
    }

    // 4. Download cover (optional, with graceful degradation)
    if metadata.CoverURL != "" {
        imageDownloader := downloader.NewHTTPImageDownloader(httpClient, 10*1024*1024)
        if _, err := imageDownloader.Download(ctx, metadata.CoverURL, coverPath, nil); err != nil {
            logWarning("Warning: Cover image download failed: %v. Audio download completed successfully.", err)
        } else {
            logSuccess("Cover image saved to: %s", coverPath)
        }
    }

    // 5. Save show notes (optional, with graceful degradation)
    if metadata.ShowNotes != "" {
        showNotesSaver := downloader.NewPlainTextShowNotesSaver()
        if _, err := showNotesSaver.Save(metadata.ShowNotes, showNotesPath); err != nil {
            logWarning("Warning: Show notes extraction failed: %v. Audio download completed successfully.", err)
        } else {
            logSuccess("Show notes saved to: %s", showNotesPath)
        }
    }

    return nil
}
```

**Contract Guarantees**:
- Audio download is attempted and required for success
- Cover and show notes failures don't prevent overall success
- Returns error only if audio download fails
- Logs detailed warnings for non-critical failures
- Respects context cancellation at all stages

---

## Testing Contracts

### Unit Test Interface

```go
package downloader_test

// MockHTTPClient is a mock Doer for testing.
type MockHTTPClient struct {
    Documents map[string]*goquery.Document
    Error     error
}

func (m *MockHTTPClient) Get(url string) (*goquery.Document, error) {
    if m.Error != nil {
        return nil, m.Error
    }
    return m.Documents[url], nil
}

// Test: ExtractURL with all metadata present
func TestExtractURL_AllMetadataPresent(t *testing.T) {
    html := `
        <html>
            <title>Test Episode</title>
            <meta property="og:audio" content="https://example.com/audio.m4a" />
            <meta property="og:image" content="https://example.com/cover.jpg" />
            <section aria-label="节目show notes">
                <p>This is the show notes content.</p>
            </section>
        </html>
    `
    doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
    mock := &MockHTTPClient{Documents: map[string]*goquery.Document{"https://example.com": doc}}

    extractor := NewHTMLExtractor(mock)
    metadata, err := extractor.ExtractURL(context.Background(), "https://example.com")

    assert.NoError(t, err)
    assert.Equal(t, "https://example.com/audio.m4a", metadata.AudioURL)
    assert.Equal(t, "https://example.com/cover.jpg", metadata.CoverURL)
    assert.Contains(t, metadata.ShowNotes, "show notes content")
    assert.Equal(t, "Test Episode", metadata.Title)
}

// Test: ExtractURL with missing cover (graceful degradation)
func TestExtractURL_MissingCover(t *testing.T) {
    // ... implementation
}

// Test: ExtractURL with missing show notes (graceful degradation)
func TestExtractURL_MissingShowNotes(t *testing.T) {
    // ... implementation
}

// Test: FormatHTMLToText preserves links
func TestFormatHTMLToText_PreservesLinks(t *testing.T) {
    saver := NewPlainTextShowNotesSaver()
    html := `<a href="https://example.com">Link Text</a>`
    result := saver.FormatHTMLToText(html)
    assert.Contains(t, result, "Link Text (URL: https://example.com)")
}

// Test: ValidateImage detects format from magic bytes
func TestValidateImage_DetectsFormat(t *testing.T) {
    downloader := NewHTTPImageDownloader(&http.Client{}, 10*1024*1024)

    // Create test JPEG file
    jpegData := []byte{0xFF, 0xD8, 0xFF, 0xE0} // JPEG magic bytes + more
    tmpfile, _ := os.CreateTemp("", "test*.jpg")
    tmpfile.Write(jpegData)
    tmpfile.Close()

    format, err := downloader.ValidateImage(tmpfile.Name())
    assert.NoError(t, err)
    assert.Equal(t, "jpg", format)
}
```

---

## Version Contract

### Interface Versioning

All interfaces follow semantic versioning:
- **Major version bump**: Breaking change to interface signature
- **Minor version bump**: New interface added (backward compatible)
- **Patch version bump**: Implementation changes only (interface unchanged)

### Current Version

- `URLExtractor`: v2.0 (breaking change from v1)
- `ImageDownloader`: v1.0 (new interface)
- `ShowNotesSaver`: v1.0 (new interface)
- `FileDownloader`: v1.0 (unchanged)

### Migration from v1 to v2

```go
// Old code (v1)
audioURL, title, err := extractor.ExtractURL(ctx, pageURL)

// New code (v2)
metadata, err := extractor.ExtractURL(ctx, pageURL)
audioURL := metadata.AudioURL
title := metadata.Title
```

**Backward Compatibility Adapter**:

```go
// LegacyURLExtractor wraps v2 extractor to provide v1 interface
type LegacyURLExtractor struct {
    v2Extractor URLExtractor
}

func (l *LegacyURLExtractor) ExtractURL(ctx context.Context, pageURL string) (string, string, error) {
    metadata, err := l.v2Extractor.ExtractURL(ctx, pageURL)
    if err != nil {
        return "", "", err
    }
    return metadata.AudioURL, metadata.Title, nil
}
```
