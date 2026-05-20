package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func NewPostgresDB(url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	cfg.MaxConns = 25
	cfg.MinConns = 5
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

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: addr, PoolSize: 10})
}

func RunMigrations(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), migrationSQL)
	return err
}

const migrationSQL = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS jobs (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employer_id      UUID NOT NULL,
    title            VARCHAR(255) NOT NULL,
    description      TEXT NOT NULL,
    location         VARCHAR(255) DEFAULT '',
    job_type         VARCHAR(50) DEFAULT 'full-time',
    category         VARCHAR(100) DEFAULT '',
    salary_min       DECIMAL(12,2) DEFAULT 0,
    salary_max       DECIMAL(12,2) DEFAULT 0,
    currency         VARCHAR(10) DEFAULT 'USD',
    requirements     TEXT DEFAULT '',
    experience_level VARCHAR(50) DEFAULT 'mid',
    status           VARCHAR(50) DEFAULT 'active',
    application_count INT DEFAULT 0,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    updated_at       TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS job_applications (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id       UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    applicant_id UUID NOT NULL,
    cover_letter TEXT DEFAULT '',
    resume_url   VARCHAR(500) DEFAULT '',
    status       VARCHAR(50) DEFAULT 'pending',
    applied_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(job_id, applicant_id)
);

CREATE INDEX IF NOT EXISTS idx_jobs_employer   ON jobs(employer_id);
CREATE INDEX IF NOT EXISTS idx_jobs_status     ON jobs(status);
CREATE INDEX IF NOT EXISTS idx_jobs_category   ON jobs(category);
CREATE INDEX IF NOT EXISTS idx_apps_job        ON job_applications(job_id);
CREATE INDEX IF NOT EXISTS idx_apps_applicant  ON job_applications(applicant_id);
`
