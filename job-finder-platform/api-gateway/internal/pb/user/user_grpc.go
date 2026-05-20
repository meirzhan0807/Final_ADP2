package userpb

import (
	"context"
	"time"

	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type RegisterRequest  struct{ Email, Password, FirstName, LastName, Role string }
type RegisterResponse struct{ UserId, Message string; Success bool }
type LoginRequest     struct{ Email, Password string }
type LoginResponse    struct{ AccessToken, RefreshToken string; User *UserResponse }
type UserResponse     struct {
	Id, Email, FirstName, LastName, Role string
	IsVerified bool
	CreatedAt, UpdatedAt time.Time
}
type GetUserRequest            struct{ UserId string }
type UpdateUserRequest         struct{ UserId, FirstName, LastName, Email string }
type DeleteUserRequest         struct{ UserId string }
type DeleteUserResponse        struct{ Success bool; Message string }
type ListUsersRequest          struct{ Page, Limit int32 }
type ListUsersResponse         struct{ Users []*UserResponse; Total int32 }
type ChangePasswordRequest     struct{ UserId, OldPassword, NewPassword string }
type ChangePasswordResponse    struct{ Success bool; Message string }
type VerifyEmailRequest        struct{ Token string }
type VerifyEmailResponse       struct{ Success bool; Message string }
type RefreshTokenRequest       struct{ RefreshToken string }
type RefreshTokenResponse      struct{ AccessToken, RefreshToken string }
type GetUserProfileRequest     struct{ UserId string }
type UserProfileResponse       struct{ UserId, Bio, Skills, ResumeUrl, Location, Phone string }
type UpdateUserProfileRequest  struct{ UserId, Bio, Skills, ResumeUrl, Location, Phone string }
type ValidateTokenRequest      struct{ Token string }
type ValidateTokenResponse     struct{ Valid bool; UserId, Role string }

type UserServiceServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	GetUser(context.Context, *GetUserRequest) (*UserResponse, error)
	UpdateUser(context.Context, *UpdateUserRequest) (*UserResponse, error)
	DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error)
	ListUsers(context.Context, *ListUsersRequest) (*ListUsersResponse, error)
	ChangePassword(context.Context, *ChangePasswordRequest) (*ChangePasswordResponse, error)
	VerifyEmail(context.Context, *VerifyEmailRequest) (*VerifyEmailResponse, error)
	RefreshToken(context.Context, *RefreshTokenRequest) (*RefreshTokenResponse, error)
	GetUserProfile(context.Context, *GetUserProfileRequest) (*UserProfileResponse, error)
	UpdateUserProfile(context.Context, *UpdateUserProfileRequest) (*UserProfileResponse, error)
	ValidateToken(context.Context, *ValidateTokenRequest) (*ValidateTokenResponse, error)
	mustEmbedUnimplementedUserServiceServer()
}

type UnimplementedUserServiceServer struct{}
func (UnimplementedUserServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) { return nil, status.Errorf(codes.Unimplemented, "not implemented") }
func (UnimplementedUserServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) { return nil, status.Errorf(codes.Unimplemented, "not implemented") }
func (UnimplementedUserServiceServer) GetUser(context.Context, *GetUserRequest) (*UserResponse, error) { return nil, status.Errorf(codes.Unimplemented, "not implemented") }
func (UnimplementedUserServiceServer) UpdateUser(context.Context, *UpdateUserRequest) (*UserResponse, error) { return nil, status.Errorf(codes.Unimplemented, "not implemented") }
func (UnimplementedUserServiceServer) DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error) { return nil, status.Errorf(codes.Unimplemented, "not implemented") }
func (UnimplementedUserServiceServer) ListUsers(context.Context, *ListUsersRequest) (*ListUsersResponse, error) { return nil, status.Errorf(codes.Unimplemented, "not implemented") }
func (UnimplementedUserServiceServer) ChangePassword(context.Context, *ChangePasswordRequest) (*ChangePasswordResponse, error) { return nil, status.Errorf(codes.Unimplemented, "not implemented") }
func (UnimplementedUserServiceServer) VerifyEmail(context.Context, *VerifyEmailRequest) (*VerifyEmailResponse, error) { return nil, status.Errorf(codes.Unimplemented, "not implemented") }
func (UnimplementedUserServiceServer) RefreshToken(context.Context, *RefreshTokenRequest) (*RefreshTokenResponse, error) { return nil, status.Errorf(codes.Unimplemented, "not implemented") }
func (UnimplementedUserServiceServer) GetUserProfile(context.Context, *GetUserProfileRequest) (*UserProfileResponse, error) { return nil, status.Errorf(codes.Unimplemented, "not implemented") }
func (UnimplementedUserServiceServer) UpdateUserProfile(context.Context, *UpdateUserProfileRequest) (*UserProfileResponse, error) { return nil, status.Errorf(codes.Unimplemented, "not implemented") }
func (UnimplementedUserServiceServer) ValidateToken(context.Context, *ValidateTokenRequest) (*ValidateTokenResponse, error) { return nil, status.Errorf(codes.Unimplemented, "not implemented") }
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

type UserServiceClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*UserResponse, error)
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UserResponse, error)
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error)
	ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (*ListUsersResponse, error)
	ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*ChangePasswordResponse, error)
	VerifyEmail(ctx context.Context, in *VerifyEmailRequest, opts ...grpc.CallOption) (*VerifyEmailResponse, error)
	RefreshToken(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*RefreshTokenResponse, error)
	GetUserProfile(ctx context.Context, in *GetUserProfileRequest, opts ...grpc.CallOption) (*UserProfileResponse, error)
	UpdateUserProfile(ctx context.Context, in *UpdateUserProfileRequest, opts ...grpc.CallOption) (*UserProfileResponse, error)
	ValidateToken(ctx context.Context, in *ValidateTokenRequest, opts ...grpc.CallOption) (*ValidateTokenResponse, error)
}

type userServiceClient struct{ cc grpc.ClientConnInterface }
func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient { return &userServiceClient{cc} }
func (c *userServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) { out := new(RegisterResponse); return out, c.cc.Invoke(ctx, "/user.UserService/Register", in, out, opts...) }
func (c *userServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) { out := new(LoginResponse); return out, c.cc.Invoke(ctx, "/user.UserService/Login", in, out, opts...) }
func (c *userServiceClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*UserResponse, error) { out := new(UserResponse); return out, c.cc.Invoke(ctx, "/user.UserService/GetUser", in, out, opts...) }
func (c *userServiceClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UserResponse, error) { out := new(UserResponse); return out, c.cc.Invoke(ctx, "/user.UserService/UpdateUser", in, out, opts...) }
func (c *userServiceClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error) { out := new(DeleteUserResponse); return out, c.cc.Invoke(ctx, "/user.UserService/DeleteUser", in, out, opts...) }
func (c *userServiceClient) ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (*ListUsersResponse, error) { out := new(ListUsersResponse); return out, c.cc.Invoke(ctx, "/user.UserService/ListUsers", in, out, opts...) }
func (c *userServiceClient) ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*ChangePasswordResponse, error) { out := new(ChangePasswordResponse); return out, c.cc.Invoke(ctx, "/user.UserService/ChangePassword", in, out, opts...) }
func (c *userServiceClient) VerifyEmail(ctx context.Context, in *VerifyEmailRequest, opts ...grpc.CallOption) (*VerifyEmailResponse, error) { out := new(VerifyEmailResponse); return out, c.cc.Invoke(ctx, "/user.UserService/VerifyEmail", in, out, opts...) }
func (c *userServiceClient) RefreshToken(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*RefreshTokenResponse, error) { out := new(RefreshTokenResponse); return out, c.cc.Invoke(ctx, "/user.UserService/RefreshToken", in, out, opts...) }
func (c *userServiceClient) GetUserProfile(ctx context.Context, in *GetUserProfileRequest, opts ...grpc.CallOption) (*UserProfileResponse, error) { out := new(UserProfileResponse); return out, c.cc.Invoke(ctx, "/user.UserService/GetUserProfile", in, out, opts...) }
func (c *userServiceClient) UpdateUserProfile(ctx context.Context, in *UpdateUserProfileRequest, opts ...grpc.CallOption) (*UserProfileResponse, error) { out := new(UserProfileResponse); return out, c.cc.Invoke(ctx, "/user.UserService/UpdateUserProfile", in, out, opts...) }
func (c *userServiceClient) ValidateToken(ctx context.Context, in *ValidateTokenRequest, opts ...grpc.CallOption) (*ValidateTokenResponse, error) { out := new(ValidateTokenResponse); return out, c.cc.Invoke(ctx, "/user.UserService/ValidateToken", in, out, opts...) }

