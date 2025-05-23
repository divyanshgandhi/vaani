package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

const (
	// DataDir is the directory where media files are stored
	DataDir = "/app/data"
	// DefaultPort is the default port the service runs on
	DefaultPort = "8083"
)

func main() {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(DataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Create router and add middleware
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Media file endpoints
	r.Route("/api/v1", func(r chi.Router) {
		// Upload file
		r.Post("/upload", handleFileUpload(logger))

		// Get file by ID
		r.Get("/media/{fileID}", handleGetFile(logger))

		// Delete file by ID
		r.Delete("/media/{fileID}", handleDeleteFile(logger))
	})

	// Serve static files
	fileServer := http.FileServer(http.Dir(DataDir))
	r.Get("/media/*", http.StripPrefix("/media", fileServer).ServeHTTP)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting media service", zap.String("port", port))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Graceful shutdown
	logger.Info("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Failed to shutdown server gracefully", zap.Error(err))
	}

	logger.Info("Server stopped")
}

// handleFileUpload handles file upload requests
func handleFileUpload(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get file from form
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "No file provided", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Generate unique filename
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
		filepath := filepath.Join(DataDir, filename)

		// Create new file
		dst, err := os.Create(filepath)
		if err != nil {
			logger.Error("Failed to create file", zap.Error(err))
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy file contents
		if _, err := io.Copy(dst, file); err != nil {
			logger.Error("Failed to copy file", zap.Error(err))
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		// Return file URL
		response := map[string]string{
			"file_id":  filename,
			"file_url": fmt.Sprintf("/media/%s", filename),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error("Failed to encode response", zap.Error(err))
		}
	}
}

// handleGetFile handles file retrieval requests
func handleGetFile(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileID := chi.URLParam(r, "fileID")
		if fileID == "" {
			http.Error(w, "File ID is required", http.StatusBadRequest)
			return
		}

		// Check if file exists
		filepath := filepath.Join(DataDir, fileID)
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// Serve the file
		http.ServeFile(w, r, filepath)
	}
}

// handleDeleteFile handles file deletion requests
func handleDeleteFile(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileID := chi.URLParam(r, "fileID")
		if fileID == "" {
			http.Error(w, "File ID is required", http.StatusBadRequest)
			return
		}

		// Check if file exists
		filepath := filepath.Join(DataDir, fileID)
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// Delete the file
		if err := os.Remove(filepath); err != nil {
			logger.Error("Failed to delete file", zap.Error(err), zap.String("file", fileID))
			http.Error(w, "Failed to delete file", http.StatusInternalServerError)
			return
		}

		// Return success response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}
