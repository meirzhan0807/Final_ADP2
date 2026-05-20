package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"jobfinder/user-service/internal/models"
	pb "jobfinder/user-service/internal/pb"
	"jobfinder/user-service/internal/services"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	svc services.UserService
}

func NewUserHandler(svc services.UserService) *UserHandler { return &UserHandler{svc: svc} }

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return nil, status.Error(codes.InvalidArgument, "all fields required")
	}
	if len(req.Password) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}
	user, err := h.svc.Register(ctx, &models.CreateUserInput{
		Email: req.Email, Password: req.Password,
		FirstName: req.FirstName, LastName: req.LastName, Role: req.Role,
	})
	if err != nil {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	return &pb.RegisterResponse{UserId: user.ID, Message: "Registration successful! Check your email.", Success: true}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}
	user, tokens, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	return &pb.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         toProto(user),
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id required")
	}
	user, err := h.svc.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return toProto(user), nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	user, err := h.svc.UpdateUser(ctx, req.UserId, &models.UpdateUserInput{
		FirstName: req.FirstName, LastName: req.LastName, Email: req.Email,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toProto(user), nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := h.svc.DeleteUser(ctx, req.UserId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteUserResponse{Success: true, Message: "User deleted"}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, total, err := h.svc.ListUsers(ctx, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var list []*pb.UserResponse
	for _, u := range users {
		list = append(list, toProto(u))
	}
	return &pb.ListUsersResponse{Users: list, Total: int32(total)}, nil
}

func (h *UserHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	if err := h.svc.ChangePassword(ctx, req.UserId, req.OldPassword, req.NewPassword); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.ChangePasswordResponse{Success: true, Message: "Password changed"}, nil
}

func (h *UserHandler) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	if err := h.svc.VerifyEmail(ctx, req.Token); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid or expired token")
	}
	return &pb.VerifyEmailResponse{Success: true, Message: "Email verified"}, nil
}

func (h *UserHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	tokens, err := h.svc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &pb.RefreshTokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (h *UserHandler) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.UserProfileResponse, error) {
	p, err := h.svc.GetProfile(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UserProfileResponse{UserId: p.UserID, Bio: p.Bio, Skills: p.Skills, ResumeUrl: p.ResumeURL, Location: p.Location, Phone: p.Phone}, nil
}

func (h *UserHandler) UpdateUserProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*pb.UserProfileResponse, error) {
	p, err := h.svc.UpdateProfile(ctx, &models.UpdateProfileInput{
		UserID: req.UserId, Bio: req.Bio, Skills: req.Skills,
		ResumeURL: req.ResumeUrl, Location: req.Location, Phone: req.Phone,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UserProfileResponse{UserId: p.UserID, Bio: p.Bio, Skills: p.Skills, ResumeUrl: p.ResumeURL, Location: p.Location, Phone: p.Phone}, nil
}

func (h *UserHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	userID, role, err := h.svc.ValidateToken(ctx, req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}
	return &pb.ValidateTokenResponse{Valid: true, UserId: userID, Role: role}, nil
}

func (h *UserHandler) mustEmbedUnimplementedUserServiceServer() {}

func toProto(u *models.User) *pb.UserResponse {
	return &pb.UserResponse{
		Id: u.ID, Email: u.Email, FirstName: u.FirstName, LastName: u.LastName,
		Role: u.Role, IsVerified: u.IsVerified, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
	}
}
