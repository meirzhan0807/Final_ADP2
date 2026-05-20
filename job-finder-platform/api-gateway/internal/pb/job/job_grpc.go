package jobpb

import (
	"context"
	"time"

	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


// ── Types ─────────────────────────────────────────────────────────────────

type CreateJobRequest struct {
	EmployerId      string   `json:"employer_id"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Location        string   `json:"location"`
	JobType         string   `json:"job_type"`
	Category        string   `json:"category"`
	SalaryMin       float64  `json:"salary_min"`
	SalaryMax       float64  `json:"salary_max"`
	Currency        string   `json:"currency"`
	Requirements    []string `json:"requirements"`
	ExperienceLevel string   `json:"experience_level"`
}

type JobResponse struct {
	Id                string    `json:"id"`
	EmployerId        string    `json:"employer_id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Location          string    `json:"location"`
	JobType           string    `json:"job_type"`
	Category          string    `json:"category"`
	SalaryMin         float64   `json:"salary_min"`
	SalaryMax         float64   `json:"salary_max"`
	Currency          string    `json:"currency"`
	Requirements      []string  `json:"requirements"`
	ExperienceLevel   string    `json:"experience_level"`
	Status            string    `json:"status"`
	ApplicationsCount int32     `json:"applications_count"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type GetJobRequest struct{ JobId string `json:"job_id"` }

type UpdateJobRequest struct {
	JobId       string  `json:"job_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	JobType     string  `json:"job_type"`
	Category    string  `json:"category"`
	SalaryMin   float64 `json:"salary_min"`
	SalaryMax   float64 `json:"salary_max"`
	Status      string  `json:"status"`
}

type DeleteJobRequest struct {
	JobId      string `json:"job_id"`
	EmployerId string `json:"employer_id"`
}

type DeleteJobResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ListJobsRequest struct {
	Page   int32  `json:"page"`
	Limit  int32  `json:"limit"`
	Status string `json:"status"`
}

type ListJobsResponse struct {
	Jobs  []*JobResponse `json:"jobs"`
	Total int32          `json:"total"`
	Page  int32          `json:"page"`
	Limit int32          `json:"limit"`
}

type SearchJobsRequest struct {
	Query           string  `json:"query"`
	Location        string  `json:"location"`
	Category        string  `json:"category"`
	JobType         string  `json:"job_type"`
	ExperienceLevel string  `json:"experience_level"`
	SalaryMin       float64 `json:"salary_min"`
	SalaryMax       float64 `json:"salary_max"`
	Page            int32   `json:"page"`
	Limit           int32   `json:"limit"`
}

type ApplyJobRequest struct {
	JobId       string `json:"job_id"`
	ApplicantId string `json:"applicant_id"`
	CoverLetter string `json:"cover_letter"`
	ResumeUrl   string `json:"resume_url"`
}

type ApplyJobResponse struct {
	ApplicationId string `json:"application_id"`
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	Status        string `json:"status"`
}

type GetApplicationsRequest struct {
	JobId string `json:"job_id"`
	Page  int32  `json:"page"`
	Limit int32  `json:"limit"`
}

type ApplicationResponse struct {
	Id          string    `json:"id"`
	JobId       string    `json:"job_id"`
	ApplicantId string    `json:"applicant_id"`
	CoverLetter string    `json:"cover_letter"`
	ResumeUrl   string    `json:"resume_url"`
	Status      string    `json:"status"`
	AppliedAt   time.Time `json:"applied_at"`
}

type GetApplicationsResponse struct {
	Applications []*ApplicationResponse `json:"applications"`
	Total        int32                  `json:"total"`
}

type UpdateApplicationStatusRequest struct {
	ApplicationId string `json:"application_id"`
	Status        string `json:"status"`
	EmployerId    string `json:"employer_id"`
}

type GetJobsByEmployerRequest struct {
	EmployerId string `json:"employer_id"`
	Page       int32  `json:"page"`
	Limit      int32  `json:"limit"`
}

type GetApplicationsByUserRequest struct {
	UserId string `json:"user_id"`
	Page   int32  `json:"page"`
	Limit  int32  `json:"limit"`
}

type GetJobCategoriesRequest  struct{}
type GetJobCategoriesResponse struct {
	Categories []string `json:"categories"`
}

// ── Server Interface ──────────────────────────────────────────────────────

type JobServiceServer interface {
	CreateJob(context.Context, *CreateJobRequest) (*JobResponse, error)
	GetJob(context.Context, *GetJobRequest) (*JobResponse, error)
	UpdateJob(context.Context, *UpdateJobRequest) (*JobResponse, error)
	DeleteJob(context.Context, *DeleteJobRequest) (*DeleteJobResponse, error)
	ListJobs(context.Context, *ListJobsRequest) (*ListJobsResponse, error)
	SearchJobs(context.Context, *SearchJobsRequest) (*ListJobsResponse, error)
	ApplyJob(context.Context, *ApplyJobRequest) (*ApplyJobResponse, error)
	GetApplications(context.Context, *GetApplicationsRequest) (*GetApplicationsResponse, error)
	UpdateApplicationStatus(context.Context, *UpdateApplicationStatusRequest) (*ApplyJobResponse, error)
	GetJobsByEmployer(context.Context, *GetJobsByEmployerRequest) (*ListJobsResponse, error)
	GetApplicationsByUser(context.Context, *GetApplicationsByUserRequest) (*GetApplicationsResponse, error)
	GetJobCategories(context.Context, *GetJobCategoriesRequest) (*GetJobCategoriesResponse, error)
	mustEmbedUnimplementedJobServiceServer()
}

type UnimplementedJobServiceServer struct{}

func (UnimplementedJobServiceServer) CreateJob(context.Context, *CreateJobRequest) (*JobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedJobServiceServer) GetJob(context.Context, *GetJobRequest) (*JobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedJobServiceServer) UpdateJob(context.Context, *UpdateJobRequest) (*JobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedJobServiceServer) DeleteJob(context.Context, *DeleteJobRequest) (*DeleteJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedJobServiceServer) ListJobs(context.Context, *ListJobsRequest) (*ListJobsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedJobServiceServer) SearchJobs(context.Context, *SearchJobsRequest) (*ListJobsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedJobServiceServer) ApplyJob(context.Context, *ApplyJobRequest) (*ApplyJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedJobServiceServer) GetApplications(context.Context, *GetApplicationsRequest) (*GetApplicationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedJobServiceServer) UpdateApplicationStatus(context.Context, *UpdateApplicationStatusRequest) (*ApplyJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedJobServiceServer) GetJobsByEmployer(context.Context, *GetJobsByEmployerRequest) (*ListJobsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedJobServiceServer) GetApplicationsByUser(context.Context, *GetApplicationsByUserRequest) (*GetApplicationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedJobServiceServer) GetJobCategories(context.Context, *GetJobCategoriesRequest) (*GetJobCategoriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}
func (UnimplementedJobServiceServer) mustEmbedUnimplementedJobServiceServer() {}

// ── Registration ──────────────────────────────────────────────────────────

func RegisterJobServiceServer(s grpc.ServiceRegistrar, srv JobServiceServer) {
	s.RegisterService(&JobService_ServiceDesc, srv)
}

var JobService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "job.JobService",
	HandlerType: (*JobServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateJob", Handler: _JobService_CreateJob_Handler},
		{MethodName: "GetJob", Handler: _JobService_GetJob_Handler},
		{MethodName: "UpdateJob", Handler: _JobService_UpdateJob_Handler},
		{MethodName: "DeleteJob", Handler: _JobService_DeleteJob_Handler},
		{MethodName: "ListJobs", Handler: _JobService_ListJobs_Handler},
		{MethodName: "SearchJobs", Handler: _JobService_SearchJobs_Handler},
		{MethodName: "ApplyJob", Handler: _JobService_ApplyJob_Handler},
		{MethodName: "GetApplications", Handler: _JobService_GetApplications_Handler},
		{MethodName: "UpdateApplicationStatus", Handler: _JobService_UpdateApplicationStatus_Handler},
		{MethodName: "GetJobsByEmployer", Handler: _JobService_GetJobsByEmployer_Handler},
		{MethodName: "GetApplicationsByUser", Handler: _JobService_GetApplicationsByUser_Handler},
		{MethodName: "GetJobCategories", Handler: _JobService_GetJobCategories_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "job/job.proto",
}

func _JobService_CreateJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateJobRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(JobServiceServer).CreateJob(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/job.JobService/CreateJob"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).CreateJob(ctx, req.(*CreateJobRequest))
	})
}
func _JobService_GetJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetJobRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(JobServiceServer).GetJob(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/job.JobService/GetJob"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).GetJob(ctx, req.(*GetJobRequest))
	})
}
func _JobService_UpdateJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateJobRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(JobServiceServer).UpdateJob(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/job.JobService/UpdateJob"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).UpdateJob(ctx, req.(*UpdateJobRequest))
	})
}
func _JobService_DeleteJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteJobRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(JobServiceServer).DeleteJob(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/job.JobService/DeleteJob"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).DeleteJob(ctx, req.(*DeleteJobRequest))
	})
}
func _JobService_ListJobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListJobsRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(JobServiceServer).ListJobs(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/job.JobService/ListJobs"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).ListJobs(ctx, req.(*ListJobsRequest))
	})
}
func _JobService_SearchJobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchJobsRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(JobServiceServer).SearchJobs(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/job.JobService/SearchJobs"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).SearchJobs(ctx, req.(*SearchJobsRequest))
	})
}
func _JobService_ApplyJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApplyJobRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(JobServiceServer).ApplyJob(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/job.JobService/ApplyJob"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).ApplyJob(ctx, req.(*ApplyJobRequest))
	})
}
func _JobService_GetApplications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetApplicationsRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(JobServiceServer).GetApplications(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/job.JobService/GetApplications"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).GetApplications(ctx, req.(*GetApplicationsRequest))
	})
}
func _JobService_UpdateApplicationStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateApplicationStatusRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(JobServiceServer).UpdateApplicationStatus(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/job.JobService/UpdateApplicationStatus"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).UpdateApplicationStatus(ctx, req.(*UpdateApplicationStatusRequest))
	})
}
func _JobService_GetJobsByEmployer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetJobsByEmployerRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(JobServiceServer).GetJobsByEmployer(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/job.JobService/GetJobsByEmployer"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).GetJobsByEmployer(ctx, req.(*GetJobsByEmployerRequest))
	})
}
func _JobService_GetApplicationsByUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetApplicationsByUserRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(JobServiceServer).GetApplicationsByUser(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/job.JobService/GetApplicationsByUser"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).GetApplicationsByUser(ctx, req.(*GetApplicationsByUserRequest))
	})
}
func _JobService_GetJobCategories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetJobCategoriesRequest)
	if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(JobServiceServer).GetJobCategories(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/job.JobService/GetJobCategories"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).GetJobCategories(ctx, req.(*GetJobCategoriesRequest))
	})
}

