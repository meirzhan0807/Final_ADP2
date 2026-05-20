package services

import (
	"context"
	"fmt"

	"jobfinder/job-service/internal/messaging"
	"jobfinder/job-service/internal/models"
	"jobfinder/job-service/internal/repositories"
)

type JobService interface {
	CreateJob(ctx context.Context, in *models.CreateJobInput) (*models.Job, error)
	GetJob(ctx context.Context, id string) (*models.Job, error)
	UpdateJob(ctx context.Context, id string, in *models.UpdateJobInput) (*models.Job, error)
	DeleteJob(ctx context.Context, id, employerID string) error
	ListJobs(ctx context.Context, page, limit int, status string) ([]*models.Job, int, error)
	SearchJobs(ctx context.Context, in *models.SearchInput) ([]*models.Job, int, error)
	GetJobsByEmployer(ctx context.Context, employerID string, page, limit int) ([]*models.Job, int, error)
	ApplyJob(ctx context.Context, jobID, applicantID, coverLetter, resumeURL string) (*models.Application, error)
	GetApplications(ctx context.Context, jobID string, page, limit int) ([]*models.Application, int, error)
	GetApplicationsByUser(ctx context.Context, userID string, page, limit int) ([]*models.Application, int, error)
	UpdateApplicationStatus(ctx context.Context, appID, status string) (*models.Application, error)
	GetCategories() []string
}

type jobService struct {
	jobRepo repositories.JobRepository
	appRepo repositories.ApplicationRepository
	nats    *messaging.NATSClient
}

func NewJobService(jr repositories.JobRepository, ar repositories.ApplicationRepository, nc *messaging.NATSClient) JobService {
	return &jobService{jobRepo: jr, appRepo: ar, nats: nc}
}

func (s *jobService) CreateJob(ctx context.Context, in *models.CreateJobInput) (*models.Job, error) {
	if in.Title == "" || in.EmployerID == "" {
		return nil, fmt.Errorf("title and employer_id required")
	}
	job, err := s.jobRepo.Create(ctx, in)
	if err != nil { return nil, err }
	s.nats.Publish("job.created", map[string]interface{}{
		"job_id": job.ID, "title": job.Title, "category": job.Category, "location": job.Location,
	})
	return job, nil
}

func (s *jobService) GetJob(ctx context.Context, id string) (*models.Job, error) {
	return s.jobRepo.GetByID(ctx, id)
}

func (s *jobService) UpdateJob(ctx context.Context, id string, in *models.UpdateJobInput) (*models.Job, error) {
	return s.jobRepo.Update(ctx, id, in)
}

func (s *jobService) DeleteJob(ctx context.Context, id, employerID string) error {
	return s.jobRepo.Delete(ctx, id, employerID)
}

func (s *jobService) ListJobs(ctx context.Context, page, limit int, status string) ([]*models.Job, int, error) {
	if page <= 0 { page = 1 }
	if limit <= 0 || limit > 100 { limit = 20 }
	return s.jobRepo.List(ctx, page, limit, status)
}

func (s *jobService) SearchJobs(ctx context.Context, in *models.SearchInput) ([]*models.Job, int, error) {
	if in.Page <= 0 { in.Page = 1 }
	if in.Limit <= 0 || in.Limit > 100 { in.Limit = 20 }
	return s.jobRepo.Search(ctx, in)
}

func (s *jobService) GetJobsByEmployer(ctx context.Context, employerID string, page, limit int) ([]*models.Job, int, error) {
	if page <= 0 { page = 1 }
	if limit <= 0 { limit = 20 }
	return s.jobRepo.GetByEmployer(ctx, employerID, page, limit)
}

func (s *jobService) ApplyJob(ctx context.Context, jobID, applicantID, coverLetter, resumeURL string) (*models.Application, error) {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil { return nil, fmt.Errorf("job not found") }
	if job.Status != "active" { return nil, fmt.Errorf("job is not accepting applications") }

	app, err := s.appRepo.Create(ctx, jobID, applicantID, coverLetter, resumeURL)
	if err != nil { return nil, err }

	s.jobRepo.IncrementApplicationCount(ctx, jobID)
	s.nats.Publish("job.applied", map[string]interface{}{
		"application_id": app.ID, "job_id": jobID, "job_title": job.Title,
		"applicant_id": applicantID, "employer_id": job.EmployerID,
	})
	return app, nil
}

func (s *jobService) GetApplications(ctx context.Context, jobID string, page, limit int) ([]*models.Application, int, error) {
	if page <= 0 { page = 1 }
	if limit <= 0 { limit = 20 }
	return s.appRepo.GetByJob(ctx, jobID, page, limit)
}

func (s *jobService) GetApplicationsByUser(ctx context.Context, userID string, page, limit int) ([]*models.Application, int, error) {
	if page <= 0 { page = 1 }
	if limit <= 0 { limit = 20 }
	return s.appRepo.GetByUser(ctx, userID, page, limit)
}

func (s *jobService) UpdateApplicationStatus(ctx context.Context, appID, status string) (*models.Application, error) {
	valid := map[string]bool{"pending": true, "reviewed": true, "accepted": true, "rejected": true}
	if !valid[status] { return nil, fmt.Errorf("invalid status") }
	app, err := s.appRepo.UpdateStatus(ctx, appID, status)
	if err != nil { return nil, err }
	s.nats.Publish("application.status_updated", map[string]interface{}{
		"application_id": appID, "applicant_id": app.ApplicantID, "status": status,
	})
	return app, nil
}

func (s *jobService) GetCategories() []string {
	return []string{"Technology", "Finance", "Healthcare", "Education", "Marketing",
		"Design", "Engineering", "Sales", "HR", "Legal", "Operations", "Customer Service",
		"Data Science", "Product Management", "Other"}
}
