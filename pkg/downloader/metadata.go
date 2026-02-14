package downloader

import "time"

// EpisodeMetadata contains all extracted metadata for a podcast episode.
type EpisodeMetadata struct {
	AudioURL        string    // Direct URL to audio file (required)
	CoverURL        string    // URL to cover image (optional)
	ShowNotes       string    // Plain text show notes (optional)
	Title           string    // Episode title (required)
	EpisodeNumber   string    // Episode number if available
	PodcastName     string    // Podcast/series name
	PublicationDate time.Time // Publication date
}
