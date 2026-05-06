package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "config-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadAppConfig_Valid(t *testing.T) {
	path := writeTempConfig(t, `{
		"state_file": "state.json",
		"environments": ["prod", "staging"],
		"interval": "10m",
		"alert_level": "critical",
		"output": {"destination": "stdout", "format": "json"}
	}`)

	cfg, err := LoadAppConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.StateFile != "state.json" {
		t.Errorf("state_file: got %q, want %q", cfg.StateFile, "state.json")
	}
	if cfg.Interval != 10*time.Minute {
		t.Errorf("interval: got %v, want 10m", cfg.Interval)
	}
	if cfg.AlertLevel != "critical" {
		t.Errorf("alert_level: got %q, want critical", cfg.AlertLevel)
	}
}

func TestLoadAppConfig_NotFound(t *testing.T) {
	_, err := LoadAppConfig("/nonexistent/config.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadAppConfig_InvalidJSON(t *testing.T) {
	path := writeTempConfig(t, `{bad json}`)
	_, err := LoadAppConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestLoadAppConfig_InvalidInterval(t *testing.T) {
	path := writeTempConfig(t, `{"state_file":"s.json","interval":"notaduration"}`)
	_, err := LoadAppConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid interval, got nil")
	}
}
