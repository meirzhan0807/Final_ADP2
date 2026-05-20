package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"jobfinder/api-gateway/internal/config"
	"jobfinder/api-gateway/internal/handlers"
	"jobfinder/api-gateway/internal/middleware"
)

func main() {
	cfg := config.Load()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	gw, err := handlers.NewGateway(cfg.UserServiceAddr, cfg.JobServiceAddr, cfg.NotificationServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RateLimit())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
	}))

	// ── Public ────────────────────────────────────────────────────────────
	r.GET("/health", gw.Health)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := r.Group("/api/v1")

	// Auth
	auth := v1.Group("/auth")
	auth.POST("/register", gw.Register)
	auth.POST("/login", gw.Login)
	auth.POST("/refresh", gw.RefreshToken)
	auth.GET("/verify-email", gw.VerifyEmail)

	// Public jobs
	v1.GET("/jobs", gw.ListJobs)
	v1.GET("/jobs/search", gw.SearchJobs)
	v1.GET("/jobs/categories", gw.GetCategories)
	v1.GET("/jobs/:id", gw.GetJob)

	// ── Protected ─────────────────────────────────────────────────────────
	authMW := middleware.Auth(cfg.JWTSecret)
	p := v1.Group("/")
	p.Use(authMW)

	// Me
	p.GET("me", gw.GetMe)
	p.PUT("me", gw.UpdateMe)
	p.PUT("me/password", gw.ChangePassword)
	p.GET("me/profile", gw.GetProfile)
	p.PUT("me/profile", gw.UpdateProfile)
	p.GET("me/applications", gw.GetMyApplications)

	// Apply
	p.POST("jobs/:id/apply", gw.ApplyJob)

	// Employer
	emp := p.Group("employer")
	emp.Use(middleware.Role("employer", "admin"))
	emp.POST("jobs", gw.CreateJob)
	emp.PUT("jobs/:id", gw.UpdateJob)
	emp.DELETE("jobs/:id", gw.DeleteJob)
	emp.GET("jobs", gw.GetMyJobs)
	emp.GET("jobs/:id/applications", gw.GetJobApplications)
	emp.PUT("applications/:id/status", gw.UpdateApplicationStatus)

	// Notifications
	n := p.Group("notifications")
	n.GET("", gw.GetNotifications)
	n.GET("unread-count", gw.GetUnreadCount)
	n.PUT(":id/read", gw.MarkNotificationRead)
	n.DELETE(":id", gw.DeleteNotification)

	// Admin
	adm := p.Group("admin")
	adm.Use(middleware.Role("admin"))
	adm.GET("users", gw.ListUsers)
	adm.DELETE("users/:id", gw.DeleteUser)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
	})

	log.Printf("API Gateway starting on :%s", cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("server: %v", err)
	}
}
