# Tasks: Podcast API Server

**Input**: Design documents from `/specs/003-podcast-api-server/`
**Prerequisites**: plan.md, spec.md, data-model.md, contracts/openapi.yaml, research.md

**Tests**: Per constitution and plan.md, tests are planned. Test tasks are included in each user story phase.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

This is a Go project using standard layout:
- Server command: `cmd/podcast-server/`
- Internal packages: `internal/server/`, `internal/taskmanager/`
- Tests: `tests/unit/`, `tests/integration/`
- Reuses: `internal/downloader/`, `internal/models/`, `internal/validator/`, `pkg/httpclient/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create new server command directory structure `cmd/podcast-server/` with main.go
- [X] T002 Create internal package directories: `internal/server/` and `internal/taskmanager/`
- [X] T003 Add `github.com/google/uuid` dependency to go.mod for UUID generation
- [X] T004 [P] Create test directory structure: `tests/unit/` and `tests/integration/`
- [X] T005 Create basic logger setup in `internal/server/logger.go` with structured logging

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Core Data Models

[-x] T006 Create DownloadTask struct in `internal/taskmanager/task.go` with UUID, URL, status, timestamps, error, progress fields
[-x] T007 Create TaskStatus enum in `internal/taskmanager/task.go` with pending, in_progress, completed, failed states
[X] T008 [P] Create PodcastEpisode struct in `internal/taskmanager/task.go` with title, paths, downloaded_at fields
[X] T009 [P] Create PodcastCatalogEntry struct in `internal/taskmanager/catalog.go` with URL, title, directory paths
[-x] T010 Create MetadataFile struct in `internal/taskmanager/metadata.go` for .metadata.json persistence

### Task Storage & Management

[-x] T011 Create in-memory task store with mutex in `internal/taskmanager/store.go` with tasks map and tasksByURL index
[-x] T012 Implement task creation method in `internal/taskmanager/store.go` with UUID generation
[-x] T013 Implement task retrieval methods in `internal/taskmanager/store.go` with GetByID and GetByURL
[-x] T014 Implement task update methods in `internal/taskmanager/store.go` with status transitions
[X] T015 [P] Implement task deletion/cleanup in `internal/taskmanager/store.go`

### Directory Scanner

[-x] T016 Implement directory scanner in `internal/taskmanager/scanner.go` with ScanDownloadsDirectory function
[-x] T017 Add .metadata.json reading logic in `internal/taskmanager/scanner.go` to recover source URLs
[-x] T018 Build in-memory catalog from scanned directories in `internal/taskmanager/scanner.go`
[X] T019 [P] Handle missing/partial podcast data gracefully in scanner with warnings

### HTTP Server Foundation

[-x] T020 Create server setup in `internal/server/server.go` with Server struct and ListenAndServe
[-x] T021 Implement request routing in `internal/server/server.go` with ServeMux setup
[-x] T022 Create error response helpers in `internal/server/response.go` with JSON error formatting
[X] T023 [P] Create middleware for logging in `internal/server/middleware.go` with request logging
[X] T024 [P] Create middleware for panic recovery in `internal/server/middleware.go` with error recovery
[X] T025 [P] Create JSON request parsing utilities in `internal/server/response.go`

### Configuration

[-x] T026 Extend config package in `internal/config/config.go` with server port, host, downloads directory options
[-x] T027 Add command-line flag parsing in `cmd/podcast-server/main.go` for port, host, downloads-dir, verbose
[-x] T028 Implement server startup in `cmd/podcast-server/main.go` with config loading and directory scanning

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Submit Podcast Download Task (Priority: P1) 🎯 MVP

**Goal**: Users can submit Xiaoyuzhou FM URLs via POST /tasks endpoint to download podcasts asynchronously

**Independent Test**: Send POST request with valid URL, verify task ID returned and download completes in background with audio, cover, and show notes saved

### Models & Validation

- [X] T029 [P] [US1] Extend URL validator in `internal/validator/url_validator.go` with XiaoyuzhouFM domain validation (support both xiaoyuzhou.fm and xiaoyuzhoufm.com)
- [X] T030 [P] [US1] Create POST request body struct in `internal/server/handlers.go` with URL field and validation tags

### Task Manager

- [X] T031 [US1] Implement task submission logic in `internal/taskmanager/manager.go` with duplicate detection
- [X] T032 [US1] Add status check for in-progress URLs in `internal/taskmanager/manager.go` to return existing tasks
- [X] T033 [US1] Add status check for completed URLs in `internal/taskmanager/manager.go` to return catalog entry
- [X] T034 [US1] Implement download executor goroutine in `internal/taskmanager/manager.go` with async processing
- [X] T035 [US1] Integrate existing downloader from `internal/downloader/` in download executor with progress callbacks
- [X] T036 [US1] Implement retry logic with exponential backoff in `internal/taskmanager/manager.go` (3 attempts, 1s→2s→4s)
- [X] T037 [US1] Save .metadata.json file on download completion in `internal/taskmanager/manager.go` with URL and paths

### HTTP Handlers

- [X] T038 [US1] Implement POST /tasks handler in `internal/server/handlers.go` with SubmitTask function
- [X] T039 [US1] Add request validation in SubmitTask handler for URL format and domain
- [X] T040 [US1] Return 202 Accepted with task response for new tasks in SubmitTask handler
- [X] T041 [US1] Return 409 Conflict for in-progress duplicate URLs in SubmitTask handler
- [X] T042 [US1] Return 303 See Other with Location header for completed downloads in SubmitTask handler
- [X] T043 [US1] Return 400 Bad Request for invalid URLs in SubmitTask handler with error details

### Integration

- [X] T044 [US1] Register POST /tasks route in `internal/server/server.go` with SubmitTask handler
- [X] T045 [US1] Wire task manager into server in `cmd/podcast-server/main.go`
- [ ] T046 [US1] Test end-to-end flow: submit URL → task created → download executes → files saved → metadata written

### Tests

- [ ] T047 [P] [US1] Create unit test for URL validation in `tests/unit/handler_test.go`
- [ ] T048 [P] [US1] Create unit test for task submission in `tests/unit/manager_test.go`
- [ ] T049 [P] [US1] Create unit test for duplicate detection in `tests/unit/manager_test.go`
- [ ] T050 [US1] Create integration test for POST /tasks endpoint in `tests/integration/api_test.go` using URLs from `examples/xiaoyuzhou_urls`
- [ ] T051 [US1] Create integration test for duplicate URL handling in `tests/integration/api_test.go`

**Validation**: All acceptance scenarios from US1 satisfied
- ✅ Valid URL → task ID returned, download starts
- ✅ Invalid URL → error response with details
- ✅ Already downloaded → existing info returned
- ✅ In-progress → existing task ID returned
- ✅ Status query → progress information

---

## Phase 4: User Story 2 - Query Download Task Status (Priority: P3)

**Goal**: Users can check download task status via GET /tasks/{id} endpoint

**Independent Test**: Submit download task, query status endpoint repeatedly, verify status transitions from pending→in_progress→completed with correct timestamps and file paths

### HTTP Handlers

- [X] T052 [P] [US2] Implement GET /tasks/{id} handler in `internal/server/handlers.go` with GetTask function
- [X] T053 [US2] Add UUID path parameter validation in GetTask handler
- [X] T054 [US2] Return 200 OK with task response for found tasks in GetTask handler
- [X] T055 [US2] Return 404 Not Found for non-existent task IDs in GetTask handler
- [X] T056 [US2] Include podcast metadata in task response when status=completed in GetTask handler

### Integration

- [X] T057 [US2] Register GET /tasks/{id} route in `internal/server/server.go` with GetTask handler

### Tests

- [ ] T058 [P] [US2] Create unit test for task retrieval in `tests/unit/handler_test.go`
- [ ] T059 [P] [US2] Create unit test for 404 handling in `tests/unit/handler_test.go`
- [ ] T060 [US2] Create integration test for GET /tasks/{id} endpoint in `tests/integration/api_test.go`
- [ ] T061 [US2] Create integration test for status transitions in `tests/integration/api_test.go`

**Validation**: All acceptance scenarios from US3 satisfied
- ✅ Submitted task → current status returned
- ✅ Completed task → success status with file paths
- ✅ Failed task → error status with reason
- ✅ Non-existent ID → 404 error

---

## Phase 5: User Story 3 - List Downloaded Podcasts (Priority: P2)

**Goal**: Users can retrieve list of all downloaded podcasts via GET /podcasts endpoint

**Independent Test**: Download multiple podcasts, call list endpoint, verify all returned with title, audio path, show notes content, and cover image path

### Catalog Management

- [X] T062 [P] [US3] Implement list retrieval in `internal/taskmanager/manager.go` with GetCatalog function
- [X] T063 [P] [US3] Add pagination support in catalog retrieval with limit/offset parameters
- [X] T064 [P] [US3] Handle show notes file reading in catalog list for content inclusion
- [X] T065 [US3] Calculate total count for pagination response in `internal/taskmanager/manager.go`

### HTTP Handlers

- [X] T066 [US3] Implement GET /podcasts handler in `internal/server/handlers.go` with ListPodcasts function
- [X] T067 [US3] Parse limit/offset query parameters in ListPodcasts handler with defaults (limit=100, offset=0)
- [X] T068 [US3] Validate limit/offset ranges in ListPodcasts handler (1≤limit≤1000, offset≥0)
- [X] T069 [US3] Return 200 OK with podcast list response in ListPodcasts handler
- [X] T070 [US3] Handle empty catalog case in ListPodcasts handler (return empty list)
- [X] T071 [US3] Handle partial data case in ListPodcasts handler (missing cover/shownotes fields)
- [X] T072 [US3] Return 400 Bad Request for invalid limit/offset in ListPodcasts handler

### Integration

- [X] T073 [US3] Register GET /podcasts route in `internal/server/server.go` with ListPodcasts handler
- [X] T074 [US3] Wire catalog manager into server startup in `cmd/podcast-server/main.go`

### Tests

- [ ] T075 [P] [US3] Create unit test for catalog listing in `tests/unit/manager_test.go`
- [ ] T076 [P] [US3] Create unit test for pagination in `tests/unit/manager_test.go`
- [ ] T077 [US3] Create integration test for GET /podcasts endpoint in `tests/integration/api_test.go`
- [ ] T078 [US3] Create integration test for empty catalog in `tests/integration/api_test.go`
- [ ] T079 [US3] Create integration test for pagination in `tests/integration/api_test.go`

**Validation**: All acceptance scenarios from US2 satisfied
- ✅ Multiple podcasts → complete list with all metadata
- ✅ No podcasts → empty list
- ✅ Partial data → available fields included, missing indicated
- ✅ Large number → results within time limits

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final refinements, performance optimization, documentation

### Performance & Optimization

- [ ] T080 [P] Profile list endpoint performance with 1000 episodes using pprof
- [ ] T081 [P] Optimize directory scanning if needed for faster startup
- [ ] T082 [P] Add response compression middleware if list responses are large
- [X] T083 Add connection timeouts and keep-alive settings in `internal/server/server.go`

### Error Handling

- [ ] T084 [P] Implement graceful error responses for disk space full scenario
- [ ] T085 [P] Handle network timeout errors from Xiaoyuzhou FM with clear messages
- [X] T086 Add structured error codes in `internal/server/response.go` for all error scenarios

### Observability

- [X] T087 [P] Add startup logging with catalog size in `cmd/podcast-server/main.go`
- [ ] T088 [P] Add request ID tracking in middleware for tracing
- [X] T089 Log all task submissions with URL and task ID in `internal/taskmanager/manager.go`
- [X] T090 Log download completion/failure with details in `internal/taskmanager/manager.go`

### Documentation

- [X] T091 [P] Update README.md with server usage instructions
- [X] T092 [P] Add API documentation link to README pointing to contracts/openapi.yaml
- [X] T093 Create example requests in docs/ or README for all three endpoints
- [X] T094 Document configuration options in quickstart.md or README

### Build & Deployment

- [X] T095 [P] Create Makefile with build, test, run targets
- [ ] T096 [P] Add .golangci.yml linting configuration if needed
- [X] T097 Verify gofmt compliance across all new files
- [X] T098 Run go vet on all new packages
- [X] T099 Build binary for current platform and test execution

---

## Dependencies

### Story Completion Order

```
Phase 1: Setup
  └─> Phase 2: Foundational (BLOCKING - must complete before any user story)
       ├─> Phase 3: US1 (Submit Download Task) 🎯 MVP
       │    ├─> Phase 4: US3 (Query Task Status) - can run in parallel with US2
       │    └─> Phase 5: US2 (List Podcasts) - can run in parallel with US3
       └─> Phase 6: Polish & Cross-Cutting
