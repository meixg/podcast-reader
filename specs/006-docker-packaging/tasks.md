# Tasks: Docker Container Packaging

**Feature**: Docker Container Packaging
**Branch**: `006-docker-packaging`
**Generated**: 2026-02-14
**Spec**: [spec.md](./spec.md) | **Plan**: [plan.md](./plan.md)

---

## Overview

### Implementation Strategy

**MVP First**: Complete User Story 1 (Build Docker Image) first to establish the core Docker infrastructure. This provides immediate value - a working Docker image that can be built and tested locally.

**Incremental Delivery**:
- Phase 1-2: Setup and foundation (health endpoint, Dockerfile)
- Phase 3: US1 - Build Docker Image (MVP - core functionality)
- Phase 4: US2 - Run Application in Container (runtime verification)
- Phase 5: US3 - Environment Configuration (flexibility)
- Phase 6: US4 - Data Persistence (volume mounts)
- Phase 7: US5 - CI/CD Automation (GitHub Actions)
- Phase 8: Polish and documentation

### User Story Dependencies

```
US1 (Build Docker Image) - FOUNDATION
    │
    ├── enables ──> US2 (Run in Container)
    │
    ├── enables ──> US3 (Environment Config)
    │
    ├── enables ──> US4 (Data Persistence)
    │
    └── enables ──> US5 (CI/CD Automation)

US2, US3, US4 can be developed in parallel after US1
US5 depends on US1 (needs working Dockerfile)
```

### Parallel Execution Opportunities

**Within US1**: T006 (Dockerfile) and T007 (.dockerignore) can be done in parallel
**Within US2**: T010 (volume test) and T011 (port mapping test) can be done in parallel
**Within US3**: T013 (PORT) and T014 (DOWNLOAD_DIR) can be done in parallel
**US2, US3, US4**: Can be implemented in parallel after US1 completes

---

## Phase 1: Project Setup

**Goal**: Initialize Docker packaging infrastructure
**Independent Test**: N/A (setup phase)

- [x] T001 Create `.github/workflows/` directory for CI/CD workflows
- [x] T002 Create `downloads/` directory with `.gitkeep` for volume mount testing

---

## Phase 2: Foundational - Health Endpoint

**Goal**: Implement health check endpoint required for container orchestration
**Independent Test**: `curl http://localhost:8080/health` returns `{"status":"healthy"}`

**Prerequisites**: None (adds to existing server)

- [x] T003 [P] Add health handler in `web/handlers/health.go`
- [x] T004 Register health route at `/health` in `cmd/server/main.go`
- [x] T005 Test health endpoint locally with `go build cmd/server/main.go`

---

## Phase 3: User Story 1 - Build Docker Image (P1)

**Goal**: Create a working Docker image that can be built locally
**Independent Test**: `docker build -t podcast-reader:test .` succeeds and `docker images` shows the image

**Prerequisites**: Phase 2 (health endpoint for HEALTHCHECK)

- [x] T006 Create `Dockerfile` with multi-stage Alpine build
- [x] T007 [P] Create `.dockerignore` to exclude unnecessary files
- [ ] T008 Build image locally: `docker build -t podcast-reader:test .` (requires Docker permissions)
- [ ] T009 Verify image size is under 100MB: `docker images podcast-reader:test`

---

## Phase 4: User Story 2 - Run Application in Container (P1)

**Goal**: Run the application inside a container and verify it works
**Independent Test**: Container starts and responds to HTTP requests on mapped port

**Prerequisites**: Phase 3 (Docker image exists)

- [ ] T010 [US2] Run container with volume mount: `docker run -d -p 8080:8080 -v $(pwd)/downloads:/app/downloads --name pr-test podcast-reader:test`
- [ ] T011 [P] [US2] Verify web interface loads at `http://localhost:8080`
- [ ] T012 [US2] Test download functionality persists to volume

---

## Phase 5: User Story 3 - Configure via Environment Variables (P2)

**Goal**: Make application configurable through environment variables
**Independent Test**: Starting container with `-e PORT=9090` changes the listening port

