# Feature Specification: Podcast Management Web Application

**Feature Branch**: `004-frontend-web-app`
**Created**: 2026-02-08
**Status**: Draft
**Input**: User description: "实现一个前端应用，提供以下页面：1. 播客列表页，列出所有已经下载的播客列表，列表中展示除了 shownotes 之外的 meta 信息，点击后在弹窗中展示 shownotes；2. 播客任务页，展示所有的任务列表，列表中展示任务信息和状态，页面中有个新建 button，点击后出现一个弹窗可以创建一个新的下载任务（输入url）。"

## Clarifications

### Session 2026-02-08

- Q: How should users switch between the podcast list page and the download tasks page? → A: Navigation bar with links/tabs for both pages (always visible)
- Q: How should the podcast list display large numbers of episodes? → A: Pagination with configurable page size (e.g., 20/50/100 per page)
- Q: How should the system handle attempts to create duplicate download tasks for the same URL? → A: Prevent duplicates and show existing task (with link/navigation to it)
- Q: How frequently should the application check for task status updates? → A: Poll every 2-3 seconds (balanced approach)
- Q: How should the application behave when the backend server is unavailable? → A: Show error message with retry button (user-initiated retry)

## User Scenarios & Testing *(mandatory)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.
  
  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - View Downloaded Podcasts (Priority: P1)

Users need to browse their downloaded podcast collection to find and listen to episodes they've already downloaded.

**Why this priority**: This is the core value proposition - users must be able to access their downloaded content. Without this, the application has no purpose.

**Independent Test**: Can be fully tested by downloading sample podcasts via CLI, then opening the web app and verifying all episodes appear with correct metadata. Delivers immediate value by providing visibility into the podcast library.

**Acceptance Scenarios**:

1. **Given** the user has downloaded 5 podcast episodes, **When** they open the podcast list page, **Then** all 5 episodes are displayed with title, podcast name, duration, file size, and download date
2. **Given** the user is viewing the podcast list, **When** they click on an episode, **Then** a modal dialog opens displaying the full show notes for that episode
3. **Given** the show notes modal is open, **When** the user clicks outside the modal or presses the close button, **Then** the modal closes and returns to the list view
4. **Given** no podcasts have been downloaded, **When** the user opens the podcast list page, **Then** an empty state message is displayed indicating no podcasts are available

---

### User Story 2 - Create Download Tasks (Priority: P2)

Users need to initiate new podcast downloads directly from the web interface without using the command line.

**Why this priority**: This enables self-service downloads and makes the application accessible to non-technical users. It's the primary action users will take to grow their library.

**Independent Test**: Can be tested by opening the tasks page, clicking the create button, entering a valid podcast URL, and verifying the task appears in the task list. Delivers value by enabling users to download new content.

**Acceptance Scenarios**:

1. **Given** the user is on the tasks page, **When** they click the "New Download" button, **Then** a modal dialog opens with a URL input field and submit button
2. **Given** the create task modal is open, **When** the user enters a valid podcast URL and clicks submit, **Then** a new download task is created and appears in the task list
3. **Given** the create task modal is open, **When** the user enters an invalid URL, **Then** an error message is displayed indicating the URL format is incorrect
4. **Given** a download task has been created, **When** the user returns to the tasks page, **Then** the new task appears in the list with "pending" or "in progress" status

---

### User Story 3 - Monitor Download Progress (Priority: P3)

Users need to track the status of their download tasks to know when episodes are ready to listen.

**Why this priority**: This provides transparency and helps users understand system activity. While important for user experience, the core functionality works without real-time monitoring.

**Independent Test**: Can be tested by creating multiple download tasks and verifying each task displays its current status (pending, downloading, completed, failed). Delivers value by keeping users informed.

**Acceptance Scenarios**:

1. **Given** the user has created multiple download tasks, **When** they view the tasks page, **Then** each task displays its current status (pending, downloading, completed, or failed)
2. **Given** a task is in "downloading" status, **When** the user views the task, **Then** the task shows progress information (percentage complete or current step)
3. **Given** a task has completed successfully, **When** the user views the task, **Then** the task shows "completed" status with completion timestamp
4. **Given** a task has failed, **When** the user views the task, **Then** the task shows "failed" status with an error message explaining what went wrong

---

### Edge Cases

- What happens when a podcast file is corrupted or missing from disk but still appears in the list?
- How does the system handle extremely long show notes (10,000+ characters) in the modal display?
- Duplicate downloads: System prevents duplicate tasks for the same URL and shows existing task with navigation link
- How does the system handle network interruptions during task creation?
- Backend unavailability: System displays error message with retry button for user-initiated reconnection
- How does the system handle special characters or emojis in podcast titles and show notes?
- Large library handling: System uses pagination with configurable page size (20/50/100 per page) to maintain performance with hundreds or thousands of episodes

