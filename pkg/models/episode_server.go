package models

import "time"

// DownloadedEpisode represents a downloaded podcast episode with metadata
type DownloadedEpisode struct {
	ID             string           `json:"id"`
	Title          string           `json:"title"`
	PodcastName    string           `json:"podcastName"`
	Duration       string           `json:"duration"`
	FileSize       int64            `json:"fileSize"`
	DownloadDate   time.Time        `json:"downloadDate"`
	ShowNotes      string           `json:"showNotes"`
	FilePath       string           `json:"filePath"`
	CoverImagePath string           `json:"coverImagePath,omitempty"`
	SourceURL      string           `json:"sourceUrl,omitempty"`
	Metadata       *PodcastMetadata `json:"metadata,omitempty"`
}

// PaginatedEpisodes represents a paginated response of episodes
type PaginatedEpisodes struct {
	Episodes   []DownloadedEpisode `json:"episodes"`
	Total      int                 `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"pageSize"`
	TotalPages int                 `json:"totalPages"`
}