**Prerequisites**: Phase 3 (Docker image exists)

- [ ] T013 [P] [US3] Update server to read `PORT` from environment with default 8080
- [ ] T014 [P] [US3] Update server to read `DOWNLOAD_DIR` from environment with default `/app/downloads`
- [ ] T015 [US3] Test custom port: `docker run -e PORT=9090 -p 9090:9090 ...`
- [ ] T016 [US3] Test custom download directory via env var

---

## Phase 6: User Story 4 - Persist Data Outside Container (P2)

**Goal**: Ensure data persists across container restarts
**Independent Test**: Download content, remove container, start new one with same volume, data is still there

**Prerequisites**: Phase 3 (Docker image exists)

- [ ] T017 [US4] Verify downloads appear on host filesystem at mounted location
- [ ] T018 [US4] Test data persistence: stop container, start new one with same volume, verify content
- [ ] T019 [US4] Test container restart preserves data

---

## Phase 7: User Story 5 - Automated CI/CD Build (P2)

**Goal**: Automate Docker image builds with GitHub Actions
**Independent Test**: Push to main branch triggers workflow and image appears in GHCR

**Prerequisites**: Phase 3 (working Dockerfile)

- [x] T020 [US5] Create `.github/workflows/docker.yml` with build job
- [x] T021 [P] [US5] Configure multi-arch build (amd64, arm64) using docker/build-push-action
- [x] T022 [P] [US5] Configure GHCR authentication with GITHUB_TOKEN
- [x] T023 [US5] Add tagging strategy: `latest` for main, version tags for releases
- [ ] T024 [US5] Test workflow on PR (build only, don't push)
- [ ] T025 [US5] Verify image is published to GHCR on main branch push

---

## Phase 8: Polish & Cross-Cutting Concerns

**Goal**: Documentation, optimization, and final touches
**Independent Test**: New user can follow README and successfully run the container

- [x] T026 Create `docker-compose.yml` for easy local development
- [x] T027 Update root `README.md` with Docker usage instructions
- [x] T028 Add Docker HEALTHCHECK instruction to Dockerfile
- [ ] T029 Verify final image size is under 100MB (requires Docker build)
- [ ] T030 Test complete workflow: build → run → download → restart → verify data (requires Docker)

---

## Task Summary

| Phase | Tasks | Story | Priority |
|-------|-------|-------|----------|
| Phase 1: Setup | 2 | - | - |
| Phase 2: Foundation | 3 | - | - |
| Phase 3: US1 | 4 | Build Docker Image | P1 |
| Phase 4: US2 | 3 | Run in Container | P1 |
| Phase 5: US3 | 4 | Environment Config | P2 |
| Phase 6: US4 | 3 | Data Persistence | P2 |
| Phase 7: US5 | 6 | CI/CD Automation | P2 |
| Phase 8: Polish | 5 | - | - |
| **Total** | **30** | | |

### Parallel Tasks

- **T006 + T007**: Dockerfile and .dockerignore (no dependencies)
- **T010 + T011**: Volume mount and port mapping tests
- **T013 + T014**: PORT and DOWNLOAD_DIR env vars
- **T021 + T022**: Multi-arch config and GHCR auth

### Suggested MVP Scope

**Complete through Phase 3 (US1)** for MVP:
- Working Dockerfile
- Local build capability
- Image size under 100MB

This provides immediate value - users can build and run the container locally even before CI/CD is implemented.

---

## File Inventory

### New Files to Create

1. `.github/workflows/docker.yml` - GitHub Actions CI/CD workflow
2. `Dockerfile` - Multi-stage Docker build
3. `.dockerignore` - Docker context exclusions
4. `docker-compose.yml` - Local development compose file
5. `internal/handlers/health.go` - Health check handler

### Files to Modify

1. `cmd/server/main.go` - Register health endpoint, read env vars
2. `README.md` - Add Docker usage documentation

### Directories to Create

1. `.github/workflows/` - GitHub Actions workflows
