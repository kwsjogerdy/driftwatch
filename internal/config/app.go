package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// AppConfig holds the top-level application configuration.
type AppConfig struct {
	StateFile   string        `json:"state_file"`
	Environments []string     `json:"environments"`
	Interval    time.Duration `json:"-"`
	IntervalRaw string        `json:"interval"`
	Output      OutputConfig  `json:"output"`
	AlertLevel  string        `json:"alert_level"`
}

// OutputConfig defines where and how reports are written.
type OutputConfig struct {
	Destination string `json:"destination"`
	Format      string `json:"format"`
	FilePath    string `json:"file_path,omitempty"`
}

// LoadAppConfig reads and parses an AppConfig from a JSON file at path.
func LoadAppConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file: %w", err)
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse json: %w", err)
	}

	if cfg.IntervalRaw != "" {
		d, err := time.ParseDuration(cfg.IntervalRaw)
		if err != nil {
			return nil, fmt.Errorf("config: invalid interval %q: %w", cfg.IntervalRaw, err)
		}
		cfg.Interval = d
	}

	return &cfg, nil
}
