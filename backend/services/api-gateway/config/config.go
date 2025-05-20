package config

import (
	"os"
	"strconv"
	"time"
)

// AppConfig holds the application configuration
type AppConfig struct {
	Port            string
	Env             string
	ShutdownTimeout time.Duration
	RateLimit       struct {
		RequestsPerSecond int
		ResetAfter        time.Duration
		Enabled           bool
	}
	LogLevel string
	Firebase struct {
		ProjectID              string
		CredentialsFile        string
		AuthEmulatorHost       string
		FirestoreEmulatorHost  string
		UseEmulator            bool
	}
}

// New creates a new AppConfig with values from environment variables
// or defaults if not set
func New() *AppConfig {
	config := &AppConfig{
		Port:            getEnv("PORT", "8080"),
		Env:             getEnv("APP_ENV", "development"),
		ShutdownTimeout: 10 * time.Second,
		LogLevel:        getEnv("LOG_LEVEL", "info"),
	}

	// Rate limiting config
	config.RateLimit.RequestsPerSecond = getEnvAsInt("RATE_LIMIT_REQUESTS_PER_SECOND", 5)
	config.RateLimit.ResetAfter = time.Duration(getEnvAsInt("RATE_LIMIT_RESET_AFTER_SECONDS", 60)) * time.Second
	config.RateLimit.Enabled = getEnvAsBool("RATE_LIMIT_ENABLED", true)

	// Firebase config
	config.Firebase.ProjectID = getEnv("FIREBASE_PROJECT_ID", "")
	config.Firebase.CredentialsFile = getEnv("FIREBASE_CREDENTIALS_FILE", "")
	config.Firebase.AuthEmulatorHost = getEnv("FIREBASE_AUTH_EMULATOR_HOST", "")
	config.Firebase.FirestoreEmulatorHost = getEnv("FIREBASE_FIRESTORE_EMULATOR_HOST", "")
	config.Firebase.UseEmulator = getEnv("FIREBASE_USE_EMULATOR", "false") == "true"

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	
	return value
}

// getEnvAsBool gets an environment variable as a boolean or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	
	return value
} 