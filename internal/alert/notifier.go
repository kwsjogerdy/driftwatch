package alert

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Severity represents the urgency level of a drift alert.
type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityCritical Severity = "CRITICAL"
)

// Alert holds a formatted drift notification.
type Alert struct {
	Timestamp   time.Time
	Environment string
	Severity    Severity
	Differences []drift.Difference
}

// Notifier writes drift alerts to a given output.
type Notifier struct {
	Writer io.Writer
}

// NewNotifier creates a Notifier that writes to stdout by default.
func NewNotifier() *Notifier {
	return &Notifier{Writer: os.Stdout}
}

// Notify formats and writes an alert for the provided differences.
// Returns the Alert that was emitted, or nil if there were no differences.
func (n *Notifier) Notify(environment string, diffs []drift.Difference) *Alert {
	if len(diffs) == 0 {
		return nil
	}

	severity := classifySeverity(diffs)
	a := &Alert{
		Timestamp:   time.Now().UTC(),
		Environment: environment,
		Severity:    severity,
		Differences: diffs,
	}

	n.write(a)
	return a
}

func (n *Notifier) write(a *Alert) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] %s — drift detected in environment %q (%d difference(s))\n",
		a.Severity, a.Timestamp.Format(time.RFC3339), a.Environment, len(a.Differences)))
	for _, d := range a.Differences {
		sb.WriteString(fmt.Sprintf("  key=%q  source=%v  target=%v  kind=%s\n",
			d.Key, d.SourceValue, d.TargetValue, d.Kind))
	}
	fmt.Fprint(n.Writer, sb.String())
}

// classifySeverity returns CRITICAL when any keys are missing, WARNING otherwise.
// A missing key indicates that configuration present in the source is absent in
// the target, which is considered more severe than a value mismatch.
func classifySeverity(diffs []drift.Difference) Severity {
	for _, d := range diffs {
		if d.Kind == drift.KindMissing {
			return SeverityCritical
		}
	}
	return SeverityWarning
}

// Summary returns a single-line human-readable description of the alert,
// suitable for use in log lines or notification titles.
func (a *Alert) Summary() string {
	return fmt.Sprintf("[%s] environment=%q differences=%d timestamp=%s",
		a.Severity, a.Environment, len(a.Differences), a.Timestamp.Format(time.RFC3339))
}
