package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"weather-cache-api/internal/cache"
)

func TestPostCache(t *testing.T) {
	c := cache.NewCache()
	h := NewHandler(c)

	body := map[string]interface{}{
		"key":   "test-key",
		"value": "test-value",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/cache", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	h.PostCache(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["status"] != "ok" || resp["key"] != "test-key" {
		t.Errorf("Unexpected response: %+v", resp)
	}
}

func TestPostCacheWithTTL(t *testing.T) {
	c := cache.NewCache()
	h := NewHandler(c)

	body := map[string]interface{}{
		"key":         "ttl-key",
		"value":       "ttl-value",
		"ttl_seconds": 1,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/cache", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	h.PostCache(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	value, ok := c.Get("ttl-key")
	if !ok || value != "ttl-value" {
		t.Errorf("Expected value not found")
	}
}

func TestPostCacheMissingKey(t *testing.T) {
	c := cache.NewCache()
	h := NewHandler(c)

	body := map[string]interface{}{
		"value": "test-value",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/cache", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	h.PostCache(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestPostCacheMissingValue(t *testing.T) {
	c := cache.NewCache()
	h := NewHandler(c)

	body := map[string]interface{}{
		"key": "test-key",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/cache", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	h.PostCache(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestPostCacheInvalidJSON(t *testing.T) {
	c := cache.NewCache()
	h := NewHandler(c)

	req := httptest.NewRequest(http.MethodPost, "/cache", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	h.PostCache(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetCacheByKey(t *testing.T) {
	c := cache.NewCache()
	h := NewHandler(c)
	c.Set("test-key", "test-value", nil)

	req := httptest.NewRequest(http.MethodGet, "/cache/test-key", nil)
	w := httptest.NewRecorder()

	h.GetCacheByKey(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["key"] != "test-key" || resp["value"] != "test-value" {
		t.Errorf("Unexpected response: %+v", resp)
	}
}

func TestGetCacheByKeyNotFound(t *testing.T) {
	c := cache.NewCache()
	h := NewHandler(c)

	req := httptest.NewRequest(http.MethodGet, "/cache/nonexistent", nil)
	w := httptest.NewRecorder()

	h.GetCacheByKey(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestGetCacheByKeyExpired(t *testing.T) {
	c := cache.NewCache()
	h := NewHandler(c)

	ttl := 50 * time.Millisecond
	c.Set("expired-key", "expired-value", &ttl)
	time.Sleep(100 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/cache/expired-key", nil)
	w := httptest.NewRecorder()

	h.GetCacheByKey(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for expired key, got %d", w.Code)
	}
}

func TestDeleteCacheByKey(t *testing.T) {
	c := cache.NewCache()
	h := NewHandler(c)
	c.Set("delete-key", "delete-value", nil)

	req := httptest.NewRequest(http.MethodDelete, "/cache/delete-key", nil)
	w := httptest.NewRecorder()

	h.DeleteCacheByKey(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["status"] != "deleted" {
		t.Errorf("Unexpected response: %+v", resp)
	}

	_, ok := c.Get("delete-key")
	if ok {
		t.Error("Key should be deleted")
	}
}

func TestDeleteCacheByKeyNotFound(t *testing.T) {
	c := cache.NewCache()
	h := NewHandler(c)

	req := httptest.NewRequest(http.MethodDelete, "/cache/nonexistent", nil)
	w := httptest.NewRecorder()

	h.DeleteCacheByKey(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestListCache(t *testing.T) {
	c := cache.NewCache()
	h := NewHandler(c)

	c.Set("key1", "value1", nil)
	c.Set("key2", "value2", nil)

	req := httptest.NewRequest(http.MethodGet, "/cache", nil)
	w := httptest.NewRecorder()

	h.ListCache(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	keys := resp["keys"].([]interface{})

	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}
}

func TestListCacheWithTTL(t *testing.T) {
	c := cache.NewCache()
	h := NewHandler(c)

	c.Set("no-ttl", "value", nil)
	ttl := 1 * time.Hour
	c.Set("with-ttl", "value", &ttl)

	req := httptest.NewRequest(http.MethodGet, "/cache", nil)
	w := httptest.NewRecorder()

	h.ListCache(w, req)

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	keys := resp["keys"].([]interface{})

	noTTL := keys[0].(map[string]interface{})
	withTTL := keys[1].(map[string]interface{})

	if noTTL["expires_at"] != nil {
		t.Error("Key without TTL should have nil expires_at")
	}

	if withTTL["expires_at"] == nil {
		t.Error("Key with TTL should have non-nil expires_at")
	}
}

func TestListCacheEmpty(t *testing.T) {
	c := cache.NewCache()
	h := NewHandler(c)

	req := httptest.NewRequest(http.MethodGet, "/cache", nil)
	w := httptest.NewRecorder()

	h.ListCache(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	keys := resp["keys"].([]interface{})

	if len(keys) != 0 {
		t.Errorf("Expected 0 keys, got %d", len(keys))
	}
}

func TestHealthCheck(t *testing.T) {
	c := cache.NewCache()
	h := NewHandler(c)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	h.HealthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", resp["status"])
	}
}
