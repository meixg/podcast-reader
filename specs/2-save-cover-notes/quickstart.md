# Quickstart Guide: Save Cover Images and Show Notes

**Feature**: 2-save-cover-notes
**Target Audience**: Developers implementing this feature
**Prerequisites**: Go 1.21+, familiarity with existing codebase

## Overview

This guide provides a step-by-step implementation path for adding cover image and show notes download functionality to the podcast downloader CLI tool.

**Estimated Implementation Time**: 4-6 hours
**Complexity**: Medium (extends existing code, adds new services)

## Phase 0: Setup & Validation (15 minutes)

### 1. Create Feature Branch

```bash
git checkout -b 2-save-cover-notes
```

### 2. Verify Prerequisites

```bash
# Check Go version
go version  # Should be 1.21+

# Run existing tests to ensure baseline
go test ./...

# Build existing CLI
go build -o podcast-downloader cmd/podcast-downloader/main.go
```

### 3. Review Existing Code

Read these files to understand the current architecture:
- `internal/downloader/url_extractor.go` - HTML extraction logic
- `internal/downloader/downloader.go` - File download logic
- `cmd/podcast-downloader/main.go` - CLI entry point

---

## Phase 1: Data Model & Interfaces (1 hour)

### Task 1.1: Create EpisodeMetadata Struct

**File**: `internal/downloader/metadata.go` (new file)

```go
package downloader

import "time"

// EpisodeMetadata contains all extractable metadata from a podcast episode page.
type EpisodeMetadata struct {
    AudioURL       string    // Direct URL to audio file (required)
    CoverURL       string    // URL to cover image (optional)
    ShowNotes      string    // Plain text show notes (optional)
    Title          string    // Episode title (required)
    EpisodeNumber  string    // Episode number if available
    PodcastName    string    // Podcast/series name
    PublicationDate time.Time // Publication date
}
```

**Test**: Create `internal/downloader/metadata_test.go` with basic struct tests.

---

### Task 1.2: Extend URLExtractor Interface

**File**: `internal/downloader/url_extractor.go` (modify existing)

**Before**:
```go
type URLExtractor interface {
    ExtractURL(ctx context.Context, pageURL string) (audioURL string, title string, err error)
}
```

**After**:
```go
type URLExtractor interface {
    ExtractURL(ctx context.Context, pageURL string) (*EpisodeMetadata, error)
}
```

**Impact**: This is a **breaking change**. Update all call sites (currently only `cmd/podcast-downloader/main.go`).

**Backward Compatibility** (optional):
```go
// LegacyURLExtractor wraps new interface for old code
type LegacyURLExtractor struct {
    newExtractor URLExtractor
}

func (l *LegacyURLExtractor) ExtractURL(ctx context.Context, pageURL string) (string, string, error) {
    metadata, err := l.newExtractor.ExtractURL(ctx, pageURL)
    if err != nil {
        return "", "", err
    }
    return metadata.AudioURL, metadata.Title, nil
}
```

---

### Task 1.3: Create New Interfaces

**File**: `internal/downloader/image_downloader.go` (new file)

```go
package downloader

import (
    "context"
    "io"
)

// ImageDownloader defines the interface for downloading cover images.
type ImageDownloader interface {
    Download(ctx context.Context, imageURL, filePath string, progress io.Writer) (int64, error)
    ValidateImage(filePath string) (format string, err error)
}
```

**File**: `internal/downloader/shownotes_saver.go` (new file)

```go
package downloader

// ShowNotesSaver defines the interface for saving show notes.
type ShowNotesSaver interface {
    Save(content, filePath string) (int64, error)
    FormatHTMLToText(htmlContent string) string
}
```

---

## Phase 2: Cover Image Extraction & Download (1.5 hours)

### Task 2.1: Implement Cover Image Extraction

**File**: `internal/downloader/url_extractor.go` (add new methods)

```go
// extractCoverURL extracts cover image URL from avatar-container.
// Cover images are guaranteed to exist on all podcast episodes.
func (e *HTMLExtractor) extractCoverURL(doc *goquery.Document) (string, error) {
    // avatar-container class selector (Xiaoyuzhou FM specific)
    // IMPORTANT: This container has TWO images:
    //   - First <img>: Episode cover (single episode artwork) ✅ This is what we want
    //   - Second <img>: Podcast account cover (channel/series artwork) ❌ Skip this
    // We use .First() to ensure we only get the episode cover
    if selection := doc.Find(".avatar-container img").First(); selection.Length() > 0 {
        if src, exists := selection.Attr("src"); exists && src != "" {
            return src, nil
        }
    }

    return "", ErrCoverNotFound
}
```

