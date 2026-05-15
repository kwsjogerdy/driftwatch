package policy

import (
	"errors"
	"fmt"
)

// Config holds the policy configuration loaded from the app config.
type Config struct {
	Enabled bool   `json:"enabled"`
	Rules   []Rule `json:"rules"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled: true,
		Rules:   []Rule{},
	}
}

// Validate checks that all rules in the config are well-formed.
func Validate(c Config) error {
	if !c.Enabled {
		return nil
	}
	if len(c.Rules) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(c.Rules))
	for i, r := range c.Rules {
		if err := r.Validate(); err != nil {
			return fmt.Errorf("rule[%d]: %w", i, err)
		}
		if _, dup := seen[r.ID]; dup {
			return fmt.Errorf("rule[%d]: duplicate id %q", i, r.ID)
		}
		seen[r.ID] = struct{}{}
	}
	return nil
}

// ActiveRules returns only the non-expired rules from the config.
func ActiveRules(c Config, isExpired func(Rule) bool) ([]Rule, error) {
	if !c.Enabled {
		return nil, nil
	}
	if err := Validate(c); err != nil {
		return nil, errors.New("invalid policy config: " + err.Error())
	}
	out := make([]Rule, 0, len(c.Rules))
	for _, r := range c.Rules {
		if !isExpired(r) {
			out = append(out, r)
		}
	}
	return out, nil
}
