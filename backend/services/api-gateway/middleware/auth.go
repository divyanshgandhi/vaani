package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/config"
	"github.com/golang-jwt/jwt/v4"
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

// JWT secret - in production, use environment variable or secure key management
var jwtSecret = []byte("your-secret-key-change-this-in-production")

// FirebaseAuthMiddleware validates Firebase JWTs in the Authorization header
func FirebaseAuthMiddleware(cfg *config.AppConfig, logger *zap.Logger) func(http.Handler) http.Handler {
	return firebaseAuthMiddlewareWithInit(cfg, logger, false)
}

// TestFirebaseAuthMiddleware is used for testing - it doesn't initialize a real Firebase client
func TestFirebaseAuthMiddleware(cfg *config.AppConfig, logger *zap.Logger) func(http.Handler) http.Handler {
	return firebaseAuthMiddlewareWithInit(cfg, logger, true)
}

// firebaseAuthMiddlewareWithInit is the internal implementation with a test mode flag
func firebaseAuthMiddlewareWithInit(cfg *config.AppConfig, logger *zap.Logger, testMode bool) func(http.Handler) http.Handler {
	// Initialize Firebase app only if not in test mode
	var app *firebase.App
	var client *auth.Client
	var err error

	if !testMode {
		// Validate configuration
		if cfg.Firebase.ProjectID == "" {
			logger.Fatal("Firebase project ID is required")
		}

		// Initialize Firebase app
		if cfg.Firebase.CredentialsFile != "" {
			// Use service account credential file if provided
			logger.Info("Using Firebase credentials file",
				zap.String("file", cfg.Firebase.CredentialsFile),
				zap.String("project_id", cfg.Firebase.ProjectID))

			opt := option.WithCredentialsFile(cfg.Firebase.CredentialsFile)
			app, err = firebase.NewApp(context.Background(), &firebase.Config{ProjectID: cfg.Firebase.ProjectID}, opt)
		} else {
			// Otherwise use application default credentials
			logger.Info("Using application default credentials",
				zap.String("project_id", cfg.Firebase.ProjectID))

			app, err = firebase.NewApp(context.Background(), &firebase.Config{ProjectID: cfg.Firebase.ProjectID})
		}

		if err != nil {
			logger.Fatal("Failed to initialize Firebase app", zap.Error(err))
		}

		// Get the Auth client
		client, err = app.Auth(context.Background())
		if err != nil {
			logger.Fatal("Failed to initialize Firebase Auth client", zap.Error(err))
		}
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

			// In test mode, always reject the token unless it's "test-valid-token"
			if testMode {
				if idToken == "test-valid-token" {
					// Add dummy user info to context
					ctx := context.WithValue(r.Context(), UserIDKey, "test-user-id")
					ctx = context.WithValue(ctx, UserEmailKey, "test@example.com")
					ctx = context.WithValue(ctx, UserPhoneKey, "+1234567890")
					next.ServeHTTP(w, r.WithContext(ctx))
				} else {
					http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				}
				return
			}

			// Try to verify as custom JWT first
			userID, email, phone, err := verifyCustomJWT(idToken, logger)
			if err != nil {
				// If custom JWT verification fails, try Firebase JWT
				token, fbErr := client.VerifyIDToken(r.Context(), idToken)
				if fbErr != nil {
					logger.Error("Failed to verify both custom and Firebase tokens",
						zap.Error(err),
						zap.Error(fbErr))
					http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
					return
				}

				// Extract user info from Firebase token
				userID = token.UID
				email, _ = token.Claims["email"].(string)
				phone, _ = token.Claims["phone_number"].(string)
			}

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

// verifyCustomJWT verifies our custom JWT tokens and returns user info
func verifyCustomJWT(tokenString string, logger *zap.Logger) (userID, email, phone string, err error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure token's signature algorithm is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", "", "", err
	}

	// Check if token is valid
	if !token.Valid {
		return "", "", "", errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", "", errors.New("invalid token claims")
	}

	// Extract phone number (which is our user ID for custom JWT)
	phoneNumber, ok := claims["phone_number"].(string)
	if !ok || phoneNumber == "" {
		return "", "", "", errors.New("phone number not found in token")
	}

	// For custom JWT, use phone number as both userID and phone
	return phoneNumber, "", phoneNumber, nil
}
