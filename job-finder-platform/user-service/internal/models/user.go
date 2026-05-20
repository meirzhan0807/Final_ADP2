package models

import "time"

type User struct {
	ID          string    `db:"id"`
	Email       string    `db:"email"`
	Password    string    `db:"password_hash"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	Role        string    `db:"role"`
	IsVerified  bool      `db:"is_verified"`
	VerifyToken string    `db:"verify_token"`
	ResetToken  string    `db:"reset_token"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type UserProfile struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Bio       string    `db:"bio"`
	Skills    string    `db:"skills"`
	ResumeURL string    `db:"resume_url"`
	Location  string    `db:"location"`
	Phone     string    `db:"phone"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CreateUserInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Role      string
}

type UpdateUserInput struct {
	FirstName string
	LastName  string
	Email     string
}

type UpdateProfileInput struct {
	UserID    string
	Bio       string
	Skills    string
	ResumeURL string
	Location  string
	Phone     string
}

type TokenClaims struct {
	UserID string
	Email  string
	Role   string
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}
