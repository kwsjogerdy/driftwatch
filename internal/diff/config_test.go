package diff_test

import (
	"testing"

	"github.com/driftwatch/internal/diff"
)

func TestDefaultConfig_Values(t *testing.T) {
	cfg := diff.DefaultConfig()
	if cfg.DefaultSeverity != "warning" {
		t.Errorf("expected default severity 'warning', got %q", cfg.DefaultSeverity)
	}
	if len(cfg.CriticalKeys) != 0 {
		t.Errorf("expected empty critical keys, got %v", cfg.CriticalKeys)
	}
	if len(cfg.IgnoreKeys) != 0 {
		t.Errorf("expected empty ignore keys, got %v", cfg.IgnoreKeys)
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := diff.DefaultConfig()
	if err := diff.Validate(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_InvalidSeverity(t *testing.T) {
	cfg := diff.DefaultConfig()
	cfg.DefaultSeverity = "unknown"
	if err := diff.Validate(cfg); err == nil {
		t.Error("expected error for invalid severity")
	}
}

func TestValidate_CriticalKeyInIgnore(t *testing.T) {
	cfg := diff.DefaultConfig()
	cfg.CriticalKeys = []string{"db_password"}
	cfg.IgnoreKeys = []string{"db_password"}
	if err := diff.Validate(cfg); err == nil {
		t.Error("expected error when critical key is also in ignore list")
	}
}

func TestApplyIgnore_FiltersKeys(t *testing.T) {
	cfg := diff.DefaultConfig()
	cfg.IgnoreKeys = []string{"skip_me", "also_skip"}

	input := map[string]string{
		"keep":      "value1",
		"skip_me":   "value2",
		"also_skip": "value3",
	}

	result := diff.ApplyIgnore(cfg, input)

	if _, ok := result["keep"]; !ok {
		t.Error("expected 'keep' key to be present")
	}
	if _, ok := result["skip_me"]; ok {
		t.Error("expected 'skip_me' to be removed")
	}
	if _, ok := result["also_skip"]; ok {
		t.Error("expected 'also_skip' to be removed")
	}
}

func TestApplyIgnore_EmptyIgnoreList(t *testing.T) {
	cfg := diff.DefaultConfig()
	input := map[string]string{"a": "1", "b": "2"}
	result := diff.ApplyIgnore(cfg, input)
	if len(result) != len(input) {
		t.Errorf("expected %d keys, got %d", len(input), len(result))
	}
}
