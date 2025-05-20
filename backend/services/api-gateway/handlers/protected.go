package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NavoDayAI/vaani/backend/services/api-gateway/middleware"
)

// UserResponse is the structure for user profile data
type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// MeHandler returns the current authenticated user's profile
func MeHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user information from the context
	userID, err := middleware.ExtractUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: Unable to identify user", http.StatusUnauthorized)
		return
	}

	// Optionally get email and phone
	email, _ := middleware.ExtractUserEmail(r)
	phone, _ := middleware.ExtractUserPhone(r)

	// Create response
	response := UserResponse{
		ID:    userID,
		Email: email,
		Phone: phone,
	}

	// Convert to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
} 