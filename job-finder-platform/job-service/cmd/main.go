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
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"jobfinder/job-service/internal/db"
	"jobfinder/job-service/internal/handlers"
	"jobfinder/job-service/internal/messaging"
	pb "jobfinder/job-service/internal/pb"
	"jobfinder/job-service/internal/repositories"
	"jobfinder/job-service/internal/services"
)

func getEnv(k, d string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return d
}

func main() {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("DB_USER", "postgres"), getEnv("DB_PASSWORD", "postgres123"),
		getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "jobdb"),
	)

	pool, err := db.NewPostgresDB(dbURL)
	if err != nil {
		log.Fatalf("DB: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(pool); err != nil {
		log.Printf("Migration warning: %v", err)
	}

	redisClient := db.NewRedisClient(getEnv("REDIS_HOST", "localhost:6379"))
	defer redisClient.Close()

	natsClient, err := messaging.NewNATSClient(getEnv("NATS_URL", "nats://localhost:4222"))
	if err != nil {
		log.Fatalf("NATS: %v", err)
	}
	defer natsClient.Close()

	jobRepo := repositories.NewJobRepository(pool)
	appRepo := repositories.NewApplicationRepository(pool)
	jobSvc := services.NewJobService(jobRepo, appRepo, natsClient)
	jobHandler := handlers.NewJobHandler(jobSvc)

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(loggingInterceptor))
	pb.RegisterJobServiceServer(grpcServer, jobHandler)
	reflection.Register(grpcServer)

	grpcPort := getEnv("GRPC_PORT", "50052")
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"status":"healthy","service":"job-service"}`))
		})
		http.ListenAndServe(":9092", mux)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Job Service gRPC on :%s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("serve: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down job-service...")
	grpcServer.GracefulStop()
}

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("[job-service] %s | %v | err=%v", info.FullMethod, time.Since(start), err)
	return resp, err
}
