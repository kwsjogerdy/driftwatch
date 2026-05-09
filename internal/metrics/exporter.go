package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ExportFormat defines the output format for exported metrics.
type ExportFormat string

const (
	ExportText ExportFormat = "text"
	ExportJSON ExportFormat = "json"
)

// ExportOptions controls how metrics are exported.
type ExportOptions struct {
	Format    ExportFormat
	Since     time.Time
	EnvFilter string
}

// Exporter writes collected metrics to an io.Writer.
type Exporter struct {
	collector *Collector
}

// NewExporter creates an Exporter backed by the given Collector.
func NewExporter(c *Collector) *Exporter {
	return &Exporter{collector: c}
}

// Export writes metrics to w according to opts.
func (e *Exporter) Export(w io.Writer, opts ExportOptions) error {
	records := e.collector.All()

	if !opts.Since.IsZero() {
		filtered := records[:0]
		for _, r := range records {
			if !r.Timestamp.Before(opts.Since) {
				filtered = append(filtered, r)
			}
		}
		records = filtered
	}

	if opts.EnvFilter != "" {
		filtered := records[:0]
		for _, r := range records {
			if r.EnvPair == opts.EnvFilter {
				filtered = append(filtered, r)
			}
		}
		records = filtered
	}

	switch opts.Format {
	case ExportJSON:
		return exportJSON(w, records)
	case ExportText, "":
		return exportText(w, records)
	default:
		return fmt.Errorf("unknown export format: %q", opts.Format)
	}
}

func exportText(w io.Writer, records []Record) error {
	if len(records) == 0 {
		_, err := fmt.Fprintln(w, "no metrics recorded")
		return err
	}
	for _, r := range records {
		_, err := fmt.Fprintf(w, "[%s] env=%s drifted=%v diffs=%d\n",
			r.Timestamp.Format(time.RFC3339), r.EnvPair, r.Drifted, r.DiffCount)
		if err != nil {
			return err
		}
	}
	return nil
}

func exportJSON(w io.Writer, records []Record) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}
