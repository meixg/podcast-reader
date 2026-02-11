package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/meixg/podcast-reader/backend/internal/models"
)

const metadataFileName = ".metadata.json"

// MetadataScanner scans for .metadata.json files
type MetadataScanner struct{}

// NewMetadataScanner creates a new metadata scanner
func NewMetadataScanner() *MetadataScanner {
	return &MetadataScanner{}
}

// ReadMetadata reads the .metadata.json file from the specified directory
// Returns nil if the file doesn't exist or is invalid
func (s *MetadataScanner) ReadMetadata(dir string) (*models.PodcastMetadata, error) {
	metadataPath := filepath.Join(dir, metadataFileName)

	// Check if file exists
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		return nil, nil // File doesn't exist, not an error
	}

	// Read file
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	// Parse JSON
	var metadata models.PodcastMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata file: %w", err)
	}

	return &metadata, nil
}

// WriteMetadata writes the metadata to a .metadata.json file in the specified directory
func (s *MetadataScanner) WriteMetadata(dir string, metadata *models.PodcastMetadata) error {
	if metadata == nil {
		return fmt.Errorf("metadata cannot be nil")
	}

	metadataPath := filepath.Join(dir, metadataFileName)

	// Marshal JSON with indentation for readability
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Write file
	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// MetadataExists checks if a .metadata.json file exists in the directory
func (s *MetadataScanner) MetadataExists(dir string) bool {
	metadataPath := filepath.Join(dir, metadataFileName)
	_, err := os.Stat(metadataPath)
	return err == nil
}
