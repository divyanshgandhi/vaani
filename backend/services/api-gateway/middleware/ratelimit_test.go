package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
)

// nopCloser is an io.ReadCloser with a no-op Close method
type nopCloser struct {
	io.Reader
}

// Close implements the io.Closer interface
func (nopCloser) Close() error { return nil }

func TestRateLimitMiddleware(t *testing.T) {
	// Create a test logger
	logger, _ := zap.NewDevelopment()

	// Create a test handler that always returns 200 OK
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply rate limit middleware - 3 requests per second, reset after 5 seconds
	rateLimitedHandler := RateLimitMiddleware(3, 5*time.Second, logger)(okHandler)

	// Create a function to make test requests
	makeRequest := func(t *testing.T) (*http.Response, error) {
		req, err := http.NewRequest("GET", "http://example.com/test", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		// Set a fixed IP address
		req.Header.Set("X-Forwarded-For", "192.168.1.1")

		// Use httptest.ResponseRecorder to capture the response
		rr := httptest.NewRecorder()
		rateLimitedHandler.ServeHTTP(rr, req)

		// Convert recorder to response
		return &http.Response{
			StatusCode: rr.Code,
			Body:       nopCloser{strings.NewReader(rr.Body.String())},
			Header:     rr.Header(),
		}, nil
	}

	// First 3 requests should succeed
	for i := 0; i < 3; i++ {
		resp, err := makeRequest(t)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	}

	// 4th request should be rate limited
	resp, err := makeRequest(t)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", resp.StatusCode)
	}

	// Check response headers
	retryAfter := resp.Header.Get("Retry-After")
	if retryAfter != "1" {
		t.Errorf("Expected Retry-After header to be 1, got %s", retryAfter)
	}
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		remoteAddr string
		expected   string
	}{
		{
			name: "X-Forwarded-For",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.195",
			},
			remoteAddr: "192.168.1.1:1234",
			expected:   "203.0.113.195",
		},
		{
			name: "X-Real-IP",
			headers: map[string]string{
				"X-Real-IP": "203.0.113.195",
			},
			remoteAddr: "192.168.1.1:1234",
			expected:   "203.0.113.195",
		},
		{
			name:       "RemoteAddr only",
			headers:    map[string]string{},
			remoteAddr: "192.168.1.1:1234",
			expected:   "192.168.1.1:1234",
		},
		{
			name: "Both headers, X-Forwarded-For takes precedence",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.195",
				"X-Real-IP":       "198.51.100.42",
			},
			remoteAddr: "192.168.1.1:1234",
			expected:   "203.0.113.195",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			req.RemoteAddr = tt.remoteAddr

			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			ip := getClientIP(req)
			if ip != tt.expected {
				t.Errorf("getClientIP() = %v, want %v", ip, tt.expected)
			}
		})
	}
}

func TestRateLimitExceededResponse(t *testing.T) {
	// Create a test logger
	logger, _ := zap.NewDevelopment()

	// Create a test handler
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Apply rate limit middleware with very low limits to trigger 429
	handler := RateLimitMiddleware(1, 5*time.Second, logger)(okHandler)

	// Create a request
	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1")

	// First request should succeed
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req)
	if rr1.Code != http.StatusOK {
		t.Errorf("First request should succeed, got status %d", rr1.Code)
	}

	// Second request should be rate limited
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req)
	if rr2.Code != http.StatusTooManyRequests {
		t.Errorf("Second request should be rate limited, got status %d", rr2.Code)
	}

	// Parse the JSON response
	var response RateLimitExceededResponse
	if err := json.NewDecoder(rr2.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	// Verify the response fields
	if response.Error != "Rate limit exceeded" {
		t.Errorf("Expected error message 'Rate limit exceeded', got '%s'", response.Error)
	}
	if response.RequestsPerSec != 1 {
		t.Errorf("Expected RequestsPerSec to be 1, got %d", response.RequestsPerSec)
	}
	if response.ResetAfter != 5 {
		t.Errorf("Expected ResetAfter to be 5, got %d", response.ResetAfter)
	}
	if response.RetryAfter != 1 {
		t.Errorf("Expected RetryAfter to be 1, got %d", response.RetryAfter)
	}
}