**Error**: Add to error definitions:
```go
var ErrCoverNotFound = errors.New("未找到封面图片")
```

---

### Task 2.2: Implement Image Downloader

**File**: `internal/downloader/image_downloader.go` (continue)

```go
package downloader

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "os"
)

// HTTPImageDownloader implements ImageDownloader with HTTP client.
type HTTPImageDownloader struct {
    client       *http.Client
    maxImageSize int64
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
    req, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
    if err != nil {
        return 0, fmt.Errorf("create request: %w", err)
    }

    resp, err := d.client.Do(req)
    if err != nil {
        return 0, fmt.Errorf("%w: %v", ErrNetworkTimeout, err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return 0, fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
    }

    // Check content length
    if resp.ContentLength > d.maxImageSize {
        return 0, fmt.Errorf("%w: size %d exceeds limit %d",
            ErrImageTooLarge, resp.ContentLength, d.maxImageSize)
    }

    out, err := os.Create(filePath)
    if err != nil {
        return 0, fmt.Errorf("%w: %v", ErrPermissionDenied, err)
    }
    defer out.Close()

    var writer io.Writer = out
    if progress != nil {
        writer = io.MultiWriter(out, progress)
    }

    bytesWritten, err := io.Copy(writer, resp.Body)
    if err != nil {
        return 0, fmt.Errorf("write file: %w", err)
    }

    return bytesWritten, nil
}

// ValidateImage validates the downloaded image.
func (d *HTTPImageDownloader) ValidateImage(filePath string) (string, error) {
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

**Errors**: Add to error definitions:
```go
var (
    ErrCoverNotFound   = errors.New("未找到封面图片")
    ErrInvalidImage    = errors.New("无效的图片文件")
    ErrImageTooLarge   = errors.New("图片文件过大")
)
```

---

### Task 2.3: Update HTMLExtractor.ExtractURL()

**File**: `internal/downloader/url_extractor.go` (modify existing method)

**Before**:
```go
func (e *HTMLExtractor) ExtractURL(ctx context.Context, pageURL string) (string, string, error) {
    doc, err := e.client.Get(pageURL)
    if err != nil {
        return "", "", fmt.Errorf("%w: %v", ErrPageNotFound, err)
    }

    title := e.extractTitle(doc)
    audioURL, err := e.extractAudioURL(doc)
    if err != nil {
        return "", title, err
    }

    return audioURL, title, nil
}
```

**After**:
```go
func (e *HTMLExtractor) ExtractURL(ctx context.Context, pageURL string) (*EpisodeMetadata, error) {
    doc, err := e.client.Get(pageURL)
    if err != nil {
        return nil, fmt.Errorf("%w: %v", ErrPageNotFound, err)
    }

    metadata := &EpisodeMetadata{}

    // Extract title (required)
    metadata.Title = e.extractTitle(doc)
    if metadata.Title == "" {
        metadata.Title = "episode" // Fallback
    }

    // Extract audio URL (required)
    metadata.AudioURL, err = e.extractAudioURL(doc)
    if err != nil {
        return nil, err
    }

    // Extract cover URL (optional)
    metadata.CoverURL, _ = e.extractCoverURL(doc)
    // No error if cover not found - graceful degradation

    // Extract show notes (optional) - add in next phase
    // metadata.ShowNotes, _ = e.extractShowNotes(doc)

    return metadata, nil
}
```

---

## Phase 3: Show Notes Extraction & Saving (1.5 hours)

### Task 3.1: Implement Show Notes Extraction

**File**: `internal/downloader/url_extractor.go` (add new method)

```go
// extractShowNotes extracts show notes using multi-fallback strategy.
func (e *HTMLExtractor) extractShowNotes(doc *goquery.Document) (string, error) {
    // Try 1: Exact aria-label match
    if selection := doc.Find("section[aria-label='节目show notes']").First(); selection.Length() > 0 {
        return selection.Html(), nil
    }

    // Try 2: Partial aria-label match (case-insensitive)
    var found bool
    var result string
    doc.Find("*[aria-label]").Each(func(i int, s *goquery.Selection) {
        if found {
            return
        }
        if ariaLabel, exists := s.Attr("aria-label"); exists {
            lowerLabel := strings.ToLower(ariaLabel)
            if strings.Contains(lowerLabel, "show notes") ||
               strings.Contains(lowerLabel, "节目说明") {
                result, _ = s.Html()
                found = true
            }
        }
    })
    if found {
        return result, nil
    }

    // Try 3: Semantic selectors
    selectors := []string{
        "article.episode-description",
        "section.description",
        "div.show-notes",
        "div.episode-content",
        "[itemprop='description']",
    }

    for _, selector := range selectors {
        if selection := doc.Find(selector).First(); selection.Length() > 0 {
            if html, err := selection.Html(); err == nil && html != "" {
                return html, nil
            }
        }
    }

    return "", ErrShowNotesNotFound
}
```

**Error**: Add to error definitions:
```go
var ErrShowNotesNotFound = errors.New("未找到节目说明")
```

---

### Task 3.2: Implement Show Notes Saver

**File**: `internal/downloader/shownotes_saver.go` (continue)

```go
package downloader

