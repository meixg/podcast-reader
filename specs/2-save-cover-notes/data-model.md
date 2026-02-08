# Data Model: Save Cover Images and Show Notes

**Feature**: 2-save-cover-notes
**Date**: 2026-02-08
**Status**: Complete

## Overview

This document defines the data structures and entities for the enhanced podcast downloader that supports cover images and show notes.

## Core Entities

### EpisodeMetadata

Represents all extractable metadata from a podcast episode page.

**Purpose**: Encapsulates all data extracted from the podcast page in a single structure, avoiding multiple HTTP requests.

**Fields**:
- `AudioURL` (string, required): Direct URL to the audio file (.m4a)
- `CoverURL` (string, optional): URL to the cover image (JPEG, PNG, or WEBP)
- `ShowNotes` (string, optional): Plain text content of show notes
- `Title` (string, required): Episode title for filename generation
- `EpisodeNumber` (string, optional): Episode number if available
- `PodcastName` (string, optional): Podcast/series name
- `PublicationDate` (time.Time, optional): Episode publication date

**Validation Rules**:
- `AudioURL` must be a valid HTTP/HTTPS URL
- `Title` must not be empty after sanitization
- `CoverURL` must be a valid HTTP/HTTPS URL if present
- `ShowNotes` length must be <= 1,000,000 characters (~1MB of text)

**State Transitions**: N/A (immutable after extraction)

**Relationships**:
- One `EpisodeMetadata` → One `DownloadSession`
- One `EpisodeMetadata` → One `AudioFile` (downloaded)
- One `EpisodeMetadata` → Zero or One `CoverImage` (downloaded)
- One `EpisodeMetadata` → Zero or One `ShowNotesFile` (created)

---

### DownloadSession

Tracks the download operation state and results.

**Purpose**: Manages the complete download workflow including audio, cover, and show notes with error tracking.

**Fields**:
- `EpisodeURL` (string, required): Original podcast episode page URL
- `Metadata` (EpisodeMetadata, required): Extracted metadata
- `AudioDownloaded` (bool, required): Whether audio file was successfully downloaded
- `CoverDownloaded` (bool, required): Whether cover image was successfully downloaded
- `ShowNotesSaved` (bool, required): Whether show notes file was successfully saved
- `AudioPath` (string, optional): Absolute path to downloaded audio file
- `CoverPath` (string, optional): Absolute path to downloaded cover image
- `ShowNotesPath` (string, optional): Absolute path to saved show notes file
- `CoverError` (error, optional): Error encountered during cover download (if any)
- `ShowNotesError` (error, optional): Error encountered during show notes extraction/saving (if any)
- `StartTime` (time.Time, required): Download start time
- `EndTime` (time.Time, optional): Download end time
- `TotalBytes` (int64, required): Total bytes downloaded (audio + cover)

**Validation Rules**:
- `EpisodeURL` must be a valid HTTP/HTTPS URL
- At least `AudioDownloaded` must be true (cover/show notes are optional)
- If `CoverDownloaded` is true, `CoverPath` must be set
- If `ShowNotesSaved` is true, `ShowNotesPath` must be set
- `EndTime` must be after `StartTime` if set

**State Transitions**:
```
[Created] → [Extracting Metadata] → [Downloading Audio] → [Downloading Cover] → [Saving Show Notes] → [Completed]
                              ↘ [Failed]                ↘ [Failed with Warnings]
```

**Relationships**:
- One `DownloadSession` → One `EpisodeMetadata`
- One `DownloadSession` → One `AudioFile`
- One `DownloadSession` → Zero or One `CoverImage`
- One `DownloadSession` → Zero or One `ShowNotesFile`

---

### AudioFile

Represents the downloaded audio file.

**Purpose**: Tracks audio file download state and validation.

**Fields**:
- `URL` (string, required): Source URL for the audio file
- `Path` (string, required): Local filesystem path
- `Size` (int64, required): File size in bytes
- `Format` (string, required): Audio format (e.g., "m4a", "mp3")
- `Validated` (bool, required): Whether file validation passed
- `ValidationMethod` (string, required): Method used for validation (e.g., "magic-bytes", "extension")

**Validation Rules**:
- `Size` must be > 0
- `Format` must be a supported audio format (m4a, mp3, wav, etc.)
- `Validated` must be true for successful downloads

**State Transitions**:
```
[Pending] → [Downloading] → [Validating] → [Complete]
              ↘ [Failed]     ↘ [Invalid]
```

