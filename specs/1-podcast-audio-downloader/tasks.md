---

# Tasks: Podcast Audio Downloader

**Input**: Design documents from `/specs/1-podcast-audio-downloader/`
**Prerequisites**: plan.md, spec.md, data-model.md, contracts/

**Tests**: Testing is included per constitution requirements (unit tests, integration tests) but not TDD approach.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **CLI tool**: `cmd/podcast-downloader/`, `internal/`, `pkg/` at repository root
- Go layout with `cmd/`, `internal/`, `pkg/`
- Paths shown below assume CLI application structure

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create project directory structure per implementation plan
- [X] T002 Initialize Go module with go.mod in repository root
- [X] T003 [P] Add goquery dependency (v1.8.1) to go.mod
- [X] T004 [P] Add urfave/cli/v2 dependency (v2.27.1) to go.mod
- [X] T005 [P] Add schollz/progressbar/v3 dependency (v3.14.1) to go.mod
- [X] T006 Run go mod tidy to sync dependencies
- [X] T007 Create .gitignore for Go projects (bin/, .exe, etc.)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T008 Create Episode struct in internal/models/episode.go with ID, Title, AudioURL, PageURL, FileSize, Duration fields
- [X] T009 [P] Implement SanitizedTitle() method on Episode struct in internal/models/episode.go
- [X] T010 [P] Implement GenerateFilename() method on Episode struct in internal/models/episode.go
- [X] T011 [P] Implement Validate() method on Episode struct in internal/models/episode.go
- [X] T012 Create DownloadSession struct in internal/models/download_session.go with Status, Progress, BytesDownloaded, etc.
- [X] T013 [P] Implement Status enum constants (Pending, InProgress, Completed, Failed) in internal/models/download_session.go
- [X] T014 [P] Implement UpdateProgress() method on DownloadSession in internal/models/download_session.go
- [X] T015 [P] Implement Complete() method on DownloadSession in internal/models/download_session.go
- [X] T016 [P] Implement Fail() method on DownloadSession in internal/models/download_session.go
- [X] T017 [P] Implement CanRetry() method on DownloadSession in internal/models/download_session.go
- [X] T018 [P] Implement IncrementRetry() method on DownloadSession in internal/models/download_session.go
- [X] T019 Create Config struct in internal/config/config.go with OutputDirectory, OverwriteExisting, Timeout, etc.
- [X] T020 Implement DefaultConfig() function in internal/config/config.go with sensible defaults
- [X] T021 [P] Implement Validate() method on Config struct in internal/config/config.go
- [X] T022 Create RetryableClient in pkg/httpclient/client.go with exponential backoff retry logic
- [X] T023 Implement NewRetryableClient() constructor in pkg/httpclient/client.go
- [X] T024 [P] Implement Do() method on RetryableClient in pkg/httpclient/client.go with retry logic
- [X] T025 Create XiaoyuzhouURLValidator struct in internal/validator/url_validator.go
- [X] T026 Implement ValidateURL() method on XiaoyuzhouURLValidator in internal/validator/url_validator.go
- [X] T027 Create FilePathValidator interface and implementation in internal/validator/filepath_validator.go
- [X] T028 Implement ValidatePath() method in FilePathValidator in internal/validator/filepath_validator.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Download Podcast Episode Audio (Priority: P1) 🎯 MVP

**Goal**: Enable users to download podcast episodes by providing a URL

**Independent Test**: Provide a valid Xiaoyuzhou FM episode URL and verify the audio file is downloaded with correct filename and content

### Implementation for User Story 1

