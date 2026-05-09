package metrics

import (
	"testing"
)

func baseSummary() Summary {
	return Summary{
		TotalRuns:    10,
		DriftCount:   3,
		ErrorCount:   1,
		DriftPercent: 30.0,
	}
}

func TestEvaluate_NoBreaches(t *testing.T) {
	e := NewEvaluator([]Threshold{
		{Field: "drift_count", WarnAt: 5, CritAt: 10},
	})
	breaches := e.Evaluate(baseSummary())
	if len(breaches) != 0 {
		t.Fatalf("expected 0 breaches, got %d", len(breaches))
	}
}

func TestEvaluate_WarnBreach(t *testing.T) {
	e := NewEvaluator([]Threshold{
		{Field: "drift_count", WarnAt: 2, CritAt: 8},
	})
	breaches := e.Evaluate(baseSummary())
	if len(breaches) != 1 {
		t.Fatalf("expected 1 breach, got %d", len(breaches))
	}
	if breaches[0].Level != ThresholdWarn {
		t.Errorf("expected warn level, got %s", breaches[0].Level)
	}
}

func TestEvaluate_CritBreach(t *testing.T) {
	e := NewEvaluator([]Threshold{
		{Field: "drift_percent", WarnAt: 10, CritAt: 25},
	})
	breaches := e.Evaluate(baseSummary())
	if len(breaches) != 1 {
		t.Fatalf("expected 1 breach, got %d", len(breaches))
	}
	if breaches[0].Level != ThresholdCritical {
		t.Errorf("expected critical level, got %s", breaches[0].Level)
	}
}

func TestEvaluate_UnknownField_Skipped(t *testing.T) {
	e := NewEvaluator([]Threshold{
		{Field: "nonexistent", WarnAt: 1, CritAt: 2},
	})
	breaches := e.Evaluate(baseSummary())
	if len(breaches) != 0 {
		t.Fatalf("expected 0 breaches for unknown field, got %d", len(breaches))
	}
}

func TestBuildThresholds_Valid(t *testing.T) {
	cfgs := []ThresholdConfig{
		{Field: "drift_count", WarnAt: 3, CritAt: 7},
		{Field: "error_count", WarnAt: 1, CritAt: 5},
	}
	ts, err := BuildThresholds(cfgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ts) != 2 {
		t.Errorf("expected 2 thresholds, got %d", len(ts))
	}
}

func TestBuildThresholds_InvalidField(t *testing.T) {
	cfgs := []ThresholdConfig{
		{Field: "unknown_field", WarnAt: 1},
	}
	_, err := BuildThresholds(cfgs)
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestBuildThresholds_CritBelowWarn(t *testing.T) {
	cfgs := []ThresholdConfig{
		{Field: "drift_count", WarnAt: 5, CritAt: 3},
	}
	_, err := BuildThresholds(cfgs)
	if err == nil {
		t.Fatal("expected error when crit_at <= warn_at")
	}
}

func TestValidateThresholdConfig_EmptyField(t *testing.T) {
	err := ValidateThresholdConfig(ThresholdConfig{Field: "", WarnAt: 1})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}
