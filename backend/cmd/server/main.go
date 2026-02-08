package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/meixg/podcast-reader/backend/internal/handlers"
	"github.com/meixg/podcast-reader/backend/internal/services"
	"github.com/meixg/podcast-reader/backend/pkg/scanner"
)

func main() {
	// Get downloads directory from environment or use default
	downloadsDir := os.Getenv("DOWNLOADS_DIR")
	if downloadsDir == "" {
		// Use relative path to downloads directory
		downloadsDir = "downloads"
	}

	// Initialize services
	episodeScanner := scanner.NewScanner(downloadsDir)
	episodeService := services.NewEpisodeService(episodeScanner)
	taskService := services.NewTaskService()
	downloadService := services.NewDownloadService(downloadsDir, taskService)

	// Set download service for task service
	taskService.SetDownloadService(downloadService)

	// Initialize handlers
	episodeHandler := handlers.NewEpisodeHandler(episodeService)
	taskHandler := handlers.NewTaskHandler(taskService)

	// Setup routes
	mux := http.NewServeMux()

	// Episode routes
	mux.HandleFunc("/api/episodes", episodeHandler.GetEpisodes)
	mux.HandleFunc("/api/episodes/", episodeHandler.GetShowNotes)

	// Task routes
	mux.HandleFunc("/api/tasks", taskHandler.HandleTasks)

	// Static file server for downloads (images, audio files)
	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Wrap with CORS middleware
	handler := corsMiddleware(mux)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	log.Printf("Scanning downloads from: %s", downloadsDir)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

// corsMiddleware adds CORS headers to all responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
