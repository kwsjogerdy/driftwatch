package suppress

import (
	"errors"
	"time"
)

// Config holds configuration for the suppression filter.
type Config struct {
	// Enabled controls whether suppression rules are applied.
	Enabled bool `json:"enabled"`

	// Rules is the list of suppression rules to apply.
	Rules []Rule `json:"rules"`

	// MaxRules is the maximum number of rules allowed (0 = unlimited).
	MaxRules int `json:"max_rules"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled:  true,
		Rules:    []Rule{},
		MaxRules: 100,
	}
}

// Validate checks the Config for logical errors.
func Validate(cfg Config) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.MaxRules < 0 {
		return errors.New("max_rules must be non-negative")
	}
	if cfg.MaxRules > 0 && len(cfg.Rules) > cfg.MaxRules {
		return errors.New("number of rules exceeds max_rules limit")
	}
	now := time.Now()
	seen := make(map[string]struct{}, len(cfg.Rules))
	for i, r := range cfg.Rules {
		if err := r.Validate(); err != nil {
			return err
		}
		if _, dup := seen[r.ID]; dup {
			return errors.New("duplicate rule id: " + r.ID)
		}
		seen[r.ID] = struct{}{}
		_ = i
		_ = now
	}
	return nil
}

// ActiveRules returns only the non-expired rules from cfg.
func ActiveRules(cfg Config, now time.Time) []Rule {
	out := make([]Rule, 0, len(cfg.Rules))
	for _, r := range cfg.Rules {
		if r.Expiry.IsZero() || r.Expiry.After(now) {
			out = append(out, r)
		}
	}
	return out
}
