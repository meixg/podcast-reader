# Implementation Plan: Docker Container Packaging

**Branch**: `006-docker-packaging` | **Date**: 2026-02-14 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/006-docker-packaging/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Create a production-ready Docker image for the podcast reader application with automated CI/CD via GitHub Actions. The image will use Alpine Linux for minimal size (<100MB), support multi-architecture (amd64/arm64), expose health checks, and publish to GitHub Container Registry on releases.

## Technical Context

**Language/Version**: Go 1.21+ (existing project standard)
**Primary Dependencies**: Docker, GitHub Actions, docker/build-push-action
**Storage**: File-based (downloads directory mounted as volume)
**Testing**: Docker build verification, container runtime tests
**Target Platform**: Linux containers (amd64, arm64)
**Project Type**: Single container web service
**Performance Goals**: Image size <100MB, startup <10 seconds
**Constraints**: Alpine Linux base, multi-arch support, GHCR publishing
**Scale/Scope**: Single-instance deployment, personal use to small team

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Go Backend Standards | ✅ PASS | Using existing Go 1.21+ codebase |
| II. Service-Oriented Architecture | ✅ PASS | Container packaging doesn't change architecture |
| III. Asynchronous Processing First | ✅ PASS | Existing async behavior preserved |
| IV. External Integration Resilience | ✅ PASS | No new external integrations |
| V. Web API First Design | ✅ PASS | Health check endpoint added per spec |

**Gate Result**: ✅ ALL CHECKS PASSED - Proceeding to Phase 0

## Project Structure

### Documentation (this feature)

```text
specs/006-docker-packaging/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
# Repository root
├── .github/
│   └── workflows/
│       └── docker.yml       # GitHub Actions CI/CD workflow
├── Dockerfile               # Multi-stage Docker build
├── .dockerignore            # Exclude files from Docker context
├── docker-compose.yml       # Optional: local development compose
├── cmd/                     # Go application entry points
│   ├── server/              # Web server main package
│   │   └── main.go
│   ├── downloader/          # CLI downloader
│   └── inspect/             # CLI inspector
├── internal/                # Internal Go packages
│   ├── config/
│   └── taskmanager/
├── frontend/                # Existing Vue.js application
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
├── go.mod                   # Go module definition (root)
├── go.sum                   # Go dependencies
└── downloads/               # Volume mount for persisted data
```

**Structure Decision**: Docker packaging adds containerization artifacts at repository root. The Go application is built from root (where go.mod resides), not from a `backend/` subdirectory. New files: `Dockerfile`, `.dockerignore`, `.github/workflows/docker.yml`, and optionally `docker-compose.yml` for local development.

## Phase 0: Research & Decisions

### Research Areas

1. **Multi-arch Docker builds with GitHub Actions**
   - Decision: Use `docker/build-push-action` with `platforms: linux/amd64,linux/arm64`
   - Rationale: Native GitHub Actions support, handles QEMU emulation automatically
   - Alternatives considered: Manual manifest creation (complex), separate workflows per arch (inefficient)

2. **Go application Docker best practices**
   - Decision: Multi-stage build with `golang:alpine` builder and `alpine` runtime
   - Rationale: Smallest image size while maintaining compatibility
   - Alternatives considered: Distroless (harder to debug), Debian slim (larger)

3. **GitHub Container Registry authentication**
   - Decision: Use `GITHUB_TOKEN` for automatic authentication
   - Rationale: No secrets management needed, automatic permissions
   - Alternatives considered: Personal access token (requires rotation), Docker Hub (separate credentials)

4. **Health check implementation**
   - Decision: Add `/health` endpoint to Go backend
   - Rationale: Constitution requires health check endpoints for monitoring
   - Alternatives considered: TCP check only (insufficient), external health checker (overkill)

### Open Questions Resolved

None - all clarifications completed in `/speckit.clarify` phase.

---

## Phase 1: Design & Contracts

### Data Model

No new data entities - this feature packages existing application.

### API Contracts

New endpoint added per FR-007:

```yaml
# /contracts/health.yaml
paths:
  /health:
    get:
      summary: Health check endpoint
      description: Returns 200 OK when application is healthy
      responses:
        '200':
          description: Service is healthy
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: "healthy"
                  timestamp:
                    type: string
                    format: date-time
```

### Configuration Schema

Environment variables (FR-003):

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Web server port |
| `DOWNLOAD_DIR` | `/app/downloads` | Podcast download directory |
| `LOG_LEVEL` | `info` | Logging level |

### Docker Image Specification

```yaml
image: ghcr.io/{owner}/podcast-reader:{tag}
tags:
  - latest        # main branch builds
  - {version}     # release tags (e.g., v1.2.3)
  - pr-{number}   # pull request builds
platforms:
  - linux/amd64
  - linux/arm64
ports:
  - 8080
volumes:
  - /app/downloads  # Persist downloaded content
healthcheck:
  test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
  interval: 30s
  timeout: 3s
  retries: 3
```

---

## Quick Start

### Build locally

```bash
docker build -t podcast-reader:latest .
```

### Run container

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/downloads:/app/downloads \
  -e PORT=8080 \
  --name podcast-reader \
  ghcr.io/{owner}/podcast-reader:latest
```

### Docker Compose

```yaml
version: '3.8'
services:
  podcast-reader:
    image: ghcr.io/{owner}/podcast-reader:latest
    ports:
      - "8080:8080"
    volumes:
      - ./downloads:/app/downloads
    environment:
      - PORT=8080
      - LOG_LEVEL=info
    restart: unless-stopped
```

---

## Generated Artifacts

| Artifact | Path | Description |
|----------|------|-------------|
| research.md | `specs/006-docker-packaging/research.md` | Phase 0 research findings |
| data-model.md | `specs/006-docker-packaging/data-model.md` | Data entities (none new) |
| contracts/ | `specs/006-docker-packaging/contracts/` | API contracts |
| quickstart.md | `specs/006-docker-packaging/quickstart.md` | Usage guide |

---

## Next Steps

Run `/speckit.tasks` to generate the implementation task list.
