package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/meixg/podcast-reader/backend/internal/models"
)

// TaskService manages download tasks in memory
type TaskService struct {
	tasks           map[string]*models.DownloadTask
	downloadService *DownloadService
	mu              sync.RWMutex
}

// NewTaskService creates a new task service
func NewTaskService() *TaskService {
	return &TaskService{
		tasks: make(map[string]*models.DownloadTask),
	}
}

// SetDownloadService sets the download service
func (s *TaskService) SetDownloadService(ds *DownloadService) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.downloadService = ds
}

// CreateTask creates a new download task and starts the download
func (s *TaskService) CreateTask(url string) (*models.DownloadTask, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for duplicate URL - only block if there's an active task (pending or downloading)
	for _, task := range s.tasks {
		if task.URL == url && (task.Status == models.TaskStatusPending || task.Status == models.TaskStatusDownloading) {
			return nil, fmt.Errorf("task already exists for this URL")
		}
	}

	// Check if already downloaded by checking file system
	if s.downloadService != nil {
		if alreadyDownloaded, err := s.downloadService.IsAlreadyDownloaded(url); err == nil && alreadyDownloaded {
			return nil, fmt.Errorf("episode already downloaded")
		}
	}

	task := &models.DownloadTask{
		ID:        uuid.New().String(),
		URL:       url,
		Status:    models.TaskStatusPending,
		CreatedAt: time.Now(),
	}

	s.tasks[task.ID] = task

	// Start download in background
	if s.downloadService != nil {
		go s.downloadService.ExecuteDownload(context.Background(), task.ID, url)
	}

	return task, nil
}

// GetTasks returns all tasks
func (s *TaskService) GetTasks() []*models.DownloadTask {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.DownloadTask, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// GetTask returns a task by ID
func (s *TaskService) GetTask(id string) (*models.DownloadTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}

// UpdateProgress updates the progress of a task
func (s *TaskService) UpdateProgress(id string, progress int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return fmt.Errorf("task not found")
	}

	task.Progress = &progress
	// Only update status to downloading if not already completed or failed
	if task.Status != models.TaskStatusCompleted && task.Status != models.TaskStatusFailed {
		task.Status = models.TaskStatusDownloading
	}
	return nil
}

// MarkCompleted marks a task as completed
func (s *TaskService) MarkCompleted(id string, episodeID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return fmt.Errorf("task not found")
	}

	now := time.Now()
	task.Status = models.TaskStatusCompleted
	task.CompletedAt = &now
	task.EpisodeID = episodeID
	progress := 100
	task.Progress = &progress
	return nil
}

// MarkFailed marks a task as failed
func (s *TaskService) MarkFailed(id string, errorMsg string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return fmt.Errorf("task not found")
	}

	now := time.Now()
	task.Status = models.TaskStatusFailed
	task.CompletedAt = &now
	task.ErrorMessage = errorMsg
	return nil
}
