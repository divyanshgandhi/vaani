package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/middleware"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/models"
	"go.uber.org/zap"
)

const (
	// MaxTextLength is the maximum allowed length for text to be processed
	MaxTextLength = 5000
)

// GenerateRequest represents the request body for the generate endpoint
type GenerateRequest struct {
	Text       string  `json:"text"`
	VoiceID    string  `json:"voice_id"`
	Speed      float64 `json:"speed,omitempty"`
	Pitch      float64 `json:"pitch,omitempty"`
	Emotion    int     `json:"emotion,omitempty"`
	Language   string  `json:"language,omitempty"`
	OutputType string  `json:"output_type,omitempty"`
}

// GenerateResponse represents the response from the generate endpoint
type GenerateResponse struct {
	JobID     string    `json:"job_id"`
	Status    string    `json:"status"`
	QueueTime time.Time `json:"queue_time"`
	Message   string    `json:"message,omitempty"`
}

// ValidationError represents an error response for validation issues
type ValidationError struct {
	Error   string `json:"error"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// GenerateHandler handles POST requests to /v1/generate
func GenerateHandler(logger *zap.Logger, jobRepo models.JobRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var req GenerateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Error("Failed to decode request body", zap.Error(err))
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "invalid_request",
				Message: "Could not parse request body",
			})
			return
		}

		// Validate text length
		if len(strings.TrimSpace(req.Text)) == 0 {
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "invalid_parameter",
				Field:   "text",
				Message: "Text cannot be empty",
			})
			return
		}
		
		if len(req.Text) > MaxTextLength {
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "invalid_parameter",
				Field:   "text",
				Message: "Text exceeds maximum length of 5000 characters",
			})
			return
		}

		// Validate voice_id is present
		if req.VoiceID == "" {
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "invalid_parameter",
				Field:   "voice_id",
				Message: "Voice ID is required",
			})
			return
		}

		// Generate a job ID
		jobID := uuid.New().String()
		
		// Get user ID from context
		userID, err := middleware.ExtractUserID(r)
		if err != nil {
			logger.Error("Failed to extract user ID", zap.Error(err))
			sendJSONError(w, http.StatusInternalServerError, ValidationError{
				Error:   "internal_error",
				Message: "Failed to process request",
			})
			return
		}
		
		// Default project ID to user ID for now
		// In a real implementation, this would come from the request or user settings
		projectID := userID
		
		// Create the job model
		job := &models.Job{
			ID:          jobID,
			ProjectID:   projectID,
			UserID:      userID,
			Text:        req.Text,
			VoiceID:     req.VoiceID,
			Speed:       req.Speed,
			Pitch:       req.Pitch,
			Emotion:     req.Emotion,
			Language:    req.Language,
			OutputType:  req.OutputType,
			TextLength:  len(req.Text),
			Status:      models.StatusQueued,
			CreatedAt:   time.Now(),
		}
		
		// Store the job in Firestore
		err = jobRepo.CreateJob(job)
		if err != nil {
			logger.Error("Failed to create job in database", 
				zap.String("job_id", jobID),
				zap.Error(err),
			)
			sendJSONError(w, http.StatusInternalServerError, ValidationError{
				Error:   "database_error",
				Message: "Failed to store job",
			})
			return
		}
		
		// Create the response
		response := GenerateResponse{
			JobID:     jobID,
			Status:    string(job.Status),
			QueueTime: job.CreatedAt,
			Message:   "Job has been queued for processing",
		}

		logger.Info("Generate job created",
			zap.String("job_id", jobID),
			zap.String("project_id", projectID),
			zap.String("user_id", userID),
			zap.String("voice_id", req.VoiceID),
			zap.Int("text_length", len(req.Text)),
		)

		// Return 202 Accepted with job ID
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(response)
	}
}

// sendJSONError sends a JSON error response
func sendJSONError(w http.ResponseWriter, statusCode int, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err)
} 