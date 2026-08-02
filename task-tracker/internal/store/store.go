package store

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/itbem-corp/task-tracker/internal/models"
)

// ErrTaskNotFound is returned when a task cannot be found by ID.
var ErrTaskNotFound = errors.New("task not found")

// TaskStore provides thread-safe in-memory storage for tasks.
type TaskStore struct {
	mu    sync.RWMutex
	tasks map[string]models.Task
}

// New creates and returns an initialized TaskStore.
func New() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]models.Task),
	}
}

// Create adds a task to the store and returns it.
func (s *TaskStore) Create(task models.Task) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks[task.ID] = task
	return task
}

// GetAll returns all tasks sorted by created_at descending (newest first).
func (s *TaskStore) GetAll() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		result = append(result, t)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})
	return result
}

// GetByID returns a single task or an error if not found.
func (s *TaskStore) GetByID(id string) (models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return models.Task{}, ErrTaskNotFound
	}
	return task, nil
}

// Update modifies an existing task's fields and returns the updated task.
func (s *TaskStore) Update(id string, req models.UpdateTaskRequest) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return models.Task{}, ErrTaskNotFound
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

	s.tasks[id] = task
	return task, nil
}

// Delete removes a task by ID. Returns an error if not found.
func (s *TaskStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrTaskNotFound
	}
	delete(s.tasks, id)
	return nil
}
