# Implementation Plan: Podcast API Server

**Branch**: `3-podcast-api-server` | **Date**: 2026-02-08 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/3-podcast-api-server/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build an HTTP API server on top of the existing podcast-downloader CLI tool to provide programmatic access to podcast downloading and library management. The server will expose REST endpoints for submitting download tasks (asynchronous processing), querying task status, and listing downloaded podcasts with metadata. Technical approach leverages existing download logic, uses Go's standard library `net/http` for the server, maintains in-memory task state with filesystem-based persistence for completed downloads, and implements retry logic with exponential backoff for external service failures.

## Technical Context

**Language/Version**: Go 1.25.5
**Primary Dependencies**: `net/http` (standard library), `github.com/PuerkitoBio/goquery` (existing), `github.com/fatih/color` (existing), existing downloader packages
**Storage**: In-memory for active tasks (lost on restart), filesystem for downloaded podcasts (scanned on startup), no database
**Testing**: Go testing package (`testing`), table-driven tests for handlers, integration tests for API endpoints
**Target Platform**: Linux/macOS/Windows (local machine or private network)
**Project Type**: Single project (HTTP server application)
**Performance Goals**: <500ms response for task submission, <100ms for status queries, <2s for listing 1000 episodes, support 10 concurrent downloads
**Constraints**: Local/private network only (no public internet exposure), single-instance server (not distributed), no authentication required, trusted environment
**Scale/Scope**: Single user/trusted environment, in-memory task tracking, filesystem-based download catalog, supports 10 concurrent requests

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Status**: ✅ PASSED (Post-Design Re-check)

### Principle Compliance Assessment

| Principle | Requirement | Implementation | Status |
|-----------|-------------|----------------|--------|
| **I. Go Backend Standards** | Go 1.21+, gofmt, error handling, goroutines | Go 1.25.5, reuses existing downloader packages, goroutines for async downloads | ✅ Compliant |
| **II. Service-Oriented Architecture** | RESTful APIs, file-based storage, independent services | REST API endpoints, in-memory + filesystem storage, separate server command | ✅ Compliant |
| **III. Asynchronous Processing First** | Return task ID immediately, status endpoints, job queues | Returns UUID on POST /tasks, GET /tasks/{id} for status, goroutine-based worker pool | ✅ Compliant |
| **IV. External Integration Resilience** | Retry logic with exponential backoff, graceful degradation | 3 retries with exponential backoff (research.md §5), graceful degradation for missing assets | ✅ Compliant |
| **V. Web API First Design** | OpenAPI documentation, JSON format, proper HTTP codes | OpenAPI 3.0 spec (contracts/openapi.yaml), JSON request/response, standard status codes | ✅ Compliant |

### Backend Standards Compliance

| Standard | Requirement | Implementation | Status |
|----------|-------------|----------------|--------|
| **Go Code Quality** | gofmt, golint, go vet, structured logging | Will use gofmt for code, structured logging planned | ✅ Planned |
| **API Design** | RESTful patterns, validation, rate limiting | GET/POST endpoints, URL validation, rate limiting deferred (local-only) | ⚠️ Partial (rate limiting out of scope) |
| **File Management** | Organized structure, atomic operations, cleanup | downloads/{title}/ structure, .metadata.json for persistence | ✅ Compliant |

### Testing Requirements Compliance

| Requirement | Implementation Plan | Status |
|-------------|-------------------|--------|
| **Backend Unit Tests** | Table-driven tests for handlers and task manager | ✅ Planned (research.md §8) |
| **Integration Tests** | Full API endpoint flows with test server | ✅ Planned (research.md §8) |
| **System Testing** | Complete download workflows, error scenarios | ✅ Planned (quickstart.md) |

### Waivers & Justifications

**Waiver #1: Rate Limiting (API Design Standard)**
- **Constitution Requirement**: "Rate limiting to prevent abuse"
- **Implementation**: Not implemented in current feature
- **Justification**: Spec Assumption #4 states "Download tasks do not require authentication or authorization (assumed trusted environment)" and server is for "local machine or private network, not exposed publicly to the internet"
- **Risk**: Low - server is not exposed to untrusted clients
- **Future Consideration**: Add rate limiting if server is deployed publicly

**Waiver #2: Graceful Shutdown (Go Code Quality)**
- **Constitution Requirement**: "Graceful shutdown handling for in-flight processing"
- **Implementation**: Not explicitly planned in initial design
- **Justification**: In-memory task state is lost on restart per Clarification #1; graceful shutdown would only minimally improve UX
- **Risk**: Medium - in-progress downloads would be interrupted anyway
- **Future Consideration**: Implement graceful shutdown for better user experience

### Summary

All core principles are satisfied. Two minor waivers granted for features out of scope (rate limiting, graceful shutdown) due to local/private network deployment model. No constitutional violations identified.

**Post-Design Re-check**: ✅ CONFIRMED - Design remains compliant with all applicable constitutional principles.

## Project Structure

### Documentation (this feature)

