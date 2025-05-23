package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NavoDayAI/vaani/backend/services/api-gateway/middleware"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/models"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// JobStatusResponse represents the response for job status queries
type JobStatusResponse struct {
	JobID        string `json:"job_id"`
	Status       string `json:"status"`
	Progress     int    `json:"progress,omitempty"`
	OutputURL    string `json:"output_url,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	CreatedAt    string `json:"created_at"`
	CompletedAt  string `json:"completed_at,omitempty"`
}

// GetJobHandler handles GET requests to /v1/jobs/{jobID}
func GetJobHandler(logger *zap.Logger, jobRepo models.JobRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get job ID from URL parameter
		jobID := chi.URLParam(r, "jobID")
		if jobID == "" {
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "invalid_parameter",
				Field:   "jobID",
				Message: "Job ID is required",
			})
			return
		}

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

		// For now, use userID as projectID (this should be improved later)
		projectID := userID

		// Get the job from the repository
		job, err := jobRepo.GetJob(projectID, jobID)
		if err != nil {
			logger.Error("Failed to get job",
				zap.String("job_id", jobID),
				zap.String("project_id", projectID),
				zap.Error(err),
			)
			sendJSONError(w, http.StatusNotFound, ValidationError{
				Error:   "job_not_found",
				Message: "Job not found",
			})
			return
		}

		// Create response
		response := JobStatusResponse{
			JobID:        job.ID,
			Status:       string(job.Status),
			OutputURL:    job.OutputURL,
			ErrorMessage: job.ErrorMessage,
			CreatedAt:    job.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}

		if !job.CompletedAt.IsZero() {
			completedAt := job.CompletedAt.Format("2006-01-02T15:04:05Z")
			response.CompletedAt = completedAt
		}

		// Calculate progress based on status
		switch job.Status {
		case models.StatusQueued:
			response.Progress = 0
		case models.StatusProcessing:
			response.Progress = 50
		case models.StatusComplete:
			response.Progress = 100
		case models.StatusFailed:
			response.Progress = 0
		}

		logger.Info("Job status retrieved",
			zap.String("job_id", jobID),
			zap.String("status", string(job.Status)),
		)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// DownloadHandler handles GET requests to /v1/download/{jobID}
func DownloadHandler(logger *zap.Logger, jobRepo models.JobRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get job ID from URL parameter
		jobID := chi.URLParam(r, "jobID")
		if jobID == "" {
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "invalid_parameter",
				Field:   "jobID",
				Message: "Job ID is required",
			})
			return
		}

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

		// For now, use userID as projectID
		projectID := userID

		// Get the job from the repository
		job, err := jobRepo.GetJob(projectID, jobID)
		if err != nil {
			logger.Error("Failed to get job for download",
				zap.String("job_id", jobID),
				zap.String("project_id", projectID),
				zap.Error(err),
			)
			sendJSONError(w, http.StatusNotFound, ValidationError{
				Error:   "job_not_found",
				Message: "Job not found",
			})
			return
		}

		// Check if job is complete and has output URL
		if job.Status != models.StatusComplete || job.OutputURL == "" {
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "job_not_ready",
				Message: "Job is not complete or has no output",
			})
			return
		}

		logger.Info("Redirecting to download URL",
			zap.String("job_id", jobID),
			zap.String("output_url", job.OutputURL),
		)

		// Redirect to the pre-signed URL
		http.Redirect(w, r, job.OutputURL, http.StatusTemporaryRedirect)
	}
}
