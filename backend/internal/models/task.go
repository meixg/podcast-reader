package models

import "time"

// TaskStatus represents the status of a download task
type TaskStatus string

const (
	TaskStatusPending            TaskStatus = "pending"
	TaskStatusDownloading        TaskStatus = "downloading"
	TaskStatusExtractingMetadata TaskStatus = "extracting_metadata"
	TaskStatusCompleted          TaskStatus = "completed"
	TaskStatusFailed             TaskStatus = "failed"
)

// DownloadTask represents a download operation with status tracking
type DownloadTask struct {
	ID           string     `json:"id"`
	URL          string     `json:"url"`
	Status       TaskStatus `json:"status"`
	CreatedAt    time.Time  `json:"createdAt"`
	CompletedAt  *time.Time `json:"completedAt,omitempty"`
	Progress     *int       `json:"progress,omitempty"`
	ErrorMessage string     `json:"errorMessage,omitempty"`
	EpisodeID    string     `json:"episodeId,omitempty"`
}

// CreateTaskRequest represents the request body for creating a task
type CreateTaskRequest struct {
	URL string `json:"url"`
}

// APIError represents a standard error response
type APIError struct {
	Error   string                 `json:"error"`
	Code    string                 `json:"code"`
	Details map[string]interface{} `json:"details,omitempty"`
}
