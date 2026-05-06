package output

import (
	"fmt"
	"strings"
)

// Config holds output configuration parsed from flags or config files.
type Config struct {
	Destination Destination
	FilePath    string
	Format      string // "text" or "json"
}

// ParseDestination converts a string to a Destination, returning an error
// for unrecognised values.
func ParseDestination(s string) (Destination, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "stdout":
		return DestStdout, nil
	case "stderr":
		return DestStderr, nil
	case "file":
		return DestFile, nil
	default:
		return "", fmt.Errorf("unsupported output destination %q; choose stdout, stderr, or file", s)
	}
}

// ValidateConfig checks that the Config is self-consistent.
func ValidateConfig(c Config) error {
	if c.Destination == DestFile && c.FilePath == "" {
		return fmt.Errorf("output destination is 'file' but no file path was provided")
	}
	switch strings.ToLower(c.Format) {
	case "text", "json":
		// valid
	case "":
		return fmt.Errorf("output format must be specified (text or json)")
	default:
		return fmt.Errorf("unsupported output format %q; choose text or json", c.Format)
	}
	return nil
}
