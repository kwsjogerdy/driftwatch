package diff

import (
	"fmt"
)

// Config controls how drift differences are classified and filtered.
type Config struct {
	// DefaultSeverity is applied to mismatched keys not listed in CriticalKeys.
	DefaultSeverity string

	// CriticalKeys are key names that should always be classified as critical drift.
	CriticalKeys []string

	// IgnoreKeys are key names that should be excluded from drift comparison.
	IgnoreKeys []string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		DefaultSeverity: "warning",
		CriticalKeys:    []string{},
		IgnoreKeys:      []string{},
	}
}

// Validate checks that the Config fields are consistent and valid.
func Validate(cfg Config) error {
	valid := map[string]bool{"warning": true, "critical": true}
	if !valid[cfg.DefaultSeverity] {
		return fmt.Errorf("invalid default_severity %q: must be 'warning' or 'critical'", cfg.DefaultSeverity)
	}

	ignoreSet := make(map[string]bool, len(cfg.IgnoreKeys))
	for _, k := range cfg.IgnoreKeys {
		ignoreSet[k] = true
	}

	for _, k := range cfg.CriticalKeys {
		if ignoreSet[k] {
			return fmt.Errorf("key %q appears in both critical_keys and ignore_keys", k)
		}
	}

	return nil
}

// ApplyIgnore returns a copy of values with any keys listed in cfg.IgnoreKeys removed.
func ApplyIgnore(cfg Config, values map[string]string) map[string]string {
	if len(cfg.IgnoreKeys) == 0 {
		copy := make(map[string]string, len(values))
		for k, v := range values {
			copy[k] = v
		}
		return copy
	}

	ignoreSet := make(map[string]bool, len(cfg.IgnoreKeys))
	for _, k := range cfg.IgnoreKeys {
		ignoreSet[k] = true
	}

	result := make(map[string]string)
	for k, v := range values {
		if !ignoreSet[k] {
			result[k] = v
		}
	}
	return result
}
