package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookConfig holds configuration for a webhook notifier.
type WebhookConfig struct {
	URL     string
	Timeout time.Duration
	Headers map[string]string
}

// WebhookPayload is the JSON body sent to the webhook endpoint.
type WebhookPayload struct {
	Timestamp  time.Time `json:"timestamp"`
	Source     string    `json:"source"`
	Target     string    `json:"target"`
	DriftCount int       `json:"drift_count"`
	Severity   string    `json:"severity"`
	Details    []string  `json:"details,omitempty"`
}

// WebhookNotifier sends drift alerts to an HTTP webhook endpoint.
type WebhookNotifier struct {
	cfg    WebhookConfig
	client *http.Client
}

// NewWebhookNotifier creates a WebhookNotifier with the given config.
func NewWebhookNotifier(cfg WebhookConfig) (*WebhookNotifier, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("webhook URL must not be empty")
	}
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &WebhookNotifier{
		cfg:    cfg,
		client: &http.Client{Timeout: timeout},
	}, nil
}

// Send posts the payload to the configured webhook URL.
func (w *WebhookNotifier) Send(payload WebhookPayload) error {
	if payload.Timestamp.IsZero() {
		payload.Timestamp = time.Now().UTC()
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, w.cfg.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range w.cfg.Headers {
		req.Header.Set(k, v)
	}
	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
