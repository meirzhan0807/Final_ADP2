package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	jobpb "jobfinder/api-gateway/internal/pb/job"
	notifpb "jobfinder/api-gateway/internal/pb/notification"
)

type Gateway struct {
	userHTTP string // e.g. "http://user-service:9091"
	job      jobpb.JobServiceClient
	notif    notifpb.NotificationServiceClient
}

func NewGateway(userAddr, jobAddr, notifAddr string) (*Gateway, error) {
	// userAddr = "user-service:50051" → http base = "http://user-service:9091"
	// Extract host from "host:port" and use port 9091
	userHost := extractHost(userAddr)
	userHTTP := fmt.Sprintf("http://%s:9091", userHost)

	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	jobConn, err := grpc.Dial(jobAddr, dialOpts...)
	if err != nil {
		return nil, err
	}
	notifConn, err := grpc.Dial(notifAddr, dialOpts...)
	if err != nil {
		return nil, err
	}

	return &Gateway{
		userHTTP: userHTTP,
		job:      jobpb.NewJobServiceClient(jobConn),
		notif:    notifpb.NewNotificationServiceClient(notifConn),
	}, nil
}

func extractHost(addr string) string {
	for i, c := range addr {
		if c == ':' {
			return addr[:i]
		}
	}
	return addr
}

// ── HTTP helper to call user-service ──────────────────────────────────────

func (g *Gateway) userReq(method, path string, body interface{}, token string) (map[string]interface{}, int, error) {
	var bodyReader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, g.userHTTP+path, bodyReader)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, resp.StatusCode, nil
}

func rctx(c *gin.Context) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go func() { <-c.Request.Context().Done(); cancel() }()
	return ctx
}

func intQ(c *gin.Context, key string, def int) int {
	v, err := strconv.Atoi(c.DefaultQuery(key, strconv.Itoa(def)))
	if err != nil {
		return def
	}
	return v
}

func floatQ(c *gin.Context, key string, def float64) float64 {
	v, err := strconv.ParseFloat(c.DefaultQuery(key, "0"), 64)
	if err != nil {
		return def
	}
	return v
}

// ── Health ────────────────────────────────────────────────────────────────

func (g *Gateway) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "api-gateway", "time": time.Now()})
}

// ── Auth ──────────────────────────────────────────────────────────────────

func (g *Gateway) Register(c *gin.Context) {
	var req struct {
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=8"`
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Role      string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Role == "" {
		req.Role = "jobseeker"
	}
	result, status, err := g.userReq("POST", "/internal/register", req, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, result)
}

func (g *Gateway) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, status, err := g.userReq("POST", "/internal/login", req, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, result)
}

func (g *Gateway) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, status, err := g.userReq("POST", "/internal/refresh-token", req, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, result)
}

func (g *Gateway) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	result, status, err := g.userReq("GET", "/internal/verify-email?token="+token, nil, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, result)
}

// ── User ──────────────────────────────────────────────────────────────────

func (g *Gateway) GetMe(c *gin.Context) {
	userID := c.GetString("user_id")
	result, status, err := g.userReq("GET", "/internal/user/"+userID, nil, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, result)
}

func (g *Gateway) UpdateMe(c *gin.Context) {
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, status, err := g.userReq("PUT", "/internal/user-update/"+c.GetString("user_id"), req, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, result)
}

func (g *Gateway) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, status, err := g.userReq("PUT", "/internal/change-password/"+c.GetString("user_id"), req, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, result)
}

func (g *Gateway) GetProfile(c *gin.Context) {
	result, status, err := g.userReq("GET", "/internal/profile/"+c.GetString("user_id"), nil, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, result)
}

func (g *Gateway) UpdateProfile(c *gin.Context) {
	var req struct {
		Bio       string `json:"bio"`
		Skills    string `json:"skills"`
		ResumeUrl string `json:"resume_url"`
		Location  string `json:"location"`
		Phone     string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, status, err := g.userReq("PUT", "/internal/profile/"+c.GetString("user_id"), req, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, result)
}

func (g *Gateway) ListUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"users": []interface{}{}, "total": 0})
}

func (g *Gateway) DeleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ── Jobs (gRPC) ───────────────────────────────────────────────────────────

