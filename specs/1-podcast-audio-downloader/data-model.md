# Data Model: Podcast Audio Downloader

**Feature**: Podcast Audio Downloader (CLI tool)
**Date**: 2026-02-08
**Purpose**: Define data structures for podcast episode metadata and download sessions

---

## Overview

The podcast audio downloader CLI operates with two primary data entities:
1. **Episode**: Represents podcast episode metadata extracted from the webpage
2. **DownloadSession**: Represents a single download operation with state tracking

---

## Entity: Episode

Represents a single podcast episode from Xiaoyuzhou FM with metadata needed for downloading and file naming.

### Fields

| Field | Type | Description | Validation | Source |
|-------|------|-------------|------------|--------|
| `ID` | `string` | Unique episode identifier from URL or page | Required, non-empty | Extracted from URL path or HTML |
| `Title` | `string` | Episode title for filename generation | Required, sanitized if invalid | Extracted from HTML `<title>` or meta tags |
| `AudioURL` | `string` | Direct URL to .m4a audio file | Required, valid URL format | Extracted from HTML `<audio src>` or `<source src>` |
| `PageURL` | `string` | Original episode page URL | Required, valid URL format | User input |
| `FileSize` | `int64` | Audio file size in bytes | Optional, -1 if unknown | HTTP `Content-Length` header |
| `Duration` | `int` | Episode duration in seconds | Optional, 0 if unknown | Extracted from metadata or HTML |

### Relationships

- **DownloadSession**: One episode can have multiple download sessions (retry attempts)

### Validation Rules

1. **ID Validation**:
   - Must not be empty
   - Must match pattern: numeric ID from URL path (e.g., "69392768281939cce65925d3")
   - Fallback: Use URL hash if ID not found

2. **Title Sanitization**:
   - Remove invalid filesystem characters: `< > : " / \ | ? *`
   - Limit to 200 characters
   - Trim whitespace
   - Fallback to episode ID if title is empty after sanitization

3. **AudioURL Validation**:
   - Must be a valid absolute URL
   - Must end with `.m4a` extension
   - Must be HTTP/HTTPS protocol
   - Reject if empty or malformed

4. **PageURL Validation**:
   - Must match Xiaoyuzhou FM domain pattern: `*.xiaoyuzhoufm.com/episode/*`
   - Must be HTTP/HTTPS protocol

### Example

```go
type Episode struct {
    ID        string
    Title     string
    AudioURL  string
    PageURL   string
    FileSize  int64
    Duration  int
}

// Example instance
episode := Episode{
    ID:       "69392768281939cce65925d3",
    Title:    "技术电台第42期: Go语言实战",
    AudioURL: "https://audio.xiaoyuzhoufm.com/episodes/69392768281939cce65925d3.m4a",
    PageURL:  "https://www.xiaoyuzhoufm.com/episode/69392768281939cce65925d3",
    FileSize: 52428800,  // 50 MB
    Duration: 1800,      // 30 minutes
}
```

---

## Entity: DownloadSession

Represents a single download operation with progress tracking and state management.

### Fields

| Field | Type | Description | Validation | Default |
|-------|------|-------------|------------|---------|
| `EpisodeID` | `string` | Reference to episode being downloaded | Required, matches Episode.ID | - |
| `FilePath` | `string` | Local file path for downloaded audio | Required, valid filesystem path | Generated from title/ID |
| `Status` | `Status` | Current download state | Required, one of: Pending, InProgress, Completed, Failed | Pending |
| `Progress` | `float64` | Download progress (0.0 to 1.0) | Required, 0 ≤ Progress ≤ 1 | 0.0 |
| `BytesDownloaded` | `int64` | Number of bytes downloaded | Required, ≥ 0 | 0 |
| `TotalBytes` | `int64` | Total file size (-1 if unknown) | Required, -1 or positive | -1 |
| `StartTime` | `time.Time` | Download start timestamp | Required | Time of creation |
| `EndTime` | `time.Time` | Download completion timestamp | Optional | Zero value if not completed |
| `Error` | `error` | Error message if download failed | Optional, nil if no error | nil |
| `RetryCount` | `int` | Number of retry attempts | Required, ≥ 0 | 0 |
| `DownloadSpeed` | `float64` | Current download speed in bytes/sec | Optional, ≥ 0 | 0 |

### Status Enum

```go
type Status int

const (
    StatusPending    Status = iota
    StatusInProgress
    StatusCompleted
    StatusFailed
)
```

**Status Transitions**:
```
Pending → InProgress → Completed
                    ↘ Failed
Pending → Failed (if validation fails)
```

### State Machine

