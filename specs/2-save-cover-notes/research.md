# Research: Save Cover Images and Show Notes

**Feature**: 2-save-cover-notes
**Date**: 2026-02-08
**Status**: Complete

## Overview

This document captures technical research and design decisions for enhancing the podcast downloader CLI to automatically download cover images and show notes alongside audio files.

## Technical Decisions

### 1. HTML Element Selection Strategy for Show Notes

**Decision**: Multi-fallback approach using aria-label attributes and semantic selectors.

**Implementation Strategy**:
1. **Primary**: Search for `<section aria-label="节目show notes">` (exact match for Xiaoyuzhou FM)
2. **Secondary**: Search for any element with `aria-label` containing "show notes" (case-insensitive)
3. **Tertiary**: Use semantic HTML selectors (`<article>`, `<section>` with description metadata, `<div class="description">`)
4. **Fallback**: Look for large text blocks near episode title/content area
5. **Failure**: Log warning with specific reason if all strategies fail

**Rationale**:
- The spec explicitly identifies `<section aria-label="节目show notes">` as the primary target
- Multi-fallback approach ensures robustness against page structure changes
- aria-label is a semantic attribute specifically designed for accessibility, making it a stable selector
- Graceful degradation aligns with constitutional requirement for external integration resilience

**Alternatives Considered**:
- **CSS Class selectors** (e.g., `.show-notes`, `.episode-description`): Rejected because class names are more likely to change during website redesigns
- **Fixed path selectors** (e.g., `div > div > section`): Rejected due to fragility - any HTML structure change breaks the selector
- **Machine learning-based extraction**: Rejected as overkill for this use case; adds complexity and dependencies

**Best Practices**:
- Prefer semantic attributes (aria-label, role) over CSS classes
- Always provide fallback strategies for web scraping
- Log specific failure reasons to aid debugging
- Test against multiple podcast pages to verify robustness

### 2. Cover Image Extraction Strategy

**Decision**: Extract cover image URL directly from `.avatar-container` class selector.

**Implementation Strategy**:
1. Find element with class `avatar-container`
2. Extract first `<img>` child (episode cover)
3. Return the image URL
4. If not found, return error (unexpected since element is guaranteed to exist)
5. Download failures should display warning but continue with audio download

**Implementation Details**:
```go
// avatar-container selector (Xiaoyuzhou FM specific)
if selection := doc.Find(".avatar-container img").First(); selection.Length() > 0 {
    if src, exists := selection.Attr("src"); exists && src != "" {
        return src, nil
    }
}
return "", ErrCoverNotFound
```

**Important Note**: The `.avatar-container` element contains **two** images:
- First `<img>`: Episode cover (single episode artwork) ✅ **This is what we want**
- Second `<img>`: Podcast account cover (channel/series artwork) ❌ **Skip this**

**Rationale**:
- The `.avatar-container` class selector is specific to Xiaoyuzhou FM's HTML structure
- Cover images are guaranteed to exist on all podcast episodes (no need for fallback)
- Selecting the **first** image ensures we get the episode-specific cover, not the podcast account cover
- No fallback needed simplifies the code and improves maintainability

**Best Practices**:
- Validate image URL before downloading (check for http/https prefix)
- Select **only the first image** from `.avatar-container img` to avoid the podcast account cover
- Handle both absolute and relative URLs
- Display warning but continue with audio download if cover download fails (graceful degradation)

### 3. Image Download and Format Handling

**Decision**: Preserve the original image format as provided by the website, with JPEG preference when multiple formats available.

**Implementation Strategy**:
1. Download image in the format returned by the image URL (preserve Content-Type)
2. Determine file extension from URL or Content-Type header
3. If multiple image URLs available (e.g., from different meta tags), prefer JPEG
4. Save with extension matching the format: `.jpg`, `.png`, `.webp`, or `.gif`
5. Validate downloaded image using magic byte detection (not just file extension)

**Rationale**:
- Preserving original format avoids unnecessary transcoding and quality loss
- JPEG is the most widely supported format and offers excellent compression for photographs
- Most podcast cover art is JPEG, so this matches real-world usage
- File extension validation is important for Windows file associations

