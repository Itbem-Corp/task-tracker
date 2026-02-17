package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/itbem-corp/task-tracker/internal/models"
	"github.com/itbem-corp/task-tracker/internal/store"
)

// RegisterRoutes mounts task CRUD endpoints on the given router.
func RegisterRoutes(r chi.Router, s *store.TaskStore) {
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", listTasks(s))
		r.Post("/", createTask(s))
		r.Get("/{id}", getTask(s))
		r.Put("/{id}", updateTask(s))
		r.Delete("/{id}", deleteTask(s))
	})

	r.Get("/health", healthCheck)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func listTasks(s *store.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks := s.GetAll()
		writeJSON(w, http.StatusOK, tasks)
	}
}

func createTask(s *store.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if err := req.Validate(); err != nil {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		task := s.Create(req.Title, req.Description, req.Status)
		writeJSON(w, http.StatusCreated, task)
	}
}

func getTask(s *store.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		task, err := s.GetByID(id)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, task)
	}
}

func updateTask(s *store.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var req models.UpdateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if err := req.Validate(); err != nil {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		task, err := s.Update(id, req)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, task)
	}
}

func deleteTask(s *store.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if err := s.Delete(id); err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, models.ErrorResponse{Error: message})
}
