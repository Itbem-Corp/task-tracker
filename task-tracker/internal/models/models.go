package models

import (
	"errors"
	"strings"
	"time"
)

// Task status constants.
const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
)

// Task represents a unit of work to be tracked.
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTaskRequest is the payload for creating a new task.
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// Validate checks that the create request has a non-empty title and a valid status (if provided).
func (r *CreateTaskRequest) Validate() error {
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("title is required")
	}
	if r.Status != "" && !IsValidStatus(r.Status) {
		return errors.New("invalid status: must be one of pending, in_progress, done")
	}
	return nil
}

// UpdateTaskRequest is the payload for updating an existing task.
type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}

// Validate checks that at least one field is set and status (if provided) is valid.
func (r *UpdateTaskRequest) Validate() error {
	if r.Title == nil && r.Description == nil && r.Status == nil {
		return errors.New("at least one field must be provided")
	}
	if r.Status != nil && !IsValidStatus(*r.Status) {
		return errors.New("invalid status: must be one of pending, in_progress, done")
	}
	return nil
}

// ErrorResponse is the standard JSON error envelope.
type ErrorResponse struct {
	Error string `json:"error"`
}

// IsValidStatus checks whether a given status string is allowed.
func IsValidStatus(s string) bool {
	switch s {
	case StatusPending, StatusInProgress, StatusDone:
		return true
	}
	return false
}