**Alternatives Considered**:
- **Convert all to JPEG**: Rejected - adds complexity (requires image processing library), potential quality loss from PNG→JPEG
- **Convert all to PNG**: Rejected - much larger file sizes, slower downloads
- **WebP only**: Rejected - limited support in older media players and OS versions

**Best Practices**:
- Use magic byte detection to validate image files (check for JPEG: `FF D8 FF`, PNG: `89 50 4E 47`)
- Set reasonable file size limits (e.g., 10MB max for cover images)
- Implement the same retry logic as audio downloads
- Display download progress for large images (>1MB)

### 4. Show Notes Text Processing and Formatting

**Decision**: Convert HTML to plain text while preserving structural elements (links, lists, timestamps).

**Implementation Strategy**:
1. Extract HTML content from the identified show notes element
2. Convert HTML to plain text using goquery's `.Text()` method
3. Preserve links by converting `<a href="...">text</a>` to `text (URL: ...)` format
4. Preserve list structure by converting `<ul>/<ol>` to bullet/numbered lists
5. Preserve timestamps and headers by detecting common patterns
6. Save as UTF-8 with BOM encoding

**Rationale**:
- Plain text is universally readable across all platforms and text editors
- Preserved links allow users to access referenced resources
- UTF-8 with BOM ensures proper display of Chinese characters and emojis on Windows
- Consistent with the spec requirement for readable, structured format

**Alternatives Considered**:
- **Save as HTML**: Rejected - requires browser to view, less portable, potential security issues
- **Markdown format**: Rejected - adds complexity, less widely supported than plain text
- **JSON format**: Rejected - not human-readable, overkill for this use case

**Best Practices**:
- Strip excessive whitespace (more than 2 consecutive newlines)
- Limit line length to ~80 characters for readability
- Handle nested HTML structures correctly (e.g., lists within lists)
- Preserve blockquotes as quoted text with `>` prefix
- Test with show notes containing emojis, Chinese characters, and special symbols

### 5. File Naming and Organization

**Decision**: Use consistent base filename across all asset types with appropriate extensions.

**Implementation Strategy**:
1. Extract episode title from HTML (same as existing logic)
2. Sanitize filename (remove invalid characters, limit length)
3. Use same base filename for all assets:
   - Audio: `{episode-title}.m4a` (existing)
   - Cover: `{episode-title}.jpg` (or .png, .webp based on format)
   - Show notes: `{episode-title}.txt`
4. Save all files in the same output directory (co-located)
5. Add optional CLI flag to organize into subdirectories by podcast name

**Rationale**:
- Consistent naming makes it easy to identify related files
- Co-located files simplify episode management and archiving
- Sorting alphabetically groups related files together
- Minimal disruption to existing user workflow

**Alternatives Considered**:
- **Separate directories for each asset type**: Rejected - makes it harder to see all episode assets at a glance
- **Timestamps in filenames**: Rejected - cluttered, harder to read
- **Show notes with `-show-notes` suffix**: Rejected - unnecessary since `.txt` extension is already distinctive

**Best Practices**:
- Sanitize filenames consistently across platforms (Windows, macOS, Linux)
- Limit filename length to avoid filesystem issues (max 255 characters)
- Preserve Chinese characters in filenames (UTF-8 encoding)
- Test with episodes having very long titles (truncate intelligently at word boundaries)

### 6. Error Handling and User Communication

**Decision**: Detailed multi-line warning messages with specific failure reasons; audio download continues regardless.

**Implementation Strategy**:
1. Try to download cover image first (fastest operation)
2. Try to extract and save show notes
3. Download audio file (primary operation)
4. For any failure in cover/show notes:
   - Log detailed warning: `"Warning: [asset type] download failed: [reason]. Audio download completed successfully."`
   - Include specific reason (e.g., "404 Not Found", "element not found", "network timeout")
5. Use structured logging with error codes for programmatic parsing
6. Return appropriate exit codes (0 for success, non-zero only if audio fails)

**Rationale**:
- Clear, actionable error messages improve user experience
- Confirming audio success reduces user anxiety when warnings appear
- Specific reasons help with debugging and reporting issues
- Graceful degradation aligns with spec requirements (FR-011, FR-012)

**Alternatives Considered**:
- **Silent failures**: Rejected - users deserve to know what was/wasn't downloaded
- **Stop on any failure**: Rejected - overkill for non-essential assets; audio is the primary goal
- **Ask user for confirmation**: Rejected - breaks automation, inconvenient for batch downloads

