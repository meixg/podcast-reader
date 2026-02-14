package models

import "time"

// Status represents the current state of a download session.
type Status int

const (
	StatusPending Status = iota
	StatusInProgress
	StatusCompleted
	StatusFailed
)

// DownloadSession represents a single download operation.
type DownloadSession struct {
	// EpisodeID references the episode being downloaded
	EpisodeID string

	// FilePath is the local filesystem path for the downloaded file
	FilePath string

	// Status is the current download state
	Status Status

	// Progress is the download progress (0.0 to 1.0)
	Progress float64

	// BytesDownloaded is the number of bytes downloaded
	BytesDownloaded int64

	// TotalBytes is the total file size (-1 if unknown)
	TotalBytes int64

	// StartTime is when the download started
	StartTime time.Time

	// EndTime is when the download completed (zero if not completed)
	EndTime time.Time

	// Error contains the error message if download failed
	Error error

	// RetryCount is the number of retry attempts
	RetryCount int

	// DownloadSpeed is the current download speed in bytes/sec
	DownloadSpeed float64
}

// UpdateProgress updates the download progress based on bytes downloaded.
func (s *DownloadSession) UpdateProgress(bytesDownloaded int64) {
	s.BytesDownloaded = bytesDownloaded
	if s.TotalBytes > 0 {
		s.Progress = float64(bytesDownloaded) / float64(s.TotalBytes)
	}
}

// Complete marks the download as completed.
func (s *DownloadSession) Complete() {
	s.Status = StatusCompleted
	s.EndTime = time.Now()
	s.Progress = 1.0
	if s.TotalBytes > 0 {
		s.BytesDownloaded = s.TotalBytes
	}
}

// Fail marks the download as failed with an error.
func (s *DownloadSession) Fail(err error) {
	s.Status = StatusFailed
	s.Error = err
	s.EndTime = time.Now()
}

// CanRetry checks if the download can be retried.
func (s *DownloadSession) CanRetry(maxRetries int) bool {
	return s.RetryCount < maxRetries
}

// IncrementRetry increments the retry count.
func (s *DownloadSession) IncrementRetry() {
	s.RetryCount++
}
