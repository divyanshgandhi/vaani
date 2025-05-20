package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NavoDayAI/vaani/backend/services/api-gateway/middleware"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/models"
	"go.uber.org/zap"
)

// MockJobRepository is a mock implementation of models.JobRepository for testing
type MockJobRepository struct {
	createJobFunc func(job *models.Job) error
}

func (m *MockJobRepository) CreateJob(job *models.Job) error {
	if m.createJobFunc != nil {
		return m.createJobFunc(job)
	}
	return nil
}

func (m *MockJobRepository) GetJob(projectID, jobID string) (*models.Job, error) {
	return nil, nil
}

func (m *MockJobRepository) UpdateJobStatus(projectID, jobID string, status models.JobStatus) error {
	return nil
}

func (m *MockJobRepository) UpdateJobOutput(projectID, jobID, outputURL string) error {
	return nil
}

func (m *MockJobRepository) ListJobs(projectID string, limit int) ([]*models.Job, error) {
	return nil, nil
}

// Create a test context with a user ID for middleware
func createTestContextWithUser(r *http.Request, userID string) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.UserContextKey, userID)
	return r.WithContext(ctx)
}

func TestGenerateHandler(t *testing.T) {
	// Create a test logger
	logger, _ := zap.NewDevelopment()

	// Create a mock repository
	mockRepo := &MockJobRepository{}

	// Define test cases
	tests := []struct {
		name           string
		requestBody    string
		userID         string
		mockCreateJob  func(job *models.Job) error
		expectedStatus int
		validateResp   func(*testing.T, []byte)
	}{
		{
			name:           "Valid request",
			requestBody:    `{"text": "Hello world", "voice_id": "test-voice"}`,
			userID:         "test-user-123",
			mockCreateJob:  nil, // Default nil implementation returns no error
			expectedStatus: http.StatusAccepted,
			validateResp: func(t *testing.T, body []byte) {
				var resp GenerateResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				
				if resp.JobID == "" {
					t.Error("Expected job_id to be present")
				}
				
				if resp.Status != "queued" {
					t.Errorf("Expected status to be 'queued', got '%s'", resp.Status)
				}
			},
		},
		{
			name:           "Database error",
			requestBody:    `{"text": "Hello world", "voice_id": "test-voice"}`,
			userID:         "test-user-123",
			mockCreateJob:  func(job *models.Job) error { return errors.New("database error") },
			expectedStatus: http.StatusInternalServerError,
			validateResp: func(t *testing.T, body []byte) {
				var resp ValidationError
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				
				if resp.Error != "database_error" {
					t.Errorf("Expected error to be 'database_error', got '%s'", resp.Error)
				}
			},
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"text": "Hello world", "voice_id": }`,
			userID:         "test-user-123",
			mockCreateJob:  nil,
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, body []byte) {
				var resp ValidationError
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				
				if resp.Error != "invalid_request" {
					t.Errorf("Expected error to be 'invalid_request', got '%s'", resp.Error)
				}
			},
		},
		{
			name:           "Empty text",
			requestBody:    `{"text": "", "voice_id": "test-voice"}`,
			userID:         "test-user-123",
			mockCreateJob:  nil,
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, body []byte) {
				var resp ValidationError
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				
				if resp.Error != "invalid_parameter" {
					t.Errorf("Expected error to be 'invalid_parameter', got '%s'", resp.Error)
				}
				
				if resp.Field != "text" {
					t.Errorf("Expected field to be 'text', got '%s'", resp.Field)
				}
			},
		},
		{
			name:           "Text too long",
			requestBody:    `{"text": "` + strings.Repeat("a", MaxTextLength+1) + `", "voice_id": "test-voice"}`,
			userID:         "test-user-123",
			mockCreateJob:  nil,
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, body []byte) {
				var resp ValidationError
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				
				if resp.Error != "invalid_parameter" {
					t.Errorf("Expected error to be 'invalid_parameter', got '%s'", resp.Error)
				}
				
				if resp.Field != "text" {
					t.Errorf("Expected field to be 'text', got '%s'", resp.Field)
				}
				
				if !strings.Contains(resp.Message, "exceeds maximum length") {
					t.Errorf("Expected message to mention maximum length, got '%s'", resp.Message)
				}
			},
		},
		{
			name:           "Missing voice_id",
			requestBody:    `{"text": "Hello world"}`,
			userID:         "test-user-123",
			mockCreateJob:  nil,
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, body []byte) {
				var resp ValidationError
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				
				if resp.Error != "invalid_parameter" {
					t.Errorf("Expected error to be 'invalid_parameter', got '%s'", resp.Error)
				}
				
				if resp.Field != "voice_id" {
					t.Errorf("Expected field to be 'voice_id', got '%s'", resp.Field)
				}
			},
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the mock
			if tc.mockCreateJob != nil {
				mockRepo.createJobFunc = tc.mockCreateJob
			} else {
				mockRepo.createJobFunc = nil
			}

			// Create a request
			req, err := http.NewRequest("POST", "/v1/generate", bytes.NewBufferString(tc.requestBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			
			// Add user to context
			if tc.userID != "" {
				req = createTestContextWithUser(req, tc.userID)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Create and call the handler
			handler := GenerateHandler(logger, mockRepo)
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			// Validate response
			tc.validateResp(t, rr.Body.Bytes())
		})
	}
} 