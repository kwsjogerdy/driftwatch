package metrics

import (
	"fmt"
	"strings"
)

// ThresholdConfig holds warn/crit thresholds for a named metric field.
type ThresholdConfig struct {
	Field    string  `json:"field"`
	Warn     float64 `json:"warn"`
	Critical float64 `json:"critical"`
}

// ValidateThresholdConfig returns an error if the config is malformed.
func ValidateThresholdConfig(cfg ThresholdConfig) error {
	if strings.TrimSpace(cfg.Field) == "" {
		return fmt.Errorf("threshold field must not be empty")
	}
	if cfg.Warn < 0 {
		return fmt.Errorf("warn threshold must be >= 0, got %.2f", cfg.Warn)
	}
	if cfg.Critical < 0 {
		return fmt.Errorf("critical threshold must be >= 0, got %.2f", cfg.Critical)
	}
	if cfg.Critical < cfg.Warn {
		return fmt.Errorf("critical threshold (%.2f) must be >= warn threshold (%.2f)", cfg.Critical, cfg.Warn)
	}
	return nil
}

// BuildThresholds validates a slice of configs and returns a map keyed by
// field name for fast lookup, or the first validation error encountered.
func BuildThresholds(cfgs []ThresholdConfig) (map[string]ThresholdConfig, error) {
	out := make(map[string]ThresholdConfig, len(cfgs))
	for _, c := range cfgs {
		if err := ValidateThresholdConfig(c); err != nil {
			return nil, fmt.Errorf("invalid threshold config for field %q: %w", c.Field, err)
		}
		out[strings.ToLower(c.Field)] = c
	}
	return out, nil
}
