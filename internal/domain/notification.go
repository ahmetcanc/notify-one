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
	ID             uuid.UUID            `json:"id"`
	BatchID        *uuid.UUID           `json:"batch_id,omitempty"`
	Recipient      string               `json:"recipient"`
	Channel        NotificationChannel  `json:"channel"`
	Content        string               `json:"content"`
	Priority       NotificationPriority `json:"priority"`
	Status         NotificationStatus   `json:"status"`
	IdempotencyKey *string              `json:"idempotency_key,omitempty"`
	ScheduledAt    *time.Time           `json:"scheduled_at,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	RetryCount     int                  `json:"retry_count"`
	NextRetryAt    time.Time            `json:"next_retry_at"`
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
