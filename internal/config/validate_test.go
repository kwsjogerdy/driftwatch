package config

import (
	"testing"
	"time"
)

func baseConfig() *AppConfig {
	return &AppConfig{
		StateFile:    "state.json",
		Environments: []string{"prod", "staging"},
		Interval:     5 * time.Minute,
		AlertLevel:   "warning",
		Output:       OutputConfig{Destination: "stdout", Format: "text"},
	}
}

func TestValidate_Valid(t *testing.T) {
	if err := Validate(baseConfig()); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_MissingStateFile(t *testing.T) {
	cfg := baseConfig()
	cfg.StateFile = ""
	if err := Validate(cfg); err != ErrMissingStateFile {
		t.Errorf("got %v, want %v", err, ErrMissingStateFile)
	}
}

func TestValidate_TooFewEnvironments(t *testing.T) {
	cfg := baseConfig()
	cfg.Environments = []string{"prod"}
	if err := Validate(cfg); err != ErrMissingEnvs {
		t.Errorf("got %v, want %v", err, ErrMissingEnvs)
	}
}

func TestValidate_InvalidAlertLevel(t *testing.T) {
	cfg := baseConfig()
	cfg.AlertLevel = "debug"
	if err := Validate(cfg); err == nil {
		t.Error("expected error for invalid alert_level, got nil")
	}
}

func TestValidate_FileDestinationMissingPath(t *testing.T) {
	cfg := baseConfig()
	cfg.Output.Destination = "file"
	cfg.Output.FilePath = ""
	if err := Validate(cfg); err == nil {
		t.Error("expected error when destination=file and file_path empty")
	}
}

func TestApplyDefaults_FillsZeroValues(t *testing.T) {
	cfg := &AppConfig{}
	ApplyDefaults(cfg)
	if cfg.Interval != 5*time.Minute {
		t.Errorf("interval default: got %v, want 5m", cfg.Interval)
	}
	if cfg.AlertLevel != "warning" {
		t.Errorf("alert_level default: got %q, want warning", cfg.AlertLevel)
	}
	if cfg.Output.Format != "text" {
		t.Errorf("format default: got %q, want text", cfg.Output.Format)
	}
}
