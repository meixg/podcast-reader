package taskmanager

import (
	"time"

	"github.com/google/uuid"
)

// TaskStatus represents the current state of a download task
type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
)

// DownloadTask represents a podcast download request
type DownloadTask struct {
	ID          string          `json:"id"`
	URL         string          `json:"url"`
	Status      TaskStatus      `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	StartedAt   *time.Time      `json:"started_at,omitempty"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
	Error       string          `json:"error,omitempty"`
	Progress    int             `json:"progress,omitempty"`
	Podcast     *PodcastEpisode `json:"podcast,omitempty"`
}

// PodcastEpisode represents a downloaded podcast with metadata
type PodcastEpisode struct {
	Title         string    `json:"title"`
	SourceURL     string    `json:"source_url"`
	AudioPath     string    `json:"audio_path"`
	CoverPath     string    `json:"cover_path,omitempty"`
	CoverFormat   string    `json:"cover_format,omitempty"`
	ShowNotesPath string    `json:"shownotes_path,omitempty"`
	DownloadedAt  time.Time `json:"downloaded_at"`
	FileSizeMB    float64   `json:"file_size_mb,omitempty"`
}

// NewDownloadTask creates a new download task with a generated UUID
func NewDownloadTask(url string) *DownloadTask {
	return &DownloadTask{
		ID:        uuid.New().String(),
		URL:       url,
		Status:    StatusPending,
		CreatedAt: time.Now(),
		Progress:  0,
	}
}
