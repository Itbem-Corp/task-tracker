package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Employee represents a company employee.
type Employee struct {
	ID         int    `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Department string `json:"department"`
	Position   string `json:"position"`
}

var employees = []Employee{
	{
		ID:         1,
		FirstName:  "Alice",
		LastName:   "Chen",
		Email:      "alice.chen@example.com",
		Department: "Engineering",
		Position:   "Senior Software Engineer",
	},
	{
		ID:         2,
		FirstName:  "Brian",
		LastName:   "Patel",
		Email:      "brian.patel@example.com",
		Department: "Marketing",
		Position:   "Marketing Manager",
	},
	{
		ID:         3,
		FirstName:  "Carmen",
		LastName:   "Rivera",
		Email:      "carmen.rivera@example.com",
		Department: "HR",
		Position:   "HR Business Partner",
	},
	{
		ID:         4,
		FirstName:  "David",
		LastName:   "Okonkwo",
		Email:      "david.okonkwo@example.com",
		Department: "Finance",
		Position:   "Financial Analyst",
	},
	{
		ID:         5,
		FirstName:  "Elena",
		LastName:   "Johansson",
		Email:      "elena.johansson@example.com",
		Department: "Design",
		Position:   "Lead UX Designer",
	},
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/employees", handleGetEmployees)
	mux.HandleFunc("GET /api/employees/{id}", handleGetEmployeeByID)
	mux.HandleFunc("GET /health", handleHealth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	server := &http.Server{
		Addr:         addr,
		Handler:      withLogging(withMethodNotAllowed(mux)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("starting server", "addr", addr)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

// handleGetEmployees returns the full list of employees as a JSON array.
func handleGetEmployees(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, employees)
}

// handleGetEmployeeByID returns a single employee matching the given ID.
func handleGetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": fmt.Sprintf("invalid employee id: %s", idStr),
		})
		return
	}

	for _, emp := range employees {
		if emp.ID == id {
			writeJSON(w, http.StatusOK, emp)
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{
		"error": fmt.Sprintf("employee with id %d not found", id),
	})
}

// handleHealth returns a simple health check response.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// writeJSON encodes v as JSON and writes it to the response with the given
// status code and an application/json Content-Type header.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

// withLogging wraps an http.Handler and logs every request.
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"duration_ms", time.Since(start).Milliseconds(),
			"remote_addr", r.RemoteAddr,
		)
	})
}

// withMethodNotAllowed wraps a ServeMux so that unregistered method/path
// combinations that match a known path return 405 instead of the default 404.
// For completely unknown paths the mux's own 404 is used.
func withMethodNotAllowed(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Let the mux handle the request normally. Go 1.22's method-aware
		// routing already returns 405 for method mismatches on registered
		// patterns, so we just ensure the response body is JSON.
		rw := &methodNotAllowedWriter{ResponseWriter: w}
		mux.ServeHTTP(rw, r)
		if rw.wroteMethodNotAllowed {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
				"error": fmt.Sprintf("method %s not allowed", r.Method),
			})
		}
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// methodNotAllowedWriter detects when the mux writes a 405 so we can replace
// the body with a JSON error message.
type methodNotAllowedWriter struct {
	http.ResponseWriter
	wroteMethodNotAllowed bool
	headerWritten         bool
}

func (w *methodNotAllowedWriter) WriteHeader(code int) {
	if code == http.StatusMethodNotAllowed {
		w.wroteMethodNotAllowed = true
		// Don't forward; we will write our own JSON response.
		return
	}
	w.headerWritten = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *methodNotAllowedWriter) Write(b []byte) (int, error) {
	if w.wroteMethodNotAllowed {
		// Swallow the default plain-text body.
		return len(b), nil
	}
	return w.ResponseWriter.Write(b)
}
