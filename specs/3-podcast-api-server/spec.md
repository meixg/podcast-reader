# Feature Specification: Podcast API Server

**Feature Branch**: `3-podcast-api-server`
**Created**: 2026-02-08
**Status**: Draft
**Input**: User description: "在现有 podcast-downloader 的基础上，实现一个 http server，提供以下 API：1. 提交一个播客的下载任务，输入是一个 xiaoyuzhou 的 url；2. 过去所有已下载播客的列表，包括 title, audioUrl, shownotes, coverImage。"

## Clarifications

### Session 2026-02-08

- Q: When the server restarts (crash, deployment, system reboot), what should happen to in-progress download tasks and their tracking data? → A: Completed downloads remain discoverable; in-progress tasks are lost and must be resubmitted after restart
- Q: What should the system do when Xiaoyuzhou FM is temporarily unavailable or returns errors during metadata extraction or download? → A: Retry failed requests up to 3 times with exponential backoff before marking task as failed; failed tasks allow resubmission of the same URL
- Q: When the server starts or when listing downloaded podcasts, how should the system discover and track completed downloads? → A: Scan the downloads directory structure on server startup to rebuild the in-memory catalog of completed downloads

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Submit Podcast Download Task (Priority: P1)

A user wants to download a podcast episode from Xiaoyuzhou FM by submitting the episode URL through a web service. Instead of using the CLI tool directly, they send the URL to an HTTP endpoint which processes the download asynchronously in the background.

**Why this priority**: This is the core functionality that enables programmatic access to podcast downloading, allowing integration with other applications and services.

**Independent Test**: Can be fully tested by sending a POST request with a Xiaoyuzhou FM URL to the download endpoint and verifying that the podcast audio, cover image, and show notes are downloaded and stored.

**Acceptance Scenarios**:

1. **Given** a valid Xiaoyuzhou FM episode URL, **When** user submits POST request to the download endpoint, **Then** system returns a task ID and begins downloading the podcast in the background
2. **Given** an invalid or malformed URL, **When** user submits POST request, **Then** system returns appropriate error response with validation details
3. **Given** a URL that has already been downloaded, **When** user submits POST request, **Then** system returns existing download information without re-downloading
4. **Given** a URL with a download task in progress, **When** user submits POST request, **Then** system returns the existing task ID and current status without creating a new task
5. **Given** a download in progress, **When** user checks task status, **Then** system returns current progress information

---

### User Story 2 - List Downloaded Podcasts (Priority: P2)

A user wants to retrieve a list of all previously downloaded podcast episodes with their metadata. They send a GET request to retrieve information about title, audio file location, show notes, and cover image for each downloaded episode.

**Why this priority**: This provides visibility into the podcast library and enables users to browse and access downloaded content programmatically.

**Independent Test**: Can be fully tested by sending a GET request to the list endpoint and verifying that all downloaded podcasts are returned with complete metadata including title, audio URL, show notes, and cover image paths.

**Acceptance Scenarios**:

1. **Given** multiple downloaded podcasts exist, **When** user sends GET request to list endpoint, **Then** system returns complete list with all metadata fields for each podcast
2. **Given** no podcasts have been downloaded, **When** user sends GET request to list endpoint, **Then** system returns empty list
3. **Given** podcasts with partial data (e.g., missing cover or show notes), **When** user retrieves list, **Then** system includes available fields and indicates missing data
4. **Given** a large number of downloaded podcasts, **When** user retrieves list, **Then** system returns results efficiently within acceptable time limits

---

### User Story 3 - Query Download Task Status (Priority: P3)

A user wants to check the status of a submitted download task to see if it has completed, is in progress, or has failed.

**Why this priority**: This provides visibility into asynchronous operations and enables better user experience for long-running downloads.

**Independent Test**: Can be fully tested by submitting a download task, then querying the task status endpoint to verify accurate status reporting.

**Acceptance Scenarios**:

1. **Given** a submitted download task, **When** user queries task status, **Then** system returns current status (pending, in-progress, completed, failed)
2. **Given** a completed download task, **When** user queries task status, **Then** system returns success status with paths to downloaded files
3. **Given** a failed download task, **When** user queries task status, **Then** system returns error status with failure reason
4. **Given** a non-existent task ID, **When** user queries status, **Then** system returns 404 error response

---

### Edge Cases

