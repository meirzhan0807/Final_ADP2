package config

import "os"

type Config struct {
	HTTPPort                string
	UserServiceAddr         string
	JobServiceAddr          string
	NotificationServiceAddr string
	JWTSecret               string
	RedisHost               string
	AppEnv                  string
}

func Load() *Config {
	return &Config{
		HTTPPort:                getEnv("HTTP_PORT", "8080"),
		UserServiceAddr:         getEnv("USER_SERVICE_ADDR", "localhost:50051"),
		JobServiceAddr:          getEnv("JOB_SERVICE_ADDR", "localhost:50052"),
		NotificationServiceAddr: getEnv("NOTIFICATION_SERVICE_ADDR", "localhost:50053"),
		JWTSecret:               getEnv("JWT_SECRET", "super-secret-jwt-key-change-in-production"),
		RedisHost:               getEnv("REDIS_HOST", "localhost:6379"),
		AppEnv:                  getEnv("APP_ENV", "development"),
	}
}

func getEnv(k, d string) string {
	if v, ok := os.LookupEnv(k); ok { return v }
	return d
}
