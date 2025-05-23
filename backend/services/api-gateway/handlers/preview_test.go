package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// TestPreviewHandlerWebSocket tests the WebSocket functionality of the preview handler
func TestPreviewHandlerWebSocket(t *testing.T) {
	// Create a test logger
	logger, _ := zap.NewDevelopment()

	// Create a test server
	server := httptest.NewServer(PreviewHandler(logger))
	defer server.Close()

	// Convert http:// URL to ws:// URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Create a test message
	testMsg := PreviewMessage{
		Type:    "text",
		Text:    "Hello world",
		VoiceID: "test-voice",
	}

	// Send the test message
	if err := ws.WriteJSON(testMsg); err != nil {
		t.Fatalf("Failed to write message: %v", err)
	}

	// Read the "start" response
	var startResp PreviewResponse
	if err := ws.ReadJSON(&startResp); err != nil {
		t.Fatalf("Failed to read start response: %v", err)
	}

	if startResp.Type != "start" {
		t.Errorf("Expected start response, got %s", startResp.Type)
	}

	// Read the audio chunks (we're sending 3 in the implementation)
	for i := 0; i < 3; i++ {
		var audioResp PreviewResponse
		if err := ws.ReadJSON(&audioResp); err != nil {
			t.Fatalf("Failed to read audio response %d: %v", i, err)
		}

		if audioResp.Type != "audio" {
			t.Errorf("Expected audio response, got %s", audioResp.Type)
		}

		if audioResp.Data == "" {
			t.Error("Expected non-empty audio data")
		}
	}

	// Read the "end" response
	var endResp PreviewResponse
	if err := ws.ReadJSON(&endResp); err != nil {
		t.Fatalf("Failed to read end response: %v", err)
	}

	if endResp.Type != "end" {
		t.Errorf("Expected end response, got %s", endResp.Type)
	}

	if endResp.Duration != len(testMsg.Text) {
		t.Errorf("Expected duration %d, got %d", len(testMsg.Text), endResp.Duration)
	}
}

// TestPreviewHandlerInvalidMessage tests that the handler properly responds to invalid messages
func TestPreviewHandlerInvalidMessage(t *testing.T) {
	// Create a test logger
	logger, _ := zap.NewDevelopment()

	// Create a test server
	server := httptest.NewServer(PreviewHandler(logger))
	defer server.Close()

	// Convert http:// URL to ws:// URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Send an invalid JSON message
	if err := ws.WriteMessage(websocket.TextMessage, []byte("invalid json")); err != nil {
		t.Fatalf("Failed to write message: %v", err)
	}

	// Read the error response
	var errorResp PreviewResponse
	if err := ws.ReadJSON(&errorResp); err != nil {
		t.Fatalf("Failed to read error response: %v", err)
	}

	if errorResp.Type != "error" {
		t.Errorf("Expected error response, got %s", errorResp.Type)
	}

	if errorResp.Error != "Invalid message format" {
		t.Errorf("Expected 'Invalid message format' error, got '%s'", errorResp.Error)
	}
}

// TestPreviewHandlerUnknownType tests that the handler properly responds to unknown message types
func TestPreviewHandlerUnknownType(t *testing.T) {
	// Create a test logger
	logger, _ := zap.NewDevelopment()

	// Create a test server
	server := httptest.NewServer(PreviewHandler(logger))
	defer server.Close()

	// Convert http:// URL to ws:// URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Create a message with an unknown type
	unknownMsg := PreviewMessage{
		Type: "unknown",
		Text: "Hello world",
	}

	// Send the unknown type message
	if err := ws.WriteJSON(unknownMsg); err != nil {
		t.Fatalf("Failed to write message: %v", err)
	}

	// Read the error response
	var errorResp PreviewResponse
	if err := ws.ReadJSON(&errorResp); err != nil {
		t.Fatalf("Failed to read error response: %v", err)
	}

	if errorResp.Type != "error" {
		t.Errorf("Expected error response, got %s", errorResp.Type)
	}

	if errorResp.Error != "Unknown message type" {
		t.Errorf("Expected 'Unknown message type' error, got '%s'", errorResp.Error)
	}
}
