package metrics

import (
	"errors"
	"fmt"
)

// ThresholdConfig holds raw configuration for metric thresholds.
type ThresholdConfig struct {
	Field  string  `json:"field"`
	WarnAt float64 `json:"warn_at"`
	CritAt float64 `json:"crit_at"`
}

// ValidateThresholdConfig checks that a ThresholdConfig is well-formed.
func ValidateThresholdConfig(c ThresholdConfig) error {
	if c.Field == "" {
		return errors.New("threshold field must not be empty")
	}
	allowed := map[string]bool{
		"total_runs":    true,
		"drift_count":   true,
		"error_count":   true,
		"drift_percent": true,
	}
	if !allowed[c.Field] {
		return fmt.Errorf("unknown threshold field %q", c.Field)
	}
	if c.WarnAt < 0 {
		return fmt.Errorf("warn_at must be >= 0, got %.2f", c.WarnAt)
	}
	if c.CritAt < 0 {
		return fmt.Errorf("crit_at must be >= 0, got %.2f", c.CritAt)
	}
	if c.CritAt > 0 && c.WarnAt > 0 && c.CritAt <= c.WarnAt {
		return fmt.Errorf("crit_at (%.2f) must be greater than warn_at (%.2f)", c.CritAt, c.WarnAt)
	}
	return nil
}

// BuildThresholds converts a slice of ThresholdConfig into Threshold values,
// validating each entry. Returns an error on the first invalid entry.
func BuildThresholds(cfgs []ThresholdConfig) ([]Threshold, error) {
	out := make([]Threshold, 0, len(cfgs))
	for i, c := range cfgs {
		if err := ValidateThresholdConfig(c); err != nil {
			return nil, fmt.Errorf("threshold[%d]: %w", i, err)
		}
		out = append(out, Threshold{Field: c.Field, WarnAt: c.WarnAt, CritAt: c.CritAt})
	}
	return out, nil
}
