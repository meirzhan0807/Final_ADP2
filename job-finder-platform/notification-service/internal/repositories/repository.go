package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"jobfinder/notification-service/internal/models"
)

type NotificationRepository interface {
	Create(ctx context.Context, userID, title, message, nType string) (*models.Notification, error)
	GetByUserID(ctx context.Context, userID string, page, limit int) ([]*models.Notification, int, error)
	MarkRead(ctx context.Context, id, userID string) error
	Delete(ctx context.Context, id, userID string) error
	GetUnreadCount(ctx context.Context, userID string) (int, error)
}

type notifRepo struct{ db *pgxpool.Pool }

func NewNotificationRepository(db *pgxpool.Pool) NotificationRepository { return &notifRepo{db: db} }

func (r *notifRepo) Create(ctx context.Context, userID, title, message, nType string) (*models.Notification, error) {
	n := &models.Notification{}
	err := r.db.QueryRow(ctx,
		`INSERT INTO notifications (id,user_id,title,message,type,is_read,created_at)
		 VALUES ($1,$2,$3,$4,$5,false,$6)
		 RETURNING id,user_id,title,message,type,is_read,created_at`,
		uuid.New().String(), userID, title, message, nType, time.Now(),
	).Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.Type, &n.IsRead, &n.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create notification: %w", err)
	}
	return n, nil
}

func (r *notifRepo) GetByUserID(ctx context.Context, userID string, page, limit int) ([]*models.Notification, int, error) {
	var total int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE user_id=$1`, userID).Scan(&total)

	rows, err := r.db.Query(ctx,
		`SELECT id,user_id,title,message,type,is_read,created_at FROM notifications
		 WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*models.Notification
	for rows.Next() {
		n := &models.Notification{}
		rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.Type, &n.IsRead, &n.CreatedAt)
		list = append(list, n)
	}
	return list, total, nil
}

func (r *notifRepo) MarkRead(ctx context.Context, id, userID string) error {
	res, err := r.db.Exec(ctx,
		`UPDATE notifications SET is_read=true WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("notification not found")
	}
	return nil
}

func (r *notifRepo) Delete(ctx context.Context, id, userID string) error {
	res, err := r.db.Exec(ctx,
		`DELETE FROM notifications WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("notification not found")
	}
	return nil
}

func (r *notifRepo) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM notifications WHERE user_id=$1 AND is_read=false`, userID,
	).Scan(&count)
	return count, err
}
