package policy_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/policy"
)

func baseConfig() policy.Config {
	return policy.DefaultConfig()
}

func TestDefaultConfig_Values(t *testing.T) {
	cfg := policy.DefaultConfig()
	if cfg.MaxRules <= 0 {
		t.Errorf("expected positive MaxRules, got %d", cfg.MaxRules)
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := baseConfig()
	cfg.Rules = []policy.Rule{
		{
			ID:        "r1",
			KeyPrefix: "db.",
			Severity:  "critical",
			ExpiresAt: time.Now().Add(time.Hour),
		},
	}
	if err := policy.Validate(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_DisabledSkipsChecks(t *testing.T) {
	cfg := baseConfig()
	cfg.Enabled = false
	cfg.MaxRules = -1 // would normally fail
	if err := policy.Validate(cfg); err != nil {
		t.Errorf("expected no error when disabled, got: %v", err)
	}
}

func TestValidate_NegativeMaxRules(t *testing.T) {
	cfg := baseConfig()
	cfg.MaxRules = -5
	if err := policy.Validate(cfg); err == nil {
		t.Error("expected error for negative MaxRules")
	}
}

func TestValidate_TooManyRules(t *testing.T) {
	cfg := baseConfig()
	cfg.MaxRules = 2
	for i := 0; i < 3; i++ {
		cfg.Rules = append(cfg.Rules, policy.Rule{
			ID:        fmt.Sprintf("r%d", i),
			KeyPrefix: "key.",
			Severity:  "warning",
			ExpiresAt: time.Now().Add(time.Hour),
		})
	}
	if err := policy.Validate(cfg); err == nil {
		t.Error("expected error when rules exceed MaxRules")
	}
}

func TestActiveRules_FiltersExpired(t *testing.T) {
	now := time.Now()
	cfg := baseConfig()
	cfg.Rules = []policy.Rule{
		{ID: "active", KeyPrefix: "a.", Severity: "warning", ExpiresAt: now.Add(time.Hour)},
		{ID: "expired", KeyPrefix: "b.", Severity: "critical", ExpiresAt: now.Add(-time.Hour)},
	}
	active := policy.ActiveRules(cfg, now)
	if len(active) != 1 {
		t.Fatalf("expected 1 active rule, got %d", len(active))
	}
	if active[0].ID != "active" {
		t.Errorf("expected rule 'active', got %q", active[0].ID)
	}
}