**Relationships**:
- Belongs to one `DownloadSession`

---

### CoverImage

Represents the downloaded cover image file.

**Purpose**: Tracks cover image download state and validation.

**Fields**:
- `URL` (string, required): Source URL for the cover image
- `Path` (string, required): Local filesystem path
- `Size` (int64, required): File size in bytes
- `Format` (string, required): Image format (e.g., "jpg", "png", "webp")
- `Width` (int, optional): Image width in pixels (if available)
- `Height` (int, optional): Image height in pixels (if available)
- `Validated` (bool, required): Whether file validation passed
- `ValidationMethod` (string, required): Method used for validation (e.g., "magic-bytes")

**Validation Rules**:
- `Size` must be > 0 and <= 10MB (configurable limit)
- `Format` must be a supported image format (jpg, png, webp, gif)
- `Validated` must be true for successful downloads
- Magic bytes must match the declared format

**State Transitions**:
```
[Pending] → [Downloading] → [Validating] → [Complete]
              ↘ [Failed]     ↘ [Invalid]
```

**Relationships**:
- Belongs to one `DownloadSession`

---

### ShowNotesFile

Represents the saved show notes text file.

**Purpose**: Tracks show notes extraction and file creation.

**Fields**:
- `Content` (string, required): Plain text content of show notes
- `Path` (string, required): Local filesystem path
- `Size` (int64, required): File size in bytes
- `Encoding` (string, required): Text encoding (always "UTF-8-BOM")
- `CharacterCount` (int, required): Number of characters in content
- `LineCount` (int, required): Number of lines in content
- `LinksPreserved` (int, required): Number of hyperlinks preserved in text
- `SourceElement` (string, optional): HTML element selector used for extraction

**Validation Rules**:
- `Content` must not be empty
- `Encoding` must be "UTF-8-BOM"
- `Size` must be > 0 and <= 1MB
- `Content` must be valid UTF-8

**State Transitions**:
```
[Pending] → [Extracting] → [Formatting] → [Writing] → [Complete]
              ↘ [Failed]     ↘ [Failed]     ↘ [Failed]
```

**Relationships**:
- Belongs to one `DownloadSession`

## Extended Models

### FilenameGenerator (Utility)

Generates consistent filenames across all asset types.

**Methods**:
- `GenerateBaseFilename(title string) string`: Sanitizes title and returns base filename
- `GetAudioFilename(baseFilename string, format string) string`: Returns audio filename
- `GetCoverFilename(baseFilename string, format string) string`: Returns cover filename
- `GetShowNotesFilename(baseFilename string) string`: Returns show notes filename

**Sanitization Rules**:
- Remove invalid characters: `< > : " / \ | ? *`
- Replace with underscore or remove (depending on platform)
- Limit length to 200 characters (leaving room for extensions)
- Preserve Chinese characters and Unicode letters
- Trim leading/trailing whitespace and dots

**Examples**:
```
Title: "Episode 42: The Future of AI"
Base filename: "Episode 42 - The Future of AI"
Audio: "Episode 42 - The Future of AI.m4a"
Cover: "Episode 42 - The Future of AI.jpg"
Show notes: "Episode 42 - The Future of AI.txt"
```

```
Title: "为什么播客很重要？"
Base filename: "为什么播客很重要？"
Audio: "为什么播客很重要？.m4a"
Cover: "为什么播客很重要？.jpg"
Show notes: "为什么播客很重要？.txt"
```

---

### HTMLExtractor (Enhanced)

Extends the existing HTML extractor to support cover and show notes extraction.

**Interface Extension**:
```go
// OLD interface (backward incompatible change)
type URLExtractor interface {
    ExtractURL(ctx context.Context, pageURL string) (audioURL string, title string, err error)
}

// NEW interface
type URLExtractor interface {
    ExtractURL(ctx context.Context, pageURL string) (*EpisodeMetadata, error)
}
```

**New Methods**:
- `extractCoverURL(doc *goquery.Document) (string, error)`: Extracts cover image URL
- `extractShowNotes(doc *goquery.Document) (string, error)`: Extracts show notes content
- `extractEpisodeNumber(doc *goquery.Document) string`: Extracts episode number
- `extractPodcastName(doc *goquery.Document) string`: Extracts podcast/series name

