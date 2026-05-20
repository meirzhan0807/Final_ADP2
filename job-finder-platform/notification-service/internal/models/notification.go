package models

import "time"

type Notification struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Title     string    `db:"title"`
	Message   string    `db:"message"`
	Type      string    `db:"type"`
	IsRead    bool      `db:"is_read"`
	CreatedAt time.Time `db:"created_at"`
}
