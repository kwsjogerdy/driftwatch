package config

import "time"

// ApplyDefaults fills in zero-value fields with sensible defaults.
// It is safe to call multiple times; only unset fields will be modified.
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

// WithDefaults returns a new AppConfig based on src with defaults applied
// to any zero-value fields. The original value is not modified.
func WithDefaults(src AppConfig) AppConfig {
	ApplyDefaults(&src)
	return src
}
