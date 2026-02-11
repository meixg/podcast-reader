# Quick Start: Podcast Metadata Extraction

**Feature**: 005-metadata-extraction

## Prerequisites

- Go 1.21+ installed
- Node.js 18+ installed (for frontend)
- Existing podcast downloader project set up

## Backend Setup

### 1. Install Dependencies

```bash
cd backend
go mod tidy
```

### 2. Build

```bash
go build -o podcast-server cmd/server/main.go
```

### 3. Run

```bash
./podcast-server
```

Server starts on `http://localhost:8080`

## Frontend Setup

### 1. Install Dependencies

```bash
cd frontend
npm install
```

### 2. Start Dev Server

```bash
npm run dev
```

Frontend available at `http://localhost:5173`

## Testing Metadata Extraction

### Download a Podcast with Metadata

```bash
curl -X POST http://localhost:8080/api/podcasts/download \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.xiaoyuzhoufm.com/episode/..."}'
```

Response:
```json
{
  "task_id": "task-abc123",
  "status": "pending"
}
```

### Check Download Status

```bash
curl http://localhost:8080/api/podcasts/task/task-abc123
```

Response when complete:
```json
{
  "task_id": "task-abc123",
  "status": "completed",
  "audio_path": "/downloads/Podcast Name/episode.m4a",
  "metadata": {
    "duration": "231分钟",
    "publish_time": "2个月前",
    "extracted_at": "2026-02-09T10:30:00Z"
  }
}
```

### List Episodes with Metadata

```bash
curl http://localhost:8080/api/podcasts
```

Response:
```json
{
  "episodes": [
    {
      "title": "Episode Title",
      "podcast_name": "Podcast Name",
      "file_path": "/downloads/...",
      "metadata": {
        "duration": "231分钟",
        "publish_time": "2个月前"
      }
    }
  ],
  "total": 1
}
```

## Verifying Metadata Files

After download, check the `.metadata.json` file:

```bash
cat downloads/Podcast\ Name/.metadata.json
```

Expected output:
```json
{
  "duration": "231分钟",
  "publish_time": "2个月前",
  "episode_title": "Episode Title",
  "podcast_name": "Podcast Name",
  "extracted_at": "2026-02-09T10:30:00Z"
}
```

## Development Workflow

### Adding New Extraction Logic

1. Edit `backend/internal/downloader/metadata_extractor.go`
2. Update selector patterns if page structure changes
3. Test with actual podcast URLs
4. Run tests: `go test ./internal/downloader/...`

### Frontend Development

1. Edit `frontend/src/components/PodcastList.vue`
2. Add metadata display columns
3. Test with mock data or running backend
4. Run lint: `npm run lint`

## Troubleshooting

### Metadata not extracted

- Check page structure hasn't changed (class names containing "info")
- Verify network connectivity to podcast page
- Check server logs for extraction errors

### Metadata file not created

- Verify write permissions in downloads directory
- Check disk space
- Review error logs

### Frontend not showing metadata

- Verify API response includes metadata field
- Check browser console for JavaScript errors
- Ensure TypeScript types are updated

## Next Steps

See [tasks.md](./tasks.md) for detailed implementation tasks.