import (
    "bytes"
    "fmt"
    "os"
    "strings"
    "unicode/utf8"

    "github.com/PuerkitoBio/goquery"
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

    // Create file
    file, err := os.Create(filePath)
    if err != nil {
        return 0, fmt.Errorf("create file: %w", err)
    }
    defer file.Close()

    // Write UTF-8 BOM
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

// FormatHTMLToText converts HTML to readable plain text.
func (s *PlainTextShowNotesSaver) FormatHTMLToText(htmlContent string) string {
    doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(htmlContent)))
    if err != nil {
        return htmlContent // Fallback
    }

    // Convert links: <a href="url">text</a> → "text (URL: url)"
    doc.Find("a").Each(func(i int, s *goquery.Selection) {
        text := s.Text()
        if href, exists := s.Attr("href"); exists && href != "" {
            replacement := fmt.Sprintf("%s (URL: %s)", text, href)
            s.ReplaceWithHtml(replacement)
        }
    })

    // Convert lists to bullet points
    doc.Find("ul").Each(func(i int, s *goquery.Selection) {
        s.Find("li").Each(func(j int, li *goquery.Selection) {
            text := li.Text()
            li.ReplaceWithHtml(fmt.Sprintf("• %s", text))
        })
        s.AfterHtml("\n")
    })

    // Convert ordered lists
    doc.Find("ol").Each(func(i int, s *goquery.Selection) {
        s.Find("li").Each(func(j int, li *goquery.Selection) {
            text := li.Text()
            li.ReplaceWithHtml(fmt.Sprintf("%d. %s", j+1, text))
        })
        s.AfterHtml("\n")
    })

    // Get text and clean whitespace
    text := doc.Text()
    text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
    text = strings.TrimSpace(text)

    return text
}
```

**Error**: Add to error definitions:
```go
var ErrInvalidEncoding = errors.New("无效的字符编码")
```

---

### Task 3.3: Update ExtractURL to Include Show Notes

**File**: `internal/downloader/url_extractor.go` (update ExtractURL method)

Add after cover extraction:
```go
// Extract show notes (optional)
metadata.ShowNotes, _ = e.extractShowNotes(doc)
```

---

## Phase 4: CLI Integration (1 hour)

### Task 4.1: Update Main Download Workflow

**File**: `cmd/podcast-downloader/main.go`

**Before**:
```go
// Extract URL and title
audioURL, title, err := extractor.ExtractURL(ctx, episodeURL)
if err != nil {
    return err
}

// Generate filename
baseFilename := sanitizeFilename(title)
audioPath := filepath.Join(*outputDir, baseFilename+".m4a")

// Download audio
downloader := downloader.NewHTTPDownloader(client, *showProgress)
progress := newProgressBar("Downloading audio...")
_, err = downloader.Download(ctx, audioURL, audioPath, progress)
if err != nil {
    return err
}
```

**After**:
```go
// Extract metadata
metadata, err := extractor.ExtractURL(ctx, episodeURL)
if err != nil {
    return err
}

// Generate filenames
baseFilename := sanitizeFilename(metadata.Title)
audioPath := filepath.Join(*outputDir, baseFilename+".m4a")
coverPath := filepath.Join(*outputDir, baseFilename+".jpg")
showNotesPath := filepath.Join(*outputDir, baseFilename+".txt")

