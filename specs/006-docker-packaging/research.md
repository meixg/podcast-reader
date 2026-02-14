# Research Findings: Docker Container Packaging

**Date**: 2026-02-14
**Feature**: Docker Container Packaging

## Research Area 1: Multi-arch Docker Builds with GitHub Actions

### Decision
Use `docker/build-push-action` with `platforms: linux/amd64,linux/arm64`

### Rationale
- Native GitHub Actions support with official Docker actions
- Automatic QEMU emulation setup for cross-platform builds
- Caching support for faster subsequent builds
- Integrated with GitHub Container Registry authentication

### Alternatives Considered

| Alternative | Pros | Cons |
|-------------|------|------|
| Manual manifest creation | Full control | Complex, error-prone, requires multiple CI steps |
| Separate workflows per arch | Parallel builds | Duplicated configuration, harder to maintain |
| Docker Buildx directly | Flexible | More complex setup, less GitHub-native |

### Key Configuration
```yaml
- uses: docker/build-push-action@v5
  with:
    platforms: linux/amd64,linux/arm64
    push: true
    tags: ghcr.io/${{ github.repository }}:latest
```

---

## Research Area 2: Go Application Docker Best Practices

### Decision
Multi-stage build with `golang:1.21-alpine` builder and `alpine:3.19` runtime

### Rationale
- Multi-stage builds produce smaller final images
- Alpine Linux base is ~5MB, keeping total image under 100MB target
- Static binary compilation eliminates libc dependencies
- Go's built-in cross-compilation makes multi-arch straightforward

### Dockerfile Pattern
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates wget
WORKDIR /app
COPY --from=builder /app/server ./
EXPOSE 8080
CMD ["./server"]
```

### Alternatives Considered

| Alternative | Size | Debuggability | Notes |
|-------------|------|---------------|-------|
| Distroless | ~20MB | Hard | Google's minimal images, no shell |
| Debian Slim | ~80MB | Easy | Larger but more compatible |
| Scratch | ~15MB | Very Hard | No shell, no certs, minimal |

---

## Research Area 3: GitHub Container Registry Authentication

### Decision
Use `GITHUB_TOKEN` for automatic authentication

### Rationale
- No additional secrets to manage or rotate
- Automatic permissions based on repository access
- Works seamlessly with `docker/login-action`
- Supports both read and write operations

### Configuration
```yaml
- uses: docker/login-action@v3
  with:
    registry: ghcr.io
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}
```

### Permissions Required
```yaml
permissions:
  contents: read
  packages: write
```

### Alternatives Considered

| Alternative | Setup | Maintenance | Notes |
|-------------|-------|-------------|-------|
| Personal Access Token | Manual | Requires rotation | Broader permissions than needed |
| Docker Hub | Separate account | Separate credentials | Additional service to manage |
| AWS ECR | Complex IAM | Requires AWS setup | Overkill for this use case |

---

## Research Area 4: Health Check Implementation

### Decision
Add `/health` HTTP endpoint to Go backend

### Rationale
- Constitution requires health check endpoints for monitoring
- HTTP health checks provide better signal than TCP alone
- Can verify application is actually ready to serve requests
- Standard pattern for container orchestration (K8s, Docker Swarm)

### Implementation
```go
// Simple health handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status": "healthy",
        "timestamp": time.Now().UTC(),
    })
}
```

### Docker HEALTHCHECK
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```

### Alternatives Considered

| Alternative | Pros | Cons |
|-------------|------|------|
| TCP port check only | Simple | Doesn't verify app is actually ready |
| Process check | Very simple | Least informative |
| External health checker | Comprehensive | Additional complexity |

---

## Research Area 5: Docker Image Tagging Strategy

### Decision
- `latest` for main branch builds
- Semantic version (e.g., `v1.2.3`) for release tags
- `pr-{number}` for pull request builds

### Rationale
- Clear distinction between development and production images
- Semantic versioning follows Go module conventions
- PR tags enable testing without polluting main tags

### Implementation
```yaml
tags: |
  ghcr.io/${{ github.repository }}:latest
  ghcr.io/${{ github.repository }}:${{ github.ref_name }}
  ghcr.io/${{ github.repository }}:pr-${{ github.event.number }}
```

---

## Summary

All research areas have been resolved with industry-standard approaches:

1. **Multi-arch builds**: GitHub Actions native support
2. **Go Docker best practices**: Alpine-based multi-stage build
3. **Registry auth**: GITHUB_TOKEN for zero-config auth
4. **Health checks**: HTTP endpoint + Docker HEALTHCHECK
5. **Tagging**: Semantic versioning with environment-specific tags

No blocking issues identified. Ready for implementation.
