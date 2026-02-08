package taskmanager

import (
	"sync"
)

// Store manages in-memory task storage with thread-safe access
type Store struct {
	mu         sync.RWMutex
	tasks      map[string]*DownloadTask // key: task ID
	tasksByURL map[string]*DownloadTask // key: URL for duplicate detection
}

// NewStore creates a new task store
func NewStore() *Store {
	return &Store{
		tasks:      make(map[string]*DownloadTask),
		tasksByURL: make(map[string]*DownloadTask),
	}
}

// Create adds a new task to the store
func (s *Store) Create(task *DownloadTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for duplicate URL
	if _, exists := s.tasksByURL[task.URL]; exists {
		return ErrDuplicateURL
	}

	s.tasks[task.ID] = task
	s.tasksByURL[task.URL] = task
	return nil
}

// GetByID retrieves a task by ID
func (s *Store) GetByID(id string) (*DownloadTask, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	return task, ok
}

// GetByURL retrieves a task by URL
func (s *Store) GetByURL(url string) (*DownloadTask, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasksByURL[url]
	return task, ok
}

// Update updates an existing task
func (s *Store) Update(task *DownloadTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[task.ID]; !exists {
		return ErrTaskNotFound
	}

	s.tasks[task.ID] = task
	// Update URL index if URL changed (shouldn't happen normally)
	s.tasksByURL[task.URL] = task
	return nil
}

// Delete removes a task from the store
func (s *Store) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, ok := s.tasks[id]; ok {
		delete(s.tasks, id)
		delete(s.tasksByURL, task.URL)
	}
}

// GetAll returns all tasks
func (s *Store) GetAll() []*DownloadTask {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*DownloadTask, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// Count returns the total number of tasks
func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.tasks)
}

// Errors
var (
	ErrDuplicateURL = &StoreError{Message: "URL already has a task"}
	ErrTaskNotFound = &StoreError{Message: "Task not found"}
)

// StoreError represents an error in task storage operations
type StoreError struct {
	Message string
}

func (e *StoreError) Error() string {
	return e.Message
}
