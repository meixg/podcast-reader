# Feature Specification: Save Cover Images and Show Notes

**Feature Branch**: `2-save-cover-notes`
**Created**: 2026-02-08
**Status**: Draft
**Input**: User description: "在现在播客音频内容时，将播客的封面和 show notes 也同步保存下来"

## Clarifications

### Session 2026-02-08

- Q: When saving show notes, which filename convention should be used? → A: Use same base filename with `.txt` extension (e.g., `episode-title.txt`)
- Q: When multiple cover image formats are available, which format should the system prefer? → A: Use the same format as the website provides, prefer JPEG when multiple available
- Q: When extracting show notes from the podcast page, what should be the HTML element selection strategy? → A: Multi-fallback: Try aria-label="节目show notes" → aria-label containing "show notes" → semantic selectors → fail
- Q: When cover image or show notes download fails, what warning message format should be displayed? → A: Detailed multi-line: "Warning: [asset type] download failed: [reason]. Audio download completed successfully."
- Q: What character encoding should be used for saving show notes text files with non-ASCII characters? → A: UTF-8 with BOM (Byte Order Mark for Windows compatibility)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Download Cover Image with Audio (Priority: P1)

When downloading a podcast episode, the user wants the cover (album art) image to be automatically saved alongside the audio file, so they can view the episode artwork in their media player and maintain a complete podcast collection.

**Why this priority**: Cover images are essential metadata for media organization and user experience. Most media players display cover art, and having it saved locally provides a complete offline experience. This is the highest priority because it's a visual element that users immediately notice when missing.

**Independent Test**: Can be fully tested by running the downloader with a podcast URL and verifying that a cover image file (JPEG/PNG) is created in the same directory as the audio file with a matching filename.

**Acceptance Scenarios**:

1. **Given** a valid podcast episode URL, **When** the download completes, **Then** a cover image file should be saved in the same directory as the audio file
2. **Given** a podcast with available cover art, **When** the cover image is downloaded, **Then** the image should be in a standard format (JPEG or PNG)
3. **Given** a downloaded episode, **When** viewing the download directory, **Then** the cover image filename should match the audio filename (e.g., `episode-name.m4a` and `episode-name.jpg`)
4. **Given** the `.avatar-container` element is guaranteed to exist on all podcast pages, **When** extracting the cover image URL, **Then** the system successfully extracts the first `<img>` element's src attribute (no fallback needed)

---

### User Story 2 - Download Show Notes with Audio (Priority: P2)

When downloading a podcast episode, the user wants the show notes (episode description, guest information, links, timestamps) to be automatically saved as a text file, so they can reference the content, visit mentioned links, and read detailed episode information while listening offline.

**Why this priority**: Show notes contain valuable metadata including guest names, topics, resource links, and timestamp references. While not as immediately visible as cover art, they provide important reference material that enhances the podcast listening experience. This is P2 because users can listen without show notes, but having them improves content discovery and follow-up research.

**Independent Test**: Can be fully tested by running the downloader with a podcast URL and verifying that a text file containing show notes is created in the same directory as the audio file, with properly formatted content including episode title, description, links, and any timestamp data.

**Acceptance Scenarios**:

1. **Given** a valid podcast episode URL with show notes, **When** the download completes, **Then** a text file containing show notes should be saved in the same directory as the audio file
2. **Given** show notes containing structured data (links, timestamps, guest lists), **When** saved to the text file, **Then** the content should be formatted in a readable, structured way
3. **Given** a downloaded episode, **When** viewing the show notes file, **Then** the filename should match the audio filename with `.txt` extension (e.g., `episode-name.m4a` and `episode-name.txt`)
4. **Given** a podcast without available show notes, **When** downloading the episode, **Then** the download should complete successfully with a warning message about missing show notes

---

### User Story 3 - Organize All Downloaded Assets (Priority: P3)

When downloading multiple podcast episodes, the user wants all related files (audio, cover image, show notes) to be organized together, so they can easily manage their podcast library and move or archive complete episodes without losing associated metadata.

**Why this priority**: File organization is important for long-term library management but is less critical than the actual downloading of cover art and show notes. Users can manually organize files if needed, but automatic organization provides better user experience. This is P3 because the feature works without it, but organization enhances usability.

**Independent Test**: Can be fully tested by downloading multiple episodes and verifying that each episode's related files are co-located and clearly associated, either through matching filenames or organized directory structure.

**Acceptance Scenarios**:

