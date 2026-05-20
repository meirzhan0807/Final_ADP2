package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "jobfinder/job-service/internal/pb"
	"jobfinder/job-service/internal/models"
	"jobfinder/job-service/internal/services"
)

type JobHandler struct {
	pb.UnimplementedJobServiceServer
	svc services.JobService
}

func NewJobHandler(svc services.JobService) *JobHandler { return &JobHandler{svc: svc} }

func (h *JobHandler) mustEmbedUnimplementedJobServiceServer() {}

func (h *JobHandler) CreateJob(ctx context.Context, req *pb.CreateJobRequest) (*pb.JobResponse, error) {
	if req.EmployerId == "" || req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "employer_id and title required")
	}
	job, err := h.svc.CreateJob(ctx, &models.CreateJobInput{
		EmployerID: req.EmployerId, Title: req.Title, Description: req.Description,
		Location: req.Location, JobType: req.JobType, Category: req.Category,
		SalaryMin: req.SalaryMin, SalaryMax: req.SalaryMax, Currency: req.Currency,
		Requirements: req.Requirements, ExperienceLevel: req.ExperienceLevel,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return jobToProto(job), nil
}

func (h *JobHandler) GetJob(ctx context.Context, req *pb.GetJobRequest) (*pb.JobResponse, error) {
	job, err := h.svc.GetJob(ctx, req.JobId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "job not found")
	}
	return jobToProto(job), nil
}

func (h *JobHandler) UpdateJob(ctx context.Context, req *pb.UpdateJobRequest) (*pb.JobResponse, error) {
	job, err := h.svc.UpdateJob(ctx, req.JobId, &models.UpdateJobInput{
		Title: req.Title, Description: req.Description, Location: req.Location,
		JobType: req.JobType, Category: req.Category, SalaryMin: req.SalaryMin,
		SalaryMax: req.SalaryMax, Status: req.Status,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return jobToProto(job), nil
}

func (h *JobHandler) DeleteJob(ctx context.Context, req *pb.DeleteJobRequest) (*pb.DeleteJobResponse, error) {
	if err := h.svc.DeleteJob(ctx, req.JobId, req.EmployerId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteJobResponse{Success: true, Message: "Job deleted"}, nil
}

func (h *JobHandler) ListJobs(ctx context.Context, req *pb.ListJobsRequest) (*pb.ListJobsResponse, error) {
	jobs, total, err := h.svc.ListJobs(ctx, int(req.Page), int(req.Limit), req.Status)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var list []*pb.JobResponse
	for _, j := range jobs {
		list = append(list, jobToProto(j))
	}
	return &pb.ListJobsResponse{Jobs: list, Total: int32(total), Page: req.Page, Limit: req.Limit}, nil
}

func (h *JobHandler) SearchJobs(ctx context.Context, req *pb.SearchJobsRequest) (*pb.ListJobsResponse, error) {
	jobs, total, err := h.svc.SearchJobs(ctx, &models.SearchInput{
		Query: req.Query, Location: req.Location, Category: req.Category,
		JobType: req.JobType, ExperienceLevel: req.ExperienceLevel,
		SalaryMin: req.SalaryMin, SalaryMax: req.SalaryMax,
		Page: int(req.Page), Limit: int(req.Limit),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var list []*pb.JobResponse
	for _, j := range jobs {
		list = append(list, jobToProto(j))
	}
	return &pb.ListJobsResponse{Jobs: list, Total: int32(total)}, nil
}

func (h *JobHandler) ApplyJob(ctx context.Context, req *pb.ApplyJobRequest) (*pb.ApplyJobResponse, error) {
	if req.JobId == "" || req.ApplicantId == "" {
		return nil, status.Error(codes.InvalidArgument, "job_id and applicant_id required")
	}
	app, err := h.svc.ApplyJob(ctx, req.JobId, req.ApplicantId, req.CoverLetter, req.ResumeUrl)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.ApplyJobResponse{ApplicationId: app.ID, Success: true, Message: "Applied successfully", Status: app.Status}, nil
}

func (h *JobHandler) GetApplications(ctx context.Context, req *pb.GetApplicationsRequest) (*pb.GetApplicationsResponse, error) {
	apps, total, err := h.svc.GetApplications(ctx, req.JobId, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var list []*pb.ApplicationResponse
	for _, a := range apps {
		list = append(list, appToProto(a))
	}
	return &pb.GetApplicationsResponse{Applications: list, Total: int32(total)}, nil
}

func (h *JobHandler) UpdateApplicationStatus(ctx context.Context, req *pb.UpdateApplicationStatusRequest) (*pb.ApplyJobResponse, error) {
	app, err := h.svc.UpdateApplicationStatus(ctx, req.ApplicationId, req.Status)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ApplyJobResponse{ApplicationId: app.ID, Success: true, Message: "Status updated", Status: app.Status}, nil
}

func (h *JobHandler) GetJobsByEmployer(ctx context.Context, req *pb.GetJobsByEmployerRequest) (*pb.ListJobsResponse, error) {
	jobs, total, err := h.svc.GetJobsByEmployer(ctx, req.EmployerId, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var list []*pb.JobResponse
	for _, j := range jobs {
		list = append(list, jobToProto(j))
	}
	return &pb.ListJobsResponse{Jobs: list, Total: int32(total)}, nil
}

func (h *JobHandler) GetApplicationsByUser(ctx context.Context, req *pb.GetApplicationsByUserRequest) (*pb.GetApplicationsResponse, error) {
	apps, total, err := h.svc.GetApplicationsByUser(ctx, req.UserId, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var list []*pb.ApplicationResponse
	for _, a := range apps {
		list = append(list, appToProto(a))
	}
	return &pb.GetApplicationsResponse{Applications: list, Total: int32(total)}, nil
}

func (h *JobHandler) GetJobCategories(ctx context.Context, req *pb.GetJobCategoriesRequest) (*pb.GetJobCategoriesResponse, error) {
	return &pb.GetJobCategoriesResponse{Categories: h.svc.GetCategories()}, nil
}

func jobToProto(j *models.Job) *pb.JobResponse {
	return &pb.JobResponse{
		Id: j.ID, EmployerId: j.EmployerID, Title: j.Title, Description: j.Description,
		Location: j.Location, JobType: j.JobType, Category: j.Category,
		SalaryMin: j.SalaryMin, SalaryMax: j.SalaryMax, Currency: j.Currency,
		Requirements: j.Requirements, ExperienceLevel: j.ExperienceLevel,
		Status: j.Status, ApplicationsCount: int32(j.ApplicationCount),
		CreatedAt: j.CreatedAt, UpdatedAt: j.UpdatedAt,
	}
}

func appToProto(a *models.Application) *pb.ApplicationResponse {
	return &pb.ApplicationResponse{
		Id: a.ID, JobId: a.JobID, ApplicantId: a.ApplicantID,
		CoverLetter: a.CoverLetter, ResumeUrl: a.ResumeURL,
		Status: a.Status, AppliedAt: a.AppliedAt,
	}
}
