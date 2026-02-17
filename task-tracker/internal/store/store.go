package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/itbem-corp/task-tracker/internal/models"
)

// TaskStore provides thread-safe in-memory storage for tasks.
type TaskStore struct {
	mu    sync.RWMutex
	tasks map[string]*models.Task
}

// New creates and returns an initialized TaskStore.
func New() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]*models.Task),
	}
}

// Create adds a new task and returns it. If status is empty, defaults to pending.
func (s *TaskStore) Create(title, description, status string) *models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	if status == "" {
		status = models.StatusPending
	}

	now := time.Now().UTC()
	task := &models.Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.tasks[task.ID] = task
	return task
}

// GetAll returns all tasks as a slice.
func (s *TaskStore) GetAll() []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		result = append(result, t)
	}
	return result
}

// GetByID returns a single task or an error if not found.
func (s *TaskStore) GetByID(id string) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	return task, nil
}

// Update modifies an existing task's fields and returns the updated task.
func (s *TaskStore) Update(id string, req models.UpdateTaskRequest) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found: %s", id)
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	task.UpdatedAt = time.Now().UTC()

	return task, nil
}

// Delete removes a task by ID. Returns an error if not found.
func (s *TaskStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return fmt.Errorf("task not found: %s", id)
	}
	delete(s.tasks, id)
	return nil
}
