# Tasks: Save Cover Images and Show Notes

**Input**: Design documents from `/specs/2-save-cover-notes/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/interfaces.md, quickstart.md

**Tests**: Tests are NOT included in this task list as they were not explicitly requested in the feature specification. Unit and integration tests can be added after implementation following Go best practices.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **CLI Tool**: Go layout with `cmd/`, `internal/`, `pkg/`
- Paths below reflect the existing CLI tool structure being extended

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and verification

- [ ] T001 Verify Go version is 1.21+ using `go version`
- [ ] T002 Run existing tests to establish baseline using `go test ./...`
- [ ] T003 Build existing CLI tool using `go build -o podcast-downloader cmd/podcast-downloader/main.go`
- [ ] T004 Review existing code structure in `internal/downloader/` and `cmd/podcast-downloader/`
- [ ] T005 Create feature branch `2-save-cover-notes` using `git checkout -b 2-save-cover-notes`

**Checkpoint**: Environment validated and ready for implementation

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core data structures and interface changes that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T006 [P] Create `EpisodeMetadata` struct in `internal/downloader/metadata.go` with fields: AudioURL, CoverURL, ShowNotes, Title, EpisodeNumber, PodcastName, PublicationDate
- [ ] T007 [P] Add new error types to `internal/downloader/url_extractor.go`: ErrCoverNotFound, ErrShowNotesNotFound, ErrInvalidImage, ErrImageTooLarge, ErrInvalidEncoding
- [ ] T008 Update `URLExtractor` interface in `internal/downloader/url_extractor.go` to return `(*EpisodeMetadata, error)` instead of `(string, string, error)` - BREAKING CHANGE
- [ ] T009 Update `HTMLExtractor.ExtractURL()` method in `internal/downloader/url_extractor.go` to return `*EpisodeMetadata` with audio URL and title (cover/show notes extraction added in later phases)
- [ ] T010 Update main download workflow in `cmd/podcast-downloader/main.go` to use new `EpisodeMetadata` return type from `ExtractURL()`
- [ ] T011 Add helper functions to `cmd/podcast-downloader/main.go`: `logWarning()` and `logSuccess()` for colored console output

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Download Cover Image with Audio (Priority: P1) 🎯 MVP

**Goal**: Automatically download and save cover images alongside podcast audio files

**Independent Test**: Run `./podcast-downloader <podcast-url>` and verify that a cover image file (JPEG/PNG/WEBP) is created in the same directory as the audio file with a matching base filename

### Implementation for User Story 1

- [ ] T012 [P] [US1] Implement `extractCoverURL()` method in `internal/downloader/url_extractor.go` that selects first `<img>` from `.avatar-container` element (Xiaoyuzhou FM specific)
- [ ] T013 [P] [US1] Create `ImageDownloader` interface in `internal/downloader/image_downloader.go` with methods: `Download()` and `ValidateImage()`
- [ ] T014 [P] [US1] Implement `HTTPImageDownloader` struct in `internal/downloader/image_downloader.go` with `Download()` method that includes retry logic and progress tracking
- [ ] T015 [US1] Implement `ValidateImage()` method in `internal/downloader/image_downloader.go` with magic byte detection for JPEG (FF D8 FF), PNG (89 50 4E 47), and WebP (RIFF....WEBP) formats
- [ ] T016 [US1] Implement `detectFormat()` helper method in `internal/downloader/image_downloader.go` to identify image format from binary header
- [ ] T017 [US1] Update `HTMLExtractor.ExtractURL()` in `internal/downloader/url_extractor.go` to call `extractCoverURL()` and populate `EpisodeMetadata.CoverURL` field
- [ ] T018 [US1] Add cover image filename generation to `cmd/podcast-downloader/main.go` using same base filename as audio with `.jpg` extension (or detected format)
- [ ] T019 [US1] Add cover image download logic to main workflow in `cmd/podcast-downloader/main.go` with graceful degradation (warning on failure, continue with audio)
- [ ] T020 [US1] Add validation step after cover download in `cmd/podcast-downloader/main.go` using `ValidateImage()` to ensure downloaded file is valid

**Checkpoint**: At this point, User Story 1 should be fully functional - cover images download alongside audio with warning messages on failure

---

## Phase 4: User Story 2 - Download Show Notes with Audio (Priority: P2)

**Goal**: Automatically extract and save show notes as UTF-8-BOM text files alongside podcast audio files

**Independent Test**: Run `./podcast-downloader <podcast-url>` and verify that a text file containing show notes is created in the same directory as the audio file with properly formatted content (links, lists, headers) and UTF-8-BOM encoding

### Implementation for User Story 2

- [ ] T021 [P] [US2] Implement `extractShowNotes()` method in `internal/downloader/url_extractor.go` with multi-fallback strategy: (1) `<section aria-label="节目show notes">`, (2) aria-label containing "show notes", (3) semantic selectors, (4) log failure if not found
- [ ] T022 [P] [US2] Create `ShowNotesSaver` interface in `internal/downloader/shownotes_saver.go` with methods: `Save()` and `FormatHTMLToText()`
- [ ] T023 [P] [US2] Implement `PlainTextShowNotesSaver` struct in `internal/downloader/shownotes_saver.go` with `Save()` method that writes UTF-8-BOM encoded text files
- [ ] T024 [US2] Implement UTF-8 validation in `PlainTextShowNotesSaver.Save()` in `internal/downloader/shownotes_saver.go` using `unicode/utf8` package
- [ ] T025 [US2] Implement `FormatHTMLToText()` method in `internal/downloader/shownotes_saver.go` to convert HTML to plain text with structure preservation (links, lists, headers)
- [ ] T026 [US2] Add HTML-to-text conversions to `FormatHTMLToText()` in `internal/downloader/shownotes_saver.go`: links → "text (URL: url)", `<ul>` → bullet points, `<ol>` → numbered lists, headers → uppercase with underline
- [ ] T027 [US2] Update `HTMLExtractor.ExtractURL()` in `internal/downloader/url_extractor.go` to call `extractShowNotes()` and populate `EpisodeMetadata.ShowNotes` field
- [ ] T028 [US2] Add show notes filename generation to `cmd/podcast-downloader/main.go` using same base filename as audio with `.txt` extension
- [ ] T029 [US2] Add show notes save logic to main workflow in `cmd/podcast-downloader/main.go` with graceful degradation (warning on failure, continue with audio)
- [ ] T030 [US2] Add whitespace cleanup in `cmd/podcast-downloader/main.go` or `FormatHTMLToText()` to strip excessive newlines (more than 2 consecutive) from show notes text

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - cover images AND show notes download alongside audio

---

## Phase 5: User Story 3 - Organize All Downloaded Assets (Priority: P3)

**Goal**: Ensure all downloaded files (audio, cover, show notes) use consistent naming for easy organization and management

**Independent Test**: Run `./podcast-downloader <podcast-url>` multiple times with different episodes, then verify that files sort alphabetically with related files adjacent to each other (episode.m4a, episode.jpg, episode.txt)

### Implementation for User Story 3

- [ ] T031 [US3] Verify filename sanitization logic in `cmd/podcast-downloader/main.go` handles Chinese characters, emojis, and special characters consistently across all file types
- [ ] T032 [US3] Confirm base filename generation in `cmd/podcast-downloader/main.go` uses same sanitized title for all three file types: audio (.m4a), cover (.jpg/.png/.webp), and show notes (.txt)
- [ ] T033 [US3] Test file organization with multiple episodes in `downloads/` directory to verify alphabetical sorting groups related files together
- [ ] T034 [US3] Update README.md in repository root with documentation showing example file organization: audio, cover, and show notes files with consistent naming

**Checkpoint**: All user stories now complete - full podcast episode download with organized file structure

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final improvements, documentation, and validation

- [ ] T035 [P] Run `gofmt -w .` to format all Go code according to Go conventions
- [ ] T036 [P] Run `go vet ./...` to check for code issues
- [ ] T037 [P] Build CLI tool using `go build -o podcast-downloader cmd/podcast-downloader/main.go` and verify no compilation errors
- [ ] T038 [P] Manual test with real podcast URL from Xiaoyuzhou FM to verify end-to-end functionality
- [ ] T039 [P] Verify cover image format detection works correctly by checking downloaded file with `file` command
- [ ] T040 [P] Verify show notes file encoding with `file episode.txt` command to confirm "UTF-8 Unicode (with BOM) text"
- [ ] T041 [P] Test graceful degradation by temporarily breaking cover download (e.g., invalid URL) and verify warning message displays but audio download continues
- [ ] T042 Update CLAUDE.md "Recent Changes" section with Feature 2 details if not already done during planning phase
- [ ] T043 Create commit with all changes for feature branch 2-save-cover-notes
- [ ] T044 Run quickstart.md verification checklist to ensure all implementation steps complete

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - User Story 1 (P1): Can start after Foundational - No dependencies on other stories
  - User Story 2 (P2): Can start after Foundational - Uses same foundation as US1 but independent functionality
  - User Story 3 (P3): Can start after Foundational - Validates and documents organization from US1+US2
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - Independent cover image download
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Independent show notes extraction
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Builds on US1+US2 but mainly validation/documentation

**Note**: US1 and US2 can proceed in parallel after Foundational phase if desired, as they operate on different aspects of the download process.

### Within Each User Story

**User Story 1 (Cover Images)**:
- T012-T013 (parallel) - Interface and extraction method
- T014-T016 (sequential) - HTTPImageDownloader implementation
- T017-T020 (sequential) - Integration into main workflow

**User Story 2 (Show Notes)**:
- T021-T022 (parallel) - Extraction method and interface
- T023-T026 (sequential) - PlainTextShowNotesSaver implementation
- T027-T030 (sequential) - Integration into main workflow

**User Story 3 (Organization)**:
- T031-T034 (sequential) - Validation and documentation

### Parallel Opportunities

**Setup Phase (Phase 1)**:
```bash
# These can run together:
Task T001: Verify Go version
Task T002: Run existing tests
Task T003: Build existing CLI tool
Task T004: Review existing code structure
```

**Foundational Phase (Phase 2)**:
```bash
# These can run together (different files):
Task T006: Create EpisodeMetadata struct
Task T007: Add new error types
```

**User Story 1 (Cover Images)**:
```bash
# These can run together (different files):
Task T012: Implement extractCoverURL() in url_extractor.go
Task T013: Create ImageDownloader interface in image_downloader.go
```

**User Story 2 (Show Notes)**:
```bash
# These can run together (different files):
Task T021: Implement extractShowNotes() in url_extractor.go
Task T022: Create ShowNotesSaver interface in shownotes_saver.go
```

**Polish Phase (Phase 6)**:
```bash
# All tasks marked [P] can run together:
Task T035: Run gofmt
Task T036: Run go vet
Task T037: Build CLI tool
Task T038: Manual test with real URL
Task T039: Verify cover image format
Task T040: Verify show notes encoding
Task T041: Test graceful degradation
```

---

## Parallel Example: User Story 1

```bash
# After Foundational phase complete, launch these in parallel:
Task T012: Implement extractCoverURL() in internal/downloader/url_extractor.go
Task T013: Create ImageDownloader interface in internal/downloader/image_downloader.go

