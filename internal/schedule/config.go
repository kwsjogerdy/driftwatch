package schedule

import (
	"errors"
	"fmt"
	"time"
)

// Config holds scheduling configuration for the drift watcher.
type Config struct {
	Interval  time.Duration
	SourceEnv string
	TargetEnv string
}

// ParseInterval parses a duration string into a time.Duration.
// Returns an error if the string is empty, invalid, or non-positive.
func ParseInterval(s string) (time.Duration, error) {
	if s == "" {
		return 0, errors.New("interval must not be empty")
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("invalid interval %q: %w", s, err)
	}
	if d <= 0 {
		return 0, fmt.Errorf("interval must be positive, got %v", d)
	}
	return d, nil
}

// ValidateConfig checks that a Config has all required fields set correctly.
func ValidateConfig(cfg Config) error {
	if cfg.Interval <= 0 {
		return fmt.Errorf("interval must be positive, got %v", cfg.Interval)
	}
	if cfg.SourceEnv == "" {
		return errors.New("source environment must not be empty")
	}
	if cfg.TargetEnv == "" {
		return errors.New("target environment must not be empty")
	}
	return nil
}
