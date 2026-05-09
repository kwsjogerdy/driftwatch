package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Format controls how metrics are rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Render writes a human-readable or JSON representation of the Summary to w.
func Render(w io.Writer, s Summary, f Format) error {
	switch f {
	case FormatJSON:
		return renderJSON(w, s)
	case FormatText:
		return renderText(w, s)
	default:
		return fmt.Errorf("metrics: unknown format %q", f)
	}
}

func renderText(w io.Writer, s Summary) error {
	lines := []string{
		fmt.Sprintf("total_runs:       %d", s.TotalRuns),
		fmt.Sprintf("drift_runs:       %d", s.DriftRuns),
		fmt.Sprintf("error_runs:       %d", s.ErrorRuns),
		fmt.Sprintf("total_drift_keys: %d", s.TotalDriftKeys),
		fmt.Sprintf("avg_duration_ms:  %d", s.AvgDuration.Milliseconds()),
	}
	_, err := fmt.Fprintln(w, strings.Join(lines, "\n"))
	return err
}

func renderJSON(w io.Writer, s Summary) error {
	payload := map[string]any{
		"total_runs":       s.TotalRuns,
		"drift_runs":       s.DriftRuns,
		"error_runs":       s.ErrorRuns,
		"total_drift_keys": s.TotalDriftKeys,
		"avg_duration_ms":  s.AvgDuration.Milliseconds(),
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
