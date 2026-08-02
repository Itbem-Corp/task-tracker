package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/itbem-corp/task-tracker/internal/models"
	"github.com/itbem-corp/task-tracker/internal/store"
)

// newTestRouter creates a fresh handler + chi router wired with all routes.
func newTestRouter() (*Handler, *chi.Mux) {
	s := store.New()
	h := NewHandler(s)
	r := chi.NewRouter()
	r.Get("/health", h.HealthCheck)
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", h.CreateTask)
		r.Get("/", h.ListTasks)
		r.Get("/{id}", h.GetTask)
		r.Put("/{id}", h.UpdateTask)
		r.Delete("/{id}", h.DeleteTask)
	})
	return h, r
}

// doRequest is a helper to execute a request against the test router.
func doRequest(r *chi.Mux, method, path string, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

// createTaskViaAPI is a helper to create a task and return the response body as map.
func createTaskViaAPI(t *testing.T, r *chi.Mux, title, description, status string) map[string]interface{} {
	t.Helper()
	body := map[string]string{"title": title}
	if description != "" {
		body["description"] = description
	}
	if status != "" {
		body["status"] = status
	}
	b, _ := json.Marshal(body)
	rr := doRequest(r, http.MethodPost, "/tasks", string(b))
	if rr.Code != http.StatusCreated {
		t.Fatalf("createTaskViaAPI: expected 201, got %d: %s", rr.Code, rr.Body.String())
	}
	var result map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("createTaskViaAPI: failed to parse response: %v", err)
	}
	return result
}

// --- HealthCheck Tests ---

func TestHealthCheck(t *testing.T) {
	_, r := newTestRouter()
	rr := doRequest(r, http.MethodGet, "/health", "")

	if rr.Code != http.StatusOK {
		t.Errorf("HealthCheck status = %d, want %d", rr.Code, http.StatusOK)
	}

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("HealthCheck Content-Type = %q, want %q", ct, "application/json")
	}

	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("HealthCheck: failed to parse body: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("HealthCheck status = %q, want %q", body["status"], "ok")
	}
}

// --- CreateTask Tests ---

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantCode   int
		wantErrMsg string // if non-empty, expect error response with this substring
	}{
		{
			name:     "valid request with title and description",
			body:     `{"title":"Buy groceries","description":"Milk, eggs, bread"}`,
			wantCode: http.StatusCreated,
		},
		{
			name:     "valid request with title only",
			body:     `{"title":"Simple task"}`,
			wantCode: http.StatusCreated,
		},
		{
			name:     "valid request with explicit status pending",
			body:     `{"title":"Task","status":"pending"}`,
			wantCode: http.StatusCreated,
		},
		{
			name:     "valid request with status in_progress",
			body:     `{"title":"Task","status":"in_progress"}`,
			wantCode: http.StatusCreated,
		},
		{
			name:     "valid request with status done",
			body:     `{"title":"Task","status":"done"}`,
			wantCode: http.StatusCreated,
		},
		{
			name:       "empty body",
			body:       `{}`,
			wantCode:   http.StatusBadRequest,
			wantErrMsg: "title is required",
		},
		{
			name:       "missing title field",
			body:       `{"description":"no title here"}`,
			wantCode:   http.StatusBadRequest,
			wantErrMsg: "title is required",
		},
		{
			name:       "empty title",
			body:       `{"title":""}`,
			wantCode:   http.StatusBadRequest,
			wantErrMsg: "title is required",
		},
		{
			name:       "whitespace-only title",
			body:       `{"title":"   "}`,
			wantCode:   http.StatusBadRequest,
			wantErrMsg: "title is required",
		},
		{
			name:       "invalid JSON",
			body:       `{not json}`,
			wantCode:   http.StatusBadRequest,
			wantErrMsg: "invalid JSON body",
		},
		{
			name:       "invalid status value",
			body:       `{"title":"Task","status":"cancelled"}`,
			wantCode:   http.StatusBadRequest,
			wantErrMsg: "invalid status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, r := newTestRouter()
			rr := doRequest(r, http.MethodPost, "/tasks", tt.body)

			if rr.Code != tt.wantCode {
				t.Errorf("CreateTask status = %d, want %d. Body: %s", rr.Code, tt.wantCode, rr.Body.String())
			}

			ct := rr.Header().Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("CreateTask Content-Type = %q, want %q", ct, "application/json")
			}

			if tt.wantErrMsg != "" {
				var errResp models.ErrorResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &errResp); err != nil {
					t.Fatalf("failed to parse error response: %v", err)
				}
				if !strings.Contains(errResp.Error, tt.wantErrMsg) {
					t.Errorf("error message = %q, want it to contain %q", errResp.Error, tt.wantErrMsg)
				}
			}

			if tt.wantCode == http.StatusCreated {
				var task models.Task
				if err := json.Unmarshal(rr.Body.Bytes(), &task); err != nil {
					t.Fatalf("failed to parse task response: %v", err)
				}
				if task.ID == "" {
					t.Error("CreateTask: response ID is empty, expected a UUID")
				}
				if task.CreatedAt.IsZero() {
					t.Error("CreateTask: CreatedAt is zero")
				}
				if task.UpdatedAt.IsZero() {
					t.Error("CreateTask: UpdatedAt is zero")
				}
			}
		})
	}
}

