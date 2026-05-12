package metrics

import "fmt"

// Severity represents the breach level of a threshold evaluation.
type Severity string

const (
	SeverityOK   Severity = "ok"
	SeverityWarn Severity = "warn"
	SeverityCrit Severity = "crit"
)

// Threshold defines warn and critical levels for a named summary field.
type Threshold struct {
	Field  string
	WarnAt float64
	CritAt float64
}

// Breach describes a single threshold violation.
type Breach struct {
	Field    string
	Value    float64
	Severity Severity
	Message  string
}

// Evaluator checks a MetricSummary against a set of Threshold rules.
type Evaluator struct {
	thresholds []Threshold
}

// NewEvaluator creates an Evaluator with the provided thresholds.
func NewEvaluator(thresholds []Threshold) *Evaluator {
	return &Evaluator{thresholds: thresholds}
}

// Evaluate returns any breaches found in the given summary.
func (e *Evaluator) Evaluate(s Summary) []Breach {
	fields := map[string]float64{
		"drift_count":    float64(s.TotalRuns),
		"missing_keys":   float64(s.TotalMissing),
		"extra_keys":     float64(s.TotalExtra),
		"changed_values": float64(s.TotalChanged),
	}

	var breaches []Breach
	for _, t := range e.thresholds {
		val, ok := fields[t.Field]
		if !ok {
			continue
		}
		sev := SeverityOK
		if t.CritAt > 0 && val >= t.CritAt {
			sev = SeverityCrit
		} else if t.WarnAt > 0 && val >= t.WarnAt {
			sev = SeverityWarn
		}
		if sev != SeverityOK {
			breaches = append(breaches, Breach{
				Field:    t.Field,
				Value:    val,
				Severity: sev,
				Message:  fmt.Sprintf("%s=%.0f breaches %s threshold (%.0f)", t.Field, val, sev, thresholdLevel(t, sev)),
			})
		}
	}
	return breaches
}

func thresholdLevel(t Threshold, sev Severity) float64 {
	if sev == SeverityCrit {
		return t.CritAt
	}
	return t.WarnAt
}