# After T012 and T013 complete, continue with:
Task T014: Implement HTTPImageDownloader with Download() method
Task T015: Implement ValidateImage() with magic byte detection
# (T014 and T015 can run in parallel as different methods)
```

---

## Implementation Strategy

### MVP First (User Story 1 Only - Recommended)

1. Complete Phase 1: Setup (15 minutes)
2. Complete Phase 2: Foundational (1 hour) ⚠️ CRITICAL
3. Complete Phase 3: User Story 1 (1.5-2 hours)
4. **STOP and VALIDATE**: Test cover image download independently with real podcast URL
5. Demo/validate MVP functionality

**MVP delivers**: Automatic cover image download alongside audio files with graceful degradation

### Incremental Delivery (Recommended)

1. Complete Setup + Foundational → Foundation ready (1.25 hours)
2. Add User Story 1 → Test independently → Cover images downloading ✅ MVP! (2 hours total)
3. Add User Story 2 → Test independently → Show notes saving ✅ (3.5 hours total)
4. Add User Story 3 → Test independently → Full file organization ✅ (4 hours total)
5. Polish Phase → Production ready (4.5 hours total)

Each story adds value without breaking previous functionality.

### Parallel Team Strategy

With multiple developers (not typical for CLI tool but possible):

1. Team completes Setup + Foundational together (1.25 hours)
2. Once Foundational is done:
   - Developer A: User Story 1 (Cover images) - 1.5 hours
   - Developer B: User Story 2 (Show notes) - 1.5 hours
3. Stories complete and integrate independently
4. Developer A or B: User Story 3 (Organization validation) - 0.5 hours
5. Both: Polish phase - 0.5 hours

**Total time with parallel execution**: ~3.5 hours vs 4.5 hours sequentially

---

## Notes

- [P] tasks = different files or methods, no blocking dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group of tasks
- Stop at any checkpoint to validate story works independently
- **Breaking change alert**: T008 changes the `URLExtractor` interface signature - update all call sites
- **No TDD approach**: Tests added during polish phase, not before implementation (spec doesn't require TDD)
- **CLI tool context**: All paths are for existing CLI tool being extended, not a new project
- **Xiaoyuzhou FM specific**: `.avatar-container` selector and show notes aria-label are specific to this platform

## Task Summary

- **Total Tasks**: 44
- **Setup Phase**: 5 tasks
- **Foundational Phase**: 6 tasks (BLOCKS all user stories)
- **User Story 1 (P1)**: 9 tasks
- **User Story 2 (P2)**: 10 tasks
- **User Story 3 (P3)**: 4 tasks
- **Polish Phase**: 10 tasks
- **Parallel Opportunities**: 24 tasks marked with [P] can run in parallel with appropriate team size

## MVP Scope (Recommended Starting Point)

**Suggested MVP**: Phase 1 + Phase 2 + Phase 3 (User Story 1)

**Delivers**: Automatic cover image download alongside audio files
**Time Estimate**: ~3.5 hours
**Independent Test**: Download podcast and verify cover image appears with audio file
**Value**: High - P1 priority, visual element users immediately notice

After MVP validates, incrementally add User Story 2 (show notes) and User Story 3 (organization validation).
