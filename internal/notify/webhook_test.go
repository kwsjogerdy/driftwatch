package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/notify"
)

func TestNewWebhookNotifier_EmptyURL(t *testing.T) {
	_, err := notify.NewWebhookNotifier(notify.WebhookConfig{})
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestParseWebhookConfig_Valid(t *testing.T) {
	cfg, err := notify.ParseWebhookConfig("https://example.com/hook", 5, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("expected 5s timeout, got %v", cfg.Timeout)
	}
}

func TestParseWebhookConfig_InvalidScheme(t *testing.T) {
	_, err := notify.ParseWebhookConfig("ftp://example.com/hook", 0, nil)
	if err == nil {
		t.Fatal("expected error for non-http scheme")
	}
}

func TestParseWebhookConfig_NegativeTimeout(t *testing.T) {
	_, err := notify.ParseWebhookConfig("https://example.com/hook", -1, nil)
	if err == nil {
		t.Fatal("expected error for negative timeout")
	}
}

func TestSend_Success(t *testing.T) {
	var received notify.WebhookPayload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json content-type")
		}
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, err := notify.NewWebhookNotifier(notify.WebhookConfig{URL: server.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	payload := notify.WebhookPayload{
		Source:     "staging",
		Target:     "production",
		DriftCount: 3,
		Severity:   "critical",
		Details:    []string{"key_a mismatch"},
	}
	if err := n.Send(payload); err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if received.DriftCount != 3 {
		t.Errorf("expected drift_count 3, got %d", received.DriftCount)
	}
	if received.Timestamp.IsZero() {
		t.Error("expected timestamp to be set")
	}
}

func TestSend_Non2xxStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := notify.NewWebhookNotifier(notify.WebhookConfig{URL: server.URL})
	err := n.Send(notify.WebhookPayload{Source: "a", Target: "b"})
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}
