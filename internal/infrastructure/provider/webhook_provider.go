package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ahmetcanc/notify-one/internal/domain"
)

type WebhookProvider struct {
	URL    string
	Client *http.Client
}

func NewWebhookProvider(url string) *WebhookProvider {
	return &WebhookProvider{
		URL:    url,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *WebhookProvider) Send(ctx context.Context, n *domain.Notification) (*domain.ProviderResult, error) {
	// Payload mapping based on assessment specs
	payload := map[string]string{
		"to":      n.Recipient,
		"channel": string(n.Channel),
		"content": n.Content,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", p.URL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted { // Expecting 202
		return nil, fmt.Errorf("provider error: %d", resp.StatusCode)
	}

	var res struct {
		MessageID string `json:"messageId"`
		Status    string `json:"status"`
	}
	json.NewDecoder(resp.Body).Decode(&res)

	return &domain.ProviderResult{ExternalID: res.MessageID, Status: res.Status}, nil
}
