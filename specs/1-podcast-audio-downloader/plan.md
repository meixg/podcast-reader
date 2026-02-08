# Implementation Plan: Podcast Audio Downloader

**Branch**: `1-podcast-audio-downloader` | **Date**: 2026-02-08 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/1-podcast-audio-downloader/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build a command-line tool that downloads podcast audio files from Xiaoyuzhou FM. The tool accepts an episode URL, extracts the direct .m4a audio link from the webpage, and downloads the file to local storage with meaningful filenames. Includes robust error handling, retry logic for network failures, and progress feedback for large downloads.

**Primary Requirements**:
- Accept Xiaoyuzhou FM episode URLs via command-line interface
- Parse HTML to extract direct audio file URLs
- Download audio files with progress tracking
- Handle errors gracefully with clear user feedback
- Support file naming and conflict resolution

## Technical Context

**Backend**: Go 1.21+ command-line application
**Frontend**: Command-line interface (CLI) - no web UI for this feature
**Primary Dependencies**:
- HTTP client: net/http (standard library)
- HTML parsing: NEEDS CLARIFICATION (goquery vs colly vs standard library)
- CLI flag parsing: NEEDS CLARIFICATION (cobra vs flag vs pflag)
- Progress display: NEEDS CLARIFICATION (progressbar vs urfave/cli vs custom)
**Storage**: File-based storage for downloaded audio files
**Testing**: Go testing package with table-driven tests
**Target Platform**: Command-line tool (Linux/macOS/Windows)
**Project Type**: CLI utility (standalone Go binary)
**Performance Goals**: Download speed comparable to browser (within 110%), handle files up to 500MB efficiently
**Constraints**:
- Must handle network interruptions gracefully
- Must respect website rate limits
- Must provide clear error messages in Chinese (since Xiaoyuzhou is Chinese service)
- Must validate downloaded files are actual audio (not HTML error pages)
**Scale/Scope**: Single-user tool, sequential downloads (no concurrent processing needed for MVP)
**External Services**: Xiaoyuzhou FM website (HTTP scraping)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Initial Evaluation**:

- [x] **Go Backend Standards**: Will use Go 1.21+ with gofmt, proper error handling, and Go modules
- [ ] **Service-Oriented Architecture**: ⚠️ **CONSTITUTIONAL EXCEPTION** - This is a CLI tool, not a web service, so REST APIs and frontend/backend separation do not apply
- [ ] **Asynchronous Processing First**: ⚠️ **CONSTITUTIONAL EXCEPTION** - CLI tool uses synchronous execution with progress feedback, not async job queues
- [x] **External Integration Resilience**: Will implement retry logic with exponential backoff and graceful error handling for Xiaoyuzhou FM
- [ ] **Web API First Design**: ⚠️ **CONSTITUTIONAL EXCEPTION** - CLI tool does not expose REST APIs or OpenAPI docs
- [ ] **Frontend Standards**: ⚠️ **CONSTITUTIONAL EXCEPTION** - No Vue.js frontend; uses command-line interface instead
- [x] **Testing Requirements**: Will include Go unit tests and integration tests

**Justification for Exceptions**:
This feature is a **standalone CLI utility** for downloading podcast files, not a web service component. The constitution's web service architecture principles (REST APIs, Vue.js frontend, async job queues) apply to the main podcast processing web application. This CLI tool is a helper utility that may eventually be integrated into the web service but is being developed first as a standalone tool for rapid prototyping and testing.

**Constitutional Principles That Apply**:
- Go Backend Standards (code quality, error handling, modules)
- External Integration Resilience (retry logic, graceful degradation)
- Testing Requirements (Go testing)

---

**Post-Design Re-evaluation**:

All technical decisions align with the constitutional principles that apply:

✅ **Go Backend Standards**: Using Go 1.21+ with gofmt, golint, go vet; comprehensive error handling; Go modules for dependency management
✅ **External Integration Resilience**: Implemented exponential backoff retry logic (3 retries), timeout configuration (30s default), graceful error handling for network failures
✅ **Testing Requirements**: Unit tests for all modules, integration tests for HTTP client, table-driven tests for validation

**Constitutional Exceptions Confirmed as Appropriate**:
- CLI tool architecture is simpler and more appropriate for a single-purpose utility
- Code in `pkg/httpclient` is structured for future reuse in web service backend
- No architectural debt introduced; tool can be refactored into web service if needed

**Final Status**: ✅ **APPROVED** - Proceed to implementation phase

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
# CLI application structure
cmd/podcast-downloader/
├── main.go               # Application entry point
└── root.go               # CLI command definitions

internal/
├── downloader/           # Download logic
│   ├── downloader.go     # Main download service
│   ├── url_extractor.go  # HTML parsing and URL extraction
│   └── progress.go       # Progress tracking
├── models/               # Data structures
│   ├── episode.go        # Episode metadata
│   └── download_session.go
├── validator/            # Input validation
│   └── url_validator.go
└── config/               # Configuration
    └── config.go

pkg/
└── httpclient/           # Reusable HTTP client with retry logic
    └── client.go

downloads/                # Default download directory (configurable)
go.mod
go.sum
```

**Structure Decision**: This is a standalone CLI tool using standard Go project layout. The `cmd/` directory contains the executable, `internal/` contains private application code, and `pkg/` contains code that could be reused in the future web service backend. No frontend is needed for this feature.

## Complexity Tracking

> **Constitutional Exceptions documented and justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Service-Oriented Architecture (no REST APIs) | CLI tool does not need web service architecture | Making this a web service would overcomplicate a simple download utility |
| Asynchronous Processing (no job queues) | CLI uses synchronous execution with progress feedback | Async processing would add complexity without user benefit for single-user tool |
| Web API First (no API documentation) | CLI tool exposes commands, not HTTP endpoints | User interacts via terminal, not API calls |
| Frontend Standards (no Vue.js) | Command-line interface is more appropriate for batch downloading | Building a web UI would be overkill for a utility tool |

**Overall Complexity Assessment**: LOW - This is a straightforward CLI tool with well-defined scope. The constitutional exceptions are justified by the different architectural pattern (CLI vs web service).
