package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Format defines the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds a formatted drift report.
type Report struct {
	GeneratedAt time.Time
	Environment string
	Diffs       []drift.Diff
	Format      Format
}

// NewReport creates a Report for the given environment and diffs.
func NewReport(env string, diffs []drift.Diff, format Format) *Report {
	return &Report{
		GeneratedAt: time.Now().UTC(),
		Environment: env,
		Diffs:       diffs,
		Format:      format,
	}
}

// Render returns the report as a string in the configured format.
func (r *Report) Render() string {
	switch r.Format {
	case FormatJSON:
		return r.renderJSON()
	default:
		return r.renderText()
	}
}

func (r *Report) renderText() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Drift Report — Environment: %s\n", r.Environment))
	sb.WriteString(fmt.Sprintf("Generated: %s\n", r.GeneratedAt.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Total diffs: %d\n", len(r.Diffs)))
	if len(r.Diffs) == 0 {
		sb.WriteString("No drift detected.\n")
		return sb.String()
	}
	sb.WriteString("\nChanges:\n")
	for _, d := range r.Diffs {
		sb.WriteString(fmt.Sprintf("  [%s] key=%q baseline=%q target=%q\n",
			strings.ToUpper(d.Type), d.Key, d.BaselineValue, d.TargetValue))
	}
	return sb.String()
}

func (r *Report) renderJSON() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`{"environment":%q,"generated_at":%q,"total_diffs":%d,"diffs":[`,
		r.Environment, r.GeneratedAt.Format(time.RFC3339), len(r.Diffs)))
	for i, d := range r.Diffs {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf(`{"type":%q,"key":%q,"baseline_value":%q,"target_value":%q}`,
			d.Type, d.Key, d.BaselineValue, d.TargetValue))
	}
	sb.WriteString("]}")  
	return sb.String()
}
