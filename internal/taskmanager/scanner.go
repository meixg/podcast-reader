package taskmanager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Scanner handles scanning the downloads directory to build the catalog
type Scanner struct {
	catalog *Catalog
	logger  *log.Logger
}

// NewScanner creates a new directory scanner
func NewScanner(catalog *Catalog, logger *log.Logger) *Scanner {
	return &Scanner{
		catalog: catalog,
		logger:  logger,
	}
}

// ScanDownloadsDirectory scans the downloads directory and builds the catalog
func (s *Scanner) ScanDownloadsDirectory(downloadsPath string) error {
	s.logger.Printf("Scanning downloads directory: %s", downloadsPath)

	// Walk through the downloads directory
	err := filepath.Walk(downloadsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the downloads directory itself and hidden files
		if path == downloadsPath {
			return nil
		}
		if strings.HasPrefix(filepath.Base(path), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Process subdirectories (each represents a podcast)
		if info.IsDir() {
			return s.scanPodcastDirectory(path)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to scan downloads directory: %w", err)
	}

	s.logger.Printf("Scan complete: found %d podcasts", s.catalog.Count())
	return nil
}

// scanPodcastDirectory scans a single podcast directory
func (s *Scanner) scanPodcastDirectory(dirPath string) error {
	// Look for .metadata.json file
	metadataPath := filepath.Join(dirPath, ".metadata.json")
	metadata, err := s.readMetadataFile(metadataPath)
	if err != nil {
		s.logger.Printf("Warning: no metadata file in %s: %v", dirPath, err)
		// Continue without metadata - can't recover URL
		return nil
	}

	// Parse downloaded timestamp
	downloadedAt, err := time.Parse(time.RFC3339, metadata.DownloadedAt)
	if err != nil {
		s.logger.Printf("Warning: invalid timestamp in %s: %v", metadataPath, err)
		downloadedAt = time.Now() // Fallback to current time
	}

	// Convert to catalog entry
	dirName := filepath.Base(dirPath)
	entry := metadata.ToPodcastCatalogEntry(dirName)
	entry.DownloadedAt = downloadedAt

	// Add to catalog
	s.catalog.Add(entry)
	s.logger.Printf("Added to catalog: %s", metadata.Title)

	return nil
}

// readMetadataFile reads and parses a .metadata.json file
func (s *Scanner) readMetadataFile(path string) (*MetadataFile, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var metadata MetadataFile
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	// Validate required fields
	if metadata.SourceURL == "" {
		return nil, fmt.Errorf("missing required field: source_url")
	}
	if metadata.Title == "" {
		return nil, fmt.Errorf("missing required field: title")
	}

	return &metadata, nil
}

// HasMissingData checks if a podcast directory has partial data (missing cover or shownotes)
func HasMissingData(dirPath string) (missingCover, missingNotes bool, err error) {
	// Check for cover image
	coverPattern := filepath.Join(dirPath, "cover.*")
	matches, _ := filepath.Glob(coverPattern)
	missingCover = len(matches) == 0

	// Check for shownotes
	notesPath := filepath.Join(dirPath, "shownotes.txt")
	if _, err := os.Stat(notesPath); os.IsNotExist(err) {
		missingNotes = true
	}

	return missingCover, missingNotes, nil
}
