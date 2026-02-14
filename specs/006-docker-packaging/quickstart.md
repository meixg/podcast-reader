# Quick Start: Docker Container Packaging

**Date**: 2026-02-14
**Feature**: Docker Container Packaging

## Prerequisites

- Docker 20.10+ installed
- (Optional) Docker Compose 2.0+ for compose-based deployment

## Local Build

### 1. Build the Docker Image

```bash
# From repository root
docker build -t podcast-reader:latest .
```

### 2. Run the Container

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/downloads:/app/downloads \
  -e PORT=8080 \
  -e LOG_LEVEL=info \
  --name podcast-reader \
  podcast-reader:latest
```

### 3. Verify Health

```bash
# Check health endpoint
curl http://localhost:8080/health

# Expected response:
# {"status":"healthy","timestamp":"2026-02-14T10:30:00Z"}
```

### 4. View Logs

```bash
docker logs -f podcast-reader
```

### 5. Stop and Remove

```bash
docker stop podcast-reader
docker rm podcast-reader
```

---

## Using Pre-built Images (GitHub Container Registry)

### 1. Pull the Image

```bash
docker pull ghcr.io/{owner}/podcast-reader:latest
```

### 2. Run

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/downloads:/app/downloads \
  --name podcast-reader \
  ghcr.io/{owner}/podcast-reader:latest
```

---

## Docker Compose (Recommended)

### 1. Create `docker-compose.yml`

```yaml
version: '3.8'

services:
  podcast-reader:
    image: ghcr.io/{owner}/podcast-reader:latest
    container_name: podcast-reader
    ports:
      - "8080:8080"
    volumes:
      - ./downloads:/app/downloads
    environment:
      - PORT=8080
      - DOWNLOAD_DIR=/app/downloads
      - LOG_LEVEL=info
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 10s
```

### 2. Start Services

```bash
docker-compose up -d
```

### 3. View Logs

```bash
docker-compose logs -f
```

### 4. Stop Services

```bash
docker-compose down
```

---

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Web server port inside container |
| `DOWNLOAD_DIR` | `/app/downloads` | Directory for downloaded podcasts |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |

### Volume Mounts

| Container Path | Purpose |
|----------------|---------|
| `/app/downloads` | Persist downloaded podcast files |

---

## Troubleshooting

### Container fails to start

```bash
# Check logs
docker logs podcast-reader

# Check if port is already in use
lsof -i :8080
```

### Permission denied on volume

```bash
# Ensure downloads directory exists and is writable
mkdir -p downloads
chmod 755 downloads
```

### Health check fails

```bash
# Check if application is responding
docker exec podcast-reader wget -q --spider http://localhost:8080/health
echo $?
# Should return 0 (success)
```

---

## Multi-Architecture Support

The published image supports both AMD64 and ARM64 architectures:

```bash
# Docker automatically pulls the correct architecture
docker pull ghcr.io/{owner}/podcast-reader:latest

# Verify architecture
docker inspect ghcr.io/{owner}/podcast-reader:latest | grep Architecture
```

---

## Updating

### Pull Latest Image

```bash
docker pull ghcr.io/{owner}/podcast-reader:latest
```

### Recreate Container (data is preserved in volume)

```bash
docker stop podcast-reader
docker rm podcast-reader
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/downloads:/app/downloads \
  --name podcast-reader \
  ghcr.io/{owner}/podcast-reader:latest
```

### With Docker Compose

```bash
docker-compose pull
docker-compose up -d
```