// ── Client ────────────────────────────────────────────────────────────────

type JobServiceClient interface {
	CreateJob(ctx context.Context, in *CreateJobRequest, opts ...grpc.CallOption) (*JobResponse, error)
	GetJob(ctx context.Context, in *GetJobRequest, opts ...grpc.CallOption) (*JobResponse, error)
	UpdateJob(ctx context.Context, in *UpdateJobRequest, opts ...grpc.CallOption) (*JobResponse, error)
	DeleteJob(ctx context.Context, in *DeleteJobRequest, opts ...grpc.CallOption) (*DeleteJobResponse, error)
	ListJobs(ctx context.Context, in *ListJobsRequest, opts ...grpc.CallOption) (*ListJobsResponse, error)
	SearchJobs(ctx context.Context, in *SearchJobsRequest, opts ...grpc.CallOption) (*ListJobsResponse, error)
	ApplyJob(ctx context.Context, in *ApplyJobRequest, opts ...grpc.CallOption) (*ApplyJobResponse, error)
	GetApplications(ctx context.Context, in *GetApplicationsRequest, opts ...grpc.CallOption) (*GetApplicationsResponse, error)
	UpdateApplicationStatus(ctx context.Context, in *UpdateApplicationStatusRequest, opts ...grpc.CallOption) (*ApplyJobResponse, error)
	GetJobsByEmployer(ctx context.Context, in *GetJobsByEmployerRequest, opts ...grpc.CallOption) (*ListJobsResponse, error)
	GetApplicationsByUser(ctx context.Context, in *GetApplicationsByUserRequest, opts ...grpc.CallOption) (*GetApplicationsResponse, error)
	GetJobCategories(ctx context.Context, in *GetJobCategoriesRequest, opts ...grpc.CallOption) (*GetJobCategoriesResponse, error)
}

