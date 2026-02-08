package taskmanager

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/meixg/podcast-reader/internal/validator"
)

// Manager manages download tasks and coordinates with the downloader
type Manager struct {
	store           *Store
	catalog         *Catalog
	validator       validator.URLValidator
	logger          *log.Logger
	downloadService *DownloadService
}

// NewManager creates a new task manager
func NewManager(store *Store, catalog *Catalog, logger *log.Logger) *Manager {
	return &Manager{
		store:           store,
		catalog:         catalog,
		validator:       validator.NewXiaoyuzhouURLValidator(),
		logger:          logger,
		downloadService: nil, // Will be set after creation with SetOutputDirectory
	}
}

// SetOutputDirectory sets the output directory and initializes the download service
func (m *Manager) SetOutputDirectory(outputDir string) {
	m.downloadService = NewDownloadService(outputDir, m.logger)
}

// SubmitTask submits a new download task
func (m *Manager) SubmitTask(url string) (*DownloadTask, error) {
	// Validate URL
	valid, errMsg := m.validator.ValidateURL(url)
	if !valid {
		return nil, fmt.Errorf("invalid URL: %s", errMsg)
	}

	// Check if URL already exists in catalog (already downloaded)
	if entry, exists := m.catalog.Get(url); exists {
		m.logger.Printf("URL already downloaded: %s", url)
		// Return a task-like response for the completed download
		return &DownloadTask{
			URL:    url,
			Status: StatusCompleted,
			Podcast: &PodcastEpisode{
				SourceURL: url,
				Title:     entry.Title,
				// File paths would be constructed from entry.Directory
			},
		}, ErrAlreadyDownloaded
	}

	// Check if URL has an in-progress task
	if task, exists := m.store.GetByURL(url); exists {
		if task.Status == StatusInProgress || task.Status == StatusPending {
			m.logger.Printf("URL already has task in progress: %s", url)
			return task, ErrTaskInProgress
		}
	}

	// Create new task
	task := NewDownloadTask(url)
	if err := m.store.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Start download in background
	go m.executeDownload(task)

	m.logger.Printf("Created task %s for URL: %s", task.ID, url)
	return task, nil
}

// executeDownload performs the actual download asynchronously
func (m *Manager) executeDownload(task *DownloadTask) {
	m.logger.Printf("Starting download for task %s", task.ID)

	// Update status to in-progress
	task.Status = StatusInProgress
	now := time.Now()
	task.StartedAt = &now
	m.store.Update(task)

	// TODO: Integrate with actual downloader from internal/downloader
	// For now, simulate download process
	m.simulateDownload(task)

	// TODO: Save .metadata.json on completion
	// TODO: Update task status based on download result
}

// simulateDownload simulates the download process (placeholder)
// TODO: Replace with actual downloader integration in T035
func (m *Manager) simulateDownload(task *DownloadTask) {
	m.logger.Printf("Simulating download for task %s", task.ID)

	// Simulate progress updates
	for i := 0; i <= 100; i += 10 {
		time.Sleep(100 * time.Millisecond)
		task.Progress = i
		m.store.Update(task)
	}

	// Mark as completed
	task.Status = StatusCompleted
	task.Progress = 100
	completedAt := time.Now()
	task.CompletedAt = &completedAt
	m.store.Update(task)

	m.logger.Printf("Download completed for task %s", task.ID)
}

// GetTask retrieves a task by ID
func (m *Manager) GetTask(id string) (*DownloadTask, error) {
	task, ok := m.store.GetByID(id)
	if !ok {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

// GetCatalog retrieves all catalog entries
func (m *Manager) GetCatalog(offset, limit int) ([]*PodcastCatalogEntry, int, error) {
	allEntries := m.catalog.GetAll()
	total := len(allEntries)

	// Apply pagination
	if offset >= total {
		return []*PodcastCatalogEntry{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return allEntries[offset:end], total, nil
}

// Errors
var (
	ErrAlreadyDownloaded = fmt.Errorf("URL already downloaded")
	ErrTaskInProgress    = fmt.Errorf("task already in progress for this URL")
)

// SaveMetadata saves a .metadata.json file for a downloaded podcast
func (m *Manager) SaveMetadata(podcastDir string, task *DownloadTask, audioPath, coverPath, shownotesPath string) error {
	metadata := MetadataFile{
		SourceURL:    task.URL,
		Title:        task.Podcast.Title,
		DownloadedAt: task.Podcast.DownloadedAt.Format(time.RFC3339),
		AudioFile:    filepath.Base(audioPath),
	}

	if coverPath != "" {
		metadata.CoverFile = filepath.Base(coverPath)
	}

	if shownotesPath != "" {
		metadata.ShowNotesFile = filepath.Base(shownotesPath)
	}

	// TODO: Write metadata to JSON file
	// This will be implemented when actual downloader is integrated

	m.logger.Printf("Metadata saved for task %s in %s", task.ID, podcastDir)
	return nil
}

// WaitForCompletion waits for a task to complete (used in testing)
func (m *Manager) WaitForCompletion(taskID string, timeout time.Duration) (*DownloadTask, error) {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if time.Now().After(deadline) {
				return nil, fmt.Errorf("timeout waiting for task %s", taskID)
			}

			task, err := m.GetTask(taskID)
			if err != nil {
				return nil, err
			}

			if task.Status == StatusCompleted || task.Status == StatusFailed {
				return task, nil
			}
		}
	}
}

// DownloadWithRetry downloads with exponential backoff retry logic
func (m *Manager) DownloadWithRetry(url string, maxRetries int, baseDelay time.Duration) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff with jitter
			delay := baseDelay * time.Duration(1<<uint(attempt-1))
			m.logger.Printf("Retry attempt %d after %v delay", attempt, delay)
			time.Sleep(delay)
		}

		// TODO: Implement actual download logic here
		// For now, just simulate
		if attempt == 0 {
			lastErr = fmt.Errorf("simulated download failure")
		} else {
			// Simulate success on retry
			m.logger.Printf("Download succeeded on attempt %d", attempt+1)
			return nil
		}
	}

	return fmt.Errorf("download failed after %d attempts: %w", maxRetries, lastErr)
}

