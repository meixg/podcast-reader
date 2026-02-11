# Data Model: Podcast Metadata Extraction and Display

**Feature**: 005-metadata-extraction

## Entities

### PodcastMetadata

Represents extracted metadata for a podcast episode.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| duration | string | No | Duration as displayed on page (e.g., "231分钟", "1小时15分钟") |
| publish_time | string | No | Relative publish time (e.g., "刚刚发布", "2个月前") |
| episode_title | string | No | Title of the episode |
| podcast_name | string | No | Name of the podcast series |
| extracted_at | string (ISO 8601) | Yes | Timestamp when metadata was extracted |

**Validation Rules**:
- All fields except `extracted_at` can be null/empty
- `extracted_at` must be valid ISO 8601 timestamp
- Duration and publish_time stored as original text, no parsing/validation

**JSON Schema**:
```json
{
  "type": "object",
  "properties": {
    "duration": { "type": ["string", "null"] },
    "publish_time": { "type": ["string", "null"] },
    "episode_title": { "type": ["string", "null"] },
    "podcast_name": { "type": ["string", "null"] },
    "extracted_at": { "type": "string", "format": "date-time" }
  },
  "required": ["extracted_at"]
}
```

### EpisodeWithMetadata

Represents a podcast episode including its metadata for API responses.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| title | string | Yes | Episode title from filename or metadata |
| podcast_name | string | Yes | Podcast series name |
| file_path | string | Yes | Absolute path to audio file |
| cover_path | string | No | Path to cover image if exists |
| shownotes_path | string | No | Path to shownotes if exists |
| metadata | PodcastMetadata | No | Extracted metadata (null if unavailable) |

### DownloadTask

Extends existing download task with metadata extraction status.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | string | Yes | Unique task identifier |
| url | string | Yes | Source podcast page URL |
| status | string | Yes | pending, downloading, extracting_metadata, completed, failed |
| audio_path | string | No | Path to downloaded audio file |
| metadata | PodcastMetadata | No | Extracted metadata |
| error | string | No | Error message if failed |
| created_at | string | Yes | Task creation timestamp |
| updated_at | string | Yes | Last update timestamp |

## Relationships

```
DownloadTask 1..1 ---> 0..1 PodcastMetadata
EpisodeWithMetadata 1..1 ---> 0..1 PodcastMetadata
```

## State Transitions

### DownloadTask Lifecycle

```
[pending] --start--> [downloading] --audio complete--> [extracting_metadata] --metadata saved--> [completed]
                                           |                           |
                                           v                           v
                                      [failed]                   [completed with null metadata]
```

**Transition Rules**:
- `extracting_metadata` state is brief; failure doesn't block completion
- If metadata extraction fails, task still completes with null metadata
- Error field populated only for audio download failures

## File Storage

### .metadata.json Location

Stored alongside audio file in podcast directory:

```
downloads/
└── Podcast Name/
    ├── episode.m4a
    ├── cover.jpg
    ├── shownotes.txt
    └── .metadata.json          # Hidden file with metadata
```

### Filename Convention

- Always named `.metadata.json` (hidden file)
- Located in same directory as audio file
- Associated with audio file by directory location

## API Data Transfer Objects

### ListEpisodesResponse

```json
{
  "episodes": [
    {
      "title": "Episode Title",
      "podcast_name": "Podcast Name",
      "file_path": "/downloads/Podcast Name/episode.m4a",
      "cover_path": "/downloads/Podcast Name/cover.jpg",
      "shownotes_path": "/downloads/Podcast Name/shownotes.txt",
      "metadata": {
        "duration": "231分钟",
        "publish_time": "2个月前",
        "episode_title": "Episode Title",
        "podcast_name": "Podcast Name",
        "extracted_at": "2026-02-09T10:30:00Z"
      }
    }
  ],
  "total": 1
}
```

### DownloadResponse (with metadata)

```json
{
  "task_id": "task-123",
  "status": "completed",
  "audio_path": "/downloads/Podcast/episode.m4a",
  "metadata": {
    "duration": "231分钟",
    "publish_time": "2个月前",
    "extracted_at": "2026-02-09T10:30:00Z"
  }
}
```

## TypeScript Types (Frontend)

```typescript
interface PodcastMetadata {
  duration?: string;
  publish_time?: string;
  episode_title?: string;
  podcast_name?: string;
  extracted_at: string;
}

interface Episode {
  title: string;
  podcast_name: string;
  file_path: string;
  cover_path?: string;
  shownotes_path?: string;
  metadata?: PodcastMetadata;
}

interface EpisodeListResponse {
  episodes: Episode[];
  total: number;
}
```
