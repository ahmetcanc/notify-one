package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ahmetcanc/notify-one/internal/domain"
	"github.com/google/uuid"
)

// NotificationRepository defines database operations
type NotificationRepository interface {
	Create(ctx context.Context, n *domain.Notification) error
	GetByID(ctx context.Context, id string) (*domain.Notification, error)
	UpdateStatus(ctx context.Context, id string, status domain.NotificationStatus) error
	GetByIdempotencyKey(ctx context.Context, key string) (*domain.Notification, error)
	List(ctx context.Context, f domain.NotificationFilter) ([]domain.Notification, error)
	BulkCancel(ctx context.Context, batchID string) error
}

type NotificationUsecase struct {
	repo  NotificationRepository
	queue QueueRepository
}

func NewNotificationUsecase(repo NotificationRepository, queue QueueRepository) *NotificationUsecase {
	return &NotificationUsecase{
		repo:  repo,
		queue: queue,
	}
}

type QueueRepository interface {
	// PushToQueue handles priority-based routing
	PushToQueue(ctx context.Context, queueName string, notificationID string) error
	GetQueueDepth(ctx context.Context, queueName string) (int64, error)
}

// Execute handles single notification logic
func (u *NotificationUsecase) Execute(ctx context.Context, n *domain.Notification) error {
	if err := u.validateNotification(*n); err != nil {
		return err
	}
	if n.IdempotencyKey != nil && *n.IdempotencyKey != "" {
		existing, err := u.repo.GetByIdempotencyKey(ctx, *n.IdempotencyKey)
		if err == nil && existing != nil {
			log.Printf("Duplicate request detected for key: %s, skipping...", *n.IdempotencyKey)
			return nil
		}
	}
	// Initialize mandatory fields before DB insert
	n.ID = uuid.New()
	n.Status = domain.StatusPending
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()

	if err := u.repo.Create(ctx, n); err != nil {
		return err
	}

	queueName := fmt.Sprintf("notification_queue_%s", n.Priority)
	return u.queue.PushToQueue(ctx, queueName, n.ID.String())
}

// ExecuteBatch handles batch notification logic
func (u *NotificationUsecase) ExecuteBatch(ctx context.Context, notifications []domain.Notification) (string, error) {
	if len(notifications) > 1000 {
		return "", fmt.Errorf("batch size exceeds maximum limit of 1000")
	}

	for _, n := range notifications {
		if err := u.validateNotification(n); err != nil {
			return "", fmt.Errorf("validation failed for recipient %s: %v", n.Recipient, err)
		}
	}

	batchID := uuid.New()
	for i := range notifications {
		n := &notifications[i]
		n.ID = uuid.New()
		n.BatchID = &batchID
		n.Status = domain.StatusPending
		n.CreatedAt = time.Now()
		n.UpdatedAt = time.Now()

		if err := u.repo.Create(ctx, n); err != nil {
			return "", err
		}

		// Route to specific priority queue
		queueName := fmt.Sprintf("notification_queue_%s", n.Priority)
		u.queue.PushToQueue(ctx, queueName, n.ID.String())
	}
	return batchID.String(), nil
}

func (u *NotificationUsecase) List(ctx context.Context, filter domain.NotificationFilter) ([]domain.Notification, error) {
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	return u.repo.List(ctx, filter)
}

func (u *NotificationUsecase) CancelBatch(ctx context.Context, batchID string) error {
	return u.repo.BulkCancel(ctx, batchID)
}

func (u *NotificationUsecase) GetSystemMetrics(ctx context.Context) (map[string]int64, error) {
	priorities := []string{"high", "normal", "low"}
	metrics := make(map[string]int64)

	for _, p := range priorities {
		queueName := fmt.Sprintf("notification_queue_%s", p)
		depth, err := u.queue.GetQueueDepth(ctx, queueName)
		if err != nil {
			metrics[p] = 0
			continue
		}
		metrics[p] = depth
	}

	return metrics, nil
}

func (u *NotificationUsecase) validateNotification(n domain.Notification) error {
	if n.Recipient == "" {
		return fmt.Errorf("recipient is required")
	}
	if n.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}

	validPriorities := map[domain.NotificationPriority]bool{
		domain.PriorityHigh:   true,
		domain.PriorityNormal: true,
		domain.PriorityLow:    true,
	}

	if !validPriorities[n.Priority] {
		return fmt.Errorf("invalid priority: %v", n.Priority)
	}

	validChannels := map[domain.NotificationChannel]bool{
		domain.ChannelSMS:   true,
		domain.ChannelEmail: true,
		domain.ChannelPush:  true,
	}

	if !validChannels[n.Channel] {
		return fmt.Errorf("invalid channel: %v", n.Channel)
	}

	if n.Channel == domain.ChannelSMS && len(n.Content) > 160 {
		return fmt.Errorf("sms content exceeds 160 characters")
	}

	return nil
}
