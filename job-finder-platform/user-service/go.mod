module jobfinder/user-service

go 1.21

require (
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.5.4
	github.com/nats-io/nats.go v1.33.1
	github.com/redis/go-redis/v9 v9.4.0
	golang.org/x/crypto v0.19.0
	google.golang.org/grpc v1.62.0
	google.golang.org/protobuf v1.33.0
	github.com/prometheus/client_golang v1.19.0
	go.uber.org/zap v1.27.0
	github.com/stretchr/testify v1.9.0
	github.com/stretchr/objx v0.5.2
	github.com/golang-migrate/migrate/v4 v4.17.0
	go.opentelemetry.io/otel v1.24.0
	go.opentelemetry.io/otel/sdk v1.24.0
	go.opentelemetry.io/otel/trace v1.24.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.24.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240221002015-b0ce06bbee7c
)
