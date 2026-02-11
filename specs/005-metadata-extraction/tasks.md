# Tasks: Podcast Metadata Extraction and Display

**Input**: Design documents from `/specs/005-metadata-extraction/`
**Prerequisites**: plan.md, spec.md, data-model.md, contracts/

**Tests**: Not explicitly requested in feature specification. Test tasks excluded.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Review existing project structure and prepare for metadata feature

- [ ] T001 Review existing backend structure in `backend/internal/downloader/`
- [ ] T002 Review existing frontend types in `frontend/src/types/podcast.ts`
- [ ] T003 Review existing API handlers in `backend/internal/handlers/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core data structures and utilities that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 [P] Create `PodcastMetadata` struct in `backend/internal/models/metadata.go`
- [x] T005 [P] Create metadata scanner in `backend/pkg/scanner/metadata_scanner.go`
- [x] T006 Add JSON serialization methods for metadata with proper null handling

**Checkpoint**: Foundation ready - metadata model and scanner available for user stories

---

## Phase 3: User Story 1 - Extract and Store Podcast Metadata (Priority: P1) 🎯 MVP

**Goal**: When downloading a podcast episode, automatically extract duration and publish time from the page and save as `.metadata.json`

**Independent Test**: Download a podcast episode and verify `.metadata.json` file is created with correct duration and publish time values

### Implementation for User Story 1

- [x] T007 [US1] Implement metadata extractor in `backend/internal/downloader/metadata_extractor.go`
- [x] T008 [US1] Implement metadata writer in `backend/internal/downloader/metadata_writer.go`
- [x] T009 [US1] Integrate metadata extraction into download flow in `backend/internal/services/download_service.go`
- [x] T010 [US1] Add error handling for extraction failures (continue download, log warning)
- [x] T011 [US1] Update download task status to include `extracting_metadata` state
- [x] T012 [US1] Add overwrite support for `.metadata.json` when re-downloading

**Checkpoint**: At this point, User Story 1 should be fully functional - downloading podcasts creates `.metadata.json` files

---

## Phase 4: User Story 2 - Display Metadata in Podcast List (Priority: P2)

**Goal**: Display duration and publish time in the podcast list interface

**Independent Test**: Open podcast list view and verify each episode displays duration and publish time (e.g., "231分钟", "2个月前")

### Backend Implementation for User Story 2

- [x] T013 [US2] Update podcast list handler to include metadata in `backend/internal/handlers/episodes.go`
- [x] T014 [US2] Update response DTO to include metadata field (Episode model updated)
- [x] T015 [US2] Handle missing metadata gracefully (null values in response)

### Frontend Implementation for User Story 2

- [x] T016 [P] [US2] Update TypeScript types in `frontend/src/types/episode.ts` to include metadata
- [x] T017 [P] [US2] Podcast service uses updated types automatically via API
- [x] T018 [US2] Update EpisodeCard component in `frontend/src/components/episodes/EpisodeCard.vue` to display duration
- [x] T019 [US2] Update EpisodeCard component to display publish time
- [x] T020 [US2] Add placeholder display ("--") for missing metadata

**Checkpoint**: At this point, User Stories 1 AND 2 should both work - downloads create metadata files, and list displays them

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T021 [P] Add logging for metadata extraction operations (in download_service.go)
- [x] T022 [P] Format Go code with gofmt across all modified files
- [x] T023 Run quickstart.md validation steps (backend builds successfully)
- [x] T024 Update CLAUDE.md with any new patterns learned (no new patterns)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-4)**: All depend on Foundational phase completion
  - User Story 1 (P1) should be completed before User Story 2 (P2) for logical flow
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
  - Creates the `.metadata.json` files that User Story 2 displays
- **User Story 2 (P2)**: Can start after User Story 1 is functional
  - Depends on metadata files existing (created by US1)
  - Backend and frontend tasks within US2 can run in parallel

### Within Each User Story

- Models/utilities before integration
- Core implementation before error handling
- Story complete before moving to next priority

### Parallel Opportunities

- T004, T005, T006 (Foundational) can run in parallel
- T007, T008 (US1 extraction/writer) can run in parallel
- T016, T017 (US2 frontend types/service) can run in parallel
- T018, T019 (US2 UI display) can run in parallel
- T021, T022, T023, T024 (Polish) can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch metadata extractor and writer together:
Task: "Implement metadata extractor in backend/internal/downloader/metadata_extractor.go"
Task: "Implement metadata writer in backend/internal/downloader/metadata_writer.go"

# Then integrate:
Task: "Integrate metadata extraction into download flow in backend/internal/downloader/downloader.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test that downloads create `.metadata.json` files correctly
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test downloads create metadata → Deploy/Demo (MVP!)
3. Add User Story 2 → Test list displays metadata → Deploy/Demo
4. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (backend extraction)
   - Developer B: User Story 2 backend (API changes)
   - Developer C: User Story 2 frontend (UI display)
3. Stories complete and integrate

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- No tests requested, but manual testing per quickstart.md is required