```text
specs/3-podcast-api-server/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   └── openapi.yaml     # OpenAPI 3.0 specification for REST API
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
cmd/
├── podcast-downloader/  # Existing CLI tool (unchanged)
│   ├── main.go
│   └── root.go
└── podcast-server/      # NEW: HTTP API server
    └── main.go          # Server entry point

internal/
├── downloader/          # Existing download logic (reused)
│   ├── downloader.go
│   ├── url_extractor.go
│   ├── metadata.go
│   ├── image_downloader.go
│   └── shownotes_saver.go
├── models/              # Existing data models (extended)
│   ├── episode.go
│   └── download_session.go
├── validator/           # Existing validators (reused)
│   ├── url_validator.go
│   └── filepath_validator.go
├── server/              # NEW: HTTP server components
│   ├── server.go        # Server setup and routing
│   ├── handlers.go      # HTTP request handlers
│   ├── middleware.go    # Middleware (logging, recovery, etc.)
│   └── response.go      # Response utilities
├── taskmanager/         # NEW: Task management
│   ├── manager.go       # Task lifecycle management
│   ├── store.go         # In-memory task store
│   └── scanner.go       # Downloads directory scanner
└── config/              # Existing configuration (extended)
    └── config.go

pkg/
└── httpclient/          # Existing HTTP client (reused)
    └── client.go

tests/
├── unit/                # NEW: Unit tests
│   ├── handler_test.go
│   ├── manager_test.go
│   └── scanner_test.go
└── integration/         # NEW: Integration tests
    └── api_test.go

downloads/              # Existing: Downloaded files (unchanged)
```

**Structure Decision**: Single project structure (Go standard layout). The HTTP server is added as a new command (`cmd/podcast-server/`) that reuses existing downloader packages. New server-specific logic is organized under `internal/server/` and `internal/taskmanager/`, maintaining clear separation of concerns. This follows Go conventions and keeps the server code modular and testable.

## Implementation Phases

### Phase 0: Research ✅ COMPLETE

**Status**: Completed
**Output**: `research.md`

**Deliverables**:
- Technology stack decisions (HTTP framework, task IDs, concurrency model)
- Integration strategy with existing codebase
- Testing and security approach
- Performance analysis

**Key Decisions**:
- Use `net/http` standard library (no external framework)
- UUID v4 for task identifiers
- Goroutines with mutex for concurrency
- Directory scanning with `.metadata.json` for URL persistence
- Exponential backoff retry (3 attempts)

### Phase 1: Design & Contracts ✅ COMPLETE

**Status**: Completed
**Outputs**: `data-model.md`, `contracts/openapi.yaml`, `quickstart.md`

**Deliverables**:
- Data model with entities (Download Task, Podcast Episode, Catalog Entry)
- OpenAPI 3.0 specification for REST API
- Development setup and usage guide
- Agent context updates (CLAUDE.md)

**Key Design Decisions**:
- In-memory task storage with filesystem catalog
- `.metadata.json` files for URL persistence across restarts
- Three REST endpoints: POST /tasks, GET /tasks/{id}, GET /podcasts
- Status polling mechanism (no WebSockets/SSE)

### Phase 2: Task Breakdown (NEXT)

**Status**: Pending
**Command**: `/speckit.tasks`
**Output**: `tasks.md`

**Will Generate**:
- Dependency-ordered implementation tasks
- Frontend/backend task breakdown (if applicable)
- Test creation tasks
- Documentation tasks

## Dependencies

### External Dependencies

| Package | Version | Purpose | Status |
|---------|---------|---------|--------|
| `net/http` | Go stdlib 1.25.5 | HTTP server | ✅ Available |
| `github.com/PuerkitoBio/goquery` | v1.8.1 | HTML parsing (existing) | ✅ Available |
| `github.com/fatih/color` | v1.18.0 | CLI colors (existing) | ✅ Available |
| `github.com/google/uuid` | TBD | UUID generation | ⚠️ To add |
| `encoding/json` | Go stdlib | JSON serialization | ✅ Available |
| `context` | Go stdlib | Request cancellation | ✅ Available |

### Internal Dependencies

| Package | Purpose | Changes Required |
|---------|---------|------------------|
| `internal/downloader/` | Download logic | Extend for progress callbacks |
| `internal/models/` | Data models | Add task models with JSON tags |
| `internal/validator/` | URL validation | Extend for domain checking |
| `pkg/httpclient/` | HTTP client | Reuse for retry logic |

## Risks & Mitigations

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| **Existing downloader not reusable** | High | Low | Refactor downloader into library during implementation |
| **Performance targets not met** | Medium | Low | In-memory storage should meet targets; profile and optimize if needed |
| **URL persistence complexity** | Medium | Medium | Use simple `.metadata.json` approach (already designed) |
| **Concurrent download conflicts** | Medium | Low | Mutex protection and duplicate detection (already designed) |
| **Filesystem scanning slow for large libraries** | Low | Low | Only scan on startup; cache in memory (acceptable per spec) |

## Open Questions

**Status**: All questions resolved

- ✅ HTTP framework selection → Use `net/http`
- ✅ Task ID format → UUID v4
- ✅ Concurrency model → Goroutines with mutex
- ✅ Restart recovery → Directory scan + `.metadata.json`
- ✅ Retry strategy → 3 attempts with exponential backoff
- ✅ URL persistence → `.metadata.json` file in each podcast directory

## Success Criteria

From spec.md, the following measurable outcomes must be achieved:

- **SC-001**: <500ms response for task submission
- **SC-002**: 95% success rate for valid URLs
- **SC-003**: <2s to list 1000 episodes
- **SC-004**: Support 10 concurrent downloads
- **SC-005**: <100ms for status queries
- **SC-006**: 90% of episodes have all assets (audio, cover, notes)
- **SC-007**: Clear error messages for 100% of failures
- **SC-008**: 99% API availability

## Next Steps

1. **Run `/speckit.tasks`** to generate implementation task breakdown
2. **Begin implementation** following the task order
3. **Write tests** alongside implementation (TDD recommended by constitution)
4. **Update documentation** as implementation evolves

---

**Plan Version**: 1.0
**Last Updated**: 2026-02-08
**Ready for Implementation**: ✅ Yes
