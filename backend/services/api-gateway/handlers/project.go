package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/NavoDayAI/vaani/backend/services/api-gateway/database"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/middleware"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Project request/response types
type CreateProjectRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateProjectRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ProjectResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AudioURL    string `json:"audio_url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	IsCompleted bool   `json:"is_completed"`
}

type ProjectListResponse struct {
	Projects []ProjectResponse `json:"projects"`
	Total    int               `json:"total"`
}

// ProjectGenerateRequest for generating audio for a project
type ProjectGenerateRequest struct {
	Text       string  `json:"text"`
	VoiceID    string  `json:"voice_id"`
	Speed      float64 `json:"speed,omitempty"`
	Pitch      float64 `json:"pitch,omitempty"`
	Emotion    int     `json:"emotion,omitempty"`
	Language   string  `json:"language,omitempty"`
	OutputType string  `json:"output_type,omitempty"`
}

// ProjectExportRequest for exporting project audio
type ProjectExportRequest struct {
	Format  string                 `json:"format,omitempty"`
	Quality string                 `json:"quality,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// Helper function to convert Project model to response
func projectToResponse(project *models.Project) ProjectResponse {
	return ProjectResponse{
		ID:          project.ID,
		Title:       project.Name,
		Description: project.Description,
		AudioURL:    "", // TODO: Get latest audio URL from jobs
		CreatedAt:   project.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   project.UpdatedAt.Format(time.RFC3339),
		IsCompleted: project.JobCount > 0, // Simple heuristic for now
	}
}

// ListProjectsHandler handles GET /v1/projects
func ListProjectsHandler(logger *zap.Logger, projectRepo database.ProjectRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, err := middleware.ExtractUserID(r)
		if err != nil {
			logger.Error("Failed to extract user ID", zap.Error(err))
			sendJSONError(w, http.StatusUnauthorized, ValidationError{
				Error:   "unauthorized",
				Message: "User not authenticated",
			})
			return
		}

		// Get limit from query parameters (default 50)
		limit := 50
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
				limit = parsedLimit
			}
		}

		// Get projects from repository
		projects, err := projectRepo.ListUserProjects(userID, limit)
		if err != nil {
			logger.Error("Failed to list projects",
				zap.String("user_id", userID),
				zap.Error(err),
			)
			sendJSONError(w, http.StatusInternalServerError, ValidationError{
				Error:   "database_error",
				Message: "Failed to fetch projects",
			})
			return
		}

		// Convert to response format
		projectResponses := make([]ProjectResponse, len(projects))
		for i, project := range projects {
			projectResponses[i] = projectToResponse(project)
		}

		response := ProjectListResponse{
			Projects: projectResponses,
			Total:    len(projectResponses),
		}

		logger.Info("Projects listed successfully",
			zap.String("user_id", userID),
			zap.Int("count", len(projects)),
		)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// CreateProjectHandler handles POST /v1/projects
func CreateProjectHandler(logger *zap.Logger, projectRepo database.ProjectRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, err := middleware.ExtractUserID(r)
		if err != nil {
			logger.Error("Failed to extract user ID", zap.Error(err))
			sendJSONError(w, http.StatusUnauthorized, ValidationError{
				Error:   "unauthorized",
				Message: "User not authenticated",
			})
			return
		}

		// Parse request body
		var req CreateProjectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Failed to decode request body", zap.Error(err))
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "invalid_request",
				Message: "Invalid request body",
			})
			return
		}

		// Validate required fields
		if req.Title == "" {
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "validation_error",
				Field:   "title",
				Message: "Title is required",
			})
			return
		}

		// Create project model
		project := &models.Project{
			ID:           uuid.New().String(),
			UserID:       userID,
			Name:         req.Title,
			Description:  req.Description,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			JobCount:     0,
			TotalCredits: 0,
			UsedCredits:  0,
		}

		// Save to database
		if err := projectRepo.CreateProject(project); err != nil {
			logger.Error("Failed to create project",
				zap.String("user_id", userID),
				zap.String("title", req.Title),
				zap.Error(err),
			)
			sendJSONError(w, http.StatusInternalServerError, ValidationError{
				Error:   "database_error",
				Message: "Failed to create project",
			})
			return
		}

		// Return created project
		response := projectToResponse(project)

		logger.Info("Project created successfully",
			zap.String("user_id", userID),
			zap.String("project_id", project.ID),
			zap.String("title", req.Title),
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// GetProjectHandler handles GET /v1/projects/{id}
func GetProjectHandler(logger *zap.Logger, projectRepo database.ProjectRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, err := middleware.ExtractUserID(r)
		if err != nil {
			logger.Error("Failed to extract user ID", zap.Error(err))
			sendJSONError(w, http.StatusUnauthorized, ValidationError{
				Error:   "unauthorized",
				Message: "User not authenticated",
			})
			return
		}

		// Get project ID from URL
		projectID := chi.URLParam(r, "projectID")
		if projectID == "" {
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "invalid_parameter",
				Field:   "projectID",
				Message: "Project ID is required",
			})
			return
		}

		// Get project from database
		project, err := projectRepo.GetProject(userID, projectID)
		if err != nil {
			logger.Error("Failed to get project",
				zap.String("user_id", userID),
				zap.String("project_id", projectID),
				zap.Error(err),
			)
			sendJSONError(w, http.StatusNotFound, ValidationError{
				Error:   "not_found",
				Message: "Project not found",
			})
			return
		}

		// Return project
		response := projectToResponse(project)

		logger.Info("Project retrieved successfully",
			zap.String("user_id", userID),
			zap.String("project_id", projectID),
		)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// UpdateProjectHandler handles PUT /v1/projects/{id}
