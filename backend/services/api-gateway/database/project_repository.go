package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/NavoDayAI/vaani/backend/services/api-gateway/models"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProjectRepository defines methods for working with projects
type ProjectRepository interface {
	CreateProject(project *models.Project) error
	GetProject(userID, projectID string) (*models.Project, error)
	UpdateProject(project *models.Project) error
	DeleteProject(userID, projectID string) error
	ListUserProjects(userID string, limit int) ([]*models.Project, error)
}

// FirestoreProjectRepository implements the ProjectRepository interface using Firestore
type FirestoreProjectRepository struct {
	firestoreClient *FirestoreClient
	logger          *zap.Logger
}

// NewFirestoreProjectRepository creates a new FirestoreProjectRepository
func NewFirestoreProjectRepository(client *FirestoreClient, logger *zap.Logger) *FirestoreProjectRepository {
	return &FirestoreProjectRepository{
		firestoreClient: client,
		logger:          logger,
	}
}

// projectsCollection returns the path to the projects collection
func (r *FirestoreProjectRepository) projectsCollection() string {
	return "projects"
}

// CreateProject creates a new project in Firestore
func (r *FirestoreProjectRepository) CreateProject(project *models.Project) error {
	ctx := context.Background()

	// Set created/updated times if not already set
	now := time.Now()
	if project.CreatedAt.IsZero() {
		project.CreatedAt = now
	}
	if project.UpdatedAt.IsZero() {
		project.UpdatedAt = now
	}

	// Initialize counters if not set
	if project.JobCount < 0 {
		project.JobCount = 0
	}
	if project.TotalCredits < 0 {
		project.TotalCredits = 0
	}
	if project.UsedCredits < 0 {
		project.UsedCredits = 0
	}

	// Define the document reference
	docRef := r.firestoreClient.client.Collection(r.projectsCollection()).Doc(project.ID)

	// Create the project document
	_, err := docRef.Set(ctx, project)
	if err != nil {
		r.logger.Error("Failed to create project in Firestore",
			zap.String("project_id", project.ID),
			zap.String("user_id", project.UserID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create project: %w", err)
	}

	r.logger.Info("Project created in Firestore",
		zap.String("project_id", project.ID),
		zap.String("user_id", project.UserID),
		zap.String("name", project.Name),
	)

	return nil
}

// GetProject retrieves a project from Firestore
func (r *FirestoreProjectRepository) GetProject(userID, projectID string) (*models.Project, error) {
	ctx := context.Background()

	// Get the project document
	docRef := r.firestoreClient.client.Collection(r.projectsCollection()).Doc(projectID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		// Check if document not found
		if status.Code(err) == codes.NotFound || strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("project not found: %s", projectID)
		}
		r.logger.Error("Failed to get project from Firestore",
			zap.String("project_id", projectID),
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Convert to Project model
	var project models.Project
	if err := doc.DataTo(&project); err != nil {
		r.logger.Error("Failed to convert Firestore document to Project",
			zap.String("project_id", projectID),
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to convert project data: %w", err)
	}

	// Check if the project belongs to the user
	if project.UserID != userID {
		return nil, fmt.Errorf("project not found: %s", projectID)
	}

	return &project, nil
}

// UpdateProject updates an existing project in Firestore
func (r *FirestoreProjectRepository) UpdateProject(project *models.Project) error {
	ctx := context.Background()

	// Update the updated_at timestamp
	project.UpdatedAt = time.Now()

	// Get the document reference
	docRef := r.firestoreClient.client.Collection(r.projectsCollection()).Doc(project.ID)

	// Update the document
	_, err := docRef.Set(ctx, project)
	if err != nil {
		r.logger.Error("Failed to update project in Firestore",
			zap.String("project_id", project.ID),
			zap.String("user_id", project.UserID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update project: %w", err)
	}

	r.logger.Info("Project updated in Firestore",
		zap.String("project_id", project.ID),
		zap.String("user_id", project.UserID),
		zap.String("name", project.Name),
	)

	return nil
}

// DeleteProject deletes a project from Firestore
func (r *FirestoreProjectRepository) DeleteProject(userID, projectID string) error {
	ctx := context.Background()

	// First verify the project belongs to the user
	project, err := r.GetProject(userID, projectID)
	if err != nil {
		return err // GetProject already returns appropriate error
	}

	if project.UserID != userID {
		return fmt.Errorf("project not found: %s", projectID)
	}

	// Delete the project document
	docRef := r.firestoreClient.client.Collection(r.projectsCollection()).Doc(projectID)
	_, err = docRef.Delete(ctx)
	if err != nil {
		r.logger.Error("Failed to delete project from Firestore",
			zap.String("project_id", projectID),
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete project: %w", err)
	}

	r.logger.Info("Project deleted from Firestore",
		zap.String("project_id", projectID),
		zap.String("user_id", userID),
	)

	return nil
}

// ListUserProjects retrieves a list of projects for a user from Firestore
func (r *FirestoreProjectRepository) ListUserProjects(userID string, limit int) ([]*models.Project, error) {
	ctx := context.Background()

	// Simplified query - just filter by user_id without ordering (to avoid composite index requirement)
	// We'll sort in memory instead
	query := r.firestoreClient.client.Collection(r.projectsCollection()).
		Where("user_id", "==", userID).
		Limit(limit)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		r.logger.Error("Failed to list projects from Firestore",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	// Convert documents to Project models
	projects := make([]*models.Project, 0, len(docs))
	for _, doc := range docs {
		var project models.Project
		if err := doc.DataTo(&project); err != nil {
			r.logger.Error("Failed to convert Firestore document to Project",
				zap.String("doc_id", doc.Ref.ID),
				zap.String("user_id", userID),
				zap.Error(err),
			)
			continue
		}
		projects = append(projects, &project)
	}

	// Sort by updated_at in memory (descending - newest first)
	for i := 0; i < len(projects)-1; i++ {
		for j := i + 1; j < len(projects); j++ {
			if projects[i].UpdatedAt.Before(projects[j].UpdatedAt) {
				projects[i], projects[j] = projects[j], projects[i]
			}
		}
	}

	r.logger.Info("Projects listed from Firestore",
		zap.String("user_id", userID),
		zap.Int("count", len(projects)),
	)

	return projects, nil
}
