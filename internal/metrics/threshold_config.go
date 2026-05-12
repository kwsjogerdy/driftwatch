package metrics

import (
	"errors"
	"fmt"
)

// ThresholdConfig holds raw configuration for drift metric thresholds.
type ThresholdConfig struct {
	Field    string  `json:"field"`
	WarnAt   float64 `json:"warn_at"`
	CritAt   float64 `json:"crit_at"`
}

// ValidateThresholdConfig checks that a ThresholdConfig is well-formed.
func ValidateThresholdConfig(c ThresholdConfig) error {
	if c.Field == "" {
		return errors.New("threshold field must not be empty")
	}
	known := map[string]bool{
		"drift_count":   true,
		"missing_keys":  true,
		"extra_keys":    true,
		"changed_values": true,
	}
	if !known[c.Field] {
		return fmt.Errorf("unknown threshold field %q", c.Field)
	}
	if c.WarnAt < 0 {
		return fmt.Errorf("warn_at must be >= 0, got %.2f", c.WarnAt)
	}
	if c.CritAt < 0 {
		return fmt.Errorf("crit_at must be >= 0, got %.2f", c.CritAt)
	}
	if c.CritAt > 0 && c.WarnAt > 0 && c.CritAt < c.WarnAt {
		return fmt.Errorf("crit_at (%.2f) must be >= warn_at (%.2f)", c.CritAt, c.WarnAt)
	}
	return nil
}

// BuildThresholds validates and converts a slice of ThresholdConfig into Threshold values.
func BuildThresholds(cfgs []ThresholdConfig) ([]Threshold, error) {
	out := make([]Threshold, 0, len(cfgs))
	for i, c := range cfgs {
		if err := ValidateThresholdConfig(c); err != nil {
			return nil, fmt.Errorf("threshold[%d]: %w", i, err)
		}
		out = append(out, Threshold{
			Field:  c.Field,
			WarnAt: c.WarnAt,
			CritAt: c.CritAt,
		})
	}
	return out, nil
}