- What happens when the Xiaoyuzhou FM website is down or returns errors? (Retry up to 3 times with exponential backoff, then mark task as failed; user can resubmit)
- What happens when the download is interrupted by network failure? (Download task fails, user can resubmit)
- What happens when the server restarts? (Completed downloads remain discoverable; in-progress tasks are lost, users must resubmit)
- What happens when disk space is insufficient for the download?
- How does the system handle concurrent download requests for the same URL? (Returns existing task/download information)
- What happens when the downloaded file is corrupted or invalid?
- How does the system handle URLs from sources other than Xiaoyuzhou FM?
- What happens when show notes or cover images are not available for an episode?
- How does the system identify and track previously downloaded URLs to prevent duplicates? (Scan downloads directory on startup to build in-memory catalog mapping URLs to downloaded files)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide an HTTP API endpoint that accepts a Xiaoyuzhou FM URL and initiates a podcast download task
- **FR-002**: System MUST return a unique task identifier for each submitted download request
- **FR-003**: System MUST process download requests asynchronously, returning the task ID immediately without waiting for download completion
- **FR-004**: System MUST download the podcast audio file, cover image, and show notes for each valid URL
- **FR-005**: System MUST provide an HTTP API endpoint that returns a list of all previously downloaded podcasts
- **FR-006**: System MUST include title, audio file path, show notes content, and cover image path for each podcast in the list
- **FR-007**: System MUST validate that submitted URLs are from Xiaoyuzhou FM (supporting both xiaoyuzhou.fm and xiaoyuzhoufm.com domains) and reject URLs from other sources
- **FR-008**: System MUST persist download task information including status, timestamps, and file locations
- **FR-009**: System MUST provide appropriate HTTP status codes and error messages for invalid requests
- **FR-010**: System MUST handle concurrent download requests without data corruption
- **FR-011**: System MUST store downloaded files in organized directories by podcast title
- **FR-012**: System MUST provide an endpoint to query the status of a download task by task ID
- **FR-013**: System MUST support filtering and pagination for the podcast list when the number of downloaded episodes is large
- **FR-014**: System MUST detect duplicate download requests by URL and return existing task or download information without creating a new download task
- **FR-015**: System MUST include task status and file location information when returning existing download information for duplicate URL submissions
- **FR-016**: System MUST retry failed requests to Xiaoyuzhou FM up to 3 times with exponential backoff before marking the task as failed
- **FR-017**: System MUST allow resubmission of failed download tasks using the same URL
- **FR-018**: System MUST scan the downloads directory on server startup to rebuild the in-memory catalog of completed downloads

### Key Entities

- **Download Task**: Represents a podcast download request with attributes including unique task ID, source URL, submission timestamp, current status (pending, in-progress, completed, failed), completion timestamp, error message (if failed), and references to downloaded files

- **Podcast Episode**: Represents a downloaded podcast with attributes including title, audio file path, cover image path, show notes content/text file path, download timestamp, and source URL

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can submit a download request and receive a task ID within 500 milliseconds
- **SC-002**: Download tasks complete successfully for 95% of valid Xiaoyuzhou FM URLs
- **SC-003**: Users can retrieve the complete list of downloaded podcasts within 2 seconds for up to 1000 episodes
- **SC-004**: System can handle at least 10 concurrent download requests without performance degradation
- **SC-005**: Task status queries return current status within 100 milliseconds
- **SC-006**: Downloaded files include audio, cover image, and show notes for 90% of episodes where these assets are available on the source
- **SC-007**: System provides clear error messages for 100% of failed requests with specific reasons (invalid URL, network error, insufficient disk space, etc.)
- **SC-008**: API is available and responsive 99% of the time when the server is running

## Assumptions

1. The HTTP server will be run on a local machine or private network, not exposed publicly to the internet
2. Downloaded files will be stored in the existing `downloads/` directory structure used by the CLI tool
3. The system will reuse existing download logic from the podcast-downloader CLI tool
4. Download tasks do not require authentication or authorization (assumed trusted environment)
5. The server will run as a single-instance application (not distributed/multi-server)
6. Show notes are stored as plain text files using the same format as the CLI tool
7. Cover images are stored in their original format (JPEG, PNG, or WebP)
8. The API will use JSON for request and response bodies
9. There is no requirement for real-time progress updates during download (status polling is sufficient)
10. No user interface is required - this is an API-only service
11. The system will maintain a record of downloaded URLs and their corresponding file locations to detect and prevent duplicate downloads
12. Download task state (in-progress tasks) is stored in-memory only and not persisted across server restarts; completed downloads are discoverable by scanning the downloads directory on server startup to rebuild the in-memory catalog
13. The system will use the same retry logic as the CLI tool (up to 3 retries with exponential backoff) when encountering errors from Xiaoyuzhou FM

## Out of Scope

- User authentication and authorization
- Real-time websocket or SSE progress updates
- Podcast discovery or search functionality
- Audio transcoding or format conversion
- Podcast playback functionality
- Web-based user interface for the API
- Integration with podcast sources other than Xiaoyuzhou FM
- Cloud storage integration (local filesystem only)
- Database persistence (in-memory or file-based storage is acceptable)
- Multi-user support or download queues per user
- Scheduled or automatic podcast downloads
