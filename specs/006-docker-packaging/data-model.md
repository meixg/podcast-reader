# Data Model: Docker Container Packaging

**Date**: 2026-02-14
**Feature**: Docker Container Packaging

## Overview

This feature does not introduce new data entities. It packages the existing podcast reader application into a Docker container. The data model remains unchanged from the existing application.

## Existing Entities (Unchanged)

The following entities continue to exist as defined in the application:

- **Episode**: Podcast episode metadata (title, URL, description, etc.)
- **DownloadSession**: Tracks active download operations
- **DownloadedFile**: Represents downloaded audio files with metadata

## Configuration Entities

### ContainerConfiguration

Runtime configuration passed via environment variables:

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `PORT` | integer | 8080 | Web server listening port |
| `DOWNLOAD_DIR` | string | "/app/downloads" | Directory for downloaded content |
| `LOG_LEVEL` | string | "info" | Logging verbosity (debug, info, warn, error) |

### DockerImage

Build-time metadata (not stored, embedded in image):

| Field | Type | Description |
|-------|------|-------------|
| `version` | string | Application version (from git tag) |
| `build_date` | string | ISO 8601 timestamp of build |
| `git_commit` | string | Short SHA of source commit |
| `platform` | string | Target architecture (amd64/arm64) |

## Volume Mounts

### Downloads Volume

| Property | Value |
|----------|-------|
| Container path | `/app/downloads` |
| Purpose | Persist downloaded podcast files |
| Required | Yes (for data durability) |
| Permissions | Read/Write |

## Notes

- No database schema changes
- No new API data models
- Configuration is environment-based, not persisted
- All state remains file-based as per existing architecture
