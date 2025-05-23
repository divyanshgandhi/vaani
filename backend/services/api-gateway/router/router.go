package router

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/NavoDayAI/vaani/backend/services/api-gateway/config"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/database"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/handlers"
	authmw "github.com/NavoDayAI/vaani/backend/services/api-gateway/middleware"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/utils/langdetect"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

// Setup initializes and returns the full application router
func Setup(cfg *config.AppConfig, logger *zap.Logger, jobRepository *database.FirestoreJobRepository) http.Handler {
	r := chi.NewRouter()

	// Middleware
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

	// Public routes (no authentication required)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Language detection endpoint (public for now, can be moved to protected later)
	r.Post("/v1/detect-language", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Text string `json:"text"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		lang, err := langdetect.Detect(request.Text)
		if err != nil {
			if errors.Is(err, langdetect.ErrTextTooShort) {
				http.Error(w, "Text too short for reliable detection", http.StatusBadRequest)
				return
			}

			logger.Error("Language detection failed", zap.Error(err))
			http.Error(w, "Failed to detect language", http.StatusInternalServerError)
			return
		}

		isIndic := langdetect.IsIndic(lang)

		response := struct {
			Language string `json:"language"`
			IsIndic  bool   `json:"is_indic"`
		}{
			Language: lang,
			IsIndic:  isIndic,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Protected routes (require authentication)
	r.Group(func(r chi.Router) {
		// Add Firebase authentication middleware
		r.Use(authmw.FirebaseAuthMiddleware(cfg, logger))

		// Add rate limiting middleware if available
		// r.Use(authmw.RateLimitMiddleware(cfg, logger))

		// TTS Generation endpoint
		r.Post("/v1/generate", handlers.GenerateHandler(logger, jobRepository))

		// Job status endpoint
		r.Get("/v1/jobs/{jobID}", handlers.GetJobHandler(logger, jobRepository))

		// Download endpoint (pre-signed URL redirect)
		r.Get("/v1/download/{jobID}", handlers.DownloadHandler(logger, jobRepository))

		// Voice cloning endpoint - placeholder for now
		r.Post("/v1/clone", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte(`{"error": "not_implemented", "message": "Voice cloning endpoint not yet implemented"}`))
		})

		// WebSocket preview endpoint
		r.Get("/v1/preview", handlers.PreviewHandler(logger))

		// User profile and quota endpoints - placeholders for now
		r.Get("/v1/profile", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte(`{"error": "not_implemented", "message": "Profile endpoint not yet implemented"}`))
		})

		r.Get("/v1/quota", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte(`{"error": "not_implemented", "message": "Quota endpoint not yet implemented"}`))
		})

		// Voice listing endpoint - placeholder for now
		r.Get("/v1/voices", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte(`{"error": "not_implemented", "message": "Voices endpoint not yet implemented"}`))
		})
	})

	return r
}
