package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"weather-cache-api/internal/cache"
	"weather-cache-api/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	c := cache.NewCache()
	h := handlers.NewHandler(c)

	// Background goroutine: clean expired entries every 10 seconds.
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			c.CleanExpired()
		}
	}()

	mux := http.NewServeMux()

	// Route /health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.HealthCheck(w, r)
	})

	// Route /cache and /cache/{key}
	mux.HandleFunc("/cache", cacheRouter(h))
	mux.HandleFunc("/cache/", cacheRouter(h))

	log.Printf("weather-cache-api starting on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// cacheRouter dispatches /cache requests based on method and path.
func cacheRouter(h *handlers.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSuffix(r.URL.Path, "/")

		// Exact /cache path → list or create
		if path == "/cache" {
			switch r.Method {
			case http.MethodGet:
				h.ListCache(w, r)
			case http.MethodPost:
				h.PostCache(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// /cache/{key} path → get or delete by key
		switch r.Method {
		case http.MethodGet:
			h.GetCacheByKey(w, r)
		case http.MethodDelete:
			h.DeleteCacheByKey(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
