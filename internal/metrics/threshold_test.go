package metrics

import (
	"testing"
)

func baseSummary() Summary {
	return Summary{
		TotalRuns:    10,
		TotalMissing: 2,
		TotalExtra:   1,
		TotalChanged: 3,
	}
}

func TestEvaluate_NoBreaches(t *testing.T) {
	e := NewEvaluator([]Threshold{
		{Field: "drift_count", WarnAt: 20, CritAt: 50},
	})
	breaches := e.Evaluate(baseSummary())
	if len(breaches) != 0 {
		t.Fatalf("expected no breaches, got %d", len(breaches))
	}
}

func TestEvaluate_WarnBreach(t *testing.T) {
	e := NewEvaluator([]Threshold{
		{Field: "changed_values", WarnAt: 2, CritAt: 10},
	})
	breaches := e.Evaluate(baseSummary())
	if len(breaches) != 1 {
		t.Fatalf("expected 1 breach, got %d", len(breaches))
	}
	if breaches[0].Severity != SeverityWarn {
		t.Errorf("expected warn, got %s", breaches[0].Severity)
	}
}

func TestEvaluate_CritBreach(t *testing.T) {
	e := NewEvaluator([]Threshold{
		{Field: "missing_keys", WarnAt: 1, CritAt: 2},
	})
	breaches := e.Evaluate(baseSummary())
	if len(breaches) != 1 {
		t.Fatalf("expected 1 breach, got %d", len(breaches))
	}
	if breaches[0].Severity != SeverityCrit {
		t.Errorf("expected crit, got %s", breaches[0].Severity)
	}
}

func TestEvaluate_UnknownField_Skipped(t *testing.T) {
	e := NewEvaluator([]Threshold{
		{Field: "nonexistent", WarnAt: 1, CritAt: 2},
	})
	breaches := e.Evaluate(baseSummary())
	if len(breaches) != 0 {
		t.Fatalf("expected no breaches for unknown field, got %d", len(breaches))
	}
}

func TestValidateThresholdConfig_Valid(t *testing.T) {
	cfg := ThresholdConfig{Field: "drift_count", WarnAt: 5, CritAt: 10}
	if err := ValidateThresholdConfig(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateThresholdConfig_EmptyField(t *testing.T) {
	cfg := ThresholdConfig{Field: "", WarnAt: 5, CritAt: 10}
	if err := ValidateThresholdConfig(cfg); err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestValidateThresholdConfig_CritLessThanWarn(t *testing.T) {
	cfg := ThresholdConfig{Field: "extra_keys", WarnAt: 10, CritAt: 5}
	if err := ValidateThresholdConfig(cfg); err == nil {
		t.Fatal("expected error when crit_at < warn_at")
	}
}

func TestBuildThresholds_ReturnsAll(t *testing.T) {
	cfgs := []ThresholdConfig{
		{Field: "drift_count", WarnAt: 3, CritAt: 6},
		{Field: "missing_keys", WarnAt: 1, CritAt: 5},
	}
	thresholds, err := BuildThresholds(cfgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(thresholds) != 2 {
		t.Errorf("expected 2 thresholds, got %d", len(thresholds))
	}
}

func TestBuildThresholds_InvalidConfig_ReturnsError(t *testing.T) {
	cfgs := []ThresholdConfig{
		{Field: "unknown_field", WarnAt: 1, CritAt: 5},
	}
	_, err := BuildThresholds(cfgs)
	if err == nil {
		t.Fatal("expected error for unknown field in BuildThresholds")
	}
}