## Requirements *(mandatory)*

### Functional Requirements

**Navigation**:

- **FR-001**: System MUST provide a navigation bar that is always visible across all pages
- **FR-002**: Navigation bar MUST include links/tabs for both the Podcast List page and Download Tasks page
- **FR-003**: System MUST indicate which page is currently active in the navigation bar

**Podcast List Page**:

- **FR-004**: System MUST display a list of all downloaded podcast episodes
- **FR-005**: System MUST show episode metadata including title, podcast name, duration, file size, and download date for each episode
- **FR-006**: System MUST exclude show notes from the list view display
- **FR-007**: System MUST implement pagination for the podcast list
- **FR-008**: System MUST allow users to configure page size (e.g., 20, 50, or 100 episodes per page)
- **FR-009**: System MUST provide pagination controls (previous, next, page numbers)
- **FR-010**: System MUST open a modal dialog when a user clicks on an episode in the list
- **FR-011**: System MUST display the full show notes content in the modal dialog
- **FR-012**: System MUST allow users to close the show notes modal by clicking outside the modal or using a close button
- **FR-013**: System MUST display an appropriate empty state message when no podcasts have been downloaded

**Download Tasks Page**:

- **FR-014**: System MUST display a list of all download tasks with their current status
- **FR-015**: System MUST show task information including URL, status (pending/downloading/completed/failed), and timestamps
- **FR-016**: System MUST provide a "New Download" button on the tasks page
- **FR-017**: System MUST open a modal dialog when the "New Download" button is clicked
- **FR-018**: System MUST provide a URL input field in the create task modal
- **FR-019**: System MUST validate the URL format before submitting a new download task
- **FR-020**: System MUST check for duplicate URLs across all existing tasks (any status)
- **FR-021**: System MUST prevent creation of duplicate download tasks for the same URL
- **FR-022**: System MUST display a message indicating the task already exists when duplicate is detected
- **FR-023**: System MUST provide a link or navigation to the existing task when duplicate is detected
- **FR-024**: System MUST display error messages for invalid URL formats
- **FR-025**: System MUST create a new download task when a valid, non-duplicate URL is submitted
- **FR-026**: System MUST update the task list to show newly created tasks

**Error Handling**:

- **FR-027**: System MUST detect when the backend server is unavailable
- **FR-028**: System MUST display a clear error message when backend is unavailable
- **FR-029**: System MUST provide a retry button to allow users to manually retry the connection
- **FR-030**: System MUST attempt to reconnect to the backend when the retry button is clicked

### Key Entities

- **Podcast Episode**: Represents a downloaded podcast episode with attributes including title, podcast name, duration, file size, download date, show notes, and file location
- **Download Task**: Represents a download operation with attributes including source URL, current status, creation timestamp, completion timestamp, progress information, and error messages (if failed)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can view their complete podcast library within 2 seconds of opening the podcast list page
- **SC-002**: Users can access show notes for any episode with no more than 2 clicks (click episode, view modal)
- **SC-003**: Users can create a new download task in under 30 seconds from opening the tasks page
- **SC-004**: 95% of users successfully create their first download task without errors or confusion
- **SC-005**: Task status updates are visible to users within 5 seconds of status changes
- **SC-006**: The application remains responsive with up to 1000 downloaded episodes in the library
- **SC-007**: Modal dialogs open and close smoothly with no perceptible lag (under 300ms)

## Dependencies and Assumptions

### Dependencies

- **Backend API**: The web application depends on a backend service that provides endpoints for:
  - Retrieving the list of downloaded podcast episodes with metadata
  - Retrieving show notes for individual episodes
  - Creating new download tasks
  - Retrieving the list of download tasks with status information
- **Existing Download System**: The CLI-based podcast downloader must continue to function and store episodes in a format accessible to the web application

### Assumptions

- **Single User**: The application is designed for single-user local use (no multi-user authentication or authorization required)
- **Local Deployment**: The web application and backend service run on the user's local machine
- **File System Access**: The backend has read access to the downloads directory where podcast files are stored
- **URL Format**: Valid podcast URLs follow the Xiaoyuzhou FM format (e.g., "https://www.xiaoyuzhoufm.com/episode/...")
- **Network Availability**: Users have internet connectivity when creating new download tasks
- **Browser Compatibility**: Users access the application through modern web browsers (last 2 versions of major browsers)
- **Data Persistence**: Download task history persists across application restarts
- **Status Update Polling**: The application polls the backend every 2-3 seconds to retrieve task status updates