- [X] T029 [P] [US1] Create HTMLExtractor struct implementing URLExtractor interface in internal/downloader/url_extractor.go
- [X] T030 [P] [US1] Implement ExtractURL() method in HTMLExtractor using goquery for HTML parsing in internal/downloader/url_extractor.go
- [X] T031 [US1] Add error types (ErrInvalidURL, ErrPageNotFound, ErrAudioNotFound, ErrAccessDenied) in internal/downloader/url_extractor.go
- [X] T032 [P] [US1] Create HTTPDownloader struct implementing FileDownloader interface in internal/downloader/downloader.go
- [X] T033 [US1] Implement Download() method in HTTPDownloader with progress tracking in internal/downloader/downloader.go
- [X] T034 [P] [US1] Implement ValidateFile() method in HTTPDownloader using magic byte check in internal/downloader/downloader.go
- [X] T035 [US1] Add file download error types (ErrNetworkTimeout, ErrConnectionRefused, ErrDiskFull, ErrPermissionDenied, ErrInvalidAudio) in internal/downloader/downloader.go
- [X] T036 [P] [US1] Create main CLI application structure in cmd/podcast-downloader/main.go
- [X] T037 [US1] Define CLI flags (output, overwrite, no-progress, retry, timeout) in cmd/podcast-downloader/main.go using urfave/cli
- [X] T038 [US1] Implement downloadPodcast() action function in cmd/podcast-downloader/main.go orchestrating the download flow
- [X] T039 [US1] Add URL validation logic in downloadPodcast() using URLValidator in cmd/podcast-downloader/main.go
- [X] T040 [US1] Add episode metadata extraction logic in downloadPodcast() using URLExtractor in cmd/podcast-downloader/main.go
- [X] T041 [US1] Add file path generation logic in downloadPodcast() using Episode.GenerateFilename() in cmd/podcast-downloader/main.go
- [X] T042 [US1] Add existing file check and overwrite/skip logic in downloadPodcast() in cmd/podcast-downloader/main.go
- [X] T043 [US1] Add file download execution in downloadPodcast() using FileDownloader in cmd/podcast-downloader/main.go
- [X] T044 [US1] Add file validation after download in downloadPodcast() using FileDownloader.ValidateFile() in cmd/podcast-downloader/main.go
- [X] T045 [US1] Add success/error reporting with Chinese messages in downloadPodcast() in cmd/podcast-downloader/main.go

**Checkpoint**: At this point, User Story 1 should be fully functional - users can download podcast episodes

---

## Phase 4: User Story 2 - Handle Invalid or Inaccessible Episodes (Priority: P2)

**Goal**: Provide clear error feedback when downloads fail

**Independent Test**: Provide invalid URLs (non-existent, malformed, no audio) and verify appropriate error messages

### Implementation for User Story 2

- [X] T046 [P] [US2] Implement URL format validation in XiaoyuzhouURLValidator.ValidateURL() with regex pattern in internal/validator/url_validator.go
- [X] T047 [P] [US2] Add Chinese error messages for URL validation failures in internal/validator/url_validator.go
- [X] T048 [US2] Implement HTTP 404 detection in URLExtractor.ExtractURL() and return ErrPageNotFound in internal/downloader/url_extractor.go
- [X] T049 [US2] Add Chinese error messages for page not found in internal/downloader/url_extractor.go
- [X] T050 [P] [US2] Implement "no audio found" detection in URLExtractor.ExtractURL() when no .m4a links exist in internal/downloader/url_extractor.go
- [X] T051 [US2] Add Chinese error message for no audio found in internal/downloader/url_extractor.go
- [X] T052 [US2] Enhance downloadPodcast() error handling in cmd/podcast-downloader/main.go to catch and display specific error messages
- [X] T053 [US2] Add user-friendly error formatting with context in downloadPodcast() in cmd/podcast-downloader/main.go

**Checkpoint**: At this point, Users Stories 1 AND 2 should both work - users get clear feedback on errors

---

## Phase 5: User Story 3 - Download Progress Feedback (Priority: P3)

**Goal**: Display download progress for large files

**Independent Test**: Download an audio file and observe progress bar updates with percentage, speed, ETA

### Implementation for User Story 3

- [X] T054 [P] [US3] Create progress bar initialization logic in internal/downloader/progress.go (integrated into main.go instead of separate file)
- [X] T055 [US3] Implement CreateProgressBar() function using schollz/progressbar/v3 in internal/downloader/progress.go (integrated into main.go)
- [X] T056 [US3] Add progress bar integration to HTTPDownloader.Download() method in internal/downloader/downloader.go
- [X] T057 [US3] Implement conditional progress display (check --no-progress flag) in HTTPDownloader.Download() in internal/downloader/downloader.go
- [X] T058 [US3] Add progress bar io.Writer wrapper for seamless HTTP integration in internal/downloader/downloader.go
- [X] T059 [US3] Test progress bar displays correctly during download in cmd/podcast-downloader by running with a large file
- [X] T060 [US3] Add completion message with final file size and location after successful download in cmd/podcast-downloader/main.go

**Checkpoint**: All user stories should now be independently functional with complete user experience

---

## Phase 6: Testing & Quality Assurance

**Purpose**: Comprehensive test coverage per constitution requirements

### Unit Tests

