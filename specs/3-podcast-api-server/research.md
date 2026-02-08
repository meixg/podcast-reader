# Research: Podcast API Server

**Feature**: 3-podcast-api-server
**Date**: 2026-02-08
**Phase**: Phase 0 - Technical Research & Decision Making

## Overview

This document captures research findings and technical decisions for building an HTTP API server on top of the existing podcast-downloader CLI tool. The server exposes REST endpoints for submitting download tasks, querying status, and listing downloaded podcasts.

## Technology Decisions

### 1. HTTP Framework: Standard Library `net/http`

**Decision**: Use Go's standard library `net/http` package instead of third-party frameworks like Gin, Echo, or Chi.

**Rationale**:
- The spec defines simple CRUD endpoints with no complex routing requirements
- Standard library is sufficient for the defined performance goals (10 concurrent requests, <500ms response time)
- Reduces external dependencies and attack surface
- Aligns with existing codebase philosophy (CLI tool uses minimal dependencies)
- Go 1.25.5 provides excellent performance out of the box
- Easy to test with `net/http/httptest`

**Alternatives Considered**:
- **Gin**: Popular framework with better routing performance, but adds dependency overhead for simple use case
- **Echo**: Similar to Gin, more features than needed
- **Chi**: Lightweight and composable, but standard library is sufficient for 3 endpoints

**Implementation Notes**:
- Use `http.ServeMux` for routing (improved in Go 1.20+ with method matching)
- Implement custom middleware for logging, recovery, and CORS if needed
- JSON encoding/decoding with `encoding/json`
- Context propagation using `context.Context` for cancellation

### 2. Task ID Generation: UUID v4

**Decision**: Use random UUIDs (version 4) for task identifiers.

**Rationale**:
- Provides uniqueness without coordination (single-instance server)
- Standard format easily understood by clients
- No sequential information leakage
- Go's `google/uuid` package is battle-tested
- Alternative: ULID would provide sortability but adds dependency

**Alternatives Considered**:
- **Sequential integers**: Requires state management, reveals usage patterns
- **ULID**: Sortable and URL-safe, but adds external dependency
- **Nano ID**: Similar to ULID, not standard in Go ecosystem

**Implementation Notes**:
- Import `github.com/google/uuid` package
- Generate with `uuid.New()`
- Validate UUID format in handlers

### 3. Concurrency Model: Goroutines with Mutex

**Decision**: Use goroutines for concurrent download execution with mutex-protected shared state for task tracking.

**Rationale**:
- Go's goroutines are lightweight and perfect for I/O-bound operations
- Spec requires support for 10 concurrent downloads (easily within goroutine capabilities)
- Mutex provides simple synchronization for in-memory task store
- No need for complex actor patterns or channels for this scale
- Existing downloader already uses goroutines internally

**Alternatives Considered**:
- **Worker pool pattern**: Good for rate limiting, but adds complexity for 10 concurrent downloads
- **Channel-based orchestration**: More idiomatic Go, but harder to reason about task state queries
- **Existing libraries**: Like ants for worker pools, overkill for this use case

**Implementation Notes**:
- Launch goroutine per download task
- Protect shared task map with `sync.RWMutex` (allow concurrent reads)
- Use `sync.WaitGroup` for graceful shutdown
- Consider buffered channel for work queue if rate limiting needed later

### 4. Downloads Directory Scanning: Recursive File Walk

**Decision**: Use `filepath.Walk` or `filepath.WalkDir` to scan downloads directory on startup.

**Rationale**:
- Standard library provides efficient directory traversal
- `WalkDir` (Go 1.16+) is more efficient than `Walk`
- Existing file structure: `downloads/{Podcast Title}/{podcast.m4a, cover.jpg, shownotes.txt}`
- Can detect directories with audio files to identify completed downloads
- No need for separate index file

**Alternatives Considered**:
- **Separate index database**: Adds persistence complexity, violates "no database" assumption
- **JSON manifest file**: Requires keeping in sync with filesystem, potential inconsistencies
- **Inotify/fsevents**: Too complex for startup-only scan

**Implementation Notes**:
- Scan on server startup before accepting requests
- Parse podcast title from directory name
- Look for `*.m4a` files (audio), `cover.{jpg,png,webp}`, `shownotes.txt`
- Build in-memory map: source URL → podcast metadata
- Handle missing/malformed files gracefully (skip or mark as partial)

### 5. Retry Strategy: Exponential Backoff with Jitter

**Decision**: Implement retry logic with exponential backoff and jitter for Xiaoyuzhou FM requests.

**Rationale**:
- Spec requires up to 3 retries with exponential backoff
- Existing `pkg/httpclient` may already have retry logic (verify during implementation)
- Exponential backoff reduces server load during outages
- Jitter prevents thundering herd problem
- Standard pattern for external service integration

**Alternatives Considered**:
- **Fixed delay retries**: Simpler but can overload downstream service
- **Circuit breaker pattern**: Overkill for single service integration
- **Existing retry libraries**: Like `avast/retry-go`, but adds dependency

**Implementation Notes**:
- Reuse existing retry logic from `pkg/httpclient` if available
- If not, implement: 1s → 2s → 4s delays with ±20% jitter
- Max 3 attempts (1 initial + 3 retries)
- Log retry attempts for observability
- Give up after max retries and mark task as failed

### 6. URL Validation: String Matching + Net URL Parsing

**Decision**: Validate URLs using string matching for domain check and `net/url` for format validation.

**Rationale**:
- Spec requires rejecting non-Xiaoyuzhou FM URLs
- Existing `internal/validator/url_validator.go` can be extended
- `net/url.Parse` validates URL format
- String matching checks for `xiaoyuzhou.fm` domain
- Fast and reliable with standard library

