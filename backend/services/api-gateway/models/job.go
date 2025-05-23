package models

import (
	"time"
)

// JobStatus represents the status of a TTS job
type JobStatus string

const (
	StatusQueued     JobStatus = "queued"
	StatusProcessing JobStatus = "processing"
	StatusComplete   JobStatus = "complete"
	StatusFailed     JobStatus = "failed"
)

// Job represents a TTS job in the system
type Job struct {
	ID           string    `firestore:"id"`
	ProjectID    string    `firestore:"project_id"`
	UserID       string    `firestore:"user_id"`
	Text         string    `firestore:"text"`
	VoiceID      string    `firestore:"voice_id"`
	Speed        float64   `firestore:"speed,omitempty"`
	Pitch        float64   `firestore:"pitch,omitempty"`
	Emotion      int       `firestore:"emotion,omitempty"`
	Language     string    `firestore:"language,omitempty"`
	OutputType   string    `firestore:"output_type,omitempty"`
	TextLength   int       `firestore:"text_length"`
	Status       JobStatus `firestore:"status"`
	CreatedAt    time.Time `firestore:"created_at"`
	StartedAt    time.Time `firestore:"started_at,omitempty"`
	CompletedAt  time.Time `firestore:"completed_at,omitempty"`
	OutputURL    string    `firestore:"output_url,omitempty"`
	ErrorMessage string    `firestore:"error_message,omitempty"`
}

// JobRepository defines methods for working with jobs
type JobRepository interface {
	CreateJob(job *Job) error
	GetJob(projectID, jobID string) (*Job, error)
	UpdateJobStatus(projectID, jobID string, status JobStatus) error
	UpdateJobOutput(projectID, jobID, outputURL string) error
	ListJobs(projectID string, limit int) ([]*Job, error)
}

// ProjectRepository defines methods for working with projects
type ProjectRepository interface {
	CreateProject(project *Project) error
	GetProject(userID, projectID string) (*Project, error)
	UpdateProject(project *Project) error
	DeleteProject(userID, projectID string) error
	ListUserProjects(userID string, limit int) ([]*Project, error)
}

// Project represents a user project
type Project struct {
	ID           string    `firestore:"id"`
	UserID       string    `firestore:"user_id"`
	Name         string    `firestore:"name"`
	Description  string    `firestore:"description,omitempty"`
	CreatedAt    time.Time `firestore:"created_at"`
	UpdatedAt    time.Time `firestore:"updated_at"`
	JobCount     int       `firestore:"job_count"`
	TotalCredits int       `firestore:"total_credits"`
	UsedCredits  int       `firestore:"used_credits"`
}
