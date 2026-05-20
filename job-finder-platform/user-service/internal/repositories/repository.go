package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"jobfinder/user-service/internal/models"
)

// ── User Repository ───────────────────────────────────────────────────────

type UserRepository interface {
	Create(ctx context.Context, input *models.CreateUserInput) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, id string, input *models.UpdateUserInput) (*models.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*models.User, int, error)
	UpdatePassword(ctx context.Context, id, hash string) error
	SetVerifyToken(ctx context.Context, id, token string) error
	VerifyEmail(ctx context.Context, token string) error
	SetResetToken(ctx context.Context, email, token string) error
}

type userRepo struct{ db *pgxpool.Pool }

func NewUserRepository(db *pgxpool.Pool) UserRepository { return &userRepo{db: db} }

func (r *userRepo) Create(ctx context.Context, in *models.CreateUserInput) (*models.User, error) {
	id := uuid.New().String()
	now := time.Now()
	u := &models.User{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO users (id,email,password_hash,first_name,last_name,role,is_verified,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,false,$7,$8)
		RETURNING id,email,password_hash,first_name,last_name,role,is_verified,created_at,updated_at`,
		id, in.Email, in.Password, in.FirstName, in.LastName, in.Role, now, now,
	).Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Role, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id,email,password_hash,first_name,last_name,role,is_verified,created_at,updated_at FROM users WHERE id=$1`, id,
	).Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Role, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return u, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id,email,password_hash,first_name,last_name,role,is_verified,created_at,updated_at FROM users WHERE email=$1`, email,
	).Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Role, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return u, nil
}

func (r *userRepo) Update(ctx context.Context, id string, in *models.UpdateUserInput) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(ctx, `
		UPDATE users SET first_name=$1,last_name=$2,email=$3,updated_at=$4 WHERE id=$5
		RETURNING id,email,password_hash,first_name,last_name,role,is_verified,created_at,updated_at`,
		in.FirstName, in.LastName, in.Email, time.Now(), id,
	).Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Role, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	return err
}

func (r *userRepo) List(ctx context.Context, page, limit int) ([]*models.User, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	rows, err := r.db.Query(ctx,
		`SELECT id,email,password_hash,first_name,last_name,role,is_verified,created_at,updated_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Role, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *userRepo) UpdatePassword(ctx context.Context, id, hash string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET password_hash=$1,updated_at=$2 WHERE id=$3`, hash, time.Now(), id)
	return err
}

func (r *userRepo) SetVerifyToken(ctx context.Context, id, token string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET verify_token=$1,updated_at=$2 WHERE id=$3`, token, time.Now(), id)
	return err
}

func (r *userRepo) VerifyEmail(ctx context.Context, token string) error {
	res, err := r.db.Exec(ctx, `UPDATE users SET is_verified=true,verify_token='',updated_at=$1 WHERE verify_token=$2`, time.Now(), token)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("invalid token")
	}
	return nil
}

func (r *userRepo) SetResetToken(ctx context.Context, email, token string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET reset_token=$1,updated_at=$2 WHERE email=$3`, token, time.Now(), email)
	return err
}

// ── Profile Repository ────────────────────────────────────────────────────

type ProfileRepository interface {
	GetByUserID(ctx context.Context, userID string) (*models.UserProfile, error)
	Upsert(ctx context.Context, in *models.UpdateProfileInput) (*models.UserProfile, error)
}

type profileRepo struct{ db *pgxpool.Pool }

func NewProfileRepository(db *pgxpool.Pool) ProfileRepository { return &profileRepo{db: db} }

func (r *profileRepo) GetByUserID(ctx context.Context, userID string) (*models.UserProfile, error) {
	p := &models.UserProfile{UserID: userID}
	err := r.db.QueryRow(ctx,
		`SELECT id,user_id,bio,skills,resume_url,location,phone,created_at,updated_at FROM user_profiles WHERE user_id=$1`, userID,
	).Scan(&p.ID, &p.UserID, &p.Bio, &p.Skills, &p.ResumeURL, &p.Location, &p.Phone, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return &models.UserProfile{UserID: userID}, nil
	}
	return p, nil
}

func (r *profileRepo) Upsert(ctx context.Context, in *models.UpdateProfileInput) (*models.UserProfile, error) {
	id := uuid.New().String()
	now := time.Now()
	p := &models.UserProfile{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO user_profiles (id,user_id,bio,skills,resume_url,location,phone,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (user_id) DO UPDATE
		SET bio=$3,skills=$4,resume_url=$5,location=$6,phone=$7,updated_at=$9
		RETURNING id,user_id,bio,skills,resume_url,location,phone,created_at,updated_at`,
		id, in.UserID, in.Bio, in.Skills, in.ResumeURL, in.Location, in.Phone, now, now,
	).Scan(&p.ID, &p.UserID, &p.Bio, &p.Skills, &p.ResumeURL, &p.Location, &p.Phone, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}
