package models

import "time"

// Task represents a unit of work to be tracked.
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTaskRequest is the payload for creating a new task.
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// UpdateTaskRequest is the payload for updating an existing task.
type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}

// Task status constants.
const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
)

// ValidStatuses returns the set of allowed status values.
func ValidStatuses() []string {
	return []string{StatusPending, StatusInProgress, StatusDone}
}

// IsValidStatus checks whether a given status string is allowed.
func IsValidStatus(s string) bool {
	for _, v := range ValidStatuses() {
		if v == s {
			return true
		}
	}
	return false
}
