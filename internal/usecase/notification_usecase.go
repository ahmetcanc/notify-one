package usecase

import (
	"context"
	"time"

	"github.com/ahmetcanc/notify-one/internal/domain"
	"github.com/google/uuid"
)

// NotificationRepository defines database operations
type NotificationRepository interface {
	Create(ctx context.Context, n *domain.Notification) error
	GetByID(ctx context.Context, id string) (*domain.Notification, error)
	UpdateStatus(ctx context.Context, id string, status domain.NotificationStatus) error
}

// QueueRepository defines the behavior for queuing notifications
type QueueRepository interface {
	PushToQueue(ctx context.Context, notificationID string) error
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

// Execute handles notification creation logic
func (u *NotificationUsecase) Execute(ctx context.Context, n *domain.Notification) error {
	// Initialize metadata
	n.ID = uuid.New()
	n.Status = domain.StatusPending
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()

	// Persist to database
	if err := u.repo.Create(ctx, n); err != nil {
		return err
	}

	return u.queue.PushToQueue(ctx, n.ID.String())
}

// ExecuteBatch handles batch notification logic
func (u *NotificationUsecase) ExecuteBatch(ctx context.Context, notifications []domain.Notification) error {
	for i := range notifications {
		n := &notifications[i]
		n.ID = uuid.New()
		n.Status = domain.StatusPending
		n.CreatedAt = time.Now()
		n.UpdatedAt = time.Now()

		if err := u.repo.Create(ctx, n); err != nil {
			return err
		}

		// Push each to queue
		if err := u.queue.PushToQueue(ctx, n.ID.String()); err != nil {
			return err
		}
	}
	return nil
}
