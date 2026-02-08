# Tasks: Podcast Management Web Application

**Input**: Design documents from `/specs/004-frontend-web-app/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are OPTIONAL - only included if explicitly requested in the feature specification.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Web app**: `backend/` and `frontend/` at repository root
- Backend: `backend/cmd/`, `backend/internal/`, `backend/pkg/`
- Frontend: `frontend/src/`, `frontend/tests/`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Create backend directory structure: backend/cmd/server/, backend/internal/{handlers,services,models}, backend/pkg/scanner/
- [ ] T002 Create frontend directory structure: frontend/src/{components,views,services,composables,types,router}/, frontend/tests/{unit,e2e}/
- [ ] T003 [P] Initialize Go module in backend/ with go.mod and required dependencies
- [ ] T004 [P] Initialize Node.js project in frontend/ with package.json (Vue 3, Vite, TypeScript, Tailwind CSS, Vue Router)
- [ ] T005 [P] Configure TypeScript in frontend/tsconfig.json with strict mode enabled
- [ ] T006 [P] Configure Vite in frontend/vite.config.ts with Vue plugin and dev server settings
- [ ] T007 [P] Configure Tailwind CSS in frontend/tailwind.config.js with content paths
- [ ] T008 [P] Configure ESLint and Prettier in frontend/ for code quality
- [ ] T009 [P] Configure Vitest in frontend/ for unit testing
- [ ] T010 [P] Configure Cypress in frontend/ for e2e testing

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T011 Create Episode model in backend/internal/models/episode.go with all fields from data-model.md
- [ ] T012 Create DownloadTask model in backend/internal/models/task.go with status enum and all fields
- [ ] T013 [P] Create TypeScript Episode interface in frontend/src/types/episode.ts
- [ ] T014 [P] Create TypeScript DownloadTask interface in frontend/src/types/task.ts
- [ ] T015 Implement file system scanner in backend/pkg/scanner/ to read downloads directory
- [ ] T016 Create in-memory task queue in backend/internal/services/task_service.go
- [ ] T017 Setup HTTP server with CORS in backend/cmd/server/main.go
- [ ] T018 [P] Create API client wrapper in frontend/src/services/api.ts with error handling
- [ ] T019 [P] Create Vue Router configuration in frontend/src/router/index.ts with routes for /episodes and /tasks
- [ ] T020 [P] Create App.vue with router-view and basic layout structure
- [ ] T021 [P] Create main.ts entry point with Vue app initialization

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - View Downloaded Podcasts (Priority: P1) 🎯 MVP

**Goal**: Users can browse their downloaded podcast collection and view episode details including show notes

**Independent Test**: Download sample podcasts via CLI, open web app, verify all episodes appear with correct metadata, click episode to view show notes in modal

### Backend Implementation for User Story 1

- [ ] T022 [P] [US1] Implement EpisodeService.GetEpisodes() with pagination in backend/internal/services/episode_service.go
- [ ] T023 [P] [US1] Implement EpisodeService.GetShowNotes() in backend/internal/services/episode_service.go
- [ ] T024 [US1] Create GET /api/episodes handler in backend/internal/handlers/episodes.go with pagination params
- [ ] T025 [US1] Create GET /api/episodes/:id/shownotes handler in backend/internal/handlers/episodes.go
- [ ] T026 [US1] Wire up episode handlers to HTTP server in backend/cmd/server/main.go

### Frontend Implementation for User Story 1

- [ ] T027 [P] [US1] Create useEpisodes composable in frontend/src/composables/useEpisodes.ts with state and API calls
- [ ] T028 [P] [US1] Create usePagination composable in frontend/src/composables/usePagination.ts
- [ ] T029 [P] [US1] Create NavigationBar component in frontend/src/components/common/NavigationBar.vue
- [ ] T030 [P] [US1] Create Modal component in frontend/src/components/common/Modal.vue
- [ ] T031 [P] [US1] Create Pagination component in frontend/src/components/common/Pagination.vue
- [ ] T032 [P] [US1] Create EpisodeCard component in frontend/src/components/episodes/EpisodeCard.vue
- [ ] T033 [US1] Create ShowNotesModal component in frontend/src/components/episodes/ShowNotesModal.vue (depends on T030)
- [ ] T034 [US1] Create EpisodeList component in frontend/src/components/episodes/EpisodeList.vue (depends on T031, T032)
- [ ] T035 [US1] Create EpisodesView page in frontend/src/views/EpisodesView.vue (depends on T027, T034)
- [ ] T036 [US1] Add empty state handling to EpisodeList component

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Create Download Tasks (Priority: P2)

**Goal**: Users can initiate new podcast downloads directly from the web interface

**Independent Test**: Open tasks page, click create button, enter valid podcast URL, verify task appears in list

### Backend Implementation for User Story 2

- [ ] T037 [P] [US2] Implement TaskService.CreateTask() with URL validation in backend/internal/services/task_service.go
- [ ] T038 [P] [US2] Implement TaskService.GetTasks() in backend/internal/services/task_service.go
- [ ] T039 [P] [US2] Add duplicate URL detection to TaskService in backend/internal/services/task_service.go
- [ ] T040 [US2] Create POST /api/tasks handler in backend/internal/handlers/tasks.go
- [ ] T041 [US2] Create GET /api/tasks handler in backend/internal/handlers/tasks.go
- [ ] T042 [US2] Wire up task handlers to HTTP server in backend/cmd/server/main.go

### Frontend Implementation for User Story 2

- [ ] T043 [P] [US2] Create useTasks composable in frontend/src/composables/useTasks.ts with state and API calls
- [ ] T044 [P] [US2] Create useTaskPolling composable in frontend/src/composables/useTaskPolling.ts for automatic updates
- [ ] T045 [P] [US2] Create TaskCard component in frontend/src/components/tasks/TaskCard.vue
- [ ] T046 [P] [US2] Create TaskStatusBadge component in frontend/src/components/tasks/TaskStatusBadge.vue
- [ ] T047 [US2] Create CreateTaskModal component in frontend/src/components/tasks/CreateTaskModal.vue (depends on T030)
- [ ] T048 [US2] Create TaskList component in frontend/src/components/tasks/TaskList.vue (depends on T045, T046)
- [ ] T049 [US2] Create TasksView page in frontend/src/views/TasksView.vue (depends on T043, T044, T048)
- [ ] T050 [US2] Add URL validation to CreateTaskModal with error messages
- [ ] T051 [US2] Add empty state handling to TaskList component

**Checkpoint**: At this point, User Story 2 should be fully functional and testable independently

---

## Phase 5: User Story 3 - Monitor Download Progress (Priority: P3)

**Goal**: Users can monitor real-time progress of active downloads

**Independent Test**: Create a download task, verify progress updates automatically, verify status changes from pending → downloading → completed

### Backend Implementation for User Story 3

- [ ] T052 [P] [US3] Implement background worker in backend/internal/services/task_worker.go to process pending tasks
- [ ] T053 [P] [US3] Add progress tracking to download operations in backend/internal/services/task_service.go
- [ ] T054 [P] [US3] Integrate existing CLI downloader with task service for actual downloads
- [ ] T055 [US3] Add task status update methods to TaskService (UpdateProgress, MarkCompleted, MarkFailed)

### Frontend Implementation for User Story 3

- [ ] T056 [P] [US3] Add progress bar component to TaskCard for downloading status
- [ ] T057 [P] [US3] Update useTaskPolling to handle progress updates in real-time
- [ ] T058 [US3] Add visual indicators for different task statuses (pending, downloading, completed, failed)
- [ ] T059 [US3] Add error message display for failed tasks in TaskCard

**Checkpoint**: At this point, User Story 3 should be fully functional and testable independently

---

## Phase 6: Final Polish & Testing

**Purpose**: Cross-cutting concerns, testing, and production readiness

- [ ] T060 [P] Add loading states to all API calls in composables
- [ ] T061 [P] Add error boundaries and fallback UI for component errors
- [ ] T062 [P] Implement responsive design for mobile/tablet viewports
- [ ] T063 [P] Add accessibility attributes (ARIA labels, keyboard navigation)
- [ ] T064 [P] Write unit tests for composables (useEpisodes, useTasks, usePagination)
- [ ] T065 [P] Write unit tests for API client wrapper
- [ ] T066 [P] Write component tests for key components (EpisodeCard, TaskCard, Modal)
- [ ] T067 [P] Write E2E tests for User Story 1 (view episodes, pagination, show notes modal)
- [ ] T068 [P] Write E2E tests for User Story 2 (create task, duplicate detection)
- [ ] T069 [P] Write E2E tests for User Story 3 (monitor progress, status updates)
- [ ] T070 [P] Add backend unit tests for EpisodeService and TaskService
- [ ] T071 [P] Add backend integration tests for HTTP handlers
- [ ] T072 Optimize bundle size and add code splitting for routes
- [ ] T073 Add production build configuration and environment variables
- [ ] T074 Create README.md with setup and deployment instructions

---

## Dependencies & Execution Order

### Critical Path (Must Complete Sequentially)

**Phase 1 → Phase 2 → User Stories**
- Phase 1 (Setup) must complete before Phase 2
- Phase 2 (Foundational) must complete before ANY user story work
- User Stories 1, 2, 3 can be implemented in parallel after Phase 2

### Within Each User Story

**Backend First, Then Frontend**
- Backend services and handlers must be complete before frontend integration
- Frontend can develop UI components in parallel with backend, but integration requires backend completion

### Component Dependencies

**User Story 1 (Episodes)**
- T030 (Modal) must complete before T033 (ShowNotesModal)
- T031 (Pagination) and T032 (EpisodeCard) must complete before T034 (EpisodeList)
- T027 (useEpisodes) and T034 (EpisodeList) must complete before T035 (EpisodesView)

**User Story 2 (Tasks)**
- T030 (Modal) must complete before T047 (CreateTaskModal)
- T045 (TaskCard) and T046 (TaskStatusBadge) must complete before T048 (TaskList)
- T043 (useTasks), T044 (useTaskPolling), and T048 (TaskList) must complete before T049 (TasksView)

---

## Parallel Opportunities

### Maximum Parallelization Strategy

**Phase 1 (Setup)**: All tasks T003-T010 can run in parallel after T001-T002 complete

**Phase 2 (Foundational)**:
- T013-T014 (TypeScript types) can run in parallel with T011-T012 (Go models)
- T018-T021 (Frontend setup) can run in parallel with T015-T017 (Backend setup)

**Phase 3 (User Story 1)**:
- Backend tasks T022-T023 can run in parallel
- Frontend tasks T027-T032 can run in parallel (all are independent components/composables)

**Phase 4 (User Story 2)**:
- Backend tasks T037-T039 can run in parallel
- Frontend tasks T043-T046 can run in parallel (all are independent components/composables)

**Phase 5 (User Story 3)**:
- Backend tasks T052-T053 can run in parallel
- Frontend tasks T056-T057 can run in parallel

**Phase 6 (Final Polish)**:
- All testing tasks T064-T071 can run in parallel
- All polish tasks T060-T063 can run in parallel

---

## Implementation Strategy

### Recommended Approach

**Week 1: Foundation**
1. Complete Phase 1 (Setup) - T001-T010
2. Complete Phase 2 (Foundational) - T011-T021
3. Verify backend server starts and frontend dev server runs

**Week 2: MVP (User Story 1)**
1. Complete Phase 3 backend - T022-T026
2. Complete Phase 3 frontend - T027-T036
3. Test end-to-end: view episodes, pagination, show notes modal

**Week 3: Task Management (User Stories 2 & 3)**
1. Complete Phase 4 backend - T037-T042
2. Complete Phase 4 frontend - T043-T051
3. Complete Phase 5 backend - T052-T055
4. Complete Phase 5 frontend - T056-T059
5. Test end-to-end: create tasks, monitor progress

**Week 4: Polish & Testing**
1. Complete Phase 6 - T060-T074
2. Run full test suite (unit + e2e)
3. Fix bugs and optimize performance
4. Deploy and document

### MVP Scope

**Minimum Viable Product = Phase 1 + Phase 2 + Phase 3 (User Story 1)**
- Total MVP tasks: T001-T036 (36 tasks)
- Delivers: Browse downloaded podcasts with pagination and show notes viewing
- Can be deployed and used independently

**Full Feature Set = All Phases**
- Total tasks: T001-T074 (74 tasks)
- Delivers: Complete podcast management with download task creation and progress monitoring

