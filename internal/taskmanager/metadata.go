package taskmanager

// MetadataFile represents the .metadata.json file in podcast directories
type MetadataFile struct {
	SourceURL     string `json:"source_url"`
	Title         string `json:"title"`
	DownloadedAt  string `json:"downloaded_at"`
	AudioFile     string `json:"audio_file"`
	CoverFile     string `json:"cover_file,omitempty"`
	ShowNotesFile string `json:"shownotes_file,omitempty"`
}

// ToPodcastCatalogEntry converts metadata to a catalog entry
// Note: DownloadedAt field needs to be parsed from string to time.Time separately
func (m *MetadataFile) ToPodcastCatalogEntry(directory string) *PodcastCatalogEntry {
	return &PodcastCatalogEntry{
		URL:          m.SourceURL,
		Title:        m.Title,
		Directory:    directory,
		AudioFile:    m.AudioFile,
		HasCover:     m.CoverFile != "",
		HasShowNotes: m.ShowNotesFile != "",
		// DownloadedAt should be set by the caller
	}
}
