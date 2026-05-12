package metrics

import (
	"fmt"
	"strings"
)

// ThresholdConfig holds user-supplied threshold configuration.
type ThresholdConfig struct {
	WarnDriftCount  int     `json:"warn_drift_count"`
	CritDriftCount  int     `json:"crit_drift_count"`
	WarnDriftRatio  float64 `json:"warn_drift_ratio"`
	CritDriftRatio  float64 `json:"crit_drift_ratio"`
	WarnRunDuration float64 `json:"warn_run_duration_seconds"`
	CritRunDuration float64 `json:"crit_run_duration_seconds"`
}

// ValidateThresholdConfig returns an error if the configuration is inconsistent.
func ValidateThresholdConfig(c ThresholdConfig) error {
	var errs []string

	if c.WarnDriftCount < 0 {
		errs = append(errs, "warn_drift_count must be >= 0")
	}
	if c.CritDriftCount < 0 {
		errs = append(errs, "crit_drift_count must be >= 0")
	}
	if c.WarnDriftCount > 0 && c.CritDriftCount > 0 && c.WarnDriftCount >= c.CritDriftCount {
		errs = append(errs, "warn_drift_count must be less than crit_drift_count")
	}
	if c.WarnDriftRatio < 0 || c.WarnDriftRatio > 1 {
		errs = append(errs, "warn_drift_ratio must be between 0 and 1")
	}
	if c.CritDriftRatio < 0 || c.CritDriftRatio > 1 {
		errs = append(errs, "crit_drift_ratio must be between 0 and 1")
	}
	if c.WarnDriftRatio > 0 && c.CritDriftRatio > 0 && c.WarnDriftRatio >= c.CritDriftRatio {
		errs = append(errs, "warn_drift_ratio must be less than crit_drift_ratio")
	}
	if c.WarnRunDuration < 0 {
		errs = append(errs, "warn_run_duration_seconds must be >= 0")
	}
	if c.CritRunDuration < 0 {
		errs = append(errs, "crit_run_duration_seconds must be >= 0")
	}
	if c.WarnRunDuration > 0 && c.CritRunDuration > 0 && c.WarnRunDuration >= c.CritRunDuration {
		errs = append(errs, "warn_run_duration_seconds must be less than crit_run_duration_seconds")
	}

	if len(errs) > 0 {
		return fmt.Errorf("invalid threshold config: %s", strings.Join(errs, "; "))
	}
	return nil
}

// BuildThresholds converts a ThresholdConfig into a slice of Threshold values
// ready for use with NewEvaluator.
func BuildThresholds(c ThresholdConfig) []Threshold {
	var out []Threshold

	if c.WarnDriftCount > 0 {
		out = append(out, Threshold{Field: "drift_count", Level: "warn", Value: float64(c.WarnDriftCount)})
	}
	if c.CritDriftCount > 0 {
		out = append(out, Threshold{Field: "drift_count", Level: "crit", Value: float64(c.CritDriftCount)})
	}
	if c.WarnDriftRatio > 0 {
		out = append(out, Threshold{Field: "drift_ratio", Level: "warn", Value: c.WarnDriftRatio})
	}
	if c.CritDriftRatio > 0 {
		out = append(out, Threshold{Field: "drift_ratio", Level: "crit", Value: c.CritDriftRatio})
	}
	if c.WarnRunDuration > 0 {
		out = append(out, Threshold{Field: "run_duration_seconds", Level: "warn", Value: c.WarnRunDuration})
	}
	if c.CritRunDuration > 0 {
		out = append(out, Threshold{Field: "run_duration_seconds", Level: "crit", Value: c.CritRunDuration})
	}

	return out
}
