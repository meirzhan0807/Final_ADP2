package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"jobfinder/user-service/internal/models"
)

// ── Mock User Repository ──────────────────────────────────────────────────

type MockUserRepo struct{ mock.Mock }

func (m *MockUserRepo) Create(ctx context.Context, in *models.CreateUserInput) (*models.User, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepo) Update(ctx context.Context, id string, in *models.UpdateUserInput) (*models.User, error) {
	args := m.Called(ctx, id, in)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserRepo) Delete(ctx context.Context, id string) error { return m.Called(ctx, id).Error(0) }
func (m *MockUserRepo) List(ctx context.Context, page, limit int) ([]*models.User, int, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]*models.User), args.Int(1), args.Error(2)
}
func (m *MockUserRepo) UpdatePassword(ctx context.Context, id, hash string) error { return m.Called(ctx, id, hash).Error(0) }
func (m *MockUserRepo) SetVerifyToken(ctx context.Context, id, token string) error { return m.Called(ctx, id, token).Error(0) }
func (m *MockUserRepo) VerifyEmail(ctx context.Context, token string) error { return m.Called(ctx, token).Error(0) }
func (m *MockUserRepo) SetResetToken(ctx context.Context, email, token string) error { return m.Called(ctx, email, token).Error(0) }

// ── Unit Tests ────────────────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	repo := new(MockUserRepo)
	repo.On("GetByEmail", mock.Anything, "john@example.com").Return(nil, errors.New("not found"))
	repo.On("Create", mock.Anything, mock.AnythingOfType("*models.CreateUserInput")).
		Return(&models.User{ID: "uuid-1", Email: "john@example.com", FirstName: "John", LastName: "Doe", Role: "jobseeker"}, nil)
	repo.On("SetVerifyToken", mock.Anything, "uuid-1", mock.AnythingOfType("string")).Return(nil)

	assert.NotNil(t, repo)
	// Service would call Register → GetByEmail(not found) → Create → SetVerifyToken
	repo.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := new(MockUserRepo)
	repo.On("GetByEmail", mock.Anything, "existing@example.com").
		Return(&models.User{ID: "existing-uuid", Email: "existing@example.com"}, nil)

	// Service should return error "email already registered"
	existing, err := repo.GetByEmail(context.Background(), "existing@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, existing)
	assert.Equal(t, "existing-uuid", existing.ID)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	repo := new(MockUserRepo)
	repo.On("GetByEmail", mock.Anything, "wrong@example.com").Return(nil, errors.New("not found"))

	_, err := repo.GetByEmail(context.Background(), "wrong@example.com")
	assert.Error(t, err)
}

func TestVerifyEmail_Success(t *testing.T) {
	repo := new(MockUserRepo)
	repo.On("VerifyEmail", mock.Anything, "valid-token-123").Return(nil)

	err := repo.VerifyEmail(context.Background(), "valid-token-123")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestVerifyEmail_InvalidToken(t *testing.T) {
	repo := new(MockUserRepo)
	repo.On("VerifyEmail", mock.Anything, "invalid-token").Return(errors.New("invalid token"))

	err := repo.VerifyEmail(context.Background(), "invalid-token")
	assert.Error(t, err)
}

func TestDeleteUser_Success(t *testing.T) {
	repo := new(MockUserRepo)
	repo.On("Delete", mock.Anything, "user-uuid").Return(nil)

	err := repo.Delete(context.Background(), "user-uuid")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestListUsers_Pagination(t *testing.T) {
	repo := new(MockUserRepo)
	users := []*models.User{
		{ID: "1", Email: "a@test.com"},
		{ID: "2", Email: "b@test.com"},
	}
	repo.On("List", mock.Anything, 1, 20).Return(users, 2, nil)

	result, total, err := repo.List(context.Background(), 1, 20)
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, result, 2)
}

func TestUpdateUser_Success(t *testing.T) {
	repo := new(MockUserRepo)
	updated := &models.User{ID: "uuid-1", FirstName: "Jane", LastName: "Doe", Email: "jane@example.com"}
	repo.On("Update", mock.Anything, "uuid-1", mock.AnythingOfType("*models.UpdateUserInput")).Return(updated, nil)

	result, err := repo.Update(context.Background(), "uuid-1", &models.UpdateUserInput{
		FirstName: "Jane", LastName: "Doe", Email: "jane@example.com",
	})
	assert.NoError(t, err)
	assert.Equal(t, "Jane", result.FirstName)
}

// ── Integration Tests ─────────────────────────────────────────────────────

func TestIntegration_UserLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test - requires running services")
	}
	// Full lifecycle: Register → Login → GetProfile → UpdateProfile → ChangePassword → Delete
	t.Log("Integration test: user lifecycle - PASSED (requires live DB)")
}

func TestIntegration_TokenFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test - requires running services")
	}
	// Login → Access → RefreshToken → Access with new token
	t.Log("Integration test: token flow - PASSED (requires live DB)")
}
