package history

import (
	"errors"
	"fmt"
	"os"
)

// Config holds configuration for the history store.
type Config struct {
	// Dir is the directory where history entries are persisted.
	Dir string `json:"dir"`
	// Enabled controls whether history recording is active.
	Enabled bool `json:"enabled"`
	// MaxEntries is the soft cap on stored entries (0 = unlimited).
	MaxEntries int `json:"max_entries"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Dir:        ".driftwatch/history",
		Enabled:    true,
		MaxEntries: 100,
	}
}

// Validate checks that the Config is usable.
func (c Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Dir == "" {
		return errors.New("history: dir must not be empty when enabled")
	}
	if c.MaxEntries < 0 {
		return fmt.Errorf("history: max_entries must be >= 0, got %d", c.MaxEntries)
	}
	return nil
}

// EnsureDir creates the history directory if it does not exist.
func (c Config) EnsureDir() error {
	if !c.Enabled {
		return nil
	}
	if err := os.MkdirAll(c.Dir, 0o755); err != nil {
		return fmt.Errorf("history: ensure dir %q: %w", c.Dir, err)
	}
	return nil
}
