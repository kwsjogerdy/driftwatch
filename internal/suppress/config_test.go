package suppress

import (
	"testing"
	"time"
)

func baseConfig() Config {
	return Config{
		Enabled:  true,
		MaxRules: 100,
		Rules: []Rule{
			{
				ID:        "rule-1",
				KeyPrefix: "db.",
				Expiry:    time.Now().Add(time.Hour),
			},
		},
	}
}

func TestDefaultConfig_Values(t *testing.T) {
	cfg := DefaultConfig()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if cfg.MaxRules != 100 {
		t.Errorf("expected MaxRules 100, got %d", cfg.MaxRules)
	}
	if len(cfg.Rules) != 0 {
		t.Errorf("expected empty Rules, got %d", len(cfg.Rules))
	}
}

func TestValidate_Valid(t *testing.T) {
	if err := Validate(baseConfig()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_DisabledSkipsChecks(t *testing.T) {
	cfg := Config{Enabled: false, MaxRules: -99}
	if err := Validate(cfg); err != nil {
		t.Fatalf("expected no error when disabled, got: %v", err)
	}
}

func TestValidate_NegativeMaxRules(t *testing.T) {
	cfg := baseConfig()
	cfg.MaxRules = -1
	if err := Validate(cfg); err == nil {
		t.Error("expected error for negative max_rules")
	}
}

func TestValidate_ExceedsMaxRules(t *testing.T) {
	cfg := baseConfig()
	cfg.MaxRules = 1
	cfg.Rules = append(cfg.Rules, Rule{
		ID:        "rule-2",
		KeyPrefix: "cache.",
		Expiry:    time.Now().Add(time.Hour),
	})
	if err := Validate(cfg); err == nil {
		t.Error("expected error when rules exceed max_rules")
	}
}

func TestValidate_DuplicateRuleID(t *testing.T) {
	cfg := baseConfig()
	cfg.Rules = append(cfg.Rules, Rule{
		ID:        "rule-1",
		KeyPrefix: "cache.",
		Expiry:    time.Now().Add(time.Hour),
	})
	if err := Validate(cfg); err == nil {
		t.Error("expected error for duplicate rule ID")
	}
}

func TestActiveRules_ExcludesExpired(t *testing.T) {
	now := time.Now()
	cfg := Config{
		Enabled: true,
		Rules: []Rule{
			{ID: "active", KeyPrefix: "db.", Expiry: now.Add(time.Hour)},
			{ID: "expired", KeyPrefix: "cache.", Expiry: now.Add(-time.Hour)},
			{ID: "no-expiry", KeyPrefix: "app."},
		},
	}
	active := ActiveRules(cfg, now)
	if len(active) != 2 {
		t.Fatalf("expected 2 active rules, got %d", len(active))
	}
	for _, r := range active {
		if r.ID == "expired" {
			t.Error("expired rule should not be in active list")
		}
	}
}
