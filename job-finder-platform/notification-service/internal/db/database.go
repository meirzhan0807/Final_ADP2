package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	cfg.MaxConns = 20
	cfg.MinConns = 3
	cfg.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}

func RunMigrations(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), migrationSQL)
	return err
}

const migrationSQL = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS notifications (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL,
    title      VARCHAR(255) NOT NULL,
    message    TEXT NOT NULL,
    type       VARCHAR(50) DEFAULT 'info',
    is_read    BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notif_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notif_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notif_created ON notifications(created_at DESC);
`