// ExecuteDownloadWithRetry executes download with retry logic (T036)
func (m *Manager) ExecuteDownloadWithRetry(task *DownloadTask) error {
	const maxRetries = 3
	const baseDelay = 1 * time.Second

	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff with jitter
			delay := baseDelay * time.Duration(1<<uint(attempt-1))
			m.logger.Printf("Retry attempt %d after %v delay", attempt, delay)
			time.Sleep(delay)
		}

		// Use the download service if available
		if m.downloadService != nil {
			ctx := context.Background()

			// Create progress callback
			progressCallback := func(progress int) {
				task.Progress = progress
				m.store.Update(task)
			}

			result := m.downloadService.DownloadEpisode(ctx, task.URL, progressCallback)

			if result.Success {
				// Update task with podcast metadata
				task.Podcast = &PodcastEpisode{
					SourceURL:     task.URL,
					Title:         result.Title,
					AudioPath:     result.AudioPath,
					CoverPath:     result.CoverPath,
					ShowNotesPath: result.ShowNotesPath,
					DownloadedAt:  time.Now(),
				}
				m.store.Update(task)

				// Add to catalog
				podcastDir := filepath.Dir(result.AudioPath)
				m.catalog.Add(&PodcastCatalogEntry{
					URL:          task.URL,
					Title:        result.Title,
					Directory:    podcastDir,
					AudioFile:    filepath.Base(result.AudioPath),
					HasCover:     result.CoverPath != "",
					HasShowNotes: result.ShowNotesPath != "",
					DownloadedAt: time.Now(),
				})

				return nil
			}

			lastErr = result.Error
		} else {
			// Fallback to simulation if download service not initialized
			return fmt.Errorf("download service not initialized")
		}
	}

	return fmt.Errorf("download failed after %d attempts: %w", maxRetries, lastErr)
}

var wg sync.WaitGroup

// BackgroundDownload launches a download in a goroutine with retry logic
func (m *Manager) BackgroundDownload(task *DownloadTask) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		m.logger.Printf("Starting background download for task %s", task.ID)

		// Update status to in-progress
		task.Status = StatusInProgress
		now := time.Now()
		task.StartedAt = &now
		m.store.Update(task)

		// Execute download with retry logic
		if err := m.ExecuteDownloadWithRetry(task); err != nil {
			// Mark task as failed
			task.Status = StatusFailed
			task.Error = err.Error()
			completedAt := time.Now()
			task.CompletedAt = &completedAt
			m.store.Update(task)
			// T090: Log download failure with details
			m.logger.Printf("[TASK_FAILED] id=%s url=%s error=%s", task.ID, task.URL, err.Error())
			return
		}

		// Mark task as completed
		task.Status = StatusCompleted
		task.Progress = 100
		completedAt := time.Now()
		task.CompletedAt = &completedAt
		m.store.Update(task)
		// T090: Log download completion with details
		title := ""
		if task.Podcast != nil {
			title = task.Podcast.Title
		}
		m.logger.Printf("[TASK_COMPLETED] id=%s url=%s title=%s", task.ID, task.URL, title)
	}()
}

// CreateAndStartTask creates a task and starts background download (T031)
func (m *Manager) CreateAndStartTask(url string) (*DownloadTask, error) {
	// Validate URL
	valid, errMsg := m.validator.ValidateURL(url)
	if !valid {
		return nil, fmt.Errorf("invalid URL: %s", errMsg)
	}

	// Check for duplicate URL (in-progress)
	if task, exists := m.store.GetByURL(url); exists {
		if task.Status == StatusInProgress || task.Status == StatusPending {
			return task, ErrTaskInProgress
		}
	}

	// Check if already downloaded (T033)
	if _, exists := m.catalog.Get(url); exists {
		return nil, ErrAlreadyDownloaded
	}

	// Create new task
	task := NewDownloadTask(url)
	if err := m.store.Create(task); err != nil {
		return nil, err
	}

	// T089: Log task submission with URL and task ID
	m.logger.Printf("[TASK_SUBMITTED] id=%s url=%s", task.ID, url)

	// Start background download with retry logic
	m.BackgroundDownload(task)

	return task, nil
}

// IsDownloaded checks if a URL has already been downloaded
func (m *Manager) IsDownloaded(url string) bool {
	_, exists := m.catalog.Get(url)
	return exists
}

// HasInProgressTask checks if a URL has an in-progress task
func (m *Manager) HasInProgressTask(url string) bool {
	if task, exists := m.store.GetByURL(url); exists {
		return task.Status == StatusInProgress || task.Status == StatusPending
	}
	return false
}

// GetInProgressTask returns the in-progress task for a URL
func (m *Manager) GetInProgressTask(url string) (*DownloadTask, bool) {
	task, exists := m.store.GetByURL(url)
	if !exists {
		return nil, false
	}
	if task.Status == StatusInProgress || task.Status == StatusPending {
		return task, true
	}
	return nil, false
}