func (g *Gateway) ListJobs(c *gin.Context) {
	resp, err := g.job.ListJobs(rctx(c), &jobpb.ListJobsRequest{
		Page: int32(intQ(c, "page", 1)), Limit: int32(intQ(c, "limit", 20)),
		Status: c.DefaultQuery("status", "active"),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (g *Gateway) SearchJobs(c *gin.Context) {
	resp, err := g.job.SearchJobs(rctx(c), &jobpb.SearchJobsRequest{
		Query: c.Query("q"), Location: c.Query("location"), Category: c.Query("category"),
		JobType: c.Query("job_type"), ExperienceLevel: c.Query("experience_level"),
		SalaryMin: floatQ(c, "salary_min", 0), SalaryMax: floatQ(c, "salary_max", 0),
		Page: int32(intQ(c, "page", 1)), Limit: int32(intQ(c, "limit", 20)),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (g *Gateway) GetJob(c *gin.Context) {
	resp, err := g.job.GetJob(rctx(c), &jobpb.GetJobRequest{JobId: c.Param("id")})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (g *Gateway) GetCategories(c *gin.Context) {
	resp, err := g.job.GetJobCategories(rctx(c), &jobpb.GetJobCategoriesRequest{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (g *Gateway) CreateJob(c *gin.Context) {
	var req struct {
		Title, Description, Location, JobType, Category, Currency, ExperienceLevel string
		SalaryMin, SalaryMax                                                        float64
		Requirements                                                                []string
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := g.job.CreateJob(rctx(c), &jobpb.CreateJobRequest{
		EmployerId: c.GetString("user_id"), Title: req.Title, Description: req.Description,
		Location: req.Location, JobType: req.JobType, Category: req.Category,
		SalaryMin: req.SalaryMin, SalaryMax: req.SalaryMax, Currency: req.Currency,
		Requirements: req.Requirements, ExperienceLevel: req.ExperienceLevel,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (g *Gateway) UpdateJob(c *gin.Context) {
	var req struct {
		Title, Description, Location, JobType, Category, Status string
		SalaryMin, SalaryMax                                     float64
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := g.job.UpdateJob(rctx(c), &jobpb.UpdateJobRequest{
		JobId: c.Param("id"), Title: req.Title, Description: req.Description,
		Location: req.Location, JobType: req.JobType, Category: req.Category,
		SalaryMin: req.SalaryMin, SalaryMax: req.SalaryMax, Status: req.Status,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (g *Gateway) DeleteJob(c *gin.Context) {
	resp, err := g.job.DeleteJob(rctx(c), &jobpb.DeleteJobRequest{
		JobId: c.Param("id"), EmployerId: c.GetString("user_id"),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (g *Gateway) ApplyJob(c *gin.Context) {
	var req struct {
		CoverLetter string `json:"cover_letter"`
		ResumeUrl   string `json:"resume_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := g.job.ApplyJob(rctx(c), &jobpb.ApplyJobRequest{
		JobId: c.Param("id"), ApplicantId: c.GetString("user_id"),
		CoverLetter: req.CoverLetter, ResumeUrl: req.ResumeUrl,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (g *Gateway) GetMyJobs(c *gin.Context) {
	resp, err := g.job.GetJobsByEmployer(rctx(c), &jobpb.GetJobsByEmployerRequest{
		EmployerId: c.GetString("user_id"),
		Page:       int32(intQ(c, "page", 1)), Limit: int32(intQ(c, "limit", 20)),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (g *Gateway) GetMyApplications(c *gin.Context) {
	resp, err := g.job.GetApplicationsByUser(rctx(c), &jobpb.GetApplicationsByUserRequest{
		UserId: c.GetString("user_id"),
		Page:   int32(intQ(c, "page", 1)), Limit: int32(intQ(c, "limit", 20)),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (g *Gateway) GetJobApplications(c *gin.Context) {
	resp, err := g.job.GetApplications(rctx(c), &jobpb.GetApplicationsRequest{
		JobId: c.Param("id"),
		Page:  int32(intQ(c, "page", 1)), Limit: int32(intQ(c, "limit", 20)),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (g *Gateway) UpdateApplicationStatus(c *gin.Context) {
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := g.job.UpdateApplicationStatus(rctx(c), &jobpb.UpdateApplicationStatusRequest{
		ApplicationId: c.Param("id"), Status: req.Status, EmployerId: c.GetString("user_id"),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// ── Notifications ─────────────────────────────────────────────────────────

func (g *Gateway) GetNotifications(c *gin.Context) {
	resp, err := g.notif.GetNotifications(rctx(c), &notifpb.GetNotificationsRequest{
		UserId: c.GetString("user_id"),
		Page:   int32(intQ(c, "page", 1)), Limit: int32(intQ(c, "limit", 20)),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (g *Gateway) GetUnreadCount(c *gin.Context) {
	resp, err := g.notif.GetUnreadCount(rctx(c), &notifpb.GetUnreadCountRequest{UserId: c.GetString("user_id")})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (g *Gateway) MarkNotificationRead(c *gin.Context) {
	resp, err := g.notif.MarkNotificationRead(rctx(c), &notifpb.MarkNotificationReadRequest{
		NotificationId: c.Param("id"), UserId: c.GetString("user_id"),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (g *Gateway) DeleteNotification(c *gin.Context) {
	resp, err := g.notif.DeleteNotification(rctx(c), &notifpb.DeleteNotificationRequest{
		NotificationId: c.Param("id"), UserId: c.GetString("user_id"),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
