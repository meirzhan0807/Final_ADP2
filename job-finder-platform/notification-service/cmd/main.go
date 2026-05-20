package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"jobfinder/notification-service/internal/db"
	"jobfinder/notification-service/internal/email"
	"jobfinder/notification-service/internal/handlers"
	pb "jobfinder/notification-service/internal/pb"
	"jobfinder/notification-service/internal/repositories"
	"jobfinder/notification-service/internal/services"
)

func getEnv(k, d string) string {
	if v, ok := os.LookupEnv(k); ok { return v }
	return d
}

func main() {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("DB_USER", "postgres"), getEnv("DB_PASSWORD", "postgres123"),
		getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "notificationdb"),
	)

	pool, err := db.NewPostgresDB(dbURL)
	if err != nil { log.Fatalf("DB: %v", err) }
	defer pool.Close()

	if err := db.RunMigrations(pool); err != nil {
		log.Printf("Migration: %v", err)
	}

	nc, err := nats.Connect(getEnv("NATS_URL", "nats://localhost:4222"),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(20),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil { log.Fatalf("NATS: %v", err) }
	defer nc.Drain()

	emailSvc := email.NewService(
		getEnv("SMTP_HOST", "smtp.gmail.com"),
		getEnv("SMTP_PORT", "587"),
		getEnv("SMTP_USER", ""),
		getEnv("SMTP_PASSWORD", ""),
		getEnv("FROM_EMAIL", "noreply@jobfinder.com"),
	)

	repo := repositories.NewNotificationRepository(pool)
	svc := services.NewNotificationService(repo, emailSvc)
	handler := handlers.NewNotificationHandler(svc)

	svc.StartEventListeners(nc)

	grpcPort := getEnv("GRPC_PORT", "50053")
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil { log.Fatalf("listen: %v", err) }

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(loggingInterceptor))
	pb.RegisterNotificationServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"status":"healthy","service":"notification-service"}`))
		})
		http.ListenAndServe(":9093", mux)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Notification Service gRPC on :%s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("serve: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down notification-service...")
	grpcServer.GracefulStop()
}

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("[notification-service] %s | %v | err=%v", info.FullMethod, time.Since(start), err)
	return resp, err
}
