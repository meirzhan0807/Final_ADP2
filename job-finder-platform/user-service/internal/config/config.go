package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL    string
	RedisHost      string
	NATSUrl        string
	GRPCPort       string
	JWTSecret      string
	SMTPHost       string
	SMTPPort       string
	SMTPUser       string
	SMTPPassword   string
	AppEnv         string
	FrontendURL    string
	JaegerEndpoint string
}

func Load() *Config {
	return &Config{
		DatabaseURL: fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", "postgres123"),
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_NAME", "userdb"),
		),
		RedisHost:      getEnv("REDIS_HOST", "localhost:6379"),
		NATSUrl:        getEnv("NATS_URL", "nats://localhost:4222"),
		GRPCPort:       getEnv("GRPC_PORT", "50051"),
		JWTSecret:      getEnv("JWT_SECRET", "super-secret-jwt-key-change-in-production"),
		SMTPHost:       getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:       getEnv("SMTP_PORT", "587"),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPassword:   getEnv("SMTP_PASSWORD", ""),
		AppEnv:         getEnv("APP_ENV", "development"),
		FrontendURL:    getEnv("FRONTEND_URL", "http://localhost:3001"),
		JaegerEndpoint: getEnv("JAEGER_ENDPOINT", "jaeger:4317"),
	}
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}
