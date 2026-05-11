package metrics

import (
	"testing"
)

func baseSummary() Summary {
	return Summary{
		TotalRuns:    10,
		DriftCount:   3,
		CleanCount:   7,
		DriftRate:    0.30,
		AvgDiffCount: 2.5,
	}
}

func TestEvaluate_NoBreaches(t *testing.T) {
	thresholds := map[string]ThresholdConfig{
		"drift_rate": {Field: "drift_rate", Warn: 0.5, Critical: 0.8},
	}
	e := NewEvaluator(thresholds)
	breaches := e.Evaluate(baseSummary())
	if len(breaches) != 0 {
		t.Fatalf("expected no breaches, got %d", len(breaches))
	}
}

func TestEvaluate_WarnBreach(t *testing.T) {
	thresholds := map[string]ThresholdConfig{
		"drift_rate": {Field: "drift_rate", Warn: 0.2, Critical: 0.8},
	}
	e := NewEvaluator(thresholds)
	breaches := e.Evaluate(baseSummary())
	if len(breaches) != 1 {
		t.Fatalf("expected 1 breach, got %d", len(breaches))
	}
	if breaches[0].Level != BreachWarn {
		t.Errorf("expected warn, got %s", breaches[0].Level)
	}
}

func TestEvaluate_CritBreach(t *testing.T) {
	thresholds := map[string]ThresholdConfig{
		"drift_rate": {Field: "drift_rate", Warn: 0.1, Critical: 0.2},
	}
	e := NewEvaluator(thresholds)
	breaches := e.Evaluate(baseSummary())
	if len(breaches) != 1 {
		t.Fatalf("expected 1 breach, got %d", len(breaches))
	}
	if breaches[0].Level != BreachCritical {
		t.Errorf("expected critical, got %s", breaches[0].Level)
	}
}

func TestEvaluate_UnknownField_Skipped(t *testing.T) {
	thresholds := map[string]ThresholdConfig{
		"nonexistent_field": {Field: "nonexistent_field", Warn: 1, Critical: 2},
	}
	e := NewEvaluator(thresholds)
	breaches := e.Evaluate(baseSummary())
	if len(breaches) != 0 {
		t.Errorf("expected no breaches for unknown field, got %d", len(breaches))
	}
}

func TestValidateThresholdConfig_Valid(t *testing.T) {
	cfg := ThresholdConfig{Field: "drift_rate", Warn: 0.3, Critical: 0.7}
	if err := ValidateThresholdConfig(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateThresholdConfig_EmptyField(t *testing.T) {
	cfg := ThresholdConfig{Field: "", Warn: 0.3, Critical: 0.7}
	if err := ValidateThresholdConfig(cfg); err == nil {
		t.Error("expected error for empty field")
	}
}

func TestValidateThresholdConfig_CritLessThanWarn(t *testing.T) {
	cfg := ThresholdConfig{Field: "drift_rate", Warn: 0.8, Critical: 0.3}
	if err := ValidateThresholdConfig(cfg); err == nil {
		t.Error("expected error when critical < warn")
	}
}

func TestBuildThresholds_Valid(t *testing.T) {
	cfgs := []ThresholdConfig{
		{Field: "drift_rate", Warn: 0.3, Critical: 0.7},
		{Field: "drift_count", Warn: 5, Critical: 10},
	}
	m, err := BuildThresholds(cfgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 2 {
		t.Errorf("expected 2 entries, got %d", len(m))
	}
}

func TestBuildThresholds_InvalidEntry_ReturnsError(t *testing.T) {
	cfgs := []ThresholdConfig{
		{Field: "", Warn: 0.3, Critical: 0.7},
	}
	_, err := BuildThresholds(cfgs)
	if err == nil {
		t.Error("expected error for invalid threshold config")
	}
}
