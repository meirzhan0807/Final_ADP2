package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"jobfinder/user-service/internal/config"
	"jobfinder/user-service/internal/messaging"
	"jobfinder/user-service/internal/models"
	"jobfinder/user-service/internal/repositories"
)

type UserService interface {
	Register(ctx context.Context, in *models.CreateUserInput) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.User, *models.AuthTokens, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, in *models.UpdateUserInput) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, page, limit int) ([]*models.User, int, error)
	ChangePassword(ctx context.Context, id, old, newPass string) error
	VerifyEmail(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, refreshToken string) (*models.AuthTokens, error)
	GetProfile(ctx context.Context, userID string) (*models.UserProfile, error)
	UpdateProfile(ctx context.Context, in *models.UpdateProfileInput) (*models.UserProfile, error)
	ValidateToken(ctx context.Context, token string) (string, string, error)
	SendPasswordResetEmail(ctx context.Context, email string) error
}

type userService struct {
	userRepo    repositories.UserRepository
	profileRepo repositories.ProfileRepository
	redis       *redis.Client
	nats        *messaging.NATSClient
	cfg         *config.Config
}

func NewUserService(ur repositories.UserRepository, pr repositories.ProfileRepository, rc *redis.Client, nc *messaging.NATSClient, cfg *config.Config) UserService {
	return &userService{userRepo: ur, profileRepo: pr, redis: rc, nats: nc, cfg: cfg}
}

func (s *userService) Register(ctx context.Context, in *models.CreateUserInput) (*models.User, error) {
	if existing, _ := s.userRepo.GetByEmail(ctx, in.Email); existing != nil {
		return nil, fmt.Errorf("email already registered")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	in.Password = string(hashed)
	if in.Role == "" {
		in.Role = "jobseeker"
	}
	user, err := s.userRepo.Create(ctx, in)
	if err != nil {
		return nil, err
	}
	token, _ := randToken(32)
	s.userRepo.SetVerifyToken(ctx, user.ID, token)
	verifyLink := fmt.Sprintf("%s/verify-email?token=%s", s.cfg.FrontendURL, token)
	s.nats.Publish("user.registered", map[string]string{
		"user_id": user.ID, "email": user.Email,
		"name": user.FirstName + " " + user.LastName, "verify_link": verifyLink,
	})
	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (*models.User, *models.AuthTokens, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}
	tokens, err := s.genTokens(user)
	if err != nil {
		return nil, nil, err
	}
	s.redis.Set(ctx, "refresh:"+user.ID, tokens.RefreshToken, 7*24*time.Hour)
	return user, tokens, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, id string, in *models.UpdateUserInput) (*models.User, error) {
	user, err := s.userRepo.Update(ctx, id, in)
	if err != nil {
		return nil, err
	}
	s.redis.Del(ctx, "user:"+id)
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return err
	}
	s.redis.Del(ctx, "user:"+id, "refresh:"+id)
	return nil
}

func (s *userService) ListUsers(ctx context.Context, page, limit int) ([]*models.User, int, error) {
	if page <= 0 { page = 1 }
	if limit <= 0 || limit > 100 { limit = 20 }
	return s.userRepo.List(ctx, page, limit)
}

func (s *userService) ChangePassword(ctx context.Context, id, old, newPass string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(old)); err != nil {
		return fmt.Errorf("incorrect old password")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.userRepo.UpdatePassword(ctx, id, string(hashed))
}

func (s *userService) VerifyEmail(ctx context.Context, token string) error {
	return s.userRepo.VerifyEmail(ctx, token)
}

func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthTokens, error) {
	claims, err := s.parseToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}
	stored, err := s.redis.Get(ctx, "refresh:"+claims.UserID).Result()
	if err != nil || stored != refreshToken {
		return nil, fmt.Errorf("token expired")
	}
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	tokens, err := s.genTokens(user)
	if err != nil {
		return nil, err
	}
	s.redis.Set(ctx, "refresh:"+user.ID, tokens.RefreshToken, 7*24*time.Hour)
	return tokens, nil
}

func (s *userService) GetProfile(ctx context.Context, userID string) (*models.UserProfile, error) {
	return s.profileRepo.GetByUserID(ctx, userID)
}

func (s *userService) UpdateProfile(ctx context.Context, in *models.UpdateProfileInput) (*models.UserProfile, error) {
	return s.profileRepo.Upsert(ctx, in)
}

func (s *userService) ValidateToken(ctx context.Context, token string) (string, string, error) {
	claims, err := s.parseToken(token)
	if err != nil {
		return "", "", fmt.Errorf("invalid token")
	}
	return claims.UserID, claims.Role, nil
}

func (s *userService) SendPasswordResetEmail(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil // don't reveal
	}
	token, _ := randToken(32)
	s.userRepo.SetResetToken(ctx, email, token)
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.cfg.FrontendURL, token)
	s.nats.Publish("user.password_reset", map[string]string{
		"email": email, "name": user.FirstName, "reset_link": resetLink,
	})
	return nil
}

func (s *userService) genTokens(user *models.User) (*models.AuthTokens, error) {
	access, err := s.sign(user, 15*time.Minute)
	if err != nil {
		return nil, err
	}
	refresh, err := s.sign(user, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}
	return &models.AuthTokens{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *userService) sign(user *models.User, dur time.Duration) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID, "email": user.Email, "role": user.Role,
		"exp": time.Now().Add(dur).Unix(), "iat": time.Now().Unix(),
	}).SignedString([]byte(s.cfg.JWTSecret))
}

func (s *userService) parseToken(tokenStr string) (*models.TokenClaims, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !t.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	cl, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return &models.TokenClaims{
		UserID: cl["user_id"].(string),
		Email:  cl["email"].(string),
		Role:   cl["role"].(string),
	}, nil
}

func randToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
