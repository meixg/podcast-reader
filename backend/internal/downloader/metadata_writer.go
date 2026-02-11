package downloader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/meixg/podcast-reader/backend/internal/models"
	"github.com/meixg/podcast-reader/backend/pkg/scanner"
)

// MetadataWriter handles writing metadata to .metadata.json files
type MetadataWriter struct {
	scanner *scanner.MetadataScanner
}

// NewMetadataWriter creates a new metadata writer
func NewMetadataWriter() *MetadataWriter {
	return &MetadataWriter{
		scanner: scanner.NewMetadataScanner(),
	}
}

// WriteMetadata writes metadata to a .metadata.json file in the specified directory
// Creates the file even if metadata is empty (with null values)
func (w *MetadataWriter) WriteMetadata(dir string, metadata *models.PodcastMetadata) error {
	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// If metadata is nil, create an empty one with just the timestamp
	if metadata == nil {
		metadata = models.NewPodcastMetadata()
	}

	// Write the metadata file
	if err := w.scanner.WriteMetadata(dir, metadata); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// MetadataExists checks if a metadata file already exists
func (w *MetadataWriter) MetadataExists(dir string) bool {
	return w.scanner.MetadataExists(dir)
}

// RemoveMetadata removes the .metadata.json file if it exists
func (w *MetadataWriter) RemoveMetadata(dir string) error {
	metadataPath := filepath.Join(dir, ".metadata.json")
	if err := os.Remove(metadataPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove metadata file: %w", err)
	}
	return nil
}
