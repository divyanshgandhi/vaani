package services

import (
	"context"
	"time"

	"github.com/NavoDayAI/vaani/backend/services/api-gateway/models"
	"github.com/NavoDayAI/vaani/backend/services/api-gateway/utils/langdetect"
	"go.uber.org/zap"
)

// JobProcessor handles background processing of TTS jobs
type JobProcessor struct {
	jobRepo   models.JobRepository
	ttsClient *TTSClient
	logger    *zap.Logger
	stopCh    chan struct{}
}

// NewJobProcessor creates a new job processor
func NewJobProcessor(jobRepo models.JobRepository, ttsClient *TTSClient, logger *zap.Logger) *JobProcessor {
	return &JobProcessor{
		jobRepo:   jobRepo,
		ttsClient: ttsClient,
		logger:    logger,
		stopCh:    make(chan struct{}),
	}
}

// Start begins processing jobs in the background
func (p *JobProcessor) Start() {
	p.logger.Info("Starting job processor")

	ticker := time.NewTicker(5 * time.Second) // Check for new jobs every 5 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.processQueuedJobs()
		case <-p.stopCh:
			p.logger.Info("Job processor stopped")
			return
		}
	}
}

// Stop stops the job processor
func (p *JobProcessor) Stop() {
	close(p.stopCh)
}

// processQueuedJobs finds and processes queued jobs
func (p *JobProcessor) processQueuedJobs() {
	// For now, we'll implement a simple approach
	// In a production system, this would use a proper job queue like Redis or Cloud Tasks

	// This is a simplified implementation - in reality you'd want to:
	// 1. Use a proper job queue system
	// 2. Handle concurrent processing
	// 3. Implement retry logic
	// 4. Add proper error handling and monitoring

	p.logger.Debug("Checking for queued jobs")
}

// ProcessJob processes a single TTS job
func (p *JobProcessor) ProcessJob(job *models.Job) error {
	ctx := context.Background()

	p.logger.Info("Processing job",
		zap.String("job_id", job.ID),
		zap.String("text", job.Text[:min(50, len(job.Text))]),
	)

	// Update job status to processing
	err := p.jobRepo.UpdateJobStatus(job.ProjectID, job.ID, models.StatusProcessing)
	if err != nil {
		p.logger.Error("Failed to update job status to processing",
			zap.String("job_id", job.ID),
			zap.Error(err),
		)
		return err
	}

	// Detect language if not provided
	language := job.Language
	if language == "" {
		detectedLang, err := langdetect.Detect(job.Text)
		if err != nil {
			p.logger.Warn("Failed to detect language, defaulting to English",
				zap.String("job_id", job.ID),
				zap.Error(err),
			)
			language = "en"
		} else {
			language = detectedLang
		}
	}

	// Prepare TTS request
	ttsReq := TTSRequest{
		Text:     job.Text,
		VoiceID:  job.VoiceID,
		Speed:    job.Speed,
		Pitch:    job.Pitch,
		Emotion:  job.Emotion,
		Language: language,
		Format:   job.OutputType,
	}

	// Generate TTS
	ttsResp, err := p.ttsClient.GenerateTTS(ctx, ttsReq)
	if err != nil {
		p.logger.Error("Failed to generate TTS",
			zap.String("job_id", job.ID),
			zap.Error(err),
		)

		// Update job status to failed
		p.jobRepo.UpdateJobStatus(job.ProjectID, job.ID, models.StatusFailed)
		return err
	}

	// Update job with output URL
	err = p.jobRepo.UpdateJobOutput(job.ProjectID, job.ID, ttsResp.AudioURL)
	if err != nil {
		p.logger.Error("Failed to update job output",
			zap.String("job_id", job.ID),
			zap.Error(err),
		)
		return err
	}

	// Update job status to complete
	err = p.jobRepo.UpdateJobStatus(job.ProjectID, job.ID, models.StatusComplete)
	if err != nil {
		p.logger.Error("Failed to update job status to complete",
			zap.String("job_id", job.ID),
			zap.Error(err),
		)
		return err
	}

	p.logger.Info("Job completed successfully",
		zap.String("job_id", job.ID),
		zap.String("output_url", ttsResp.AudioURL),
	)

	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