**Alternatives Considered**:
- **Regular expressions**: Powerful but error-prone for URL validation
- **Allow-list of patterns**: Similar to string matching, more complex
- **DNS resolution**: Overkill and slow for validation

**Implementation Notes**:
- Check URL contains `xiaoyuzhou.fm` or `www.xiaoyuzhou.fm`
- Parse with `net/url.Parse` to ensure valid format
- Return 400 Bad Request for invalid URLs
- Return 400 or 422 for valid but non-Xiaoyuzhou URLs

### 7. Error Handling: Structured JSON Responses

**Decision**: Return errors as structured JSON with consistent format.

**Rationale**:
- Spec requires clear error messages with specific reasons
- JSON is already the chosen response format (Assumption #8)
- Standard practice for REST APIs
- Easy for clients to parse and display
- Supports internationalization if needed later

**Response Format**:
```json
{
  "error": {
    "code": "INVALID_URL",
    "message": "The provided URL is not valid: missing scheme",
    "details": "URL must start with http:// or https://"
  }
}
```

**Implementation Notes**:
- Define error codes as constants
- Use HTTP status codes appropriately (400, 404, 500, 503)
- Include error context in logs (don't expose internals to clients)
- Wrap underlying errors with context using `fmt.Errorf`

### 8. Progress Tracking: Status Polling

**Decision**: Use status polling mechanism instead of real-time updates (as specified in Assumption #9).

**Rationale**:
- Spec explicitly states "status polling is sufficient"
- Simpler implementation than WebSockets or SSE
- No need for persistent connections
- Existing CLI tool doesn't have progress updates
- 100ms response time for status queries (SC-005) makes polling viable

**Alternatives Considered**:
- **Server-Sent Events (SSE)**: Real-time but adds complexity
- **WebSockets**: Bidirectional but overkill for one-way status
- **Long polling**: Complexity for marginal benefit

**Implementation Notes**:
- Provide `GET /tasks/{id}` endpoint
- Return current status: pending, in-progress, completed, failed
- Include progress percentage if available from downloader
- Clients poll at reasonable interval (e.g., 1-2 seconds)

## Integration with Existing Code

### Reusing Downloader Logic

**Decision**: Directly import and use existing `internal/downloader` packages.

**Rationale**:
- Spec Assumption #3: "The system will reuse existing download logic"
- Avoids code duplication
- Proven functionality from CLI tool
- Minimal changes required to existing packages

**Required Changes**:
- Extract download logic from CLI command into reusable function
- Add progress callback for status updates
- Ensure error handling works in server context (not just CLI)
- May need to remove CLI-specific dependencies (color, progressbar) from core logic

### Adapting Models

**Decision**: Extend existing models without breaking changes.

**Existing Models**:
- `internal/models/episode.go`: Episode metadata
- `internal/models/download_session.go`: Download session tracking

**Required Additions**:
- Task model with status field
- URL-to-file mapping for duplicate detection
- JSON serialization tags for API responses

## Testing Strategy

### Unit Tests

**Scope**:
- Handler functions (request parsing, response formatting)
- Task manager (state transitions, duplicate detection)
- URL validator
- Directory scanner

**Tools**:
- Go `testing` package
- Table-driven tests for multiple scenarios
- `httptest` for HTTP handler testing

### Integration Tests

**Scope**:
- Full API endpoint flows
- Task lifecycle (submit → in-progress → complete)
- Duplicate detection
- Error scenarios

**Tools**:
- Start test server
- Make real HTTP requests
- Verify filesystem changes
- Mock Xiaoyuzhou FM responses

## Security Considerations

### Local/Private Network Only

**Decision**: No authentication or TLS required (per Assumption #4).

**Rationale**:
- Spec assumes trusted environment
- Not exposed to public internet
- Simplifies implementation

**Recommendations**:
- Document security assumption clearly
- Consider adding basic API key if deployment changes
- Add warning if server listens on 0.0.0.0

### Input Validation

**Required**:
- URL format validation
- Domain whitelist (xiaoyuzhou.fm only)
- Path traversal prevention in file paths
- Request size limits

### Rate Limiting

**Decision**: Not implemented initially (out of scope per spec).

**Future Consideration**:
- Add if server is exposed to untrusted clients
- Simple in-memory rate limiter would suffice

## Performance Considerations

### Response Time Targets

**From Spec**:
- Task submission: <500ms (SC-001)
- Status query: <100ms (SC-005)
- List podcasts: <2s for 1000 episodes (SC-003)

**Strategies**:
- In-memory task storage (fast lookups)
- Efficient directory scanning (cache on startup)
- Pagination for list endpoint
- Async download processing (return immediately)

### Concurrent Downloads

**Target**: Support 10 concurrent downloads (SC-004)

**Implementation**:
- Goroutine per download (lightweight)
- No explicit limit needed at Go level
- Consider semaphore if limiting external service load

### Memory Usage

**Estimates**:
- Task state: ~1KB per task
- 10 concurrent tasks: ~10KB
- Download catalog: ~500 bytes per podcast
- 1000 podcasts: ~500KB
- Total: ~1MB for task + catalog state

**Conclusion**: Memory usage is negligible for modern systems.

## Open Questions

All technical questions have been resolved through this research phase. No NEEDS CLARIFICATION items remain.

## Next Steps

Proceed to **Phase 1: Design & Contracts**
- Generate `data-model.md` with entity definitions
- Create `contracts/openapi.yaml` with API specification
- Write `quickstart.md` with development setup instructions
