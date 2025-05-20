package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// PreviewMessage represents a message sent by the client
type PreviewMessage struct {
	Type    string  `json:"type"`
	Text    string  `json:"text,omitempty"`
	VoiceID string  `json:"voice_id,omitempty"`
	Speed   float64 `json:"speed,omitempty"`
	Pitch   float64 `json:"pitch,omitempty"`
	Emotion int     `json:"emotion,omitempty"`
}

// PreviewResponse represents a response sent to the client
type PreviewResponse struct {
	Type     string `json:"type"`
	Data     string `json:"data,omitempty"`
	Error    string `json:"error,omitempty"`
	Duration int    `json:"duration,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, this should be more restrictive
		return true
	},
}

// PreviewHandler handles WebSocket connections for the preview endpoint
func PreviewHandler(logger *zap.Logger) http.HandlerFunc {
	// This is a simple dummy audio chunk in base64
	// In a real implementation, this would be actual audio data
	dummyAudio := base64.StdEncoding.EncodeToString([]byte("DUMMY_AUDIO_DATA"))

	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("Failed to upgrade connection", zap.Error(err))
			return
		}
		defer conn.Close()

		logger.Info("Client connected to preview WebSocket")

		// Read messages from the client
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.Error("WebSocket connection closed unexpectedly", zap.Error(err))
				} else {
					logger.Info("WebSocket connection closed", zap.Error(err))
				}
				break
			}

			// Parse the message
			var previewMsg PreviewMessage
			if err := json.Unmarshal(message, &previewMsg); err != nil {
				logger.Error("Failed to parse message", zap.Error(err))
				
				// Send error back to client
				errorResp := PreviewResponse{
					Type:  "error",
					Error: "Invalid message format",
				}
				
				if err := conn.WriteJSON(errorResp); err != nil {
					logger.Error("Failed to send error message", zap.Error(err))
					return
				}
				
				continue
			}
			
			// Handle the message based on its type
			switch previewMsg.Type {
			case "text":
				// Log the received text
				logger.Info("Preview text received", 
					zap.String("text", previewMsg.Text),
					zap.String("voice_id", previewMsg.VoiceID))
				
				// Send back a simple acknowledgment
				ackResp := PreviewResponse{
					Type: "start",
				}
				
				if err := conn.WriteJSON(ackResp); err != nil {
					logger.Error("Failed to send ack message", zap.Error(err))
					return
				}
				
				// In a real implementation, this would stream audio chunks as they're generated
				// For now, we'll simulate streaming by sending a few chunks with delays
				
				// Calculate a fake "duration" based on text length (1ms per character)
				duration := len(previewMsg.Text)
				numChunks := 3 // Arbitrary number of chunks
				
				for i := 0; i < numChunks; i++ {
					// Wait a bit to simulate processing time
					time.Sleep(200 * time.Millisecond)
					
					// Send audio chunk
					audioResp := PreviewResponse{
						Type:     "audio",
						Data:     dummyAudio,
						Duration: duration / numChunks,
					}
					
					if err := conn.WriteJSON(audioResp); err != nil {
						logger.Error("Failed to send audio chunk", zap.Error(err))
						return
					}
				}
				
				// Send end marker
				endResp := PreviewResponse{
					Type:     "end",
					Duration: duration,
				}
				
				if err := conn.WriteJSON(endResp); err != nil {
					logger.Error("Failed to send end message", zap.Error(err))
					return
				}
				
			default:
				// Handle unknown message type
				logger.Warn("Unknown message type", zap.String("type", previewMsg.Type))
				
				errorResp := PreviewResponse{
					Type:  "error",
					Error: "Unknown message type",
				}
				
				if err := conn.WriteJSON(errorResp); err != nil {
					logger.Error("Failed to send error message", zap.Error(err))
					return
				}
			}
		}
	}
} 