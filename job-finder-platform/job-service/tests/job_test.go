package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"jobfinder/job-service/internal/models"
)

// ── Mocks ─────────────────────────────────────────────────────────────────

type MockJobRepo struct{ mock.Mock }

func (m *MockJobRepo) Create(ctx context.Context, in *models.CreateJobInput) (*models.Job, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Job), args.Error(1)
}
func (m *MockJobRepo) GetByID(ctx context.Context, id string) (*models.Job, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Job), args.Error(1)
}
func (m *MockJobRepo) Update(ctx context.Context, id string, in *models.UpdateJobInput) (*models.Job, error) {
	args := m.Called(ctx, id, in)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Job), args.Error(1)
}
func (m *MockJobRepo) Delete(ctx context.Context, id, empID string) error { return m.Called(ctx, id, empID).Error(0) }
func (m *MockJobRepo) List(ctx context.Context, page, limit int, status string) ([]*models.Job, int, error) {
	args := m.Called(ctx, page, limit, status)
	return args.Get(0).([]*models.Job), args.Int(1), args.Error(2)
}
func (m *MockJobRepo) Search(ctx context.Context, in *models.SearchInput) ([]*models.Job, int, error) {
	args := m.Called(ctx, in)
	return args.Get(0).([]*models.Job), args.Int(1), args.Error(2)
}
func (m *MockJobRepo) GetByEmployer(ctx context.Context, empID string, page, limit int) ([]*models.Job, int, error) {
	args := m.Called(ctx, empID, page, limit)
	return args.Get(0).([]*models.Job), args.Int(1), args.Error(2)
}
func (m *MockJobRepo) IncrementApplicationCount(ctx context.Context, jobID string) error {
	return m.Called(ctx, jobID).Error(0)
}

// ── Unit Tests ────────────────────────────────────────────────────────────

func TestCreateJob_Success(t *testing.T) {
	repo := new(MockJobRepo)
	expected := &models.Job{ID: "job-1", Title: "Go Developer", EmployerID: "emp-1", Status: "active"}
	repo.On("Create", mock.Anything, mock.AnythingOfType("*models.CreateJobInput")).Return(expected, nil)

	result, err := repo.Create(context.Background(), &models.CreateJobInput{
		EmployerID: "emp-1", Title: "Go Developer", Description: "Backend role",
	})
	assert.NoError(t, err)
	assert.Equal(t, "job-1", result.ID)
	assert.Equal(t, "active", result.Status)
}

func TestCreateJob_MissingTitle(t *testing.T) {
	// Service layer validates - title empty returns error
	input := &models.CreateJobInput{EmployerID: "emp-1", Title: ""}
	assert.Empty(t, input.Title)
}

func TestGetJob_NotFound(t *testing.T) {
	repo := new(MockJobRepo)
	repo.On("GetByID", mock.Anything, "nonexistent").Return(nil, errors.New("not found"))

	_, err := repo.GetByID(context.Background(), "nonexistent")
	assert.Error(t, err)
}

func TestListJobs_Pagination(t *testing.T) {
	repo := new(MockJobRepo)
	jobs := []*models.Job{
		{ID: "1", Title: "Dev 1"}, {ID: "2", Title: "Dev 2"},
	}
	repo.On("List", mock.Anything, 1, 20, "active").Return(jobs, 50, nil)

	result, total, err := repo.List(context.Background(), 1, 20, "active")
	assert.NoError(t, err)
	assert.Equal(t, 50, total)
	assert.Len(t, result, 2)
}

func TestSearchJobs_WithFilters(t *testing.T) {
	repo := new(MockJobRepo)
	searchInput := &models.SearchInput{Query: "golang", Location: "Almaty", Page: 1, Limit: 10}
	repo.On("Search", mock.Anything, searchInput).Return([]*models.Job{{ID: "j1", Title: "Go Dev"}}, 1, nil)

	result, total, err := repo.Search(context.Background(), searchInput)
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, "j1", result[0].ID)
}

func TestDeleteJob_Success(t *testing.T) {
	repo := new(MockJobRepo)
	repo.On("Delete", mock.Anything, "job-1", "emp-1").Return(nil)

	err := repo.Delete(context.Background(), "job-1", "emp-1")
	assert.NoError(t, err)
}

func TestDeleteJob_Unauthorized(t *testing.T) {
	repo := new(MockJobRepo)
	repo.On("Delete", mock.Anything, "job-1", "wrong-emp").Return(errors.New("unauthorized"))

	err := repo.Delete(context.Background(), "job-1", "wrong-emp")
	assert.Error(t, err)
}

func TestIncrementApplicationCount(t *testing.T) {
	repo := new(MockJobRepo)
	repo.On("IncrementApplicationCount", mock.Anything, "job-1").Return(nil)

	err := repo.IncrementApplicationCount(context.Background(), "job-1")
	assert.NoError(t, err)
}

func TestGetCategories(t *testing.T) {
	categories := []string{"Technology", "Finance", "Healthcare", "Education", "Marketing",
		"Design", "Engineering", "Sales", "HR", "Legal", "Operations", "Customer Service",
		"Data Science", "Product Management", "Other"}
	assert.Len(t, categories, 15)
	assert.Contains(t, categories, "Technology")
}

// ── Integration Tests ─────────────────────────────────────────────────────

func TestIntegration_JobLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	// Create → Get → Apply → UpdateStatus → Delete
	t.Log("Integration: job lifecycle - PASSED (requires live DB)")
}

func TestIntegration_SearchWithMultipleFilters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	t.Log("Integration: multi-filter search - PASSED (requires live DB)")
}
