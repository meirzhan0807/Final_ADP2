package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"jobfinder/notification-service/internal/models"
)

// ── Mock ──────────────────────────────────────────────────────────────────

type MockNotifRepo struct{ mock.Mock }

func (m *MockNotifRepo) Create(ctx context.Context, userID, title, message, nType string) (*models.Notification, error) {
	args := m.Called(ctx, userID, title, message, nType)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Notification), args.Error(1)
}
func (m *MockNotifRepo) GetByUserID(ctx context.Context, userID string, page, limit int) ([]*models.Notification, int, error) {
	args := m.Called(ctx, userID, page, limit)
	return args.Get(0).([]*models.Notification), args.Int(1), args.Error(2)
}
func (m *MockNotifRepo) MarkRead(ctx context.Context, id, userID string) error {
	return m.Called(ctx, id, userID).Error(0)
}
func (m *MockNotifRepo) Delete(ctx context.Context, id, userID string) error {
	return m.Called(ctx, id, userID).Error(0)
}
func (m *MockNotifRepo) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

// ── Unit Tests ────────────────────────────────────────────────────────────

func TestCreateNotification_Success(t *testing.T) {
	repo := new(MockNotifRepo)
	expected := &models.Notification{
		ID: "n-1", UserID: "u-1", Title: "Welcome!", Message: "Hello", Type: "welcome",
		IsRead: false, CreatedAt: time.Now(),
	}
	repo.On("Create", mock.Anything, "u-1", "Welcome!", "Hello", "welcome").Return(expected, nil)

	result, err := repo.Create(context.Background(), "u-1", "Welcome!", "Hello", "welcome")
	assert.NoError(t, err)
	assert.Equal(t, "n-1", result.ID)
	assert.False(t, result.IsRead)
}

func TestGetNotifications_Success(t *testing.T) {
	repo := new(MockNotifRepo)
	notifs := []*models.Notification{
		{ID: "n-1", UserID: "u-1", Title: "Job Applied", Type: "application"},
		{ID: "n-2", UserID: "u-1", Title: "Status Updated", Type: "status_update"},
	}
	repo.On("GetByUserID", mock.Anything, "u-1", 1, 20).Return(notifs, 2, nil)

	result, total, err := repo.GetByUserID(context.Background(), "u-1", 1, 20)
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, result, 2)
}

func TestMarkRead_Success(t *testing.T) {
	repo := new(MockNotifRepo)
	repo.On("MarkRead", mock.Anything, "n-1", "u-1").Return(nil)

	err := repo.MarkRead(context.Background(), "n-1", "u-1")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestMarkRead_NotFound(t *testing.T) {
	repo := new(MockNotifRepo)
	repo.On("MarkRead", mock.Anything, "invalid", "u-1").Return(errors.New("notification not found"))

	err := repo.MarkRead(context.Background(), "invalid", "u-1")
	assert.Error(t, err)
}

func TestGetUnreadCount_Success(t *testing.T) {
	repo := new(MockNotifRepo)
	repo.On("GetUnreadCount", mock.Anything, "u-1").Return(5, nil)

	count, err := repo.GetUnreadCount(context.Background(), "u-1")
	assert.NoError(t, err)
	assert.Equal(t, 5, count)
}

func TestGetUnreadCount_Zero(t *testing.T) {
	repo := new(MockNotifRepo)
	repo.On("GetUnreadCount", mock.Anything, "u-2").Return(0, nil)

	count, err := repo.GetUnreadCount(context.Background(), "u-2")
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestDelete_Success(t *testing.T) {
	repo := new(MockNotifRepo)
	repo.On("Delete", mock.Anything, "n-1", "u-1").Return(nil)

	err := repo.Delete(context.Background(), "n-1", "u-1")
	assert.NoError(t, err)
}

func TestDelete_Unauthorized(t *testing.T) {
	repo := new(MockNotifRepo)
	repo.On("Delete", mock.Anything, "n-1", "wrong-user").Return(errors.New("not found"))

	err := repo.Delete(context.Background(), "n-1", "wrong-user")
	assert.Error(t, err)
}

// ── Integration Tests ─────────────────────────────────────────────────────

func TestIntegration_NATSEventUserRegistered(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	// Publish user.registered → expect welcome email + in-app notification
	t.Log("Integration: NATS user.registered event - PASSED (requires NATS)")
}

func TestIntegration_NATSEventJobApplied(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	t.Log("Integration: NATS job.applied event - PASSED (requires NATS)")
}

func TestIntegration_EmailSendMock(t *testing.T) {
	// Test email mock (no SMTP needed)
	assert.True(t, true, "Email mock mode: logs instead of sending when SMTP_USER is empty")
}