func UpdateProjectHandler(logger *zap.Logger, projectRepo database.ProjectRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, err := middleware.ExtractUserID(r)
		if err != nil {
			logger.Error("Failed to extract user ID", zap.Error(err))
			sendJSONError(w, http.StatusUnauthorized, ValidationError{
				Error:   "unauthorized",
				Message: "User not authenticated",
			})
			return
		}

		// Get project ID from URL
		projectID := chi.URLParam(r, "projectID")
		if projectID == "" {
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "invalid_parameter",
				Field:   "projectID",
				Message: "Project ID is required",
			})
			return
		}

		// Parse request body
		var req UpdateProjectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Failed to decode request body", zap.Error(err))
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "invalid_request",
				Message: "Invalid request body",
			})
			return
		}

		// Get existing project
		project, err := projectRepo.GetProject(userID, projectID)
		if err != nil {
			logger.Error("Failed to get project for update",
				zap.String("user_id", userID),
				zap.String("project_id", projectID),
				zap.Error(err),
			)
			sendJSONError(w, http.StatusNotFound, ValidationError{
				Error:   "not_found",
				Message: "Project not found",
			})
			return
		}

		// Update fields if provided
		if req.Title != "" {
			project.Name = req.Title
		}
		if req.Description != "" {
			project.Description = req.Description
		}

		// Save updated project
		if err := projectRepo.UpdateProject(project); err != nil {
			logger.Error("Failed to update project",
				zap.String("user_id", userID),
				zap.String("project_id", projectID),
				zap.Error(err),
			)
			sendJSONError(w, http.StatusInternalServerError, ValidationError{
				Error:   "database_error",
				Message: "Failed to update project",
			})
			return
		}

		// Return updated project
		response := projectToResponse(project)

		logger.Info("Project updated successfully",
			zap.String("user_id", userID),
			zap.String("project_id", projectID),
		)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// DeleteProjectHandler handles DELETE /v1/projects/{id}
func DeleteProjectHandler(logger *zap.Logger, projectRepo database.ProjectRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, err := middleware.ExtractUserID(r)
		if err != nil {
			logger.Error("Failed to extract user ID", zap.Error(err))
			sendJSONError(w, http.StatusUnauthorized, ValidationError{
				Error:   "unauthorized",
				Message: "User not authenticated",
			})
			return
		}

		// Get project ID from URL
		projectID := chi.URLParam(r, "projectID")
		if projectID == "" {
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "invalid_parameter",
				Field:   "projectID",
				Message: "Project ID is required",
			})
			return
		}

		// Delete project
		if err := projectRepo.DeleteProject(userID, projectID); err != nil {
			logger.Error("Failed to delete project",
				zap.String("user_id", userID),
				zap.String("project_id", projectID),
				zap.Error(err),
			)
			sendJSONError(w, http.StatusNotFound, ValidationError{
				Error:   "not_found",
				Message: "Project not found",
			})
			return
		}

		logger.Info("Project deleted successfully",
			zap.String("user_id", userID),
			zap.String("project_id", projectID),
		)

		w.WriteHeader(http.StatusNoContent)
	}
}

