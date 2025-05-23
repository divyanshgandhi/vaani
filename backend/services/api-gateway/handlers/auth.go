package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/NavoDayAI/vaani/backend/services/api-gateway/services"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

// OTP storage - in production, use Redis or database
var (
	otpStore = make(map[string]OTPData)
	otpMutex = sync.RWMutex{}
)

type OTPData struct {
	OTP       string
	ExpiresAt time.Time
	Attempts  int
}

type SendOTPRequest struct {
	PhoneNumber string `json:"phoneNumber"`
}

type SendOTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type VerifyOTPRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	OTP         string `json:"otp"`
}

type VerifyOTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

// JWT secret - in production, use environment variable or secure key management
var jwtSecret = []byte("your-secret-key-change-this-in-production")

// SendOTPHandler handles sending OTP to phone number
func SendOTPHandler(logger *zap.Logger) http.HandlerFunc {
	// Initialize SMS service
	smsService := services.NewSMSService(logger)

	return func(w http.ResponseWriter, r *http.Request) {
		var req SendOTPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Failed to decode request", zap.Error(err))
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.PhoneNumber == "" {
			http.Error(w, "Phone number is required", http.StatusBadRequest)
			return
		}

		// Generate 6-digit OTP
		otp := generateOTP()

		// Store OTP with 5-minute expiration
		otpMutex.Lock()
		otpStore[req.PhoneNumber] = OTPData{
			OTP:       otp,
			ExpiresAt: time.Now().Add(5 * time.Minute),
			Attempts:  0,
		}
		otpMutex.Unlock()

		// Send OTP via SMS
		if err := smsService.SendOTP(req.PhoneNumber, otp); err != nil {
			logger.Error("Failed to send OTP SMS",
				zap.String("phone", req.PhoneNumber),
				zap.Error(err))

			// Clean up stored OTP on send failure
			otpMutex.Lock()
			delete(otpStore, req.PhoneNumber)
			otpMutex.Unlock()

			http.Error(w, "Failed to send OTP. Please try again.", http.StatusInternalServerError)
			return
		}

		logger.Info("OTP sent successfully",
			zap.String("phone", req.PhoneNumber))

		response := SendOTPResponse{
			Success: true,
			Message: "OTP sent successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error("Failed to encode response", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

// VerifyOTPHandler handles OTP verification and returns JWT token
func VerifyOTPHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req VerifyOTPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Failed to decode request", zap.Error(err))
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.PhoneNumber == "" || req.OTP == "" {
			http.Error(w, "Phone number and OTP are required", http.StatusBadRequest)
			return
		}

		otpMutex.Lock()
		otpData, exists := otpStore[req.PhoneNumber]
		if exists {
			otpData.Attempts++
			otpStore[req.PhoneNumber] = otpData
		}
		otpMutex.Unlock()

		if !exists {
			response := VerifyOTPResponse{
				Success: false,
				Message: "No OTP found for this phone number",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Check if OTP has expired
		if time.Now().After(otpData.ExpiresAt) {
			// Clean up expired OTP
			otpMutex.Lock()
			delete(otpStore, req.PhoneNumber)
			otpMutex.Unlock()

			response := VerifyOTPResponse{
				Success: false,
				Message: "OTP has expired",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Check max attempts
		if otpData.Attempts > 3 {
			// Clean up after too many attempts
			otpMutex.Lock()
			delete(otpStore, req.PhoneNumber)
			otpMutex.Unlock()

			response := VerifyOTPResponse{
				Success: false,
				Message: "Too many attempts. Please request a new OTP",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Verify OTP
		if req.OTP != otpData.OTP {
			response := VerifyOTPResponse{
				Success: false,
				Message: "Invalid OTP",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			return
		}

		// OTP is valid, generate JWT token
		token, err := generateJWT(req.PhoneNumber)
		if err != nil {
			logger.Error("Failed to generate JWT", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Clean up used OTP
		otpMutex.Lock()
		delete(otpStore, req.PhoneNumber)
		otpMutex.Unlock()

		response := VerifyOTPResponse{
			Success: true,
			Message: "OTP verified successfully",
			Token:   token,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error("Failed to encode response", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

// generateOTP creates a 6-digit OTP
func generateOTP() string {
	otp := rand.Intn(900000) + 100000 // Generate number between 100000 and 999999
	return strconv.Itoa(otp)
}

// generateJWT creates a JWT token for the authenticated user
func generateJWT(phoneNumber string) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"phone_number": phoneNumber,
		"exp":          time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
		"iat":          time.Now().Unix(),
		"iss":          "vaani-api",
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
