package models

import (
	"testing"
)

func TestCreateTaskRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateTaskRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid request with all fields",
			req:     CreateTaskRequest{Title: "Buy groceries", Description: "Milk, eggs, bread", Status: "pending"},
			wantErr: false,
		},
		{
			name:    "valid request with title only",
			req:     CreateTaskRequest{Title: "Buy groceries"},
			wantErr: false,
		},
		{
			name:    "valid request with status pending",
			req:     CreateTaskRequest{Title: "Task", Status: StatusPending},
			wantErr: false,
		},
		{
			name:    "valid request with status in_progress",
			req:     CreateTaskRequest{Title: "Task", Status: StatusInProgress},
			wantErr: false,
		},
		{
			name:    "valid request with status done",
			req:     CreateTaskRequest{Title: "Task", Status: StatusDone},
			wantErr: false,
		},
		{
			name:    "empty title returns error",
			req:     CreateTaskRequest{Title: ""},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name:    "whitespace-only title returns error",
			req:     CreateTaskRequest{Title: "   "},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name:    "tab-only title returns error",
			req:     CreateTaskRequest{Title: "\t\n"},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name:    "invalid status returns error",
			req:     CreateTaskRequest{Title: "Task", Status: "invalid"},
			wantErr: true,
			errMsg:  "invalid status",
		},
		{
			name:    "status with wrong case returns error",
			req:     CreateTaskRequest{Title: "Task", Status: "Pending"},
			wantErr: true,
			errMsg:  "invalid status",
		},
		{
			name:    "empty status is allowed (defaults later)",
			req:     CreateTaskRequest{Title: "Task", Status: ""},
			wantErr: false,
		},
		{
			name:    "title with leading/trailing spaces is valid",
			req:     CreateTaskRequest{Title: "  valid title  "},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error containing %q, got nil", tt.errMsg)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %q, want it to contain %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestUpdateTaskRequest_Validate(t *testing.T) {
	strPtr := func(s string) *string { return &s }

	tests := []struct {
		name    string
		req     UpdateTaskRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid update with title only",
			req:     UpdateTaskRequest{Title: strPtr("New Title")},
			wantErr: false,
		},
		{
			name:    "valid update with description only",
			req:     UpdateTaskRequest{Description: strPtr("New desc")},
			wantErr: false,
		},
		{
			name:    "valid update with status only",
			req:     UpdateTaskRequest{Status: strPtr(StatusInProgress)},
			wantErr: false,
		},
		{
			name:    "valid update with all fields",
			req:     UpdateTaskRequest{Title: strPtr("Title"), Description: strPtr("Desc"), Status: strPtr(StatusDone)},
			wantErr: false,
		},
		{
			name:    "valid update with status pending",
			req:     UpdateTaskRequest{Status: strPtr(StatusPending)},
			wantErr: false,
		},
		{
			name:    "valid update with status done",
			req:     UpdateTaskRequest{Status: strPtr(StatusDone)},
			wantErr: false,
		},
		{
			name:    "empty update all nil returns error",
			req:     UpdateTaskRequest{},
			wantErr: true,
			errMsg:  "at least one field must be provided",
		},
		{
			name:    "invalid status returns error",
			req:     UpdateTaskRequest{Status: strPtr("cancelled")},
			wantErr: true,
			errMsg:  "invalid status",
		},
		{
			name:    "invalid status with valid title returns error",
			req:     UpdateTaskRequest{Title: strPtr("Title"), Status: strPtr("nope")},
			wantErr: true,
			errMsg:  "invalid status",
		},
		{
			name:    "empty string title is allowed (not nil)",
			req:     UpdateTaskRequest{Title: strPtr("")},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error containing %q, got nil", tt.errMsg)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %q, want it to contain %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{name: "pending is valid", status: StatusPending, want: true},
		{name: "in_progress is valid", status: StatusInProgress, want: true},
		{name: "done is valid", status: StatusDone, want: true},
		{name: "empty string is invalid", status: "", want: false},
		{name: "random string is invalid", status: "random", want: false},
		{name: "uppercase Pending is invalid", status: "Pending", want: false},
		{name: "DONE uppercase is invalid", status: "DONE", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidStatus(tt.status)
			if got != tt.want {
				t.Errorf("IsValidStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

// contains checks if substr is within s (helper for error message matching).
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
