package taskmanager

import "time"

// PodcastCatalogEntry represents an entry in the download catalog
type PodcastCatalogEntry struct {
	URL          string    `json:"url"`
	Title        string    `json:"title"`
	Directory    string    `json:"directory"`
	AudioFile    string    `json:"audio_file"`
	HasCover     bool      `json:"has_cover"`
	HasShowNotes bool      `json:"has_shownotes"`
	DownloadedAt time.Time `json:"downloaded_at"`
}

// Catalog represents the in-memory index of downloaded podcasts
type Catalog struct {
	entries map[string]*PodcastCatalogEntry // key: URL
}

// NewCatalog creates a new empty catalog
func NewCatalog() *Catalog {
	return &Catalog{
		entries: make(map[string]*PodcastCatalogEntry),
	}
}

// Add adds an entry to the catalog
func (c *Catalog) Add(entry *PodcastCatalogEntry) {
	c.entries[entry.URL] = entry
}

// Get retrieves an entry by URL
func (c *Catalog) Get(url string) (*PodcastCatalogEntry, bool) {
	entry, ok := c.entries[url]
	return entry, ok
}

// GetAll returns all catalog entries
func (c *Catalog) GetAll() []*PodcastCatalogEntry {
	entries := make([]*PodcastCatalogEntry, 0, len(c.entries))
	for _, entry := range c.entries {
		entries = append(entries, entry)
	}
	return entries
}

// Count returns the total number of entries
func (c *Catalog) Count() int {
	return len(c.entries)
}