type jobServiceClient struct{ cc grpc.ClientConnInterface }

func NewJobServiceClient(cc grpc.ClientConnInterface) JobServiceClient { return &jobServiceClient{cc} }

func (c *jobServiceClient) CreateJob(ctx context.Context, in *CreateJobRequest, opts ...grpc.CallOption) (*JobResponse, error) {
	out := new(JobResponse); return out, c.cc.Invoke(ctx, "/job.JobService/CreateJob", in, out, opts...)
}
func (c *jobServiceClient) GetJob(ctx context.Context, in *GetJobRequest, opts ...grpc.CallOption) (*JobResponse, error) {
	out := new(JobResponse); return out, c.cc.Invoke(ctx, "/job.JobService/GetJob", in, out, opts...)
}
func (c *jobServiceClient) UpdateJob(ctx context.Context, in *UpdateJobRequest, opts ...grpc.CallOption) (*JobResponse, error) {
	out := new(JobResponse); return out, c.cc.Invoke(ctx, "/job.JobService/UpdateJob", in, out, opts...)
}
func (c *jobServiceClient) DeleteJob(ctx context.Context, in *DeleteJobRequest, opts ...grpc.CallOption) (*DeleteJobResponse, error) {
	out := new(DeleteJobResponse); return out, c.cc.Invoke(ctx, "/job.JobService/DeleteJob", in, out, opts...)
}
func (c *jobServiceClient) ListJobs(ctx context.Context, in *ListJobsRequest, opts ...grpc.CallOption) (*ListJobsResponse, error) {
	out := new(ListJobsResponse); return out, c.cc.Invoke(ctx, "/job.JobService/ListJobs", in, out, opts...)
}
func (c *jobServiceClient) SearchJobs(ctx context.Context, in *SearchJobsRequest, opts ...grpc.CallOption) (*ListJobsResponse, error) {
	out := new(ListJobsResponse); return out, c.cc.Invoke(ctx, "/job.JobService/SearchJobs", in, out, opts...)
}
func (c *jobServiceClient) ApplyJob(ctx context.Context, in *ApplyJobRequest, opts ...grpc.CallOption) (*ApplyJobResponse, error) {
	out := new(ApplyJobResponse); return out, c.cc.Invoke(ctx, "/job.JobService/ApplyJob", in, out, opts...)
}
func (c *jobServiceClient) GetApplications(ctx context.Context, in *GetApplicationsRequest, opts ...grpc.CallOption) (*GetApplicationsResponse, error) {
	out := new(GetApplicationsResponse); return out, c.cc.Invoke(ctx, "/job.JobService/GetApplications", in, out, opts...)
}
func (c *jobServiceClient) UpdateApplicationStatus(ctx context.Context, in *UpdateApplicationStatusRequest, opts ...grpc.CallOption) (*ApplyJobResponse, error) {
	out := new(ApplyJobResponse); return out, c.cc.Invoke(ctx, "/job.JobService/UpdateApplicationStatus", in, out, opts...)
}
func (c *jobServiceClient) GetJobsByEmployer(ctx context.Context, in *GetJobsByEmployerRequest, opts ...grpc.CallOption) (*ListJobsResponse, error) {
	out := new(ListJobsResponse); return out, c.cc.Invoke(ctx, "/job.JobService/GetJobsByEmployer", in, out, opts...)
}
func (c *jobServiceClient) GetApplicationsByUser(ctx context.Context, in *GetApplicationsByUserRequest, opts ...grpc.CallOption) (*GetApplicationsResponse, error) {
	out := new(GetApplicationsResponse); return out, c.cc.Invoke(ctx, "/job.JobService/GetApplicationsByUser", in, out, opts...)
}
func (c *jobServiceClient) GetJobCategories(ctx context.Context, in *GetJobCategoriesRequest, opts ...grpc.CallOption) (*GetJobCategoriesResponse, error) {
	out := new(GetJobCategoriesResponse); return out, c.cc.Invoke(ctx, "/job.JobService/GetJobCategories", in, out, opts...)
}

