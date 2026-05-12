package metrics

import (
	"testing"
)

func TestValidateThresholdConfig_Valid(t *testing.T) {
	cfg := ThresholdConfig{
		WarnDriftCount:  2,
		CritDriftCount:  5,
		WarnDriftRatio:  0.2,
		CritDriftRatio:  0.5,
		WarnRunDuration: 1.0,
		CritRunDuration: 5.0,
	}
	if err := ValidateThresholdConfig(cfg); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateThresholdConfig_NegativeCount(t *testing.T) {
	cfg := ThresholdConfig{WarnDriftCount: -1}
	if err := ValidateThresholdConfig(cfg); err == nil {
		t.Fatal("expected error for negative warn_drift_count")
	}
}

func TestValidateThresholdConfig_WarnNotLessThanCrit_Count(t *testing.T) {
	cfg := ThresholdConfig{WarnDriftCount: 5, CritDriftCount: 5}
	if err := ValidateThresholdConfig(cfg); err == nil {
		t.Fatal("expected error when warn_drift_count >= crit_drift_count")
	}
}

func TestValidateThresholdConfig_RatioOutOfRange(t *testing.T) {
	cfg := ThresholdConfig{WarnDriftRatio: 1.5}
	if err := ValidateThresholdConfig(cfg); err == nil {
		t.Fatal("expected error for ratio > 1")
	}
}

func TestValidateThresholdConfig_WarnNotLessThanCrit_Ratio(t *testing.T) {
	cfg := ThresholdConfig{WarnDriftRatio: 0.6, CritDriftRatio: 0.4}
	if err := ValidateThresholdConfig(cfg); err == nil {
		t.Fatal("expected error when warn_drift_ratio >= crit_drift_ratio")
	}
}

func TestValidateThresholdConfig_NegativeDuration(t *testing.T) {
	cfg := ThresholdConfig{WarnRunDuration: -0.5}
	if err := ValidateThresholdConfig(cfg); err == nil {
		t.Fatal("expected error for negative warn_run_duration_seconds")
	}
}

func TestBuildThresholds_Empty(t *testing.T) {
	thresholds := BuildThresholds(ThresholdConfig{})
	if len(thresholds) != 0 {
		t.Fatalf("expected 0 thresholds, got %d", len(thresholds))
	}
}

func TestBuildThresholds_FullConfig(t *testing.T) {
	cfg := ThresholdConfig{
		WarnDriftCount:  2,
		CritDriftCount:  5,
		WarnDriftRatio:  0.2,
		CritDriftRatio:  0.5,
		WarnRunDuration: 1.0,
		CritRunDuration: 5.0,
	}
	thresholds := BuildThresholds(cfg)
	if len(thresholds) != 6 {
		t.Fatalf("expected 6 thresholds, got %d", len(thresholds))
	}
}

func TestBuildThresholds_PartialConfig(t *testing.T) {
	cfg := ThresholdConfig{
		WarnDriftCount: 3,
		CritDriftCount: 8,
	}
	thresholds := BuildThresholds(cfg)
	if len(thresholds) != 2 {
		t.Fatalf("expected 2 thresholds, got %d", len(thresholds))
	}
	if thresholds[0].Field != "drift_count" || thresholds[0].Level != "warn" {
		t.Errorf("unexpected first threshold: %+v", thresholds[0])
	}
	if thresholds[1].Field != "drift_count" || thresholds[1].Level != "crit" {
		t.Errorf("unexpected second threshold: %+v", thresholds[1])
	}
}
