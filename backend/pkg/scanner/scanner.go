package scanner

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/meixg/podcast-reader/backend/internal/models"
)

// Scanner scans the downloads directory for podcast episodes
type Scanner struct {
	downloadsDir    string
	metadataScanner *MetadataScanner
}

// NewScanner creates a new scanner instance
func NewScanner(downloadsDir string) *Scanner {
	return &Scanner{
		downloadsDir:    downloadsDir,
		metadataScanner: NewMetadataScanner(),
	}
}

// ScanEpisodes scans the downloads directory and returns all episodes
func (s *Scanner) ScanEpisodes() ([]models.Episode, error) {
	var episodes []models.Episode

	err := filepath.Walk(s.downloadsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process audio files
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".m4a" && ext != ".mp3" {
			return nil
		}

		episode, err := s.parseEpisode(path, info)
		if err != nil {
			// Log error but continue scanning
			fmt.Printf("Error parsing episode %s: %v\n", path, err)
			return nil
		}

		episodes = append(episodes, episode)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %w", err)
	}

	return episodes, nil
}

// parseEpisode extracts episode metadata from a file
func (s *Scanner) parseEpisode(audioPath string, info os.FileInfo) (models.Episode, error) {
	// Generate ID from file path
	id := s.generateID(audioPath)

	// Get podcast name from parent directory
	podcastName := filepath.Base(filepath.Dir(audioPath))

	// Get title from filename (without extension)
	title := strings.TrimSuffix(filepath.Base(audioPath), filepath.Ext(audioPath))

	// Look for cover image
	dir := filepath.Dir(audioPath)
	coverPath := s.findCoverImage(dir)

	// Look for show notes
	showNotes := s.readShowNotes(dir)

	// Read metadata if available
	metadata, _ := s.metadataScanner.ReadMetadata(dir)

	episode := models.Episode{
		ID:             id,
		Title:          title,
		PodcastName:    podcastName,
		Duration:       "00:00", // TODO: Extract from audio metadata
		FileSize:       info.Size(),
		DownloadDate:   info.ModTime(),
		ShowNotes:      showNotes,
		FilePath:       audioPath,
		CoverImagePath: coverPath,
		Metadata:       metadata,
	}

	return episode, nil
}

// generateID creates a unique ID from the file path
func (s *Scanner) generateID(path string) string {
	hash := md5.Sum([]byte(path))
	return hex.EncodeToString(hash[:])
}

// findCoverImage looks for a cover image in the directory
func (s *Scanner) findCoverImage(dir string) string {
	coverNames := []string{"cover.jpg", "cover.png", "cover.webp"}
	for _, name := range coverNames {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// readShowNotes reads the show notes file if it exists
func (s *Scanner) readShowNotes(dir string) string {
	notesPath := filepath.Join(dir, "shownotes.txt")
	data, err := os.ReadFile(notesPath)
	if err != nil {
		return ""
	}
	return string(data)
}
