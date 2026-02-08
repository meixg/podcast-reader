# Data Model: Podcast Management Web Application

**Feature**: 004-frontend-web-app
**Date**: 2026-02-08

## Overview

This document defines the data structures for the podcast management web application, including entities, relationships, and validation rules.

## Core Entities

### Episode

Represents a downloaded podcast episode with metadata.

**Attributes**:
- `id` (string): Unique identifier (file path or hash)
- `title` (string): Episode title
- `podcastName` (string): Name of the podcast series
- `duration` (string): Episode duration (e.g., "45:30")
- `fileSize` (number): File size in bytes
- `downloadDate` (string): ISO 8601 timestamp
- `showNotes` (string): Full episode description/show notes
- `filePath` (string): Absolute path to audio file
- `coverImagePath` (string, optional): Path to cover image

**Validation Rules**:
- `id`: Required, non-empty
- `title`: Required, max 500 characters
- `podcastName`: Required, max 200 characters
- `duration`: Required, format "HH:MM:SS" or "MM:SS"
- `fileSize`: Required, positive integer
- `downloadDate`: Required, valid ISO 8601 timestamp
- `showNotes`: Optional, max 50,000 characters
- `filePath`: Required, valid file path

**State Transitions**: None (immutable after download)

### DownloadTask

Represents a download operation with status tracking.

**Attributes**:
- `id` (string): Unique task identifier (UUID)
- `url` (string): Source podcast episode URL
- `status` (enum): Current task status
  - `pending`: Task created, not started
  - `downloading`: Download in progress
  - `completed`: Download finished successfully
  - `failed`: Download failed with error
- `createdAt` (string): ISO 8601 timestamp when task was created
- `completedAt` (string, optional): ISO 8601 timestamp when task finished
- `progress` (number, optional): Download progress percentage (0-100)
- `errorMessage` (string, optional): Error description if status is `failed`
- `episodeId` (string, optional): ID of resulting episode if completed

**Validation Rules**:
- `id`: Required, valid UUID format
- `url`: Required, valid URL format, must match Xiaoyuzhou FM pattern
- `status`: Required, one of: pending, downloading, completed, failed
- `createdAt`: Required, valid ISO 8601 timestamp
- `completedAt`: Optional, valid ISO 8601 timestamp, only if status is completed/failed
- `progress`: Optional, integer 0-100, only if status is downloading
- `errorMessage`: Optional, max 1000 characters, only if status is failed
- `episodeId`: Optional, non-empty, only if status is completed

**State Transitions**:
```
pending → downloading → completed
                     → failed
```

## Relationships

### Episode ← DownloadTask

- **Type**: One-to-One (optional)
- **Description**: A completed DownloadTask may reference the resulting Episode
- **Cardinality**:
  - One DownloadTask can create zero or one Episode
  - One Episode is created by exactly one DownloadTask
- **Implementation**: DownloadTask.episodeId references Episode.id

## Pagination Models

### PaginatedEpisodes

Response model for paginated episode list.

**Attributes**:
- `episodes` (Episode[]): Array of episodes for current page
- `total` (number): Total number of episodes across all pages
- `page` (number): Current page number (1-indexed)
- `pageSize` (number): Number of items per page
- `totalPages` (number): Total number of pages

**Validation Rules**:
- `episodes`: Required, array of valid Episode objects
- `total`: Required, non-negative integer
- `page`: Required, positive integer, <= totalPages
- `pageSize`: Required, one of: 20, 50, 100
- `totalPages`: Required, non-negative integer, calculated as ceil(total / pageSize)

## Error Models

### APIError

Standard error response format for all API endpoints.

**Attributes**:
- `error` (string): Human-readable error message
- `code` (string): Machine-readable error code
- `details` (object, optional): Additional error context

**Error Codes**:
- `INVALID_URL`: URL format validation failed
- `DUPLICATE_TASK`: Task already exists for this URL
- `NOT_FOUND`: Requested resource not found
- `SERVER_ERROR`: Internal server error

**Example**:
```json
{
  "error": "A download task for this URL already exists",
  "code": "DUPLICATE_TASK",
  "details": {
    "existingTaskId": "abc-123",
    "url": "https://www.xiaoyuzhoufm.com/episode/..."
  }
}
```

## TypeScript Type Definitions

### Frontend Types

```typescript
// types/episode.ts
export interface Episode {
  id: string;
  title: string;
  podcastName: string;
  duration: string;
  fileSize: number;
  downloadDate: string;
  showNotes: string;
  filePath: string;
  coverImagePath?: string;
}

export interface PaginatedEpisodes {
  episodes: Episode[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}
```

```typescript
// types/task.ts
export type TaskStatus = 'pending' | 'downloading' | 'completed' | 'failed';

export interface DownloadTask {
  id: string;
  url: string;
  status: TaskStatus;
  createdAt: string;
  completedAt?: string;
  progress?: number;
  errorMessage?: string;
  episodeId?: string;
}

export interface CreateTaskRequest {
  url: string;
}

export interface APIError {
  error: string;
  code: string;
  details?: Record<string, any>;
}
```

### Go Struct Definitions

```go
// internal/models/episode.go
package models

import "time"

type Episode struct {
    ID              string    `json:"id"`
    Title           string    `json:"title"`
    PodcastName     string    `json:"podcastName"`
    Duration        string    `json:"duration"`
    FileSize        int64     `json:"fileSize"`
    DownloadDate    time.Time `json:"downloadDate"`
    ShowNotes       string    `json:"showNotes"`
    FilePath        string    `json:"filePath"`
    CoverImagePath  string    `json:"coverImagePath,omitempty"`
}

type PaginatedEpisodes struct {
    Episodes   []Episode `json:"episodes"`
    Total      int       `json:"total"`
    Page       int       `json:"page"`
    PageSize   int       `json:"pageSize"`
    TotalPages int       `json:"totalPages"`
}
```

```go
// internal/models/task.go
package models

import "time"

type TaskStatus string

const (
    TaskStatusPending     TaskStatus = "pending"
    TaskStatusDownloading TaskStatus = "downloading"
    TaskStatusCompleted   TaskStatus = "completed"
    TaskStatusFailed      TaskStatus = "failed"
)

type DownloadTask struct {
    ID           string     `json:"id"`
    URL          string     `json:"url"`
    Status       TaskStatus `json:"status"`
    CreatedAt    time.Time  `json:"createdAt"`
    CompletedAt  *time.Time `json:"completedAt,omitempty"`
    Progress     *int       `json:"progress,omitempty"`
    ErrorMessage string     `json:"errorMessage,omitempty"`
    EpisodeID    string     `json:"episodeId,omitempty"`
}

type CreateTaskRequest struct {
    URL string `json:"url"`
}

type APIError struct {
    Error   string                 `json:"error"`
    Code    string                 `json:"code"`
    Details map[string]interface{} `json:"details,omitempty"`
}
```
