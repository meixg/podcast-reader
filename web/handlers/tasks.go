package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/meixg/podcast-reader/pkg/models"
	"github.com/meixg/podcast-reader/web/services"
)

// TaskHandler handles task-related HTTP requests
type TaskHandler struct {
	service *services.TaskService
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

// HandleTasks handles both GET and POST /api/tasks
func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	default:
		h.sendError(w, "Method not allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
	}
}

// getTasks handles GET /api/tasks
func (h *TaskHandler) getTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.service.GetTasks()
	h.sendJSON(w, tasks, http.StatusOK)
}

// createTask handles POST /api/tasks
func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", "INVALID_REQUEST", http.StatusBadRequest)
		return
	}

	// Validate URL
	if req.URL == "" {
		h.sendError(w, "URL is required", "INVALID_URL", http.StatusBadRequest)
		return
	}

	if !strings.Contains(req.URL, "xiaoyuzhoufm.com") {
		h.sendError(w, "Invalid URL format. Must be a Xiaoyuzhou FM URL", "INVALID_URL", http.StatusBadRequest)
		return
	}

	// Create task
	task, err := h.service.CreateTask(req.URL)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			h.sendError(w, "A download task for this URL already exists", "DUPLICATE_TASK", http.StatusBadRequest)
		} else {
			h.sendError(w, "Failed to create task", "SERVER_ERROR", http.StatusInternalServerError)
		}
		return
	}

	h.sendJSON(w, task, http.StatusCreated)
}

// Helper methods
func (h *TaskHandler) sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *TaskHandler) sendError(w http.ResponseWriter, message, code string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.APIError{
		Error: message,
		Code:  code,
	})
}
