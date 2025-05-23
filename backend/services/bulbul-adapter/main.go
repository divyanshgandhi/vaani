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
	// Sarvam TTS API URL - Updated to use correct Sarvam endpoint
	SarvamTTSAPIURL = "https://api.sarvam.ai/text-to-speech"
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

// SarvamTTSRequest represents the request format for Sarvam TTS API
type SarvamTTSRequest struct {
	Inputs              []string `json:"inputs"`
	TargetLanguageCode  string   `json:"target_language_code"`
	Speaker             string   `json:"speaker,omitempty"`
	Pitch               float64  `json:"pitch,omitempty"`
	Pace                float64  `json:"pace,omitempty"`
	Loudness            float64  `json:"loudness,omitempty"`
	SpeechSampleRate    int      `json:"speech_sample_rate,omitempty"`
	EnablePreprocessing bool     `json:"enable_preprocessing"`
	Model               string   `json:"model"`
}

// SarvamTTSResponse represents the response format from Sarvam TTS API
type SarvamTTSResponse struct {
	Audios []string `json:"audios"`
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

		// Call Sarvam TTS API
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

// callSarvamAPI calls the Sarvam TTS API with correct format
func callSarvamAPI(request TTSRequest, apiKey string, logger *zap.Logger) (*TTSResponse, error) {
	// Map our voice IDs to Sarvam speakers
	speaker := mapVoiceToSarvamSpeaker(request.VoiceID)

	// Map our language codes to Sarvam language codes
	languageCode := mapLanguageToSarvamCode(request.Language)

	// Prepare Sarvam API request according to documentation
	sarvamRequest := SarvamTTSRequest{
		Inputs:              []string{request.Text},
		TargetLanguageCode:  languageCode,
		Speaker:             speaker,
		Pitch:               request.Pitch,
		Pace:                request.Speed,
		Loudness:            1.5,   // Default loudness
		SpeechSampleRate:    22050, // Default sample rate
		EnablePreprocessing: true,
		Model:               "bulbul:v1",
	}

	// Marshal request
	requestBody, err := json.Marshal(sarvamRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	logger.Debug("Sarvam API request", zap.String("body", string(requestBody)))

	// Create HTTP request
	req, err := http.NewRequest("POST", SarvamTTSAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set correct headers according to Sarvam documentation
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-subscription-key", apiKey)

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var sarvamResponse SarvamTTSResponse
	if err := json.NewDecoder(resp.Body).Decode(&sarvamResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Sarvam API error (status %d): %v", resp.StatusCode, sarvamResponse)
	}

	// Validate response
	if len(sarvamResponse.Audios) == 0 {
		return nil, fmt.Errorf("no audio generated")
	}

	// The audio is base64 encoded, we need to save it and return URL
	// For now, return the base64 data directly (in production, save to storage)
	audioURL := fmt.Sprintf("data:audio/wav;base64,%s", sarvamResponse.Audios[0])

	// Estimate duration (characters / 20 = rough seconds)
	duration := float64(len(request.Text)) / 20.0

	return &TTSResponse{
		AudioURL: audioURL,
		Duration: duration,
		Format:   "wav",
	}, nil
}

// mapVoiceToSarvamSpeaker maps our voice IDs to Sarvam speaker names
func mapVoiceToSarvamSpeaker(voiceID string) string {
	switch voiceID {
	case "hi_female_1":
		return "meera"
	case "hi_male_1":
		return "arvind"
	case "ta_female_1":
		return "pavithra"
	case "bn_female_1":
		return "maitreyi"
	default:
		return "meera" // Default to meera
	}
}

// mapLanguageToSarvamCode maps our language codes to Sarvam language codes
func mapLanguageToSarvamCode(language string) string {
	switch language {
	case "hi":
		return "hi-IN"
	case "ta":
		return "ta-IN"
	case "bn":
		return "bn-IN"
	case "en":
		return "en-IN"
	case "gu":
		return "gu-IN"
	case "kn":
		return "kn-IN"
	case "ml":
		return "ml-IN"
	case "mr":
		return "mr-IN"
	case "od":
		return "od-IN"
	case "pa":
		return "pa-IN"
	case "te":
		return "te-IN"
	default:
		return "hi-IN" // Default to Hindi
	}
}

// handleGetVoices returns a list of supported voices
func handleGetVoices(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Return voices compatible with Sarvam TTS
		voices := []map[string]interface{}{
			{
				"id":          "hi_female_1",
				"name":        "Meera",
				"language":    "hi",
				"gender":      "female",
				"preview_url": "https://storage.vaani.app/samples/hi_female_1.mp3",
			},
			{
				"id":          "hi_male_1",
				"name":        "Arvind",
				"language":    "hi",
				"gender":      "male",
				"preview_url": "https://storage.vaani.app/samples/hi_male_1.mp3",
			},
			{
				"id":          "ta_female_1",
				"name":        "Pavithra",
				"language":    "ta",
				"gender":      "female",
				"preview_url": "https://storage.vaani.app/samples/ta_female_1.mp3",
			},
			{
				"id":          "bn_female_1",
				"name":        "Maitreyi",
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
