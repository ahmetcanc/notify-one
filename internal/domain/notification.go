package domain

import (
	"time"

	"github.com/google/uuid"
)

type NotificationChannel string
type NotificationStatus string
type NotificationPriority string

const (
	ChannelSMS   NotificationChannel = "sms"
	ChannelEmail NotificationChannel = "email"
	ChannelPush  NotificationChannel = "push"

	StatusPending    NotificationStatus = "pending"
	StatusProcessing NotificationStatus = "processing"
	StatusSent       NotificationStatus = "sent"
	StatusFailed     NotificationStatus = "failed"
	StatusCancelled  NotificationStatus = "cancelled"

	PriorityLow    NotificationPriority = "low"
	PriorityNormal NotificationPriority = "normal"
	PriorityHigh   NotificationPriority = "high"
)

type Notification struct {
	ID             uuid.UUID            `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	BatchID        *uuid.UUID           `json:"batch_id,omitempty" example:"722e8400-e29b-41d4-a716-446655440000"`
	Recipient      string               `json:"recipient" example:"+905554443322"`
	Channel        NotificationChannel  `json:"channel" example:"sms" enums:"sms,email,push"`
	Content        string               `json:"content" example:"Successful launch! 🚀"`
	Priority       NotificationPriority `json:"priority" example:"high" enums:"low,normal,high"`
	Status         NotificationStatus   `json:"status" example:"pending"`
	IdempotencyKey *string              `json:"idempotency_key,omitempty" example:"unique-key-123"`
	ScheduledAt    *time.Time           `json:"scheduled_at,omitempty" example:"2026-03-06T10:00:00Z"`
	CreatedAt      time.Time            `json:"created_at" example:"2026-03-06T02:50:00Z"`
	UpdatedAt      time.Time            `json:"updated_at" example:"2026-03-06T02:50:00Z"`
	RetryCount     int                  `json:"retry_count" example:"0"`
	NextRetryAt    time.Time            `json:"next_retry_at" example:"2026-03-06T03:00:00Z"`
}

type NotificationFilter struct {
	Status    string
	Channel   string
	Priority  string
	BatchID   *string
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}