**Best Practices**:
- Use color coding for warnings (yellow) to distinguish from errors (red)
- Include troubleshooting hints for common failures
- Log full error details to a file for debugging (optional `--verbose` flag)
- Rate limit warnings to avoid spam when downloading multiple episodes

### 7. Character Encoding for Show Notes

**Decision**: UTF-8 with BOM (Byte Order Mark) encoding.

**Implementation Strategy**:
1. Write show notes files with UTF-8 encoding
2. Prepend UTF-8 BOM (`0xEF 0xBB 0xBF`) at the start of the file
3. Use Go's `utf8` package for encoding validation
4. Test with various character sets (Chinese, emojis, accented characters)

**Rationale**:
- UTF-8 is the universal standard for text encoding
- BOM ensures Windows applications (Notepad, etc.) correctly identify the encoding
- Unix/Linux systems ignore the BOM, so no compatibility issues
- Spec requirement explicitly states "UTF-8 with BOM for Windows compatibility"

**Alternatives Considered**:
- **UTF-8 without BOM**: Rejected - Windows Notepad and some older tools misidentify the encoding
- **UTF-16**: Rejected - larger file size, not universally supported on all platforms
- **System default encoding**: Rejected - causes corruption when moving files between systems

**Best Practices**:
- Always validate UTF-8 encoding before writing
- Handle replacement characters for invalid sequences
- Test show notes with Chinese, Japanese, Korean, and emoji characters
- Document the encoding choice in the README/user guide

### 8. Progress Display for Cover Downloads

**Decision**: Reuse existing progress bar infrastructure for cover image downloads.

**Implementation Strategy**:
1. Use the same `progressbar/v3` library as audio downloads
2. Show progress for images larger than 1MB (suppress for smaller images to avoid flicker)
3. Use a single progress bar when downloading sequentially (audio → cover → show notes)
4. Add label to progress bar indicating current asset (e.g., "Downloading audio...", "Downloading cover...")

**Rationale**:
- Consistent user experience across all download types
- Large images can take several seconds on slow connections
- Reusing existing code reduces complexity and duplication
- Spec requirement FR-010 explicitly asks for progress feedback

**Alternatives Considered**:
- **No progress for images**: Rejected - violates FR-010, poor UX for large images
- **Separate progress bar for each asset**: Rejected - visual clutter, harder to read
- **Spinner instead of progress bar**: Rejected - less informative for large files

**Best Practices**:
- Detect content length from HTTP headers to enable accurate progress
- Handle cases where content length is unknown (show indeterminate progress)
- Clear the progress line after completion (use `\r` carriage return)
- Suppress progress output when redirecting to a file (`--quiet` mode)

### 9. Extension Strategy for URLExtractor Interface

**Decision**: Extend the existing `URLExtractor` interface to return cover image URL and show notes content.

**Implementation Strategy**:
1. Extend `URLExtractor.ExtractURL()` method signature (breaking change):
   - Current: `(audioURL string, title string, error)`
   - New: `(audioURL string, coverURL string, showNotes string, title string, error)`
2. Alternatively, create a new `EpisodeMetadata` struct to return all values:
   ```go
   type EpisodeMetadata struct {
       AudioURL   string
       CoverURL   string
       ShowNotes  string
       Title      string
   }
   ```
3. Update `HTMLExtractor` implementation accordingly
4. Maintain backward compatibility by providing a wrapper function if needed

**Rationale**:
- Single HTML parsing pass is more efficient than multiple requests
- Consistent with single-responsibility principle (extractor does all extraction)
- Returning a struct is cleaner than multiple return values and allows future extensions

**Alternatives Considered**:
- **Separate extractor interfaces**: Rejected - would require multiple HTTP requests to the same page
- **Store intermediate state in struct**: Rejected - adds complexity, harder to test
- **Separate CLI commands**: Rejected - poor UX, users want single command to download everything

**Best Practices**:
- Use Go's error wrapping to provide context: `fmt.Errorf("extract show notes: %w", err)`
- Cache the HTML document to avoid re-parsing for each extraction type
- Write unit tests for each extraction method independently
- Document the interface change in migration notes

### 10. Testing Strategy