// ── Proto interface methods ──────────────────────────────────────────────
func (m *CreateJobRequest) Reset()         { *m = CreateJobRequest{} }
func (m *CreateJobRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *CreateJobRequest) ProtoMessage()   {}

func (m *JobResponse) Reset()         { *m = JobResponse{} }
func (m *JobResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *JobResponse) ProtoMessage()   {}

func (m *GetJobRequest) Reset()         { *m = GetJobRequest{} }
func (m *GetJobRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *GetJobRequest) ProtoMessage()   {}

func (m *UpdateJobRequest) Reset()         { *m = UpdateJobRequest{} }
func (m *UpdateJobRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *UpdateJobRequest) ProtoMessage()   {}

func (m *DeleteJobRequest) Reset()         { *m = DeleteJobRequest{} }
func (m *DeleteJobRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *DeleteJobRequest) ProtoMessage()   {}

func (m *DeleteJobResponse) Reset()         { *m = DeleteJobResponse{} }
func (m *DeleteJobResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *DeleteJobResponse) ProtoMessage()   {}

func (m *ListJobsRequest) Reset()         { *m = ListJobsRequest{} }
func (m *ListJobsRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *ListJobsRequest) ProtoMessage()   {}

func (m *ListJobsResponse) Reset()         { *m = ListJobsResponse{} }
func (m *ListJobsResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *ListJobsResponse) ProtoMessage()   {}

