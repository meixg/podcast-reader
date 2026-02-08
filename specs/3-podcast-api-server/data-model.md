# Data Model: Podcast API Server

**Feature**: 3-podcast-api-server
**Date**: 2026-02-08
**Phase**: Phase 1 - Design & Contracts

## Overview

This document defines the data entities and their relationships for the Podcast API Server. The system uses an in-memory task store for active downloads and filesystem-based storage for completed podcasts.

## Entity Definitions

### 1. Download Task

**Purpose**: Tracks the lifecycle of a podcast download request from submission to completion.

**Fields**:

| Field | Type | Description | Validation | Notes |
|-------|------|-------------|------------|-------|
| `id` | string (UUID) | Unique task identifier | Valid UUID v4 format | Auto-generated on submission |
| `url` | string | Source Xiaoyuzhou FM URL | Valid URL, xiaoyuzhou.fm domain | Used for duplicate detection |
| `status` | enum | Current task state | One of: pending, in_progress, completed, failed | State transitions are unidirectional |
| `created_at` | datetime | Task submission timestamp | ISO 8601 format | Auto-generated |
| `started_at` | datetime | Download start timestamp | ISO 8601 format | Null until status = in_progress |
| `completed_at` | datetime | Task completion timestamp | ISO 8601 format | Null until status = completed or failed |
| `error` | string | Error message if failed | Non-empty when status = failed | Null otherwise |
| `progress` | integer | Download progress percentage | 0-100 | Optional, if available from downloader |
| `podcast` | Podcast Episode (embedded) | Downloaded podcast metadata | Null until status = completed | Populated after successful download |

**State Transitions**:

```
pending → in_progress → completed
                    ↘ failed
```

**Constraints**:
- `id` is immutable and unique
- `url` is immutable
- `status` transitions follow the diagram above (no rollback)
- `completed_at` must be >= `started_at` >= `created_at`

**Indexes**:
- Primary: `id`
- Secondary: `url` (for duplicate detection)

### 2. Podcast Episode

**Purpose**: Represents a successfully downloaded podcast with all its assets and metadata.

**Fields**:

| Field | Type | Description | Validation | Notes |
|-------|------|-------------|------------|-------|
| `title` | string | Podcast episode title | Non-empty if status = completed | Sanitized from source |
| `source_url` | string | Original Xiaoyuzhou FM URL | Valid URL, xiaoyuzhou.fm domain | Immutable |
| `audio_path` | string | Absolute path to audio file | Valid file path, .m4a extension | Relative to downloads directory optional |
| `cover_path` | string | Absolute path to cover image | Valid file path or empty | Optional (some episodes lack covers) |
| `cover_format` | enum | Cover image format | One of: jpg, png, webp, empty | Derived from file extension |
| `shownotes_path` | string | Absolute path to show notes file | Valid file path or empty | Optional (some episodes lack notes) |
| `downloaded_at` | datetime | Download completion timestamp | ISO 8601 format | When files were written to disk |
| `file_size_mb` | number | Audio file size in megabytes | Non-negative | Optional, for display purposes |

**Validation Rules**:
- At least `audio_path` must be non-empty (required asset)
- `cover_path` and `shownotes_path` are optional (graceful degradation)
- All paths must exist on filesystem when returned
- Paths are absolute for API responses

**Relationships**:
- Embedded in completed `DownloadTask`
- One-to-one with files in `downloads/{title}/` directory

### 3. Podcast Catalog Entry

**Purpose**: In-memory index of downloaded podcasts for duplicate detection and list endpoint.

**Fields**:

| Field | Type | Description | Validation | Notes |
|-------|------|-------------|------------|-------|
| `url` | string | Source URL (key) | Valid URL, unique | Primary key |
| `title` | string | Podcast title | Non-empty | From directory name |
| `directory` | string | Relative path to podcast directory | Valid subdirectory of downloads/ | e.g., "Podcast Title" |
| `audio_file` | string | Audio filename | "podcast.m4a" or similar | From filesystem scan |
| `has_cover` | boolean | Cover image exists | true/false | Derived from filesystem |
| `has_shownotes` | boolean | Show notes file exists | true/false | Derived from filesystem |
| `downloaded_at` | datetime | File modification time | Valid timestamp | From filesystem mtime |

**Constraints**:
- `url` is unique (enforced by map structure)
- Entries are populated on server startup via directory scan
- Updated when downloads complete

**Indexes**:
- Primary: `url` (map key)

## Relationships

