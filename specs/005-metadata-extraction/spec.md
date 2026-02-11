# Feature Specification: Podcast Metadata Extraction and Display

**Feature Branch**: `005-metadata-extraction`
**Created**: 2026-02-09
**Status**: Draft
**Input**: User description: "抓取播客阶段，增加 .metadata.json 文件，从页面中提取播客时长和发布时间，格式分别为"231分钟"，"2个月前"，存放在 class 名称包含 info 的 DOM 元素下。然后将这些信息展示在播客列表中。"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Extract and Store Podcast Metadata (Priority: P1)

When a user downloads a podcast episode, the system automatically extracts additional metadata (duration and publish time) from the podcast page and saves it alongside the audio file in a structured format.

**Why this priority**: This is the core functionality that enables richer podcast information to be captured at download time. Without this, the display story cannot function.

**Independent Test**: Can be fully tested by downloading a podcast episode and verifying that a `.metadata.json` file is created in the podcast directory with the correct duration and publish time values.

**Acceptance Scenarios**:

1. **Given** a podcast episode page URL with duration and publish time (e.g., "231分钟", "2个月前") in a DOM element with class containing "info", **When** the user initiates a download, **Then** a `.metadata.json` file is created containing the extracted duration and publish time in their original formats
2. **Given** a podcast episode without accessible metadata, **When** the download completes, **Then** the `.metadata.json` file is still created with null or empty values for missing fields
3. **Given** an existing podcast download, **When** the user re-downloads with overwrite enabled, **Then** the `.metadata.json` file is updated with fresh metadata from the page

---

### User Story 2 - Display Metadata in Podcast List (Priority: P2)

Users can view the duration and publish time of each podcast episode directly in the podcast list interface, making it easier to browse and select episodes.

**Why this priority**: This provides immediate user value by surfacing the extracted metadata in a readable format, improving the browsing experience.

**Independent Test**: Can be fully tested by opening the podcast list view and verifying that each episode displays its duration and relative publish time (e.g., "231分钟", "刚刚发布", "2个月前").

**Acceptance Scenarios**:

1. **Given** multiple podcast episodes with metadata files, **When** the user views the podcast list, **Then** each episode displays its duration and publish time alongside the episode title
2. **Given** a podcast episode without metadata, **When** displayed in the list, **Then** the metadata fields show placeholder text (e.g., "--") or are hidden
3. **Given** a long podcast list, **When** the user scrolls through episodes, **Then** the metadata loads efficiently without significant performance degradation

---

### Edge Cases

- **Metadata extraction fails after audio download succeeds**: System continues with download completion, saves `.metadata.json` with null values for missing fields, and logs a warning
- What happens when the podcast page structure changes and the info class element is not found?
- How does the system handle podcasts with extremely long durations (e.g., "3000分钟")?
- What happens when the publish time format is unexpected (not in relative format like "2个月前")?
- How does the system handle network failures during metadata extraction (after audio download succeeds)?
- What happens when the `.metadata.json` file is corrupted or manually edited with invalid JSON?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST extract podcast duration from DOM elements with class names containing "info"
- **FR-002**: System MUST extract podcast publish time from DOM elements with class names containing "info"
- **FR-003**: System MUST create a `.metadata.json` file in the podcast's download directory
- **FR-004**: The `.metadata.json` file MUST contain duration in the original text format as displayed on the page (e.g., "231分钟", "1小时15分钟")
- **FR-005**: The `.metadata.json` file MUST contain publish time in the original text format as displayed on the page (e.g., "刚刚发布", "1小时前", "昨天", "2个月前")
- **FR-006**: System MUST display duration and publish time in the podcast list interface
- **FR-007**: System MUST handle missing metadata gracefully by creating `.metadata.json` with null values and logging a warning, without failing the audio download
- **FR-008**: System MUST overwrite existing `.metadata.json` when re-downloading with overwrite flag

### Key Entities *(include if feature involves data)*

- **PodcastMetadata**: Represents the extracted metadata for a podcast episode
  - Duration: The length of the podcast in original text format as displayed on the page (e.g., "231分钟", "1小时15分钟")
  - PublishTime: The relative time since publication in original text format (e.g., "刚刚发布", "1小时前", "昨天", "2个月前")
  - EpisodeTitle: The title of the podcast episode
  - PodcastName: The name of the podcast series
- **DownloadSession**: Represents the download process including metadata extraction
  - AudioFile: The downloaded audio file path
  - MetadataFile: The path to the `.metadata.json` file
  - ExtractionStatus: Success/failure status of metadata extraction

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of successfully downloaded podcasts have an associated `.metadata.json` file
- **SC-002**: Metadata extraction succeeds for 95%+ of podcast pages with standard structure
- **SC-003**: Users can identify podcast duration and recency without opening individual episode pages
- **SC-004**: Podcast list displays metadata for all episodes with available metadata within 1 second of page load
- **SC-005**: Missing metadata does not cause errors or crashes in the podcast list display

## Clarifications

### Session 2026-02-09

- Q: How should the system handle metadata extraction failure when audio download succeeds? → A: Continue download, save metadata file with null values, log warning

## Assumptions

- The podcast page structure from Xiaoyuzhou FM remains consistent with class names containing "info" for metadata
- Duration and publish time are displayed as text content within the identified DOM elements
- The metadata extraction happens during the download phase, not as a separate background process
- The podcast list interface has sufficient space to display duration and publish time alongside existing information
- The system captures metadata in whatever format is displayed on the page (e.g., duration as "231分钟" or "1小时15分钟", publish time as "刚刚发布", "1小时前", "昨天", or "2个月前")
