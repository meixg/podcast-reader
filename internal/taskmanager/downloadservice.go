package taskmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/meixg/podcast-reader/internal/downloader"
)

// DownloadService handles the complete download workflow for podcast episodes.
type DownloadService struct {
	extractor       downloader.URLExtractor
	fileDownloader  downloader.FileDownloader
	imageDownloader *downloader.HTTPImageDownloader
	showNotesSaver  *downloader.PlainTextShowNotesSaver
	outputDirectory string
	logger          *log.Logger
}

// NewDownloadService creates a new download service with all required components.
func NewDownloadService(outputDir string, logger *log.Logger) *DownloadService {
	// Create HTTP client with appropriate timeouts
	metadataClient := downloader.NewHTTPClient(30 * time.Second)
	downloadClient := &http.Client{Timeout: 1 * time.Hour}
	imageClient := &http.Client{Timeout: 2 * time.Minute}

	return &DownloadService{
		extractor:       downloader.NewHTMLExtractor(metadataClient),
		fileDownloader:  downloader.NewHTTPDownloader(downloadClient, false),
		imageDownloader: downloader.NewHTTPImageDownloader(imageClient, 10*1024*1024),
		showNotesSaver:  downloader.NewPlainTextShowNotesSaver(),
		outputDirectory: outputDir,
		logger:          logger,
	}
}

// DownloadResult contains the result of a download operation.
type DownloadResult struct {
	Success       bool
	AudioPath     string
	CoverPath     string
	ShowNotesPath string
	Title         string
	Error         error
}

// DownloadEpisode downloads a podcast episode from the given URL.
func (s *DownloadService) DownloadEpisode(ctx context.Context, url string, progressCallback func(int)) *DownloadResult {
	result := &DownloadResult{}

	// Step 1: Extract metadata
	s.logger.Printf("Extracting metadata from: %s", url)
	metadata, err := s.extractor.ExtractURL(ctx, url)
	if err != nil {
		result.Error = fmt.Errorf("failed to extract metadata: %w", err)
		result.Success = false
		return result
	}

	result.Title = metadata.Title

	// Step 2: Create podcast directory
	podcastTitle := sanitizeDirectoryName(metadata.Title)
	podcastDir := filepath.Join(s.outputDirectory, podcastTitle)

	if err := os.MkdirAll(podcastDir, 0755); err != nil {
		result.Error = fmt.Errorf("failed to create directory: %w", err)
		result.Success = false
		return result
	}

	// Step 3: Download audio file
	audioPath := filepath.Join(podcastDir, "podcast.m4a")
	s.logger.Printf("Downloading audio to: %s", audioPath)

	if progressCallback != nil {
		progressCallback(10)
	}

	bytesWritten, err := s.fileDownloader.Download(ctx, metadata.AudioURL, audioPath, nil)
	if err != nil {
		result.Error = fmt.Errorf("audio download failed: %w", err)
		result.Success = false
		return result
	}

	if progressCallback != nil {
		progressCallback(60)
	}

	s.logger.Printf("Audio downloaded: %d bytes", bytesWritten)

	// Step 4: Validate audio file
	if err := s.fileDownloader.ValidateFile(audioPath); err != nil {
		os.Remove(audioPath)
		result.Error = fmt.Errorf("audio validation failed: %w", err)
		result.Success = false
		return result
	}

	result.AudioPath = audioPath

	// Step 5: Download cover image (optional, graceful degradation)
	if metadata.CoverURL != "" {
		coverPath := filepath.Join(podcastDir, "cover.jpg")
		s.logger.Printf("Downloading cover to: %s", coverPath)

		_, err := s.imageDownloader.Download(ctx, metadata.CoverURL, coverPath, nil)
		if err != nil {
			s.logger.Printf("Warning: Cover image download failed: %v (continuing anyway)", err)
		} else {
			result.CoverPath = coverPath
			s.logger.Printf("Cover saved: %s", coverPath)
		}
	}

	if progressCallback != nil {
		progressCallback(80)
	}

	// Step 6: Save show notes (optional, graceful degradation)
	if metadata.ShowNotes != "" {
		showNotesPath := filepath.Join(podcastDir, "shownotes.txt")
		s.logger.Printf("Saving show notes to: %s", showNotesPath)

		err := s.showNotesSaver.Save(metadata.ShowNotes, showNotesPath)
		if err != nil {
			s.logger.Printf("Warning: Show notes save failed: %v (continuing anyway)", err)
		} else {
			result.ShowNotesPath = showNotesPath
			s.logger.Printf("Show notes saved: %s", showNotesPath)
		}
	}

	if progressCallback != nil {
		progressCallback(90)
	}

	// Step 7: Save metadata file
	metadataPath := filepath.Join(podcastDir, ".metadata.json")
	if err := s.saveMetadataFile(metadataPath, url, metadata, result); err != nil {
		s.logger.Printf("Warning: Failed to save metadata file: %v", err)
		// Don't fail the download if metadata saving fails
	}

	if progressCallback != nil {
		progressCallback(100)
	}

	result.Success = true
	return result
}

// saveMetadataFile saves the .metadata.json file with download information.
func (s *DownloadService) saveMetadataFile(path, url string, metadata *downloader.EpisodeMetadata, result *DownloadResult) error {
	metadataFile := MetadataFile{
		SourceURL:    url,
		Title:        metadata.Title,
		DownloadedAt: time.Now().Format(time.RFC3339),
		AudioFile:    "podcast.m4a",
	}

	if result.CoverPath != "" {
		metadataFile.CoverFile = "cover.jpg"
	}

	if result.ShowNotesPath != "" {
		metadataFile.ShowNotesFile = "shownotes.txt"
	}

	data, err := json.MarshalIndent(metadataFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// sanitizeDirectoryName creates a filesystem-safe directory name from a title.
func sanitizeDirectoryName(title string) string {
	// Remove invalid filesystem characters
	invalidChars := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	result := title
	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}

	// Limit length
	if len(result) > 200 {
		result = result[:200]
	}

	// Trim whitespace
	result = strings.TrimSpace(result)

	// Fallback if empty
	if result == "" {
		result = "unknown_podcast"
	}

	return result
}