```
Download Task (1) ──(1:1)──> Podcast Episode (when completed)
Download Task (1) ──(N:1)──> Podcast Catalog Entry (for duplicate detection)
Podcast Catalog Entry (1) ──(1:1)──> Filesystem Directory
```

## Storage

### In-Memory Storage

**Structures**:

```go
// Task store (in-memory, lost on restart)
tasks map[string]*DownloadTask  // key: task ID
tasksByURL map[string]*DownloadTask  // key: URL (for duplicates)

// Podcast catalog (rebuilt on startup)
catalog map[string]*PodcastCatalogEntry  // key: URL
```

**Thread Safety**:
- Protected by `sync.RWMutex`
- Allows concurrent reads (status queries, list endpoint)
- Exclusive writes (task submission, updates)

### Filesystem Storage

**Structure**:

```
downloads/
├── Podcast Title 1/
│   ├── podcast.m4a          # Audio file (required)
│   ├── cover.jpg            # Cover image (optional)
│   └── shownotes.txt        # Show notes (optional)
├── Podcast Title 2/
│   ├── podcast.m4a
│   └── cover.png
└── ...
```

**Startup Scan Process**:
1. Recursively walk `downloads/` directory
2. For each subdirectory:
   - Check for `*.m4a` file → audio exists
   - Check for `cover.{jpg,png,webp}` → cover exists
   - Check for `shownotes.txt` → notes exist
3. Extract source URL from... [NEEDS CLARIFICATION: How to recover URL from downloaded files?]

**Issue**: Current file structure doesn't store source URL, making duplicate detection impossible after restart.

**Resolution Options**:

| Option | Description | Pros | Cons |
|--------|-------------|------|------|
| A | Add `.metadata.json` file in each directory with URL | Reliable URL recovery, easy to implement | Changes file structure, requires CLI tool update |
| B | Parse shownotes.txt for embedded URL | No file structure change | URL may not be in notes, unreliable |
| C | Rebuild catalog only from in-memory tasks (no filesystem scan) | Simple, no file changes | Lost on restart, violates FR-018 |
| D | Use directory name as key instead of URL | Simple | Title collisions, can't detect same URL with different titles |

**Decision**: Option A - Add `.metadata.json` file

**Rationale**:
- FR-014 requires duplicate detection by URL
- Clarification #3 established: scan downloads directory on startup
- Need persistent URL mapping across restarts
- Minimal change to existing structure
- Can store additional metadata later (tags, ratings, etc.)

**Updated File Structure**:

```
downloads/
├── Podcast Title 1/
│   ├── podcast.m4a
│   ├── cover.jpg
│   ├── shownotes.txt
│   └── .metadata.json     # NEW: URL and metadata
└── ...
```

**Metadata File Format**:

```json
{
  "source_url": "https://www.xiaoyuzhou.fm/episode/...",
  "title": "Podcast Title 1",
  "downloaded_at": "2026-02-08T10:30:00Z",
  "audio_file": "podcast.m4a",
  "cover_file": "cover.jpg",
  "shownotes_file": "shownotes.txt"
}
```

## Validation Rules

### URL Validation

**Format**:
```
^https?://([a-z0-9-]+\.)?xiaoyuzhou\.fm(?:m)?\.com/episode/.+$
```

**Domain Variants Supported**:
- `xiaoyuzhou.fm` - primary domain
- `www.xiaoyuzhou.fm` - primary with www
- `xiaoyuzhoufm.com` - alternative domain
- `www.xiaoyuzhoufm.com` - alternative with www

**Examples**:
- ✅ `https://www.xiaoyuzhou.fm/episode/12345678`
- ✅ `http://xiaoyuzhou.fm/episode/abc123`
- ✅ `https://www.xiaoyuzhoufm.com/episode/12345678` (alternative domain)
- ✅ `http://xiaoyuzhoufm.com/episode/abc123` (alternative domain)
- ❌ `https://google.com/episode/123` (wrong domain)
- ❌ `xiaoyuzhou.fm/episode/123` (missing scheme)
- ❌ `https://www.xiaoyuzhou.fm/podcast/123` (wrong path)

**Test URLs Available**:
See `examples/xiaoyuzhou_urls` for 9 valid Xiaoyuzhou FM episode URLs for testing.

### Task ID Validation

**Format**: UUID v4
```
^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$
```

### Path Validation

