package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NavoDayAI/vaani/backend/services/api-gateway/config"
	"go.uber.org/zap"
)

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		expectedVal string
	}{
		{
			name:        "Valid Bearer Token",
			authHeader:  "Bearer validTokenHere",
			expectedVal: "validTokenHere",
		},
		{
			name:        "No Bearer Prefix",
			authHeader:  "validTokenHere",
			expectedVal: "",
		},
		{
			name:        "Empty Header",
			authHeader:  "",
			expectedVal: "",
		},
		{
			name:        "Bearer With No Token",
			authHeader:  "Bearer ",
			expectedVal: "",
		},
		{
			name:        "Extra Parts",
			authHeader:  "Bearer token extra",
			expectedVal: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractToken(tt.authHeader)
			if result != tt.expectedVal {
				t.Errorf("extractToken(%q) = %q, want %q", tt.authHeader, result, tt.expectedVal)
			}
		})
	}
}

func TestFirebaseAuthMiddlewareUnauthorized(t *testing.T) {
	// Create a simple HTTP handler that we'll wrap with our middleware
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a test logger
	logger, _ := zap.NewDevelopment()

	// Create a mock config
	cfg := &config.AppConfig{}
	cfg.Firebase.ProjectID = "test-project"

	// Use the test middleware that doesn't try to initialize Firebase
	middleware := TestFirebaseAuthMiddleware(cfg, logger)

	// Create a request with no Authorization header
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// This should run the middleware and reject the request
	handler := middleware(nextHandler)
	handler.ServeHTTP(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestFirebaseAuthMiddlewareAuthorized(t *testing.T) {
	// Create a simple HTTP handler that we'll wrap with our middleware
	var capturedUserID string
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := ExtractUserID(r)
		if err == nil {
			capturedUserID = userID
		}
		w.WriteHeader(http.StatusOK)
	})

	// Create a test logger
	logger, _ := zap.NewDevelopment()

	// Create a mock config
	cfg := &config.AppConfig{}
	cfg.Firebase.ProjectID = "test-project"

	// Use the test middleware that doesn't try to initialize Firebase
	middleware := TestFirebaseAuthMiddleware(cfg, logger)

	// Create a request with a valid test token
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer test-valid-token")
	rr := httptest.NewRecorder()

	// This should run the middleware and accept the request
	handler := middleware(nextHandler)
	handler.ServeHTTP(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check that the user ID was extracted
	if capturedUserID != "test-user-id" {
		t.Errorf("Expected user ID 'test-user-id', got %q", capturedUserID)
	}
}
