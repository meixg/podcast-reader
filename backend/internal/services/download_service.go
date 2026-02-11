package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/meixg/podcast-reader/backend/internal/downloader"
	"github.com/meixg/podcast-reader/backend/internal/models"
)

// DownloadService handles the complete podcast download workflow
type DownloadService struct {
	downloadsDir        string
	httpClient          *http.Client
	urlExtractor        downloader.URLExtractor
	fileDownloader      downloader.FileDownloader
	imageDownloader     downloader.ImageDownloader
	metadataExtractor   *downloader.MetadataExtractor
	metadataWriter      *downloader.MetadataWriter
	taskService         *TaskService
}

// NewDownloadService creates a new download service
func NewDownloadService(downloadsDir string, taskService *TaskService) *DownloadService {
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	// Create HTTP client wrapper for goquery
	httpClientWrapper := &httpClientDoer{client: httpClient}

	return &DownloadService{
		downloadsDir:      downloadsDir,
		httpClient:        httpClient,
		urlExtractor:      downloader.NewHTMLExtractor(httpClientWrapper),
		fileDownloader:    downloader.NewHTTPDownloader(httpClient, false),
		imageDownloader:   downloader.NewHTTPImageDownloader(httpClient, 10*1024*1024), // 10MB max
		metadataExtractor: downloader.NewMetadataExtractor(httpClientWrapper),
		metadataWriter:    downloader.NewMetadataWriter(),
		taskService:       taskService,
	}
}

// httpClientDoer wraps http.Client to implement the Doer interface
type httpClientDoer struct {
	client *http.Client
}

// Get implements the Doer interface for goquery
func (h *httpClientDoer) Get(url string) (*goquery.Document, error) {
	resp, err := h.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, url)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}

// ExecuteDownload executes the complete download workflow for a task
func (s *DownloadService) ExecuteDownload(ctx context.Context, taskID, url string) {
	// Update task status to downloading
	s.taskService.UpdateProgress(taskID, 0)

	// Step 1: Extract metadata (30% progress)
	metadata, err := s.extractMetadata(ctx, url)
	if err != nil {
		s.taskService.MarkFailed(taskID, fmt.Sprintf("提取元数据失败: %v", err))
		return
	}
	s.taskService.UpdateProgress(taskID, 30)

	// Step 2: Create podcast directory (40% progress)
	podcastDir, err := s.createPodcastDir(metadata.Title)
	if err != nil {
		s.taskService.MarkFailed(taskID, fmt.Sprintf("创建目录失败: %v", err))
		return
	}
	s.taskService.UpdateProgress(taskID, 40)

	// Step 3: Download audio file (40-90% progress)
	audioPath := filepath.Join(podcastDir, "podcast.m4a")
	err = s.downloadAudio(ctx, metadata.AudioURL, audioPath, taskID)
	if err != nil {
		s.taskService.MarkFailed(taskID, fmt.Sprintf("下载音频失败: %v", err))
		return
	}
	s.taskService.UpdateProgress(taskID, 90)

	// Step 4: Download cover image (95% progress)
	if metadata.CoverURL != "" {
		coverPath := filepath.Join(podcastDir, "cover.jpg")
		if err := s.downloadCover(ctx, metadata.CoverURL, coverPath); err != nil {
			log.Printf("Warning: Failed to download cover: %v", err)
		}
	}
	s.taskService.UpdateProgress(taskID, 95)

	// Step 5: Save show notes (95% progress)
	if metadata.ShowNotes != "" {
		shownotesPath := filepath.Join(podcastDir, "shownotes.txt")
		if err := s.saveShowNotes(metadata.ShowNotes, shownotesPath); err != nil {
			log.Printf("Warning: Failed to save show notes: %v", err)
		}
	}
	s.taskService.UpdateProgress(taskID, 95)

	// Step 6: Extract and save metadata (98% progress)
	// Continue even if metadata extraction fails
	s.taskService.UpdateTaskStatus(taskID, "extracting_metadata")
	if err := s.extractAndSaveMetadata(ctx, url, podcastDir); err != nil {
		log.Printf("Warning: Failed to extract metadata: %v", err)
		// Continue - metadata extraction failure doesn't block download
	}
	s.taskService.UpdateProgress(taskID, 98)

	// Step 7: Mark as completed (100% progress)
	s.taskService.MarkCompleted(taskID, "")
	s.taskService.UpdateProgress(taskID, 100)

	log.Printf("Download completed: %s -> %s", metadata.Title, podcastDir)
}

