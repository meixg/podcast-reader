package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/meixg/podcast-reader/backend/internal/models"
	"github.com/meixg/podcast-reader/backend/internal/services"
)

// EpisodeHandler handles episode-related HTTP requests
type EpisodeHandler struct {
	service *services.EpisodeService
}

// NewEpisodeHandler creates a new episode handler
func NewEpisodeHandler(service *services.EpisodeService) *EpisodeHandler {
	return &EpisodeHandler{
		service: service,
	}
}

// GetEpisodes handles GET /api/episodes
func (h *EpisodeHandler) GetEpisodes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	page := h.parseIntParam(r, "page", 1)
	pageSize := h.parseIntParam(r, "pageSize", 20)

	// Validate page size
	if pageSize != 20 && pageSize != 50 && pageSize != 100 {
		h.sendError(w, "Invalid page size. Must be 20, 50, or 100", "INVALID_PARAMETER", http.StatusBadRequest)
		return
	}

	// Get episodes
	result, err := h.service.GetEpisodes(page, pageSize)
	if err != nil {
		h.sendError(w, "Failed to get episodes", "SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	h.sendJSON(w, result, http.StatusOK)
}

// GetShowNotes handles GET /api/episodes/:id/shownotes
func (h *EpisodeHandler) GetShowNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	// Extract episode ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/episodes/")
	episodeID := strings.TrimSuffix(path, "/shownotes")

	if episodeID == "" {
		h.sendError(w, "Episode ID required", "INVALID_PARAMETER", http.StatusBadRequest)
		return
	}

	// Get show notes
	showNotes, err := h.service.GetShowNotes(episodeID)
	if err != nil {
		h.sendError(w, "Episode not found", "NOT_FOUND", http.StatusNotFound)
		return
	}

	h.sendJSON(w, map[string]string{"showNotes": showNotes}, http.StatusOK)
}

// Helper methods
func (h *EpisodeHandler) parseIntParam(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func (h *EpisodeHandler) sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *EpisodeHandler) sendError(w http.ResponseWriter, message, code string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.APIError{
		Error: message,
		Code:  code,
	})
}
