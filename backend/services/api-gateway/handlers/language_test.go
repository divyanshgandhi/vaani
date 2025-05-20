package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NavoDayAI/vaani/backend/services/api-gateway/utils/langdetect"
	"go.uber.org/zap"
)

func TestDetectLanguageHandler(t *testing.T) {
	// Create a test logger
	logger, _ := zap.NewDevelopment()

	// Define test cases
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		validateResponse func(*testing.T, *LanguageDetectionResponse)
	}{
		{
			name: "Valid English Text",
			requestBody: map[string]interface{}{
				"text": "This is a sample text in English to test language detection.",
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, resp *LanguageDetectionResponse) {
				if resp.Language != langdetect.English {
					t.Errorf("Expected language %s, got %s", langdetect.English, resp.Language)
				}
				if resp.IsIndic {
					t.Errorf("Expected IsIndic to be false for English")
				}
			},
		},
		{
			name: "Valid Hindi Text",
			requestBody: map[string]interface{}{
				"text": "यह हिंदी में एक नमूना पाठ है जिसका उपयोग भाषा का पता लगाने के लिए किया जाता है।",
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, resp *LanguageDetectionResponse) {
				if resp.Language != langdetect.Hindi {
					t.Errorf("Expected language %s, got %s", langdetect.Hindi, resp.Language)
				}
				if !resp.IsIndic {
					t.Errorf("Expected IsIndic to be true for Hindi")
				}
			},
		},
		{
			name: "Text Too Short",
			requestBody: map[string]interface{}{
				"text": "Hi",
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, resp *LanguageDetectionResponse) {
				if resp.Language != langdetect.Unknown {
					t.Errorf("Expected language %s for short text, got %s", langdetect.Unknown, resp.Language)
				}
			},
		},
		{
			name: "Empty Text",
			requestBody: map[string]interface{}{
				"text": "",
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: nil,
		},
		{
			name: "Invalid JSON",
			requestBody: nil, // Will cause JSON marshaling to produce invalid JSON
			expectedStatus: http.StatusBadRequest,
			validateResponse: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			var reqBody []byte
			var err error
			if tt.requestBody != nil {
				reqBody, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			} else {
				reqBody = []byte("{invalid json")
			}

			req, err := http.NewRequest("POST", "/detect-language", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Create handler
			handler := DetectLanguageHandler(logger)

			// Handle request
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
				return
			}

			// For successful responses, validate the response body
			if tt.expectedStatus == http.StatusOK && tt.validateResponse != nil {
				var response LanguageDetectionResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				// Check the expected response
				tt.validateResponse(t, &response)
				
				// Check that supported languages are included
				if len(response.SupportedLangs) == 0 {
					t.Errorf("Expected supported languages to be included in response")
				}
			}
		})
	}
} 