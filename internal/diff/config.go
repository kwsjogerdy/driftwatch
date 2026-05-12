package diff

import (
	"errors"
	"strings"
)

// Config controls the behaviour of the diff Builder.
type Config struct {
	// CriticalKeys lists keys whose drift should be treated as critical.
	CriticalKeys []string `json:"critical_keys"`
	// IgnoreKeys lists keys to exclude from comparison entirely.
	IgnoreKeys []string `json:"ignore_keys"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		CriticalKeys: []string{},
		IgnoreKeys:   []string{},
	}
}

// Validate checks that the Config is consistent.
func Validate(c Config) error {
	critSet := make(map[string]bool, len(c.CriticalKeys))
	for _, k := range c.CriticalKeys {
		if strings.TrimSpace(k) == "" {
			return errors.New("diff: critical_keys contains a blank entry")
		}
		critSet[k] = true
	}
	for _, k := range c.IgnoreKeys {
		if strings.TrimSpace(k) == "" {
			return errors.New("diff: ignore_keys contains a blank entry")
		}
		if critSet[k] {
			return errors.New("diff: key \"" + k + "\" appears in both critical_keys and ignore_keys")
		}
	}
	return nil
}

// ApplyIgnore returns a copy of m with all ignore_keys removed.
func ApplyIgnore(m map[string]interface{}, ignore []string) map[string]interface{} {
	if len(ignore) == 0 {
		return m
	}
	skip := make(map[string]bool, len(ignore))
	for _, k := range ignore {
		skip[k] = true
	}
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		if !skip[k] {
			out[k] = v
		}
	}
	return out
}