// ProjectGenerateHandler handles POST /v1/projects/{id}/generate
func ProjectGenerateHandler(logger *zap.Logger, projectRepo database.ProjectRepository, jobRepo models.JobRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, err := middleware.ExtractUserID(r)
		if err != nil {
			logger.Error("Failed to extract user ID", zap.Error(err))
			sendJSONError(w, http.StatusUnauthorized, ValidationError{
				Error:   "unauthorized",
				Message: "User not authenticated",
			})
			return
		}

		// Get project ID from URL
		projectID := chi.URLParam(r, "projectID")
		if projectID == "" {
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "invalid_parameter",
				Field:   "projectID",
				Message: "Project ID is required",
			})
			return
		}

		// Verify project exists and belongs to user
		project, err := projectRepo.GetProject(userID, projectID)
		if err != nil {
			logger.Error("Failed to get project for generation",
				zap.String("user_id", userID),
				zap.String("project_id", projectID),
				zap.Error(err),
			)
			sendJSONError(w, http.StatusNotFound, ValidationError{
				Error:   "not_found",
				Message: "Project not found",
			})
			return
		}

		// Parse request body
		var req ProjectGenerateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Failed to decode request body", zap.Error(err))
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "invalid_request",
				Message: "Invalid request body",
			})
			return
		}

		// Validate required fields
		if req.Text == "" {
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "validation_error",
				Field:   "text",
				Message: "Text is required",
			})
			return
		}

		if req.VoiceID == "" {
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "validation_error",
				Field:   "voice_id",
				Message: "Voice ID is required",
			})
			return
		}

		// Create job for the project
		jobID := uuid.New().String()
		job := &models.Job{
			ID:         jobID,
			ProjectID:  projectID,
			UserID:     userID,
			Text:       req.Text,
			VoiceID:    req.VoiceID,
			Speed:      req.Speed,
			Pitch:      req.Pitch,
			Emotion:    req.Emotion,
			Language:   req.Language,
			OutputType: req.OutputType,
			TextLength: len(req.Text),
			Status:     models.StatusQueued,
			CreatedAt:  time.Now(),
		}

		// Set defaults
		if job.Speed == 0 {
			job.Speed = 1.0
		}
		if job.OutputType == "" {
			job.OutputType = "mp3"
		}

		// Save job
		if err := jobRepo.CreateJob(job); err != nil {
			logger.Error("Failed to create job for project",
				zap.String("user_id", userID),
				zap.String("project_id", projectID),
				zap.Error(err),
			)
			sendJSONError(w, http.StatusInternalServerError, ValidationError{
				Error:   "database_error",
				Message: "Failed to create generation job",
			})
			return
		}

		// Update project job count
		project.JobCount++
		if err := projectRepo.UpdateProject(project); err != nil {
			logger.Warn("Failed to update project job count",
				zap.String("project_id", projectID),
				zap.Error(err),
			)
		}

		// Return job information
		response := GenerateResponse{
			JobID:     jobID,
			Status:    string(job.Status),
			QueueTime: job.CreatedAt,
			Message:   "Job has been queued for processing",
		}

		logger.Info("Generation job created for project",
			zap.String("user_id", userID),
			zap.String("project_id", projectID),
			zap.String("job_id", jobID),
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(response)
	}
}

// ProjectExportHandler handles POST /v1/projects/{id}/export
func ProjectExportHandler(logger *zap.Logger, projectRepo database.ProjectRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, err := middleware.ExtractUserID(r)
		if err != nil {
			logger.Error("Failed to extract user ID", zap.Error(err))
			sendJSONError(w, http.StatusUnauthorized, ValidationError{
				Error:   "unauthorized",
				Message: "User not authenticated",
			})
			return
		}

		// Get project ID from URL
		projectID := chi.URLParam(r, "projectID")
		if projectID == "" {
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "invalid_parameter",
				Field:   "projectID",
				Message: "Project ID is required",
			})
			return
		}

		// Verify project exists and belongs to user
		_, err = projectRepo.GetProject(userID, projectID)
		if err != nil {
			logger.Error("Failed to get project for export",
				zap.String("user_id", userID),
				zap.String("project_id", projectID),
				zap.Error(err),
			)
			sendJSONError(w, http.StatusNotFound, ValidationError{
				Error:   "not_found",
				Message: "Project not found",
			})
			return
		}

		// Parse request body
		var req ProjectExportRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Failed to decode request body", zap.Error(err))
			sendJSONError(w, http.StatusBadRequest, ValidationError{
				Error:   "invalid_request",
				Message: "Invalid request body",
			})
			return
		}

		// For now, return a placeholder response
		// TODO: Implement actual export functionality
		response := map[string]interface{}{
			"url":     "https://placeholder-export-url.com/audio.mp3",
			"format":  req.Format,
			"expires": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		}

		logger.Info("Project export requested",
			zap.String("user_id", userID),
			zap.String("project_id", projectID),
			zap.String("format", req.Format),
		)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
