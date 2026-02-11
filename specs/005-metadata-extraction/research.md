# Research: Podcast Metadata Extraction and Display

**Date**: 2026-02-09
**Feature**: 005-metadata-extraction

## Research Areas

### 1. HTML Parsing for Metadata Extraction

**Context**: Extract duration and publish time from Xiaoyuzhou FM podcast pages where these values are contained in DOM elements with class names containing "info".

**Decision**: Use existing `github.com/PuerkitoBio/goquery` library

**Rationale**:
- Already used in the project for HTML parsing (url_extractor.go)
- jQuery-like syntax familiar to developers
- Handles malformed HTML gracefully
- No additional dependencies needed

**Implementation Pattern**:
```go
// Find elements with class containing "info"
doc.Find("[class*='info']").Each(func(i int, s *goquery.Selection) {
    text := strings.TrimSpace(s.Text())
    // Parse duration: look for patterns like "231分钟", "1小时15分钟"
    // Parse publish time: look for patterns like "刚刚发布", "2个月前"
})
```

**Alternatives considered**:
- `golang.org/x/net/html`: Lower-level, more verbose, no CSS selectors
- `github.com/antchfx/htmlquery`: XPath-based, steeper learning curve

### 2. JSON File Storage Format

**Decision**: Use standard JSON with struct tags for serialization

**Rationale**:
- Native Go support via `encoding/json`
- Human-readable for debugging
- Easy to extend with new fields
- Works well with file-based storage approach

**File Format**:
```json
{
  "duration": "231分钟",
  "publish_time": "2个月前",
  "episode_title": "Episode Title",
  "podcast_name": "Podcast Name",
  "extracted_at": "2026-02-09T10:30:00Z"
}
```

### 3. Error Handling Strategy

**Decision**: Continue download on metadata extraction failure, save file with null values

**Rationale**:
- Primary goal is audio download; metadata is enhancement
- Users can retry download to refresh metadata
- Logging allows debugging extraction issues

**Implementation**:
- Extract metadata before/during audio download
- If extraction fails, log warning with URL and error
- Still create `.metadata.json` with null fields
- Audio download proceeds normally

### 4. Frontend Data Display

**Decision**: Extend existing PodcastList component with metadata columns

**Rationale**:
- Minimal UI changes needed
- Consistent with existing list display pattern
- Metadata displayed as-is (no conversion needed)

**Display Format**:
- Duration: Display original text (e.g., "231分钟")
- Publish Time: Display original text (e.g., "2个月前")
- Missing data: Show "--" or hide column

### 5. API Contract for Metadata

**Decision**: Extend existing podcast list endpoint with optional metadata field

**Rationale**:
- Backward compatible (metadata is optional)
- Single API call gets all needed data
- Frontend can choose to display or ignore

**API Response Structure**:
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
  ]
}
```

## Open Questions Resolved

None - all technical decisions clear from existing codebase patterns and feature requirements.

## References

- Existing url_extractor.go for goquery usage patterns
- CLAUDE.md for project technology stack
- Feature spec for functional requirements
