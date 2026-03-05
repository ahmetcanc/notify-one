package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahmetcanc/notify-one/internal/domain"
	"github.com/ahmetcanc/notify-one/internal/usecase"
)

type mockRepo struct{}

func (m *mockRepo) Create(ctx context.Context, n *domain.Notification) error { return nil }
func (m *mockRepo) GetByID(ctx context.Context, id string) (*domain.Notification, error) {
	return nil, nil
}
func (m *mockRepo) UpdateStatus(ctx context.Context, id string, status domain.NotificationStatus) error {
	return nil
}
func (m *mockRepo) GetByIdempotencyKey(ctx context.Context, key string) (*domain.Notification, error) {
	return nil, nil
}
func (m *mockRepo) List(ctx context.Context, f domain.NotificationFilter) ([]domain.Notification, error) {
	return nil, nil
}
func (m *mockRepo) BulkCancel(ctx context.Context, batchID string) error { return nil }

type mockQueue struct{}

func (m *mockQueue) PushToQueue(ctx context.Context, q string, id string) error { return nil }
func (m *mockQueue) GetQueueDepth(ctx context.Context, q string) (int64, error) { return 0, nil }

func TestNotificationHandler_Send(t *testing.T) {
	uc := usecase.NewNotificationUsecase(&mockRepo{}, &mockQueue{})

	h := NewNotificationHandler(uc)

	t.Run("Success - Valid Request", func(t *testing.T) {
		payload := map[string]interface{}{
			"recipient": "+905554443322",
			"channel":   "sms",
			"content":   "Handler Test Message",
			"priority":  "high",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/notifications", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		h.Send(w, req)

		if w.Code != http.StatusAccepted {
			t.Errorf("Want: 202, Got: %d", w.Code)
		}
	})
}

func TestNotificationHandler_List(t *testing.T) {
	h := NewNotificationHandler(usecase.NewNotificationUsecase(&mockRepo{}, &mockQueue{}))

	t.Run("Success - Filtering", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/notifications?status=sent&limit=10", nil)
		w := httptest.NewRecorder()

		h.List(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Want: 200, Got: %d", w.Code)
		}
	})
}

func TestNotificationHandler_Cancel_Fail(t *testing.T) {
	h := NewNotificationHandler(usecase.NewNotificationUsecase(&mockRepo{}, &mockQueue{}))

	t.Run("Fail - Missing BatchID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/notifications/batch//cancel", nil)
		w := httptest.NewRecorder()

		h.Cancel(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Want: 400, Got: %d", w.Code)
		}
	})
}