**Rules**:
- Must be absolute path
- Must not contain `..` (no directory traversal)
- Audio file must have `.m4a` extension
- Cover file must have `.jpg`, `.png`, or `.webp` extension
- Show notes file must be `.txt` extension
- All paths must exist when returned in API response

## State Machine

### Download Task Lifecycle

```
┌─────────┐
│ PENDING │ (initial state after submission)
└────┬────┘
     │
     │ Start download
     ▼
┌─────────────┐
│ IN_PROGRESS │ (downloading audio, cover, notes)
└─────┬───────┘
      │
      ├───┐
      │   │ Success
      │   ▼
      │ ┌───────────┐
      │ │ COMPLETED │ (all assets downloaded)
      │ └───────────┘
      │
      │ Failure (network, disk space, validation)
      ▼
  ┌─────────┐
  │ FAILED  │ (error message set)
  └─────────┘
```

**State Transition Triggers**:
- `pending → in_progress`: Download starts (goroutine spawned)
- `in_progress → completed`: All assets successfully downloaded
- `in_progress → failed`: Retry limit exceeded or unrecoverable error

**Allowed State Queries**:
- All states: Status query returns current state
- `completed`: Returns podcast episode with file paths
- `failed`: Returns error message

## Serialization

### JSON Representation

**Download Task Response**:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "url": "https://www.xiaoyuzhou.fm/episode/12345678",
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

**Podcast List Item**:

```json
{
  "title": "Example Podcast Episode",
  "source_url": "https://www.xiaoyuzhou.fm/episode/12345678",
  "audio_path": "/home/user/downloads/Example Podcast Episode/podcast.m4a",
  "shownotes": "This is the show notes content...",
  "cover_image": "/home/user/downloads/Example Podcast Episode/cover.jpg"
}
```

**Error Response**:

```json
{
  "error": {
    "code": "INVALID_URL",
    "message": "The provided URL is not from Xiaoyuzhou FM",
    "details": "URL must contain xiaoyuzhou.fm domain"
  }
}
```

## Go Struct Definitions

```go
package taskmanager

import "time"

// TaskStatus represents the current state of a download task
type TaskStatus string

const (
    StatusPending    TaskStatus = "pending"
    StatusInProgress TaskStatus = "in_progress"
    StatusCompleted  TaskStatus = "completed"
    StatusFailed     TaskStatus = "failed"
)

// DownloadTask represents a podcast download request
type DownloadTask struct {
    ID          string        `json:"id"`
    URL         string        `json:"url"`
    Status      TaskStatus    `json:"status"`
    CreatedAt   time.Time     `json:"created_at"`
    StartedAt   *time.Time    `json:"started_at,omitempty"`
    CompletedAt *time.Time    `json:"completed_at,omitempty"`
    Error       string        `json:"error,omitempty"`
    Progress    int           `json:"progress,omitempty"`
    Podcast     *PodcastEpisode `json:"podcast,omitempty"`
}

// PodcastEpisode represents a downloaded podcast with metadata
type PodcastEpisode struct {
    Title          string    `json:"title"`
    SourceURL      string    `json:"source_url"`
    AudioPath      string    `json:"audio_path"`
    CoverPath      string    `json:"cover_path,omitempty"`
    CoverFormat    string    `json:"cover_format,omitempty"`
    ShowNotesPath  string    `json:"shownotes_path,omitempty"`
    DownloadedAt   time.Time `json:"downloaded_at"`
    FileSizeMB     float64   `json:"file_size_mb,omitempty"`
}

// PodcastCatalogEntry represents an entry in the download catalog
type PodcastCatalogEntry struct {
    URL          string    `json:"url"`
    Title        string    `json:"title"`
    Directory    string    `json:"directory"`
    AudioFile    string    `json:"audio_file"`
    HasCover     bool      `json:"has_cover"`
    HasShowNotes bool      `json:"has_shownotes"`
    DownloadedAt time.Time `json:"downloaded_at"`
}

// MetadataFile represents the .metadata.json file in podcast directories
type MetadataFile struct {
    SourceURL      string `json:"source_url"`
    Title          string `json:"title"`
    DownloadedAt   string `json:"downloaded_at"`
    AudioFile      string `json:"audio_file"`
    CoverFile      string `json:"cover_file,omitempty"`
    ShowNotesFile  string `json:"shownotes_file,omitempty"`
}
```

## Considerations for Future Enhancements

**Potential Additions** (out of scope for current feature):
- User-defined tags/ratings
- Playlist or collection support
- Download statistics (total size, count by podcast)
- Search/filter functionality
- Batch download operations
