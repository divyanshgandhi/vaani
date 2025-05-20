package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/models"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FirestoreJobRepository implements the JobRepository interface using Firestore
type FirestoreJobRepository struct {
	firestoreClient *FirestoreClient
	logger          *zap.Logger
}

// NewFirestoreJobRepository creates a new FirestoreJobRepository
func NewFirestoreJobRepository(client *FirestoreClient, logger *zap.Logger) *FirestoreJobRepository {
	return &FirestoreJobRepository{
		firestoreClient: client,
		logger:          logger,
	}
}

// jobsCollection returns the path to the jobs collection for a project
func (r *FirestoreJobRepository) jobsCollection(projectID string) string {
	return fmt.Sprintf("projects/%s/jobs", projectID)
}

// CreateJob creates a new job in Firestore
func (r *FirestoreJobRepository) CreateJob(job *models.Job) error {
	ctx := context.Background()

	// Set created time if not already set
	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now()
	}

	// Set status to queued if not already set
	if job.Status == "" {
		job.Status = models.StatusQueued
	}

	// Define the document reference
	docRef := r.firestoreClient.client.Collection(r.jobsCollection(job.ProjectID)).Doc(job.ID)

	// Create the job document
	_, err := docRef.Set(ctx, job)
	if err != nil {
		r.logger.Error("Failed to create job in Firestore",
			zap.String("job_id", job.ID),
			zap.String("project_id", job.ProjectID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create job: %w", err)
	}

	r.logger.Info("Job created in Firestore",
		zap.String("job_id", job.ID),
		zap.String("project_id", job.ProjectID),
		zap.String("status", string(job.Status)),
	)

	return nil
}

// GetJob retrieves a job from Firestore
func (r *FirestoreJobRepository) GetJob(projectID, jobID string) (*models.Job, error) {
	ctx := context.Background()

	// Get the job document
	docRef := r.firestoreClient.client.Collection(r.jobsCollection(projectID)).Doc(jobID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		// Check if document not found
		if status.Code(err) == codes.NotFound || strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("job not found: %s", jobID)
		}
		r.logger.Error("Failed to get job from Firestore",
			zap.String("job_id", jobID),
			zap.String("project_id", projectID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	// Convert to Job model
	var job models.Job
	if err := doc.DataTo(&job); err != nil {
		r.logger.Error("Failed to convert Firestore document to Job",
			zap.String("job_id", jobID),
			zap.String("project_id", projectID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to convert job data: %w", err)
	}

	return &job, nil
}

// UpdateJobStatus updates the status of a job in Firestore
func (r *FirestoreJobRepository) UpdateJobStatus(projectID, jobID string, status models.JobStatus) error {
	ctx := context.Background()

	// Get the document reference
	docRef := r.firestoreClient.client.Collection(r.jobsCollection(projectID)).Doc(jobID)

	// Prepare update data based on status
	updates := map[string]interface{}{
		"status": status,
	}

	// Add timestamp based on status
	switch status {
	case models.StatusProcessing:
		updates["started_at"] = time.Now()
	case models.StatusComplete, models.StatusFailed:
		updates["completed_at"] = time.Now()
	}

	// Update the document
	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: "status", Value: status},
	})
	if err != nil {
		// Check if document not found
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("job not found: %s", jobID)
		}
		r.logger.Error("Failed to update job status in Firestore",
			zap.String("job_id", jobID),
			zap.String("project_id", projectID),
			zap.String("status", string(status)),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update job status: %w", err)
	}

	r.logger.Info("Job status updated in Firestore",
		zap.String("job_id", jobID),
		zap.String("project_id", projectID),
		zap.String("status", string(status)),
	)

	return nil
}

// UpdateJobOutput updates the output URL of a job in Firestore
func (r *FirestoreJobRepository) UpdateJobOutput(projectID, jobID, outputURL string) error {
	ctx := context.Background()

	// Get the document reference
	docRef := r.firestoreClient.client.Collection(r.jobsCollection(projectID)).Doc(jobID)

	// Update the document
	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: "output_url", Value: outputURL},
		{Path: "status", Value: models.StatusComplete},
		{Path: "completed_at", Value: time.Now()},
	})
	if err != nil {
		// Check if document not found
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("job not found: %s", jobID)
		}
		r.logger.Error("Failed to update job output in Firestore",
			zap.String("job_id", jobID),
			zap.String("project_id", projectID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update job output: %w", err)
	}

	r.logger.Info("Job output updated in Firestore",
		zap.String("job_id", jobID),
		zap.String("project_id", projectID),
	)

	return nil
}

// ListJobs retrieves a list of jobs for a project from Firestore
func (r *FirestoreJobRepository) ListJobs(projectID string, limit int) ([]*models.Job, error) {
	ctx := context.Background()

	// Query the jobs collection
	query := r.firestoreClient.client.Collection(r.jobsCollection(projectID)).
		OrderBy("created_at", firestore.Desc).
		Limit(limit)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		r.logger.Error("Failed to list jobs from Firestore",
			zap.String("project_id", projectID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}

	// Convert documents to Job models
	jobs := make([]*models.Job, 0, len(docs))
	for _, doc := range docs {
		var job models.Job
		if err := doc.DataTo(&job); err != nil {
			r.logger.Error("Failed to convert Firestore document to Job",
				zap.String("doc_id", doc.Ref.ID),
				zap.Error(err),
			)
			continue
		}
		jobs = append(jobs, &job)
	}

	return jobs, nil
} 