// ── Proto interface methods ──────────────────────────────────────────────
func (m *RegisterRequest) Reset()         { *m = RegisterRequest{} }
func (m *RegisterRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *RegisterRequest) ProtoMessage()   {}

func (m *RegisterResponse) Reset()         { *m = RegisterResponse{} }
func (m *RegisterResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *RegisterResponse) ProtoMessage()   {}

func (m *LoginRequest) Reset()         { *m = LoginRequest{} }
func (m *LoginRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *LoginRequest) ProtoMessage()   {}

func (m *LoginResponse) Reset()         { *m = LoginResponse{} }
func (m *LoginResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *LoginResponse) ProtoMessage()   {}

func (m *UserResponse) Reset()         { *m = UserResponse{} }
func (m *UserResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *UserResponse) ProtoMessage()   {}

func (m *GetUserRequest) Reset()         { *m = GetUserRequest{} }
func (m *GetUserRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *GetUserRequest) ProtoMessage()   {}

func (m *UpdateUserRequest) Reset()         { *m = UpdateUserRequest{} }
func (m *UpdateUserRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *UpdateUserRequest) ProtoMessage()   {}

func (m *DeleteUserRequest) Reset()         { *m = DeleteUserRequest{} }
func (m *DeleteUserRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *DeleteUserRequest) ProtoMessage()   {}

func (m *DeleteUserResponse) Reset()         { *m = DeleteUserResponse{} }
func (m *DeleteUserResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *DeleteUserResponse) ProtoMessage()   {}

func (m *ListUsersRequest) Reset()         { *m = ListUsersRequest{} }
func (m *ListUsersRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *ListUsersRequest) ProtoMessage()   {}

func (m *ListUsersResponse) Reset()         { *m = ListUsersResponse{} }
func (m *ListUsersResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *ListUsersResponse) ProtoMessage()   {}

func (m *ChangePasswordRequest) Reset()         { *m = ChangePasswordRequest{} }
func (m *ChangePasswordRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *ChangePasswordRequest) ProtoMessage()   {}

func (m *ChangePasswordResponse) Reset()         { *m = ChangePasswordResponse{} }
func (m *ChangePasswordResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *ChangePasswordResponse) ProtoMessage()   {}

func (m *VerifyEmailRequest) Reset()         { *m = VerifyEmailRequest{} }
func (m *VerifyEmailRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *VerifyEmailRequest) ProtoMessage()   {}

func (m *VerifyEmailResponse) Reset()         { *m = VerifyEmailResponse{} }
func (m *VerifyEmailResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *VerifyEmailResponse) ProtoMessage()   {}

func (m *RefreshTokenRequest) Reset()         { *m = RefreshTokenRequest{} }
func (m *RefreshTokenRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *RefreshTokenRequest) ProtoMessage()   {}

func (m *RefreshTokenResponse) Reset()         { *m = RefreshTokenResponse{} }
func (m *RefreshTokenResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *RefreshTokenResponse) ProtoMessage()   {}

func (m *GetUserProfileRequest) Reset()         { *m = GetUserProfileRequest{} }
func (m *GetUserProfileRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *GetUserProfileRequest) ProtoMessage()   {}

func (m *UserProfileResponse) Reset()         { *m = UserProfileResponse{} }
func (m *UserProfileResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *UserProfileResponse) ProtoMessage()   {}

func (m *UpdateUserProfileRequest) Reset()         { *m = UpdateUserProfileRequest{} }
func (m *UpdateUserProfileRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *UpdateUserProfileRequest) ProtoMessage()   {}

func (m *ValidateTokenRequest) Reset()         { *m = ValidateTokenRequest{} }
func (m *ValidateTokenRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *ValidateTokenRequest) ProtoMessage()   {}

func (m *ValidateTokenResponse) Reset()         { *m = ValidateTokenResponse{} }
func (m *ValidateTokenResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *ValidateTokenResponse) ProtoMessage()   {}

