package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"jobfinder/user-service/internal/config"
	"jobfinder/user-service/internal/db"
	"jobfinder/user-service/internal/handlers"
	"jobfinder/user-service/internal/messaging"
	"jobfinder/user-service/internal/models"
	pb "jobfinder/user-service/internal/pb"
	"jobfinder/user-service/internal/repositories"
	"jobfinder/user-service/internal/services"
)

func main() {
	cfg := config.Load()

	pool, err := db.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connect: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(pool); err != nil {
		log.Printf("Migration warning: %v", err)
	}

	redisClient := db.NewRedisClient(cfg.RedisHost)
	defer redisClient.Close()

	natsClient, err := messaging.NewNATSClient(cfg.NATSUrl)
	if err != nil {
		log.Fatalf("NATS connect: %v", err)
	}
	defer natsClient.Close()

	userRepo := repositories.NewUserRepository(pool)
	profileRepo := repositories.NewProfileRepository(pool)
	userSvc := services.NewUserService(userRepo, profileRepo, redisClient, natsClient, cfg)
	userHandler := handlers.NewUserHandler(userSvc)

	// gRPC server
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(loggingInterceptor))
	pb.RegisterUserServiceServer(grpcServer, userHandler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	// HTTP server — metrics + internal REST API
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"healthy","service":"user-service"}`))
		})

		// ── Internal REST endpoints (called by api-gateway) ──
		mux.HandleFunc("/internal/register", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			var req struct {
				Email     string `json:"email"`
				Password  string `json:"password"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Role      string `json:"role"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				jsonErr(w, "invalid json", http.StatusBadRequest)
				return
			}
			if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
				jsonErr(w, "all fields required", http.StatusBadRequest)
				return
			}
			if len(req.Password) < 8 {
				jsonErr(w, "password must be at least 8 characters", http.StatusBadRequest)
				return
			}
			if req.Role == "" {
				req.Role = "jobseeker"
			}
			user, err := userSvc.Register(r.Context(), &models.CreateUserInput{
				Email: req.Email, Password: req.Password,
				FirstName: req.FirstName, LastName: req.LastName, Role: req.Role,
			})
			if err != nil {
				jsonErr(w, err.Error(), http.StatusConflict)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"user_id": user.ID, "message": "Registration successful!", "success": true,
			})
		})

		mux.HandleFunc("/internal/login", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			var req struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				jsonErr(w, "invalid json", http.StatusBadRequest)
				return
			}
			user, tokens, err := userSvc.Login(r.Context(), req.Email, req.Password)
			if err != nil {
				jsonErr(w, "invalid credentials", http.StatusUnauthorized)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"access_token":  tokens.AccessToken,
				"refresh_token": tokens.RefreshToken,
				"user": map[string]interface{}{
					"id": user.ID, "email": user.Email,
					"first_name": user.FirstName, "last_name": user.LastName,
					"role": user.Role, "is_verified": user.IsVerified,
				},
			})
		})

		mux.HandleFunc("/internal/user/", func(w http.ResponseWriter, r *http.Request) {
			userID := strings.TrimPrefix(r.URL.Path, "/internal/user/")
			if userID == "" {
				jsonErr(w, "user_id required", http.StatusBadRequest)
				return
			}
			user, err := userSvc.GetUser(r.Context(), userID)
			if err != nil {
				jsonErr(w, "user not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id": user.ID, "email": user.Email,
				"first_name": user.FirstName, "last_name": user.LastName,
				"role": user.Role, "is_verified": user.IsVerified,
			})
		})

		mux.HandleFunc("/internal/user-update/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			userID := strings.TrimPrefix(r.URL.Path, "/internal/user-update/")
			var req struct {
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Email     string `json:"email"`
			}
			json.NewDecoder(r.Body).Decode(&req)
			user, err := userSvc.UpdateUser(r.Context(), userID, &models.UpdateUserInput{
				FirstName: req.FirstName, LastName: req.LastName, Email: req.Email,
			})
			if err != nil {
				jsonErr(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id": user.ID, "email": user.Email,
				"first_name": user.FirstName, "last_name": user.LastName, "role": user.Role,
			})
		})

		mux.HandleFunc("/internal/change-password/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			userID := strings.TrimPrefix(r.URL.Path, "/internal/change-password/")
			var req struct {
				OldPassword string `json:"old_password"`
				NewPassword string `json:"new_password"`
			}
			json.NewDecoder(r.Body).Decode(&req)
			if err := userSvc.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
				jsonErr(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Password changed"})
		})

		mux.HandleFunc("/internal/profile/", func(w http.ResponseWriter, r *http.Request) {
			userID := strings.TrimPrefix(r.URL.Path, "/internal/profile/")
			if r.Method == http.MethodGet {
				p, err := userSvc.GetProfile(r.Context(), userID)
				if err != nil {
					jsonErr(w, err.Error(), http.StatusNotFound)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"user_id": p.UserID, "bio": p.Bio, "skills": p.Skills,
					"resume_url": p.ResumeURL, "location": p.Location, "phone": p.Phone,
				})
			} else if r.Method == http.MethodPut {
				var req struct {
					Bio       string `json:"bio"`
					Skills    string `json:"skills"`
					ResumeUrl string `json:"resume_url"`
					Location  string `json:"location"`
					Phone     string `json:"phone"`
				}
				json.NewDecoder(r.Body).Decode(&req)
				p, err := userSvc.UpdateProfile(r.Context(), &models.UpdateProfileInput{
					UserID: userID, Bio: req.Bio, Skills: req.Skills,
					ResumeURL: req.ResumeUrl, Location: req.Location, Phone: req.Phone,
				})
				if err != nil {
					jsonErr(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"user_id": p.UserID, "bio": p.Bio, "skills": p.Skills,
					"location": p.Location, "phone": p.Phone,
				})
			}
		})

		mux.HandleFunc("/internal/validate-token", func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			token = strings.TrimPrefix(token, "Bearer ")
			userID, role, err := userSvc.ValidateToken(r.Context(), token)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{"valid": false})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"valid": true, "user_id": userID, "role": role})
		})

		mux.HandleFunc("/internal/refresh-token", func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				RefreshToken string `json:"refresh_token"`
			}
			json.NewDecoder(r.Body).Decode(&req)
			tokens, err := userSvc.RefreshToken(r.Context(), req.RefreshToken)
			if err != nil {
				jsonErr(w, "invalid refresh token", http.StatusUnauthorized)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"access_token": tokens.AccessToken, "refresh_token": tokens.RefreshToken,
			})
		})

		mux.HandleFunc("/internal/verify-email", func(w http.ResponseWriter, r *http.Request) {
			token := r.URL.Query().Get("token")
			if err := userSvc.VerifyEmail(r.Context(), token); err != nil {
				jsonErr(w, "invalid token", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
		})

		mux.HandleFunc("/internal/applications", func(w http.ResponseWriter, r *http.Request) {
			// placeholder — job-service handles this
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"applications": []interface{}{}})
		})

		log.Printf("User Service HTTP on :9091")
		http.ListenAndServe(":9091", mux)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("User Service gRPC on :%s", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("serve: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down user-service...")
	grpcServer.GracefulStop()
}

func jsonErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("[user-service] %s | %v | err=%v", info.FullMethod, time.Since(start), err)
	return resp, err
}
