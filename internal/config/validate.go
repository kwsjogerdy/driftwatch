package config

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrMissingStateFile  = errors.New("config: state_file is required")
	ErrMissingEnvs       = errors.New("config: at least two environments are required")
	ErrInvalidInterval   = errors.New("config: interval must be positive")
	ErrInvalidAlertLevel = errors.New("config: alert_level must be 'warning' or 'critical'")
	ErrInvalidFormat     = errors.New("config: output format must be 'text' or 'json'")
)

var validAlertLevels = map[string]bool{"warning": true, "critical": true}
var validFormats = map[string]bool{"text": true, "json": true}

// Validate checks that the AppConfig has all required fields with valid values.
func Validate(cfg *AppConfig) error {
	if cfg.StateFile == "" {
		return ErrMissingStateFile
	}

	if len(cfg.Environments) < 2 {
		return ErrMissingEnvs
	}

	if cfg.Interval != 0 && cfg.Interval < time.Second {
		return fmt.Errorf("%w: got %s", ErrInvalidInterval, cfg.Interval)
	}

	if cfg.AlertLevel != "" && !validAlertLevels[cfg.AlertLevel] {
		return fmt.Errorf("%w: got %q", ErrInvalidAlertLevel, cfg.AlertLevel)
	}

	if cfg.Output.Format != "" && !validFormats[cfg.Output.Format] {
		return fmt.Errorf("%w: got %q", ErrInvalidFormat, cfg.Output.Format)
	}

	if cfg.Output.Destination == "file" && cfg.Output.FilePath == "" {
		return errors.New("config: file_path is required when destination is 'file'")
	}

	return nil
}
