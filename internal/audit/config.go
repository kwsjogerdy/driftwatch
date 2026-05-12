package audit

import (
	"errors"
	"strings"
)

// Config holds configuration for the audit logger.
type Config struct {
	Enabled     bool   `json:"enabled"`
	Format      string `json:"format"`       // "text" or "json"
	Destination string `json:"destination"` // "stdout", "stderr", or a file path
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled:     true,
		Format:      "text",
		Destination: "stdout",
	}
}

// Validate checks that the Config fields are acceptable.
func Validate(c Config) error {
	if !c.Enabled {
		return nil
	}
	f := strings.ToLower(c.Format)
	if f != "text" && f != "json" {
		return errors.New("audit: format must be \"text\" or \"json\"")
	}
	d := strings.ToLower(c.Destination)
	if d == "" {
		return errors.New("audit: destination must not be empty")
	}
	if d != "stdout" && d != "stderr" && !strings.HasPrefix(c.Destination, "/") && !strings.Contains(c.Destination, "/") {
		return errors.New("audit: destination must be \"stdout\", \"stderr\", or a valid file path")
	}
	return nil
}

// NormaliseFormat returns the lowercase format, defaulting to "text".
func NormaliseFormat(c Config) string {
	f := strings.ToLower(c.Format)
	if f == "" {
		return "text"
	}
	return f
}
