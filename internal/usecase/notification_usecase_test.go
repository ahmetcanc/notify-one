package usecase

import (
	"context"
	"testing"

	"github.com/ahmetcanc/notify-one/internal/domain"
	"github.com/google/uuid"
)

type MockRepository struct {
	CreateFunc              func(n *domain.Notification) error
	GetByIdempotencyKeyFunc func(key string) (*domain.Notification, error)
}

func (m *MockRepository) Create(ctx context.Context, n *domain.Notification) error {
	return m.CreateFunc(n)
}
func (m *MockRepository) GetByID(ctx context.Context, id string) (*domain.Notification, error) {
	return nil, nil
}
func (m *MockRepository) UpdateStatus(ctx context.Context, id string, status domain.NotificationStatus) error {
	return nil
}
func (m *MockRepository) GetByIdempotencyKey(ctx context.Context, key string) (*domain.Notification, error) {
	return m.GetByIdempotencyKeyFunc(key)
}
func (m *MockRepository) List(ctx context.Context, f domain.NotificationFilter) ([]domain.Notification, error) {
	return nil, nil
}
func (m *MockRepository) BulkCancel(ctx context.Context, batchID string) error { return nil }

type MockQueue struct {
	PushToQueueFunc func(queueName string, id string) error
}

func (m *MockQueue) PushToQueue(ctx context.Context, q string, id string) error {
	return m.PushToQueueFunc(q, id)
}
func (m *MockQueue) GetQueueDepth(ctx context.Context, q string) (int64, error) { return 0, nil }

func TestExecute_Success(t *testing.T) {
	// Setup
	mockRepo := &MockRepository{
		CreateFunc:              func(n *domain.Notification) error { return nil },
		GetByIdempotencyKeyFunc: func(key string) (*domain.Notification, error) { return nil, nil },
	}
	var capturedQueue, capturedID string
	mockQueue := &MockQueue{
		PushToQueueFunc: func(q string, id string) error {
			capturedQueue = q
			capturedID = id
			return nil
		},
	}
	uc := NewNotificationUsecase(mockRepo, mockQueue)

	n := &domain.Notification{
		Recipient: "+905551112233",
		Channel:   domain.ChannelSMS,
		Content:   "Test message",
		Priority:  domain.PriorityHigh,
	}

	// Action
	err := uc.Execute(context.Background(), n)

	// Assertions
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if n.ID == uuid.Nil {
		t.Error("Not created ID")
	}
	if n.Status != domain.StatusPending {
		t.Errorf("Wrong status: %v", n.Status)
	}
	if capturedQueue != "notification_queue_high" {
		t.Errorf("Wrong queue: %s", capturedQueue)
	}
	if capturedID != n.ID.String() {
		t.Error("Waited ID, but got different ID")
	}
}

func TestExecute_ValidationError(t *testing.T) {
	uc := NewNotificationUsecase(&MockRepository{}, &MockQueue{})

	n := &domain.Notification{Content: "Mesaj", Priority: domain.PriorityNormal}

	err := uc.Execute(context.Background(), n)

	if err == nil {
		t.Error("Waited validation error, but got nil")
	}
}

func TestExecuteBatch_LimitError(t *testing.T) {
	uc := NewNotificationUsecase(&MockRepository{}, &MockQueue{})

	oversizedBatch := make([]domain.Notification, 1001)
	for i := range oversizedBatch {
		oversizedBatch[i] = domain.Notification{
			Recipient: "test@test.com",
			Content:   "Batch Test",
			Priority:  domain.PriorityNormal,
			Channel:   domain.ChannelEmail,
		}
	}

	_, err := uc.ExecuteBatch(context.Background(), oversizedBatch)

	if err == nil {
		t.Error("1000 limitini aşan batch için hata bekleniyordu ama alınmadı")
	}
}