```
┌──────────┐
│ Pending  │ ──[Start Download]──┐
└──────────┘                     │
                                 ▼
                          ┌──────────┐
                          │InProgress│ ──[Success]──┐
┌──────────┐                └──────────┘              │
│  Failed  │ ←[Retry Limit]───[Network Error]         ▼
└──────────┘                                       ┌───────────┐
     │                                              │ Completed │
     └──[User Retry]──[Validation Error]────────────┴───────────┘
```

### Relationships

- **Episode**: Many-to-one (multiple sessions can reference the same episode)
- **File**: One-to-one (each session creates or overwrites one file)

### Validation Rules

1. **FilePath Validation**:
   - Must be absolute or relative path with .m4a extension
   - Parent directory must exist or be creatable
   - Must not conflict with existing files unless overwrite mode enabled

2. **Progress Validation**:
   - Must be between 0.0 and 1.0
   - Must be 1.0 when Status is Completed
   - Must be < 1.0 when Status is InProgress

3. **BytesDownloaded Validation**:
   - Must not exceed TotalBytes (if TotalBytes > 0)
   - Must equal TotalBytes when Status is Completed (if TotalBytes > 0)

4. **RetryCount Validation**:
   - Must not exceed maximum retry limit (3)
   - Increments on each retry attempt

### Example

```go
type DownloadSession struct {
    EpisodeID       string
    FilePath        string
    Status          Status
    Progress        float64
    BytesDownloaded int64
    TotalBytes      int64
    StartTime       time.Time
    EndTime         time.Time
    Error           error
    RetryCount      int
    DownloadSpeed   float64
}

// Example instance (in progress)
session := DownloadSession{
    EpisodeID:       "69392768281939cce65925d3",
    FilePath:        "./downloads/技术电台第42期_Go语言实战.m4a",
    Status:          StatusInProgress,
    Progress:        0.45,  // 45% complete
    BytesDownloaded: 23592960,  // 22.5 MB
    TotalBytes:      52428800,  // 50 MB
    StartTime:       time.Now().Add(-2 * time.Minute),
    RetryCount:      0,
    DownloadSpeed:   196608,  // 192 KB/s
}
```

---

## File Naming Convention

### Pattern

```
{sanitized_title}_{episode_id}.m4a
```

### Examples

| Original Title | Sanitized Title | Episode ID | Final Filename |
|----------------|-----------------|------------|----------------|
| "技术电台第42期: Go语言实战" | "技术电台第42期_ Go语言实战" | 69392768281939cce65925d3 | "技术电台第42期_ Go语言实战_69392768281939cce65925d3.m4a" |
| "Episode #42: What's New?" | "Episode_42_ What_s_New_" | abc123 | "Episode_42_ What_s_New__abc123.m4a" |
| "Test<>:File|?Name*" | "Test___File__Name_" | xyz789 | "Test___File__Name__xyz789.m4a" |

### Fallback

If title is empty or becomes empty after sanitization:
```
{episode_id}.m4a
```

Example: `"69392768281939cce65925d3.m4a"`

---

## Directory Structure

### Default Layout

```
downloads/
├── 技术电台第42期_ Go语言实战_69392768281939cce65925d3.m4a
├── Episode_42_ What_s_New__abc123.m4a
└── 69392768281939cce65925d3.m4a
```

### Configuration

Users can override default directory via CLI flag:
```bash
podcast-downloader --output ~/music/podcasts <url>
```

---

## Error Handling

### Episode Extraction Errors

| Error | Description | User Action |
|-------|-------------|-------------|
| `ErrInvalidURL` | URL format is invalid | Check URL syntax |
| `ErrPageNotFound` | Episode page returns 404 | Verify episode exists |
| `ErrAudioNotFound` | No audio link found on page | Contact website support |
| `ErrAccessDenied` | Page requires authentication | Cannot handle (assumption #2) |

### Download Errors

| Error | Description | Retryable |
|-------|-------------|-----------|
| `ErrNetworkTimeout` | Request timeout | Yes (up to 3 retries) |
| `ErrConnectionRefused` | Cannot connect to server | Yes (up to 3 retries) |
| `ErrDiskFull` | Insufficient disk space | No (user action required) |
| `ErrPermissionDenied` | Cannot write to directory | No (user action required) |
| `ErrInvalidAudio` | Downloaded file is not audio | No (server issue) |

---

## Persistence Model

**Note**: This CLI tool does not persist Episode or DownloadSession data to disk or database. All data is ephemeral and stored in memory during execution.

**Rationale**:
- Simplifies the tool for MVP
- Users can see download status via progress bar
- File system provides persistence (downloaded files)
- Future web service can add database persistence if needed

---

## Next Steps

With data models defined:
1. ✅ Episode entity with metadata extraction
2. ✅ DownloadSession entity with progress tracking
3. ✅ File naming and validation rules
4. ✅ Error handling strategy
5. ✅ Status state machine

Proceed to **contracts/** and **quickstart.md** generation.