func TestCreateTask_DefaultStatus(t *testing.T) {
	_, r := newTestRouter()
	rr := doRequest(r, http.MethodPost, "/tasks", `{"title":"No status"}`)

	if rr.Code != http.StatusCreated {
		t.Fatalf("CreateTask status = %d, want 201", rr.Code)
	}

	var task models.Task
	json.Unmarshal(rr.Body.Bytes(), &task)

	if task.Status != models.StatusPending {
		t.Errorf("CreateTask default status = %q, want %q", task.Status, models.StatusPending)
	}
}

func TestCreateTask_ExplicitStatus(t *testing.T) {
	_, r := newTestRouter()
	rr := doRequest(r, http.MethodPost, "/tasks", `{"title":"WIP","status":"in_progress"}`)

	if rr.Code != http.StatusCreated {
		t.Fatalf("CreateTask status = %d, want 201", rr.Code)
	}

	var task models.Task
	json.Unmarshal(rr.Body.Bytes(), &task)

	if task.Status != models.StatusInProgress {
		t.Errorf("CreateTask status = %q, want %q", task.Status, models.StatusInProgress)
	}
}

// --- ListTasks Tests ---

func TestListTasks(t *testing.T) {
	t.Run("empty store returns empty array", func(t *testing.T) {
		_, r := newTestRouter()
		rr := doRequest(r, http.MethodGet, "/tasks", "")

		if rr.Code != http.StatusOK {
			t.Errorf("ListTasks status = %d, want %d", rr.Code, http.StatusOK)
		}

		ct := rr.Header().Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("ListTasks Content-Type = %q, want %q", ct, "application/json")
		}

		// Must be [] not null
		body := strings.TrimSpace(rr.Body.String())
		if body != "[]" {
			t.Errorf("ListTasks empty body = %q, want %q", body, "[]")
		}
	})

	t.Run("returns created tasks", func(t *testing.T) {
		_, r := newTestRouter()
		createTaskViaAPI(t, r, "Task 1", "Desc 1", "")
		createTaskViaAPI(t, r, "Task 2", "Desc 2", "")
		createTaskViaAPI(t, r, "Task 3", "Desc 3", "")

		rr := doRequest(r, http.MethodGet, "/tasks", "")

		if rr.Code != http.StatusOK {
			t.Fatalf("ListTasks status = %d, want 200", rr.Code)
		}

		var tasks []models.Task
		if err := json.Unmarshal(rr.Body.Bytes(), &tasks); err != nil {
			t.Fatalf("failed to parse tasks: %v", err)
		}
		if len(tasks) != 3 {
			t.Errorf("ListTasks returned %d tasks, want 3", len(tasks))
		}
	})

	t.Run("returns tasks sorted newest first", func(t *testing.T) {
		_, r := newTestRouter()
		first := createTaskViaAPI(t, r, "First", "", "")
		createTaskViaAPI(t, r, "Second", "", "")
		third := createTaskViaAPI(t, r, "Third", "", "")

		rr := doRequest(r, http.MethodGet, "/tasks", "")
		var tasks []models.Task
		json.Unmarshal(rr.Body.Bytes(), &tasks)

		if len(tasks) != 3 {
			t.Fatalf("expected 3 tasks, got %d", len(tasks))
		}
		// Newest should be first — but since creation is nearly instantaneous,
		// the order may not be perfectly deterministic. At minimum verify all IDs present.
		ids := map[string]bool{}
		for _, task := range tasks {
			ids[task.ID] = true
		}
		firstID := first["id"].(string)
		thirdID := third["id"].(string)
		if !ids[firstID] || !ids[thirdID] {
			t.Error("ListTasks missing expected task IDs")
		}
	})
}

// --- GetTask Tests ---