**Decision**: Table-driven unit tests for extraction logic, integration tests for full download workflow.

**Implementation Strategy**:
1. **Unit tests**:
   - Test cover image extraction with various HTML structures
   - Test show notes extraction with multiple fallback strategies
   - Test image format detection and validation
   - Test text encoding handling (UTF-8, Chinese characters, emojis)
   - Test filename sanitization and edge cases

2. **Integration tests**:
   - Test full download workflow with real podcast URLs
   - Test error scenarios (404, network failures, missing elements)
   - Test CLI flag combinations and output directory handling
   - Test concurrent downloads (if supported)

3. **Test fixtures**:
   - Sample HTML files for different podcast page structures
   - Sample image files (JPEG, PNG, WebP) for validation tests
   - Sample show notes with various character encodings

**Rationale**:
- Table-driven tests are idiomatic Go and make it easy to add test cases
- Integration tests ensure all components work together
- Test fixtures ensure consistent test data and avoid external dependencies

**Best Practices**:
- Use mocking for HTTP calls in unit tests (avoid real network requests)
- Test both success and failure paths
- Aim for >80% code coverage on new code
- Add benchmarks for performance-critical code (e.g., HTML parsing)

## Dependencies Analysis

### New Dependencies Required
**None** - All functionality can be implemented using existing dependencies:
- `goquery` (v1.8.1): HTML parsing and element selection
- Existing HTTP client: Downloading images and handling retries
- Standard library: File I/O, text processing, encoding

### Existing Dependencies to Extend Usage
- **goquery**: Enhanced usage for aria-label selectors and text extraction
- **cli/v2**: No changes needed (CLI flags already support output directory)
- **progressbar/v3**: Reuse for cover image download progress

## Migration Strategy

### Backward Compatibility
- Existing CLI usage remains unchanged: `./podcast-downloader <url>`
- New functionality is automatic - no new CLI flags required
- Users can opt-out with flags if needed (future enhancement)

### Code Changes Required
1. Extend `URLExtractor` interface and `HTMLExtractor` implementation
2. Create new services: `ImageDownloader`, `ShowNotesSaver`
3. Update main download workflow to call new services
4. Add error handling and warning messages
5. Update tests to cover new functionality

### Rollout Plan
1. Implement new services in isolation (unit tests)
2. Integrate with existing downloader (integration tests)
3. Manual testing with real podcast URLs
4. Update documentation (README, usage examples)
5. Release as patch version (backward compatible enhancement)

## Performance Considerations

### Expected Performance Impact
- **Cover image download**: +2-5 seconds per episode (depending on image size)
- **Show notes extraction**: +0.5-1 second per episode (HTML parsing is fast)
- **Total overhead**: +3-6 seconds per episode (acceptable for CLI tool)

### Optimization Opportunities
- Parallel downloads: Download cover and show notes concurrently with audio (future enhancement)
- Caching: Cache extracted metadata to avoid re-parsing on retries
- Compression: Compress show notes files if very large (future enhancement)

## Security Considerations

### Input Validation
- Validate image URLs to prevent SSRF (Server-Side Request Forgery)
- Limit file size to prevent disk space exhaustion
- Sanitize filenames to prevent path traversal attacks

### Output Safety
- Use magic byte validation to ensure downloaded files are actual images
- Escape special characters in show notes to prevent markup injection
- Set appropriate file permissions (0644 for downloaded files)

## Open Questions Resolved

All questions from the spec's Technical Context section have been resolved through this research:
- ✅ HTML element selection strategy (multi-fallback approach)
- ✅ Cover image extraction and format handling (preserve original, prefer JPEG)
- ✅ Show notes text processing (HTML to plain text with structure preservation)
- ✅ File naming and organization (consistent base filenames)
- ✅ Error handling and user communication (detailed warnings, graceful degradation)
- ✅ Character encoding (UTF-8 with BOM)
- ✅ Progress display (reuse progressbar/v3)
- ✅ Interface extension strategy (return struct with all metadata)
- ✅ Testing strategy (table-driven unit tests + integration tests)

## Next Steps

Proceed to **Phase 1**: Design & Contracts
1. Generate data-model.md with entities and relationships
2. Create contracts/ (N/A for CLI tool, but document interfaces)
3. Generate quickstart.md with implementation guide
4. Update agent context files with new technology decisions
