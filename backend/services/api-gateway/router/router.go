package router

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/NavoDayAI/vaani/backend/services/api-gateway/config"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/database"
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

	// Public routes
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Language detection endpoint
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

	// Add secured routes and other endpoints here
	// This is where you would add authentication middleware and protected routes

	return r
}