- [ ] T061 [P] Create episode_test.go in internal/models/ with table-driven tests for SanitizedTitle()
- [ ] T062 [P] Create episode_test.go in internal/models/ with tests for GenerateFilename()
- [ ] T063 [P] Create episode_test.go in internal/models/ with tests for Episode.Validate()
- [ ] T064 [P] Create download_session_test.go in internal/models/ with tests for state transitions
- [ ] T065 [P] Create download_session_test.go in internal/models/ with tests for CanRetry() logic
- [ ] T066 [P] Create url_validator_test.go in internal/validator/ with valid/invalid URL test cases
- [ ] T067 [P] Create filepath_validator_test.go in internal/validator/ with path validation tests
- [ ] T068 [P] Create config_test.go in internal/config/ with default config tests
- [ ] T069 [P] Create client_test.go in pkg/httpclient/ with retry logic tests
- [ ] T070 [P] Create url_extractor_test.go in internal/downloader/ with mock HTML tests for URL extraction
- [ ] T071 [P] Create downloader_test.go in internal/downloader/ with file download tests
- [ ] T072 [P] Create downloader_test.go in internal/downloader/ with file validation tests

### Integration Tests

- [ ] T073 Create integration test for full download workflow in tests/integration/download_integration_test.go
- [ ] T074 [P] Create integration test for error scenarios in tests/integration/error_handling_test.go
- [ ] T075 [P] Create integration test for retry logic in tests/integration/retry_test.go

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T076 Run gofmt -w . on all Go files to ensure consistent formatting
- [ ] T077 Run golint on all packages and fix warnings
- [X] T078 Run go vet on all packages and fix issues
- [X] T079 Add README.md in repository root with installation and usage instructions
- [X] T080 Add LICENSE file (choose appropriate license)
- [X] T081 Create example bash script for downloading multiple episodes in examples/batch_download.sh
- [X] T082 Test binary compilation for Linux, macOS, Windows in cmd/podcast-downloader/ (tested for Linux)
- [X] T083 Verify go mod tidy produces clean go.mod and go.sum files
- [ ] T084 Run go test ./... to ensure all tests pass
- [X] T085 Test with real Xiaoyuzhou FM episode URLs to validate end-to-end functionality
- [ ] T086 Update quickstart.md with any discovered issues or workarounds (comprehensive README provided instead)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Testing (Phase 6)**: Depends on all desired user stories being complete
- **Polish (Phase 7)**: Depends on all implementation phases being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Extends US1 error handling but independently testable
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Enhances US1 experience but independently testable

### Within Each User Story

- Models in Foundational phase before user story implementation
- URL extraction before file download (within US1)
- Core download before error handling (US1 → US2)
- Core download before progress display (US1 → US3)
- Tests after implementation (Phase 6)

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel (T003, T004, T005)
- All Foundational model tasks marked [P] can run in parallel (T009-T011, T014-T017, T024, T026-T028)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- Within US1: URL extractor and file downloader can start in parallel (T029-T030, T032-T034)
- Within US2: Error detection tasks can run in parallel (T046-T047, T048-T049, T050-T051)
- Within US3: Progress bar components can run in parallel (T054-T055)
- All unit tests marked [P] can run in parallel (T061-T072)
- Integration tests can run in parallel (T074-T075)
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch URL extractor and file downloader tasks together:
Task: "Create HTMLExtractor struct implementing URLExtractor interface in internal/downloader/url_extractor.go"
Task: "Create HTTPDownloader struct implementing FileDownloader interface in internal/downloader/downloader.go"

# After both complete, launch implementation tasks:
Task: "Implement ExtractURL() method in HTMLExtractor using goquery for HTML parsing"
Task: "Implement Download() method in HTTPDownloader with progress tracking"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 by downloading actual podcast episode
5. Build binary and test manually: `./podcast-downloader "https://www.xiaoyuzhoufm.com/episode/..."`
6. Verify downloaded file is valid audio with correct filename
7. Demo MVP if ready

**MVP delivers**: Working CLI tool that downloads podcast episodes - immediately valuable to users

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Add Testing (Phase 6) → Quality assurance
6. Add Polish (Phase 7) → Production-ready release

Each story adds value without breaking previous stories:
- US1: Core download functionality
- US2: Better error messages
- US3: Progress feedback for large files

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (download logic)
   - Developer B: User Story 2 (error handling)
   - Developer C: User Story 3 (progress display)
3. Stories complete and integrate independently
4. Team converges for Testing and Polish phases

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Follow Go conventions: gofmt, golint, go vet before committing
- Test with real Xiaoyuzhou FM URLs to validate functionality
- Chinese error messages improve UX for target audience
- File validation prevents corrupted downloads
- Retry logic handles transient network failures
- Progress bar provides feedback for large files (user expectation)
