package schedule_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/schedule"
)

func TestParseInterval_Valid(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{"30s", 30 * time.Second},
		{"5m", 5 * time.Minute},
		{"2h", 2 * time.Hour},
		{"1m30s", 90 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := schedule.ParseInterval(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseInterval_Invalid(t *testing.T) {
	inputs := []string{"", "abc", "10", "-5s"}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			_, err := schedule.ParseInterval(input)
			if err == nil {
				t.Errorf("expected error for input %q, got nil", input)
			}
		})
	}
}

func TestValidateConfig_Valid(t *testing.T) {
	cfg := schedule.Config{
		Interval:    30 * time.Second,
		SourceEnv:   "production",
		TargetEnv:   "staging",
	}
	if err := schedule.ValidateConfig(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateConfig_ZeroInterval(t *testing.T) {
	cfg := schedule.Config{
		Interval:  0,
		SourceEnv: "production",
		TargetEnv: "staging",
	}
	if err := schedule.ValidateConfig(cfg); err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestValidateConfig_MissingEnvs(t *testing.T) {
	tests := []schedule.Config{
		{Interval: time.Minute, SourceEnv: "", TargetEnv: "staging"},
		{Interval: time.Minute, SourceEnv: "production", TargetEnv: ""},
	}
	for _, cfg := range tests {
		if err := schedule.ValidateConfig(cfg); err == nil {
			t.Errorf("expected error for config %+v", cfg)
		}
	}
}

func TestValidateConfig_SameEnvs(t *testing.T) {
	cfg := schedule.Config{
		Interval:  time.Minute,
		SourceEnv: "production",
		TargetEnv: "production",
	}
	if err := schedule.ValidateConfig(cfg); err == nil {
		t.Fatal("expected error when SourceEnv and TargetEnv are the same")
	}
}
