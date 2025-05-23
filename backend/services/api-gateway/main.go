package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NavoDayAI/vaani/backend/services/api-gateway/config"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/database"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/handlers"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/router"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/services"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/utils/langdetect"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load .env file if present
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file:", err)
	}

	// Initialize logging
	logger, err := createLogger()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Load configuration
	cfg := config.New()
	logger.Info("Config loaded",
		zap.String("port", cfg.Port),
		zap.String("env", cfg.Env),
		zap.String("firebase_project_id", cfg.Firebase.ProjectID))

	var firestoreClient *database.FirestoreClient
	var jobRepository *database.FirestoreJobRepository
	var projectRepository *database.FirestoreProjectRepository
	var ttsClient *services.TTSClient
	var jobProcessor *services.JobProcessor

	// Try to initialize Firestore client
	firestoreClient, err = database.NewFirestoreClient(cfg, logger)
	if err != nil {
		logger.Warn("Failed to initialize Firestore client, proceeding with limited functionality",
			zap.Error(err))
	} else {
		defer firestoreClient.Close()
		// Initialize repositories
		jobRepository = database.NewFirestoreJobRepository(firestoreClient, logger)
		projectRepository = database.NewFirestoreProjectRepository(firestoreClient, logger)
		logger.Info("Firestore client and repositories initialized successfully")

		// Initialize TTS client
		bulbulURL := getEnv("BULBUL_ADAPTER_URL", "http://localhost:8082")
		orpheusURL := getEnv("ORPHEUS_SERVICE_URL", "http://localhost:8081")
		mediaURL := getEnv("MEDIA_SERVICE_URL", "http://localhost:8083")

		ttsClient = services.NewTTSClient(bulbulURL, orpheusURL, mediaURL, logger)
		logger.Info("TTS client initialized",
			zap.String("bulbul_url", bulbulURL),
			zap.String("orpheus_url", orpheusURL),
			zap.String("media_url", mediaURL))

		// Initialize job processor
		jobProcessor = services.NewJobProcessor(jobRepository, ttsClient, logger)

		// Start job processor in background
		go jobProcessor.Start()
		logger.Info("Job processor started")
	}

	// Set up the router with or without Firebase support
	var r http.Handler
	if firestoreClient != nil && jobRepository != nil && projectRepository != nil {
		// Use full router with auth and repositories
		r = router.Setup(cfg, logger, jobRepository, projectRepository)
		logger.Info("Using full router configuration with authentication")
	} else {
		// Use limited router without Firebase dependencies
		r = setupBasicRouter(logger)
		logger.Info("Using basic router configuration without authentication")
	}

	// Create the server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		logger.Info("Starting server", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-quit
	logger.Info("Shutting down server", zap.String("signal", sig.String()))

	// Stop job processor if it was started
	if jobProcessor != nil {
		jobProcessor.Stop()
		logger.Info("Job processor stopped")
	}

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server stopped")
}

// setupBasicRouter creates a simple router without Firebase dependencies
// This is used as a fallback if Firebase initialization fails
func setupBasicRouter(logger *zap.Logger) http.Handler {
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

	// Auth endpoints
	r.Post("/v1/auth/send-otp", handlers.SendOTPHandler(logger))
	r.Post("/v1/auth/verify-otp", handlers.VerifyOTPHandler(logger))

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

	return r
}

func createLogger() (*zap.Logger, error) {
	env := os.Getenv("APP_ENV")

	if env == "production" {
		return zap.NewProduction()
	} else {
		return zap.NewDevelopment()
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
