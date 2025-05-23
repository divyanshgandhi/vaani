package main

import (
	"bytes"
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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

const (
	// DefaultPort is the default port the service runs on
	DefaultPort = "8082"
	// Default Bulbul API URL
	DefaultBulbulAPIURL = "https://api.bulbul.ai/v1/tts"
)

// TTSRequest represents a text-to-speech request
type TTSRequest struct {
	Text     string  `json:"text"`
	VoiceID  string  `json:"voice_id"`
	Speed    float64 `json:"speed,omitempty"`
	Pitch    float64 `json:"pitch,omitempty"`
	Emotion  int     `json:"emotion,omitempty"`
	Language string  `json:"language"`
}

// TTSResponse represents a text-to-speech response
type TTSResponse struct {
	AudioURL string  `json:"audio_url"`
	Duration float64 `json:"duration,omitempty"`
	Format   string  `json:"format,omitempty"`
}

func main() {
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

	// TTS endpoints
	r.Route("/api/v1", func(r chi.Router) {
		// Generate TTS
		r.Post("/tts", handleTTS(logger))

		// Get supported voices
		r.Get("/voices", handleGetVoices(logger))
	})

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
		logger.Info("Starting bulbul adapter service", zap.String("port", port))
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

// handleTTS handles text-to-speech requests
func handleTTS(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request
		var request TTSRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate request
		if request.Text == "" {
			http.Error(w, "Text is required", http.StatusBadRequest)
			return
		}

		// Log request
		logger.Info("Received TTS request",
			zap.String("language", request.Language),
			zap.String("voice_id", request.VoiceID),
			zap.Int("text_length", len(request.Text)),
		)

		// Get Sarvam API key from environment
		apiKey := os.Getenv("SARVAM_API_KEY")
		if apiKey == "" {
			logger.Warn("SARVAM_API_KEY not set, using mock response")
			// Return mock response for development
			audioURL := fmt.Sprintf("https://storage.vaani.app/media/%s_%s.mp3",
				request.Language, time.Now().Format("20060102150405"))

			response := TTSResponse{
				AudioURL: audioURL,
				Duration: float64(len(request.Text)) / 20.0,
				Format:   "mp3",
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				logger.Error("Failed to encode response", zap.Error(err))
			}
			return
		}

		// Call Sarvam Bulbul API
		response, err := callSarvamAPI(request, apiKey, logger)
		if err != nil {
			logger.Error("Failed to call Sarvam API", zap.Error(err))
			http.Error(w, "Failed to generate TTS", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error("Failed to encode response", zap.Error(err))
		}
	}
}

// callSarvamAPI calls the Sarvam Bulbul API
func callSarvamAPI(request TTSRequest, apiKey string, logger *zap.Logger) (*TTSResponse, error) {
	// Prepare Sarvam API request
	sarvamRequest := map[string]interface{}{
		"text":     request.Text,
		"language": request.Language,
		"speaker":  request.VoiceID,
	}

	// Add emotion/style if provided
	if request.Emotion > 0 {
		// Map emotion (0-100) to Sarvam style
		style := "neutral"
		if request.Emotion > 70 {
			style = "excited"
		} else if request.Emotion > 40 {
			style = "happy"
		}
		sarvamRequest["style"] = style
	}

	// Marshal request
	requestBody, err := json.Marshal(sarvamRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", DefaultBulbulAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Sarvam-Key", apiKey)

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var sarvamResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&sarvamResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Sarvam API error: %v", sarvamResponse)
	}

	// Extract audio URL from response
	audioURL, ok := sarvamResponse["audio_url"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid response format: missing audio_url")
	}

	// Extract duration if available
	duration := float64(len(request.Text)) / 20.0 // Default estimate
	if d, ok := sarvamResponse["duration"].(float64); ok {
		duration = d
	}

	return &TTSResponse{
		AudioURL: audioURL,
		Duration: duration,
		Format:   "mp3",
	}, nil
}

// handleGetVoices returns a list of supported voices
func handleGetVoices(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// In a real implementation, this would dynamically fetch from Bulbul API
		// For now, we'll return a hardcoded list of Indic voices
		voices := []map[string]interface{}{
			{
				"id":          "hi_female_1",
				"name":        "Priya",
				"language":    "hi",
				"gender":      "female",
				"preview_url": "https://storage.vaani.app/samples/hi_female_1.mp3",
			},
			{
				"id":          "hi_male_1",
				"name":        "Rahul",
				"language":    "hi",
				"gender":      "male",
				"preview_url": "https://storage.vaani.app/samples/hi_male_1.mp3",
			},
			{
				"id":          "ta_female_1",
				"name":        "Anjali",
				"language":    "ta",
				"gender":      "female",
				"preview_url": "https://storage.vaani.app/samples/ta_female_1.mp3",
			},
			{
				"id":          "bn_female_1",
				"name":        "Meera",
				"language":    "bn",
				"gender":      "female",
				"preview_url": "https://storage.vaani.app/samples/bn_female_1.mp3",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(voices); err != nil {
			logger.Error("Failed to encode response", zap.Error(err))
		}
	}
}
