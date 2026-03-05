package database

import (
	"context"
	"fmt"
	"time"

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

	query := `SELECT id, recipient, channel, content, status, priority, 
                     retry_count, next_retry_at, created_at, updated_at 
              FROM notifications WHERE id = $1`

	err = r.db.QueryRow(ctx, query, parsedID).Scan(
		&n.ID, &n.Recipient, &n.Channel, &n.Content, &n.Status, &n.Priority,
		&n.RetryCount, &n.NextRetryAt, &n.CreatedAt, &n.UpdatedAt,
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

// GetByIdempotencyKey finds a notification by its unique idempotency key
func (r *PostgresNotificationRepository) GetByIdempotencyKey(ctx context.Context, key string) (*domain.Notification, error) {
	var n domain.Notification
	query := `SELECT id, recipient, channel, status FROM notifications WHERE idempotency_key = $1 LIMIT 1`

	err := r.db.QueryRow(ctx, query, key).Scan(&n.ID, &n.Recipient, &n.Channel, &n.Status)
	if err != nil {
		return nil, err // pgx.ErrNoRows durumunda üst katman bunu yönetir
	}
	return &n, nil
}

// List notifications with dynamic filters and pagination
func (r *PostgresNotificationRepository) List(ctx context.Context, f domain.NotificationFilter) ([]domain.Notification, error) {
	query := `SELECT id, recipient, channel, content, priority, status, batch_id, created_at 
			  FROM notifications WHERE 1=1`
	args := []interface{}{}
	counter := 1

	if f.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", counter)
		args = append(args, f.Status)
		counter++
	}
	if f.Channel != "" {
		query += fmt.Sprintf(" AND channel = $%d", counter)
		args = append(args, f.Channel)
		counter++
	}
	if f.BatchID != nil {
		query += fmt.Sprintf(" AND batch_id = $%d", counter)
		args = append(args, *f.BatchID)
		counter++
	}

	if f.StartDate != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", counter)
		args = append(args, *f.StartDate)
		counter++
	}
	if f.EndDate != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", counter)
		args = append(args, *f.EndDate)
		counter++
	}

	if f.Priority != "" {
		query += fmt.Sprintf(" AND priority = $%d", counter)
		args = append(args, f.Priority)
		counter++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", counter, counter+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		var n domain.Notification
		err := rows.Scan(&n.ID, &n.Recipient, &n.Channel, &n.Content, &n.Priority, &n.Status, &n.BatchID, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

// BulkCancel updates pending notifications to cancelled
func (r *PostgresNotificationRepository) BulkCancel(ctx context.Context, batchID string) error {
	query := `UPDATE notifications SET status = 'cancelled', updated_at = NOW() 
			  WHERE batch_id = $1 AND status = 'pending'`
	_, err := r.db.Exec(ctx, query, batchID)
	return err
}

// UpdateRetryInfo updates notification retry metadata
func (r *PostgresNotificationRepository) UpdateRetryInfo(ctx context.Context, id uuid.UUID, retryCount int, nextRetryAt time.Time) error {
	query := `
        UPDATE notifications 
        SET retry_count = $1, next_retry_at = $2, updated_at = NOW() 
        WHERE id = $3`

	_, err := r.db.Exec(ctx, query, retryCount, nextRetryAt, id)
	return err
}