```

### User Story Dependencies

- **US1 (P1)**: No dependencies on other user stories - foundational for the system
- **US2 (P3)**: Depends on US1 (needs tasks to query)
- **US3 (P2)**: Independent of US1/US2 (can be implemented in parallel after Phase 2)

**Note**: US2 (Query Status) is P3 but depends on US1. US3 (List Podcasts) is P2 but independent. US3 can be implemented before US2 if desired.

---

## Parallel Execution Opportunities

### Within Phases

**Phase 1 (Setup)**: T004, T005 can run in parallel after T001-T003

**Phase 2 (Foundational)**:
- T008, T009, T010 (models) can run in parallel after T006-T007
- T015, T019, T023, T024, T025 can run in parallel

**Phase 3 (US1)**:
- T029, T030 can run in parallel
- T047, T048, T049 (unit tests) can run in parallel
- T050, T051 (integration tests) can run in parallel

**Phase 4 (US2)**:
- T052, T058, T059 (unit tests) can run in parallel

**Phase 5 (US3)**:
- T062, T063, T064, T065 can run in parallel
- T075, T076 (unit tests) can run in parallel

**Phase 6 (Polish)**: Most tasks can run in parallel

### Across User Stories

**After Phase 2 completes**: US2 (Phase 4) and US3 (Phase 5) can be developed in parallel by different team members

---

## Implementation Strategy

### MVP Scope (Minimum Viable Product)

**Deliverable**: Just Phase 3 (User Story 1) - Submit download tasks

**Value**: Users can programmatically submit podcast download tasks via HTTP API

**Tasks**: T001-T051 (51 tasks)

**Timeline Estimate**: Complete US1 first to deliver core value, then add US2 and US3

### Incremental Delivery

1. **Sprint 1**: MVP = US1 (submit tasks)
2. **Sprint 2**: Add US3 (list podcasts) - independent, provides visibility
3. **Sprint 3**: Add US2 (query status) - depends on US1, completes the API

### Success Criteria

Each user story is independently testable and delivers value:
- US1: Core download functionality
- US3: Library visibility (no US1 dependency)
- US2: Progress tracking (requires US1)

---

## Task Summary

- **Total Tasks**: 99
- **Setup Phase**: 5 tasks
- **Foundational Phase**: 23 tasks (BLOCKING)
- **US1 Phase**: 23 tasks (MVP)
- **US2 Phase**: 10 tasks
- **US3 Phase**: 18 tasks
- **Polish Phase**: 20 tasks

**Parallelizable Tasks**: 45 tasks marked with [P]

**Critical Path**: T001-T028 → T029-T046 → T052-T057 (for US1 + US2)

**MVP (US1 only)**: 51 tasks

---

**Tasks Version**: 1.0
**Generated**: 2026-02-08
**Ready for Implementation**: ✅ Yes
