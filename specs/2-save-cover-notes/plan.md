# Implementation Plan: Save Cover Images and Show Notes

**Branch**: `2-save-cover-notes` | **Date**: 2026-02-08 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/2-save-cover-notes/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Enhance the existing podcast downloader CLI tool to automatically download and save cover images and show notes alongside audio files. The feature extends the current HTML extraction logic to retrieve cover image URLs and show notes content from podcast pages, then downloads and saves these assets with consistent filename conventions. This is a CLI enhancement (constitutional exception justified as utility tool extension), not a web application.

## Technical Context

**Backend**: Go 1.21+ CLI application (existing podcast-downloader tool)
**Frontend**: None (CLI tool)
**Primary Dependencies**:
- github.com/PuerkitoBio/goquery (v1.8.1) - HTML parsing (existing)
- github.com/urfave/cli/v2 (v2.27.1) - CLI framework (existing)
- github.com/schollz/progressbar/v3 (v3.14.1) - Progress display (existing)
**Storage**: File-based storage (downloads/ directory with audio, images, and text files)
**Testing**: Go testing with table-driven tests
**Target Platform**: CLI tool (Linux/macOS/Windows)
**Project Type**: CLI utility tool enhancement (constitutional exception for utility tools)
**Performance Goals**:
- Cover image download: <5 seconds for images under 2MB
- Show notes extraction: <2 seconds per episode
- No degradation to existing audio download performance
**Constraints**:
- Must maintain backward compatibility with existing CLI interface
- Graceful degradation: audio download continues even if cover/show notes fail
- Must handle HTML structure variations in podcast pages
**Scale/Scope**:
- Single-user CLI tool (no concurrent users)
- Typical usage: 1-10 episodes per session
- Handle large show notes (10,000+ characters) without performance issues
**External Services**: Xiaoyuzhou FM website scraping (existing)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **Go Backend Standards**: Code uses gofmt, follows Go conventions, includes proper error handling
- [x] **Service-Oriented Architecture**: Clear separation of concerns (extraction, download, file management)
- [ ] **Asynchronous Processing First**: N/A - CLI tool uses sequential processing (constitutional exception justified for CLI utility)
- [x] **External Integration Resilience**: Retry logic, graceful degradation for failed downloads
- [ ] **Web API First Design**: N/A - CLI tool (constitutional exception for utility tools)
- [ ] **Frontend Standards**: N/A - CLI tool (no frontend)
- [x] **Testing Requirements**: Unit tests with table-driven tests, integration tests for download flows

**Constitutional Exceptions**:
1. **Service-Oriented Architecture (Web API First)**: This is a CLI tool enhancement, not a web service. The constitution allows exceptions for utility tools where CLI interfaces are more appropriate.
2. **Asynchronous Processing First**: CLI tool uses sequential processing which is simpler and more appropriate for single-user utility tools. Async complexity not justified for this use case.

## Project Structure

### Documentation (this feature)

```text
specs/2-save-cover-notes/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (N/A for CLI tool)
└── tasks.md             # Phase 2 output (NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
cmd/podcast-downloader/
├── main.go              # CLI application entry point (existing)

internal/
├── downloader/          # Download logic (existing)
│   ├── downloader.go    # Main download service (existing)
│   ├── url_extractor.go # HTML parsing and URL extraction (existing - to be extended)
│   ├── http_client.go   # HTTP client with retry logic (existing)
│   ├── cover_extractor.go   # NEW: Cover image URL extraction
│   ├── shownotes_extractor.go # NEW: Show notes content extraction
│   ├── image_downloader.go   # NEW: Cover image download service
│   └── shownotes_saver.go    # NEW: Show notes file writer
├── models/              # Data structures (existing)
│   ├── episode.go       # Episode metadata (existing - to be extended)
│   └── download_session.go   # Download session tracking (existing - to be extended)
├── validator/           # Input validation (existing)
│   └── url_validator.go # URL validation (existing)
└── config/              # Configuration (existing)
    └── config.go        # CLI configuration (existing)

pkg/
└── httpclient/          # Reusable HTTP client with retry logic (existing)
    └── client.go        # HTTP client implementation (existing)

downloads/               # Default download directory (existing)
```

**Structure Decision**: Extending the existing CLI tool structure by adding new extraction and download services to the `internal/downloader/` package. This maintains consistency with the current architecture while adding the new functionality. No frontend or web service components needed.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| N/A | No constitutional violations - all exceptions are pre-justified in the constitution for CLI utility tools | N/A |