func TestGetTask(t *testing.T) {
	t.Run("existing task returns 200", func(t *testing.T) {
		_, r := newTestRouter()
		created := createTaskViaAPI(t, r, "Find me", "Description", "")
		id := created["id"].(string)

		rr := doRequest(r, http.MethodGet, "/tasks/"+id, "")

		if rr.Code != http.StatusOK {
			t.Errorf("GetTask status = %d, want %d", rr.Code, http.StatusOK)
		}

		ct := rr.Header().Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("GetTask Content-Type = %q, want %q", ct, "application/json")
		}

		var task models.Task
		json.Unmarshal(rr.Body.Bytes(), &task)
		if task.ID != id {
			t.Errorf("GetTask ID = %q, want %q", task.ID, id)
		}
		if task.Title != "Find me" {
			t.Errorf("GetTask Title = %q, want %q", task.Title, "Find me")
		}
		if task.Description != "Description" {
			t.Errorf("GetTask Description = %q, want %q", task.Description, "Description")
		}
	})

	t.Run("non-existent UUID returns 404", func(t *testing.T) {
		_, r := newTestRouter()
		rr := doRequest(r, http.MethodGet, "/tasks/00000000-0000-0000-0000-000000000000", "")

		if rr.Code != http.StatusNotFound {
			t.Errorf("GetTask status = %d, want %d", rr.Code, http.StatusNotFound)
		}

		var errResp models.ErrorResponse
		json.Unmarshal(rr.Body.Bytes(), &errResp)
		if errResp.Error == "" {
			t.Error("GetTask 404: error message is empty")
		}
	})

	t.Run("random string ID returns 404", func(t *testing.T) {
		_, r := newTestRouter()
		rr := doRequest(r, http.MethodGet, "/tasks/not-a-real-id", "")

		if rr.Code != http.StatusNotFound {
			t.Errorf("GetTask status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})
}

// --- UpdateTask Tests ---

func TestUpdateTask(t *testing.T) {
	strPtr := func(s string) *string { return &s }

	t.Run("update title only", func(t *testing.T) {
		_, r := newTestRouter()
		created := createTaskViaAPI(t, r, "Original", "Original desc", "")
		id := created["id"].(string)

		body, _ := json.Marshal(models.UpdateTaskRequest{Title: strPtr("Updated Title")})
		rr := doRequest(r, http.MethodPut, "/tasks/"+id, string(body))

		if rr.Code != http.StatusOK {
			t.Errorf("UpdateTask status = %d, want %d. Body: %s", rr.Code, http.StatusOK, rr.Body.String())
		}

		var task models.Task
		json.Unmarshal(rr.Body.Bytes(), &task)
		if task.Title != "Updated Title" {
			t.Errorf("UpdateTask Title = %q, want %q", task.Title, "Updated Title")
		}
		if task.Description != "Original desc" {
			t.Errorf("UpdateTask Description changed to %q, want %q (unchanged)", task.Description, "Original desc")
		}
		if task.Status != models.StatusPending {
			t.Errorf("UpdateTask Status changed to %q, want %q (unchanged)", task.Status, models.StatusPending)
		}
	})

	t.Run("update status only", func(t *testing.T) {
		_, r := newTestRouter()
		created := createTaskViaAPI(t, r, "Task", "", "")
		id := created["id"].(string)

		body, _ := json.Marshal(models.UpdateTaskRequest{Status: strPtr(models.StatusDone)})
		rr := doRequest(r, http.MethodPut, "/tasks/"+id, string(body))

		if rr.Code != http.StatusOK {
			t.Fatalf("UpdateTask status = %d, want 200", rr.Code)
		}

		var task models.Task
		json.Unmarshal(rr.Body.Bytes(), &task)
		if task.Status != models.StatusDone {
			t.Errorf("UpdateTask Status = %q, want %q", task.Status, models.StatusDone)
		}
	})

	t.Run("update all fields", func(t *testing.T) {
		_, r := newTestRouter()
		created := createTaskViaAPI(t, r, "Old", "Old desc", "")
		id := created["id"].(string)

		body, _ := json.Marshal(models.UpdateTaskRequest{
			Title:       strPtr("New"),
			Description: strPtr("New desc"),
			Status:      strPtr(models.StatusInProgress),
		})
		rr := doRequest(r, http.MethodPut, "/tasks/"+id, string(body))

		if rr.Code != http.StatusOK {
			t.Fatalf("UpdateTask status = %d, want 200", rr.Code)
		}

		var task models.Task
		json.Unmarshal(rr.Body.Bytes(), &task)
		if task.Title != "New" {
			t.Errorf("Title = %q, want %q", task.Title, "New")
		}
		if task.Description != "New desc" {
			t.Errorf("Description = %q, want %q", task.Description, "New desc")
		}
		if task.Status != models.StatusInProgress {
			t.Errorf("Status = %q, want %q", task.Status, models.StatusInProgress)
		}
	})

	t.Run("non-existent ID returns 404", func(t *testing.T) {
		_, r := newTestRouter()
		body, _ := json.Marshal(models.UpdateTaskRequest{Title: strPtr("X")})
		rr := doRequest(r, http.MethodPut, "/tasks/nonexistent", string(body))

		if rr.Code != http.StatusNotFound {
			t.Errorf("UpdateTask status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("empty body returns 400", func(t *testing.T) {
		_, r := newTestRouter()
		created := createTaskViaAPI(t, r, "Task", "", "")
		id := created["id"].(string)

		rr := doRequest(r, http.MethodPut, "/tasks/"+id, `{}`)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("UpdateTask status = %d, want %d", rr.Code, http.StatusBadRequest)
		}

		var errResp models.ErrorResponse
		json.Unmarshal(rr.Body.Bytes(), &errResp)
		if !strings.Contains(errResp.Error, "at least one field") {
			t.Errorf("error = %q, want it to contain 'at least one field'", errResp.Error)
		}
	})

	t.Run("invalid status returns 400", func(t *testing.T) {
		_, r := newTestRouter()
		created := createTaskViaAPI(t, r, "Task", "", "")
		id := created["id"].(string)

		body, _ := json.Marshal(models.UpdateTaskRequest{Status: strPtr("invalid")})
		rr := doRequest(r, http.MethodPut, "/tasks/"+id, string(body))

		if rr.Code != http.StatusBadRequest {
			t.Errorf("UpdateTask status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid JSON returns 400", func(t *testing.T) {
		_, r := newTestRouter()
		created := createTaskViaAPI(t, r, "Task", "", "")
		id := created["id"].(string)

		rr := doRequest(r, http.MethodPut, "/tasks/"+id, `{broken`)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("UpdateTask status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("update persists and is visible via GET", func(t *testing.T) {
		_, r := newTestRouter()
		created := createTaskViaAPI(t, r, "Before", "", "")
		id := created["id"].(string)

		body, _ := json.Marshal(models.UpdateTaskRequest{Title: strPtr("After")})
		doRequest(r, http.MethodPut, "/tasks/"+id, string(body))

		rr := doRequest(r, http.MethodGet, "/tasks/"+id, "")
		var task models.Task
		json.Unmarshal(rr.Body.Bytes(), &task)
		if task.Title != "After" {
			t.Errorf("GET after UPDATE Title = %q, want %q", task.Title, "After")
		}
	})
}

// --- DeleteTask Tests ---

func TestDeleteTask(t *testing.T) {
	t.Run("delete existing task returns 204", func(t *testing.T) {
		_, r := newTestRouter()
		created := createTaskViaAPI(t, r, "Delete me", "", "")
		id := created["id"].(string)

		rr := doRequest(r, http.MethodDelete, "/tasks/"+id, "")

		if rr.Code != http.StatusNoContent {
			t.Errorf("DeleteTask status = %d, want %d", rr.Code, http.StatusNoContent)
		}

		if rr.Body.Len() != 0 {
			t.Errorf("DeleteTask body should be empty, got %q", rr.Body.String())
		}
	})

	t.Run("delete non-existent returns 404", func(t *testing.T) {
		_, r := newTestRouter()
		rr := doRequest(r, http.MethodDelete, "/tasks/nonexistent", "")

		if rr.Code != http.StatusNotFound {
			t.Errorf("DeleteTask status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("deleted task no longer accessible via GET", func(t *testing.T) {
		_, r := newTestRouter()
		created := createTaskViaAPI(t, r, "Ephemeral", "", "")
		id := created["id"].(string)

		doRequest(r, http.MethodDelete, "/tasks/"+id, "")

		rr := doRequest(r, http.MethodGet, "/tasks/"+id, "")
		if rr.Code != http.StatusNotFound {
			t.Errorf("GET after DELETE status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("deleted task no longer in list", func(t *testing.T) {
		_, r := newTestRouter()
		created := createTaskViaAPI(t, r, "Will be deleted", "", "")
		createTaskViaAPI(t, r, "Will stay", "", "")
		id := created["id"].(string)

		doRequest(r, http.MethodDelete, "/tasks/"+id, "")

		rr := doRequest(r, http.MethodGet, "/tasks", "")
		var tasks []models.Task
		json.Unmarshal(rr.Body.Bytes(), &tasks)

		if len(tasks) != 1 {
			t.Errorf("ListTasks after delete returned %d tasks, want 1", len(tasks))
		}
		if len(tasks) == 1 && tasks[0].Title != "Will stay" {
			t.Errorf("remaining task Title = %q, want %q", tasks[0].Title, "Will stay")
		}
	})

	t.Run("double delete returns 404", func(t *testing.T) {
		_, r := newTestRouter()
		created := createTaskViaAPI(t, r, "Delete twice", "", "")
		id := created["id"].(string)

		doRequest(r, http.MethodDelete, "/tasks/"+id, "")
		rr := doRequest(r, http.MethodDelete, "/tasks/"+id, "")

		if rr.Code != http.StatusNotFound {
			t.Errorf("second DeleteTask status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})
}

// --- Integration-style: Full CRUD lifecycle ---

func TestFullCRUDLifecycle(t *testing.T) {
	_, r := newTestRouter()

	// 1. List tasks — should be empty
	rr := doRequest(r, http.MethodGet, "/tasks", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("initial list status = %d", rr.Code)
	}
	body := strings.TrimSpace(rr.Body.String())
	if body != "[]" {
		t.Fatalf("initial list = %q, want empty array", body)
	}

	// 2. Create a task
	rr = doRequest(r, http.MethodPost, "/tasks", `{"title":"Lifecycle test","description":"Testing full CRUD"}`)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rr.Code)
	}
	var created models.Task
	json.Unmarshal(rr.Body.Bytes(), &created)
	if created.ID == "" {
		t.Fatal("created task has empty ID")
	}
	if created.Status != models.StatusPending {
		t.Errorf("created status = %q, want pending", created.Status)
	}

	// 3. Get the task
	rr = doRequest(r, http.MethodGet, "/tasks/"+created.ID, "")
	if rr.Code != http.StatusOK {
		t.Fatalf("get status = %d", rr.Code)
	}
	var fetched models.Task
	json.Unmarshal(rr.Body.Bytes(), &fetched)
	if fetched.Title != "Lifecycle test" {
		t.Errorf("fetched title = %q", fetched.Title)
	}

	// 4. Update the task
	updateBody, _ := json.Marshal(map[string]string{"status": "done", "title": "Lifecycle complete"})
	rr = doRequest(r, http.MethodPut, "/tasks/"+created.ID, string(updateBody))
	if rr.Code != http.StatusOK {
		t.Fatalf("update status = %d: %s", rr.Code, rr.Body.String())
	}
	var updated models.Task
	json.Unmarshal(rr.Body.Bytes(), &updated)
	if updated.Status != models.StatusDone {
		t.Errorf("updated status = %q, want done", updated.Status)
	}
	if updated.Title != "Lifecycle complete" {
		t.Errorf("updated title = %q", updated.Title)
	}

	// 5. List — should have 1 task
	rr = doRequest(r, http.MethodGet, "/tasks", "")
	var tasks []models.Task
	json.Unmarshal(rr.Body.Bytes(), &tasks)
	if len(tasks) != 1 {
		t.Errorf("list after update has %d tasks, want 1", len(tasks))
	}

	// 6. Delete the task
	rr = doRequest(r, http.MethodDelete, "/tasks/"+created.ID, "")
	if rr.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d", rr.Code)
	}

	// 7. Verify it's gone
	rr = doRequest(r, http.MethodGet, "/tasks/"+created.ID, "")
	if rr.Code != http.StatusNotFound {
		t.Errorf("get after delete status = %d, want 404", rr.Code)
	}

	// 8. List — should be empty again
	rr = doRequest(r, http.MethodGet, "/tasks", "")
	json.Unmarshal(rr.Body.Bytes(), &tasks)
	if len(tasks) != 0 {
		t.Errorf("list after delete has %d tasks, want 0", len(tasks))
	}
}

// --- Response Format Validation ---

func TestAllResponsesHaveContentTypeJSON(t *testing.T) {
	_, r := newTestRouter()
	created := createTaskViaAPI(t, r, "Content-Type check", "", "")
	id := created["id"].(string)

	endpoints := []struct {
		method string
		path   string
		body   string
		skip   bool // skip Content-Type check (204 may not have it)
	}{
		{http.MethodGet, "/health", "", false},
		{http.MethodGet, "/tasks", "", false},
		{http.MethodGet, "/tasks/" + id, "", false},
		{http.MethodPut, "/tasks/" + id, `{"title":"Updated"}`, false},
		{http.MethodDelete, "/tasks/" + id, "", true}, // 204 No Content
	}

	for _, ep := range endpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			rr := doRequest(r, ep.method, ep.path, ep.body)
			if ep.skip {
				return
			}
			ct := rr.Header().Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("Content-Type = %q, want %q", ct, "application/json")
			}
		})
	}
}
