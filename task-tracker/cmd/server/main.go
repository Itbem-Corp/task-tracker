package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/itbem-corp/task-tracker/internal/handlers"
	"github.com/itbem-corp/task-tracker/internal/middleware"
	"github.com/itbem-corp/task-tracker/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	taskStore := store.New()

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.Logger)
	r.Use(chimw.Recoverer)

	handlers.RegisterRoutes(r, taskStore)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("task-tracker listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