// Download audio (required)
audioDownloader := downloader.NewHTTPDownloader(client, *showProgress)
progress := newProgressBar("Downloading audio...")
_, err = audioDownloader.Download(ctx, metadata.AudioURL, audioPath, progress)
if err != nil {
    return err
}
logSuccess("Audio saved to: %s", audioPath)

// Download cover image (guaranteed to exist, but download may fail)
if metadata.CoverURL != "" {
    imageDownloader := downloader.NewHTTPImageDownloader(client, 10*1024*1024)
    if _, err := imageDownloader.Download(ctx, metadata.CoverURL, coverPath, nil); err != nil {
        logWarning("Warning: Cover image download failed: %v. Audio download completed successfully.", err)
    } else {
        logSuccess("Cover image saved to: %s", coverPath)
    }
}

// Save show notes (optional, with graceful degradation)
if metadata.ShowNotes != "" {
    showNotesSaver := downloader.NewPlainTextShowNotesSaver()
    if _, err := showNotesSaver.Save(metadata.ShowNotes, showNotesPath); err != nil {
        logWarning("Warning: Show notes extraction failed: %v. Audio download completed successfully.", err)
    } else {
        logSuccess("Show notes saved to: %s", showNotesPath)
    }
}
```

**Helper Functions** (add to main.go):
```go
func logWarning(format string, args ...interface{}) {
    fmt.Printf("\033[33m%s\033[0m\n", fmt.Sprintf(format, args...))
}

func logSuccess(format string, args ...interface{}) {
    fmt.Printf("\033[32m%s\033[0m\n", fmt.Sprintf(format, args...))
}
```

---

## Phase 5: Testing (1 hour)

### Task 5.1: Write Unit Tests

**File**: `internal/downloader/url_extractor_test.go`

```go
func TestExtractCoverURL(t *testing.T) {
    tests := []struct {
        name     string
        html     string
        expected string
        wantErr  bool
    }{
        {
            name: "og:image present",
            html: `<meta property="og:image" content="https://example.com/cover.jpg" />`,
            expected: "https://example.com/cover.jpg",
            wantErr: false,
        },
        {
            name:     "no cover image",
            html:     `<html><body>No image</body></html>`,
            expected: "",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            doc, _ := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
            extractor := NewHTMLExtractor(nil)
            result, err := extractor.extractCoverURL(doc)
            if (err != nil) != tt.wantErr {
                t.Errorf("extractCoverURL() error = %v, wantErr %v", err, tt.wantErr)
            }
            if result != tt.expected {
                t.Errorf("extractCoverURL() = %v, want %v", result, tt.expected)
            }
        })
    }
}

func TestExtractShowNotes(t *testing.T) {
    tests := []struct {
        name     string
        html     string
        contains string
        wantErr  bool
    }{
        {
            name: "aria-label exact match",
            html: `<section aria-label="节目show notes"><p>Show notes content</p></section>`,
            contains: "Show notes content",
            wantErr:  false,
        },
        {
            name:     "no show notes",
            html:     `<html><body>No notes</body></html>`,
            contains: "",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            doc, _ := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
            extractor := NewHTMLExtractor(nil)
            result, err := extractor.extractShowNotes(doc)
            if (err != nil) != tt.wantErr {
                t.Errorf("extractShowNotes() error = %v, wantErr %v", err, tt.wantErr)
            }
            if !tt.wantErr && !strings.Contains(result, tt.contains) {
                t.Errorf("extractShowNotes() = %v, want contains %v", result, tt.contains)
            }
        })
    }
}
```

---

### Task 5.2: Write Integration Tests

**File**: `cmd/podcast-downloader/main_test.go`

```go
func TestDownloadEpisodeWithCoverAndShowNotes(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    // Use a real test podcast URL (or mock server)
    testURL := "https://www.xiaoyuzhoufm.com/episode/test"

    // Create temp directory
    tmpDir, _ := os.MkdirTemp("", "podcast-test-*")
    defer os.RemoveAll(tmpDir)

    // Run download
    err := downloadEpisode(context.Background(), testURL, tmpDir)
    if err != nil {
        t.Fatalf("downloadEpisode() error = %v", err)
    }

    // Verify files exist
    files, _ := os.ReadDir(tmpDir)
    hasAudio := false
    hasCover := false
    hasShowNotes := false

    for _, f := range files {
        switch filepath.Ext(f.Name()) {
        case ".m4a":
            hasAudio = true
        case ".jpg", ".png", ".webp":
            hasCover = true
        case ".txt":
            hasShowNotes = true
        }
    }

    if !hasAudio {
        t.Error("Audio file was not created")
    }
    // Cover and show notes are optional, so don't fail if missing
}
```

---

### Task 5.3: Manual Testing

```bash
# Build the CLI
go build -o podcast-downloader cmd/podcast-downloader/main.go

