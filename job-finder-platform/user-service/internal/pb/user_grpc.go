// Code generated manually - gRPC service definitions for User Service
package userpb

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)


// ── Request/Response types ────────────────────────────────────────────────

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

type RegisterResponse struct {
	UserId  string `json:"user_id"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	User         *UserResponse `json:"user"`
}

type UserResponse struct {
	Id         string    `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type GetUserRequest struct {
	UserId string `json:"user_id"`
}

type UpdateUserRequest struct {
	UserId    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type DeleteUserRequest struct {
	UserId string `json:"user_id"`
}

type DeleteUserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ListUsersRequest struct {
	Page  int32 `json:"page"`
	Limit int32 `json:"limit"`
}

type ListUsersResponse struct {
	Users []*UserResponse `json:"users"`
	Total int32           `json:"total"`
}

type ChangePasswordRequest struct {
	UserId      string `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangePasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

type VerifyEmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type GetUserProfileRequest struct {
	UserId string `json:"user_id"`
}

type UserProfileResponse struct {
	UserId    string `json:"user_id"`
	Bio       string `json:"bio"`
	Skills    string `json:"skills"`
	ResumeUrl string `json:"resume_url"`
	Location  string `json:"location"`
	Phone     string `json:"phone"`
}

type UpdateUserProfileRequest struct {
	UserId    string `json:"user_id"`
	Bio       string `json:"bio"`
	Skills    string `json:"skills"`
	ResumeUrl string `json:"resume_url"`
	Location  string `json:"location"`
	Phone     string `json:"phone"`
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	Valid  bool   `json:"valid"`
	UserId string `json:"user_id"`
	Role   string `json:"role"`
}

// ── Service interface ─────────────────────────────────────────────────────

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

func (UnimplementedUserServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedUserServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedUserServiceServer) GetUser(context.Context, *GetUserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedUserServiceServer) UpdateUser(context.Context, *UpdateUserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedUserServiceServer) DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedUserServiceServer) ListUsers(context.Context, *ListUsersRequest) (*ListUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUsers not implemented")
}
func (UnimplementedUserServiceServer) ChangePassword(context.Context, *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePassword not implemented")
}
func (UnimplementedUserServiceServer) VerifyEmail(context.Context, *VerifyEmailRequest) (*VerifyEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyEmail not implemented")
}
func (UnimplementedUserServiceServer) RefreshToken(context.Context, *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshToken not implemented")
}
func (UnimplementedUserServiceServer) GetUserProfile(context.Context, *GetUserProfileRequest) (*UserProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserProfile not implemented")
}
func (UnimplementedUserServiceServer) UpdateUserProfile(context.Context, *UpdateUserProfileRequest) (*UserProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserProfile not implemented")
}
func (UnimplementedUserServiceServer) ValidateToken(context.Context, *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateToken not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// ── gRPC registration ─────────────────────────────────────────────────────

const UserService_Register_FullMethodName         = "/user.UserService/Register"
const UserService_Login_FullMethodName            = "/user.UserService/Login"
const UserService_GetUser_FullMethodName          = "/user.UserService/GetUser"
const UserService_UpdateUser_FullMethodName       = "/user.UserService/UpdateUser"
const UserService_DeleteUser_FullMethodName       = "/user.UserService/DeleteUser"
const UserService_ListUsers_FullMethodName        = "/user.UserService/ListUsers"
const UserService_ChangePassword_FullMethodName   = "/user.UserService/ChangePassword"
const UserService_VerifyEmail_FullMethodName      = "/user.UserService/VerifyEmail"
const UserService_RefreshToken_FullMethodName     = "/user.UserService/RefreshToken"
const UserService_GetUserProfile_FullMethodName   = "/user.UserService/GetUserProfile"
const UserService_UpdateUserProfile_FullMethodName = "/user.UserService/UpdateUserProfile"
const UserService_ValidateToken_FullMethodName    = "/user.UserService/ValidateToken"

var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Register", Handler: _UserService_Register_Handler},
		{MethodName: "Login", Handler: _UserService_Login_Handler},
		{MethodName: "GetUser", Handler: _UserService_GetUser_Handler},
		{MethodName: "UpdateUser", Handler: _UserService_UpdateUser_Handler},
		{MethodName: "DeleteUser", Handler: _UserService_DeleteUser_Handler},
		{MethodName: "ListUsers", Handler: _UserService_ListUsers_Handler},
		{MethodName: "ChangePassword", Handler: _UserService_ChangePassword_Handler},
		{MethodName: "VerifyEmail", Handler: _UserService_VerifyEmail_Handler},
		{MethodName: "RefreshToken", Handler: _UserService_RefreshToken_Handler},
		{MethodName: "GetUserProfile", Handler: _UserService_GetUserProfile_Handler},
		{MethodName: "UpdateUserProfile", Handler: _UserService_UpdateUserProfile_Handler},
		{MethodName: "ValidateToken", Handler: _UserService_ValidateToken_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user/user.proto",
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_Register_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_Login_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_GetUser_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_UpdateUser_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_DeleteUser_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ListUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ListUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_ListUsers_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ListUsers(ctx, req.(*ListUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ChangePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ChangePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_ChangePassword_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ChangePassword(ctx, req.(*ChangePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_VerifyEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).VerifyEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_VerifyEmail_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).VerifyEmail(ctx, req.(*VerifyEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_RefreshToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RefreshToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_RefreshToken_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).RefreshToken(ctx, req.(*RefreshTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUserProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_GetUserProfile_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserProfile(ctx, req.(*GetUserProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateUserProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateUserProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_UpdateUserProfile_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateUserProfile(ctx, req.(*UpdateUserProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ValidateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ValidateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: UserService_ValidateToken_FullMethodName}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ValidateToken(ctx, req.(*ValidateTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ── gRPC Client ───────────────────────────────────────────────────────────

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

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, UserService_Register_FullMethodName, in, out, opts...)
	return out, err
}
func (c *userServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, UserService_Login_FullMethodName, in, out, opts...)
	return out, err
}
func (c *userServiceClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, UserService_GetUser_FullMethodName, in, out, opts...)
	return out, err
}
func (c *userServiceClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, UserService_UpdateUser_FullMethodName, in, out, opts...)
	return out, err
}
func (c *userServiceClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error) {
	out := new(DeleteUserResponse)
	err := c.cc.Invoke(ctx, UserService_DeleteUser_FullMethodName, in, out, opts...)
	return out, err
}
func (c *userServiceClient) ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (*ListUsersResponse, error) {
	out := new(ListUsersResponse)
	err := c.cc.Invoke(ctx, UserService_ListUsers_FullMethodName, in, out, opts...)
	return out, err
}
func (c *userServiceClient) ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*ChangePasswordResponse, error) {
	out := new(ChangePasswordResponse)
	err := c.cc.Invoke(ctx, UserService_ChangePassword_FullMethodName, in, out, opts...)
	return out, err
}
func (c *userServiceClient) VerifyEmail(ctx context.Context, in *VerifyEmailRequest, opts ...grpc.CallOption) (*VerifyEmailResponse, error) {
	out := new(VerifyEmailResponse)
	err := c.cc.Invoke(ctx, UserService_VerifyEmail_FullMethodName, in, out, opts...)
	return out, err
}
func (c *userServiceClient) RefreshToken(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*RefreshTokenResponse, error) {
	out := new(RefreshTokenResponse)
	err := c.cc.Invoke(ctx, UserService_RefreshToken_FullMethodName, in, out, opts...)
	return out, err
}
func (c *userServiceClient) GetUserProfile(ctx context.Context, in *GetUserProfileRequest, opts ...grpc.CallOption) (*UserProfileResponse, error) {
	out := new(UserProfileResponse)
	err := c.cc.Invoke(ctx, UserService_GetUserProfile_FullMethodName, in, out, opts...)
	return out, err
}
func (c *userServiceClient) UpdateUserProfile(ctx context.Context, in *UpdateUserProfileRequest, opts ...grpc.CallOption) (*UserProfileResponse, error) {
	out := new(UserProfileResponse)
	err := c.cc.Invoke(ctx, UserService_UpdateUserProfile_FullMethodName, in, out, opts...)
	return out, err
}
func (c *userServiceClient) ValidateToken(ctx context.Context, in *ValidateTokenRequest, opts ...grpc.CallOption) (*ValidateTokenResponse, error) {
	out := new(ValidateTokenResponse)
	err := c.cc.Invoke(ctx, UserService_ValidateToken_FullMethodName, in, out, opts...)
	return out, err
}

// ── Proto interface methods (required for gRPC encoding) ────────────────
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

