# Implementation Plan: Podcast Metadata Extraction and Display

**Branch**: `005-metadata-extraction` | **Date**: 2026-02-09 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/005-metadata-extraction/spec.md`

## Summary

Extend the podcast downloader to extract episode metadata (duration and publish time) from Xiaoyuzhou FM podcast pages during download, save it as `.metadata.json`, and display this information in the podcast list interface. This builds on the existing Go-based CLI downloader and web service architecture.

## Technical Context

**Language/Version**: Go 1.21+ (existing project standard)
**Primary Dependencies**:
- `github.com/PuerkitoBio/goquery` (existing) - HTML parsing for metadata extraction
- Standard library: `encoding/json`, `net/http`, `os`
**Storage**: File-based (existing downloads directory structure)
**Testing**: Go testing (`go test`), table-driven tests
**Target Platform**: Linux server (backend), Web browser (frontend display)
**Project Type**: Web application (backend + frontend)
**Performance Goals**: Metadata extraction < 500ms per episode, list display < 1 second
**Constraints**: No database, file-based storage only, graceful degradation when metadata unavailable
**Scale/Scope**: Single-user local deployment, hundreds of podcast episodes

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Go Backend Standards | ✓ Pass | Using existing Go 1.21+ with gofmt, error handling |
| II. Service-Oriented Architecture | ✓ Pass | Backend API + frontend display, file-based storage |
| III. Asynchronous Processing First | ✓ Pass | Metadata extraction part of download flow, no separate async needed |
| IV. External Integration Resilience | ✓ Pass | HTML scraping with retry logic, graceful degradation on failure |
| V. Web API First Design | ✓ Pass | Extend existing API endpoints for metadata |

**All gates pass.** Proceeding to Phase 0.

## Project Structure

### Documentation (this feature)

```text
specs/005-metadata-extraction/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
backend/
├── internal/
│   ├── downloader/      # Extend existing downloader
│   │   ├── metadata_extractor.go   # NEW: Extract duration/publish time from HTML
│   │   └── metadata_writer.go      # NEW: Write .metadata.json files
│   ├── models/
│   │   └── metadata.go             # NEW: PodcastMetadata struct
│   └── handlers/
│       └── podcast_handler.go      # MODIFY: Add metadata to list response
├── pkg/
│   └── scanner/
│       └── metadata_scanner.go     # NEW: Scan .metadata.json files

downloads/               # Existing structure
├── Podcast Title/
│   ├── podcast.m4a
│   ├── cover.jpg
│   ├── shownotes.txt
│   └── .metadata.json              # NEW: Metadata file

frontend/                # Existing Vue 3 + TypeScript project
├── src/
│   ├── components/
│   │   └── PodcastList.vue         # MODIFY: Display duration and publish time
│   ├── types/
│   │   └── podcast.ts              # MODIFY: Add metadata fields
│   └── services/
│       └── podcastService.ts       # MODIFY: Handle metadata in API responses
```

**Structure Decision**: Extend the existing backend/frontend split. Add metadata extraction logic to the downloader package, create scanner for reading metadata files, and extend frontend types and components to display the new fields.

## Complexity Tracking

> No constitution violations. This is a straightforward extension of existing functionality.
