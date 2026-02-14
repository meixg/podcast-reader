package models

import "time"

// PodcastMetadata represents extracted metadata for a podcast episode
type PodcastMetadata struct {
	Duration     string    `json:"duration"`      // Duration as displayed on page (e.g., "231分钟", "1小时15分钟")
	PublishTime  string    `json:"publish_time"`  // Relative publish time (e.g., "刚刚发布", "2个月前")
	EpisodeTitle string    `json:"episode_title"` // Title of the episode
	PodcastName  string    `json:"podcast_name"`  // Name of the podcast series
	SourceURL    string    `json:"source_url"`    // Original page URL (e.g., https://www.xiaoyuzhoufm.com/episode/...)
	ExtractedAt  time.Time `json:"extracted_at"`  // Timestamp when metadata was extracted
}

// EpisodeWithMetadata represents a podcast episode including its metadata for API responses
type EpisodeWithMetadata struct {
	Episode
	Metadata *PodcastMetadata `json:"metadata,omitempty"`
}

// NewPodcastMetadata creates a new PodcastMetadata with the current timestamp
func NewPodcastMetadata() *PodcastMetadata {
	return &PodcastMetadata{
		ExtractedAt: time.Now().UTC(),
	}
}

// IsEmpty returns true if the metadata has no meaningful content
func (m *PodcastMetadata) IsEmpty() bool {
	if m == nil {
		return true
	}
	return m.Duration == "" && m.PublishTime == "" && m.EpisodeTitle == "" && m.PodcastName == ""
}
