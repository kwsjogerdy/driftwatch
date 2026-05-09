package baseline

import (
	"errors"
	"fmt"
	"time"
)

// Config holds configuration for baseline snapshot behaviour.
type Config struct {
	// Dir is the directory where baseline snapshots are persisted.
	Dir string `json:"dir"`

	// MaxAge is the maximum age of a baseline before it is considered stale.
	// A zero value means baselines never expire.
	MaxAge time.Duration `json:"max_age"`

	// AutoUpdate controls whether the baseline is automatically updated
	// after a successful drift-free comparison.
	AutoUpdate bool `json:"auto_update"`

	// Environments lists the environment names that baseline tracking is
	// enabled for. An empty slice means all environments are tracked.
	Environments []string `json:"environments"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Dir:          ".driftwatch/baselines",
		MaxAge:       7 * 24 * time.Hour, // one week
		AutoUpdate:   false,
		Environments: []string{},
	}
}

// Validate checks that the Config fields are consistent and usable.
// It returns a descriptive error for the first problem found.
func Validate(c Config) error {
	if c.Dir == "" {
		return errors.New("baseline: dir must not be empty")
	}

	if c.MaxAge < 0 {
		return fmt.Errorf("baseline: max_age must be non-negative, got %s", c.MaxAge)
	}

	seen := make(map[string]struct{}, len(c.Environments))
	for i, env := range c.Environments {
		if env == "" {
			return fmt.Errorf("baseline: environments[%d] must not be blank", i)
		}
		if _, dup := seen[env]; dup {
			return fmt.Errorf("baseline: duplicate environment %q in environments list", env)
		}
		seen[env] = struct{}{}
	}

	return nil
}

// IsStale reports whether a baseline recorded at recordedAt should be
// considered stale given the Config's MaxAge. A zero MaxAge never expires.
func (c Config) IsStale(recordedAt time.Time) bool {
	if c.MaxAge == 0 {
		return false
	}
	return time.Since(recordedAt) > c.MaxAge
}

// TracksEnvironment reports whether the given environment name is covered by
// this Config. When the Environments list is empty every environment is tracked.
func (c Config) TracksEnvironment(env string) bool {
	if len(c.Environments) == 0 {
		return true
	}
	for _, e := range c.Environments {
		if e == env {
			return true
		}
	}
	return false
}
