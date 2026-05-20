package notifpb

import (
	"context"
	"time"

	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


// ── Types ─────────────────────────────────────────────────────────────────

type SendEmailRequest struct {
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	IsHtml   bool   `json:"is_html"`
	FromName string `json:"from_name"`
}
type SendEmailResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	MessageId string `json:"message_id"`
}
type WelcomeEmailRequest struct {
	To               string `json:"to"`
	UserName         string `json:"user_name"`
	VerificationLink string `json:"verification_link"`
}
type JobApplicationEmailRequest struct {
	EmployerEmail string `json:"employer_email"`
	ApplicantName string `json:"applicant_name"`
	JobTitle      string `json:"job_title"`
	ApplicationId string `json:"application_id"`
}
type ApplicationStatusEmailRequest struct {
	ApplicantEmail string `json:"applicant_email"`
	ApplicantName  string `json:"applicant_name"`
	JobTitle       string `json:"job_title"`
	Status         string `json:"status"`
	Message        string `json:"message"`
}
type PasswordResetEmailRequest struct {
	To        string `json:"to"`
	UserName  string `json:"user_name"`
	ResetLink string `json:"reset_link"`
}
type GetNotificationsRequest struct {
	UserId string `json:"user_id"`
	Page   int32  `json:"page"`
	Limit  int32  `json:"limit"`
}
type NotificationResponse struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
type GetNotificationsResponse struct {
	Notifications []*NotificationResponse `json:"notifications"`
	Total         int32                   `json:"total"`
}
type MarkNotificationReadRequest struct {
	NotificationId string `json:"notification_id"`
	UserId         string `json:"user_id"`
}
type MarkNotificationReadResponse struct{ Success bool `json:"success"` }
type GetUnreadCountRequest struct{ UserId string `json:"user_id"` }
type GetUnreadCountResponse struct{ Count int32 `json:"count"` }
type DeleteNotificationRequest struct {
	NotificationId string `json:"notification_id"`
	UserId         string `json:"user_id"`
}
type DeleteNotificationResponse struct{ Success bool `json:"success"` }
type SendBulkEmailRequest struct {
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
}
type JobAlertEmailRequest struct {
	To          string   `json:"to"`
	UserName    string   `json:"user_name"`
	JobTitles   []string `json:"job_titles"`
	SearchQuery string   `json:"search_query"`
}
type CreateNotificationRequest struct {
	UserId  string `json:"user_id"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// ── Server Interface ──────────────────────────────────────────────────────

type NotificationServiceServer interface {
	SendEmail(context.Context, *SendEmailRequest) (*SendEmailResponse, error)
	SendWelcomeEmail(context.Context, *WelcomeEmailRequest) (*SendEmailResponse, error)
	SendJobApplicationEmail(context.Context, *JobApplicationEmailRequest) (*SendEmailResponse, error)
	SendApplicationStatusEmail(context.Context, *ApplicationStatusEmailRequest) (*SendEmailResponse, error)
	SendPasswordResetEmail(context.Context, *PasswordResetEmailRequest) (*SendEmailResponse, error)
	GetNotifications(context.Context, *GetNotificationsRequest) (*GetNotificationsResponse, error)
	MarkNotificationRead(context.Context, *MarkNotificationReadRequest) (*MarkNotificationReadResponse, error)
	GetUnreadCount(context.Context, *GetUnreadCountRequest) (*GetUnreadCountResponse, error)
	DeleteNotification(context.Context, *DeleteNotificationRequest) (*DeleteNotificationResponse, error)
	SendBulkEmail(context.Context, *SendBulkEmailRequest) (*SendEmailResponse, error)
	SendJobAlertEmail(context.Context, *JobAlertEmailRequest) (*SendEmailResponse, error)
	CreateNotification(context.Context, *CreateNotificationRequest) (*NotificationResponse, error)
	mustEmbedUnimplementedNotificationServiceServer()
}

type UnimplementedNotificationServiceServer struct{}

func (UnimplementedNotificationServiceServer) SendEmail(context.Context, *SendEmailRequest) (*SendEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedNotificationServiceServer) SendWelcomeEmail(context.Context, *WelcomeEmailRequest) (*SendEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedNotificationServiceServer) SendJobApplicationEmail(context.Context, *JobApplicationEmailRequest) (*SendEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedNotificationServiceServer) SendApplicationStatusEmail(context.Context, *ApplicationStatusEmailRequest) (*SendEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedNotificationServiceServer) SendPasswordResetEmail(context.Context, *PasswordResetEmailRequest) (*SendEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedNotificationServiceServer) GetNotifications(context.Context, *GetNotificationsRequest) (*GetNotificationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedNotificationServiceServer) MarkNotificationRead(context.Context, *MarkNotificationReadRequest) (*MarkNotificationReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedNotificationServiceServer) GetUnreadCount(context.Context, *GetUnreadCountRequest) (*GetUnreadCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedNotificationServiceServer) DeleteNotification(context.Context, *DeleteNotificationRequest) (*DeleteNotificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedNotificationServiceServer) SendBulkEmail(context.Context, *SendBulkEmailRequest) (*SendEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedNotificationServiceServer) SendJobAlertEmail(context.Context, *JobAlertEmailRequest) (*SendEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedNotificationServiceServer) CreateNotification(context.Context, *CreateNotificationRequest) (*NotificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {}

// ── Registration & Handlers ───────────────────────────────────────────────

func RegisterNotificationServiceServer(s grpc.ServiceRegistrar, srv NotificationServiceServer) {
	s.RegisterService(&NotificationService_ServiceDesc, srv)
}

var NotificationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notification.NotificationService",
	HandlerType: (*NotificationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "SendEmail", Handler: _Notif_SendEmail_Handler},
		{MethodName: "SendWelcomeEmail", Handler: _Notif_SendWelcomeEmail_Handler},
		{MethodName: "SendJobApplicationEmail", Handler: _Notif_SendJobApplicationEmail_Handler},
		{MethodName: "SendApplicationStatusEmail", Handler: _Notif_SendApplicationStatusEmail_Handler},
		{MethodName: "SendPasswordResetEmail", Handler: _Notif_SendPasswordResetEmail_Handler},
		{MethodName: "GetNotifications", Handler: _Notif_GetNotifications_Handler},
		{MethodName: "MarkNotificationRead", Handler: _Notif_MarkNotificationRead_Handler},
		{MethodName: "GetUnreadCount", Handler: _Notif_GetUnreadCount_Handler},
		{MethodName: "DeleteNotification", Handler: _Notif_DeleteNotification_Handler},
		{MethodName: "SendBulkEmail", Handler: _Notif_SendBulkEmail_Handler},
		{MethodName: "SendJobAlertEmail", Handler: _Notif_SendJobAlertEmail_Handler},
		{MethodName: "CreateNotification", Handler: _Notif_CreateNotification_Handler},
	},
	Streams: []grpc.StreamDesc{}, Metadata: "notification/notification.proto",
}

func _Notif_SendEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendEmailRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(NotificationServiceServer).SendEmail(ctx, in) }
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/notification.NotificationService/SendEmail"}, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(NotificationServiceServer).SendEmail(ctx, req.(*SendEmailRequest)) })
}
func _Notif_SendWelcomeEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WelcomeEmailRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(NotificationServiceServer).SendWelcomeEmail(ctx, in) }
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/notification.NotificationService/SendWelcomeEmail"}, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(NotificationServiceServer).SendWelcomeEmail(ctx, req.(*WelcomeEmailRequest)) })
}
func _Notif_SendJobApplicationEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobApplicationEmailRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(NotificationServiceServer).SendJobApplicationEmail(ctx, in) }
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/notification.NotificationService/SendJobApplicationEmail"}, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(NotificationServiceServer).SendJobApplicationEmail(ctx, req.(*JobApplicationEmailRequest)) })
}
func _Notif_SendApplicationStatusEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApplicationStatusEmailRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(NotificationServiceServer).SendApplicationStatusEmail(ctx, in) }
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/notification.NotificationService/SendApplicationStatusEmail"}, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(NotificationServiceServer).SendApplicationStatusEmail(ctx, req.(*ApplicationStatusEmailRequest)) })
}
func _Notif_SendPasswordResetEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PasswordResetEmailRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(NotificationServiceServer).SendPasswordResetEmail(ctx, in) }
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/notification.NotificationService/SendPasswordResetEmail"}, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(NotificationServiceServer).SendPasswordResetEmail(ctx, req.(*PasswordResetEmailRequest)) })
}
func _Notif_GetNotifications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNotificationsRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(NotificationServiceServer).GetNotifications(ctx, in) }
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/notification.NotificationService/GetNotifications"}, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(NotificationServiceServer).GetNotifications(ctx, req.(*GetNotificationsRequest)) })
}
func _Notif_MarkNotificationRead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MarkNotificationReadRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(NotificationServiceServer).MarkNotificationRead(ctx, in) }
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/notification.NotificationService/MarkNotificationRead"}, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(NotificationServiceServer).MarkNotificationRead(ctx, req.(*MarkNotificationReadRequest)) })
}
func _Notif_GetUnreadCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUnreadCountRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(NotificationServiceServer).GetUnreadCount(ctx, in) }
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/notification.NotificationService/GetUnreadCount"}, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(NotificationServiceServer).GetUnreadCount(ctx, req.(*GetUnreadCountRequest)) })
}
func _Notif_DeleteNotification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteNotificationRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(NotificationServiceServer).DeleteNotification(ctx, in) }
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/notification.NotificationService/DeleteNotification"}, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(NotificationServiceServer).DeleteNotification(ctx, req.(*DeleteNotificationRequest)) })
}
func _Notif_SendBulkEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendBulkEmailRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(NotificationServiceServer).SendBulkEmail(ctx, in) }
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/notification.NotificationService/SendBulkEmail"}, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(NotificationServiceServer).SendBulkEmail(ctx, req.(*SendBulkEmailRequest)) })
}
func _Notif_SendJobAlertEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobAlertEmailRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(NotificationServiceServer).SendJobAlertEmail(ctx, in) }
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/notification.NotificationService/SendJobAlertEmail"}, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(NotificationServiceServer).SendJobAlertEmail(ctx, req.(*JobAlertEmailRequest)) })
}
func _Notif_CreateNotification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateNotificationRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(NotificationServiceServer).CreateNotification(ctx, in) }
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/notification.NotificationService/CreateNotification"}, func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(NotificationServiceServer).CreateNotification(ctx, req.(*CreateNotificationRequest)) })
}

// ── Client ────────────────────────────────────────────────────────────────

type NotificationServiceClient interface {
	SendEmail(ctx context.Context, in *SendEmailRequest, opts ...grpc.CallOption) (*SendEmailResponse, error)
	SendWelcomeEmail(ctx context.Context, in *WelcomeEmailRequest, opts ...grpc.CallOption) (*SendEmailResponse, error)
	SendJobApplicationEmail(ctx context.Context, in *JobApplicationEmailRequest, opts ...grpc.CallOption) (*SendEmailResponse, error)
	SendApplicationStatusEmail(ctx context.Context, in *ApplicationStatusEmailRequest, opts ...grpc.CallOption) (*SendEmailResponse, error)
	SendPasswordResetEmail(ctx context.Context, in *PasswordResetEmailRequest, opts ...grpc.CallOption) (*SendEmailResponse, error)
	GetNotifications(ctx context.Context, in *GetNotificationsRequest, opts ...grpc.CallOption) (*GetNotificationsResponse, error)
	MarkNotificationRead(ctx context.Context, in *MarkNotificationReadRequest, opts ...grpc.CallOption) (*MarkNotificationReadResponse, error)
	GetUnreadCount(ctx context.Context, in *GetUnreadCountRequest, opts ...grpc.CallOption) (*GetUnreadCountResponse, error)
	DeleteNotification(ctx context.Context, in *DeleteNotificationRequest, opts ...grpc.CallOption) (*DeleteNotificationResponse, error)
	SendBulkEmail(ctx context.Context, in *SendBulkEmailRequest, opts ...grpc.CallOption) (*SendEmailResponse, error)
	SendJobAlertEmail(ctx context.Context, in *JobAlertEmailRequest, opts ...grpc.CallOption) (*SendEmailResponse, error)
	CreateNotification(ctx context.Context, in *CreateNotificationRequest, opts ...grpc.CallOption) (*NotificationResponse, error)
}

type notificationServiceClient struct{ cc grpc.ClientConnInterface }

func NewNotificationServiceClient(cc grpc.ClientConnInterface) NotificationServiceClient {
	return &notificationServiceClient{cc}
}

func (c *notificationServiceClient) SendEmail(ctx context.Context, in *SendEmailRequest, opts ...grpc.CallOption) (*SendEmailResponse, error) {
	out := new(SendEmailResponse); return out, c.cc.Invoke(ctx, "/notification.NotificationService/SendEmail", in, out, opts...)
}
func (c *notificationServiceClient) SendWelcomeEmail(ctx context.Context, in *WelcomeEmailRequest, opts ...grpc.CallOption) (*SendEmailResponse, error) {
	out := new(SendEmailResponse); return out, c.cc.Invoke(ctx, "/notification.NotificationService/SendWelcomeEmail", in, out, opts...)
}
func (c *notificationServiceClient) SendJobApplicationEmail(ctx context.Context, in *JobApplicationEmailRequest, opts ...grpc.CallOption) (*SendEmailResponse, error) {
	out := new(SendEmailResponse); return out, c.cc.Invoke(ctx, "/notification.NotificationService/SendJobApplicationEmail", in, out, opts...)
}
func (c *notificationServiceClient) SendApplicationStatusEmail(ctx context.Context, in *ApplicationStatusEmailRequest, opts ...grpc.CallOption) (*SendEmailResponse, error) {
	out := new(SendEmailResponse); return out, c.cc.Invoke(ctx, "/notification.NotificationService/SendApplicationStatusEmail", in, out, opts...)
}
func (c *notificationServiceClient) SendPasswordResetEmail(ctx context.Context, in *PasswordResetEmailRequest, opts ...grpc.CallOption) (*SendEmailResponse, error) {
	out := new(SendEmailResponse); return out, c.cc.Invoke(ctx, "/notification.NotificationService/SendPasswordResetEmail", in, out, opts...)
}
func (c *notificationServiceClient) GetNotifications(ctx context.Context, in *GetNotificationsRequest, opts ...grpc.CallOption) (*GetNotificationsResponse, error) {
	out := new(GetNotificationsResponse); return out, c.cc.Invoke(ctx, "/notification.NotificationService/GetNotifications", in, out, opts...)
}
func (c *notificationServiceClient) MarkNotificationRead(ctx context.Context, in *MarkNotificationReadRequest, opts ...grpc.CallOption) (*MarkNotificationReadResponse, error) {
	out := new(MarkNotificationReadResponse); return out, c.cc.Invoke(ctx, "/notification.NotificationService/MarkNotificationRead", in, out, opts...)
}
func (c *notificationServiceClient) GetUnreadCount(ctx context.Context, in *GetUnreadCountRequest, opts ...grpc.CallOption) (*GetUnreadCountResponse, error) {
	out := new(GetUnreadCountResponse); return out, c.cc.Invoke(ctx, "/notification.NotificationService/GetUnreadCount", in, out, opts...)
}
func (c *notificationServiceClient) DeleteNotification(ctx context.Context, in *DeleteNotificationRequest, opts ...grpc.CallOption) (*DeleteNotificationResponse, error) {
	out := new(DeleteNotificationResponse); return out, c.cc.Invoke(ctx, "/notification.NotificationService/DeleteNotification", in, out, opts...)
}
func (c *notificationServiceClient) SendBulkEmail(ctx context.Context, in *SendBulkEmailRequest, opts ...grpc.CallOption) (*SendEmailResponse, error) {
	out := new(SendEmailResponse); return out, c.cc.Invoke(ctx, "/notification.NotificationService/SendBulkEmail", in, out, opts...)
}
func (c *notificationServiceClient) SendJobAlertEmail(ctx context.Context, in *JobAlertEmailRequest, opts ...grpc.CallOption) (*SendEmailResponse, error) {
	out := new(SendEmailResponse); return out, c.cc.Invoke(ctx, "/notification.NotificationService/SendJobAlertEmail", in, out, opts...)
}
func (c *notificationServiceClient) CreateNotification(ctx context.Context, in *CreateNotificationRequest, opts ...grpc.CallOption) (*NotificationResponse, error) {
	out := new(NotificationResponse); return out, c.cc.Invoke(ctx, "/notification.NotificationService/CreateNotification", in, out, opts...)
}

// ── Proto interface methods ──────────────────────────────────────────────
func (m *SendEmailRequest) Reset()         { *m = SendEmailRequest{} }
func (m *SendEmailRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *SendEmailRequest) ProtoMessage()   {}

func (m *SendEmailResponse) Reset()         { *m = SendEmailResponse{} }
func (m *SendEmailResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *SendEmailResponse) ProtoMessage()   {}

func (m *WelcomeEmailRequest) Reset()         { *m = WelcomeEmailRequest{} }
func (m *WelcomeEmailRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *WelcomeEmailRequest) ProtoMessage()   {}

func (m *JobApplicationEmailRequest) Reset()         { *m = JobApplicationEmailRequest{} }
func (m *JobApplicationEmailRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *JobApplicationEmailRequest) ProtoMessage()   {}

func (m *ApplicationStatusEmailRequest) Reset()         { *m = ApplicationStatusEmailRequest{} }
func (m *ApplicationStatusEmailRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *ApplicationStatusEmailRequest) ProtoMessage()   {}

func (m *PasswordResetEmailRequest) Reset()         { *m = PasswordResetEmailRequest{} }
func (m *PasswordResetEmailRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *PasswordResetEmailRequest) ProtoMessage()   {}

func (m *GetNotificationsRequest) Reset()         { *m = GetNotificationsRequest{} }
func (m *GetNotificationsRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *GetNotificationsRequest) ProtoMessage()   {}

func (m *NotificationResponse) Reset()         { *m = NotificationResponse{} }
func (m *NotificationResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *NotificationResponse) ProtoMessage()   {}

func (m *GetNotificationsResponse) Reset()         { *m = GetNotificationsResponse{} }
func (m *GetNotificationsResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *GetNotificationsResponse) ProtoMessage()   {}

func (m *MarkNotificationReadRequest) Reset()         { *m = MarkNotificationReadRequest{} }
func (m *MarkNotificationReadRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *MarkNotificationReadRequest) ProtoMessage()   {}

func (m *MarkNotificationReadResponse) Reset()         { *m = MarkNotificationReadResponse{} }
func (m *MarkNotificationReadResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *MarkNotificationReadResponse) ProtoMessage()   {}

func (m *GetUnreadCountRequest) Reset()         { *m = GetUnreadCountRequest{} }
func (m *GetUnreadCountRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *GetUnreadCountRequest) ProtoMessage()   {}

func (m *GetUnreadCountResponse) Reset()         { *m = GetUnreadCountResponse{} }
func (m *GetUnreadCountResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *GetUnreadCountResponse) ProtoMessage()   {}

func (m *DeleteNotificationRequest) Reset()         { *m = DeleteNotificationRequest{} }
func (m *DeleteNotificationRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *DeleteNotificationRequest) ProtoMessage()   {}

func (m *DeleteNotificationResponse) Reset()         { *m = DeleteNotificationResponse{} }
func (m *DeleteNotificationResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *DeleteNotificationResponse) ProtoMessage()   {}

func (m *SendBulkEmailRequest) Reset()         { *m = SendBulkEmailRequest{} }
func (m *SendBulkEmailRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *SendBulkEmailRequest) ProtoMessage()   {}

func (m *JobAlertEmailRequest) Reset()         { *m = JobAlertEmailRequest{} }
func (m *JobAlertEmailRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *JobAlertEmailRequest) ProtoMessage()   {}

func (m *CreateNotificationRequest) Reset()         { *m = CreateNotificationRequest{} }
func (m *CreateNotificationRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *CreateNotificationRequest) ProtoMessage()   {}

