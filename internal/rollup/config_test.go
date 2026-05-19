package rollup_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/rollup"
)

func TestDefaultConfig_Values(t *testing.T) {
	cfg := rollup.DefaultConfig()
	if cfg.Period == "" {
		t.Fatal("expected non-empty default period")
	}
	if cfg.TopKeysLimit <= 0 {
		t.Fatalf("expected positive TopKeysLimit, got %d", cfg.TopKeysLimit)
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := rollup.DefaultConfig()
	if err := rollup.Validate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_InvalidPeriod(t *testing.T) {
	cfg := rollup.DefaultConfig()
	cfg.Period = "fortnight"
	if err := rollup.Validate(cfg); err == nil {
		t.Fatal("expected error for invalid period")
	}
}

func TestValidate_ZeroTopKeysLimit(t *testing.T) {
	cfg := rollup.DefaultConfig()
	cfg.TopKeysLimit = 0
	if err := rollup.Validate(cfg); err == nil {
		t.Fatal("expected error for zero TopKeysLimit")
	}
}

func TestValidate_NegativeTopKeysLimit(t *testing.T) {
	cfg := rollup.DefaultConfig()
	cfg.TopKeysLimit = -1
	if err := rollup.Validate(cfg); err == nil {
		t.Fatal("expected error for negative TopKeysLimit")
	}
}

func TestPeriodFor_Daily(t *testing.T) {
	cfg := rollup.DefaultConfig()
	cfg.Period = "daily"
	now := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)
	start, end := rollup.PeriodFor(cfg, now)
	if !start.Equal(time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("unexpected start: %v", start)
	}
	if !end.Equal(time.Date(2024, 6, 16, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("unexpected end: %v", end)
	}
}

func TestPeriodFor_Weekly(t *testing.T) {
	cfg := rollup.DefaultConfig()
	cfg.Period = "weekly"
	// Wednesday 2024-06-19
	now := time.Date(2024, 6, 19, 10, 0, 0, 0, time.UTC)
	start, end := rollup.PeriodFor(cfg, now)
	if end.Sub(start) != 7*24*time.Hour {
		t.Errorf("expected 7-day window, got %v", end.Sub(start))
	}
	if start.After(now) {
		t.Errorf("start %v is after now %v", start, now)
	}
}

func TestPeriodFor_Monthly(t *testing.T) {
	cfg := rollup.DefaultConfig()
	cfg.Period = "monthly"
	now := time.Date(2024, 6, 15, 8, 0, 0, 0, time.UTC)
	start, end := rollup.PeriodFor(cfg, now)
	if !start.Equal(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("unexpected start: %v", start)
	}
	if !end.Equal(time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("unexpected end: %v", end)
	}
}
