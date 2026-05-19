package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// LoadConfig reads a policy configuration from the given file path.
// It applies defaults, validates the result, and returns the config.
func LoadConfig(path string) (Config, error) {
	return LoadConfigWithTime(path, time.Now())
}

// LoadConfigWithTime is like LoadConfig but accepts an explicit reference
// time, making it easier to test expiry behaviour.
func LoadConfigWithTime(path string, now time.Time) (Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, fmt.Errorf("policy config not found: %s", path)
		}
		return cfg, fmt.Errorf("reading policy config: %w", err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing policy config: %w", err)
	}

	if err := Validate(cfg); err != nil {
		return cfg, fmt.Errorf("invalid policy config: %w", err)
	}

	return cfg, nil
}
