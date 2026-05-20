package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"jobfinder/notification-service/internal/email"
	"jobfinder/notification-service/internal/models"
	"jobfinder/notification-service/internal/repositories"
)

type NotificationService interface {
	SendEmail(ctx context.Context, to, subject, body string, isHTML bool) error
	SendWelcomeEmail(ctx context.Context, to, name, link string) error
	SendJobApplicationEmail(ctx context.Context, employerEmail, applicantName, jobTitle, appID string) error
	SendApplicationStatusEmail(ctx context.Context, email, name, jobTitle, status, message string) error
	SendPasswordResetEmail(ctx context.Context, to, name, link string) error
	SendJobAlertEmail(ctx context.Context, to, name string, titles []string, query string) error
	SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error
	CreateNotification(ctx context.Context, userID, title, message, nType string) (*models.Notification, error)
	GetNotifications(ctx context.Context, userID string, page, limit int) ([]*models.Notification, int, error)
	MarkRead(ctx context.Context, id, userID string) error
	Delete(ctx context.Context, id, userID string) error
	GetUnreadCount(ctx context.Context, userID string) (int, error)
	StartEventListeners(nc *nats.Conn)
}

type notifService struct {
	repo  repositories.NotificationRepository
	email *email.Service
}

func NewNotificationService(repo repositories.NotificationRepository, emailSvc *email.Service) NotificationService {
	return &notifService{repo: repo, email: emailSvc}
}

func (s *notifService) SendEmail(_ context.Context, to, subject, body string, isHTML bool) error {
	return s.email.Send(to, subject, body, isHTML)
}
func (s *notifService) SendWelcomeEmail(_ context.Context, to, name, link string) error {
	return s.email.SendWelcome(to, name, link)
}
func (s *notifService) SendJobApplicationEmail(_ context.Context, empEmail, appName, jobTitle, appID string) error {
	return s.email.SendJobApplied(empEmail, appName, jobTitle, appID)
}
func (s *notifService) SendApplicationStatusEmail(_ context.Context, to, name, jobTitle, status, msg string) error {
	return s.email.SendStatusUpdate(to, name, jobTitle, status, msg)
}
func (s *notifService) SendPasswordResetEmail(_ context.Context, to, name, link string) error {
	return s.email.SendPasswordReset(to, name, link)
}
func (s *notifService) SendJobAlertEmail(_ context.Context, to, name string, titles []string, query string) error {
	return s.email.SendJobAlert(to, name, titles, query)
}
func (s *notifService) SendBulkEmail(_ context.Context, recipients []string, subject, body string) error {
	for _, r := range recipients {
		if err := s.email.Send(r, subject, body, true); err != nil {
			log.Printf("bulk email to %s failed: %v", r, err)
		}
	}
	return nil
}
func (s *notifService) CreateNotification(ctx context.Context, userID, title, message, nType string) (*models.Notification, error) {
	return s.repo.Create(ctx, userID, title, message, nType)
}
func (s *notifService) GetNotifications(ctx context.Context, userID string, page, limit int) ([]*models.Notification, int, error) {
	if page <= 0 { page = 1 }
	if limit <= 0 { limit = 20 }
	return s.repo.GetByUserID(ctx, userID, page, limit)
}
func (s *notifService) MarkRead(ctx context.Context, id, userID string) error {
	return s.repo.MarkRead(ctx, id, userID)
}
func (s *notifService) Delete(ctx context.Context, id, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}
func (s *notifService) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

func (s *notifService) StartEventListeners(nc *nats.Conn) {
	ctx := context.Background()

	nc.Subscribe("user.registered", func(msg *nats.Msg) {
		var ev map[string]string
		if err := json.Unmarshal(msg.Data, &ev); err != nil { return }
		s.email.SendWelcome(ev["email"], ev["name"], ev["verify_link"])
		s.repo.Create(ctx, ev["user_id"], "Welcome to JobFinder!", "Your account was created. Please verify your email.", "welcome")
	})

	nc.Subscribe("user.password_reset", func(msg *nats.Msg) {
		var ev map[string]string
		if err := json.Unmarshal(msg.Data, &ev); err != nil { return }
		s.email.SendPasswordReset(ev["email"], ev["name"], ev["reset_link"])
	})

	nc.Subscribe("job.applied", func(msg *nats.Msg) {
		var ev map[string]interface{}
		if err := json.Unmarshal(msg.Data, &ev); err != nil { return }
		applicantID, _ := ev["applicant_id"].(string)
		jobTitle, _ := ev["job_title"].(string)
		appID, _ := ev["application_id"].(string)
		if applicantID != "" {
			s.repo.Create(ctx, applicantID,
				"Application Submitted",
				"Your application for '"+jobTitle+"' (ID: "+appID+") was submitted successfully.",
				"application")
		}
	})

	nc.Subscribe("application.status_updated", func(msg *nats.Msg) {
		var ev map[string]interface{}
		if err := json.Unmarshal(msg.Data, &ev); err != nil { return }
		applicantID, _ := ev["applicant_id"].(string)
		status, _ := ev["status"].(string)
		if applicantID != "" {
			s.repo.Create(ctx, applicantID,
				"Application Status Updated",
				"Your application status changed to: "+status,
				"status_update")
		}
	})

	nc.Subscribe("job.created", func(msg *nats.Msg) {
		var ev map[string]interface{}
		if err := json.Unmarshal(msg.Data, &ev); err != nil { return }
		log.Printf("[notification] New job created: %v", ev["title"])
	})

	log.Println("NATS event listeners started")
}