1. **Given** multiple downloaded episodes, **When** viewing the download directory, **Then** each episode's audio, cover, and show notes files should be clearly associated through matching base filenames
2. **Given** a user moving an episode to archive, **When** they select the audio file, **Then** it should be easy to identify and move the associated cover and show notes files
3. **Given** downloaded files with consistent naming, **When** sorting files alphabetically in the directory, **Then** related files should appear adjacent to each other

---

### Edge Cases

- The `.avatar-container` element always contains two images: first is episode cover (saved), second is podcast account cover (ignored)
- How does the system handle show notes with extremely long content (e.g., 10,000+ character transcripts)?
- What happens when cover image URLs are invalid or return 404 errors after extraction? (Answer: warning displayed, audio download continues)
- How does the system handle special characters or non-ASCII characters in show notes when saving to text files?
- What happens when the cover image download succeeds but the show notes extraction fails? (Answer: audio and cover are saved, warning shown for show notes)
- How does the system handle show notes with complex HTML formatting (tables, lists, embedded media)?
- What happens when the user specifies a custom output directory - should all files still be co-located?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: When downloading a podcast episode, the system MUST extract the cover image URL from the `.avatar-container` element's first `<img>` child (episode cover)
- **FR-002**: The system MUST download the cover image and save it to the same directory as the audio file
- **FR-003**: The system MUST save cover images in the same format as provided by the website (JPEG, PNG, or WEBP), preferring JPEG when multiple formats are available
- **FR-004**: The system MUST name the cover image file using the same base filename as the audio file (e.g., `episode-title.jpg` for `episode-title.m4a`)
- **FR-005**: When downloading a podcast episode, the system MUST extract show notes content from the podcast page using a multi-fallback strategy: (1) search for `<section aria-label="节目show notes">`, (2) search for any element with aria-label containing "show notes", (3) use semantic HTML selectors (article, section with description metadata), and (4) log failure if no show notes found
- **FR-006**: The system MUST save show notes as a text file in the same directory as the audio file
- **FR-007**: The system MUST name the show notes file using the same base filename as the audio file with `.txt` extension (e.g., `episode-title.txt` for `episode-title.m4a`)
- **FR-008**: The system MUST convert HTML-formatted show notes into readable plain text format and save the file using UTF-8 encoding with BOM (Byte Order Mark) to ensure proper display of non-ASCII characters (e.g., Chinese, emojis) across all platforms
- **FR-009**: The system MUST preserve important structural elements in show notes (links, lists, timestamps, section headers)
- **FR-010**: The system MUST display download progress feedback for cover images (especially for large image files)
- **FR-011**: If cover image download fails (network error, invalid URL, etc.), the system MUST display a detailed warning message: "Warning: Cover image download failed: [specific reason]. Audio download completed successfully." and continue with the audio download (the `.avatar-container` element is guaranteed to exist, but network issues may still occur)
- **FR-012**: If show notes extraction fails, the system MUST display a detailed multi-line warning message: "Warning: Show notes extraction failed: [specific reason]. Audio download completed successfully." and continue with the audio download
- **FR-013**: The system MUST validate that downloaded cover images are valid image files (not HTML error pages or corrupted data)
- **FR-014**: The system MUST handle network errors during cover image download with retry logic (same as audio downloads)
- **FR-015**: The system MUST support the existing CLI flags for file naming and output directories for cover and show notes files

### Key Entities

- **Podcast Episode**: Represents a single podcast episode with associated metadata including audio file URL, cover image URL, and show notes content
- **Cover Image**: Visual artwork associated with a podcast episode or podcast series, typically in JPEG or PNG format
- **Show Notes**: Textual metadata describing the episode content, including episode title, description, guest information, topics, timestamp references, and related links
- **Download Session**: A complete download operation that includes audio file, cover image, and show notes text file

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can successfully download cover images for 95% of podcast episodes that have available cover art
- **SC-002**: Users can successfully download show notes for 95% of podcast episodes that have available show notes content
- **SC-003**: Cover image downloads complete within 5 seconds for standard-size images (under 2MB)
- **SC-004**: Show notes are preserved with 100% text accuracy (no character encoding corruption or loss)
- **SC-005**: All related files for an episode (audio, cover, show notes) use consistent base filenames for easy association
- **SC-006**: Download failures for cover images or show notes do not prevent audio download from completing (graceful degradation)
- **SC-007**: Users can identify and manage all files associated with a podcast episode without confusion about file relationships
- **SC-008**: The system handles missing cover art or show notes gracefully with clear warning messages, not silent failures
