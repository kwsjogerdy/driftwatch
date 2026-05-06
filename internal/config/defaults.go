package config

import "time"

// ApplyDefaults fills in zero-value fields with sensible defaults.
func ApplyDefaults(cfg *AppConfig) {
	if cfg.Interval == 0 {
		cfg.Interval = 5 * time.Minute
	}

	if cfg.AlertLevel == "" {
		cfg.AlertLevel = "warning"
	}

	if cfg.Output.Destination == "" {
		cfg.Output.Destination = "stdout"
	}

	if cfg.Output.Format == "" {
		cfg.Output.Format = "text"
	}
}
