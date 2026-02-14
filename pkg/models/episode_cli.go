package models

import (
	"fmt"
	"regexp"
	"strings"
)

// Episode represents a podcast episode with metadata.
type Episode struct {
	// ID is the unique episode identifier
	ID string

	// Title is the episode title (may be empty if not found)
	Title string

	// AudioURL is the direct URL to the .m4a audio file
	AudioURL string

	// PageURL is the original episode page URL
	PageURL string

	// FileSize is the audio file size in bytes (-1 if unknown)
	FileSize int64

	// Duration is the episode duration in seconds (0 if unknown)
	Duration int
}

// SanitizedTitle returns a filesystem-safe version of the title.
// Falls back to episode ID if title is empty after sanitization.
func (e *Episode) SanitizedTitle() string {
	if e.Title == "" {
		return e.ID
	}

	// Remove invalid filesystem characters: < > : " / \ | ? *
	reg := regexp.MustCompile(`[<>:"/\\|?*]`)
	sanitized := reg.ReplaceAllString(e.Title, "_")

	// Limit to 200 characters
	if len(sanitized) > 200 {
		sanitized = sanitized[:200]
	}

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	// Fallback to ID if empty after sanitization
	if sanitized == "" {
		return e.ID
	}

	return sanitized
}

// GenerateFilename generates a filename in the format:
// {sanitized_title}_{episode_id}.m4a
func (e *Episode) GenerateFilename() string {
	sanitizedTitle := e.SanitizedTitle()
	return fmt.Sprintf("%s_%s.m4a", sanitizedTitle, e.ID)
}

// Validate checks if all required fields are populated and valid.
func (e *Episode) Validate() error {
	if e.ID == "" {
		return fmt.Errorf("episode ID is required")
	}

	if e.PageURL == "" {
		return fmt.Errorf("page URL is required")
	}

	if e.AudioURL == "" {
		return fmt.Errorf("audio URL is required")
	}

	// Remove query parameters and fragment before checking extension
	audioURL := e.AudioURL
	if idx := strings.Index(audioURL, "?"); idx != -1 {
		audioURL = audioURL[:idx]
	}
	if idx := strings.Index(audioURL, "#"); idx != -1 {
		audioURL = audioURL[:idx]
	}

	// Check if AudioURL ends with .m4a
	if !strings.HasSuffix(audioURL, ".m4a") {
		return fmt.Errorf("audio URL must be .m4a format")
	}

	return nil
}
