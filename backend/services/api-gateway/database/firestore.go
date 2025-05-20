package database

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/config"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

// FirestoreClient wraps the Firestore client and provides utility methods
type FirestoreClient struct {
	client *firestore.Client
	logger *zap.Logger
}

// NewFirestoreClient creates and returns a new FirestoreClient
func NewFirestoreClient(cfg *config.AppConfig, logger *zap.Logger) (*FirestoreClient, error) {
	ctx := context.Background()
	
	// Validate configuration
	if cfg.Firebase.ProjectID == "" {
		return nil, fmt.Errorf("project id is required to access Firestore")
	}
	
	// Initialize Firebase app
	var app *firebase.App
	var err error
	
	// Create Firebase config with project ID
	fbConfig := &firebase.Config{
		ProjectID: cfg.Firebase.ProjectID,
	}
	
	if cfg.Firebase.CredentialsFile != "" {
		// Use service account credentials file
		logger.Info("Using Firebase credentials file", 
			zap.String("file", cfg.Firebase.CredentialsFile),
			zap.String("project_id", cfg.Firebase.ProjectID))
		
		opt := option.WithCredentialsFile(cfg.Firebase.CredentialsFile)
		app, err = firebase.NewApp(ctx, fbConfig, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
		}
	} else {
		// Use application default credentials
		logger.Info("Using application default credentials",
			zap.String("project_id", cfg.Firebase.ProjectID))
		
		app, err = firebase.NewApp(ctx, fbConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
		}
	}
	
	// Initialize Firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firestore client: %w", err)
	}
	
	return &FirestoreClient{
		client: client,
		logger: logger,
	}, nil
}

// Close closes the Firestore client
func (fc *FirestoreClient) Close() error {
	if fc.client != nil {
		return fc.client.Close()
	}
	return nil
}

// GetClient returns the underlying Firestore client
func (fc *FirestoreClient) GetClient() *firestore.Client {
	return fc.client
} 