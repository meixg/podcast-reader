package services

import (
	"fmt"
	"math"
	"sort"

	"github.com/meixg/podcast-reader/backend/internal/models"
	"github.com/meixg/podcast-reader/backend/pkg/scanner"
)

// EpisodeService manages episode operations
type EpisodeService struct {
	scanner *scanner.Scanner
}

// NewEpisodeService creates a new episode service
func NewEpisodeService(scanner *scanner.Scanner) *EpisodeService {
	return &EpisodeService{
		scanner: scanner,
	}
}

// GetEpisodes returns paginated episodes
func (s *EpisodeService) GetEpisodes(page, pageSize int) (*models.PaginatedEpisodes, error) {
	// Scan all episodes
	episodes, err := s.scanner.ScanEpisodes()
	if err != nil {
		return nil, fmt.Errorf("failed to scan episodes: %w", err)
	}

	// Sort by download date (newest first)
	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].DownloadDate.After(episodes[j].DownloadDate)
	})

	total := len(episodes)
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	// Calculate pagination
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		return &models.PaginatedEpisodes{
			Episodes:   []models.Episode{},
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		}, nil
	}

	if end > total {
		end = total
	}

	return &models.PaginatedEpisodes{
		Episodes:   episodes[start:end],
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetShowNotes returns show notes for a specific episode
func (s *EpisodeService) GetShowNotes(episodeID string) (string, error) {
	episodes, err := s.scanner.ScanEpisodes()
	if err != nil {
		return "", fmt.Errorf("failed to scan episodes: %w", err)
	}

	for _, episode := range episodes {
		if episode.ID == episodeID {
			return episode.ShowNotes, nil
		}
	}

	return "", fmt.Errorf("episode not found")
}