// IsAlreadyDownloaded checks if an episode has already been downloaded
func (s *DownloadService) IsAlreadyDownloaded(url string) (bool, error) {
	// Extract metadata to get the title
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	metadata, err := s.extractMetadata(ctx, url)
	if err != nil {
		return false, err
	}

	// Generate the podcast directory path
	safeTitle := sanitizeFilename(metadata.Title)
	podcastDir := filepath.Join(s.downloadsDir, safeTitle)

	// Check if directory exists and contains podcast.m4a
	audioPath := filepath.Join(podcastDir, "podcast.m4a")
	if _, err := os.Stat(audioPath); err == nil {
		return true, nil
	}

	return false, nil
}

// extractMetadata extracts episode metadata from the URL
func (s *DownloadService) extractMetadata(ctx context.Context, url string) (*downloader.EpisodeMetadata, error) {
	metadata, err := s.urlExtractor.ExtractURL(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to extract metadata: %w", err)
	}
	return metadata, nil
}

// createPodcastDir creates a directory for the podcast episode
func (s *DownloadService) createPodcastDir(title string) (string, error) {
	// Sanitize title for filename
	safeTitle := sanitizeFilename(title)
	podcastDir := filepath.Join(s.downloadsDir, safeTitle)

	if err := os.MkdirAll(podcastDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	return podcastDir, nil
}

// downloadAudio downloads the audio file with progress tracking
func (s *DownloadService) downloadAudio(ctx context.Context, audioURL, destPath string, taskID string) error {
	// Create a progress writer that updates the task
	progressWriter := &taskProgressWriter{
		taskService: s.taskService,
		taskID:      taskID,
		minProgress: 40,
		maxProgress: 90,
	}

	bytesWritten, err := s.fileDownloader.Download(ctx, audioURL, destPath, progressWriter)
	if err != nil {
		return err
	}

	// Validate the downloaded file
	if err := s.fileDownloader.ValidateFile(destPath); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("invalid audio file: %w", err)
	}

	log.Printf("Downloaded audio: %s (%d bytes)", destPath, bytesWritten)
	return nil
}

// downloadCover downloads the cover image
func (s *DownloadService) downloadCover(ctx context.Context, coverURL, destPath string) error {
	_, err := s.imageDownloader.Download(ctx, coverURL, destPath, nil)
	if err != nil {
		return err
	}

	log.Printf("Downloaded cover: %s", destPath)
	return nil
}

// saveShowNotes saves show notes to a text file
func (s *DownloadService) saveShowNotes(htmlContent, destPath string) error {
	// Convert HTML to plain text
	textContent := convertHTMLToText(htmlContent)

	// Add UTF-8 BOM for proper Chinese character display
	content := "\uFEFF" + textContent

	if err := os.WriteFile(destPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write show notes: %w", err)
	}

	log.Printf("Saved show notes: %s", destPath)
	return nil
}

// taskProgressWriter writes progress updates to the task
type taskProgressWriter struct {
	taskService *TaskService
	taskID      string
	minProgress int
	maxProgress int
	totalBytes  int64
	written     int64
}

func (w *taskProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	w.written += int64(n)

	// Calculate progress percentage based on expected size
	// Since we don't know the total size upfront, we'll just increment
	progress := w.minProgress + (w.maxProgress-w.minProgress)/2 // Placeholder
	w.taskService.UpdateProgress(w.taskID, progress)

	return n, nil
}

// sanitizeFilename removes invalid characters from filenames
func sanitizeFilename(name string) string {
	// Remove invalid characters
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := name
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "")
	}

	// Limit length
	if len(result) > 200 {
		result = result[:200]
	}

	return strings.TrimSpace(result)
}

// convertHTMLToText converts HTML content to plain text
func convertHTMLToText(html string) string {
	// Simple HTML to text conversion
	// Remove HTML tags and decode entities
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html // Return original if parsing fails
	}

	// Get text content
	text := doc.Text()

	// Clean up whitespace
	lines := strings.Split(text, "\n")
	var cleanedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return strings.Join(cleanedLines, "\n\n")
}

// extractAndSaveMetadata extracts metadata from the page and saves it to .metadata.json
func (s *DownloadService) extractAndSaveMetadata(ctx context.Context, pageURL, podcastDir string) error {
	// Extract metadata from page
	metadata, err := s.metadataExtractor.ExtractMetadata(ctx, pageURL)
	if err != nil {
		// Log warning but don't fail - write empty metadata file
		log.Printf("Warning: Metadata extraction failed for %s: %v", pageURL, err)
		metadata = models.NewPodcastMetadata()
	}

	// Write metadata file (even if empty)
	if err := s.metadataWriter.WriteMetadata(podcastDir, metadata); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	log.Printf("Metadata saved: %s/.metadata.json", podcastDir)
	return nil
}