**Fallback Strategy**:
```go
func (e *HTMLExtractor) extractShowNotes(doc *goquery.Document) (string, error) {
    // Try 1: Exact aria-label match
    if selection := doc.Find("section[aria-label='节目show notes']").First(); selection.Length() > 0 {
        return selection.Text(), nil
    }

    // Try 2: Partial aria-label match (case-insensitive)
    doc.Find("*[aria-label]").Each(func(i int, s *goquery.Selection) {
        if ariaLabel, exists := s.Attr("aria-label"); exists {
            if containsIgnoreCase(ariaLabel, "show notes") {
                // Return first match
            }
        }
    })

    // Try 3: Semantic selectors
    selectors := []string{
        "article.episode-description",
        "section.description",
        "div.show-notes",
        "div.episode-content",
    }
    for _, selector := range selectors {
        // ...
    }

    // Fallback: Return error with specific reason
    return "", ErrShowNotesNotFound
}
```

---

### ImageDownloader (New)

Downloads cover images with progress tracking and validation.

**Interface**:
```go
type ImageDownloader interface {
    Download(ctx context.Context, imageURL, filePath string, progress io.Writer) (int64, error)
    ValidateImage(filePath string) (format string, err error)
    DetectFormat(data []byte) string
}
```

**Methods**:
- `Download()`: Downloads image with retry logic and progress display
- `ValidateImage()`: Validates downloaded file using magic byte detection
- `DetectFormat()`: Detects image format from binary data

**Magic Byte Detection**:
```go
func (d *HTTPImageDownloader) DetectFormat(data []byte) string {
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
    if data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 &&
       len(data) > 11 && data[8] == 0x57 && data[9] == 0x45 && data[10] == 0x42 && data[11] == 0x50 {
        return "webp"
    }

    return ""
}
```

---

### ShowNotesSaver (New)

Saves show notes content to text files with proper encoding.

**Interface**:
```go
type ShowNotesSaver interface {
    Save(content, filePath string) (int64, error)
    FormatHTMLToText(htmlContent string) string
    WriteWithBOM(file *os.File, content string) error
}
```

**Methods**:
- `Save()`: Saves formatted show notes to file with UTF-8-BOM encoding
- `FormatHTMLToText()`: Converts HTML to plain text while preserving structure
- `WriteWithBOM()`: Writes content with UTF-8 BOM prefix

**HTML to Text Conversion**:
```go
func (s *PlainTextShowNotesSaver) FormatHTMLToText(htmlContent string) string {
    doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
    if err != nil {
        return htmlContent // Fallback to original
    }

    var builder strings.Builder

    // Convert links: <a href="url">text</a> → "text (URL: url)"
    doc.Find("a").Each(func(i int, s *goquery.Selection) {
        text := s.Text()
        if href, exists := s.Attr("href"); exists {
            s.ReplaceWithHtml(fmt.Sprintf("%s (URL: %s)", text, href))
        }
    })

    // Convert lists to bullet/numbered lists
    // Convert headers to uppercase with underline
    // Preserve blockquotes with "> " prefix
    // Strip excessive whitespace

    return builder.String()
}
```

## Entity Relationship Diagram

```
┌─────────────────────┐
│  EpisodeMetadata    │
│  - AudioURL         │
│  - CoverURL         │
│  - ShowNotes        │
│  - Title            │
└──────────┬──────────┘
           │
           │ 1
           │
           ▼
┌─────────────────────┐       ┌─────────────────┐
│  DownloadSession    │──────→│   AudioFile     │
│  - EpisodeURL       │   1   │ - URL           │
│  - AudioPath        │       │ - Path          │
│  - CoverPath        │       │ - Size          │
│  - ShowNotesPath    │       └─────────────────┘
│  - State            │
└──────────┬──────────┘
           │
           │ 1
           │
           ├─────────────────────┐
           │                     │
           ▼                     ▼
    ┌──────────┐          ┌──────────────┐
    │CoverImage│ 0..1     │ShowNotesFile │ 0..1
    │- URL     │          │ - Content    │
    │- Path    │          │ - Path       │
    │- Size    │          │ - Encoding   │
    └──────────┘          └──────────────┘
```

## Data Flow

```
1. User runs: ./podcast-downloader <episode-url>
   ↓
2. DownloadSession created with EpisodeURL
   ↓
3. HTMLExtractor.ExtractURL() → EpisodeMetadata
   ├── extractAudioURL()
   ├── extractCoverURL() (new)
   ├── extractShowNotes() (new)
   └── extractTitle()
   ↓
4. Sequential Downloads:
   ├── AudioDownloader.Download() → AudioFile
   ├── ImageDownloader.Download() → CoverImage (with error handling)
   └── ShowNotesSaver.Save() → ShowNotesFile (with error handling)
   ↓
5. DownloadSession marked complete
   ↓
6. Display results with warnings for any failures
```