func (m *SearchJobsRequest) Reset()         { *m = SearchJobsRequest{} }
func (m *SearchJobsRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *SearchJobsRequest) ProtoMessage()   {}

func (m *ApplyJobRequest) Reset()         { *m = ApplyJobRequest{} }
func (m *ApplyJobRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *ApplyJobRequest) ProtoMessage()   {}

func (m *ApplyJobResponse) Reset()         { *m = ApplyJobResponse{} }
func (m *ApplyJobResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *ApplyJobResponse) ProtoMessage()   {}

func (m *GetApplicationsRequest) Reset()         { *m = GetApplicationsRequest{} }
func (m *GetApplicationsRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *GetApplicationsRequest) ProtoMessage()   {}

func (m *ApplicationResponse) Reset()         { *m = ApplicationResponse{} }
func (m *ApplicationResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *ApplicationResponse) ProtoMessage()   {}

func (m *GetApplicationsResponse) Reset()         { *m = GetApplicationsResponse{} }
func (m *GetApplicationsResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *GetApplicationsResponse) ProtoMessage()   {}

func (m *UpdateApplicationStatusRequest) Reset()         { *m = UpdateApplicationStatusRequest{} }
func (m *UpdateApplicationStatusRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *UpdateApplicationStatusRequest) ProtoMessage()   {}

func (m *GetJobsByEmployerRequest) Reset()         { *m = GetJobsByEmployerRequest{} }
func (m *GetJobsByEmployerRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *GetJobsByEmployerRequest) ProtoMessage()   {}

func (m *GetApplicationsByUserRequest) Reset()         { *m = GetApplicationsByUserRequest{} }
func (m *GetApplicationsByUserRequest) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *GetApplicationsByUserRequest) ProtoMessage()   {}

func (m *GetJobCategoriesResponse) Reset()         { *m = GetJobCategoriesResponse{} }
func (m *GetJobCategoriesResponse) String() string  { return fmt.Sprintf("%+v", *m) }
func (m *GetJobCategoriesResponse) ProtoMessage()   {}