# Test with a real podcast URL
./podcast-downloader "https://www.xiaoyuzhoufm.com/episode/..."

# Check output directory
ls -la downloads/

# Verify files
# - episode-title.m4a (audio)
# - episode-title.jpg (cover)
# - episode-title.txt (show notes)

# Verify show notes encoding
file downloads/episode-title.txt
# Should show: UTF-8 Unicode (with BOM) text

# Verify cover image
file downloads/episode-title.jpg
# Should show: JPEG image data
```

---

## Phase 6: Documentation & Polish (30 minutes)

### Task 6.1: Update README

**File**: `README.md` (in project root)

Add section:
```markdown
## Downloaded Files

When you download a podcast episode, the tool now automatically saves:

1. **Audio file** (`.m4a`) - The episode audio
2. **Cover image** (`.jpg`, `.png`, or `.webp`) - Episode artwork
3. **Show notes** (`.txt`) - Episode description with links and timestamps

All files use the same base filename and are saved in the same directory for easy organization.

### Example

```
./podcast-downloader "https://www.xiaoyuzhoufm.com/episode/12345"

# Creates:
# downloads/Episode-Title.m4a
# downloads/Episode-Title.jpg
# downloads/Episode-Title.txt
```
```

---

### Task 6.2: Add Error Handling Examples

Update help text:
```bash
./podcast-downloader --help
```

Add to output:
```
If cover image or show notes cannot be downloaded, the tool will display a warning
but continue with the audio download. The audio file is always saved successfully
if the URL is valid.
```

---

## Verification Checklist

Before creating a pull request, verify:

- [ ] All unit tests pass: `go test ./...`
- [ ] Integration tests pass: `go test -v ./cmd/podcast-downloader/`
- [ ] CLI builds successfully: `go build -o podcast-downloader cmd/podcast-downloader/main.go`
- [ ] Manual test with real URL succeeds
- [ ] Cover image downloads and validates correctly
- [ ] Show notes save with UTF-8-BOM encoding
- [ ] Warning messages display correctly for failures
- [ ] Audio download still works (backward compatibility)
- [ ] Code follows gofmt: `gofmt -w .`
- [ ] No lint errors: `golint ./...`
- [ ] No vet errors: `go vet ./...`

---

## Troubleshooting

### Issue: Cover image downloads but fails validation

**Solution**: Check magic byte detection in `detectFormat()`. Some sites may redirect to HTML error pages. Add response content-type validation:

```go
contentType := resp.Header.Get("Content-Type")
if !strings.HasPrefix(contentType, "image/") {
    return "", fmt.Errorf("not an image: %s", contentType)
}
```

---

### Issue: Show notes contain HTML tags

**Solution**: The `FormatHTMLToText()` method may not handle all HTML structures. Add additional selectors or use a library like `github.com/k3a/html2text`.

---

### Issue: Chinese characters display incorrectly in Windows Notepad

**Solution**: Verify UTF-8 BOM is being written. Check:

```go
// Write UTF-8 BOM
bom := []byte{0xEF, 0xBB, 0xBF}
if _, err := file.Write(bom); err != nil {
    return 0, err
}
```

---

## Next Steps

After implementation:

1. **Create Pull Request**: Target branch `main`
2. **Update CLAUDE.md**: Run `/speckit.implement` to update project documentation
3. **Test Coverage**: Aim for >80% coverage on new code
4. **Performance Testing**: Test with multiple episodes in sequence

---

## References

- **Research**: `specs/2-save-cover-notes/research.md`
- **Data Model**: `specs/2-save-cover-notes/data-model.md`
- **Interfaces**: `specs/2-save-cover-notes/contracts/interfaces.md`
- **Spec**: `specs/2-save-cover-notes/spec.md`
- **Existing Code**: `internal/downloader/`
