package domain

import "context"

// ProviderResult holds external service response
type ProviderResult struct {
	ExternalID string
	Status     string
}

// NotificationProvider defines external delivery behavior
type NotificationProvider interface {
	Send(ctx context.Context, n *Notification) (*ProviderResult, error)
}
