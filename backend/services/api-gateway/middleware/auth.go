package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/config"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type contextKey string

const (
	// UserIDKey is the key used to store the user ID in the request context
	UserIDKey contextKey = "userID"
	// UserEmailKey is the key used to store the user email in the request context
	UserEmailKey contextKey = "userEmail"
	// UserPhoneKey is the key used to store the user phone in the request context
	UserPhoneKey contextKey = "userPhone"
)

// FirebaseAuthMiddleware validates Firebase JWTs in the Authorization header
func FirebaseAuthMiddleware(cfg *config.AppConfig, logger *zap.Logger) func(http.Handler) http.Handler {
	// Initialize Firebase
	var app *firebase.App
	var err error

	if cfg.Firebase.CredentialsFile != "" {
		// Use service account credential file if provided
		opt := option.WithCredentialsFile(cfg.Firebase.CredentialsFile)
		app, err = firebase.NewApp(context.Background(), &firebase.Config{ProjectID: cfg.Firebase.ProjectID}, opt)
	} else {
		// Otherwise use application default credentials
		app, err = firebase.NewApp(context.Background(), &firebase.Config{ProjectID: cfg.Firebase.ProjectID})
	}

	if err != nil {
		logger.Fatal("Failed to initialize Firebase app", zap.Error(err))
	}

	// Get the Auth client
	client, err := app.Auth(context.Background())
	if err != nil {
		logger.Fatal("Failed to initialize Firebase Auth client", zap.Error(err))
	}

	// Return the middleware
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the authorization header
			authHeader := r.Header.Get("Authorization")
			idToken := extractToken(authHeader)

			if idToken == "" {
				http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
				return
			}

			// Verify the token
			token, err := client.VerifyIDToken(r.Context(), idToken)
			if err != nil {
				logger.Error("Failed to verify ID token", zap.Error(err))
				http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				return
			}

			// Extract user info
			userID := token.UID
			email, _ := token.Claims["email"].(string)
			phone, _ := token.Claims["phone_number"].(string)

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserEmailKey, email)
			ctx = context.WithValue(ctx, UserPhoneKey, phone)

			// Continue with the request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ExtractUserID retrieves the user ID from the request context
func ExtractUserID(r *http.Request) (string, error) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}

// ExtractUserEmail retrieves the user email from the request context
func ExtractUserEmail(r *http.Request) (string, error) {
	email, ok := r.Context().Value(UserEmailKey).(string)
	if !ok || email == "" {
		return "", errors.New("user email not found in context")
	}
	return email, nil
}

// ExtractUserPhone retrieves the user phone from the request context
func ExtractUserPhone(r *http.Request) (string, error) {
	phone, ok := r.Context().Value(UserPhoneKey).(string)
	if !ok || phone == "" {
		return "", errors.New("user phone not found in context")
	}
	return phone, nil
}

// extractToken extracts the JWT from the Authorization header
func extractToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
} 