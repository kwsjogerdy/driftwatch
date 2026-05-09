package notify

import (
	"fmt"
	"net/url"
	"time"
)

// DefaultWebhookTimeout is used when no timeout is specified.
const DefaultWebhookTimeout = 10 * time.Second

// ParseWebhookConfig builds and validates a WebhookConfig from raw values.
func ParseWebhookConfig(rawURL string, timeoutSec int, headers map[string]string) (WebhookConfig, error) {
	if rawURL == "" {
		return WebhookConfig{}, fmt.Errorf("webhook url is required")
	}
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return WebhookConfig{}, fmt.Errorf("webhook url is invalid: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return WebhookConfig{}, fmt.Errorf("webhook url scheme must be http or https, got %q", parsed.Scheme)
	}
	timeout := DefaultWebhookTimeout
	if timeoutSec > 0 {
		timeout = time.Duration(timeoutSec) * time.Second
	}
	if timeoutSec < 0 {
		return WebhookConfig{}, fmt.Errorf("webhook timeout must be non-negative")
	}
	h := make(map[string]string)
	for k, v := range headers {
		if k == "" {
			return WebhookConfig{}, fmt.Errorf("webhook header key must not be empty")
		}
		h[k] = v
	}
	return WebhookConfig{
		URL:     rawURL,
		Timeout: timeout,
		Headers: h,
	}, nil
}
