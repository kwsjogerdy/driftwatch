package suppress

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// LoadConfig reads a suppression Config from a JSON file at path.
// If path is empty, DefaultConfig is returned.
func LoadConfig(path string) (Config, error) {
	if path == "" {
		return DefaultConfig(), nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, fmt.Errorf("suppress config not found: %s", path)
		}
		return Config{}, fmt.Errorf("reading suppress config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing suppress config: %w", err)
	}
	if err := Validate(cfg); err != nil {
		return Config{}, fmt.Errorf("invalid suppress config: %w", err)
	}
	return cfg, nil
}

// LoadConfigWithTime is like LoadConfig but uses the provided time for expiry
// filtering, returning only active rules.
func LoadConfigWithTime(path string, now time.Time) (Config, []Rule, error) {
	cfg, err := LoadConfig(path)
	if err != nil {
		return Config{}, nil, err
	}
	return cfg, ActiveRules(cfg, now), nil
}
