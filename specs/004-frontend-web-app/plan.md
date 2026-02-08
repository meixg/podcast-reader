# Implementation Plan: Podcast Management Web Application

**Branch**: `004-frontend-web-app` | **Date**: 2026-02-08 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/004-frontend-web-app/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build a web-based frontend application for managing podcast downloads with two main pages: (1) Podcast List page displaying downloaded episodes with metadata and show notes modal, (2) Download Tasks page for creating and monitoring download operations. The application uses Vue 3 + TypeScript + Vite + Tailwind CSS for the frontend, communicating with a Go backend API that manages podcast downloads and file storage.

## Technical Context

**Language/Version**:
- Frontend: TypeScript 5.x with Vue 3.4+
- Backend: Go 1.21+ (existing)

**Primary Dependencies**:
- Frontend: Vue 3, Vite 5.x, Tailwind CSS 3.x, TypeScript
- Backend: Go standard library (net/http), existing downloader packages

**Storage**: File-based (existing downloads directory, in-memory task queue)

**Testing**:
- Frontend: Vitest (unit), Cypress (e2e)
- Backend: Go testing package (existing)

**Target Platform**: Web browsers (last 2 versions), local deployment

**Project Type**: Web application (frontend + backend)

**Performance Goals**:
- Page load < 2 seconds (SC-001)
- Modal open < 300ms (SC-007)
- Status updates within 5 seconds (SC-005)
- Support 1000+ episodes (SC-006)

**Constraints**:
- Single-user local deployment
- No authentication required
- Poll backend every 2-3 seconds
- Pagination required (20/50/100 per page)

**Scale/Scope**:
- 2 main pages (Podcast List, Download Tasks)
- ~10-15 Vue components
- 4-6 backend API endpoints
- Support up to 1000 episodes

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### ✅ Go Backend Standards
- Using Go 1.21+ with existing downloader packages
- Will follow gofmt, Go conventions, comprehensive error handling
- Backend API will be RESTful with proper status codes

### ✅ Service-Oriented Architecture
- Clear separation: Vue.js frontend + Go backend
- Backend exposes RESTful APIs, frontend handles UI/UX
- File-based storage for downloads (existing pattern)

### ✅ Asynchronous Processing First
- Download tasks are asynchronous (existing pattern)
- Status endpoints for tracking (will be added)
- Task queue for managing downloads (existing)

### ✅ Web API First Design
- All functionality accessible via REST APIs
- JSON request/response format
- Proper HTTP status codes and CORS for frontend

### ✅ Vue.js Architecture
- Vue 3 with Composition API ✓
- TypeScript for type safety ✓
- Vite for build tooling ✓
- Component-based architecture ✓

### ✅ Styling
- Tailwind CSS for utility-first styling ✓
- Responsive design with mobile-first approach ✓

### ✅ State Management
- Vue 3 Composition API with reactive state ✓
- No external state management library (as per constitution) ✓
- Proper loading states and error handling ✓

### ✅ Build Quality
- ESLint + Prettier for code formatting ✓
- TypeScript strict mode ✓
- Vitest for unit tests, Cypress for e2e ✓

**Result**: ✅ All constitutional requirements satisfied. No violations to justify.

## Project Structure

### Documentation (this feature)

```text
specs/004-frontend-web-app/
├── spec.md              # Feature specification
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   └── api.yaml         # OpenAPI specification
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
backend/
├── cmd/
│   └── server/
│       └── main.go           # HTTP server entry point
├── internal/
│   ├── handlers/             # HTTP handlers for API endpoints
│   │   ├── episodes.go       # GET /api/episodes, GET /api/episodes/:id/shownotes
│   │   └── tasks.go          # GET /api/tasks, POST /api/tasks
│   ├── services/             # Business logic
│   │   ├── episode_service.go
│   │   └── task_service.go
│   └── models/               # Data structures
│       ├── episode.go
│       └── task.go
├── pkg/                      # Reusable packages
│   └── scanner/              # File system scanner for downloads
└── storage/                  # File-based storage
    └── downloads/            # Existing podcast downloads

frontend/
├── src/
│   ├── components/           # Vue components
│   │   ├── common/
│   │   │   ├── NavigationBar.vue
│   │   │   ├── Modal.vue
│   │   │   └── Pagination.vue
│   │   ├── episodes/
│   │   │   ├── EpisodeList.vue
│   │   │   ├── EpisodeCard.vue
│   │   │   └── ShowNotesModal.vue
│   │   └── tasks/
│   │       ├── TaskList.vue
│   │       ├── TaskCard.vue
│   │       └── CreateTaskModal.vue
│   ├── views/                # Page components
│   │   ├── EpisodesView.vue
│   │   └── TasksView.vue
│   ├── services/             # API client
│   │   └── api.ts
│   ├── composables/          # Vue 3 composables
│   │   ├── useEpisodes.ts
│   │   ├── useTasks.ts
│   │   └── usePagination.ts
│   ├── types/                # TypeScript types
│   │   ├── episode.ts
│   │   └── task.ts
│   ├── router/               # Vue Router
│   │   └── index.ts
│   ├── App.vue
│   └── main.ts
├── public/
├── tests/
│   ├── unit/                 # Vitest unit tests
│   └── e2e/                  # Cypress e2e tests
├── package.json
├── vite.config.ts
├── tailwind.config.js
└── tsconfig.json
```

**Structure Decision**: Web application structure with separate backend/ and frontend/ directories. Backend uses Go with standard project layout (cmd/, internal/, pkg/). Frontend uses Vue 3 with Vite, following component-based architecture with clear separation of concerns (components, views, services, composables).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
