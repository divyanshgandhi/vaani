package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// TTSClient handles communication with TTS services
type TTSClient struct {
	bulbulURL  string
	orpheusURL string
	mediaURL   string
	httpClient *http.Client
	logger     *zap.Logger
}

// TTSRequest represents a request to TTS services
type TTSRequest struct {
	Text     string  `json:"text"`
	VoiceID  string  `json:"voice_id"`
	Speed    float64 `json:"speed,omitempty"`
	Pitch    float64 `json:"pitch,omitempty"`
	Emotion  int     `json:"emotion,omitempty"`
	Language string  `json:"language"`
	Format   string  `json:"format,omitempty"`
}

// TTSResponse represents a response from TTS services
type TTSResponse struct {
	AudioURL string  `json:"audio_url"`
	Duration float64 `json:"duration,omitempty"`
	Format   string  `json:"format,omitempty"`
	Error    string  `json:"error,omitempty"`
}

// NewTTSClient creates a new TTS client
func NewTTSClient(bulbulURL, orpheusURL, mediaURL string, logger *zap.Logger) *TTSClient {
	return &TTSClient{
		bulbulURL:  bulbulURL,
		orpheusURL: orpheusURL,
		mediaURL:   mediaURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// GenerateTTS generates TTS audio using the appropriate service based on language
func (c *TTSClient) GenerateTTS(ctx context.Context, req TTSRequest) (*TTSResponse, error) {
	// Determine which service to use based on language
	if c.isIndicLanguage(req.Language) {
		return c.generateWithBulbul(ctx, req)
	}
	return c.generateWithOrpheus(ctx, req)
}

// generateWithBulbul generates TTS using the Bulbul adapter
func (c *TTSClient) generateWithBulbul(ctx context.Context, req TTSRequest) (*TTSResponse, error) {
	c.logger.Info("Generating TTS with Bulbul",
		zap.String("language", req.Language),
		zap.String("voice_id", req.VoiceID),
	)

	// Prepare request body
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.bulbulURL+"/api/v1/tts", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Bulbul: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var ttsResp TTSResponse
	if err := json.NewDecoder(resp.Body).Decode(&ttsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bulbul API error: %s", ttsResp.Error)
	}

	return &ttsResp, nil
}

// generateWithOrpheus generates TTS using the Orpheus inference service
func (c *TTSClient) generateWithOrpheus(ctx context.Context, req TTSRequest) (*TTSResponse, error) {
	c.logger.Info("Generating TTS with Orpheus",
		zap.String("language", req.Language),
		zap.String("voice_id", req.VoiceID),
	)

	// Prepare request body
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.orpheusURL+"/api/v1/tts", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Orpheus: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var ttsResp TTSResponse
	if err := json.NewDecoder(resp.Body).Decode(&ttsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Orpheus API error: %s", ttsResp.Error)
	}

	return &ttsResp, nil
}

// isIndicLanguage checks if the language is an Indic language
func (c *TTSClient) isIndicLanguage(lang string) bool {
	indicLanguages := map[string]bool{
		"hi": true, // Hindi
		"ta": true, // Tamil
		"te": true, // Telugu
		"ml": true, // Malayalam
		"kn": true, // Kannada
		"mr": true, // Marathi
		"gu": true, // Gujarati
		"bn": true, // Bengali
		"pa": true, // Punjabi
		"or": true, // Odia
	}
	return indicLanguages[lang]
}
