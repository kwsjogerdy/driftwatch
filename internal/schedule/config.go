package schedule

import (
	"errors"
	"time"
)

// DefaultInterval is the fallback watch interval when none is specified.
const DefaultInterval = 5 * time.Minute

// ParseInterval parses a duration string and returns a validated interval.
// Falls back to DefaultInterval if the input is empty.
func ParseInterval(raw string) (time.Duration, error) {
	if raw == "" {
		return DefaultInterval, nil
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, errors.New("invalid interval: " + err.Error())
	}
	if d <= 0 {
		return 0, errors.New("interval must be positive")
	}
	return d, nil
}

// ValidateConfig checks that required WatchConfig fields are populated.
func ValidateConfig(cfg WatchConfig) error {
	if cfg.SourceEnv == "" {
		return errors.New("source environment must not be empty")
	}
	if cfg.TargetEnv == "" {
		return errors.New("target environment must not be empty")
	}
	if cfg.StateFile == "" {
		return errors.New("state file path must not be empty")
	}
	if cfg.SourceEnv == cfg.TargetEnv {
		return errors.New("source and target environments must differ")
	}
	if cfg.Interval <= 0 {
		return errors.New("interval must be positive")
	}
	return nil
}
