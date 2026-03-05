package database

import (
	"context"
	"fmt"

	"github.com/ahmetcanc/notify-one/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresNotificationRepository struct {
	db *pgxpool.Pool
}

func NewPostgresNotificationRepository(db *pgxpool.Pool) *PostgresNotificationRepository {
	return &PostgresNotificationRepository{db: db}
}

// Create persists a new notification into the database
func (r *PostgresNotificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	query := `
		INSERT INTO notifications (
			id, batch_id, recipient, channel, content, priority, status, idempotency_key, scheduled_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.Exec(ctx, query,
		n.ID, n.BatchID, n.Recipient, n.Channel, n.Content, n.Priority, n.Status, n.IdempotencyKey, n.ScheduledAt, n.CreatedAt, n.UpdatedAt,
	)
	return err
}

// GetByID retrieves a notification record
func (r *PostgresNotificationRepository) GetByID(ctx context.Context, id string) (*domain.Notification, error) {
	var n domain.Notification

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid format: %w", err)
	}

	query := `SELECT id, recipient, channel, content, status, priority, created_at, updated_at 
			  FROM notifications WHERE id = $1`

	err = r.db.QueryRow(ctx, query, parsedID).Scan(
		&n.ID, &n.Recipient, &n.Channel, &n.Content, &n.Status, &n.Priority, &n.CreatedAt, &n.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &n, nil
}

// UpdateStatus updates the status of a specific notification
func (r *PostgresNotificationRepository) UpdateStatus(ctx context.Context, id string, status domain.NotificationStatus) error {
	query := `UPDATE notifications SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}
