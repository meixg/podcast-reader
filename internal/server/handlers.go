package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/meixg/podcast-reader/internal/taskmanager"
)

// SubmitTaskRequest represents the request body for submitting a download task
type SubmitTaskRequest struct {
	URL string `json:"url"`
}

// Validate validates the request
func (r *SubmitTaskRequest) Validate() error {
	// URL validation is done separately by validator
	if r.URL == "" {
		return ErrEmptyURL
	}
	return nil
}

// SubmitTaskHandler handles POST /tasks requests
func SubmitTaskHandler(manager *taskmanager.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only accept POST requests
		if r.Method != http.MethodPost {
			WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST method is allowed", "")
			return
		}

		// Parse request body
		var req SubmitTaskRequest
		if err := ParseJSONRequest(r, &req); err != nil {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Failed to parse request body", err.Error())
			return
		}

		// Validate request
		if err := req.Validate(); err != nil {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Request validation failed", err.Error())
			return
		}

		// Submit task (T038, T031, T032, T033)
		task, err := manager.CreateAndStartTask(req.URL)

		// Handle already downloaded (T033, T042)
		if err == taskmanager.ErrAlreadyDownloaded {
			// Return 303 See Other with location to a hypothetical task endpoint
			// Since we don't have a task ID for downloaded items, just return 200 with the info
			WriteJSON(w, http.StatusOK, task)
			return
		}

		// Handle in-progress duplicate (T032, T041)
		if err == taskmanager.ErrTaskInProgress {
			WriteJSON(w, http.StatusConflict, task)
			return
		}

		// Handle other errors (invalid URL, etc.) (T043, T039)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "INVALID_URL", "The provided URL is not valid", err.Error())
			return
		}

		// Return 202 Accepted with task response (T040)
		WriteJSON(w, http.StatusAccepted, task)
	}
}

// Errors
var (
	ErrEmptyURL = &RequestError{Message: "URL field is required"}
)

// RequestError represents an error in request parsing
type RequestError struct {
	Message string
}

func (e *RequestError) Error() string {
	return e.Message
}

// GetTaskHandler handles GET /tasks/{id} requests (T052)
func GetTaskHandler(manager *taskmanager.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only accept GET requests
		if r.Method != http.MethodGet {
			WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET method is allowed", "")
			return
		}

		// Extract task ID from URL path (T053)
		// URL format: /tasks/{id}
		path := r.URL.Path
		prefix := "/tasks/"
		if !strings.HasPrefix(path, prefix) {
			WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Invalid URL path format", "Expected /tasks/{id}")
			return
		}

		taskID := strings.TrimPrefix(path, prefix)
		if taskID == "" {
			WriteError(w, http.StatusBadRequest, "MISSING_TASK_ID", "Task ID is required", "")
			return
		}

		// Validate UUID format (T053)
		if _, err := uuid.Parse(taskID); err != nil {
			WriteError(w, http.StatusBadRequest, "INVALID_TASK_ID", "Task ID must be a valid UUID", err.Error())
			return
		}

		// Retrieve task (T054)
		task, err := manager.GetTask(taskID)
		if err != nil {
			// Handle non-existent task (T055)
			WriteError(w, http.StatusNotFound, "TASK_NOT_FOUND", "Task not found", "")
			return
		}

		// Return 200 OK with task response (T054, T056)
		// Podcast metadata is included when status=completed (T056)
		WriteJSON(w, http.StatusOK, task)
	}
}

// TasksHandler handles both POST /tasks and GET /tasks/{id} (T044, T057)
func TasksHandler(manager *taskmanager.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Route based on path pattern
		// POST /tasks - submit new task
		// GET /tasks/{id} - get task status
		if path == "/tasks" {
			if r.Method == http.MethodPost {
				SubmitTaskHandler(manager)(w, r)
				return
			}
		}

		if strings.HasPrefix(path, "/tasks/") && path != "/tasks" {
			if r.Method == http.MethodGet {
				GetTaskHandler(manager)(w, r)
				return
			}
		}

		// No matching route
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found", "")
	}
}

// ListPodcastsResponse represents the response for listing podcasts
type ListPodcastsResponse struct {
	Podcasts []*PodcastCatalogEntry `json:"podcasts"`
	Total    int                    `json:"total"`
	Limit    int                    `json:"limit,omitempty"`
	Offset   int                    `json:"offset,omitempty"`
}

// PodcastCatalogEntry represents a podcast in the catalog (T066)
type PodcastCatalogEntry struct {
	URL          string `json:"url"`
	Title        string `json:"title"`
	Directory    string `json:"directory"`
	AudioFile    string `json:"audio_file"`
	HasCover     bool   `json:"has_cover"`
	HasShowNotes bool   `json:"has_shownotes"`
}

// ListPodcastsHandler handles GET /podcasts requests (T066)
func ListPodcastsHandler(manager *taskmanager.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only accept GET requests
		if r.Method != http.MethodGet {
			WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET method is allowed", "")
			return
		}

		// Parse limit/offset query parameters (T067)
		query := r.URL.Query()
		limitStr := query.Get("limit")
		offsetStr := query.Get("offset")

		// Set defaults (T067)
		limit := 100
		offset := 0

		// Parse limit
		if limitStr != "" {
			if _, err := fmt.Sscanf(limitStr, "%d", &limit); err != nil {
				WriteError(w, http.StatusBadRequest, "INVALID_LIMIT", "Limit must be a valid integer", err.Error())
				return
			}
		}

		// Parse offset
		if offsetStr != "" {
			if _, err := fmt.Sscanf(offsetStr, "%d", &offset); err != nil {
				WriteError(w, http.StatusBadRequest, "INVALID_OFFSET", "Offset must be a valid integer", err.Error())
				return
			}
		}

		// Validate limit/offset ranges (T068)
		if limit < 1 || limit > 1000 {
			WriteError(w, http.StatusBadRequest, "INVALID_LIMIT", "Limit must be between 1 and 1000", "")
			return
		}

		if offset < 0 {
			WriteError(w, http.StatusBadRequest, "INVALID_OFFSET", "Offset must be >= 0", "")
			return
		}

		// Get catalog from manager (T062, T063)
		entries, total, err := manager.GetCatalog(offset, limit)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve catalog", err.Error())
			return
		}

		// Convert to response format (T066, T071)
		podcasts := make([]*PodcastCatalogEntry, 0, len(entries))
		for _, entry := range entries {
			podcasts = append(podcasts, &PodcastCatalogEntry{
				URL:          entry.URL,
				Title:        entry.Title,
				Directory:    entry.Directory,
				AudioFile:    entry.AudioFile,
				HasCover:     entry.HasCover,
				HasShowNotes: entry.HasShowNotes,
			})
		}

		// Return 200 OK with podcast list (T069, T070)
		response := ListPodcastsResponse{
			Podcasts: podcasts,
			Total:    total,
			Limit:    limit,
			Offset:   offset,
		}

		WriteJSON(w, http.StatusOK, response)
	}
}
