package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "jobfinder/notification-service/internal/pb"
	"jobfinder/notification-service/internal/services"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
	svc services.NotificationService
}

func NewNotificationHandler(svc services.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) mustEmbedUnimplementedNotificationServiceServer() {}

func (h *NotificationHandler) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {
	if req.To == "" || req.Subject == "" {
		return nil, status.Error(codes.InvalidArgument, "to and subject required")
	}
	if err := h.svc.SendEmail(ctx, req.To, req.Subject, req.Body, req.IsHtml); err != nil {
		return &pb.SendEmailResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.SendEmailResponse{Success: true, Message: "sent"}, nil
}

func (h *NotificationHandler) SendWelcomeEmail(ctx context.Context, req *pb.WelcomeEmailRequest) (*pb.SendEmailResponse, error) {
	if err := h.svc.SendWelcomeEmail(ctx, req.To, req.UserName, req.VerificationLink); err != nil {
		return &pb.SendEmailResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.SendEmailResponse{Success: true, Message: "welcome email sent"}, nil
}

func (h *NotificationHandler) SendJobApplicationEmail(ctx context.Context, req *pb.JobApplicationEmailRequest) (*pb.SendEmailResponse, error) {
	if err := h.svc.SendJobApplicationEmail(ctx, req.EmployerEmail, req.ApplicantName, req.JobTitle, req.ApplicationId); err != nil {
		return &pb.SendEmailResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.SendEmailResponse{Success: true, Message: "job application email sent"}, nil
}

func (h *NotificationHandler) SendApplicationStatusEmail(ctx context.Context, req *pb.ApplicationStatusEmailRequest) (*pb.SendEmailResponse, error) {
	if err := h.svc.SendApplicationStatusEmail(ctx, req.ApplicantEmail, req.ApplicantName, req.JobTitle, req.Status, req.Message); err != nil {
		return &pb.SendEmailResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.SendEmailResponse{Success: true, Message: "status email sent"}, nil
}

func (h *NotificationHandler) SendPasswordResetEmail(ctx context.Context, req *pb.PasswordResetEmailRequest) (*pb.SendEmailResponse, error) {
	if err := h.svc.SendPasswordResetEmail(ctx, req.To, req.UserName, req.ResetLink); err != nil {
		return &pb.SendEmailResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.SendEmailResponse{Success: true, Message: "password reset email sent"}, nil
}

func (h *NotificationHandler) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id required")
	}
	items, total, err := h.svc.GetNotifications(ctx, req.UserId, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var list []*pb.NotificationResponse
	for _, n := range items {
		list = append(list, &pb.NotificationResponse{
			Id: n.ID, UserId: n.UserID, Title: n.Title, Message: n.Message,
			Type: n.Type, IsRead: n.IsRead, CreatedAt: n.CreatedAt,
		})
	}
	return &pb.GetNotificationsResponse{Notifications: list, Total: int32(total)}, nil
}

func (h *NotificationHandler) MarkNotificationRead(ctx context.Context, req *pb.MarkNotificationReadRequest) (*pb.MarkNotificationReadResponse, error) {
	if err := h.svc.MarkRead(ctx, req.NotificationId, req.UserId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.MarkNotificationReadResponse{Success: true}, nil
}

func (h *NotificationHandler) GetUnreadCount(ctx context.Context, req *pb.GetUnreadCountRequest) (*pb.GetUnreadCountResponse, error) {
	count, err := h.svc.GetUnreadCount(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetUnreadCountResponse{Count: int32(count)}, nil
}

func (h *NotificationHandler) DeleteNotification(ctx context.Context, req *pb.DeleteNotificationRequest) (*pb.DeleteNotificationResponse, error) {
	if err := h.svc.Delete(ctx, req.NotificationId, req.UserId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteNotificationResponse{Success: true}, nil
}

func (h *NotificationHandler) SendBulkEmail(ctx context.Context, req *pb.SendBulkEmailRequest) (*pb.SendEmailResponse, error) {
	if err := h.svc.SendBulkEmail(ctx, req.Recipients, req.Subject, req.Body); err != nil {
		return &pb.SendEmailResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.SendEmailResponse{Success: true, Message: "bulk emails sent"}, nil
}

func (h *NotificationHandler) SendJobAlertEmail(ctx context.Context, req *pb.JobAlertEmailRequest) (*pb.SendEmailResponse, error) {
	if err := h.svc.SendJobAlertEmail(ctx, req.To, req.UserName, req.JobTitles, req.SearchQuery); err != nil {
		return &pb.SendEmailResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.SendEmailResponse{Success: true, Message: "job alert email sent"}, nil
}

func (h *NotificationHandler) CreateNotification(ctx context.Context, req *pb.CreateNotificationRequest) (*pb.NotificationResponse, error) {
	n, err := h.svc.CreateNotification(ctx, req.UserId, req.Title, req.Message, req.Type)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.NotificationResponse{
		Id: n.ID, UserId: n.UserID, Title: n.Title,
		Message: n.Message, Type: n.Type, IsRead: n.IsRead, CreatedAt: n.CreatedAt,
	}, nil
}