## Error Handling Model

### Error Types

```go
var (
    // Existing errors
    ErrInvalidURL     = errors.New("无效的URL")
    ErrPageNotFound   = errors.New("页面不存在")
    ErrAudioNotFound  = errors.New("未找到音频文件")

    // New errors for cover images
    ErrCoverNotFound        = errors.New("未找到封面图片")
    ErrCoverDownloadFailed  = errors.New("封面图片下载失败")
    ErrInvalidImage         = errors.New("无效的图片文件")
    ErrImageTooLarge        = errors.New("图片文件过大")

    // New errors for show notes
    ErrShowNotesNotFound    = errors.New("未找到节目说明")
    ErrShowNotesExtractionFailed = errors.New("节目说明提取失败")
    ErrInvalidEncoding      = errors.New("无效的字符编码")
)
```

### Error Wrapping

All errors are wrapped with context:
```go
return "", fmt.Errorf("提取封面图片: %w", ErrCoverNotFound)
return "", fmt.Errorf("下载封面图片 %s: %w", url, ErrCoverDownloadFailed)
```

### Warning Display

```go
if session.CoverError != nil {
    logWarning("Warning: Cover image download failed: %s. Audio download completed successfully.", session.CoverError)
}
if session.ShowNotesError != nil {
    logWarning("Warning: Show notes extraction failed: %s. Audio download completed successfully.", session.ShowNotesError)
}
```

## Storage Model

### File Organization

```
downloads/                           # Default output directory
├── Episode 01 - Introduction/       # Optional: Subdirectory by episode (future)
│   ├── episode-01.m4a
│   ├── episode-01.jpg
│   └── episode-01.txt
├── Episode 02 - Deep Dive/
│   ├── episode-02.m4a
│   ├── episode-02.png
│   └── episode-02.txt
└── Episode 03 - Interview/
    ├── episode-03.m4a
    ├── episode-03.webp
    └── episode-03.txt
```

**Current behavior** (v1.x):
```
downloads/
├── episode-01.m4a
├── episode-01.jpg
├── episode-01.txt
├── episode-02.m4a
├── episode-02.jpg
└── episode-02.txt
```

### File Permissions

- Audio files: `0644` (rw-r--r--)
- Cover images: `0644` (rw-r--r--)
- Show notes: `0644` (rw-r--r--)

### File Naming Edge Cases

| Scenario | Title | Sanitized Base | Audio | Cover | Notes |
|----------|-------|----------------|-------|-------|-------|
| Chinese characters | "为什么播客很重要？" | "为什么播客很重要？" | .m4a | .jpg | .txt |
| Special chars | "Episode <1>:" | "Episode 1_" | .m4a | .jpg | .txt |
| Very long title | "A".Repeat(300) | "A".Repeat(200) | .m4a | .jpg | .txt |
| Emoji | "🎵 Episode" | "🎵 Episode" | .m4a | .jpg | .txt |
| Empty after sanitize | "<>!@#" | "episode" | .m4a | .jpg | .txt |
| Multiple dots | "...test..." | "test" | .m4a | .jpg | .txt |

## Migration Path

### Breaking Changes

1. **URLExtractor interface change**:
```go
// Old
ExtractURL(ctx context.Context, pageURL string) (audioURL string, title string, error)

// New
ExtractURL(ctx context.Context, pageURL string) (*EpisodeMetadata, error)
```

**Migration strategy**:
- Create `ExtractURLV2()` method with new signature
- Deprecate old `ExtractURL()` method
- Provide adapter function for backward compatibility

### Non-Breaking Additions

- New services: `ImageDownloader`, `ShowNotesSaver`
- New error types
- Enhanced filename generation (backward compatible)

---

## Glossary

- **Cover Image**: Album art or thumbnail image associated with a podcast episode
- **Show Notes**: Episode description, guest information, topics, links, and timestamps
- **EpisodeMetadata**: All extractable data from a podcast page in one structure
- **BOM**: Byte Order Mark, a Unicode character (U+FEFF) used to signal encoding
- **Magic Bytes**: File signature at the beginning of a file indicating its format
- **Graceful Degradation**: System continues functioning when non-critical components fail
