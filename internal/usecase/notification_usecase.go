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
}

type NotificationUsecase struct {
	repo NotificationRepository
}

func NewNotificationUsecase(repo NotificationRepository) *NotificationUsecase {
	return &NotificationUsecase{
		repo: repo,
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

	// TODO: Push to Redis queue for worker

	return nil
}
