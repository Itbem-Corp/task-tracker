package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"weather-cache-api/internal/cache"
)

// Handler holds dependencies for HTTP endpoint handlers.
type Handler struct {
	cache *cache.Cache
}

// NewHandler creates a Handler with the given cache.
func NewHandler(c *cache.Cache) *Handler {
	return &Handler{cache: c}
}

// postCacheRequest is the expected JSON body for POST /cache.
type postCacheRequest struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	TTLSeconds *int   `json:"ttl_seconds,omitempty"`
}

// PostCache handles POST /cache - stores a key-value pair with optional TTL.
func (h *Handler) PostCache(w http.ResponseWriter, r *http.Request) {
	var req postCacheRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	if req.Key == "" || req.Value == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "key and value are required"})
		return
	}

	var ttl *time.Duration
	if req.TTLSeconds != nil {
		d := time.Duration(*req.TTLSeconds) * time.Second
		ttl = &d
	}

	h.cache.Set(req.Key, req.Value, ttl)

	writeJSON(w, http.StatusCreated, map[string]string{
		"status": "ok",
		"key":    req.Key,
	})
}

// GetCacheByKey handles GET /cache/{key} - retrieves a cached value.
func (h *Handler) GetCacheByKey(w http.ResponseWriter, r *http.Request) {
	key := extractKey(r.URL.Path)
	if key == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "key is required"})
		return
	}

	value, ok := h.cache.Get(key)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "key not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"key":   key,
		"value": value,
	})
}

// DeleteCacheByKey handles DELETE /cache/{key} - removes a cached key.
func (h *Handler) DeleteCacheByKey(w http.ResponseWriter, r *http.Request) {
	key := extractKey(r.URL.Path)
	if key == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "key is required"})
		return
	}

	if !h.cache.Delete(key) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "key not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "deleted",
		"key":    key,
	})
}

// keyInfoResponse is the JSON representation of a key in the list response.
type keyInfoResponse struct {
	Key       string  `json:"key"`
	ExpiresAt *string `json:"expires_at"`
}

// ListCache handles GET /cache - returns all non-expired keys.
func (h *Handler) ListCache(w http.ResponseWriter, r *http.Request) {
	keys := h.cache.Keys()

	items := make([]keyInfoResponse, len(keys))
	for i, ki := range keys {
		items[i] = keyInfoResponse{Key: ki.Key}
		if ki.ExpiresAt != nil {
			s := ki.ExpiresAt.Format(time.RFC3339)
			items[i].ExpiresAt = &s
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"keys": items,
	})
}

// HealthCheck handles GET /health - returns service health status.
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// extractKey parses the cache key from a URL path like /cache/mykey.
func extractKey(path string) string {
	trimmed := strings.TrimPrefix(path, "/cache/")
	if trimmed == path || trimmed == "" {
		return ""
	}
	return trimmed
}

// writeJSON encodes v as JSON and writes it to w with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
